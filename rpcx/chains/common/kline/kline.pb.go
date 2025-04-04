// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        v3.20.3
// source: common/kline.proto

package kline

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

type IntervalData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Chain         string                 `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	C             float64                `protobuf:"fixed64,2,opt,name=c,proto3" json:"c,omitempty"`
	H             float64                `protobuf:"fixed64,3,opt,name=h,proto3" json:"h,omitempty"`
	L             float64                `protobuf:"fixed64,4,opt,name=l,proto3" json:"l,omitempty"`
	O             float64                `protobuf:"fixed64,5,opt,name=o,proto3" json:"o,omitempty"`
	Timestamp     int64                  `protobuf:"varint,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	V             float64                `protobuf:"fixed64,7,opt,name=v,proto3" json:"v,omitempty"`
	QuoteAddress  string                 `protobuf:"bytes,8,opt,name=quote_address,json=quoteAddress,proto3" json:"quote_address,omitempty"`    //token address
	MarketAddress string                 `protobuf:"bytes,9,opt,name=market_address,json=marketAddress,proto3" json:"market_address,omitempty"` //market address
	Interval      string                 `protobuf:"bytes,10,opt,name=interval,proto3" json:"interval,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *IntervalData) Reset() {
	*x = IntervalData{}
	mi := &file_common_kline_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IntervalData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntervalData) ProtoMessage() {}

func (x *IntervalData) ProtoReflect() protoreflect.Message {
	mi := &file_common_kline_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntervalData.ProtoReflect.Descriptor instead.
func (*IntervalData) Descriptor() ([]byte, []int) {
	return file_common_kline_proto_rawDescGZIP(), []int{0}
}

func (x *IntervalData) GetChain() string {
	if x != nil {
		return x.Chain
	}
	return ""
}

func (x *IntervalData) GetC() float64 {
	if x != nil {
		return x.C
	}
	return 0
}

func (x *IntervalData) GetH() float64 {
	if x != nil {
		return x.H
	}
	return 0
}

func (x *IntervalData) GetL() float64 {
	if x != nil {
		return x.L
	}
	return 0
}

func (x *IntervalData) GetO() float64 {
	if x != nil {
		return x.O
	}
	return 0
}

func (x *IntervalData) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *IntervalData) GetV() float64 {
	if x != nil {
		return x.V
	}
	return 0
}

func (x *IntervalData) GetQuoteAddress() string {
	if x != nil {
		return x.QuoteAddress
	}
	return ""
}

func (x *IntervalData) GetMarketAddress() string {
	if x != nil {
		return x.MarketAddress
	}
	return ""
}

func (x *IntervalData) GetInterval() string {
	if x != nil {
		return x.Interval
	}
	return ""
}

var File_common_kline_proto protoreflect.FileDescriptor

var file_common_kline_proto_rawDesc = []byte{
	0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6b, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6b, 0x6c, 0x69, 0x6e, 0x65, 0x22, 0xf0, 0x01, 0x0a, 0x0c,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x44, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05,
	0x63, 0x68, 0x61, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x12, 0x0c, 0x0a, 0x01, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x63,
	0x12, 0x0c, 0x0a, 0x01, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x68, 0x12, 0x0c,
	0x0a, 0x01, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x6c, 0x12, 0x0c, 0x0a, 0x01,
	0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01, 0x6f, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x0c, 0x0a, 0x01, 0x76, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x01, 0x76, 0x12, 0x23, 0x0a, 0x0d, 0x71, 0x75, 0x6f, 0x74, 0x65, 0x5f,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x71,
	0x75, 0x6f, 0x74, 0x65, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x6d,
	0x61, 0x72, 0x6b, 0x65, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x42, 0x35,
	0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x69, 0x6d,
	0x61, 0x6e, 0x63, 0x65, 0x2d, 0x61, 0x69, 0x2f, 0x73, 0x6d, 0x64, 0x78, 0x2f, 0x72, 0x70, 0x63,
	0x78, 0x2f, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f,
	0x6b, 0x6c, 0x69, 0x6e, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_kline_proto_rawDescOnce sync.Once
	file_common_kline_proto_rawDescData = file_common_kline_proto_rawDesc
)

func file_common_kline_proto_rawDescGZIP() []byte {
	file_common_kline_proto_rawDescOnce.Do(func() {
		file_common_kline_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_kline_proto_rawDescData)
	})
	return file_common_kline_proto_rawDescData
}

var file_common_kline_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_kline_proto_goTypes = []any{
	(*IntervalData)(nil), // 0: kline.IntervalData
}
var file_common_kline_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_kline_proto_init() }
func file_common_kline_proto_init() {
	if File_common_kline_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_kline_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_kline_proto_goTypes,
		DependencyIndexes: file_common_kline_proto_depIdxs,
		MessageInfos:      file_common_kline_proto_msgTypes,
	}.Build()
	File_common_kline_proto = out.File
	file_common_kline_proto_rawDesc = nil
	file_common_kline_proto_goTypes = nil
	file_common_kline_proto_depIdxs = nil
}
