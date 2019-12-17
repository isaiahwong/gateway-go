package server

import (
	"context"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/isaiahwong/gateway-go/proto-gen/auth"

	"github.com/isaiahwong/gateway-go/proto-gen/payment"

	"google.golang.org/grpc"
)

func GetProtos() map[string]func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error {
	services := map[string]func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error{}

	services["AuthService"] = auth.RegisterAuthServiceHandler

	services["PaymentService"] = payment.RegisterPaymentServiceHandler
	return services
}
