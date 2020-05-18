package messages

import (
	encoding_binary "encoding/binary"
	fmt "fmt"
	actor "github.com/AsynkronIT/protoactor-go/actor"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
	models "gitlab.com/alphaticks/alphac/models"
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

type ExecutionType int32

const (
	New            ExecutionType = 0
	Fill           ExecutionType = 1
	Done           ExecutionType = 2
	Canceled       ExecutionType = 3
	Replaced       ExecutionType = 4
	PendingCancel  ExecutionType = 5
	Stopped        ExecutionType = 6
	Rejected       ExecutionType = 7
	Suspended      ExecutionType = 8
	PendingNew     ExecutionType = 9
	Calculated     ExecutionType = 10
	Expired        ExecutionType = 11
	PendingReplace ExecutionType = 12
	Trade          ExecutionType = 13
	OrderStatus    ExecutionType = 14
	Settlement     ExecutionType = 15
)

var ExecutionType_name = map[int32]string{
	0:  "New",
	1:  "Fill",
	2:  "Done",
	3:  "Canceled",
	4:  "Replaced",
	5:  "PendingCancel",
	6:  "Stopped",
	7:  "Rejected",
	8:  "Suspended",
	9:  "PendingNew",
	10: "Calculated",
	11: "Expired",
	12: "PendingReplace",
	13: "Trade",
	14: "OrderStatus",
	15: "Settlement",
}

var ExecutionType_value = map[string]int32{
	"New":            0,
	"Fill":           1,
	"Done":           2,
	"Canceled":       3,
	"Replaced":       4,
	"PendingCancel":  5,
	"Stopped":        6,
	"Rejected":       7,
	"Suspended":      8,
	"PendingNew":     9,
	"Calculated":     10,
	"Expired":        11,
	"PendingReplace": 12,
	"Trade":          13,
	"OrderStatus":    14,
	"Settlement":     15,
}

func (ExecutionType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{0}
}

type OrderRejectionReason int32

const (
	UnknownSymbol                  OrderRejectionReason = 0
	UnsupportedOrderCharacteristic OrderRejectionReason = 1
	IncorrectQuantity              OrderRejectionReason = 2
	ExchangeClosed                 OrderRejectionReason = 3
	DuplicateOrder                 OrderRejectionReason = 4
	InvalidAccount                 OrderRejectionReason = 5
	Other                          OrderRejectionReason = 6
)

var OrderRejectionReason_name = map[int32]string{
	0: "UnknownSymbol",
	1: "UnsupportedOrderCharacteristic",
	2: "IncorrectQuantity",
	3: "ExchangeClosed",
	4: "DuplicateOrder",
	5: "InvalidAccount",
	6: "Other",
}

var OrderRejectionReason_value = map[string]int32{
	"UnknownSymbol":                  0,
	"UnsupportedOrderCharacteristic": 1,
	"IncorrectQuantity":              2,
	"ExchangeClosed":                 3,
	"DuplicateOrder":                 4,
	"InvalidAccount":                 5,
	"Other":                          6,
}

func (OrderRejectionReason) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{1}
}

type BusinessRequestReject struct {
	RequestID uint64 `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Reason    string `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
}

func (m *BusinessRequestReject) Reset()      { *m = BusinessRequestReject{} }
func (*BusinessRequestReject) ProtoMessage() {}
func (*BusinessRequestReject) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{0}
}
func (m *BusinessRequestReject) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BusinessRequestReject) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BusinessRequestReject.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BusinessRequestReject) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BusinessRequestReject.Merge(m, src)
}
func (m *BusinessRequestReject) XXX_Size() int {
	return m.Size()
}
func (m *BusinessRequestReject) XXX_DiscardUnknown() {
	xxx_messageInfo_BusinessRequestReject.DiscardUnknown(m)
}

var xxx_messageInfo_BusinessRequestReject proto.InternalMessageInfo

func (m *BusinessRequestReject) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *BusinessRequestReject) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type MarketDataRequest struct {
	RequestID   uint64                      `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe   bool                        `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
	Subscriber  *actor.PID                  `protobuf:"bytes,3,opt,name=subscriber,proto3" json:"subscriber,omitempty"`
	Instrument  *models.Instrument          `protobuf:"bytes,4,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Aggregation models.OrderBookAggregation `protobuf:"varint,5,opt,name=aggregation,proto3,enum=models.OrderBookAggregation" json:"aggregation,omitempty"`
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

func (m *MarketDataRequest) GetSubscriber() *actor.PID {
	if m != nil {
		return m.Subscriber
	}
	return nil
}

func (m *MarketDataRequest) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *MarketDataRequest) GetAggregation() models.OrderBookAggregation {
	if m != nil {
		return m.Aggregation
	}
	return models.L2
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
	RequestID  uint64                    `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID uint64                    `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	SnapshotL2 *models.OBL2Snapshot      `protobuf:"bytes,3,opt,name=snapshotL2,proto3" json:"snapshotL2,omitempty"`
	SnapshotL3 *models.OBL3Snapshot      `protobuf:"bytes,4,opt,name=snapshotL3,proto3" json:"snapshotL3,omitempty"`
	Trades     []*models.AggregatedTrade `protobuf:"bytes,5,rep,name=trades,proto3" json:"trades,omitempty"`
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

func (m *MarketDataSnapshot) GetSnapshotL2() *models.OBL2Snapshot {
	if m != nil {
		return m.SnapshotL2
	}
	return nil
}

func (m *MarketDataSnapshot) GetSnapshotL3() *models.OBL3Snapshot {
	if m != nil {
		return m.SnapshotL3
	}
	return nil
}

func (m *MarketDataSnapshot) GetTrades() []*models.AggregatedTrade {
	if m != nil {
		return m.Trades
	}
	return nil
}

type MarketDataIncrementalRefresh struct {
	RequestID  uint64                    `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID uint64                    `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	UpdateL2   *models.OBL2Update        `protobuf:"bytes,3,opt,name=updateL2,proto3" json:"updateL2,omitempty"`
	UpdateL3   *models.OBL3Update        `protobuf:"bytes,4,opt,name=updateL3,proto3" json:"updateL3,omitempty"`
	Trades     []*models.AggregatedTrade `protobuf:"bytes,5,rep,name=trades,proto3" json:"trades,omitempty"`
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

func (m *MarketDataIncrementalRefresh) GetUpdateL2() *models.OBL2Update {
	if m != nil {
		return m.UpdateL2
	}
	return nil
}

func (m *MarketDataIncrementalRefresh) GetUpdateL3() *models.OBL3Update {
	if m != nil {
		return m.UpdateL3
	}
	return nil
}

func (m *MarketDataIncrementalRefresh) GetTrades() []*models.AggregatedTrade {
	if m != nil {
		return m.Trades
	}
	return nil
}

type SecurityListRequest struct {
	RequestID  uint64     `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe  bool       `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
	Subscriber *actor.PID `protobuf:"bytes,3,opt,name=subscriber,proto3" json:"subscriber,omitempty"`
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

func (m *SecurityListRequest) GetSubscriber() *actor.PID {
	if m != nil {
		return m.Subscriber
	}
	return nil
}

type SecurityList struct {
	RequestID  uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID uint64             `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Error      string             `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	Securities []*models.Security `protobuf:"bytes,4,rep,name=securities,proto3" json:"securities,omitempty"`
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

func (m *SecurityList) GetSecurities() []*models.Security {
	if m != nil {
		return m.Securities
	}
	return nil
}

type ExecutionReport struct {
	SeqNum               uint64               `protobuf:"varint,1,opt,name=seq_num,json=seqNum,proto3" json:"seq_num,omitempty"`
	OrderID              string               `protobuf:"bytes,2,opt,name=orderID,proto3" json:"orderID,omitempty"`
	ClientOrderID        *types.StringValue   `protobuf:"bytes,3,opt,name=client_orderID,json=clientOrderID,proto3" json:"client_orderID,omitempty"`
	ExecutionID          string               `protobuf:"bytes,4,opt,name=executionID,proto3" json:"executionID,omitempty"`
	ExecutionType        ExecutionType        `protobuf:"varint,5,opt,name=execution_type,json=executionType,proto3,enum=messages.ExecutionType" json:"execution_type,omitempty"`
	OrderStatus          models.OrderStatus   `protobuf:"varint,6,opt,name=order_status,json=orderStatus,proto3,enum=models.OrderStatus" json:"order_status,omitempty"`
	Instrument           *models.Instrument   `protobuf:"bytes,7,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Side                 models.Side          `protobuf:"varint,8,opt,name=side,proto3,enum=models.Side" json:"side,omitempty"`
	LeavesQuantity       float64              `protobuf:"fixed64,9,opt,name=leaves_quantity,json=leavesQuantity,proto3" json:"leaves_quantity,omitempty"`
	CumQuantity          float64              `protobuf:"fixed64,10,opt,name=cum_quantity,json=cumQuantity,proto3" json:"cum_quantity,omitempty"`
	TransactionTime      *types.Timestamp     `protobuf:"bytes,11,opt,name=transaction_time,json=transactionTime,proto3" json:"transaction_time,omitempty"`
	TradeID              *types.StringValue   `protobuf:"bytes,12,opt,name=tradeID,proto3" json:"tradeID,omitempty"`
	OrderRejectionReason OrderRejectionReason `protobuf:"varint,13,opt,name=order_rejection_reason,json=orderRejectionReason,proto3,enum=messages.OrderRejectionReason" json:"order_rejection_reason,omitempty"`
}

func (m *ExecutionReport) Reset()      { *m = ExecutionReport{} }
func (*ExecutionReport) ProtoMessage() {}
func (*ExecutionReport) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{7}
}
func (m *ExecutionReport) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExecutionReport) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExecutionReport.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExecutionReport) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecutionReport.Merge(m, src)
}
func (m *ExecutionReport) XXX_Size() int {
	return m.Size()
}
func (m *ExecutionReport) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecutionReport.DiscardUnknown(m)
}

var xxx_messageInfo_ExecutionReport proto.InternalMessageInfo

func (m *ExecutionReport) GetSeqNum() uint64 {
	if m != nil {
		return m.SeqNum
	}
	return 0
}

func (m *ExecutionReport) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *ExecutionReport) GetClientOrderID() *types.StringValue {
	if m != nil {
		return m.ClientOrderID
	}
	return nil
}

func (m *ExecutionReport) GetExecutionID() string {
	if m != nil {
		return m.ExecutionID
	}
	return ""
}

func (m *ExecutionReport) GetExecutionType() ExecutionType {
	if m != nil {
		return m.ExecutionType
	}
	return New
}

func (m *ExecutionReport) GetOrderStatus() models.OrderStatus {
	if m != nil {
		return m.OrderStatus
	}
	return models.New
}

func (m *ExecutionReport) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *ExecutionReport) GetSide() models.Side {
	if m != nil {
		return m.Side
	}
	return models.Buy
}

func (m *ExecutionReport) GetLeavesQuantity() float64 {
	if m != nil {
		return m.LeavesQuantity
	}
	return 0
}

func (m *ExecutionReport) GetCumQuantity() float64 {
	if m != nil {
		return m.CumQuantity
	}
	return 0
}

func (m *ExecutionReport) GetTransactionTime() *types.Timestamp {
	if m != nil {
		return m.TransactionTime
	}
	return nil
}

func (m *ExecutionReport) GetTradeID() *types.StringValue {
	if m != nil {
		return m.TradeID
	}
	return nil
}

func (m *ExecutionReport) GetOrderRejectionReason() OrderRejectionReason {
	if m != nil {
		return m.OrderRejectionReason
	}
	return UnknownSymbol
}

type OrderStatusRequest struct {
	RequestID     uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe     bool               `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
	Subscriber    *actor.PID         `protobuf:"bytes,3,opt,name=subscriber,proto3" json:"subscriber,omitempty"`
	OrderID       *types.StringValue `protobuf:"bytes,4,opt,name=orderID,proto3" json:"orderID,omitempty"`
	ClientOrderID *types.StringValue `protobuf:"bytes,5,opt,name=client_orderID,json=clientOrderID,proto3" json:"client_orderID,omitempty"`
	Instrument    *models.Instrument `protobuf:"bytes,6,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Account       *models.Account    `protobuf:"bytes,7,opt,name=account,proto3" json:"account,omitempty"`
}

func (m *OrderStatusRequest) Reset()      { *m = OrderStatusRequest{} }
func (*OrderStatusRequest) ProtoMessage() {}
func (*OrderStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{8}
}
func (m *OrderStatusRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderStatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderStatusRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderStatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderStatusRequest.Merge(m, src)
}
func (m *OrderStatusRequest) XXX_Size() int {
	return m.Size()
}
func (m *OrderStatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderStatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OrderStatusRequest proto.InternalMessageInfo

func (m *OrderStatusRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *OrderStatusRequest) GetSubscribe() bool {
	if m != nil {
		return m.Subscribe
	}
	return false
}

func (m *OrderStatusRequest) GetSubscriber() *actor.PID {
	if m != nil {
		return m.Subscriber
	}
	return nil
}

func (m *OrderStatusRequest) GetOrderID() *types.StringValue {
	if m != nil {
		return m.OrderID
	}
	return nil
}

func (m *OrderStatusRequest) GetClientOrderID() *types.StringValue {
	if m != nil {
		return m.ClientOrderID
	}
	return nil
}

func (m *OrderStatusRequest) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *OrderStatusRequest) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

type OrderList struct {
	RequestID  uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID uint64          `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Error      string          `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	Orders     []*models.Order `protobuf:"bytes,4,rep,name=orders,proto3" json:"orders,omitempty"`
}

func (m *OrderList) Reset()      { *m = OrderList{} }
func (*OrderList) ProtoMessage() {}
func (*OrderList) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{9}
}
func (m *OrderList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderList.Merge(m, src)
}
func (m *OrderList) XXX_Size() int {
	return m.Size()
}
func (m *OrderList) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderList.DiscardUnknown(m)
}

var xxx_messageInfo_OrderList proto.InternalMessageInfo

func (m *OrderList) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *OrderList) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *OrderList) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *OrderList) GetOrders() []*models.Order {
	if m != nil {
		return m.Orders
	}
	return nil
}

type PositionsRequest struct {
	RequestID  uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe  bool               `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
	Subscriber *actor.PID         `protobuf:"bytes,3,opt,name=subscriber,proto3" json:"subscriber,omitempty"`
	Instrument *models.Instrument `protobuf:"bytes,4,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Account    *models.Account    `protobuf:"bytes,5,opt,name=account,proto3" json:"account,omitempty"`
}

func (m *PositionsRequest) Reset()      { *m = PositionsRequest{} }
func (*PositionsRequest) ProtoMessage() {}
func (*PositionsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{10}
}
func (m *PositionsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PositionsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PositionsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PositionsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PositionsRequest.Merge(m, src)
}
func (m *PositionsRequest) XXX_Size() int {
	return m.Size()
}
func (m *PositionsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PositionsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PositionsRequest proto.InternalMessageInfo

func (m *PositionsRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *PositionsRequest) GetSubscribe() bool {
	if m != nil {
		return m.Subscribe
	}
	return false
}

func (m *PositionsRequest) GetSubscriber() *actor.PID {
	if m != nil {
		return m.Subscriber
	}
	return nil
}

func (m *PositionsRequest) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *PositionsRequest) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

type PositionList struct {
	RequestID  uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe  bool               `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
	Subscriber *actor.PID         `protobuf:"bytes,3,opt,name=subscriber,proto3" json:"subscriber,omitempty"`
	ResponseID uint64             `protobuf:"varint,4,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Error      string             `protobuf:"bytes,5,opt,name=error,proto3" json:"error,omitempty"`
	Positions  []*models.Position `protobuf:"bytes,6,rep,name=positions,proto3" json:"positions,omitempty"`
}

func (m *PositionList) Reset()      { *m = PositionList{} }
func (*PositionList) ProtoMessage() {}
func (*PositionList) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{11}
}
func (m *PositionList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PositionList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PositionList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PositionList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PositionList.Merge(m, src)
}
func (m *PositionList) XXX_Size() int {
	return m.Size()
}
func (m *PositionList) XXX_DiscardUnknown() {
	xxx_messageInfo_PositionList.DiscardUnknown(m)
}

var xxx_messageInfo_PositionList proto.InternalMessageInfo

func (m *PositionList) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *PositionList) GetSubscribe() bool {
	if m != nil {
		return m.Subscribe
	}
	return false
}

func (m *PositionList) GetSubscriber() *actor.PID {
	if m != nil {
		return m.Subscriber
	}
	return nil
}

func (m *PositionList) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *PositionList) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *PositionList) GetPositions() []*models.Position {
	if m != nil {
		return m.Positions
	}
	return nil
}

type NewOrderSingle struct {
	ClientOrderID string             `protobuf:"bytes,1,opt,name=client_orderID,json=clientOrderID,proto3" json:"client_orderID,omitempty"`
	Instrument    *models.Instrument `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Account       *models.Account    `protobuf:"bytes,3,opt,name=account,proto3" json:"account,omitempty"`
	OrderType     models.OrderType   `protobuf:"varint,4,opt,name=order_type,json=orderType,proto3,enum=models.OrderType" json:"order_type,omitempty"`
	OrderSide     models.Side        `protobuf:"varint,5,opt,name=order_side,json=orderSide,proto3,enum=models.Side" json:"order_side,omitempty"`
	TimeInForce   models.TimeInForce `protobuf:"varint,6,opt,name=time_in_force,json=timeInForce,proto3,enum=models.TimeInForce" json:"time_in_force,omitempty"`
	Quantity      float64            `protobuf:"fixed64,7,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Price         *types.DoubleValue `protobuf:"bytes,8,opt,name=price,proto3" json:"price,omitempty"`
}

func (m *NewOrderSingle) Reset()      { *m = NewOrderSingle{} }
func (*NewOrderSingle) ProtoMessage() {}
func (*NewOrderSingle) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{12}
}
func (m *NewOrderSingle) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NewOrderSingle) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NewOrderSingle.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NewOrderSingle) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewOrderSingle.Merge(m, src)
}
func (m *NewOrderSingle) XXX_Size() int {
	return m.Size()
}
func (m *NewOrderSingle) XXX_DiscardUnknown() {
	xxx_messageInfo_NewOrderSingle.DiscardUnknown(m)
}

var xxx_messageInfo_NewOrderSingle proto.InternalMessageInfo

func (m *NewOrderSingle) GetClientOrderID() string {
	if m != nil {
		return m.ClientOrderID
	}
	return ""
}

func (m *NewOrderSingle) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *NewOrderSingle) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

func (m *NewOrderSingle) GetOrderType() models.OrderType {
	if m != nil {
		return m.OrderType
	}
	return models.Market
}

func (m *NewOrderSingle) GetOrderSide() models.Side {
	if m != nil {
		return m.OrderSide
	}
	return models.Buy
}

func (m *NewOrderSingle) GetTimeInForce() models.TimeInForce {
	if m != nil {
		return m.TimeInForce
	}
	return models.Session
}

func (m *NewOrderSingle) GetQuantity() float64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *NewOrderSingle) GetPrice() *types.DoubleValue {
	if m != nil {
		return m.Price
	}
	return nil
}

func init() {
	proto.RegisterEnum("messages.ExecutionType", ExecutionType_name, ExecutionType_value)
	proto.RegisterEnum("messages.OrderRejectionReason", OrderRejectionReason_name, OrderRejectionReason_value)
	proto.RegisterType((*BusinessRequestReject)(nil), "messages.BusinessRequestReject")
	proto.RegisterType((*MarketDataRequest)(nil), "messages.MarketDataRequest")
	proto.RegisterType((*MarketDataRequestReject)(nil), "messages.MarketDataRequestReject")
	proto.RegisterType((*MarketDataSnapshot)(nil), "messages.MarketDataSnapshot")
	proto.RegisterType((*MarketDataIncrementalRefresh)(nil), "messages.MarketDataIncrementalRefresh")
	proto.RegisterType((*SecurityListRequest)(nil), "messages.SecurityListRequest")
	proto.RegisterType((*SecurityList)(nil), "messages.SecurityList")
	proto.RegisterType((*ExecutionReport)(nil), "messages.ExecutionReport")
	proto.RegisterType((*OrderStatusRequest)(nil), "messages.OrderStatusRequest")
	proto.RegisterType((*OrderList)(nil), "messages.OrderList")
	proto.RegisterType((*PositionsRequest)(nil), "messages.PositionsRequest")
	proto.RegisterType((*PositionList)(nil), "messages.PositionList")
	proto.RegisterType((*NewOrderSingle)(nil), "messages.NewOrderSingle")
}

func init() { proto.RegisterFile("executor_messages.proto", fileDescriptor_350c53ba9303a7e6) }

var fileDescriptor_350c53ba9303a7e6 = []byte{
	// 1404 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x56, 0x4f, 0x6f, 0x1b, 0xbb,
	0x11, 0xd7, 0xda, 0xfa, 0x3b, 0xfa, 0xe3, 0x35, 0xe3, 0xc4, 0x82, 0x61, 0x6c, 0x55, 0x01, 0x41,
	0xdd, 0x14, 0x95, 0x5b, 0x3b, 0x75, 0x6e, 0x01, 0x62, 0xcb, 0x01, 0x04, 0x24, 0x71, 0x4a, 0x39,
	0xbd, 0x0a, 0xf4, 0xee, 0x44, 0xde, 0x7a, 0x45, 0xae, 0x49, 0x6e, 0x1c, 0x1f, 0x0a, 0xf4, 0xd4,
	0x73, 0x0f, 0x3d, 0xf5, 0x13, 0x14, 0xe8, 0x17, 0xe9, 0x31, 0xc7, 0x14, 0xe8, 0xa1, 0x71, 0x2e,
	0xbd, 0x14, 0x0d, 0xfa, 0x01, 0x1e, 0x1e, 0x96, 0xbb, 0x2b, 0xad, 0x1c, 0x27, 0x71, 0x5e, 0xde,
	0x0b, 0xf0, 0x6e, 0x9c, 0xe1, 0x6f, 0x86, 0x9c, 0x3f, 0xfc, 0x0d, 0x61, 0x15, 0x5f, 0xa2, 0x1b,
	0x69, 0x21, 0x47, 0x13, 0x54, 0x8a, 0x8d, 0x51, 0xf5, 0x42, 0x29, 0xb4, 0x20, 0xd5, 0x4c, 0x5e,
	0xfb, 0xc9, 0x58, 0x88, 0x71, 0x80, 0x9b, 0x46, 0x7f, 0x14, 0x3d, 0xdf, 0xd4, 0xfe, 0x04, 0x95,
	0x66, 0x93, 0x30, 0x81, 0xae, 0x39, 0x97, 0x01, 0x67, 0x92, 0x85, 0x21, 0xca, 0xd4, 0xd5, 0xda,
	0x6f, 0xc6, 0xbe, 0x0e, 0xd8, 0x51, 0xcf, 0x15, 0x93, 0x4d, 0x16, 0x84, 0xc7, 0x4c, 0xfb, 0xee,
	0x89, 0x4a, 0x96, 0xee, 0xe6, 0x44, 0x78, 0x18, 0xa8, 0xcd, 0x09, 0x93, 0x27, 0xa8, 0x47, 0x1e,
	0xd3, 0x2c, 0x35, 0xbb, 0x77, 0x1d, 0x33, 0x85, 0x6e, 0x24, 0x7d, 0x7d, 0x9e, 0x37, 0xfc, 0xf5,
	0x75, 0x0c, 0x99, 0xeb, 0x8a, 0x88, 0xeb, 0xd4, 0x64, 0x67, 0xec, 0xeb, 0xe3, 0x28, 0x31, 0x79,
	0xa0, 0xce, 0xf9, 0x89, 0x14, 0x7c, 0x70, 0x98, 0x44, 0xc4, 0x5c, 0x2d, 0xe4, 0x2f, 0xc7, 0x62,
	0xd3, 0x2c, 0x12, 0x5d, 0x1a, 0x5a, 0xf7, 0x31, 0xdc, 0xdc, 0x8d, 0x94, 0xcf, 0x51, 0x29, 0x8a,
	0xa7, 0x11, 0x2a, 0x4d, 0xf1, 0xf7, 0xe8, 0x6a, 0xb2, 0x0e, 0x35, 0x99, 0x28, 0x06, 0xfd, 0xb6,
	0xd5, 0xb1, 0x36, 0x8a, 0x74, 0xa6, 0x20, 0xb7, 0xa0, 0x2c, 0x91, 0x29, 0xc1, 0xdb, 0x0b, 0x1d,
	0x6b, 0xa3, 0x46, 0x53, 0xa9, 0xfb, 0x3f, 0x0b, 0x96, 0x1f, 0x9b, 0x44, 0xf4, 0x99, 0x66, 0xa9,
	0xc7, 0x4f, 0xf8, 0x5a, 0x87, 0x9a, 0x8a, 0x8e, 0x94, 0x2b, 0xfd, 0x23, 0x34, 0xee, 0xaa, 0x74,
	0xa6, 0x20, 0x77, 0x00, 0xa6, 0x82, 0x6c, 0x2f, 0x76, 0xac, 0x8d, 0xfa, 0x16, 0xf4, 0x4c, 0x24,
	0xbd, 0xa7, 0x83, 0x3e, 0xcd, 0xed, 0x92, 0x2d, 0x00, 0x9f, 0x2b, 0x2d, 0xa3, 0x09, 0x72, 0xdd,
	0x2e, 0x1a, 0x2c, 0xe9, 0x25, 0xf9, 0xea, 0x0d, 0xa6, 0x3b, 0x34, 0x87, 0x22, 0xf7, 0xa1, 0xce,
	0xc6, 0x63, 0x89, 0x63, 0xa6, 0x7d, 0xc1, 0xdb, 0xa5, 0x8e, 0xb5, 0xd1, 0xda, 0x5a, 0xcf, 0x8c,
	0x0e, 0xa4, 0x87, 0x72, 0x57, 0x88, 0x93, 0x07, 0x33, 0x0c, 0xcd, 0x1b, 0x74, 0x0f, 0x60, 0xf5,
	0xbd, 0x80, 0xbf, 0x34, 0x85, 0x64, 0xe6, 0x71, 0xc8, 0x59, 0xa8, 0x8e, 0xc5, 0xa7, 0x9c, 0x39,
	0x00, 0x12, 0x55, 0x28, 0xb8, 0xc2, 0x41, 0xdf, 0x38, 0x2c, 0xd2, 0x9c, 0x86, 0xdc, 0x05, 0x50,
	0xa9, 0xa7, 0x47, 0x5b, 0x69, 0x16, 0x57, 0xa6, 0x41, 0xee, 0x3e, 0xda, 0xca, 0xce, 0xa1, 0x39,
	0xdc, 0x9c, 0xd5, 0x76, 0x9a, 0xcf, 0xbc, 0xd5, 0xf6, 0x15, 0x56, 0xdb, 0x64, 0x13, 0xca, 0x5a,
	0x32, 0x0f, 0x55, 0xbb, 0xd4, 0x59, 0xdc, 0xa8, 0x6f, 0xad, 0x66, 0x16, 0x59, 0x0e, 0xd1, 0x3b,
	0x8c, 0xf7, 0x69, 0x0a, 0xeb, 0xfe, 0xd7, 0x82, 0xf5, 0x59, 0xc4, 0x03, 0xee, 0x4a, 0x8c, 0x4b,
	0xc3, 0x02, 0x8a, 0xcf, 0x25, 0xaa, 0xe3, 0x2f, 0x8c, 0xbd, 0x07, 0xd5, 0x28, 0xf4, 0x98, 0xc6,
	0x69, 0xe4, 0x24, 0x1f, 0xf9, 0x33, 0xb3, 0x47, 0xa7, 0x98, 0x1c, 0x7e, 0xfb, 0x72, 0x0f, 0xc5,
	0x31, 0x5f, 0xc2, 0x7f, 0x87, 0x78, 0xff, 0x00, 0x37, 0x86, 0xe9, 0xab, 0x7f, 0xe4, 0xc7, 0xdd,
	0xf2, 0x55, 0x5f, 0x49, 0xf7, 0x2f, 0x16, 0x34, 0xf2, 0xe7, 0x7f, 0x61, 0x7a, 0x57, 0xa0, 0x84,
	0x52, 0x8a, 0xe4, 0xd4, 0x1a, 0x4d, 0x04, 0xf2, 0x2b, 0x80, 0x94, 0xd9, 0x7c, 0x54, 0xed, 0xa2,
	0x49, 0x8c, 0x9d, 0x25, 0x26, 0x3b, 0x9d, 0xe6, 0x30, 0xdd, 0xbf, 0x96, 0x60, 0x69, 0xdf, 0x70,
	0x79, 0xfc, 0xc6, 0x30, 0x14, 0x52, 0x93, 0x55, 0xa8, 0x28, 0x3c, 0x1d, 0xf1, 0x68, 0x92, 0xde,
	0xab, 0xac, 0xf0, 0xf4, 0x49, 0x34, 0x21, 0x6d, 0xa8, 0x88, 0xf8, 0x69, 0xa6, 0x37, 0xaa, 0xd1,
	0x4c, 0x24, 0x7b, 0xd0, 0x72, 0x03, 0x1f, 0xb9, 0x1e, 0x65, 0x80, 0x24, 0x1b, 0xeb, 0xbd, 0x84,
	0xe4, 0x7b, 0x19, 0xc9, 0xf7, 0x86, 0x5a, 0xfa, 0x7c, 0xfc, 0x3b, 0x16, 0x44, 0x48, 0x9b, 0x89,
	0xcd, 0x41, 0xea, 0xa4, 0x03, 0x75, 0xcc, 0xae, 0x32, 0xe8, 0x9b, 0x2e, 0xa8, 0xd1, 0xbc, 0x8a,
	0xdc, 0x87, 0xd6, 0x54, 0x1c, 0xe9, 0xf3, 0x10, 0x53, 0xe6, 0x58, 0xed, 0x4d, 0xc7, 0xd0, 0x34,
	0x98, 0xc3, 0xf3, 0x10, 0x69, 0x13, 0xf3, 0x22, 0xd9, 0x81, 0x86, 0xb9, 0xdf, 0x48, 0x69, 0xa6,
	0x23, 0xd5, 0x2e, 0x1b, 0xeb, 0x1b, 0x73, 0xbc, 0x33, 0x34, 0x5b, 0xb4, 0x2e, 0x66, 0xc2, 0x25,
	0x8a, 0xab, 0x5c, 0x8b, 0xe2, 0x3a, 0x50, 0x54, 0xbe, 0x87, 0xed, 0xaa, 0x39, 0xa3, 0x31, 0xad,
	0x82, 0xef, 0x21, 0x35, 0x3b, 0xe4, 0x67, 0xb0, 0x14, 0x20, 0x7b, 0x81, 0x6a, 0x74, 0x1a, 0x31,
	0xae, 0x7d, 0x7d, 0xde, 0xae, 0x75, 0xac, 0x0d, 0x8b, 0xb6, 0x12, 0xf5, 0x6f, 0x53, 0x2d, 0xf9,
	0x29, 0x34, 0xdc, 0x68, 0x32, 0x43, 0x81, 0x41, 0xd5, 0xdd, 0x68, 0x32, 0x85, 0xec, 0x83, 0xad,
	0x25, 0xe3, 0x8a, 0xb9, 0x49, 0x6e, 0xfc, 0x09, 0xb6, 0xeb, 0xe6, 0x9e, 0x6b, 0xef, 0x95, 0xe0,
	0x30, 0x1b, 0xc4, 0x74, 0x29, 0x67, 0x13, 0x6b, 0xc9, 0x0e, 0x54, 0xcc, 0x73, 0x19, 0xf4, 0xdb,
	0x8d, 0x6b, 0x14, 0x30, 0x03, 0x93, 0x43, 0xb8, 0x95, 0x24, 0x56, 0x1a, 0x12, 0x8e, 0xaf, 0x90,
	0xd2, 0x6c, 0xd3, 0x84, 0xef, 0xcc, 0x0a, 0x64, 0x92, 0x4c, 0x33, 0x18, 0x35, 0x28, 0xba, 0x22,
	0xae, 0xd0, 0x76, 0x5f, 0x2f, 0x00, 0xc9, 0xd7, 0xe4, 0x2b, 0x0f, 0xb6, 0x9d, 0x59, 0xbb, 0x17,
	0xaf, 0x93, 0x8c, 0x0f, 0x3f, 0x86, 0xd2, 0xe7, 0x3f, 0x86, 0xf9, 0x96, 0x2b, 0x5f, 0xab, 0xe5,
	0x7e, 0x0e, 0x95, 0xf4, 0x7f, 0x92, 0xf6, 0xe8, 0xd2, 0x94, 0x14, 0x13, 0x35, 0xcd, 0xf6, 0xbb,
	0x7f, 0xb2, 0xa0, 0x66, 0x8e, 0xfa, 0xc1, 0xb8, 0xe8, 0x36, 0x94, 0x4d, 0xf8, 0x19, 0x0f, 0x35,
	0xe7, 0x5e, 0x19, 0x4d, 0x37, 0xbb, 0xff, 0xb4, 0xc0, 0x7e, 0x2a, 0x94, 0x1f, 0x97, 0x5d, 0xfd,
	0x18, 0xbe, 0x2e, 0xb9, 0x24, 0x97, 0x3e, 0x91, 0xe4, 0x7f, 0x59, 0xd0, 0xc8, 0x62, 0xbb, 0x46,
	0x9e, 0xbf, 0xbf, 0xb8, 0xe6, 0x2b, 0x56, 0xfc, 0x70, 0xc5, 0x4a, 0xf9, 0x8a, 0xf5, 0xa0, 0x16,
	0x66, 0x95, 0x68, 0x97, 0xe7, 0x87, 0x47, 0x16, 0x06, 0x9d, 0x41, 0xba, 0xdf, 0x2c, 0x40, 0xeb,
	0x09, 0x9e, 0x25, 0x2f, 0xd4, 0xe7, 0xe3, 0x00, 0xc9, 0xed, 0xf7, 0x5a, 0xdf, 0x32, 0x27, 0x7c,
	0xb4, 0xb9, 0x17, 0x3e, 0x37, 0xef, 0x8b, 0x1f, 0xcf, 0x7b, 0x3c, 0x06, 0x13, 0x36, 0x32, 0x23,
	0xa2, 0x68, 0x18, 0x68, 0x79, 0xae, 0xfd, 0xcc, 0x70, 0xa8, 0x89, 0x6c, 0x49, 0x7e, 0x91, 0x59,
	0x18, 0xca, 0x2e, 0x5d, 0x41, 0xd9, 0x09, 0x38, 0x5e, 0x92, 0x7b, 0xd0, 0x8c, 0xf9, 0x75, 0xe4,
	0xf3, 0xd1, 0x73, 0x21, 0x5d, 0xbc, 0x3c, 0x46, 0x62, 0x26, 0x1d, 0xf0, 0x87, 0xf1, 0x16, 0xad,
	0xeb, 0x99, 0x40, 0xd6, 0xa0, 0x3a, 0xe5, 0xf0, 0x8a, 0xe1, 0xf0, 0xa9, 0x4c, 0xb6, 0xa0, 0x14,
	0x4a, 0xdf, 0x4d, 0xe6, 0xc5, 0x55, 0x5c, 0xd1, 0x17, 0xd1, 0x51, 0x80, 0x09, 0x57, 0x24, 0xd0,
	0x3b, 0xff, 0xb7, 0xa0, 0x39, 0x37, 0xef, 0x48, 0x05, 0x16, 0x9f, 0xe0, 0x99, 0x5d, 0x20, 0x55,
	0x28, 0x3e, 0xf4, 0x83, 0xc0, 0xb6, 0xe2, 0x55, 0x5f, 0x70, 0xb4, 0x17, 0x48, 0x03, 0xaa, 0x7b,
	0x8c, 0xbb, 0x18, 0xa0, 0x67, 0x2f, 0xc6, 0x12, 0xc5, 0x30, 0x60, 0x2e, 0x7a, 0x76, 0x91, 0x2c,
	0x43, 0xf3, 0x29, 0x72, 0xcf, 0xe7, 0xe3, 0x04, 0x62, 0x97, 0x48, 0x1d, 0x2a, 0x43, 0x2d, 0xc2,
	0x10, 0x3d, 0xbb, 0x9c, 0xa0, 0x63, 0x76, 0x46, 0xcf, 0xae, 0x90, 0x26, 0xd4, 0x86, 0x91, 0x0a,
	0x91, 0x7b, 0xe8, 0xd9, 0x55, 0xd2, 0x02, 0x48, 0x8d, 0xe3, 0xc3, 0x6b, 0xb1, 0xbc, 0xc7, 0x02,
	0x37, 0x0a, 0xe2, 0x5f, 0x98, 0x0d, 0xb1, 0xa7, 0xfd, 0x97, 0xa1, 0x2f, 0xd1, 0xb3, 0xeb, 0x84,
	0x40, 0x2b, 0x05, 0xa7, 0xc7, 0xdb, 0x0d, 0x52, 0x83, 0x92, 0xf9, 0xac, 0xd9, 0x4d, 0xb2, 0x04,
	0xf5, 0x1c, 0xe5, 0xdb, 0xad, 0xd8, 0xd9, 0x10, 0xb5, 0x0e, 0xcc, 0xef, 0xd4, 0x5e, 0xba, 0xf3,
	0x77, 0x0b, 0x56, 0xae, 0x9a, 0x21, 0x71, 0x08, 0xcf, 0xf8, 0x09, 0x17, 0x67, 0x7c, 0x78, 0x3e,
	0x39, 0x12, 0x81, 0x5d, 0x20, 0x5d, 0x70, 0x9e, 0x71, 0x15, 0x85, 0xf1, 0xb7, 0x06, 0x3d, 0x63,
	0xb5, 0x77, 0xcc, 0x24, 0x73, 0x35, 0x4a, 0x5f, 0x69, 0xdf, 0xb5, 0x2d, 0x72, 0x13, 0x96, 0x07,
	0xdc, 0x15, 0x52, 0xa2, 0xab, 0xb3, 0x71, 0x6a, 0x2f, 0xc4, 0xd7, 0xdc, 0x7f, 0xe9, 0x1e, 0x33,
	0x3e, 0xc6, 0xbd, 0x40, 0x28, 0x93, 0x32, 0x02, 0xad, 0x7e, 0x14, 0x06, 0xbe, 0xcb, 0x34, 0x1a,
	0x67, 0x76, 0x31, 0xd6, 0x0d, 0xf8, 0x0b, 0x16, 0xf8, 0x5e, 0xda, 0x86, 0x76, 0x29, 0x0e, 0xe7,
	0x40, 0x1f, 0xa3, 0xb4, 0xcb, 0xbb, 0x77, 0x5f, 0xbd, 0x71, 0x0a, 0xaf, 0xdf, 0x38, 0x85, 0x77,
	0x6f, 0x1c, 0xeb, 0x8f, 0x17, 0x8e, 0xf5, 0xb7, 0x0b, 0xc7, 0xfa, 0xc7, 0x85, 0x63, 0xbd, 0xba,
	0x70, 0xac, 0x7f, 0x5f, 0x38, 0xd6, 0x7f, 0x2e, 0x9c, 0xc2, 0xbb, 0x0b, 0xc7, 0xfa, 0xf3, 0x5b,
	0xa7, 0xf0, 0xea, 0xad, 0x53, 0x78, 0xfd, 0xd6, 0x29, 0x1c, 0x95, 0x4d, 0xd5, 0xb7, 0xbf, 0x0d,
	0x00, 0x00, 0xff, 0xff, 0x8d, 0xe7, 0x95, 0xea, 0x67, 0x0f, 0x00, 0x00,
}

func (x ExecutionType) String() string {
	s, ok := ExecutionType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x OrderRejectionReason) String() string {
	s, ok := OrderRejectionReason_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (this *BusinessRequestReject) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*BusinessRequestReject)
	if !ok {
		that2, ok := that.(BusinessRequestReject)
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
	if !this.Subscriber.Equal(that1.Subscriber) {
		return false
	}
	if !this.Instrument.Equal(that1.Instrument) {
		return false
	}
	if this.Aggregation != that1.Aggregation {
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
	if !this.Subscriber.Equal(that1.Subscriber) {
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
func (this *ExecutionReport) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ExecutionReport)
	if !ok {
		that2, ok := that.(ExecutionReport)
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
	if this.SeqNum != that1.SeqNum {
		return false
	}
	if this.OrderID != that1.OrderID {
		return false
	}
	if !this.ClientOrderID.Equal(that1.ClientOrderID) {
		return false
	}
	if this.ExecutionID != that1.ExecutionID {
		return false
	}
	if this.ExecutionType != that1.ExecutionType {
		return false
	}
	if this.OrderStatus != that1.OrderStatus {
		return false
	}
	if !this.Instrument.Equal(that1.Instrument) {
		return false
	}
	if this.Side != that1.Side {
		return false
	}
	if this.LeavesQuantity != that1.LeavesQuantity {
		return false
	}
	if this.CumQuantity != that1.CumQuantity {
		return false
	}
	if !this.TransactionTime.Equal(that1.TransactionTime) {
		return false
	}
	if !this.TradeID.Equal(that1.TradeID) {
		return false
	}
	if this.OrderRejectionReason != that1.OrderRejectionReason {
		return false
	}
	return true
}
func (this *OrderStatusRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OrderStatusRequest)
	if !ok {
		that2, ok := that.(OrderStatusRequest)
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
	if !this.Subscriber.Equal(that1.Subscriber) {
		return false
	}
	if !this.OrderID.Equal(that1.OrderID) {
		return false
	}
	if !this.ClientOrderID.Equal(that1.ClientOrderID) {
		return false
	}
	if !this.Instrument.Equal(that1.Instrument) {
		return false
	}
	if !this.Account.Equal(that1.Account) {
		return false
	}
	return true
}
func (this *OrderList) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OrderList)
	if !ok {
		that2, ok := that.(OrderList)
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
	if len(this.Orders) != len(that1.Orders) {
		return false
	}
	for i := range this.Orders {
		if !this.Orders[i].Equal(that1.Orders[i]) {
			return false
		}
	}
	return true
}
func (this *PositionsRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PositionsRequest)
	if !ok {
		that2, ok := that.(PositionsRequest)
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
	if !this.Subscriber.Equal(that1.Subscriber) {
		return false
	}
	if !this.Instrument.Equal(that1.Instrument) {
		return false
	}
	if !this.Account.Equal(that1.Account) {
		return false
	}
	return true
}
func (this *PositionList) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PositionList)
	if !ok {
		that2, ok := that.(PositionList)
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
	if !this.Subscriber.Equal(that1.Subscriber) {
		return false
	}
	if this.ResponseID != that1.ResponseID {
		return false
	}
	if this.Error != that1.Error {
		return false
	}
	if len(this.Positions) != len(that1.Positions) {
		return false
	}
	for i := range this.Positions {
		if !this.Positions[i].Equal(that1.Positions[i]) {
			return false
		}
	}
	return true
}
func (this *NewOrderSingle) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*NewOrderSingle)
	if !ok {
		that2, ok := that.(NewOrderSingle)
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
	if this.ClientOrderID != that1.ClientOrderID {
		return false
	}
	if !this.Instrument.Equal(that1.Instrument) {
		return false
	}
	if !this.Account.Equal(that1.Account) {
		return false
	}
	if this.OrderType != that1.OrderType {
		return false
	}
	if this.OrderSide != that1.OrderSide {
		return false
	}
	if this.TimeInForce != that1.TimeInForce {
		return false
	}
	if this.Quantity != that1.Quantity {
		return false
	}
	if !this.Price.Equal(that1.Price) {
		return false
	}
	return true
}
func (this *BusinessRequestReject) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&messages.BusinessRequestReject{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Reason: "+fmt.Sprintf("%#v", this.Reason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *MarketDataRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.MarketDataRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Subscribe: "+fmt.Sprintf("%#v", this.Subscribe)+",\n")
	if this.Subscriber != nil {
		s = append(s, "Subscriber: "+fmt.Sprintf("%#v", this.Subscriber)+",\n")
	}
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	s = append(s, "Aggregation: "+fmt.Sprintf("%#v", this.Aggregation)+",\n")
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
	s := make([]string, 0, 7)
	s = append(s, "&messages.SecurityListRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Subscribe: "+fmt.Sprintf("%#v", this.Subscribe)+",\n")
	if this.Subscriber != nil {
		s = append(s, "Subscriber: "+fmt.Sprintf("%#v", this.Subscriber)+",\n")
	}
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
func (this *ExecutionReport) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 17)
	s = append(s, "&messages.ExecutionReport{")
	s = append(s, "SeqNum: "+fmt.Sprintf("%#v", this.SeqNum)+",\n")
	s = append(s, "OrderID: "+fmt.Sprintf("%#v", this.OrderID)+",\n")
	if this.ClientOrderID != nil {
		s = append(s, "ClientOrderID: "+fmt.Sprintf("%#v", this.ClientOrderID)+",\n")
	}
	s = append(s, "ExecutionID: "+fmt.Sprintf("%#v", this.ExecutionID)+",\n")
	s = append(s, "ExecutionType: "+fmt.Sprintf("%#v", this.ExecutionType)+",\n")
	s = append(s, "OrderStatus: "+fmt.Sprintf("%#v", this.OrderStatus)+",\n")
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	s = append(s, "Side: "+fmt.Sprintf("%#v", this.Side)+",\n")
	s = append(s, "LeavesQuantity: "+fmt.Sprintf("%#v", this.LeavesQuantity)+",\n")
	s = append(s, "CumQuantity: "+fmt.Sprintf("%#v", this.CumQuantity)+",\n")
	if this.TransactionTime != nil {
		s = append(s, "TransactionTime: "+fmt.Sprintf("%#v", this.TransactionTime)+",\n")
	}
	if this.TradeID != nil {
		s = append(s, "TradeID: "+fmt.Sprintf("%#v", this.TradeID)+",\n")
	}
	s = append(s, "OrderRejectionReason: "+fmt.Sprintf("%#v", this.OrderRejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OrderStatusRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 11)
	s = append(s, "&messages.OrderStatusRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Subscribe: "+fmt.Sprintf("%#v", this.Subscribe)+",\n")
	if this.Subscriber != nil {
		s = append(s, "Subscriber: "+fmt.Sprintf("%#v", this.Subscriber)+",\n")
	}
	if this.OrderID != nil {
		s = append(s, "OrderID: "+fmt.Sprintf("%#v", this.OrderID)+",\n")
	}
	if this.ClientOrderID != nil {
		s = append(s, "ClientOrderID: "+fmt.Sprintf("%#v", this.ClientOrderID)+",\n")
	}
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	if this.Account != nil {
		s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OrderList) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.OrderList{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	s = append(s, "Error: "+fmt.Sprintf("%#v", this.Error)+",\n")
	if this.Orders != nil {
		s = append(s, "Orders: "+fmt.Sprintf("%#v", this.Orders)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *PositionsRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.PositionsRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Subscribe: "+fmt.Sprintf("%#v", this.Subscribe)+",\n")
	if this.Subscriber != nil {
		s = append(s, "Subscriber: "+fmt.Sprintf("%#v", this.Subscriber)+",\n")
	}
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	if this.Account != nil {
		s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *PositionList) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 10)
	s = append(s, "&messages.PositionList{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Subscribe: "+fmt.Sprintf("%#v", this.Subscribe)+",\n")
	if this.Subscriber != nil {
		s = append(s, "Subscriber: "+fmt.Sprintf("%#v", this.Subscriber)+",\n")
	}
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	s = append(s, "Error: "+fmt.Sprintf("%#v", this.Error)+",\n")
	if this.Positions != nil {
		s = append(s, "Positions: "+fmt.Sprintf("%#v", this.Positions)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *NewOrderSingle) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 12)
	s = append(s, "&messages.NewOrderSingle{")
	s = append(s, "ClientOrderID: "+fmt.Sprintf("%#v", this.ClientOrderID)+",\n")
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	if this.Account != nil {
		s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
	}
	s = append(s, "OrderType: "+fmt.Sprintf("%#v", this.OrderType)+",\n")
	s = append(s, "OrderSide: "+fmt.Sprintf("%#v", this.OrderSide)+",\n")
	s = append(s, "TimeInForce: "+fmt.Sprintf("%#v", this.TimeInForce)+",\n")
	s = append(s, "Quantity: "+fmt.Sprintf("%#v", this.Quantity)+",\n")
	if this.Price != nil {
		s = append(s, "Price: "+fmt.Sprintf("%#v", this.Price)+",\n")
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
func (m *BusinessRequestReject) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BusinessRequestReject) MarshalTo(dAtA []byte) (int, error) {
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
	if m.Subscriber != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Subscriber.Size()))
		n1, err := m.Subscriber.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.Instrument != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Instrument.Size()))
		n2, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if m.Aggregation != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Aggregation))
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
	if m.Subscriber != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Subscriber.Size()))
		n7, err := m.Subscriber.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n7
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

func (m *ExecutionReport) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExecutionReport) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.SeqNum != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SeqNum))
	}
	if len(m.OrderID) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.OrderID)))
		i += copy(dAtA[i:], m.OrderID)
	}
	if m.ClientOrderID != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ClientOrderID.Size()))
		n8, err := m.ClientOrderID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n8
	}
	if len(m.ExecutionID) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.ExecutionID)))
		i += copy(dAtA[i:], m.ExecutionID)
	}
	if m.ExecutionType != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ExecutionType))
	}
	if m.OrderStatus != 0 {
		dAtA[i] = 0x30
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderStatus))
	}
	if m.Instrument != nil {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Instrument.Size()))
		n9, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n9
	}
	if m.Side != 0 {
		dAtA[i] = 0x40
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Side))
	}
	if m.LeavesQuantity != 0 {
		dAtA[i] = 0x49
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.LeavesQuantity))))
		i += 8
	}
	if m.CumQuantity != 0 {
		dAtA[i] = 0x51
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.CumQuantity))))
		i += 8
	}
	if m.TransactionTime != nil {
		dAtA[i] = 0x5a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.TransactionTime.Size()))
		n10, err := m.TransactionTime.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n10
	}
	if m.TradeID != nil {
		dAtA[i] = 0x62
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.TradeID.Size()))
		n11, err := m.TradeID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n11
	}
	if m.OrderRejectionReason != 0 {
		dAtA[i] = 0x68
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderRejectionReason))
	}
	return i, nil
}

func (m *OrderStatusRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderStatusRequest) MarshalTo(dAtA []byte) (int, error) {
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
	if m.Subscriber != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Subscriber.Size()))
		n12, err := m.Subscriber.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n12
	}
	if m.OrderID != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderID.Size()))
		n13, err := m.OrderID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n13
	}
	if m.ClientOrderID != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ClientOrderID.Size()))
		n14, err := m.ClientOrderID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n14
	}
	if m.Instrument != nil {
		dAtA[i] = 0x32
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Instrument.Size()))
		n15, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n15
	}
	if m.Account != nil {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n16, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n16
	}
	return i, nil
}

func (m *OrderList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderList) MarshalTo(dAtA []byte) (int, error) {
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
	if len(m.Orders) > 0 {
		for _, msg := range m.Orders {
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

func (m *PositionsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PositionsRequest) MarshalTo(dAtA []byte) (int, error) {
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
	if m.Subscriber != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Subscriber.Size()))
		n17, err := m.Subscriber.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n17
	}
	if m.Instrument != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Instrument.Size()))
		n18, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n18
	}
	if m.Account != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n19, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n19
	}
	return i, nil
}

func (m *PositionList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PositionList) MarshalTo(dAtA []byte) (int, error) {
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
	if m.Subscriber != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Subscriber.Size()))
		n20, err := m.Subscriber.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n20
	}
	if m.ResponseID != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
	}
	if len(m.Error) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.Error)))
		i += copy(dAtA[i:], m.Error)
	}
	if len(m.Positions) > 0 {
		for _, msg := range m.Positions {
			dAtA[i] = 0x32
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

func (m *NewOrderSingle) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NewOrderSingle) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ClientOrderID) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.ClientOrderID)))
		i += copy(dAtA[i:], m.ClientOrderID)
	}
	if m.Instrument != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Instrument.Size()))
		n21, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n21
	}
	if m.Account != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n22, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n22
	}
	if m.OrderType != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderType))
	}
	if m.OrderSide != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderSide))
	}
	if m.TimeInForce != 0 {
		dAtA[i] = 0x30
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.TimeInForce))
	}
	if m.Quantity != 0 {
		dAtA[i] = 0x39
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Quantity))))
		i += 8
	}
	if m.Price != nil {
		dAtA[i] = 0x42
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Price.Size()))
		n23, err := m.Price.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n23
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
func (m *BusinessRequestReject) Size() (n int) {
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
	if m.Subscriber != nil {
		l = m.Subscriber.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Instrument != nil {
		l = m.Instrument.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Aggregation != 0 {
		n += 1 + sovExecutorMessages(uint64(m.Aggregation))
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
	if m.Subscriber != nil {
		l = m.Subscriber.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
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

func (m *ExecutionReport) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SeqNum != 0 {
		n += 1 + sovExecutorMessages(uint64(m.SeqNum))
	}
	l = len(m.OrderID)
	if l > 0 {
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.ClientOrderID != nil {
		l = m.ClientOrderID.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	l = len(m.ExecutionID)
	if l > 0 {
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.ExecutionType != 0 {
		n += 1 + sovExecutorMessages(uint64(m.ExecutionType))
	}
	if m.OrderStatus != 0 {
		n += 1 + sovExecutorMessages(uint64(m.OrderStatus))
	}
	if m.Instrument != nil {
		l = m.Instrument.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Side != 0 {
		n += 1 + sovExecutorMessages(uint64(m.Side))
	}
	if m.LeavesQuantity != 0 {
		n += 9
	}
	if m.CumQuantity != 0 {
		n += 9
	}
	if m.TransactionTime != nil {
		l = m.TransactionTime.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.TradeID != nil {
		l = m.TradeID.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.OrderRejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.OrderRejectionReason))
	}
	return n
}

func (m *OrderStatusRequest) Size() (n int) {
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
	if m.Subscriber != nil {
		l = m.Subscriber.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.OrderID != nil {
		l = m.OrderID.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.ClientOrderID != nil {
		l = m.ClientOrderID.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Instrument != nil {
		l = m.Instrument.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Account != nil {
		l = m.Account.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *OrderList) Size() (n int) {
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
	if len(m.Orders) > 0 {
		for _, e := range m.Orders {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	return n
}

func (m *PositionsRequest) Size() (n int) {
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
	if m.Subscriber != nil {
		l = m.Subscriber.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Instrument != nil {
		l = m.Instrument.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Account != nil {
		l = m.Account.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *PositionList) Size() (n int) {
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
	if m.Subscriber != nil {
		l = m.Subscriber.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.ResponseID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.ResponseID))
	}
	l = len(m.Error)
	if l > 0 {
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if len(m.Positions) > 0 {
		for _, e := range m.Positions {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	return n
}

func (m *NewOrderSingle) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ClientOrderID)
	if l > 0 {
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Instrument != nil {
		l = m.Instrument.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Account != nil {
		l = m.Account.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.OrderType != 0 {
		n += 1 + sovExecutorMessages(uint64(m.OrderType))
	}
	if m.OrderSide != 0 {
		n += 1 + sovExecutorMessages(uint64(m.OrderSide))
	}
	if m.TimeInForce != 0 {
		n += 1 + sovExecutorMessages(uint64(m.TimeInForce))
	}
	if m.Quantity != 0 {
		n += 9
	}
	if m.Price != nil {
		l = m.Price.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
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
func (this *BusinessRequestReject) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&BusinessRequestReject{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Reason:` + fmt.Sprintf("%v", this.Reason) + `,`,
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
		`Subscriber:` + strings.Replace(fmt.Sprintf("%v", this.Subscriber), "PID", "actor.PID", 1) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`Aggregation:` + fmt.Sprintf("%v", this.Aggregation) + `,`,
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
		`SnapshotL2:` + strings.Replace(fmt.Sprintf("%v", this.SnapshotL2), "OBL2Snapshot", "models.OBL2Snapshot", 1) + `,`,
		`SnapshotL3:` + strings.Replace(fmt.Sprintf("%v", this.SnapshotL3), "OBL3Snapshot", "models.OBL3Snapshot", 1) + `,`,
		`Trades:` + strings.Replace(fmt.Sprintf("%v", this.Trades), "AggregatedTrade", "models.AggregatedTrade", 1) + `,`,
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
		`UpdateL2:` + strings.Replace(fmt.Sprintf("%v", this.UpdateL2), "OBL2Update", "models.OBL2Update", 1) + `,`,
		`UpdateL3:` + strings.Replace(fmt.Sprintf("%v", this.UpdateL3), "OBL3Update", "models.OBL3Update", 1) + `,`,
		`Trades:` + strings.Replace(fmt.Sprintf("%v", this.Trades), "AggregatedTrade", "models.AggregatedTrade", 1) + `,`,
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
		`Subscriber:` + strings.Replace(fmt.Sprintf("%v", this.Subscriber), "PID", "actor.PID", 1) + `,`,
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
		`Securities:` + strings.Replace(fmt.Sprintf("%v", this.Securities), "Security", "models.Security", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ExecutionReport) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ExecutionReport{`,
		`SeqNum:` + fmt.Sprintf("%v", this.SeqNum) + `,`,
		`OrderID:` + fmt.Sprintf("%v", this.OrderID) + `,`,
		`ClientOrderID:` + strings.Replace(fmt.Sprintf("%v", this.ClientOrderID), "StringValue", "types.StringValue", 1) + `,`,
		`ExecutionID:` + fmt.Sprintf("%v", this.ExecutionID) + `,`,
		`ExecutionType:` + fmt.Sprintf("%v", this.ExecutionType) + `,`,
		`OrderStatus:` + fmt.Sprintf("%v", this.OrderStatus) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`Side:` + fmt.Sprintf("%v", this.Side) + `,`,
		`LeavesQuantity:` + fmt.Sprintf("%v", this.LeavesQuantity) + `,`,
		`CumQuantity:` + fmt.Sprintf("%v", this.CumQuantity) + `,`,
		`TransactionTime:` + strings.Replace(fmt.Sprintf("%v", this.TransactionTime), "Timestamp", "types.Timestamp", 1) + `,`,
		`TradeID:` + strings.Replace(fmt.Sprintf("%v", this.TradeID), "StringValue", "types.StringValue", 1) + `,`,
		`OrderRejectionReason:` + fmt.Sprintf("%v", this.OrderRejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrderStatusRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OrderStatusRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Subscribe:` + fmt.Sprintf("%v", this.Subscribe) + `,`,
		`Subscriber:` + strings.Replace(fmt.Sprintf("%v", this.Subscriber), "PID", "actor.PID", 1) + `,`,
		`OrderID:` + strings.Replace(fmt.Sprintf("%v", this.OrderID), "StringValue", "types.StringValue", 1) + `,`,
		`ClientOrderID:` + strings.Replace(fmt.Sprintf("%v", this.ClientOrderID), "StringValue", "types.StringValue", 1) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrderList) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OrderList{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Error:` + fmt.Sprintf("%v", this.Error) + `,`,
		`Orders:` + strings.Replace(fmt.Sprintf("%v", this.Orders), "Order", "models.Order", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *PositionsRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&PositionsRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Subscribe:` + fmt.Sprintf("%v", this.Subscribe) + `,`,
		`Subscriber:` + strings.Replace(fmt.Sprintf("%v", this.Subscriber), "PID", "actor.PID", 1) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *PositionList) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&PositionList{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Subscribe:` + fmt.Sprintf("%v", this.Subscribe) + `,`,
		`Subscriber:` + strings.Replace(fmt.Sprintf("%v", this.Subscriber), "PID", "actor.PID", 1) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Error:` + fmt.Sprintf("%v", this.Error) + `,`,
		`Positions:` + strings.Replace(fmt.Sprintf("%v", this.Positions), "Position", "models.Position", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *NewOrderSingle) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&NewOrderSingle{`,
		`ClientOrderID:` + fmt.Sprintf("%v", this.ClientOrderID) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`OrderType:` + fmt.Sprintf("%v", this.OrderType) + `,`,
		`OrderSide:` + fmt.Sprintf("%v", this.OrderSide) + `,`,
		`TimeInForce:` + fmt.Sprintf("%v", this.TimeInForce) + `,`,
		`Quantity:` + fmt.Sprintf("%v", this.Quantity) + `,`,
		`Price:` + strings.Replace(fmt.Sprintf("%v", this.Price), "DoubleValue", "types.DoubleValue", 1) + `,`,
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
func (m *BusinessRequestReject) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: BusinessRequestReject: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BusinessRequestReject: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Subscriber", wireType)
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
			if m.Subscriber == nil {
				m.Subscriber = &actor.PID{}
			}
			if err := m.Subscriber.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
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
				m.Instrument = &models.Instrument{}
			}
			if err := m.Instrument.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Aggregation", wireType)
			}
			m.Aggregation = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Aggregation |= models.OrderBookAggregation(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
				m.SnapshotL2 = &models.OBL2Snapshot{}
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
				m.SnapshotL3 = &models.OBL3Snapshot{}
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
			m.Trades = append(m.Trades, &models.AggregatedTrade{})
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
				m.UpdateL2 = &models.OBL2Update{}
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
				m.UpdateL3 = &models.OBL3Update{}
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
			m.Trades = append(m.Trades, &models.AggregatedTrade{})
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
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subscriber", wireType)
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
			if m.Subscriber == nil {
				m.Subscriber = &actor.PID{}
			}
			if err := m.Subscriber.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
			m.Securities = append(m.Securities, &models.Security{})
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
func (m *ExecutionReport) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: ExecutionReport: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExecutionReport: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SeqNum", wireType)
			}
			m.SeqNum = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SeqNum |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderID", wireType)
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
			m.OrderID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientOrderID", wireType)
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
			if m.ClientOrderID == nil {
				m.ClientOrderID = &types.StringValue{}
			}
			if err := m.ClientOrderID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExecutionID", wireType)
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
			m.ExecutionID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExecutionType", wireType)
			}
			m.ExecutionType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExecutionType |= ExecutionType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderStatus", wireType)
			}
			m.OrderStatus = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrderStatus |= models.OrderStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
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
				m.Instrument = &models.Instrument{}
			}
			if err := m.Instrument.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Side", wireType)
			}
			m.Side = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Side |= models.Side(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
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
		case 10:
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
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransactionTime", wireType)
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
			if m.TransactionTime == nil {
				m.TransactionTime = &types.Timestamp{}
			}
			if err := m.TransactionTime.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TradeID", wireType)
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
			if m.TradeID == nil {
				m.TradeID = &types.StringValue{}
			}
			if err := m.TradeID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 13:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderRejectionReason", wireType)
			}
			m.OrderRejectionReason = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrderRejectionReason |= OrderRejectionReason(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
func (m *OrderStatusRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OrderStatusRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderStatusRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Subscriber", wireType)
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
			if m.Subscriber == nil {
				m.Subscriber = &actor.PID{}
			}
			if err := m.Subscriber.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderID", wireType)
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
			if m.OrderID == nil {
				m.OrderID = &types.StringValue{}
			}
			if err := m.OrderID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientOrderID", wireType)
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
			if m.ClientOrderID == nil {
				m.ClientOrderID = &types.StringValue{}
			}
			if err := m.ClientOrderID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
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
				m.Instrument = &models.Instrument{}
			}
			if err := m.Instrument.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
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
			if m.Account == nil {
				m.Account = &models.Account{}
			}
			if err := m.Account.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *OrderList) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OrderList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderList: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Orders", wireType)
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
			m.Orders = append(m.Orders, &models.Order{})
			if err := m.Orders[len(m.Orders)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *PositionsRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: PositionsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PositionsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Subscriber", wireType)
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
			if m.Subscriber == nil {
				m.Subscriber = &actor.PID{}
			}
			if err := m.Subscriber.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
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
				m.Instrument = &models.Instrument{}
			}
			if err := m.Instrument.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
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
			if m.Account == nil {
				m.Account = &models.Account{}
			}
			if err := m.Account.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *PositionList) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: PositionList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PositionList: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Subscriber", wireType)
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
			if m.Subscriber == nil {
				m.Subscriber = &actor.PID{}
			}
			if err := m.Subscriber.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
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
		case 5:
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
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Positions", wireType)
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
			m.Positions = append(m.Positions, &models.Position{})
			if err := m.Positions[len(m.Positions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *NewOrderSingle) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: NewOrderSingle: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NewOrderSingle: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientOrderID", wireType)
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
			m.ClientOrderID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
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
				m.Instrument = &models.Instrument{}
			}
			if err := m.Instrument.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
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
			if m.Account == nil {
				m.Account = &models.Account{}
			}
			if err := m.Account.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderType", wireType)
			}
			m.OrderType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrderType |= models.OrderType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderSide", wireType)
			}
			m.OrderSide = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrderSide |= models.Side(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TimeInForce", wireType)
			}
			m.TimeInForce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TimeInForce |= models.TimeInForce(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
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
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
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
			if m.Price == nil {
				m.Price = &types.DoubleValue{}
			}
			if err := m.Price.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
