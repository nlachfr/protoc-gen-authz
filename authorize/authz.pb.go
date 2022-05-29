// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: authorize/authz.proto

package authorize

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FileRule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Globals *FileRule_Globals      `protobuf:"bytes,1,opt,name=globals,proto3" json:"globals,omitempty"`
	Rules   map[string]*MethodRule `protobuf:"bytes,2,rep,name=rules,proto3" json:"rules,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *FileRule) Reset() {
	*x = FileRule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authorize_authz_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileRule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileRule) ProtoMessage() {}

func (x *FileRule) ProtoReflect() protoreflect.Message {
	mi := &file_authorize_authz_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileRule.ProtoReflect.Descriptor instead.
func (*FileRule) Descriptor() ([]byte, []int) {
	return file_authorize_authz_proto_rawDescGZIP(), []int{0}
}

func (x *FileRule) GetGlobals() *FileRule_Globals {
	if x != nil {
		return x.Globals
	}
	return nil
}

func (x *FileRule) GetRules() map[string]*MethodRule {
	if x != nil {
		return x.Rules
	}
	return nil
}

type MethodRule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Expr string `protobuf:"bytes,1,opt,name=expr,proto3" json:"expr,omitempty"`
}

func (x *MethodRule) Reset() {
	*x = MethodRule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authorize_authz_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MethodRule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MethodRule) ProtoMessage() {}

func (x *MethodRule) ProtoReflect() protoreflect.Message {
	mi := &file_authorize_authz_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MethodRule.ProtoReflect.Descriptor instead.
func (*MethodRule) Descriptor() ([]byte, []int) {
	return file_authorize_authz_proto_rawDescGZIP(), []int{1}
}

func (x *MethodRule) GetExpr() string {
	if x != nil {
		return x.Expr
	}
	return ""
}

type AuthorizationContext struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Peer     *AuthorizationContext_Peer                     `protobuf:"bytes,1,opt,name=peer,proto3" json:"peer,omitempty"`
	Metadata map[string]*AuthorizationContext_MetadataValue `protobuf:"bytes,2,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *AuthorizationContext) Reset() {
	*x = AuthorizationContext{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authorize_authz_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthorizationContext) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorizationContext) ProtoMessage() {}

func (x *AuthorizationContext) ProtoReflect() protoreflect.Message {
	mi := &file_authorize_authz_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorizationContext.ProtoReflect.Descriptor instead.
func (*AuthorizationContext) Descriptor() ([]byte, []int) {
	return file_authorize_authz_proto_rawDescGZIP(), []int{2}
}

func (x *AuthorizationContext) GetPeer() *AuthorizationContext_Peer {
	if x != nil {
		return x.Peer
	}
	return nil
}

func (x *AuthorizationContext) GetMetadata() map[string]*AuthorizationContext_MetadataValue {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type FileRule_Globals struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Functions map[string]string `protobuf:"bytes,1,rep,name=functions,proto3" json:"functions,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Constants map[string]string `protobuf:"bytes,2,rep,name=constants,proto3" json:"constants,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *FileRule_Globals) Reset() {
	*x = FileRule_Globals{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authorize_authz_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileRule_Globals) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileRule_Globals) ProtoMessage() {}

func (x *FileRule_Globals) ProtoReflect() protoreflect.Message {
	mi := &file_authorize_authz_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileRule_Globals.ProtoReflect.Descriptor instead.
func (*FileRule_Globals) Descriptor() ([]byte, []int) {
	return file_authorize_authz_proto_rawDescGZIP(), []int{0, 0}
}

func (x *FileRule_Globals) GetFunctions() map[string]string {
	if x != nil {
		return x.Functions
	}
	return nil
}

func (x *FileRule_Globals) GetConstants() map[string]string {
	if x != nil {
		return x.Constants
	}
	return nil
}

type AuthorizationContext_Peer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr     string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	AuthInfo string `protobuf:"bytes,2,opt,name=auth_info,json=authInfo,proto3" json:"auth_info,omitempty"`
}

func (x *AuthorizationContext_Peer) Reset() {
	*x = AuthorizationContext_Peer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authorize_authz_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthorizationContext_Peer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorizationContext_Peer) ProtoMessage() {}

func (x *AuthorizationContext_Peer) ProtoReflect() protoreflect.Message {
	mi := &file_authorize_authz_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorizationContext_Peer.ProtoReflect.Descriptor instead.
func (*AuthorizationContext_Peer) Descriptor() ([]byte, []int) {
	return file_authorize_authz_proto_rawDescGZIP(), []int{2, 0}
}

func (x *AuthorizationContext_Peer) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *AuthorizationContext_Peer) GetAuthInfo() string {
	if x != nil {
		return x.AuthInfo
	}
	return ""
}

type AuthorizationContext_MetadataValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []string `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *AuthorizationContext_MetadataValue) Reset() {
	*x = AuthorizationContext_MetadataValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authorize_authz_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthorizationContext_MetadataValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorizationContext_MetadataValue) ProtoMessage() {}

func (x *AuthorizationContext_MetadataValue) ProtoReflect() protoreflect.Message {
	mi := &file_authorize_authz_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorizationContext_MetadataValue.ProtoReflect.Descriptor instead.
func (*AuthorizationContext_MetadataValue) Descriptor() ([]byte, []int) {
	return file_authorize_authz_proto_rawDescGZIP(), []int{2, 1}
}

func (x *AuthorizationContext_MetadataValue) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

var file_authorize_authz_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*FileRule)(nil),
		Field:         1145,
		Name:          "authorize.file",
		Tag:           "bytes,1145,opt,name=file",
		Filename:      "authorize/authz.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*MethodRule)(nil),
		Field:         1145,
		Name:          "authorize.method",
		Tag:           "bytes,1145,opt,name=method",
		Filename:      "authorize/authz.proto",
	},
}

// Extension fields to descriptorpb.FileOptions.
var (
	// optional authorize.FileRule file = 1145;
	E_File = &file_authorize_authz_proto_extTypes[0]
)

// Extension fields to descriptorpb.MethodOptions.
var (
	// optional authorize.MethodRule method = 1145;
	E_Method = &file_authorize_authz_proto_extTypes[1]
)

var File_authorize_authz_proto protoreflect.FileDescriptor

var file_authorize_authz_proto_rawDesc = []byte{
	0x0a, 0x15, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x2f, 0x61, 0x75, 0x74, 0x68,
	0x7a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69,
	0x7a, 0x65, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe4, 0x03, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x75, 0x6c,
	0x65, 0x12, 0x35, 0x0a, 0x07, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x2e, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x73, 0x52,
	0x07, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x73, 0x12, 0x34, 0x0a, 0x05, 0x72, 0x75, 0x6c, 0x65,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x69, 0x7a, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x2e, 0x52, 0x75, 0x6c,
	0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x1a, 0x99,
	0x02, 0x0a, 0x07, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x73, 0x12, 0x48, 0x0a, 0x09, 0x66, 0x75,
	0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e,
	0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x75,
	0x6c, 0x65, 0x2e, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x73, 0x2e, 0x46, 0x75, 0x6e, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x66, 0x75, 0x6e, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x48, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x69, 0x7a, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x2e, 0x47, 0x6c, 0x6f,
	0x62, 0x61, 0x6c, 0x73, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x73, 0x1a, 0x3c,
	0x0a, 0x0e, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3c, 0x0a, 0x0e,
	0x43, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x4f, 0x0a, 0x0a, 0x52, 0x75,
	0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2b, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x69, 0x7a, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x75, 0x6c, 0x65,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x20, 0x0a, 0x0a, 0x4d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x65, 0x78, 0x70,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x65, 0x78, 0x70, 0x72, 0x22, 0xe9, 0x02,
	0x0a, 0x14, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12, 0x38, 0x0a, 0x04, 0x70, 0x65, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65,
	0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x78, 0x74, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x52, 0x04, 0x70, 0x65, 0x65, 0x72,
	0x12, 0x49, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x2e, 0x41,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x74,
	0x65, 0x78, 0x74, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x37, 0x0a, 0x04, 0x50,
	0x65, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x5f,
	0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x61, 0x75, 0x74, 0x68,
	0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x27, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x1a, 0x6a, 0x0a,
	0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x43, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2d, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x2e, 0x41, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74,
	0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x3a, 0x46, 0x0a, 0x04, 0x66, 0x69, 0x6c,
	0x65, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0xf9, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69,
	0x7a, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x04, 0x66, 0x69, 0x6c,
	0x65, 0x3a, 0x4e, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x1e, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xf9, 0x08, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x2e, 0x4d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x4e, 0x65, 0x61, 0x6b, 0x78, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65,
	0x6e, 0x2d, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_authorize_authz_proto_rawDescOnce sync.Once
	file_authorize_authz_proto_rawDescData = file_authorize_authz_proto_rawDesc
)

func file_authorize_authz_proto_rawDescGZIP() []byte {
	file_authorize_authz_proto_rawDescOnce.Do(func() {
		file_authorize_authz_proto_rawDescData = protoimpl.X.CompressGZIP(file_authorize_authz_proto_rawDescData)
	})
	return file_authorize_authz_proto_rawDescData
}

var file_authorize_authz_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_authorize_authz_proto_goTypes = []interface{}{
	(*FileRule)(nil),                  // 0: authorize.FileRule
	(*MethodRule)(nil),                // 1: authorize.MethodRule
	(*AuthorizationContext)(nil),      // 2: authorize.AuthorizationContext
	(*FileRule_Globals)(nil),          // 3: authorize.FileRule.Globals
	nil,                               // 4: authorize.FileRule.RulesEntry
	nil,                               // 5: authorize.FileRule.Globals.FunctionsEntry
	nil,                               // 6: authorize.FileRule.Globals.ConstantsEntry
	(*AuthorizationContext_Peer)(nil), // 7: authorize.AuthorizationContext.Peer
	(*AuthorizationContext_MetadataValue)(nil), // 8: authorize.AuthorizationContext.MetadataValue
	nil,                                // 9: authorize.AuthorizationContext.MetadataEntry
	(*descriptorpb.FileOptions)(nil),   // 10: google.protobuf.FileOptions
	(*descriptorpb.MethodOptions)(nil), // 11: google.protobuf.MethodOptions
}
var file_authorize_authz_proto_depIdxs = []int32{
	3,  // 0: authorize.FileRule.globals:type_name -> authorize.FileRule.Globals
	4,  // 1: authorize.FileRule.rules:type_name -> authorize.FileRule.RulesEntry
	7,  // 2: authorize.AuthorizationContext.peer:type_name -> authorize.AuthorizationContext.Peer
	9,  // 3: authorize.AuthorizationContext.metadata:type_name -> authorize.AuthorizationContext.MetadataEntry
	5,  // 4: authorize.FileRule.Globals.functions:type_name -> authorize.FileRule.Globals.FunctionsEntry
	6,  // 5: authorize.FileRule.Globals.constants:type_name -> authorize.FileRule.Globals.ConstantsEntry
	1,  // 6: authorize.FileRule.RulesEntry.value:type_name -> authorize.MethodRule
	8,  // 7: authorize.AuthorizationContext.MetadataEntry.value:type_name -> authorize.AuthorizationContext.MetadataValue
	10, // 8: authorize.file:extendee -> google.protobuf.FileOptions
	11, // 9: authorize.method:extendee -> google.protobuf.MethodOptions
	0,  // 10: authorize.file:type_name -> authorize.FileRule
	1,  // 11: authorize.method:type_name -> authorize.MethodRule
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	10, // [10:12] is the sub-list for extension type_name
	8,  // [8:10] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_authorize_authz_proto_init() }
func file_authorize_authz_proto_init() {
	if File_authorize_authz_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_authorize_authz_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileRule); i {
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
		file_authorize_authz_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MethodRule); i {
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
		file_authorize_authz_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthorizationContext); i {
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
		file_authorize_authz_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileRule_Globals); i {
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
		file_authorize_authz_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthorizationContext_Peer); i {
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
		file_authorize_authz_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthorizationContext_MetadataValue); i {
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
			RawDescriptor: file_authorize_authz_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_authorize_authz_proto_goTypes,
		DependencyIndexes: file_authorize_authz_proto_depIdxs,
		MessageInfos:      file_authorize_authz_proto_msgTypes,
		ExtensionInfos:    file_authorize_authz_proto_extTypes,
	}.Build()
	File_authorize_authz_proto = out.File
	file_authorize_authz_proto_rawDesc = nil
	file_authorize_authz_proto_goTypes = nil
	file_authorize_authz_proto_depIdxs = nil
}
