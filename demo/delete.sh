#!/bin/sh

kubectl delete deployment vault-kubernetes-authenticator vault-kubernetes-synchronizer vault-kubernetes-token-renewer
kubectl delete secret first second third

kubectl delete deployment vault-dev-server
kubectl delete service vault-dev-server
kubectl delete configmap bootstrap-scripts
kubectl delete secret kubernetes-ca-cert pem-keys
