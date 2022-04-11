package main

import (
	"context"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	vault "github.com/hashicorp/vault/api"
	kv "github.com/postfinance/vaultkv"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	quote         = "The fool don‘t think he is wise, but the wise man knows himself to be a fool."
	trivialString = "h"
)

func TestDecode(t *testing.T) {
	t.Run("not encoded", func(t *testing.T) {
		res, err := decode(quote)
		assert.NoError(t, err)
		assert.Equal(t, quote, string(res))
	})

	t.Run("base64 encoded", func(t *testing.T) {
		str := "base64:" + base64.StdEncoding.EncodeToString([]byte(quote))
		res, err := decode(str)
		assert.NoError(t, err)
		assert.Equal(t, quote, string(res))
	})

	t.Run("base64 encoded", func(t *testing.T) {
		str := "base64:" + base64.StdEncoding.EncodeToString([]byte(trivialString))
		res, err := decode(str)
		assert.NoError(t, err)
		assert.Equal(t, trivialString, string(res))
	})

	t.Run("base64 decode fails", func(t *testing.T) {
		str := "base64:" + quote
		_, err := decode(str)
		assert.Error(t, err)
	})

	t.Run("unknown encoding", func(t *testing.T) {
		str := "base32:" + base32.StdEncoding.EncodeToString([]byte(quote))
		res, err := decode(str)
		assert.NoError(t, err)
		assert.Equal(t, str, string(res))
	})
}

func TestDecodeAcceptance(t *testing.T) {

}

func TestSplitLabels(t *testing.T) {
	labels := "s1=batman,s2,s3=superman,s4=,s5,"

	exp := map[string]string{
		"s1": "batman",
		"s3": "superman",
		"s4": "", // guess who? the invisible man.
	}

	res := splitLabels(labels)

	require.True(t, len(exp) == len(res))

	for k, v := range res {
		assert.Equal(t, v, exp[k])
	}
}

func TestMergeLabels(t *testing.T) {
	existing := map[string]string{
		"e1": "batman",
		"e2": "superman",
	}

	configured := map[string]string{
		"c1": "wonder woman",
		"e2": "supergirl",
	}

	exp := map[string]string{
		"e1": "batman",
		"c1": "wonder woman",
		"e2": "supergirl",
	}

	res := mergeLabels(existing, configured)

	require.True(t, len(exp) == len(res))

	for k, v := range res {
		assert.Equal(t, v, exp[k])
	}
}

//nolint:funlen // complex integration testing
func TestIntegration(t *testing.T) {
	const (
		rootToken = "90b03685-e17b-7e5e-13a0-e14e45baeb2f" // nolint: gosec
	)

	// Setup Vault
	pool, err := dockertest.NewPool("unix:///var/run/docker.sock")
	require.NoError(t, err, "could not connect to Docker")

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("vault", "latest", []string{
		"VAULT_DEV_ROOT_TOKEN_ID=" + rootToken,
		"VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200",
	})
	require.NoError(t, err, "could not start container")

	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	})

	host := os.Getenv("DOCKER_HOST")
	if host == "" {
		host = "localhost"
	}

	if host != "localhost" && !strings.Contains(host, ".") {
		host += ".pnet.ch"
	}

	vaultAddr := fmt.Sprintf("http://%s:%s", host, resource.GetPort("8200/tcp"))

	_ = os.Setenv("VAULT_ADDR", vaultAddr)
	_ = os.Setenv("VAULT_TOKEN", rootToken)

	fmt.Println("VAULT_ADDR:", vaultAddr)

	vaultConfig := vault.DefaultConfig()
	require.NoError(t, vaultConfig.ReadEnvironment(), "failed to read config from environment")

	vaultClient, err := vault.NewClient(vaultConfig)
	require.NoError(t, err, "failed to create Vault client")

	err = pool.Retry(func() error {
		_, err = vaultClient.Sys().ListMounts()
		return err
	})
	require.NoError(t, err, "failed to connect to Vault")

	// Setup Kubernetes (kind)
	u, err := user.Current()
	require.NoError(t, err, "failed to get current user")

	configPath := filepath.Join(u.HomeDir, ".kube", "config")

	c, err := clientcmd.BuildConfigFromFlags("", configPath)
	require.NoError(t, err, "failed to get Kubernetes config from: %s", configPath)

	clientset, err := kubernetes.NewForConfig(c)
	require.NoError(t, err, "failed to create client set")

	_, err = clientset.DiscoveryClient.ServerVersion()
	require.NoError(t, err, "failed to connect cluster")

	type testStruct struct {
		Bool            bool
		Int             int
		Float           float32
		String          string
		StringDecode1st string
		Byte            []byte
		SliceOfInt      []int
		SliceOfFloat    []float64
		SliceOfString   []string
	}

	testData := testStruct{
		Bool:          true,
		Int:           42,
		Float:         42.42,
		String:        "an ordinary a string",
		Byte:          []byte("a string as byte slice"),
		SliceOfInt:    []int{1, 2, 3},
		SliceOfFloat:  []float64{1.1, 2.22, 3.333},
		SliceOfString: []string{"A", "B", "C"},
	}

	testData.StringDecode1st = fmt.Sprintf("base64:%s", base64.StdEncoding.EncodeToString([]byte(testData.String)))

	secretPath := "secret/data/test"

	t.Run("decode data from vault", func(t *testing.T) {
		jsonEncoded, err := json.Marshal(testData)
		require.NoError(t, err)
		require.NotEmpty(t, jsonEncoded)

		inputData := map[string]interface{}{
			"data": map[string]interface{}{
				"bool":              testData.Bool,
				"int":               testData.Int,
				"float":             testData.Float,
				"string":            testData.String,
				"stringDecode1st":   testData.StringDecode1st,
				"byte":              testData.Byte,
				"sliceOfInt":        testData.SliceOfInt,
				"sliceOfFloat":      testData.SliceOfFloat,
				"sliceOfString":     testData.SliceOfString,
				"structJSONEncoded": jsonEncoded,
				"struct":            testData,
			},
		}

		_, err = vaultClient.Logical().Write(secretPath, inputData)
		require.NoError(t, err, "failed to write secret %s", secretPath)

		s, err := vaultClient.Logical().Read(secretPath)
		require.NoError(t, err, "failed to read secret %s", secretPath)

		secrets := s.Data["data"].(map[string]interface{})
		v, err := decode(secrets["bool"])
		require.NoError(t, err, "failed to decode bool")
		assert.Equal(t, []byte(fmt.Sprintf("%v", testData.Bool)), v)

		v, err = decode(secrets["int"])
		require.NoError(t, err, "failed to decode int: %v", v)
		assert.Equal(t, []byte(fmt.Sprintf("%v", testData.Int)), v)

		v, err = decode(secrets["float"])
		require.NoError(t, err, "failed to decode float: %v", v)
		assert.Equal(t, []byte(fmt.Sprintf("%.2f", testData.Float)), v)

		v, err = decode(secrets["string"])
		require.NoError(t, err, "failed to decode string: %v", v)
		assert.Equal(t, []byte(testData.String), v)

		v, err = decode(secrets["stringDecode1st"])
		require.NoError(t, err, "failed to decode stringDecode1st: %v", v)
		assert.Equal(t, []byte(testData.String), v, string(v))

		v, err = decode(secrets["byte"])
		require.NoError(t, err, "failed to decode []byte: %v", v)
		bytes, err := base64.StdEncoding.DecodeString(string(v))
		require.NoError(t, err, "failed to decode base64 encoded []byte")
		assert.Equal(t, testData.Byte, bytes)

		v, err = decode(secrets["sliceOfInt"])
		require.NoError(t, err, "failed to decode sliceOfInt: %v", v)
		sliceOfInt, err := json.Marshal(testData.SliceOfInt)
		require.NoError(t, err, "failed to json.Marshal sliceOfInt")
		assert.Equal(t, sliceOfInt, v)

		v, err = decode(secrets["sliceOfFloat"])
		require.NoError(t, err, "failed to decode sliceOfFloat: %v", v)
		sliceOfFloat, err := json.Marshal(testData.SliceOfFloat)
		require.NoError(t, err, "failed to json.Marshal sliceOfFloat")
		assert.Equal(t, sliceOfFloat, v)

		v, err = decode(secrets["sliceOfString"])
		require.NoError(t, err, "failed to decode sliceOfString: %v", v)
		sliceOfString, err := json.Marshal(testData.SliceOfString)
		require.NoError(t, err, "failed to json.Marshal sliceOfString")
		assert.Equal(t, sliceOfString, v)

		v, err = decode(secrets["structJSONEncoded"])
		require.NoError(t, err, "failed to decode structJSONEncoded: %v", v)
		bytes, err = base64.StdEncoding.DecodeString(string(v))
		require.NoError(t, err, "failed to decode base64 encoded []byte")
		assert.Equal(t, jsonEncoded, bytes)

		v, err = decode(secrets["struct"])
		require.NoError(t, err, "failed to decode struct: %v", v)
		act := testStruct{}
		require.NoError(t, json.Unmarshal(v, &act), "failed to json.Marshal struct")
		assert.Equal(t, testData, act)
	})

	t.Run("synchronize with vault", func(t *testing.T) {
		secretClient, err := kv.New(vaultClient, secretPath)
		require.NoError(t, err, "failed to create kv.Client for %s", secretPath)

		c := &syncConfig{
			Secrets: map[string]string{
				"test": "secret/data/test",
			},
			SecretPrefix: "v3t-",
			Namespace:    "default",
			k8sClientset: clientset,
			secretClients: map[string]*kv.Client{
				secretClient.Mount: secretClient,
			},
			annotation: vaultAnnotation,
		}

		require.NoError(t, c.synchronize(), "failed to synchronize secrets")

		_, err = clientset.CoreV1().Secrets("default").Get(context.TODO(), "v3t-test", v1.GetOptions{})
		require.NoError(t, err, "failed to get k8s secret v3t-test")

		//nolint:godox // to be solved
		/*
				TODO: how to deal with []byte and other already base64 encoded strings?
			for k, v := range s.Data {
				t.Log(k, string(v))
			}
		*/
	})
}
