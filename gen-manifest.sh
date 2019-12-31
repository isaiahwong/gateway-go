#!/bin/bash
# TODO: Update CN with dynamic properties

basedir="$(dirname "$0")/manifest"
keydir="$(dirname "$0")/webhook-cert"
out="gateway.yaml"

# K8S Manifest values
namespace="default"
service_name="gateway-service"

# Create manifest file
cp $basedir/template.yaml $basedir/$out

if [[ $1 == "production" ]]
  then
    echo 'gen production'
    basedir="$(dirname "$0")/release"
  else
    # Replaces ingress.yaml service. USE ONLY IN DEV
    find $basedir -name "ingress.yaml" \
    -exec sed -i '' -e \
    "s/serviceName:.*/serviceName: ${service_name}/g;" {} +;
fi

set -euo pipefail

# Deletes any existing key directories 
[[ -d $keydir ]] && rm -r $keydir
mkdir $keydir

# Generate keys into a temporary directory.
echo "Generating TLS keys ..."

chmod 0700 $keydir
cd $keydir

# Generate the CA cert and private key
openssl req -nodes -new -x509 -keyout ca.key -out ca.crt -subj "/CN=$service_name.$namespace.svc"
# Generate the private key for the webhook server
openssl genrsa -out webhook-server-tls.key 2048
# Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
openssl req -new -key webhook-server-tls.key -subj "/CN=$service_name.$namespace.svc" \
    | openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -out webhook-server-tls.crt

cd ..

ca_pem_b64="$(openssl base64 -A <"${keydir}/ca.crt")"
tls_crt="$(cat ${keydir}/webhook-server-tls.crt | base64)"
tls_key="$(cat ${keydir}/webhook-server-tls.key | base64)"

# Replaces the respective values crt, key and caBundle in output file
find $basedir -name $out \
-exec sed -i '' -e \
"s/tls.crt:.*/tls.crt: ${tls_crt}/g;\
 s/tls.key:.*/tls.key: ${tls_key}/g;\
 s/caBundle:.*/caBundle: ${ca_pem_b64}/g;\
 s/{{SERVICE_NAME}}/${service_name}/g;\
 s/{{NAMESPACE}}/${namespace}/g;" {} +;

# Deletes key dir
rm -rf $keydir

echo "$basedir/$out has been interpolated with keys"
