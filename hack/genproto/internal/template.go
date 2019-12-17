package internal

const ProtoTemplate = `
package server

import (
	"context"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	{{range $i, $e := .}}
		"github.com/isaiahwong/gateway-go/proto-gen/{{ $e.Package }}"
	{{end}}
	"google.golang.org/grpc"
)

func GetProtos() map[string]func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error {
	services := map[string]func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error{}
	{{range $i, $e := .}}
		services["{{ $e.OriginalName }}"] = {{ $e.Package }}.Register{{ $e.OriginalName}}Handler
	{{end}}
	return services
}
`
