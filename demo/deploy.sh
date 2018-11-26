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

cd $(dirname $0)

kubectl config use-context ${CLUSTER}

# setup rbac
envsubst < k8s/rbac.yaml | kubectl --namespace=${NAMESPACE} apply -f -

# setup vault
bash vault/setup-secrets.sh
bash vault/setup-configmap.sh
envsubst < vault/deployment.yaml | kubectl --namespace=${NAMESPACE} apply -f -

# wait for vault
while ! $(kubectl --namespace=${NAMESPACE} get pods --field-selector=status.phase=Running 2> /dev/null | grep -q '^vault-dev-server'); do
    echo "> vault-dev-server not yet running - waiting"
    sleep 5
done

# setup demo deployments
{
    envsubst < k8s/authenticator/deployment.yaml
    envsubst < k8s/synchronizer/deployment.yaml
    envsubst < k8s/token-renewer/deployment.yaml
} | kubectl --namespace=${NAMESPACE} apply -f -


