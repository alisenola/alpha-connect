// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: asset_data.proto

package models

import (
	models "gitlab.com/alphaticks/xchanger/models"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ProtocolAsset struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProtocolAssetID uint64                  `protobuf:"varint,1,opt,name=protocol_assetID,json=protocolAssetID,proto3" json:"protocol_assetID,omitempty"`
	Protocol        *models.Protocol        `protobuf:"bytes,2,opt,name=protocol,proto3" json:"protocol,omitempty"`
	Asset           *models.Asset           `protobuf:"bytes,3,opt,name=asset,proto3" json:"asset,omitempty"`
	Chain           *models.Chain           `protobuf:"bytes,4,opt,name=chain,proto3" json:"chain,omitempty"`
	CreationBlock   *wrapperspb.UInt64Value `protobuf:"bytes,5,opt,name=creation_block,json=creationBlock,proto3" json:"creation_block,omitempty"`
	CreationDate    *timestamppb.Timestamp  `protobuf:"bytes,6,opt,name=creation_date,json=creationDate,proto3" json:"creation_date,omitempty"`
	ContractAddress *wrapperspb.StringValue `protobuf:"bytes,7,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
	Decimals        *wrapperspb.UInt32Value `protobuf:"bytes,8,opt,name=decimals,proto3" json:"decimals,omitempty"`
}

func (x *ProtocolAsset) Reset() {
	*x = ProtocolAsset{}
	if protoimpl.UnsafeEnabled {
		mi := &file_asset_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtocolAsset) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtocolAsset) ProtoMessage() {}

func (x *ProtocolAsset) ProtoReflect() protoreflect.Message {
	mi := &file_asset_data_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtocolAsset.ProtoReflect.Descriptor instead.
func (*ProtocolAsset) Descriptor() ([]byte, []int) {
	return file_asset_data_proto_rawDescGZIP(), []int{0}
}

func (x *ProtocolAsset) GetProtocolAssetID() uint64 {
	if x != nil {
		return x.ProtocolAssetID
	}
	return 0
}

func (x *ProtocolAsset) GetProtocol() *models.Protocol {
	if x != nil {
		return x.Protocol
	}
	return nil
}

func (x *ProtocolAsset) GetAsset() *models.Asset {
	if x != nil {
		return x.Asset
	}
	return nil
}

func (x *ProtocolAsset) GetChain() *models.Chain {
	if x != nil {
		return x.Chain
	}
	return nil
}

func (x *ProtocolAsset) GetCreationBlock() *wrapperspb.UInt64Value {
	if x != nil {
		return x.CreationBlock
	}
	return nil
}

func (x *ProtocolAsset) GetCreationDate() *timestamppb.Timestamp {
	if x != nil {
		return x.CreationDate
	}
	return nil
}

func (x *ProtocolAsset) GetContractAddress() *wrapperspb.StringValue {
	if x != nil {
		return x.ContractAddress
	}
	return nil
}

func (x *ProtocolAsset) GetDecimals() *wrapperspb.UInt32Value {
	if x != nil {
		return x.Decimals
	}
	return nil
}

type MarketableProtocolAsset struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProtocolAsset             *ProtocolAsset   `protobuf:"bytes,1,opt,name=protocol_asset,json=protocolAsset,proto3" json:"protocol_asset,omitempty"`
	Market                    *models.Exchange `protobuf:"bytes,2,opt,name=market,proto3" json:"market,omitempty"`
	MarketableProtocolAssetID uint64           `protobuf:"varint,3,opt,name=marketable_protocol_assetID,json=marketableProtocolAssetID,proto3" json:"marketable_protocol_assetID,omitempty"`
}

func (x *MarketableProtocolAsset) Reset() {
	*x = MarketableProtocolAsset{}
	if protoimpl.UnsafeEnabled {
		mi := &file_asset_data_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarketableProtocolAsset) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarketableProtocolAsset) ProtoMessage() {}

func (x *MarketableProtocolAsset) ProtoReflect() protoreflect.Message {
	mi := &file_asset_data_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarketableProtocolAsset.ProtoReflect.Descriptor instead.
func (*MarketableProtocolAsset) Descriptor() ([]byte, []int) {
	return file_asset_data_proto_rawDescGZIP(), []int{1}
}

func (x *MarketableProtocolAsset) GetProtocolAsset() *ProtocolAsset {
	if x != nil {
		return x.ProtocolAsset
	}
	return nil
}

func (x *MarketableProtocolAsset) GetMarket() *models.Exchange {
	if x != nil {
		return x.Market
	}
	return nil
}

func (x *MarketableProtocolAsset) GetMarketableProtocolAssetID() uint64 {
	if x != nil {
		return x.MarketableProtocolAssetID
	}
	return 0
}

var File_asset_data_proto protoreflect.FileDescriptor

var file_asset_data_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x73, 0x73, 0x65, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x1a, 0x34, 0x67, 0x69, 0x74, 0x6c,
	0x61, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x74, 0x69, 0x63, 0x6b,
	0x73, 0x2f, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x72, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x2f, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xbb, 0x03, 0x0a, 0x0d, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x41, 0x73,
	0x73, 0x65, 0x74, 0x12, 0x29, 0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x5f,
	0x61, 0x73, 0x73, 0x65, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x41, 0x73, 0x73, 0x65, 0x74, 0x49, 0x44, 0x12, 0x2c,
	0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x23, 0x0a, 0x05,
	0x61, 0x73, 0x73, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x05, 0x61, 0x73, 0x73, 0x65,
	0x74, 0x12, 0x23, 0x0a, 0x05, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0d, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x52,
	0x05, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x12, 0x43, 0x0a, 0x0e, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0d, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x3f, 0x0a, 0x0d, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x65, 0x12, 0x47, 0x0a, 0x10,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x38, 0x0a, 0x08, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c,
	0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x33, 0x32,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x73, 0x22,
	0xc1, 0x01, 0x0a, 0x17, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x3c, 0x0a, 0x0e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x41, 0x73, 0x73, 0x65, 0x74, 0x52, 0x0d, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x28, 0x0a, 0x06, 0x6d, 0x61, 0x72,
	0x6b, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x2e, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x06, 0x6d, 0x61, 0x72,
	0x6b, 0x65, 0x74, 0x12, 0x3e, 0x0a, 0x1b, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74,
	0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x19, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x41, 0x73, 0x73, 0x65,
	0x74, 0x49, 0x44, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x74, 0x69, 0x63, 0x6b, 0x73, 0x2f, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x2d, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_asset_data_proto_rawDescOnce sync.Once
	file_asset_data_proto_rawDescData = file_asset_data_proto_rawDesc
)

func file_asset_data_proto_rawDescGZIP() []byte {
	file_asset_data_proto_rawDescOnce.Do(func() {
		file_asset_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_asset_data_proto_rawDescData)
	})
	return file_asset_data_proto_rawDescData
}

var file_asset_data_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_asset_data_proto_goTypes = []interface{}{
	(*ProtocolAsset)(nil),           // 0: models.ProtocolAsset
	(*MarketableProtocolAsset)(nil), // 1: models.MarketableProtocolAsset
	(*models.Protocol)(nil),         // 2: models.Protocol
	(*models.Asset)(nil),            // 3: models.Asset
	(*models.Chain)(nil),            // 4: models.Chain
	(*wrapperspb.UInt64Value)(nil),  // 5: google.protobuf.UInt64Value
	(*timestamppb.Timestamp)(nil),   // 6: google.protobuf.Timestamp
	(*wrapperspb.StringValue)(nil),  // 7: google.protobuf.StringValue
	(*wrapperspb.UInt32Value)(nil),  // 8: google.protobuf.UInt32Value
	(*models.Exchange)(nil),         // 9: models.Exchange
}
var file_asset_data_proto_depIdxs = []int32{
	2, // 0: models.ProtocolAsset.protocol:type_name -> models.Protocol
	3, // 1: models.ProtocolAsset.asset:type_name -> models.Asset
	4, // 2: models.ProtocolAsset.chain:type_name -> models.Chain
	5, // 3: models.ProtocolAsset.creation_block:type_name -> google.protobuf.UInt64Value
	6, // 4: models.ProtocolAsset.creation_date:type_name -> google.protobuf.Timestamp
	7, // 5: models.ProtocolAsset.contract_address:type_name -> google.protobuf.StringValue
	8, // 6: models.ProtocolAsset.decimals:type_name -> google.protobuf.UInt32Value
	0, // 7: models.MarketableProtocolAsset.protocol_asset:type_name -> models.ProtocolAsset
	9, // 8: models.MarketableProtocolAsset.market:type_name -> models.Exchange
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_asset_data_proto_init() }
func file_asset_data_proto_init() {
	if File_asset_data_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_asset_data_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtocolAsset); i {
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
		file_asset_data_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarketableProtocolAsset); i {
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
			RawDescriptor: file_asset_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_asset_data_proto_goTypes,
		DependencyIndexes: file_asset_data_proto_depIdxs,
		MessageInfos:      file_asset_data_proto_msgTypes,
	}.Build()
	File_asset_data_proto = out.File
	file_asset_data_proto_rawDesc = nil
	file_asset_data_proto_goTypes = nil
	file_asset_data_proto_depIdxs = nil
}
