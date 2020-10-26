// authenticator provides authentication with Vault
//
// the received Vault token will be stored in VAULT_TOKEN_PATH
//
// authenticator is meant to be used in an init container on Kubernetes
package main

import (
	"fmt"
	"log"
	"os"

	k8s "github.com/postfinance/vaultk8s"
)

func main() {
	c, err := k8s.NewFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	token, err := c.Authenticate()
	if err != nil {
		if c.AllowFail {
			log.Println(fmt.Errorf("authentication failed - ALLOW_FAIL is set therefore pod will continue: %w", err))
			os.Exit(0)
		}

		log.Fatal(fmt.Errorf("authentication failed: %w", err))
	}

	log.Printf("successfully authenticated to vault")

	if err := c.StoreToken(token); err != nil {
		log.Fatal(err)
	}

	log.Printf("successfully stored vault token at %s", c.TokenPath)
}
