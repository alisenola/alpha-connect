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

type ExecutionInstruction int32

const (
	ParticipateDoNotInitiate ExecutionInstruction = 0
	ReduceOnly               ExecutionInstruction = 1
)

var ExecutionInstruction_name = map[int32]string{
	0: "ParticipateDoNotInitiate",
	1: "ReduceOnly",
}

var ExecutionInstruction_value = map[string]int32{
	"ParticipateDoNotInitiate": 0,
	"ReduceOnly":               1,
}

func (ExecutionInstruction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f22c063896c39066, []int{1}
}

type TradeType int32

const (
	Regular             TradeType = 0
	Block               TradeType = 1
	ExchangeForPhysical TradeType = 2
	Transfer            TradeType = 3
)

var TradeType_name = map[int32]string{
	0: "Regular",
	1: "Block",
	2: "ExchangeForPhysical",
	3: "Transfer",
}

var TradeType_value = map[string]int32{
	"Regular":             0,
	"Block":               1,
	"ExchangeForPhysical": 2,
	"Transfer":            3,
}

func (TradeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f22c063896c39066, []int{2}
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
	return fileDescriptor_f22c063896c39066, []int{3}
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
	return fileDescriptor_f22c063896c39066, []int{4}
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
	return fileDescriptor_f22c063896c39066, []int{5}
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
	return fileDescriptor_f22c063896c39066, []int{6}
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
	MaxLimitQuantity  *types.DoubleValue `protobuf:"bytes,18,opt,name=max_limit_quantity,json=maxLimitQuantity,proto3" json:"max_limit_quantity,omitempty"`
	MaxMarketQuantity *types.DoubleValue `protobuf:"bytes,19,opt,name=max_market_quantity,json=maxMarketQuantity,proto3" json:"max_market_quantity,omitempty"`
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

func (m *Security) GetMaxLimitQuantity() *types.DoubleValue {
	if m != nil {
		return m.MaxLimitQuantity
	}
	return nil
}

func (m *Security) GetMaxMarketQuantity() *types.DoubleValue {
	if m != nil {
		return m.MaxMarketQuantity
	}
	return nil
}

type Order struct {
	OrderID               string                 `protobuf:"bytes,1,opt,name=orderID,proto3" json:"orderID,omitempty"`
	ClientOrderID         string                 `protobuf:"bytes,2,opt,name=client_orderID,json=clientOrderID,proto3" json:"client_orderID,omitempty"`
	Instrument            *Instrument            `protobuf:"bytes,3,opt,name=instrument,proto3" json:"instrument,omitempty"`
	OrderStatus           OrderStatus            `protobuf:"varint,4,opt,name=order_status,json=orderStatus,proto3,enum=models.OrderStatus" json:"order_status,omitempty"`
	OrderType             OrderType              `protobuf:"varint,5,opt,name=order_type,json=orderType,proto3,enum=models.OrderType" json:"order_type,omitempty"`
	Side                  Side                   `protobuf:"varint,6,opt,name=side,proto3,enum=models.Side" json:"side,omitempty"`
	TimeInForce           TimeInForce            `protobuf:"varint,7,opt,name=time_in_force,json=timeInForce,proto3,enum=models.TimeInForce" json:"time_in_force,omitempty"`
	LeavesQuantity        float64                `protobuf:"fixed64,8,opt,name=leaves_quantity,json=leavesQuantity,proto3" json:"leaves_quantity,omitempty"`
	CumQuantity           float64                `protobuf:"fixed64,9,opt,name=cum_quantity,json=cumQuantity,proto3" json:"cum_quantity,omitempty"`
	Price                 *types.DoubleValue     `protobuf:"bytes,10,opt,name=price,proto3" json:"price,omitempty"`
	ExecutionInstructions []ExecutionInstruction `protobuf:"varint,11,rep,packed,name=execution_instructions,json=executionInstructions,proto3,enum=models.ExecutionInstruction" json:"execution_instructions,omitempty"`
	Tag                   string                 `protobuf:"bytes,12,opt,name=tag,proto3" json:"tag,omitempty"`
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

func (m *Order) GetExecutionInstructions() []ExecutionInstruction {
	if m != nil {
		return m.ExecutionInstructions
	}
	return nil
}

func (m *Order) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

type Position struct {
	Account    string      `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
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

func (m *Position) GetAccount() string {
	if m != nil {
		return m.Account
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
	Account  string        `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	Asset    *models.Asset `protobuf:"bytes,2,opt,name=asset,proto3" json:"asset,omitempty"`
	Quantity float64       `protobuf:"fixed64,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
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

func (m *Balance) GetAccount() string {
	if m != nil {
		return m.Account
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

type TradeCapture struct {
	Side            Side               `protobuf:"varint,1,opt,name=side,proto3,enum=models.Side" json:"side,omitempty"`
	Type            TradeType          `protobuf:"varint,2,opt,name=type,proto3,enum=models.TradeType" json:"type,omitempty"`
	Price           float64            `protobuf:"fixed64,3,opt,name=price,proto3" json:"price,omitempty"`
	Quantity        float64            `protobuf:"fixed64,4,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Commission      float64            `protobuf:"fixed64,5,opt,name=commission,proto3" json:"commission,omitempty"`
	CommissionAsset *models.Asset      `protobuf:"bytes,6,opt,name=commission_asset,json=commissionAsset,proto3" json:"commission_asset,omitempty"`
	TradeID         string             `protobuf:"bytes,7,opt,name=tradeID,proto3" json:"tradeID,omitempty"`
	Instrument      *Instrument        `protobuf:"bytes,8,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Trade_LinkID    *types.StringValue `protobuf:"bytes,9,opt,name=trade_LinkID,json=tradeLinkID,proto3" json:"trade_LinkID,omitempty"`
	OrderID         *types.StringValue `protobuf:"bytes,10,opt,name=orderID,proto3" json:"orderID,omitempty"`
	ClientOrderID   *types.StringValue `protobuf:"bytes,11,opt,name=client_orderID,json=clientOrderID,proto3" json:"client_orderID,omitempty"`
	TransactionTime *types.Timestamp   `protobuf:"bytes,12,opt,name=transaction_time,json=transactionTime,proto3" json:"transaction_time,omitempty"`
}

func (m *TradeCapture) Reset()      { *m = TradeCapture{} }
func (*TradeCapture) ProtoMessage() {}
func (*TradeCapture) Descriptor() ([]byte, []int) {
	return fileDescriptor_f22c063896c39066, []int{5}
}
func (m *TradeCapture) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TradeCapture) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TradeCapture.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TradeCapture) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TradeCapture.Merge(m, src)
}
func (m *TradeCapture) XXX_Size() int {
	return m.Size()
}
func (m *TradeCapture) XXX_DiscardUnknown() {
	xxx_messageInfo_TradeCapture.DiscardUnknown(m)
}

var xxx_messageInfo_TradeCapture proto.InternalMessageInfo

func (m *TradeCapture) GetSide() Side {
	if m != nil {
		return m.Side
	}
	return Buy
}

func (m *TradeCapture) GetType() TradeType {
	if m != nil {
		return m.Type
	}
	return Regular
}

func (m *TradeCapture) GetPrice() float64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *TradeCapture) GetQuantity() float64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *TradeCapture) GetCommission() float64 {
	if m != nil {
		return m.Commission
	}
	return 0
}

func (m *TradeCapture) GetCommissionAsset() *models.Asset {
	if m != nil {
		return m.CommissionAsset
	}
	return nil
}

func (m *TradeCapture) GetTradeID() string {
	if m != nil {
		return m.TradeID
	}
	return ""
}

func (m *TradeCapture) GetInstrument() *Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *TradeCapture) GetTrade_LinkID() *types.StringValue {
	if m != nil {
		return m.Trade_LinkID
	}
	return nil
}

func (m *TradeCapture) GetOrderID() *types.StringValue {
	if m != nil {
		return m.OrderID
	}
	return nil
}

func (m *TradeCapture) GetClientOrderID() *types.StringValue {
	if m != nil {
		return m.ClientOrderID
	}
	return nil
}

func (m *TradeCapture) GetTransactionTime() *types.Timestamp {
	if m != nil {
		return m.TransactionTime
	}
	return nil
}

func init() {
	proto.RegisterEnum("models.InstrumentStatus", InstrumentStatus_name, InstrumentStatus_value)
	proto.RegisterEnum("models.ExecutionInstruction", ExecutionInstruction_name, ExecutionInstruction_value)
	proto.RegisterEnum("models.TradeType", TradeType_name, TradeType_value)
	proto.RegisterEnum("models.Side", Side_name, Side_value)
	proto.RegisterEnum("models.OrderStatus", OrderStatus_name, OrderStatus_value)
	proto.RegisterEnum("models.OrderType", OrderType_name, OrderType_value)
	proto.RegisterEnum("models.TimeInForce", TimeInForce_name, TimeInForce_value)
	proto.RegisterType((*Instrument)(nil), "models.Instrument")
	proto.RegisterType((*Security)(nil), "models.Security")
	proto.RegisterType((*Order)(nil), "models.Order")
	proto.RegisterType((*Position)(nil), "models.Position")
	proto.RegisterType((*Balance)(nil), "models.Balance")
	proto.RegisterType((*TradeCapture)(nil), "models.TradeCapture")
}

func init() { proto.RegisterFile("security_data.proto", fileDescriptor_f22c063896c39066) }

var fileDescriptor_f22c063896c39066 = []byte{
	// 1501 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x56, 0x4f, 0x6f, 0xe4, 0x48,
	0x15, 0x6f, 0xa7, 0xff, 0xa4, 0xfb, 0xf5, 0xbf, 0x4a, 0x65, 0x77, 0x31, 0xd1, 0x60, 0x42, 0xd0,
	0x8a, 0x28, 0x82, 0x64, 0x15, 0xa2, 0x01, 0xa4, 0x95, 0x56, 0x9b, 0x74, 0x02, 0xbd, 0x64, 0x27,
	0xc1, 0x1d, 0xb8, 0x5a, 0x15, 0xbb, 0xd2, 0x29, 0x52, 0x76, 0x79, 0xca, 0xe5, 0x99, 0xf4, 0x0d,
	0x89, 0x2f, 0xc0, 0x27, 0xe0, 0x3c, 0xe2, 0xc0, 0x91, 0xcf, 0xc0, 0x71, 0x8e, 0x23, 0x71, 0x61,
	0x32, 0x17, 0x8e, 0xf3, 0x11, 0x50, 0x55, 0xd9, 0xdd, 0x9d, 0x4e, 0x66, 0xa6, 0x6f, 0x7e, 0xaf,
	0x7e, 0xbf, 0x57, 0xef, 0x7f, 0x19, 0xd6, 0x33, 0x1a, 0xe6, 0x92, 0xa9, 0x49, 0x10, 0x11, 0x45,
	0x76, 0x53, 0x29, 0x94, 0xc0, 0x8d, 0x58, 0x44, 0x94, 0x67, 0x1b, 0xfb, 0x63, 0xa6, 0x38, 0xb9,
	0xdc, 0x0d, 0x45, 0xbc, 0x47, 0x78, 0x7a, 0x4d, 0x14, 0x0b, 0x6f, 0xb2, 0xbd, 0xdb, 0xf0, 0x9a,
	0x24, 0x63, 0x2a, 0xf7, 0x2c, 0x6c, 0xcf, 0x90, 0x32, 0xcb, 0xdd, 0xf0, 0xc6, 0x42, 0x8c, 0x39,
	0xb5, 0xca, 0xcb, 0xfc, 0x6a, 0xef, 0xa5, 0x24, 0x69, 0x4a, 0x65, 0x79, 0xfe, 0xe3, 0xc5, 0x73,
	0xc5, 0x62, 0x9a, 0x29, 0x12, 0xa7, 0x16, 0xb0, 0xf5, 0x4f, 0x07, 0x60, 0x98, 0x64, 0x4a, 0xe6,
	0x31, 0x4d, 0x14, 0xfe, 0x1a, 0xa0, 0x74, 0x71, 0x38, 0x70, 0x9d, 0x4d, 0x67, 0xbb, 0xbd, 0xff,
	0x64, 0xd7, 0x1a, 0xd9, 0x2d, 0x8d, 0xec, 0xfe, 0x71, 0x98, 0xa8, 0xa7, 0x07, 0x7f, 0x22, 0x3c,
	0xa7, 0xfe, 0x1c, 0x1e, 0xff, 0x1c, 0x9a, 0xb4, 0x70, 0xd7, 0x5d, 0x31, 0x5c, 0xb4, 0x6b, 0xbd,
	0xde, 0x3d, 0x2e, 0xf4, 0xfe, 0x14, 0x81, 0x0f, 0xa0, 0x91, 0x4d, 0xe2, 0x4b, 0xc1, 0xdd, 0xea,
	0x07, 0xee, 0x19, 0x29, 0xc9, 0x92, 0xb1, 0xbd, 0xa7, 0xc0, 0x6e, 0xfd, 0x63, 0x15, 0x9a, 0xa3,
	0xe2, 0x4a, 0xec, 0x3d, 0x70, 0xb7, 0x76, 0xcf, 0xa1, 0x2d, 0xe8, 0x94, 0xd2, 0xc5, 0x24, 0xb5,
	0x4e, 0xb5, 0xfc, 0x7b, 0xba, 0x7b, 0x4e, 0x57, 0x3f, 0xe9, 0xf4, 0x17, 0x53, 0xa7, 0x6b, 0xc6,
	0x56, 0x21, 0xe1, 0x5f, 0x00, 0xe4, 0x49, 0x44, 0x25, 0x9f, 0xb0, 0x64, 0xec, 0xae, 0x1a, 0x3b,
	0xdd, 0xd2, 0xce, 0xb7, 0x59, 0x46, 0x95, 0x3f, 0x07, 0xc0, 0x07, 0xd0, 0x7b, 0x9e, 0x0b, 0x45,
	0x83, 0x30, 0x97, 0x92, 0x26, 0xe1, 0xc4, 0x6d, 0x3e, 0x46, 0xe9, 0x1a, 0xd0, 0x51, 0x81, 0xc1,
	0x5f, 0x41, 0x23, 0x53, 0x44, 0xe5, 0x99, 0xdb, 0xda, 0x74, 0xb6, 0x7b, 0xfb, 0x6e, 0x89, 0x9e,
	0x55, 0x70, 0x64, 0xce, 0xfd, 0x02, 0x87, 0x7f, 0x04, 0xc0, 0xb2, 0x80, 0x25, 0x2f, 0xa8, 0xcc,
	0xa8, 0x0b, 0x9b, 0xce, 0x76, 0xd3, 0x6f, 0xb1, 0x6c, 0x68, 0x15, 0xf8, 0x14, 0xd6, 0x63, 0x96,
	0x04, 0xa9, 0x64, 0x21, 0x0d, 0x58, 0x12, 0x4a, 0xaa, 0x6d, 0xb8, 0xed, 0x0f, 0xd4, 0x63, 0x20,
	0xf2, 0x4b, 0x4e, 0x6d, 0x3d, 0xd6, 0x62, 0x96, 0x9c, 0x6b, 0xde, 0xb0, 0xa4, 0xe1, 0xdf, 0x40,
	0x4b, 0x8a, 0x3c, 0x89, 0x02, 0x2e, 0x94, 0xdb, 0x59, 0xc2, 0x46, 0xd3, 0xc0, 0x4f, 0x85, 0xa1,
	0xc6, 0xe4, 0x86, 0xca, 0xe0, 0x8a, 0x52, 0xb7, 0xbb, 0x0c, 0xd5, 0xc0, 0x4f, 0x28, 0xd5, 0x54,
	0x35, 0xa5, 0xf6, 0x96, 0xa1, 0xaa, 0x92, 0xfa, 0x35, 0x40, 0x9c, 0x73, 0xc5, 0x52, 0xce, 0xa8,
	0x74, 0xfb, 0x4b, 0x70, 0xe7, 0xf0, 0xf8, 0x1b, 0xe8, 0xc6, 0x44, 0x4d, 0xc7, 0x99, 0xba, 0xc8,
	0x18, 0xd8, 0x78, 0x60, 0xe0, 0xa2, 0x9c, 0x39, 0xbf, 0x53, 0x12, 0x06, 0x44, 0x51, 0x7c, 0x02,
	0xfd, 0xb2, 0x13, 0x47, 0xf9, 0xa5, 0x69, 0xd0, 0xb5, 0x25, 0x26, 0x61, 0x91, 0x84, 0xbf, 0x03,
	0x1c, 0x93, 0xdb, 0x80, 0xb3, 0x98, 0xa9, 0xe0, 0x79, 0x4e, 0x12, 0xc5, 0xd4, 0xc4, 0xc5, 0x4b,
	0x84, 0x83, 0x62, 0x72, 0x7b, 0xaa, 0x69, 0x7f, 0x28, 0x58, 0xa6, 0x23, 0xc8, 0x6d, 0x10, 0x13,
	0x79, 0x43, 0xe7, 0x8c, 0xad, 0x2f, 0xd5, 0x11, 0xe4, 0xf6, 0x7b, 0xc3, 0x2b, 0xad, 0x6d, 0xbd,
	0xaa, 0x41, 0xfd, 0x4c, 0x46, 0x54, 0x62, 0x17, 0x56, 0x85, 0xfe, 0x28, 0xc6, 0xb4, 0xe5, 0x97,
	0x22, 0xfe, 0x12, 0x7a, 0x21, 0x67, 0x34, 0x51, 0x41, 0x09, 0xb0, 0x53, 0xda, 0xb5, 0xda, 0xb3,
	0x02, 0xb6, 0x0f, 0xc0, 0xa6, 0x5d, 0x5e, 0x0c, 0x2a, 0x7e, 0xd8, 0xff, 0xfe, 0x1c, 0x0a, 0x3f,
	0x85, 0x8e, 0xb1, 0x19, 0x14, 0x53, 0x53, 0x33, 0x53, 0xb3, 0x5e, 0xb2, 0x8c, 0xe9, 0x62, 0x60,
	0xda, 0x62, 0x26, 0xe0, 0xaf, 0x00, 0x2c, 0x4f, 0xe9, 0x9a, 0xd4, 0x0d, 0x6b, 0xed, 0x1e, 0x4b,
	0xe7, 0xdd, 0x6f, 0x89, 0xf2, 0x13, 0x6f, 0x42, 0x2d, 0x63, 0x11, 0x75, 0x1b, 0x06, 0xdb, 0x29,
	0xb1, 0x23, 0x16, 0x51, 0xdf, 0x9c, 0xe0, 0x5f, 0x41, 0x57, 0xef, 0xde, 0x80, 0x25, 0xc1, 0x95,
	0x90, 0x21, 0x35, 0x3b, 0x62, 0xce, 0x19, 0xdd, 0x24, 0xc3, 0xe4, 0x44, 0x1f, 0xf9, 0x6d, 0x35,
	0x13, 0xf0, 0xcf, 0xa0, 0xcf, 0x29, 0x79, 0x41, 0xb3, 0x59, 0x35, 0xf4, 0xae, 0x70, 0xfc, 0x9e,
	0x55, 0x4f, 0x4b, 0xf7, 0x13, 0xe8, 0x84, 0x79, 0x3c, 0x43, 0xb5, 0x0c, 0xaa, 0x1d, 0xe6, 0xf1,
	0x14, 0xb2, 0x0f, 0x75, 0x33, 0xeb, 0x66, 0x13, 0x7c, 0xaa, 0x9e, 0x16, 0x8a, 0x47, 0xf0, 0x05,
	0xbd, 0xa5, 0x61, 0xae, 0x98, 0x48, 0x02, 0x9b, 0xdc, 0x50, 0x7f, 0x67, 0x6e, 0x7b, 0xb3, 0xba,
	0xdd, 0xdb, 0x7f, 0x32, 0xdb, 0x96, 0x05, 0x6a, 0x38, 0x03, 0xf9, 0x9f, 0xd3, 0x47, 0xb4, 0x19,
	0x46, 0x50, 0x55, 0x64, 0x6c, 0x96, 0x44, 0xcb, 0xd7, 0x9f, 0x5b, 0x7f, 0x77, 0xa0, 0x79, 0x2e,
	0x32, 0xa6, 0xcf, 0x75, 0xb7, 0x90, 0x30, 0x14, 0x79, 0xa2, 0xca, 0x6e, 0x29, 0xc4, 0x85, 0x36,
	0x58, 0x59, 0xaa, 0x0d, 0x36, 0xa0, 0x39, 0x4d, 0x4a, 0xd5, 0x24, 0x65, 0x2a, 0xe3, 0xcf, 0xa0,
	0x1e, 0x4a, 0x91, 0xd9, 0xde, 0x68, 0xfa, 0x56, 0xc0, 0x18, 0x6a, 0xa1, 0xc8, 0x94, 0x29, 0xbd,
	0xe3, 0x9b, 0xef, 0xad, 0x08, 0x56, 0x0f, 0x09, 0x27, 0x49, 0x48, 0x3f, 0xe2, 0xde, 0x4f, 0xa1,
	0x4e, 0xf4, 0xe6, 0x2e, 0x3c, 0x5b, 0x58, 0xe7, 0xf6, 0xec, 0x63, 0xfe, 0x6c, 0xfd, 0xab, 0x06,
	0x9d, 0x0b, 0x49, 0x22, 0x7a, 0x44, 0x52, 0x95, 0xcb, 0x59, 0x67, 0x39, 0x1f, 0xec, 0xac, 0x2f,
	0xa1, 0xa6, 0xca, 0xc7, 0x6d, 0xae, 0x4f, 0x8d, 0x15, 0xd3, 0xa7, 0xe6, 0x58, 0x47, 0x6a, 0x6b,
	0x6f, 0xaf, 0x2c, 0xaa, 0x3b, 0xef, 0x4b, 0x6d, 0x21, 0x37, 0x1e, 0x40, 0x28, 0xe2, 0x98, 0x65,
	0x19, 0x13, 0x49, 0x91, 0x8b, 0x39, 0x0d, 0xfe, 0x35, 0xa0, 0x99, 0x14, 0xd8, 0xb8, 0x1b, 0x8f,
	0xc5, 0xdd, 0x9f, 0xc1, 0x8c, 0x42, 0x27, 0x50, 0x69, 0xf7, 0x86, 0x03, 0x33, 0x06, 0x2d, 0xbf,
	0x14, 0x17, 0xea, 0xdb, 0x5c, 0xaa, 0xbe, 0xdf, 0x40, 0xc7, 0xd0, 0x83, 0x53, 0x96, 0xdc, 0x0c,
	0x07, 0xa6, 0xf1, 0x3f, 0xb5, 0x44, 0xdb, 0x86, 0x61, 0x09, 0xf8, 0xe9, 0x6c, 0x39, 0xc1, 0x12,
	0xdc, 0xe9, 0xea, 0x3a, 0x7a, 0xb0, 0xba, 0xda, 0x4b, 0xd0, 0x17, 0x16, 0xdb, 0x31, 0x20, 0x25,
	0x49, 0x92, 0x11, 0x33, 0x1a, 0x81, 0x1e, 0xfd, 0xe2, 0xf1, 0xfc, 0xd8, 0x4b, 0xd2, 0x9f, 0xe3,
	0x68, 0xed, 0xce, 0x5f, 0x1d, 0x40, 0x8b, 0xbf, 0x01, 0xb8, 0x07, 0x70, 0x2e, 0xa9, 0xee, 0x04,
	0x96, 0x8c, 0x51, 0x05, 0xb7, 0x61, 0xb5, 0x14, 0x1c, 0xdc, 0x87, 0xf6, 0xb9, 0xc8, 0x54, 0xa9,
	0x58, 0xc1, 0x1d, 0x68, 0x1e, 0x27, 0xd1, 0xd9, 0xd5, 0x80, 0x4c, 0x50, 0x15, 0x37, 0xa1, 0xf6,
	0x3b, 0xc2, 0x15, 0xaa, 0x61, 0x04, 0x9d, 0x6f, 0xed, 0xe0, 0x7e, 0x4f, 0x54, 0x78, 0x8d, 0xea,
	0xb8, 0x05, 0xf5, 0x43, 0x49, 0xc9, 0x0d, 0x6a, 0x68, 0xd2, 0x80, 0x65, 0xe4, 0x92, 0xd3, 0x08,
	0xad, 0xee, 0x0c, 0xe0, 0xb3, 0xc7, 0xd6, 0x00, 0x7e, 0x02, 0xee, 0x39, 0x91, 0x8a, 0x85, 0x2c,
	0x25, 0x8a, 0x0e, 0xc4, 0x33, 0xa1, 0x86, 0x09, 0x53, 0x8c, 0x28, 0x8a, 0x2a, 0xda, 0x4d, 0x9f,
	0x46, 0x79, 0x48, 0xcf, 0x12, 0x3e, 0x41, 0xce, 0xce, 0x77, 0xd0, 0x9a, 0x76, 0xaf, 0xf6, 0xd9,
	0xa7, 0xe3, 0x9c, 0x13, 0x89, 0x2a, 0xe6, 0x62, 0x2e, 0xc2, 0x1b, 0xe4, 0xe0, 0x1f, 0xc0, 0x7a,
	0xf9, 0x7f, 0x76, 0x22, 0xe4, 0xf9, 0xf5, 0x24, 0x63, 0x21, 0xe1, 0x36, 0x8c, 0x0b, 0x9d, 0x9c,
	0x2b, 0x2a, 0x51, 0x75, 0xe7, 0x87, 0x50, 0xd3, 0xb3, 0x82, 0x57, 0xa1, 0x7a, 0x98, 0x4f, 0x50,
	0x45, 0xc7, 0x35, 0xa2, 0x9c, 0x23, 0x67, 0xe7, 0x3f, 0x0e, 0xb4, 0xe7, 0xde, 0x00, 0x0d, 0x79,
	0x46, 0x5f, 0xa2, 0x0a, 0x5e, 0x87, 0xbe, 0xf1, 0x96, 0x70, 0x3e, 0x39, 0x61, 0x5c, 0x87, 0xe6,
	0x60, 0x80, 0x46, 0xf1, 0xbd, 0xa2, 0x6d, 0x0c, 0x44, 0x42, 0x51, 0x55, 0x5f, 0x76, 0xa4, 0x77,
	0x82, 0xd6, 0xd7, 0xb4, 0xe4, 0xd3, 0x94, 0x93, 0x90, 0x46, 0xa8, 0x8e, 0xd7, 0xa0, 0x7b, 0x4e,
	0x13, 0x9d, 0x5c, 0x0b, 0x41, 0x0d, 0x1d, 0xcc, 0x48, 0x89, 0x34, 0xd5, 0xc9, 0xb2, 0xe8, 0x3f,
	0xd3, 0x50, 0xd1, 0x08, 0x35, 0x71, 0x17, 0x5a, 0xa3, 0x3c, 0x4b, 0x69, 0x12, 0xd1, 0x08, 0xb5,
	0x4c, 0xe9, 0x2c, 0x59, 0xfb, 0x04, 0x5a, 0x3e, 0x22, 0x3c, 0xcc, 0x39, 0xd1, 0xf0, 0xb6, 0xb6,
	0x74, 0x7c, 0x9b, 0x32, 0x49, 0x23, 0xd4, 0xc1, 0x18, 0x7a, 0x05, 0xb8, 0xb8, 0x1e, 0x75, 0x77,
	0x22, 0x68, 0x4d, 0x9f, 0x2a, 0xed, 0xbc, 0x7d, 0x9a, 0x6d, 0x0e, 0xcd, 0x9b, 0x8f, 0x1c, 0x93,
	0x0b, 0x25, 0x52, 0xb4, 0x62, 0x6e, 0x57, 0x22, 0xb5, 0x07, 0x55, 0x6d, 0xd0, 0x7c, 0x0e, 0xaf,
	0x2e, 0x44, 0x1e, 0x5e, 0x9b, 0xe0, 0xd6, 0xa1, 0x6f, 0x6d, 0xcc, 0x94, 0xf5, 0x9d, 0x31, 0xb4,
	0xe7, 0x5e, 0x2e, 0x13, 0x1f, 0x35, 0x83, 0x8e, 0x2a, 0xda, 0xc8, 0x6f, 0x85, 0x88, 0x2e, 0x18,
	0xe7, 0x45, 0x02, 0x1c, 0xd3, 0x4b, 0xea, 0xe2, 0x9a, 0x9e, 0xa5, 0x34, 0xb1, 0x5d, 0xf7, 0x39,
	0xac, 0x0d, 0xe3, 0x98, 0x46, 0xba, 0x17, 0xce, 0x64, 0x01, 0xac, 0xea, 0x78, 0x75, 0xba, 0xcf,
	0xe4, 0xef, 0x19, 0xe7, 0xa8, 0x76, 0x78, 0xf0, 0xfa, 0xad, 0x57, 0x79, 0xf3, 0xd6, 0xab, 0xbc,
	0x7f, 0xeb, 0x39, 0x7f, 0xb9, 0xf3, 0x9c, 0x57, 0x77, 0x9e, 0xf3, 0xef, 0x3b, 0xcf, 0x79, 0x7d,
	0xe7, 0x39, 0xff, 0xbd, 0xf3, 0x9c, 0xff, 0xdd, 0x79, 0x95, 0xf7, 0x77, 0x9e, 0xf3, 0xb7, 0x77,
	0x5e, 0xe5, 0xf5, 0x3b, 0xaf, 0xf2, 0xe6, 0x9d, 0x57, 0xb9, 0x6c, 0x98, 0xd1, 0xf9, 0xe5, 0xff,
	0x03, 0x00, 0x00, 0xff, 0xff, 0x10, 0x7e, 0xf5, 0x47, 0x79, 0x0d, 0x00, 0x00,
}

func (x InstrumentStatus) String() string {
	s, ok := InstrumentStatus_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x ExecutionInstruction) String() string {
	s, ok := ExecutionInstruction_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x TradeType) String() string {
	s, ok := TradeType_name[int32(x)]
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
	if !this.MaxLimitQuantity.Equal(that1.MaxLimitQuantity) {
		return false
	}
	if !this.MaxMarketQuantity.Equal(that1.MaxMarketQuantity) {
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
	if len(this.ExecutionInstructions) != len(that1.ExecutionInstructions) {
		return false
	}
	for i := range this.ExecutionInstructions {
		if this.ExecutionInstructions[i] != that1.ExecutionInstructions[i] {
			return false
		}
	}
	if this.Tag != that1.Tag {
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
	if this.Account != that1.Account {
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
	if this.Account != that1.Account {
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
func (this *TradeCapture) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*TradeCapture)
	if !ok {
		that2, ok := that.(TradeCapture)
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
	if this.Side != that1.Side {
		return false
	}
	if this.Type != that1.Type {
		return false
	}
	if this.Price != that1.Price {
		return false
	}
	if this.Quantity != that1.Quantity {
		return false
	}
	if this.Commission != that1.Commission {
		return false
	}
	if !this.CommissionAsset.Equal(that1.CommissionAsset) {
		return false
	}
	if this.TradeID != that1.TradeID {
		return false
	}
	if !this.Instrument.Equal(that1.Instrument) {
		return false
	}
	if !this.Trade_LinkID.Equal(that1.Trade_LinkID) {
		return false
	}
	if !this.OrderID.Equal(that1.OrderID) {
		return false
	}
	if !this.ClientOrderID.Equal(that1.ClientOrderID) {
		return false
	}
	if !this.TransactionTime.Equal(that1.TransactionTime) {
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
	s := make([]string, 0, 21)
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
	if this.MaxLimitQuantity != nil {
		s = append(s, "MaxLimitQuantity: "+fmt.Sprintf("%#v", this.MaxLimitQuantity)+",\n")
	}
	if this.MaxMarketQuantity != nil {
		s = append(s, "MaxMarketQuantity: "+fmt.Sprintf("%#v", this.MaxMarketQuantity)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Order) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 16)
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
	s = append(s, "ExecutionInstructions: "+fmt.Sprintf("%#v", this.ExecutionInstructions)+",\n")
	s = append(s, "Tag: "+fmt.Sprintf("%#v", this.Tag)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Position) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&models.Position{")
	s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
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
	s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
	if this.Asset != nil {
		s = append(s, "Asset: "+fmt.Sprintf("%#v", this.Asset)+",\n")
	}
	s = append(s, "Quantity: "+fmt.Sprintf("%#v", this.Quantity)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *TradeCapture) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 16)
	s = append(s, "&models.TradeCapture{")
	s = append(s, "Side: "+fmt.Sprintf("%#v", this.Side)+",\n")
	s = append(s, "Type: "+fmt.Sprintf("%#v", this.Type)+",\n")
	s = append(s, "Price: "+fmt.Sprintf("%#v", this.Price)+",\n")
	s = append(s, "Quantity: "+fmt.Sprintf("%#v", this.Quantity)+",\n")
	s = append(s, "Commission: "+fmt.Sprintf("%#v", this.Commission)+",\n")
	if this.CommissionAsset != nil {
		s = append(s, "CommissionAsset: "+fmt.Sprintf("%#v", this.CommissionAsset)+",\n")
	}
	s = append(s, "TradeID: "+fmt.Sprintf("%#v", this.TradeID)+",\n")
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	if this.Trade_LinkID != nil {
		s = append(s, "Trade_LinkID: "+fmt.Sprintf("%#v", this.Trade_LinkID)+",\n")
	}
	if this.OrderID != nil {
		s = append(s, "OrderID: "+fmt.Sprintf("%#v", this.OrderID)+",\n")
	}
	if this.ClientOrderID != nil {
		s = append(s, "ClientOrderID: "+fmt.Sprintf("%#v", this.ClientOrderID)+",\n")
	}
	if this.TransactionTime != nil {
		s = append(s, "TransactionTime: "+fmt.Sprintf("%#v", this.TransactionTime)+",\n")
	}
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
	if m.MaxMarketQuantity != nil {
		{
			size, err := m.MaxMarketQuantity.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x9a
	}
	if m.MaxLimitQuantity != nil {
		{
			size, err := m.MaxLimitQuantity.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x92
	}
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
	if len(m.Tag) > 0 {
		i -= len(m.Tag)
		copy(dAtA[i:], m.Tag)
		i = encodeVarintSecurityData(dAtA, i, uint64(len(m.Tag)))
		i--
		dAtA[i] = 0x62
	}
	if len(m.ExecutionInstructions) > 0 {
		dAtA17 := make([]byte, len(m.ExecutionInstructions)*10)
		var j16 int
		for _, num := range m.ExecutionInstructions {
			for num >= 1<<7 {
				dAtA17[j16] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j16++
			}
			dAtA17[j16] = uint8(num)
			j16++
		}
		i -= j16
		copy(dAtA[i:], dAtA17[:j16])
		i = encodeVarintSecurityData(dAtA, i, uint64(j16))
		i--
		dAtA[i] = 0x5a
	}
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
	if len(m.Account) > 0 {
		i -= len(m.Account)
		copy(dAtA[i:], m.Account)
		i = encodeVarintSecurityData(dAtA, i, uint64(len(m.Account)))
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
	if len(m.Account) > 0 {
		i -= len(m.Account)
		copy(dAtA[i:], m.Account)
		i = encodeVarintSecurityData(dAtA, i, uint64(len(m.Account)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TradeCapture) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TradeCapture) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TradeCapture) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.TransactionTime != nil {
		{
			size, err := m.TransactionTime.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x62
	}
	if m.ClientOrderID != nil {
		{
			size, err := m.ClientOrderID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x5a
	}
	if m.OrderID != nil {
		{
			size, err := m.OrderID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x52
	}
	if m.Trade_LinkID != nil {
		{
			size, err := m.Trade_LinkID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x4a
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
		dAtA[i] = 0x42
	}
	if len(m.TradeID) > 0 {
		i -= len(m.TradeID)
		copy(dAtA[i:], m.TradeID)
		i = encodeVarintSecurityData(dAtA, i, uint64(len(m.TradeID)))
		i--
		dAtA[i] = 0x3a
	}
	if m.CommissionAsset != nil {
		{
			size, err := m.CommissionAsset.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSecurityData(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.Commission != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Commission))))
		i--
		dAtA[i] = 0x29
	}
	if m.Quantity != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Quantity))))
		i--
		dAtA[i] = 0x21
	}
	if m.Price != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Price))))
		i--
		dAtA[i] = 0x19
	}
	if m.Type != 0 {
		i = encodeVarintSecurityData(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x10
	}
	if m.Side != 0 {
		i = encodeVarintSecurityData(dAtA, i, uint64(m.Side))
		i--
		dAtA[i] = 0x8
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
	if m.MaxLimitQuantity != nil {
		l = m.MaxLimitQuantity.Size()
		n += 2 + l + sovSecurityData(uint64(l))
	}
	if m.MaxMarketQuantity != nil {
		l = m.MaxMarketQuantity.Size()
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
	if len(m.ExecutionInstructions) > 0 {
		l = 0
		for _, e := range m.ExecutionInstructions {
			l += sovSecurityData(uint64(e))
		}
		n += 1 + sovSecurityData(uint64(l)) + l
	}
	l = len(m.Tag)
	if l > 0 {
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
	l = len(m.Account)
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
	l = len(m.Account)
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

func (m *TradeCapture) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Side != 0 {
		n += 1 + sovSecurityData(uint64(m.Side))
	}
	if m.Type != 0 {
		n += 1 + sovSecurityData(uint64(m.Type))
	}
	if m.Price != 0 {
		n += 9
	}
	if m.Quantity != 0 {
		n += 9
	}
	if m.Commission != 0 {
		n += 9
	}
	if m.CommissionAsset != nil {
		l = m.CommissionAsset.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	l = len(m.TradeID)
	if l > 0 {
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.Instrument != nil {
		l = m.Instrument.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.Trade_LinkID != nil {
		l = m.Trade_LinkID.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.OrderID != nil {
		l = m.OrderID.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.ClientOrderID != nil {
		l = m.ClientOrderID.Size()
		n += 1 + l + sovSecurityData(uint64(l))
	}
	if m.TransactionTime != nil {
		l = m.TransactionTime.Size()
		n += 1 + l + sovSecurityData(uint64(l))
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
		`MaxLimitQuantity:` + strings.Replace(fmt.Sprintf("%v", this.MaxLimitQuantity), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`MaxMarketQuantity:` + strings.Replace(fmt.Sprintf("%v", this.MaxMarketQuantity), "DoubleValue", "types.DoubleValue", 1) + `,`,
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
		`ExecutionInstructions:` + fmt.Sprintf("%v", this.ExecutionInstructions) + `,`,
		`Tag:` + fmt.Sprintf("%v", this.Tag) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Position) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Position{`,
		`Account:` + fmt.Sprintf("%v", this.Account) + `,`,
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
		`Account:` + fmt.Sprintf("%v", this.Account) + `,`,
		`Asset:` + strings.Replace(fmt.Sprintf("%v", this.Asset), "Asset", "models.Asset", 1) + `,`,
		`Quantity:` + fmt.Sprintf("%v", this.Quantity) + `,`,
		`}`,
	}, "")
	return s
}
func (this *TradeCapture) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&TradeCapture{`,
		`Side:` + fmt.Sprintf("%v", this.Side) + `,`,
		`Type:` + fmt.Sprintf("%v", this.Type) + `,`,
		`Price:` + fmt.Sprintf("%v", this.Price) + `,`,
		`Quantity:` + fmt.Sprintf("%v", this.Quantity) + `,`,
		`Commission:` + fmt.Sprintf("%v", this.Commission) + `,`,
		`CommissionAsset:` + strings.Replace(fmt.Sprintf("%v", this.CommissionAsset), "Asset", "models.Asset", 1) + `,`,
		`TradeID:` + fmt.Sprintf("%v", this.TradeID) + `,`,
		`Instrument:` + strings.Replace(this.Instrument.String(), "Instrument", "Instrument", 1) + `,`,
		`Trade_LinkID:` + strings.Replace(fmt.Sprintf("%v", this.Trade_LinkID), "StringValue", "types.StringValue", 1) + `,`,
		`OrderID:` + strings.Replace(fmt.Sprintf("%v", this.OrderID), "StringValue", "types.StringValue", 1) + `,`,
		`ClientOrderID:` + strings.Replace(fmt.Sprintf("%v", this.ClientOrderID), "StringValue", "types.StringValue", 1) + `,`,
		`TransactionTime:` + strings.Replace(fmt.Sprintf("%v", this.TransactionTime), "Timestamp", "types.Timestamp", 1) + `,`,
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
		case 18:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxLimitQuantity", wireType)
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
			if m.MaxLimitQuantity == nil {
				m.MaxLimitQuantity = &types.DoubleValue{}
			}
			if err := m.MaxLimitQuantity.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 19:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxMarketQuantity", wireType)
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
			if m.MaxMarketQuantity == nil {
				m.MaxMarketQuantity = &types.DoubleValue{}
			}
			if err := m.MaxMarketQuantity.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
		case 11:
			if wireType == 0 {
				var v ExecutionInstruction
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSecurityData
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= ExecutionInstruction(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.ExecutionInstructions = append(m.ExecutionInstructions, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSecurityData
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthSecurityData
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthSecurityData
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				if elementCount != 0 && len(m.ExecutionInstructions) == 0 {
					m.ExecutionInstructions = make([]ExecutionInstruction, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v ExecutionInstruction
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSecurityData
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= ExecutionInstruction(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.ExecutionInstructions = append(m.ExecutionInstructions, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field ExecutionInstructions", wireType)
			}
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tag", wireType)
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
			m.Tag = string(dAtA[iNdEx:postIndex])
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
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
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
			m.Account = string(dAtA[iNdEx:postIndex])
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
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
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
			m.Account = string(dAtA[iNdEx:postIndex])
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
func (m *TradeCapture) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: TradeCapture: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TradeCapture: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
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
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSecurityData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= TradeType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
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
		case 4:
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
		case 5:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Commission", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Commission = float64(math.Float64frombits(v))
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CommissionAsset", wireType)
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
			if m.CommissionAsset == nil {
				m.CommissionAsset = &models.Asset{}
			}
			if err := m.CommissionAsset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TradeID", wireType)
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
			m.TradeID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
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
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Trade_LinkID", wireType)
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
			if m.Trade_LinkID == nil {
				m.Trade_LinkID = &types.StringValue{}
			}
			if err := m.Trade_LinkID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderID", wireType)
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
			if m.OrderID == nil {
				m.OrderID = &types.StringValue{}
			}
			if err := m.OrderID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientOrderID", wireType)
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
			if m.ClientOrderID == nil {
				m.ClientOrderID = &types.StringValue{}
			}
			if err := m.ClientOrderID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransactionTime", wireType)
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
			if m.TransactionTime == nil {
				m.TransactionTime = &types.Timestamp{}
			}
			if err := m.TransactionTime.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
