
// Code generated by github.com/isaiahwong/gateway-go/hack/genproto. DO NOT EDIT.
package protogen

import (
	"context"
	runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"github.com/isaiahwong/gateway-go/protogen/accounts/v1"
)

type ServiceDesc struct {
	ServiceName    string
	PackageSVC     string
	Package        string
	CurrentPackage string
	Handler        func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
}

var _Services = map[string]ServiceDesc{
	"api.accounts.v1.AccountsService": ServiceDesc{
		ServiceName: "AccountsService",
		PackageSVC: "api.accounts.v1.AccountsService",
		Package:        "api.accounts.v1",
		CurrentPackage: "accounts",
		Handler: accounts.RegisterAccountsServiceHandler,
	},
}

// Returns generated protos that have been generated with  protoc-gen-grpc-gateway
func GetProtos() map[string]ServiceDesc {
	return _Services
}
