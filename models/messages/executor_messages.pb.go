package messages

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	models1 "gitlab.com/alphaticks/alphac/models"
	_ "gitlab.com/alphaticks/gorderbook"
	models "gitlab.com/alphaticks/xchanger/models"
	io "io"
	math "math"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Instrument struct {
	SecurityID uint64           `protobuf:"varint,1,opt,name=securityID,proto3" json:"securityID,omitempty"`
	Exchange   *models.Exchange `protobuf:"bytes,2,opt,name=exchange,proto3" json:"exchange,omitempty"`
	Symbol     string           `protobuf:"bytes,3,opt,name=symbol,proto3" json:"symbol,omitempty"`
}

func (m *Instrument) Reset()      { *m = Instrument{} }
func (*Instrument) ProtoMessage() {}
func (*Instrument) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{0}
}
func (m *Instrument) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Instrument) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Instrument.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Instrument) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Instrument.Merge(m, src)
}
func (m *Instrument) XXX_Size() int {
	return m.Size()
}
func (m *Instrument) XXX_DiscardUnknown() {
	xxx_messageInfo_Instrument.DiscardUnknown(m)
}

var xxx_messageInfo_Instrument proto.InternalMessageInfo

func (m *Instrument) GetSecurityID() uint64 {
	if m != nil {
		return m.SecurityID
	}
	return 0
}

func (m *Instrument) GetExchange() *models.Exchange {
	if m != nil {
		return m.Exchange
	}
	return nil
}

func (m *Instrument) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

type MarketDataRequest struct {
	RequestID  uint64      `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe  bool        `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
	Instrument *Instrument `protobuf:"bytes,3,opt,name=instrument,proto3" json:"instrument,omitempty"`
}

func (m *MarketDataRequest) Reset()      { *m = MarketDataRequest{} }
func (*MarketDataRequest) ProtoMessage() {}
func (*MarketDataRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{1}
}
func (m *MarketDataRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MarketDataRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MarketDataRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MarketDataRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MarketDataRequest.Merge(m, src)
}
func (m *MarketDataRequest) XXX_Size() int {
	return m.Size()
}
func (m *MarketDataRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MarketDataRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MarketDataRequest proto.InternalMessageInfo

func (m *MarketDataRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *MarketDataRequest) GetSubscribe() bool {
	if m != nil {
		return m.Subscribe
	}
	return false
}

func (m *MarketDataRequest) GetInstrument() *Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

type MarketDataRequestReject struct {
	RequestID uint64 `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Reason    string `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
}

func (m *MarketDataRequestReject) Reset()      { *m = MarketDataRequestReject{} }
func (*MarketDataRequestReject) ProtoMessage() {}
func (*MarketDataRequestReject) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{2}
}
func (m *MarketDataRequestReject) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MarketDataRequestReject) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MarketDataRequestReject.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MarketDataRequestReject) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MarketDataRequestReject.Merge(m, src)
}
func (m *MarketDataRequestReject) XXX_Size() int {
	return m.Size()
}
func (m *MarketDataRequestReject) XXX_DiscardUnknown() {
	xxx_messageInfo_MarketDataRequestReject.DiscardUnknown(m)
}

var xxx_messageInfo_MarketDataRequestReject proto.InternalMessageInfo

func (m *MarketDataRequestReject) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *MarketDataRequestReject) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type MarketDataSnapshot struct {
	RequestID  uint64                     `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID uint64                     `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	SnapshotL2 *models1.OBL2Snapshot      `protobuf:"bytes,3,opt,name=snapshotL2,proto3" json:"snapshotL2,omitempty"`
	SnapshotL3 *models1.OBL3Snapshot      `protobuf:"bytes,4,opt,name=snapshotL3,proto3" json:"snapshotL3,omitempty"`
	Trades     []*models1.AggregatedTrade `protobuf:"bytes,5,rep,name=trades,proto3" json:"trades,omitempty"`
}

func (m *MarketDataSnapshot) Reset()      { *m = MarketDataSnapshot{} }
func (*MarketDataSnapshot) ProtoMessage() {}
func (*MarketDataSnapshot) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{3}
}
func (m *MarketDataSnapshot) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MarketDataSnapshot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MarketDataSnapshot.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MarketDataSnapshot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MarketDataSnapshot.Merge(m, src)
}
func (m *MarketDataSnapshot) XXX_Size() int {
	return m.Size()
}
func (m *MarketDataSnapshot) XXX_DiscardUnknown() {
	xxx_messageInfo_MarketDataSnapshot.DiscardUnknown(m)
}

var xxx_messageInfo_MarketDataSnapshot proto.InternalMessageInfo

func (m *MarketDataSnapshot) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *MarketDataSnapshot) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *MarketDataSnapshot) GetSnapshotL2() *models1.OBL2Snapshot {
	if m != nil {
		return m.SnapshotL2
	}
	return nil
}

func (m *MarketDataSnapshot) GetSnapshotL3() *models1.OBL3Snapshot {
	if m != nil {
		return m.SnapshotL3
	}
	return nil
}

func (m *MarketDataSnapshot) GetTrades() []*models1.AggregatedTrade {
	if m != nil {
		return m.Trades
	}
	return nil
}

type MarketDataIncrementalRefresh struct {
	RequestID  uint64                     `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID uint64                     `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	UpdateL2   *models1.OBL2Update        `protobuf:"bytes,3,opt,name=updateL2,proto3" json:"updateL2,omitempty"`
	UpdateL3   *models1.OBL3Update        `protobuf:"bytes,4,opt,name=updateL3,proto3" json:"updateL3,omitempty"`
	Trades     []*models1.AggregatedTrade `protobuf:"bytes,5,rep,name=trades,proto3" json:"trades,omitempty"`
}

func (m *MarketDataIncrementalRefresh) Reset()      { *m = MarketDataIncrementalRefresh{} }
func (*MarketDataIncrementalRefresh) ProtoMessage() {}
func (*MarketDataIncrementalRefresh) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{4}
}
func (m *MarketDataIncrementalRefresh) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MarketDataIncrementalRefresh) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MarketDataIncrementalRefresh.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MarketDataIncrementalRefresh) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MarketDataIncrementalRefresh.Merge(m, src)
}
func (m *MarketDataIncrementalRefresh) XXX_Size() int {
	return m.Size()
}
func (m *MarketDataIncrementalRefresh) XXX_DiscardUnknown() {
	xxx_messageInfo_MarketDataIncrementalRefresh.DiscardUnknown(m)
}

var xxx_messageInfo_MarketDataIncrementalRefresh proto.InternalMessageInfo

func (m *MarketDataIncrementalRefresh) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *MarketDataIncrementalRefresh) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *MarketDataIncrementalRefresh) GetUpdateL2() *models1.OBL2Update {
	if m != nil {
		return m.UpdateL2
	}
	return nil
}

func (m *MarketDataIncrementalRefresh) GetUpdateL3() *models1.OBL3Update {
	if m != nil {
		return m.UpdateL3
	}
	return nil
}

func (m *MarketDataIncrementalRefresh) GetTrades() []*models1.AggregatedTrade {
	if m != nil {
		return m.Trades
	}
	return nil
}

type SecurityListRequest struct {
	RequestID uint64 `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe bool   `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
}

func (m *SecurityListRequest) Reset()      { *m = SecurityListRequest{} }
func (*SecurityListRequest) ProtoMessage() {}
func (*SecurityListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{5}
}
func (m *SecurityListRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SecurityListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SecurityListRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SecurityListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SecurityListRequest.Merge(m, src)
}
func (m *SecurityListRequest) XXX_Size() int {
	return m.Size()
}
func (m *SecurityListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SecurityListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SecurityListRequest proto.InternalMessageInfo

func (m *SecurityListRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *SecurityListRequest) GetSubscribe() bool {
	if m != nil {
		return m.Subscribe
	}
	return false
}

type SecurityList struct {
	RequestID  uint64              `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID uint64              `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Error      string              `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	Securities []*models1.Security `protobuf:"bytes,4,rep,name=securities,proto3" json:"securities,omitempty"`
}

func (m *SecurityList) Reset()      { *m = SecurityList{} }
func (*SecurityList) ProtoMessage() {}
func (*SecurityList) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{6}
}
func (m *SecurityList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SecurityList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SecurityList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SecurityList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SecurityList.Merge(m, src)
}
func (m *SecurityList) XXX_Size() int {
	return m.Size()
}
func (m *SecurityList) XXX_DiscardUnknown() {
	xxx_messageInfo_SecurityList.DiscardUnknown(m)
}

var xxx_messageInfo_SecurityList proto.InternalMessageInfo

func (m *SecurityList) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *SecurityList) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *SecurityList) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *SecurityList) GetSecurities() []*models1.Security {
	if m != nil {
		return m.Securities
	}
	return nil
}

func init() {
	proto.RegisterType((*Instrument)(nil), "messages.Instrument")
	proto.RegisterType((*MarketDataRequest)(nil), "messages.MarketDataRequest")
	proto.RegisterType((*MarketDataRequestReject)(nil), "messages.MarketDataRequestReject")
	proto.RegisterType((*MarketDataSnapshot)(nil), "messages.MarketDataSnapshot")
	proto.RegisterType((*MarketDataIncrementalRefresh)(nil), "messages.MarketDataIncrementalRefresh")
	proto.RegisterType((*SecurityListRequest)(nil), "messages.SecurityListRequest")
	proto.RegisterType((*SecurityList)(nil), "messages.SecurityList")
}

func init() { proto.RegisterFile("executor_messages.proto", fileDescriptor_350c53ba9303a7e6) }

var fileDescriptor_350c53ba9303a7e6 = []byte{
	// 589 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0xbd, 0x6e, 0xd4, 0x40,
	0x10, 0xbe, 0x25, 0x97, 0xd3, 0x65, 0x43, 0x01, 0x26, 0x4a, 0xac, 0x28, 0x5a, 0x4e, 0xae, 0xae,
	0x00, 0x1b, 0xf9, 0xc2, 0x03, 0x10, 0x85, 0xe2, 0xa4, 0xa0, 0x88, 0x0d, 0xd4, 0xd1, 0xda, 0x9e,
	0xf8, 0x4c, 0x6c, 0xaf, 0xd9, 0x5d, 0x4b, 0x49, 0x47, 0x43, 0x4f, 0xc1, 0x43, 0xf0, 0x28, 0x94,
	0x29, 0x53, 0x12, 0xa7, 0xa1, 0x41, 0xe4, 0x11, 0x90, 0xd7, 0xbf, 0x90, 0x88, 0x20, 0xd2, 0xed,
	0xcc, 0x7e, 0xdf, 0x37, 0xf3, 0xad, 0x66, 0x16, 0x6f, 0xc0, 0x09, 0xf8, 0xb9, 0xe2, 0xe2, 0x30,
	0x01, 0x29, 0x59, 0x08, 0xd2, 0xce, 0x04, 0x57, 0xdc, 0x18, 0x37, 0xf1, 0xe6, 0xd3, 0x30, 0x52,
	0x8b, 0xdc, 0xb3, 0x7d, 0x9e, 0x38, 0x21, 0x0f, 0xb9, 0xa3, 0x01, 0x5e, 0x7e, 0xa4, 0x23, 0x1d,
	0xe8, 0x53, 0x45, 0xdc, 0x7c, 0x1c, 0x72, 0x1e, 0xc6, 0xd0, 0xa1, 0x54, 0x94, 0x80, 0x54, 0x2c,
	0xc9, 0x6a, 0x40, 0xa9, 0x17, 0xb3, 0x4a, 0x8f, 0xc5, 0xd9, 0x82, 0xa9, 0xc8, 0x3f, 0x96, 0x4e,
	0xc8, 0x45, 0x00, 0xc2, 0xe3, 0xfc, 0xb8, 0xa2, 0xd6, 0x8d, 0x6c, 0x3e, 0xbf, 0x19, 0xae, 0x8f,
	0xbe, 0x93, 0xf0, 0x00, 0x62, 0xe9, 0x24, 0x4c, 0x1c, 0x83, 0x3a, 0x0c, 0x98, 0x62, 0x35, 0xcd,
	0xfd, 0x17, 0x9a, 0x04, 0x3f, 0x17, 0x91, 0x3a, 0xfd, 0x3b, 0xe7, 0xc4, 0x5f, 0xb0, 0x34, 0x04,
	0xd1, 0xb0, 0xfa, 0xed, 0x59, 0x02, 0xe3, 0x79, 0x2a, 0x95, 0xc8, 0x13, 0x48, 0x95, 0x41, 0x30,
	0x6e, 0x34, 0xe7, 0xbb, 0x26, 0x9a, 0xa0, 0xe9, 0x90, 0xf6, 0x32, 0xc6, 0x13, 0x3c, 0x86, 0x5a,
	0xce, 0xbc, 0x37, 0x41, 0xd3, 0x55, 0xf7, 0x81, 0x5d, 0xa9, 0xda, 0x2f, 0xeb, 0x3c, 0x6d, 0x11,
	0xc6, 0x3a, 0x1e, 0xc9, 0xd3, 0xc4, 0xe3, 0xb1, 0xb9, 0x34, 0x41, 0xd3, 0x15, 0x5a, 0x47, 0xd6,
	0x47, 0x84, 0x1f, 0xbe, 0xd2, 0x8e, 0x77, 0x99, 0x62, 0x14, 0xde, 0xe7, 0x20, 0x95, 0xb1, 0x85,
	0x57, 0x44, 0x75, 0x6c, 0x4b, 0x77, 0x89, 0xf2, 0x56, 0xe6, 0x9e, 0xf4, 0x45, 0xe4, 0x55, 0xa5,
	0xc7, 0xb4, 0x4b, 0x18, 0xdb, 0x18, 0x47, 0xad, 0x0b, 0x5d, 0x6d, 0xd5, 0x5d, 0xb3, 0xdb, 0x91,
	0xe8, 0x1c, 0xd2, 0x1e, 0xce, 0xda, 0xc7, 0x1b, 0xd7, 0xda, 0xa0, 0xf0, 0x0e, 0xfc, 0xdb, 0x9a,
	0x59, 0xc7, 0x23, 0x01, 0x4c, 0xf2, 0x54, 0x77, 0xb2, 0x42, 0xeb, 0xc8, 0xfa, 0x89, 0xb0, 0xd1,
	0x29, 0x1e, 0xa4, 0x2c, 0x93, 0x0b, 0x7e, 0x9b, 0x18, 0xc1, 0x58, 0x80, 0xcc, 0x78, 0x2a, 0x61,
	0xbe, 0xab, 0x05, 0x87, 0xb4, 0x97, 0x29, 0xbd, 0xc9, 0x5a, 0x69, 0xcf, 0xed, 0xbc, 0x55, 0xaf,
	0xbe, 0xbf, 0xb3, 0xe7, 0x36, 0x75, 0x68, 0x0f, 0xf7, 0x1b, 0x6b, 0x66, 0x0e, 0xaf, 0xb1, 0x66,
	0x37, 0xb0, 0x66, 0x86, 0x83, 0x47, 0x4a, 0xb0, 0x00, 0xa4, 0xb9, 0x3c, 0x59, 0x9a, 0xae, 0xba,
	0x1b, 0x0d, 0xe3, 0x45, 0x18, 0x0a, 0x08, 0x99, 0x82, 0xe0, 0x4d, 0x79, 0x4f, 0x6b, 0x98, 0xf5,
	0x03, 0xe1, 0xad, 0xce, 0xf1, 0x3c, 0xf5, 0x05, 0x94, 0x4f, 0xcb, 0x62, 0x0a, 0x47, 0x02, 0xe4,
	0xe2, 0x8e, 0xde, 0x6d, 0x3c, 0xce, 0xb3, 0x80, 0x29, 0x68, 0x9d, 0x1b, 0x7d, 0xe7, 0x6f, 0xf5,
	0x1d, 0x6d, 0x31, 0x3d, 0x7c, 0xe3, 0xb9, 0x8f, 0x9f, 0xfd, 0x81, 0xff, 0x0f, 0xbf, 0xaf, 0xf1,
	0xa3, 0x83, 0x7a, 0x1d, 0xf6, 0xa2, 0x72, 0x5a, 0xee, 0x3c, 0xbb, 0xd6, 0x67, 0x84, 0xef, 0xf7,
	0x35, 0xef, 0xf8, 0x64, 0x6b, 0x78, 0x19, 0x84, 0xe0, 0xa2, 0xde, 0xb9, 0x2a, 0x30, 0x9e, 0xb5,
	0x8b, 0x1d, 0x81, 0x34, 0x87, 0xda, 0x6c, 0xbb, 0xba, 0x4d, 0x75, 0xda, 0xc3, 0xec, 0x6c, 0x9f,
	0x5d, 0x90, 0xc1, 0xf9, 0x05, 0x19, 0x5c, 0x5d, 0x10, 0xf4, 0xa1, 0x20, 0xe8, 0x4b, 0x41, 0xd0,
	0xd7, 0x82, 0xa0, 0xb3, 0x82, 0xa0, 0x6f, 0x05, 0x41, 0xdf, 0x0b, 0x32, 0xb8, 0x2a, 0x08, 0xfa,
	0x74, 0x49, 0x06, 0x67, 0x97, 0x64, 0x70, 0x7e, 0x49, 0x06, 0xde, 0x48, 0xff, 0x2a, 0xb3, 0x5f,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x9d, 0xd7, 0x79, 0x6f, 0x98, 0x05, 0x00, 0x00,
}

func (this *Instrument) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Instrument)
	if !ok {
		that2, ok := that.(Instrument)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.SecurityID != that1.SecurityID {
		return false
	}
	if !this.Exchange.Equal(that1.Exchange) {
		return false
	}
	if this.Symbol != that1.Symbol {
		return false
	}
	return true
}
func (this *MarketDataRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MarketDataRequest)
	if !ok {
		that2, ok := that.(MarketDataRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.RequestID != that1.RequestID {
		return false
	}
	if this.Subscribe != that1.Subscribe {
		return false
	}
	if !this.Instrument.Equal(that1.Instrument) {
		return false
	}
	return true
}
func (this *MarketDataRequestReject) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MarketDataRequestReject)
	if !ok {
		that2, ok := that.(MarketDataRequestReject)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.RequestID != that1.RequestID {
		return false
	}
	if this.Reason != that1.Reason {
		return false
	}
	return true
}
func (this *MarketDataSnapshot) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MarketDataSnapshot)
	if !ok {
		that2, ok := that.(MarketDataSnapshot)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.RequestID != that1.RequestID {
		return false
	}
	if this.ResponseID != that1.ResponseID {
		return false
	}
	if !this.SnapshotL2.Equal(that1.SnapshotL2) {
		return false
	}
	if !this.SnapshotL3.Equal(that1.SnapshotL3) {
		return false
	}
	if len(this.Trades) != len(that1.Trades) {
		return false
	}
	for i := range this.Trades {
		if !this.Trades[i].Equal(that1.Trades[i]) {
			return false
		}
	}
	return true
}
func (this *MarketDataIncrementalRefresh) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MarketDataIncrementalRefresh)
	if !ok {
		that2, ok := that.(MarketDataIncrementalRefresh)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.RequestID != that1.RequestID {
		return false
	}
	if this.ResponseID != that1.ResponseID {
		return false
	}
	if !this.UpdateL2.Equal(that1.UpdateL2) {
		return false
	}
	if !this.UpdateL3.Equal(that1.UpdateL3) {
		return false
	}
	if len(this.Trades) != len(that1.Trades) {
		return false
	}
	for i := range this.Trades {
		if !this.Trades[i].Equal(that1.Trades[i]) {
			return false
		}
	}
	return true
}
func (this *SecurityListRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SecurityListRequest)
	if !ok {
		that2, ok := that.(SecurityListRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.RequestID != that1.RequestID {
		return false
	}
	if this.Subscribe != that1.Subscribe {
		return false
	}
	return true
}
func (this *SecurityList) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SecurityList)
	if !ok {
		that2, ok := that.(SecurityList)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.RequestID != that1.RequestID {
		return false
	}
	if this.ResponseID != that1.ResponseID {
		return false
	}
	if this.Error != that1.Error {
		return false
	}
	if len(this.Securities) != len(that1.Securities) {
		return false
	}
	for i := range this.Securities {
		if !this.Securities[i].Equal(that1.Securities[i]) {
			return false
		}
	}
	return true
}
func (this *Instrument) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.Instrument{")
	s = append(s, "SecurityID: "+fmt.Sprintf("%#v", this.SecurityID)+",\n")
	if this.Exchange != nil {
		s = append(s, "Exchange: "+fmt.Sprintf("%#v", this.Exchange)+",\n")
	}
	s = append(s, "Symbol: "+fmt.Sprintf("%#v", this.Symbol)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *MarketDataRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.MarketDataRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Subscribe: "+fmt.Sprintf("%#v", this.Subscribe)+",\n")
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *MarketDataRequestReject) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&messages.MarketDataRequestReject{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Reason: "+fmt.Sprintf("%#v", this.Reason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *MarketDataSnapshot) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.MarketDataSnapshot{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.SnapshotL2 != nil {
		s = append(s, "SnapshotL2: "+fmt.Sprintf("%#v", this.SnapshotL2)+",\n")
	}
	if this.SnapshotL3 != nil {
		s = append(s, "SnapshotL3: "+fmt.Sprintf("%#v", this.SnapshotL3)+",\n")
	}
	if this.Trades != nil {
		s = append(s, "Trades: "+fmt.Sprintf("%#v", this.Trades)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *MarketDataIncrementalRefresh) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.MarketDataIncrementalRefresh{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.UpdateL2 != nil {
		s = append(s, "UpdateL2: "+fmt.Sprintf("%#v", this.UpdateL2)+",\n")
	}
	if this.UpdateL3 != nil {
		s = append(s, "UpdateL3: "+fmt.Sprintf("%#v", this.UpdateL3)+",\n")
	}
	if this.Trades != nil {
		s = append(s, "Trades: "+fmt.Sprintf("%#v", this.Trades)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *SecurityListRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&messages.SecurityListRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Subscribe: "+fmt.Sprintf("%#v", this.Subscribe)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *SecurityList) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.SecurityList{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	s = append(s, "Error: "+fmt.Sprintf("%#v", this.Error)+",\n")
	if this.Securities != nil {
		s = append(s, "Securities: "+fmt.Sprintf("%#v", this.Securities)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringExecutorMessages(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *Instrument) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Instrument) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.SecurityID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SecurityID))
	}
	if m.Exchange != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Exchange.Size()))
		n1, err := m.Exchange.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if len(m.Symbol) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.Symbol)))
		i += copy(dAtA[i:], m.Symbol)
	}
	return i, nil
}

func (m *MarketDataRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MarketDataRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RequestID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
	}
	if m.Subscribe {
		dAtA[i] = 0x10
		i++
		if m.Subscribe {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.Instrument != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Instrument.Size()))
		n2, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func (m *MarketDataRequestReject) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MarketDataRequestReject) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RequestID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
	}
	if len(m.Reason) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.Reason)))
		i += copy(dAtA[i:], m.Reason)
	}
	return i, nil
}

func (m *MarketDataSnapshot) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MarketDataSnapshot) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RequestID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
	}
	if m.ResponseID != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
	}
	if m.SnapshotL2 != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SnapshotL2.Size()))
		n3, err := m.SnapshotL2.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	if m.SnapshotL3 != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SnapshotL3.Size()))
		n4, err := m.SnapshotL3.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n4
	}
	if len(m.Trades) > 0 {
		for _, msg := range m.Trades {
			dAtA[i] = 0x2a
			i++
			i = encodeVarintExecutorMessages(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *MarketDataIncrementalRefresh) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MarketDataIncrementalRefresh) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RequestID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
	}
	if m.ResponseID != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
	}
	if m.UpdateL2 != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.UpdateL2.Size()))
		n5, err := m.UpdateL2.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n5
	}
	if m.UpdateL3 != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.UpdateL3.Size()))
		n6, err := m.UpdateL3.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n6
	}
	if len(m.Trades) > 0 {
		for _, msg := range m.Trades {
			dAtA[i] = 0x2a
			i++
			i = encodeVarintExecutorMessages(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *SecurityListRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SecurityListRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RequestID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
	}
	if m.Subscribe {
		dAtA[i] = 0x10
		i++
		if m.Subscribe {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	return i, nil
}

func (m *SecurityList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SecurityList) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RequestID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
	}
	if m.ResponseID != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
	}
	if len(m.Error) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.Error)))
		i += copy(dAtA[i:], m.Error)
	}
	if len(m.Securities) > 0 {
		for _, msg := range m.Securities {
			dAtA[i] = 0x22
			i++
			i = encodeVarintExecutorMessages(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeVarintExecutorMessages(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Instrument) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SecurityID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.SecurityID))
	}
	if m.Exchange != nil {
		l = m.Exchange.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	l = len(m.Symbol)
	if l > 0 {
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *MarketDataRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
	}
	if m.Subscribe {
		n += 2
	}
	if m.Instrument != nil {
		l = m.Instrument.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *MarketDataRequestReject) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
	}
	l = len(m.Reason)
	if l > 0 {
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *MarketDataSnapshot) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
	}
	if m.ResponseID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.ResponseID))
	}
	if m.SnapshotL2 != nil {
		l = m.SnapshotL2.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.SnapshotL3 != nil {
		l = m.SnapshotL3.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if len(m.Trades) > 0 {
		for _, e := range m.Trades {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	return n
}

func (m *MarketDataIncrementalRefresh) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
	}
	if m.ResponseID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.ResponseID))
	}
	if m.UpdateL2 != nil {
		l = m.UpdateL2.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.UpdateL3 != nil {
		l = m.UpdateL3.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if len(m.Trades) > 0 {
		for _, e := range m.Trades {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	return n
}

func (m *SecurityListRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
	}
	if m.Subscribe {
		n += 2
	}
	return n
}

func (m *SecurityList) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
	}
	if m.ResponseID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.ResponseID))
	}
	l = len(m.Error)
	if l > 0 {
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if len(m.Securities) > 0 {
		for _, e := range m.Securities {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	return n
}

func sovExecutorMessages(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozExecutorMessages(x uint64) (n int) {
	return sovExecutorMessages(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Instrument) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Instrument{`,
		`SecurityID:` + fmt.Sprintf("%v", this.SecurityID) + `,`,
		`Exchange:` + strings.Replace(fmt.Sprintf("%v", this.Exchange), "Exchange", "models.Exchange", 1) + `,`,
		`Symbol:` + fmt.Sprintf("%v", this.Symbol) + `,`,
		`}`,
	}, "")
	return s
}
func (this *MarketDataRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&MarketDataRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Subscribe:` + fmt.Sprintf("%v", this.Subscribe) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "Instrument", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *MarketDataRequestReject) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&MarketDataRequestReject{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Reason:` + fmt.Sprintf("%v", this.Reason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *MarketDataSnapshot) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&MarketDataSnapshot{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`SnapshotL2:` + strings.Replace(fmt.Sprintf("%v", this.SnapshotL2), "OBL2Snapshot", "models1.OBL2Snapshot", 1) + `,`,
		`SnapshotL3:` + strings.Replace(fmt.Sprintf("%v", this.SnapshotL3), "OBL3Snapshot", "models1.OBL3Snapshot", 1) + `,`,
		`Trades:` + strings.Replace(fmt.Sprintf("%v", this.Trades), "AggregatedTrade", "models1.AggregatedTrade", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *MarketDataIncrementalRefresh) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&MarketDataIncrementalRefresh{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`UpdateL2:` + strings.Replace(fmt.Sprintf("%v", this.UpdateL2), "OBL2Update", "models1.OBL2Update", 1) + `,`,
		`UpdateL3:` + strings.Replace(fmt.Sprintf("%v", this.UpdateL3), "OBL3Update", "models1.OBL3Update", 1) + `,`,
		`Trades:` + strings.Replace(fmt.Sprintf("%v", this.Trades), "AggregatedTrade", "models1.AggregatedTrade", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *SecurityListRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&SecurityListRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Subscribe:` + fmt.Sprintf("%v", this.Subscribe) + `,`,
		`}`,
	}, "")
	return s
}
func (this *SecurityList) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&SecurityList{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Error:` + fmt.Sprintf("%v", this.Error) + `,`,
		`Securities:` + strings.Replace(fmt.Sprintf("%v", this.Securities), "Security", "models1.Security", 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringExecutorMessages(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Instrument) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowExecutorMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Instrument: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Instrument: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SecurityID", wireType)
			}
			m.SecurityID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SecurityID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Exchange", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Exchange == nil {
				m.Exchange = &models.Exchange{}
			}
			if err := m.Exchange.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Symbol", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Symbol = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipExecutorMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MarketDataRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowExecutorMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MarketDataRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MarketDataRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequestID", wireType)
			}
			m.RequestID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RequestID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subscribe", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Subscribe = bool(v != 0)
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Instrument", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Instrument == nil {
				m.Instrument = &Instrument{}
			}
			if err := m.Instrument.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipExecutorMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MarketDataRequestReject) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowExecutorMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MarketDataRequestReject: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MarketDataRequestReject: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequestID", wireType)
			}
			m.RequestID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RequestID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Reason", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Reason = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipExecutorMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MarketDataSnapshot) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowExecutorMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MarketDataSnapshot: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MarketDataSnapshot: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequestID", wireType)
			}
			m.RequestID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RequestID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ResponseID", wireType)
			}
			m.ResponseID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ResponseID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SnapshotL2", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SnapshotL2 == nil {
				m.SnapshotL2 = &models1.OBL2Snapshot{}
			}
			if err := m.SnapshotL2.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SnapshotL3", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SnapshotL3 == nil {
				m.SnapshotL3 = &models1.OBL3Snapshot{}
			}
			if err := m.SnapshotL3.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Trades", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Trades = append(m.Trades, &models1.AggregatedTrade{})
			if err := m.Trades[len(m.Trades)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipExecutorMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MarketDataIncrementalRefresh) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowExecutorMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MarketDataIncrementalRefresh: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MarketDataIncrementalRefresh: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequestID", wireType)
			}
			m.RequestID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RequestID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ResponseID", wireType)
			}
			m.ResponseID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ResponseID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpdateL2", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.UpdateL2 == nil {
				m.UpdateL2 = &models1.OBL2Update{}
			}
			if err := m.UpdateL2.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpdateL3", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.UpdateL3 == nil {
				m.UpdateL3 = &models1.OBL3Update{}
			}
			if err := m.UpdateL3.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Trades", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Trades = append(m.Trades, &models1.AggregatedTrade{})
			if err := m.Trades[len(m.Trades)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipExecutorMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SecurityListRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowExecutorMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SecurityListRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SecurityListRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequestID", wireType)
			}
			m.RequestID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RequestID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subscribe", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Subscribe = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipExecutorMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SecurityList) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowExecutorMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SecurityList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SecurityList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequestID", wireType)
			}
			m.RequestID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RequestID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ResponseID", wireType)
			}
			m.ResponseID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ResponseID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Error", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Error = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Securities", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Securities = append(m.Securities, &models1.Security{})
			if err := m.Securities[len(m.Securities)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipExecutorMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthExecutorMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipExecutorMessages(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowExecutorMessages
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthExecutorMessages
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthExecutorMessages
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowExecutorMessages
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipExecutorMessages(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthExecutorMessages
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthExecutorMessages = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowExecutorMessages   = fmt.Errorf("proto: integer overflow")
)
