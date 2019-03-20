// synchronizer synchronizes Vault secrets with Kubernetes secrets
//
// synchronizer expects a valid Vault token in VAULT_TOKEN_PATH (see authenticator)
// all Kubernetes secrets receive an annotation to identify and delete them as synchronized secrets when they are no longer needed
//
// synchronizer is meant to be used in an init container on Kubernetes.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/postfinance/vault/k8s"
	"github.com/postfinance/vault/kv"
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
		log.Println(errors.Wrap(err, "cannot synchronize secrets - all secrets seems to be available therefore pod creation will continue"))
		os.Exit(0)
	}
	c.vault.UseToken(token)

	if err := c.prepare(); err != nil {
		log.Fatal(errors.Wrap(err, "failed to prepare synchronization of secrets"))
	}

	if err := c.synchronize(); err != nil {
		log.Fatal(errors.Wrap(err, "failed to synchronize secrets"))
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
}

func newFromEnvironment() (*syncConfig, error) {
	var err error
	c := &syncConfig{}
	c.vault, err = k8s.NewFromEnvironment()
	if err != nil {
		return nil, err
	}
	c.Secrets = make(map[string]string)
	for _, item := range strings.Split(os.Getenv("VAULT_SECRETS"), ",") {
		if len(item) == 0 {
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
	content, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return nil, errors.Wrap(err, "could not get namespace")
	}
	c.Namespace = strings.TrimSpace(string(content))
	// connect to kubernetes
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get k8s config")
	}
	c.k8sClientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get k8s k8sClientset")
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
		_, err := sc.k8sClientset.CoreV1().Secrets(sc.Namespace).Get(k, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("secret %s does not exist", k)
		}
	}
	return nil
}

// synchronize secret from vault to the current kubernetes namespace
func (sc *syncConfig) synchronize() error {
	// create/update the secrets
	annotations := make(map[string]string)
	for k, v := range sc.Secrets {
		// get secret from vault
		log.Println("read", v, "from vault")
		s, err := sc.secretClients[strings.SplitN(v, "/", 2)[0]].Read(v)
		if err != nil {
			return err
		}
		if s == nil {
			continue
		}
		// convert data
		data := make(map[string][]byte)
		for k, v := range s {
			data[k] = []byte(v.(string))
		}
		// create/update k8s secret
		annotations[vaultAnnotation] = v
		secret := &corev1.Secret{}
		secret.Name = fmt.Sprintf("%s%s", sc.SecretPrefix, k)
		secret.Data = data
		secret.Annotations = annotations
		// create (insert) or update the secret
		_, err = sc.k8sClientset.CoreV1().Secrets(sc.Namespace).Get(secret.Name, metav1.GetOptions{})
		if apierr.IsNotFound(err) {
			log.Println("create secret", secret.Name, "from vault secret", v)
			if _, err := sc.k8sClientset.CoreV1().Secrets(sc.Namespace).Create(secret); err != nil {
				return err
			}
			continue
		}
		log.Println("update secret", secret.Name, "from vault secret", v)
		if _, err = sc.k8sClientset.CoreV1().Secrets(sc.Namespace).Update(secret); err != nil {
			return err
		}
	}
	// delete obsolete secrets
	secretList, err := sc.k8sClientset.CoreV1().Secrets(sc.Namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Println(errors.Wrap(err, "cleanup of unused vault secrets failed"))
		os.Exit(0)
	}
	for _, s := range secretList.Items {
		// only secrets from vault
		if _, ok := s.Annotations[vaultAnnotation]; !ok {
			continue
		}
		// only if vault secret is not in secrets
		if _, ok := sc.Secrets[strings.TrimPrefix(s.Name, sc.SecretPrefix)]; ok {
			continue
		}
		log.Println("delete secret", s.Name)
		if err := sc.k8sClientset.CoreV1().Secrets(sc.Namespace).Delete(s.Name, &metav1.DeleteOptions{}); err != nil {
			log.Println(errors.Wrapf(err, "delete obsolete vault secret %s failed", s.Name))
		}
	}
	return nil
}

// prepare
func (sc *syncConfig) prepare() error {
	sc.secretClients = make(map[string]*kv.Client)
	secrets := make(map[string]string)
	for k, v := range sc.Secrets {
		mount := strings.SplitN(v, "/", 2)[0]
		// ensure kv.Client for mount
		if _, ok := sc.secretClients[mount]; !ok {
			secretClient, err := kv.New(sc.vault.Client(), mount+"/")
			if err != nil {
				return err
			}
			sc.secretClients[mount] = secretClient
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
		// TODO: check for secret == nil
		for _, k := range keys {
			secrets[k] = path.Join(v, k)
		}
	}
	sc.Secrets = secrets
	return nil
}
