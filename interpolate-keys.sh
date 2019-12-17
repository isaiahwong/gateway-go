#!/bin/bash
# TODO: Update CN with dynamic properties

basedir="$(dirname "$0")/manifest"
keydir="$(dirname "$0")/webhook-cert"

if [[ $1 == "production" ]]
 then
  echo 'gen production'
  basedir="$(dirname "$0")/release"
fi

# set -euo pipefail

[[ -d $keydir ]] && rm -r $keydir
mkdir $keydir

# Generate keys into a temporary directory.
echo "Generating TLS keys ..."
# "${basedir}/gen-key.sh" "$keydir"

chmod 0700 $keydir
cd $keydir

# Generate the CA cert and private key
openssl req -nodes -new -x509 -keyout ca.key -out ca.crt -subj "/CN=gateway-service.default.svc"
# Generate the private key for the webhook server
openssl genrsa -out webhook-server-tls.key 2048
# Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
openssl req -new -key webhook-server-tls.key -subj "/CN=gateway-service.default.svc" \
    | openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -out webhook-server-tls.crt

cd ..

ca_pem_b64="$(openssl base64 -A <"${keydir}/ca.crt")"
tls_crt="$(cat ${keydir}/webhook-server-tls.crt | base64)"
tls_key="$(cat ${keydir}/webhook-server-tls.key | base64)"

# Replaces the crt and key in gateway.yaml
find $basedir -name "gateway.yaml" -exec sed -i '' -e "s/tls.crt:.*/tls.crt: ${tls_crt}/g; s/tls.key:.*/tls.key: ${tls_key}/g;" {} +;

# Replaces the cabundle in gateway.yaml
find $basedir -name "gateway.yaml" -exec sed -i '' -e"s/caBundle:.*/caBundle: ${ca_pem_b64}/g;" {} +;

rm -rf webhook-cert

echo "$basedir/gateway.yaml has been interpolated with keys"
