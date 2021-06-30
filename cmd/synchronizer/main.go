// synchronizer synchronizes Vault secrets with Kubernetes secrets
//
// synchronizer expects a valid Vault token in VAULT_TOKEN_PATH (see authenticator)
// all Kubernetes secrets receive an annotation to identify and delete them as synchronized secrets when they are no longer needed
//
// synchronizer is meant to be used in an init container on Kubernetes.
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	k8s "github.com/postfinance/vaultk8s"
	kv "github.com/postfinance/vaultkv"
	corev1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	vaultAnnotation = "vault-secret"
)

func main() {
	c, err := newFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	token, err := c.vault.LoadToken()
	if err != nil {
		if err := c.checkSecrets(); err != nil {
			log.Fatal(err)
		}
		// you get only here ...
		// IF  ALLOW_FAIL=true was set for vault-kubernetes-authenticator
		// AND vault-kubernetes-authenticator failed to authenticate
		// AND all testable secrets are present
		log.Println(fmt.Errorf("cannot synchronize secrets - all secrets seems to be available therefore pod creation will continue: %w", err))
		os.Exit(0)
	}

	c.vault.UseToken(token)

	if err := c.prepare(); err != nil {
		log.Fatal(fmt.Errorf("failed to prepare synchronization of secrets: %w", err))
	}

	if err := c.synchronize(); err != nil {
		log.Fatal(fmt.Errorf("failed to synchronize secrets: %w", err))
	}

	log.Printf("secrets successfully synchronized")
}

type syncConfig struct {
	Secrets       map[string]string // key = kubernetes secret name, value = vault secret name
	SecretPrefix  string            // prefix for kubernetes secret name
	Namespace     string
	k8sClientset  *kubernetes.Clientset
	secretClients map[string]*kv.Client
	vault         *k8s.Vault
	annotation    string
	labels        map[string]string
}

func newFromEnvironment() (*syncConfig, error) {
	var err error

	c := &syncConfig{}

	c.vault, err = k8s.NewFromEnvironment()
	if err != nil {
		return nil, err
	}

	c.annotation = getEnv("SYNCHRONIZER_ANNOTATION", vaultAnnotation)

	log.Println("Using annotation [", c.annotation, "] to detect managed secrets")

	c.labels = splitLabels(getEnv("SYNCHRONIZER_LABELS", ""))

	c.Secrets = make(map[string]string)

	for _, item := range strings.Split(os.Getenv("VAULT_SECRETS"), ",") {
		if item == "" {
			continue
		}

		s := strings.SplitN(item, ":", 2)

		switch {
		case strings.HasSuffix(s[0], "/"):
			c.Secrets[s[0]] = s[0]
		case len(s) > 1:
			c.Secrets[s[1]] = s[0]
		default:
			c.Secrets[path.Base(s[0])] = s[0]
		}
	}

	if len(c.Secrets) == 0 {
		return nil, fmt.Errorf("no secrets to synchronize - check VAULT_SECRETS")
	}

	c.SecretPrefix = os.Getenv("SECRET_PREFIX")

	// current kubernetes namespace
	content, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return nil, fmt.Errorf("could not get namespace: %w", err)
	}

	c.Namespace = strings.TrimSpace(string(content))

	// connect to kubernetes
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s config: %w", err)
	}

	c.k8sClientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s k8sClientset: %w", err)
	}

	return c, nil
}

// checkSecrets check the existence of a secret and not the content
func (sc *syncConfig) checkSecrets() error {
	// check secrets
	for k, v := range sc.Secrets {
		if strings.HasSuffix(v, "/") {
			log.Printf("WARNING: cannot check existence of secrets from vault path %s without connection to vault\n", v)
			continue
		}

		log.Println("check k8s secret", k, "from vault secret", v)

		_, err := sc.k8sClientset.CoreV1().Secrets(sc.Namespace).Get(context.Background(), k, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("could not fetch secret %s from namespace %s: %s", k, sc.Namespace, err.Error())
		}
	}

	return nil
}

// synchronize secret from vault to the current kubernetes namespace
// nolint: gocognit, gocyclo, funlen
func (sc *syncConfig) synchronize() error {
	// create/update the secrets
	annotations := make(map[string]string)

	for k, v := range sc.Secrets {
		var clientExists bool

		var mount string

		for m := range sc.secretClients {
			if strings.HasPrefix(v, m) {
				clientExists = true
				mount = m

				break
			}
		}

		if !clientExists {
			return fmt.Errorf("no client exists for %s", v)
		}
		// get secret from vault

		log.Println("read", v, "from vault")

		s, err := sc.secretClients[mount].Read(v)
		if err != nil {
			return err
		}

		if s == nil {
			log.Println("secret", v, "not found")
			continue
		}
		// convert data
		data := make(map[string][]byte)

		for k, v := range s {
			w, err := decode(v.(string))
			if err != nil {
				return err
			}

			data[k] = w
		}
		// create/update k8s secret
		annotations[sc.annotation] = v
		secret := &corev1.Secret{}
		secret.Name = fmt.Sprintf("%s%s", sc.SecretPrefix, k)
		secret.Data = data
		secret.Annotations = annotations
		// create (insert) or update the secret
		existing, err := sc.k8sClientset.CoreV1().Secrets(sc.Namespace).Get(context.Background(), secret.Name, metav1.GetOptions{})
		if err != nil {
			if apierr.IsNotFound(err) {
				log.Println("create secret", secret.Name, "from vault secret", v)

				if _, err := sc.k8sClientset.CoreV1().Secrets(sc.Namespace).Create(context.Background(), secret, metav1.CreateOptions{}); err != nil {
					return err
				}

				continue
			}

			return err
		}

		if _, ok := existing.Annotations[sc.annotation]; !ok {
			log.Println("WARNING: ignoring secret", secret.Name, "- not managed by synchronizer")
			continue
		}

		secret.Labels = mergeLabels(existing.Labels, sc.labels)

		log.Println("update secret", secret.Name, "from vault secret", v)

		if _, err = sc.k8sClientset.CoreV1().Secrets(sc.Namespace).Update(context.Background(), secret, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}
	// delete obsolete secrets
	secretList, err := sc.k8sClientset.CoreV1().Secrets(sc.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Println(fmt.Errorf("cleanup of unused vault secrets failed: %w", err))
		os.Exit(0)
	}

	for i := range secretList.Items {
		s := secretList.Items[i]
		// only secrets from vault
		if _, ok := s.Annotations[sc.annotation]; !ok {
			continue
		}
		// only if vault secret is not in secrets
		if _, ok := sc.Secrets[strings.TrimPrefix(s.Name, sc.SecretPrefix)]; ok {
			continue
		}

		log.Println("delete secret", s.Name)

		if err := sc.k8sClientset.CoreV1().Secrets(sc.Namespace).Delete(context.Background(), s.Name, metav1.DeleteOptions{}); err != nil {
			log.Println(fmt.Errorf("delete obsolete vault secret %s failed: %w", s.Name, err))
		}
	}

	return nil
}

// prepare
func (sc *syncConfig) prepare() error {
	sc.secretClients = make(map[string]*kv.Client)
	secrets := make(map[string]string)

	for k, v := range sc.Secrets {
		var clientExists bool

		var mount string

		// check if we already have a client for this prefix
		for m := range sc.secretClients {
			if strings.HasPrefix(v, m) {
				clientExists = true
				mount = m
				log.Println("found client for mount: ", mount)

				break
			}
		}

		if !clientExists {
			secretClient, err := kv.New(sc.vault.Client(), v)
			if err != nil {
				return err
			}

			mount = secretClient.Mount
			sc.secretClients[mount] = secretClient

			log.Printf("created v%d client for mount: %s", secretClient.Version, mount)
		}
		// v is a secret
		if !strings.HasSuffix(v, "/") {
			secrets[k] = v
			continue
		}
		// v is a path -> get all secrets from v
		keys, err := sc.secretClients[mount].List(v)
		if err != nil {
			return err
		}

		if keys == nil {
			continue
		}

		// nolint: godox // TODO: check for secret == nil
		for _, k := range keys {
			secrets[k] = path.Join(v, k)
		}
	}

	sc.Secrets = secrets

	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func decode(s string) ([]byte, error) {
	switch {
	case strings.HasPrefix(s, "base64:"):
		return base64.StdEncoding.DecodeString(strings.TrimPrefix(s, "base64:"))
	default:
		return []byte(s), nil
	}
}

// splitLabels splits a string like "a=b,c=d,e=f" and returns a map
func splitLabels(s string) map[string]string {
	labels := make(map[string]string)

	for _, l := range strings.Split(s, ",") {
		p := strings.SplitN(l, "=", 2)

		if len(p) != 2 {
			continue
		}

		labels[p[0]] = p[1]
	}

	return labels
}

// mergeLabels merges existing labels with configured labels
// existing labels will be overwritten if a configured label with the same key exists
func mergeLabels(existing, configured map[string]string) map[string]string {
	labels := make(map[string]string)

	for k, v := range existing {
		labels[k] = v
	}

	for k, v := range configured {
		labels[k] = v
	}

	return labels
}
