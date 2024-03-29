$(shell PATH=$PATH:$GOPATH/bin)
BUILD_ID := $(shell git rev-parse --short HEAD 2>/dev/null || echo no-commit-id)
IMAGE_NAME := registry.gitlab.com/eco_system/gateway
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

set-image:
	kubectl set image deployments/gateway-deployment gateway=$(IMAGE_NAME)

set-image-latest:
	kubectl set image deployments/gateway-deployment gateway=$(IMAGE_NAME):latest

set-image-sha:
	kubectl set image deployments/gateway-deployment gateway=$(IMAGE_NAME):${BUILD_ID}

gen-manifest:
	./tools/gen-manifest.sh gen-cert --image $(IMAGE_NAME)

gen-manifest-release:
	./tools/gen-manifest.sh gen-cert --release true --image $(IMAGE_NAME)

genhealth:
	protoc --go_out=plugins=grpc:proto -I $(PROTO_DIR) $(PROTO_DIR)/health.proto

genproto:
	go run main.go -b -m accounts/map.json

genproto-manual:
	if [ ! -d "api/go" ]; then \
			mkdir api/go; \
	fi

	if [ ! -d "api/go/gen" ]; then \
			mkdir api/go/gen; \
	fi

	protoc \
		-I./api/ecosystem \
		-I./api/third_party/googleapis \
		--go_out ./api/go/gen \
		--go-grpc_out ./api/go/gen \
		--go_opt paths=source_relative \
		--go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./api/go/gen \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		./api/ecosystem/accounts/v1/*.proto
		

apply-registry-secret:
	kubectl create secret docker-registry gitlab-auth --docker-server=registry.gitlab.com --docker-username=$DUSER --docker-password=$DPW --docker-email=iisaiah.dev@gmail.com

