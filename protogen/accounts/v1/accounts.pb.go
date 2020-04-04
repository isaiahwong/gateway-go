// Code generated by protoc-gen-go. DO NOT EDIT.
// source: accounts/v1/accounts.proto

package accounts

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type IntrospectRequest struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Scope                string   `protobuf:"bytes,2,opt,name=scope,proto3" json:"scope,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IntrospectRequest) Reset()         { *m = IntrospectRequest{} }
func (m *IntrospectRequest) String() string { return proto.CompactTextString(m) }
func (*IntrospectRequest) ProtoMessage()    {}
func (*IntrospectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8726aced901bdecf, []int{0}
}

func (m *IntrospectRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IntrospectRequest.Unmarshal(m, b)
}
func (m *IntrospectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IntrospectRequest.Marshal(b, m, deterministic)
}
func (m *IntrospectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IntrospectRequest.Merge(m, src)
}
func (m *IntrospectRequest) XXX_Size() int {
	return xxx_messageInfo_IntrospectRequest.Size(m)
}
func (m *IntrospectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IntrospectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IntrospectRequest proto.InternalMessageInfo

func (m *IntrospectRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *IntrospectRequest) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

type IntrospectResponse struct {
	Active               bool     `protobuf:"varint,1,opt,name=active,proto3" json:"active,omitempty"`
	Aud                  []string `protobuf:"bytes,2,rep,name=aud,proto3" json:"aud,omitempty"`
	ClientId             string   `protobuf:"bytes,3,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Exp                  int64    `protobuf:"varint,4,opt,name=exp,proto3" json:"exp,omitempty"`
	Ext                  []byte   `protobuf:"bytes,5,opt,name=ext,proto3" json:"ext,omitempty"`
	Iat                  int64    `protobuf:"varint,6,opt,name=iat,proto3" json:"iat,omitempty"`
	Iss                  string   `protobuf:"bytes,7,opt,name=iss,proto3" json:"iss,omitempty"`
	Nbf                  int64    `protobuf:"varint,8,opt,name=nbf,proto3" json:"nbf,omitempty"`
	ObfuscatedSubject    string   `protobuf:"bytes,9,opt,name=obfuscated_subject,json=obfuscatedSubject,proto3" json:"obfuscated_subject,omitempty"`
	Scope                string   `protobuf:"bytes,10,opt,name=scope,proto3" json:"scope,omitempty"`
	Sub                  string   `protobuf:"bytes,11,opt,name=sub,proto3" json:"sub,omitempty"`
	TokenType            string   `protobuf:"bytes,12,opt,name=token_type,json=tokenType,proto3" json:"token_type,omitempty"`
	Username             string   `protobuf:"bytes,13,opt,name=username,proto3" json:"username,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IntrospectResponse) Reset()         { *m = IntrospectResponse{} }
func (m *IntrospectResponse) String() string { return proto.CompactTextString(m) }
func (*IntrospectResponse) ProtoMessage()    {}
func (*IntrospectResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8726aced901bdecf, []int{1}
}

func (m *IntrospectResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IntrospectResponse.Unmarshal(m, b)
}
func (m *IntrospectResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IntrospectResponse.Marshal(b, m, deterministic)
}
func (m *IntrospectResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IntrospectResponse.Merge(m, src)
}
func (m *IntrospectResponse) XXX_Size() int {
	return xxx_messageInfo_IntrospectResponse.Size(m)
}
func (m *IntrospectResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_IntrospectResponse.DiscardUnknown(m)
}

var xxx_messageInfo_IntrospectResponse proto.InternalMessageInfo

func (m *IntrospectResponse) GetActive() bool {
	if m != nil {
		return m.Active
	}
	return false
}

func (m *IntrospectResponse) GetAud() []string {
	if m != nil {
		return m.Aud
	}
	return nil
}

func (m *IntrospectResponse) GetClientId() string {
	if m != nil {
		return m.ClientId
	}
	return ""
}

func (m *IntrospectResponse) GetExp() int64 {
	if m != nil {
		return m.Exp
	}
	return 0
}

func (m *IntrospectResponse) GetExt() []byte {
	if m != nil {
		return m.Ext
	}
	return nil
}

func (m *IntrospectResponse) GetIat() int64 {
	if m != nil {
		return m.Iat
	}
	return 0
}

func (m *IntrospectResponse) GetIss() string {
	if m != nil {
		return m.Iss
	}
	return ""
}

func (m *IntrospectResponse) GetNbf() int64 {
	if m != nil {
		return m.Nbf
	}
	return 0
}

func (m *IntrospectResponse) GetObfuscatedSubject() string {
	if m != nil {
		return m.ObfuscatedSubject
	}
	return ""
}

func (m *IntrospectResponse) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

func (m *IntrospectResponse) GetSub() string {
	if m != nil {
		return m.Sub
	}
	return ""
}

func (m *IntrospectResponse) GetTokenType() string {
	if m != nil {
		return m.TokenType
	}
	return ""
}

func (m *IntrospectResponse) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func init() {
	proto.RegisterType((*IntrospectRequest)(nil), "api.accounts.v2.IntrospectRequest")
	proto.RegisterType((*IntrospectResponse)(nil), "api.accounts.v2.IntrospectResponse")
}

func init() {
	proto.RegisterFile("accounts/v1/accounts.proto", fileDescriptor_8726aced901bdecf)
}

var fileDescriptor_8726aced901bdecf = []byte{
	// 343 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xc1, 0x6a, 0xe3, 0x30,
	0x10, 0x86, 0xd7, 0xf1, 0x26, 0x6b, 0xcf, 0x66, 0xc9, 0x46, 0x2c, 0x8b, 0xc8, 0xb2, 0x60, 0xdc,
	0x8b, 0x2f, 0x75, 0x68, 0xfa, 0x00, 0xa5, 0xbd, 0xe5, 0xea, 0x14, 0x0a, 0xbd, 0x04, 0x59, 0x9e,
	0x80, 0xd3, 0x56, 0x52, 0x2d, 0xc9, 0x34, 0x4f, 0xd8, 0xd7, 0x2a, 0x92, 0x92, 0x26, 0xb4, 0xd0,
	0xdb, 0xcc, 0x37, 0xbf, 0x7f, 0x33, 0xf3, 0x0b, 0x66, 0x8c, 0x73, 0x69, 0x85, 0xd1, 0xf3, 0xfe,
	0x62, 0x7e, 0xa8, 0x4b, 0xd5, 0x49, 0x23, 0xc9, 0x84, 0xa9, 0xb6, 0x7c, 0x67, 0xfd, 0x22, 0xbf,
	0x82, 0xe9, 0x52, 0x98, 0x4e, 0x6a, 0x85, 0xdc, 0x54, 0xf8, 0x6c, 0x51, 0x1b, 0xf2, 0x07, 0x86,
	0x46, 0x3e, 0xa0, 0xa0, 0x51, 0x16, 0x15, 0x69, 0x15, 0x1a, 0x47, 0x35, 0x97, 0x0a, 0xe9, 0x20,
	0x50, 0xdf, 0xe4, 0xaf, 0x03, 0x20, 0xa7, 0x0e, 0x5a, 0x49, 0xa1, 0x91, 0xfc, 0x85, 0x11, 0xe3,
	0xa6, 0xed, 0xd1, 0x7b, 0x24, 0xd5, 0xbe, 0x23, 0xbf, 0x21, 0x66, 0xb6, 0xa1, 0x83, 0x2c, 0x2e,
	0xd2, 0xca, 0x95, 0xe4, 0x1f, 0xa4, 0xfc, 0xb1, 0x45, 0x61, 0xd6, 0x6d, 0x43, 0x63, 0x6f, 0x9d,
	0x04, 0xb0, 0x6c, 0x9c, 0x1c, 0x5f, 0x14, 0xfd, 0x9e, 0x45, 0x45, 0x5c, 0xb9, 0x32, 0x10, 0x43,
	0x87, 0x59, 0x54, 0x8c, 0x1d, 0x31, 0x8e, 0xb4, 0xcc, 0xd0, 0x51, 0xd0, 0xb4, 0x2c, 0x10, 0xad,
	0xe9, 0x0f, 0x6f, 0xe6, 0x4a, 0x47, 0x44, 0xbd, 0xa1, 0x49, 0xd0, 0x88, 0x7a, 0x43, 0xce, 0x81,
	0xc8, 0x7a, 0x63, 0x35, 0x67, 0x06, 0x9b, 0xb5, 0xb6, 0xf5, 0x16, 0xb9, 0xa1, 0xa9, 0xff, 0x64,
	0x7a, 0x9c, 0xac, 0xc2, 0xe0, 0xb8, 0x3c, 0x9c, 0x2c, 0xef, 0x6c, 0xb5, 0xad, 0xe9, 0xcf, 0xf0,
	0x23, 0x6d, 0x6b, 0xf2, 0x1f, 0xc0, 0x5f, 0x6b, 0x6d, 0x76, 0x0a, 0xe9, 0xd8, 0x0f, 0x52, 0x4f,
	0x6e, 0x77, 0x0a, 0xc9, 0x0c, 0x12, 0xab, 0xb1, 0x13, 0xec, 0x09, 0xe9, 0xaf, 0xb0, 0xeb, 0xa1,
	0x5f, 0x6c, 0x61, 0x72, 0xbd, 0x4f, 0x66, 0x85, 0x5d, 0xdf, 0x72, 0x24, 0x77, 0x00, 0xc7, 0xdb,
	0x92, 0xbc, 0xfc, 0x90, 0x5e, 0xf9, 0x29, 0xba, 0xd9, 0xd9, 0x97, 0x9a, 0x10, 0x4e, 0xfe, 0xed,
	0x06, 0xee, 0x93, 0x83, 0xa6, 0x1e, 0xf9, 0xa7, 0x71, 0xf9, 0x16, 0x00, 0x00, 0xff, 0xff, 0x38,
	0x91, 0xdd, 0x3c, 0x38, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AccountsServiceClient is the client API for AccountsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AccountsServiceClient interface {
	Introspect(ctx context.Context, in *IntrospectRequest, opts ...grpc.CallOption) (*IntrospectResponse, error)
}

type accountsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountsServiceClient(cc grpc.ClientConnInterface) AccountsServiceClient {
	return &accountsServiceClient{cc}
}

func (c *accountsServiceClient) Introspect(ctx context.Context, in *IntrospectRequest, opts ...grpc.CallOption) (*IntrospectResponse, error) {
	out := new(IntrospectResponse)
	err := c.cc.Invoke(ctx, "/api.accounts.v2.AccountsService/Introspect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountsServiceServer is the server API for AccountsService service.
type AccountsServiceServer interface {
	Introspect(context.Context, *IntrospectRequest) (*IntrospectResponse, error)
}

// UnimplementedAccountsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAccountsServiceServer struct {
}

func (*UnimplementedAccountsServiceServer) Introspect(ctx context.Context, req *IntrospectRequest) (*IntrospectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Introspect not implemented")
}

func RegisterAccountsServiceServer(s *grpc.Server, srv AccountsServiceServer) {
	s.RegisterService(&_AccountsService_serviceDesc, srv)
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
		FullMethod: "/api.accounts.v2.AccountsService/Introspect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountsServiceServer).Introspect(ctx, req.(*IntrospectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AccountsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.accounts.v2.AccountsService",
	HandlerType: (*AccountsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Introspect",
			Handler:    _AccountsService_Introspect_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "accounts/v1/accounts.proto",
}
