#!/bin/sh

if [ $# -eq 0 ]; then
    echo "usage: $0 namespace"
    exit 1
fi
NAMESPACE=$1

kubectl --namespace=${NAMESPACE} delete deployment vault-kubernetes-authenticator vault-kubernetes-synchronizer vault-kubernetes-token-renewer
kubectl --namespace=${NAMESPACE} delete secret first second third

kubectl --namespace=${NAMESPACE} delete deployment vault-dev-server
kubectl --namespace=${NAMESPACE} delete service vault-dev-server
kubectl --namespace=${NAMESPACE} delete configmap bootstrap-scripts
kubectl --namespace=${NAMESPACE} delete secret kubernetes-ca-cert pem-keys
