package models

import (
	bytes "bytes"
	encoding_binary "encoding/binary"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
	gorderbook "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	io "io"
	math "math"
	math_bits "math/bits"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

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
	FundingRate             StatType = 10
)

var StatType_name = map[int32]string{
	0:  "IndexValue",
	1:  "OpeningPrice",
	2:  "ClosingPrice",
	3:  "SettlementPrice",
	4:  "TradingSessionHighPrice",
	5:  "TradingSessionLowPrice",
	6:  "TradingSessionVWAPPrice",
	7:  "Imbalance",
	8:  "TradeVolume",
	9:  "OpenInterest",
	10: "FundingRate",
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
	"FundingRate":             10,
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
		n, err := m.MarshalToSizedBuffer(b)
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
		n, err := m.MarshalToSizedBuffer(b)
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
		n, err := m.MarshalToSizedBuffer(b)
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
		n, err := m.MarshalToSizedBuffer(b)
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
		n, err := m.MarshalToSizedBuffer(b)
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
		n, err := m.MarshalToSizedBuffer(b)
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

type UPV3Snapshot struct {
	Ticks                   []*gorderbook.UPV3Tick     `protobuf:"bytes,1,rep,name=ticks,proto3" json:"ticks,omitempty"`
	Positions               []*gorderbook.UPV3Position `protobuf:"bytes,2,rep,name=positions,proto3" json:"positions,omitempty"`
	Liquidity               []byte                     `protobuf:"bytes,3,opt,name=liquidity,proto3" json:"liquidity,omitempty"`
	SqrtPrice               []byte                     `protobuf:"bytes,4,opt,name=sqrt_price,json=sqrtPrice,proto3" json:"sqrt_price,omitempty"`
	FeeGrowthGlobal_0X128   []byte                     `protobuf:"bytes,5,opt,name=fee_growth_global_0x128,json=feeGrowthGlobal0x128,proto3" json:"fee_growth_global_0x128,omitempty"`
	FeeGrowthGlobal_1X128   []byte                     `protobuf:"bytes,6,opt,name=fee_growth_global_1x128,json=feeGrowthGlobal1x128,proto3" json:"fee_growth_global_1x128,omitempty"`
	ProtocolFees_0          []byte                     `protobuf:"bytes,7,opt,name=protocol_fees_0,json=protocolFees0,proto3" json:"protocol_fees_0,omitempty"`
	ProtocolFees_1          []byte                     `protobuf:"bytes,8,opt,name=protocol_fees_1,json=protocolFees1,proto3" json:"protocol_fees_1,omitempty"`
	TotalValueLockedToken_0 []byte                     `protobuf:"bytes,9,opt,name=total_value_locked_token_0,json=totalValueLockedToken0,proto3" json:"total_value_locked_token_0,omitempty"`
	TotalValueLockedToken_1 []byte                     `protobuf:"bytes,10,opt,name=total_value_locked_token_1,json=totalValueLockedToken1,proto3" json:"total_value_locked_token_1,omitempty"`
	Tick                    int32                      `protobuf:"varint,11,opt,name=tick,proto3" json:"tick,omitempty"`
	FeeTier                 int32                      `protobuf:"varint,12,opt,name=fee_tier,json=feeTier,proto3" json:"fee_tier,omitempty"`
	Timestamp               *types.Timestamp           `protobuf:"bytes,13,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (m *UPV3Snapshot) Reset()      { *m = UPV3Snapshot{} }
func (*UPV3Snapshot) ProtoMessage() {}
func (*UPV3Snapshot) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{6}
}
func (m *UPV3Snapshot) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UPV3Snapshot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UPV3Snapshot.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UPV3Snapshot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UPV3Snapshot.Merge(m, src)
}
func (m *UPV3Snapshot) XXX_Size() int {
	return m.Size()
}
func (m *UPV3Snapshot) XXX_DiscardUnknown() {
	xxx_messageInfo_UPV3Snapshot.DiscardUnknown(m)
}

var xxx_messageInfo_UPV3Snapshot proto.InternalMessageInfo

func (m *UPV3Snapshot) GetTicks() []*gorderbook.UPV3Tick {
	if m != nil {
		return m.Ticks
	}
	return nil
}

func (m *UPV3Snapshot) GetPositions() []*gorderbook.UPV3Position {
	if m != nil {
		return m.Positions
	}
	return nil
}

func (m *UPV3Snapshot) GetLiquidity() []byte {
	if m != nil {
		return m.Liquidity
	}
	return nil
}

func (m *UPV3Snapshot) GetSqrtPrice() []byte {
	if m != nil {
		return m.SqrtPrice
	}
	return nil
}

func (m *UPV3Snapshot) GetFeeGrowthGlobal_0X128() []byte {
	if m != nil {
		return m.FeeGrowthGlobal_0X128
	}
	return nil
}

func (m *UPV3Snapshot) GetFeeGrowthGlobal_1X128() []byte {
	if m != nil {
		return m.FeeGrowthGlobal_1X128
	}
	return nil
}

func (m *UPV3Snapshot) GetProtocolFees_0() []byte {
	if m != nil {
		return m.ProtocolFees_0
	}
	return nil
}

func (m *UPV3Snapshot) GetProtocolFees_1() []byte {
	if m != nil {
		return m.ProtocolFees_1
	}
	return nil
}

func (m *UPV3Snapshot) GetTotalValueLockedToken_0() []byte {
	if m != nil {
		return m.TotalValueLockedToken_0
	}
	return nil
}

func (m *UPV3Snapshot) GetTotalValueLockedToken_1() []byte {
	if m != nil {
		return m.TotalValueLockedToken_1
	}
	return nil
}

func (m *UPV3Snapshot) GetTick() int32 {
	if m != nil {
		return m.Tick
	}
	return 0
}

func (m *UPV3Snapshot) GetFeeTier() int32 {
	if m != nil {
		return m.FeeTier
	}
	return 0
}

func (m *UPV3Snapshot) GetTimestamp() *types.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

type UPV3Update struct {
	Initialize      *gorderbook.UPV3Initialize      `protobuf:"bytes,1,opt,name=initialize,proto3" json:"initialize,omitempty"`
	Mint            *gorderbook.UPV3Mint            `protobuf:"bytes,2,opt,name=mint,proto3" json:"mint,omitempty"`
	Burn            *gorderbook.UPV3Burn            `protobuf:"bytes,3,opt,name=burn,proto3" json:"burn,omitempty"`
	Swap            *gorderbook.UPV3Swap            `protobuf:"bytes,4,opt,name=swap,proto3" json:"swap,omitempty"`
	Collect         *gorderbook.UPV3Collect         `protobuf:"bytes,5,opt,name=collect,proto3" json:"collect,omitempty"`
	Flash           *gorderbook.UPV3Flash           `protobuf:"bytes,6,opt,name=flash,proto3" json:"flash,omitempty"`
	SetFeeProtocol  *gorderbook.UPV3SetFeeProtocol  `protobuf:"bytes,7,opt,name=set_fee_protocol,json=setFeeProtocol,proto3" json:"set_fee_protocol,omitempty"`
	CollectProtocol *gorderbook.UPV3CollectProtocol `protobuf:"bytes,8,opt,name=collect_protocol,json=collectProtocol,proto3" json:"collect_protocol,omitempty"`
	Removed         bool                            `protobuf:"varint,9,opt,name=removed,proto3" json:"removed,omitempty"`
	Block           uint64                          `protobuf:"varint,10,opt,name=block,proto3" json:"block,omitempty"`
}

func (m *UPV3Update) Reset()      { *m = UPV3Update{} }
func (*UPV3Update) ProtoMessage() {}
func (*UPV3Update) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{7}
}
func (m *UPV3Update) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UPV3Update) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UPV3Update.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UPV3Update) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UPV3Update.Merge(m, src)
}
func (m *UPV3Update) XXX_Size() int {
	return m.Size()
}
func (m *UPV3Update) XXX_DiscardUnknown() {
	xxx_messageInfo_UPV3Update.DiscardUnknown(m)
}

var xxx_messageInfo_UPV3Update proto.InternalMessageInfo

func (m *UPV3Update) GetInitialize() *gorderbook.UPV3Initialize {
	if m != nil {
		return m.Initialize
	}
	return nil
}

func (m *UPV3Update) GetMint() *gorderbook.UPV3Mint {
	if m != nil {
		return m.Mint
	}
	return nil
}

func (m *UPV3Update) GetBurn() *gorderbook.UPV3Burn {
	if m != nil {
		return m.Burn
	}
	return nil
}

func (m *UPV3Update) GetSwap() *gorderbook.UPV3Swap {
	if m != nil {
		return m.Swap
	}
	return nil
}

func (m *UPV3Update) GetCollect() *gorderbook.UPV3Collect {
	if m != nil {
		return m.Collect
	}
	return nil
}

func (m *UPV3Update) GetFlash() *gorderbook.UPV3Flash {
	if m != nil {
		return m.Flash
	}
	return nil
}

func (m *UPV3Update) GetSetFeeProtocol() *gorderbook.UPV3SetFeeProtocol {
	if m != nil {
		return m.SetFeeProtocol
	}
	return nil
}

func (m *UPV3Update) GetCollectProtocol() *gorderbook.UPV3CollectProtocol {
	if m != nil {
		return m.CollectProtocol
	}
	return nil
}

func (m *UPV3Update) GetRemoved() bool {
	if m != nil {
		return m.Removed
	}
	return false
}

func (m *UPV3Update) GetBlock() uint64 {
	if m != nil {
		return m.Block
	}
	return 0
}

type NftTransfer struct {
	Transfer *gorderbook.NftTransfer `protobuf:"bytes,1,opt,name=transfer,proto3" json:"transfer,omitempty"`
	Removed  bool                    `protobuf:"varint,2,opt,name=removed,proto3" json:"removed,omitempty"`
	Block    uint64                  `protobuf:"varint,3,opt,name=block,proto3" json:"block,omitempty"`
}

func (m *NftTransfer) Reset()      { *m = NftTransfer{} }
func (*NftTransfer) ProtoMessage() {}
func (*NftTransfer) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{8}
}
func (m *NftTransfer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NftTransfer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NftTransfer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NftTransfer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NftTransfer.Merge(m, src)
}
func (m *NftTransfer) XXX_Size() int {
	return m.Size()
}
func (m *NftTransfer) XXX_DiscardUnknown() {
	xxx_messageInfo_NftTransfer.DiscardUnknown(m)
}

var xxx_messageInfo_NftTransfer proto.InternalMessageInfo

func (m *NftTransfer) GetTransfer() *gorderbook.NftTransfer {
	if m != nil {
		return m.Transfer
	}
	return nil
}

func (m *NftTransfer) GetRemoved() bool {
	if m != nil {
		return m.Removed
	}
	return false
}

func (m *NftTransfer) GetBlock() uint64 {
	if m != nil {
		return m.Block
	}
	return 0
}

type Trade struct {
	Price    float64 `protobuf:"fixed64,1,opt,name=price,proto3" json:"price,omitempty"`
	Quantity float64 `protobuf:"fixed64,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	ID       uint64  `protobuf:"varint,3,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (m *Trade) Reset()      { *m = Trade{} }
func (*Trade) ProtoMessage() {}
func (*Trade) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe080a1ca5f6ba79, []int{9}
}
func (m *Trade) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Trade) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Trade.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_fe080a1ca5f6ba79, []int{10}
}
func (m *AggregatedTrade) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AggregatedTrade) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AggregatedTrade.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_fe080a1ca5f6ba79, []int{11}
}
func (m *Liquidation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Liquidation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Liquidation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_fe080a1ca5f6ba79, []int{12}
}
func (m *Funding) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Funding) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Funding.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_fe080a1ca5f6ba79, []int{13}
}
func (m *Stat) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Stat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Stat.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	proto.RegisterType((*UPV3Snapshot)(nil), "models.UPV3Snapshot")
	proto.RegisterType((*UPV3Update)(nil), "models.UPV3Update")
	proto.RegisterType((*NftTransfer)(nil), "models.NftTransfer")
	proto.RegisterType((*Trade)(nil), "models.Trade")
	proto.RegisterType((*AggregatedTrade)(nil), "models.AggregatedTrade")
	proto.RegisterType((*Liquidation)(nil), "models.Liquidation")
	proto.RegisterType((*Funding)(nil), "models.Funding")
	proto.RegisterType((*Stat)(nil), "models.Stat")
}

func init() { proto.RegisterFile("market_data.proto", fileDescriptor_fe080a1ca5f6ba79) }

var fileDescriptor_fe080a1ca5f6ba79 = []byte{
	// 1324 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xdc, 0x56, 0x4f, 0x6f, 0x1b, 0x45,
	0x14, 0xf7, 0x3a, 0xeb, 0x7f, 0xcf, 0x4e, 0xb2, 0x0c, 0x55, 0xbb, 0xa4, 0x65, 0x1b, 0x2c, 0x51,
	0x45, 0x45, 0x75, 0x62, 0xa7, 0xa0, 0xaa, 0x08, 0x55, 0x49, 0xab, 0x14, 0x4b, 0x81, 0x84, 0x4d,
	0x9a, 0x1e, 0x57, 0x63, 0xef, 0x78, 0x33, 0xf2, 0x7a, 0x67, 0xbb, 0x33, 0x6e, 0x1a, 0x84, 0x10,
	0x17, 0xee, 0x7c, 0x01, 0x6e, 0x20, 0x71, 0xe0, 0xc0, 0x89, 0xcf, 0x50, 0x71, 0xea, 0xb1, 0x07,
	0x84, 0xa8, 0x7b, 0xe1, 0xc0, 0xa1, 0x1f, 0x01, 0xcd, 0xec, 0x1f, 0xdb, 0x24, 0x2e, 0x44, 0x45,
	0x08, 0x71, 0xda, 0x79, 0xef, 0xfd, 0x7e, 0x6f, 0xde, 0xbc, 0xf7, 0xf6, 0xcd, 0xc0, 0x6b, 0x03,
	0x1c, 0xf5, 0x89, 0x70, 0x5c, 0x2c, 0x70, 0x23, 0x8c, 0x98, 0x60, 0xa8, 0x38, 0x60, 0x2e, 0xf1,
	0xf9, 0xd2, 0x35, 0x8f, 0x8a, 0xc3, 0x61, 0xa7, 0xd1, 0x65, 0x83, 0x55, 0x8f, 0x79, 0x6c, 0x55,
	0x99, 0x3b, 0xc3, 0x9e, 0x92, 0x94, 0xa0, 0x56, 0x31, 0x6d, 0xe9, 0xb2, 0xc7, 0x98, 0xe7, 0x93,
	0x31, 0x4a, 0xd0, 0x01, 0xe1, 0x02, 0x0f, 0xc2, 0x04, 0x60, 0xfd, 0x19, 0x70, 0x14, 0xe1, 0x30,
	0x24, 0x11, 0x4f, 0xec, 0xb7, 0x3c, 0x2a, 0x7c, 0x1c, 0xef, 0x87, 0xfd, 0xf0, 0x10, 0x0b, 0xda,
	0xed, 0xf3, 0x55, 0x8f, 0x45, 0x2e, 0x89, 0x3a, 0x8c, 0xf5, 0x27, 0x96, 0x8d, 0x38, 0xc6, 0xd8,
	0x57, 0xe2, 0xa0, 0xfe, 0x19, 0xc0, 0xce, 0xe6, 0x76, 0xf3, 0x5e, 0xe8, 0x62, 0x41, 0xd0, 0x1b,
	0x50, 0xee, 0x10, 0x2e, 0x9c, 0x0e, 0x75, 0x4d, 0x6d, 0x59, 0x5b, 0xd1, 0xec, 0x92, 0x94, 0x37,
	0xa9, 0x9b, 0x99, 0x30, 0xef, 0x9b, 0xf9, 0xb1, 0x69, 0x83, 0xf7, 0xd1, 0x0d, 0xa8, 0x64, 0x71,
	0x9b, 0x73, 0xcb, 0xda, 0x4a, 0xb5, 0xb5, 0xd4, 0x88, 0x03, 0x6f, 0xa4, 0x81, 0x37, 0xf6, 0x53,
	0x84, 0x3d, 0x06, 0xd7, 0x3f, 0x87, 0x9a, 0xdc, 0x7d, 0x2f, 0xc0, 0x21, 0x3f, 0x64, 0xe2, 0x5f,
	0xdf, 0xff, 0x6b, 0x4d, 0x1d, 0xbf, 0x95, 0x1c, 0xff, 0x16, 0x14, 0x7d, 0xf2, 0x90, 0xf8, 0xdc,
	0xd4, 0x96, 0xe7, 0x56, 0xaa, 0xad, 0xb7, 0x1a, 0x27, 0xb2, 0xd7, 0xd8, 0x91, 0x8a, 0x4d, 0xc6,
	0xfa, 0xdb, 0x12, 0xb9, 0xa9, 0x3f, 0xfe, 0xe5, 0x72, 0xce, 0x4e, 0x68, 0xd3, 0x91, 0xe4, 0xcf,
	0x10, 0x09, 0x3a, 0x07, 0x05, 0x11, 0x61, 0x97, 0x98, 0xfa, 0xb2, 0xb6, 0x52, 0xb6, 0x63, 0xa1,
	0xfe, 0x53, 0x5e, 0x25, 0xa8, 0x95, 0x25, 0xe8, 0x7d, 0xd0, 0x3b, 0xd4, 0xe5, 0x66, 0xfe, 0x6c,
	0xf1, 0x29, 0x92, 0x24, 0x63, 0xde, 0xe7, 0xe6, 0xdc, 0x19, 0xc9, 0x92, 0x34, 0x7d, 0x34, 0xfd,
	0x2c, 0x47, 0xbb, 0x0d, 0x0b, 0xb2, 0x2b, 0x9d, 0x30, 0x22, 0x5d, 0xca, 0x29, 0x0b, 0xcc, 0x82,
	0xa2, 0x5f, 0x3a, 0x41, 0xbf, 0xd7, 0x0e, 0xc4, 0x7b, 0xd7, 0x0f, 0xb0, 0x3f, 0x24, 0xf6, 0xbc,
	0xe4, 0xec, 0xa6, 0x14, 0xb4, 0x01, 0xf3, 0x3e, 0x13, 0x13, 0x3e, 0x8a, 0x7f, 0xc3, 0x47, 0xcd,
	0x67, 0x22, 0x73, 0x51, 0xff, 0x31, 0xaf, 0x8a, 0xbd, 0x9e, 0x14, 0xbb, 0x35, 0x95, 0x4a, 0x73,
	0x66, 0x36, 0x26, 0x33, 0xd8, 0x9a, 0xca, 0xe0, 0x5f, 0x72, 0xfe, 0x27, 0x89, 0xfb, 0x41, 0x53,
	0x5d, 0xb8, 0x9e, 0x75, 0xe1, 0x7f, 0x3e, 0x75, 0xf5, 0x9f, 0x75, 0xa8, 0xdd, 0xdb, 0x3d, 0x18,
	0x87, 0xdc, 0x84, 0x82, 0x1a, 0x8d, 0xc9, 0x9f, 0x7d, 0xf1, 0x94, 0xfd, 0x25, 0x7e, 0x9f, 0x76,
	0xfb, 0x76, 0x8c, 0x44, 0x1f, 0x40, 0x25, 0x64, 0x9c, 0x0a, 0xca, 0x82, 0xf4, 0xa8, 0x97, 0x67,
	0xd0, 0x76, 0x13, 0x9c, 0x3d, 0x66, 0xa0, 0x4b, 0x50, 0xf1, 0xe9, 0x83, 0x21, 0x75, 0xa9, 0x38,
	0x56, 0x53, 0xa9, 0x66, 0x8f, 0x15, 0xe8, 0x4d, 0x00, 0xfe, 0x20, 0x92, 0x75, 0xa1, 0xdd, 0xf8,
	0xa7, 0xaf, 0xd9, 0x15, 0xa9, 0xd9, 0x95, 0x0a, 0xf4, 0x2e, 0x5c, 0xe8, 0x11, 0xe2, 0x78, 0x11,
	0x3b, 0x12, 0x87, 0x8e, 0xe7, 0xb3, 0x0e, 0xf6, 0x9d, 0xb5, 0x47, 0xcd, 0xd6, 0x0d, 0xd5, 0x03,
	0x35, 0xfb, 0x5c, 0x8f, 0x90, 0xbb, 0xca, 0x7a, 0x57, 0x19, 0x95, 0xed, 0x74, 0x5a, 0x53, 0xd1,
	0x8a, 0xa7, 0xd2, 0x94, 0x0d, 0x5d, 0x81, 0x45, 0x95, 0xce, 0x2e, 0xf3, 0x9d, 0x1e, 0x21, 0xdc,
	0x59, 0x33, 0x4b, 0x0a, 0x3e, 0x9f, 0xaa, 0xb7, 0x08, 0xe1, 0x6b, 0x27, 0x71, 0x4d, 0xb3, 0x7c,
	0x12, 0xd7, 0x44, 0x37, 0x61, 0x49, 0x30, 0x81, 0x7d, 0xe7, 0xa1, 0xec, 0x26, 0xc7, 0x67, 0xdd,
	0x3e, 0x71, 0x1d, 0xc1, 0xfa, 0x24, 0x70, 0xd6, 0xcc, 0x8a, 0xa2, 0x9c, 0x57, 0x08, 0xd5, 0x6e,
	0xdb, 0xca, 0xbe, 0x2f, 0xcd, 0x6b, 0x2f, 0xe5, 0x36, 0x4d, 0x78, 0x09, 0xb7, 0x89, 0x10, 0xe8,
	0xb2, 0x74, 0x66, 0x75, 0x59, 0x5b, 0x29, 0xd8, 0x6a, 0x2d, 0xef, 0x0d, 0x99, 0x12, 0x41, 0x49,
	0x64, 0xd6, 0x94, 0xbe, 0xd4, 0x23, 0x64, 0x9f, 0x92, 0x68, 0xba, 0xbd, 0xe6, 0xcf, 0xd2, 0x5e,
	0xdf, 0xe8, 0x00, 0xb2, 0xee, 0xc9, 0x28, 0xd9, 0x00, 0xa0, 0x01, 0x15, 0x14, 0xfb, 0xf4, 0x53,
	0xa2, 0x2e, 0xae, 0xd3, 0xc7, 0xab, 0xa4, 0xb4, 0x33, 0xa0, 0x3d, 0x41, 0x42, 0xab, 0xa0, 0x0f,
	0x68, 0x20, 0x92, 0x4b, 0x63, 0x56, 0x7b, 0x7e, 0x44, 0x03, 0x61, 0x2b, 0xa0, 0x24, 0x74, 0x86,
	0x51, 0x90, 0xdc, 0x77, 0xb3, 0x08, 0x9b, 0xc3, 0x28, 0xb0, 0x15, 0x50, 0x12, 0xf8, 0x11, 0x4e,
	0xff, 0xa3, 0x59, 0x84, 0xbd, 0x23, 0x1c, 0xda, 0x0a, 0x88, 0x6e, 0x40, 0xa9, 0xcb, 0x7c, 0x9f,
	0x74, 0x45, 0x32, 0x77, 0xac, 0x19, 0x9c, 0xdb, 0x31, 0xca, 0x4e, 0xe1, 0xa8, 0x05, 0x85, 0x9e,
	0x8f, 0xf9, 0xe1, 0xc4, 0xac, 0x39, 0x8d, 0xb7, 0x25, 0x31, 0x76, 0x0c, 0x45, 0x3b, 0x60, 0x70,
	0x22, 0x64, 0x5b, 0x39, 0x69, 0x33, 0xa9, 0x26, 0xac, 0xb6, 0xde, 0x9e, 0x15, 0x2a, 0x11, 0x5b,
	0x84, 0xec, 0x26, 0x60, 0x7b, 0x81, 0x4f, 0xc9, 0xe8, 0x13, 0x30, 0x92, 0x78, 0xc6, 0x0e, 0xcb,
	0xca, 0xe1, 0x95, 0x97, 0x9f, 0x23, 0xf3, 0xb8, 0xd8, 0x9d, 0x56, 0x20, 0x13, 0x4a, 0x11, 0x19,
	0xb0, 0x87, 0xc4, 0x55, 0x4d, 0x5c, 0xb6, 0x53, 0x51, 0x5e, 0xdf, 0x1d, 0xd9, 0xaa, 0xaa, 0x41,
	0x75, 0x3b, 0x16, 0xea, 0xc7, 0x50, 0xfd, 0xb8, 0x27, 0xf6, 0x23, 0x1c, 0xf0, 0x1e, 0x89, 0xd0,
	0x4d, 0x28, 0x8b, 0x64, 0x9d, 0x34, 0xc9, 0x69, 0x19, 0x9d, 0x60, 0xd8, 0x19, 0x7e, 0x72, 0xeb,
	0xfc, 0x8c, 0xad, 0xe7, 0x26, 0xb7, 0x6e, 0x43, 0x61, 0x5f, 0x3e, 0x21, 0xa4, 0x39, 0x9e, 0x31,
	0xf1, 0x7b, 0x2a, 0x16, 0xd0, 0x12, 0x94, 0x1f, 0x0c, 0x71, 0x20, 0xe4, 0x6c, 0x8a, 0x5f, 0x53,
	0x99, 0x8c, 0x16, 0x20, 0xdf, 0xbe, 0x93, 0x78, 0xcb, 0xb7, 0xef, 0xd4, 0xbf, 0xd7, 0x60, 0x71,
	0xc3, 0xf3, 0x22, 0xe2, 0x61, 0x41, 0xdc, 0xd8, 0xab, 0x01, 0x73, 0xe9, 0x1b, 0xad, 0x6c, 0xcb,
	0xe5, 0x2b, 0x3c, 0x7d, 0x96, 0xa1, 0x8a, 0x53, 0xf7, 0xd9, 0xc6, 0x93, 0x2a, 0xf4, 0x0e, 0x14,
	0xd5, 0x7b, 0x88, 0x9b, 0xba, 0x1a, 0xc3, 0xf3, 0x69, 0xae, 0x54, 0x30, 0xe9, 0x1b, 0x2c, 0x86,
	0xd4, 0xbf, 0xd5, 0xa0, 0xba, 0xad, 0xe6, 0x2c, 0x96, 0x83, 0xf8, 0x1f, 0x0d, 0xd5, 0x84, 0x92,
	0xaa, 0x57, 0x16, 0x66, 0x2a, 0x8e, 0xd3, 0xac, 0xcf, 0x4a, 0x73, 0x61, 0x3a, 0xcd, 0xf5, 0xfb,
	0x50, 0xda, 0x1a, 0x06, 0x2e, 0x0d, 0xbc, 0xe9, 0x80, 0xb4, 0xb3, 0x04, 0x84, 0x40, 0x8f, 0xb0,
	0x20, 0x49, 0x0d, 0xd5, 0xba, 0xfe, 0xa5, 0x06, 0xfa, 0x9e, 0xc0, 0xe2, 0x15, 0xdc, 0x5e, 0x83,
	0x0a, 0x17, 0x58, 0x38, 0xe2, 0x38, 0x8c, 0x7d, 0x2f, 0xb4, 0x8c, 0x34, 0xe7, 0xd2, 0xf5, 0xfe,
	0x71, 0x48, 0xec, 0x32, 0x4f, 0x56, 0xf2, 0xf0, 0x6a, 0x5a, 0xab, 0xa4, 0x68, 0x76, 0x2c, 0x5c,
	0xfd, 0x5d, 0x83, 0x72, 0x0a, 0x46, 0x0b, 0x00, 0xed, 0xc0, 0x25, 0x8f, 0xd4, 0xd0, 0x36, 0x72,
	0xc8, 0x80, 0xda, 0x4e, 0x48, 0x02, 0x1a, 0x78, 0xea, 0xc2, 0x33, 0x34, 0xa9, 0xb9, 0xed, 0x33,
	0x9e, 0x69, 0xf2, 0xe8, 0x75, 0x58, 0xdc, 0x23, 0x42, 0xf8, 0x64, 0x40, 0x82, 0xf8, 0x5e, 0x34,
	0xe6, 0xd0, 0x45, 0xb8, 0x20, 0xab, 0x4e, 0x03, 0x6f, 0x8f, 0x70, 0xf9, 0x3c, 0xf9, 0x90, 0x7a,
	0x87, 0xb1, 0x51, 0x47, 0x4b, 0x70, 0x7e, 0xda, 0xb8, 0xcd, 0x8e, 0x62, 0x5b, 0xe1, 0x24, 0xf1,
	0xe0, 0xfe, 0xc6, 0x6e, 0x6c, 0x2c, 0xa2, 0x79, 0xa8, 0xb4, 0x07, 0x1d, 0xec, 0xe3, 0xa0, 0x4b,
	0x8c, 0x12, 0x5a, 0x84, 0xaa, 0x6a, 0xad, 0x03, 0xe6, 0x0f, 0x07, 0xc4, 0x28, 0xa7, 0xe1, 0xb6,
	0x03, 0x41, 0x22, 0xc2, 0x85, 0x51, 0x91, 0x90, 0xa4, 0x7c, 0x36, 0x16, 0xc4, 0x80, 0xab, 0x0d,
	0x38, 0x97, 0x3d, 0x9f, 0xd3, 0xdf, 0x45, 0xf6, 0x5f, 0x11, 0xf2, 0xdb, 0x4d, 0x23, 0xa7, 0xbe,
	0x2d, 0x43, 0x53, 0xdf, 0x75, 0x23, 0xbf, 0x79, 0xfd, 0xc9, 0x33, 0x2b, 0xf7, 0xf4, 0x99, 0x95,
	0x7b, 0xf1, 0xcc, 0xd2, 0xbe, 0x18, 0x59, 0xda, 0x77, 0x23, 0x4b, 0x7b, 0x3c, 0xb2, 0xb4, 0x27,
	0x23, 0x4b, 0xfb, 0x75, 0x64, 0x69, 0xbf, 0x8d, 0xac, 0xdc, 0x8b, 0x91, 0xa5, 0x7d, 0xf5, 0xdc,
	0xca, 0x3d, 0x79, 0x6e, 0xe5, 0x9e, 0x3e, 0xb7, 0x72, 0x9d, 0xa2, 0x2a, 0xdc, 0xfa, 0x1f, 0x01,
	0x00, 0x00, 0xff, 0xff, 0xc5, 0x32, 0xb0, 0x6b, 0x84, 0x0e, 0x00, 0x00,
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
func (this *UPV3Snapshot) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*UPV3Snapshot)
	if !ok {
		that2, ok := that.(UPV3Snapshot)
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
	if len(this.Ticks) != len(that1.Ticks) {
		return false
	}
	for i := range this.Ticks {
		if !this.Ticks[i].Equal(that1.Ticks[i]) {
			return false
		}
	}
	if len(this.Positions) != len(that1.Positions) {
		return false
	}
	for i := range this.Positions {
		if !this.Positions[i].Equal(that1.Positions[i]) {
			return false
		}
	}
	if !bytes.Equal(this.Liquidity, that1.Liquidity) {
		return false
	}
	if !bytes.Equal(this.SqrtPrice, that1.SqrtPrice) {
		return false
	}
	if !bytes.Equal(this.FeeGrowthGlobal_0X128, that1.FeeGrowthGlobal_0X128) {
		return false
	}
	if !bytes.Equal(this.FeeGrowthGlobal_1X128, that1.FeeGrowthGlobal_1X128) {
		return false
	}
	if !bytes.Equal(this.ProtocolFees_0, that1.ProtocolFees_0) {
		return false
	}
	if !bytes.Equal(this.ProtocolFees_1, that1.ProtocolFees_1) {
		return false
	}
	if !bytes.Equal(this.TotalValueLockedToken_0, that1.TotalValueLockedToken_0) {
		return false
	}
	if !bytes.Equal(this.TotalValueLockedToken_1, that1.TotalValueLockedToken_1) {
		return false
	}
	if this.Tick != that1.Tick {
		return false
	}
	if this.FeeTier != that1.FeeTier {
		return false
	}
	if !this.Timestamp.Equal(that1.Timestamp) {
		return false
	}
	return true
}
func (this *UPV3Update) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*UPV3Update)
	if !ok {
		that2, ok := that.(UPV3Update)
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
	if !this.Initialize.Equal(that1.Initialize) {
		return false
	}
	if !this.Mint.Equal(that1.Mint) {
		return false
	}
	if !this.Burn.Equal(that1.Burn) {
		return false
	}
	if !this.Swap.Equal(that1.Swap) {
		return false
	}
	if !this.Collect.Equal(that1.Collect) {
		return false
	}
	if !this.Flash.Equal(that1.Flash) {
		return false
	}
	if !this.SetFeeProtocol.Equal(that1.SetFeeProtocol) {
		return false
	}
	if !this.CollectProtocol.Equal(that1.CollectProtocol) {
		return false
	}
	if this.Removed != that1.Removed {
		return false
	}
	if this.Block != that1.Block {
		return false
	}
	return true
}
func (this *NftTransfer) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*NftTransfer)
	if !ok {
		that2, ok := that.(NftTransfer)
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
	if !this.Transfer.Equal(that1.Transfer) {
		return false
	}
	if this.Removed != that1.Removed {
		return false
	}
	if this.Block != that1.Block {
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
		vs := make([]gorderbook.OrderBookLevel, len(this.Levels))
		for i := range vs {
			vs[i] = this.Levels[i]
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
		vs := make([]gorderbook.OrderBookLevel, len(this.Bids))
		for i := range vs {
			vs[i] = this.Bids[i]
		}
		s = append(s, "Bids: "+fmt.Sprintf("%#v", vs)+",\n")
	}
	if this.Asks != nil {
		vs := make([]gorderbook.OrderBookLevel, len(this.Asks))
		for i := range vs {
			vs[i] = this.Asks[i]
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
		vs := make([]gorderbook.Order, len(this.Bids))
		for i := range vs {
			vs[i] = this.Bids[i]
		}
		s = append(s, "Bids: "+fmt.Sprintf("%#v", vs)+",\n")
	}
	if this.Asks != nil {
		vs := make([]gorderbook.Order, len(this.Asks))
		for i := range vs {
			vs[i] = this.Asks[i]
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
		vs := make([]gorderbook.Order, len(this.Bids))
		for i := range vs {
			vs[i] = this.Bids[i]
		}
		s = append(s, "Bids: "+fmt.Sprintf("%#v", vs)+",\n")
	}
	if this.Asks != nil {
		vs := make([]gorderbook.Order, len(this.Asks))
		for i := range vs {
			vs[i] = this.Asks[i]
		}
		s = append(s, "Asks: "+fmt.Sprintf("%#v", vs)+",\n")
	}
	if this.Timestamp != nil {
		s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *UPV3Snapshot) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 17)
	s = append(s, "&models.UPV3Snapshot{")
	if this.Ticks != nil {
		s = append(s, "Ticks: "+fmt.Sprintf("%#v", this.Ticks)+",\n")
	}
	if this.Positions != nil {
		s = append(s, "Positions: "+fmt.Sprintf("%#v", this.Positions)+",\n")
	}
	s = append(s, "Liquidity: "+fmt.Sprintf("%#v", this.Liquidity)+",\n")
	s = append(s, "SqrtPrice: "+fmt.Sprintf("%#v", this.SqrtPrice)+",\n")
	s = append(s, "FeeGrowthGlobal_0X128: "+fmt.Sprintf("%#v", this.FeeGrowthGlobal_0X128)+",\n")
	s = append(s, "FeeGrowthGlobal_1X128: "+fmt.Sprintf("%#v", this.FeeGrowthGlobal_1X128)+",\n")
	s = append(s, "ProtocolFees_0: "+fmt.Sprintf("%#v", this.ProtocolFees_0)+",\n")
	s = append(s, "ProtocolFees_1: "+fmt.Sprintf("%#v", this.ProtocolFees_1)+",\n")
	s = append(s, "TotalValueLockedToken_0: "+fmt.Sprintf("%#v", this.TotalValueLockedToken_0)+",\n")
	s = append(s, "TotalValueLockedToken_1: "+fmt.Sprintf("%#v", this.TotalValueLockedToken_1)+",\n")
	s = append(s, "Tick: "+fmt.Sprintf("%#v", this.Tick)+",\n")
	s = append(s, "FeeTier: "+fmt.Sprintf("%#v", this.FeeTier)+",\n")
	if this.Timestamp != nil {
		s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *UPV3Update) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 14)
	s = append(s, "&models.UPV3Update{")
	if this.Initialize != nil {
		s = append(s, "Initialize: "+fmt.Sprintf("%#v", this.Initialize)+",\n")
	}
	if this.Mint != nil {
		s = append(s, "Mint: "+fmt.Sprintf("%#v", this.Mint)+",\n")
	}
	if this.Burn != nil {
		s = append(s, "Burn: "+fmt.Sprintf("%#v", this.Burn)+",\n")
	}
	if this.Swap != nil {
		s = append(s, "Swap: "+fmt.Sprintf("%#v", this.Swap)+",\n")
	}
	if this.Collect != nil {
		s = append(s, "Collect: "+fmt.Sprintf("%#v", this.Collect)+",\n")
	}
	if this.Flash != nil {
		s = append(s, "Flash: "+fmt.Sprintf("%#v", this.Flash)+",\n")
	}
	if this.SetFeeProtocol != nil {
		s = append(s, "SetFeeProtocol: "+fmt.Sprintf("%#v", this.SetFeeProtocol)+",\n")
	}
	if this.CollectProtocol != nil {
		s = append(s, "CollectProtocol: "+fmt.Sprintf("%#v", this.CollectProtocol)+",\n")
	}
	s = append(s, "Removed: "+fmt.Sprintf("%#v", this.Removed)+",\n")
	s = append(s, "Block: "+fmt.Sprintf("%#v", this.Block)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *NftTransfer) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&models.NftTransfer{")
	if this.Transfer != nil {
		s = append(s, "Transfer: "+fmt.Sprintf("%#v", this.Transfer)+",\n")
	}
	s = append(s, "Removed: "+fmt.Sprintf("%#v", this.Removed)+",\n")
	s = append(s, "Block: "+fmt.Sprintf("%#v", this.Block)+",\n")
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
		vs := make([]Trade, len(this.Trades))
		for i := range vs {
			vs[i] = this.Trades[i]
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
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OBL1Update) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OBL1Update) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Timestamp != nil {
		{
			size, err := m.Timestamp.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.BestAsk != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.BestAsk))))
		i--
		dAtA[i] = 0x11
	}
	if m.BestBid != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.BestBid))))
		i--
		dAtA[i] = 0x9
	}
	return len(dAtA) - i, nil
}

func (m *OBL1Snapshot) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OBL1Snapshot) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OBL1Snapshot) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Timestamp != nil {
		{
			size, err := m.Timestamp.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.BestAsk != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.BestAsk))))
		i--
		dAtA[i] = 0x11
	}
	if m.BestBid != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.BestBid))))
		i--
		dAtA[i] = 0x9
	}
	return len(dAtA) - i, nil
}

func (m *OBL2Update) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OBL2Update) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OBL2Update) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Trade {
		i--
		if m.Trade {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if m.Timestamp != nil {
		{
			size, err := m.Timestamp.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Levels) > 0 {
		for iNdEx := len(m.Levels) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Levels[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMarketData(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *OBL2Snapshot) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OBL2Snapshot) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OBL2Snapshot) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LotPrecision != nil {
		{
			size, err := m.LotPrecision.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.TickPrecision != nil {
		{
			size, err := m.TickPrecision.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.Timestamp != nil {
		{
			size, err := m.Timestamp.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if len(m.Asks) > 0 {
		for iNdEx := len(m.Asks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Asks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMarketData(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Bids) > 0 {
		for iNdEx := len(m.Bids) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Bids[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMarketData(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	return len(dAtA) - i, nil
}

func (m *OBL3Update) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OBL3Update) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OBL3Update) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LotPrecision != nil {
		{
			size, err := m.LotPrecision.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.TickPrecision != nil {
		{
			size, err := m.TickPrecision.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.Timestamp != nil {
		{
			size, err := m.Timestamp.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if len(m.Asks) > 0 {
		for iNdEx := len(m.Asks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Asks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMarketData(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Bids) > 0 {
		for iNdEx := len(m.Bids) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Bids[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMarketData(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	return len(dAtA) - i, nil
}

func (m *OBL3Snapshot) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OBL3Snapshot) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OBL3Snapshot) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Timestamp != nil {
		{
			size, err := m.Timestamp.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if len(m.Asks) > 0 {
		for iNdEx := len(m.Asks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Asks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMarketData(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Bids) > 0 {
		for iNdEx := len(m.Bids) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Bids[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMarketData(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	return len(dAtA) - i, nil
}

func (m *UPV3Snapshot) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UPV3Snapshot) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UPV3Snapshot) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Timestamp != nil {
		{
			size, err := m.Timestamp.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x6a
	}
	if m.FeeTier != 0 {
		i = encodeVarintMarketData(dAtA, i, uint64(m.FeeTier))
		i--
		dAtA[i] = 0x60
	}
	if m.Tick != 0 {
		i = encodeVarintMarketData(dAtA, i, uint64(m.Tick))
		i--
		dAtA[i] = 0x58
	}
	if len(m.TotalValueLockedToken_1) > 0 {
		i -= len(m.TotalValueLockedToken_1)
		copy(dAtA[i:], m.TotalValueLockedToken_1)
		i = encodeVarintMarketData(dAtA, i, uint64(len(m.TotalValueLockedToken_1)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.TotalValueLockedToken_0) > 0 {
		i -= len(m.TotalValueLockedToken_0)
		copy(dAtA[i:], m.TotalValueLockedToken_0)
		i = encodeVarintMarketData(dAtA, i, uint64(len(m.TotalValueLockedToken_0)))
		i--
		dAtA[i] = 0x4a
	}
	if len(m.ProtocolFees_1) > 0 {
		i -= len(m.ProtocolFees_1)
		copy(dAtA[i:], m.ProtocolFees_1)
		i = encodeVarintMarketData(dAtA, i, uint64(len(m.ProtocolFees_1)))
		i--
		dAtA[i] = 0x42
	}
	if len(m.ProtocolFees_0) > 0 {
		i -= len(m.ProtocolFees_0)
		copy(dAtA[i:], m.ProtocolFees_0)
		i = encodeVarintMarketData(dAtA, i, uint64(len(m.ProtocolFees_0)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.FeeGrowthGlobal_1X128) > 0 {
		i -= len(m.FeeGrowthGlobal_1X128)
		copy(dAtA[i:], m.FeeGrowthGlobal_1X128)
		i = encodeVarintMarketData(dAtA, i, uint64(len(m.FeeGrowthGlobal_1X128)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.FeeGrowthGlobal_0X128) > 0 {
		i -= len(m.FeeGrowthGlobal_0X128)
		copy(dAtA[i:], m.FeeGrowthGlobal_0X128)
		i = encodeVarintMarketData(dAtA, i, uint64(len(m.FeeGrowthGlobal_0X128)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.SqrtPrice) > 0 {
		i -= len(m.SqrtPrice)
		copy(dAtA[i:], m.SqrtPrice)
		i = encodeVarintMarketData(dAtA, i, uint64(len(m.SqrtPrice)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Liquidity) > 0 {
		i -= len(m.Liquidity)
		copy(dAtA[i:], m.Liquidity)
		i = encodeVarintMarketData(dAtA, i, uint64(len(m.Liquidity)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Positions) > 0 {
		for iNdEx := len(m.Positions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Positions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMarketData(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Ticks) > 0 {
		for iNdEx := len(m.Ticks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Ticks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMarketData(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *UPV3Update) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UPV3Update) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UPV3Update) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Block != 0 {
		i = encodeVarintMarketData(dAtA, i, uint64(m.Block))
		i--
		dAtA[i] = 0x50
	}
	if m.Removed {
		i--
		if m.Removed {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x48
	}
	if m.CollectProtocol != nil {
		{
			size, err := m.CollectProtocol.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x42
	}
	if m.SetFeeProtocol != nil {
		{
			size, err := m.SetFeeProtocol.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x3a
	}
	if m.Flash != nil {
		{
			size, err := m.Flash.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.Collect != nil {
		{
			size, err := m.Collect.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.Swap != nil {
		{
			size, err := m.Swap.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Burn != nil {
		{
			size, err := m.Burn.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Mint != nil {
		{
			size, err := m.Mint.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Initialize != nil {
		{
			size, err := m.Initialize.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *NftTransfer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NftTransfer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NftTransfer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Block != 0 {
		i = encodeVarintMarketData(dAtA, i, uint64(m.Block))
		i--
		dAtA[i] = 0x18
	}
	if m.Removed {
		i--
		if m.Removed {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.Transfer != nil {
		{
			size, err := m.Transfer.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Trade) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Trade) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Trade) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ID != 0 {
		i = encodeVarintMarketData(dAtA, i, uint64(m.ID))
		i--
		dAtA[i] = 0x18
	}
	if m.Quantity != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Quantity))))
		i--
		dAtA[i] = 0x11
	}
	if m.Price != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Price))))
		i--
		dAtA[i] = 0x9
	}
	return len(dAtA) - i, nil
}

func (m *AggregatedTrade) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AggregatedTrade) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AggregatedTrade) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Trades) > 0 {
		for iNdEx := len(m.Trades) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Trades[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMarketData(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.AggregateID != 0 {
		i = encodeVarintMarketData(dAtA, i, uint64(m.AggregateID))
		i--
		dAtA[i] = 0x18
	}
	if m.Timestamp != nil {
		{
			size, err := m.Timestamp.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Bid {
		i--
		if m.Bid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Liquidation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Liquidation) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Liquidation) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Quantity != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Quantity))))
		i--
		dAtA[i] = 0x29
	}
	if m.Price != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Price))))
		i--
		dAtA[i] = 0x21
	}
	if m.OrderID != 0 {
		i = encodeVarintMarketData(dAtA, i, uint64(m.OrderID))
		i--
		dAtA[i] = 0x18
	}
	if m.Timestamp != nil {
		{
			size, err := m.Timestamp.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Bid {
		i--
		if m.Bid {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Funding) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Funding) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Funding) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Rate != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Rate))))
		i--
		dAtA[i] = 0x11
	}
	if m.Timestamp != nil {
		{
			size, err := m.Timestamp.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Stat) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Stat) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Stat) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Value != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Value))))
		i--
		dAtA[i] = 0x19
	}
	if m.StatType != 0 {
		i = encodeVarintMarketData(dAtA, i, uint64(m.StatType))
		i--
		dAtA[i] = 0x10
	}
	if m.Timestamp != nil {
		{
			size, err := m.Timestamp.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMarketData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintMarketData(dAtA []byte, offset int, v uint64) int {
	offset -= sovMarketData(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
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

func (m *UPV3Snapshot) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Ticks) > 0 {
		for _, e := range m.Ticks {
			l = e.Size()
			n += 1 + l + sovMarketData(uint64(l))
		}
	}
	if len(m.Positions) > 0 {
		for _, e := range m.Positions {
			l = e.Size()
			n += 1 + l + sovMarketData(uint64(l))
		}
	}
	l = len(m.Liquidity)
	if l > 0 {
		n += 1 + l + sovMarketData(uint64(l))
	}
	l = len(m.SqrtPrice)
	if l > 0 {
		n += 1 + l + sovMarketData(uint64(l))
	}
	l = len(m.FeeGrowthGlobal_0X128)
	if l > 0 {
		n += 1 + l + sovMarketData(uint64(l))
	}
	l = len(m.FeeGrowthGlobal_1X128)
	if l > 0 {
		n += 1 + l + sovMarketData(uint64(l))
	}
	l = len(m.ProtocolFees_0)
	if l > 0 {
		n += 1 + l + sovMarketData(uint64(l))
	}
	l = len(m.ProtocolFees_1)
	if l > 0 {
		n += 1 + l + sovMarketData(uint64(l))
	}
	l = len(m.TotalValueLockedToken_0)
	if l > 0 {
		n += 1 + l + sovMarketData(uint64(l))
	}
	l = len(m.TotalValueLockedToken_1)
	if l > 0 {
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.Tick != 0 {
		n += 1 + sovMarketData(uint64(m.Tick))
	}
	if m.FeeTier != 0 {
		n += 1 + sovMarketData(uint64(m.FeeTier))
	}
	if m.Timestamp != nil {
		l = m.Timestamp.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	return n
}

func (m *UPV3Update) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Initialize != nil {
		l = m.Initialize.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.Mint != nil {
		l = m.Mint.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.Burn != nil {
		l = m.Burn.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.Swap != nil {
		l = m.Swap.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.Collect != nil {
		l = m.Collect.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.Flash != nil {
		l = m.Flash.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.SetFeeProtocol != nil {
		l = m.SetFeeProtocol.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.CollectProtocol != nil {
		l = m.CollectProtocol.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.Removed {
		n += 2
	}
	if m.Block != 0 {
		n += 1 + sovMarketData(uint64(m.Block))
	}
	return n
}

func (m *NftTransfer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Transfer != nil {
		l = m.Transfer.Size()
		n += 1 + l + sovMarketData(uint64(l))
	}
	if m.Removed {
		n += 2
	}
	if m.Block != 0 {
		n += 1 + sovMarketData(uint64(m.Block))
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
	return (math_bits.Len64(x|1) + 6) / 7
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
	repeatedStringForLevels := "[]OrderBookLevel{"
	for _, f := range this.Levels {
		repeatedStringForLevels += fmt.Sprintf("%v", f) + ","
	}
	repeatedStringForLevels += "}"
	s := strings.Join([]string{`&OBL2Update{`,
		`Levels:` + repeatedStringForLevels + `,`,
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
	repeatedStringForBids := "[]OrderBookLevel{"
	for _, f := range this.Bids {
		repeatedStringForBids += fmt.Sprintf("%v", f) + ","
	}
	repeatedStringForBids += "}"
	repeatedStringForAsks := "[]OrderBookLevel{"
	for _, f := range this.Asks {
		repeatedStringForAsks += fmt.Sprintf("%v", f) + ","
	}
	repeatedStringForAsks += "}"
	s := strings.Join([]string{`&OBL2Snapshot{`,
		`Bids:` + repeatedStringForBids + `,`,
		`Asks:` + repeatedStringForAsks + `,`,
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
	repeatedStringForBids := "[]Order{"
	for _, f := range this.Bids {
		repeatedStringForBids += fmt.Sprintf("%v", f) + ","
	}
	repeatedStringForBids += "}"
	repeatedStringForAsks := "[]Order{"
	for _, f := range this.Asks {
		repeatedStringForAsks += fmt.Sprintf("%v", f) + ","
	}
	repeatedStringForAsks += "}"
	s := strings.Join([]string{`&OBL3Update{`,
		`Bids:` + repeatedStringForBids + `,`,
		`Asks:` + repeatedStringForAsks + `,`,
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
	repeatedStringForBids := "[]Order{"
	for _, f := range this.Bids {
		repeatedStringForBids += fmt.Sprintf("%v", f) + ","
	}
	repeatedStringForBids += "}"
	repeatedStringForAsks := "[]Order{"
	for _, f := range this.Asks {
		repeatedStringForAsks += fmt.Sprintf("%v", f) + ","
	}
	repeatedStringForAsks += "}"
	s := strings.Join([]string{`&OBL3Snapshot{`,
		`Bids:` + repeatedStringForBids + `,`,
		`Asks:` + repeatedStringForAsks + `,`,
		`Timestamp:` + strings.Replace(fmt.Sprintf("%v", this.Timestamp), "Timestamp", "types.Timestamp", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *UPV3Snapshot) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForTicks := "[]*UPV3Tick{"
	for _, f := range this.Ticks {
		repeatedStringForTicks += strings.Replace(fmt.Sprintf("%v", f), "UPV3Tick", "gorderbook.UPV3Tick", 1) + ","
	}
	repeatedStringForTicks += "}"
	repeatedStringForPositions := "[]*UPV3Position{"
	for _, f := range this.Positions {
		repeatedStringForPositions += strings.Replace(fmt.Sprintf("%v", f), "UPV3Position", "gorderbook.UPV3Position", 1) + ","
	}
	repeatedStringForPositions += "}"
	s := strings.Join([]string{`&UPV3Snapshot{`,
		`Ticks:` + repeatedStringForTicks + `,`,
		`Positions:` + repeatedStringForPositions + `,`,
		`Liquidity:` + fmt.Sprintf("%v", this.Liquidity) + `,`,
		`SqrtPrice:` + fmt.Sprintf("%v", this.SqrtPrice) + `,`,
		`FeeGrowthGlobal_0X128:` + fmt.Sprintf("%v", this.FeeGrowthGlobal_0X128) + `,`,
		`FeeGrowthGlobal_1X128:` + fmt.Sprintf("%v", this.FeeGrowthGlobal_1X128) + `,`,
		`ProtocolFees_0:` + fmt.Sprintf("%v", this.ProtocolFees_0) + `,`,
		`ProtocolFees_1:` + fmt.Sprintf("%v", this.ProtocolFees_1) + `,`,
		`TotalValueLockedToken_0:` + fmt.Sprintf("%v", this.TotalValueLockedToken_0) + `,`,
		`TotalValueLockedToken_1:` + fmt.Sprintf("%v", this.TotalValueLockedToken_1) + `,`,
		`Tick:` + fmt.Sprintf("%v", this.Tick) + `,`,
		`FeeTier:` + fmt.Sprintf("%v", this.FeeTier) + `,`,
		`Timestamp:` + strings.Replace(fmt.Sprintf("%v", this.Timestamp), "Timestamp", "types.Timestamp", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *UPV3Update) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&UPV3Update{`,
		`Initialize:` + strings.Replace(fmt.Sprintf("%v", this.Initialize), "UPV3Initialize", "gorderbook.UPV3Initialize", 1) + `,`,
		`Mint:` + strings.Replace(fmt.Sprintf("%v", this.Mint), "UPV3Mint", "gorderbook.UPV3Mint", 1) + `,`,
		`Burn:` + strings.Replace(fmt.Sprintf("%v", this.Burn), "UPV3Burn", "gorderbook.UPV3Burn", 1) + `,`,
		`Swap:` + strings.Replace(fmt.Sprintf("%v", this.Swap), "UPV3Swap", "gorderbook.UPV3Swap", 1) + `,`,
		`Collect:` + strings.Replace(fmt.Sprintf("%v", this.Collect), "UPV3Collect", "gorderbook.UPV3Collect", 1) + `,`,
		`Flash:` + strings.Replace(fmt.Sprintf("%v", this.Flash), "UPV3Flash", "gorderbook.UPV3Flash", 1) + `,`,
		`SetFeeProtocol:` + strings.Replace(fmt.Sprintf("%v", this.SetFeeProtocol), "UPV3SetFeeProtocol", "gorderbook.UPV3SetFeeProtocol", 1) + `,`,
		`CollectProtocol:` + strings.Replace(fmt.Sprintf("%v", this.CollectProtocol), "UPV3CollectProtocol", "gorderbook.UPV3CollectProtocol", 1) + `,`,
		`Removed:` + fmt.Sprintf("%v", this.Removed) + `,`,
		`Block:` + fmt.Sprintf("%v", this.Block) + `,`,
		`}`,
	}, "")
	return s
}
func (this *NftTransfer) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&NftTransfer{`,
		`Transfer:` + strings.Replace(fmt.Sprintf("%v", this.Transfer), "NftTransfer", "gorderbook.NftTransfer", 1) + `,`,
		`Removed:` + fmt.Sprintf("%v", this.Removed) + `,`,
		`Block:` + fmt.Sprintf("%v", this.Block) + `,`,
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
	repeatedStringForTrades := "[]Trade{"
	for _, f := range this.Trades {
		repeatedStringForTrades += strings.Replace(strings.Replace(f.String(), "Trade", "Trade", 1), `&`, ``, 1) + ","
	}
	repeatedStringForTrades += "}"
	s := strings.Join([]string{`&AggregatedTrade{`,
		`Bid:` + fmt.Sprintf("%v", this.Bid) + `,`,
		`Timestamp:` + strings.Replace(fmt.Sprintf("%v", this.Timestamp), "Timestamp", "types.Timestamp", 1) + `,`,
		`AggregateID:` + fmt.Sprintf("%v", this.AggregateID) + `,`,
		`Trades:` + repeatedStringForTrades + `,`,
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *UPV3Snapshot) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: UPV3Snapshot: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UPV3Snapshot: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ticks", wireType)
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
			m.Ticks = append(m.Ticks, &gorderbook.UPV3Tick{})
			if err := m.Ticks[len(m.Ticks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Positions", wireType)
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
			m.Positions = append(m.Positions, &gorderbook.UPV3Position{})
			if err := m.Positions[len(m.Positions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Liquidity", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Liquidity = append(m.Liquidity[:0], dAtA[iNdEx:postIndex]...)
			if m.Liquidity == nil {
				m.Liquidity = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SqrtPrice", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SqrtPrice = append(m.SqrtPrice[:0], dAtA[iNdEx:postIndex]...)
			if m.SqrtPrice == nil {
				m.SqrtPrice = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeGrowthGlobal_0X128", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FeeGrowthGlobal_0X128 = append(m.FeeGrowthGlobal_0X128[:0], dAtA[iNdEx:postIndex]...)
			if m.FeeGrowthGlobal_0X128 == nil {
				m.FeeGrowthGlobal_0X128 = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeGrowthGlobal_1X128", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FeeGrowthGlobal_1X128 = append(m.FeeGrowthGlobal_1X128[:0], dAtA[iNdEx:postIndex]...)
			if m.FeeGrowthGlobal_1X128 == nil {
				m.FeeGrowthGlobal_1X128 = []byte{}
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProtocolFees_0", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProtocolFees_0 = append(m.ProtocolFees_0[:0], dAtA[iNdEx:postIndex]...)
			if m.ProtocolFees_0 == nil {
				m.ProtocolFees_0 = []byte{}
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProtocolFees_1", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProtocolFees_1 = append(m.ProtocolFees_1[:0], dAtA[iNdEx:postIndex]...)
			if m.ProtocolFees_1 == nil {
				m.ProtocolFees_1 = []byte{}
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalValueLockedToken_0", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TotalValueLockedToken_0 = append(m.TotalValueLockedToken_0[:0], dAtA[iNdEx:postIndex]...)
			if m.TotalValueLockedToken_0 == nil {
				m.TotalValueLockedToken_0 = []byte{}
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalValueLockedToken_1", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMarketData
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMarketData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TotalValueLockedToken_1 = append(m.TotalValueLockedToken_1[:0], dAtA[iNdEx:postIndex]...)
			if m.TotalValueLockedToken_1 == nil {
				m.TotalValueLockedToken_1 = []byte{}
			}
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tick", wireType)
			}
			m.Tick = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Tick |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeTier", wireType)
			}
			m.FeeTier = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FeeTier |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 13:
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *UPV3Update) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: UPV3Update: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UPV3Update: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Initialize", wireType)
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
			if m.Initialize == nil {
				m.Initialize = &gorderbook.UPV3Initialize{}
			}
			if err := m.Initialize.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Mint", wireType)
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
			if m.Mint == nil {
				m.Mint = &gorderbook.UPV3Mint{}
			}
			if err := m.Mint.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Burn", wireType)
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
			if m.Burn == nil {
				m.Burn = &gorderbook.UPV3Burn{}
			}
			if err := m.Burn.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Swap", wireType)
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
			if m.Swap == nil {
				m.Swap = &gorderbook.UPV3Swap{}
			}
			if err := m.Swap.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Collect", wireType)
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
			if m.Collect == nil {
				m.Collect = &gorderbook.UPV3Collect{}
			}
			if err := m.Collect.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Flash", wireType)
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
			if m.Flash == nil {
				m.Flash = &gorderbook.UPV3Flash{}
			}
			if err := m.Flash.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SetFeeProtocol", wireType)
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
			if m.SetFeeProtocol == nil {
				m.SetFeeProtocol = &gorderbook.UPV3SetFeeProtocol{}
			}
			if err := m.SetFeeProtocol.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollectProtocol", wireType)
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
			if m.CollectProtocol == nil {
				m.CollectProtocol = &gorderbook.UPV3CollectProtocol{}
			}
			if err := m.CollectProtocol.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Removed", wireType)
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
			m.Removed = bool(v != 0)
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Block", wireType)
			}
			m.Block = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Block |= uint64(b&0x7F) << shift
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *NftTransfer) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: NftTransfer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NftTransfer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Transfer", wireType)
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
			if m.Transfer == nil {
				m.Transfer = &gorderbook.NftTransfer{}
			}
			if err := m.Transfer.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Removed", wireType)
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
			m.Removed = bool(v != 0)
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Block", wireType)
			}
			m.Block = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Block |= uint64(b&0x7F) << shift
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
	depth := 0
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
		case 1:
			iNdEx += 8
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
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMarketData
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMarketData
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMarketData        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMarketData          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMarketData = fmt.Errorf("proto: unexpected end of group")
)
