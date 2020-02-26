# WIP
Porting gateway node to gateway go

# Gateway-go
Gateway-go provides a single [`REST`](https://en.m.wikipedia.org/wiki/Representational_state_transfer ) entry point which forwards requests to your microservices in your [Kubernetes][k8s] cluster. It discovers your [services][k8s-service] that resides in your cluster to which are exposed via **HTTP** or **gRPC**. Incoming requests will then be proxied to the respective services. Gateway works on (local) and GCP's Google Kubernetes Engine. 

Gateway-go utilises [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) to generate gRPC stubs that map to [`google.api.http`](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto#L46) annotations. Portions of it’s [`README.md`](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/README.md) has been included in the documentation.


# Table of Contents
* [Installation](#Installation)
* [Gateway](#gateway)
* [Directory Layout](#Directory-Layout)
* [Technology Stack](#technology-stack)

# Installation
The gateway-go requires a local installation of the Google protocol buffers compiler `protoc`.

[Link: protoc](https://github.com/protocolbuffers/protobuf/releases)

Install the following packages with `go get -u`

```
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/golang/protobuf/protoc-gen-go
```

Ensure `$GOBIN` is configured

# Usage
Example proto files are included in the project.


1. Define your protos in `proto/api`

You can separate your proto messages into different files. In the example below, where datatype Payment is in `proto/api/payment/schema.proto` 

`proto/api/payment`:

```
syntax = "proto3";

package api.payment;

option go_package = "payment";

import "payment/schema.proto";
import "payment/paypal.proto";
import "payment/stripe.proto";

import "google/api/annotations.proto";

service PaymentService {
  rpc CreatePayment(CreatePaymentRequest) returns (CreatePaymentResponse) {
    option (google.api.http) = {
      post: "/v1/payment/create",
      body: "*"
    };
  };

  message CreatePaymentRequest {
   string user = 1;
   string email = 2;
  }

  message CreatePaymentResponse {
   bool success = 1;
   Payment payment = 2;
  }
}

```

2\. Generate gRPC stub

The following generates gRPC code for Golang based on path/to/your_service.proto:

```
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I ./proto/third_party/googleapis\
  --go_out=plugins=grpc:. path/to/your_service.proto
```

The following code generates all proto files in the folder `proto/api/payment`

```
protoc -I./proto/api -I./proto/third_party/googleapis --go_out=plugins=grpc:./protogen ./proto/api/payment/*.proto
```


# Apply gateway deployment, services and environment files
kubectl apply -f ./k8s/

```

# Webhook
....

# Configuring the gateway with environment variables
The environment variables can be defined in the `k8s/env.yaml ` or `/src/.env` file.  
*Environment variables defined in the `env.yaml` will override the `.env` file*  

```
...
    spec:
      serviceAccountName: api-gateway
      containers:
        - name: gateway
          image: registry.gitlab.com/isaiahwong/kinddd/api/gateway
          imagePullPolicy: IfNotPresent
          ports:
            - name: gateway-port
              containerPort: 5000
          # Edit this portion
          env:
          - name: NODE_ENV
            value: "development"
          - name: MAINTENANCE_MODE
            value: "false"
          - name: SVC_DISCOVERY_INTERVAL
            value: "5000"
          - name: ENABLE_CONSOLE_LOGS_IN_PROD
            value: "false"
          - name: ENABLE_CONSOLE_LOGS_IN_TEST
            value: "true"
...
```

## Environment Variable types

| Variable | Default Value | Description | 
| -------- | ------------- | ----------- | 
| `PORT` | `5000` | Defines which port the `gateway` will run on. <br/> **Note**: *It is not advisable to change the `PORT` when using it with kubernetes. If you have to, do remember to amend the port that binds to the container* |
| `NODE_ENV` | `Development` | Defines if the application is running in production or development. |
| `SVC_DISCOVERY_INTERVAL` | `5000` Milli | How often the gateway will poll Kubernetes [Service Discovery][k8s-svc-discovery] |


# Preparing services to connect with the gateway
```
apiVersion: v1
kind: Service
metadata:
  name: payment-service
  labels:
    resourceType: api-service
  annotations:
    config: '{
      "expose": true,
      "serviceType": "resource",
      "path": "payment",
      "apiVersion": "v1",
      "authentication": {
        "required": "true",
        "exclude": [
          "/api/v1/payment/stripe/webhook/paymentintent",
          "/api/v1/payment/stripe/webhook/test",
          "/api/v1/payment/paypal/webhook/order",
          "/api/v1/payment/paypal/webhook/test"
        ]
      }
    }'
```

# Clean up
```
skaffold delete 

./clean.sh
```

# Directory Layout
```
.
├── /docs/                      # Documentation for the gateway
│
├── /k8s/                       # Development Kubernetes manifest files    
│   ├── /api/                   # Api Services k8s
│   ├── /nginx-ingress/         # Kubernetes ingress config files
│   ├── gateway.yaml            # Gateway deployment and service
│   ├── env.yaml                # Gateway env variables
│   └── ...                     # Other config files 
│
├── /locales/                   # locales configs
│   ├── /en/                    # English locales
│
├── /proto/                     # Protocol Buffers Descriptions (Services)
│
├── /release/                   # Production Kubernetes manifest files    
│   ├── /api/                   # Api Services k8s
│   ├── /nginx-ingress/         # Kubernetes ingress config files
│   ├── gateway.yaml            # Gateway deployment and service
│   ├── env.yaml                # Gateway env variables
│   └── ...                     # Other config files 
│
├── /test/                      # Test files (WIP)
│
├── /src/                       # Api Gateway source code
│
├── clean.sh                    #  Clean dangling images
│
└── skaffold.yaml               #  Skaffold config for development
```

Running without Kubernetes i.e `npm run dev`  
Rename `/src/.env.example` to `/src/.env`
```
# Main

# K8S
SVC_DISCOVERY_INTERVAL=5000

# Config
MAINTENANCE_MODE=false

ENABLE_CONSOLE_LOGS_IN_PROD=false
ENABLE_CONSOLE_LOGS_IN_TEST=true
```

# Dependencies
`Gateway` utilise Kubernetes's [Service Discovery][k8s-svc-discovery] and depends on Kubernetes Client - [GoDaddy Client][godaddy-client]

[Isaiah]: https://www.iisaiah.com
[brew]: https://brew.sh/
[minikube]: https://github.com/kubernetes/minikube/releases/  
[vbox]: https://www.virtualbox.org/wiki/Downloads
[express]: https://github.com/expressjs/express

[node]: https://github.com/nodejs/node
[skaffold]: https://github.com/GoogleContainerTools/skaffold
[mailer]: https://nodemailer.com/

[godaddy-client]: https://github.com/godaddy/kubernetes-client
[ingress-nginx]: https://github.com/kubernetes/ingress-nginx
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[docker-desktop]: https://www.docker.com/products/docker-desktop
[k8s-service]: https://kubernetes.io/docs/concepts/services-networking/service/
[k8s]: https://github.com/kubernetes/kubernetes
[k8s-svc-discovery]: https://kubernetes.io/docs/tasks/administer-cluster/access-cluster-api/
