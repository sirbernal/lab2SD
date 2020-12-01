// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.1
// source: proto/NodeService.proto

package proto

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type AliveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *AliveRequest) Reset() {
	*x = AliveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_NodeService_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AliveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AliveRequest) ProtoMessage() {}

func (x *AliveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_NodeService_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AliveRequest.ProtoReflect.Descriptor instead.
func (*AliveRequest) Descriptor() ([]byte, []int) {
	return file_proto_NodeService_proto_rawDescGZIP(), []int{0}
}

func (x *AliveRequest) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type AliveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *AliveResponse) Reset() {
	*x = AliveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_NodeService_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AliveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AliveResponse) ProtoMessage() {}

func (x *AliveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_NodeService_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AliveResponse.ProtoReflect.Descriptor instead.
func (*AliveResponse) Descriptor() ([]byte, []int) {
	return file_proto_NodeService_proto_rawDescGZIP(), []int{1}
}

func (x *AliveResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type PropuestaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Prop []int64 `protobuf:"varint,4,rep,packed,name=prop,proto3" json:"prop,omitempty"`
	Name string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *PropuestaRequest) Reset() {
	*x = PropuestaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_NodeService_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PropuestaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PropuestaRequest) ProtoMessage() {}

func (x *PropuestaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_NodeService_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PropuestaRequest.ProtoReflect.Descriptor instead.
func (*PropuestaRequest) Descriptor() ([]byte, []int) {
	return file_proto_NodeService_proto_rawDescGZIP(), []int{2}
}

func (x *PropuestaRequest) GetProp() []int64 {
	if x != nil {
		return x.Prop
	}
	return nil
}

func (x *PropuestaRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type PropuestaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg  bool    `protobuf:"varint,1,opt,name=msg,proto3" json:"msg,omitempty"`
	Prop []int64 `protobuf:"varint,4,rep,packed,name=prop,proto3" json:"prop,omitempty"`
}

func (x *PropuestaResponse) Reset() {
	*x = PropuestaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_NodeService_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PropuestaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PropuestaResponse) ProtoMessage() {}

func (x *PropuestaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_NodeService_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PropuestaResponse.ProtoReflect.Descriptor instead.
func (*PropuestaResponse) Descriptor() ([]byte, []int) {
	return file_proto_NodeService_proto_rawDescGZIP(), []int{3}
}

func (x *PropuestaResponse) GetMsg() bool {
	if x != nil {
		return x.Msg
	}
	return false
}

func (x *PropuestaResponse) GetProp() []int64 {
	if x != nil {
		return x.Prop
	}
	return nil
}

type DistribucionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chunk []byte `protobuf:"bytes,1,opt,name=chunk,proto3" json:"chunk,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *DistribucionRequest) Reset() {
	*x = DistribucionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_NodeService_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DistribucionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DistribucionRequest) ProtoMessage() {}

func (x *DistribucionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_NodeService_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DistribucionRequest.ProtoReflect.Descriptor instead.
func (*DistribucionRequest) Descriptor() ([]byte, []int) {
	return file_proto_NodeService_proto_rawDescGZIP(), []int{4}
}

func (x *DistribucionRequest) GetChunk() []byte {
	if x != nil {
		return x.Chunk
	}
	return nil
}

func (x *DistribucionRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type DistribucionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resp string `protobuf:"bytes,1,opt,name=resp,proto3" json:"resp,omitempty"`
}

func (x *DistribucionResponse) Reset() {
	*x = DistribucionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_NodeService_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DistribucionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DistribucionResponse) ProtoMessage() {}

func (x *DistribucionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_NodeService_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DistribucionResponse.ProtoReflect.Descriptor instead.
func (*DistribucionResponse) Descriptor() ([]byte, []int) {
	return file_proto_NodeService_proto_rawDescGZIP(), []int{5}
}

func (x *DistribucionResponse) GetResp() string {
	if x != nil {
		return x.Resp
	}
	return ""
}

type RicandAgraRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RicandAgraRequest) Reset() {
	*x = RicandAgraRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_NodeService_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RicandAgraRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RicandAgraRequest) ProtoMessage() {}

func (x *RicandAgraRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_NodeService_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RicandAgraRequest.ProtoReflect.Descriptor instead.
func (*RicandAgraRequest) Descriptor() ([]byte, []int) {
	return file_proto_NodeService_proto_rawDescGZIP(), []int{6}
}

func (x *RicandAgraRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type RicandAgraResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resp string `protobuf:"bytes,1,opt,name=resp,proto3" json:"resp,omitempty"`
	Id   int64  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RicandAgraResponse) Reset() {
	*x = RicandAgraResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_NodeService_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RicandAgraResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RicandAgraResponse) ProtoMessage() {}

func (x *RicandAgraResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_NodeService_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RicandAgraResponse.ProtoReflect.Descriptor instead.
func (*RicandAgraResponse) Descriptor() ([]byte, []int) {
	return file_proto_NodeService_proto_rawDescGZIP(), []int{7}
}

func (x *RicandAgraResponse) GetResp() string {
	if x != nil {
		return x.Resp
	}
	return ""
}

func (x *RicandAgraResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

var File_proto_NodeService_proto protoreflect.FileDescriptor

var file_proto_NodeService_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x4e, 0x6f, 0x64, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6e, 0x6f, 0x64, 0x65, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x20, 0x0a, 0x0c, 0x41, 0x6c, 0x69, 0x76, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x21, 0x0a, 0x0d, 0x41, 0x6c, 0x69,
	0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x3a, 0x0a, 0x10,
	0x50, 0x72, 0x6f, 0x70, 0x75, 0x65, 0x73, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x72, 0x6f, 0x70, 0x18, 0x04, 0x20, 0x03, 0x28, 0x03, 0x52, 0x04,
	0x70, 0x72, 0x6f, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x39, 0x0a, 0x11, 0x50, 0x72, 0x6f, 0x70,
	0x75, 0x65, 0x73, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12,
	0x12, 0x0a, 0x04, 0x70, 0x72, 0x6f, 0x70, 0x18, 0x04, 0x20, 0x03, 0x28, 0x03, 0x52, 0x04, 0x70,
	0x72, 0x6f, 0x70, 0x22, 0x3f, 0x0a, 0x13, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x63,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x68,
	0x75, 0x6e, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x22, 0x2a, 0x0a, 0x14, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75,
	0x63, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x72, 0x65, 0x73, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x65, 0x73, 0x70,
	0x22, 0x23, 0x0a, 0x11, 0x52, 0x69, 0x63, 0x61, 0x6e, 0x64, 0x41, 0x67, 0x72, 0x61, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x38, 0x0a, 0x12, 0x52, 0x69, 0x63, 0x61, 0x6e, 0x64, 0x41,
	0x67, 0x72, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x72,
	0x65, 0x73, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x65, 0x73, 0x70, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x32,
	0xcd, 0x02, 0x0a, 0x0b, 0x4e, 0x6f, 0x64, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x42, 0x0a, 0x05, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x12, 0x1a, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x4e, 0x0a, 0x09, 0x50, 0x72, 0x6f, 0x70, 0x75, 0x65, 0x73, 0x74, 0x61,
	0x12, 0x1e, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x50, 0x72, 0x6f, 0x70, 0x75, 0x65, 0x73, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1f, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x50, 0x72, 0x6f, 0x70, 0x75, 0x65, 0x73, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x57, 0x0a, 0x0c, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x63,
	0x69, 0x6f, 0x6e, 0x12, 0x21, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x63, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x63, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x0a,
	0x52, 0x69, 0x63, 0x61, 0x6e, 0x64, 0x41, 0x67, 0x72, 0x61, 0x12, 0x1f, 0x2e, 0x6e, 0x6f, 0x64,
	0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x69, 0x63, 0x61, 0x6e, 0x64,
	0x41, 0x67, 0x72, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x6e, 0x6f,
	0x64, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x69, 0x63, 0x61, 0x6e,
	0x64, 0x41, 0x67, 0x72, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x14, 0x5a, 0x12, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x3b,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_NodeService_proto_rawDescOnce sync.Once
	file_proto_NodeService_proto_rawDescData = file_proto_NodeService_proto_rawDesc
)

func file_proto_NodeService_proto_rawDescGZIP() []byte {
	file_proto_NodeService_proto_rawDescOnce.Do(func() {
		file_proto_NodeService_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_NodeService_proto_rawDescData)
	})
	return file_proto_NodeService_proto_rawDescData
}

var file_proto_NodeService_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_NodeService_proto_goTypes = []interface{}{
	(*AliveRequest)(nil),         // 0: node_service.AliveRequest
	(*AliveResponse)(nil),        // 1: node_service.AliveResponse
	(*PropuestaRequest)(nil),     // 2: node_service.PropuestaRequest
	(*PropuestaResponse)(nil),    // 3: node_service.PropuestaResponse
	(*DistribucionRequest)(nil),  // 4: node_service.DistribucionRequest
	(*DistribucionResponse)(nil), // 5: node_service.DistribucionResponse
	(*RicandAgraRequest)(nil),    // 6: node_service.RicandAgraRequest
	(*RicandAgraResponse)(nil),   // 7: node_service.RicandAgraResponse
}
var file_proto_NodeService_proto_depIdxs = []int32{
	0, // 0: node_service.NodeService.Alive:input_type -> node_service.AliveRequest
	2, // 1: node_service.NodeService.Propuesta:input_type -> node_service.PropuestaRequest
	4, // 2: node_service.NodeService.Distribucion:input_type -> node_service.DistribucionRequest
	6, // 3: node_service.NodeService.RicandAgra:input_type -> node_service.RicandAgraRequest
	1, // 4: node_service.NodeService.Alive:output_type -> node_service.AliveResponse
	3, // 5: node_service.NodeService.Propuesta:output_type -> node_service.PropuestaResponse
	5, // 6: node_service.NodeService.Distribucion:output_type -> node_service.DistribucionResponse
	7, // 7: node_service.NodeService.RicandAgra:output_type -> node_service.RicandAgraResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_NodeService_proto_init() }
func file_proto_NodeService_proto_init() {
	if File_proto_NodeService_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_NodeService_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AliveRequest); i {
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
		file_proto_NodeService_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AliveResponse); i {
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
		file_proto_NodeService_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PropuestaRequest); i {
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
		file_proto_NodeService_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PropuestaResponse); i {
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
		file_proto_NodeService_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DistribucionRequest); i {
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
		file_proto_NodeService_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DistribucionResponse); i {
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
		file_proto_NodeService_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RicandAgraRequest); i {
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
		file_proto_NodeService_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RicandAgraResponse); i {
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
			RawDescriptor: file_proto_NodeService_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_NodeService_proto_goTypes,
		DependencyIndexes: file_proto_NodeService_proto_depIdxs,
		MessageInfos:      file_proto_NodeService_proto_msgTypes,
	}.Build()
	File_proto_NodeService_proto = out.File
	file_proto_NodeService_proto_rawDesc = nil
	file_proto_NodeService_proto_goTypes = nil
	file_proto_NodeService_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// NodeServiceClient is the client API for NodeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NodeServiceClient interface {
	Alive(ctx context.Context, in *AliveRequest, opts ...grpc.CallOption) (*AliveResponse, error)
	Propuesta(ctx context.Context, in *PropuestaRequest, opts ...grpc.CallOption) (*PropuestaResponse, error)
	Distribucion(ctx context.Context, in *DistribucionRequest, opts ...grpc.CallOption) (*DistribucionResponse, error)
	RicandAgra(ctx context.Context, in *RicandAgraRequest, opts ...grpc.CallOption) (*RicandAgraResponse, error)
}

type nodeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeServiceClient(cc grpc.ClientConnInterface) NodeServiceClient {
	return &nodeServiceClient{cc}
}

func (c *nodeServiceClient) Alive(ctx context.Context, in *AliveRequest, opts ...grpc.CallOption) (*AliveResponse, error) {
	out := new(AliveResponse)
	err := c.cc.Invoke(ctx, "/node_service.NodeService/Alive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServiceClient) Propuesta(ctx context.Context, in *PropuestaRequest, opts ...grpc.CallOption) (*PropuestaResponse, error) {
	out := new(PropuestaResponse)
	err := c.cc.Invoke(ctx, "/node_service.NodeService/Propuesta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServiceClient) Distribucion(ctx context.Context, in *DistribucionRequest, opts ...grpc.CallOption) (*DistribucionResponse, error) {
	out := new(DistribucionResponse)
	err := c.cc.Invoke(ctx, "/node_service.NodeService/Distribucion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeServiceClient) RicandAgra(ctx context.Context, in *RicandAgraRequest, opts ...grpc.CallOption) (*RicandAgraResponse, error) {
	out := new(RicandAgraResponse)
	err := c.cc.Invoke(ctx, "/node_service.NodeService/RicandAgra", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeServiceServer is the server API for NodeService service.
type NodeServiceServer interface {
	Alive(context.Context, *AliveRequest) (*AliveResponse, error)
	Propuesta(context.Context, *PropuestaRequest) (*PropuestaResponse, error)
	Distribucion(context.Context, *DistribucionRequest) (*DistribucionResponse, error)
	RicandAgra(context.Context, *RicandAgraRequest) (*RicandAgraResponse, error)
}

// UnimplementedNodeServiceServer can be embedded to have forward compatible implementations.
type UnimplementedNodeServiceServer struct {
}

func (*UnimplementedNodeServiceServer) Alive(context.Context, *AliveRequest) (*AliveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Alive not implemented")
}
func (*UnimplementedNodeServiceServer) Propuesta(context.Context, *PropuestaRequest) (*PropuestaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Propuesta not implemented")
}
func (*UnimplementedNodeServiceServer) Distribucion(context.Context, *DistribucionRequest) (*DistribucionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Distribucion not implemented")
}
func (*UnimplementedNodeServiceServer) RicandAgra(context.Context, *RicandAgraRequest) (*RicandAgraResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RicandAgra not implemented")
}

func RegisterNodeServiceServer(s *grpc.Server, srv NodeServiceServer) {
	s.RegisterService(&_NodeService_serviceDesc, srv)
}

func _NodeService_Alive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AliveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).Alive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/node_service.NodeService/Alive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).Alive(ctx, req.(*AliveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeService_Propuesta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PropuestaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).Propuesta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/node_service.NodeService/Propuesta",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).Propuesta(ctx, req.(*PropuestaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeService_Distribucion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DistribucionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).Distribucion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/node_service.NodeService/Distribucion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).Distribucion(ctx, req.(*DistribucionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeService_RicandAgra_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RicandAgraRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).RicandAgra(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/node_service.NodeService/RicandAgra",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).RicandAgra(ctx, req.(*RicandAgraRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _NodeService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "node_service.NodeService",
	HandlerType: (*NodeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Alive",
			Handler:    _NodeService_Alive_Handler,
		},
		{
			MethodName: "Propuesta",
			Handler:    _NodeService_Propuesta_Handler,
		},
		{
			MethodName: "Distribucion",
			Handler:    _NodeService_Distribucion_Handler,
		},
		{
			MethodName: "RicandAgra",
			Handler:    _NodeService_RicandAgra_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/NodeService.proto",
}
