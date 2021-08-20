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
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
	UnsupportedRequest             RejectionReason = 24
	UnknownAccount                 RejectionReason = 25
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
	24: "UnsupportedRequest",
	25: "UnknownAccount",
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
	"UnsupportedRequest":             24,
	"UnknownAccount":                 25,
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
	return fileDescriptor_350c53ba9303a7e6, []int{4}
}

type UpdateReason int32

const (
	FundingFee      UpdateReason = 0
	BalanceTransfer UpdateReason = 1
	Deposit         UpdateReason = 2
	Withdraw        UpdateReason = 3
	Order           UpdateReason = 4
	Unknown         UpdateReason = 5
)

var UpdateReason_name = map[int32]string{
	0: "FundingFee",
	1: "BalanceTransfer",
	2: "Deposit",
	3: "Withdraw",
	4: "Order",
	5: "Unknown",
}

var UpdateReason_value = map[string]int32{
	"FundingFee":      0,
	"BalanceTransfer": 1,
	"Deposit":         2,
	"Withdraw":        3,
	"Order":           4,
	"Unknown":         5,
}

func (UpdateReason) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{5}
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
		n, err := m.MarshalToSizedBuffer(b)
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
		n, err := m.MarshalToSizedBuffer(b)
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

type MarketStatisticsRequest struct {
	RequestID  uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Instrument *models.Instrument `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Statistics []models.StatType  `protobuf:"varint,3,rep,packed,name=statistics,proto3,enum=models.StatType" json:"statistics,omitempty"`
}

func (m *MarketStatisticsRequest) Reset()      { *m = MarketStatisticsRequest{} }
func (*MarketStatisticsRequest) ProtoMessage() {}
func (*MarketStatisticsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{2}
}
func (m *MarketStatisticsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MarketStatisticsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MarketStatisticsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MarketStatisticsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MarketStatisticsRequest.Merge(m, src)
}
func (m *MarketStatisticsRequest) XXX_Size() int {
	return m.Size()
}
func (m *MarketStatisticsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MarketStatisticsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MarketStatisticsRequest proto.InternalMessageInfo

func (m *MarketStatisticsRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *MarketStatisticsRequest) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *MarketStatisticsRequest) GetStatistics() []models.StatType {
	if m != nil {
		return m.Statistics
	}
	return nil
}

type MarketStatisticsResponse struct {
	RequestID       uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64          `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Statistics      []*models.Stat  `protobuf:"bytes,3,rep,name=statistics,proto3" json:"statistics,omitempty"`
	Success         bool            `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *MarketStatisticsResponse) Reset()      { *m = MarketStatisticsResponse{} }
func (*MarketStatisticsResponse) ProtoMessage() {}
func (*MarketStatisticsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{3}
}
func (m *MarketStatisticsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MarketStatisticsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MarketStatisticsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MarketStatisticsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MarketStatisticsResponse.Merge(m, src)
}
func (m *MarketStatisticsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MarketStatisticsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MarketStatisticsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MarketStatisticsResponse proto.InternalMessageInfo

func (m *MarketStatisticsResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *MarketStatisticsResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *MarketStatisticsResponse) GetStatistics() []*models.Stat {
	if m != nil {
		return m.Statistics
	}
	return nil
}

func (m *MarketStatisticsResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *MarketStatisticsResponse) GetRejectionReason() RejectionReason {
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
	return fileDescriptor_350c53ba9303a7e6, []int{4}
}
func (m *MarketDataRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MarketDataRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MarketDataRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{5}
}
func (m *MarketDataResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MarketDataResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MarketDataResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{6}
}
func (m *MarketDataIncrementalRefresh) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MarketDataIncrementalRefresh) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MarketDataIncrementalRefresh.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{7}
}
func (m *AccountDataRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountDataRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccountDataRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	RequestID       uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64             `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Securities      []*models.Security `protobuf:"bytes,3,rep,name=securities,proto3" json:"securities,omitempty"`
	Orders          []*models.Order    `protobuf:"bytes,4,rep,name=orders,proto3" json:"orders,omitempty"`
	Positions       []*models.Position `protobuf:"bytes,5,rep,name=positions,proto3" json:"positions,omitempty"`
	Balances        []*models.Balance  `protobuf:"bytes,6,rep,name=balances,proto3" json:"balances,omitempty"`
	MakerFee        *types.DoubleValue `protobuf:"bytes,7,opt,name=maker_fee,json=makerFee,proto3" json:"maker_fee,omitempty"`
	TakerFee        *types.DoubleValue `protobuf:"bytes,8,opt,name=taker_fee,json=takerFee,proto3" json:"taker_fee,omitempty"`
	Success         bool               `protobuf:"varint,9,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason    `protobuf:"varint,10,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
	SeqNum          uint64             `protobuf:"varint,11,opt,name=seq_num,json=seqNum,proto3" json:"seq_num,omitempty"`
}

func (m *AccountDataResponse) Reset()      { *m = AccountDataResponse{} }
func (*AccountDataResponse) ProtoMessage() {}
func (*AccountDataResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{8}
}
func (m *AccountDataResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountDataResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccountDataResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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

func (m *AccountDataResponse) GetSecurities() []*models.Security {
	if m != nil {
		return m.Securities
	}
	return nil
}

func (m *AccountDataResponse) GetOrders() []*models.Order {
	if m != nil {
		return m.Orders
	}
	return nil
}

func (m *AccountDataResponse) GetPositions() []*models.Position {
	if m != nil {
		return m.Positions
	}
	return nil
}

func (m *AccountDataResponse) GetBalances() []*models.Balance {
	if m != nil {
		return m.Balances
	}
	return nil
}

func (m *AccountDataResponse) GetMakerFee() *types.DoubleValue {
	if m != nil {
		return m.MakerFee
	}
	return nil
}

func (m *AccountDataResponse) GetTakerFee() *types.DoubleValue {
	if m != nil {
		return m.TakerFee
	}
	return nil
}

func (m *AccountDataResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *AccountDataResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
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
	return fileDescriptor_350c53ba9303a7e6, []int{9}
}
func (m *AccountDataIncrementalRefresh) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountDataIncrementalRefresh) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccountDataIncrementalRefresh.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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

type TradeCaptureReportRequest struct {
	RequestID uint64                    `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Account   *models.Account           `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	Filter    *TradeCaptureReportFilter `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (m *TradeCaptureReportRequest) Reset()      { *m = TradeCaptureReportRequest{} }
func (*TradeCaptureReportRequest) ProtoMessage() {}
func (*TradeCaptureReportRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{10}
}
func (m *TradeCaptureReportRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TradeCaptureReportRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TradeCaptureReportRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TradeCaptureReportRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TradeCaptureReportRequest.Merge(m, src)
}
func (m *TradeCaptureReportRequest) XXX_Size() int {
	return m.Size()
}
func (m *TradeCaptureReportRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TradeCaptureReportRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TradeCaptureReportRequest proto.InternalMessageInfo

func (m *TradeCaptureReportRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *TradeCaptureReportRequest) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

func (m *TradeCaptureReportRequest) GetFilter() *TradeCaptureReportFilter {
	if m != nil {
		return m.Filter
	}
	return nil
}

type TradeCaptureReportFilter struct {
	OrderID       *types.StringValue `protobuf:"bytes,1,opt,name=orderID,proto3" json:"orderID,omitempty"`
	ClientOrderID *types.StringValue `protobuf:"bytes,2,opt,name=client_orderID,json=clientOrderID,proto3" json:"client_orderID,omitempty"`
	Instrument    *models.Instrument `protobuf:"bytes,3,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Side          *SideValue         `protobuf:"bytes,4,opt,name=side,proto3" json:"side,omitempty"`
	From          *types.Timestamp   `protobuf:"bytes,5,opt,name=from,proto3" json:"from,omitempty"`
	To            *types.Timestamp   `protobuf:"bytes,6,opt,name=to,proto3" json:"to,omitempty"`
}

func (m *TradeCaptureReportFilter) Reset()      { *m = TradeCaptureReportFilter{} }
func (*TradeCaptureReportFilter) ProtoMessage() {}
func (*TradeCaptureReportFilter) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{11}
}
func (m *TradeCaptureReportFilter) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TradeCaptureReportFilter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TradeCaptureReportFilter.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TradeCaptureReportFilter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TradeCaptureReportFilter.Merge(m, src)
}
func (m *TradeCaptureReportFilter) XXX_Size() int {
	return m.Size()
}
func (m *TradeCaptureReportFilter) XXX_DiscardUnknown() {
	xxx_messageInfo_TradeCaptureReportFilter.DiscardUnknown(m)
}

var xxx_messageInfo_TradeCaptureReportFilter proto.InternalMessageInfo

func (m *TradeCaptureReportFilter) GetOrderID() *types.StringValue {
	if m != nil {
		return m.OrderID
	}
	return nil
}

func (m *TradeCaptureReportFilter) GetClientOrderID() *types.StringValue {
	if m != nil {
		return m.ClientOrderID
	}
	return nil
}

func (m *TradeCaptureReportFilter) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *TradeCaptureReportFilter) GetSide() *SideValue {
	if m != nil {
		return m.Side
	}
	return nil
}

func (m *TradeCaptureReportFilter) GetFrom() *types.Timestamp {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *TradeCaptureReportFilter) GetTo() *types.Timestamp {
	if m != nil {
		return m.To
	}
	return nil
}

type TradeCaptureReport struct {
	RequestID       uint64                 `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64                 `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Trades          []*models.TradeCapture `protobuf:"bytes,3,rep,name=trades,proto3" json:"trades,omitempty"`
	Success         bool                   `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason        `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *TradeCaptureReport) Reset()      { *m = TradeCaptureReport{} }
func (*TradeCaptureReport) ProtoMessage() {}
func (*TradeCaptureReport) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{12}
}
func (m *TradeCaptureReport) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TradeCaptureReport) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TradeCaptureReport.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TradeCaptureReport) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TradeCaptureReport.Merge(m, src)
}
func (m *TradeCaptureReport) XXX_Size() int {
	return m.Size()
}
func (m *TradeCaptureReport) XXX_DiscardUnknown() {
	xxx_messageInfo_TradeCaptureReport.DiscardUnknown(m)
}

var xxx_messageInfo_TradeCaptureReport proto.InternalMessageInfo

func (m *TradeCaptureReport) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *TradeCaptureReport) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *TradeCaptureReport) GetTrades() []*models.TradeCapture {
	if m != nil {
		return m.Trades
	}
	return nil
}

func (m *TradeCaptureReport) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *TradeCaptureReport) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type SecurityDefinitionRequest struct {
	RequestID  uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Instrument *models.Instrument `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
}

func (m *SecurityDefinitionRequest) Reset()      { *m = SecurityDefinitionRequest{} }
func (*SecurityDefinitionRequest) ProtoMessage() {}
func (*SecurityDefinitionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{13}
}
func (m *SecurityDefinitionRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SecurityDefinitionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SecurityDefinitionRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{14}
}
func (m *SecurityDefinitionResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SecurityDefinitionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SecurityDefinitionResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{15}
}
func (m *SecurityListRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SecurityListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SecurityListRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{16}
}
func (m *SecurityList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SecurityList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SecurityList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{17}
}
func (m *ExecutionReport) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExecutionReport) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExecutionReport.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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

type AccountUpdate struct {
	UpdateReason UpdateReason   `protobuf:"varint,1,opt,name=update_reason,json=updateReason,proto3,enum=messages.UpdateReason" json:"update_reason,omitempty"`
	Asset        *models1.Asset `protobuf:"bytes,2,opt,name=asset,proto3" json:"asset,omitempty"`
	Balance      float64        `protobuf:"fixed64,3,opt,name=balance,proto3" json:"balance,omitempty"`
}

func (m *AccountUpdate) Reset()      { *m = AccountUpdate{} }
func (*AccountUpdate) ProtoMessage() {}
func (*AccountUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{18}
}
func (m *AccountUpdate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccountUpdate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AccountUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountUpdate.Merge(m, src)
}
func (m *AccountUpdate) XXX_Size() int {
	return m.Size()
}
func (m *AccountUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_AccountUpdate proto.InternalMessageInfo

func (m *AccountUpdate) GetUpdateReason() UpdateReason {
	if m != nil {
		return m.UpdateReason
	}
	return FundingFee
}

func (m *AccountUpdate) GetAsset() *models1.Asset {
	if m != nil {
		return m.Asset
	}
	return nil
}

func (m *AccountUpdate) GetBalance() float64 {
	if m != nil {
		return m.Balance
	}
	return 0
}

type SideValue struct {
	Value models.Side `protobuf:"varint,1,opt,name=value,proto3,enum=models.Side" json:"value,omitempty"`
}

func (m *SideValue) Reset()      { *m = SideValue{} }
func (*SideValue) ProtoMessage() {}
func (*SideValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{19}
}
func (m *SideValue) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SideValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SideValue.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{20}
}
func (m *OrderStatusValue) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderStatusValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderStatusValue.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{21}
}
func (m *OrderFilter) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderFilter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderFilter.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{22}
}
func (m *OrderStatusRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderStatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderStatusRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{23}
}
func (m *OrderList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{24}
}
func (m *PositionsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PositionsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PositionsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{25}
}
func (m *PositionList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PositionList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PositionList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{26}
}
func (m *BalancesRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BalancesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BalancesRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{27}
}
func (m *BalanceList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BalanceList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BalanceList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{28}
}
func (m *NewOrder) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NewOrder) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NewOrder.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{29}
}
func (m *NewOrderSingleRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NewOrderSingleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NewOrderSingleRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{30}
}
func (m *NewOrderSingleResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NewOrderSingleResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NewOrderSingleResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{31}
}
func (m *NewOrderBulkRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NewOrderBulkRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NewOrderBulkRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{32}
}
func (m *NewOrderBulkResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NewOrderBulkResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NewOrderBulkResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{33}
}
func (m *OrderUpdate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderUpdate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{34}
}
func (m *OrderReplaceRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderReplaceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderReplaceRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	OrderID         string          `protobuf:"bytes,3,opt,name=orderID,proto3" json:"orderID,omitempty"`
	Success         bool            `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *OrderReplaceResponse) Reset()      { *m = OrderReplaceResponse{} }
func (*OrderReplaceResponse) ProtoMessage() {}
func (*OrderReplaceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{35}
}
func (m *OrderReplaceResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderReplaceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderReplaceResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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

func (m *OrderReplaceResponse) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
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
	return fileDescriptor_350c53ba9303a7e6, []int{36}
}
func (m *OrderBulkReplaceRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderBulkReplaceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderBulkReplaceRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{37}
}
func (m *OrderBulkReplaceResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderBulkReplaceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderBulkReplaceResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{38}
}
func (m *OrderCancelRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderCancelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderCancelRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{39}
}
func (m *OrderCancelResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderCancelResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderCancelResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{40}
}
func (m *OrderMassCancelRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderMassCancelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderMassCancelRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	return fileDescriptor_350c53ba9303a7e6, []int{41}
}
func (m *OrderMassCancelResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderMassCancelResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderMassCancelResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
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
	proto.RegisterEnum("messages.UpdateReason", UpdateReason_name, UpdateReason_value)
	proto.RegisterType((*HistoricalLiquidationsRequest)(nil), "messages.HistoricalLiquidationsRequest")
	proto.RegisterType((*HistoricalLiquidationsResponse)(nil), "messages.HistoricalLiquidationsResponse")
	proto.RegisterType((*MarketStatisticsRequest)(nil), "messages.MarketStatisticsRequest")
	proto.RegisterType((*MarketStatisticsResponse)(nil), "messages.MarketStatisticsResponse")
	proto.RegisterType((*MarketDataRequest)(nil), "messages.MarketDataRequest")
	proto.RegisterType((*MarketDataResponse)(nil), "messages.MarketDataResponse")
	proto.RegisterType((*MarketDataIncrementalRefresh)(nil), "messages.MarketDataIncrementalRefresh")
	proto.RegisterType((*AccountDataRequest)(nil), "messages.AccountDataRequest")
	proto.RegisterType((*AccountDataResponse)(nil), "messages.AccountDataResponse")
	proto.RegisterType((*AccountDataIncrementalRefresh)(nil), "messages.AccountDataIncrementalRefresh")
	proto.RegisterType((*TradeCaptureReportRequest)(nil), "messages.TradeCaptureReportRequest")
	proto.RegisterType((*TradeCaptureReportFilter)(nil), "messages.TradeCaptureReportFilter")
	proto.RegisterType((*TradeCaptureReport)(nil), "messages.TradeCaptureReport")
	proto.RegisterType((*SecurityDefinitionRequest)(nil), "messages.SecurityDefinitionRequest")
	proto.RegisterType((*SecurityDefinitionResponse)(nil), "messages.SecurityDefinitionResponse")
	proto.RegisterType((*SecurityListRequest)(nil), "messages.SecurityListRequest")
	proto.RegisterType((*SecurityList)(nil), "messages.SecurityList")
	proto.RegisterType((*ExecutionReport)(nil), "messages.ExecutionReport")
	proto.RegisterType((*AccountUpdate)(nil), "messages.AccountUpdate")
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
	// 3063 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xdc, 0x5a, 0xcd, 0x6f, 0x1b, 0xc7,
	0xd9, 0xe7, 0x8a, 0x94, 0x48, 0x3e, 0xa4, 0xa4, 0xd5, 0xc8, 0xb6, 0x68, 0x46, 0x66, 0x94, 0xcd,
	0x9b, 0x37, 0x8e, 0xf2, 0x5a, 0x8a, 0xe9, 0xc4, 0x79, 0xdf, 0x37, 0xb5, 0x01, 0x59, 0x94, 0x10,
	0x02, 0x8a, 0xad, 0xae, 0xa4, 0x36, 0x40, 0x0f, 0xc4, 0x68, 0x39, 0xa4, 0xb6, 0x5a, 0xee, 0xd2,
	0xbb, 0xb3, 0xb1, 0x75, 0x28, 0xd0, 0x5e, 0x53, 0xa0, 0xc8, 0x3f, 0xd0, 0xa2, 0x87, 0x7e, 0x1d,
	0x0b, 0x14, 0x68, 0x7b, 0xe8, 0xa5, 0xa7, 0xb6, 0x87, 0x16, 0x3e, 0x15, 0x09, 0x50, 0x14, 0x89,
	0x02, 0x14, 0xb9, 0x35, 0xb7, 0xf6, 0x58, 0xcc, 0xc7, 0x2e, 0x67, 0x49, 0x4a, 0xa2, 0x62, 0x2a,
	0x4d, 0x7a, 0xdb, 0x9d, 0xf9, 0x3d, 0x33, 0xf3, 0x7c, 0xce, 0xf3, 0x3c, 0xbb, 0xb0, 0x40, 0x1e,
	0x13, 0x2b, 0xa4, 0x9e, 0xdf, 0xe8, 0x90, 0x20, 0xc0, 0x6d, 0x12, 0xac, 0x74, 0x7d, 0x8f, 0x7a,
	0x28, 0x17, 0xbd, 0x97, 0x9f, 0x6d, 0x7b, 0x5e, 0xdb, 0x21, 0xab, 0x7c, 0x7c, 0x3f, 0x6c, 0xad,
	0x52, 0xbb, 0x43, 0x02, 0x8a, 0x3b, 0x5d, 0x01, 0x2d, 0x57, 0xfa, 0x01, 0x8f, 0x7c, 0xdc, 0xed,
	0x12, 0x5f, 0x2e, 0x55, 0xae, 0xb6, 0x6d, 0xea, 0xe0, 0xfd, 0x15, 0xcb, 0xeb, 0xac, 0x62, 0xa7,
	0x7b, 0x80, 0xa9, 0x6d, 0x1d, 0x06, 0xab, 0x8f, 0xad, 0x03, 0xec, 0xb6, 0x89, 0xbf, 0xda, 0xf1,
	0x9a, 0xc4, 0x09, 0x04, 0x79, 0x44, 0xf3, 0x95, 0xe1, 0x34, 0xfc, 0xf1, 0x86, 0xe5, 0xb9, 0x2e,
	0xb1, 0x68, 0x44, 0xd8, 0xc1, 0xfe, 0x21, 0xa1, 0x8d, 0x26, 0xa6, 0x58, 0x52, 0xdf, 0x3d, 0x07,
	0x75, 0x40, 0xac, 0xd0, 0xb7, 0xe9, 0x91, 0x4a, 0x7f, 0xe7, 0x1c, 0xf4, 0xd8, 0xb2, 0xbc, 0xd0,
	0x4d, 0x6c, 0x7f, 0xbb, 0x6d, 0xd3, 0x83, 0x50, 0x90, 0xaf, 0x05, 0x47, 0xee, 0xa1, 0xef, 0xb9,
	0xf5, 0x5d, 0xc1, 0x20, 0xb6, 0xa8, 0xe7, 0xdf, 0x68, 0x7b, 0xab, 0xfc, 0x21, 0xc1, 0xb4, 0xf1,
	0x47, 0x0d, 0xae, 0xbd, 0x69, 0x07, 0xd4, 0xf3, 0x6d, 0x0b, 0x3b, 0x5b, 0xf6, 0xc3, 0xd0, 0x6e,
	0x62, 0x6a, 0x7b, 0x6e, 0x60, 0x92, 0x87, 0x21, 0x09, 0x28, 0x5a, 0x84, 0xbc, 0x2f, 0x1e, 0xeb,
	0xb5, 0x92, 0xb6, 0xa4, 0x5d, 0xcf, 0x98, 0xbd, 0x01, 0x54, 0x05, 0xb0, 0xdd, 0x80, 0xfa, 0x61,
	0x87, 0xb8, 0xb4, 0x34, 0xb1, 0xa4, 0x5d, 0x2f, 0x54, 0xd1, 0x8a, 0x38, 0xe7, 0x4a, 0x3d, 0x9e,
	0x31, 0x15, 0x14, 0x5a, 0x81, 0x4c, 0xcb, 0xf7, 0x3a, 0xa5, 0x34, 0x47, 0x97, 0x57, 0x84, 0x2e,
	0x57, 0x22, 0x5d, 0xae, 0xec, 0x46, 0xca, 0x36, 0x39, 0x0e, 0x2d, 0xc3, 0x04, 0xf5, 0x4a, 0x99,
	0x33, 0xd1, 0x13, 0xd4, 0x33, 0xfe, 0xa1, 0x41, 0xe5, 0x24, 0x7e, 0x82, 0xae, 0xe7, 0x06, 0xe4,
	0x0c, 0x86, 0x2a, 0x00, 0xbe, 0x44, 0xd6, 0x6b, 0x9c, 0xa1, 0x8c, 0xa9, 0x8c, 0xa0, 0xd7, 0xa1,
	0xe8, 0x28, 0xab, 0x96, 0xd2, 0x4b, 0xe9, 0xeb, 0x85, 0xea, 0x7c, 0xc4, 0xb2, 0xb2, 0xa3, 0x99,
	0x00, 0xa2, 0x12, 0x64, 0x83, 0xd0, 0xb2, 0x48, 0x10, 0x70, 0x56, 0x72, 0x66, 0xf4, 0x8a, 0x6a,
	0xa0, 0xfb, 0xe4, 0x9b, 0xc4, 0x62, 0xb8, 0x86, 0x4f, 0x70, 0xe0, 0xb9, 0xa5, 0xc9, 0x25, 0xed,
	0xfa, 0x4c, 0xf5, 0xea, 0x4a, 0xec, 0x22, 0x66, 0x84, 0x30, 0x39, 0xc0, 0x9c, 0xf5, 0x93, 0x03,
	0xc6, 0xf7, 0x35, 0x58, 0x78, 0x8b, 0x9b, 0xe5, 0x0e, 0xc5, 0xd4, 0x0e, 0xa8, 0x6d, 0x5d, 0xa0,
	0x0e, 0x5f, 0x01, 0x08, 0xe2, 0x6d, 0xb8, 0x10, 0x66, 0xaa, 0x7a, 0x44, 0xc3, 0x0e, 0xb0, 0x7b,
	0xd4, 0x25, 0xa6, 0x82, 0x31, 0xfe, 0xa6, 0x41, 0x69, 0xf0, 0x7c, 0x63, 0xd1, 0xc9, 0xff, 0x0c,
	0x1c, 0xa6, 0x50, 0x2d, 0xaa, 0x87, 0x51, 0x0f, 0x72, 0xe1, 0x8a, 0xf8, 0xbb, 0x06, 0x73, 0x82,
	0xd1, 0x1a, 0xa6, 0x78, 0x34, 0x15, 0x2c, 0x42, 0x3e, 0x08, 0xf7, 0x03, 0xcb, 0xb7, 0xf7, 0x09,
	0x67, 0x30, 0x67, 0xf6, 0x06, 0xd0, 0x32, 0x40, 0xfc, 0xe2, 0x4b, 0xb7, 0x81, 0x15, 0xee, 0xcd,
	0x2b, 0xdb, 0xf5, 0x9a, 0xa9, 0xcc, 0xf6, 0x29, 0x33, 0x33, 0x92, 0x32, 0xef, 0x42, 0x01, 0xb7,
	0xdb, 0x3e, 0x69, 0x73, 0x53, 0x95, 0x2c, 0x2f, 0x46, 0x44, 0x0f, 0xfc, 0x26, 0xf1, 0xef, 0x79,
	0xde, 0xe1, 0x5a, 0x0f, 0x63, 0xaa, 0x04, 0xc6, 0x0f, 0xd2, 0x80, 0x54, 0x8e, 0xc7, 0xa2, 0xd4,
	0x57, 0x01, 0x02, 0x17, 0x77, 0x83, 0x03, 0x8f, 0x6e, 0xdd, 0x94, 0x4c, 0x5f, 0x8a, 0xcf, 0x74,
	0x6f, 0xeb, 0xe6, 0x8e, 0x9c, 0x35, 0x15, 0x5c, 0x82, 0xaa, 0x2a, 0xd9, 0x57, 0xa9, 0xaa, 0x43,
	0xa8, 0xaa, 0x09, 0xaa, 0x5b, 0x9c, 0xff, 0x24, 0xd5, 0xad, 0x21, 0x54, 0xb7, 0xd0, 0x2a, 0x4c,
	0x51, 0x1f, 0x37, 0x49, 0x50, 0x9a, 0xe2, 0x26, 0xb7, 0x10, 0x51, 0x44, 0x82, 0x22, 0xcd, 0x5d,
	0x36, 0x6f, 0x4a, 0x18, 0x5a, 0x80, 0x6c, 0x40, 0x1e, 0x36, 0xdc, 0xb0, 0x53, 0xca, 0x72, 0x7e,
	0xa7, 0x02, 0xf2, 0xf0, 0x7e, 0xd8, 0x51, 0x4d, 0x32, 0x77, 0xb6, 0x49, 0xe6, 0xcf, 0x6d, 0x92,
	0xbf, 0x49, 0xc3, 0x62, 0x4f, 0x41, 0x75, 0xd7, 0xf2, 0x09, 0x53, 0x3c, 0x76, 0x4c, 0xd2, 0xf2,
	0x49, 0x70, 0xf0, 0x94, 0xaa, 0x52, 0xf8, 0x4a, 0x27, 0xf8, 0x5a, 0x81, 0x5c, 0xd8, 0x6d, 0x62,
	0x4a, 0xb6, 0x6e, 0xf6, 0x9b, 0x22, 0xd3, 0xe0, 0x1e, 0x9f, 0x33, 0x63, 0x8c, 0x82, 0xaf, 0x4a,
	0x2d, 0xa8, 0xf8, 0x6a, 0x1f, 0xbe, 0xaa, 0xe0, 0x6f, 0x95, 0xa6, 0x06, 0xf0, 0xb7, 0xfa, 0xf0,
	0xaa, 0xc6, 0xb2, 0xa3, 0x69, 0xec, 0x35, 0x28, 0x28, 0x41, 0x9c, 0x2b, 0xe7, 0x84, 0x60, 0xaf,
	0xe2, 0xd0, 0x4b, 0x90, 0x6d, 0x85, 0x6e, 0xd3, 0x76, 0xdb, 0x5c, 0x59, 0x85, 0xea, 0x6c, 0x44,
	0xb2, 0x29, 0x86, 0xcd, 0x68, 0x1e, 0x19, 0x30, 0xc9, 0x62, 0x53, 0x50, 0x82, 0x21, 0x61, 0x4b,
	0x4c, 0x19, 0x3f, 0xd5, 0x00, 0xad, 0x89, 0x3b, 0xff, 0xdf, 0x11, 0x52, 0x5e, 0x82, 0xac, 0xcc,
	0x38, 0xa4, 0x12, 0x63, 0x6e, 0xe4, 0xa1, 0xcc, 0x68, 0xde, 0xf8, 0x67, 0x1a, 0xe6, 0x13, 0x27,
	0x1d, 0x4b, 0x28, 0x60, 0x97, 0x8d, 0x48, 0x99, 0x6c, 0x12, 0xc5, 0xf7, 0xde, 0x65, 0x23, 0x93,
	0x29, 0x53, 0xc1, 0xa0, 0x17, 0x60, 0xca, 0x63, 0x61, 0x8b, 0x85, 0x78, 0x86, 0x9e, 0x4e, 0x04,
	0x33, 0x53, 0x4e, 0xa2, 0x15, 0xc8, 0x77, 0xbd, 0xc0, 0x16, 0x37, 0xf9, 0x64, 0x72, 0xdd, 0x6d,
	0x39, 0x61, 0xf6, 0x20, 0xe8, 0x65, 0xc8, 0xed, 0x63, 0x07, 0xbb, 0x56, 0xec, 0xf3, 0xb1, 0x28,
	0xee, 0x89, 0x71, 0x33, 0x06, 0xa0, 0xff, 0x83, 0x7c, 0x07, 0x1f, 0x12, 0xbf, 0xd1, 0x22, 0x84,
	0xfb, 0x7b, 0xa1, 0xba, 0x38, 0x90, 0xbd, 0xd4, 0xbc, 0x70, 0xdf, 0x21, 0x5f, 0xc3, 0x4e, 0x48,
	0xcc, 0x1c, 0x87, 0x6f, 0x12, 0xc2, 0x48, 0x69, 0x4c, 0x9a, 0x1b, 0x85, 0x94, 0x46, 0xa4, 0x4a,
	0x28, 0xc9, 0x9f, 0x1d, 0x4a, 0xe0, 0xbc, 0xa1, 0x44, 0xf5, 0xf5, 0x82, 0xea, 0xeb, 0xc6, 0x7b,
	0x1a, 0x5c, 0x53, 0x54, 0x3f, 0xf6, 0x20, 0x73, 0x13, 0xa6, 0x7c, 0xd2, 0xf5, 0x7c, 0x2a, 0xad,
	0x55, 0x39, 0xf4, 0x06, 0x2f, 0x28, 0xf8, 0x19, 0x19, 0xc0, 0x94, 0x40, 0xe3, 0x87, 0x1a, 0x5c,
	0xe5, 0xfe, 0xbc, 0x8e, 0xbb, 0x34, 0xf4, 0x89, 0x9c, 0x1e, 0xc9, 0x7d, 0x14, 0xa3, 0x9f, 0x38,
	0xdd, 0xe8, 0xd1, 0xff, 0xc3, 0x54, 0xcb, 0x76, 0x68, 0xec, 0x47, 0x46, 0xef, 0x64, 0x83, 0xbb,
	0x6f, 0x72, 0xa4, 0x29, 0x29, 0x8c, 0x27, 0x13, 0x50, 0x3a, 0x09, 0x84, 0x6e, 0x43, 0x96, 0x1b,
	0xaa, 0x3c, 0xdf, 0x30, 0x23, 0xd8, 0xa1, 0xbe, 0xed, 0xb6, 0x85, 0x11, 0x44, 0x60, 0xb4, 0x0e,
	0x33, 0x96, 0x63, 0x13, 0x97, 0x36, 0x22, 0xf2, 0x89, 0x11, 0xc8, 0xa7, 0x05, 0xcd, 0x03, 0xb9,
	0x48, 0x32, 0x91, 0x48, 0x8f, 0x94, 0x48, 0xbc, 0x08, 0x99, 0xc0, 0x6e, 0x12, 0x19, 0x26, 0xe6,
	0x7b, 0x72, 0xd8, 0xb1, 0x9b, 0xd2, 0x52, 0x39, 0x20, 0x2e, 0x01, 0x26, 0xcf, 0x55, 0x02, 0x4c,
	0x8d, 0x54, 0x02, 0x7c, 0xa8, 0x01, 0x1a, 0x14, 0xe9, 0x53, 0xa7, 0x98, 0xd1, 0xcd, 0x21, 0xc2,
	0x4f, 0x9c, 0x1d, 0x24, 0x76, 0x8a, 0xae, 0x8d, 0x8b, 0x4e, 0x31, 0x3b, 0x70, 0x35, 0x0a, 0x7b,
	0x35, 0xd2, 0xb2, 0x5d, 0x11, 0xa8, 0x2e, 0x2a, 0xd9, 0x37, 0x3e, 0xd1, 0xa0, 0x3c, 0x6c, 0xbf,
	0x31, 0x25, 0xef, 0xb9, 0xa8, 0x1e, 0x96, 0x56, 0x36, 0x18, 0xda, 0x63, 0xc4, 0x85, 0x4b, 0xf6,
	0x5b, 0x30, 0x1f, 0xed, 0xba, 0x65, 0x07, 0xf4, 0x73, 0xbe, 0x6a, 0x8d, 0xbf, 0x6a, 0x50, 0x54,
	0xf7, 0xff, 0xdc, 0x2f, 0xce, 0x8b, 0x96, 0xef, 0x2f, 0xb3, 0x30, 0xdb, 0x17, 0xae, 0xd5, 0x2b,
	0x45, 0xeb, 0x4f, 0x8b, 0xd5, 0x00, 0x96, 0x3f, 0x2d, 0xc2, 0xa5, 0xcf, 0x1f, 0xe1, 0x96, 0xa0,
	0x40, 0xa2, 0xa3, 0xd4, 0x6b, 0x9c, 0xdf, 0xbc, 0xa9, 0x0e, 0xa1, 0xbb, 0x30, 0x13, 0xbf, 0x36,
	0xe8, 0x51, 0x97, 0x48, 0x8e, 0x17, 0x86, 0xdc, 0x3d, 0xbc, 0xe0, 0x9d, 0x26, 0xea, 0x2b, 0xba,
	0x0d, 0x45, 0x7e, 0xbe, 0x06, 0xcb, 0xe3, 0xc2, 0x80, 0x07, 0xb0, 0x99, 0x5e, 0xfe, 0xc8, 0x0f,
	0xb2, 0xc3, 0xa7, 0xcc, 0x82, 0xd7, 0x7b, 0xe9, 0x73, 0xd2, 0xec, 0x88, 0xb1, 0x77, 0xd6, 0x21,
	0xf8, 0x1d, 0x12, 0x34, 0x1e, 0x86, 0xd8, 0xa5, 0xcc, 0x9d, 0x58, 0x02, 0xa0, 0x99, 0x33, 0x62,
	0xf8, 0xab, 0x72, 0x14, 0x3d, 0x07, 0x45, 0x2b, 0xec, 0xf4, 0x50, 0xc0, 0x51, 0x05, 0x2b, 0xec,
	0xc4, 0x90, 0x0d, 0xd0, 0xa9, 0x8f, 0xdd, 0x00, 0x0b, 0x6d, 0x53, 0xbb, 0x43, 0xf8, 0x6d, 0x7f,
	0x7a, 0xf0, 0x9d, 0x55, 0x68, 0xd8, 0x28, 0xbb, 0xbf, 0x78, 0x40, 0xac, 0xd7, 0x4a, 0xc5, 0x51,
	0xee, 0x2f, 0x09, 0x46, 0x6f, 0x00, 0xb4, 0x6c, 0xc7, 0x69, 0x74, 0x7d, 0xdb, 0x22, 0xa5, 0xe9,
	0x11, 0xf2, 0x9f, 0x3c, 0xc3, 0x6f, 0x33, 0x38, 0x5a, 0x83, 0x69, 0x4e, 0x1c, 0xf3, 0x37, 0x33,
	0x02, 0x7d, 0x91, 0x91, 0xc4, 0xec, 0xb3, 0xfd, 0x09, 0x69, 0xe0, 0x0e, 0xbf, 0xfe, 0x67, 0x47,
	0xda, 0x9f, 0x90, 0x35, 0x0e, 0x47, 0xaf, 0x40, 0x91, 0x11, 0x5b, 0xa1, 0xef, 0x13, 0xd7, 0x3a,
	0x2a, 0xe9, 0x9c, 0x3c, 0x4e, 0x40, 0xd7, 0x82, 0x80, 0x50, 0xb3, 0xd0, 0x22, 0x64, 0x5d, 0x22,
	0x58, 0x04, 0x64, 0x14, 0xdc, 0xbe, 0xe6, 0xb8, 0x85, 0xcc, 0xf5, 0xec, 0x6b, 0x93, 0x10, 0x6e,
	0x59, 0xd9, 0x96, 0x78, 0x40, 0xab, 0xc0, 0x36, 0x6b, 0xec, 0xe3, 0xc0, 0x0e, 0x4a, 0x88, 0xc3,
	0x51, 0x02, 0x7e, 0x8f, 0xcd, 0x98, 0x6c, 0x49, 0xfe, 0x34, 0xd4, 0x71, 0xe7, 0xcf, 0xed, 0xb8,
	0xdf, 0xd5, 0x60, 0x5a, 0x66, 0x3e, 0xa2, 0xac, 0x42, 0x6f, 0xc0, 0xb4, 0x28, 0xac, 0xa2, 0x45,
	0x35, 0xbe, 0xe8, 0x95, 0xde, 0xa2, 0xb2, 0xfe, 0x12, 0x2b, 0x16, 0x43, 0xe5, 0x0d, 0x3d, 0x0f,
	0x93, 0x98, 0x49, 0x42, 0xde, 0x40, 0x7d, 0xe2, 0x11, 0x73, 0xcc, 0xff, 0x65, 0x36, 0xcd, 0xdd,
	0x5b, 0x33, 0xa3, 0x57, 0x63, 0x15, 0xf2, 0x71, 0x4a, 0xc1, 0x4a, 0xa8, 0x77, 0xd8, 0x83, 0x3c,
	0x40, 0xaf, 0x84, 0xb2, 0x9b, 0xc4, 0x14, 0x53, 0xc6, 0x1d, 0xd0, 0x15, 0x6f, 0x13, 0x74, 0x2f,
	0x25, 0xe9, 0x86, 0xba, 0xa5, 0x24, 0xff, 0xc9, 0x04, 0x14, 0xf8, 0xf0, 0x7f, 0x7c, 0x66, 0x76,
	0xa7, 0x2f, 0x64, 0x45, 0x19, 0x5a, 0x4c, 0xd0, 0x2f, 0xc6, 0x44, 0xe4, 0x32, 0xfe, 0xac, 0x01,
	0x52, 0xe5, 0xf7, 0x85, 0x2d, 0x55, 0xd1, 0x8d, 0x38, 0x6b, 0x17, 0x2c, 0x5e, 0xee, 0x63, 0xb1,
	0x2f, 0x51, 0xff, 0xb9, 0x06, 0x79, 0x3e, 0x3e, 0xc2, 0xb5, 0xac, 0x5c, 0xa2, 0x13, 0xc9, 0x4b,
	0xb4, 0x57, 0x97, 0xa6, 0x4f, 0xab, 0x4b, 0x87, 0xb9, 0x6c, 0xe6, 0xdc, 0x2e, 0xfb, 0x81, 0x06,
	0x7a, 0x54, 0xc5, 0x06, 0x5f, 0x86, 0x3e, 0xa4, 0xa2, 0xbd, 0xc9, 0x33, 0x1a, 0x0d, 0x7f, 0xd1,
	0xa0, 0x18, 0xf1, 0x36, 0x86, 0x44, 0x29, 0xd1, 0x08, 0x48, 0x9f, 0xdd, 0x08, 0xb8, 0xe8, 0x34,
	0xe9, 0x0f, 0x1a, 0xcc, 0xca, 0x8e, 0xc2, 0xe7, 0xae, 0xb9, 0x38, 0x34, 0x67, 0x4e, 0x09, 0xcd,
	0xe7, 0x50, 0xd5, 0xfb, 0x1a, 0x14, 0x24, 0x2f, 0x63, 0xd0, 0x94, 0xda, 0x82, 0x49, 0x9f, 0xd5,
	0x82, 0xb9, 0x68, 0x35, 0xfd, 0x28, 0x0d, 0xb9, 0xfb, 0xe4, 0x11, 0x77, 0x5e, 0xf4, 0xc2, 0x40,
	0x6c, 0xd7, 0x78, 0x46, 0x79, 0x6a, 0xf4, 0x1e, 0xf9, 0x6b, 0x8b, 0x08, 0xca, 0x4a, 0x0e, 0x3a,
	0x97, 0x08, 0x1d, 0x3c, 0x47, 0xc8, 0x7b, 0xd1, 0x23, 0x7a, 0x39, 0xa2, 0xe0, 0x51, 0x7f, 0x6a,
	0xc8, 0xc5, 0x28, 0xc0, 0xec, 0x11, 0xbd, 0x0e, 0xd3, 0x2c, 0xc5, 0x6b, 0xd8, 0x6e, 0xa3, 0xe5,
	0xf9, 0x96, 0xe8, 0x56, 0x29, 0x17, 0x22, 0x4b, 0xe6, 0xea, 0xee, 0x26, 0x9b, 0x32, 0x0b, 0xb4,
	0xf7, 0x82, 0xca, 0x90, 0x8b, 0xd3, 0xac, 0x1c, 0xbf, 0xa1, 0xe3, 0x77, 0x54, 0x85, 0x49, 0x91,
	0xbf, 0xe5, 0x47, 0xc8, 0x9f, 0x04, 0x14, 0xed, 0xc1, 0x95, 0x5e, 0xbe, 0x2d, 0xf8, 0xb7, 0x84,
	0x4f, 0x02, 0xff, 0xc2, 0x54, 0x19, 0x92, 0x77, 0xd7, 0x7b, 0x30, 0xf3, 0x32, 0x19, 0x32, 0x1a,
	0x18, 0xef, 0x6a, 0x70, 0x39, 0x52, 0xd3, 0x8e, 0xed, 0xb6, 0x1d, 0x32, 0xf6, 0x1e, 0xd0, 0x75,
	0x98, 0xe4, 0xf2, 0xec, 0x5d, 0xc7, 0xd1, 0x41, 0xa3, 0x8d, 0x4d, 0x01, 0x30, 0x9e, 0x68, 0x70,
	0xa5, 0xff, 0x30, 0x63, 0x29, 0xa4, 0x15, 0x63, 0x4f, 0x27, 0x8d, 0x5d, 0xa9, 0xa3, 0x32, 0xc9,
	0x3a, 0x6a, 0x3c, 0x6e, 0xf0, 0xae, 0x06, 0xf3, 0x11, 0x4b, 0xf7, 0x42, 0xe7, 0x70, 0xec, 0xd2,
	0x5d, 0xee, 0xbb, 0x36, 0x87, 0x89, 0x57, 0x22, 0x98, 0x7c, 0x2f, 0x25, 0x0f, 0x33, 0x16, 0xe9,
	0x96, 0x21, 0x27, 0x85, 0x26, 0x0e, 0x91, 0x37, 0xe3, 0xf7, 0x0b, 0x0f, 0x33, 0xdf, 0x89, 0xb2,
	0x4f, 0x99, 0x79, 0x7f, 0xd6, 0xec, 0x73, 0x0b, 0xe6, 0x3d, 0xdf, 0x6e, 0x37, 0x3e, 0x43, 0x0a,
	0x3a, 0xc7, 0x08, 0xd7, 0x13, 0x81, 0xec, 0x7f, 0x15, 0xe7, 0x4f, 0x8f, 0xd2, 0xa3, 0x1e, 0x0c,
	0x0d, 0x99, 0x91, 0x43, 0x83, 0xf1, 0x5b, 0x0d, 0xe6, 0x85, 0xa2, 0x49, 0xd7, 0xc1, 0x16, 0xb9,
	0xb8, 0x4f, 0xdb, 0x8a, 0x5d, 0xa6, 0xcf, 0xce, 0x21, 0x45, 0x55, 0x23, 0x39, 0xe9, 0xcf, 0x21,
	0x65, 0x01, 0x24, 0x41, 0xc6, 0x9f, 0x34, 0xb8, 0x94, 0xe4, 0x61, 0x5c, 0x8e, 0xaf, 0x76, 0x41,
	0x14, 0xf7, 0xbe, 0x68, 0xc3, 0xfc, 0x9d, 0x06, 0x0b, 0x8a, 0xa3, 0x7d, 0x71, 0x14, 0xb3, 0x0a,
	0x59, 0x21, 0xf3, 0xe8, 0x03, 0xd0, 0x09, 0x9a, 0x89, 0x50, 0xc6, 0xaf, 0x35, 0x28, 0x0d, 0x72,
	0x72, 0xc1, 0x71, 0x79, 0x3c, 0x69, 0xfe, 0xf7, 0x26, 0x64, 0xc9, 0xb5, 0xce, 0x92, 0x1e, 0x67,
	0x34, 0xf9, 0xdf, 0x4e, 0xb6, 0xd6, 0x9e, 0xa2, 0x80, 0x4d, 0x3f, 0x6d, 0x01, 0x3b, 0xf6, 0xda,
	0xe0, 0x17, 0x51, 0xa8, 0x88, 0x04, 0xf2, 0xa5, 0x50, 0xe3, 0x7b, 0x1a, 0x5c, 0xe1, 0xa7, 0x7e,
	0x0b, 0x07, 0xc1, 0x79, 0x54, 0x79, 0x8e, 0x7b, 0xf4, 0x46, 0xdf, 0x97, 0xaa, 0x33, 0x6a, 0xde,
	0x5f, 0x45, 0xee, 0xad, 0x1e, 0xe9, 0xcb, 0x20, 0xcc, 0xe5, 0x4f, 0x34, 0x98, 0x4e, 0x74, 0x66,
	0x51, 0x16, 0xd2, 0xf7, 0xc9, 0x23, 0x3d, 0x85, 0x72, 0x90, 0xa9, 0x79, 0x2e, 0xd1, 0x35, 0x54,
	0x84, 0x9c, 0x60, 0x8a, 0x34, 0xf5, 0x09, 0xf6, 0x26, 0xfd, 0xbe, 0xa9, 0xa7, 0xd1, 0x1c, 0x4c,
	0x6f, 0x13, 0xfe, 0x85, 0x5e, 0x40, 0xf4, 0x0c, 0x2a, 0x40, 0x76, 0x87, 0x7a, 0xdd, 0x2e, 0x69,
	0xea, 0x93, 0x02, 0xcd, 0xf6, 0x24, 0x4d, 0x7d, 0x0a, 0x4d, 0x43, 0x7e, 0x27, 0x0c, 0xba, 0xc4,
	0x6d, 0x92, 0xa6, 0x9e, 0x45, 0x33, 0x00, 0x92, 0x98, 0x6d, 0x99, 0x63, 0xef, 0xeb, 0xd8, 0xb1,
	0x42, 0x07, 0x33, 0x78, 0x9e, 0xad, 0xb4, 0xf1, 0xb8, 0x6b, 0xfb, 0xa4, 0xa9, 0x03, 0x42, 0x30,
	0x23, 0xc1, 0x72, 0x7b, 0xbd, 0x80, 0xf2, 0x30, 0xc9, 0xbf, 0x2b, 0xe9, 0x45, 0x34, 0x2b, 0xaf,
	0x7e, 0xd1, 0x4f, 0xd1, 0xa7, 0xd9, 0x62, 0x3b, 0x84, 0x52, 0x87, 0x7f, 0x5c, 0xd5, 0x67, 0x96,
	0x3f, 0xc8, 0xc0, 0x6c, 0x9f, 0x3c, 0x18, 0xfd, 0x03, 0x7a, 0x40, 0x7c, 0x3d, 0xc5, 0x18, 0xd9,
	0x73, 0x0f, 0x5d, 0xef, 0x91, 0xbb, 0x73, 0xd4, 0xd9, 0xf7, 0x1c, 0x5d, 0x43, 0x97, 0x61, 0x2e,
	0x1a, 0x92, 0x9f, 0x00, 0xea, 0x35, 0x7d, 0x02, 0x19, 0x50, 0xd9, 0x73, 0x83, 0xb0, 0xdb, 0xf5,
	0x7c, 0x4a, 0x9a, 0xc2, 0x81, 0x0e, 0xb0, 0x8f, 0x2d, 0x4a, 0x7c, 0xfe, 0xeb, 0x94, 0x9e, 0x66,
	0xa4, 0x75, 0xd7, 0xf2, 0x7c, 0x9f, 0x58, 0x34, 0x6a, 0x97, 0xea, 0x19, 0xc6, 0xc3, 0x86, 0xfc,
	0xb7, 0x72, 0xdd, 0xf1, 0x02, 0x2e, 0x21, 0x04, 0x33, 0xb5, 0xb0, 0xeb, 0xd8, 0x16, 0xa6, 0x84,
	0x2f, 0xa6, 0x4f, 0xb1, 0xb1, 0xba, 0xfb, 0x0e, 0x76, 0xec, 0xa6, 0x34, 0x4d, 0x3d, 0x8b, 0xe6,
	0x61, 0x76, 0xd7, 0xf3, 0xb6, 0x30, 0x25, 0xbb, 0x9e, 0x94, 0x75, 0x0e, 0xe9, 0x50, 0x94, 0x47,
	0x14, 0xa4, 0x79, 0x54, 0x82, 0x4b, 0x62, 0x76, 0xcd, 0xf1, 0x09, 0x6e, 0x1e, 0x49, 0x99, 0xe9,
	0x80, 0x2e, 0x81, 0x5e, 0xb3, 0x5b, 0x2d, 0xe2, 0x13, 0x97, 0x0a, 0x1e, 0x03, 0xbd, 0xa0, 0x6c,
	0x25, 0xbd, 0x48, 0x2f, 0x32, 0x64, 0x74, 0xcc, 0xb5, 0xed, 0xfa, 0x86, 0xef, 0x7b, 0xbe, 0x3e,
	0xcd, 0xf6, 0x92, 0x48, 0xb1, 0xd7, 0x0c, 0xe3, 0xd2, 0xc4, 0x94, 0x6c, 0xd9, 0x1d, 0x9b, 0x6e,
	0x3c, 0xb6, 0x08, 0x61, 0x6a, 0x9d, 0x45, 0xcf, 0xc0, 0x82, 0x22, 0xa0, 0x1d, 0x51, 0x31, 0x77,
	0x99, 0xd8, 0x75, 0x9d, 0xd1, 0xbc, 0x65, 0x07, 0x81, 0xed, 0xb6, 0x7b, 0x01, 0x4c, 0x9f, 0x63,
	0x96, 0xf1, 0xe6, 0xee, 0xee, 0xb6, 0xd8, 0x0b, 0xa1, 0x05, 0x98, 0xbf, 0xcf, 0xbf, 0x7b, 0x30,
	0x45, 0xe3, 0x7d, 0x47, 0x4a, 0x66, 0x1e, 0x5d, 0x01, 0x74, 0xdf, 0x73, 0x05, 0x87, 0xbd, 0xf1,
	0x4b, 0x42, 0x57, 0xf1, 0x9e, 0xc2, 0x3f, 0xf5, 0xcb, 0x4c, 0x1a, 0xfd, 0xba, 0x62, 0x56, 0xae,
	0x5f, 0x41, 0xcf, 0xc2, 0x33, 0x03, 0x33, 0xbd, 0x8a, 0x4d, 0x5f, 0x60, 0x3b, 0x29, 0x80, 0x48,
	0x38, 0x25, 0x26, 0x30, 0x29, 0xf2, 0x48, 0x37, 0x57, 0x97, 0xeb, 0x90, 0x95, 0xfd, 0x67, 0x66,
	0x76, 0x26, 0x69, 0x33, 0x13, 0xf6, 0xfc, 0x23, 0x3d, 0xc5, 0xfc, 0x69, 0x17, 0x3f, 0xd6, 0x35,
	0xa6, 0xbf, 0x2d, 0xcf, 0xc2, 0xce, 0xba, 0xd7, 0xe9, 0x30, 0x01, 0x78, 0xae, 0x3e, 0xc1, 0x64,
	0x1a, 0x49, 0x7a, 0x93, 0x90, 0x40, 0x4f, 0x2f, 0xbf, 0x06, 0xb9, 0xa8, 0x37, 0xcd, 0x9c, 0x67,
	0x6d, 0x3f, 0xf0, 0x9c, 0x90, 0x12, 0x3d, 0xc5, 0xbc, 0x61, 0x9b, 0xf8, 0x7b, 0xae, 0x4d, 0x75,
	0x4d, 0xb8, 0x8e, 0x6f, 0x11, 0x97, 0xe2, 0x36, 0xd1, 0x27, 0x96, 0x6b, 0x70, 0x69, 0x58, 0xa5,
	0x87, 0x16, 0xa1, 0xb4, 0x8d, 0x7d, 0x6a, 0x5b, 0x76, 0x17, 0x53, 0x52, 0xf3, 0xee, 0x7b, 0xb4,
	0xee, 0xda, 0xd4, 0xc6, 0x7c, 0x49, 0x7e, 0xd8, 0x66, 0x68, 0x91, 0x07, 0xae, 0x73, 0xa4, 0x6b,
	0xcb, 0x2d, 0x28, 0xaa, 0xbd, 0x68, 0x36, 0x2f, 0x7f, 0xc4, 0xd9, 0x24, 0x0c, 0x3f, 0x1f, 0x77,
	0x5b, 0x76, 0x7d, 0xec, 0x06, 0x2d, 0xe2, 0xeb, 0x1a, 0x3b, 0x57, 0x8d, 0xf0, 0x96, 0x8f, 0x88,
	0x0e, 0x5f, 0xb7, 0xe9, 0x41, 0xd3, 0xc7, 0x8f, 0xf4, 0x34, 0xf7, 0x2f, 0xae, 0x20, 0x1e, 0x15,
	0xa4, 0xd8, 0xf4, 0xc9, 0xea, 0x8f, 0x73, 0x90, 0xdb, 0x90, 0x7f, 0x37, 0xa3, 0x3d, 0x80, 0xde,
	0x3f, 0x57, 0xe8, 0x99, 0x5e, 0xf4, 0x1a, 0xf8, 0x39, 0xb0, 0xfc, 0xdf, 0xc3, 0x26, 0x07, 0xff,
	0xa0, 0x30, 0x52, 0xaf, 0x68, 0xe8, 0x6d, 0x28, 0x28, 0xbf, 0x59, 0xa0, 0xc5, 0x1e, 0xe9, 0xe0,
	0x2f, 0x42, 0xe5, 0x17, 0x87, 0xce, 0x9e, 0xb0, 0x32, 0x06, 0x34, 0xf8, 0x95, 0x17, 0x3d, 0xaf,
	0xf4, 0x8a, 0x4f, 0xfa, 0xe6, 0x5c, 0xfe, 0xaf, 0xd3, 0x41, 0x22, 0xee, 0x1b, 0x29, 0xb4, 0xc1,
	0x82, 0x57, 0xfc, 0xb1, 0xf1, 0xda, 0x20, 0x95, 0xf2, 0xd1, 0xb5, 0x7c, 0x65, 0xf8, 0xb4, 0x91,
	0x42, 0x77, 0x60, 0xea, 0x81, 0xe8, 0x94, 0x2e, 0x0e, 0x6d, 0x4c, 0x47, 0x2b, 0xcc, 0xf7, 0xcd,
	0x4a, 0xf2, 0x35, 0xc8, 0xc7, 0x7d, 0x51, 0xa4, 0xb4, 0xb6, 0xfb, 0x9b, 0xa5, 0xea, 0x09, 0xd4,
	0x66, 0xa3, 0x91, 0x42, 0x77, 0x21, 0x17, 0xf5, 0xe7, 0x90, 0x72, 0x31, 0xf5, 0xf5, 0xec, 0xca,
	0x97, 0x07, 0xa6, 0x24, 0xfd, 0x1e, 0xcc, 0x24, 0x9b, 0x00, 0xe8, 0xd9, 0xc1, 0x9a, 0x36, 0xd1,
	0xab, 0x28, 0x2f, 0x9d, 0x0c, 0x88, 0xe5, 0xfb, 0x00, 0x8a, 0x6a, 0xed, 0xab, 0x4a, 0x78, 0x48,
	0x81, 0x5e, 0xae, 0x9c, 0x34, 0xad, 0x2e, 0xa8, 0x56, 0x2c, 0xea, 0x82, 0x43, 0xaa, 0x31, 0x75,
	0xc1, 0x61, 0x85, 0x8e, 0x91, 0x42, 0xdf, 0x90, 0x1f, 0x62, 0x94, 0x3c, 0x1b, 0x3d, 0xd7, 0x47,
	0x35, 0x58, 0x4d, 0x94, 0x8d, 0xd3, 0x20, 0xf1, 0xe2, 0x5b, 0xf2, 0xb2, 0x14, 0x71, 0x74, 0xc0,
	0x38, 0x12, 0x59, 0x55, 0xf9, 0xda, 0x09, 0xb3, 0xf1, 0x6a, 0x6f, 0xc3, 0x6c, 0x5f, 0xf6, 0x83,
	0x96, 0xfa, 0x68, 0x06, 0x72, 0xb5, 0xf2, 0x73, 0xa7, 0x20, 0xa2, 0x95, 0xef, 0xbd, 0xfa, 0xe4,
	0xa3, 0x4a, 0xea, 0xfd, 0x8f, 0x2a, 0xa9, 0x4f, 0x3f, 0xaa, 0x68, 0xdf, 0x3e, 0xae, 0x68, 0x3f,
	0x3b, 0xae, 0x68, 0xbf, 0x3f, 0xae, 0x68, 0x4f, 0x8e, 0x2b, 0xda, 0x87, 0xc7, 0x15, 0xed, 0x93,
	0xe3, 0x4a, 0xea, 0xd3, 0xe3, 0x8a, 0xf6, 0xde, 0xc7, 0x95, 0xd4, 0x93, 0x8f, 0x2b, 0xa9, 0xf7,
	0x3f, 0xae, 0xa4, 0xf6, 0xa7, 0x78, 0x6e, 0x7d, 0xeb, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xef,
	0xc3, 0x28, 0xe9, 0x41, 0x31, 0x00, 0x00,
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
func (x UpdateReason) String() string {
	s, ok := UpdateReason_name[int32(x)]
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
func (this *MarketStatisticsRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MarketStatisticsRequest)
	if !ok {
		that2, ok := that.(MarketStatisticsRequest)
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
	if len(this.Statistics) != len(that1.Statistics) {
		return false
	}
	for i := range this.Statistics {
		if this.Statistics[i] != that1.Statistics[i] {
			return false
		}
	}
	return true
}
func (this *MarketStatisticsResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MarketStatisticsResponse)
	if !ok {
		that2, ok := that.(MarketStatisticsResponse)
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
	if len(this.Statistics) != len(that1.Statistics) {
		return false
	}
	for i := range this.Statistics {
		if !this.Statistics[i].Equal(that1.Statistics[i]) {
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
	if len(this.Securities) != len(that1.Securities) {
		return false
	}
	for i := range this.Securities {
		if !this.Securities[i].Equal(that1.Securities[i]) {
			return false
		}
	}
	if len(this.Orders) != len(that1.Orders) {
		return false
	}
	for i := range this.Orders {
		if !this.Orders[i].Equal(that1.Orders[i]) {
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
	if len(this.Balances) != len(that1.Balances) {
		return false
	}
	for i := range this.Balances {
		if !this.Balances[i].Equal(that1.Balances[i]) {
			return false
		}
	}
	if !this.MakerFee.Equal(that1.MakerFee) {
		return false
	}
	if !this.TakerFee.Equal(that1.TakerFee) {
		return false
	}
	if this.Success != that1.Success {
		return false
	}
	if this.RejectionReason != that1.RejectionReason {
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
func (this *TradeCaptureReportRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*TradeCaptureReportRequest)
	if !ok {
		that2, ok := that.(TradeCaptureReportRequest)
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
func (this *TradeCaptureReportFilter) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*TradeCaptureReportFilter)
	if !ok {
		that2, ok := that.(TradeCaptureReportFilter)
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
	if !this.From.Equal(that1.From) {
		return false
	}
	if !this.To.Equal(that1.To) {
		return false
	}
	return true
}
func (this *TradeCaptureReport) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*TradeCaptureReport)
	if !ok {
		that2, ok := that.(TradeCaptureReport)
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
	if len(this.Trades) != len(that1.Trades) {
		return false
	}
	for i := range this.Trades {
		if !this.Trades[i].Equal(that1.Trades[i]) {
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
func (this *AccountUpdate) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AccountUpdate)
	if !ok {
		that2, ok := that.(AccountUpdate)
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
	if this.UpdateReason != that1.UpdateReason {
		return false
	}
	if !this.Asset.Equal(that1.Asset) {
		return false
	}
	if this.Balance != that1.Balance {
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
	if this.OrderID != that1.OrderID {
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
func (this *MarketStatisticsRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.MarketStatisticsRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	s = append(s, "Statistics: "+fmt.Sprintf("%#v", this.Statistics)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *MarketStatisticsResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.MarketStatisticsResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Statistics != nil {
		s = append(s, "Statistics: "+fmt.Sprintf("%#v", this.Statistics)+",\n")
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
	s := make([]string, 0, 15)
	s = append(s, "&messages.AccountDataResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Securities != nil {
		s = append(s, "Securities: "+fmt.Sprintf("%#v", this.Securities)+",\n")
	}
	if this.Orders != nil {
		s = append(s, "Orders: "+fmt.Sprintf("%#v", this.Orders)+",\n")
	}
	if this.Positions != nil {
		s = append(s, "Positions: "+fmt.Sprintf("%#v", this.Positions)+",\n")
	}
	if this.Balances != nil {
		s = append(s, "Balances: "+fmt.Sprintf("%#v", this.Balances)+",\n")
	}
	if this.MakerFee != nil {
		s = append(s, "MakerFee: "+fmt.Sprintf("%#v", this.MakerFee)+",\n")
	}
	if this.TakerFee != nil {
		s = append(s, "TakerFee: "+fmt.Sprintf("%#v", this.TakerFee)+",\n")
	}
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
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
func (this *TradeCaptureReportRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.TradeCaptureReportRequest{")
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
func (this *TradeCaptureReportFilter) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 10)
	s = append(s, "&messages.TradeCaptureReportFilter{")
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
	if this.From != nil {
		s = append(s, "From: "+fmt.Sprintf("%#v", this.From)+",\n")
	}
	if this.To != nil {
		s = append(s, "To: "+fmt.Sprintf("%#v", this.To)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *TradeCaptureReport) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.TradeCaptureReport{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Trades != nil {
		s = append(s, "Trades: "+fmt.Sprintf("%#v", this.Trades)+",\n")
	}
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
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
func (this *AccountUpdate) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.AccountUpdate{")
	s = append(s, "UpdateReason: "+fmt.Sprintf("%#v", this.UpdateReason)+",\n")
	if this.Asset != nil {
		s = append(s, "Asset: "+fmt.Sprintf("%#v", this.Asset)+",\n")
	}
	s = append(s, "Balance: "+fmt.Sprintf("%#v", this.Balance)+",\n")
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
	s := make([]string, 0, 9)
	s = append(s, "&messages.OrderReplaceResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	s = append(s, "OrderID: "+fmt.Sprintf("%#v", this.OrderID)+",\n")
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

// UnimplementedExecutorServer can be embedded to have forward compatible implementations.
type UnimplementedExecutorServer struct {
}

func (*UnimplementedExecutorServer) MarketData(req *MarketDataRequest, srv Executor_MarketDataServer) error {
	return status.Errorf(codes.Unimplemented, "method MarketData not implemented")
}
func (*UnimplementedExecutorServer) AccountData(req *AccountDataRequest, srv Executor_AccountDataServer) error {
	return status.Errorf(codes.Unimplemented, "method AccountData not implemented")
}
func (*UnimplementedExecutorServer) SecurityDefinition(ctx context.Context, req *SecurityDefinitionRequest) (*SecurityDefinitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SecurityDefinition not implemented")
}
func (*UnimplementedExecutorServer) Securities(ctx context.Context, req *SecurityListRequest) (*SecurityList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Securities not implemented")
}
func (*UnimplementedExecutorServer) Orders(ctx context.Context, req *OrderStatusRequest) (*OrderList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Orders not implemented")
}
func (*UnimplementedExecutorServer) Positions(ctx context.Context, req *PositionsRequest) (*PositionList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Positions not implemented")
}
func (*UnimplementedExecutorServer) Balances(ctx context.Context, req *BalancesRequest) (*BalanceList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Balances not implemented")
}
func (*UnimplementedExecutorServer) NewOrderSingle(ctx context.Context, req *NewOrderSingleRequest) (*NewOrderSingleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewOrderSingle not implemented")
}
func (*UnimplementedExecutorServer) NewOrderBulk(ctx context.Context, req *NewOrderBulkRequest) (*NewOrderBulkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewOrderBulk not implemented")
}
func (*UnimplementedExecutorServer) OrderReplace(ctx context.Context, req *OrderReplaceRequest) (*OrderReplaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderReplace not implemented")
}
func (*UnimplementedExecutorServer) OrderBulkReplace(ctx context.Context, req *OrderBulkReplaceRequest) (*OrderBulkReplaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderBulkReplace not implemented")
}
func (*UnimplementedExecutorServer) OrderCancel(ctx context.Context, req *OrderCancelRequest) (*OrderCancelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderCancel not implemented")
}
func (*UnimplementedExecutorServer) OrderMassCancel(ctx context.Context, req *OrderMassCancelRequest) (*OrderMassCancelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderMassCancel not implemented")
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
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HistoricalLiquidationsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HistoricalLiquidationsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.To != nil {
		{
			size, err := m.To.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.From != nil {
		{
			size, err := m.From.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Instrument != nil {
		{
			size, err := m.Instrument.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *HistoricalLiquidationsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HistoricalLiquidationsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HistoricalLiquidationsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x28
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if len(m.Liquidations) > 0 {
		for iNdEx := len(m.Liquidations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Liquidations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *MarketStatisticsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MarketStatisticsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MarketStatisticsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Statistics) > 0 {
		dAtA5 := make([]byte, len(m.Statistics)*10)
		var j4 int
		for _, num := range m.Statistics {
			for num >= 1<<7 {
				dAtA5[j4] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j4++
			}
			dAtA5[j4] = uint8(num)
			j4++
		}
		i -= j4
		copy(dAtA[i:], dAtA5[:j4])
		i = encodeVarintExecutorMessages(dAtA, i, uint64(j4))
		i--
		dAtA[i] = 0x1a
	}
	if m.Instrument != nil {
		{
			size, err := m.Instrument.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *MarketStatisticsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MarketStatisticsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MarketStatisticsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x28
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if len(m.Statistics) > 0 {
		for iNdEx := len(m.Statistics) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Statistics[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *MarketDataRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MarketDataRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MarketDataRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Aggregation != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Aggregation))
		i--
		dAtA[i] = 0x28
	}
	if m.Instrument != nil {
		{
			size, err := m.Instrument.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Subscriber != nil {
		{
			size, err := m.Subscriber.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Subscribe {
		i--
		if m.Subscribe {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *MarketDataResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MarketDataResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MarketDataResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x48
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x40
	}
	if m.SeqNum != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SeqNum))
		i--
		dAtA[i] = 0x38
	}
	if len(m.Trades) > 0 {
		for iNdEx := len(m.Trades) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Trades[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if m.SnapshotL3 != nil {
		{
			size, err := m.SnapshotL3.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.SnapshotL2 != nil {
		{
			size, err := m.SnapshotL2.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.SnapshotL1 != nil {
		{
			size, err := m.SnapshotL1.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *MarketDataIncrementalRefresh) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MarketDataIncrementalRefresh) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MarketDataIncrementalRefresh) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Stats) > 0 {
		for iNdEx := len(m.Stats) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Stats[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x52
		}
	}
	if m.Funding != nil {
		{
			size, err := m.Funding.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x4a
	}
	if m.Liquidation != nil {
		{
			size, err := m.Liquidation.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x42
	}
	if len(m.Trades) > 0 {
		for iNdEx := len(m.Trades) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Trades[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if m.UpdateL3 != nil {
		{
			size, err := m.UpdateL3.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.UpdateL2 != nil {
		{
			size, err := m.UpdateL2.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.UpdateL1 != nil {
		{
			size, err := m.UpdateL1.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.SeqNum != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SeqNum))
		i--
		dAtA[i] = 0x18
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *AccountDataRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountDataRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AccountDataRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Account != nil {
		{
			size, err := m.Account.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Subscriber != nil {
		{
			size, err := m.Subscriber.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Subscribe {
		i--
		if m.Subscribe {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *AccountDataResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountDataResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AccountDataResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SeqNum != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SeqNum))
		i--
		dAtA[i] = 0x58
	}
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x50
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x48
	}
	if m.TakerFee != nil {
		{
			size, err := m.TakerFee.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x42
	}
	if m.MakerFee != nil {
		{
			size, err := m.MakerFee.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x3a
	}
	if len(m.Balances) > 0 {
		for iNdEx := len(m.Balances) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Balances[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.Positions) > 0 {
		for iNdEx := len(m.Positions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Positions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Orders) > 0 {
		for iNdEx := len(m.Orders) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Orders[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Securities) > 0 {
		for iNdEx := len(m.Securities) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Securities[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *AccountDataIncrementalRefresh) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountDataIncrementalRefresh) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AccountDataIncrementalRefresh) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Report != nil {
		{
			size, err := m.Report.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TradeCaptureReportRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TradeCaptureReportRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TradeCaptureReportRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Filter != nil {
		{
			size, err := m.Filter.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Account != nil {
		{
			size, err := m.Account.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TradeCaptureReportFilter) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TradeCaptureReportFilter) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TradeCaptureReportFilter) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.To != nil {
		{
			size, err := m.To.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.From != nil {
		{
			size, err := m.From.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.Side != nil {
		{
			size, err := m.Side.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Instrument != nil {
		{
			size, err := m.Instrument.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.ClientOrderID != nil {
		{
			size, err := m.ClientOrderID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.OrderID != nil {
		{
			size, err := m.OrderID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TradeCaptureReport) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TradeCaptureReport) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TradeCaptureReport) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x28
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if len(m.Trades) > 0 {
		for iNdEx := len(m.Trades) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Trades[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SecurityDefinitionRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SecurityDefinitionRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SecurityDefinitionRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Instrument != nil {
		{
			size, err := m.Instrument.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SecurityDefinitionResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SecurityDefinitionResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SecurityDefinitionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x28
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if m.Security != nil {
		{
			size, err := m.Security.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SecurityListRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SecurityListRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SecurityListRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Subscriber != nil {
		{
			size, err := m.Subscriber.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Subscribe {
		i--
		if m.Subscribe {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SecurityList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SecurityList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SecurityList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x28
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if len(m.Securities) > 0 {
		for iNdEx := len(m.Securities) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Securities[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ExecutionReport) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExecutionReport) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExecutionReport) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x98
	}
	if m.FeeBasis != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.FeeBasis))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x90
	}
	if m.FeeType != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.FeeType))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x88
	}
	if m.FeeCurrency != nil {
		{
			size, err := m.FeeCurrency.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x82
	}
	if m.FeeAmount != nil {
		{
			size, err := m.FeeAmount.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x7a
	}
	if m.FillQuantity != nil {
		{
			size, err := m.FillQuantity.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x72
	}
	if m.FillPrice != nil {
		{
			size, err := m.FillPrice.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x6a
	}
	if m.TradeID != nil {
		{
			size, err := m.TradeID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x62
	}
	if m.TransactionTime != nil {
		{
			size, err := m.TransactionTime.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x5a
	}
	if m.CumQuantity != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.CumQuantity))))
		i--
		dAtA[i] = 0x51
	}
	if m.LeavesQuantity != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.LeavesQuantity))))
		i--
		dAtA[i] = 0x49
	}
	if m.Instrument != nil {
		{
			size, err := m.Instrument.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x3a
	}
	if m.OrderStatus != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderStatus))
		i--
		dAtA[i] = 0x30
	}
	if m.ExecutionType != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ExecutionType))
		i--
		dAtA[i] = 0x28
	}
	if len(m.ExecutionID) > 0 {
		i -= len(m.ExecutionID)
		copy(dAtA[i:], m.ExecutionID)
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.ExecutionID)))
		i--
		dAtA[i] = 0x22
	}
	if m.ClientOrderID != nil {
		{
			size, err := m.ClientOrderID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.OrderID) > 0 {
		i -= len(m.OrderID)
		copy(dAtA[i:], m.OrderID)
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.OrderID)))
		i--
		dAtA[i] = 0x12
	}
	if m.SeqNum != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SeqNum))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *AccountUpdate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountUpdate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AccountUpdate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Balance != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Balance))))
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
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.UpdateReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.UpdateReason))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SideValue) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SideValue) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SideValue) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Value != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Value))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *OrderStatusValue) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderStatusValue) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderStatusValue) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Value != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Value))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *OrderFilter) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderFilter) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderFilter) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.OrderStatus != nil {
		{
			size, err := m.OrderStatus.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.Side != nil {
		{
			size, err := m.Side.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Instrument != nil {
		{
			size, err := m.Instrument.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.ClientOrderID != nil {
		{
			size, err := m.ClientOrderID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.OrderID != nil {
		{
			size, err := m.OrderID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *OrderStatusRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderStatusRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderStatusRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Filter != nil {
		{
			size, err := m.Filter.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.Account != nil {
		{
			size, err := m.Account.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Subscriber != nil {
		{
			size, err := m.Subscriber.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Subscribe {
		i--
		if m.Subscribe {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *OrderList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Orders) > 0 {
		for iNdEx := len(m.Orders) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Orders[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *PositionsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PositionsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PositionsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Account != nil {
		{
			size, err := m.Account.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.Instrument != nil {
		{
			size, err := m.Instrument.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Subscriber != nil {
		{
			size, err := m.Subscriber.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Subscribe {
		i--
		if m.Subscribe {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *PositionList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PositionList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PositionList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x28
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if len(m.Positions) > 0 {
		for iNdEx := len(m.Positions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Positions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *BalancesRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BalancesRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BalancesRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Account != nil {
		{
			size, err := m.Account.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.Asset != nil {
		{
			size, err := m.Asset.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Subscriber != nil {
		{
			size, err := m.Subscriber.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Subscribe {
		i--
		if m.Subscribe {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *BalanceList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BalanceList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BalanceList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x28
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if len(m.Balances) > 0 {
		for iNdEx := len(m.Balances) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Balances[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *NewOrder) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NewOrder) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NewOrder) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ExecutionInstructions) > 0 {
		dAtA57 := make([]byte, len(m.ExecutionInstructions)*10)
		var j56 int
		for _, num := range m.ExecutionInstructions {
			for num >= 1<<7 {
				dAtA57[j56] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j56++
			}
			dAtA57[j56] = uint8(num)
			j56++
		}
		i -= j56
		copy(dAtA[i:], dAtA57[:j56])
		i = encodeVarintExecutorMessages(dAtA, i, uint64(j56))
		i--
		dAtA[i] = 0x52
	}
	if m.Price != nil {
		{
			size, err := m.Price.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x4a
	}
	if m.Quantity != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Quantity))))
		i--
		dAtA[i] = 0x41
	}
	if m.TimeInForce != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.TimeInForce))
		i--
		dAtA[i] = 0x38
	}
	if m.OrderSide != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderSide))
		i--
		dAtA[i] = 0x30
	}
	if m.OrderType != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.OrderType))
		i--
		dAtA[i] = 0x28
	}
	if m.Instrument != nil {
		{
			size, err := m.Instrument.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.ClientOrderID) > 0 {
		i -= len(m.ClientOrderID)
		copy(dAtA[i:], m.ClientOrderID)
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.ClientOrderID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *NewOrderSingleRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NewOrderSingleRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NewOrderSingleRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Order != nil {
		{
			size, err := m.Order.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Account != nil {
		{
			size, err := m.Account.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *NewOrderSingleResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NewOrderSingleResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NewOrderSingleResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x28
	}
	if len(m.OrderID) > 0 {
		i -= len(m.OrderID)
		copy(dAtA[i:], m.OrderID)
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.OrderID)))
		i--
		dAtA[i] = 0x22
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *NewOrderBulkRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NewOrderBulkRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NewOrderBulkRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Orders) > 0 {
		for iNdEx := len(m.Orders) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Orders[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.Account != nil {
		{
			size, err := m.Account.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *NewOrderBulkResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NewOrderBulkResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NewOrderBulkResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x28
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if len(m.OrderIDs) > 0 {
		for iNdEx := len(m.OrderIDs) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.OrderIDs[iNdEx])
			copy(dAtA[i:], m.OrderIDs[iNdEx])
			i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.OrderIDs[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *OrderUpdate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderUpdate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderUpdate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Quantity != nil {
		{
			size, err := m.Quantity.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.OrigClientOrderID != nil {
		{
			size, err := m.OrigClientOrderID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.OrderID != nil {
		{
			size, err := m.OrderID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *OrderReplaceRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderReplaceRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderReplaceRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Update != nil {
		{
			size, err := m.Update.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Account != nil {
		{
			size, err := m.Account.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Instrument != nil {
		{
			size, err := m.Instrument.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *OrderReplaceResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderReplaceResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderReplaceResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x28
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if len(m.OrderID) > 0 {
		i -= len(m.OrderID)
		copy(dAtA[i:], m.OrderID)
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.OrderID)))
		i--
		dAtA[i] = 0x1a
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *OrderBulkReplaceRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderBulkReplaceRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderBulkReplaceRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Updates) > 0 {
		for iNdEx := len(m.Updates) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Updates[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.Account != nil {
		{
			size, err := m.Account.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Instrument != nil {
		{
			size, err := m.Instrument.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *OrderBulkReplaceResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderBulkReplaceResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderBulkReplaceResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x20
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *OrderCancelRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderCancelRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderCancelRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Account != nil {
		{
			size, err := m.Account.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.Instrument != nil {
		{
			size, err := m.Instrument.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.ClientOrderID != nil {
		{
			size, err := m.ClientOrderID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.OrderID != nil {
		{
			size, err := m.OrderID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *OrderCancelResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderCancelResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderCancelResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x20
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *OrderMassCancelRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderMassCancelRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderMassCancelRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Filter != nil {
		{
			size, err := m.Filter.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Account != nil {
		{
			size, err := m.Account.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *OrderMassCancelResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderMassCancelResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderMassCancelResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x20
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if m.ResponseID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ResponseID))
		i--
		dAtA[i] = 0x10
	}
	if m.RequestID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RequestID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintExecutorMessages(dAtA []byte, offset int, v uint64) int {
	offset -= sovExecutorMessages(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
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

func (m *MarketStatisticsRequest) Size() (n int) {
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
	if len(m.Statistics) > 0 {
		l = 0
		for _, e := range m.Statistics {
			l += sovExecutorMessages(uint64(e))
		}
		n += 1 + sovExecutorMessages(uint64(l)) + l
	}
	return n
}

func (m *MarketStatisticsResponse) Size() (n int) {
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
	if len(m.Statistics) > 0 {
		for _, e := range m.Statistics {
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
	if len(m.Securities) > 0 {
		for _, e := range m.Securities {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	if len(m.Orders) > 0 {
		for _, e := range m.Orders {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	if len(m.Positions) > 0 {
		for _, e := range m.Positions {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	if len(m.Balances) > 0 {
		for _, e := range m.Balances {
			l = e.Size()
			n += 1 + l + sovExecutorMessages(uint64(l))
		}
	}
	if m.MakerFee != nil {
		l = m.MakerFee.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.TakerFee != nil {
		l = m.TakerFee.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Success {
		n += 2
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
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

func (m *TradeCaptureReportRequest) Size() (n int) {
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

func (m *TradeCaptureReportFilter) Size() (n int) {
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

func (m *TradeCaptureReport) Size() (n int) {
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
	if len(m.Trades) > 0 {
		for _, e := range m.Trades {
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

func (m *AccountUpdate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.UpdateReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.UpdateReason))
	}
	if m.Asset != nil {
		l = m.Asset.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Balance != 0 {
		n += 9
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
	l = len(m.OrderID)
	if l > 0 {
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
	return (math_bits.Len64(x|1) + 6) / 7
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
	repeatedStringForLiquidations := "[]*Liquidation{"
	for _, f := range this.Liquidations {
		repeatedStringForLiquidations += strings.Replace(fmt.Sprintf("%v", f), "Liquidation", "models.Liquidation", 1) + ","
	}
	repeatedStringForLiquidations += "}"
	s := strings.Join([]string{`&HistoricalLiquidationsResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Liquidations:` + repeatedStringForLiquidations + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *MarketStatisticsRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&MarketStatisticsRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`Statistics:` + fmt.Sprintf("%v", this.Statistics) + `,`,
		`}`,
	}, "")
	return s
}
func (this *MarketStatisticsResponse) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForStatistics := "[]*Stat{"
	for _, f := range this.Statistics {
		repeatedStringForStatistics += strings.Replace(fmt.Sprintf("%v", f), "Stat", "models.Stat", 1) + ","
	}
	repeatedStringForStatistics += "}"
	s := strings.Join([]string{`&MarketStatisticsResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Statistics:` + repeatedStringForStatistics + `,`,
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
	repeatedStringForTrades := "[]*AggregatedTrade{"
	for _, f := range this.Trades {
		repeatedStringForTrades += strings.Replace(fmt.Sprintf("%v", f), "AggregatedTrade", "models.AggregatedTrade", 1) + ","
	}
	repeatedStringForTrades += "}"
	s := strings.Join([]string{`&MarketDataResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`SnapshotL1:` + strings.Replace(fmt.Sprintf("%v", this.SnapshotL1), "OBL1Snapshot", "models.OBL1Snapshot", 1) + `,`,
		`SnapshotL2:` + strings.Replace(fmt.Sprintf("%v", this.SnapshotL2), "OBL2Snapshot", "models.OBL2Snapshot", 1) + `,`,
		`SnapshotL3:` + strings.Replace(fmt.Sprintf("%v", this.SnapshotL3), "OBL3Snapshot", "models.OBL3Snapshot", 1) + `,`,
		`Trades:` + repeatedStringForTrades + `,`,
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
	repeatedStringForTrades := "[]*AggregatedTrade{"
	for _, f := range this.Trades {
		repeatedStringForTrades += strings.Replace(fmt.Sprintf("%v", f), "AggregatedTrade", "models.AggregatedTrade", 1) + ","
	}
	repeatedStringForTrades += "}"
	repeatedStringForStats := "[]*Stat{"
	for _, f := range this.Stats {
		repeatedStringForStats += strings.Replace(fmt.Sprintf("%v", f), "Stat", "models.Stat", 1) + ","
	}
	repeatedStringForStats += "}"
	s := strings.Join([]string{`&MarketDataIncrementalRefresh{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`SeqNum:` + fmt.Sprintf("%v", this.SeqNum) + `,`,
		`UpdateL1:` + strings.Replace(fmt.Sprintf("%v", this.UpdateL1), "OBL1Update", "models.OBL1Update", 1) + `,`,
		`UpdateL2:` + strings.Replace(fmt.Sprintf("%v", this.UpdateL2), "OBL2Update", "models.OBL2Update", 1) + `,`,
		`UpdateL3:` + strings.Replace(fmt.Sprintf("%v", this.UpdateL3), "OBL3Update", "models.OBL3Update", 1) + `,`,
		`Trades:` + repeatedStringForTrades + `,`,
		`Liquidation:` + strings.Replace(fmt.Sprintf("%v", this.Liquidation), "Liquidation", "models.Liquidation", 1) + `,`,
		`Funding:` + strings.Replace(fmt.Sprintf("%v", this.Funding), "Funding", "models.Funding", 1) + `,`,
		`Stats:` + repeatedStringForStats + `,`,
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
	repeatedStringForSecurities := "[]*Security{"
	for _, f := range this.Securities {
		repeatedStringForSecurities += strings.Replace(fmt.Sprintf("%v", f), "Security", "models.Security", 1) + ","
	}
	repeatedStringForSecurities += "}"
	repeatedStringForOrders := "[]*Order{"
	for _, f := range this.Orders {
		repeatedStringForOrders += strings.Replace(fmt.Sprintf("%v", f), "Order", "models.Order", 1) + ","
	}
	repeatedStringForOrders += "}"
	repeatedStringForPositions := "[]*Position{"
	for _, f := range this.Positions {
		repeatedStringForPositions += strings.Replace(fmt.Sprintf("%v", f), "Position", "models.Position", 1) + ","
	}
	repeatedStringForPositions += "}"
	repeatedStringForBalances := "[]*Balance{"
	for _, f := range this.Balances {
		repeatedStringForBalances += strings.Replace(fmt.Sprintf("%v", f), "Balance", "models.Balance", 1) + ","
	}
	repeatedStringForBalances += "}"
	s := strings.Join([]string{`&AccountDataResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Securities:` + repeatedStringForSecurities + `,`,
		`Orders:` + repeatedStringForOrders + `,`,
		`Positions:` + repeatedStringForPositions + `,`,
		`Balances:` + repeatedStringForBalances + `,`,
		`MakerFee:` + strings.Replace(fmt.Sprintf("%v", this.MakerFee), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`TakerFee:` + strings.Replace(fmt.Sprintf("%v", this.TakerFee), "DoubleValue", "types.DoubleValue", 1) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
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
		`Report:` + strings.Replace(this.Report.String(), "ExecutionReport", "ExecutionReport", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *TradeCaptureReportRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&TradeCaptureReportRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`Filter:` + strings.Replace(this.Filter.String(), "TradeCaptureReportFilter", "TradeCaptureReportFilter", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *TradeCaptureReportFilter) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&TradeCaptureReportFilter{`,
		`OrderID:` + strings.Replace(fmt.Sprintf("%v", this.OrderID), "StringValue", "types.StringValue", 1) + `,`,
		`ClientOrderID:` + strings.Replace(fmt.Sprintf("%v", this.ClientOrderID), "StringValue", "types.StringValue", 1) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`Side:` + strings.Replace(this.Side.String(), "SideValue", "SideValue", 1) + `,`,
		`From:` + strings.Replace(fmt.Sprintf("%v", this.From), "Timestamp", "types.Timestamp", 1) + `,`,
		`To:` + strings.Replace(fmt.Sprintf("%v", this.To), "Timestamp", "types.Timestamp", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *TradeCaptureReport) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForTrades := "[]*TradeCapture{"
	for _, f := range this.Trades {
		repeatedStringForTrades += strings.Replace(fmt.Sprintf("%v", f), "TradeCapture", "models.TradeCapture", 1) + ","
	}
	repeatedStringForTrades += "}"
	s := strings.Join([]string{`&TradeCaptureReport{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Trades:` + repeatedStringForTrades + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
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
	repeatedStringForSecurities := "[]*Security{"
	for _, f := range this.Securities {
		repeatedStringForSecurities += strings.Replace(fmt.Sprintf("%v", f), "Security", "models.Security", 1) + ","
	}
	repeatedStringForSecurities += "}"
	s := strings.Join([]string{`&SecurityList{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Securities:` + repeatedStringForSecurities + `,`,
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
func (this *AccountUpdate) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AccountUpdate{`,
		`UpdateReason:` + fmt.Sprintf("%v", this.UpdateReason) + `,`,
		`Asset:` + strings.Replace(fmt.Sprintf("%v", this.Asset), "Asset", "models1.Asset", 1) + `,`,
		`Balance:` + fmt.Sprintf("%v", this.Balance) + `,`,
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
		`Side:` + strings.Replace(this.Side.String(), "SideValue", "SideValue", 1) + `,`,
		`OrderStatus:` + strings.Replace(this.OrderStatus.String(), "OrderStatusValue", "OrderStatusValue", 1) + `,`,
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
		`Filter:` + strings.Replace(this.Filter.String(), "OrderFilter", "OrderFilter", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *OrderList) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForOrders := "[]*Order{"
	for _, f := range this.Orders {
		repeatedStringForOrders += strings.Replace(fmt.Sprintf("%v", f), "Order", "models.Order", 1) + ","
	}
	repeatedStringForOrders += "}"
	s := strings.Join([]string{`&OrderList{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`Orders:` + repeatedStringForOrders + `,`,
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
	repeatedStringForPositions := "[]*Position{"
	for _, f := range this.Positions {
		repeatedStringForPositions += strings.Replace(fmt.Sprintf("%v", f), "Position", "models.Position", 1) + ","
	}
	repeatedStringForPositions += "}"
	s := strings.Join([]string{`&PositionList{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Positions:` + repeatedStringForPositions + `,`,
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
	repeatedStringForBalances := "[]*Balance{"
	for _, f := range this.Balances {
		repeatedStringForBalances += strings.Replace(fmt.Sprintf("%v", f), "Balance", "models.Balance", 1) + ","
	}
	repeatedStringForBalances += "}"
	s := strings.Join([]string{`&BalanceList{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Balances:` + repeatedStringForBalances + `,`,
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
		`Order:` + strings.Replace(this.Order.String(), "NewOrder", "NewOrder", 1) + `,`,
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
	repeatedStringForOrders := "[]*NewOrder{"
	for _, f := range this.Orders {
		repeatedStringForOrders += strings.Replace(f.String(), "NewOrder", "NewOrder", 1) + ","
	}
	repeatedStringForOrders += "}"
	s := strings.Join([]string{`&NewOrderBulkRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`Orders:` + repeatedStringForOrders + `,`,
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
		`Update:` + strings.Replace(this.Update.String(), "OrderUpdate", "OrderUpdate", 1) + `,`,
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
		`OrderID:` + fmt.Sprintf("%v", this.OrderID) + `,`,
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
	repeatedStringForUpdates := "[]*OrderUpdate{"
	for _, f := range this.Updates {
		repeatedStringForUpdates += strings.Replace(f.String(), "OrderUpdate", "OrderUpdate", 1) + ","
	}
	repeatedStringForUpdates += "}"
	s := strings.Join([]string{`&OrderBulkReplaceRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`Updates:` + repeatedStringForUpdates + `,`,
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
		`Filter:` + strings.Replace(this.Filter.String(), "OrderFilter", "OrderFilter", 1) + `,`,
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *MarketStatisticsRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MarketStatisticsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MarketStatisticsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
			if wireType == 0 {
				var v models.StatType
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowExecutorMessages
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= models.StatType(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Statistics = append(m.Statistics, v)
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
				if elementCount != 0 && len(m.Statistics) == 0 {
					m.Statistics = make([]models.StatType, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v models.StatType
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowExecutorMessages
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= models.StatType(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Statistics = append(m.Statistics, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Statistics", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipExecutorMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *MarketStatisticsResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MarketStatisticsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MarketStatisticsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Statistics", wireType)
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
			m.Statistics = append(m.Statistics, &models.Stat{})
			if err := m.Statistics[len(m.Statistics)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
		case 5:
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
		case 6:
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
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MakerFee", wireType)
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
			if m.MakerFee == nil {
				m.MakerFee = &types.DoubleValue{}
			}
			if err := m.MakerFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TakerFee", wireType)
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
			if m.TakerFee == nil {
				m.TakerFee = &types.DoubleValue{}
			}
			if err := m.TakerFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
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
		case 10:
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
		case 11:
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *TradeCaptureReportRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: TradeCaptureReportRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TradeCaptureReportRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
				m.Filter = &TradeCaptureReportFilter{}
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *TradeCaptureReportFilter) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: TradeCaptureReportFilter: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TradeCaptureReportFilter: illegal tag %d (wire type %d)", fieldNum, wire)
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
		case 6:
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *TradeCaptureReport) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: TradeCaptureReport: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TradeCaptureReport: illegal tag %d (wire type %d)", fieldNum, wire)
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
			m.Trades = append(m.Trades, &models.TradeCapture{})
			if err := m.Trades[len(m.Trades)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *AccountUpdate) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: AccountUpdate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccountUpdate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpdateReason", wireType)
			}
			m.UpdateReason = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UpdateReason |= UpdateReason(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
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
		case 3:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Balance", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Balance = float64(math.Float64frombits(v))
		default:
			iNdEx = preIndex
			skippy, err := skipExecutorMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
	depth := 0
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
		case 1:
			iNdEx += 8
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
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupExecutorMessages
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthExecutorMessages
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthExecutorMessages        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowExecutorMessages          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupExecutorMessages = fmt.Errorf("proto: unexpected end of group")
)
