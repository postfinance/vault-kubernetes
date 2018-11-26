#!/bin/bash

FILE="bootstrap-dev-vault"
envsubst < $(dirname $0)/bootstrap-dev-vault.envsubst > ${FILE}
kubectl --namespace=${NAMESPACE} create configmap bootstrap-scripts --from-file=${FILE}
rm -rf ${FILE}
