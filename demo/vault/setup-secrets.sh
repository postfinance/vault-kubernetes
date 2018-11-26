#!/bin/bash

FILE=truststore.pem

# create secret kubernetes_ca_cert from truststore.crt
if [ ! -f "${TRUSTSTORE}" ]; then
    echo "truststore \"${TRUSTSTORE}\" for kubernetes api server does not exist - environment variable TRUSTSTORE set?"
    exit 1
fi
cp ${TRUSTSTORE} ${FILE}
kubectl --namespace=${NAMESPACE} create secret generic kubernetes-ca-cert --from-file=${FILE}
rm -f ${FILE}

# extract ca certificate of service account and create secret pem-keys
kubectl --namespace=${NAMESPACE} get secrets $(kubectl get sa vault-auth -o=jsonpath="{.secrets[0].name}") -o=jsonpath="{.data['ca\.crt']}" | base64 --decode > ${FILE}
kubectl --namespace=${NAMESPACE} create secret generic pem-keys --from-file=${FILE}
rm -f ${FILE}

