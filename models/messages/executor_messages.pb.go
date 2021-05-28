package messages

import (
	context "context"
	encoding_binary "encoding/binary"
	fmt "fmt"
	actor "github.com/AsynkronIT/protoactor-go/actor"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
	models "gitlab.com/alphaticks/alpha-connect/models"
	models1 "gitlab.com/alphaticks/xchanger/models"
	grpc "google.golang.org/grpc"
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
	InvalidOrder                   RejectionReason = 14
	RateLimitExceeded              RejectionReason = 15
	UnsupportedSubscription        RejectionReason = 16
	MissingInstrument              RejectionReason = 17
	HTTPError                      RejectionReason = 18
	NonReplaceableOrder            RejectionReason = 19
	NonCancelableOrder             RejectionReason = 20
	UnsupportedFilter              RejectionReason = 21
	UnsupportedOrderType           RejectionReason = 22
	UnsupportedOrderTimeInForce    RejectionReason = 23
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
	14: "InvalidOrder",
	15: "RateLimitExceeded",
	16: "UnsupportedSubscription",
	17: "MissingInstrument",
	18: "HTTPError",
	19: "NonReplaceableOrder",
	20: "NonCancelableOrder",
	21: "UnsupportedFilter",
	22: "UnsupportedOrderType",
	23: "UnsupportedOrderTimeInForce",
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
	"InvalidOrder":                   14,
	"RateLimitExceeded":              15,
	"UnsupportedSubscription":        16,
	"MissingInstrument":              17,
	"HTTPError":                      18,
	"NonReplaceableOrder":            19,
	"NonCancelableOrder":             20,
	"UnsupportedFilter":              21,
	"UnsupportedOrderType":           22,
	"UnsupportedOrderTimeInForce":    23,
}

func (RejectionReason) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{1}
}

type FeeType int32

const (
	Regulatory      FeeType = 0
	Tax             FeeType = 1
	LocalCommission FeeType = 2
	ExchangeFees    FeeType = 3
)

var FeeType_name = map[int32]string{
	0: "Regulatory",
	1: "Tax",
	2: "LocalCommission",
	3: "ExchangeFees",
}

var FeeType_value = map[string]int32{
	"Regulatory":      0,
	"Tax":             1,
	"LocalCommission": 2,
	"ExchangeFees":    3,
}

func (FeeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{2}
}

type FeeBasis int32

const (
	Absolute   FeeBasis = 0
	PerUnit    FeeBasis = 1
	Percentage FeeBasis = 2
)

var FeeBasis_name = map[int32]string{
	0: "Absolute",
	1: "PerUnit",
	2: "Percentage",
}

var FeeBasis_value = map[string]int32{
	"Absolute":   0,
	"PerUnit":    1,
	"Percentage": 2,
}

func (FeeBasis) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{3}
}

type ExecutionInstruction int32

const (
	ParticipateDoNotInitiate ExecutionInstruction = 0
)

var ExecutionInstruction_name = map[int32]string{
	0: "ParticipateDoNotInitiate",
}

var ExecutionInstruction_value = map[string]int32{
	"ParticipateDoNotInitiate": 0,
}

func (ExecutionInstruction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{4}
}

type HistoricalLiquidationsRequest struct {
	RequestID  uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Instrument *models.Instrument `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
	From       *types.Timestamp   `protobuf:"bytes,3,opt,name=from,proto3" json:"from,omitempty"`
	To         *types.Timestamp   `protobuf:"bytes,4,opt,name=to,proto3" json:"to,omitempty"`
}

func (m *HistoricalLiquidationsRequest) Reset()      { *m = HistoricalLiquidationsRequest{} }
func (*HistoricalLiquidationsRequest) ProtoMessage() {}
func (*HistoricalLiquidationsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{0}
}
func (m *HistoricalLiquidationsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HistoricalLiquidationsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HistoricalLiquidationsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HistoricalLiquidationsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HistoricalLiquidationsRequest.Merge(m, src)
}
func (m *HistoricalLiquidationsRequest) XXX_Size() int {
	return m.Size()
}
func (m *HistoricalLiquidationsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HistoricalLiquidationsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HistoricalLiquidationsRequest proto.InternalMessageInfo

func (m *HistoricalLiquidationsRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *HistoricalLiquidationsRequest) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *HistoricalLiquidationsRequest) GetFrom() *types.Timestamp {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *HistoricalLiquidationsRequest) GetTo() *types.Timestamp {
	if m != nil {
		return m.To
	}
	return nil
}

type HistoricalLiquidationsResponse struct {
	RequestID       uint64                `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64                `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Liquidations    []*models.Liquidation `protobuf:"bytes,3,rep,name=liquidations,proto3" json:"liquidations,omitempty"`
	Success         bool                  `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason       `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *HistoricalLiquidationsResponse) Reset()      { *m = HistoricalLiquidationsResponse{} }
func (*HistoricalLiquidationsResponse) ProtoMessage() {}
func (*HistoricalLiquidationsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{1}
}
func (m *HistoricalLiquidationsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HistoricalLiquidationsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HistoricalLiquidationsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HistoricalLiquidationsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HistoricalLiquidationsResponse.Merge(m, src)
}
func (m *HistoricalLiquidationsResponse) XXX_Size() int {
	return m.Size()
}
func (m *HistoricalLiquidationsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HistoricalLiquidationsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HistoricalLiquidationsResponse proto.InternalMessageInfo

func (m *HistoricalLiquidationsResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *HistoricalLiquidationsResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *HistoricalLiquidationsResponse) GetLiquidations() []*models.Liquidation {
	if m != nil {
		return m.Liquidations
	}
	return nil
}

func (m *HistoricalLiquidationsResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *HistoricalLiquidationsResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
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
	return fileDescriptor_350c53ba9303a7e6, []int{2}
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
	return models.L1
}

type MarketDataResponse struct {
	RequestID       uint64                    `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64                    `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	SnapshotL1      *models.OBL1Snapshot      `protobuf:"bytes,3,opt,name=snapshotL1,proto3" json:"snapshotL1,omitempty"`
	SnapshotL2      *models.OBL2Snapshot      `protobuf:"bytes,4,opt,name=snapshotL2,proto3" json:"snapshotL2,omitempty"`
	SnapshotL3      *models.OBL3Snapshot      `protobuf:"bytes,5,opt,name=snapshotL3,proto3" json:"snapshotL3,omitempty"`
	Trades          []*models.AggregatedTrade `protobuf:"bytes,6,rep,name=trades,proto3" json:"trades,omitempty"`
	SeqNum          uint64                    `protobuf:"varint,7,opt,name=seq_num,json=seqNum,proto3" json:"seq_num,omitempty"`
	Success         bool                      `protobuf:"varint,8,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason           `protobuf:"varint,9,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *MarketDataResponse) Reset()      { *m = MarketDataResponse{} }
func (*MarketDataResponse) ProtoMessage() {}
func (*MarketDataResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{3}
}
func (m *MarketDataResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MarketDataResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MarketDataResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MarketDataResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MarketDataResponse.Merge(m, src)
}
func (m *MarketDataResponse) XXX_Size() int {
	return m.Size()
}
func (m *MarketDataResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MarketDataResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MarketDataResponse proto.InternalMessageInfo

func (m *MarketDataResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *MarketDataResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *MarketDataResponse) GetSnapshotL1() *models.OBL1Snapshot {
	if m != nil {
		return m.SnapshotL1
	}
	return nil
}

func (m *MarketDataResponse) GetSnapshotL2() *models.OBL2Snapshot {
	if m != nil {
		return m.SnapshotL2
	}
	return nil
}

func (m *MarketDataResponse) GetSnapshotL3() *models.OBL3Snapshot {
	if m != nil {
		return m.SnapshotL3
	}
	return nil
}

func (m *MarketDataResponse) GetTrades() []*models.AggregatedTrade {
	if m != nil {
		return m.Trades
	}
	return nil
}

func (m *MarketDataResponse) GetSeqNum() uint64 {
	if m != nil {
		return m.SeqNum
	}
	return 0
}

func (m *MarketDataResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *MarketDataResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type MarketDataIncrementalRefresh struct {
	RequestID   uint64                    `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID  uint64                    `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	SeqNum      uint64                    `protobuf:"varint,3,opt,name=seq_num,json=seqNum,proto3" json:"seq_num,omitempty"`
	UpdateL1    *models.OBL1Update        `protobuf:"bytes,4,opt,name=updateL1,proto3" json:"updateL1,omitempty"`
	UpdateL2    *models.OBL2Update        `protobuf:"bytes,5,opt,name=updateL2,proto3" json:"updateL2,omitempty"`
	UpdateL3    *models.OBL3Update        `protobuf:"bytes,6,opt,name=updateL3,proto3" json:"updateL3,omitempty"`
	Trades      []*models.AggregatedTrade `protobuf:"bytes,7,rep,name=trades,proto3" json:"trades,omitempty"`
	Liquidation *models.Liquidation       `protobuf:"bytes,8,opt,name=liquidation,proto3" json:"liquidation,omitempty"`
	Funding     *models.Funding           `protobuf:"bytes,9,opt,name=funding,proto3" json:"funding,omitempty"`
	Stats       []*models.Stat            `protobuf:"bytes,10,rep,name=stats,proto3" json:"stats,omitempty"`
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

func (m *MarketDataIncrementalRefresh) GetSeqNum() uint64 {
	if m != nil {
		return m.SeqNum
	}
	return 0
}

func (m *MarketDataIncrementalRefresh) GetUpdateL1() *models.OBL1Update {
	if m != nil {
		return m.UpdateL1
	}
	return nil
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

func (m *MarketDataIncrementalRefresh) GetLiquidation() *models.Liquidation {
	if m != nil {
		return m.Liquidation
	}
	return nil
}

func (m *MarketDataIncrementalRefresh) GetFunding() *models.Funding {
	if m != nil {
		return m.Funding
	}
	return nil
}

func (m *MarketDataIncrementalRefresh) GetStats() []*models.Stat {
	if m != nil {
		return m.Stats
	}
	return nil
}

type AccountDataRequest struct {
	RequestID  uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe  bool            `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
	Subscriber *actor.PID      `protobuf:"bytes,3,opt,name=subscriber,proto3" json:"subscriber,omitempty"`
	Account    *models.Account `protobuf:"bytes,4,opt,name=account,proto3" json:"account,omitempty"`
}

func (m *AccountDataRequest) Reset()      { *m = AccountDataRequest{} }
func (*AccountDataRequest) ProtoMessage() {}
func (*AccountDataRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{5}
}
func (m *AccountDataRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountDataRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccountDataRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AccountDataRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountDataRequest.Merge(m, src)
}
func (m *AccountDataRequest) XXX_Size() int {
	return m.Size()
}
func (m *AccountDataRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountDataRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AccountDataRequest proto.InternalMessageInfo

func (m *AccountDataRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *AccountDataRequest) GetSubscribe() bool {
	if m != nil {
		return m.Subscribe
	}
	return false
}

func (m *AccountDataRequest) GetSubscriber() *actor.PID {
	if m != nil {
		return m.Subscriber
	}
	return nil
}

func (m *AccountDataRequest) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

type AccountDataResponse struct {
	RequestID  uint64        `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID uint64        `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Orders     *OrderList    `protobuf:"bytes,3,opt,name=orders,proto3" json:"orders,omitempty"`
	Positions  *PositionList `protobuf:"bytes,4,opt,name=positions,proto3" json:"positions,omitempty"`
	Balances   *BalanceList  `protobuf:"bytes,5,opt,name=balances,proto3" json:"balances,omitempty"`
	SeqNum     uint64        `protobuf:"varint,6,opt,name=seq_num,json=seqNum,proto3" json:"seq_num,omitempty"`
}

func (m *AccountDataResponse) Reset()      { *m = AccountDataResponse{} }
func (*AccountDataResponse) ProtoMessage() {}
func (*AccountDataResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{6}
}
func (m *AccountDataResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountDataResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccountDataResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AccountDataResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountDataResponse.Merge(m, src)
}
func (m *AccountDataResponse) XXX_Size() int {
	return m.Size()
}
func (m *AccountDataResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountDataResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AccountDataResponse proto.InternalMessageInfo

func (m *AccountDataResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *AccountDataResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *AccountDataResponse) GetOrders() *OrderList {
	if m != nil {
		return m.Orders
	}
	return nil
}

func (m *AccountDataResponse) GetPositions() *PositionList {
	if m != nil {
		return m.Positions
	}
	return nil
}

func (m *AccountDataResponse) GetBalances() *BalanceList {
	if m != nil {
		return m.Balances
	}
	return nil
}

func (m *AccountDataResponse) GetSeqNum() uint64 {
	if m != nil {
		return m.SeqNum
	}
	return 0
}

type AccountDataIncrementalRefresh struct {
	RequestID  uint64           `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID uint64           `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Report     *ExecutionReport `protobuf:"bytes,3,opt,name=report,proto3" json:"report,omitempty"`
}

func (m *AccountDataIncrementalRefresh) Reset()      { *m = AccountDataIncrementalRefresh{} }
func (*AccountDataIncrementalRefresh) ProtoMessage() {}
func (*AccountDataIncrementalRefresh) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{7}
}
func (m *AccountDataIncrementalRefresh) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountDataIncrementalRefresh) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccountDataIncrementalRefresh.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AccountDataIncrementalRefresh) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountDataIncrementalRefresh.Merge(m, src)
}
func (m *AccountDataIncrementalRefresh) XXX_Size() int {
	return m.Size()
}
func (m *AccountDataIncrementalRefresh) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountDataIncrementalRefresh.DiscardUnknown(m)
}

var xxx_messageInfo_AccountDataIncrementalRefresh proto.InternalMessageInfo

func (m *AccountDataIncrementalRefresh) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *AccountDataIncrementalRefresh) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *AccountDataIncrementalRefresh) GetReport() *ExecutionReport {
	if m != nil {
		return m.Report
	}
	return nil
}

type SecurityDefinitionRequest struct {
	RequestID  uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Instrument *models.Instrument `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
}

func (m *SecurityDefinitionRequest) Reset()      { *m = SecurityDefinitionRequest{} }
func (*SecurityDefinitionRequest) ProtoMessage() {}
func (*SecurityDefinitionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{8}
}
func (m *SecurityDefinitionRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SecurityDefinitionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SecurityDefinitionRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SecurityDefinitionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SecurityDefinitionRequest.Merge(m, src)
}
func (m *SecurityDefinitionRequest) XXX_Size() int {
	return m.Size()
}
func (m *SecurityDefinitionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SecurityDefinitionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SecurityDefinitionRequest proto.InternalMessageInfo

func (m *SecurityDefinitionRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *SecurityDefinitionRequest) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

type SecurityDefinitionResponse struct {
	RequestID       uint64           `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64           `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Security        *models.Security `protobuf:"bytes,3,opt,name=security,proto3" json:"security,omitempty"`
	Success         bool             `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason  `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *SecurityDefinitionResponse) Reset()      { *m = SecurityDefinitionResponse{} }
func (*SecurityDefinitionResponse) ProtoMessage() {}
func (*SecurityDefinitionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{9}
}
func (m *SecurityDefinitionResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SecurityDefinitionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SecurityDefinitionResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SecurityDefinitionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SecurityDefinitionResponse.Merge(m, src)
}
func (m *SecurityDefinitionResponse) XXX_Size() int {
	return m.Size()
}
func (m *SecurityDefinitionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SecurityDefinitionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SecurityDefinitionResponse proto.InternalMessageInfo

func (m *SecurityDefinitionResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *SecurityDefinitionResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *SecurityDefinitionResponse) GetSecurity() *models.Security {
	if m != nil {
		return m.Security
	}
	return nil
}

func (m *SecurityDefinitionResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *SecurityDefinitionResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type SecurityListRequest struct {
	RequestID  uint64     `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe  bool       `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
	Subscriber *actor.PID `protobuf:"bytes,3,opt,name=subscriber,proto3" json:"subscriber,omitempty"`
}

func (m *SecurityListRequest) Reset()      { *m = SecurityListRequest{} }
func (*SecurityListRequest) ProtoMessage() {}
func (*SecurityListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{10}
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
	RequestID       uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64             `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Securities      []*models.Security `protobuf:"bytes,3,rep,name=securities,proto3" json:"securities,omitempty"`
	Success         bool               `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason    `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *SecurityList) Reset()      { *m = SecurityList{} }
func (*SecurityList) ProtoMessage() {}
func (*SecurityList) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{11}
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

func (m *SecurityList) GetSecurities() []*models.Security {
	if m != nil {
		return m.Securities
	}
	return nil
}

func (m *SecurityList) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *SecurityList) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
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
	FeeAmount       *types.DoubleValue `protobuf:"bytes,15,opt,name=fee_amount,json=feeAmount,proto3" json:"fee_amount,omitempty"`
	FeeCurrency     *models1.Asset     `protobuf:"bytes,16,opt,name=fee_currency,json=feeCurrency,proto3" json:"fee_currency,omitempty"`
	FeeType         FeeType            `protobuf:"varint,17,opt,name=fee_type,json=feeType,proto3,enum=messages.FeeType" json:"fee_type,omitempty"`
	FeeBasis        FeeBasis           `protobuf:"varint,18,opt,name=fee_basis,json=feeBasis,proto3,enum=messages.FeeBasis" json:"fee_basis,omitempty"`
	RejectionReason RejectionReason    `protobuf:"varint,19,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *ExecutionReport) Reset()      { *m = ExecutionReport{} }
func (*ExecutionReport) ProtoMessage() {}
func (*ExecutionReport) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{12}
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

func (m *ExecutionReport) GetFeeAmount() *types.DoubleValue {
	if m != nil {
		return m.FeeAmount
	}
	return nil
}

func (m *ExecutionReport) GetFeeCurrency() *models1.Asset {
	if m != nil {
		return m.FeeCurrency
	}
	return nil
}

func (m *ExecutionReport) GetFeeType() FeeType {
	if m != nil {
		return m.FeeType
	}
	return Regulatory
}

func (m *ExecutionReport) GetFeeBasis() FeeBasis {
	if m != nil {
		return m.FeeBasis
	}
	return Absolute
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
	return fileDescriptor_350c53ba9303a7e6, []int{13}
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
	return fileDescriptor_350c53ba9303a7e6, []int{14}
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
	return fileDescriptor_350c53ba9303a7e6, []int{15}
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
	return fileDescriptor_350c53ba9303a7e6, []int{16}
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
	return fileDescriptor_350c53ba9303a7e6, []int{17}
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
	return fileDescriptor_350c53ba9303a7e6, []int{18}
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
	RequestID       uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64             `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Positions       []*models.Position `protobuf:"bytes,3,rep,name=positions,proto3" json:"positions,omitempty"`
	Success         bool               `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason    `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *PositionList) Reset()      { *m = PositionList{} }
func (*PositionList) ProtoMessage() {}
func (*PositionList) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{19}
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

func (m *PositionList) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *PositionList) GetPositions() []*models.Position {
	if m != nil {
		return m.Positions
	}
	return nil
}

func (m *PositionList) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *PositionList) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type BalancesRequest struct {
	RequestID  uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe  bool            `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
	Subscriber *actor.PID      `protobuf:"bytes,3,opt,name=subscriber,proto3" json:"subscriber,omitempty"`
	Asset      *models1.Asset  `protobuf:"bytes,4,opt,name=asset,proto3" json:"asset,omitempty"`
	Account    *models.Account `protobuf:"bytes,5,opt,name=account,proto3" json:"account,omitempty"`
}

func (m *BalancesRequest) Reset()      { *m = BalancesRequest{} }
func (*BalancesRequest) ProtoMessage() {}
func (*BalancesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{20}
}
func (m *BalancesRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BalancesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BalancesRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BalancesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BalancesRequest.Merge(m, src)
}
func (m *BalancesRequest) XXX_Size() int {
	return m.Size()
}
func (m *BalancesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_BalancesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_BalancesRequest proto.InternalMessageInfo

func (m *BalancesRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *BalancesRequest) GetSubscribe() bool {
	if m != nil {
		return m.Subscribe
	}
	return false
}

func (m *BalancesRequest) GetSubscriber() *actor.PID {
	if m != nil {
		return m.Subscriber
	}
	return nil
}

func (m *BalancesRequest) GetAsset() *models1.Asset {
	if m != nil {
		return m.Asset
	}
	return nil
}

func (m *BalancesRequest) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

type BalanceList struct {
	RequestID       uint64            `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64            `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Balances        []*models.Balance `protobuf:"bytes,3,rep,name=balances,proto3" json:"balances,omitempty"`
	Success         bool              `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason   `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *BalanceList) Reset()      { *m = BalanceList{} }
func (*BalanceList) ProtoMessage() {}
func (*BalanceList) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{21}
}
func (m *BalanceList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BalanceList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BalanceList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BalanceList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BalanceList.Merge(m, src)
}
func (m *BalanceList) XXX_Size() int {
	return m.Size()
}
func (m *BalanceList) XXX_DiscardUnknown() {
	xxx_messageInfo_BalanceList.DiscardUnknown(m)
}

var xxx_messageInfo_BalanceList proto.InternalMessageInfo

func (m *BalanceList) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *BalanceList) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *BalanceList) GetBalances() []*models.Balance {
	if m != nil {
		return m.Balances
	}
	return nil
}

func (m *BalanceList) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *BalanceList) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type NewOrder struct {
	ClientOrderID         string                 `protobuf:"bytes,1,opt,name=client_orderID,json=clientOrderID,proto3" json:"client_orderID,omitempty"`
	Instrument            *models.Instrument     `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
	OrderType             models.OrderType       `protobuf:"varint,5,opt,name=order_type,json=orderType,proto3,enum=models.OrderType" json:"order_type,omitempty"`
	OrderSide             models.Side            `protobuf:"varint,6,opt,name=order_side,json=orderSide,proto3,enum=models.Side" json:"order_side,omitempty"`
	TimeInForce           models.TimeInForce     `protobuf:"varint,7,opt,name=time_in_force,json=timeInForce,proto3,enum=models.TimeInForce" json:"time_in_force,omitempty"`
	Quantity              float64                `protobuf:"fixed64,8,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Price                 *types.DoubleValue     `protobuf:"bytes,9,opt,name=price,proto3" json:"price,omitempty"`
	ExecutionInstructions []ExecutionInstruction `protobuf:"varint,10,rep,packed,name=execution_instructions,json=executionInstructions,proto3,enum=messages.ExecutionInstruction" json:"execution_instructions,omitempty"`
}

func (m *NewOrder) Reset()      { *m = NewOrder{} }
func (*NewOrder) ProtoMessage() {}
func (*NewOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{22}
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

func (m *NewOrder) GetExecutionInstructions() []ExecutionInstruction {
	if m != nil {
		return m.ExecutionInstructions
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
	return fileDescriptor_350c53ba9303a7e6, []int{23}
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
	ResponseID      uint64          `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Success         bool            `protobuf:"varint,3,opt,name=success,proto3" json:"success,omitempty"`
	OrderID         string          `protobuf:"bytes,4,opt,name=orderID,proto3" json:"orderID,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *NewOrderSingleResponse) Reset()      { *m = NewOrderSingleResponse{} }
func (*NewOrderSingleResponse) ProtoMessage() {}
func (*NewOrderSingleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{24}
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

func (m *NewOrderSingleResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
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
	return fileDescriptor_350c53ba9303a7e6, []int{25}
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
	ResponseID      uint64          `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	OrderIDs        []string        `protobuf:"bytes,3,rep,name=orderIDs,proto3" json:"orderIDs,omitempty"`
	Success         bool            `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *NewOrderBulkResponse) Reset()      { *m = NewOrderBulkResponse{} }
func (*NewOrderBulkResponse) ProtoMessage() {}
func (*NewOrderBulkResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{26}
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

func (m *NewOrderBulkResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *NewOrderBulkResponse) GetOrderIDs() []string {
	if m != nil {
		return m.OrderIDs
	}
	return nil
}

func (m *NewOrderBulkResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *NewOrderBulkResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type OrderUpdate struct {
	OrderID           *types.StringValue `protobuf:"bytes,1,opt,name=orderID,proto3" json:"orderID,omitempty"`
	OrigClientOrderID *types.StringValue `protobuf:"bytes,2,opt,name=orig_client_orderID,json=origClientOrderID,proto3" json:"orig_client_orderID,omitempty"`
	Quantity          *types.DoubleValue `protobuf:"bytes,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Price             *types.DoubleValue `protobuf:"bytes,4,opt,name=price,proto3" json:"price,omitempty"`
}

func (m *OrderUpdate) Reset()      { *m = OrderUpdate{} }
func (*OrderUpdate) ProtoMessage() {}
func (*OrderUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{27}
}
func (m *OrderUpdate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderUpdate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderUpdate.Merge(m, src)
}
func (m *OrderUpdate) XXX_Size() int {
	return m.Size()
}
func (m *OrderUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_OrderUpdate proto.InternalMessageInfo

func (m *OrderUpdate) GetOrderID() *types.StringValue {
	if m != nil {
		return m.OrderID
	}
	return nil
}

func (m *OrderUpdate) GetOrigClientOrderID() *types.StringValue {
	if m != nil {
		return m.OrigClientOrderID
	}
	return nil
}

func (m *OrderUpdate) GetQuantity() *types.DoubleValue {
	if m != nil {
		return m.Quantity
	}
	return nil
}

func (m *OrderUpdate) GetPrice() *types.DoubleValue {
	if m != nil {
		return m.Price
	}
	return nil
}

type OrderReplaceRequest struct {
	RequestID  uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Instrument *models.Instrument `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Account    *models.Account    `protobuf:"bytes,3,opt,name=account,proto3" json:"account,omitempty"`
	Update     *OrderUpdate       `protobuf:"bytes,4,opt,name=update,proto3" json:"update,omitempty"`
}

func (m *OrderReplaceRequest) Reset()      { *m = OrderReplaceRequest{} }
func (*OrderReplaceRequest) ProtoMessage() {}
func (*OrderReplaceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{28}
}
func (m *OrderReplaceRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderReplaceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderReplaceRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderReplaceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderReplaceRequest.Merge(m, src)
}
func (m *OrderReplaceRequest) XXX_Size() int {
	return m.Size()
}
func (m *OrderReplaceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderReplaceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OrderReplaceRequest proto.InternalMessageInfo

func (m *OrderReplaceRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *OrderReplaceRequest) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *OrderReplaceRequest) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

func (m *OrderReplaceRequest) GetUpdate() *OrderUpdate {
	if m != nil {
		return m.Update
	}
	return nil
}

type OrderReplaceResponse struct {
	RequestID       uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64          `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Success         bool            `protobuf:"varint,3,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,4,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *OrderReplaceResponse) Reset()      { *m = OrderReplaceResponse{} }
func (*OrderReplaceResponse) ProtoMessage() {}
func (*OrderReplaceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{29}
}
func (m *OrderReplaceResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderReplaceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderReplaceResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderReplaceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderReplaceResponse.Merge(m, src)
}
func (m *OrderReplaceResponse) XXX_Size() int {
	return m.Size()
}
func (m *OrderReplaceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderReplaceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_OrderReplaceResponse proto.InternalMessageInfo

func (m *OrderReplaceResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *OrderReplaceResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *OrderReplaceResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *OrderReplaceResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type OrderBulkReplaceRequest struct {
	RequestID  uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Instrument *models.Instrument `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Account    *models.Account    `protobuf:"bytes,3,opt,name=account,proto3" json:"account,omitempty"`
	Updates    []*OrderUpdate     `protobuf:"bytes,4,rep,name=updates,proto3" json:"updates,omitempty"`
}

func (m *OrderBulkReplaceRequest) Reset()      { *m = OrderBulkReplaceRequest{} }
func (*OrderBulkReplaceRequest) ProtoMessage() {}
func (*OrderBulkReplaceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{30}
}
func (m *OrderBulkReplaceRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderBulkReplaceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderBulkReplaceRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderBulkReplaceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderBulkReplaceRequest.Merge(m, src)
}
func (m *OrderBulkReplaceRequest) XXX_Size() int {
	return m.Size()
}
func (m *OrderBulkReplaceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderBulkReplaceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OrderBulkReplaceRequest proto.InternalMessageInfo

func (m *OrderBulkReplaceRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *OrderBulkReplaceRequest) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *OrderBulkReplaceRequest) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

func (m *OrderBulkReplaceRequest) GetUpdates() []*OrderUpdate {
	if m != nil {
		return m.Updates
	}
	return nil
}

type OrderBulkReplaceResponse struct {
	RequestID       uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64          `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Success         bool            `protobuf:"varint,3,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,4,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *OrderBulkReplaceResponse) Reset()      { *m = OrderBulkReplaceResponse{} }
func (*OrderBulkReplaceResponse) ProtoMessage() {}
func (*OrderBulkReplaceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{31}
}
func (m *OrderBulkReplaceResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderBulkReplaceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderBulkReplaceResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderBulkReplaceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderBulkReplaceResponse.Merge(m, src)
}
func (m *OrderBulkReplaceResponse) XXX_Size() int {
	return m.Size()
}
func (m *OrderBulkReplaceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderBulkReplaceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_OrderBulkReplaceResponse proto.InternalMessageInfo

func (m *OrderBulkReplaceResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *OrderBulkReplaceResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *OrderBulkReplaceResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *OrderBulkReplaceResponse) GetRejectionReason() RejectionReason {
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
	return fileDescriptor_350c53ba9303a7e6, []int{32}
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
	ResponseID      uint64          `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Success         bool            `protobuf:"varint,3,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,4,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *OrderCancelResponse) Reset()      { *m = OrderCancelResponse{} }
func (*OrderCancelResponse) ProtoMessage() {}
func (*OrderCancelResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{33}
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

func (m *OrderCancelResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
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
	RequestID uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Account   *models.Account `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	Filter    *OrderFilter    `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (m *OrderMassCancelRequest) Reset()      { *m = OrderMassCancelRequest{} }
func (*OrderMassCancelRequest) ProtoMessage() {}
func (*OrderMassCancelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{34}
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

func (m *OrderMassCancelRequest) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

func (m *OrderMassCancelRequest) GetFilter() *OrderFilter {
	if m != nil {
		return m.Filter
	}
	return nil
}

type OrderMassCancelResponse struct {
	RequestID       uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64          `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Success         bool            `protobuf:"varint,3,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,4,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *OrderMassCancelResponse) Reset()      { *m = OrderMassCancelResponse{} }
func (*OrderMassCancelResponse) ProtoMessage() {}
func (*OrderMassCancelResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{35}
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

func (m *OrderMassCancelResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
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
	proto.RegisterEnum("messages.FeeType", FeeType_name, FeeType_value)
	proto.RegisterEnum("messages.FeeBasis", FeeBasis_name, FeeBasis_value)
	proto.RegisterEnum("messages.ExecutionInstruction", ExecutionInstruction_name, ExecutionInstruction_value)
	proto.RegisterType((*HistoricalLiquidationsRequest)(nil), "messages.HistoricalLiquidationsRequest")
	proto.RegisterType((*HistoricalLiquidationsResponse)(nil), "messages.HistoricalLiquidationsResponse")
	proto.RegisterType((*MarketDataRequest)(nil), "messages.MarketDataRequest")
	proto.RegisterType((*MarketDataResponse)(nil), "messages.MarketDataResponse")
	proto.RegisterType((*MarketDataIncrementalRefresh)(nil), "messages.MarketDataIncrementalRefresh")
	proto.RegisterType((*AccountDataRequest)(nil), "messages.AccountDataRequest")
	proto.RegisterType((*AccountDataResponse)(nil), "messages.AccountDataResponse")
	proto.RegisterType((*AccountDataIncrementalRefresh)(nil), "messages.AccountDataIncrementalRefresh")
	proto.RegisterType((*SecurityDefinitionRequest)(nil), "messages.SecurityDefinitionRequest")
	proto.RegisterType((*SecurityDefinitionResponse)(nil), "messages.SecurityDefinitionResponse")
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
	proto.RegisterType((*BalancesRequest)(nil), "messages.BalancesRequest")
	proto.RegisterType((*BalanceList)(nil), "messages.BalanceList")
	proto.RegisterType((*NewOrder)(nil), "messages.NewOrder")
	proto.RegisterType((*NewOrderSingleRequest)(nil), "messages.NewOrderSingleRequest")
	proto.RegisterType((*NewOrderSingleResponse)(nil), "messages.NewOrderSingleResponse")
	proto.RegisterType((*NewOrderBulkRequest)(nil), "messages.NewOrderBulkRequest")
	proto.RegisterType((*NewOrderBulkResponse)(nil), "messages.NewOrderBulkResponse")
	proto.RegisterType((*OrderUpdate)(nil), "messages.OrderUpdate")
	proto.RegisterType((*OrderReplaceRequest)(nil), "messages.OrderReplaceRequest")
	proto.RegisterType((*OrderReplaceResponse)(nil), "messages.OrderReplaceResponse")
	proto.RegisterType((*OrderBulkReplaceRequest)(nil), "messages.OrderBulkReplaceRequest")
	proto.RegisterType((*OrderBulkReplaceResponse)(nil), "messages.OrderBulkReplaceResponse")
	proto.RegisterType((*OrderCancelRequest)(nil), "messages.OrderCancelRequest")
	proto.RegisterType((*OrderCancelResponse)(nil), "messages.OrderCancelResponse")
	proto.RegisterType((*OrderMassCancelRequest)(nil), "messages.OrderMassCancelRequest")
	proto.RegisterType((*OrderMassCancelResponse)(nil), "messages.OrderMassCancelResponse")
}

func init() { proto.RegisterFile("executor_messages.proto", fileDescriptor_350c53ba9303a7e6) }

var fileDescriptor_350c53ba9303a7e6 = []byte{
	// 2792 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x5a, 0xcd, 0x6f, 0x1b, 0xc7,
	0x15, 0xe7, 0x8a, 0x92, 0x48, 0x3e, 0x52, 0xd4, 0x6a, 0x24, 0x4b, 0x0c, 0x23, 0x33, 0xca, 0xa6,
	0x69, 0x1c, 0xa5, 0xa1, 0x6c, 0xda, 0x71, 0x02, 0xb4, 0x31, 0x20, 0x8b, 0x32, 0x42, 0x40, 0xb1,
	0xd5, 0x95, 0x54, 0x04, 0xe8, 0x81, 0x18, 0x2d, 0x47, 0xd4, 0xd6, 0xcb, 0x5d, 0x7a, 0x67, 0x36,
	0xb6, 0x0e, 0x05, 0xda, 0x1e, 0x73, 0x28, 0xf2, 0x17, 0xf4, 0xd4, 0xaf, 0x4b, 0x81, 0x02, 0x2d,
	0xda, 0x1e, 0x7a, 0xe9, 0xa9, 0xed, 0xa1, 0x80, 0x4f, 0x45, 0x0a, 0x14, 0x45, 0xa3, 0x5c, 0x7c,
	0x6b, 0x6e, 0xbd, 0x16, 0xf3, 0xb1, 0xcb, 0xe1, 0x87, 0x2c, 0x2a, 0x96, 0xd2, 0xf8, 0xb6, 0x33,
	0xf3, 0x7b, 0x33, 0xf3, 0x3e, 0xe7, 0xbd, 0x99, 0x85, 0x25, 0xf2, 0x88, 0x38, 0x11, 0x0b, 0xc2,
	0x66, 0x87, 0x50, 0x8a, 0xdb, 0x84, 0x56, 0xbb, 0x61, 0xc0, 0x02, 0x94, 0x8d, 0xdb, 0xe5, 0x97,
	0xda, 0x41, 0xd0, 0xf6, 0xc8, 0x9a, 0xe8, 0xdf, 0x8f, 0x0e, 0xd6, 0x98, 0xdb, 0x21, 0x94, 0xe1,
	0x4e, 0x57, 0x42, 0xcb, 0x95, 0x41, 0xc0, 0xc3, 0x10, 0x77, 0xbb, 0x24, 0x54, 0x53, 0x95, 0x6b,
	0x6d, 0x97, 0x79, 0x78, 0xbf, 0xea, 0x04, 0x9d, 0x35, 0xec, 0x75, 0x0f, 0x31, 0x73, 0x9d, 0xfb,
	0x74, 0xed, 0x91, 0x73, 0x88, 0xfd, 0x36, 0x09, 0xd7, 0x3a, 0x41, 0x8b, 0x78, 0x54, 0x92, 0xc7,
	0x34, 0xdf, 0x1a, 0x4d, 0x23, 0x3e, 0xdf, 0x74, 0x02, 0xdf, 0x27, 0x0e, 0x8b, 0x09, 0x3b, 0x38,
	0xbc, 0x4f, 0x58, 0xb3, 0x85, 0x19, 0x56, 0xd4, 0xb7, 0xce, 0x40, 0x4d, 0x89, 0x13, 0x85, 0x2e,
	0x3b, 0xd2, 0xe9, 0xdf, 0x39, 0x03, 0x3d, 0x76, 0x9c, 0x20, 0xf2, 0x99, 0xa2, 0xbc, 0xd9, 0x76,
	0xd9, 0x61, 0x24, 0x29, 0xd7, 0xe9, 0x91, 0x7f, 0x3f, 0x0c, 0xfc, 0xc6, 0xae, 0xe4, 0x0d, 0x3b,
	0x2c, 0x08, 0xdf, 0x6c, 0x07, 0x6b, 0xe2, 0xa3, 0x8f, 0x5f, 0xeb, 0x6f, 0x06, 0x5c, 0x7e, 0xcf,
	0xa5, 0x2c, 0x08, 0x5d, 0x07, 0x7b, 0x5b, 0xee, 0x83, 0xc8, 0x6d, 0x61, 0xe6, 0x06, 0x3e, 0xb5,
	0xc9, 0x83, 0x88, 0x50, 0x86, 0x96, 0x21, 0x17, 0xca, 0xcf, 0x46, 0xbd, 0x64, 0xac, 0x18, 0x57,
	0x26, 0xed, 0x5e, 0x07, 0xaa, 0x01, 0xb8, 0x3e, 0x65, 0x61, 0xd4, 0x21, 0x3e, 0x2b, 0x4d, 0xac,
	0x18, 0x57, 0xf2, 0x35, 0x54, 0x95, 0x5b, 0xac, 0x36, 0x92, 0x11, 0x5b, 0x43, 0xa1, 0x2a, 0x4c,
	0x1e, 0x84, 0x41, 0xa7, 0x94, 0x16, 0xe8, 0x72, 0x55, 0xaa, 0xb1, 0x1a, 0xab, 0xb1, 0xba, 0x1b,
	0xeb, 0xd9, 0x16, 0x38, 0xb4, 0x0a, 0x13, 0x2c, 0x28, 0x4d, 0x9e, 0x8a, 0x9e, 0x60, 0x81, 0xf5,
	0x5f, 0x03, 0x2a, 0x27, 0xf1, 0x43, 0xbb, 0x81, 0x4f, 0xc9, 0x29, 0x0c, 0x55, 0x00, 0x42, 0x85,
	0x6c, 0xd4, 0x05, 0x43, 0x93, 0xb6, 0xd6, 0x83, 0xde, 0x86, 0x82, 0xa7, 0xcd, 0x5a, 0x4a, 0xaf,
	0xa4, 0xaf, 0xe4, 0x6b, 0xf3, 0x31, 0xcb, 0xda, 0x8a, 0x76, 0x1f, 0x10, 0x95, 0x20, 0x43, 0x23,
	0xc7, 0x21, 0x94, 0x0a, 0x56, 0xb2, 0x76, 0xdc, 0x44, 0x75, 0x30, 0x43, 0xf2, 0x3d, 0xe2, 0x70,
	0x5c, 0x33, 0x24, 0x98, 0x06, 0x7e, 0x69, 0x6a, 0xc5, 0xb8, 0x52, 0xac, 0xbd, 0x50, 0x4d, 0xbc,
	0xc3, 0x8e, 0x11, 0xb6, 0x00, 0xd8, 0xb3, 0x61, 0x7f, 0x87, 0xf5, 0x1f, 0x03, 0xe6, 0xde, 0x17,
	0x16, 0x59, 0xc7, 0x0c, 0x8f, 0xa7, 0xbd, 0x65, 0xc8, 0xd1, 0x68, 0x9f, 0x3a, 0xa1, 0xbb, 0x4f,
	0x04, 0xaf, 0x59, 0xbb, 0xd7, 0x81, 0x56, 0x01, 0x92, 0x46, 0xa8, 0xb4, 0x05, 0x55, 0x61, 0x44,
	0xd5, 0xed, 0x46, 0xdd, 0xd6, 0x46, 0x07, 0xec, 0x60, 0x72, 0x2c, 0x3b, 0xb8, 0x05, 0x79, 0xdc,
	0x6e, 0x87, 0xa4, 0x2d, 0x24, 0xa4, 0x58, 0x5e, 0x8e, 0x89, 0xee, 0x85, 0x2d, 0x12, 0xde, 0x0e,
	0x82, 0xfb, 0xeb, 0x3d, 0x8c, 0xad, 0x13, 0x58, 0x3f, 0x49, 0x03, 0xd2, 0x39, 0x3e, 0x17, 0xfd,
	0xde, 0x00, 0xa0, 0x3e, 0xee, 0xd2, 0xc3, 0x80, 0x6d, 0x5d, 0x53, 0x4c, 0x2f, 0x24, 0x7b, 0xba,
	0xbd, 0x75, 0x6d, 0x47, 0x8d, 0xda, 0x1a, 0xae, 0x8f, 0xaa, 0xa6, 0xd8, 0xd7, 0xa9, 0x6a, 0x23,
	0xa8, 0x6a, 0x7d, 0x54, 0xd7, 0x05, 0xff, 0xfd, 0x54, 0xd7, 0x47, 0x50, 0x5d, 0x47, 0x6b, 0x30,
	0xcd, 0x42, 0xdc, 0x22, 0xb4, 0x34, 0x2d, 0x6c, 0x6f, 0x29, 0xa6, 0x88, 0x05, 0x45, 0x5a, 0xbb,
	0x7c, 0xdc, 0x56, 0x30, 0xb4, 0x04, 0x19, 0x4a, 0x1e, 0x34, 0xfd, 0xa8, 0x53, 0xca, 0x08, 0x7e,
	0xa7, 0x29, 0x79, 0x70, 0x37, 0xea, 0xe8, 0x26, 0x99, 0x3d, 0xdd, 0x24, 0x73, 0x67, 0x36, 0xc9,
	0x3f, 0xa6, 0x61, 0xb9, 0xa7, 0xa0, 0x86, 0xef, 0x84, 0x84, 0x2b, 0x1e, 0x7b, 0x36, 0x39, 0x08,
	0x09, 0x3d, 0x7c, 0x46, 0x55, 0x69, 0x7c, 0xa5, 0xfb, 0xf8, 0xaa, 0x42, 0x36, 0xea, 0xb6, 0x30,
	0x23, 0x5b, 0xd7, 0x06, 0x4d, 0x91, 0x6b, 0x70, 0x4f, 0x8c, 0xd9, 0x09, 0x46, 0xc3, 0xd7, 0x94,
	0x16, 0x74, 0x7c, 0x6d, 0x00, 0x5f, 0xd3, 0xf0, 0xd7, 0x4b, 0xd3, 0x43, 0xf8, 0xeb, 0x03, 0x78,
	0x5d, 0x63, 0x99, 0xf1, 0x34, 0xf6, 0x16, 0xe4, 0xb5, 0xd8, 0x21, 0x94, 0x73, 0x42, 0x8c, 0xd1,
	0x71, 0xe8, 0x75, 0xc8, 0x1c, 0x44, 0x7e, 0xcb, 0xf5, 0xdb, 0x42, 0x59, 0xf9, 0xda, 0x6c, 0x4c,
	0x72, 0x47, 0x76, 0xdb, 0xf1, 0x38, 0xb2, 0x60, 0x8a, 0x32, 0xcc, 0x68, 0x09, 0xc4, 0x8e, 0x0a,
	0x31, 0x70, 0x87, 0x61, 0x66, 0xcb, 0x21, 0xeb, 0x17, 0x06, 0xa0, 0x75, 0x79, 0xca, 0xfc, 0x3f,
	0x42, 0xca, 0xeb, 0x90, 0x51, 0x67, 0x9c, 0x52, 0x62, 0xc2, 0x8d, 0xda, 0x94, 0x1d, 0x8f, 0x5b,
	0x3f, 0x9a, 0x80, 0xf9, 0xbe, 0x9d, 0x9e, 0x4b, 0x28, 0x78, 0x03, 0xa6, 0x03, 0x1e, 0x84, 0xa8,
	0xda, 0xe8, 0x7c, 0xcf, 0xf4, 0x45, 0x70, 0xda, 0x72, 0x29, 0xb3, 0x15, 0x04, 0xdd, 0x80, 0x5c,
	0x37, 0xa0, 0xae, 0x3c, 0x14, 0xe4, 0x7e, 0x17, 0x7b, 0xf8, 0x6d, 0x35, 0x24, 0x48, 0x7a, 0x40,
	0x74, 0x0d, 0xb2, 0xfb, 0xd8, 0xc3, 0xbe, 0x43, 0xa8, 0xb2, 0xbc, 0x4b, 0x3d, 0xa2, 0xdb, 0x72,
	0x44, 0xd0, 0x24, 0x30, 0xdd, 0xea, 0xa7, 0x75, 0xab, 0xb7, 0x3e, 0x36, 0xe0, 0xb2, 0x26, 0x84,
	0x73, 0x77, 0xb7, 0x6b, 0x30, 0x1d, 0x92, 0x6e, 0x10, 0x32, 0x25, 0x0e, 0x2d, 0x12, 0x6c, 0x8a,
	0x64, 0x4e, 0x38, 0x3e, 0x07, 0xd8, 0x0a, 0x68, 0x75, 0xe0, 0x85, 0x1d, 0x95, 0xe6, 0xd4, 0xc9,
	0x81, 0xeb, 0xbb, 0x12, 0x73, 0x41, 0x89, 0x85, 0xf5, 0xc4, 0x80, 0xf2, 0xa8, 0xf5, 0xce, 0xc5,
	0x1a, 0xbe, 0x01, 0xd9, 0x38, 0x65, 0x53, 0x02, 0x30, 0x13, 0xa7, 0x51, 0xfd, 0x76, 0x82, 0xb8,
	0xf0, 0xd3, 0xfe, 0xfb, 0x30, 0x1f, 0xaf, 0x2a, 0xec, 0xe3, 0xcb, 0xf5, 0x4d, 0xeb, 0x5f, 0x06,
	0x14, 0xf4, 0xf5, 0x9f, 0x51, 0xb6, 0x57, 0x01, 0x94, 0xe4, 0x5c, 0x12, 0xa7, 0x54, 0xc3, 0xd2,
	0xd5, 0x30, 0x17, 0x2e, 0xdf, 0xdf, 0x65, 0x60, 0x76, 0xc0, 0xaa, 0x75, 0xcf, 0x33, 0x06, 0xcf,
	0x51, 0x11, 0x05, 0x14, 0x6f, 0x39, 0x3b, 0x6e, 0xa2, 0x0d, 0x28, 0x3a, 0x9e, 0x4b, 0x7c, 0xd6,
	0x8c, 0x01, 0x52, 0xae, 0xcb, 0x43, 0x69, 0xec, 0x0e, 0x0b, 0x5d, 0xbf, 0xfd, 0x1d, 0xec, 0x45,
	0xc4, 0x9e, 0x91, 0x34, 0xf7, 0xd4, 0x24, 0x2b, 0x90, 0x27, 0xf1, 0x56, 0x1a, 0x75, 0xc1, 0x6f,
	0xce, 0xd6, 0xbb, 0xd0, 0x2d, 0x28, 0x26, 0xcd, 0x26, 0x3b, 0xea, 0x12, 0xc5, 0xf1, 0xd2, 0x08,
	0x17, 0xdd, 0x3d, 0xea, 0x12, 0x7b, 0x86, 0xe8, 0x4d, 0x74, 0x13, 0x0a, 0x62, 0x7f, 0x4d, 0x1e,
	0xf8, 0x23, 0x2a, 0x02, 0x4b, 0xb1, 0x77, 0xe0, 0x88, 0x8d, 0xec, 0x88, 0x21, 0x3b, 0x1f, 0xf4,
	0x1a, 0x03, 0x4e, 0x9a, 0x19, 0x2b, 0xeb, 0x7b, 0x0d, 0x66, 0x3d, 0x82, 0x3f, 0x24, 0xb4, 0xf9,
	0x20, 0xc2, 0x3e, 0xe3, 0xee, 0xc4, 0x0f, 0x2b, 0xc3, 0x2e, 0xca, 0xee, 0x6f, 0xab, 0x5e, 0xf4,
	0x32, 0x14, 0x9c, 0xa8, 0xd3, 0x43, 0x81, 0x40, 0xe5, 0x9d, 0xa8, 0x93, 0x40, 0x36, 0xc1, 0x64,
	0x21, 0xf6, 0x29, 0x96, 0xda, 0xe6, 0x05, 0x62, 0x29, 0x7f, 0x6a, 0x9d, 0x30, 0xab, 0xd1, 0xf0,
	0x5e, 0x74, 0x13, 0x32, 0xe2, 0xe0, 0x6d, 0xd4, 0x4b, 0x85, 0x31, 0xd4, 0x13, 0x83, 0xd1, 0x37,
	0x01, 0x0e, 0x5c, 0xcf, 0x6b, 0x76, 0x43, 0xd7, 0x21, 0xa5, 0x99, 0x13, 0x48, 0xeb, 0x41, 0xb4,
	0xef, 0x11, 0x49, 0x9a, 0xe3, 0xf8, 0x6d, 0x0e, 0x47, 0xeb, 0x30, 0x23, 0x88, 0x13, 0xfe, 0x8a,
	0x63, 0xd0, 0x17, 0x38, 0x49, 0xc2, 0x3e, 0x5f, 0x9f, 0x90, 0x26, 0xee, 0x88, 0x43, 0x72, 0x76,
	0xac, 0xf5, 0x09, 0x59, 0x17, 0x70, 0x74, 0x15, 0x0a, 0x9c, 0xd8, 0x89, 0xc2, 0x90, 0xf8, 0xce,
	0x51, 0xc9, 0x14, 0xe4, 0x33, 0xc9, 0x19, 0x4b, 0x29, 0x61, 0x76, 0xfe, 0x80, 0x90, 0x0d, 0x85,
	0xe0, 0x11, 0x90, 0x53, 0x08, 0xfb, 0x9a, 0x13, 0x16, 0x32, 0xd7, 0xb3, 0xaf, 0x3b, 0x84, 0x08,
	0xcb, 0xca, 0x1c, 0xc8, 0x0f, 0xb4, 0x06, 0x7c, 0xb1, 0xe6, 0x3e, 0xa6, 0x2e, 0x2d, 0x21, 0x01,
	0x47, 0x7d, 0xf0, 0xdb, 0x7c, 0xc4, 0xe6, 0x53, 0x8a, 0xaf, 0x91, 0x8e, 0x3b, 0x7f, 0x66, 0xc7,
	0x5d, 0x83, 0xdc, 0x8e, 0xdb, 0x92, 0xec, 0xf2, 0x2c, 0xe7, 0x43, 0xfe, 0x21, 0xfc, 0xb5, 0xa8,
	0x65, 0x39, 0x6e, 0x8b, 0xd8, 0x72, 0xc8, 0x7a, 0x17, 0x4c, 0xcd, 0xbe, 0x25, 0xdd, 0xeb, 0xfd,
	0x74, 0x23, 0x1d, 0x41, 0x91, 0xff, 0x7c, 0x02, 0xf2, 0xa2, 0xfb, 0x8e, 0xeb, 0x31, 0x12, 0x72,
	0x5b, 0x8a, 0x5d, 0xdd, 0x18, 0xc7, 0x96, 0x4e, 0x8e, 0x14, 0x13, 0x67, 0x8f, 0x14, 0xfd, 0xfe,
	0x98, 0x1e, 0xd3, 0x1f, 0x27, 0xa9, 0xdb, 0x22, 0x2a, 0x67, 0xd1, 0x72, 0x9c, 0x44, 0x8c, 0xb6,
	0x00, 0xa0, 0x77, 0x07, 0x82, 0xc4, 0x94, 0x72, 0xb4, 0xfe, 0xa4, 0x48, 0x13, 0x63, 0x5f, 0xac,
	0xb0, 0xfe, 0x6e, 0x00, 0xd2, 0xe5, 0xf7, 0x95, 0xcd, 0x26, 0xd1, 0x9b, 0x30, 0x7d, 0x20, 0x94,
	0x39, 0x9c, 0x92, 0x69, 0x9a, 0xb6, 0x15, 0xc8, 0xfa, 0xb5, 0x01, 0xb9, 0x24, 0x1f, 0x3c, 0x85,
	0x1f, 0xed, 0xd8, 0x9a, 0xe8, 0x3f, 0xb6, 0x5e, 0xd5, 0x92, 0xcd, 0xb4, 0xee, 0x88, 0x62, 0xea,
	0x24, 0xcd, 0x1c, 0xe5, 0x24, 0x93, 0x67, 0x76, 0x92, 0x7f, 0x18, 0x60, 0xc6, 0x29, 0x29, 0x7d,
	0x1e, 0xae, 0x0a, 0x34, 0xed, 0x4d, 0x9d, 0x52, 0x0b, 0xfc, 0xd3, 0x80, 0x82, 0x9e, 0x6e, 0x3f,
	0x63, 0x6a, 0x52, 0xd5, 0xf3, 0xfa, 0x81, 0xcc, 0x24, 0x5e, 0x46, 0xcf, 0xe8, 0x2f, 0x3a, 0x31,
	0xf9, 0xab, 0x01, 0xb3, 0xaa, 0x30, 0xf8, 0xd2, 0x35, 0xf7, 0x0a, 0x4c, 0x61, 0x7e, 0x2c, 0x28,
	0xa5, 0x0d, 0x9c, 0x15, 0x72, 0xec, 0x2c, 0xaa, 0xfa, 0xc4, 0x80, 0xbc, 0x56, 0xe4, 0x3c, 0x73,
	0xb9, 0xd6, 0xab, 0xa5, 0xa4, 0xa2, 0x92, 0x95, 0xd5, 0x22, 0x5a, 0x15, 0x75, 0xd1, 0x6a, 0xfa,
	0x69, 0x1a, 0xb2, 0x77, 0xc9, 0x43, 0xe1, 0xbc, 0xe8, 0xd5, 0xa1, 0xd8, 0x6e, 0x88, 0x1c, 0xee,
	0xa9, 0xd1, 0x7b, 0xbc, 0xbb, 0xd4, 0xab, 0x00, 0x32, 0x28, 0x6b, 0x59, 0xdf, 0x5c, 0x5f, 0xe8,
	0x10, 0xa7, 0x72, 0x2e, 0x88, 0x3f, 0xd1, 0x1b, 0x31, 0x85, 0x88, 0xfa, 0xd3, 0x23, 0x0e, 0x46,
	0x09, 0xe6, 0x9f, 0xe8, 0x6d, 0x98, 0xe1, 0x49, 0x55, 0xd3, 0xf5, 0x9b, 0x07, 0x41, 0xe8, 0x10,
	0x91, 0xe3, 0x69, 0x07, 0x22, 0x4f, 0x9f, 0x1a, 0xfe, 0x1d, 0x3e, 0x64, 0xe7, 0x59, 0xaf, 0x81,
	0xca, 0x90, 0x4d, 0x12, 0x9b, 0xac, 0x48, 0xdc, 0x92, 0x36, 0xaa, 0xc1, 0x94, 0xcc, 0x98, 0x72,
	0x63, 0x64, 0x2c, 0x12, 0x8a, 0xf6, 0x60, 0xb1, 0x97, 0xe1, 0x4a, 0xfe, 0x1d, 0xe9, 0x93, 0xb0,
	0x92, 0xbe, 0x52, 0xac, 0x55, 0x46, 0x64, 0xba, 0x8d, 0x1e, 0xcc, 0xbe, 0x44, 0x46, 0xf4, 0x52,
	0xeb, 0x23, 0x03, 0x2e, 0xc5, 0x6a, 0xda, 0x71, 0xfd, 0xb6, 0x47, 0xc6, 0xf3, 0x29, 0xcd, 0xc8,
	0x27, 0x4e, 0x39, 0x4d, 0xae, 0xc0, 0x94, 0x90, 0x67, 0xef, 0x38, 0x8e, 0x37, 0x1a, 0x2f, 0x6c,
	0x4b, 0x80, 0xf5, 0xd8, 0x80, 0xc5, 0xc1, 0xcd, 0x9c, 0x4b, 0xe9, 0xaa, 0x19, 0x7b, 0xba, 0xdf,
	0xd8, 0xb5, 0xca, 0x65, 0xb2, 0xbf, 0x72, 0x39, 0x1f, 0x37, 0xf8, 0xc8, 0x80, 0xf9, 0x98, 0xa5,
	0xdb, 0x91, 0x77, 0xff, 0xdc, 0xa5, 0xbb, 0x3a, 0x70, 0x6c, 0x8e, 0x12, 0xaf, 0x42, 0x70, 0xf9,
	0x2e, 0xf4, 0x6f, 0xe6, 0x5c, 0xa4, 0x5b, 0x86, 0xac, 0x12, 0x9a, 0xdc, 0x44, 0xce, 0x4e, 0xda,
	0x17, 0x1e, 0x66, 0x7e, 0x18, 0x67, 0x9f, 0xf2, 0xce, 0xf1, 0x0b, 0x67, 0x9f, 0x5b, 0x30, 0x1f,
	0x84, 0x6e, 0xbb, 0xf9, 0x05, 0x52, 0xd0, 0x39, 0x4e, 0xb8, 0xd1, 0x17, 0xc8, 0xde, 0xd1, 0x9c,
	0x3f, 0x3d, 0x86, 0x8f, 0x8f, 0x08, 0x0d, 0x93, 0x63, 0x87, 0x06, 0xeb, 0x4f, 0x06, 0xcc, 0x4b,
	0x45, 0x93, 0xae, 0x87, 0x1d, 0x72, 0x71, 0x0f, 0x57, 0x9a, 0x5d, 0xa6, 0x4f, 0xcf, 0x21, 0xe5,
	0xf5, 0xaf, 0xe2, 0x64, 0x30, 0x87, 0x54, 0x77, 0xc4, 0x0a, 0x64, 0xfd, 0xd6, 0x80, 0x85, 0x7e,
	0x1e, 0x2e, 0xd8, 0xf1, 0xcf, 0x27, 0x8f, 0xfc, 0xb3, 0x01, 0x4b, 0x9a, 0x3b, 0x7d, 0x75, 0xc4,
	0xbf, 0x06, 0x19, 0x29, 0x59, 0xee, 0x77, 0xe9, 0x93, 0xe5, 0x1f, 0xa3, 0xac, 0x3f, 0x18, 0x50,
	0x1a, 0xe6, 0xe4, 0xb9, 0x50, 0xc2, 0x8f, 0x27, 0x54, 0x61, 0xb5, 0xc1, 0x53, 0x1b, 0x6f, 0x3c,
	0xf9, 0xdf, 0xec, 0xbf, 0xb2, 0x7a, 0x86, 0x32, 0x35, 0xfd, 0xac, 0x65, 0xea, 0xb9, 0x57, 0x00,
	0xbf, 0x89, 0x03, 0x42, 0x2c, 0x90, 0xe7, 0x42, 0x8d, 0x1f, 0x1b, 0xb0, 0x28, 0x76, 0xfd, 0x3e,
	0xa6, 0xf4, 0x2c, 0xaa, 0x3c, 0xc3, 0x69, 0xd9, 0xab, 0x6c, 0xd3, 0xe3, 0x54, 0xb6, 0xbf, 0x8f,
	0xdd, 0x5b, 0xdf, 0xd2, 0xf3, 0x20, 0xcc, 0xd5, 0x27, 0x06, 0xcc, 0xf4, 0xdd, 0x78, 0xa2, 0x0c,
	0xa4, 0xef, 0x92, 0x87, 0x66, 0x0a, 0x65, 0x61, 0xb2, 0x1e, 0xf8, 0xc4, 0x34, 0x50, 0x01, 0xb2,
	0x92, 0x29, 0xd2, 0x32, 0x27, 0x78, 0x4b, 0xf9, 0x7d, 0xcb, 0x4c, 0xa3, 0x39, 0x98, 0xd9, 0x26,
	0xe2, 0xa9, 0x4c, 0x42, 0xcc, 0x49, 0x94, 0x87, 0xcc, 0x0e, 0x0b, 0xba, 0x5d, 0xd2, 0x32, 0xa7,
	0x24, 0x9a, 0xaf, 0x49, 0x5a, 0xe6, 0x34, 0x9a, 0x81, 0xdc, 0x4e, 0x44, 0xbb, 0xc4, 0x6f, 0x91,
	0x96, 0x99, 0x41, 0x45, 0x00, 0x45, 0xcc, 0x97, 0xcc, 0xf2, 0xf6, 0x06, 0xf6, 0x9c, 0xc8, 0xc3,
	0x1c, 0x9e, 0xe3, 0x33, 0x6d, 0x3e, 0xea, 0xba, 0x21, 0x69, 0x99, 0x80, 0x10, 0x14, 0x15, 0x58,
	0x2d, 0x6f, 0xe6, 0x51, 0x0e, 0xa6, 0xc4, 0x83, 0xa0, 0x59, 0x40, 0xb3, 0xea, 0x80, 0x97, 0xb7,
	0x26, 0xe6, 0x0c, 0x9f, 0x6c, 0x87, 0x30, 0xe6, 0x89, 0xb7, 0x1d, 0xb3, 0xb8, 0xfa, 0xab, 0x49,
	0x98, 0x1d, 0x90, 0x07, 0xa7, 0xbf, 0xc7, 0x0e, 0x49, 0x68, 0xa6, 0x38, 0x23, 0x7b, 0xfe, 0x7d,
	0x3f, 0x78, 0xe8, 0xef, 0x1c, 0x75, 0xf6, 0x03, 0xcf, 0x34, 0xd0, 0x25, 0x98, 0x8b, 0xbb, 0xd4,
	0xd5, 0x7a, 0xa3, 0x6e, 0x4e, 0x20, 0x0b, 0x2a, 0x7b, 0x3e, 0x8d, 0xba, 0xdd, 0x20, 0x64, 0xa4,
	0x25, 0x1d, 0xe8, 0x10, 0x87, 0xd8, 0x61, 0x24, 0x74, 0x29, 0x73, 0x1d, 0x33, 0xcd, 0x49, 0x1b,
	0xbe, 0x13, 0x84, 0x21, 0x71, 0x58, 0x7c, 0x0d, 0x69, 0x4e, 0x72, 0x1e, 0x36, 0xd5, 0x6f, 0x35,
	0x1b, 0x5e, 0x40, 0x85, 0x84, 0x10, 0x14, 0xeb, 0x51, 0xd7, 0x73, 0x1d, 0xcc, 0x88, 0x98, 0xcc,
	0x9c, 0xe6, 0x7d, 0x0d, 0xff, 0x43, 0xec, 0xb9, 0x2d, 0x65, 0x9a, 0x66, 0x06, 0xcd, 0xc3, 0xec,
	0x6e, 0x10, 0x6c, 0x61, 0x46, 0x76, 0x03, 0x25, 0xeb, 0x2c, 0x32, 0xa1, 0xa0, 0xb6, 0x28, 0x49,
	0x73, 0xa8, 0x04, 0x0b, 0x72, 0x74, 0xdd, 0x0b, 0x09, 0x6e, 0x1d, 0x29, 0x99, 0x99, 0x80, 0x16,
	0xc0, 0xac, 0xbb, 0x07, 0x07, 0x24, 0x24, 0x3e, 0x93, 0x3c, 0x52, 0x33, 0xaf, 0x2d, 0xa5, 0xbc,
	0xc8, 0x2c, 0x70, 0x64, 0xbc, 0xcd, 0xf5, 0xed, 0xc6, 0x66, 0x18, 0x06, 0xa1, 0x39, 0xc3, 0xd7,
	0x52, 0x48, 0xb9, 0x56, 0x91, 0x73, 0x69, 0x63, 0x46, 0xb6, 0xdc, 0x8e, 0xcb, 0x36, 0x1f, 0x39,
	0x84, 0x70, 0xb5, 0xce, 0xa2, 0x17, 0x61, 0x49, 0x13, 0xd0, 0x8e, 0xac, 0x8b, 0xbb, 0x5c, 0xec,
	0xa6, 0xc9, 0x69, 0xde, 0x77, 0x29, 0x75, 0xfd, 0x76, 0x2f, 0x80, 0x99, 0x73, 0xdc, 0x32, 0xde,
	0xdb, 0xdd, 0xdd, 0x96, 0x6b, 0x21, 0xb4, 0x04, 0xf3, 0x77, 0xc5, 0x7b, 0x02, 0x57, 0x34, 0xde,
	0xf7, 0x94, 0x64, 0xe6, 0xd1, 0x22, 0xa0, 0xbb, 0x81, 0x2f, 0x39, 0xec, 0xf5, 0x2f, 0x48, 0x5d,
	0x25, 0x6b, 0x4a, 0xff, 0x34, 0x2f, 0x71, 0x69, 0x0c, 0xea, 0x8a, 0x5b, 0xb9, 0xb9, 0x88, 0x5e,
	0x82, 0x17, 0x87, 0x46, 0x7a, 0x75, 0x99, 0xb9, 0xb4, 0xda, 0x80, 0x8c, 0xba, 0xab, 0xe5, 0xa6,
	0x64, 0x93, 0x36, 0x37, 0xcb, 0x20, 0x3c, 0x32, 0x53, 0xdc, 0x47, 0x76, 0xf1, 0x23, 0xd3, 0xe0,
	0x3a, 0xd9, 0x0a, 0x1c, 0xec, 0x6d, 0x04, 0x9d, 0x0e, 0x67, 0x2a, 0xf0, 0xcd, 0x09, 0x2e, 0xa7,
	0x58, 0x7a, 0x77, 0x08, 0xa1, 0x66, 0x7a, 0xf5, 0x2d, 0xc8, 0xc6, 0xf7, 0xb8, 0xdc, 0x21, 0xd6,
	0xf7, 0x69, 0xe0, 0x45, 0x8c, 0x98, 0x29, 0x6e, 0xe1, 0xdb, 0x24, 0xdc, 0xf3, 0x5d, 0x66, 0x1a,
	0xd2, 0x1d, 0x42, 0x87, 0xf8, 0x0c, 0xb7, 0x89, 0x39, 0xb1, 0x7a, 0x03, 0x16, 0x46, 0xd5, 0x68,
	0x68, 0x19, 0x4a, 0xdb, 0x38, 0x64, 0xae, 0xe3, 0x76, 0x31, 0x23, 0xf5, 0xe0, 0x6e, 0xc0, 0x1a,
	0xbe, 0xcb, 0x5c, 0xcc, 0xa7, 0xac, 0xfd, 0x2c, 0x0b, 0xd9, 0x4d, 0xf5, 0xd3, 0x18, 0xda, 0x03,
	0xe8, 0xfd, 0x58, 0x80, 0x5e, 0xec, 0x45, 0x86, 0xa1, 0x3f, 0x60, 0xca, 0x5f, 0x1f, 0x35, 0x38,
	0xfc, 0x38, 0x6a, 0xa5, 0xae, 0x1a, 0xe8, 0x03, 0xc8, 0x6b, 0x2f, 0xa8, 0x68, 0xb9, 0x47, 0x3a,
	0xfc, 0x0e, 0x5e, 0x7e, 0x6d, 0xe4, 0xe8, 0x09, 0x33, 0x63, 0x40, 0xc3, 0x2f, 0x93, 0xe8, 0x15,
	0xed, 0xb6, 0xf5, 0xa4, 0x77, 0xd2, 0xf2, 0xd7, 0x9e, 0x0e, 0x92, 0x31, 0xd5, 0x4a, 0xa1, 0x4d,
	0x1e, 0x18, 0x92, 0x07, 0xb2, 0xcb, 0xc3, 0x54, 0xda, 0x43, 0x61, 0x79, 0x71, 0xf4, 0xb0, 0x95,
	0x42, 0xef, 0xc2, 0xf4, 0x3d, 0x79, 0xd7, 0xb8, 0x3c, 0xf2, 0x6a, 0x37, 0x9e, 0x61, 0xd4, 0x6b,
	0xb8, 0x95, 0x42, 0xeb, 0x90, 0x4b, 0x6e, 0x16, 0x51, 0x79, 0xf8, 0x05, 0x9c, 0x8e, 0xd8, 0x81,
	0x7e, 0x5d, 0x67, 0xa5, 0xd0, 0x2d, 0xc8, 0xc6, 0x37, 0x5c, 0xe8, 0x85, 0xa1, 0xe7, 0xf0, 0x64,
	0x82, 0xd1, 0x2f, 0xe5, 0x56, 0x0a, 0xed, 0x41, 0xb1, 0xbf, 0x8c, 0x46, 0x2f, 0x0d, 0x57, 0x85,
	0x7d, 0xd5, 0x7e, 0x79, 0xe5, 0x64, 0x40, 0x22, 0xdf, 0x7b, 0x50, 0xd0, 0xab, 0x47, 0x5d, 0xc2,
	0x23, 0x4a, 0xdc, 0x72, 0xe5, 0xa4, 0x61, 0x7d, 0x42, 0x3d, 0xe7, 0xd7, 0x27, 0x1c, 0x51, 0xcf,
	0xe8, 0x13, 0x8e, 0x2a, 0x15, 0xac, 0x14, 0xfa, 0xae, 0x7a, 0xca, 0xd0, 0x72, 0x58, 0xf4, 0xf2,
	0x00, 0xd5, 0x70, 0xa6, 0x5e, 0xb6, 0x9e, 0x06, 0x49, 0x26, 0xdf, 0x52, 0x07, 0x91, 0x8c, 0x51,
	0x43, 0xc6, 0xd1, 0x97, 0xb1, 0x94, 0x2f, 0x9f, 0x30, 0x9a, 0xcc, 0xf6, 0x01, 0xcc, 0x0e, 0x64,
	0x16, 0x68, 0x65, 0x80, 0x66, 0x28, 0x0f, 0x2a, 0xbf, 0xfc, 0x14, 0x44, 0x3c, 0xf3, 0xed, 0x1b,
	0x8f, 0x3f, 0xad, 0xa4, 0x3e, 0xf9, 0xb4, 0x92, 0xfa, 0xfc, 0xd3, 0x8a, 0xf1, 0x83, 0xe3, 0x8a,
	0xf1, 0xcb, 0xe3, 0x8a, 0xf1, 0x97, 0xe3, 0x8a, 0xf1, 0xf8, 0xb8, 0x62, 0xfc, 0xfb, 0xb8, 0x62,
	0x3c, 0x39, 0xae, 0xa4, 0x3e, 0x3f, 0xae, 0x18, 0x1f, 0x7f, 0x56, 0x49, 0x3d, 0xfe, 0xac, 0x92,
	0xfa, 0xe4, 0xb3, 0x4a, 0x6a, 0x7f, 0x5a, 0xe4, 0xad, 0xd7, 0xff, 0x17, 0x00, 0x00, 0xff, 0xff,
	0xe4, 0x65, 0xf5, 0xef, 0x98, 0x2a, 0x00, 0x00,
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
func (x FeeType) String() string {
	s, ok := FeeType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (x FeeBasis) String() string {
	s, ok := FeeBasis_name[int32(x)]
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
func (this *HistoricalLiquidationsRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HistoricalLiquidationsRequest)
	if !ok {
		that2, ok := that.(HistoricalLiquidationsRequest)
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
	if !this.From.Equal(that1.From) {
		return false
	}
	if !this.To.Equal(that1.To) {
		return false
	}
	return true
}
func (this *HistoricalLiquidationsResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HistoricalLiquidationsResponse)
	if !ok {
		that2, ok := that.(HistoricalLiquidationsResponse)
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
	if len(this.Liquidations) != len(that1.Liquidations) {
		return false
	}
	for i := range this.Liquidations {
		if !this.Liquidations[i].Equal(that1.Liquidations[i]) {
			return false
		}
	}
	if this.Success != that1.Success {
		return false
	}
	if this.RejectionReason != that1.RejectionReason {
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
func (this *MarketDataResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MarketDataResponse)
	if !ok {
		that2, ok := that.(MarketDataResponse)
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
	if !this.SnapshotL1.Equal(that1.SnapshotL1) {
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
	if this.SeqNum != that1.SeqNum {
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
	if this.SeqNum != that1.SeqNum {
		return false
	}
	if !this.UpdateL1.Equal(that1.UpdateL1) {
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
	if !this.Liquidation.Equal(that1.Liquidation) {
		return false
	}
	if !this.Funding.Equal(that1.Funding) {
		return false
	}
	if len(this.Stats) != len(that1.Stats) {
		return false
	}
	for i := range this.Stats {
		if !this.Stats[i].Equal(that1.Stats[i]) {
			return false
		}
	}
	return true
}
func (this *AccountDataRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AccountDataRequest)
	if !ok {
		that2, ok := that.(AccountDataRequest)
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
	return true
}
func (this *AccountDataResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AccountDataResponse)
	if !ok {
		that2, ok := that.(AccountDataResponse)
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
	if !this.Orders.Equal(that1.Orders) {
		return false
	}
	if !this.Positions.Equal(that1.Positions) {
		return false
	}
	if !this.Balances.Equal(that1.Balances) {
		return false
	}
	if this.SeqNum != that1.SeqNum {
		return false
	}
	return true
}
func (this *AccountDataIncrementalRefresh) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AccountDataIncrementalRefresh)
	if !ok {
		that2, ok := that.(AccountDataIncrementalRefresh)
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
	if !this.Report.Equal(that1.Report) {
		return false
	}
	return true
}
func (this *SecurityDefinitionRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SecurityDefinitionRequest)
	if !ok {
		that2, ok := that.(SecurityDefinitionRequest)
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
	return true
}
func (this *SecurityDefinitionResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SecurityDefinitionResponse)
	if !ok {
		that2, ok := that.(SecurityDefinitionResponse)
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
	if !this.Security.Equal(that1.Security) {
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
	if len(this.Securities) != len(that1.Securities) {
		return false
	}
	for i := range this.Securities {
		if !this.Securities[i].Equal(that1.Securities[i]) {
			return false
		}
	}
	if this.Success != that1.Success {
		return false
	}
	if this.RejectionReason != that1.RejectionReason {
		return false
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
	if !this.FeeAmount.Equal(that1.FeeAmount) {
		return false
	}
	if !this.FeeCurrency.Equal(that1.FeeCurrency) {
		return false
	}
	if this.FeeType != that1.FeeType {
		return false
	}
	if this.FeeBasis != that1.FeeBasis {
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
	if this.ResponseID != that1.ResponseID {
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
	if this.Success != that1.Success {
		return false
	}
	if this.RejectionReason != that1.RejectionReason {
		return false
	}
	return true
}
func (this *BalancesRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*BalancesRequest)
	if !ok {
		that2, ok := that.(BalancesRequest)
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
	if !this.Asset.Equal(that1.Asset) {
		return false
	}
	if !this.Account.Equal(that1.Account) {
		return false
	}
	return true
}
func (this *BalanceList) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*BalanceList)
	if !ok {
		that2, ok := that.(BalanceList)
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
	if len(this.Balances) != len(that1.Balances) {
		return false
	}
	for i := range this.Balances {
		if !this.Balances[i].Equal(that1.Balances[i]) {
			return false
		}
	}
	if this.Success != that1.Success {
		return false
	}
	if this.RejectionReason != that1.RejectionReason {
		return false
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
	if len(this.ExecutionInstructions) != len(that1.ExecutionInstructions) {
		return false
	}
	for i := range this.ExecutionInstructions {
		if this.ExecutionInstructions[i] != that1.ExecutionInstructions[i] {
			return false
		}
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
	if this.ResponseID != that1.ResponseID {
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
	if this.ResponseID != that1.ResponseID {
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
	if this.Success != that1.Success {
		return false
	}
	if this.RejectionReason != that1.RejectionReason {
		return false
	}
	return true
}
func (this *OrderUpdate) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OrderUpdate)
	if !ok {
		that2, ok := that.(OrderUpdate)
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
	if !this.OrigClientOrderID.Equal(that1.OrigClientOrderID) {
		return false
	}
	if !this.Quantity.Equal(that1.Quantity) {
		return false
	}
	if !this.Price.Equal(that1.Price) {
		return false
	}
	return true
}
func (this *OrderReplaceRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OrderReplaceRequest)
	if !ok {
		that2, ok := that.(OrderReplaceRequest)
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
	if !this.Update.Equal(that1.Update) {
		return false
	}
	return true
}
func (this *OrderReplaceResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OrderReplaceResponse)
	if !ok {
		that2, ok := that.(OrderReplaceResponse)
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
	if this.Success != that1.Success {
		return false
	}
	if this.RejectionReason != that1.RejectionReason {
		return false
	}
	return true
}
func (this *OrderBulkReplaceRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OrderBulkReplaceRequest)
	if !ok {
		that2, ok := that.(OrderBulkReplaceRequest)
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
	if len(this.Updates) != len(that1.Updates) {
		return false
	}
	for i := range this.Updates {
		if !this.Updates[i].Equal(that1.Updates[i]) {
			return false
		}
	}
	return true
}
func (this *OrderBulkReplaceResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OrderBulkReplaceResponse)
	if !ok {
		that2, ok := that.(OrderBulkReplaceResponse)
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
	if this.Success != that1.Success {
		return false
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
	if this.ResponseID != that1.ResponseID {
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
	if !this.Account.Equal(that1.Account) {
		return false
	}
	if !this.Filter.Equal(that1.Filter) {
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
	if this.ResponseID != that1.ResponseID {
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
func (this *HistoricalLiquidationsRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.HistoricalLiquidationsRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	if this.From != nil {
		s = append(s, "From: "+fmt.Sprintf("%#v", this.From)+",\n")
	}
	if this.To != nil {
		s = append(s, "To: "+fmt.Sprintf("%#v", this.To)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *HistoricalLiquidationsResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.HistoricalLiquidationsResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Liquidations != nil {
		s = append(s, "Liquidations: "+fmt.Sprintf("%#v", this.Liquidations)+",\n")
	}
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
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
func (this *MarketDataResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 13)
	s = append(s, "&messages.MarketDataResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.SnapshotL1 != nil {
		s = append(s, "SnapshotL1: "+fmt.Sprintf("%#v", this.SnapshotL1)+",\n")
	}
	if this.SnapshotL2 != nil {
		s = append(s, "SnapshotL2: "+fmt.Sprintf("%#v", this.SnapshotL2)+",\n")
	}
	if this.SnapshotL3 != nil {
		s = append(s, "SnapshotL3: "+fmt.Sprintf("%#v", this.SnapshotL3)+",\n")
	}
	if this.Trades != nil {
		s = append(s, "Trades: "+fmt.Sprintf("%#v", this.Trades)+",\n")
	}
	s = append(s, "SeqNum: "+fmt.Sprintf("%#v", this.SeqNum)+",\n")
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *MarketDataIncrementalRefresh) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 14)
	s = append(s, "&messages.MarketDataIncrementalRefresh{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	s = append(s, "SeqNum: "+fmt.Sprintf("%#v", this.SeqNum)+",\n")
	if this.UpdateL1 != nil {
		s = append(s, "UpdateL1: "+fmt.Sprintf("%#v", this.UpdateL1)+",\n")
	}
	if this.UpdateL2 != nil {
		s = append(s, "UpdateL2: "+fmt.Sprintf("%#v", this.UpdateL2)+",\n")
	}
	if this.UpdateL3 != nil {
		s = append(s, "UpdateL3: "+fmt.Sprintf("%#v", this.UpdateL3)+",\n")
	}
	if this.Trades != nil {
		s = append(s, "Trades: "+fmt.Sprintf("%#v", this.Trades)+",\n")
	}
	if this.Liquidation != nil {
		s = append(s, "Liquidation: "+fmt.Sprintf("%#v", this.Liquidation)+",\n")
	}
	if this.Funding != nil {
		s = append(s, "Funding: "+fmt.Sprintf("%#v", this.Funding)+",\n")
	}
	if this.Stats != nil {
		s = append(s, "Stats: "+fmt.Sprintf("%#v", this.Stats)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *AccountDataRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.AccountDataRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Subscribe: "+fmt.Sprintf("%#v", this.Subscribe)+",\n")
	if this.Subscriber != nil {
		s = append(s, "Subscriber: "+fmt.Sprintf("%#v", this.Subscriber)+",\n")
	}
	if this.Account != nil {
		s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *AccountDataResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 10)
	s = append(s, "&messages.AccountDataResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Orders != nil {
		s = append(s, "Orders: "+fmt.Sprintf("%#v", this.Orders)+",\n")
	}
	if this.Positions != nil {
		s = append(s, "Positions: "+fmt.Sprintf("%#v", this.Positions)+",\n")
	}
	if this.Balances != nil {
		s = append(s, "Balances: "+fmt.Sprintf("%#v", this.Balances)+",\n")
	}
	s = append(s, "SeqNum: "+fmt.Sprintf("%#v", this.SeqNum)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *AccountDataIncrementalRefresh) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.AccountDataIncrementalRefresh{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Report != nil {
		s = append(s, "Report: "+fmt.Sprintf("%#v", this.Report)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *SecurityDefinitionRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&messages.SecurityDefinitionRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *SecurityDefinitionResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.SecurityDefinitionResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Security != nil {
		s = append(s, "Security: "+fmt.Sprintf("%#v", this.Security)+",\n")
	}
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
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
	s := make([]string, 0, 9)
	s = append(s, "&messages.SecurityList{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Securities != nil {
		s = append(s, "Securities: "+fmt.Sprintf("%#v", this.Securities)+",\n")
	}
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ExecutionReport) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 22)
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
	if this.FeeAmount != nil {
		s = append(s, "FeeAmount: "+fmt.Sprintf("%#v", this.FeeAmount)+",\n")
	}
	if this.FeeCurrency != nil {
		s = append(s, "FeeCurrency: "+fmt.Sprintf("%#v", this.FeeCurrency)+",\n")
	}
	s = append(s, "FeeType: "+fmt.Sprintf("%#v", this.FeeType)+",\n")
	s = append(s, "FeeBasis: "+fmt.Sprintf("%#v", this.FeeBasis)+",\n")
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
	s := make([]string, 0, 9)
	s = append(s, "&messages.PositionList{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Positions != nil {
		s = append(s, "Positions: "+fmt.Sprintf("%#v", this.Positions)+",\n")
	}
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *BalancesRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.BalancesRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Subscribe: "+fmt.Sprintf("%#v", this.Subscribe)+",\n")
	if this.Subscriber != nil {
		s = append(s, "Subscriber: "+fmt.Sprintf("%#v", this.Subscriber)+",\n")
	}
	if this.Asset != nil {
		s = append(s, "Asset: "+fmt.Sprintf("%#v", this.Asset)+",\n")
	}
	if this.Account != nil {
		s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *BalanceList) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.BalanceList{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Balances != nil {
		s = append(s, "Balances: "+fmt.Sprintf("%#v", this.Balances)+",\n")
	}
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *NewOrder) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 12)
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
	s = append(s, "ExecutionInstructions: "+fmt.Sprintf("%#v", this.ExecutionInstructions)+",\n")
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
	s := make([]string, 0, 9)
	s = append(s, "&messages.NewOrderSingleResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
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
	s := make([]string, 0, 9)
	s = append(s, "&messages.NewOrderBulkResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	s = append(s, "OrderIDs: "+fmt.Sprintf("%#v", this.OrderIDs)+",\n")
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OrderUpdate) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.OrderUpdate{")
	if this.OrderID != nil {
		s = append(s, "OrderID: "+fmt.Sprintf("%#v", this.OrderID)+",\n")
	}
	if this.OrigClientOrderID != nil {
		s = append(s, "OrigClientOrderID: "+fmt.Sprintf("%#v", this.OrigClientOrderID)+",\n")
	}
	if this.Quantity != nil {
		s = append(s, "Quantity: "+fmt.Sprintf("%#v", this.Quantity)+",\n")
	}
	if this.Price != nil {
		s = append(s, "Price: "+fmt.Sprintf("%#v", this.Price)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OrderReplaceRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.OrderReplaceRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	if this.Account != nil {
		s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
	}
	if this.Update != nil {
		s = append(s, "Update: "+fmt.Sprintf("%#v", this.Update)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OrderReplaceResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.OrderReplaceResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OrderBulkReplaceRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.OrderBulkReplaceRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	if this.Account != nil {
		s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
	}
	if this.Updates != nil {
		s = append(s, "Updates: "+fmt.Sprintf("%#v", this.Updates)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OrderBulkReplaceResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.OrderBulkReplaceResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
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
	s := make([]string, 0, 8)
	s = append(s, "&messages.OrderCancelResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
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
	if this.Account != nil {
		s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
	}
	if this.Filter != nil {
		s = append(s, "Filter: "+fmt.Sprintf("%#v", this.Filter)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *OrderMassCancelResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.OrderMassCancelResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
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

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ExecutorClient is the client API for Executor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ExecutorClient interface {
	MarketData(ctx context.Context, in *MarketDataRequest, opts ...grpc.CallOption) (Executor_MarketDataClient, error)
	AccountData(ctx context.Context, in *AccountDataRequest, opts ...grpc.CallOption) (Executor_AccountDataClient, error)
	SecurityDefinition(ctx context.Context, in *SecurityDefinitionRequest, opts ...grpc.CallOption) (*SecurityDefinitionResponse, error)
	Securities(ctx context.Context, in *SecurityListRequest, opts ...grpc.CallOption) (*SecurityList, error)
	Orders(ctx context.Context, in *OrderStatusRequest, opts ...grpc.CallOption) (*OrderList, error)
	Positions(ctx context.Context, in *PositionsRequest, opts ...grpc.CallOption) (*PositionList, error)
	Balances(ctx context.Context, in *BalancesRequest, opts ...grpc.CallOption) (*BalanceList, error)
	NewOrderSingle(ctx context.Context, in *NewOrderSingleRequest, opts ...grpc.CallOption) (*NewOrderSingleResponse, error)
	NewOrderBulk(ctx context.Context, in *NewOrderBulkRequest, opts ...grpc.CallOption) (*NewOrderBulkResponse, error)
	OrderReplace(ctx context.Context, in *OrderReplaceRequest, opts ...grpc.CallOption) (*OrderReplaceResponse, error)
	OrderBulkReplace(ctx context.Context, in *OrderBulkReplaceRequest, opts ...grpc.CallOption) (*OrderBulkReplaceResponse, error)
	OrderCancel(ctx context.Context, in *OrderCancelRequest, opts ...grpc.CallOption) (*OrderCancelResponse, error)
	OrderMassCancel(ctx context.Context, in *OrderMassCancelRequest, opts ...grpc.CallOption) (*OrderMassCancelResponse, error)
}

type executorClient struct {
	cc *grpc.ClientConn
}

func NewExecutorClient(cc *grpc.ClientConn) ExecutorClient {
	return &executorClient{cc}
}

func (c *executorClient) MarketData(ctx context.Context, in *MarketDataRequest, opts ...grpc.CallOption) (Executor_MarketDataClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Executor_serviceDesc.Streams[0], "/messages.Executor/MarketData", opts...)
	if err != nil {
		return nil, err
	}
	x := &executorMarketDataClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Executor_MarketDataClient interface {
	Recv() (*MarketDataIncrementalRefresh, error)
	grpc.ClientStream
}

type executorMarketDataClient struct {
	grpc.ClientStream
}

func (x *executorMarketDataClient) Recv() (*MarketDataIncrementalRefresh, error) {
	m := new(MarketDataIncrementalRefresh)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *executorClient) AccountData(ctx context.Context, in *AccountDataRequest, opts ...grpc.CallOption) (Executor_AccountDataClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Executor_serviceDesc.Streams[1], "/messages.Executor/AccountData", opts...)
	if err != nil {
		return nil, err
	}
	x := &executorAccountDataClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Executor_AccountDataClient interface {
	Recv() (*AccountDataIncrementalRefresh, error)
	grpc.ClientStream
}

type executorAccountDataClient struct {
	grpc.ClientStream
}

func (x *executorAccountDataClient) Recv() (*AccountDataIncrementalRefresh, error) {
	m := new(AccountDataIncrementalRefresh)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *executorClient) SecurityDefinition(ctx context.Context, in *SecurityDefinitionRequest, opts ...grpc.CallOption) (*SecurityDefinitionResponse, error) {
	out := new(SecurityDefinitionResponse)
	err := c.cc.Invoke(ctx, "/messages.Executor/SecurityDefinition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) Securities(ctx context.Context, in *SecurityListRequest, opts ...grpc.CallOption) (*SecurityList, error) {
	out := new(SecurityList)
	err := c.cc.Invoke(ctx, "/messages.Executor/Securities", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) Orders(ctx context.Context, in *OrderStatusRequest, opts ...grpc.CallOption) (*OrderList, error) {
	out := new(OrderList)
	err := c.cc.Invoke(ctx, "/messages.Executor/Orders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) Positions(ctx context.Context, in *PositionsRequest, opts ...grpc.CallOption) (*PositionList, error) {
	out := new(PositionList)
	err := c.cc.Invoke(ctx, "/messages.Executor/Positions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) Balances(ctx context.Context, in *BalancesRequest, opts ...grpc.CallOption) (*BalanceList, error) {
	out := new(BalanceList)
	err := c.cc.Invoke(ctx, "/messages.Executor/Balances", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) NewOrderSingle(ctx context.Context, in *NewOrderSingleRequest, opts ...grpc.CallOption) (*NewOrderSingleResponse, error) {
	out := new(NewOrderSingleResponse)
	err := c.cc.Invoke(ctx, "/messages.Executor/NewOrderSingle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) NewOrderBulk(ctx context.Context, in *NewOrderBulkRequest, opts ...grpc.CallOption) (*NewOrderBulkResponse, error) {
	out := new(NewOrderBulkResponse)
	err := c.cc.Invoke(ctx, "/messages.Executor/NewOrderBulk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) OrderReplace(ctx context.Context, in *OrderReplaceRequest, opts ...grpc.CallOption) (*OrderReplaceResponse, error) {
	out := new(OrderReplaceResponse)
	err := c.cc.Invoke(ctx, "/messages.Executor/OrderReplace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) OrderBulkReplace(ctx context.Context, in *OrderBulkReplaceRequest, opts ...grpc.CallOption) (*OrderBulkReplaceResponse, error) {
	out := new(OrderBulkReplaceResponse)
	err := c.cc.Invoke(ctx, "/messages.Executor/OrderBulkReplace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) OrderCancel(ctx context.Context, in *OrderCancelRequest, opts ...grpc.CallOption) (*OrderCancelResponse, error) {
	out := new(OrderCancelResponse)
	err := c.cc.Invoke(ctx, "/messages.Executor/OrderCancel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) OrderMassCancel(ctx context.Context, in *OrderMassCancelRequest, opts ...grpc.CallOption) (*OrderMassCancelResponse, error) {
	out := new(OrderMassCancelResponse)
	err := c.cc.Invoke(ctx, "/messages.Executor/OrderMassCancel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExecutorServer is the server API for Executor service.
type ExecutorServer interface {
	MarketData(*MarketDataRequest, Executor_MarketDataServer) error
	AccountData(*AccountDataRequest, Executor_AccountDataServer) error
	SecurityDefinition(context.Context, *SecurityDefinitionRequest) (*SecurityDefinitionResponse, error)
	Securities(context.Context, *SecurityListRequest) (*SecurityList, error)
	Orders(context.Context, *OrderStatusRequest) (*OrderList, error)
	Positions(context.Context, *PositionsRequest) (*PositionList, error)
	Balances(context.Context, *BalancesRequest) (*BalanceList, error)
	NewOrderSingle(context.Context, *NewOrderSingleRequest) (*NewOrderSingleResponse, error)
	NewOrderBulk(context.Context, *NewOrderBulkRequest) (*NewOrderBulkResponse, error)
	OrderReplace(context.Context, *OrderReplaceRequest) (*OrderReplaceResponse, error)
	OrderBulkReplace(context.Context, *OrderBulkReplaceRequest) (*OrderBulkReplaceResponse, error)
	OrderCancel(context.Context, *OrderCancelRequest) (*OrderCancelResponse, error)
	OrderMassCancel(context.Context, *OrderMassCancelRequest) (*OrderMassCancelResponse, error)
}

func RegisterExecutorServer(s *grpc.Server, srv ExecutorServer) {
	s.RegisterService(&_Executor_serviceDesc, srv)
}

func _Executor_MarketData_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(MarketDataRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExecutorServer).MarketData(m, &executorMarketDataServer{stream})
}

type Executor_MarketDataServer interface {
	Send(*MarketDataIncrementalRefresh) error
	grpc.ServerStream
}

type executorMarketDataServer struct {
	grpc.ServerStream
}

func (x *executorMarketDataServer) Send(m *MarketDataIncrementalRefresh) error {
	return x.ServerStream.SendMsg(m)
}

func _Executor_AccountData_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(AccountDataRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExecutorServer).AccountData(m, &executorAccountDataServer{stream})
}

type Executor_AccountDataServer interface {
	Send(*AccountDataIncrementalRefresh) error
	grpc.ServerStream
}

type executorAccountDataServer struct {
	grpc.ServerStream
}

func (x *executorAccountDataServer) Send(m *AccountDataIncrementalRefresh) error {
	return x.ServerStream.SendMsg(m)
}

func _Executor_SecurityDefinition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecurityDefinitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).SecurityDefinition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.Executor/SecurityDefinition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).SecurityDefinition(ctx, req.(*SecurityDefinitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_Securities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecurityListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).Securities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.Executor/Securities",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).Securities(ctx, req.(*SecurityListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_Orders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).Orders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.Executor/Orders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).Orders(ctx, req.(*OrderStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_Positions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PositionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).Positions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.Executor/Positions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).Positions(ctx, req.(*PositionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_Balances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).Balances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.Executor/Balances",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).Balances(ctx, req.(*BalancesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_NewOrderSingle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewOrderSingleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).NewOrderSingle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.Executor/NewOrderSingle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).NewOrderSingle(ctx, req.(*NewOrderSingleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_NewOrderBulk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewOrderBulkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).NewOrderBulk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.Executor/NewOrderBulk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).NewOrderBulk(ctx, req.(*NewOrderBulkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_OrderReplace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderReplaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).OrderReplace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.Executor/OrderReplace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).OrderReplace(ctx, req.(*OrderReplaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_OrderBulkReplace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderBulkReplaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).OrderBulkReplace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.Executor/OrderBulkReplace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).OrderBulkReplace(ctx, req.(*OrderBulkReplaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_OrderCancel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderCancelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).OrderCancel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.Executor/OrderCancel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).OrderCancel(ctx, req.(*OrderCancelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_OrderMassCancel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderMassCancelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).OrderMassCancel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.Executor/OrderMassCancel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).OrderMassCancel(ctx, req.(*OrderMassCancelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Executor_serviceDesc = grpc.ServiceDesc{
	ServiceName: "messages.Executor",
	HandlerType: (*ExecutorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SecurityDefinition",
			Handler:    _Executor_SecurityDefinition_Handler,
		},
		{
			MethodName: "Securities",
			Handler:    _Executor_Securities_Handler,
		},
		{
			MethodName: "Orders",
			Handler:    _Executor_Orders_Handler,
		},
		{
			MethodName: "Positions",
			Handler:    _Executor_Positions_Handler,
		},
		{
			MethodName: "Balances",
			Handler:    _Executor_Balances_Handler,
		},
		{
			MethodName: "NewOrderSingle",
			Handler:    _Executor_NewOrderSingle_Handler,
		},
		{
			MethodName: "NewOrderBulk",
			Handler:    _Executor_NewOrderBulk_Handler,
		},
		{
			MethodName: "OrderReplace",
			Handler:    _Executor_OrderReplace_Handler,
		},
		{
			MethodName: "OrderBulkReplace",
			Handler:    _Executor_OrderBulkReplace_Handler,
		},
		{
			MethodName: "OrderCancel",
			Handler:    _Executor_OrderCancel_Handler,
		},
		{
			MethodName: "OrderMassCancel",
			Handler:    _Executor_OrderMassCancel_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MarketData",
			Handler:       _Executor_MarketData_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "AccountData",
			Handler:       _Executor_AccountData_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "executor_messages.proto",
}

func (m *HistoricalLiquidationsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HistoricalLiquidationsRequest) MarshalTo(dAtA []byte) (int, error) {
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
		n1, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.From != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.From.Size()))
		n2, err := m.From.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if m.To != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.To.Size()))
		n3, err := m.To.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	return i, nil
}

func (m *HistoricalLiquidationsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HistoricalLiquidationsResponse) MarshalTo(dAtA []byte) (int, error) {
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
	if len(m.Liquidations) > 0 {
		for _, msg := range m.Liquidations {
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
	if m.Success {
		dAtA[i] = 0x20
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
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
		n4, err := m.Subscriber.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n4
	}
	if m.Instrument != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Instrument.Size()))
		n5, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n5
	}
	if m.Aggregation != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Aggregation))
	}
	return i, nil
}

func (m *MarketDataResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MarketDataResponse) MarshalTo(dAtA []byte) (int, error) {
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
	if m.SnapshotL1 != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SnapshotL1.Size()))
		n6, err := m.SnapshotL1.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n6
	}
	if m.SnapshotL2 != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SnapshotL2.Size()))
		n7, err := m.SnapshotL2.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n7
	}
	if m.SnapshotL3 != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SnapshotL3.Size()))
		n8, err := m.SnapshotL3.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n8
	}
	if len(m.Trades) > 0 {
		for _, msg := range m.Trades {
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
	if m.SeqNum != 0 {
		dAtA[i] = 0x38
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SeqNum))
	}
	if m.Success {
		dAtA[i] = 0x40
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x48
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
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
	if m.SeqNum != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SeqNum))
	}
	if m.UpdateL1 != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.UpdateL1.Size()))
		n9, err := m.UpdateL1.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n9
	}
	if m.UpdateL2 != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.UpdateL2.Size()))
		n10, err := m.UpdateL2.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n10
	}
	if m.UpdateL3 != nil {
		dAtA[i] = 0x32
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.UpdateL3.Size()))
		n11, err := m.UpdateL3.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n11
	}
	if len(m.Trades) > 0 {
		for _, msg := range m.Trades {
			dAtA[i] = 0x3a
			i++
			i = encodeVarintExecutorMessages(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Liquidation != nil {
		dAtA[i] = 0x42
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Liquidation.Size()))
		n12, err := m.Liquidation.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n12
	}
	if m.Funding != nil {
		dAtA[i] = 0x4a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Funding.Size()))
		n13, err := m.Funding.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n13
	}
	if len(m.Stats) > 0 {
		for _, msg := range m.Stats {
			dAtA[i] = 0x52
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

func (m *AccountDataRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountDataRequest) MarshalTo(dAtA []byte) (int, error) {
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
		n14, err := m.Subscriber.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n14
	}
	if m.Account != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n15, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n15
	}
	return i, nil
}

func (m *AccountDataResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountDataResponse) MarshalTo(dAtA []byte) (int, error) {
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
	if m.Orders != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Orders.Size()))
		n16, err := m.Orders.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n16
	}
	if m.Positions != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Positions.Size()))
		n17, err := m.Positions.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n17
	}
	if m.Balances != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Balances.Size()))
		n18, err := m.Balances.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n18
	}
	if m.SeqNum != 0 {
		dAtA[i] = 0x30
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SeqNum))
	}
	return i, nil
}

func (m *AccountDataIncrementalRefresh) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountDataIncrementalRefresh) MarshalTo(dAtA []byte) (int, error) {
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
	if m.Report != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Report.Size()))
		n19, err := m.Report.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n19
	}
	return i, nil
}

func (m *SecurityDefinitionRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SecurityDefinitionRequest) MarshalTo(dAtA []byte) (int, error) {
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
		n20, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n20
	}
	return i, nil
}

func (m *SecurityDefinitionResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SecurityDefinitionResponse) MarshalTo(dAtA []byte) (int, error) {
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
	if m.Security != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Security.Size()))
		n21, err := m.Security.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n21
	}
	if m.Success {
		dAtA[i] = 0x20
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
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
		n22, err := m.Subscriber.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n22
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
	if len(m.Securities) > 0 {
		for _, msg := range m.Securities {
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
	if m.Success {
		dAtA[i] = 0x20
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
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
		n23, err := m.ClientOrderID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n23
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
		n24, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n24
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
		n25, err := m.TransactionTime.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n25
	}
	if m.TradeID != nil {
		dAtA[i] = 0x62
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.TradeID.Size()))
		n26, err := m.TradeID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n26
	}
	if m.FillPrice != nil {
		dAtA[i] = 0x6a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.FillPrice.Size()))
		n27, err := m.FillPrice.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n27
	}
	if m.FillQuantity != nil {
		dAtA[i] = 0x72
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.FillQuantity.Size()))
		n28, err := m.FillQuantity.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n28
	}
	if m.FeeAmount != nil {
		dAtA[i] = 0x7a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.FeeAmount.Size()))
		n29, err := m.FeeAmount.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n29
	}
	if m.FeeCurrency != nil {
		dAtA[i] = 0x82
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.FeeCurrency.Size()))
		n30, err := m.FeeCurrency.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n30
	}
	if m.FeeType != 0 {
		dAtA[i] = 0x88
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.FeeType))
	}
	if m.FeeBasis != 0 {
		dAtA[i] = 0x90
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.FeeBasis))
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x98
		i++
		dAtA[i] = 0x1
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
		n31, err := m.OrderID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n31
	}
	if m.ClientOrderID != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ClientOrderID.Size()))
		n32, err := m.ClientOrderID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n32
	}
	if m.Instrument != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Instrument.Size()))
		n33, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n33
	}
	if m.Side != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Side.Size()))
		n34, err := m.Side.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n34
	}
	if m.OrderStatus != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderStatus.Size()))
		n35, err := m.OrderStatus.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n35
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
		n36, err := m.Subscriber.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n36
	}
	if m.Account != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n37, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n37
	}
	if m.Filter != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Filter.Size()))
		n38, err := m.Filter.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n38
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
		n39, err := m.Subscriber.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n39
	}
	if m.Instrument != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Instrument.Size()))
		n40, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n40
	}
	if m.Account != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n41, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n41
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
	if m.ResponseID != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
	}
	if len(m.Positions) > 0 {
		for _, msg := range m.Positions {
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
	if m.Success {
		dAtA[i] = 0x20
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
	}
	return i, nil
}

func (m *BalancesRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BalancesRequest) MarshalTo(dAtA []byte) (int, error) {
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
		n42, err := m.Subscriber.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n42
	}
	if m.Asset != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Asset.Size()))
		n43, err := m.Asset.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n43
	}
	if m.Account != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n44, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n44
	}
	return i, nil
}

func (m *BalanceList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BalanceList) MarshalTo(dAtA []byte) (int, error) {
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
	if len(m.Balances) > 0 {
		for _, msg := range m.Balances {
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
	if m.Success {
		dAtA[i] = 0x20
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
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
		n45, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n45
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
		n46, err := m.Price.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n46
	}
	if len(m.ExecutionInstructions) > 0 {
		dAtA48 := make([]byte, len(m.ExecutionInstructions)*10)
		var j47 int
		for _, num := range m.ExecutionInstructions {
			for num >= 1<<7 {
				dAtA48[j47] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j47++
			}
			dAtA48[j47] = uint8(num)
			j47++
		}
		dAtA[i] = 0x52
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(j47))
		i += copy(dAtA[i:], dAtA48[:j47])
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
		n49, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n49
	}
	if m.Order != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Order.Size()))
		n50, err := m.Order.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n50
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
	if m.ResponseID != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
	}
	if m.Success {
		dAtA[i] = 0x18
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if len(m.OrderID) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.OrderID)))
		i += copy(dAtA[i:], m.OrderID)
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x28
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
		n51, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n51
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
	if m.ResponseID != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
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
	if m.Success {
		dAtA[i] = 0x20
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
	}
	return i, nil
}

func (m *OrderUpdate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderUpdate) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.OrderID != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderID.Size()))
		n52, err := m.OrderID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n52
	}
	if m.OrigClientOrderID != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrigClientOrderID.Size()))
		n53, err := m.OrigClientOrderID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n53
	}
	if m.Quantity != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Quantity.Size()))
		n54, err := m.Quantity.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n54
	}
	if m.Price != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Price.Size()))
		n55, err := m.Price.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n55
	}
	return i, nil
}

func (m *OrderReplaceRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderReplaceRequest) MarshalTo(dAtA []byte) (int, error) {
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
		n56, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n56
	}
	if m.Account != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n57, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n57
	}
	if m.Update != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Update.Size()))
		n58, err := m.Update.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n58
	}
	return i, nil
}

func (m *OrderReplaceResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderReplaceResponse) MarshalTo(dAtA []byte) (int, error) {
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
	if m.Success {
		dAtA[i] = 0x18
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
	}
	return i, nil
}

func (m *OrderBulkReplaceRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderBulkReplaceRequest) MarshalTo(dAtA []byte) (int, error) {
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
		n59, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n59
	}
	if m.Account != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n60, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n60
	}
	if len(m.Updates) > 0 {
		for _, msg := range m.Updates {
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

func (m *OrderBulkReplaceResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderBulkReplaceResponse) MarshalTo(dAtA []byte) (int, error) {
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
	if m.Success {
		dAtA[i] = 0x18
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
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
		n61, err := m.OrderID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n61
	}
	if m.ClientOrderID != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ClientOrderID.Size()))
		n62, err := m.ClientOrderID.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n62
	}
	if m.Instrument != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Instrument.Size()))
		n63, err := m.Instrument.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n63
	}
	if m.Account != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n64, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n64
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
	if m.ResponseID != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
	}
	if m.Success {
		dAtA[i] = 0x18
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x20
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
	if m.Account != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Account.Size()))
		n65, err := m.Account.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n65
	}
	if m.Filter != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Filter.Size()))
		n66, err := m.Filter.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n66
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
	if m.ResponseID != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
	}
	if m.Success {
		dAtA[i] = 0x18
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.RejectionReason != 0 {
		dAtA[i] = 0x20
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
func (m *HistoricalLiquidationsRequest) Size() (n int) {
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
	if m.From != nil {
		l = m.From.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.To != nil {
		l = m.To.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *HistoricalLiquidationsResponse) Size() (n int) {
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
	if len(m.Liquidations) > 0 {
		for _, e := range m.Liquidations {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	if m.Success {
		n += 2
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
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

func (m *MarketDataResponse) Size() (n int) {
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
	if m.SnapshotL1 != nil {
		l = m.SnapshotL1.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
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
	if m.SeqNum != 0 {
		n += 1 + sovExecutorMessages(uint64(m.SeqNum))
	}
	if m.Success {
		n += 2
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
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
	if m.SeqNum != 0 {
		n += 1 + sovExecutorMessages(uint64(m.SeqNum))
	}
	if m.UpdateL1 != nil {
		l = m.UpdateL1.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
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
	if m.Liquidation != nil {
		l = m.Liquidation.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Funding != nil {
		l = m.Funding.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if len(m.Stats) > 0 {
		for _, e := range m.Stats {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	return n
}

func (m *AccountDataRequest) Size() (n int) {
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
	return n
}

func (m *AccountDataResponse) Size() (n int) {
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
	if m.Orders != nil {
		l = m.Orders.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Positions != nil {
		l = m.Positions.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Balances != nil {
		l = m.Balances.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.SeqNum != 0 {
		n += 1 + sovExecutorMessages(uint64(m.SeqNum))
	}
	return n
}

func (m *AccountDataIncrementalRefresh) Size() (n int) {
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
	if m.Report != nil {
		l = m.Report.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *SecurityDefinitionRequest) Size() (n int) {
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
	return n
}

func (m *SecurityDefinitionResponse) Size() (n int) {
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
	if m.Security != nil {
		l = m.Security.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Success {
		n += 2
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
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
	if len(m.Securities) > 0 {
		for _, e := range m.Securities {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	if m.Success {
		n += 2
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
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
	if m.FeeAmount != nil {
		l = m.FeeAmount.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.FeeCurrency != nil {
		l = m.FeeCurrency.Size()
		n += 2 + l + sovExecutorMessages(uint64(l))
	}
	if m.FeeType != 0 {
		n += 2 + sovExecutorMessages(uint64(m.FeeType))
	}
	if m.FeeBasis != 0 {
		n += 2 + sovExecutorMessages(uint64(m.FeeBasis))
	}
	if m.RejectionReason != 0 {
		n += 2 + sovExecutorMessages(uint64(m.RejectionReason))
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
	if m.ResponseID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.ResponseID))
	}
	if len(m.Positions) > 0 {
		for _, e := range m.Positions {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	if m.Success {
		n += 2
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
	}
	return n
}

func (m *BalancesRequest) Size() (n int) {
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
	if m.Asset != nil {
		l = m.Asset.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Account != nil {
		l = m.Account.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *BalanceList) Size() (n int) {
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
	if len(m.Balances) > 0 {
		for _, e := range m.Balances {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	if m.Success {
		n += 2
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
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
	if len(m.ExecutionInstructions) > 0 {
		l = 0
		for _, e := range m.ExecutionInstructions {
			l += sovExecutorMessages(uint64(e))
		}
		n += 1 + sovExecutorMessages(uint64(l)) + l
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
	if m.ResponseID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.ResponseID))
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
	if m.ResponseID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.ResponseID))
	}
	if len(m.OrderIDs) > 0 {
		for _, s := range m.OrderIDs {
			l = len(s)
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	if m.Success {
		n += 2
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
	}
	return n
}

func (m *OrderUpdate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.OrderID != nil {
		l = m.OrderID.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.OrigClientOrderID != nil {
		l = m.OrigClientOrderID.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Quantity != nil {
		l = m.Quantity.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Price != nil {
		l = m.Price.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *OrderReplaceRequest) Size() (n int) {
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
	if m.Update != nil {
		l = m.Update.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *OrderReplaceResponse) Size() (n int) {
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
	if m.Success {
		n += 2
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
	}
	return n
}

func (m *OrderBulkReplaceRequest) Size() (n int) {
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
	if len(m.Updates) > 0 {
		for _, e := range m.Updates {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	return n
}

func (m *OrderBulkReplaceResponse) Size() (n int) {
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
	if m.Success {
		n += 2
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
	if m.ResponseID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.ResponseID))
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

func (m *OrderMassCancelResponse) Size() (n int) {
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
func (this *HistoricalLiquidationsRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&HistoricalLiquidationsRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`From:` + strings.Replace(fmt.Sprintf("%v", this.From), "Timestamp", "types.Timestamp", 1) + `,`,
		`To:` + strings.Replace(fmt.Sprintf("%v", this.To), "Timestamp", "types.Timestamp", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *HistoricalLiquidationsResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&HistoricalLiquidationsResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Liquidations:` + strings.Replace(fmt.Sprintf("%v", this.Liquidations), "Liquidation", "models.Liquidation", 1) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
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
func (this *MarketDataResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&MarketDataResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`SnapshotL1:` + strings.Replace(fmt.Sprintf("%v", this.SnapshotL1), "OBL1Snapshot", "models.OBL1Snapshot", 1) + `,`,
		`SnapshotL2:` + strings.Replace(fmt.Sprintf("%v", this.SnapshotL2), "OBL2Snapshot", "models.OBL2Snapshot", 1) + `,`,
		`SnapshotL3:` + strings.Replace(fmt.Sprintf("%v", this.SnapshotL3), "OBL3Snapshot", "models.OBL3Snapshot", 1) + `,`,
		`Trades:` + strings.Replace(fmt.Sprintf("%v", this.Trades), "AggregatedTrade", "models.AggregatedTrade", 1) + `,`,
		`SeqNum:` + fmt.Sprintf("%v", this.SeqNum) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
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
		`SeqNum:` + fmt.Sprintf("%v", this.SeqNum) + `,`,
		`UpdateL1:` + strings.Replace(fmt.Sprintf("%v", this.UpdateL1), "OBL1Update", "models.OBL1Update", 1) + `,`,
		`UpdateL2:` + strings.Replace(fmt.Sprintf("%v", this.UpdateL2), "OBL2Update", "models.OBL2Update", 1) + `,`,
		`UpdateL3:` + strings.Replace(fmt.Sprintf("%v", this.UpdateL3), "OBL3Update", "models.OBL3Update", 1) + `,`,
		`Trades:` + strings.Replace(fmt.Sprintf("%v", this.Trades), "AggregatedTrade", "models.AggregatedTrade", 1) + `,`,
		`Liquidation:` + strings.Replace(fmt.Sprintf("%v", this.Liquidation), "Liquidation", "models.Liquidation", 1) + `,`,
		`Funding:` + strings.Replace(fmt.Sprintf("%v", this.Funding), "Funding", "models.Funding", 1) + `,`,
		`Stats:` + strings.Replace(fmt.Sprintf("%v", this.Stats), "Stat", "models.Stat", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AccountDataRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AccountDataRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Subscribe:` + fmt.Sprintf("%v", this.Subscribe) + `,`,
		`Subscriber:` + strings.Replace(fmt.Sprintf("%v", this.Subscriber), "PID", "actor.PID", 1) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AccountDataResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AccountDataResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Orders:` + strings.Replace(fmt.Sprintf("%v", this.Orders), "OrderList", "OrderList", 1) + `,`,
		`Positions:` + strings.Replace(fmt.Sprintf("%v", this.Positions), "PositionList", "PositionList", 1) + `,`,
		`Balances:` + strings.Replace(fmt.Sprintf("%v", this.Balances), "BalanceList", "BalanceList", 1) + `,`,
		`SeqNum:` + fmt.Sprintf("%v", this.SeqNum) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AccountDataIncrementalRefresh) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AccountDataIncrementalRefresh{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Report:` + strings.Replace(fmt.Sprintf("%v", this.Report), "ExecutionReport", "ExecutionReport", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *SecurityDefinitionRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&SecurityDefinitionRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *SecurityDefinitionResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&SecurityDefinitionResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Security:` + strings.Replace(fmt.Sprintf("%v", this.Security), "Security", "models.Security", 1) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
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
		`Securities:` + strings.Replace(fmt.Sprintf("%v", this.Securities), "Security", "models.Security", 1) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
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
		`FeeAmount:` + strings.Replace(fmt.Sprintf("%v", this.FeeAmount), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`FeeCurrency:` + strings.Replace(fmt.Sprintf("%v", this.FeeCurrency), "Asset", "models1.Asset", 1) + `,`,
		`FeeType:` + fmt.Sprintf("%v", this.FeeType) + `,`,
		`FeeBasis:` + fmt.Sprintf("%v", this.FeeBasis) + `,`,
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
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Positions:` + strings.Replace(fmt.Sprintf("%v", this.Positions), "Position", "models.Position", 1) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *BalancesRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&BalancesRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Subscribe:` + fmt.Sprintf("%v", this.Subscribe) + `,`,
		`Subscriber:` + strings.Replace(fmt.Sprintf("%v", this.Subscriber), "PID", "actor.PID", 1) + `,`,
		`Asset:` + strings.Replace(fmt.Sprintf("%v", this.Asset), "Asset", "models1.Asset", 1) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *BalanceList) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&BalanceList{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Balances:` + strings.Replace(fmt.Sprintf("%v", this.Balances), "Balance", "models.Balance", 1) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
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
		`ExecutionInstructions:` + fmt.Sprintf("%v", this.ExecutionInstructions) + `,`,
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
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
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
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`OrderIDs:` + fmt.Sprintf("%v", this.OrderIDs) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrderUpdate) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OrderUpdate{`,
		`OrderID:` + strings.Replace(fmt.Sprintf("%v", this.OrderID), "StringValue", "types.StringValue", 1) + `,`,
		`OrigClientOrderID:` + strings.Replace(fmt.Sprintf("%v", this.OrigClientOrderID), "StringValue", "types.StringValue", 1) + `,`,
		`Quantity:` + strings.Replace(fmt.Sprintf("%v", this.Quantity), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`Price:` + strings.Replace(fmt.Sprintf("%v", this.Price), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrderReplaceRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OrderReplaceRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`Update:` + strings.Replace(fmt.Sprintf("%v", this.Update), "OrderUpdate", "OrderUpdate", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrderReplaceResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OrderReplaceResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrderBulkReplaceRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OrderBulkReplaceRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`Updates:` + strings.Replace(fmt.Sprintf("%v", this.Updates), "OrderUpdate", "OrderUpdate", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrderBulkReplaceResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&OrderBulkReplaceResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
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
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
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
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`Filter:` + strings.Replace(fmt.Sprintf("%v", this.Filter), "OrderFilter", "OrderFilter", 1) + `,`,
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
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
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
func (m *HistoricalLiquidationsRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: HistoricalLiquidationsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HistoricalLiquidationsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
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
			if m.From == nil {
				m.From = &types.Timestamp{}
			}
			if err := m.From.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field To", wireType)
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
			if m.To == nil {
				m.To = &types.Timestamp{}
			}
			if err := m.To.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *HistoricalLiquidationsResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: HistoricalLiquidationsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HistoricalLiquidationsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Liquidations", wireType)
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
			m.Liquidations = append(m.Liquidations, &models.Liquidation{})
			if err := m.Liquidations[len(m.Liquidations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
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
		case 5:
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
func (m *MarketDataResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MarketDataResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MarketDataResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field SnapshotL1", wireType)
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
			if m.SnapshotL1 == nil {
				m.SnapshotL1 = &models.OBL1Snapshot{}
			}
			if err := m.SnapshotL1.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
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
		case 5:
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
		case 6:
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
		case 7:
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
		case 8:
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
		case 9:
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
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpdateL1", wireType)
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
			if m.UpdateL1 == nil {
				m.UpdateL1 = &models.OBL1Update{}
			}
			if err := m.UpdateL1.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
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
		case 6:
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
		case 7:
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
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Liquidation", wireType)
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
			if m.Liquidation == nil {
				m.Liquidation = &models.Liquidation{}
			}
			if err := m.Liquidation.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Funding", wireType)
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
			if m.Funding == nil {
				m.Funding = &models.Funding{}
			}
			if err := m.Funding.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stats", wireType)
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
			m.Stats = append(m.Stats, &models.Stat{})
			if err := m.Stats[len(m.Stats)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *AccountDataRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: AccountDataRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccountDataRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *AccountDataResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: AccountDataResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccountDataResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
			if m.Orders == nil {
				m.Orders = &OrderList{}
			}
			if err := m.Orders.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
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
			if m.Positions == nil {
				m.Positions = &PositionList{}
			}
			if err := m.Positions.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Balances", wireType)
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
			if m.Balances == nil {
				m.Balances = &BalanceList{}
			}
			if err := m.Balances.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
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
func (m *AccountDataIncrementalRefresh) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: AccountDataIncrementalRefresh: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccountDataIncrementalRefresh: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Report", wireType)
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
			if m.Report == nil {
				m.Report = &ExecutionReport{}
			}
			if err := m.Report.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *SecurityDefinitionRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: SecurityDefinitionRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SecurityDefinitionRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *SecurityDefinitionResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: SecurityDefinitionResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SecurityDefinitionResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Security", wireType)
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
			if m.Security == nil {
				m.Security = &models.Security{}
			}
			if err := m.Security.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
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
		case 5:
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
		case 4:
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
		case 5:
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeAmount", wireType)
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
			if m.FeeAmount == nil {
				m.FeeAmount = &types.DoubleValue{}
			}
			if err := m.FeeAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 16:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeCurrency", wireType)
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
			if m.FeeCurrency == nil {
				m.FeeCurrency = &models1.Asset{}
			}
			if err := m.FeeCurrency.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 17:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeType", wireType)
			}
			m.FeeType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FeeType |= FeeType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 18:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeBasis", wireType)
			}
			m.FeeBasis = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FeeBasis |= FeeBasis(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 19:
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
		case 4:
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
		case 5:
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
func (m *BalancesRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: BalancesRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BalancesRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Asset", wireType)
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
			if m.Asset == nil {
				m.Asset = &models1.Asset{}
			}
			if err := m.Asset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *BalanceList) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: BalanceList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BalanceList: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Balances", wireType)
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
			m.Balances = append(m.Balances, &models.Balance{})
			if err := m.Balances[len(m.Balances)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
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
		case 5:
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
		case 10:
			if wireType == 0 {
				var v ExecutionInstruction
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowExecutorMessages
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
						return ErrIntOverflowExecutorMessages
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
					return ErrInvalidLengthExecutorMessages
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthExecutorMessages
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
							return ErrIntOverflowExecutorMessages
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
		case 4:
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
		case 5:
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
		case 5:
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
func (m *OrderUpdate) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OrderUpdate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderUpdate: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field OrigClientOrderID", wireType)
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
			if m.OrigClientOrderID == nil {
				m.OrigClientOrderID = &types.StringValue{}
			}
			if err := m.OrigClientOrderID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Quantity", wireType)
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
			if m.Quantity == nil {
				m.Quantity = &types.DoubleValue{}
			}
			if err := m.Quantity.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
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
func (m *OrderReplaceRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OrderReplaceRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderReplaceRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Update", wireType)
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
			if m.Update == nil {
				m.Update = &OrderUpdate{}
			}
			if err := m.Update.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *OrderReplaceResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OrderReplaceResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderReplaceResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *OrderBulkReplaceRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OrderBulkReplaceRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderBulkReplaceRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Updates", wireType)
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
			m.Updates = append(m.Updates, &OrderUpdate{})
			if err := m.Updates[len(m.Updates)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *OrderBulkReplaceResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OrderBulkReplaceResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderBulkReplaceResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
