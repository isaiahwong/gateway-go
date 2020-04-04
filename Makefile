$(shell PATH=$PATH:$GOPATH/bin)
BUILD_ID := $(shell git rev-parse --short HEAD 2>/dev/null || echo no-commit-id)
IMAGE_NAME := registry.gitlab.com/isaiahwong/cluster/gateway
VERSION := 0.0.1

PROTO_DIR := ../../pb

.DEFAULT_GOAL := help

help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build:
	docker build -t $(IMAGE_NAME):latest . --rm=true

build-sha: 
	docker build -t $(IMAGE_NAME):$(BUILD_ID) . --rm=true

push: 
	docker push $(IMAGE_NAME):latest

push-sha:
	docker push $(IMAGE_NAME):$(BUILD_ID)

build-all:
	docker build -t $(IMAGE_NAME):latest -t $(IMAGE_NAME):$(BUILD_ID) . --rm=true

push-all:
	make push
	make push-sha

build-push:
	make build-all
	make push-all

clean:
	docker rmi $( docker images | grep '<none>') --force 2>/dev/null

gen-manifest-release:
	./tools/gen-manifest.sh gen-cert --release true --image registry.gitlab.com/isaiahwong/cluster/gateway

genhealth:
	protoc --go_out=plugins=grpc:proto -I $(PROTO_DIR) $(PROTO_DIR)/health.proto

genproto:
	go run main.go -b -m proto/_map.json

genproto-manual:
	if [ ! -d "protogen" ]; then \
			mkdir protogen; \
	fi

	protoc \
		-I./proto/accounts-proto/api \
		-I$(GOPATH)/src \
		-I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true,allow_repeated_fields_in_body=true:protogen \
		--swagger_out=logtostderr=true:protogen \
		--go_out=plugins=grpc:./protogen \
		./proto/accounts-proto/api/accounts/v1/*.proto

