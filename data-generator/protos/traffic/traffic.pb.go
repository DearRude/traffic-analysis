// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.24.4
// source: traffic.proto

package traffic

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

type Point struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lat float64 `protobuf:"fixed64,1,opt,name=lat,proto3" json:"lat,omitempty"`
	Lon float64 `protobuf:"fixed64,2,opt,name=lon,proto3" json:"lon,omitempty"`
}

func (x *Point) Reset() {
	*x = Point{}
	if protoimpl.UnsafeEnabled {
		mi := &file_traffic_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Point) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Point) ProtoMessage() {}

func (x *Point) ProtoReflect() protoreflect.Message {
	mi := &file_traffic_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Point.ProtoReflect.Descriptor instead.
func (*Point) Descriptor() ([]byte, []int) {
	return file_traffic_proto_rawDescGZIP(), []int{0}
}

func (x *Point) GetLat() float64 {
	if x != nil {
		return x.Lat
	}
	return 0
}

func (x *Point) GetLon() float64 {
	if x != nil {
		return x.Lon
	}
	return 0
}

type LineTraffic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         uint32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Length     float64  `protobuf:"fixed64,2,opt,name=length,proto3" json:"length,omitempty"`
	Timestamp  int64    `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	City       string   `protobuf:"bytes,4,opt,name=city,proto3" json:"city,omitempty"`
	RoadClass  string   `protobuf:"bytes,5,opt,name=road_class,json=roadClass,proto3" json:"road_class,omitempty"`
	Congestion string   `protobuf:"bytes,6,opt,name=congestion,proto3" json:"congestion,omitempty"`
	Geometry   []*Point `protobuf:"bytes,7,rep,name=geometry,proto3" json:"geometry,omitempty"`
}

func (x *LineTraffic) Reset() {
	*x = LineTraffic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_traffic_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LineTraffic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LineTraffic) ProtoMessage() {}

func (x *LineTraffic) ProtoReflect() protoreflect.Message {
	mi := &file_traffic_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LineTraffic.ProtoReflect.Descriptor instead.
func (*LineTraffic) Descriptor() ([]byte, []int) {
	return file_traffic_proto_rawDescGZIP(), []int{1}
}

func (x *LineTraffic) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LineTraffic) GetLength() float64 {
	if x != nil {
		return x.Length
	}
	return 0
}

func (x *LineTraffic) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *LineTraffic) GetCity() string {
	if x != nil {
		return x.City
	}
	return ""
}

func (x *LineTraffic) GetRoadClass() string {
	if x != nil {
		return x.RoadClass
	}
	return ""
}

func (x *LineTraffic) GetCongestion() string {
	if x != nil {
		return x.Congestion
	}
	return ""
}

func (x *LineTraffic) GetGeometry() []*Point {
	if x != nil {
		return x.Geometry
	}
	return nil
}

type LineTraffics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Traffics []*LineTraffic `protobuf:"bytes,1,rep,name=traffics,proto3" json:"traffics,omitempty"`
}

func (x *LineTraffics) Reset() {
	*x = LineTraffics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_traffic_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LineTraffics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LineTraffics) ProtoMessage() {}

func (x *LineTraffics) ProtoReflect() protoreflect.Message {
	mi := &file_traffic_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LineTraffics.ProtoReflect.Descriptor instead.
func (*LineTraffics) Descriptor() ([]byte, []int) {
	return file_traffic_proto_rawDescGZIP(), []int{2}
}

func (x *LineTraffics) GetTraffics() []*LineTraffic {
	if x != nil {
		return x.Traffics
	}
	return nil
}

var File_traffic_proto protoreflect.FileDescriptor

var file_traffic_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x2b, 0x0a, 0x05, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x61, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6c, 0x61, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6c, 0x6f, 0x6e, 0x22, 0xca, 0x01, 0x0a,
	0x0b, 0x4c, 0x69, 0x6e, 0x65, 0x54, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06,
	0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x6c, 0x65,
	0x6e, 0x67, 0x74, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x63, 0x69, 0x74, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x6f, 0x61, 0x64, 0x5f, 0x63,
	0x6c, 0x61, 0x73, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x6f, 0x61, 0x64,
	0x43, 0x6c, 0x61, 0x73, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x67, 0x65, 0x73, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x67, 0x65,
	0x73, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x08, 0x67, 0x65, 0x6f, 0x6d, 0x65, 0x74, 0x72,
	0x79, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52,
	0x08, 0x67, 0x65, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x22, 0x38, 0x0a, 0x0c, 0x4c, 0x69, 0x6e,
	0x65, 0x54, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x73, 0x12, 0x28, 0x0a, 0x08, 0x74, 0x72, 0x61,
	0x66, 0x66, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x4c, 0x69,
	0x6e, 0x65, 0x54, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x52, 0x08, 0x74, 0x72, 0x61, 0x66, 0x66,
	0x69, 0x63, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_traffic_proto_rawDescOnce sync.Once
	file_traffic_proto_rawDescData = file_traffic_proto_rawDesc
)

func file_traffic_proto_rawDescGZIP() []byte {
	file_traffic_proto_rawDescOnce.Do(func() {
		file_traffic_proto_rawDescData = protoimpl.X.CompressGZIP(file_traffic_proto_rawDescData)
	})
	return file_traffic_proto_rawDescData
}

var file_traffic_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_traffic_proto_goTypes = []interface{}{
	(*Point)(nil),        // 0: Point
	(*LineTraffic)(nil),  // 1: LineTraffic
	(*LineTraffics)(nil), // 2: LineTraffics
}
var file_traffic_proto_depIdxs = []int32{
	0, // 0: LineTraffic.geometry:type_name -> Point
	1, // 1: LineTraffics.traffics:type_name -> LineTraffic
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_traffic_proto_init() }
func file_traffic_proto_init() {
	if File_traffic_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_traffic_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Point); i {
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
		file_traffic_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LineTraffic); i {
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
		file_traffic_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LineTraffics); i {
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
			RawDescriptor: file_traffic_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_traffic_proto_goTypes,
		DependencyIndexes: file_traffic_proto_depIdxs,
		MessageInfos:      file_traffic_proto_msgTypes,
	}.Build()
	File_traffic_proto = out.File
	file_traffic_proto_rawDesc = nil
	file_traffic_proto_goTypes = nil
	file_traffic_proto_depIdxs = nil
}
