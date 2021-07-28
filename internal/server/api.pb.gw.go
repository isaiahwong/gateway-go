// Code generated by gitlab.com/eco_system/gateway/hack/genproto. DO NOT EDIT.
package server

import (
	"context"
	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.com/eco_system/gateway/api/go/gen/accounts/v1"
	"gitlab.com/eco_system/gateway/api/go/gen/fitness/v1"
	"google.golang.org/grpc"
)

type ServiceDesc struct {
	ServiceName    string
	PackageSVC     string
	Package        string
	CurrentPackage string
	Handler        func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
}

var _Services = map[string]ServiceDesc{
	"api.accounts.v1.AccountsService": {
		ServiceName:    "AccountsService",
		PackageSVC:     "api.accounts.v1.AccountsService",
		Package:        "api.accounts.v1",
		CurrentPackage: "accounts",
		Handler:        accounts.RegisterAccountsServiceHandler,
	},
	"api.fitness.v1.FitnessService": {
		ServiceName:    "FitnessService",
		PackageSVC:     "api.fitness.v1.FitnessService",
		Package:        "api.fitness.v1",
		CurrentPackage: "fitness",
		Handler:        fitness.RegisterFitnessServiceHandler,
	},
}

// Returns generated protos that have been generated with  protoc-gen-grpc-gateway
func GetProtos() map[string]ServiceDesc {
	return _Services
}
