package models

import (
	encoding_binary "encoding/binary"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
	models "gitlab.com/alphaticks/xchanger/models"
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

type InstrumentStatus int32

const (
	PreTrading   InstrumentStatus = 0
	Trading      InstrumentStatus = 1
	PostTrading  InstrumentStatus = 2
	EndOfDay     InstrumentStatus = 3
	Halt         InstrumentStatus = 4
	AuctionMatch InstrumentStatus = 5
	Break        InstrumentStatus = 6
	Disabled     InstrumentStatus = 7
)

var InstrumentStatus_name = map[int32]string{
	0: "PreTrading",
	1: "Trading",
	2: "PostTrading",
	3: "EndOfDay",
	4: "Halt",
	5: "AuctionMatch",
	6: "Break",
	7: "Disabled",
}

var InstrumentStatus_value = map[string]int32{
	"PreTrading":   0,
	"Trading":      1,
	"PostTrading":  2,
	"EndOfDay":     3,
	"Halt":         4,
	"AuctionMatch": 5,
	"Break":        6,
	"Disabled":     7,
}

func (InstrumentStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f22c063896c39066, []int{0}
}

type Side int32

const (
	Buy  Side = 0
	Sell Side = 1
)

var Side_name = map[int32]string{
	0: "Buy",
	1: "Sell",
}

var Side_value = map[string]int32{
	"Buy":  0,
	"Sell": 1,
}

func (Side) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f22c063896c39066, []int{1}
}

type OrderStatus int32

const (
	New             OrderStatus = 0
	PartiallyFilled OrderStatus = 1
	Filled          OrderStatus = 2
	Done            OrderStatus = 3
	Canceled        OrderStatus = 4
	Replaced        OrderStatus = 5
	PendingCancel   OrderStatus = 6
	Stopped         OrderStatus = 7
	Rejected        OrderStatus = 8
	Suspended       OrderStatus = 9
	PendingNew      OrderStatus = 10
	Calculated      OrderStatus = 11
	Expired         OrderStatus = 12
	PendingReplace  OrderStatus = 13
)

var OrderStatus_name = map[int32]string{
	0:  "New",
	1:  "PartiallyFilled",
	2:  "Filled",
	3:  "Done",
	4:  "Canceled",
	5:  "Replaced",
	6:  "PendingCancel",
	7:  "Stopped",
	8:  "Rejected",
	9:  "Suspended",
	10: "PendingNew",
	11: "Calculated",
	12: "Expired",
	13: "PendingReplace",
}

var OrderStatus_value = map[string]int32{
	"New":             0,
	"PartiallyFilled": 1,
	"Filled":          2,
	"Done":            3,
	"Canceled":        4,
	"Replaced":        5,
	"PendingCancel":   6,
	"Stopped":         7,
	"Rejected":        8,
	"Suspended":       9,
	"PendingNew":      10,
	"Calculated":      11,
	"Expired":         12,
	"PendingReplace":  13,
}

func (OrderStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f22c063896c39066, []int{2}
}

type OrderType int32

const (
	Market          OrderType = 0
	Limit           OrderType = 1
	Stop            OrderType = 2
	StopLimit       OrderType = 3
	LimitIfTouched  OrderType = 4
	MarketIfTouched OrderType = 5
)

var OrderType_name = map[int32]string{
	0: "Market",
	1: "Limit",
	2: "Stop",
	3: "StopLimit",
	4: "LimitIfTouched",
	5: "MarketIfTouched",
}

var OrderType_value = map[string]int32{
	"Market":          0,
	"Limit":           1,
	"Stop":            2,
	"StopLimit":       3,
	"LimitIfTouched":  4,
	"MarketIfTouched": 5,
}

func (OrderType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f22c063896c39066, []int{3}
}

type TimeInForce int32

const (
	Session           TimeInForce = 0
	GoodTillCancel    TimeInForce = 1
	AtTheOpening      TimeInForce = 2
	ImmediateOrCancel TimeInForce = 3
	FillOrKill        TimeInForce = 4
)

var TimeInForce_name = map[int32]string{
	0: "Session",
	1: "GoodTillCancel",
	2: "AtTheOpening",
	3: "ImmediateOrCancel",
	4: "FillOrKill",
}

var TimeInForce_value = map[string]int32{
	"Session":           0,
	"GoodTillCancel":    1,
	"AtTheOpening":      2,
	"ImmediateOrCancel": 3,
	"FillOrKill":        4,
}

func (TimeInForce) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f22c063896c39066, []int{4}
}

type Instrument struct {
	SecurityID *types.UInt64Value `protobuf:"bytes,1,opt,name=securityID,proto3" json:"securityID,omitempty"`
	Exchange   *models.Exchange   `protobuf:"bytes,2,opt,name=exchange,proto3" json:"exchange,omitempty"`
	Symbol     *types.StringValue `protobuf:"bytes,3,opt,name=symbol,proto3" json:"symbol,omitempty"`
}

func (m *Instrument) Reset()      { *m = Instrument{} }
func (*Instrument) ProtoMessage() {}
func (*Instrument) Descriptor() ([]byte, []int) {
	return fileDescriptor_f22c063896c39066, []int{0}
}
func (m *Instrument) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Instrument) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Instrument.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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

func (m *Instrument) GetSecurityID() *types.UInt64Value {
	if m != nil {
		return m.SecurityID
	}
	return nil
}

func (m *Instrument) GetExchange() *models.Exchange {
	if m != nil {
		return m.Exchange
	}
	return nil
}

func (m *Instrument) GetSymbol() *types.StringValue {
	if m != nil {
		return m.Symbol
	}
	return nil
}

type Security struct {
	SecurityID        uint64             `protobuf:"varint,1,opt,name=securityID,proto3" json:"securityID,omitempty"`
	SecurityType      string             `protobuf:"bytes,2,opt,name=securityType,proto3" json:"securityType,omitempty"`
	Exchange          *models.Exchange   `protobuf:"bytes,3,opt,name=exchange,proto3" json:"exchange,omitempty"`
	Symbol            string             `protobuf:"bytes,4,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Underlying        *models.Asset      `protobuf:"bytes,7,opt,name=underlying,proto3" json:"underlying,omitempty"`
	QuoteCurrency     *models.Asset      `protobuf:"bytes,8,opt,name=quote_currency,json=quoteCurrency,proto3" json:"quote_currency,omitempty"`
	Status            InstrumentStatus   `protobuf:"varint,9,opt,name=status,proto3,enum=models.InstrumentStatus" json:"status,omitempty"`
	IsInverse         bool               `protobuf:"varint,10,opt,name=is_inverse,json=isInverse,proto3" json:"is_inverse,omitempty"`
	MinPriceIncrement *types.DoubleValue `protobuf:"bytes,11,opt,name=min_price_increment,json=minPriceIncrement,proto3" json:"min_price_increment,omitempty"`
	RoundLot          *types.DoubleValue `protobuf:"bytes,12,opt,name=round_lot,json=roundLot,proto3" json:"round_lot,omitempty"`
	MakerFee          *types.DoubleValue `protobuf:"bytes,13,opt,name=maker_fee,json=makerFee,proto3" json:"maker_fee,omitempty"`
	TakerFee          *types.DoubleValue `protobuf:"bytes,14,opt,name=taker_fee,json=takerFee,proto3" json:"taker_fee,omitempty"`
	Multiplier        *types.DoubleValue `protobuf:"bytes,15,opt,name=multiplier,proto3" json:"multiplier,omitempty"`
	MaturityDate      *types.Timestamp   `protobuf:"bytes,16,opt,name=maturity_date,json=maturityDate,proto3" json:"maturity_date,omitempty"`
	SecuritySubType   *types.StringValue `protobuf:"bytes,17,opt,name=securitySubType,proto3" json:"securitySubType,omitempty"`
}

func (m *Security) Reset()      { *m = Security{} }
func (*Security) ProtoMessage() {}
func (*Security) Descriptor() ([]byte, []int) {
	return fileDescriptor_f22c063896c39066, []int{1}
}
func (m *Security) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Security) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Security.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Security) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Security.Merge(m, src)
}
func (m *Security) XXX_Size() int {
	return m.Size()
}
func (m *Security) XXX_DiscardUnknown() {
	xxx_messageInfo_Security.DiscardUnknown(m)
}

var xxx_messageInfo_Security proto.InternalMessageInfo

func (m *Security) GetSecurityID() uint64 {
	if m != nil {
		return m.SecurityID
	}
	return 0
}

func (m *Security) GetSecurityType() string {
	if m != nil {
		return m.SecurityType
	}
	return ""
}

func (m *Security) GetExchange() *models.Exchange {
	if m != nil {
		return m.Exchange
	}
	return nil
}

func (m *Security) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *Security) GetUnderlying() *models.Asset {
	if m != nil {
		return m.Underlying
	}
	return nil
}

func (m *Security) GetQuoteCurrency() *models.Asset {
	if m != nil {
		return m.QuoteCurrency
	}
	return nil
}

func (m *Security) GetStatus() InstrumentStatus {
	if m != nil {
		return m.Status
	}
	return PreTrading
}

func (m *Security) GetIsInverse() bool {
	if m != nil {
		return m.IsInverse
	}
	return false
}

func (m *Security) GetMinPriceIncrement() *types.DoubleValue {
	if m != nil {
		return m.MinPriceIncrement
	}
	return nil
}

func (m *Security) GetRoundLot() *types.DoubleValue {
	if m != nil {
		return m.RoundLot
	}
	return nil
}

func (m *Security) GetMakerFee() *types.DoubleValue {
	if m != nil {
		return m.MakerFee
	}
	return nil
}

func (m *Security) GetTakerFee() *types.DoubleValue {
	if m != nil {
		return m.TakerFee
	}
	return nil
}

func (m *Security) GetMultiplier() *types.DoubleValue {
	if m != nil {
		return m.Multiplier
	}
	return nil
}

func (m *Security) GetMaturityDate() *types.Timestamp {
	if m != nil {
		return m.MaturityDate
	}
	return nil
}

func (m *Security) GetSecuritySubType() *types.StringValue {
	if m != nil {
		return m.SecuritySubType
	}
	return nil
}

type Order struct {
	OrderID        string             `protobuf:"bytes,1,opt,name=orderID,proto3" json:"orderID,omitempty"`
	ClientOrderID  string             `protobuf:"bytes,2,opt,name=client_orderID,json=clientOrderID,proto3" json:"client_orderID,omitempty"`
	Instrument     *Instrument        `protobuf:"bytes,3,opt,name=instrument,proto3" json:"instrument,omitempty"`
	OrderStatus    OrderStatus        `protobuf:"varint,4,opt,name=order_status,json=orderStatus,proto3,enum=models.OrderStatus" json:"order_status,omitempty"`
	OrderType      OrderType          `protobuf:"varint,5,opt,name=order_type,json=orderType,proto3,enum=models.OrderType" json:"order_type,omitempty"`
	Side           Side               `protobuf:"varint,6,opt,name=side,proto3,enum=models.Side" json:"side,omitempty"`
	TimeInForce    TimeInForce        `protobuf:"varint,7,opt,name=time_in_force,json=timeInForce,proto3,enum=models.TimeInForce" json:"time_in_force,omitempty"`
	LeavesQuantity float64            `protobuf:"fixed64,8,opt,name=leaves_quantity,json=leavesQuantity,proto3" json:"leaves_quantity,omitempty"`
	CumQuantity    float64            `protobuf:"fixed64,9,opt,name=cum_quantity,json=cumQuantity,proto3" json:"cum_quantity,omitempty"`
	Price          *types.DoubleValue `protobuf:"bytes,10,opt,name=price,proto3" json:"price,omitempty"`
}

func (m *Order) Reset()      { *m = Order{} }
func (*Order) ProtoMessage() {}
func (*Order) Descriptor() ([]byte, []int) {
	return fileDescriptor_f22c063896c39066, []int{2}
}
func (m *Order) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Order) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Order.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Order) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Order.Merge(m, src)
}
func (m *Order) XXX_Size() int {
	return m.Size()
}
func (m *Order) XXX_DiscardUnknown() {
	xxx_messageInfo_Order.DiscardUnknown(m)
}

var xxx_messageInfo_Order proto.InternalMessageInfo

func (m *Order) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *Order) GetClientOrderID() string {
	if m != nil {
		return m.ClientOrderID
	}
	return ""
}

func (m *Order) GetInstrument() *Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *Order) GetOrderStatus() OrderStatus {
	if m != nil {
		return m.OrderStatus
	}
	return New
}

func (m *Order) GetOrderType() OrderType {
	if m != nil {
		return m.OrderType
	}
	return Market
}

func (m *Order) GetSide() Side {
	if m != nil {
		return m.Side
	}
	return Buy
}

func (m *Order) GetTimeInForce() TimeInForce {
	if m != nil {
		return m.TimeInForce
	}
	return Session
}

func (m *Order) GetLeavesQuantity() float64 {
	if m != nil {
		return m.LeavesQuantity
	}
	return 0
}

func (m *Order) GetCumQuantity() float64 {
	if m != nil {
		return m.CumQuantity
	}
	return 0
}

func (m *Order) GetPrice() *types.DoubleValue {
	if m != nil {
		return m.Price
	}
	return nil
}

type Position struct {
	AccountID  string      `protobuf:"bytes,1,opt,name=accountID,proto3" json:"accountID,omitempty"`
	Instrument *Instrument `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Quantity   float64     `protobuf:"fixed64,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Cross      bool        `protobuf:"varint,4,opt,name=cross,proto3" json:"cross,omitempty"`
	Cost       float64     `protobuf:"fixed64,5,opt,name=cost,proto3" json:"cost,omitempty"`
}

func (m *Position) Reset()      { *m = Position{} }
func (*Position) ProtoMessage() {}
func (*Position) Descriptor() ([]byte, []int) {
	return fileDescriptor_f22c063896c39066, []int{3}
}
func (m *Position) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Position) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Position.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Position) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Position.Merge(m, src)
}
func (m *Position) XXX_Size() int {
	return m.Size()
}
func (m *Position) XXX_DiscardUnknown() {
	xxx_messageInfo_Position.DiscardUnknown(m)
}

var xxx_messageInfo_Position proto.InternalMessageInfo

func (m *Position) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *Position) GetInstrument() *Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *Position) GetQuantity() float64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *Position) GetCross() bool {
	if m != nil {
		return m.Cross
	}
	return false
}

func (m *Position) GetCost() float64 {
	if m != nil {
		return m.Cost
	}
	return 0
}

type Balance struct {
	AccountID string        `protobuf:"bytes,1,opt,name=accountID,proto3" json:"accountID,omitempty"`
	Asset     *models.Asset `protobuf:"bytes,2,opt,name=asset,proto3" json:"asset,omitempty"`
	Quantity  float64       `protobuf:"fixed64,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
}

func (m *Balance) Reset()      { *m = Balance{} }
func (*Balance) ProtoMessage() {}
func (*Balance) Descriptor() ([]byte, []int) {
	return fileDescriptor_f22c063896c39066, []int{4}
}
func (m *Balance) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Balance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Balance.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Balance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Balance.Merge(m, src)
}
func (m *Balance) XXX_Size() int {
	return m.Size()
}
func (m *Balance) XXX_DiscardUnknown() {
	xxx_messageInfo_Balance.DiscardUnknown(m)
}

var xxx_messageInfo_Balance proto.InternalMessageInfo

func (m *Balance) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *Balance) GetAsset() *models.Asset {
	if m != nil {
		return m.Asset
	}
	return nil
}

func (m *Balance) GetQuantity() float64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func init() {
	proto.RegisterEnum("models.InstrumentStatus", InstrumentStatus_name, InstrumentStatus_value)
	proto.RegisterEnum("models.Side", Side_name, Side_value)
	proto.RegisterEnum("models.OrderStatus", OrderStatus_name, OrderStatus_value)
	proto.RegisterEnum("models.OrderType", OrderType_name, OrderType_value)
	proto.RegisterEnum("models.TimeInForce", TimeInForce_name, TimeInForce_value)
	proto.RegisterType((*Instrument)(nil), "models.Instrument")
	proto.RegisterType((*Security)(nil), "models.Security")
	proto.RegisterType((*Order)(nil), "models.Order")
	proto.RegisterType((*Position)(nil), "models.Position")
	proto.RegisterType((*Balance)(nil), "models.Balance")
}

func init() { proto.RegisterFile("security_data.proto", fileDescriptor_f22c063896c39066) }

var fileDescriptor_f22c063896c39066 = []byte{
	// 1182 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xcd, 0x6e, 0x1c, 0x45,
	0x17, 0x9d, 0xf6, 0xfc, 0x78, 0xe6, 0xce, 0x8f, 0xdb, 0xe5, 0xef, 0x43, 0x8d, 0x15, 0x1a, 0x63,
	0x84, 0x88, 0x2c, 0x18, 0x47, 0xc6, 0x0a, 0x42, 0x8a, 0x84, 0xe2, 0x38, 0x86, 0x11, 0x09, 0x36,
	0x3d, 0x86, 0xed, 0xa8, 0xa6, 0xfb, 0x7a, 0x5c, 0xb8, 0xba, 0xaa, 0x53, 0x5d, 0x9d, 0x64, 0x76,
	0x48, 0xbc, 0x40, 0x9e, 0x81, 0x15, 0x2b, 0xde, 0x81, 0x1d, 0xcb, 0x2c, 0x23, 0xb1, 0x21, 0x93,
	0x0d, 0xcb, 0x3c, 0x02, 0xaa, 0xea, 0xee, 0xf1, 0xd8, 0x21, 0xd8, 0xbb, 0xba, 0xb7, 0xce, 0x39,
	0x75, 0xaa, 0xee, 0xbd, 0x05, 0x6b, 0x29, 0x86, 0x99, 0x62, 0x7a, 0x3a, 0x8a, 0xa8, 0xa6, 0xfd,
	0x44, 0x49, 0x2d, 0x49, 0x23, 0x96, 0x11, 0xf2, 0x74, 0x7d, 0x67, 0xc2, 0x34, 0xa7, 0xe3, 0x7e,
	0x28, 0xe3, 0x6d, 0xca, 0x93, 0x53, 0xaa, 0x59, 0x78, 0x96, 0x6e, 0x3f, 0x0d, 0x4f, 0xa9, 0x98,
	0xa0, 0xda, 0xce, 0x61, 0xdb, 0x96, 0x94, 0xe6, 0xdc, 0x75, 0x7f, 0x22, 0xe5, 0x84, 0x63, 0x9e,
	0x1c, 0x67, 0x27, 0xdb, 0x4f, 0x14, 0x4d, 0x12, 0x54, 0xe5, 0xfe, 0xfb, 0x97, 0xf7, 0x35, 0x8b,
	0x31, 0xd5, 0x34, 0x4e, 0x72, 0xc0, 0xe6, 0x6f, 0x0e, 0xc0, 0x40, 0xa4, 0x5a, 0x65, 0x31, 0x0a,
	0x4d, 0xee, 0x00, 0x94, 0x16, 0x07, 0xfb, 0x9e, 0xb3, 0xe1, 0xdc, 0x6c, 0xef, 0xdc, 0xe8, 0xe7,
	0x22, 0xfd, 0x52, 0xa4, 0xff, 0xfd, 0x40, 0xe8, 0xdb, 0xbb, 0x3f, 0x50, 0x9e, 0x61, 0xb0, 0x80,
	0x27, 0x9f, 0x40, 0x13, 0x0b, 0xbb, 0xde, 0x92, 0xe5, 0xba, 0xfd, 0xdc, 0x75, 0xff, 0x7e, 0x91,
	0x0f, 0xe6, 0x08, 0xb2, 0x0b, 0x8d, 0x74, 0x1a, 0x8f, 0x25, 0xf7, 0xaa, 0x6f, 0x39, 0x67, 0xa8,
	0x15, 0x13, 0x93, 0xfc, 0x9c, 0x02, 0xbb, 0xf9, 0xac, 0x01, 0xcd, 0x61, 0x71, 0x24, 0xf1, 0xdf,
	0xb0, 0x5b, 0xbb, 0x60, 0x68, 0x13, 0x3a, 0x65, 0x74, 0x3c, 0x4d, 0x72, 0x53, 0xad, 0xe0, 0x42,
	0xee, 0x82, 0xe9, 0xea, 0x95, 0xa6, 0xdf, 0x99, 0x9b, 0xae, 0x59, 0xad, 0x22, 0x22, 0x9f, 0x02,
	0x64, 0x22, 0x42, 0xc5, 0xa7, 0x4c, 0x4c, 0xbc, 0x65, 0xab, 0xd3, 0x2d, 0x75, 0xee, 0xa6, 0x29,
	0xea, 0x60, 0x01, 0x40, 0x76, 0xa1, 0xf7, 0x28, 0x93, 0x1a, 0x47, 0x61, 0xa6, 0x14, 0x8a, 0x70,
	0xea, 0x35, 0xff, 0x8d, 0xd2, 0xb5, 0xa0, 0x7b, 0x05, 0x86, 0xdc, 0x82, 0x46, 0xaa, 0xa9, 0xce,
	0x52, 0xaf, 0xb5, 0xe1, 0xdc, 0xec, 0xed, 0x78, 0x25, 0xfa, 0xbc, 0x82, 0x43, 0xbb, 0x1f, 0x14,
	0x38, 0xf2, 0x1e, 0x00, 0x4b, 0x47, 0x4c, 0x3c, 0x46, 0x95, 0xa2, 0x07, 0x1b, 0xce, 0xcd, 0x66,
	0xd0, 0x62, 0xe9, 0x20, 0x4f, 0x90, 0x07, 0xb0, 0x16, 0x33, 0x31, 0x4a, 0x14, 0x0b, 0x71, 0xc4,
	0x44, 0xa8, 0xd0, 0x68, 0x78, 0xed, 0xb7, 0xd4, 0x63, 0x5f, 0x66, 0x63, 0x8e, 0x79, 0x3d, 0x56,
	0x63, 0x26, 0x8e, 0x0c, 0x6f, 0x50, 0xd2, 0xc8, 0x17, 0xd0, 0x52, 0x32, 0x13, 0xd1, 0x88, 0x4b,
	0xed, 0x75, 0xae, 0xa1, 0xd1, 0xb4, 0xf0, 0x07, 0xd2, 0x52, 0x63, 0x7a, 0x86, 0x6a, 0x74, 0x82,
	0xe8, 0x75, 0xaf, 0x43, 0xb5, 0xf0, 0x03, 0x44, 0x43, 0xd5, 0x73, 0x6a, 0xef, 0x3a, 0x54, 0x5d,
	0x52, 0xef, 0x00, 0xc4, 0x19, 0xd7, 0x2c, 0xe1, 0x0c, 0x95, 0xb7, 0x72, 0x0d, 0xee, 0x02, 0x9e,
	0x7c, 0x09, 0xdd, 0x98, 0xea, 0xf9, 0x38, 0xa3, 0xe7, 0x5a, 0x81, 0xf5, 0x37, 0x04, 0x8e, 0xcb,
	0x99, 0x0b, 0x3a, 0x25, 0x61, 0x9f, 0x6a, 0x24, 0x07, 0xb0, 0x52, 0x76, 0xe2, 0x30, 0x1b, 0xdb,
	0x06, 0x5d, 0xbd, 0xc6, 0x24, 0x5c, 0x26, 0x6d, 0xfe, 0x5e, 0x85, 0xfa, 0xa1, 0x8a, 0x50, 0x11,
	0x0f, 0x96, 0xa5, 0x59, 0x14, 0xc3, 0xd0, 0x0a, 0xca, 0x90, 0x7c, 0x04, 0xbd, 0x90, 0x33, 0x14,
	0x7a, 0x54, 0x02, 0xf2, 0x59, 0xe8, 0xe6, 0xd9, 0xc3, 0x02, 0xb6, 0x03, 0xc0, 0xe6, 0xbd, 0x54,
	0x8c, 0x03, 0x79, 0xb3, 0xcb, 0x82, 0x05, 0x14, 0xb9, 0x0d, 0x1d, 0xab, 0x39, 0x2a, 0x7a, 0xb3,
	0x66, 0x7b, 0x73, 0xad, 0x64, 0x59, 0xe9, 0xa2, 0x2d, 0xdb, 0xf2, 0x3c, 0x20, 0xb7, 0x00, 0x72,
	0x9e, 0x36, 0x37, 0xaf, 0x5b, 0xd6, 0xea, 0x05, 0x96, 0xb9, 0x5d, 0xd0, 0x92, 0xe5, 0x92, 0x6c,
	0x40, 0x2d, 0x65, 0x11, 0x7a, 0x0d, 0x8b, 0xed, 0x94, 0xd8, 0x21, 0x8b, 0x30, 0xb0, 0x3b, 0xe4,
	0x73, 0xe8, 0x9a, 0x1f, 0x6e, 0xc4, 0xc4, 0xe8, 0x44, 0xaa, 0x10, 0xed, 0x24, 0x2e, 0x98, 0x31,
	0xa5, 0x18, 0x88, 0x03, 0xb3, 0x15, 0xb4, 0xf5, 0x79, 0x40, 0x3e, 0x86, 0x15, 0x8e, 0xf4, 0x31,
	0xa6, 0xa3, 0x47, 0x19, 0x15, 0x9a, 0xe9, 0x7c, 0x22, 0x9d, 0xa0, 0x97, 0xa7, 0xbf, 0x2b, 0xb2,
	0xe4, 0x03, 0xe8, 0x84, 0x59, 0x7c, 0x8e, 0x6a, 0x59, 0x54, 0x3b, 0xcc, 0xe2, 0x39, 0x64, 0x07,
	0xea, 0x76, 0xa2, 0xec, 0xbc, 0x5d, 0xd5, 0x51, 0x39, 0x74, 0xf3, 0x17, 0x07, 0x9a, 0x47, 0x32,
	0x65, 0x9a, 0x49, 0x41, 0x6e, 0x40, 0x8b, 0x86, 0xa1, 0xcc, 0x84, 0x9e, 0x17, 0xf2, 0x3c, 0x71,
	0xa9, 0x46, 0x4b, 0xd7, 0xaa, 0xd1, 0x3a, 0x34, 0xe7, 0x8e, 0xab, 0xd6, 0xf1, 0x3c, 0x26, 0xff,
	0x83, 0x7a, 0xa8, 0x64, 0x9a, 0x17, 0xae, 0x19, 0xe4, 0x01, 0x21, 0x50, 0x0b, 0x65, 0xaa, 0x6d,
	0x5d, 0x9c, 0xc0, 0xae, 0x37, 0x4f, 0x61, 0x79, 0x8f, 0x72, 0x2a, 0x42, 0xbc, 0xc2, 0xe2, 0x87,
	0x50, 0xa7, 0xe6, 0x03, 0x2b, 0xdc, 0x5d, 0xfa, 0xd5, 0xf2, 0xbd, 0xff, 0xf2, 0xb4, 0xf5, 0xb3,
	0x03, 0xee, 0xe5, 0x4f, 0x8d, 0xf4, 0x00, 0x8e, 0x14, 0x1e, 0x2b, 0x1a, 0x31, 0x31, 0x71, 0x2b,
	0xa4, 0x0d, 0xcb, 0x65, 0xe0, 0x90, 0x15, 0x68, 0x1f, 0xc9, 0x54, 0x97, 0x89, 0x25, 0xd2, 0x81,
	0xe6, 0x7d, 0x11, 0x1d, 0x9e, 0xec, 0xd3, 0xa9, 0x5b, 0x25, 0x4d, 0xa8, 0x7d, 0x4d, 0xb9, 0x76,
	0x6b, 0xc4, 0x85, 0xce, 0xdd, 0x2c, 0x34, 0xef, 0xfc, 0x90, 0xea, 0xf0, 0xd4, 0xad, 0x93, 0x16,
	0xd4, 0xf7, 0x14, 0xd2, 0x33, 0xb7, 0x61, 0x48, 0xfb, 0x2c, 0xa5, 0x63, 0x8e, 0x91, 0xbb, 0xbc,
	0xf5, 0x2e, 0xd4, 0x4c, 0x6f, 0x91, 0x65, 0xa8, 0xee, 0x65, 0x53, 0xb7, 0x62, 0x54, 0x86, 0xc8,
	0xb9, 0xeb, 0x6c, 0xfd, 0xe9, 0x40, 0x7b, 0xa1, 0xb3, 0x0d, 0xe4, 0x5b, 0x7c, 0xe2, 0x56, 0xc8,
	0x1a, 0xac, 0x1c, 0x51, 0xa5, 0x19, 0xe5, 0x7c, 0x7a, 0xc0, 0xb8, 0x11, 0x72, 0x08, 0x40, 0xa3,
	0x58, 0x2f, 0x19, 0x8d, 0x7d, 0x29, 0xd0, 0xad, 0x9a, 0xc3, 0xee, 0x99, 0xc7, 0x34, 0xf9, 0x9a,
	0x89, 0x02, 0x4c, 0x38, 0x0d, 0x31, 0x72, 0xeb, 0x64, 0x15, 0xba, 0x47, 0x28, 0xcc, 0x55, 0x72,
	0x88, 0xdb, 0x30, 0xd7, 0x1d, 0x6a, 0x99, 0x24, 0xc6, 0x5a, 0x8e, 0xfe, 0x11, 0x43, 0x8d, 0x91,
	0xdb, 0x24, 0x5d, 0x68, 0x0d, 0xb3, 0x34, 0x41, 0x11, 0x61, 0xe4, 0xb6, 0xec, 0x43, 0xe5, 0x64,
	0xe3, 0x09, 0x4c, 0x7c, 0x8f, 0xf2, 0x30, 0xe3, 0xd4, 0xc0, 0xdb, 0x46, 0xe9, 0xfe, 0xd3, 0x84,
	0x29, 0x8c, 0xdc, 0x0e, 0x21, 0xd0, 0x2b, 0xc0, 0xc5, 0xf1, 0x6e, 0x77, 0x2b, 0x82, 0xd6, 0x7c,
	0x00, 0x8d, 0xf9, 0x87, 0x54, 0x9d, 0xa1, 0x76, 0x2b, 0xe6, 0xa9, 0x1e, 0xb0, 0x98, 0x69, 0xd7,
	0xb1, 0x6f, 0xa1, 0x65, 0xe2, 0x2e, 0xd9, 0xd3, 0xb5, 0x4c, 0xf2, 0x8d, 0xaa, 0x11, 0xb4, 0xcb,
	0xc1, 0xc9, 0xb1, 0xcc, 0xc2, 0x53, 0x7b, 0xb9, 0x35, 0x58, 0xc9, 0x35, 0xce, 0x93, 0xf5, 0xad,
	0x09, 0xb4, 0x17, 0xe6, 0xd1, 0xde, 0x0f, 0xd3, 0x94, 0x49, 0xe1, 0x56, 0x8c, 0xc8, 0x57, 0x52,
	0x46, 0xc7, 0x8c, 0xf3, 0xe2, 0x01, 0x1c, 0x5b, 0x39, 0x7d, 0x7c, 0x8a, 0x87, 0x09, 0x8a, 0xbc,
	0xc6, 0xff, 0x87, 0xd5, 0x41, 0x1c, 0x63, 0xc4, 0xa8, 0xc6, 0x43, 0x55, 0x00, 0xab, 0xe6, 0xbe,
	0xe6, 0xb9, 0x0f, 0xd5, 0x37, 0x8c, 0x73, 0xb7, 0xb6, 0xb7, 0xfb, 0xfc, 0xa5, 0x5f, 0x79, 0xf1,
	0xd2, 0xaf, 0xbc, 0x7e, 0xe9, 0x3b, 0x3f, 0xcd, 0x7c, 0xe7, 0xd7, 0x99, 0xef, 0xfc, 0x31, 0xf3,
	0x9d, 0xe7, 0x33, 0xdf, 0xf9, 0x6b, 0xe6, 0x3b, 0x7f, 0xcf, 0xfc, 0xca, 0xeb, 0x99, 0xef, 0x3c,
	0x7b, 0xe5, 0x57, 0x9e, 0xbf, 0xf2, 0x2b, 0x2f, 0x5e, 0xf9, 0x95, 0x71, 0xc3, 0xce, 0xeb, 0x67,
	0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0xb2, 0x39, 0x2a, 0x31, 0xb5, 0x09, 0x00, 0x00,
}

func (x InstrumentStatus) String() string {
	s, ok := InstrumentStatus_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x Side) String() string {
	s, ok := Side_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x OrderStatus) String() string {
	s, ok := OrderStatus_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x OrderType) String() string {
	s, ok := OrderType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x TimeInForce) String() string {
	s, ok := TimeInForce_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
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
	if !this.SecurityID.Equal(that1.SecurityID) {
		return false
	}
	if !this.Exchange.Equal(that1.Exchange) {
		return false
	}
	if !this.Symbol.Equal(that1.Symbol) {
		return false
	}
	return true
}
func (this *Security) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Security)
	if !ok {
		that2, ok := that.(Security)
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
	if this.SecurityType != that1.SecurityType {
		return false
	}
	if !this.Exchange.Equal(that1.Exchange) {
		return false
	}
	if this.Symbol != that1.Symbol {
		return false
	}
	if !this.Underlying.Equal(that1.Underlying) {
		return false
	}
	if !this.QuoteCurrency.Equal(that1.QuoteCurrency) {
		return false
	}
	if this.Status != that1.Status {
		return false
	}
	if this.IsInverse != that1.IsInverse {
		return false
	}
	if !this.MinPriceIncrement.Equal(that1.MinPriceIncrement) {
		return false
	}
	if !this.RoundLot.Equal(that1.RoundLot) {
		return false
	}
	if !this.MakerFee.Equal(that1.MakerFee) {
		return false
	}
	if !this.TakerFee.Equal(that1.TakerFee) {
		return false
	}
	if !this.Multiplier.Equal(that1.Multiplier) {
		return false
	}
	if !this.MaturityDate.Equal(that1.MaturityDate) {
		return false
	}
	if !this.SecuritySubType.Equal(that1.SecuritySubType) {
		return false
	}
	return true
}
func (this *Order) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Order)
	if !ok {
		that2, ok := that.(Order)
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
	if this.OrderID != that1.OrderID {
		return false
	}
	if this.ClientOrderID != that1.ClientOrderID {
		return false
	}
	if !this.Instrument.Equal(that1.Instrument) {
		return false
	}
	if this.OrderStatus != that1.OrderStatus {
		return false
	}
	if this.OrderType != that1.OrderType {
		return false
	}
	if this.Side != that1.Side {
		return false
	}
	if this.TimeInForce != that1.TimeInForce {
		return false
	}
	if this.LeavesQuantity != that1.LeavesQuantity {
		return false
	}
	if this.CumQuantity != that1.CumQuantity {
		return false
	}
	if !this.Price.Equal(that1.Price) {
		return false
	}
	return true
}
func (this *Position) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Position)
	if !ok {
		that2, ok := that.(Position)
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
	if this.AccountID != that1.AccountID {
		return false
	}
	if !this.Instrument.Equal(that1.Instrument) {
		return false
	}
	if this.Quantity != that1.Quantity {
		return false
	}
	if this.Cross != that1.Cross {
		return false
	}
	if this.Cost != that1.Cost {
		return false
	}
	return true
}
func (this *Balance) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Balance)
	if !ok {
		that2, ok := that.(Balance)
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
	if this.AccountID != that1.AccountID {
		return false
	}
	if !this.Asset.Equal(that1.Asset) {
		return false
	}
	if this.Quantity != that1.Quantity {
		return false
	}
	return true
}
func (this *Instrument) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&models.Instrument{")
	if this.SecurityID != nil {
		s = append(s, "SecurityID: "+fmt.Sprintf("%#v", this.SecurityID)+",\n")
	}
	if this.Exchange != nil {
		s = append(s, "Exchange: "+fmt.Sprintf("%#v", this.Exchange)+",\n")
	}
	if this.Symbol != nil {
		s = append(s, "Symbol: "+fmt.Sprintf("%#v", this.Symbol)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Security) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 19)
	s = append(s, "&models.Security{")
	s = append(s, "SecurityID: "+fmt.Sprintf("%#v", this.SecurityID)+",\n")
	s = append(s, "SecurityType: "+fmt.Sprintf("%#v", this.SecurityType)+",\n")
	if this.Exchange != nil {
		s = append(s, "Exchange: "+fmt.Sprintf("%#v", this.Exchange)+",\n")
	}
	s = append(s, "Symbol: "+fmt.Sprintf("%#v", this.Symbol)+",\n")
	if this.Underlying != nil {
		s = append(s, "Underlying: "+fmt.Sprintf("%#v", this.Underlying)+",\n")
	}
	if this.QuoteCurrency != nil {
		s = append(s, "QuoteCurrency: "+fmt.Sprintf("%#v", this.QuoteCurrency)+",\n")
	}
	s = append(s, "Status: "+fmt.Sprintf("%#v", this.Status)+",\n")
	s = append(s, "IsInverse: "+fmt.Sprintf("%#v", this.IsInverse)+",\n")
	if this.MinPriceIncrement != nil {
		s = append(s, "MinPriceIncrement: "+fmt.Sprintf("%#v", this.MinPriceIncrement)+",\n")
	}
	if this.RoundLot != nil {
		s = append(s, "RoundLot: "+fmt.Sprintf("%#v", this.RoundLot)+",\n")
	}
	if this.MakerFee != nil {
		s = append(s, "MakerFee: "+fmt.Sprintf("%#v", this.MakerFee)+",\n")
	}
	if this.TakerFee != nil {
		s = append(s, "TakerFee: "+fmt.Sprintf("%#v", this.TakerFee)+",\n")
	}
	if this.Multiplier != nil {
		s = append(s, "Multiplier: "+fmt.Sprintf("%#v", this.Multiplier)+",\n")
	}
	if this.MaturityDate != nil {
		s = append(s, "MaturityDate: "+fmt.Sprintf("%#v", this.MaturityDate)+",\n")
	}
	if this.SecuritySubType != nil {
		s = append(s, "SecuritySubType: "+fmt.Sprintf("%#v", this.SecuritySubType)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Order) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 14)
	s = append(s, "&models.Order{")
	s = append(s, "OrderID: "+fmt.Sprintf("%#v", this.OrderID)+",\n")
	s = append(s, "ClientOrderID: "+fmt.Sprintf("%#v", this.ClientOrderID)+",\n")
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	s = append(s, "OrderStatus: "+fmt.Sprintf("%#v", this.OrderStatus)+",\n")
	s = append(s, "OrderType: "+fmt.Sprintf("%#v", this.OrderType)+",\n")
	s = append(s, "Side: "+fmt.Sprintf("%#v", this.Side)+",\n")
	s = append(s, "TimeInForce: "+fmt.Sprintf("%#v", this.TimeInForce)+",\n")
	s = append(s, "LeavesQuantity: "+fmt.Sprintf("%#v", this.LeavesQuantity)+",\n")
	s = append(s, "CumQuantity: "+fmt.Sprintf("%#v", this.CumQuantity)+",\n")
	if this.Price != nil {
		s = append(s, "Price: "+fmt.Sprintf("%#v", this.Price)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Position) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&models.Position{")
	s = append(s, "AccountID: "+fmt.Sprintf("%#v", this.AccountID)+",\n")
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	s = append(s, "Quantity: "+fmt.Sprintf("%#v", this.Quantity)+",\n")
	s = append(s, "Cross: "+fmt.Sprintf("%#v", this.Cross)+",\n")
	s = append(s, "Cost: "+fmt.Sprintf("%#v", this.Cost)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Balance) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&models.Balance{")
	s = append(s, "AccountID: "+fmt.Sprintf("%#v", this.AccountID)+",\n")
	if this.Asset != nil {
		s = append(s, "Asset: "+fmt.Sprintf("%#v", this.Asset)+",\n")
	}
	s = append(s, "Quantity: "+fmt.Sprintf("%#v", this.Quantity)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringSecurityData(v interface{}, typ string) string {
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
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Instrument) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Instrument) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Symbol != nil {
		{
			size, err := m.Symbol.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Exchange != nil {
		{
			size, err := m.Exchange.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.SecurityID != nil {
		{
			size, err := m.SecurityID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Security) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Security) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Security) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SecuritySubType != nil {
		{
			size, err := m.SecuritySubType.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x8a
	}
	if m.MaturityDate != nil {
		{
			size, err := m.MaturityDate.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x82
	}
	if m.Multiplier != nil {
		{
			size, err := m.Multiplier.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x7a
	}
	if m.TakerFee != nil {
		{
			size, err := m.TakerFee.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x72
	}
	if m.MakerFee != nil {
		{
			size, err := m.MakerFee.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x6a
	}
	if m.RoundLot != nil {
		{
			size, err := m.RoundLot.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x62
	}
	if m.MinPriceIncrement != nil {
		{
			size, err := m.MinPriceIncrement.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x5a
	}
	if m.IsInverse {
		i--
		if m.IsInverse {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x50
	}
	if m.Status != 0 {
		i = encodeVarintSecurityData(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x48
	}
	if m.QuoteCurrency != nil {
		{
			size, err := m.QuoteCurrency.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x42
	}
	if m.Underlying != nil {
		{
			size, err := m.Underlying.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x3a
	}
	if len(m.Symbol) > 0 {
		i -= len(m.Symbol)
		copy(dAtA[i:], m.Symbol)
		i = encodeVarintSecurityData(dAtA, i, uint64(len(m.Symbol)))
		i--
		dAtA[i] = 0x22
	}
	if m.Exchange != nil {
		{
			size, err := m.Exchange.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.SecurityType) > 0 {
		i -= len(m.SecurityType)
		copy(dAtA[i:], m.SecurityType)
		i = encodeVarintSecurityData(dAtA, i, uint64(len(m.SecurityType)))
		i--
		dAtA[i] = 0x12
	}
	if m.SecurityID != 0 {
		i = encodeVarintSecurityData(dAtA, i, uint64(m.SecurityID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Order) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Order) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Order) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Price != nil {
		{
			size, err := m.Price.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x52
	}
	if m.CumQuantity != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.CumQuantity))))
		i--
		dAtA[i] = 0x49
	}
	if m.LeavesQuantity != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.LeavesQuantity))))
		i--
		dAtA[i] = 0x41
	}
	if m.TimeInForce != 0 {
		i = encodeVarintSecurityData(dAtA, i, uint64(m.TimeInForce))
		i--
		dAtA[i] = 0x38
	}
	if m.Side != 0 {
		i = encodeVarintSecurityData(dAtA, i, uint64(m.Side))
		i--
		dAtA[i] = 0x30
	}
	if m.OrderType != 0 {
		i = encodeVarintSecurityData(dAtA, i, uint64(m.OrderType))
		i--
		dAtA[i] = 0x28
	}
	if m.OrderStatus != 0 {
		i = encodeVarintSecurityData(dAtA, i, uint64(m.OrderStatus))
		i--
		dAtA[i] = 0x20
	}
	if m.Instrument != nil {
		{
			size, err := m.Instrument.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ClientOrderID) > 0 {
		i -= len(m.ClientOrderID)
		copy(dAtA[i:], m.ClientOrderID)
		i = encodeVarintSecurityData(dAtA, i, uint64(len(m.ClientOrderID)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.OrderID) > 0 {
		i -= len(m.OrderID)
		copy(dAtA[i:], m.OrderID)
		i = encodeVarintSecurityData(dAtA, i, uint64(len(m.OrderID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Position) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Position) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Position) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Cost != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Cost))))
		i--
		dAtA[i] = 0x29
	}
	if m.Cross {
		i--
		if m.Cross {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if m.Quantity != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Quantity))))
		i--
		dAtA[i] = 0x19
	}
	if m.Instrument != nil {
		{
			size, err := m.Instrument.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.AccountID) > 0 {
		i -= len(m.AccountID)
		copy(dAtA[i:], m.AccountID)
		i = encodeVarintSecurityData(dAtA, i, uint64(len(m.AccountID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Balance) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Balance) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Balance) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Quantity != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Quantity))))
		i--
		dAtA[i] = 0x19
	}
	if m.Asset != nil {
		{
			size, err := m.Asset.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.AccountID) > 0 {
		i -= len(m.AccountID)
		copy(dAtA[i:], m.AccountID)
		i = encodeVarintSecurityData(dAtA, i, uint64(len(m.AccountID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSecurityData(dAtA []byte, offset int, v uint64) int {
	offset -= sovSecurityData(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Instrument) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SecurityID != nil {
		l = m.SecurityID.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.Exchange != nil {
		l = m.Exchange.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.Symbol != nil {
		l = m.Symbol.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	return n
}

func (m *Security) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SecurityID != 0 {
		n += 1 + sovSecurityData(uint64(m.SecurityID))
	}
	l = len(m.SecurityType)
	if l > 0 {
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.Exchange != nil {
		l = m.Exchange.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	l = len(m.Symbol)
	if l > 0 {
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.Underlying != nil {
		l = m.Underlying.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.QuoteCurrency != nil {
		l = m.QuoteCurrency.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.Status != 0 {
		n += 1 + sovSecurityData(uint64(m.Status))
	}
	if m.IsInverse {
		n += 2
	}
	if m.MinPriceIncrement != nil {
		l = m.MinPriceIncrement.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.RoundLot != nil {
		l = m.RoundLot.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.MakerFee != nil {
		l = m.MakerFee.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.TakerFee != nil {
		l = m.TakerFee.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.Multiplier != nil {
		l = m.Multiplier.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.MaturityDate != nil {
		l = m.MaturityDate.Size()
		n += 2 + l + sovSecurityData(uint64(l))
	}
	if m.SecuritySubType != nil {
		l = m.SecuritySubType.Size()
		n += 2 + l + sovSecurityData(uint64(l))
	}
	return n
}

func (m *Order) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.OrderID)
	if l > 0 {
		n += 1 + l + sovSecurityData(uint64(l))
	}
	l = len(m.ClientOrderID)
	if l > 0 {
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.Instrument != nil {
		l = m.Instrument.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.OrderStatus != 0 {
		n += 1 + sovSecurityData(uint64(m.OrderStatus))
	}
	if m.OrderType != 0 {
		n += 1 + sovSecurityData(uint64(m.OrderType))
	}
	if m.Side != 0 {
		n += 1 + sovSecurityData(uint64(m.Side))
	}
	if m.TimeInForce != 0 {
		n += 1 + sovSecurityData(uint64(m.TimeInForce))
	}
	if m.LeavesQuantity != 0 {
		n += 9
	}
	if m.CumQuantity != 0 {
		n += 9
	}
	if m.Price != nil {
		l = m.Price.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	return n
}

func (m *Position) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.AccountID)
	if l > 0 {
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.Instrument != nil {
		l = m.Instrument.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.Quantity != 0 {
		n += 9
	}
	if m.Cross {
		n += 2
	}
	if m.Cost != 0 {
		n += 9
	}
	return n
}

func (m *Balance) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.AccountID)
	if l > 0 {
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.Asset != nil {
		l = m.Asset.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.Quantity != 0 {
		n += 9
	}
	return n
}

func sovSecurityData(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSecurityData(x uint64) (n int) {
	return sovSecurityData(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Instrument) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Instrument{`,
		`SecurityID:` + strings.Replace(fmt.Sprintf("%v", this.SecurityID), "UInt64Value", "types.UInt64Value", 1) + `,`,
		`Exchange:` + strings.Replace(fmt.Sprintf("%v", this.Exchange), "Exchange", "models.Exchange", 1) + `,`,
		`Symbol:` + strings.Replace(fmt.Sprintf("%v", this.Symbol), "StringValue", "types.StringValue", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Security) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Security{`,
		`SecurityID:` + fmt.Sprintf("%v", this.SecurityID) + `,`,
		`SecurityType:` + fmt.Sprintf("%v", this.SecurityType) + `,`,
		`Exchange:` + strings.Replace(fmt.Sprintf("%v", this.Exchange), "Exchange", "models.Exchange", 1) + `,`,
		`Symbol:` + fmt.Sprintf("%v", this.Symbol) + `,`,
		`Underlying:` + strings.Replace(fmt.Sprintf("%v", this.Underlying), "Asset", "models.Asset", 1) + `,`,
		`QuoteCurrency:` + strings.Replace(fmt.Sprintf("%v", this.QuoteCurrency), "Asset", "models.Asset", 1) + `,`,
		`Status:` + fmt.Sprintf("%v", this.Status) + `,`,
		`IsInverse:` + fmt.Sprintf("%v", this.IsInverse) + `,`,
		`MinPriceIncrement:` + strings.Replace(fmt.Sprintf("%v", this.MinPriceIncrement), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`RoundLot:` + strings.Replace(fmt.Sprintf("%v", this.RoundLot), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`MakerFee:` + strings.Replace(fmt.Sprintf("%v", this.MakerFee), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`TakerFee:` + strings.Replace(fmt.Sprintf("%v", this.TakerFee), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`Multiplier:` + strings.Replace(fmt.Sprintf("%v", this.Multiplier), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`MaturityDate:` + strings.Replace(fmt.Sprintf("%v", this.MaturityDate), "Timestamp", "types.Timestamp", 1) + `,`,
		`SecuritySubType:` + strings.Replace(fmt.Sprintf("%v", this.SecuritySubType), "StringValue", "types.StringValue", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Order) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Order{`,
		`OrderID:` + fmt.Sprintf("%v", this.OrderID) + `,`,
		`ClientOrderID:` + fmt.Sprintf("%v", this.ClientOrderID) + `,`,
		`Instrument:` + strings.Replace(this.Instrument.String(), "Instrument", "Instrument", 1) + `,`,
		`OrderStatus:` + fmt.Sprintf("%v", this.OrderStatus) + `,`,
		`OrderType:` + fmt.Sprintf("%v", this.OrderType) + `,`,
		`Side:` + fmt.Sprintf("%v", this.Side) + `,`,
		`TimeInForce:` + fmt.Sprintf("%v", this.TimeInForce) + `,`,
		`LeavesQuantity:` + fmt.Sprintf("%v", this.LeavesQuantity) + `,`,
		`CumQuantity:` + fmt.Sprintf("%v", this.CumQuantity) + `,`,
		`Price:` + strings.Replace(fmt.Sprintf("%v", this.Price), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Position) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Position{`,
		`AccountID:` + fmt.Sprintf("%v", this.AccountID) + `,`,
		`Instrument:` + strings.Replace(this.Instrument.String(), "Instrument", "Instrument", 1) + `,`,
		`Quantity:` + fmt.Sprintf("%v", this.Quantity) + `,`,
		`Cross:` + fmt.Sprintf("%v", this.Cross) + `,`,
		`Cost:` + fmt.Sprintf("%v", this.Cost) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Balance) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Balance{`,
		`AccountID:` + fmt.Sprintf("%v", this.AccountID) + `,`,
		`Asset:` + strings.Replace(fmt.Sprintf("%v", this.Asset), "Asset", "models.Asset", 1) + `,`,
		`Quantity:` + fmt.Sprintf("%v", this.Quantity) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringSecurityData(v interface{}) string {
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
				return ErrIntOverflowSecurityData
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SecurityID", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SecurityID == nil {
				m.SecurityID = &types.UInt64Value{}
			}
			if err := m.SecurityID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Exchange", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
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
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Symbol == nil {
				m.Symbol = &types.StringValue{}
			}
			if err := m.Symbol.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSecurityData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSecurityData
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
func (m *Security) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSecurityData
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
			return fmt.Errorf("proto: Security: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Security: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SecurityID", wireType)
			}
			m.SecurityID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return fmt.Errorf("proto: wrong wireType = %d for field SecurityType", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SecurityType = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Exchange", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
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
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Symbol", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Symbol = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Underlying", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Underlying == nil {
				m.Underlying = &models.Asset{}
			}
			if err := m.Underlying.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field QuoteCurrency", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.QuoteCurrency == nil {
				m.QuoteCurrency = &models.Asset{}
			}
			if err := m.QuoteCurrency.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= InstrumentStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsInverse", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
			m.IsInverse = bool(v != 0)
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinPriceIncrement", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.MinPriceIncrement == nil {
				m.MinPriceIncrement = &types.DoubleValue{}
			}
			if err := m.MinPriceIncrement.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RoundLot", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RoundLot == nil {
				m.RoundLot = &types.DoubleValue{}
			}
			if err := m.RoundLot.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MakerFee", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.MakerFee == nil {
				m.MakerFee = &types.DoubleValue{}
			}
			if err := m.MakerFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TakerFee", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.TakerFee == nil {
				m.TakerFee = &types.DoubleValue{}
			}
			if err := m.TakerFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 15:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Multiplier", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Multiplier == nil {
				m.Multiplier = &types.DoubleValue{}
			}
			if err := m.Multiplier.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 16:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaturityDate", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.MaturityDate == nil {
				m.MaturityDate = &types.Timestamp{}
			}
			if err := m.MaturityDate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 17:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SecuritySubType", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SecuritySubType == nil {
				m.SecuritySubType = &types.StringValue{}
			}
			if err := m.SecuritySubType.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSecurityData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSecurityData
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
func (m *Order) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSecurityData
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
			return fmt.Errorf("proto: Order: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Order: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OrderID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientOrderID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClientOrderID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Instrument", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
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
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderStatus", wireType)
			}
			m.OrderStatus = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrderStatus |= OrderStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderType", wireType)
			}
			m.OrderType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrderType |= OrderType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Side", wireType)
			}
			m.Side = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Side |= Side(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TimeInForce", wireType)
			}
			m.TimeInForce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TimeInForce |= TimeInForce(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field LeavesQuantity", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.LeavesQuantity = float64(math.Float64frombits(v))
		case 9:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field CumQuantity", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.CumQuantity = float64(math.Float64frombits(v))
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Price == nil {
				m.Price = &types.DoubleValue{}
			}
			if err := m.Price.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSecurityData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSecurityData
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
func (m *Position) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSecurityData
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
			return fmt.Errorf("proto: Position: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Position: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AccountID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AccountID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Instrument", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
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
		case 3:
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
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cross", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
			m.Cross = bool(v != 0)
		case 5:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cost", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Cost = float64(math.Float64frombits(v))
		default:
			iNdEx = preIndex
			skippy, err := skipSecurityData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSecurityData
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
func (m *Balance) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSecurityData
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
			return fmt.Errorf("proto: Balance: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Balance: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AccountID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AccountID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asset", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
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
				return ErrInvalidLengthSecurityData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSecurityData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Asset == nil {
				m.Asset = &models.Asset{}
			}
			if err := m.Asset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
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
			skippy, err := skipSecurityData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSecurityData
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
func skipSecurityData(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSecurityData
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
					return 0, ErrIntOverflowSecurityData
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
					return 0, ErrIntOverflowSecurityData
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
				return 0, ErrInvalidLengthSecurityData
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSecurityData
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSecurityData
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSecurityData        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSecurityData          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSecurityData = fmt.Errorf("proto: unexpected end of group")
)
