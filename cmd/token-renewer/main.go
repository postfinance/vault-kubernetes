// token-renewer renews Vault secrets before they expiration
//
// token-renewer expects a valid Vault token in VAULT_TOKEN_PATH (see Authenticator)
// token-renewer is also able to re-authenticate, if enabled. This can also be used to authenticate initially and to renounce authenticator.
//
// token-renewer is meant to be used as a sidecar container on Kubernetes.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/postfinance/vault/k8s"
)

func main() {
	c, err := k8s.NewFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// get token and re-authenticate if enabled
	token, err := c.GetToken()
	if err != nil {
		log.Fatal(err)
	}
	// store token in case re-authentication was done
	if err := c.StoreToken(token); err != nil {
		log.Fatal(err)
	}

	renewer, err := c.NewRenewer(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("start renewer loop")
	go renewer.Renew()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case err := <-renewer.DoneCh():
			if err != nil {
				log.Fatal(errors.Wrap(err, "failed to renew token"))
			}
			os.Exit(0)
		case <-renewer.RenewCh():
			log.Println("token renewed")
		case <-exit:
			log.Println("signal received - stop execution")
			renewer.Stop()
			goto Exit
		}
	}
Exit:
}
