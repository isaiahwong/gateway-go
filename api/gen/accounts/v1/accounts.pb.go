// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.0
// source: accounts/v1/accounts.proto

package accounts

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type IntrospectRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Scope string `protobuf:"bytes,2,opt,name=scope,proto3" json:"scope,omitempty"`
}

func (x *IntrospectRequest) Reset() {
	*x = IntrospectRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounts_v1_accounts_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IntrospectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntrospectRequest) ProtoMessage() {}

func (x *IntrospectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_accounts_v1_accounts_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntrospectRequest.ProtoReflect.Descriptor instead.
func (*IntrospectRequest) Descriptor() ([]byte, []int) {
	return file_accounts_v1_accounts_proto_rawDescGZIP(), []int{0}
}

func (x *IntrospectRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *IntrospectRequest) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

type IntrospectResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Active            bool     `protobuf:"varint,1,opt,name=active,proto3" json:"active,omitempty"`
	Aud               []string `protobuf:"bytes,2,rep,name=aud,proto3" json:"aud,omitempty"`
	ClientId          string   `protobuf:"bytes,3,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Exp               int64    `protobuf:"varint,4,opt,name=exp,proto3" json:"exp,omitempty"`
	Ext               []byte   `protobuf:"bytes,5,opt,name=ext,proto3" json:"ext,omitempty"`
	Iat               int64    `protobuf:"varint,6,opt,name=iat,proto3" json:"iat,omitempty"`
	Iss               string   `protobuf:"bytes,7,opt,name=iss,proto3" json:"iss,omitempty"`
	Nbf               int64    `protobuf:"varint,8,opt,name=nbf,proto3" json:"nbf,omitempty"`
	ObfuscatedSubject string   `protobuf:"bytes,9,opt,name=obfuscated_subject,json=obfuscatedSubject,proto3" json:"obfuscated_subject,omitempty"`
	Scope             string   `protobuf:"bytes,10,opt,name=scope,proto3" json:"scope,omitempty"`
	Sub               string   `protobuf:"bytes,11,opt,name=sub,proto3" json:"sub,omitempty"`
	TokenType         string   `protobuf:"bytes,12,opt,name=token_type,json=tokenType,proto3" json:"token_type,omitempty"`
	Username          string   `protobuf:"bytes,13,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *IntrospectResponse) Reset() {
	*x = IntrospectResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_accounts_v1_accounts_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IntrospectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntrospectResponse) ProtoMessage() {}

func (x *IntrospectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_accounts_v1_accounts_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntrospectResponse.ProtoReflect.Descriptor instead.
func (*IntrospectResponse) Descriptor() ([]byte, []int) {
	return file_accounts_v1_accounts_proto_rawDescGZIP(), []int{1}
}

func (x *IntrospectResponse) GetActive() bool {
	if x != nil {
		return x.Active
	}
	return false
}

func (x *IntrospectResponse) GetAud() []string {
	if x != nil {
		return x.Aud
	}
	return nil
}

func (x *IntrospectResponse) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *IntrospectResponse) GetExp() int64 {
	if x != nil {
		return x.Exp
	}
	return 0
}

func (x *IntrospectResponse) GetExt() []byte {
	if x != nil {
		return x.Ext
	}
	return nil
}

func (x *IntrospectResponse) GetIat() int64 {
	if x != nil {
		return x.Iat
	}
	return 0
}

func (x *IntrospectResponse) GetIss() string {
	if x != nil {
		return x.Iss
	}
	return ""
}

func (x *IntrospectResponse) GetNbf() int64 {
	if x != nil {
		return x.Nbf
	}
	return 0
}

func (x *IntrospectResponse) GetObfuscatedSubject() string {
	if x != nil {
		return x.ObfuscatedSubject
	}
	return ""
}

func (x *IntrospectResponse) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *IntrospectResponse) GetSub() string {
	if x != nil {
		return x.Sub
	}
	return ""
}

func (x *IntrospectResponse) GetTokenType() string {
	if x != nil {
		return x.TokenType
	}
	return ""
}

func (x *IntrospectResponse) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

var File_accounts_v1_accounts_proto protoreflect.FileDescriptor

var file_accounts_v1_accounts_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x61, 0x70,
	0x69, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3f, 0x0a, 0x11, 0x49,
	0x6e, 0x74, 0x72, 0x6f, 0x73, 0x70, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x22, 0xc7, 0x02, 0x0a,
	0x12, 0x49, 0x6e, 0x74, 0x72, 0x6f, 0x73, 0x70, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x61,
	0x75, 0x64, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x61, 0x75, 0x64, 0x12, 0x1b, 0x0a,
	0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x78,
	0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x65, 0x78, 0x70, 0x12, 0x10, 0x0a, 0x03,
	0x65, 0x78, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x65, 0x78, 0x74, 0x12, 0x10,
	0x0a, 0x03, 0x69, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x69, 0x61, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x69, 0x73, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x69,
	0x73, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x6e, 0x62, 0x66, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x03, 0x6e, 0x62, 0x66, 0x12, 0x2d, 0x0a, 0x12, 0x6f, 0x62, 0x66, 0x75, 0x73, 0x63, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x11, 0x6f, 0x62, 0x66, 0x75, 0x73, 0x63, 0x61, 0x74, 0x65, 0x64, 0x53, 0x75, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x75, 0x62,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x75, 0x62, 0x12, 0x1d, 0x0a, 0x0a, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x32, 0x88, 0x01, 0x0a, 0x0f, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x75, 0x0a, 0x0a, 0x49, 0x6e,
	0x74, 0x72, 0x6f, 0x73, 0x70, 0x65, 0x63, 0x74, 0x12, 0x22, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x74, 0x72, 0x6f,
	0x73, 0x70, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49,
	0x6e, 0x74, 0x72, 0x6f, 0x73, 0x70, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x22, 0x13, 0x2f, 0x76, 0x31, 0x2f, 0x61,
	0x75, 0x74, 0x68, 0x2f, 0x69, 0x6e, 0x74, 0x72, 0x6f, 0x73, 0x70, 0x65, 0x63, 0x74, 0x3a, 0x01,
	0x2a, 0x42, 0x16, 0x5a, 0x14, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31,
	0x3b, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_accounts_v1_accounts_proto_rawDescOnce sync.Once
	file_accounts_v1_accounts_proto_rawDescData = file_accounts_v1_accounts_proto_rawDesc
)

func file_accounts_v1_accounts_proto_rawDescGZIP() []byte {
	file_accounts_v1_accounts_proto_rawDescOnce.Do(func() {
		file_accounts_v1_accounts_proto_rawDescData = protoimpl.X.CompressGZIP(file_accounts_v1_accounts_proto_rawDescData)
	})
	return file_accounts_v1_accounts_proto_rawDescData
}

var file_accounts_v1_accounts_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_accounts_v1_accounts_proto_goTypes = []interface{}{
	(*IntrospectRequest)(nil),  // 0: api.accounts.v1.IntrospectRequest
	(*IntrospectResponse)(nil), // 1: api.accounts.v1.IntrospectResponse
}
var file_accounts_v1_accounts_proto_depIdxs = []int32{
	0, // 0: api.accounts.v1.AccountsService.Introspect:input_type -> api.accounts.v1.IntrospectRequest
	1, // 1: api.accounts.v1.AccountsService.Introspect:output_type -> api.accounts.v1.IntrospectResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_accounts_v1_accounts_proto_init() }
func file_accounts_v1_accounts_proto_init() {
	if File_accounts_v1_accounts_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_accounts_v1_accounts_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IntrospectRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_accounts_v1_accounts_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IntrospectResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_accounts_v1_accounts_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_accounts_v1_accounts_proto_goTypes,
		DependencyIndexes: file_accounts_v1_accounts_proto_depIdxs,
		MessageInfos:      file_accounts_v1_accounts_proto_msgTypes,
	}.Build()
	File_accounts_v1_accounts_proto = out.File
	file_accounts_v1_accounts_proto_rawDesc = nil
	file_accounts_v1_accounts_proto_goTypes = nil
	file_accounts_v1_accounts_proto_depIdxs = nil
}