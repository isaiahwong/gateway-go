FROM golang:1.13-alpine as builder

RUN apk add --update nodejs npm
RUN apk add --update npm
# RUN apk add --update curl
# RUN apk add --update git
# RUN apk add --update unzip

# RUN go get -u github.com/golang/protobuf/protoc-gen-go
# RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
# RUN curl -Lo protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.11.4/protoc-3.11.4-linux-x86_64.zip
# RUN unzip protoc.zip

# RUN export PATH=$PATH:$(pwd)/bin

WORKDIR /gateway

COPY go.mod . 
COPY go.sum .
# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

RUN go run main.go -b

WORKDIR /gateway/cmd
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/gateway

FROM alpine
COPY --from=builder /go/bin/gateway /go/bin/gateway

ENTRYPOINT ["/go/bin/gateway"]