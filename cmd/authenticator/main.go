// authenticator provides authentication with Vault
//
// the received Vault token will be stored in VAULT_TOKEN_PATH
//
// authenticator is meant to be used in an init container on Kubernetes.
package main

import (
	"log"
	"os"

	"github.com/postfinance/vault-kubernetes/pkg/auth"
	"github.com/pkg/errors"
)

func main() {
	c, err := auth.NewConfigFromEnvironment()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to get config"))
	}

	token, err := c.Authenticate()
	if err != nil {
		if c.AllowFail {
			log.Println(errors.Wrap(err, "authentication failed - ALLOW_FAIL is set therefore pod will continue"))
			os.Exit(0)
		} else {
			log.Fatal(errors.Wrap(err, "authentication failed"))
		}
	}
	log.Printf("successfully authenticated to vault")

	if err := c.StoreToken(token); err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully stored vault token at %s", c.VaultTokenPath)

	os.Exit(0)
}
