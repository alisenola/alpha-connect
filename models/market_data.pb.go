package models

import (
	encoding_binary "encoding/binary"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
	gorderbook "gitlab.com/alphaticks/gorderbook"
	io "io"
	math "math"
	reflect "reflect"
	strconv "strconv"
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

type StatType int32

const (
	IndexValue              StatType = 0
	OpeningPrice            StatType = 1
	ClosingPrice            StatType = 2
	SettlementPrice         StatType = 3
	TradingSessionHighPrice StatType = 4
	TradingSessionLowPrice  StatType = 5
	TradingSessionVWAPPrice StatType = 6
	Imbalance               StatType = 7
	TradeVolume             StatType = 8
	OpenInterest            StatType = 9
)

var StatType_name = map[int32]string{
	0: "IndexValue",
	1: "OpeningPrice",
	2: "ClosingPrice",
	3: "SettlementPrice",
	4: "TradingSessionHighPrice",
	5: "TradingSessionLowPrice",
	6: "TradingSessionVWAPPrice",
	7: "Imbalance",
	8: "TradeVolume",
	9: "OpenInterest",
}

var StatType_value = map[string]int32{
	"IndexValue":              0,
	"OpeningPrice":            1,
	"ClosingPrice":            2,
	"SettlementPrice":         3,
	"TradingSessionHighPrice": 4,
	"TradingSessionLowPrice":  5,
	"TradingSessionVWAPPrice": 6,
	"Imbalance":               7,
	"TradeVolume":             8,
	"OpenInterest":            9,
}

func (StatType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{0}
}

type OrderBookAggregation int32

const (
	L1 OrderBookAggregation = 0
	L2 OrderBookAggregation = 1
	L3 OrderBookAggregation = 2
)

var OrderBookAggregation_name = map[int32]string{
	0: "L1",
	1: "L2",
	2: "L3",
}

var OrderBookAggregation_value = map[string]int32{
	"L1": 0,
	"L2": 1,
	"L3": 2,
}

func (OrderBookAggregation) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{1}
}

type OBL1Update struct {
	BestBid   float64          `protobuf:"fixed64,1,opt,name=best_bid,json=bestBid,proto3" json:"best_bid,omitempty"`
	BestAsk   float64          `protobuf:"fixed64,2,opt,name=best_ask,json=bestAsk,proto3" json:"best_ask,omitempty"`
	Timestamp *types.Timestamp `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (m *OBL1Update) Reset()      { *m = OBL1Update{} }
func (*OBL1Update) ProtoMessage() {}
func (*OBL1Update) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{0}
}
func (m *OBL1Update) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OBL1Update) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OBL1Update.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OBL1Update) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OBL1Update.Merge(m, src)
}
func (m *OBL1Update) XXX_Size() int {
	return m.Size()
}
func (m *OBL1Update) XXX_DiscardUnknown() {
	xxx_messageInfo_OBL1Update.DiscardUnknown(m)
}

var xxx_messageInfo_OBL1Update proto.InternalMessageInfo

func (m *OBL1Update) GetBestBid() float64 {
	if m != nil {
		return m.BestBid
	}
	return 0
}

func (m *OBL1Update) GetBestAsk() float64 {
	if m != nil {
		return m.BestAsk
	}
	return 0
}

func (m *OBL1Update) GetTimestamp() *types.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

type OBL1Snapshot struct {
	BestBid   float64          `protobuf:"fixed64,1,opt,name=best_bid,json=bestBid,proto3" json:"best_bid,omitempty"`
	BestAsk   float64          `protobuf:"fixed64,2,opt,name=best_ask,json=bestAsk,proto3" json:"best_ask,omitempty"`
	Timestamp *types.Timestamp `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (m *OBL1Snapshot) Reset()      { *m = OBL1Snapshot{} }
func (*OBL1Snapshot) ProtoMessage() {}
func (*OBL1Snapshot) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{1}
}
func (m *OBL1Snapshot) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OBL1Snapshot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OBL1Snapshot.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OBL1Snapshot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OBL1Snapshot.Merge(m, src)
}
func (m *OBL1Snapshot) XXX_Size() int {
	return m.Size()
}
func (m *OBL1Snapshot) XXX_DiscardUnknown() {
	xxx_messageInfo_OBL1Snapshot.DiscardUnknown(m)
}

var xxx_messageInfo_OBL1Snapshot proto.InternalMessageInfo

func (m *OBL1Snapshot) GetBestBid() float64 {
	if m != nil {
		return m.BestBid
	}
	return 0
}

func (m *OBL1Snapshot) GetBestAsk() float64 {
	if m != nil {
		return m.BestAsk
	}
	return 0
}

func (m *OBL1Snapshot) GetTimestamp() *types.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

type OBL2Update struct {
	Levels    []gorderbook.OrderBookLevel `protobuf:"bytes,1,rep,name=levels,proto3" json:"levels"`
	Timestamp *types.Timestamp            `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Trade     bool                        `protobuf:"varint,4,opt,name=trade,proto3" json:"trade,omitempty"`
}

func (m *OBL2Update) Reset()      { *m = OBL2Update{} }
func (*OBL2Update) ProtoMessage() {}
func (*OBL2Update) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{2}
}
func (m *OBL2Update) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OBL2Update) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OBL2Update.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OBL2Update) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OBL2Update.Merge(m, src)
}
func (m *OBL2Update) XXX_Size() int {
	return m.Size()
}
func (m *OBL2Update) XXX_DiscardUnknown() {
	xxx_messageInfo_OBL2Update.DiscardUnknown(m)
}

var xxx_messageInfo_OBL2Update proto.InternalMessageInfo

func (m *OBL2Update) GetLevels() []gorderbook.OrderBookLevel {
	if m != nil {
		return m.Levels
	}
	return nil
}

func (m *OBL2Update) GetTimestamp() *types.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *OBL2Update) GetTrade() bool {
	if m != nil {
		return m.Trade
	}
	return false
}

type OBL2Snapshot struct {
	Bids          []gorderbook.OrderBookLevel `protobuf:"bytes,2,rep,name=bids,proto3" json:"bids"`
	Asks          []gorderbook.OrderBookLevel `protobuf:"bytes,3,rep,name=asks,proto3" json:"asks"`
	Timestamp     *types.Timestamp            `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	TickPrecision *types.UInt64Value          `protobuf:"bytes,5,opt,name=tick_precision,json=tickPrecision,proto3" json:"tick_precision,omitempty"`
	LotPrecision  *types.UInt64Value          `protobuf:"bytes,6,opt,name=lot_precision,json=lotPrecision,proto3" json:"lot_precision,omitempty"`
}

func (m *OBL2Snapshot) Reset()      { *m = OBL2Snapshot{} }
func (*OBL2Snapshot) ProtoMessage() {}
func (*OBL2Snapshot) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{3}
}
func (m *OBL2Snapshot) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OBL2Snapshot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OBL2Snapshot.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OBL2Snapshot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OBL2Snapshot.Merge(m, src)
}
func (m *OBL2Snapshot) XXX_Size() int {
	return m.Size()
}
func (m *OBL2Snapshot) XXX_DiscardUnknown() {
	xxx_messageInfo_OBL2Snapshot.DiscardUnknown(m)
}

var xxx_messageInfo_OBL2Snapshot proto.InternalMessageInfo

func (m *OBL2Snapshot) GetBids() []gorderbook.OrderBookLevel {
	if m != nil {
		return m.Bids
	}
	return nil
}

func (m *OBL2Snapshot) GetAsks() []gorderbook.OrderBookLevel {
	if m != nil {
		return m.Asks
	}
	return nil
}

func (m *OBL2Snapshot) GetTimestamp() *types.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *OBL2Snapshot) GetTickPrecision() *types.UInt64Value {
	if m != nil {
		return m.TickPrecision
	}
	return nil
}

func (m *OBL2Snapshot) GetLotPrecision() *types.UInt64Value {
	if m != nil {
		return m.LotPrecision
	}
	return nil
}

type OBL3Update struct {
	Bids          []gorderbook.Order `protobuf:"bytes,2,rep,name=bids,proto3" json:"bids"`
	Asks          []gorderbook.Order `protobuf:"bytes,3,rep,name=asks,proto3" json:"asks"`
	Timestamp     *types.Timestamp   `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	TickPrecision *types.UInt64Value `protobuf:"bytes,5,opt,name=tick_precision,json=tickPrecision,proto3" json:"tick_precision,omitempty"`
	LotPrecision  *types.UInt64Value `protobuf:"bytes,6,opt,name=lot_precision,json=lotPrecision,proto3" json:"lot_precision,omitempty"`
}

func (m *OBL3Update) Reset()      { *m = OBL3Update{} }
func (*OBL3Update) ProtoMessage() {}
func (*OBL3Update) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{4}
}
func (m *OBL3Update) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OBL3Update) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OBL3Update.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OBL3Update) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OBL3Update.Merge(m, src)
}
func (m *OBL3Update) XXX_Size() int {
	return m.Size()
}
func (m *OBL3Update) XXX_DiscardUnknown() {
	xxx_messageInfo_OBL3Update.DiscardUnknown(m)
}

var xxx_messageInfo_OBL3Update proto.InternalMessageInfo

func (m *OBL3Update) GetBids() []gorderbook.Order {
	if m != nil {
		return m.Bids
	}
	return nil
}

func (m *OBL3Update) GetAsks() []gorderbook.Order {
	if m != nil {
		return m.Asks
	}
	return nil
}

func (m *OBL3Update) GetTimestamp() *types.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *OBL3Update) GetTickPrecision() *types.UInt64Value {
	if m != nil {
		return m.TickPrecision
	}
	return nil
}

func (m *OBL3Update) GetLotPrecision() *types.UInt64Value {
	if m != nil {
		return m.LotPrecision
	}
	return nil
}

type OBL3Snapshot struct {
	Bids      []gorderbook.Order `protobuf:"bytes,2,rep,name=bids,proto3" json:"bids"`
	Asks      []gorderbook.Order `protobuf:"bytes,3,rep,name=asks,proto3" json:"asks"`
	Timestamp *types.Timestamp   `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (m *OBL3Snapshot) Reset()      { *m = OBL3Snapshot{} }
func (*OBL3Snapshot) ProtoMessage() {}
func (*OBL3Snapshot) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{5}
}
func (m *OBL3Snapshot) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OBL3Snapshot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OBL3Snapshot.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OBL3Snapshot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OBL3Snapshot.Merge(m, src)
}
func (m *OBL3Snapshot) XXX_Size() int {
	return m.Size()
}
func (m *OBL3Snapshot) XXX_DiscardUnknown() {
	xxx_messageInfo_OBL3Snapshot.DiscardUnknown(m)
}

var xxx_messageInfo_OBL3Snapshot proto.InternalMessageInfo

func (m *OBL3Snapshot) GetBids() []gorderbook.Order {
	if m != nil {
		return m.Bids
	}
	return nil
}

func (m *OBL3Snapshot) GetAsks() []gorderbook.Order {
	if m != nil {
		return m.Asks
	}
	return nil
}

func (m *OBL3Snapshot) GetTimestamp() *types.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

type Trade struct {
	Price    float64 `protobuf:"fixed64,1,opt,name=price,proto3" json:"price,omitempty"`
	Quantity float64 `protobuf:"fixed64,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	ID       uint64  `protobuf:"varint,3,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (m *Trade) Reset()      { *m = Trade{} }
func (*Trade) ProtoMessage() {}
func (*Trade) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{6}
}
func (m *Trade) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Trade) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Trade.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Trade) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Trade.Merge(m, src)
}
func (m *Trade) XXX_Size() int {
	return m.Size()
}
func (m *Trade) XXX_DiscardUnknown() {
	xxx_messageInfo_Trade.DiscardUnknown(m)
}

var xxx_messageInfo_Trade proto.InternalMessageInfo

func (m *Trade) GetPrice() float64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *Trade) GetQuantity() float64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *Trade) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

type AggregatedTrade struct {
	Bid         bool             `protobuf:"varint,1,opt,name=bid,proto3" json:"bid,omitempty"`
	Timestamp   *types.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	AggregateID uint64           `protobuf:"varint,3,opt,name=aggregateID,proto3" json:"aggregateID,omitempty"`
	Trades      []Trade          `protobuf:"bytes,4,rep,name=trades,proto3" json:"trades"`
}

func (m *AggregatedTrade) Reset()      { *m = AggregatedTrade{} }
func (*AggregatedTrade) ProtoMessage() {}
func (*AggregatedTrade) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{7}
}
func (m *AggregatedTrade) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AggregatedTrade) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AggregatedTrade.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AggregatedTrade) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AggregatedTrade.Merge(m, src)
}
func (m *AggregatedTrade) XXX_Size() int {
	return m.Size()
}
func (m *AggregatedTrade) XXX_DiscardUnknown() {
	xxx_messageInfo_AggregatedTrade.DiscardUnknown(m)
}

var xxx_messageInfo_AggregatedTrade proto.InternalMessageInfo

func (m *AggregatedTrade) GetBid() bool {
	if m != nil {
		return m.Bid
	}
	return false
}

func (m *AggregatedTrade) GetTimestamp() *types.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *AggregatedTrade) GetAggregateID() uint64 {
	if m != nil {
		return m.AggregateID
	}
	return 0
}

func (m *AggregatedTrade) GetTrades() []Trade {
	if m != nil {
		return m.Trades
	}
	return nil
}

type Liquidation struct {
	Bid       bool             `protobuf:"varint,1,opt,name=bid,proto3" json:"bid,omitempty"`
	Timestamp *types.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	OrderID   uint64           `protobuf:"varint,3,opt,name=orderID,proto3" json:"orderID,omitempty"`
	Price     float64          `protobuf:"fixed64,4,opt,name=price,proto3" json:"price,omitempty"`
	Quantity  float64          `protobuf:"fixed64,5,opt,name=quantity,proto3" json:"quantity,omitempty"`
}

func (m *Liquidation) Reset()      { *m = Liquidation{} }
func (*Liquidation) ProtoMessage() {}
func (*Liquidation) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{8}
}
func (m *Liquidation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Liquidation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Liquidation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Liquidation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Liquidation.Merge(m, src)
}
func (m *Liquidation) XXX_Size() int {
	return m.Size()
}
func (m *Liquidation) XXX_DiscardUnknown() {
	xxx_messageInfo_Liquidation.DiscardUnknown(m)
}

var xxx_messageInfo_Liquidation proto.InternalMessageInfo

func (m *Liquidation) GetBid() bool {
	if m != nil {
		return m.Bid
	}
	return false
}

func (m *Liquidation) GetTimestamp() *types.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Liquidation) GetOrderID() uint64 {
	if m != nil {
		return m.OrderID
	}
	return 0
}

func (m *Liquidation) GetPrice() float64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *Liquidation) GetQuantity() float64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

type Funding struct {
	Timestamp *types.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Rate      float64          `protobuf:"fixed64,2,opt,name=rate,proto3" json:"rate,omitempty"`
}

func (m *Funding) Reset()      { *m = Funding{} }
func (*Funding) ProtoMessage() {}
func (*Funding) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{9}
}
func (m *Funding) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Funding) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Funding.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Funding) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Funding.Merge(m, src)
}
func (m *Funding) XXX_Size() int {
	return m.Size()
}
func (m *Funding) XXX_DiscardUnknown() {
	xxx_messageInfo_Funding.DiscardUnknown(m)
}

var xxx_messageInfo_Funding proto.InternalMessageInfo

func (m *Funding) GetTimestamp() *types.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Funding) GetRate() float64 {
	if m != nil {
		return m.Rate
	}
	return 0
}

type Stat struct {
	Timestamp *types.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	StatType  StatType         `protobuf:"varint,2,opt,name=stat_type,json=statType,proto3,enum=models.StatType" json:"stat_type,omitempty"`
	Value     float64          `protobuf:"fixed64,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *Stat) Reset()      { *m = Stat{} }
func (*Stat) ProtoMessage() {}
func (*Stat) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{10}
}
func (m *Stat) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Stat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Stat.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Stat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stat.Merge(m, src)
}
func (m *Stat) XXX_Size() int {
	return m.Size()
}
func (m *Stat) XXX_DiscardUnknown() {
	xxx_messageInfo_Stat.DiscardUnknown(m)
}

var xxx_messageInfo_Stat proto.InternalMessageInfo

func (m *Stat) GetTimestamp() *types.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Stat) GetStatType() StatType {
	if m != nil {
		return m.StatType
	}
	return IndexValue
}

func (m *Stat) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func init() {
	proto.RegisterEnum("models.StatType", StatType_name, StatType_value)
	proto.RegisterEnum("models.OrderBookAggregation", OrderBookAggregation_name, OrderBookAggregation_value)
	proto.RegisterType((*OBL1Update)(nil), "models.OBL1Update")
	proto.RegisterType((*OBL1Snapshot)(nil), "models.OBL1Snapshot")
	proto.RegisterType((*OBL2Update)(nil), "models.OBL2Update")
	proto.RegisterType((*OBL2Snapshot)(nil), "models.OBL2Snapshot")
	proto.RegisterType((*OBL3Update)(nil), "models.OBL3Update")
	proto.RegisterType((*OBL3Snapshot)(nil), "models.OBL3Snapshot")
	proto.RegisterType((*Trade)(nil), "models.Trade")
	proto.RegisterType((*AggregatedTrade)(nil), "models.AggregatedTrade")
	proto.RegisterType((*Liquidation)(nil), "models.Liquidation")
	proto.RegisterType((*Funding)(nil), "models.Funding")
	proto.RegisterType((*Stat)(nil), "models.Stat")
}

func init() { proto.RegisterFile("market_data.proto", fileDescriptor_fe080a1ca5f6ba79) }

var fileDescriptor_fe080a1ca5f6ba79 = []byte{
	// 853 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xdc, 0x55, 0x4f, 0x6b, 0xdc, 0x46,
	0x14, 0xd7, 0x68, 0xb5, 0xeb, 0xf5, 0xf3, 0x3f, 0x65, 0x1a, 0x5a, 0x75, 0x5b, 0x94, 0x45, 0x27,
	0x13, 0x13, 0x99, 0xac, 0x4d, 0xc9, 0xd5, 0x9b, 0x50, 0xba, 0xb0, 0x60, 0x23, 0x3b, 0xce, 0xd1,
	0xcc, 0xae, 0xa6, 0xf2, 0xb0, 0x92, 0x46, 0xd1, 0xcc, 0x26, 0x35, 0x94, 0xd2, 0x4b, 0xef, 0x3d,
	0xf5, 0x03, 0x14, 0x0a, 0x39, 0xf4, 0x23, 0xf4, 0x03, 0xe4, 0xe8, 0x63, 0xa0, 0x50, 0xea, 0xf5,
	0xa5, 0xc7, 0x7c, 0x84, 0x32, 0xa3, 0x3f, 0xeb, 0xad, 0x09, 0xd8, 0xa4, 0x94, 0x92, 0x93, 0xe6,
	0xcd, 0xfb, 0xbd, 0x37, 0xbf, 0xf9, 0xbd, 0x79, 0x4f, 0x70, 0x27, 0x21, 0xf9, 0x84, 0xca, 0x93,
	0x90, 0x48, 0xe2, 0x67, 0x39, 0x97, 0x1c, 0xb7, 0x12, 0x1e, 0xd2, 0x58, 0x74, 0x1e, 0x44, 0x4c,
	0x9e, 0x4e, 0x47, 0xfe, 0x98, 0x27, 0xdb, 0x11, 0x8f, 0xf8, 0xb6, 0x76, 0x8f, 0xa6, 0x5f, 0x6b,
	0x4b, 0x1b, 0x7a, 0x55, 0x84, 0x75, 0xee, 0x45, 0x9c, 0x47, 0x31, 0x9d, 0xa3, 0x24, 0x4b, 0xa8,
	0x90, 0x24, 0xc9, 0x4a, 0x80, 0xfb, 0x4f, 0xc0, 0xcb, 0x9c, 0x64, 0x19, 0xcd, 0x45, 0xe9, 0x57,
	0xe7, 0xc5, 0xa4, 0x38, 0x8f, 0xc4, 0xd9, 0x29, 0x91, 0x6c, 0x3c, 0x11, 0xdb, 0x11, 0xcf, 0x43,
	0x9a, 0x8f, 0x38, 0x9f, 0x14, 0x91, 0x25, 0xdc, 0xfb, 0x16, 0x60, 0xbf, 0x3f, 0x7c, 0xf8, 0x34,
	0x0b, 0x89, 0xa4, 0xf8, 0x53, 0x68, 0x8f, 0xa8, 0x90, 0x27, 0x23, 0x16, 0x3a, 0xa8, 0x8b, 0x36,
	0x51, 0xb0, 0xa4, 0xec, 0x3e, 0x0b, 0x6b, 0x17, 0x11, 0x13, 0xc7, 0x9c, 0xbb, 0xf6, 0xc4, 0x04,
	0x3f, 0x82, 0xe5, 0x9a, 0xa5, 0xd3, 0xe8, 0xa2, 0xcd, 0x95, 0x5e, 0xc7, 0x2f, 0x68, 0xfa, 0x15,
	0x4d, 0xff, 0xa8, 0x42, 0x04, 0x73, 0xb0, 0xf7, 0x1d, 0xac, 0xaa, 0xd3, 0x0f, 0x53, 0x92, 0x89,
	0x53, 0x2e, 0xff, 0xf3, 0xf3, 0x7f, 0x42, 0xfa, 0xfa, 0xbd, 0xf2, 0xfa, 0x8f, 0xa0, 0x15, 0xd3,
	0x17, 0x34, 0x16, 0x0e, 0xea, 0x36, 0xca, 0x2c, 0x95, 0x6c, 0xfe, 0xbe, 0x5a, 0xf5, 0x39, 0x9f,
	0x0c, 0x15, 0xa4, 0x6f, 0xbd, 0xfe, 0xe3, 0x9e, 0x11, 0x94, 0xf8, 0x45, 0x0a, 0xe6, 0x2d, 0x28,
	0xe0, 0xbb, 0xd0, 0x94, 0x39, 0x09, 0xa9, 0x63, 0x75, 0xd1, 0x66, 0x3b, 0x28, 0x0c, 0xef, 0x37,
	0x53, 0x2b, 0xd3, 0xab, 0x95, 0xd9, 0x05, 0x6b, 0xc4, 0x42, 0xe1, 0x98, 0x37, 0x24, 0xa6, 0xd1,
	0x2a, 0x8a, 0x88, 0x89, 0x70, 0x1a, 0x37, 0x8d, 0x52, 0xe8, 0xc5, 0xcb, 0x58, 0xb7, 0xb9, 0xcc,
	0x63, 0x58, 0x57, 0xcf, 0xed, 0x24, 0xcb, 0xe9, 0x98, 0x09, 0xc6, 0x53, 0xa7, 0xa9, 0xc3, 0x3f,
	0xbf, 0x16, 0xfe, 0x74, 0x90, 0xca, 0x2f, 0x76, 0x8f, 0x49, 0x3c, 0xa5, 0xc1, 0x9a, 0x8a, 0x39,
	0xa8, 0x42, 0xf0, 0x1e, 0xac, 0xc5, 0x5c, 0x5e, 0xc9, 0xd1, 0xba, 0x41, 0x8e, 0xd5, 0x98, 0xcb,
	0x3a, 0x85, 0xf7, 0xca, 0xd4, 0x75, 0xdd, 0x29, 0xeb, 0xba, 0xb5, 0x20, 0xde, 0x9d, 0xeb, 0x32,
	0x5c, 0xd5, 0x6c, 0x6b, 0x41, 0xb3, 0x77, 0x83, 0x3f, 0x10, 0xa9, 0x7e, 0x46, 0xfa, 0xa5, 0xed,
	0xd4, 0x2f, 0xed, 0x7f, 0x28, 0x96, 0x37, 0x80, 0xe6, 0x91, 0xea, 0x0b, 0xd5, 0x2d, 0x59, 0xce,
	0xc6, 0xb4, 0x9c, 0x0e, 0x85, 0x81, 0x3b, 0xd0, 0x7e, 0x3e, 0x25, 0xa9, 0x64, 0xf2, 0xac, 0x9c,
	0x0d, 0xb5, 0x8d, 0xd7, 0xc1, 0x1c, 0x3c, 0xd1, 0x53, 0xc1, 0x0a, 0xcc, 0xc1, 0x13, 0xef, 0x57,
	0x04, 0x1b, 0x7b, 0x51, 0x94, 0xd3, 0x88, 0x48, 0x1a, 0x16, 0x59, 0x6d, 0x68, 0x54, 0x13, 0xa7,
	0x1d, 0xa8, 0xe5, 0x7b, 0xf4, 0x73, 0x17, 0x56, 0x48, 0x95, 0xbe, 0x3e, 0xf8, 0xea, 0x16, 0xde,
	0x82, 0x96, 0x6e, 0x72, 0xe1, 0x58, 0x5a, 0xb5, 0x35, 0xbf, 0xf8, 0x55, 0xf8, 0x9a, 0x4c, 0x35,
	0x58, 0x0a, 0x88, 0xf7, 0x0b, 0x82, 0x95, 0x21, 0x7b, 0x3e, 0x65, 0x21, 0x91, 0xaa, 0xe2, 0xff,
	0x26, 0x55, 0x07, 0x96, 0x74, 0xb9, 0x6a, 0x9a, 0x95, 0x39, 0x97, 0xd9, 0x7a, 0x97, 0xcc, 0xcd,
	0x45, 0x99, 0xbd, 0x67, 0xb0, 0xf4, 0xe5, 0x34, 0x0d, 0x59, 0x1a, 0x2d, 0x12, 0x42, 0xb7, 0x21,
	0x84, 0xc1, 0xca, 0x89, 0xa4, 0x65, 0x0d, 0xf5, 0xda, 0xfb, 0x01, 0x81, 0x75, 0x28, 0x89, 0x7c,
	0x8f, 0xb4, 0x0f, 0x60, 0x59, 0x48, 0x22, 0x4f, 0xe4, 0x59, 0x56, 0xe4, 0x5e, 0xef, 0xd9, 0x95,
	0xe6, 0x2a, 0xf5, 0xd1, 0x59, 0x46, 0x83, 0xb6, 0x28, 0x57, 0xea, 0xf2, 0x2f, 0x54, 0xa3, 0x68,
	0x51, 0x50, 0x50, 0x18, 0xf7, 0x7f, 0x47, 0xd0, 0xae, 0xc0, 0x78, 0x1d, 0x60, 0x90, 0x86, 0xf4,
	0x1b, 0xdd, 0x50, 0xb6, 0x81, 0x6d, 0x58, 0xdd, 0xcf, 0x68, 0xca, 0xd2, 0xe8, 0x40, 0x29, 0x65,
	0x23, 0xb5, 0xf3, 0x38, 0xe6, 0xa2, 0xde, 0x31, 0xf1, 0x47, 0xb0, 0x71, 0x48, 0xa5, 0x8c, 0x69,
	0x42, 0x53, 0x59, 0x6c, 0x36, 0xf0, 0x67, 0xf0, 0x89, 0xaa, 0x3a, 0x4b, 0xa3, 0x43, 0x2a, 0x54,
	0x3f, 0x7e, 0xc5, 0xa2, 0xd3, 0xc2, 0x69, 0xe1, 0x0e, 0x7c, 0xbc, 0xe8, 0x1c, 0xf2, 0x97, 0x85,
	0xaf, 0x79, 0x3d, 0xf0, 0xf8, 0xd9, 0xde, 0x41, 0xe1, 0x6c, 0xe1, 0x35, 0x58, 0x1e, 0x24, 0x23,
	0x12, 0x93, 0x74, 0x4c, 0xed, 0x25, 0xbc, 0x01, 0x2b, 0xfa, 0x69, 0x1d, 0xf3, 0x78, 0x9a, 0x50,
	0xbb, 0x5d, 0xd1, 0x1d, 0xa4, 0x92, 0xe6, 0x54, 0x48, 0x7b, 0xf9, 0xbe, 0x0f, 0x77, 0xeb, 0x1f,
	0x42, 0xd5, 0x1d, 0xea, 0xb9, 0xb5, 0xc0, 0x1c, 0x3e, 0xb4, 0x0d, 0xfd, 0xed, 0xd9, 0x48, 0x7f,
	0x77, 0x6c, 0xb3, 0xbf, 0x7b, 0x7e, 0xe1, 0x1a, 0x6f, 0x2e, 0x5c, 0xe3, 0xed, 0x85, 0x8b, 0xbe,
	0x9f, 0xb9, 0xe8, 0xd5, 0xcc, 0x45, 0xaf, 0x67, 0x2e, 0x3a, 0x9f, 0xb9, 0xe8, 0xcf, 0x99, 0x8b,
	0xfe, 0x9a, 0xb9, 0xc6, 0xdb, 0x99, 0x8b, 0x7e, 0xbc, 0x74, 0x8d, 0xf3, 0x4b, 0xd7, 0x78, 0x73,
	0xe9, 0x1a, 0xa3, 0x96, 0xae, 0xd3, 0xce, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xc7, 0x3e, 0x7b,
	0xef, 0x2f, 0x09, 0x00, 0x00,
}

func (x StatType) String() string {
	s, ok := StatType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x OrderBookAggregation) String() string {
	s, ok := OrderBookAggregation_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (this *OBL1Update) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OBL1Update)
	if !ok {
		that2, ok := that.(OBL1Update)
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
	if this.BestBid != that1.BestBid {
		return false
	}
	if this.BestAsk != that1.BestAsk {
		return false
	}
	if !this.Timestamp.Equal(that1.Timestamp) {
		return false
	}
	return true
}
func (this *OBL1Snapshot) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OBL1Snapshot)
	if !ok {
		that2, ok := that.(OBL1Snapshot)
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
	if this.BestBid != that1.BestBid {
		return false
	}
	if this.BestAsk != that1.BestAsk {
		return false
	}
	if !this.Timestamp.Equal(that1.Timestamp) {
		return false
	}
	return true
}
func (this *OBL2Update) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OBL2Update)
	if !ok {
		that2, ok := that.(OBL2Update)
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
	if len(this.Levels) != len(that1.Levels) {
		return false
	}
	for i := range this.Levels {
		if !this.Levels[i].Equal(&that1.Levels[i]) {
			return false
		}
	}
	if !this.Timestamp.Equal(that1.Timestamp) {
		return false
	}
	if this.Trade != that1.Trade {
		return false
	}
	return true
}
func (this *OBL2Snapshot) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OBL2Snapshot)
	if !ok {
		that2, ok := that.(OBL2Snapshot)
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
	if len(this.Bids) != len(that1.Bids) {
		return false
	}
	for i := range this.Bids {
		if !this.Bids[i].Equal(&that1.Bids[i]) {
			return false
		}
	}
	if len(this.Asks) != len(that1.Asks) {
		return false
	}
	for i := range this.Asks {
		if !this.Asks[i].Equal(&that1.Asks[i]) {
			return false
		}
	}
	if !this.Timestamp.Equal(that1.Timestamp) {
		return false
	}
	if !this.TickPrecision.Equal(that1.TickPrecision) {
		return false
	}
	if !this.LotPrecision.Equal(that1.LotPrecision) {
		return false
	}
	return true
}
func (this *OBL3Update) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OBL3Update)
	if !ok {
		that2, ok := that.(OBL3Update)
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
	if len(this.Bids) != len(that1.Bids) {
		return false
	}
	for i := range this.Bids {
		if !this.Bids[i].Equal(&that1.Bids[i]) {
			return false
		}
	}
	if len(this.Asks) != len(that1.Asks) {
		return false
	}
	for i := range this.Asks {
		if !this.Asks[i].Equal(&that1.Asks[i]) {
			return false
		}
	}
	if !this.Timestamp.Equal(that1.Timestamp) {
		return false
	}
	if !this.TickPrecision.Equal(that1.TickPrecision) {
		return false
	}
	if !this.LotPrecision.Equal(that1.LotPrecision) {
		return false
	}
	return true
}
func (this *OBL3Snapshot) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OBL3Snapshot)
	if !ok {
		that2, ok := that.(OBL3Snapshot)
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
	if len(this.Bids) != len(that1.Bids) {
		return false
	}
	for i := range this.Bids {
		if !this.Bids[i].Equal(&that1.Bids[i]) {
			return false
		}
	}
	if len(this.Asks) != len(that1.Asks) {
		return false
	}
	for i := range this.Asks {
		if !this.Asks[i].Equal(&that1.Asks[i]) {
			return false
		}
	}
	if !this.Timestamp.Equal(that1.Timestamp) {
		return false
	}
	return true
}
func (this *Trade) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Trade)
	if !ok {
		that2, ok := that.(Trade)
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
	if this.Price != that1.Price {
		return false
	}
	if this.Quantity != that1.Quantity {
		return false
	}
	if this.ID != that1.ID {
		return false
	}
	return true
}
func (this *AggregatedTrade) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AggregatedTrade)
	if !ok {
		that2, ok := that.(AggregatedTrade)
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
	if this.Bid != that1.Bid {
		return false
	}
	if !this.Timestamp.Equal(that1.Timestamp) {
		return false
	}
	if this.AggregateID != that1.AggregateID {
		return false
	}
	if len(this.Trades) != len(that1.Trades) {
		return false
	}
	for i := range this.Trades {
		if !this.Trades[i].Equal(&that1.Trades[i]) {
			return false
		}
	}
	return true
}
func (this *Liquidation) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Liquidation)
	if !ok {
		that2, ok := that.(Liquidation)
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
	if this.Bid != that1.Bid {
		return false
	}
	if !this.Timestamp.Equal(that1.Timestamp) {
		return false
	}
	if this.OrderID != that1.OrderID {
		return false
	}
	if this.Price != that1.Price {
		return false
	}
	if this.Quantity != that1.Quantity {
		return false
	}
	return true
}
func (this *Funding) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Funding)
	if !ok {
		that2, ok := that.(Funding)
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
	if !this.Timestamp.Equal(that1.Timestamp) {
		return false
	}
	if this.Rate != that1.Rate {
		return false
	}
	return true
}
func (this *Stat) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Stat)
	if !ok {
		that2, ok := that.(Stat)
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
	if !this.Timestamp.Equal(that1.Timestamp) {
		return false
	}
	if this.StatType != that1.StatType {
		return false
	}
	if this.Value != that1.Value {
		return false
	}
	return true
}
func (this *OBL1Update) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&models.OBL1Update{")
	s = append(s, "BestBid: "+fmt.Sprintf("%#v", this.BestBid)+",\n")
	s = append(s, "BestAsk: "+fmt.Sprintf("%#v", this.BestAsk)+",\n")
	if this.Timestamp != nil {
		s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OBL1Snapshot) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&models.OBL1Snapshot{")
	s = append(s, "BestBid: "+fmt.Sprintf("%#v", this.BestBid)+",\n")
	s = append(s, "BestAsk: "+fmt.Sprintf("%#v", this.BestAsk)+",\n")
	if this.Timestamp != nil {
		s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OBL2Update) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&models.OBL2Update{")
	if this.Levels != nil {
		vs := make([]*gorderbook.OrderBookLevel, len(this.Levels))
		for i := range vs {
			vs[i] = &this.Levels[i]
		}
		s = append(s, "Levels: "+fmt.Sprintf("%#v", vs)+",\n")
	}
	if this.Timestamp != nil {
		s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	}
	s = append(s, "Trade: "+fmt.Sprintf("%#v", this.Trade)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OBL2Snapshot) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&models.OBL2Snapshot{")
	if this.Bids != nil {
		vs := make([]*gorderbook.OrderBookLevel, len(this.Bids))
		for i := range vs {
			vs[i] = &this.Bids[i]
		}
		s = append(s, "Bids: "+fmt.Sprintf("%#v", vs)+",\n")
	}
	if this.Asks != nil {
		vs := make([]*gorderbook.OrderBookLevel, len(this.Asks))
		for i := range vs {
			vs[i] = &this.Asks[i]
		}
		s = append(s, "Asks: "+fmt.Sprintf("%#v", vs)+",\n")
	}
	if this.Timestamp != nil {
		s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	}
	if this.TickPrecision != nil {
		s = append(s, "TickPrecision: "+fmt.Sprintf("%#v", this.TickPrecision)+",\n")
	}
	if this.LotPrecision != nil {
		s = append(s, "LotPrecision: "+fmt.Sprintf("%#v", this.LotPrecision)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OBL3Update) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&models.OBL3Update{")
	if this.Bids != nil {
		vs := make([]*gorderbook.Order, len(this.Bids))
		for i := range vs {
			vs[i] = &this.Bids[i]
		}
		s = append(s, "Bids: "+fmt.Sprintf("%#v", vs)+",\n")
	}
	if this.Asks != nil {
		vs := make([]*gorderbook.Order, len(this.Asks))
		for i := range vs {
			vs[i] = &this.Asks[i]
		}
		s = append(s, "Asks: "+fmt.Sprintf("%#v", vs)+",\n")
	}
	if this.Timestamp != nil {
		s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	}
	if this.TickPrecision != nil {
		s = append(s, "TickPrecision: "+fmt.Sprintf("%#v", this.TickPrecision)+",\n")
	}
	if this.LotPrecision != nil {
		s = append(s, "LotPrecision: "+fmt.Sprintf("%#v", this.LotPrecision)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OBL3Snapshot) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&models.OBL3Snapshot{")
	if this.Bids != nil {
		vs := make([]*gorderbook.Order, len(this.Bids))
		for i := range vs {
			vs[i] = &this.Bids[i]
		}
		s = append(s, "Bids: "+fmt.Sprintf("%#v", vs)+",\n")
	}
	if this.Asks != nil {
		vs := make([]*gorderbook.Order, len(this.Asks))
		for i := range vs {
			vs[i] = &this.Asks[i]
		}
		s = append(s, "Asks: "+fmt.Sprintf("%#v", vs)+",\n")
	}
	if this.Timestamp != nil {
		s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Trade) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&models.Trade{")
	s = append(s, "Price: "+fmt.Sprintf("%#v", this.Price)+",\n")
	s = append(s, "Quantity: "+fmt.Sprintf("%#v", this.Quantity)+",\n")
	s = append(s, "ID: "+fmt.Sprintf("%#v", this.ID)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *AggregatedTrade) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&models.AggregatedTrade{")
	s = append(s, "Bid: "+fmt.Sprintf("%#v", this.Bid)+",\n")
	if this.Timestamp != nil {
		s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	}
	s = append(s, "AggregateID: "+fmt.Sprintf("%#v", this.AggregateID)+",\n")
	if this.Trades != nil {
		vs := make([]*Trade, len(this.Trades))
		for i := range vs {
			vs[i] = &this.Trades[i]
		}
		s = append(s, "Trades: "+fmt.Sprintf("%#v", vs)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Liquidation) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&models.Liquidation{")
	s = append(s, "Bid: "+fmt.Sprintf("%#v", this.Bid)+",\n")
	if this.Timestamp != nil {
		s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	}
	s = append(s, "OrderID: "+fmt.Sprintf("%#v", this.OrderID)+",\n")
	s = append(s, "Price: "+fmt.Sprintf("%#v", this.Price)+",\n")
	s = append(s, "Quantity: "+fmt.Sprintf("%#v", this.Quantity)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Funding) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&models.Funding{")
	if this.Timestamp != nil {
		s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	}
	s = append(s, "Rate: "+fmt.Sprintf("%#v", this.Rate)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Stat) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&models.Stat{")
	if this.Timestamp != nil {
		s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	}
	s = append(s, "StatType: "+fmt.Sprintf("%#v", this.StatType)+",\n")
	s = append(s, "Value: "+fmt.Sprintf("%#v", this.Value)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringMarketData(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *OBL1Update) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OBL1Update) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.BestBid != 0 {
		dAtA[i] = 0x9
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.BestBid))))
		i += 8
	}
	if m.BestAsk != 0 {
		dAtA[i] = 0x11
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.BestAsk))))
		i += 8
	}
	if m.Timestamp != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.Timestamp.Size()))
		n1, err := m.Timestamp.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	return i, nil
}

func (m *OBL1Snapshot) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OBL1Snapshot) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.BestBid != 0 {
		dAtA[i] = 0x9
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.BestBid))))
		i += 8
	}
	if m.BestAsk != 0 {
		dAtA[i] = 0x11
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.BestAsk))))
		i += 8
	}
	if m.Timestamp != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.Timestamp.Size()))
		n2, err := m.Timestamp.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func (m *OBL2Update) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OBL2Update) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Levels) > 0 {
		for _, msg := range m.Levels {
			dAtA[i] = 0xa
			i++
			i = encodeVarintMarketData(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Timestamp != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.Timestamp.Size()))
		n3, err := m.Timestamp.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	if m.Trade {
		dAtA[i] = 0x20
		i++
		if m.Trade {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	return i, nil
}

func (m *OBL2Snapshot) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OBL2Snapshot) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Bids) > 0 {
		for _, msg := range m.Bids {
			dAtA[i] = 0x12
			i++
			i = encodeVarintMarketData(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Asks) > 0 {
		for _, msg := range m.Asks {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintMarketData(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Timestamp != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.Timestamp.Size()))
		n4, err := m.Timestamp.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n4
	}
	if m.TickPrecision != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.TickPrecision.Size()))
		n5, err := m.TickPrecision.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n5
	}
	if m.LotPrecision != nil {
		dAtA[i] = 0x32
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.LotPrecision.Size()))
		n6, err := m.LotPrecision.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n6
	}
	return i, nil
}

func (m *OBL3Update) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OBL3Update) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Bids) > 0 {
		for _, msg := range m.Bids {
			dAtA[i] = 0x12
			i++
			i = encodeVarintMarketData(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Asks) > 0 {
		for _, msg := range m.Asks {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintMarketData(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Timestamp != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.Timestamp.Size()))
		n7, err := m.Timestamp.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n7
	}
	if m.TickPrecision != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.TickPrecision.Size()))
		n8, err := m.TickPrecision.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n8
	}
	if m.LotPrecision != nil {
		dAtA[i] = 0x32
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.LotPrecision.Size()))
		n9, err := m.LotPrecision.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n9
	}
	return i, nil
}

func (m *OBL3Snapshot) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OBL3Snapshot) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Bids) > 0 {
		for _, msg := range m.Bids {
			dAtA[i] = 0x12
			i++
			i = encodeVarintMarketData(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Asks) > 0 {
		for _, msg := range m.Asks {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintMarketData(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Timestamp != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.Timestamp.Size()))
		n10, err := m.Timestamp.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n10
	}
	return i, nil
}

func (m *Trade) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Trade) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Price != 0 {
		dAtA[i] = 0x9
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Price))))
		i += 8
	}
	if m.Quantity != 0 {
		dAtA[i] = 0x11
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Quantity))))
		i += 8
	}
	if m.ID != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.ID))
	}
	return i, nil
}

func (m *AggregatedTrade) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AggregatedTrade) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Bid {
		dAtA[i] = 0x8
		i++
		if m.Bid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.Timestamp != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.Timestamp.Size()))
		n11, err := m.Timestamp.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n11
	}
	if m.AggregateID != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.AggregateID))
	}
	if len(m.Trades) > 0 {
		for _, msg := range m.Trades {
			dAtA[i] = 0x22
			i++
			i = encodeVarintMarketData(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *Liquidation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Liquidation) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Bid {
		dAtA[i] = 0x8
		i++
		if m.Bid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.Timestamp != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.Timestamp.Size()))
		n12, err := m.Timestamp.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n12
	}
	if m.OrderID != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.OrderID))
	}
	if m.Price != 0 {
		dAtA[i] = 0x21
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Price))))
		i += 8
	}
	if m.Quantity != 0 {
		dAtA[i] = 0x29
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Quantity))))
		i += 8
	}
	return i, nil
}

func (m *Funding) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Funding) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Timestamp != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.Timestamp.Size()))
		n13, err := m.Timestamp.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n13
	}
	if m.Rate != 0 {
		dAtA[i] = 0x11
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Rate))))
		i += 8
	}
	return i, nil
}

func (m *Stat) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Stat) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Timestamp != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.Timestamp.Size()))
		n14, err := m.Timestamp.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n14
	}
	if m.StatType != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintMarketData(dAtA, i, uint64(m.StatType))
	}
	if m.Value != 0 {
		dAtA[i] = 0x19
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Value))))
		i += 8
	}
	return i, nil
}

func encodeVarintMarketData(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *OBL1Update) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BestBid != 0 {
		n += 9
	}
	if m.BestAsk != 0 {
		n += 9
	}
	if m.Timestamp != nil {
		l = m.Timestamp.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	return n
}

func (m *OBL1Snapshot) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BestBid != 0 {
		n += 9
	}
	if m.BestAsk != 0 {
		n += 9
	}
	if m.Timestamp != nil {
		l = m.Timestamp.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	return n
}

func (m *OBL2Update) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Levels) > 0 {
		for _, e := range m.Levels {
			l = e.Size()
			n += 1 + l + sovMarketData(uint64(l))
		}
	}
	if m.Timestamp != nil {
		l = m.Timestamp.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.Trade {
		n += 2
	}
	return n
}

func (m *OBL2Snapshot) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Bids) > 0 {
		for _, e := range m.Bids {
			l = e.Size()
			n += 1 + l + sovMarketData(uint64(l))
		}
	}
	if len(m.Asks) > 0 {
		for _, e := range m.Asks {
			l = e.Size()
			n += 1 + l + sovMarketData(uint64(l))
		}
	}
	if m.Timestamp != nil {
		l = m.Timestamp.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.TickPrecision != nil {
		l = m.TickPrecision.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.LotPrecision != nil {
		l = m.LotPrecision.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	return n
}

func (m *OBL3Update) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Bids) > 0 {
		for _, e := range m.Bids {
			l = e.Size()
			n += 1 + l + sovMarketData(uint64(l))
		}
	}
	if len(m.Asks) > 0 {
		for _, e := range m.Asks {
			l = e.Size()
			n += 1 + l + sovMarketData(uint64(l))
		}
	}
	if m.Timestamp != nil {
		l = m.Timestamp.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.TickPrecision != nil {
		l = m.TickPrecision.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.LotPrecision != nil {
		l = m.LotPrecision.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	return n
}

func (m *OBL3Snapshot) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Bids) > 0 {
		for _, e := range m.Bids {
			l = e.Size()
			n += 1 + l + sovMarketData(uint64(l))
		}
	}
	if len(m.Asks) > 0 {
		for _, e := range m.Asks {
			l = e.Size()
			n += 1 + l + sovMarketData(uint64(l))
		}
	}
	if m.Timestamp != nil {
		l = m.Timestamp.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	return n
}

func (m *Trade) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Price != 0 {
		n += 9
	}
	if m.Quantity != 0 {
		n += 9
	}
	if m.ID != 0 {
		n += 1 + sovMarketData(uint64(m.ID))
	}
	return n
}

func (m *AggregatedTrade) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Bid {
		n += 2
	}
	if m.Timestamp != nil {
		l = m.Timestamp.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.AggregateID != 0 {
		n += 1 + sovMarketData(uint64(m.AggregateID))
	}
	if len(m.Trades) > 0 {
		for _, e := range m.Trades {
			l = e.Size()
			n += 1 + l + sovMarketData(uint64(l))
		}
	}
	return n
}

func (m *Liquidation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Bid {
		n += 2
	}
	if m.Timestamp != nil {
		l = m.Timestamp.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.OrderID != 0 {
		n += 1 + sovMarketData(uint64(m.OrderID))
	}
	if m.Price != 0 {
		n += 9
	}
	if m.Quantity != 0 {
		n += 9
	}
	return n
}

func (m *Funding) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Timestamp != nil {
		l = m.Timestamp.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.Rate != 0 {
		n += 9
	}
	return n
}

func (m *Stat) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Timestamp != nil {
		l = m.Timestamp.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.StatType != 0 {
		n += 1 + sovMarketData(uint64(m.StatType))
	}
	if m.Value != 0 {
		n += 9
	}
	return n
}

func sovMarketData(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozMarketData(x uint64) (n int) {
	return sovMarketData(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *OBL1Update) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OBL1Update{`,
		`BestBid:` + fmt.Sprintf("%v", this.BestBid) + `,`,
		`BestAsk:` + fmt.Sprintf("%v", this.BestAsk) + `,`,
		`Timestamp:` + strings.Replace(fmt.Sprintf("%v", this.Timestamp), "Timestamp", "types.Timestamp", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OBL1Snapshot) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OBL1Snapshot{`,
		`BestBid:` + fmt.Sprintf("%v", this.BestBid) + `,`,
		`BestAsk:` + fmt.Sprintf("%v", this.BestAsk) + `,`,
		`Timestamp:` + strings.Replace(fmt.Sprintf("%v", this.Timestamp), "Timestamp", "types.Timestamp", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OBL2Update) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OBL2Update{`,
		`Levels:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Levels), "OrderBookLevel", "gorderbook.OrderBookLevel", 1), `&`, ``, 1) + `,`,
		`Timestamp:` + strings.Replace(fmt.Sprintf("%v", this.Timestamp), "Timestamp", "types.Timestamp", 1) + `,`,
		`Trade:` + fmt.Sprintf("%v", this.Trade) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OBL2Snapshot) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OBL2Snapshot{`,
		`Bids:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Bids), "OrderBookLevel", "gorderbook.OrderBookLevel", 1), `&`, ``, 1) + `,`,
		`Asks:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Asks), "OrderBookLevel", "gorderbook.OrderBookLevel", 1), `&`, ``, 1) + `,`,
		`Timestamp:` + strings.Replace(fmt.Sprintf("%v", this.Timestamp), "Timestamp", "types.Timestamp", 1) + `,`,
		`TickPrecision:` + strings.Replace(fmt.Sprintf("%v", this.TickPrecision), "UInt64Value", "types.UInt64Value", 1) + `,`,
		`LotPrecision:` + strings.Replace(fmt.Sprintf("%v", this.LotPrecision), "UInt64Value", "types.UInt64Value", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OBL3Update) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OBL3Update{`,
		`Bids:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Bids), "Order", "gorderbook.Order", 1), `&`, ``, 1) + `,`,
		`Asks:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Asks), "Order", "gorderbook.Order", 1), `&`, ``, 1) + `,`,
		`Timestamp:` + strings.Replace(fmt.Sprintf("%v", this.Timestamp), "Timestamp", "types.Timestamp", 1) + `,`,
		`TickPrecision:` + strings.Replace(fmt.Sprintf("%v", this.TickPrecision), "UInt64Value", "types.UInt64Value", 1) + `,`,
		`LotPrecision:` + strings.Replace(fmt.Sprintf("%v", this.LotPrecision), "UInt64Value", "types.UInt64Value", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OBL3Snapshot) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OBL3Snapshot{`,
		`Bids:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Bids), "Order", "gorderbook.Order", 1), `&`, ``, 1) + `,`,
		`Asks:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Asks), "Order", "gorderbook.Order", 1), `&`, ``, 1) + `,`,
		`Timestamp:` + strings.Replace(fmt.Sprintf("%v", this.Timestamp), "Timestamp", "types.Timestamp", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Trade) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Trade{`,
		`Price:` + fmt.Sprintf("%v", this.Price) + `,`,
		`Quantity:` + fmt.Sprintf("%v", this.Quantity) + `,`,
		`ID:` + fmt.Sprintf("%v", this.ID) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AggregatedTrade) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AggregatedTrade{`,
		`Bid:` + fmt.Sprintf("%v", this.Bid) + `,`,
		`Timestamp:` + strings.Replace(fmt.Sprintf("%v", this.Timestamp), "Timestamp", "types.Timestamp", 1) + `,`,
		`AggregateID:` + fmt.Sprintf("%v", this.AggregateID) + `,`,
		`Trades:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Trades), "Trade", "Trade", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Liquidation) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Liquidation{`,
		`Bid:` + fmt.Sprintf("%v", this.Bid) + `,`,
		`Timestamp:` + strings.Replace(fmt.Sprintf("%v", this.Timestamp), "Timestamp", "types.Timestamp", 1) + `,`,
		`OrderID:` + fmt.Sprintf("%v", this.OrderID) + `,`,
		`Price:` + fmt.Sprintf("%v", this.Price) + `,`,
		`Quantity:` + fmt.Sprintf("%v", this.Quantity) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Funding) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Funding{`,
		`Timestamp:` + strings.Replace(fmt.Sprintf("%v", this.Timestamp), "Timestamp", "types.Timestamp", 1) + `,`,
		`Rate:` + fmt.Sprintf("%v", this.Rate) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Stat) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Stat{`,
		`Timestamp:` + strings.Replace(fmt.Sprintf("%v", this.Timestamp), "Timestamp", "types.Timestamp", 1) + `,`,
		`StatType:` + fmt.Sprintf("%v", this.StatType) + `,`,
		`Value:` + fmt.Sprintf("%v", this.Value) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringMarketData(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *OBL1Update) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarketData
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
			return fmt.Errorf("proto: OBL1Update: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OBL1Update: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field BestBid", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.BestBid = float64(math.Float64frombits(v))
		case 2:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field BestAsk", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.BestAsk = float64(math.Float64frombits(v))
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Timestamp == nil {
				m.Timestamp = &types.Timestamp{}
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMarketData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMarketData
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMarketData
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
func (m *OBL1Snapshot) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarketData
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
			return fmt.Errorf("proto: OBL1Snapshot: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OBL1Snapshot: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field BestBid", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.BestBid = float64(math.Float64frombits(v))
		case 2:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field BestAsk", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.BestAsk = float64(math.Float64frombits(v))
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Timestamp == nil {
				m.Timestamp = &types.Timestamp{}
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMarketData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMarketData
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMarketData
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
func (m *OBL2Update) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarketData
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
			return fmt.Errorf("proto: OBL2Update: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OBL2Update: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Levels", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Levels = append(m.Levels, gorderbook.OrderBookLevel{})
			if err := m.Levels[len(m.Levels)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Timestamp == nil {
				m.Timestamp = &types.Timestamp{}
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Trade", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
			m.Trade = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipMarketData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMarketData
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMarketData
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
func (m *OBL2Snapshot) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarketData
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
			return fmt.Errorf("proto: OBL2Snapshot: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OBL2Snapshot: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bids", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Bids = append(m.Bids, gorderbook.OrderBookLevel{})
			if err := m.Bids[len(m.Bids)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Asks = append(m.Asks, gorderbook.OrderBookLevel{})
			if err := m.Asks[len(m.Asks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Timestamp == nil {
				m.Timestamp = &types.Timestamp{}
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TickPrecision", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.TickPrecision == nil {
				m.TickPrecision = &types.UInt64Value{}
			}
			if err := m.TickPrecision.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LotPrecision", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.LotPrecision == nil {
				m.LotPrecision = &types.UInt64Value{}
			}
			if err := m.LotPrecision.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMarketData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMarketData
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMarketData
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
func (m *OBL3Update) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarketData
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
			return fmt.Errorf("proto: OBL3Update: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OBL3Update: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bids", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Bids = append(m.Bids, gorderbook.Order{})
			if err := m.Bids[len(m.Bids)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Asks = append(m.Asks, gorderbook.Order{})
			if err := m.Asks[len(m.Asks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Timestamp == nil {
				m.Timestamp = &types.Timestamp{}
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TickPrecision", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.TickPrecision == nil {
				m.TickPrecision = &types.UInt64Value{}
			}
			if err := m.TickPrecision.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LotPrecision", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.LotPrecision == nil {
				m.LotPrecision = &types.UInt64Value{}
			}
			if err := m.LotPrecision.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMarketData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMarketData
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMarketData
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
func (m *OBL3Snapshot) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarketData
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
			return fmt.Errorf("proto: OBL3Snapshot: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OBL3Snapshot: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bids", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Bids = append(m.Bids, gorderbook.Order{})
			if err := m.Bids[len(m.Bids)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Asks = append(m.Asks, gorderbook.Order{})
			if err := m.Asks[len(m.Asks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Timestamp == nil {
				m.Timestamp = &types.Timestamp{}
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMarketData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMarketData
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMarketData
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
func (m *Trade) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarketData
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
			return fmt.Errorf("proto: Trade: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Trade: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Price = float64(math.Float64frombits(v))
		case 2:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Quantity", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Quantity = float64(math.Float64frombits(v))
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMarketData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMarketData
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMarketData
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
func (m *AggregatedTrade) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarketData
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
			return fmt.Errorf("proto: AggregatedTrade: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AggregatedTrade: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
			m.Bid = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Timestamp == nil {
				m.Timestamp = &types.Timestamp{}
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AggregateID", wireType)
			}
			m.AggregateID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AggregateID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Trades", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Trades = append(m.Trades, Trade{})
			if err := m.Trades[len(m.Trades)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMarketData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMarketData
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMarketData
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
func (m *Liquidation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarketData
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
			return fmt.Errorf("proto: Liquidation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Liquidation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bid", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
			m.Bid = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Timestamp == nil {
				m.Timestamp = &types.Timestamp{}
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderID", wireType)
			}
			m.OrderID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrderID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Price = float64(math.Float64frombits(v))
		case 5:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Quantity", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Quantity = float64(math.Float64frombits(v))
		default:
			iNdEx = preIndex
			skippy, err := skipMarketData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMarketData
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMarketData
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
func (m *Funding) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarketData
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
			return fmt.Errorf("proto: Funding: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Funding: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Timestamp == nil {
				m.Timestamp = &types.Timestamp{}
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rate", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Rate = float64(math.Float64frombits(v))
		default:
			iNdEx = preIndex
			skippy, err := skipMarketData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMarketData
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMarketData
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
func (m *Stat) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarketData
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
			return fmt.Errorf("proto: Stat: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Stat: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
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
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Timestamp == nil {
				m.Timestamp = &types.Timestamp{}
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StatType", wireType)
			}
			m.StatType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StatType |= StatType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Value = float64(math.Float64frombits(v))
		default:
			iNdEx = preIndex
			skippy, err := skipMarketData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMarketData
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMarketData
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
func skipMarketData(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMarketData
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
					return 0, ErrIntOverflowMarketData
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
					return 0, ErrIntOverflowMarketData
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
				return 0, ErrInvalidLengthMarketData
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthMarketData
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowMarketData
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
				next, err := skipMarketData(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthMarketData
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
	ErrInvalidLengthMarketData = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMarketData   = fmt.Errorf("proto: integer overflow")
)
