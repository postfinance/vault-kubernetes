#!/bin/bash

set -e

if [ $# -eq 0 ]; then
    echo "usage: ${0} profile"
    exit 1
fi

if [ ! -r $1 ]; then
    echo "ERROR: failed to read profile \"$1\""
    exit 1
fi

source $1

kubectl config use-context ${CLUSTER}

# setup rbac
envsubst < ../demo/k8s/rbac.yaml | kubectl sudo --namespace=${NAMESPACE} apply -f -

# setup vault
bash ../demo/vault/setup-secrets.sh
bash ../demo/vault/setup-configmap.sh
envsubst < ../demo/vault/deployment.yaml | kubectl --namespace=${NAMESPACE} apply -f -

# wait for vault
while ! $(kubectl get pods --field-selector=status.phase=Running 2> /dev/null | grep -q '^vault-dev-server'); do
    echo "> vault-dev-server not yet running - waiting"
    sleep 5
done

# setup demo deployments
{
    envsubst < ../demo/k8s/authenticator/deployment.yaml
    envsubst < ../demo/k8s/synchronizer/deployment.yaml
    envsubst < ../demo/k8s/token-renewer/deployment.yaml
} | kubectl --namespace=${NAMESPACE} apply -f -


