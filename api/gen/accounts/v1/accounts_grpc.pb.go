// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package accounts

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AccountsServiceClient is the client API for AccountsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountsServiceClient interface {
	LoginWithChallenge(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HydraResponse, error)
	ConsentWithChallenge(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*RedirectResponse, error)
	Introspect(ctx context.Context, in *IntrospectRequest, opts ...grpc.CallOption) (*IntrospectResponse, error)
	AccountExists(ctx context.Context, in *AccountExistsRequest, opts ...grpc.CallOption) (*AccountExistsResponse, error)
	IsAuthenticated(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*AuthenticateResponse, error)
	SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*RedirectResponse, error)
	Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*RedirectResponse, error)
	EmailExists(ctx context.Context, in *EmailExistsRequest, opts ...grpc.CallOption) (*EmailExistsResponse, error)
}

type accountsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountsServiceClient(cc grpc.ClientConnInterface) AccountsServiceClient {
	return &accountsServiceClient{cc}
}

func (c *accountsServiceClient) LoginWithChallenge(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HydraResponse, error) {
	out := new(HydraResponse)
	err := c.cc.Invoke(ctx, "/api.accounts.v1.AccountsService/LoginWithChallenge", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountsServiceClient) ConsentWithChallenge(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*RedirectResponse, error) {
	out := new(RedirectResponse)
	err := c.cc.Invoke(ctx, "/api.accounts.v1.AccountsService/ConsentWithChallenge", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountsServiceClient) Introspect(ctx context.Context, in *IntrospectRequest, opts ...grpc.CallOption) (*IntrospectResponse, error) {
	out := new(IntrospectResponse)
	err := c.cc.Invoke(ctx, "/api.accounts.v1.AccountsService/Introspect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountsServiceClient) AccountExists(ctx context.Context, in *AccountExistsRequest, opts ...grpc.CallOption) (*AccountExistsResponse, error) {
	out := new(AccountExistsResponse)
	err := c.cc.Invoke(ctx, "/api.accounts.v1.AccountsService/AccountExists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountsServiceClient) IsAuthenticated(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*AuthenticateResponse, error) {
	out := new(AuthenticateResponse)
	err := c.cc.Invoke(ctx, "/api.accounts.v1.AccountsService/IsAuthenticated", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountsServiceClient) SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*RedirectResponse, error) {
	out := new(RedirectResponse)
	err := c.cc.Invoke(ctx, "/api.accounts.v1.AccountsService/SignUp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountsServiceClient) Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*RedirectResponse, error) {
	out := new(RedirectResponse)
	err := c.cc.Invoke(ctx, "/api.accounts.v1.AccountsService/Authenticate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountsServiceClient) EmailExists(ctx context.Context, in *EmailExistsRequest, opts ...grpc.CallOption) (*EmailExistsResponse, error) {
	out := new(EmailExistsResponse)
	err := c.cc.Invoke(ctx, "/api.accounts.v1.AccountsService/EmailExists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountsServiceServer is the server API for AccountsService service.
// All implementations must embed UnimplementedAccountsServiceServer
// for forward compatibility
type AccountsServiceServer interface {
	LoginWithChallenge(context.Context, *Empty) (*HydraResponse, error)
	ConsentWithChallenge(context.Context, *Empty) (*RedirectResponse, error)
	Introspect(context.Context, *IntrospectRequest) (*IntrospectResponse, error)
	AccountExists(context.Context, *AccountExistsRequest) (*AccountExistsResponse, error)
	IsAuthenticated(context.Context, *Empty) (*AuthenticateResponse, error)
	SignUp(context.Context, *SignUpRequest) (*RedirectResponse, error)
	Authenticate(context.Context, *AuthenticateRequest) (*RedirectResponse, error)
	EmailExists(context.Context, *EmailExistsRequest) (*EmailExistsResponse, error)
	mustEmbedUnimplementedAccountsServiceServer()
}

// UnimplementedAccountsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAccountsServiceServer struct {
}

func (UnimplementedAccountsServiceServer) LoginWithChallenge(context.Context, *Empty) (*HydraResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginWithChallenge not implemented")
}
func (UnimplementedAccountsServiceServer) ConsentWithChallenge(context.Context, *Empty) (*RedirectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConsentWithChallenge not implemented")
}
func (UnimplementedAccountsServiceServer) Introspect(context.Context, *IntrospectRequest) (*IntrospectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Introspect not implemented")
}
func (UnimplementedAccountsServiceServer) AccountExists(context.Context, *AccountExistsRequest) (*AccountExistsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AccountExists not implemented")
}
func (UnimplementedAccountsServiceServer) IsAuthenticated(context.Context, *Empty) (*AuthenticateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsAuthenticated not implemented")
}
func (UnimplementedAccountsServiceServer) SignUp(context.Context, *SignUpRequest) (*RedirectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}
func (UnimplementedAccountsServiceServer) Authenticate(context.Context, *AuthenticateRequest) (*RedirectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
}
func (UnimplementedAccountsServiceServer) EmailExists(context.Context, *EmailExistsRequest) (*EmailExistsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmailExists not implemented")
}
func (UnimplementedAccountsServiceServer) mustEmbedUnimplementedAccountsServiceServer() {}

// UnsafeAccountsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountsServiceServer will
// result in compilation errors.
type UnsafeAccountsServiceServer interface {
	mustEmbedUnimplementedAccountsServiceServer()
}

func RegisterAccountsServiceServer(s grpc.ServiceRegistrar, srv AccountsServiceServer) {
	s.RegisterService(&AccountsService_ServiceDesc, srv)
}

func _AccountsService_LoginWithChallenge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountsServiceServer).LoginWithChallenge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.accounts.v1.AccountsService/LoginWithChallenge",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountsServiceServer).LoginWithChallenge(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountsService_ConsentWithChallenge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountsServiceServer).ConsentWithChallenge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.accounts.v1.AccountsService/ConsentWithChallenge",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountsServiceServer).ConsentWithChallenge(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountsService_Introspect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IntrospectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountsServiceServer).Introspect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.accounts.v1.AccountsService/Introspect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountsServiceServer).Introspect(ctx, req.(*IntrospectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountsService_AccountExists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountExistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountsServiceServer).AccountExists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.accounts.v1.AccountsService/AccountExists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountsServiceServer).AccountExists(ctx, req.(*AccountExistsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountsService_IsAuthenticated_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountsServiceServer).IsAuthenticated(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.accounts.v1.AccountsService/IsAuthenticated",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountsServiceServer).IsAuthenticated(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountsService_SignUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignUpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountsServiceServer).SignUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.accounts.v1.AccountsService/SignUp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountsServiceServer).SignUp(ctx, req.(*SignUpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountsService_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountsServiceServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.accounts.v1.AccountsService/Authenticate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountsServiceServer).Authenticate(ctx, req.(*AuthenticateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountsService_EmailExists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailExistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountsServiceServer).EmailExists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.accounts.v1.AccountsService/EmailExists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountsServiceServer).EmailExists(ctx, req.(*EmailExistsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountsService_ServiceDesc is the grpc.ServiceDesc for AccountsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccountsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.accounts.v1.AccountsService",
	HandlerType: (*AccountsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LoginWithChallenge",
			Handler:    _AccountsService_LoginWithChallenge_Handler,
		},
		{
			MethodName: "ConsentWithChallenge",
			Handler:    _AccountsService_ConsentWithChallenge_Handler,
		},
		{
			MethodName: "Introspect",
			Handler:    _AccountsService_Introspect_Handler,
		},
		{
			MethodName: "AccountExists",
			Handler:    _AccountsService_AccountExists_Handler,
		},
		{
			MethodName: "IsAuthenticated",
			Handler:    _AccountsService_IsAuthenticated_Handler,
		},
		{
			MethodName: "SignUp",
			Handler:    _AccountsService_SignUp_Handler,
		},
		{
			MethodName: "Authenticate",
			Handler:    _AccountsService_Authenticate_Handler,
		},
		{
			MethodName: "EmailExists",
			Handler:    _AccountsService_EmailExists_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "accounts/v1/accounts.proto",
}
