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
	Done           ExecutionType = 1
	Canceled       ExecutionType = 2
	Replaced       ExecutionType = 3
	PendingCancel  ExecutionType = 4
	Stopped        ExecutionType = 5
	Rejected       ExecutionType = 6
	Suspended      ExecutionType = 7
	PendingNew     ExecutionType = 8
	Calculated     ExecutionType = 9
	Expired        ExecutionType = 10
	PendingReplace ExecutionType = 11
	Trade          ExecutionType = 12
	OrderStatus    ExecutionType = 13
	Settlement     ExecutionType = 14
)

var ExecutionType_name = map[int32]string{
	0:  "New",
	1:  "Done",
	2:  "Canceled",
	3:  "Replaced",
	4:  "PendingCancel",
	5:  "Stopped",
	6:  "Rejected",
	7:  "Suspended",
	8:  "PendingNew",
	9:  "Calculated",
	10: "Expired",
	11: "PendingReplace",
	12: "Trade",
	13: "OrderStatus",
	14: "Settlement",
}

var ExecutionType_value = map[string]int32{
	"New":            0,
	"Done":           1,
	"Canceled":       2,
	"Replaced":       3,
	"PendingCancel":  4,
	"Stopped":        5,
	"Rejected":       6,
	"Suspended":      7,
	"PendingNew":     8,
	"Calculated":     9,
	"Expired":        10,
	"PendingReplace": 11,
	"Trade":          12,
	"OrderStatus":    13,
	"Settlement":     14,
}

func (ExecutionType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{0}
}

type RejectionReason int32

const (
	Other                          RejectionReason = 0
	UnknownSymbol                  RejectionReason = 1
	UnknownSecurityID              RejectionReason = 2
	UnsupportedOrderCharacteristic RejectionReason = 3
	IncorrectQuantity              RejectionReason = 4
	ExchangeClosed                 RejectionReason = 5
	DuplicateOrder                 RejectionReason = 6
	InvalidAccount                 RejectionReason = 7
	TooLateToCancel                RejectionReason = 8
	UnknownOrder                   RejectionReason = 9
	CancelAlreadyPending           RejectionReason = 10
	DifferentSymbols               RejectionReason = 11
	InvalidRequest                 RejectionReason = 12
	ExchangeAPIError               RejectionReason = 13
)

var RejectionReason_name = map[int32]string{
	0:  "Other",
	1:  "UnknownSymbol",
	2:  "UnknownSecurityID",
	3:  "UnsupportedOrderCharacteristic",
	4:  "IncorrectQuantity",
	5:  "ExchangeClosed",
	6:  "DuplicateOrder",
	7:  "InvalidAccount",
	8:  "TooLateToCancel",
	9:  "UnknownOrder",
	10: "CancelAlreadyPending",
	11: "DifferentSymbols",
	12: "InvalidRequest",
	13: "ExchangeAPIError",
}

var RejectionReason_value = map[string]int32{
	"Other":                          0,
	"UnknownSymbol":                  1,
	"UnknownSecurityID":              2,
	"UnsupportedOrderCharacteristic": 3,
	"IncorrectQuantity":              4,
	"ExchangeClosed":                 5,
	"DuplicateOrder":                 6,
	"InvalidAccount":                 7,
	"TooLateToCancel":                8,
	"UnknownOrder":                   9,
	"CancelAlreadyPending":           10,
	"DifferentSymbols":               11,
	"InvalidRequest":                 12,
	"ExchangeAPIError":               13,
}

func (RejectionReason) EnumDescriptor() ([]byte, []int) {
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
	SeqNum          uint64             `protobuf:"varint,1,opt,name=seq_num,json=seqNum,proto3" json:"seq_num,omitempty"`
	OrderID         string             `protobuf:"bytes,2,opt,name=orderID,proto3" json:"orderID,omitempty"`
	ClientOrderID   *types.StringValue `protobuf:"bytes,3,opt,name=client_orderID,json=clientOrderID,proto3" json:"client_orderID,omitempty"`
	ExecutionID     string             `protobuf:"bytes,4,opt,name=executionID,proto3" json:"executionID,omitempty"`
	ExecutionType   ExecutionType      `protobuf:"varint,5,opt,name=execution_type,json=executionType,proto3,enum=messages.ExecutionType" json:"execution_type,omitempty"`
	OrderStatus     models.OrderStatus `protobuf:"varint,6,opt,name=order_status,json=orderStatus,proto3,enum=models.OrderStatus" json:"order_status,omitempty"`
	Instrument      *models.Instrument `protobuf:"bytes,7,opt,name=instrument,proto3" json:"instrument,omitempty"`
	LeavesQuantity  float64            `protobuf:"fixed64,9,opt,name=leaves_quantity,json=leavesQuantity,proto3" json:"leaves_quantity,omitempty"`
	CumQuantity     float64            `protobuf:"fixed64,10,opt,name=cum_quantity,json=cumQuantity,proto3" json:"cum_quantity,omitempty"`
	TransactionTime *types.Timestamp   `protobuf:"bytes,11,opt,name=transaction_time,json=transactionTime,proto3" json:"transaction_time,omitempty"`
	TradeID         *types.StringValue `protobuf:"bytes,12,opt,name=tradeID,proto3" json:"tradeID,omitempty"`
	FillPrice       *types.DoubleValue `protobuf:"bytes,13,opt,name=fill_price,json=fillPrice,proto3" json:"fill_price,omitempty"`
	FillQuantity    *types.DoubleValue `protobuf:"bytes,14,opt,name=fill_quantity,json=fillQuantity,proto3" json:"fill_quantity,omitempty"`
	RejectionReason RejectionReason    `protobuf:"varint,15,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
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

func (m *ExecutionReport) GetFillPrice() *types.DoubleValue {
	if m != nil {
		return m.FillPrice
	}
	return nil
}

func (m *ExecutionReport) GetFillQuantity() *types.DoubleValue {
	if m != nil {
		return m.FillQuantity
	}
	return nil
}

func (m *ExecutionReport) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type SideValue struct {
	Value models.Side `protobuf:"varint,1,opt,name=value,proto3,enum=models.Side" json:"value,omitempty"`
}

func (m *SideValue) Reset()      { *m = SideValue{} }
func (*SideValue) ProtoMessage() {}
func (*SideValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{8}
}
func (m *SideValue) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SideValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SideValue.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SideValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SideValue.Merge(m, src)
}
func (m *SideValue) XXX_Size() int {
	return m.Size()
}
func (m *SideValue) XXX_DiscardUnknown() {
	xxx_messageInfo_SideValue.DiscardUnknown(m)
}

var xxx_messageInfo_SideValue proto.InternalMessageInfo

func (m *SideValue) GetValue() models.Side {
	if m != nil {
		return m.Value
	}
	return models.Buy
}

type OrderStatusValue struct {
	Value models.OrderStatus `protobuf:"varint,1,opt,name=value,proto3,enum=models.OrderStatus" json:"value,omitempty"`
}

func (m *OrderStatusValue) Reset()      { *m = OrderStatusValue{} }
func (*OrderStatusValue) ProtoMessage() {}
func (*OrderStatusValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{9}
}
func (m *OrderStatusValue) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderStatusValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderStatusValue.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderStatusValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderStatusValue.Merge(m, src)
}
func (m *OrderStatusValue) XXX_Size() int {
	return m.Size()
}
func (m *OrderStatusValue) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderStatusValue.DiscardUnknown(m)
}

var xxx_messageInfo_OrderStatusValue proto.InternalMessageInfo

func (m *OrderStatusValue) GetValue() models.OrderStatus {
	if m != nil {
		return m.Value
	}
	return models.New
}

type OrderFilter struct {
	OrderID       *types.StringValue `protobuf:"bytes,1,opt,name=orderID,proto3" json:"orderID,omitempty"`
	ClientOrderID *types.StringValue `protobuf:"bytes,2,opt,name=client_orderID,json=clientOrderID,proto3" json:"client_orderID,omitempty"`
	Instrument    *models.Instrument `protobuf:"bytes,3,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Side          *SideValue         `protobuf:"bytes,4,opt,name=side,proto3" json:"side,omitempty"`
	OrderStatus   *OrderStatusValue  `protobuf:"bytes,5,opt,name=order_status,json=orderStatus,proto3" json:"order_status,omitempty"`
}

func (m *OrderFilter) Reset()      { *m = OrderFilter{} }
func (*OrderFilter) ProtoMessage() {}
func (*OrderFilter) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{10}
}
func (m *OrderFilter) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderFilter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderFilter.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderFilter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderFilter.Merge(m, src)
}
func (m *OrderFilter) XXX_Size() int {
	return m.Size()
}
func (m *OrderFilter) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderFilter.DiscardUnknown(m)
}

var xxx_messageInfo_OrderFilter proto.InternalMessageInfo

func (m *OrderFilter) GetOrderID() *types.StringValue {
	if m != nil {
		return m.OrderID
	}
	return nil
}

func (m *OrderFilter) GetClientOrderID() *types.StringValue {
	if m != nil {
		return m.ClientOrderID
	}
	return nil
}

func (m *OrderFilter) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *OrderFilter) GetSide() *SideValue {
	if m != nil {
		return m.Side
	}
	return nil
}

func (m *OrderFilter) GetOrderStatus() *OrderStatusValue {
	if m != nil {
		return m.OrderStatus
	}
	return nil
}

type OrderStatusRequest struct {
	RequestID  uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe  bool            `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
	Subscriber *actor.PID      `protobuf:"bytes,3,opt,name=subscriber,proto3" json:"subscriber,omitempty"`
	Account    *models.Account `protobuf:"bytes,4,opt,name=account,proto3" json:"account,omitempty"`
	Filter     *OrderFilter    `protobuf:"bytes,5,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (m *OrderStatusRequest) Reset()      { *m = OrderStatusRequest{} }
func (*OrderStatusRequest) ProtoMessage() {}
func (*OrderStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{11}
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

func (m *OrderStatusRequest) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

func (m *OrderStatusRequest) GetFilter() *OrderFilter {
	if m != nil {
		return m.Filter
	}
	return nil
}

type OrderList struct {
	RequestID       uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Success         bool            `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	Orders          []*models.Order `protobuf:"bytes,3,rep,name=orders,proto3" json:"orders,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,4,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *OrderList) Reset()      { *m = OrderList{} }
func (*OrderList) ProtoMessage() {}
func (*OrderList) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{12}
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

func (m *OrderList) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *OrderList) GetOrders() []*models.Order {
	if m != nil {
		return m.Orders
	}
	return nil
}

func (m *OrderList) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
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
	return fileDescriptor_350c53ba9303a7e6, []int{13}
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
	return fileDescriptor_350c53ba9303a7e6, []int{14}
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

type NewOrder struct {
	ClientOrderID string             `protobuf:"bytes,1,opt,name=client_orderID,json=clientOrderID,proto3" json:"client_orderID,omitempty"`
	Instrument    *models.Instrument `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
	OrderType     models.OrderType   `protobuf:"varint,5,opt,name=order_type,json=orderType,proto3,enum=models.OrderType" json:"order_type,omitempty"`
	OrderSide     models.Side        `protobuf:"varint,6,opt,name=order_side,json=orderSide,proto3,enum=models.Side" json:"order_side,omitempty"`
	TimeInForce   models.TimeInForce `protobuf:"varint,7,opt,name=time_in_force,json=timeInForce,proto3,enum=models.TimeInForce" json:"time_in_force,omitempty"`
	Quantity      float64            `protobuf:"fixed64,8,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Price         *types.DoubleValue `protobuf:"bytes,9,opt,name=price,proto3" json:"price,omitempty"`
}

func (m *NewOrder) Reset()      { *m = NewOrder{} }
func (*NewOrder) ProtoMessage() {}
func (*NewOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{15}
}
func (m *NewOrder) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NewOrder) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NewOrder.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NewOrder) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewOrder.Merge(m, src)
}
func (m *NewOrder) XXX_Size() int {
	return m.Size()
}
func (m *NewOrder) XXX_DiscardUnknown() {
	xxx_messageInfo_NewOrder.DiscardUnknown(m)
}

var xxx_messageInfo_NewOrder proto.InternalMessageInfo

func (m *NewOrder) GetClientOrderID() string {
	if m != nil {
		return m.ClientOrderID
	}
	return ""
}

func (m *NewOrder) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *NewOrder) GetOrderType() models.OrderType {
	if m != nil {
		return m.OrderType
	}
	return models.Market
}

func (m *NewOrder) GetOrderSide() models.Side {
	if m != nil {
		return m.OrderSide
	}
	return models.Buy
}

func (m *NewOrder) GetTimeInForce() models.TimeInForce {
	if m != nil {
		return m.TimeInForce
	}
	return models.Session
}

func (m *NewOrder) GetQuantity() float64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *NewOrder) GetPrice() *types.DoubleValue {
	if m != nil {
		return m.Price
	}
	return nil
}

type NewOrderSingleRequest struct {
	RequestID uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Account   *models.Account `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	Order     *NewOrder       `protobuf:"bytes,3,opt,name=order,proto3" json:"order,omitempty"`
}

func (m *NewOrderSingleRequest) Reset()      { *m = NewOrderSingleRequest{} }
func (*NewOrderSingleRequest) ProtoMessage() {}
func (*NewOrderSingleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{16}
}
func (m *NewOrderSingleRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NewOrderSingleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NewOrderSingleRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NewOrderSingleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewOrderSingleRequest.Merge(m, src)
}
func (m *NewOrderSingleRequest) XXX_Size() int {
	return m.Size()
}
func (m *NewOrderSingleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NewOrderSingleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NewOrderSingleRequest proto.InternalMessageInfo

func (m *NewOrderSingleRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *NewOrderSingleRequest) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

func (m *NewOrderSingleRequest) GetOrder() *NewOrder {
	if m != nil {
		return m.Order
	}
	return nil
}

type NewOrderSingleResponse struct {
	RequestID       uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Success         bool            `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	OrderID         string          `protobuf:"bytes,3,opt,name=orderID,proto3" json:"orderID,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,4,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *NewOrderSingleResponse) Reset()      { *m = NewOrderSingleResponse{} }
func (*NewOrderSingleResponse) ProtoMessage() {}
func (*NewOrderSingleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{17}
}
func (m *NewOrderSingleResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NewOrderSingleResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NewOrderSingleResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NewOrderSingleResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewOrderSingleResponse.Merge(m, src)
}
func (m *NewOrderSingleResponse) XXX_Size() int {
	return m.Size()
}
func (m *NewOrderSingleResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NewOrderSingleResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NewOrderSingleResponse proto.InternalMessageInfo

func (m *NewOrderSingleResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *NewOrderSingleResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *NewOrderSingleResponse) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *NewOrderSingleResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type NewOrderBulkRequest struct {
	RequestID uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Account   *models.Account `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	Orders    []*NewOrder     `protobuf:"bytes,3,rep,name=orders,proto3" json:"orders,omitempty"`
}

func (m *NewOrderBulkRequest) Reset()      { *m = NewOrderBulkRequest{} }
func (*NewOrderBulkRequest) ProtoMessage() {}
func (*NewOrderBulkRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{18}
}
func (m *NewOrderBulkRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NewOrderBulkRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NewOrderBulkRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NewOrderBulkRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewOrderBulkRequest.Merge(m, src)
}
func (m *NewOrderBulkRequest) XXX_Size() int {
	return m.Size()
}
func (m *NewOrderBulkRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NewOrderBulkRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NewOrderBulkRequest proto.InternalMessageInfo

func (m *NewOrderBulkRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *NewOrderBulkRequest) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

func (m *NewOrderBulkRequest) GetOrders() []*NewOrder {
	if m != nil {
		return m.Orders
	}
	return nil
}

type NewOrderBulkResponse struct {
	RequestID       uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Success         bool            `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	OrderIDs        []string        `protobuf:"bytes,3,rep,name=orderIDs,proto3" json:"orderIDs,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,4,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *NewOrderBulkResponse) Reset()      { *m = NewOrderBulkResponse{} }
func (*NewOrderBulkResponse) ProtoMessage() {}
func (*NewOrderBulkResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{19}
}
func (m *NewOrderBulkResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NewOrderBulkResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NewOrderBulkResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NewOrderBulkResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewOrderBulkResponse.Merge(m, src)
}
func (m *NewOrderBulkResponse) XXX_Size() int {
	return m.Size()
}
func (m *NewOrderBulkResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NewOrderBulkResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NewOrderBulkResponse proto.InternalMessageInfo

func (m *NewOrderBulkResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *NewOrderBulkResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *NewOrderBulkResponse) GetOrderIDs() []string {
	if m != nil {
		return m.OrderIDs
	}
	return nil
}

func (m *NewOrderBulkResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type OrderCancelRequest struct {
	RequestID     uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	OrderID       *types.StringValue `protobuf:"bytes,2,opt,name=orderID,proto3" json:"orderID,omitempty"`
	ClientOrderID *types.StringValue `protobuf:"bytes,3,opt,name=client_orderID,json=clientOrderID,proto3" json:"client_orderID,omitempty"`
	Instrument    *models.Instrument `protobuf:"bytes,4,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Account       *models.Account    `protobuf:"bytes,5,opt,name=account,proto3" json:"account,omitempty"`
}

func (m *OrderCancelRequest) Reset()      { *m = OrderCancelRequest{} }
func (*OrderCancelRequest) ProtoMessage() {}
func (*OrderCancelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{20}
}
func (m *OrderCancelRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderCancelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderCancelRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderCancelRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderCancelRequest.Merge(m, src)
}
func (m *OrderCancelRequest) XXX_Size() int {
	return m.Size()
}
func (m *OrderCancelRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderCancelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OrderCancelRequest proto.InternalMessageInfo

func (m *OrderCancelRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *OrderCancelRequest) GetOrderID() *types.StringValue {
	if m != nil {
		return m.OrderID
	}
	return nil
}

func (m *OrderCancelRequest) GetClientOrderID() *types.StringValue {
	if m != nil {
		return m.ClientOrderID
	}
	return nil
}

func (m *OrderCancelRequest) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *OrderCancelRequest) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

type OrderCancelResponse struct {
	RequestID       uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Success         bool            `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,3,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *OrderCancelResponse) Reset()      { *m = OrderCancelResponse{} }
func (*OrderCancelResponse) ProtoMessage() {}
func (*OrderCancelResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{21}
}
func (m *OrderCancelResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderCancelResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderCancelResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderCancelResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderCancelResponse.Merge(m, src)
}
func (m *OrderCancelResponse) XXX_Size() int {
	return m.Size()
}
func (m *OrderCancelResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderCancelResponse.DiscardUnknown(m)
}

var xxx_messageInfo_OrderCancelResponse proto.InternalMessageInfo

func (m *OrderCancelResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *OrderCancelResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *OrderCancelResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type OrderMassCancelRequest struct {
	RequestID  uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Instrument *models.Instrument `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Account    *models.Account    `protobuf:"bytes,3,opt,name=account,proto3" json:"account,omitempty"`
}

func (m *OrderMassCancelRequest) Reset()      { *m = OrderMassCancelRequest{} }
func (*OrderMassCancelRequest) ProtoMessage() {}
func (*OrderMassCancelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{22}
}
func (m *OrderMassCancelRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderMassCancelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderMassCancelRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderMassCancelRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderMassCancelRequest.Merge(m, src)
}
func (m *OrderMassCancelRequest) XXX_Size() int {
	return m.Size()
}
func (m *OrderMassCancelRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderMassCancelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OrderMassCancelRequest proto.InternalMessageInfo

func (m *OrderMassCancelRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *OrderMassCancelRequest) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *OrderMassCancelRequest) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

type OrderMassCancelResponse struct {
	RequestID       uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Success         bool            `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,3,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *OrderMassCancelResponse) Reset()      { *m = OrderMassCancelResponse{} }
func (*OrderMassCancelResponse) ProtoMessage() {}
func (*OrderMassCancelResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{23}
}
func (m *OrderMassCancelResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderMassCancelResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderMassCancelResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderMassCancelResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderMassCancelResponse.Merge(m, src)
}
func (m *OrderMassCancelResponse) XXX_Size() int {
	return m.Size()
}
func (m *OrderMassCancelResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderMassCancelResponse.DiscardUnknown(m)
}

var xxx_messageInfo_OrderMassCancelResponse proto.InternalMessageInfo

func (m *OrderMassCancelResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *OrderMassCancelResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *OrderMassCancelResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

func init() {
	proto.RegisterEnum("messages.ExecutionType", ExecutionType_name, ExecutionType_value)
	proto.RegisterEnum("messages.RejectionReason", RejectionReason_name, RejectionReason_value)
	proto.RegisterType((*BusinessRequestReject)(nil), "messages.BusinessRequestReject")
	proto.RegisterType((*MarketDataRequest)(nil), "messages.MarketDataRequest")
	proto.RegisterType((*MarketDataRequestReject)(nil), "messages.MarketDataRequestReject")
	proto.RegisterType((*MarketDataSnapshot)(nil), "messages.MarketDataSnapshot")
	proto.RegisterType((*MarketDataIncrementalRefresh)(nil), "messages.MarketDataIncrementalRefresh")
	proto.RegisterType((*SecurityListRequest)(nil), "messages.SecurityListRequest")
	proto.RegisterType((*SecurityList)(nil), "messages.SecurityList")
	proto.RegisterType((*ExecutionReport)(nil), "messages.ExecutionReport")
	proto.RegisterType((*SideValue)(nil), "messages.SideValue")
	proto.RegisterType((*OrderStatusValue)(nil), "messages.OrderStatusValue")
	proto.RegisterType((*OrderFilter)(nil), "messages.OrderFilter")
	proto.RegisterType((*OrderStatusRequest)(nil), "messages.OrderStatusRequest")
	proto.RegisterType((*OrderList)(nil), "messages.OrderList")
	proto.RegisterType((*PositionsRequest)(nil), "messages.PositionsRequest")
	proto.RegisterType((*PositionList)(nil), "messages.PositionList")
	proto.RegisterType((*NewOrder)(nil), "messages.NewOrder")
	proto.RegisterType((*NewOrderSingleRequest)(nil), "messages.NewOrderSingleRequest")
	proto.RegisterType((*NewOrderSingleResponse)(nil), "messages.NewOrderSingleResponse")
	proto.RegisterType((*NewOrderBulkRequest)(nil), "messages.NewOrderBulkRequest")
	proto.RegisterType((*NewOrderBulkResponse)(nil), "messages.NewOrderBulkResponse")
	proto.RegisterType((*OrderCancelRequest)(nil), "messages.OrderCancelRequest")
	proto.RegisterType((*OrderCancelResponse)(nil), "messages.OrderCancelResponse")
	proto.RegisterType((*OrderMassCancelRequest)(nil), "messages.OrderMassCancelRequest")
	proto.RegisterType((*OrderMassCancelResponse)(nil), "messages.OrderMassCancelResponse")
}

func init() { proto.RegisterFile("executor_messages.proto", fileDescriptor_350c53ba9303a7e6) }

var fileDescriptor_350c53ba9303a7e6 = []byte{
	// 1760 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x58, 0x4d, 0x6f, 0xdb, 0xc8,
	0x19, 0x36, 0xf5, 0xcd, 0x57, 0x1f, 0x66, 0xc6, 0x4e, 0xac, 0x1a, 0x86, 0xea, 0x0a, 0x58, 0xd4,
	0x9b, 0x62, 0xe5, 0xad, 0xbd, 0xcd, 0x1e, 0x8a, 0x5d, 0xc0, 0x8e, 0xbc, 0x80, 0x80, 0x7c, 0xb8,
	0x94, 0xd3, 0xab, 0x30, 0x26, 0x5f, 0xcb, 0xac, 0x29, 0x92, 0x99, 0x19, 0x26, 0xf1, 0xa1, 0x40,
	0xcf, 0x7b, 0x28, 0x0a, 0xb4, 0x45, 0xff, 0x41, 0xd1, 0x63, 0x7f, 0x41, 0xcf, 0xbd, 0x14, 0xc8,
	0xa9, 0xd8, 0x02, 0x3d, 0x34, 0xce, 0x25, 0x97, 0xa2, 0xfb, 0x13, 0x0a, 0xce, 0x0c, 0x29, 0x4a,
	0x71, 0x36, 0x4a, 0x6c, 0x04, 0xed, 0x8d, 0xf3, 0xce, 0xf3, 0xcc, 0xcc, 0xfb, 0xcc, 0xfb, 0x31,
	0x12, 0xac, 0xe1, 0x33, 0x74, 0x62, 0x11, 0xb2, 0xd1, 0x04, 0x39, 0xa7, 0x63, 0xe4, 0xbd, 0x88,
	0x85, 0x22, 0x24, 0xb5, 0x74, 0xbc, 0xfe, 0xfd, 0x71, 0x18, 0x8e, 0x7d, 0xdc, 0x96, 0xf6, 0xe3,
	0xf8, 0x64, 0x5b, 0x78, 0x13, 0xe4, 0x82, 0x4e, 0x22, 0x05, 0x5d, 0xef, 0xcc, 0x03, 0x9e, 0x32,
	0x1a, 0x45, 0xc8, 0xf4, 0x52, 0xeb, 0x3f, 0x19, 0x7b, 0xc2, 0xa7, 0xc7, 0x3d, 0x27, 0x9c, 0x6c,
	0x53, 0x3f, 0x3a, 0xa5, 0xc2, 0x73, 0xce, 0xb8, 0xfa, 0x74, 0xb6, 0x27, 0xa1, 0x8b, 0x3e, 0xdf,
	0x9e, 0x50, 0x76, 0x86, 0x62, 0xe4, 0x52, 0x41, 0x35, 0xed, 0xf3, 0x45, 0x68, 0x1c, 0x9d, 0x98,
	0x79, 0xe2, 0x3c, 0x4f, 0xfc, 0xf1, 0x22, 0x44, 0xea, 0x38, 0x61, 0x1c, 0x08, 0x4d, 0xb9, 0x33,
	0xf6, 0xc4, 0x69, 0xac, 0x28, 0x7b, 0xfc, 0x3c, 0x38, 0x63, 0x61, 0x30, 0x38, 0x52, 0x1e, 0x51,
	0x47, 0x84, 0xec, 0x93, 0x71, 0xb8, 0x2d, 0x3f, 0x94, 0x4d, 0xbb, 0xd6, 0xbd, 0x0f, 0x37, 0xf7,
	0x63, 0xee, 0x05, 0xc8, 0xb9, 0x8d, 0x8f, 0x63, 0xe4, 0xc2, 0xc6, 0x5f, 0xa0, 0x23, 0xc8, 0x06,
	0x98, 0x4c, 0x19, 0x06, 0xfd, 0xb6, 0xb1, 0x69, 0x6c, 0x95, 0xec, 0xa9, 0x81, 0xdc, 0x82, 0x0a,
	0x43, 0xca, 0xc3, 0xa0, 0x5d, 0xd8, 0x34, 0xb6, 0x4c, 0x5b, 0x8f, 0xba, 0xff, 0x31, 0xe0, 0xc6,
	0x7d, 0x29, 0x44, 0x9f, 0x0a, 0xaa, 0x57, 0x7c, 0xcb, 0x5a, 0x1b, 0x60, 0xf2, 0xf8, 0x98, 0x3b,
	0xcc, 0x3b, 0x46, 0xb9, 0x5c, 0xcd, 0x9e, 0x1a, 0xc8, 0x6d, 0x80, 0x6c, 0xc0, 0xda, 0xc5, 0x4d,
	0x63, 0xab, 0xbe, 0x03, 0x3d, 0xe9, 0x49, 0xef, 0x70, 0xd0, 0xb7, 0x73, 0xb3, 0x64, 0x07, 0xc0,
	0x0b, 0xb8, 0x60, 0xf1, 0x04, 0x03, 0xd1, 0x2e, 0x49, 0x2c, 0xe9, 0x29, 0xbd, 0x7a, 0x83, 0x6c,
	0xc6, 0xce, 0xa1, 0xc8, 0x97, 0x50, 0xa7, 0xe3, 0x31, 0xc3, 0x31, 0x15, 0x5e, 0x18, 0xb4, 0xcb,
	0x9b, 0xc6, 0x56, 0x6b, 0x67, 0x23, 0x25, 0x3d, 0x64, 0x2e, 0xb2, 0xfd, 0x30, 0x3c, 0xdb, 0x9b,
	0x62, 0xec, 0x3c, 0xa1, 0xfb, 0x10, 0xd6, 0x5e, 0x73, 0xf8, 0xaa, 0x12, 0x92, 0xe9, 0x8a, 0xc3,
	0x80, 0x46, 0xfc, 0x34, 0x7c, 0xdb, 0x62, 0x1d, 0x00, 0x86, 0x3c, 0x0a, 0x03, 0x8e, 0x83, 0xbe,
	0x5c, 0xb0, 0x64, 0xe7, 0x2c, 0xe4, 0x33, 0x00, 0xae, 0x57, 0xba, 0xb7, 0xa3, 0x55, 0x5c, 0xcd,
	0x9c, 0xdc, 0xbf, 0xb7, 0x93, 0xee, 0x63, 0xe7, 0x70, 0x33, 0xac, 0x5d, 0xad, 0x67, 0x9e, 0xb5,
	0x7b, 0x09, 0x6b, 0x97, 0x6c, 0x43, 0x45, 0x30, 0xea, 0x22, 0x6f, 0x97, 0x37, 0x8b, 0x5b, 0xf5,
	0x9d, 0xb5, 0x94, 0x91, 0x6a, 0x88, 0xee, 0x51, 0x32, 0x6f, 0x6b, 0x58, 0xf7, 0xdf, 0x06, 0x6c,
	0x4c, 0x3d, 0x1e, 0x04, 0x0e, 0xc3, 0xe4, 0x6a, 0xa8, 0x6f, 0xe3, 0x09, 0x43, 0x7e, 0x7a, 0x45,
	0xdf, 0x7b, 0x50, 0x8b, 0x23, 0x97, 0x0a, 0xcc, 0x3c, 0x27, 0x79, 0xcf, 0x1f, 0xc9, 0x39, 0x3b,
	0xc3, 0xe4, 0xf0, 0xbb, 0xf3, 0x31, 0x94, 0xf8, 0x3c, 0x87, 0x7f, 0x0f, 0x7f, 0x7f, 0x09, 0x2b,
	0x43, 0x9d, 0xf5, 0xf7, 0xbc, 0x24, 0x5a, 0x3e, 0x68, 0x96, 0x74, 0x7f, 0x67, 0x40, 0x23, 0xbf,
	0xff, 0x15, 0xe5, 0x5d, 0x85, 0x32, 0x32, 0x16, 0xaa, 0x5d, 0x4d, 0x5b, 0x0d, 0xc8, 0xa7, 0x00,
	0xba, 0xb2, 0x79, 0xc8, 0xdb, 0x25, 0x29, 0x8c, 0x95, 0x0a, 0x93, 0xee, 0x6e, 0xe7, 0x30, 0xdd,
	0x57, 0x65, 0x58, 0x3e, 0x90, 0xb5, 0x3c, 0xc9, 0x31, 0x8c, 0x42, 0x26, 0xc8, 0x1a, 0x54, 0x39,
	0x3e, 0x1e, 0x05, 0xf1, 0x44, 0x9f, 0xab, 0xc2, 0xf1, 0xf1, 0x83, 0x78, 0x42, 0xda, 0x50, 0x0d,
	0x93, 0xd4, 0xd4, 0x27, 0x32, 0xed, 0x74, 0x48, 0xee, 0x42, 0xcb, 0xf1, 0x3d, 0x0c, 0xc4, 0x28,
	0x05, 0x28, 0x35, 0x36, 0x7a, 0xaa, 0xc8, 0xf7, 0xd2, 0x22, 0xdf, 0x1b, 0x0a, 0xe6, 0x05, 0xe3,
	0x9f, 0x53, 0x3f, 0x46, 0xbb, 0xa9, 0x38, 0x0f, 0xf5, 0x22, 0x9b, 0x50, 0xc7, 0xf4, 0x28, 0x83,
	0xbe, 0x8c, 0x02, 0xd3, 0xce, 0x9b, 0xc8, 0x97, 0xd0, 0xca, 0x86, 0x23, 0x71, 0x1e, 0xa1, 0xae,
	0x1c, 0x6b, 0xbd, 0xac, 0x0d, 0x65, 0xce, 0x1c, 0x9d, 0x47, 0x68, 0x37, 0x31, 0x3f, 0x24, 0x77,
	0xa0, 0x21, 0xcf, 0x37, 0xe2, 0x82, 0x8a, 0x98, 0xb7, 0x2b, 0x92, 0xbd, 0x32, 0x53, 0x77, 0x86,
	0x72, 0xca, 0xae, 0x87, 0xd3, 0xc1, 0x5c, 0x89, 0xab, 0x2e, 0x54, 0xe2, 0x7e, 0x08, 0xcb, 0x3e,
	0xd2, 0x27, 0xc8, 0x47, 0x8f, 0x63, 0x1a, 0x08, 0x4f, 0x9c, 0xb7, 0xcd, 0x4d, 0x63, 0xcb, 0xb0,
	0x5b, 0xca, 0xfc, 0x33, 0x6d, 0x25, 0x3f, 0x80, 0x86, 0x13, 0x4f, 0xa6, 0x28, 0x90, 0xa8, 0xba,
	0x13, 0x4f, 0x32, 0xc8, 0x01, 0x58, 0x82, 0xd1, 0x80, 0x53, 0x47, 0x79, 0xee, 0x4d, 0xb0, 0x5d,
	0x97, 0xa7, 0x58, 0x7f, 0x4d, 0xe0, 0xa3, 0xb4, 0xcd, 0xda, 0xcb, 0x39, 0x4e, 0x62, 0x25, 0x77,
	0xa0, 0x2a, 0x93, 0x61, 0xd0, 0x6f, 0x37, 0x16, 0xb8, 0x9e, 0x14, 0x4c, 0x7e, 0x0a, 0x70, 0xe2,
	0xf9, 0xfe, 0x28, 0x62, 0x9e, 0x83, 0xed, 0xe6, 0x1b, 0xa8, 0xfd, 0x30, 0x3e, 0xf6, 0x51, 0x51,
	0xcd, 0x04, 0x7f, 0x98, 0xc0, 0xc9, 0x1e, 0x34, 0x25, 0x39, 0xf3, 0xaf, 0xb5, 0x00, 0xbf, 0x91,
	0x50, 0x32, 0xf7, 0xfb, 0x60, 0x31, 0x59, 0xdc, 0x13, 0xe7, 0x75, 0xf9, 0x5e, 0x96, 0x57, 0xf7,
	0xbd, 0xe9, 0xc5, 0xdb, 0x29, 0xc2, 0x96, 0x00, 0x7b, 0x99, 0xcd, 0x1a, 0xba, 0xdb, 0x60, 0x0e,
	0x3d, 0x57, 0x6d, 0x40, 0xba, 0x50, 0x7e, 0x92, 0x7c, 0xc8, 0x08, 0x6f, 0xed, 0x34, 0xb2, 0x24,
	0xf1, 0x5c, 0xb4, 0xd5, 0x54, 0xf7, 0x0b, 0xb0, 0x72, 0x11, 0xa1, 0x78, 0x1f, 0xcf, 0xf2, 0x2e,
	0x0d, 0x1d, 0x4d, 0xff, 0x63, 0x01, 0xea, 0xd2, 0xfc, 0x95, 0xe7, 0x0b, 0x64, 0x89, 0xfa, 0x69,
	0x72, 0x18, 0x8b, 0xa8, 0xff, 0xe6, 0xdc, 0x2a, 0xbc, 0x7b, 0x6e, 0xcd, 0x46, 0x70, 0x71, 0xc1,
	0x08, 0x2e, 0x71, 0xcf, 0x45, 0x5d, 0x8e, 0x57, 0xa6, 0x52, 0x67, 0x32, 0xda, 0x12, 0x40, 0xbe,
	0x98, 0x4b, 0xab, 0xb2, 0x0e, 0xcd, 0x8c, 0x30, 0x2f, 0xe3, 0x4c, 0x76, 0x75, 0xff, 0x6e, 0x00,
	0xc9, 0xeb, 0xf7, 0x81, 0xdf, 0x2f, 0x1f, 0x43, 0x55, 0xbf, 0xea, 0xb4, 0xa7, 0xcb, 0x59, 0x2b,
	0x51, 0x66, 0x3b, 0x9d, 0x27, 0x9f, 0x40, 0xe5, 0x44, 0x5e, 0xa6, 0x76, 0xf1, 0xe6, 0x9c, 0x8b,
	0xea, 0xa6, 0x6d, 0x0d, 0xea, 0xfe, 0xd9, 0x00, 0x53, 0xda, 0x17, 0x28, 0xf8, 0x6d, 0xa8, 0xf2,
	0xd8, 0x71, 0x90, 0x73, 0xed, 0x4d, 0x3a, 0x24, 0x1f, 0x41, 0x45, 0xaa, 0xc5, 0xdb, 0x45, 0x59,
	0xd0, 0x9b, 0x33, 0x31, 0x67, 0xeb, 0xc9, 0x4b, 0x93, 0xa4, 0xf4, 0xce, 0x49, 0xf2, 0x0f, 0x03,
	0xac, 0xc3, 0x90, 0x7b, 0x89, 0x89, 0xff, 0x3f, 0xbc, 0x24, 0x73, 0xb7, 0x57, 0xfe, 0xee, 0xdb,
	0xeb, 0xfe, 0xd3, 0x80, 0x46, 0xea, 0xdb, 0x02, 0x37, 0x72, 0x7d, 0x7e, 0xcd, 0x36, 0xf3, 0xd2,
	0x9b, 0x9b, 0x79, 0x39, 0xdf, 0xcc, 0x7b, 0x60, 0x46, 0xe9, 0x4d, 0xb4, 0x2b, 0xb3, 0xbd, 0x3c,
	0x75, 0xc3, 0x9e, 0x42, 0xba, 0x7f, 0x2b, 0x40, 0xed, 0x01, 0x3e, 0x95, 0x51, 0x41, 0x3e, 0x7a,
	0xad, 0x68, 0x18, 0x72, 0xed, 0xef, 0x2c, 0x0b, 0x85, 0x85, 0x14, 0xff, 0x14, 0x40, 0x65, 0x7b,
	0xae, 0x01, 0xdf, 0x98, 0x89, 0x49, 0xd9, 0x7a, 0xcd, 0x30, 0xfd, 0x24, 0x3f, 0x4a, 0x19, 0xb2,
	0x9c, 0x54, 0x2e, 0xa9, 0xb8, 0x0a, 0x9c, 0x7c, 0x92, 0xcf, 0xa1, 0x99, 0xf4, 0xb7, 0x91, 0x17,
	0x8c, 0x4e, 0x42, 0xe6, 0xa0, 0x6c, 0xb7, 0xb9, 0x4a, 0x9b, 0x74, 0xb2, 0x41, 0xf0, 0x55, 0x32,
	0x65, 0xd7, 0xc5, 0x74, 0x40, 0xd6, 0xa1, 0x96, 0xf5, 0x98, 0x9a, 0xec, 0xa1, 0xd9, 0x98, 0xec,
	0x40, 0x59, 0x35, 0x2f, 0x73, 0x81, 0xe6, 0xa3, 0xa0, 0xdd, 0xaf, 0x0d, 0xb8, 0x99, 0xea, 0x39,
	0xf4, 0x82, 0xb1, 0x8f, 0x8b, 0xe5, 0x43, 0x2e, 0x22, 0x0b, 0x6f, 0xa9, 0x27, 0x5b, 0x50, 0x96,
	0x8e, 0x4f, 0x0b, 0x72, 0x9a, 0xa8, 0xe9, 0xc6, 0xb6, 0x02, 0x24, 0xa5, 0xe4, 0xd6, 0xfc, 0x61,
	0x54, 0xfc, 0xbc, 0x77, 0x5d, 0xc9, 0xbd, 0xe6, 0x8a, 0xb3, 0xaf, 0xb9, 0xeb, 0x29, 0x25, 0x5f,
	0x1b, 0xb0, 0x92, 0x1e, 0x79, 0x3f, 0xf6, 0xcf, 0xae, 0x5d, 0xbd, 0xdb, 0x73, 0x85, 0xf1, 0x32,
	0xf9, 0x34, 0x22, 0xd1, 0x6f, 0x75, 0xf6, 0x30, 0x57, 0x54, 0x6f, 0x1d, 0x6a, 0x5a, 0x2e, 0xb5,
	0xbd, 0x69, 0x67, 0xe3, 0x6b, 0xd2, 0xef, 0xd7, 0x05, 0xdd, 0x16, 0xef, 0xd2, 0xc0, 0x41, 0x7f,
	0x31, 0xf9, 0xee, 0xcc, 0x3e, 0xd1, 0xaf, 0xf0, 0xc8, 0x28, 0x5e, 0xf5, 0x91, 0x71, 0xed, 0xf5,
	0xfb, 0xb7, 0x06, 0xac, 0xcc, 0x08, 0x72, 0xc5, 0x2b, 0xbc, 0xec, 0x9a, 0x8a, 0xef, 0x7c, 0x4d,
	0xbf, 0x37, 0xe0, 0x96, 0x3c, 0xd5, 0x7d, 0xca, 0xf9, 0xbb, 0x5c, 0xd5, 0xfb, 0xd4, 0xde, 0x9c,
	0x5a, 0xc5, 0xb7, 0xa8, 0xf5, 0x07, 0x03, 0xd6, 0x5e, 0x3b, 0xd7, 0xff, 0x82, 0x62, 0xb7, 0x5f,
	0x19, 0xd0, 0x9c, 0xf9, 0x99, 0x46, 0xaa, 0x50, 0x7c, 0x80, 0x4f, 0xad, 0x25, 0x52, 0x83, 0x52,
	0x3f, 0x0c, 0xd0, 0x32, 0x48, 0x03, 0x6a, 0xea, 0xd0, 0xe8, 0x5a, 0x85, 0x64, 0x64, 0x63, 0xe4,
	0x53, 0x07, 0x5d, 0xab, 0x48, 0x6e, 0x40, 0xf3, 0x10, 0x03, 0xd7, 0x0b, 0xc6, 0x0a, 0x62, 0x95,
	0x48, 0x1d, 0xaa, 0x43, 0x11, 0x46, 0x11, 0xba, 0x56, 0x59, 0xa1, 0x93, 0x3d, 0xd1, 0xb5, 0x2a,
	0xa4, 0x09, 0xe6, 0x30, 0xe6, 0x11, 0x06, 0x2e, 0xba, 0x56, 0x95, 0xb4, 0x00, 0x34, 0x39, 0xd9,
	0xb2, 0x96, 0x8c, 0xef, 0x52, 0xdf, 0x89, 0x7d, 0x9a, 0xc0, 0xcd, 0x64, 0xa5, 0x83, 0x67, 0x91,
	0xc7, 0xd0, 0xb5, 0x80, 0x10, 0x68, 0x69, 0xb0, 0xde, 0xde, 0xaa, 0x13, 0x13, 0xca, 0xf2, 0x9f,
	0x05, 0xab, 0x41, 0x96, 0xf5, 0x0b, 0x5f, 0x3d, 0x5c, 0xad, 0x66, 0xb2, 0xd8, 0x10, 0x85, 0xf0,
	0xe5, 0x5f, 0x29, 0x56, 0xeb, 0xf6, 0x5f, 0x0a, 0xb0, 0x3c, 0xa7, 0x47, 0xc2, 0x7f, 0x28, 0x4e,
	0x91, 0x59, 0x4b, 0x89, 0x23, 0x8f, 0x82, 0xb3, 0x20, 0x7c, 0x1a, 0x0c, 0xcf, 0x27, 0xc7, 0xa1,
	0x6f, 0x19, 0xe4, 0x26, 0xdc, 0x48, 0x4d, 0xfa, 0xf7, 0xfa, 0xa0, 0x6f, 0x15, 0x48, 0x17, 0x3a,
	0x8f, 0x02, 0x1e, 0x47, 0xc9, 0x0f, 0x74, 0x74, 0x55, 0x16, 0x9c, 0x52, 0x46, 0x1d, 0x81, 0xcc,
	0xe3, 0xc2, 0x73, 0xac, 0x62, 0x42, 0x1d, 0x04, 0x4e, 0xc8, 0x18, 0x3a, 0x22, 0xfd, 0xed, 0x64,
	0x95, 0x12, 0x1f, 0x0e, 0x9e, 0x39, 0xa7, 0x34, 0x18, 0xe3, 0x5d, 0x3f, 0xe4, 0x52, 0x21, 0x02,
	0xad, 0x7e, 0x1c, 0xf9, 0x9e, 0x43, 0x05, 0xca, 0xc5, 0xac, 0x4a, 0x62, 0x1b, 0x04, 0x4f, 0xa8,
	0xef, 0xb9, 0x3a, 0x96, 0xac, 0x2a, 0x59, 0x81, 0xe5, 0xa3, 0x30, 0xbc, 0x47, 0x05, 0x1e, 0x85,
	0x5a, 0xeb, 0x1a, 0xb1, 0xa0, 0xa1, 0x8f, 0xa8, 0xa8, 0x26, 0x69, 0xc3, 0xaa, 0x9a, 0xdd, 0xf3,
	0x19, 0x52, 0xf7, 0x5c, 0x6b, 0x66, 0x01, 0x59, 0x05, 0xab, 0xef, 0x9d, 0x9c, 0x20, 0xc3, 0x40,
	0x28, 0x1f, 0xb9, 0x55, 0xcf, 0x6d, 0xa5, 0x53, 0xc5, 0x6a, 0x24, 0xc8, 0xf4, 0x98, 0x7b, 0x87,
	0x83, 0x83, 0xe4, 0x09, 0x64, 0x35, 0xf7, 0x3f, 0x7b, 0xfe, 0xa2, 0xb3, 0xf4, 0xcd, 0x8b, 0xce,
	0xd2, 0xb7, 0x2f, 0x3a, 0xc6, 0xaf, 0x2e, 0x3a, 0xc6, 0x9f, 0x2e, 0x3a, 0xc6, 0x5f, 0x2f, 0x3a,
	0xc6, 0xf3, 0x8b, 0x8e, 0xf1, 0xaf, 0x8b, 0x8e, 0xf1, 0xea, 0xa2, 0xb3, 0xf4, 0xed, 0x45, 0xc7,
	0xf8, 0xcd, 0xcb, 0xce, 0xd2, 0xf3, 0x97, 0x9d, 0xa5, 0x6f, 0x5e, 0x76, 0x96, 0x8e, 0x2b, 0xb2,
	0x5a, 0xed, 0xfe, 0x37, 0x00, 0x00, 0xff, 0xff, 0x81, 0x23, 0xc9, 0xe3, 0xa7, 0x16, 0x00, 0x00,
}

func (x ExecutionType) String() string {
	s, ok := ExecutionType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x RejectionReason) String() string {
	s, ok := RejectionReason_name[int32(x)]
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
	if !this.FillPrice.Equal(that1.FillPrice) {
		return false
	}
	if !this.FillQuantity.Equal(that1.FillQuantity) {
		return false
	}
	if this.RejectionReason != that1.RejectionReason {
		return false
	}
	return true
}
func (this *SideValue) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SideValue)
	if !ok {
		that2, ok := that.(SideValue)
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
	if this.Value != that1.Value {
		return false
	}
	return true
}
func (this *OrderStatusValue) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OrderStatusValue)
	if !ok {
		that2, ok := that.(OrderStatusValue)
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
	if this.Value != that1.Value {
		return false
	}
	return true
}
func (this *OrderFilter) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OrderFilter)
	if !ok {
		that2, ok := that.(OrderFilter)
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
	if !this.OrderID.Equal(that1.OrderID) {
		return false
	}
	if !this.ClientOrderID.Equal(that1.ClientOrderID) {
		return false
	}
	if !this.Instrument.Equal(that1.Instrument) {
		return false
	}
	if !this.Side.Equal(that1.Side) {
		return false
	}
	if !this.OrderStatus.Equal(that1.OrderStatus) {
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
	if !this.Account.Equal(that1.Account) {
		return false
	}
	if !this.Filter.Equal(that1.Filter) {
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
	if this.Success != that1.Success {
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
	if this.RejectionReason != that1.RejectionReason {
		return false
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
func (this *NewOrder) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*NewOrder)
	if !ok {
		that2, ok := that.(NewOrder)
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
func (this *NewOrderSingleRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*NewOrderSingleRequest)
	if !ok {
		that2, ok := that.(NewOrderSingleRequest)
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
	if !this.Account.Equal(that1.Account) {
		return false
	}
	if !this.Order.Equal(that1.Order) {
		return false
	}
	return true
}
func (this *NewOrderSingleResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*NewOrderSingleResponse)
	if !ok {
		that2, ok := that.(NewOrderSingleResponse)
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
	if this.Success != that1.Success {
		return false
	}
	if this.OrderID != that1.OrderID {
		return false
	}
	if this.RejectionReason != that1.RejectionReason {
		return false
	}
	return true
}
func (this *NewOrderBulkRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*NewOrderBulkRequest)
	if !ok {
		that2, ok := that.(NewOrderBulkRequest)
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
	if !this.Account.Equal(that1.Account) {
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
func (this *NewOrderBulkResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*NewOrderBulkResponse)
	if !ok {
		that2, ok := that.(NewOrderBulkResponse)
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
	if this.Success != that1.Success {
		return false
	}
	if len(this.OrderIDs) != len(that1.OrderIDs) {
		return false
	}
	for i := range this.OrderIDs {
		if this.OrderIDs[i] != that1.OrderIDs[i] {
			return false
		}
	}
	if this.RejectionReason != that1.RejectionReason {
		return false
	}
	return true
}
func (this *OrderCancelRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OrderCancelRequest)
	if !ok {
		that2, ok := that.(OrderCancelRequest)
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
func (this *OrderCancelResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OrderCancelResponse)
	if !ok {
		that2, ok := that.(OrderCancelResponse)
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
	if this.Success != that1.Success {
		return false
	}
	if this.RejectionReason != that1.RejectionReason {
		return false
	}
	return true
}
func (this *OrderMassCancelRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OrderMassCancelRequest)
	if !ok {
		that2, ok := that.(OrderMassCancelRequest)
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
	if !this.Instrument.Equal(that1.Instrument) {
		return false
	}
	if !this.Account.Equal(that1.Account) {
		return false
	}
	return true
}
func (this *OrderMassCancelResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OrderMassCancelResponse)
	if !ok {
		that2, ok := that.(OrderMassCancelResponse)
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
	if this.Success != that1.Success {
		return false
	}
	if this.RejectionReason != that1.RejectionReason {
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
	s := make([]string, 0, 18)
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
	s = append(s, "LeavesQuantity: "+fmt.Sprintf("%#v", this.LeavesQuantity)+",\n")
	s = append(s, "CumQuantity: "+fmt.Sprintf("%#v", this.CumQuantity)+",\n")
	if this.TransactionTime != nil {
		s = append(s, "TransactionTime: "+fmt.Sprintf("%#v", this.TransactionTime)+",\n")
	}
	if this.TradeID != nil {
		s = append(s, "TradeID: "+fmt.Sprintf("%#v", this.TradeID)+",\n")
	}
	if this.FillPrice != nil {
		s = append(s, "FillPrice: "+fmt.Sprintf("%#v", this.FillPrice)+",\n")
	}
	if this.FillQuantity != nil {
		s = append(s, "FillQuantity: "+fmt.Sprintf("%#v", this.FillQuantity)+",\n")
	}
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *SideValue) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&messages.SideValue{")
	s = append(s, "Value: "+fmt.Sprintf("%#v", this.Value)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OrderStatusValue) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&messages.OrderStatusValue{")
	s = append(s, "Value: "+fmt.Sprintf("%#v", this.Value)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OrderFilter) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.OrderFilter{")
	if this.OrderID != nil {
		s = append(s, "OrderID: "+fmt.Sprintf("%#v", this.OrderID)+",\n")
	}
	if this.ClientOrderID != nil {
		s = append(s, "ClientOrderID: "+fmt.Sprintf("%#v", this.ClientOrderID)+",\n")
	}
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	if this.Side != nil {
		s = append(s, "Side: "+fmt.Sprintf("%#v", this.Side)+",\n")
	}
	if this.OrderStatus != nil {
		s = append(s, "OrderStatus: "+fmt.Sprintf("%#v", this.OrderStatus)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OrderStatusRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.OrderStatusRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Subscribe: "+fmt.Sprintf("%#v", this.Subscribe)+",\n")
	if this.Subscriber != nil {
		s = append(s, "Subscriber: "+fmt.Sprintf("%#v", this.Subscriber)+",\n")
	}
	if this.Account != nil {
		s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
	}
	if this.Filter != nil {
		s = append(s, "Filter: "+fmt.Sprintf("%#v", this.Filter)+",\n")
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
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	if this.Orders != nil {
		s = append(s, "Orders: "+fmt.Sprintf("%#v", this.Orders)+",\n")
	}
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
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
func (this *NewOrder) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 11)
	s = append(s, "&messages.NewOrder{")
	s = append(s, "ClientOrderID: "+fmt.Sprintf("%#v", this.ClientOrderID)+",\n")
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
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
func (this *NewOrderSingleRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.NewOrderSingleRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	if this.Account != nil {
		s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
	}
	if this.Order != nil {
		s = append(s, "Order: "+fmt.Sprintf("%#v", this.Order)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *NewOrderSingleResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.NewOrderSingleResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "OrderID: "+fmt.Sprintf("%#v", this.OrderID)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *NewOrderBulkRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.NewOrderBulkRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	if this.Account != nil {
		s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
	}
	if this.Orders != nil {
		s = append(s, "Orders: "+fmt.Sprintf("%#v", this.Orders)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *NewOrderBulkResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.NewOrderBulkResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "OrderIDs: "+fmt.Sprintf("%#v", this.OrderIDs)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OrderCancelRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.OrderCancelRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
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
func (this *OrderCancelResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.OrderCancelResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OrderMassCancelRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.OrderMassCancelRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	if this.Account != nil {
		s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OrderMassCancelResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.OrderMassCancelResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
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
	if m.FillPrice != nil {
		dAtA[i] = 0x6a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.FillPrice.Size()))
		n12, err := m.FillPrice.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n12
	}
	if m.FillQuantity != nil {
		dAtA[i] = 0x72
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.FillQuantity.Size()))
		n13, err := m.FillQuantity.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n13
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x78
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
	}
	return i, nil
}

func (m *SideValue) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SideValue) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Value != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Value))
	}
	return i, nil
}

func (m *OrderStatusValue) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderStatusValue) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Value != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Value))
	}
	return i, nil
}

func (m *OrderFilter) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderFilter) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.OrderID != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderID.Size()))
		n14, err := m.OrderID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n14
	}
	if m.ClientOrderID != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ClientOrderID.Size()))
		n15, err := m.ClientOrderID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n15
	}
	if m.Instrument != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Instrument.Size()))
		n16, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n16
	}
	if m.Side != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Side.Size()))
		n17, err := m.Side.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n17
	}
	if m.OrderStatus != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderStatus.Size()))
		n18, err := m.OrderStatus.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n18
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
		n19, err := m.Subscriber.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n19
	}
	if m.Account != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n20, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n20
	}
	if m.Filter != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Filter.Size()))
		n21, err := m.Filter.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n21
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
	if m.Success {
		dAtA[i] = 0x10
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if len(m.Orders) > 0 {
		for _, msg := range m.Orders {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintExecutorMessages(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
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
		n22, err := m.Subscriber.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n22
	}
	if m.Instrument != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Instrument.Size()))
		n23, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n23
	}
	if m.Account != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n24, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n24
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
		n25, err := m.Subscriber.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n25
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

func (m *NewOrder) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NewOrder) MarshalTo(dAtA []byte) (int, error) {
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
		n26, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n26
	}
	if m.OrderType != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderType))
	}
	if m.OrderSide != 0 {
		dAtA[i] = 0x30
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderSide))
	}
	if m.TimeInForce != 0 {
		dAtA[i] = 0x38
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.TimeInForce))
	}
	if m.Quantity != 0 {
		dAtA[i] = 0x41
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Quantity))))
		i += 8
	}
	if m.Price != nil {
		dAtA[i] = 0x4a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Price.Size()))
		n27, err := m.Price.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n27
	}
	return i, nil
}

func (m *NewOrderSingleRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NewOrderSingleRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RequestID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
	}
	if m.Account != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n28, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n28
	}
	if m.Order != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Order.Size()))
		n29, err := m.Order.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n29
	}
	return i, nil
}

func (m *NewOrderSingleResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NewOrderSingleResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RequestID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
	}
	if m.Success {
		dAtA[i] = 0x10
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if len(m.OrderID) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.OrderID)))
		i += copy(dAtA[i:], m.OrderID)
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
	}
	return i, nil
}

func (m *NewOrderBulkRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NewOrderBulkRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RequestID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
	}
	if m.Account != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n30, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n30
	}
	if len(m.Orders) > 0 {
		for _, msg := range m.Orders {
			dAtA[i] = 0x1a
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

func (m *NewOrderBulkResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NewOrderBulkResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RequestID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
	}
	if m.Success {
		dAtA[i] = 0x10
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if len(m.OrderIDs) > 0 {
		for _, s := range m.OrderIDs {
			dAtA[i] = 0x1a
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
	}
	return i, nil
}

func (m *OrderCancelRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderCancelRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RequestID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
	}
	if m.OrderID != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderID.Size()))
		n31, err := m.OrderID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n31
	}
	if m.ClientOrderID != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ClientOrderID.Size()))
		n32, err := m.ClientOrderID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n32
	}
	if m.Instrument != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Instrument.Size()))
		n33, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n33
	}
	if m.Account != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n34, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n34
	}
	return i, nil
}

func (m *OrderCancelResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderCancelResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RequestID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
	}
	if m.Success {
		dAtA[i] = 0x10
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
	}
	return i, nil
}

func (m *OrderMassCancelRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderMassCancelRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RequestID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
	}
	if m.Instrument != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Instrument.Size()))
		n35, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n35
	}
	if m.Account != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n36, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n36
	}
	return i, nil
}

func (m *OrderMassCancelResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderMassCancelResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.RequestID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
	}
	if m.Success {
		dAtA[i] = 0x10
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
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
	if m.FillPrice != nil {
		l = m.FillPrice.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.FillQuantity != nil {
		l = m.FillQuantity.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
	}
	return n
}

func (m *SideValue) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Value != 0 {
		n += 1 + sovExecutorMessages(uint64(m.Value))
	}
	return n
}

func (m *OrderStatusValue) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Value != 0 {
		n += 1 + sovExecutorMessages(uint64(m.Value))
	}
	return n
}

func (m *OrderFilter) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
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
	if m.Side != nil {
		l = m.Side.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.OrderStatus != nil {
		l = m.OrderStatus.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
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
	if m.Account != nil {
		l = m.Account.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Filter != nil {
		l = m.Filter.Size()
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
	if m.Success {
		n += 2
	}
	if len(m.Orders) > 0 {
		for _, e := range m.Orders {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
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

func (m *NewOrder) Size() (n int) {
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

func (m *NewOrderSingleRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
	}
	if m.Account != nil {
		l = m.Account.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Order != nil {
		l = m.Order.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *NewOrderSingleResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
	}
	if m.Success {
		n += 2
	}
	l = len(m.OrderID)
	if l > 0 {
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
	}
	return n
}

func (m *NewOrderBulkRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
	}
	if m.Account != nil {
		l = m.Account.Size()
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

func (m *NewOrderBulkResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
	}
	if m.Success {
		n += 2
	}
	if len(m.OrderIDs) > 0 {
		for _, s := range m.OrderIDs {
			l = len(s)
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
	}
	return n
}

func (m *OrderCancelRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
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

func (m *OrderCancelResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
	}
	if m.Success {
		n += 2
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
	}
	return n
}

func (m *OrderMassCancelRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
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

func (m *OrderMassCancelResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
	}
	if m.Success {
		n += 2
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
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
		`LeavesQuantity:` + fmt.Sprintf("%v", this.LeavesQuantity) + `,`,
		`CumQuantity:` + fmt.Sprintf("%v", this.CumQuantity) + `,`,
		`TransactionTime:` + strings.Replace(fmt.Sprintf("%v", this.TransactionTime), "Timestamp", "types.Timestamp", 1) + `,`,
		`TradeID:` + strings.Replace(fmt.Sprintf("%v", this.TradeID), "StringValue", "types.StringValue", 1) + `,`,
		`FillPrice:` + strings.Replace(fmt.Sprintf("%v", this.FillPrice), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`FillQuantity:` + strings.Replace(fmt.Sprintf("%v", this.FillQuantity), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *SideValue) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&SideValue{`,
		`Value:` + fmt.Sprintf("%v", this.Value) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrderStatusValue) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OrderStatusValue{`,
		`Value:` + fmt.Sprintf("%v", this.Value) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrderFilter) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OrderFilter{`,
		`OrderID:` + strings.Replace(fmt.Sprintf("%v", this.OrderID), "StringValue", "types.StringValue", 1) + `,`,
		`ClientOrderID:` + strings.Replace(fmt.Sprintf("%v", this.ClientOrderID), "StringValue", "types.StringValue", 1) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`Side:` + strings.Replace(fmt.Sprintf("%v", this.Side), "SideValue", "SideValue", 1) + `,`,
		`OrderStatus:` + strings.Replace(fmt.Sprintf("%v", this.OrderStatus), "OrderStatusValue", "OrderStatusValue", 1) + `,`,
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
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`Filter:` + strings.Replace(fmt.Sprintf("%v", this.Filter), "OrderFilter", "OrderFilter", 1) + `,`,
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
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`Orders:` + strings.Replace(fmt.Sprintf("%v", this.Orders), "Order", "models.Order", 1) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
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
func (this *NewOrder) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&NewOrder{`,
		`ClientOrderID:` + fmt.Sprintf("%v", this.ClientOrderID) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`OrderType:` + fmt.Sprintf("%v", this.OrderType) + `,`,
		`OrderSide:` + fmt.Sprintf("%v", this.OrderSide) + `,`,
		`TimeInForce:` + fmt.Sprintf("%v", this.TimeInForce) + `,`,
		`Quantity:` + fmt.Sprintf("%v", this.Quantity) + `,`,
		`Price:` + strings.Replace(fmt.Sprintf("%v", this.Price), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *NewOrderSingleRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&NewOrderSingleRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`Order:` + strings.Replace(fmt.Sprintf("%v", this.Order), "NewOrder", "NewOrder", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *NewOrderSingleResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&NewOrderSingleResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`OrderID:` + fmt.Sprintf("%v", this.OrderID) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *NewOrderBulkRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&NewOrderBulkRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`Orders:` + strings.Replace(fmt.Sprintf("%v", this.Orders), "NewOrder", "NewOrder", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *NewOrderBulkResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&NewOrderBulkResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`OrderIDs:` + fmt.Sprintf("%v", this.OrderIDs) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrderCancelRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OrderCancelRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`OrderID:` + strings.Replace(fmt.Sprintf("%v", this.OrderID), "StringValue", "types.StringValue", 1) + `,`,
		`ClientOrderID:` + strings.Replace(fmt.Sprintf("%v", this.ClientOrderID), "StringValue", "types.StringValue", 1) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrderCancelResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OrderCancelResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrderMassCancelRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OrderMassCancelRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrderMassCancelResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OrderMassCancelResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FillPrice", wireType)
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
			if m.FillPrice == nil {
				m.FillPrice = &types.DoubleValue{}
			}
			if err := m.FillPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FillQuantity", wireType)
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
			if m.FillQuantity == nil {
				m.FillQuantity = &types.DoubleValue{}
			}
			if err := m.FillQuantity.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 15:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RejectionReason", wireType)
			}
			m.RejectionReason = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RejectionReason |= RejectionReason(b&0x7F) << shift
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
func (m *SideValue) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: SideValue: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SideValue: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			m.Value = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Value |= models.Side(b&0x7F) << shift
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
func (m *OrderStatusValue) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OrderStatusValue: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderStatusValue: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			m.Value = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Value |= models.OrderStatus(b&0x7F) << shift
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
func (m *OrderFilter) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OrderFilter: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderFilter: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
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
		case 2:
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
				m.Instrument = &models.Instrument{}
			}
			if err := m.Instrument.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Side", wireType)
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
			if m.Side == nil {
				m.Side = &SideValue{}
			}
			if err := m.Side.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderStatus", wireType)
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
			if m.OrderStatus == nil {
				m.OrderStatus = &OrderStatusValue{}
			}
			if err := m.OrderStatus.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Filter", wireType)
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
			if m.Filter == nil {
				m.Filter = &OrderFilter{}
			}
			if err := m.Filter.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
				return fmt.Errorf("proto: wrong wireType = %d for field Success", wireType)
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
			m.Success = bool(v != 0)
		case 3:
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
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RejectionReason", wireType)
			}
			m.RejectionReason = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RejectionReason |= RejectionReason(b&0x7F) << shift
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
func (m *NewOrder) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: NewOrder: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NewOrder: illegal tag %d (wire type %d)", fieldNum, wire)
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
		case 5:
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
		case 6:
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
		case 7:
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
		case 8:
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
		case 9:
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
func (m *NewOrderSingleRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: NewOrderSingleRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NewOrderSingleRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Order", wireType)
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
			if m.Order == nil {
				m.Order = &NewOrder{}
			}
			if err := m.Order.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *NewOrderSingleResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: NewOrderSingleResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NewOrderSingleResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Success", wireType)
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
			m.Success = bool(v != 0)
		case 3:
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
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RejectionReason", wireType)
			}
			m.RejectionReason = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RejectionReason |= RejectionReason(b&0x7F) << shift
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
func (m *NewOrderBulkRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: NewOrderBulkRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NewOrderBulkRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
		case 3:
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
			m.Orders = append(m.Orders, &NewOrder{})
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
func (m *NewOrderBulkResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: NewOrderBulkResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NewOrderBulkResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Success", wireType)
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
			m.Success = bool(v != 0)
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderIDs", wireType)
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
			m.OrderIDs = append(m.OrderIDs, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RejectionReason", wireType)
			}
			m.RejectionReason = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RejectionReason |= RejectionReason(b&0x7F) << shift
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
func (m *OrderCancelRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OrderCancelRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderCancelRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *OrderCancelResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OrderCancelResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderCancelResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Success", wireType)
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
			m.Success = bool(v != 0)
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RejectionReason", wireType)
			}
			m.RejectionReason = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RejectionReason |= RejectionReason(b&0x7F) << shift
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
func (m *OrderMassCancelRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OrderMassCancelRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderMassCancelRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *OrderMassCancelResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OrderMassCancelResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderMassCancelResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Success", wireType)
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
			m.Success = bool(v != 0)
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RejectionReason", wireType)
			}
			m.RejectionReason = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RejectionReason |= RejectionReason(b&0x7F) << shift
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
