[![Go Report Card](https://goreportcard.com/badge/github.com/postfinance/vault-kubernetes)](https://goreportcard.com/report/github.com/postfinance/vault-kubernetes)
[![Build Status](https://travis-ci.org/postfinance/vault-kubernetes.svg?branch=master)](https://travis-ci.org/postfinance/vault-kubernetes)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Credits](#credits)
- [Scenarios](#scenarios)
    - [Scenario 1 - Get a Vault token for one time use](#scenario-1-get-a-vault-token-for-one-time-use)
    - [Scenario 2 - Sync Vault secrets to Kubernetes secrets](#scenario-2-sync-vault-secrets-to-kubernetes-secrets)
    - [Scenario 3 - Get a Vault token for use during the lifetime of a pod](#scenario-3-get-a-vault-token-for-use-during-the-lifetime-of-a-pod)
- [Issues](#issues)
- [Vault client configuration](#vault-client-configuration)
- [Init Container _vault-kubernetes-authenticator_](#init-container-_vault-kubernetes-authenticator_)
    - [Configuration](#configuration)
    - [Example](#example)
- [Init Container _vault-kubernetes-synchronizer_](#init-container-_vault-kubernetes-synchronizer_)
    - [Secret Mapping](#secret-mapping)
    - [Configuration](#configuration-1)
    - [Error handling](#error-handling)
    - [Example](#example-1)
    - [Example - with failed authentication](#example-with-failed-authentication)
- [Sidecar _vault-kubernetes-token-renewer_](#sidecar-_vault-kubernetes-token-renewer_)
    - [Configuration](#configuration-2)
    - [Example](#example-2)
- [Build](#build)
- [Demo](#demo)
- [Links](#links)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


# Credits
[Based on the work of Seth Vargo](https://github.com/sethvargo/vault-kubernetes-authenticator)


# Scenarios

## Scenario 1 - Get a Vault token for one time use

Start the Init Container _vault-kubernetes-authenticator_ to authenticate to Vault and get a Vault token.

The Vault token will expire after the given TTL.

## Scenario 2 - Sync Vault secrets to Kubernetes secrets

Start the Init Container _vault-kubernetes-authenticator_ to authenticate to Vault and get a Vault token.

After successful completion, start the Init Container _vault-kubernetes-synchronizer_ to synchronize secrets to Kubernetes.

The Vault token will expire after the given TTL.

## Scenario 3 - Get a Vault token for use during the lifetime of a pod

Start the Init Container _vault-kubernetes-authenticator_ to authenticate to Vault and get a Vault token.

After successful completion start the Sidecar Container _vault-kubernetes-token-renewer_ to regularly renew your Vault token.


# Issues

_vault-kubernetes-token-renewer_ container will be restarted if the token renewal fails (for restartPolicy=always). When the token cannot be renewed (e.g. the token is in the meantime expired):
- let the pod terminate and restart. On restart _vault-kubernetes-authenticator_ will issue a new token. A possible solution could be to use [Share Process Namespace between Containers in a Pod](https://kubernetes.io/docs/tasks/configure-pod-container/share-process-namespace) (Kubernetes 1.12 beta) and [Container Lifecycle Hooks](https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks)
- let _vault-kubernetes-token-renewer_ re-authenticate and update VAULT_TOKEN_PATH if the token in VAULT_TOKEN_PATH is invalid. The token consumer needs to observe VAULT_TOKEN_PATH for changes (inotify) or read VAULT_TOKEN_PATH on every connect to Vault (isn't a thing because VAULT_TOKEN_PATH is usually in-memory). This can be done independent from the previous case because the token will be valid on after pod creation

removed `go.sum` from repo due to issue with go version and k8s.io/client-go:
```
go: verifying k8s.io/client-go@v9.0.0+incompatible: checksum mismatch
```


# Vault client configuration

The usual environment variables for Vault will be used:

- VAULT_ADDR
- VAULT_CACERT
- VAULT_CAPATH
- VAULT_CLIENT_CERT
- VAULT_CLIENT_KEY
- VAULT_CLIENT_TIMEOUT
- VAULT_SKIP_VERIFY
- VAULT_TLS_SERVER_NAME
- VAULT_WRAP_TTL
- VAULT_MAX_RETRIES
- VAULT_TOKEN
- VAULT_MFA
- VAULT_RATE_LIMIT

> see https://godoc.org/github.com/hashicorp/vault/api#Config.ReadEnvironment

> the minimal configuration is VAULT_ADDR with VAULT_SKIP_VERIFY=true


# Init Container _vault-kubernetes-authenticator_

## Configuration

- VAULT_ROLE - Required the name of the Vault role to use for authentication.

- VAULT_TOKEN_PATH - the destination path on disk to store the token. Usually this is a shared volume.

- VAULT_K8S_MOUNT_PATH - the name of the mount where the Kubernetes auth method is enabled. This defaults to auth/kubernetes, but if you changed the mount path you will need to set this value to that path (vault auth enable -path=k8s kubernetes -> VAULT_K8S_MOUNT_PATH=auth/k8s)

- SERVICE_ACCOUNT_PATH - the path on disk where the Kubernetes service account jtw token lives. This defaults to /var/run/secrets/kubernetes.io/serviceaccount/token.

- ALLOW_FAIL - the container will successfully terminate even if the authentication to Vault failed, no token will be written to VAULT_TOKEN_PATH. **This condition needs to be handeled in the succeeding container.** (default: "false")

## Example

```
$ k logs vault-kubernetes-authenticator-5675d58d95-4wd8v -c vault-kubernetes-authenticator
2018/11/26 14:56:29 successfully authenticated to vault
2018/11/26 14:56:29 successfully stored vault token at /home/vault/.vault-token

$ k exec -ti vault-kubernetes-authenticator-5675d58d95-4wd8v sh
~ $ VAULT_TOKEN=$(cat /home/vault/.vault-token)
~ $ echo $VAULT_TOKEN
8Pj0EzFLWQv8uWcjbP9hF1MB
~ $
```


# Init Container _vault-kubernetes-synchronizer_

Depends on Init Container _vault-kubernetes-authenticator_

- each Kubernetes secrets created by _vault-kubernetes-synchronizer_ get the annotation `vault-secret: <vault secret path>`

- obsolete secrets created by _vault-kubernetes-synchronizer_ will be deleted

## Secret Mapping

| Mapping                     | Vault                 | Kubernetes  | Remark             |
|-----------------------------|:----------------------|:------------|:-------------------|
| secret/k8s/first            | secret/k8s/first      | first       | Vault KV version 1 |
| secret/k8s/first:third      | secret/k8s/first      | third       | Vault KV version 1 |
| secret/data/k8s/first       | secret/data/k8s/first | first       | Vault KV version 2 |
| secret/data/k8s/first:third | secret/data/k8s/first | third       | Vault KV version 2 |

> you have to provide the correct secret path in Vault
> for KV version 1 the path starts with secret/
> for KV version 2 the path starts with secret/data

> labels/names in Kubernetes will be validated according to [RFC-1123](https://tools.ietf.org/html/rfc1123)

## Configuration

- VAULT_TOKEN_PATH - the destination path on disk to store the token. Usually this is a shared volume.

- VAULT_SECRETS - comma separated list of secrets (see Secret Mapping)

- SECRET_PREFIX - prefix for synchronized secrets (e.g. for SECRET_PREFIX="v3t_" Vault secret "first" will get secret "v3t_first" in k8s)

> set ALLOW_FAIL="true" for _vault-kubernetes-authenticator_

## Error handling

If Vault authentication fails in _vault-kubernetes-authenticator_ and ALLOW_FAIL="true" has been set for _vault-kubernetes-authenticator_ the failed authentication will be handeled as follows:
- all secrets in VAULT_SECRETS are available in the namespace (the content of the secrets will not be considered)- _vault-kubernetes-synchronizer_ issues a warning and terminates successfullly.
- any secret from VAULT_SECRETS is missing in the namespace _vault-secret-synchronizer_ fails.

## Example

Two secrets in Vault:
```
$ vault kv get secret/k8s/first
====== Metadata ======
...
=== Data ===
Key    Value
---    -----
one    12345678
two    23456781
$ vault kv get secret/k8s/second
====== Metadata ======
...
===== Data =====
Key       Value
---       -----
green     lantern
poison    ivy
```

Configure the two secrets for synchronisation with the environment variable VAULT_SECRETS:
```
$ vi deployment.yaml
...
    - name: VAULT_SECRETS
      value: secret/data/k8s/first,secret/data/k8s/second
...
```

```
$ k logs vault-kubernetes-synchronizer-6875c88858-t6hdw -c vault-kubernetes-authenticator
2018/11/26 14:56:30 successfully authenticated to vault
2018/11/26 14:56:30 successfully stored vault token at /home/vault/.vault-token

$ k logs vault-kubernetes-synchronizer-6875c88858-t6hdw -c vault-kubernetes-synchronizer
2018/11/26 14:56:31 read secret/data/k8s-np/appl-vault-dev-e1/first from vault
2018/11/26 14:56:31 create secret third from vault secret secret/data/k8s-np/appl-vault-dev-e1/first
2018/11/26 14:56:31 read secret/data/k8s-np/appl-vault-dev-e1/first from vault
2018/11/26 14:56:31 create secret first from vault secret secret/data/k8s-np/appl-vault-dev-e1/first
2018/11/26 14:56:31 read secret/data/k8s-np/appl-vault-dev-e1/second from vault
2018/11/26 14:56:31 create secret second from vault secret secret/data/k8s-np/appl-vault-dev-e1/second
2018/11/26 14:56:31 secrets successfully synchronized

$ k get secrets | grep -e first -e second -e third
first                                Opaque                                2      16m
second                               Opaque                                2      16m
third                                Opaque                                2      16m

$ k describe secrets first second third
Name:         first
Namespace:    vault-test
Labels:       <none>
Annotations:  vault-secret=secret/data/k8s/first

Type:  Opaque

Data
====
one:  8 bytes
two:  8 bytes


Name:         second
Namespace:    vault-test
Labels:       <none>
Annotations:  vault-secret=secret/data/k8s/second

Type:  Opaque

Data
====
poison:  3 bytes
green:   7 bytes


Name:         third
Namespace:    vault-test
Labels:       <none>
Annotations:  vault-secret=secret/data/k8s/first

Type:  Opaque

Data
====
one:  8 bytes
two:  8 bytes
```

## Example - with failed authentication

ALLOW_FAIL="false" set for _vault-kubernetes-authenticator_
```
$ k logs vault-kubernetes-synchronizer-6875c88858-mbdsp -c vault-kubernetes-authenticator
2018/11/26 15:26:01 authentication failed: login failed with role from environment variable VAULT_ROLE: "k8s-np-appl-vault-dev-e1-auth": Put http://vault-dev-server.appl-vault-dev-e1.svc.cluster.local:8200/v1/auth/k8s-np/login: dial tcp 10.127.21.136:8200: i/o timeout

$ k logs vault-kubernetes-synchronizer-6875c88858-mbdsp -c vault-kubernetes-synchronizer
Error from server (BadRequest): container "vault-kubernetes-synchronizer" in pod "vault-kubernetes-synchronizer-6875c88858-mbdsp" is waiting to start: PodInitializing

$ k get pods
NAME                                             READY   STATUS                  RESTARTS   AGE
vault-kubernetes-synchronizer-6875c88858-mbdsp   0/1     Init:CrashLoopBackOff   3          7m40s
```

ALLOW_FAIL="true" set for _vault-kubernetes-authenticator_
```
$ k logs vault-kubernetes-synchronizer-7d5f65895-2pf4j -c vault-kubernetes-authenticator -f
2018/11/26 15:36:53 authentication failed - ALLOW_FAIL is set therefore pod will continue: login failed with role from environment variable VAULT_ROLE: "k8s-np-appl-vault-dev-e1-auth": Put http://vault-dev-server.appl-vault-dev-e1.svc.cluster.local:8200/v1/auth/k8s-np/login: dial tcp 10.127.21.136:8200: i/o timeout

$ k logs vault-kubernetes-synchronizer-7d5f65895-2pf4j -c vault-kubernetes-synchronizer
2018/11/26 15:36:55 check secret second from vault secret secret/data/k8s-np/appl-vault-dev-e1/second
2018/11/26 15:36:55 check secret third from vault secret secret/data/k8s-np/appl-vault-dev-e1/first
2018/11/26 15:36:55 check secret first from vault secret secret/data/k8s-np/appl-vault-dev-e1/first
2018/11/26 15:36:55 cannot synchronize secrets - all secrets seems to be available therefore pod creation will continue: could not get vault token: open /home/vault/.vault-token: no such file or directory

$ k get pods
NAME                                            READY   STATUS    RESTARTS   AGE
vault-kubernetes-synchronizer-7d5f65895-2pf4j   1/1     Running   0          5m18s
```


# Sidecar _vault-kubernetes-token-renewer_

Depends on Init Container _vault-kubernetes-authenticator_

- renew the Vault token regularly

## Configuration

- VAULT_TOKEN_PATH - the destination path on disk to store the token. Usually this is a shared volume.
- VAULT_REAUTH - re-authenticate if the token is invalid (default: "false")
- VAULT_TTL - requested token ttl (can be overwritten by Vault)

## Example

```
$ k logs vault-kubernetes-token-renewer-844488f7bc-c6ztf -c vault-kubernetes-authenticator
2018/11/26 14:56:30 successfully authenticated to vault
2018/11/26 14:56:30 successfully stored vault token at /home/vault/.vault-token

$ k logs vault-kubernetes-token-renewer-844488f7bc-c6ztf  -c vault-kubernetes-token-renewer
2018/11/26 14:56:32 start renewer loop
2018/11/26 14:56:32 token renewed
```


# Build

Install [mage](https://magefile.org/)

> The `DOCKER_TARGET` environment variable will be used to tag and push the images. If not set, the images will not be tagged and pushed.

```
$ export GO111MODULE=on
$ export DOCKER_TARGET="registry.example.com/repopath"
$ mage buildAllImages
```


# Demo

- Edit `profile`

```
cd demo
./deploy.sh profile
...
./delete.sh namespace
```


# Links

- [Using HashiCorp Vault with Kubernetes (Cloud Next '18)](https://www.youtube.com/watch?v=B16YTeSs1hI)
- [Github - vault-kubernetes-authenticator](https://github.com/sethvargo/vault-kubernetes-authenticator)
- [Vault - Kubernetes Auth Method](https://www.vaultproject.io/docs/auth/kubernetes.html)
- [Kubernetes - Init Containers](https://kubernetes.io/docs/concepts/workloads/pods/init-containers)
