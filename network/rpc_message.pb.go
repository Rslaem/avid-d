// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.11.2
// source: rpc_message.proto

package network

import (
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

type MessageType int32

const (
	MessageType_NONE  MessageType = 0
	MessageType_VID   MessageType = 1
	MessageType_PB    MessageType = 2
	MessageType_SMVBA MessageType = 3
)

// Enum value maps for MessageType.
var (
	MessageType_name = map[int32]string{
		0: "NONE",
		1: "VID",
		2: "PB",
		3: "SMVBA",
	}
	MessageType_value = map[string]int32{
		"NONE":  0,
		"VID":   1,
		"PB":    2,
		"SMVBA": 3,
	}
)

func (x MessageType) Enum() *MessageType {
	p := new(MessageType)
	*p = x
	return p
}

func (x MessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc_message_proto_enumTypes[0].Descriptor()
}

func (MessageType) Type() protoreflect.EnumType {
	return &file_rpc_message_proto_enumTypes[0]
}

func (x MessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageType.Descriptor instead.
func (MessageType) EnumDescriptor() ([]byte, []int) {
	return file_rpc_message_proto_rawDescGZIP(), []int{0}
}

type RpcMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type MessageType `protobuf:"varint,5,opt,name=type,proto3,enum=main.MessageType" json:"type,omitempty"`
	From int32       `protobuf:"varint,1,opt,name=from,proto3" json:"from,omitempty"`
	Dest int32       `protobuf:"varint,2,opt,name=dest,proto3" json:"dest,omitempty"`
	Size int32       `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	Data []byte      `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *RpcMessage) Reset() {
	*x = RpcMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RpcMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RpcMessage) ProtoMessage() {}

func (x *RpcMessage) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RpcMessage.ProtoReflect.Descriptor instead.
func (*RpcMessage) Descriptor() ([]byte, []int) {
	return file_rpc_message_proto_rawDescGZIP(), []int{0}
}

func (x *RpcMessage) GetType() MessageType {
	if x != nil {
		return x.Type
	}
	return MessageType_NONE
}

func (x *RpcMessage) GetFrom() int32 {
	if x != nil {
		return x.From
	}
	return 0
}

func (x *RpcMessage) GetDest() int32 {
	if x != nil {
		return x.Dest
	}
	return 0
}

func (x *RpcMessage) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *RpcMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_rpc_message_proto protoreflect.FileDescriptor

var file_rpc_message_proto_rawDesc = []byte{
	0x0a, 0x11, 0x72, 0x70, 0x63, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x83, 0x01, 0x0a, 0x0a, 0x52, 0x70,
	0x63, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x25, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x66,
	0x72, 0x6f, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x64, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x2a,
	0x33, 0x0a, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08,
	0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x56, 0x49, 0x44, 0x10,
	0x01, 0x12, 0x06, 0x0a, 0x02, 0x50, 0x42, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x4d, 0x56,
	0x42, 0x41, 0x10, 0x03, 0x42, 0x03, 0x5a, 0x01, 0x2e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_rpc_message_proto_rawDescOnce sync.Once
	file_rpc_message_proto_rawDescData = file_rpc_message_proto_rawDesc
)

func file_rpc_message_proto_rawDescGZIP() []byte {
	file_rpc_message_proto_rawDescOnce.Do(func() {
		file_rpc_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_message_proto_rawDescData)
	})
	return file_rpc_message_proto_rawDescData
}

var file_rpc_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_rpc_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_rpc_message_proto_goTypes = []interface{}{
	(MessageType)(0),   // 0: main.messageType
	(*RpcMessage)(nil), // 1: main.RpcMessage
}
var file_rpc_message_proto_depIdxs = []int32{
	0, // 0: main.RpcMessage.type:type_name -> main.messageType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpc_message_proto_init() }
func file_rpc_message_proto_init() {
	if File_rpc_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RpcMessage); i {
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
			RawDescriptor: file_rpc_message_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_message_proto_goTypes,
		DependencyIndexes: file_rpc_message_proto_depIdxs,
		EnumInfos:         file_rpc_message_proto_enumTypes,
		MessageInfos:      file_rpc_message_proto_msgTypes,
	}.Build()
	File_rpc_message_proto = out.File
	file_rpc_message_proto_rawDesc = nil
	file_rpc_message_proto_goTypes = nil
	file_rpc_message_proto_depIdxs = nil
}