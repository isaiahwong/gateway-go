#!/bin/bash
set -e

dir="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
basedir="$dir/../k8s"
keydir="$dir/webhook-cert"

check_root() {
  local root_secret
  root_secret=$(kubectl get secret webhook-server-tls -o yaml \
  | sed -n 's/^.*tls.crt: //p')
  if [ -z "${root_secret}" ]; then
    echo "Root secret is empty. Are you using the self-signed CA?"
    return
  fi

  echo "Fetching root cert from istio-system namespace..."
  kubectl get secret webhook-server-tls -o yaml | \
    awk '/ca.crt/ {print $2}' | base64 --decode > ca.crt
  if [[ ! -f ./ca.crt ]]; then
    echo "failed to get cacert, check the istio installation namespace."
    return
  fi

  local root_date
  local root_sec
  root_date=$(openssl x509 -in ca.crt -noout -enddate | cut -f2 -d'=')
  if [[ "$(uname)" == "Darwin" ]]; then
    root_sec=$(date -jf "%b  %e %k:%M:%S %Y %Z" "${root_date}" '+%s')
  else
    root_sec=$(date -d "${root_date}" '+%s')
  fi

  local now_sec
  local days_left
  now_sec=$(date '+%s')
  days_left=$(echo "(${root_sec} - ${now_sec}) / (3600 * 24)" | bc)

  rm ca.crt

  cat << EOF
Your Root Cert will expire after
   ${root_date}
Current time is
  $(date)


=====YOU HAVE ${days_left} DAYS BEFORE THE ROOT CERT EXPIRES!=====

EOF
}

generate_keys() {
  # Generate keys into a temporary directory.
  echo "Generating TLS keys ..."
  # Generate the CA cert and private key
  openssl req -nodes -new -x509 -sha256 -days 3650 -keyout ca.key -out ca.crt -subj "/O=$service"
  # Generate the private key for the webhook server
  openssl genrsa -out webhook-server-tls.key 2048
  # Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
  openssl req -new -key webhook-server-tls.key -subj "/CN=$service.$namespace.svc" \
      | openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -out webhook-server-tls.crt
}

gen_cert() {
  if [[ $release == true ]]
  then
    echo 'Release production'
    # Create manifest file
    tmpdir="$dir/../release"
    cp $basedir/template.yaml $tmpdir/$out
    basedir=$tmpdir
  else
    # Create manifest file
    cp $basedir/template.yaml $basedir/$out
  fi

  if [[ $nginx == true ]]
  then
  # Replaces ingress.yaml service. USE ONLY IN DEV
      find $basedir -name "ingress.yaml" \
      -exec sed -i '' -e \
      "s/serviceName:.*/serviceName: ${service}/g;" {} +;
  fi


    # Deletes any existing key directories 
  [[ -d $keydir ]] && rm -r $keydir
  mkdir $keydir

  chmod 0700 $keydir
  cd $keydir

  generate_keys

  ca_pem_b64="$(openssl base64 -A <"ca.crt")"
  tls_crt="$(cat webhook-server-tls.crt | base64)"
  tls_key="$(cat webhook-server-tls.key | base64)"

  cd ..

  # Replaces the respective values crt, key and caBundle in output file
  find $basedir -name $out \
  -exec sed -i '' -e \
  "s/{{TLS_CRT}}/${tls_crt}/g;\
  s/{{TLS_KEY}}/${tls_key}/g;\
  s/{{CA_CRT}}/${ca_pem_b64}/g;\
  s/{{SERVICE_NAME}}/${service}/g;\
  s/{{NAMESPACE}}/${namespace}/g;" {} +;

  # Deletes key dir
  rm -rf $keydir

  echo "$basedir/$out has been interpolated with keys"
}

case $1 in
  check-root)
    check_root
    ;;

  gen-cert)
  while [[ $# -gt 0 ]]; do
    case ${1} in
      --service)
          service="$2"
          shift
          ;;
      --namespace)
          namespace="$2"
          shift
          ;;
      --release)
          release="$2"
          shift
          ;;
      --nginx)
          nginx="$2"
          shift
          ;;
      --out)
          out="$2"
          shift
          ;;
      *)
      esac
      shift
    done
    [ -z ${service} ] && service=gateway-service
    [ -z ${namespace} ] && namespace=default
    [ -z ${nginx} ] && nginx=false
    [ -z ${release} ] && release=false
    [ -z ${out} ] && out=gateway.yaml

    if [ ! -x "$(command -v openssl)" ]; then
        echo "openssl not found"
        exit 1
    fi
    gen_cert
    ;;
  *)
    echo "Usage: check-root | gen-cert

check-root
  Check the expiration date of the current root certificate.

gen-cert
  This will replace the current mutating webhook
  certificate with a new 10-year lifetime root certificate.
"

esac
