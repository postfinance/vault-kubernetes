// Package auth provides authentication with Vault on Kubernetes
//
// Authentication is done with the Kubernetes Auth Method by Vault.
//
// See also ``Kubernetes Auth Method`` from the Vault documentation
// https://www.vaultproject.io/docs/auth/kubernetes.html
package auth

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

// Config represents the configuration to get a valid Vault token
type Config struct {
	VaultRole          string
	VaultTokenPath     string
	VaultReAuth        bool
	VaultTTL           int
	VaultK8SMountPath  string
	ServiceAccountPath string
	AllowFail          bool
	vault              *api.Client
}

// NewConfigFromEnvironment returns a initialized Config for authentication
func NewConfigFromEnvironment() (*Config, error) {
	c := &Config{}
	c.VaultRole = os.Getenv("VAULT_ROLE")
	c.VaultTokenPath = os.Getenv("VAULT_TOKEN_PATH")
	if c.VaultTokenPath == "" {
		return nil, fmt.Errorf("missing VAULT_TOKEN_PATH")
	}
	if s := os.Getenv("VAULT_REAUTH"); s != "" {
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, errors.Wrap(err, "1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False are valid values for ALLOW_FAIL")
		}
		c.VaultReAuth = b
	}
	if s := os.Getenv("VAULT_TTL"); s != "" {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.Wrap(err, "%s is not a valid integer for VAULT_TTL")
		}
		c.VaultTTL = i
	}
	c.VaultK8SMountPath = os.Getenv("VAULT_K8S_MOUNT_PATH")
	if c.VaultK8SMountPath == "" {
		c.VaultK8SMountPath = "auth/kubernetes/login"
	}
	c.ServiceAccountPath = os.Getenv("SERVICE_ACCOUNT_PATH")
	if c.ServiceAccountPath == "" {
		c.ServiceAccountPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	}
	if s := os.Getenv("ALLOW_FAIL"); s != "" {
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, errors.Wrap(err, "1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False are valid values for ALLOW_FAIL")
		}
		c.AllowFail = b
	}
	// create vault client
	vaultConfig := api.DefaultConfig()
	if err := vaultConfig.ReadEnvironment(); err != nil {
		return nil, errors.Wrap(err, "failed to read environment for vault")
	}
	var err error
	c.vault, err = api.NewClient(vaultConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create vault client")
	}
	return c, nil
}

// Authenticate to vault
func (c *Config) Authenticate() (string, error) {
	var empty string
	// read jwt of serviceaccount
	content, err := ioutil.ReadFile(c.ServiceAccountPath)
	if err != nil {
		return empty, errors.Wrap(err, "failed to read jwt token")
	}
	jwt := string(bytes.TrimSpace(content))

	// authenticate
	data := make(map[string]interface{})
	data["role"] = c.VaultRole
	data["jwt"] = jwt
	s, err := c.vault.Logical().Write(c.VaultK8SMountPath, data)
	if err != nil {
		return empty, errors.Wrapf(err, "login failed with role from environment variable VAULT_ROLE: %q", c.VaultRole)
	}
	if len(s.Warnings) > 0 {
		return empty, fmt.Errorf("login failed with: %s", strings.Join(s.Warnings, " - "))
	}
	return s.Auth.ClientToken, nil
}

// LoadToken from VaultTokenPath
func (c *Config) LoadToken() (string, error) {
	content, err := ioutil.ReadFile(c.VaultTokenPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to load token")
	}
	return string(content), nil
}

// StoreToken in VaultTokenPath
func (c *Config) StoreToken(token string) error {
	if err := ioutil.WriteFile(c.VaultTokenPath, []byte(token), 0644); err != nil {
		return errors.Wrap(err, "failed to save token")
	}
	return nil
}

// GetToken tries to load the vault token from VaultTokenPath
// if token is not available, invalid or not renewable
// and VaultReAuth is true, try to re-authenticate
func (c *Config) GetToken() (string, error) {
	var empty string
	token, err := c.LoadToken()
	if err != nil {
		if c.VaultReAuth {
			return c.Authenticate()
		}
		return empty, errors.Wrapf(err, "failed to load token form: %s", c.VaultTokenPath)
	}
	c.vault.SetToken(token)
	if _, err = c.vault.Auth().Token().RenewSelf(c.VaultTTL); err != nil {
		if c.VaultReAuth {
			return c.Authenticate()
		}
		return empty, errors.Wrap(err, "failed to renew token")
	}
	return token, nil
}

// NewRenewer returns a *api.Renewer to renew the vault token regularly
func (c *Config) NewRenewer(token string) (*api.Renewer, error) {
	c.vault.SetToken(token)
	// renew the token to get a secret usable for renewer
	secret, err := c.vault.Auth().Token().RenewSelf(c.VaultTTL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to renew-self token")
	}
	renewer, err := c.vault.NewRenewer(&api.RenewerInput{Secret: secret})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token renewer")
	}
	return renewer, nil
}
