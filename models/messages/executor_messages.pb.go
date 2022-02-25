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
	UnknownExchange                RejectionReason = 26
	GraphQLError                   RejectionReason = 27
	ABIError                       RejectionReason = 28
	EthRPCError                    RejectionReason = 29
	UnknowProtocol                 RejectionReason = 30
	MissingProtocolAsset           RejectionReason = 31
	UnknownProtocolAsset           RejectionReason = 32
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
	26: "UnknownExchange",
	27: "GraphQLError",
	28: "ABIError",
	29: "EthRPCError",
	30: "UnknowProtocol",
	31: "MissingProtocolAsset",
	32: "UnknownProtocolAsset",
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
	"UnknownExchange":                26,
	"GraphQLError":                   27,
	"ABIError":                       28,
	"EthRPCError":                    29,
	"UnknowProtocol":                 30,
	"MissingProtocolAsset":           31,
	"UnknownProtocolAsset":           32,
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

type AccountMovementType int32

const (
	Unknown      AccountMovementType = 0
	FundingFee   AccountMovementType = 1
	Exchange     AccountMovementType = 2
	Deposit      AccountMovementType = 3
	Withdrawal   AccountMovementType = 4
	RealizedPnl  AccountMovementType = 5
	Commission   AccountMovementType = 6
	WelcomeBonus AccountMovementType = 7
)

var AccountMovementType_name = map[int32]string{
	0: "Unknown",
	1: "FundingFee",
	2: "Exchange",
	3: "Deposit",
	4: "Withdrawal",
	5: "RealizedPnl",
	6: "Commission",
	7: "WelcomeBonus",
}

var AccountMovementType_value = map[string]int32{
	"Unknown":      0,
	"FundingFee":   1,
	"Exchange":     2,
	"Deposit":      3,
	"Withdrawal":   4,
	"RealizedPnl":  5,
	"Commission":   6,
	"WelcomeBonus": 7,
}

func (AccountMovementType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{4}
}

type HistoricalOpenInterestsRequest struct {
	RequestID  uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Instrument *models.Instrument `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
	From       *types.Timestamp   `protobuf:"bytes,3,opt,name=from,proto3" json:"from,omitempty"`
	To         *types.Timestamp   `protobuf:"bytes,4,opt,name=to,proto3" json:"to,omitempty"`
}

func (m *HistoricalOpenInterestsRequest) Reset()      { *m = HistoricalOpenInterestsRequest{} }
func (*HistoricalOpenInterestsRequest) ProtoMessage() {}
func (*HistoricalOpenInterestsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{0}
}
func (m *HistoricalOpenInterestsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HistoricalOpenInterestsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HistoricalOpenInterestsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HistoricalOpenInterestsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HistoricalOpenInterestsRequest.Merge(m, src)
}
func (m *HistoricalOpenInterestsRequest) XXX_Size() int {
	return m.Size()
}
func (m *HistoricalOpenInterestsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HistoricalOpenInterestsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HistoricalOpenInterestsRequest proto.InternalMessageInfo

func (m *HistoricalOpenInterestsRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *HistoricalOpenInterestsRequest) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *HistoricalOpenInterestsRequest) GetFrom() *types.Timestamp {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *HistoricalOpenInterestsRequest) GetTo() *types.Timestamp {
	if m != nil {
		return m.To
	}
	return nil
}

type HistoricalOpenInterestsResponse struct {
	RequestID       uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64          `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Interests       []*models.Stat  `protobuf:"bytes,3,rep,name=interests,proto3" json:"interests,omitempty"`
	Success         bool            `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *HistoricalOpenInterestsResponse) Reset()      { *m = HistoricalOpenInterestsResponse{} }
func (*HistoricalOpenInterestsResponse) ProtoMessage() {}
func (*HistoricalOpenInterestsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{1}
}
func (m *HistoricalOpenInterestsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HistoricalOpenInterestsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HistoricalOpenInterestsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HistoricalOpenInterestsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HistoricalOpenInterestsResponse.Merge(m, src)
}
func (m *HistoricalOpenInterestsResponse) XXX_Size() int {
	return m.Size()
}
func (m *HistoricalOpenInterestsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HistoricalOpenInterestsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HistoricalOpenInterestsResponse proto.InternalMessageInfo

func (m *HistoricalOpenInterestsResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *HistoricalOpenInterestsResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *HistoricalOpenInterestsResponse) GetInterests() []*models.Stat {
	if m != nil {
		return m.Interests
	}
	return nil
}

func (m *HistoricalOpenInterestsResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *HistoricalOpenInterestsResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
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
	return fileDescriptor_350c53ba9303a7e6, []int{2}
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
	return fileDescriptor_350c53ba9303a7e6, []int{3}
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

type HistoricalUnipoolV3DataRequest struct {
	RequestID  uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Instrument *models.Instrument `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
	Start      uint64             `protobuf:"varint,3,opt,name=start,proto3" json:"start,omitempty"`
	End        uint64             `protobuf:"varint,4,opt,name=end,proto3" json:"end,omitempty"`
}

func (m *HistoricalUnipoolV3DataRequest) Reset()      { *m = HistoricalUnipoolV3DataRequest{} }
func (*HistoricalUnipoolV3DataRequest) ProtoMessage() {}
func (*HistoricalUnipoolV3DataRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{4}
}
func (m *HistoricalUnipoolV3DataRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HistoricalUnipoolV3DataRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HistoricalUnipoolV3DataRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HistoricalUnipoolV3DataRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HistoricalUnipoolV3DataRequest.Merge(m, src)
}
func (m *HistoricalUnipoolV3DataRequest) XXX_Size() int {
	return m.Size()
}
func (m *HistoricalUnipoolV3DataRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HistoricalUnipoolV3DataRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HistoricalUnipoolV3DataRequest proto.InternalMessageInfo

func (m *HistoricalUnipoolV3DataRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *HistoricalUnipoolV3DataRequest) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *HistoricalUnipoolV3DataRequest) GetStart() uint64 {
	if m != nil {
		return m.Start
	}
	return 0
}

func (m *HistoricalUnipoolV3DataRequest) GetEnd() uint64 {
	if m != nil {
		return m.End
	}
	return 0
}

type HistoricalUnipoolV3DataResponse struct {
	RequestID       uint64               `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64               `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Events          []*models.UPV3Update `protobuf:"bytes,3,rep,name=events,proto3" json:"events,omitempty"`
	SeqNum          uint64               `protobuf:"varint,4,opt,name=seq_num,json=seqNum,proto3" json:"seq_num,omitempty"`
	Success         bool                 `protobuf:"varint,5,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason      `protobuf:"varint,6,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *HistoricalUnipoolV3DataResponse) Reset()      { *m = HistoricalUnipoolV3DataResponse{} }
func (*HistoricalUnipoolV3DataResponse) ProtoMessage() {}
func (*HistoricalUnipoolV3DataResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{5}
}
func (m *HistoricalUnipoolV3DataResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HistoricalUnipoolV3DataResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HistoricalUnipoolV3DataResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HistoricalUnipoolV3DataResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HistoricalUnipoolV3DataResponse.Merge(m, src)
}
func (m *HistoricalUnipoolV3DataResponse) XXX_Size() int {
	return m.Size()
}
func (m *HistoricalUnipoolV3DataResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HistoricalUnipoolV3DataResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HistoricalUnipoolV3DataResponse proto.InternalMessageInfo

func (m *HistoricalUnipoolV3DataResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *HistoricalUnipoolV3DataResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *HistoricalUnipoolV3DataResponse) GetEvents() []*models.UPV3Update {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *HistoricalUnipoolV3DataResponse) GetSeqNum() uint64 {
	if m != nil {
		return m.SeqNum
	}
	return 0
}

func (m *HistoricalUnipoolV3DataResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *HistoricalUnipoolV3DataResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type HistoricalProtocolAssetTransferRequest struct {
	RequestID     uint64                `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ProtocolAsset *models.ProtocolAsset `protobuf:"bytes,2,opt,name=protocol_asset,json=protocolAsset,proto3" json:"protocol_asset,omitempty"`
	Start         uint64                `protobuf:"varint,3,opt,name=start,proto3" json:"start,omitempty"`
	Stop          uint64                `protobuf:"varint,4,opt,name=stop,proto3" json:"stop,omitempty"`
}

func (m *HistoricalProtocolAssetTransferRequest) Reset() {
	*m = HistoricalProtocolAssetTransferRequest{}
}
func (*HistoricalProtocolAssetTransferRequest) ProtoMessage() {}
func (*HistoricalProtocolAssetTransferRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{6}
}
func (m *HistoricalProtocolAssetTransferRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HistoricalProtocolAssetTransferRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HistoricalProtocolAssetTransferRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HistoricalProtocolAssetTransferRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HistoricalProtocolAssetTransferRequest.Merge(m, src)
}
func (m *HistoricalProtocolAssetTransferRequest) XXX_Size() int {
	return m.Size()
}
func (m *HistoricalProtocolAssetTransferRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HistoricalProtocolAssetTransferRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HistoricalProtocolAssetTransferRequest proto.InternalMessageInfo

func (m *HistoricalProtocolAssetTransferRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *HistoricalProtocolAssetTransferRequest) GetProtocolAsset() *models.ProtocolAsset {
	if m != nil {
		return m.ProtocolAsset
	}
	return nil
}

func (m *HistoricalProtocolAssetTransferRequest) GetStart() uint64 {
	if m != nil {
		return m.Start
	}
	return 0
}

func (m *HistoricalProtocolAssetTransferRequest) GetStop() uint64 {
	if m != nil {
		return m.Stop
	}
	return 0
}

type HistoricalProtocolAssetTransferResponse struct {
	RequestID       uint64                        `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64                        `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Update          []*models.ProtocolAssetUpdate `protobuf:"bytes,3,rep,name=update,proto3" json:"update,omitempty"`
	SeqNum          uint64                        `protobuf:"varint,4,opt,name=seq_num,json=seqNum,proto3" json:"seq_num,omitempty"`
	Success         bool                          `protobuf:"varint,5,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason               `protobuf:"varint,6,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *HistoricalProtocolAssetTransferResponse) Reset() {
	*m = HistoricalProtocolAssetTransferResponse{}
}
func (*HistoricalProtocolAssetTransferResponse) ProtoMessage() {}
func (*HistoricalProtocolAssetTransferResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{7}
}
func (m *HistoricalProtocolAssetTransferResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HistoricalProtocolAssetTransferResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HistoricalProtocolAssetTransferResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HistoricalProtocolAssetTransferResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HistoricalProtocolAssetTransferResponse.Merge(m, src)
}
func (m *HistoricalProtocolAssetTransferResponse) XXX_Size() int {
	return m.Size()
}
func (m *HistoricalProtocolAssetTransferResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HistoricalProtocolAssetTransferResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HistoricalProtocolAssetTransferResponse proto.InternalMessageInfo

func (m *HistoricalProtocolAssetTransferResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *HistoricalProtocolAssetTransferResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *HistoricalProtocolAssetTransferResponse) GetUpdate() []*models.ProtocolAssetUpdate {
	if m != nil {
		return m.Update
	}
	return nil
}

func (m *HistoricalProtocolAssetTransferResponse) GetSeqNum() uint64 {
	if m != nil {
		return m.SeqNum
	}
	return 0
}

func (m *HistoricalProtocolAssetTransferResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *HistoricalProtocolAssetTransferResponse) GetRejectionReason() RejectionReason {
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
	return fileDescriptor_350c53ba9303a7e6, []int{8}
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
	return fileDescriptor_350c53ba9303a7e6, []int{9}
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
	Stats       []models.StatType           `protobuf:"varint,6,rep,packed,name=stats,proto3,enum=models.StatType" json:"stats,omitempty"`
}

func (m *MarketDataRequest) Reset()      { *m = MarketDataRequest{} }
func (*MarketDataRequest) ProtoMessage() {}
func (*MarketDataRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{10}
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

func (m *MarketDataRequest) GetStats() []models.StatType {
	if m != nil {
		return m.Stats
	}
	return nil
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
	return fileDescriptor_350c53ba9303a7e6, []int{11}
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
	return fileDescriptor_350c53ba9303a7e6, []int{12}
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

type UnipoolV3DataRequest struct {
	RequestID  uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe  bool               `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
	Subscriber *actor.PID         `protobuf:"bytes,3,opt,name=subscriber,proto3" json:"subscriber,omitempty"`
	Instrument *models.Instrument `protobuf:"bytes,4,opt,name=instrument,proto3" json:"instrument,omitempty"`
}

func (m *UnipoolV3DataRequest) Reset()      { *m = UnipoolV3DataRequest{} }
func (*UnipoolV3DataRequest) ProtoMessage() {}
func (*UnipoolV3DataRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{13}
}
func (m *UnipoolV3DataRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UnipoolV3DataRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UnipoolV3DataRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UnipoolV3DataRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnipoolV3DataRequest.Merge(m, src)
}
func (m *UnipoolV3DataRequest) XXX_Size() int {
	return m.Size()
}
func (m *UnipoolV3DataRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UnipoolV3DataRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UnipoolV3DataRequest proto.InternalMessageInfo

func (m *UnipoolV3DataRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *UnipoolV3DataRequest) GetSubscribe() bool {
	if m != nil {
		return m.Subscribe
	}
	return false
}

func (m *UnipoolV3DataRequest) GetSubscriber() *actor.PID {
	if m != nil {
		return m.Subscriber
	}
	return nil
}

func (m *UnipoolV3DataRequest) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

type UnipoolV3DataResponse struct {
	RequestID       uint64               `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64               `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Update          []*models.UPV3Update `protobuf:"bytes,3,rep,name=update,proto3" json:"update,omitempty"`
	SeqNum          uint64               `protobuf:"varint,4,opt,name=seq_num,json=seqNum,proto3" json:"seq_num,omitempty"`
	Success         bool                 `protobuf:"varint,5,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason      `protobuf:"varint,6,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *UnipoolV3DataResponse) Reset()      { *m = UnipoolV3DataResponse{} }
func (*UnipoolV3DataResponse) ProtoMessage() {}
func (*UnipoolV3DataResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{14}
}
func (m *UnipoolV3DataResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UnipoolV3DataResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UnipoolV3DataResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UnipoolV3DataResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnipoolV3DataResponse.Merge(m, src)
}
func (m *UnipoolV3DataResponse) XXX_Size() int {
	return m.Size()
}
func (m *UnipoolV3DataResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UnipoolV3DataResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UnipoolV3DataResponse proto.InternalMessageInfo

func (m *UnipoolV3DataResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *UnipoolV3DataResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *UnipoolV3DataResponse) GetUpdate() []*models.UPV3Update {
	if m != nil {
		return m.Update
	}
	return nil
}

func (m *UnipoolV3DataResponse) GetSeqNum() uint64 {
	if m != nil {
		return m.SeqNum
	}
	return 0
}

func (m *UnipoolV3DataResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *UnipoolV3DataResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type ProtocolAssetTransferRequest struct {
	RequestID     uint64                `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe     bool                  `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
	Subscriber    *actor.PID            `protobuf:"bytes,3,opt,name=subscriber,proto3" json:"subscriber,omitempty"`
	ProtocolAsset *models.ProtocolAsset `protobuf:"bytes,4,opt,name=protocol_asset,json=protocolAsset,proto3" json:"protocol_asset,omitempty"`
}

func (m *ProtocolAssetTransferRequest) Reset()      { *m = ProtocolAssetTransferRequest{} }
func (*ProtocolAssetTransferRequest) ProtoMessage() {}
func (*ProtocolAssetTransferRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{15}
}
func (m *ProtocolAssetTransferRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProtocolAssetTransferRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProtocolAssetTransferRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProtocolAssetTransferRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProtocolAssetTransferRequest.Merge(m, src)
}
func (m *ProtocolAssetTransferRequest) XXX_Size() int {
	return m.Size()
}
func (m *ProtocolAssetTransferRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProtocolAssetTransferRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProtocolAssetTransferRequest proto.InternalMessageInfo

func (m *ProtocolAssetTransferRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *ProtocolAssetTransferRequest) GetSubscribe() bool {
	if m != nil {
		return m.Subscribe
	}
	return false
}

func (m *ProtocolAssetTransferRequest) GetSubscriber() *actor.PID {
	if m != nil {
		return m.Subscriber
	}
	return nil
}

func (m *ProtocolAssetTransferRequest) GetProtocolAsset() *models.ProtocolAsset {
	if m != nil {
		return m.ProtocolAsset
	}
	return nil
}

type ProtocolAssetTransferResponse struct {
	RequestID       uint64                        `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64                        `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Update          []*models.ProtocolAssetUpdate `protobuf:"bytes,3,rep,name=update,proto3" json:"update,omitempty"`
	SeqNum          uint64                        `protobuf:"varint,4,opt,name=seq_num,json=seqNum,proto3" json:"seq_num,omitempty"`
	Success         bool                          `protobuf:"varint,5,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason               `protobuf:"varint,6,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *ProtocolAssetTransferResponse) Reset()      { *m = ProtocolAssetTransferResponse{} }
func (*ProtocolAssetTransferResponse) ProtoMessage() {}
func (*ProtocolAssetTransferResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{16}
}
func (m *ProtocolAssetTransferResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProtocolAssetTransferResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProtocolAssetTransferResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProtocolAssetTransferResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProtocolAssetTransferResponse.Merge(m, src)
}
func (m *ProtocolAssetTransferResponse) XXX_Size() int {
	return m.Size()
}
func (m *ProtocolAssetTransferResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ProtocolAssetTransferResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ProtocolAssetTransferResponse proto.InternalMessageInfo

func (m *ProtocolAssetTransferResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *ProtocolAssetTransferResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *ProtocolAssetTransferResponse) GetUpdate() []*models.ProtocolAssetUpdate {
	if m != nil {
		return m.Update
	}
	return nil
}

func (m *ProtocolAssetTransferResponse) GetSeqNum() uint64 {
	if m != nil {
		return m.SeqNum
	}
	return 0
}

func (m *ProtocolAssetTransferResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *ProtocolAssetTransferResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type UnipoolV3DataIncrementalRefresh struct {
	RequestID  uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID uint64             `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	SeqNum     uint64             `protobuf:"varint,3,opt,name=seq_num,json=seqNum,proto3" json:"seq_num,omitempty"`
	Update     *models.UPV3Update `protobuf:"bytes,4,opt,name=update,proto3" json:"update,omitempty"`
}

func (m *UnipoolV3DataIncrementalRefresh) Reset()      { *m = UnipoolV3DataIncrementalRefresh{} }
func (*UnipoolV3DataIncrementalRefresh) ProtoMessage() {}
func (*UnipoolV3DataIncrementalRefresh) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{17}
}
func (m *UnipoolV3DataIncrementalRefresh) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UnipoolV3DataIncrementalRefresh) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UnipoolV3DataIncrementalRefresh.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UnipoolV3DataIncrementalRefresh) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnipoolV3DataIncrementalRefresh.Merge(m, src)
}
func (m *UnipoolV3DataIncrementalRefresh) XXX_Size() int {
	return m.Size()
}
func (m *UnipoolV3DataIncrementalRefresh) XXX_DiscardUnknown() {
	xxx_messageInfo_UnipoolV3DataIncrementalRefresh.DiscardUnknown(m)
}

var xxx_messageInfo_UnipoolV3DataIncrementalRefresh proto.InternalMessageInfo

func (m *UnipoolV3DataIncrementalRefresh) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *UnipoolV3DataIncrementalRefresh) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *UnipoolV3DataIncrementalRefresh) GetSeqNum() uint64 {
	if m != nil {
		return m.SeqNum
	}
	return 0
}

func (m *UnipoolV3DataIncrementalRefresh) GetUpdate() *models.UPV3Update {
	if m != nil {
		return m.Update
	}
	return nil
}

type ProtocolAssetDataIncrementalRefresh struct {
	RequestID  uint64                      `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID uint64                      `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	SeqNum     uint64                      `protobuf:"varint,3,opt,name=seq_num,json=seqNum,proto3" json:"seq_num,omitempty"`
	Update     *models.ProtocolAssetUpdate `protobuf:"bytes,4,opt,name=update,proto3" json:"update,omitempty"`
}

func (m *ProtocolAssetDataIncrementalRefresh) Reset()      { *m = ProtocolAssetDataIncrementalRefresh{} }
func (*ProtocolAssetDataIncrementalRefresh) ProtoMessage() {}
func (*ProtocolAssetDataIncrementalRefresh) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{18}
}
func (m *ProtocolAssetDataIncrementalRefresh) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProtocolAssetDataIncrementalRefresh) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProtocolAssetDataIncrementalRefresh.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProtocolAssetDataIncrementalRefresh) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProtocolAssetDataIncrementalRefresh.Merge(m, src)
}
func (m *ProtocolAssetDataIncrementalRefresh) XXX_Size() int {
	return m.Size()
}
func (m *ProtocolAssetDataIncrementalRefresh) XXX_DiscardUnknown() {
	xxx_messageInfo_ProtocolAssetDataIncrementalRefresh.DiscardUnknown(m)
}

var xxx_messageInfo_ProtocolAssetDataIncrementalRefresh proto.InternalMessageInfo

func (m *ProtocolAssetDataIncrementalRefresh) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *ProtocolAssetDataIncrementalRefresh) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *ProtocolAssetDataIncrementalRefresh) GetSeqNum() uint64 {
	if m != nil {
		return m.SeqNum
	}
	return 0
}

func (m *ProtocolAssetDataIncrementalRefresh) GetUpdate() *models.ProtocolAssetUpdate {
	if m != nil {
		return m.Update
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
	return fileDescriptor_350c53ba9303a7e6, []int{19}
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
	return fileDescriptor_350c53ba9303a7e6, []int{20}
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
	return fileDescriptor_350c53ba9303a7e6, []int{21}
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

type AccountMovement struct {
	Asset      *models1.Asset      `protobuf:"bytes,1,opt,name=asset,proto3" json:"asset,omitempty"`
	Change     float64             `protobuf:"fixed64,2,opt,name=change,proto3" json:"change,omitempty"`
	Type       AccountMovementType `protobuf:"varint,3,opt,name=type,proto3,enum=messages.AccountMovementType" json:"type,omitempty"`
	Subtype    string              `protobuf:"bytes,4,opt,name=subtype,proto3" json:"subtype,omitempty"`
	MovementID string              `protobuf:"bytes,5,opt,name=movementID,proto3" json:"movementID,omitempty"`
	Time       *types.Timestamp    `protobuf:"bytes,6,opt,name=time,proto3" json:"time,omitempty"`
}

func (m *AccountMovement) Reset()      { *m = AccountMovement{} }
func (*AccountMovement) ProtoMessage() {}
func (*AccountMovement) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{22}
}
func (m *AccountMovement) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountMovement) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccountMovement.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AccountMovement) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountMovement.Merge(m, src)
}
func (m *AccountMovement) XXX_Size() int {
	return m.Size()
}
func (m *AccountMovement) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountMovement.DiscardUnknown(m)
}

var xxx_messageInfo_AccountMovement proto.InternalMessageInfo

func (m *AccountMovement) GetAsset() *models1.Asset {
	if m != nil {
		return m.Asset
	}
	return nil
}

func (m *AccountMovement) GetChange() float64 {
	if m != nil {
		return m.Change
	}
	return 0
}

func (m *AccountMovement) GetType() AccountMovementType {
	if m != nil {
		return m.Type
	}
	return Unknown
}

func (m *AccountMovement) GetSubtype() string {
	if m != nil {
		return m.Subtype
	}
	return ""
}

func (m *AccountMovement) GetMovementID() string {
	if m != nil {
		return m.MovementID
	}
	return ""
}

func (m *AccountMovement) GetTime() *types.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

type AccountMovementRequest struct {
	RequestID uint64                 `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Account   *models.Account        `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	Type      AccountMovementType    `protobuf:"varint,3,opt,name=type,proto3,enum=messages.AccountMovementType" json:"type,omitempty"`
	Filter    *AccountMovementFilter `protobuf:"bytes,4,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (m *AccountMovementRequest) Reset()      { *m = AccountMovementRequest{} }
func (*AccountMovementRequest) ProtoMessage() {}
func (*AccountMovementRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{23}
}
func (m *AccountMovementRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountMovementRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccountMovementRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AccountMovementRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountMovementRequest.Merge(m, src)
}
func (m *AccountMovementRequest) XXX_Size() int {
	return m.Size()
}
func (m *AccountMovementRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountMovementRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AccountMovementRequest proto.InternalMessageInfo

func (m *AccountMovementRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *AccountMovementRequest) GetAccount() *models.Account {
	if m != nil {
		return m.Account
	}
	return nil
}

func (m *AccountMovementRequest) GetType() AccountMovementType {
	if m != nil {
		return m.Type
	}
	return Unknown
}

func (m *AccountMovementRequest) GetFilter() *AccountMovementFilter {
	if m != nil {
		return m.Filter
	}
	return nil
}

type AccountMovementFilter struct {
	Instrument *models.Instrument `protobuf:"bytes,1,opt,name=instrument,proto3" json:"instrument,omitempty"`
	From       *types.Timestamp   `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To         *types.Timestamp   `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
}

func (m *AccountMovementFilter) Reset()      { *m = AccountMovementFilter{} }
func (*AccountMovementFilter) ProtoMessage() {}
func (*AccountMovementFilter) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{24}
}
func (m *AccountMovementFilter) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountMovementFilter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccountMovementFilter.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AccountMovementFilter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountMovementFilter.Merge(m, src)
}
func (m *AccountMovementFilter) XXX_Size() int {
	return m.Size()
}
func (m *AccountMovementFilter) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountMovementFilter.DiscardUnknown(m)
}

var xxx_messageInfo_AccountMovementFilter proto.InternalMessageInfo

func (m *AccountMovementFilter) GetInstrument() *models.Instrument {
	if m != nil {
		return m.Instrument
	}
	return nil
}

func (m *AccountMovementFilter) GetFrom() *types.Timestamp {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *AccountMovementFilter) GetTo() *types.Timestamp {
	if m != nil {
		return m.To
	}
	return nil
}

type AccountMovementResponse struct {
	RequestID       uint64             `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64             `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Movements       []*AccountMovement `protobuf:"bytes,3,rep,name=movements,proto3" json:"movements,omitempty"`
	Success         bool               `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason    `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *AccountMovementResponse) Reset()      { *m = AccountMovementResponse{} }
func (*AccountMovementResponse) ProtoMessage() {}
func (*AccountMovementResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{25}
}
func (m *AccountMovementResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountMovementResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccountMovementResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AccountMovementResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountMovementResponse.Merge(m, src)
}
func (m *AccountMovementResponse) XXX_Size() int {
	return m.Size()
}
func (m *AccountMovementResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountMovementResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AccountMovementResponse proto.InternalMessageInfo

func (m *AccountMovementResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *AccountMovementResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *AccountMovementResponse) GetMovements() []*AccountMovement {
	if m != nil {
		return m.Movements
	}
	return nil
}

func (m *AccountMovementResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *AccountMovementResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type TradeCaptureReportRequest struct {
	RequestID uint64                    `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Account   *models.Account           `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	Filter    *TradeCaptureReportFilter `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (m *TradeCaptureReportRequest) Reset()      { *m = TradeCaptureReportRequest{} }
func (*TradeCaptureReportRequest) ProtoMessage() {}
func (*TradeCaptureReportRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{26}
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
	FromID        *types.StringValue `protobuf:"bytes,7,opt,name=fromID,proto3" json:"fromID,omitempty"`
}

func (m *TradeCaptureReportFilter) Reset()      { *m = TradeCaptureReportFilter{} }
func (*TradeCaptureReportFilter) ProtoMessage() {}
func (*TradeCaptureReportFilter) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{27}
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

func (m *TradeCaptureReportFilter) GetFromID() *types.StringValue {
	if m != nil {
		return m.FromID
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
	return fileDescriptor_350c53ba9303a7e6, []int{28}
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
	return fileDescriptor_350c53ba9303a7e6, []int{29}
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
	return fileDescriptor_350c53ba9303a7e6, []int{30}
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
	return fileDescriptor_350c53ba9303a7e6, []int{31}
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
	return fileDescriptor_350c53ba9303a7e6, []int{32}
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
	return fileDescriptor_350c53ba9303a7e6, []int{33}
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
	Type    AccountMovementType `protobuf:"varint,1,opt,name=type,proto3,enum=messages.AccountMovementType" json:"type,omitempty"`
	Asset   *models1.Asset      `protobuf:"bytes,2,opt,name=asset,proto3" json:"asset,omitempty"`
	Balance float64             `protobuf:"fixed64,3,opt,name=balance,proto3" json:"balance,omitempty"`
}

func (m *AccountUpdate) Reset()      { *m = AccountUpdate{} }
func (*AccountUpdate) ProtoMessage() {}
func (*AccountUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{34}
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

func (m *AccountUpdate) GetType() AccountMovementType {
	if m != nil {
		return m.Type
	}
	return Unknown
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
	return fileDescriptor_350c53ba9303a7e6, []int{35}
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
	return fileDescriptor_350c53ba9303a7e6, []int{36}
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
	return fileDescriptor_350c53ba9303a7e6, []int{37}
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
	return fileDescriptor_350c53ba9303a7e6, []int{38}
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
	ResponseID      uint64          `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	Success         bool            `protobuf:"varint,3,opt,name=success,proto3" json:"success,omitempty"`
	Orders          []*models.Order `protobuf:"bytes,4,rep,name=orders,proto3" json:"orders,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *OrderList) Reset()      { *m = OrderList{} }
func (*OrderList) ProtoMessage() {}
func (*OrderList) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{39}
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

func (m *OrderList) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
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
	return fileDescriptor_350c53ba9303a7e6, []int{40}
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
	return fileDescriptor_350c53ba9303a7e6, []int{41}
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
	return fileDescriptor_350c53ba9303a7e6, []int{42}
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
	return fileDescriptor_350c53ba9303a7e6, []int{43}
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
	ClientOrderID         string                        `protobuf:"bytes,1,opt,name=client_orderID,json=clientOrderID,proto3" json:"client_orderID,omitempty"`
	Instrument            *models.Instrument            `protobuf:"bytes,2,opt,name=instrument,proto3" json:"instrument,omitempty"`
	OrderType             models.OrderType              `protobuf:"varint,5,opt,name=order_type,json=orderType,proto3,enum=models.OrderType" json:"order_type,omitempty"`
	OrderSide             models.Side                   `protobuf:"varint,6,opt,name=order_side,json=orderSide,proto3,enum=models.Side" json:"order_side,omitempty"`
	TimeInForce           models.TimeInForce            `protobuf:"varint,7,opt,name=time_in_force,json=timeInForce,proto3,enum=models.TimeInForce" json:"time_in_force,omitempty"`
	Quantity              float64                       `protobuf:"fixed64,8,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Price                 *types.DoubleValue            `protobuf:"bytes,9,opt,name=price,proto3" json:"price,omitempty"`
	ExecutionInstructions []models.ExecutionInstruction `protobuf:"varint,10,rep,packed,name=execution_instructions,json=executionInstructions,proto3,enum=models.ExecutionInstruction" json:"execution_instructions,omitempty"`
	Tag                   string                        `protobuf:"bytes,11,opt,name=tag,proto3" json:"tag,omitempty"`
}

func (m *NewOrder) Reset()      { *m = NewOrder{} }
func (*NewOrder) ProtoMessage() {}
func (*NewOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{44}
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

func (m *NewOrder) GetExecutionInstructions() []models.ExecutionInstruction {
	if m != nil {
		return m.ExecutionInstructions
	}
	return nil
}

func (m *NewOrder) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

type NewOrderSingleRequest struct {
	RequestID uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Account   *models.Account `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	Order     *NewOrder       `protobuf:"bytes,3,opt,name=order,proto3" json:"order,omitempty"`
}

func (m *NewOrderSingleRequest) Reset()      { *m = NewOrderSingleRequest{} }
func (*NewOrderSingleRequest) ProtoMessage() {}
func (*NewOrderSingleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{45}
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
	return fileDescriptor_350c53ba9303a7e6, []int{46}
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
	return fileDescriptor_350c53ba9303a7e6, []int{47}
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
	return fileDescriptor_350c53ba9303a7e6, []int{48}
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
	return fileDescriptor_350c53ba9303a7e6, []int{49}
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
	return fileDescriptor_350c53ba9303a7e6, []int{50}
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
	return fileDescriptor_350c53ba9303a7e6, []int{51}
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
	return fileDescriptor_350c53ba9303a7e6, []int{52}
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
	return fileDescriptor_350c53ba9303a7e6, []int{53}
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
	return fileDescriptor_350c53ba9303a7e6, []int{54}
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
	return fileDescriptor_350c53ba9303a7e6, []int{55}
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
	return fileDescriptor_350c53ba9303a7e6, []int{56}
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
	return fileDescriptor_350c53ba9303a7e6, []int{57}
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

type TransferDataRequest struct {
	RequestID  uint64     `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe  bool       `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
	Subscriber *actor.PID `protobuf:"bytes,3,opt,name=subscriber,proto3" json:"subscriber,omitempty"`
}

func (m *TransferDataRequest) Reset()      { *m = TransferDataRequest{} }
func (*TransferDataRequest) ProtoMessage() {}
func (*TransferDataRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{58}
}
func (m *TransferDataRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TransferDataRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TransferDataRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TransferDataRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferDataRequest.Merge(m, src)
}
func (m *TransferDataRequest) XXX_Size() int {
	return m.Size()
}
func (m *TransferDataRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferDataRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TransferDataRequest proto.InternalMessageInfo

func (m *TransferDataRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *TransferDataRequest) GetSubscribe() bool {
	if m != nil {
		return m.Subscribe
	}
	return false
}

func (m *TransferDataRequest) GetSubscriber() *actor.PID {
	if m != nil {
		return m.Subscriber
	}
	return nil
}

type TransferDataResponse struct {
	RequestID       uint64          `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64          `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	SeqNum          uint64          `protobuf:"varint,3,opt,name=seq_num,json=seqNum,proto3" json:"seq_num,omitempty"`
	Success         bool            `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *TransferDataResponse) Reset()      { *m = TransferDataResponse{} }
func (*TransferDataResponse) ProtoMessage() {}
func (*TransferDataResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{59}
}
func (m *TransferDataResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TransferDataResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TransferDataResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TransferDataResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferDataResponse.Merge(m, src)
}
func (m *TransferDataResponse) XXX_Size() int {
	return m.Size()
}
func (m *TransferDataResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferDataResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TransferDataResponse proto.InternalMessageInfo

func (m *TransferDataResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *TransferDataResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *TransferDataResponse) GetSeqNum() uint64 {
	if m != nil {
		return m.SeqNum
	}
	return 0
}

func (m *TransferDataResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *TransferDataResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type ProtocolAssetDefinitionRequest struct {
	RequestID       uint64 `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ProtocolAssetID uint64 `protobuf:"varint,2,opt,name=protocol_asset_iD,json=protocolAssetID,proto3" json:"protocol_asset_iD,omitempty"`
}

func (m *ProtocolAssetDefinitionRequest) Reset()      { *m = ProtocolAssetDefinitionRequest{} }
func (*ProtocolAssetDefinitionRequest) ProtoMessage() {}
func (*ProtocolAssetDefinitionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{60}
}
func (m *ProtocolAssetDefinitionRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProtocolAssetDefinitionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProtocolAssetDefinitionRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProtocolAssetDefinitionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProtocolAssetDefinitionRequest.Merge(m, src)
}
func (m *ProtocolAssetDefinitionRequest) XXX_Size() int {
	return m.Size()
}
func (m *ProtocolAssetDefinitionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProtocolAssetDefinitionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProtocolAssetDefinitionRequest proto.InternalMessageInfo

func (m *ProtocolAssetDefinitionRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *ProtocolAssetDefinitionRequest) GetProtocolAssetID() uint64 {
	if m != nil {
		return m.ProtocolAssetID
	}
	return 0
}

type ProtocolAssetDefinitionResponse struct {
	RequestID       uint64                `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64                `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	ProtocolAsset   *models.ProtocolAsset `protobuf:"bytes,3,opt,name=protocol_asset,json=protocolAsset,proto3" json:"protocol_asset,omitempty"`
	Success         bool                  `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason       `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *ProtocolAssetDefinitionResponse) Reset()      { *m = ProtocolAssetDefinitionResponse{} }
func (*ProtocolAssetDefinitionResponse) ProtoMessage() {}
func (*ProtocolAssetDefinitionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{61}
}
func (m *ProtocolAssetDefinitionResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProtocolAssetDefinitionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProtocolAssetDefinitionResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProtocolAssetDefinitionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProtocolAssetDefinitionResponse.Merge(m, src)
}
func (m *ProtocolAssetDefinitionResponse) XXX_Size() int {
	return m.Size()
}
func (m *ProtocolAssetDefinitionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ProtocolAssetDefinitionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ProtocolAssetDefinitionResponse proto.InternalMessageInfo

func (m *ProtocolAssetDefinitionResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *ProtocolAssetDefinitionResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *ProtocolAssetDefinitionResponse) GetProtocolAsset() *models.ProtocolAsset {
	if m != nil {
		return m.ProtocolAsset
	}
	return nil
}

func (m *ProtocolAssetDefinitionResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *ProtocolAssetDefinitionResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type ProtocolAssetListRequest struct {
	RequestID  uint64     `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Subscribe  bool       `protobuf:"varint,2,opt,name=subscribe,proto3" json:"subscribe,omitempty"`
	Subscriber *actor.PID `protobuf:"bytes,3,opt,name=subscriber,proto3" json:"subscriber,omitempty"`
}

func (m *ProtocolAssetListRequest) Reset()      { *m = ProtocolAssetListRequest{} }
func (*ProtocolAssetListRequest) ProtoMessage() {}
func (*ProtocolAssetListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{62}
}
func (m *ProtocolAssetListRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProtocolAssetListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProtocolAssetListRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProtocolAssetListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProtocolAssetListRequest.Merge(m, src)
}
func (m *ProtocolAssetListRequest) XXX_Size() int {
	return m.Size()
}
func (m *ProtocolAssetListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProtocolAssetListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProtocolAssetListRequest proto.InternalMessageInfo

func (m *ProtocolAssetListRequest) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *ProtocolAssetListRequest) GetSubscribe() bool {
	if m != nil {
		return m.Subscribe
	}
	return false
}

func (m *ProtocolAssetListRequest) GetSubscriber() *actor.PID {
	if m != nil {
		return m.Subscriber
	}
	return nil
}

type ProtocolAssetListResponse struct {
	RequestID       uint64                  `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64                  `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	ProtocolAssets  []*models.ProtocolAsset `protobuf:"bytes,3,rep,name=protocol_assets,json=protocolAssets,proto3" json:"protocol_assets,omitempty"`
	Success         bool                    `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason         `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *ProtocolAssetListResponse) Reset()      { *m = ProtocolAssetListResponse{} }
func (*ProtocolAssetListResponse) ProtoMessage() {}
func (*ProtocolAssetListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{63}
}
func (m *ProtocolAssetListResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProtocolAssetListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProtocolAssetListResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProtocolAssetListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProtocolAssetListResponse.Merge(m, src)
}
func (m *ProtocolAssetListResponse) XXX_Size() int {
	return m.Size()
}
func (m *ProtocolAssetListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ProtocolAssetListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ProtocolAssetListResponse proto.InternalMessageInfo

func (m *ProtocolAssetListResponse) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *ProtocolAssetListResponse) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *ProtocolAssetListResponse) GetProtocolAssets() []*models.ProtocolAsset {
	if m != nil {
		return m.ProtocolAssets
	}
	return nil
}

func (m *ProtocolAssetListResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *ProtocolAssetListResponse) GetRejectionReason() RejectionReason {
	if m != nil {
		return m.RejectionReason
	}
	return Other
}

type ProtocolAssetList struct {
	RequestID       uint64                  `protobuf:"varint,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ResponseID      uint64                  `protobuf:"varint,2,opt,name=responseID,proto3" json:"responseID,omitempty"`
	ProtocolAssets  []*models.ProtocolAsset `protobuf:"bytes,3,rep,name=protocol_assets,json=protocolAssets,proto3" json:"protocol_assets,omitempty"`
	Success         bool                    `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	RejectionReason RejectionReason         `protobuf:"varint,5,opt,name=rejection_reason,json=rejectionReason,proto3,enum=messages.RejectionReason" json:"rejection_reason,omitempty"`
}

func (m *ProtocolAssetList) Reset()      { *m = ProtocolAssetList{} }
func (*ProtocolAssetList) ProtoMessage() {}
func (*ProtocolAssetList) Descriptor() ([]byte, []int) {
	return fileDescriptor_350c53ba9303a7e6, []int{64}
}
func (m *ProtocolAssetList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProtocolAssetList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProtocolAssetList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProtocolAssetList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProtocolAssetList.Merge(m, src)
}
func (m *ProtocolAssetList) XXX_Size() int {
	return m.Size()
}
func (m *ProtocolAssetList) XXX_DiscardUnknown() {
	xxx_messageInfo_ProtocolAssetList.DiscardUnknown(m)
}

var xxx_messageInfo_ProtocolAssetList proto.InternalMessageInfo

func (m *ProtocolAssetList) GetRequestID() uint64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *ProtocolAssetList) GetResponseID() uint64 {
	if m != nil {
		return m.ResponseID
	}
	return 0
}

func (m *ProtocolAssetList) GetProtocolAssets() []*models.ProtocolAsset {
	if m != nil {
		return m.ProtocolAssets
	}
	return nil
}

func (m *ProtocolAssetList) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *ProtocolAssetList) GetRejectionReason() RejectionReason {
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
	proto.RegisterEnum("messages.AccountMovementType", AccountMovementType_name, AccountMovementType_value)
	proto.RegisterType((*HistoricalOpenInterestsRequest)(nil), "messages.HistoricalOpenInterestsRequest")
	proto.RegisterType((*HistoricalOpenInterestsResponse)(nil), "messages.HistoricalOpenInterestsResponse")
	proto.RegisterType((*HistoricalLiquidationsRequest)(nil), "messages.HistoricalLiquidationsRequest")
	proto.RegisterType((*HistoricalLiquidationsResponse)(nil), "messages.HistoricalLiquidationsResponse")
	proto.RegisterType((*HistoricalUnipoolV3DataRequest)(nil), "messages.HistoricalUnipoolV3DataRequest")
	proto.RegisterType((*HistoricalUnipoolV3DataResponse)(nil), "messages.HistoricalUnipoolV3DataResponse")
	proto.RegisterType((*HistoricalProtocolAssetTransferRequest)(nil), "messages.HistoricalProtocolAssetTransferRequest")
	proto.RegisterType((*HistoricalProtocolAssetTransferResponse)(nil), "messages.HistoricalProtocolAssetTransferResponse")
	proto.RegisterType((*MarketStatisticsRequest)(nil), "messages.MarketStatisticsRequest")
	proto.RegisterType((*MarketStatisticsResponse)(nil), "messages.MarketStatisticsResponse")
	proto.RegisterType((*MarketDataRequest)(nil), "messages.MarketDataRequest")
	proto.RegisterType((*MarketDataResponse)(nil), "messages.MarketDataResponse")
	proto.RegisterType((*MarketDataIncrementalRefresh)(nil), "messages.MarketDataIncrementalRefresh")
	proto.RegisterType((*UnipoolV3DataRequest)(nil), "messages.UnipoolV3DataRequest")
	proto.RegisterType((*UnipoolV3DataResponse)(nil), "messages.UnipoolV3DataResponse")
	proto.RegisterType((*ProtocolAssetTransferRequest)(nil), "messages.ProtocolAssetTransferRequest")
	proto.RegisterType((*ProtocolAssetTransferResponse)(nil), "messages.ProtocolAssetTransferResponse")
	proto.RegisterType((*UnipoolV3DataIncrementalRefresh)(nil), "messages.UnipoolV3DataIncrementalRefresh")
	proto.RegisterType((*ProtocolAssetDataIncrementalRefresh)(nil), "messages.ProtocolAssetDataIncrementalRefresh")
	proto.RegisterType((*AccountDataRequest)(nil), "messages.AccountDataRequest")
	proto.RegisterType((*AccountDataResponse)(nil), "messages.AccountDataResponse")
	proto.RegisterType((*AccountDataIncrementalRefresh)(nil), "messages.AccountDataIncrementalRefresh")
	proto.RegisterType((*AccountMovement)(nil), "messages.AccountMovement")
	proto.RegisterType((*AccountMovementRequest)(nil), "messages.AccountMovementRequest")
	proto.RegisterType((*AccountMovementFilter)(nil), "messages.AccountMovementFilter")
	proto.RegisterType((*AccountMovementResponse)(nil), "messages.AccountMovementResponse")
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
	proto.RegisterType((*TransferDataRequest)(nil), "messages.TransferDataRequest")
	proto.RegisterType((*TransferDataResponse)(nil), "messages.TransferDataResponse")
	proto.RegisterType((*ProtocolAssetDefinitionRequest)(nil), "messages.ProtocolAssetDefinitionRequest")
	proto.RegisterType((*ProtocolAssetDefinitionResponse)(nil), "messages.ProtocolAssetDefinitionResponse")
	proto.RegisterType((*ProtocolAssetListRequest)(nil), "messages.ProtocolAssetListRequest")
	proto.RegisterType((*ProtocolAssetListResponse)(nil), "messages.ProtocolAssetListResponse")
	proto.RegisterType((*ProtocolAssetList)(nil), "messages.ProtocolAssetList")
}

func init() { proto.RegisterFile("executor_messages.proto", fileDescriptor_350c53ba9303a7e6) }

var fileDescriptor_350c53ba9303a7e6 = []byte{
	// 3693 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x5c, 0x4d, 0x6c, 0x1b, 0x57,
	0x7e, 0xe7, 0x90, 0x14, 0x3f, 0xfe, 0x94, 0xc4, 0xd1, 0x93, 0x6c, 0xd1, 0xb4, 0x4c, 0x3b, 0x93,
	0x6e, 0xe2, 0x68, 0x1b, 0x29, 0xa6, 0xb2, 0x49, 0xdb, 0xdd, 0x18, 0x90, 0x45, 0xb9, 0x4b, 0x40,
	0xb1, 0xb5, 0x23, 0x29, 0xbb, 0x40, 0x0f, 0xc4, 0x68, 0xf8, 0x44, 0xcd, 0x7a, 0x38, 0x43, 0xcf,
	0x3c, 0xda, 0x56, 0x81, 0x05, 0xba, 0x45, 0xd1, 0xc3, 0x16, 0x28, 0x72, 0x29, 0x0a, 0x14, 0x68,
	0xd1, 0x43, 0xd1, 0xf6, 0xd2, 0x02, 0xfd, 0x40, 0xdb, 0x43, 0x2f, 0x05, 0x0a, 0xb4, 0x3d, 0x64,
	0x91, 0x53, 0x91, 0x02, 0x45, 0xb1, 0x71, 0x80, 0x36, 0xc0, 0x5e, 0xd2, 0x53, 0x3f, 0x80, 0x02,
	0xc5, 0xfb, 0x1a, 0xbe, 0xe1, 0x87, 0x44, 0xc6, 0xa4, 0xbd, 0x09, 0xf6, 0xc6, 0x99, 0xf7, 0xff,
	0xbf, 0xf7, 0xff, 0x7e, 0xbf, 0xf7, 0x7f, 0x23, 0xc1, 0x2a, 0x7e, 0x82, 0xed, 0x2e, 0xf1, 0x83,
	0x46, 0x1b, 0x87, 0xa1, 0xd5, 0xc2, 0xe1, 0x46, 0x27, 0xf0, 0x89, 0x8f, 0x72, 0xf2, 0xb9, 0x7c,
	0xbd, 0xe5, 0xfb, 0x2d, 0x17, 0x6f, 0xb2, 0xf7, 0xc7, 0xdd, 0x93, 0x4d, 0xe2, 0xb4, 0x71, 0x48,
	0xac, 0x76, 0x87, 0x93, 0x96, 0x2b, 0xfd, 0x04, 0x8f, 0x03, 0xab, 0xd3, 0xc1, 0x81, 0x98, 0xaa,
	0x5c, 0x6d, 0x39, 0xc4, 0xb5, 0x8e, 0x37, 0x6c, 0xbf, 0xbd, 0x69, 0xb9, 0x9d, 0x53, 0x8b, 0x38,
	0xf6, 0x83, 0x70, 0xf3, 0x89, 0x7d, 0x6a, 0x79, 0x2d, 0x1c, 0x6c, 0xb6, 0xfd, 0x26, 0x76, 0x43,
	0xce, 0x2e, 0x79, 0xbe, 0x31, 0x9c, 0x87, 0xfd, 0x7c, 0xdd, 0xf6, 0x3d, 0x0f, 0xdb, 0x44, 0x32,
	0xb6, 0xad, 0xe0, 0x01, 0x26, 0x8d, 0xa6, 0x45, 0x2c, 0xc1, 0x7d, 0x7b, 0x02, 0xee, 0x10, 0xdb,
	0xdd, 0xc0, 0x21, 0x67, 0x2a, 0xff, 0x3b, 0x13, 0xf0, 0x5b, 0xb6, 0xed, 0x77, 0xbd, 0xd8, 0xf2,
	0x5f, 0x9f, 0x84, 0x3d, 0x0c, 0xe3, 0xb2, 0xbf, 0xd5, 0x72, 0xc8, 0x69, 0x97, 0x33, 0x6f, 0x87,
	0x67, 0xde, 0x83, 0xc0, 0xf7, 0xea, 0x87, 0xdc, 0x3a, 0x96, 0x4d, 0xfc, 0xe0, 0xf5, 0x96, 0xbf,
	0xc9, 0x7e, 0xc4, 0x2c, 0x66, 0xfc, 0x50, 0x83, 0xca, 0x37, 0x9d, 0x90, 0xf8, 0x81, 0x63, 0x5b,
	0xee, 0xfd, 0x0e, 0xf6, 0xea, 0x1e, 0xc1, 0x01, 0x0e, 0x49, 0x68, 0xe2, 0x87, 0x5d, 0x1c, 0x12,
	0xb4, 0x06, 0xf9, 0x80, 0xff, 0xac, 0xd7, 0x4a, 0xda, 0x0d, 0xed, 0x66, 0xda, 0xec, 0xbd, 0x40,
	0x55, 0x00, 0xc7, 0x0b, 0x49, 0xd0, 0x6d, 0x63, 0x8f, 0x94, 0x92, 0x37, 0xb4, 0x9b, 0x85, 0x2a,
	0xda, 0xe0, 0x62, 0x6e, 0xd4, 0xa3, 0x11, 0x53, 0xa1, 0x42, 0x1b, 0x90, 0x3e, 0x09, 0xfc, 0x76,
	0x29, 0xc5, 0xa8, 0xcb, 0x1b, 0x3c, 0x12, 0x36, 0x64, 0x24, 0x6c, 0x1c, 0xca, 0x50, 0x31, 0x19,
	0x1d, 0x5a, 0x87, 0x24, 0xf1, 0x4b, 0xe9, 0x0b, 0xa9, 0x93, 0xc4, 0x37, 0x7e, 0xac, 0xc1, 0xf5,
	0x91, 0x0a, 0x85, 0x1d, 0xdf, 0x0b, 0xf1, 0x05, 0x1a, 0x55, 0x00, 0x02, 0x41, 0x59, 0xaf, 0x31,
	0x8d, 0xd2, 0xa6, 0xf2, 0x06, 0xad, 0x43, 0xde, 0x91, 0x53, 0x96, 0x52, 0x37, 0x52, 0x37, 0x0b,
	0xd5, 0x79, 0xa9, 0xf0, 0x01, 0xb1, 0x88, 0xd9, 0x1b, 0x46, 0x25, 0xc8, 0x86, 0x5d, 0xdb, 0xc6,
	0x61, 0xc8, 0xc4, 0xcf, 0x99, 0xf2, 0x11, 0xd5, 0x40, 0x0f, 0xf0, 0x77, 0xb1, 0x4d, 0x1c, 0xdf,
	0x6b, 0x04, 0xd8, 0x0a, 0x7d, 0xaf, 0x34, 0x77, 0x43, 0xbb, 0xb9, 0x58, 0xbd, 0xb2, 0x11, 0x25,
	0x95, 0x29, 0x29, 0x4c, 0x46, 0x60, 0x16, 0x83, 0xf8, 0x0b, 0xe3, 0x03, 0x0d, 0xae, 0xf5, 0xb4,
	0xdd, 0x73, 0x1e, 0x76, 0x9d, 0xa6, 0x45, 0xc7, 0xbf, 0xa0, 0xde, 0xfb, 0xaf, 0x58, 0x38, 0xc6,
	0xf5, 0x99, 0x8a, 0xf3, 0xde, 0x86, 0x79, 0x57, 0x99, 0x55, 0xf8, 0x6f, 0x59, 0xaa, 0xac, 0xac,
	0x68, 0xc6, 0x08, 0x67, 0xee, 0xc9, 0xdf, 0x89, 0x69, 0x7e, 0xe4, 0x39, 0x1d, 0xdf, 0x77, 0xdf,
	0xdb, 0xaa, 0x59, 0xc4, 0x9a, 0x9d, 0x2b, 0x57, 0x60, 0x2e, 0x24, 0x56, 0x40, 0x98, 0x2f, 0xd3,
	0x26, 0x7f, 0x40, 0x3a, 0xa4, 0xb0, 0xd7, 0x64, 0x6a, 0xa6, 0x4d, 0xfa, 0xd3, 0xf8, 0xb5, 0xa4,
	0x9a, 0x54, 0x7d, 0xc2, 0x4d, 0x29, 0xa9, 0x32, 0xf8, 0x11, 0xf6, 0xa2, 0x8c, 0x8a, 0x24, 0x3f,
	0xda, 0x7f, 0x6f, 0xeb, 0xa8, 0xd3, 0xb4, 0x08, 0x36, 0x05, 0x05, 0x5a, 0x85, 0x6c, 0x88, 0x1f,
	0x36, 0xbc, 0x6e, 0x5b, 0xc8, 0x98, 0x09, 0xf1, 0xc3, 0x7b, 0xdd, 0xb6, 0xea, 0xa3, 0xb9, 0x8b,
	0x7d, 0x94, 0x99, 0xd8, 0x47, 0x7f, 0xaa, 0xc1, 0x2b, 0x3d, 0x33, 0xec, 0xd3, 0x10, 0xb6, 0x7d,
	0x77, 0x9b, 0xd6, 0xe2, 0xc3, 0xc0, 0xf2, 0xc2, 0x13, 0x1c, 0x8c, 0xe7, 0xab, 0x6f, 0xc0, 0x62,
	0x47, 0x70, 0x37, 0x58, 0x29, 0x17, 0xfe, 0xba, 0x24, 0xb5, 0x8e, 0xcd, 0x6d, 0x2e, 0x74, 0xd4,
	0xc7, 0x11, 0x5e, 0x43, 0x90, 0x0e, 0x89, 0xdf, 0x11, 0x26, 0x61, 0xbf, 0x8d, 0xdf, 0x4a, 0xc2,
	0xab, 0x17, 0x0a, 0x3c, 0x15, 0xff, 0x6d, 0x41, 0xa6, 0xcb, 0xbc, 0x24, 0xfc, 0x77, 0x75, 0xa8,
	0x26, 0xd2, 0x91, 0x9c, 0xf4, 0xc5, 0x39, 0xf2, 0x77, 0x35, 0x58, 0x7d, 0x97, 0xed, 0xff, 0xb4,
	0x60, 0x3b, 0x21, 0x71, 0xec, 0x19, 0x16, 0xcc, 0x37, 0x00, 0xc2, 0x68, 0x19, 0x66, 0x9f, 0xc5,
	0xaa, 0xae, 0xee, 0x18, 0x87, 0x67, 0x1d, 0x6c, 0x2a, 0x34, 0xc6, 0xbf, 0x6b, 0x50, 0x1a, 0x94,
	0x6f, 0x2a, 0x8e, 0xfa, 0xd9, 0x01, 0x61, 0xfa, 0xb7, 0x2f, 0x65, 0x7c, 0xe6, 0x55, 0xef, 0xb7,
	0x93, 0xb0, 0xc4, 0x15, 0x1d, 0xbf, 0xd0, 0xad, 0x41, 0x3e, 0xec, 0x1e, 0x87, 0x76, 0xe0, 0x1c,
	0x63, 0xa6, 0x60, 0xce, 0xec, 0xbd, 0x40, 0xeb, 0x00, 0xd1, 0x43, 0x20, 0xf6, 0x28, 0xd8, 0x60,
	0xc8, 0x67, 0x63, 0xbf, 0x5e, 0x33, 0x95, 0xd1, 0x3e, 0x67, 0xa6, 0xc7, 0x72, 0xe6, 0x6d, 0x28,
	0x58, 0xad, 0x56, 0x80, 0x5b, 0x6c, 0x5f, 0x10, 0x2a, 0xaf, 0x49, 0xa6, 0xfb, 0x41, 0x13, 0x07,
	0x77, 0x7c, 0xff, 0xc1, 0x76, 0x8f, 0xc6, 0x54, 0x19, 0xd0, 0x2b, 0x2c, 0x79, 0x49, 0x58, 0xca,
	0x8c, 0x88, 0x03, 0x3e, 0x6c, 0xfc, 0x5e, 0x0a, 0x90, 0x6a, 0x99, 0xa9, 0x38, 0xff, 0x4d, 0x80,
	0xd0, 0xb3, 0x3a, 0xe1, 0xa9, 0x4f, 0xf6, 0x6e, 0x09, 0xe3, 0xac, 0x44, 0xb2, 0xdf, 0xd9, 0xbb,
	0x75, 0x20, 0x46, 0x4d, 0x85, 0x2e, 0xc6, 0x55, 0x15, 0x66, 0x52, 0xb9, 0xaa, 0x43, 0xb8, 0xaa,
	0x31, 0xae, 0x2d, 0x66, 0xa7, 0x38, 0xd7, 0xd6, 0x10, 0xae, 0x2d, 0xb4, 0x09, 0x19, 0x12, 0x58,
	0x4d, 0xcc, 0xed, 0x53, 0xa8, 0xae, 0x4a, 0x0e, 0x69, 0x50, 0xdc, 0x3c, 0xa4, 0xe3, 0xa6, 0x20,
	0x53, 0x6b, 0x48, 0x76, 0x54, 0x0d, 0xc9, 0x5d, 0x1c, 0xba, 0xf9, 0x89, 0x43, 0xf7, 0x6f, 0x53,
	0xb0, 0xd6, 0x73, 0x50, 0xdd, 0xb3, 0x03, 0x4c, 0x03, 0xc4, 0x72, 0x4d, 0x7c, 0x12, 0xe0, 0xf0,
	0xf4, 0x19, 0x5d, 0xa5, 0xe8, 0x95, 0x8a, 0xe9, 0xb5, 0x01, 0x39, 0x5e, 0x3e, 0xf7, 0x6e, 0xf5,
	0x87, 0x2c, 0xf5, 0xa0, 0x28, 0xb1, 0x11, 0x8d, 0x42, 0x5f, 0x15, 0x5e, 0x50, 0xe9, 0xab, 0x7d,
	0xf4, 0x55, 0x85, 0x7e, 0x8b, 0x55, 0xd6, 0x38, 0xfd, 0x56, 0x1f, 0xbd, 0xea, 0xb1, 0xec, 0x78,
	0x1e, 0xfb, 0x1a, 0x14, 0x14, 0x64, 0xc5, 0x9c, 0x33, 0x02, 0x81, 0xa9, 0x74, 0xe8, 0x35, 0xc8,
	0x9e, 0x74, 0xbd, 0xa6, 0xe3, 0xb5, 0x98, 0xb3, 0x0a, 0xd5, 0xa2, 0x64, 0xb9, 0xcb, 0x5f, 0x9b,
	0x72, 0x1c, 0x19, 0x32, 0xc7, 0x60, 0x48, 0x79, 0x13, 0xf9, 0xf5, 0xe7, 0x1a, 0xac, 0x7c, 0x0e,
	0x94, 0xf5, 0x42, 0x8b, 0x8f, 0xf1, 0xbf, 0x1a, 0x5c, 0x9a, 0x11, 0xfa, 0x8a, 0xed, 0xde, 0x43,
	0xd1, 0xd7, 0x8b, 0xde, 0xb4, 0xff, 0x5e, 0x83, 0xb5, 0x67, 0xc0, 0x5c, 0xd3, 0xf3, 0xdc, 0x20,
	0x7a, 0x4b, 0x8f, 0x8f, 0xde, 0x8c, 0xdf, 0x48, 0xc2, 0xb5, 0x9f, 0x22, 0x31, 0xe9, 0xd4, 0x3f,
	0xd0, 0xe0, 0x7a, 0x2c, 0xa2, 0x9f, 0x5f, 0x21, 0xed, 0x05, 0x7d, 0x5f, 0xf2, 0x0d, 0x06, 0xbd,
	0xf1, 0x67, 0x1a, 0xbc, 0x1c, 0xb3, 0xdf, 0xf3, 0x15, 0x75, 0xab, 0x4f, 0xd4, 0x71, 0x7c, 0x6a,
	0xfc, 0x91, 0x06, 0x68, 0x9b, 0xb7, 0x99, 0x5e, 0x44, 0x7d, 0x7b, 0x0d, 0xb2, 0xa2, 0xc9, 0x25,
	0x84, 0x8e, 0xea, 0xb5, 0x10, 0xca, 0x94, 0xe3, 0xc6, 0x7f, 0xa7, 0x60, 0x39, 0x26, 0xe9, 0x54,
	0x12, 0x81, 0xc2, 0x6e, 0xde, 0xa5, 0x73, 0xb0, 0x44, 0xba, 0x3d, 0xb8, 0x25, 0xfa, 0x77, 0xa6,
	0x42, 0x83, 0xbe, 0x02, 0x19, 0x9f, 0x02, 0x38, 0x0a, 0x76, 0x29, 0xf5, 0x42, 0x0c, 0xd6, 0x99,
	0x62, 0x10, 0x6d, 0x40, 0xbe, 0xe3, 0x87, 0x0e, 0x6f, 0x20, 0xcc, 0xc5, 0xe7, 0xdd, 0x17, 0x03,
	0x66, 0x8f, 0x04, 0x7d, 0x15, 0x72, 0xc7, 0x96, 0x6b, 0x79, 0x76, 0x84, 0x6a, 0x22, 0x53, 0xdc,
	0xe1, 0xef, 0xcd, 0x88, 0x00, 0xfd, 0x3c, 0xe4, 0xdb, 0xd6, 0x03, 0x1c, 0x34, 0x4e, 0x30, 0x66,
	0x88, 0xa6, 0x50, 0x5d, 0x1b, 0x68, 0x9a, 0xd4, 0xfc, 0xee, 0xb1, 0x8b, 0xdf, 0xb3, 0xdc, 0x2e,
	0x36, 0x73, 0x8c, 0xfc, 0x2e, 0xc6, 0x94, 0x95, 0x44, 0xac, 0xb9, 0x71, 0x58, 0x89, 0x64, 0x55,
	0xd2, 0x3c, 0x7f, 0x71, 0x9a, 0xc3, 0xa4, 0x69, 0xae, 0x46, 0x76, 0x41, 0x8d, 0x6c, 0xe3, 0x7d,
	0x0d, 0xae, 0x29, 0xae, 0x9f, 0x7a, 0x4a, 0xdd, 0x82, 0x4c, 0x80, 0x3b, 0xbe, 0x38, 0x2c, 0x17,
	0x54, 0xa1, 0x77, 0x59, 0x0f, 0x9b, 0xc9, 0x48, 0x09, 0x4c, 0x41, 0x68, 0xfc, 0x58, 0x83, 0xa2,
	0x10, 0xe9, 0x5d, 0xff, 0x11, 0x93, 0x06, 0xbd, 0x0c, 0x73, 0xbc, 0xd2, 0x6b, 0x6c, 0x96, 0x28,
	0x30, 0x78, 0x85, 0xe7, 0x63, 0xe8, 0x32, 0x64, 0x78, 0x73, 0x9a, 0xc9, 0xa1, 0x99, 0xe2, 0x09,
	0xdd, 0x82, 0x34, 0x39, 0xeb, 0x60, 0x26, 0xc1, 0x62, 0xf5, 0x5a, 0x4f, 0x82, 0xbe, 0x55, 0x18,
	0xfc, 0x67, 0xa4, 0xdc, 0x1f, 0xc7, 0x8c, 0x8b, 0x26, 0x4f, 0xde, 0x94, 0x8f, 0x54, 0xe1, 0xb6,
	0xa0, 0xaf, 0xd7, 0x58, 0x4d, 0xce, 0x9b, 0xca, 0x1b, 0xb4, 0x01, 0x69, 0xe2, 0xb4, 0xb1, 0x80,
	0x6e, 0xe7, 0x76, 0xe7, 0x28, 0x9d, 0xf1, 0x81, 0x06, 0x97, 0xfb, 0xe4, 0x18, 0xaf, 0x52, 0x28,
	0xf9, 0x9d, 0x3c, 0x3f, 0xbf, 0x3f, 0x8f, 0x01, 0xde, 0x86, 0xcc, 0x89, 0xe3, 0x12, 0x1c, 0x88,
	0xe2, 0x71, 0x7d, 0x24, 0xd3, 0x5d, 0x46, 0x66, 0x0a, 0x72, 0x5a, 0xf5, 0x2e, 0x0d, 0xa5, 0xe8,
	0x03, 0x5c, 0xda, 0x44, 0xbd, 0xce, 0xe4, 0x44, 0xbd, 0xce, 0xd4, 0x58, 0xbd, 0xce, 0xcf, 0x34,
	0x58, 0x1d, 0xb0, 0xfc, 0x94, 0x9a, 0x9c, 0x79, 0x19, 0x11, 0xb2, 0xf0, 0x5d, 0x19, 0x69, 0x3f,
	0xb3, 0x47, 0x3b, 0xf3, 0xe3, 0xfe, 0xef, 0x6b, 0x70, 0x85, 0x1d, 0x06, 0x76, 0xac, 0x0e, 0xe9,
	0x06, 0x58, 0x64, 0xde, 0xb4, 0xe3, 0xed, 0x17, 0xa2, 0xe0, 0xe1, 0x9e, 0x30, 0x7a, 0x22, 0x0e,
	0xae, 0xde, 0x17, 0x3f, 0xbf, 0x9e, 0x82, 0xd2, 0x28, 0x22, 0xf4, 0x16, 0x64, 0xd9, 0x1e, 0x20,
	0xe4, 0x1b, 0x56, 0x5f, 0x0f, 0x48, 0xe0, 0x78, 0x2d, 0x5e, 0x5f, 0x25, 0x31, 0xda, 0x81, 0x45,
	0xdb, 0x75, 0xb0, 0x47, 0x1a, 0x92, 0x3d, 0x39, 0x06, 0xfb, 0x02, 0xe7, 0xb9, 0x2f, 0x26, 0x89,
	0xc7, 0x6f, 0x6a, 0xac, 0xf8, 0x7d, 0x15, 0xd2, 0xa1, 0xd3, 0x94, 0xb0, 0x61, 0xb9, 0x67, 0x87,
	0x03, 0xa7, 0x29, 0x36, 0x01, 0x46, 0x10, 0x05, 0xfa, 0xdc, 0x44, 0x81, 0x9e, 0x19, 0x27, 0xd0,
	0xd1, 0x9b, 0x90, 0xa1, 0x3c, 0xf5, 0xda, 0xc8, 0xfd, 0x4c, 0xd5, 0x5a, 0xd0, 0x1a, 0x3f, 0xd2,
	0x00, 0x0d, 0x3a, 0xe2, 0x99, 0xbb, 0x5f, 0xf2, 0xb0, 0xca, 0xd3, 0x22, 0x6a, 0x48, 0xc4, 0x56,
	0x92, 0x27, 0xd5, 0x59, 0xa7, 0x43, 0x1b, 0xae, 0x48, 0x1c, 0x52, 0xc3, 0x27, 0x8e, 0xc7, 0x91,
	0xc3, 0xac, 0xfa, 0x90, 0xc6, 0xa7, 0x1a, 0x94, 0x87, 0xad, 0x37, 0xa5, 0xbe, 0x62, 0x4e, 0xde,
	0x89, 0x8a, 0xd8, 0x1c, 0xc4, 0x5a, 0x11, 0xc5, 0xcc, 0x2d, 0xfb, 0x3d, 0x58, 0x96, 0xab, 0xee,
	0x39, 0x21, 0x79, 0xce, 0xd8, 0xd7, 0xf8, 0x37, 0x0d, 0xe6, 0xd5, 0xf5, 0x9f, 0x3b, 0x92, 0x9d,
	0xb5, 0x7d, 0xff, 0x2a, 0x0b, 0xc5, 0x3e, 0xfc, 0xa4, 0x62, 0x3c, 0xad, 0xff, 0x0c, 0xa9, 0x96,
	0xbd, 0xfc, 0x79, 0x75, 0x31, 0x35, 0x79, 0x5d, 0xbc, 0x01, 0x05, 0x2c, 0x45, 0xa9, 0xd7, 0x04,
	0x5e, 0x52, 0x5f, 0xa1, 0xdb, 0xb0, 0x18, 0x3d, 0x36, 0x18, 0x12, 0xe1, 0x1a, 0xaf, 0x0e, 0x01,
	0x83, 0x0c, 0x83, 0x2c, 0x60, 0xf5, 0x11, 0xbd, 0x05, 0xf3, 0x4c, 0xbe, 0x46, 0x48, 0x2c, 0xd2,
	0x0d, 0xc5, 0x31, 0x77, 0x39, 0x76, 0x3a, 0x38, 0x60, 0x43, 0x66, 0xc1, 0xef, 0x3d, 0xf4, 0x25,
	0x69, 0x76, 0xcc, 0x8a, 0x5d, 0x74, 0xb1, 0xf5, 0x08, 0x87, 0x8d, 0x87, 0x5d, 0xcb, 0x23, 0x34,
	0x9d, 0xf2, 0x0c, 0x4d, 0x2e, 0xf2, 0xd7, 0xdf, 0x12, 0x6f, 0xd1, 0x4b, 0x30, 0x6f, 0x77, 0xdb,
	0x3d, 0x2a, 0x60, 0x54, 0x05, 0xbb, 0xdb, 0x8e, 0x48, 0x76, 0x41, 0x27, 0x81, 0xe5, 0x85, 0x16,
	0xf7, 0x36, 0xc3, 0x85, 0x85, 0x0b, 0x4b, 0x76, 0x51, 0xe1, 0xa1, 0x6f, 0xe9, 0xae, 0xc7, 0x0a,
	0x62, 0xbd, 0x56, 0x9a, 0x1f, 0x67, 0xd7, 0x13, 0xc4, 0xe8, 0xeb, 0x00, 0x27, 0x8e, 0xeb, 0x36,
	0x3a, 0x81, 0x63, 0xe3, 0xd2, 0xc2, 0x18, 0x07, 0x92, 0x3c, 0xa5, 0xdf, 0xa7, 0xe4, 0x68, 0x1b,
	0x16, 0x18, 0x73, 0xa4, 0xdf, 0xe2, 0x18, 0xfc, 0xf3, 0x94, 0x25, 0x52, 0x9f, 0xae, 0x8f, 0x71,
	0xc3, 0x6a, 0x33, 0xd0, 0x50, 0x1c, 0x6b, 0x7d, 0x8c, 0xb7, 0x19, 0x39, 0x7a, 0x03, 0xe6, 0x29,
	0xb3, 0xdd, 0x0d, 0x02, 0xec, 0xd9, 0x67, 0x25, 0x7d, 0x18, 0xf0, 0x2f, 0x9c, 0x60, 0xbc, 0x23,
	0x28, 0x68, 0x05, 0xa4, 0x1c, 0x2c, 0xbe, 0x96, 0x58, 0x84, 0x2c, 0xf5, 0xe2, 0xeb, 0x2e, 0xc6,
	0x2c, 0xb2, 0xb2, 0x27, 0xfc, 0x07, 0xda, 0x04, 0xba, 0x58, 0xe3, 0xd8, 0x0a, 0x9d, 0xb0, 0x84,
	0x18, 0x39, 0x8a, 0x91, 0xdf, 0xa1, 0x23, 0x26, 0x9d, 0x92, 0xfd, 0x1a, 0x9a, 0xb8, 0xcb, 0x13,
	0x27, 0xee, 0xf7, 0x35, 0x58, 0x10, 0x78, 0x89, 0xb7, 0x0b, 0x22, 0x70, 0xae, 0x8d, 0x0f, 0xce,
	0xa3, 0xd3, 0x50, 0xf2, 0x9c, 0xd3, 0x50, 0x09, 0xb2, 0xe2, 0x50, 0xcb, 0x92, 0x5a, 0x33, 0xe5,
	0xa3, 0xb1, 0x09, 0xf9, 0x08, 0x7e, 0x20, 0x03, 0xe6, 0x1e, 0xd1, 0x1f, 0x62, 0xfd, 0x5e, 0xaf,
	0xd6, 0x69, 0x62, 0x93, 0x0f, 0x19, 0xef, 0x80, 0xae, 0xe4, 0x18, 0xe7, 0x7b, 0x2d, 0xce, 0x37,
	0x34, 0x19, 0x05, 0xfb, 0x1f, 0x26, 0xa1, 0xc0, 0x5e, 0x7f, 0xe9, 0x51, 0xdc, 0x3b, 0x7d, 0x85,
	0x4a, 0xa2, 0xb9, 0x88, 0xa1, 0xdf, 0x8c, 0xb1, 0x7a, 0x65, 0xfc, 0xb3, 0x06, 0x48, 0xb5, 0xdf,
	0x4f, 0x6c, 0xc7, 0x08, 0xbd, 0x1e, 0x21, 0xfc, 0x39, 0xd9, 0x7a, 0x8d, 0xa9, 0xd8, 0x07, 0xea,
	0x3f, 0xd4, 0x20, 0xcf, 0xde, 0x4f, 0x61, 0x33, 0x56, 0xb6, 0xd6, 0x54, 0x7c, 0x6b, 0x1d, 0xb3,
	0x7d, 0x34, 0x9d, 0x1d, 0xf8, 0x5f, 0x34, 0xd0, 0x65, 0xb3, 0x29, 0xfc, 0x22, 0x5c, 0x9c, 0x2a,
	0xde, 0x9d, 0xbb, 0xa0, 0x1f, 0xf8, 0xaf, 0x1a, 0xcc, 0x4b, 0xdd, 0xa6, 0xe0, 0xb1, 0x58, 0xbf,
	0x2e, 0x75, 0x71, 0xbf, 0x6e, 0xd6, 0xe0, 0xe9, 0x9f, 0x34, 0x28, 0x8a, 0xc6, 0xdf, 0x73, 0xf7,
	0x5c, 0x54, 0xba, 0xd3, 0xe7, 0x94, 0xee, 0x09, 0x5c, 0xf5, 0x91, 0x06, 0x05, 0xa1, 0xcb, 0x14,
	0x3c, 0xa5, 0x76, 0x4a, 0x53, 0x17, 0x75, 0x4a, 0x67, 0xed, 0xa6, 0xbf, 0x4c, 0x41, 0xee, 0x1e,
	0x7e, 0xcc, 0x92, 0x17, 0x7d, 0x65, 0xa0, 0xf6, 0x6b, 0x0c, 0x67, 0x9e, 0x5b, 0xdd, 0xc7, 0xfe,
	0x3c, 0x84, 0x17, 0x6d, 0x05, 0x99, 0x2e, 0xc5, 0x4a, 0x07, 0xdb, 0x7a, 0xf3, 0xbe, 0xfc, 0x89,
	0xbe, 0x2a, 0x39, 0xd8, 0xae, 0x90, 0x19, 0xb2, 0x71, 0x72, 0x62, 0xfa, 0x13, 0xbd, 0x0d, 0x0b,
	0x14, 0xf8, 0x35, 0x1c, 0xaf, 0x71, 0xe2, 0x07, 0x36, 0x6f, 0x2a, 0x2b, 0x1b, 0x26, 0x85, 0x78,
	0x75, 0xef, 0x2e, 0x1d, 0x32, 0x0b, 0xa4, 0xf7, 0x80, 0xca, 0x90, 0x8b, 0xc0, 0x57, 0x8e, 0xed,
	0xe0, 0xd1, 0x33, 0xaa, 0xc2, 0x1c, 0x47, 0x75, 0xf9, 0x31, 0x50, 0x15, 0x27, 0x45, 0x07, 0x70,
	0xb9, 0x87, 0xc2, 0xb9, 0xfe, 0x36, 0xcf, 0x49, 0x60, 0x9f, 0x42, 0x44, 0x1f, 0x51, 0x44, 0x58,
	0xbc, 0xde, 0x23, 0x32, 0x2f, 0xe1, 0x21, 0x6f, 0x43, 0xa4, 0x43, 0x8a, 0x58, 0x2d, 0x86, 0x6a,
	0xf3, 0x26, 0xfd, 0x69, 0xfc, 0x40, 0x83, 0x4b, 0xd2, 0x6d, 0x07, 0x8e, 0xd7, 0x72, 0xf1, 0xd4,
	0xfb, 0x4b, 0x37, 0x61, 0x8e, 0xd9, 0xb7, 0xb7, 0x7d, 0xcb, 0xa0, 0x92, 0x0b, 0x9b, 0x9c, 0x80,
	0x6e, 0x3c, 0x97, 0xfb, 0x85, 0x99, 0xca, 0x71, 0x7b, 0xf4, 0x2e, 0xa4, 0x9c, 0xb6, 0xd2, 0xf1,
	0xd3, 0xd6, 0x74, 0xd2, 0xe2, 0x07, 0x1a, 0x2c, 0x4b, 0x95, 0xee, 0x74, 0xdd, 0x07, 0x53, 0xb7,
	0xee, 0x7a, 0xb4, 0x8d, 0x46, 0x97, 0xd1, 0x03, 0xe6, 0x15, 0x14, 0xd4, 0xbe, 0x2b, 0x71, 0x61,
	0xa6, 0x62, 0xdd, 0x32, 0xe4, 0x84, 0xd1, 0xb8, 0x10, 0x79, 0x33, 0x7a, 0x9e, 0x79, 0xd9, 0xf9,
	0xbe, 0x44, 0xab, 0x02, 0x9f, 0x7f, 0x5e, 0xb4, 0xba, 0x07, 0xcb, 0x7e, 0xe0, 0xb4, 0x1a, 0x9f,
	0x03, 0xb2, 0x2e, 0x51, 0xc6, 0x9d, 0x58, 0x61, 0xfb, 0x39, 0xa5, 0x18, 0xa4, 0xc6, 0xb9, 0x5a,
	0x1a, 0x2c, 0x15, 0xe9, 0xb1, 0x4b, 0x85, 0xf1, 0x77, 0x1a, 0x2c, 0x73, 0x47, 0xe3, 0x8e, 0x6b,
	0xd9, 0x78, 0x76, 0xdf, 0xe6, 0x29, 0x71, 0x99, 0xba, 0x18, 0x73, 0xc6, 0x2e, 0x61, 0xfb, 0x31,
	0x67, 0xdf, 0xf5, 0xeb, 0x0f, 0x35, 0x58, 0x89, 0xeb, 0x30, 0xad, 0xc4, 0x57, 0x7b, 0x25, 0x4a,
	0x7a, 0xcf, 0x3a, 0x30, 0xff, 0x41, 0x83, 0x55, 0x25, 0xd1, 0x7e, 0x72, 0x1c, 0xb3, 0x09, 0x59,
	0x6e, 0x73, 0x09, 0xbc, 0x47, 0x78, 0x46, 0x52, 0x19, 0x7f, 0xa3, 0x41, 0x69, 0x50, 0x93, 0x19,
	0xd7, 0xe5, 0x61, 0x4e, 0x48, 0x4f, 0xec, 0x84, 0xdf, 0x4c, 0x8a, 0x23, 0xda, 0x0e, 0x05, 0x41,
	0xee, 0x78, 0xf6, 0x7f, 0x2b, 0xde, 0x80, 0x7b, 0x86, 0x03, 0x6f, 0xea, 0x59, 0x0f, 0xbc, 0x53,
	0x3f, 0x2b, 0xfc, 0x85, 0x2c, 0x15, 0xd2, 0x20, 0x5f, 0x08, 0x37, 0xbe, 0xaf, 0xc1, 0x65, 0x26,
	0xf5, 0xbb, 0x56, 0x18, 0x4e, 0xe2, 0xca, 0x09, 0xf6, 0xd1, 0xd7, 0xfb, 0x6e, 0xc1, 0x2e, 0x38,
	0x23, 0xff, 0xb5, 0x4c, 0x6f, 0x55, 0xa4, 0x2f, 0x84, 0x31, 0xbf, 0x07, 0xcb, 0xf2, 0x1b, 0xaa,
	0x17, 0xf0, 0xa1, 0x8b, 0xf1, 0x81, 0x06, 0x2b, 0xf1, 0xf5, 0xa7, 0x62, 0xb5, 0x91, 0x1f, 0x03,
	0xcd, 0xba, 0xce, 0x7f, 0x17, 0x2a, 0xf1, 0x4f, 0x9d, 0x26, 0xbc, 0x9a, 0x5a, 0x87, 0xa5, 0xf8,
	0xe7, 0x71, 0x0d, 0x47, 0xea, 0x57, 0x8c, 0x7d, 0x0a, 0x57, 0xaf, 0x19, 0xff, 0xa7, 0xc1, 0xf5,
	0x91, 0x8b, 0x4d, 0xc5, 0x8c, 0x83, 0x1f, 0xeb, 0xa5, 0x26, 0xf8, 0x53, 0x8b, 0x59, 0xdb, 0xfa,
	0x57, 0x35, 0x28, 0xc5, 0x04, 0x78, 0x11, 0xb7, 0x55, 0xff, 0xa3, 0xc1, 0x95, 0x21, 0x42, 0x4c,
	0xc5, 0xfc, 0xb7, 0xa1, 0x18, 0x37, 0xbf, 0x44, 0xf5, 0x23, 0xec, 0xbf, 0x18, 0xb3, 0xff, 0xec,
	0xd1, 0xf6, 0x7f, 0x6a, 0xb0, 0x34, 0xa0, 0xfb, 0x97, 0x5b, 0xe7, 0xf5, 0x4f, 0x35, 0x58, 0x88,
	0xdd, 0x77, 0xa1, 0x2c, 0xa4, 0xee, 0xe1, 0xc7, 0x7a, 0x02, 0xe5, 0x20, 0x5d, 0xf3, 0x3d, 0xac,
	0x6b, 0x68, 0x1e, 0x72, 0x7c, 0x13, 0xc0, 0x4d, 0x3d, 0x49, 0x9f, 0x04, 0x4e, 0x6a, 0xea, 0x29,
	0xb4, 0x04, 0x0b, 0xfb, 0x98, 0x7d, 0x6a, 0xcd, 0x49, 0xf4, 0x34, 0x2a, 0x40, 0xf6, 0x80, 0xf8,
	0x9d, 0x0e, 0x6e, 0xea, 0x73, 0x9c, 0x9a, 0xae, 0x89, 0x9b, 0x7a, 0x06, 0x2d, 0x40, 0xfe, 0xa0,
	0x1b, 0x76, 0xb0, 0xd7, 0xc4, 0x4d, 0x3d, 0x8b, 0x16, 0x01, 0x04, 0x33, 0x5d, 0x32, 0x47, 0x9f,
	0x77, 0x2c, 0xd7, 0xee, 0xba, 0x16, 0x25, 0xcf, 0xd3, 0x99, 0x76, 0x9f, 0x74, 0x9c, 0x00, 0x37,
	0x75, 0x40, 0x08, 0x16, 0x05, 0xb1, 0x58, 0x5e, 0x2f, 0xa0, 0x3c, 0xcc, 0xb1, 0xdb, 0x7a, 0x7d,
	0x1e, 0x15, 0xc5, 0x51, 0x89, 0xf7, 0xab, 0xf5, 0x05, 0x3a, 0xd9, 0x01, 0x26, 0xc4, 0x65, 0x37,
	0x16, 0xfa, 0xe2, 0xfa, 0x7f, 0xcc, 0x41, 0xb1, 0xcf, 0x1e, 0x94, 0xff, 0x3e, 0x39, 0xc5, 0x81,
	0x9e, 0xa0, 0x8a, 0x1c, 0x79, 0x0f, 0x3c, 0xff, 0xb1, 0x77, 0x70, 0xd6, 0x3e, 0xf6, 0x5d, 0x5d,
	0x43, 0x97, 0x60, 0x49, 0xbe, 0x12, 0x17, 0xab, 0xf5, 0x9a, 0x9e, 0x44, 0x06, 0x54, 0x8e, 0xbc,
	0xb0, 0xdb, 0xe9, 0xf8, 0x01, 0xc1, 0x4d, 0x0e, 0x38, 0x4e, 0xad, 0xc0, 0xb2, 0x09, 0x0e, 0xd8,
	0xdf, 0xca, 0xe8, 0x29, 0xca, 0x5a, 0xf7, 0x6c, 0x3f, 0x08, 0xb0, 0x4d, 0xe4, 0x25, 0x94, 0x9e,
	0xa6, 0x3a, 0xec, 0x8a, 0xbf, 0x5a, 0xde, 0x71, 0xfd, 0x90, 0x59, 0x08, 0xc1, 0x62, 0xad, 0xdb,
	0x71, 0x1d, 0xdb, 0x22, 0x98, 0x4d, 0xa6, 0x67, 0xe8, 0xbb, 0xba, 0xf7, 0xc8, 0x72, 0x9d, 0xa6,
	0xd8, 0xca, 0xf5, 0x2c, 0x5a, 0x86, 0xe2, 0xa1, 0xef, 0xef, 0x59, 0x04, 0x1f, 0xfa, 0xc2, 0xd6,
	0x39, 0xa4, 0xc3, 0xbc, 0x10, 0x91, 0xb3, 0xe6, 0x51, 0x09, 0x56, 0xf8, 0xe8, 0xb6, 0x1b, 0x60,
	0xab, 0x79, 0x26, 0x6c, 0xa6, 0x03, 0x5a, 0x01, 0xbd, 0xe6, 0x9c, 0x9c, 0xe0, 0x00, 0x7b, 0x84,
	0xeb, 0x18, 0xea, 0x05, 0x65, 0x29, 0x51, 0x6b, 0xf4, 0x79, 0x4a, 0x29, 0xc5, 0xdc, 0xde, 0xaf,
	0xef, 0x06, 0x81, 0x1f, 0xe8, 0x0b, 0x74, 0x2d, 0x41, 0xc9, 0xd7, 0x5a, 0xa4, 0x5a, 0x9a, 0x16,
	0xc1, 0x7b, 0x4e, 0xdb, 0x21, 0xbb, 0x4f, 0x6c, 0x8c, 0xa9, 0x5b, 0x8b, 0xe8, 0x2a, 0xac, 0x2a,
	0x06, 0x3a, 0xe0, 0xd5, 0xa5, 0x43, 0xcd, 0xae, 0xeb, 0x94, 0xe7, 0x5d, 0x27, 0x0c, 0x1d, 0xaf,
	0xd5, 0x03, 0x7c, 0xfa, 0x12, 0x8d, 0x8c, 0x6f, 0x1e, 0x1e, 0xee, 0xf3, 0xb5, 0x10, 0x5a, 0x85,
	0xe5, 0x7b, 0xec, 0x36, 0x99, 0x3a, 0xda, 0x3a, 0x76, 0x85, 0x65, 0x96, 0xd1, 0x65, 0x40, 0xf7,
	0x7c, 0x8f, 0x6b, 0xd8, 0x7b, 0xbf, 0xc2, 0x7d, 0x15, 0xad, 0xc9, 0xf1, 0x8c, 0x7e, 0x89, 0x5a,
	0xa3, 0xdf, 0x57, 0x34, 0xca, 0xf5, 0xcb, 0xe8, 0x3a, 0x5c, 0x1d, 0x18, 0xe9, 0x75, 0xbc, 0xf4,
	0x55, 0xba, 0x92, 0x42, 0x20, 0x8d, 0x53, 0xa2, 0x06, 0x13, 0x26, 0x97, 0xbe, 0xb9, 0x42, 0x7d,
	0x23, 0xde, 0x49, 0xbb, 0xe9, 0x65, 0x6a, 0xaf, 0x5f, 0x0c, 0xac, 0xce, 0xe9, 0xb7, 0xf6, 0xb8,
	0x56, 0x57, 0x69, 0x32, 0x6c, 0xdf, 0x11, 0xf6, 0x5c, 0xa3, 0x11, 0xbb, 0x4b, 0x4e, 0xcd, 0xfd,
	0x1d, 0xfe, 0xe2, 0x5a, 0x6f, 0x66, 0x59, 0x13, 0xf4, 0x0a, 0x55, 0x40, 0x98, 0x2b, 0x56, 0x28,
	0xf4, 0xeb, 0x5c, 0x35, 0xb6, 0x66, 0x7c, 0xe4, 0xc6, 0x7a, 0x1d, 0xb2, 0xe2, 0x8e, 0x91, 0x26,
	0x81, 0x89, 0x5b, 0x34, 0xa1, 0xfc, 0xe0, 0x4c, 0x4f, 0xd0, 0xec, 0x3e, 0xb4, 0x9e, 0xe8, 0x1a,
	0x95, 0x78, 0xcf, 0xb7, 0x2d, 0x77, 0xc7, 0x6f, 0xb7, 0xe9, 0xfc, 0xbe, 0xa7, 0x27, 0xa9, 0xc4,
	0x52, 0xfe, 0xbb, 0x18, 0x87, 0x7a, 0x6a, 0xfd, 0x6b, 0x90, 0x93, 0xf7, 0x8f, 0x4c, 0xfa, 0xe3,
	0xd0, 0x77, 0xbb, 0x04, 0xeb, 0x09, 0x9a, 0x9b, 0xfb, 0x38, 0x38, 0xf2, 0x1c, 0xa2, 0x6b, 0x3c,
	0x91, 0x03, 0x1b, 0x7b, 0xc4, 0x6a, 0x61, 0x3d, 0xb9, 0xfe, 0xbe, 0x16, 0x7d, 0xc5, 0xab, 0xde,
	0x19, 0x52, 0x26, 0x21, 0xb3, 0x9e, 0xa0, 0x4c, 0xe2, 0xcf, 0x35, 0xee, 0x62, 0x51, 0x66, 0x22,
	0xeb, 0x25, 0x29, 0x69, 0x0d, 0xb3, 0x46, 0xbc, 0x9e, 0xa2, 0xa4, 0xdf, 0x76, 0xc8, 0x69, 0x33,
	0xb0, 0x1e, 0x5b, 0xb4, 0xc4, 0x14, 0xa1, 0x60, 0x62, 0xcb, 0x75, 0x7e, 0x19, 0x37, 0xf7, 0x3d,
	0x57, 0x9f, 0x63, 0x95, 0xa3, 0xa7, 0x49, 0x86, 0x6a, 0xf2, 0x6d, 0xec, 0xda, 0x7e, 0x1b, 0xdf,
	0xf1, 0xbd, 0x6e, 0xa8, 0x67, 0xab, 0x7f, 0x92, 0xeb, 0x05, 0xf5, 0xae, 0xf8, 0x97, 0x05, 0xe8,
	0x08, 0xa0, 0xf7, 0x77, 0x3b, 0xe8, 0x6a, 0xaf, 0x70, 0x0e, 0xfc, 0x21, 0x5a, 0xf9, 0x95, 0x61,
	0x83, 0x83, 0xdf, 0xa8, 0x1a, 0x89, 0x37, 0x34, 0xf4, 0x1d, 0x28, 0x28, 0x1f, 0xb2, 0xa2, 0xb5,
	0x81, 0x8b, 0x54, 0x75, 0xe2, 0x57, 0x87, 0x8e, 0x8e, 0x98, 0xd9, 0x02, 0x34, 0xf8, 0xd9, 0x0e,
	0x7a, 0x59, 0xb9, 0x06, 0x1c, 0xf5, 0x11, 0x51, 0xf9, 0x67, 0xce, 0x27, 0xe2, 0x5b, 0x96, 0x91,
	0x40, 0xbb, 0xb4, 0x6e, 0x46, 0x5f, 0x8f, 0x5c, 0x1b, 0xe4, 0x52, 0x70, 0x49, 0xf9, 0xf2, 0xf0,
	0x61, 0x23, 0x81, 0xde, 0x81, 0xcc, 0x7d, 0x7e, 0xc9, 0xb5, 0x36, 0xf4, 0xce, 0x51, 0xce, 0xb0,
	0xdc, 0x37, 0x2a, 0xd8, 0xb7, 0x21, 0x1f, 0x5d, 0x69, 0x21, 0xe5, 0xd6, 0xb2, 0xff, 0x9e, 0x4b,
	0x95, 0x40, 0xbd, 0x27, 0x32, 0x12, 0xe8, 0x36, 0xe4, 0xe4, 0xd5, 0x0a, 0x52, 0xf6, 0xc4, 0xbe,
	0xeb, 0x96, 0xf2, 0xa5, 0x81, 0x21, 0xc1, 0x7f, 0x04, 0x8b, 0xf1, 0x7e, 0x2d, 0xba, 0x3e, 0xd8,
	0x7e, 0x8c, 0xb5, 0x95, 0xcb, 0x37, 0x46, 0x13, 0x44, 0xf6, 0xbd, 0x0f, 0xf3, 0x6a, 0x9b, 0x52,
	0xb5, 0xf0, 0x90, 0x5e, 0x6a, 0xb9, 0x32, 0x6a, 0x58, 0x9d, 0x50, 0x6d, 0x2e, 0xa9, 0x13, 0x0e,
	0x69, 0x9c, 0xa9, 0x13, 0x0e, 0xeb, 0x49, 0x19, 0x09, 0xf4, 0x4b, 0xe2, 0x8e, 0x5d, 0x69, 0x89,
	0xa0, 0x97, 0xfa, 0xb8, 0x06, 0x1b, 0x3f, 0x65, 0xe3, 0x3c, 0x92, 0x68, 0xf2, 0x3d, 0xb1, 0x4f,
	0xf3, 0x12, 0x3e, 0x10, 0x1c, 0xb1, 0x03, 0x70, 0xf9, 0xda, 0x88, 0xd1, 0x68, 0xb6, 0xef, 0x40,
	0xb1, 0xef, 0xa0, 0x8a, 0x6e, 0xf4, 0xf1, 0x0c, 0x1c, 0xab, 0xcb, 0x2f, 0x9d, 0x43, 0x21, 0x67,
	0xbe, 0xf3, 0xe6, 0x87, 0x1f, 0x57, 0x12, 0x1f, 0x7d, 0x5c, 0x49, 0x7c, 0xf6, 0x71, 0x45, 0xfb,
	0x95, 0xa7, 0x15, 0xed, 0x8f, 0x9f, 0x56, 0xb4, 0x7f, 0x7c, 0x5a, 0xd1, 0x3e, 0x7c, 0x5a, 0xd1,
	0x7e, 0xf4, 0xb4, 0xa2, 0x7d, 0xfa, 0xb4, 0x92, 0xf8, 0xec, 0x69, 0x45, 0x7b, 0xff, 0x93, 0x4a,
	0xe2, 0xc3, 0x4f, 0x2a, 0x89, 0x8f, 0x3e, 0xa9, 0x24, 0x8e, 0x33, 0x0c, 0xbf, 0x6d, 0xfd, 0x7f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x48, 0x74, 0x6b, 0x18, 0x16, 0x45, 0x00, 0x00,
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
func (x AccountMovementType) String() string {
	s, ok := AccountMovementType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (this *HistoricalOpenInterestsRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HistoricalOpenInterestsRequest)
	if !ok {
		that2, ok := that.(HistoricalOpenInterestsRequest)
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
func (this *HistoricalOpenInterestsResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HistoricalOpenInterestsResponse)
	if !ok {
		that2, ok := that.(HistoricalOpenInterestsResponse)
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
	if len(this.Interests) != len(that1.Interests) {
		return false
	}
	for i := range this.Interests {
		if !this.Interests[i].Equal(that1.Interests[i]) {
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
func (this *HistoricalUnipoolV3DataRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HistoricalUnipoolV3DataRequest)
	if !ok {
		that2, ok := that.(HistoricalUnipoolV3DataRequest)
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
	if this.Start != that1.Start {
		return false
	}
	if this.End != that1.End {
		return false
	}
	return true
}
func (this *HistoricalUnipoolV3DataResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HistoricalUnipoolV3DataResponse)
	if !ok {
		that2, ok := that.(HistoricalUnipoolV3DataResponse)
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
	if len(this.Events) != len(that1.Events) {
		return false
	}
	for i := range this.Events {
		if !this.Events[i].Equal(that1.Events[i]) {
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
func (this *HistoricalProtocolAssetTransferRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HistoricalProtocolAssetTransferRequest)
	if !ok {
		that2, ok := that.(HistoricalProtocolAssetTransferRequest)
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
	if !this.ProtocolAsset.Equal(that1.ProtocolAsset) {
		return false
	}
	if this.Start != that1.Start {
		return false
	}
	if this.Stop != that1.Stop {
		return false
	}
	return true
}
func (this *HistoricalProtocolAssetTransferResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HistoricalProtocolAssetTransferResponse)
	if !ok {
		that2, ok := that.(HistoricalProtocolAssetTransferResponse)
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
	if len(this.Update) != len(that1.Update) {
		return false
	}
	for i := range this.Update {
		if !this.Update[i].Equal(that1.Update[i]) {
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
	if len(this.Stats) != len(that1.Stats) {
		return false
	}
	for i := range this.Stats {
		if this.Stats[i] != that1.Stats[i] {
			return false
		}
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
func (this *UnipoolV3DataRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*UnipoolV3DataRequest)
	if !ok {
		that2, ok := that.(UnipoolV3DataRequest)
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
	return true
}
func (this *UnipoolV3DataResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*UnipoolV3DataResponse)
	if !ok {
		that2, ok := that.(UnipoolV3DataResponse)
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
	if len(this.Update) != len(that1.Update) {
		return false
	}
	for i := range this.Update {
		if !this.Update[i].Equal(that1.Update[i]) {
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
func (this *ProtocolAssetTransferRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ProtocolAssetTransferRequest)
	if !ok {
		that2, ok := that.(ProtocolAssetTransferRequest)
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
	if !this.ProtocolAsset.Equal(that1.ProtocolAsset) {
		return false
	}
	return true
}
func (this *ProtocolAssetTransferResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ProtocolAssetTransferResponse)
	if !ok {
		that2, ok := that.(ProtocolAssetTransferResponse)
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
	if len(this.Update) != len(that1.Update) {
		return false
	}
	for i := range this.Update {
		if !this.Update[i].Equal(that1.Update[i]) {
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
func (this *UnipoolV3DataIncrementalRefresh) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*UnipoolV3DataIncrementalRefresh)
	if !ok {
		that2, ok := that.(UnipoolV3DataIncrementalRefresh)
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
	if !this.Update.Equal(that1.Update) {
		return false
	}
	return true
}
func (this *ProtocolAssetDataIncrementalRefresh) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ProtocolAssetDataIncrementalRefresh)
	if !ok {
		that2, ok := that.(ProtocolAssetDataIncrementalRefresh)
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
	if !this.Update.Equal(that1.Update) {
		return false
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
func (this *AccountMovement) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AccountMovement)
	if !ok {
		that2, ok := that.(AccountMovement)
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
	if !this.Asset.Equal(that1.Asset) {
		return false
	}
	if this.Change != that1.Change {
		return false
	}
	if this.Type != that1.Type {
		return false
	}
	if this.Subtype != that1.Subtype {
		return false
	}
	if this.MovementID != that1.MovementID {
		return false
	}
	if !this.Time.Equal(that1.Time) {
		return false
	}
	return true
}
func (this *AccountMovementRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AccountMovementRequest)
	if !ok {
		that2, ok := that.(AccountMovementRequest)
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
	if this.Type != that1.Type {
		return false
	}
	if !this.Filter.Equal(that1.Filter) {
		return false
	}
	return true
}
func (this *AccountMovementFilter) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AccountMovementFilter)
	if !ok {
		that2, ok := that.(AccountMovementFilter)
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
func (this *AccountMovementResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AccountMovementResponse)
	if !ok {
		that2, ok := that.(AccountMovementResponse)
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
	if len(this.Movements) != len(that1.Movements) {
		return false
	}
	for i := range this.Movements {
		if !this.Movements[i].Equal(that1.Movements[i]) {
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
	if !this.FromID.Equal(that1.FromID) {
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
	if this.Type != that1.Type {
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
	if this.ResponseID != that1.ResponseID {
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
	if this.Tag != that1.Tag {
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
func (this *TransferDataRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*TransferDataRequest)
	if !ok {
		that2, ok := that.(TransferDataRequest)
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
func (this *TransferDataResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*TransferDataResponse)
	if !ok {
		that2, ok := that.(TransferDataResponse)
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
	if this.Success != that1.Success {
		return false
	}
	if this.RejectionReason != that1.RejectionReason {
		return false
	}
	return true
}
func (this *ProtocolAssetDefinitionRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ProtocolAssetDefinitionRequest)
	if !ok {
		that2, ok := that.(ProtocolAssetDefinitionRequest)
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
	if this.ProtocolAssetID != that1.ProtocolAssetID {
		return false
	}
	return true
}
func (this *ProtocolAssetDefinitionResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ProtocolAssetDefinitionResponse)
	if !ok {
		that2, ok := that.(ProtocolAssetDefinitionResponse)
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
	if !this.ProtocolAsset.Equal(that1.ProtocolAsset) {
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
func (this *ProtocolAssetListRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ProtocolAssetListRequest)
	if !ok {
		that2, ok := that.(ProtocolAssetListRequest)
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
func (this *ProtocolAssetListResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ProtocolAssetListResponse)
	if !ok {
		that2, ok := that.(ProtocolAssetListResponse)
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
	if len(this.ProtocolAssets) != len(that1.ProtocolAssets) {
		return false
	}
	for i := range this.ProtocolAssets {
		if !this.ProtocolAssets[i].Equal(that1.ProtocolAssets[i]) {
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
func (this *ProtocolAssetList) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ProtocolAssetList)
	if !ok {
		that2, ok := that.(ProtocolAssetList)
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
	if len(this.ProtocolAssets) != len(that1.ProtocolAssets) {
		return false
	}
	for i := range this.ProtocolAssets {
		if !this.ProtocolAssets[i].Equal(that1.ProtocolAssets[i]) {
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
func (this *HistoricalOpenInterestsRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.HistoricalOpenInterestsRequest{")
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
func (this *HistoricalOpenInterestsResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.HistoricalOpenInterestsResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Interests != nil {
		s = append(s, "Interests: "+fmt.Sprintf("%#v", this.Interests)+",\n")
	}
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
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
func (this *HistoricalUnipoolV3DataRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.HistoricalUnipoolV3DataRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	s = append(s, "Start: "+fmt.Sprintf("%#v", this.Start)+",\n")
	s = append(s, "End: "+fmt.Sprintf("%#v", this.End)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *HistoricalUnipoolV3DataResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 10)
	s = append(s, "&messages.HistoricalUnipoolV3DataResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Events != nil {
		s = append(s, "Events: "+fmt.Sprintf("%#v", this.Events)+",\n")
	}
	s = append(s, "SeqNum: "+fmt.Sprintf("%#v", this.SeqNum)+",\n")
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *HistoricalProtocolAssetTransferRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.HistoricalProtocolAssetTransferRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	if this.ProtocolAsset != nil {
		s = append(s, "ProtocolAsset: "+fmt.Sprintf("%#v", this.ProtocolAsset)+",\n")
	}
	s = append(s, "Start: "+fmt.Sprintf("%#v", this.Start)+",\n")
	s = append(s, "Stop: "+fmt.Sprintf("%#v", this.Stop)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *HistoricalProtocolAssetTransferResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 10)
	s = append(s, "&messages.HistoricalProtocolAssetTransferResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Update != nil {
		s = append(s, "Update: "+fmt.Sprintf("%#v", this.Update)+",\n")
	}
	s = append(s, "SeqNum: "+fmt.Sprintf("%#v", this.SeqNum)+",\n")
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
	s := make([]string, 0, 10)
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
	s = append(s, "Stats: "+fmt.Sprintf("%#v", this.Stats)+",\n")
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
func (this *UnipoolV3DataRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.UnipoolV3DataRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Subscribe: "+fmt.Sprintf("%#v", this.Subscribe)+",\n")
	if this.Subscriber != nil {
		s = append(s, "Subscriber: "+fmt.Sprintf("%#v", this.Subscriber)+",\n")
	}
	if this.Instrument != nil {
		s = append(s, "Instrument: "+fmt.Sprintf("%#v", this.Instrument)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *UnipoolV3DataResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 10)
	s = append(s, "&messages.UnipoolV3DataResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Update != nil {
		s = append(s, "Update: "+fmt.Sprintf("%#v", this.Update)+",\n")
	}
	s = append(s, "SeqNum: "+fmt.Sprintf("%#v", this.SeqNum)+",\n")
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ProtocolAssetTransferRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.ProtocolAssetTransferRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Subscribe: "+fmt.Sprintf("%#v", this.Subscribe)+",\n")
	if this.Subscriber != nil {
		s = append(s, "Subscriber: "+fmt.Sprintf("%#v", this.Subscriber)+",\n")
	}
	if this.ProtocolAsset != nil {
		s = append(s, "ProtocolAsset: "+fmt.Sprintf("%#v", this.ProtocolAsset)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ProtocolAssetTransferResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 10)
	s = append(s, "&messages.ProtocolAssetTransferResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Update != nil {
		s = append(s, "Update: "+fmt.Sprintf("%#v", this.Update)+",\n")
	}
	s = append(s, "SeqNum: "+fmt.Sprintf("%#v", this.SeqNum)+",\n")
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *UnipoolV3DataIncrementalRefresh) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.UnipoolV3DataIncrementalRefresh{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	s = append(s, "SeqNum: "+fmt.Sprintf("%#v", this.SeqNum)+",\n")
	if this.Update != nil {
		s = append(s, "Update: "+fmt.Sprintf("%#v", this.Update)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ProtocolAssetDataIncrementalRefresh) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.ProtocolAssetDataIncrementalRefresh{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	s = append(s, "SeqNum: "+fmt.Sprintf("%#v", this.SeqNum)+",\n")
	if this.Update != nil {
		s = append(s, "Update: "+fmt.Sprintf("%#v", this.Update)+",\n")
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
func (this *AccountMovement) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 10)
	s = append(s, "&messages.AccountMovement{")
	if this.Asset != nil {
		s = append(s, "Asset: "+fmt.Sprintf("%#v", this.Asset)+",\n")
	}
	s = append(s, "Change: "+fmt.Sprintf("%#v", this.Change)+",\n")
	s = append(s, "Type: "+fmt.Sprintf("%#v", this.Type)+",\n")
	s = append(s, "Subtype: "+fmt.Sprintf("%#v", this.Subtype)+",\n")
	s = append(s, "MovementID: "+fmt.Sprintf("%#v", this.MovementID)+",\n")
	if this.Time != nil {
		s = append(s, "Time: "+fmt.Sprintf("%#v", this.Time)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *AccountMovementRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&messages.AccountMovementRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	if this.Account != nil {
		s = append(s, "Account: "+fmt.Sprintf("%#v", this.Account)+",\n")
	}
	s = append(s, "Type: "+fmt.Sprintf("%#v", this.Type)+",\n")
	if this.Filter != nil {
		s = append(s, "Filter: "+fmt.Sprintf("%#v", this.Filter)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *AccountMovementFilter) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.AccountMovementFilter{")
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
func (this *AccountMovementResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.AccountMovementResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.Movements != nil {
		s = append(s, "Movements: "+fmt.Sprintf("%#v", this.Movements)+",\n")
	}
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
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
	s := make([]string, 0, 11)
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
	if this.FromID != nil {
		s = append(s, "FromID: "+fmt.Sprintf("%#v", this.FromID)+",\n")
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
	s = append(s, "Type: "+fmt.Sprintf("%#v", this.Type)+",\n")
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
	s := make([]string, 0, 9)
	s = append(s, "&messages.OrderList{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
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
	s := make([]string, 0, 13)
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
	s = append(s, "Tag: "+fmt.Sprintf("%#v", this.Tag)+",\n")
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
func (this *TransferDataRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.TransferDataRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Subscribe: "+fmt.Sprintf("%#v", this.Subscribe)+",\n")
	if this.Subscriber != nil {
		s = append(s, "Subscriber: "+fmt.Sprintf("%#v", this.Subscriber)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *TransferDataResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.TransferDataResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	s = append(s, "SeqNum: "+fmt.Sprintf("%#v", this.SeqNum)+",\n")
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ProtocolAssetDefinitionRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&messages.ProtocolAssetDefinitionRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ProtocolAssetID: "+fmt.Sprintf("%#v", this.ProtocolAssetID)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ProtocolAssetDefinitionResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.ProtocolAssetDefinitionResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.ProtocolAsset != nil {
		s = append(s, "ProtocolAsset: "+fmt.Sprintf("%#v", this.ProtocolAsset)+",\n")
	}
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ProtocolAssetListRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&messages.ProtocolAssetListRequest{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "Subscribe: "+fmt.Sprintf("%#v", this.Subscribe)+",\n")
	if this.Subscriber != nil {
		s = append(s, "Subscriber: "+fmt.Sprintf("%#v", this.Subscriber)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ProtocolAssetListResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.ProtocolAssetListResponse{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.ProtocolAssets != nil {
		s = append(s, "ProtocolAssets: "+fmt.Sprintf("%#v", this.ProtocolAssets)+",\n")
	}
	s = append(s, "Success: "+fmt.Sprintf("%#v", this.Success)+",\n")
	s = append(s, "RejectionReason: "+fmt.Sprintf("%#v", this.RejectionReason)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ProtocolAssetList) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&messages.ProtocolAssetList{")
	s = append(s, "RequestID: "+fmt.Sprintf("%#v", this.RequestID)+",\n")
	s = append(s, "ResponseID: "+fmt.Sprintf("%#v", this.ResponseID)+",\n")
	if this.ProtocolAssets != nil {
		s = append(s, "ProtocolAssets: "+fmt.Sprintf("%#v", this.ProtocolAssets)+",\n")
	}
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

// ExchangeExecutorClient is the client API for ExchangeExecutor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ExchangeExecutorClient interface {
	MarketData(ctx context.Context, in *MarketDataRequest, opts ...grpc.CallOption) (ExchangeExecutor_MarketDataClient, error)
	AccountData(ctx context.Context, in *AccountDataRequest, opts ...grpc.CallOption) (ExchangeExecutor_AccountDataClient, error)
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

type exchangeExecutorClient struct {
	cc *grpc.ClientConn
}

func NewExchangeExecutorClient(cc *grpc.ClientConn) ExchangeExecutorClient {
	return &exchangeExecutorClient{cc}
}

func (c *exchangeExecutorClient) MarketData(ctx context.Context, in *MarketDataRequest, opts ...grpc.CallOption) (ExchangeExecutor_MarketDataClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ExchangeExecutor_serviceDesc.Streams[0], "/messages.ExchangeExecutor/MarketData", opts...)
	if err != nil {
		return nil, err
	}
	x := &exchangeExecutorMarketDataClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ExchangeExecutor_MarketDataClient interface {
	Recv() (*MarketDataIncrementalRefresh, error)
	grpc.ClientStream
}

type exchangeExecutorMarketDataClient struct {
	grpc.ClientStream
}

func (x *exchangeExecutorMarketDataClient) Recv() (*MarketDataIncrementalRefresh, error) {
	m := new(MarketDataIncrementalRefresh)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *exchangeExecutorClient) AccountData(ctx context.Context, in *AccountDataRequest, opts ...grpc.CallOption) (ExchangeExecutor_AccountDataClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ExchangeExecutor_serviceDesc.Streams[1], "/messages.ExchangeExecutor/AccountData", opts...)
	if err != nil {
		return nil, err
	}
	x := &exchangeExecutorAccountDataClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ExchangeExecutor_AccountDataClient interface {
	Recv() (*AccountDataIncrementalRefresh, error)
	grpc.ClientStream
}

type exchangeExecutorAccountDataClient struct {
	grpc.ClientStream
}

func (x *exchangeExecutorAccountDataClient) Recv() (*AccountDataIncrementalRefresh, error) {
	m := new(AccountDataIncrementalRefresh)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *exchangeExecutorClient) SecurityDefinition(ctx context.Context, in *SecurityDefinitionRequest, opts ...grpc.CallOption) (*SecurityDefinitionResponse, error) {
	out := new(SecurityDefinitionResponse)
	err := c.cc.Invoke(ctx, "/messages.ExchangeExecutor/SecurityDefinition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exchangeExecutorClient) Securities(ctx context.Context, in *SecurityListRequest, opts ...grpc.CallOption) (*SecurityList, error) {
	out := new(SecurityList)
	err := c.cc.Invoke(ctx, "/messages.ExchangeExecutor/Securities", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exchangeExecutorClient) Orders(ctx context.Context, in *OrderStatusRequest, opts ...grpc.CallOption) (*OrderList, error) {
	out := new(OrderList)
	err := c.cc.Invoke(ctx, "/messages.ExchangeExecutor/Orders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exchangeExecutorClient) Positions(ctx context.Context, in *PositionsRequest, opts ...grpc.CallOption) (*PositionList, error) {
	out := new(PositionList)
	err := c.cc.Invoke(ctx, "/messages.ExchangeExecutor/Positions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exchangeExecutorClient) Balances(ctx context.Context, in *BalancesRequest, opts ...grpc.CallOption) (*BalanceList, error) {
	out := new(BalanceList)
	err := c.cc.Invoke(ctx, "/messages.ExchangeExecutor/Balances", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exchangeExecutorClient) NewOrderSingle(ctx context.Context, in *NewOrderSingleRequest, opts ...grpc.CallOption) (*NewOrderSingleResponse, error) {
	out := new(NewOrderSingleResponse)
	err := c.cc.Invoke(ctx, "/messages.ExchangeExecutor/NewOrderSingle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exchangeExecutorClient) NewOrderBulk(ctx context.Context, in *NewOrderBulkRequest, opts ...grpc.CallOption) (*NewOrderBulkResponse, error) {
	out := new(NewOrderBulkResponse)
	err := c.cc.Invoke(ctx, "/messages.ExchangeExecutor/NewOrderBulk", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exchangeExecutorClient) OrderReplace(ctx context.Context, in *OrderReplaceRequest, opts ...grpc.CallOption) (*OrderReplaceResponse, error) {
	out := new(OrderReplaceResponse)
	err := c.cc.Invoke(ctx, "/messages.ExchangeExecutor/OrderReplace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exchangeExecutorClient) OrderBulkReplace(ctx context.Context, in *OrderBulkReplaceRequest, opts ...grpc.CallOption) (*OrderBulkReplaceResponse, error) {
	out := new(OrderBulkReplaceResponse)
	err := c.cc.Invoke(ctx, "/messages.ExchangeExecutor/OrderBulkReplace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exchangeExecutorClient) OrderCancel(ctx context.Context, in *OrderCancelRequest, opts ...grpc.CallOption) (*OrderCancelResponse, error) {
	out := new(OrderCancelResponse)
	err := c.cc.Invoke(ctx, "/messages.ExchangeExecutor/OrderCancel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exchangeExecutorClient) OrderMassCancel(ctx context.Context, in *OrderMassCancelRequest, opts ...grpc.CallOption) (*OrderMassCancelResponse, error) {
	out := new(OrderMassCancelResponse)
	err := c.cc.Invoke(ctx, "/messages.ExchangeExecutor/OrderMassCancel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExchangeExecutorServer is the server API for ExchangeExecutor service.
type ExchangeExecutorServer interface {
	MarketData(*MarketDataRequest, ExchangeExecutor_MarketDataServer) error
	AccountData(*AccountDataRequest, ExchangeExecutor_AccountDataServer) error
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

// UnimplementedExchangeExecutorServer can be embedded to have forward compatible implementations.
type UnimplementedExchangeExecutorServer struct {
}

func (*UnimplementedExchangeExecutorServer) MarketData(req *MarketDataRequest, srv ExchangeExecutor_MarketDataServer) error {
	return status.Errorf(codes.Unimplemented, "method MarketData not implemented")
}
func (*UnimplementedExchangeExecutorServer) AccountData(req *AccountDataRequest, srv ExchangeExecutor_AccountDataServer) error {
	return status.Errorf(codes.Unimplemented, "method AccountData not implemented")
}
func (*UnimplementedExchangeExecutorServer) SecurityDefinition(ctx context.Context, req *SecurityDefinitionRequest) (*SecurityDefinitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SecurityDefinition not implemented")
}
func (*UnimplementedExchangeExecutorServer) Securities(ctx context.Context, req *SecurityListRequest) (*SecurityList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Securities not implemented")
}
func (*UnimplementedExchangeExecutorServer) Orders(ctx context.Context, req *OrderStatusRequest) (*OrderList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Orders not implemented")
}
func (*UnimplementedExchangeExecutorServer) Positions(ctx context.Context, req *PositionsRequest) (*PositionList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Positions not implemented")
}
func (*UnimplementedExchangeExecutorServer) Balances(ctx context.Context, req *BalancesRequest) (*BalanceList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Balances not implemented")
}
func (*UnimplementedExchangeExecutorServer) NewOrderSingle(ctx context.Context, req *NewOrderSingleRequest) (*NewOrderSingleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewOrderSingle not implemented")
}
func (*UnimplementedExchangeExecutorServer) NewOrderBulk(ctx context.Context, req *NewOrderBulkRequest) (*NewOrderBulkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewOrderBulk not implemented")
}
func (*UnimplementedExchangeExecutorServer) OrderReplace(ctx context.Context, req *OrderReplaceRequest) (*OrderReplaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderReplace not implemented")
}
func (*UnimplementedExchangeExecutorServer) OrderBulkReplace(ctx context.Context, req *OrderBulkReplaceRequest) (*OrderBulkReplaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderBulkReplace not implemented")
}
func (*UnimplementedExchangeExecutorServer) OrderCancel(ctx context.Context, req *OrderCancelRequest) (*OrderCancelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderCancel not implemented")
}
func (*UnimplementedExchangeExecutorServer) OrderMassCancel(ctx context.Context, req *OrderMassCancelRequest) (*OrderMassCancelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderMassCancel not implemented")
}

func RegisterExchangeExecutorServer(s *grpc.Server, srv ExchangeExecutorServer) {
	s.RegisterService(&_ExchangeExecutor_serviceDesc, srv)
}

func _ExchangeExecutor_MarketData_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(MarketDataRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExchangeExecutorServer).MarketData(m, &exchangeExecutorMarketDataServer{stream})
}

type ExchangeExecutor_MarketDataServer interface {
	Send(*MarketDataIncrementalRefresh) error
	grpc.ServerStream
}

type exchangeExecutorMarketDataServer struct {
	grpc.ServerStream
}

func (x *exchangeExecutorMarketDataServer) Send(m *MarketDataIncrementalRefresh) error {
	return x.ServerStream.SendMsg(m)
}

func _ExchangeExecutor_AccountData_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(AccountDataRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExchangeExecutorServer).AccountData(m, &exchangeExecutorAccountDataServer{stream})
}

type ExchangeExecutor_AccountDataServer interface {
	Send(*AccountDataIncrementalRefresh) error
	grpc.ServerStream
}

type exchangeExecutorAccountDataServer struct {
	grpc.ServerStream
}

func (x *exchangeExecutorAccountDataServer) Send(m *AccountDataIncrementalRefresh) error {
	return x.ServerStream.SendMsg(m)
}

func _ExchangeExecutor_SecurityDefinition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecurityDefinitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangeExecutorServer).SecurityDefinition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.ExchangeExecutor/SecurityDefinition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangeExecutorServer).SecurityDefinition(ctx, req.(*SecurityDefinitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExchangeExecutor_Securities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecurityListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangeExecutorServer).Securities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.ExchangeExecutor/Securities",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangeExecutorServer).Securities(ctx, req.(*SecurityListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExchangeExecutor_Orders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangeExecutorServer).Orders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.ExchangeExecutor/Orders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangeExecutorServer).Orders(ctx, req.(*OrderStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExchangeExecutor_Positions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PositionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangeExecutorServer).Positions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.ExchangeExecutor/Positions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangeExecutorServer).Positions(ctx, req.(*PositionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExchangeExecutor_Balances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangeExecutorServer).Balances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.ExchangeExecutor/Balances",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangeExecutorServer).Balances(ctx, req.(*BalancesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExchangeExecutor_NewOrderSingle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewOrderSingleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangeExecutorServer).NewOrderSingle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.ExchangeExecutor/NewOrderSingle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangeExecutorServer).NewOrderSingle(ctx, req.(*NewOrderSingleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExchangeExecutor_NewOrderBulk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewOrderBulkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangeExecutorServer).NewOrderBulk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.ExchangeExecutor/NewOrderBulk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangeExecutorServer).NewOrderBulk(ctx, req.(*NewOrderBulkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExchangeExecutor_OrderReplace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderReplaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangeExecutorServer).OrderReplace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.ExchangeExecutor/OrderReplace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangeExecutorServer).OrderReplace(ctx, req.(*OrderReplaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExchangeExecutor_OrderBulkReplace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderBulkReplaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangeExecutorServer).OrderBulkReplace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.ExchangeExecutor/OrderBulkReplace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangeExecutorServer).OrderBulkReplace(ctx, req.(*OrderBulkReplaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExchangeExecutor_OrderCancel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderCancelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangeExecutorServer).OrderCancel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.ExchangeExecutor/OrderCancel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangeExecutorServer).OrderCancel(ctx, req.(*OrderCancelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExchangeExecutor_OrderMassCancel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderMassCancelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangeExecutorServer).OrderMassCancel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.ExchangeExecutor/OrderMassCancel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangeExecutorServer).OrderMassCancel(ctx, req.(*OrderMassCancelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ExchangeExecutor_serviceDesc = grpc.ServiceDesc{
	ServiceName: "messages.ExchangeExecutor",
	HandlerType: (*ExchangeExecutorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SecurityDefinition",
			Handler:    _ExchangeExecutor_SecurityDefinition_Handler,
		},
		{
			MethodName: "Securities",
			Handler:    _ExchangeExecutor_Securities_Handler,
		},
		{
			MethodName: "Orders",
			Handler:    _ExchangeExecutor_Orders_Handler,
		},
		{
			MethodName: "Positions",
			Handler:    _ExchangeExecutor_Positions_Handler,
		},
		{
			MethodName: "Balances",
			Handler:    _ExchangeExecutor_Balances_Handler,
		},
		{
			MethodName: "NewOrderSingle",
			Handler:    _ExchangeExecutor_NewOrderSingle_Handler,
		},
		{
			MethodName: "NewOrderBulk",
			Handler:    _ExchangeExecutor_NewOrderBulk_Handler,
		},
		{
			MethodName: "OrderReplace",
			Handler:    _ExchangeExecutor_OrderReplace_Handler,
		},
		{
			MethodName: "OrderBulkReplace",
			Handler:    _ExchangeExecutor_OrderBulkReplace_Handler,
		},
		{
			MethodName: "OrderCancel",
			Handler:    _ExchangeExecutor_OrderCancel_Handler,
		},
		{
			MethodName: "OrderMassCancel",
			Handler:    _ExchangeExecutor_OrderMassCancel_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MarketData",
			Handler:       _ExchangeExecutor_MarketData_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "AccountData",
			Handler:       _ExchangeExecutor_AccountData_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "executor_messages.proto",
}

func (m *HistoricalOpenInterestsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HistoricalOpenInterestsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HistoricalOpenInterestsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *HistoricalOpenInterestsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HistoricalOpenInterestsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HistoricalOpenInterestsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if len(m.Interests) > 0 {
		for iNdEx := len(m.Interests) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Interests[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *HistoricalUnipoolV3DataRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HistoricalUnipoolV3DataRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HistoricalUnipoolV3DataRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.End != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.End))
		i--
		dAtA[i] = 0x20
	}
	if m.Start != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Start))
		i--
		dAtA[i] = 0x18
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

func (m *HistoricalUnipoolV3DataResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HistoricalUnipoolV3DataResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HistoricalUnipoolV3DataResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x30
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if m.SeqNum != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SeqNum))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Events) > 0 {
		for iNdEx := len(m.Events) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Events[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *HistoricalProtocolAssetTransferRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HistoricalProtocolAssetTransferRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HistoricalProtocolAssetTransferRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Stop != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Stop))
		i--
		dAtA[i] = 0x20
	}
	if m.Start != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Start))
		i--
		dAtA[i] = 0x18
	}
	if m.ProtocolAsset != nil {
		{
			size, err := m.ProtocolAsset.MarshalToSizedBuffer(dAtA[:i])
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

func (m *HistoricalProtocolAssetTransferResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HistoricalProtocolAssetTransferResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HistoricalProtocolAssetTransferResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x30
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if m.SeqNum != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SeqNum))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Update) > 0 {
		for iNdEx := len(m.Update) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Update[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
		dAtA10 := make([]byte, len(m.Statistics)*10)
		var j9 int
		for _, num := range m.Statistics {
			for num >= 1<<7 {
				dAtA10[j9] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j9++
			}
			dAtA10[j9] = uint8(num)
			j9++
		}
		i -= j9
		copy(dAtA[i:], dAtA10[:j9])
		i = encodeVarintExecutorMessages(dAtA, i, uint64(j9))
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
	if len(m.Stats) > 0 {
		dAtA13 := make([]byte, len(m.Stats)*10)
		var j12 int
		for _, num := range m.Stats {
			for num >= 1<<7 {
				dAtA13[j12] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j12++
			}
			dAtA13[j12] = uint8(num)
			j12++
		}
		i -= j12
		copy(dAtA[i:], dAtA13[:j12])
		i = encodeVarintExecutorMessages(dAtA, i, uint64(j12))
		i--
		dAtA[i] = 0x32
	}
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

func (m *UnipoolV3DataRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UnipoolV3DataRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UnipoolV3DataRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *UnipoolV3DataResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UnipoolV3DataResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UnipoolV3DataResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x30
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if m.SeqNum != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SeqNum))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Update) > 0 {
		for iNdEx := len(m.Update) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Update[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *ProtocolAssetTransferRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProtocolAssetTransferRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProtocolAssetTransferRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ProtocolAsset != nil {
		{
			size, err := m.ProtocolAsset.MarshalToSizedBuffer(dAtA[:i])
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

func (m *ProtocolAssetTransferResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProtocolAssetTransferResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProtocolAssetTransferResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RejectionReason != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.RejectionReason))
		i--
		dAtA[i] = 0x30
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if m.SeqNum != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.SeqNum))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Update) > 0 {
		for iNdEx := len(m.Update) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Update[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *UnipoolV3DataIncrementalRefresh) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UnipoolV3DataIncrementalRefresh) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UnipoolV3DataIncrementalRefresh) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *ProtocolAssetDataIncrementalRefresh) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProtocolAssetDataIncrementalRefresh) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProtocolAssetDataIncrementalRefresh) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *AccountMovement) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountMovement) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AccountMovement) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Time != nil {
		{
			size, err := m.Time.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if len(m.MovementID) > 0 {
		i -= len(m.MovementID)
		copy(dAtA[i:], m.MovementID)
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.MovementID)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Subtype) > 0 {
		i -= len(m.Subtype)
		copy(dAtA[i:], m.Subtype)
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.Subtype)))
		i--
		dAtA[i] = 0x22
	}
	if m.Type != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x18
	}
	if m.Change != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Change))))
		i--
		dAtA[i] = 0x11
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
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AccountMovementRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountMovementRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AccountMovementRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
		dAtA[i] = 0x22
	}
	if m.Type != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x18
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

func (m *AccountMovementFilter) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountMovementFilter) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AccountMovementFilter) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
		dAtA[i] = 0x1a
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
		dAtA[i] = 0x12
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
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AccountMovementResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountMovementResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AccountMovementResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if len(m.Movements) > 0 {
		for iNdEx := len(m.Movements) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Movements[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if m.FromID != nil {
		{
			size, err := m.FromID.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintExecutorMessages(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x3a
	}
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
	if m.Type != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.Type))
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
		dAtA[i] = 0x28
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
	if len(m.Tag) > 0 {
		i -= len(m.Tag)
		copy(dAtA[i:], m.Tag)
		i = encodeVarintExecutorMessages(dAtA, i, uint64(len(m.Tag)))
		i--
		dAtA[i] = 0x5a
	}
	if len(m.ExecutionInstructions) > 0 {
		dAtA78 := make([]byte, len(m.ExecutionInstructions)*10)
		var j77 int
		for _, num := range m.ExecutionInstructions {
			for num >= 1<<7 {
				dAtA78[j77] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j77++
			}
			dAtA78[j77] = uint8(num)
			j77++
		}
		i -= j77
		copy(dAtA[i:], dAtA78[:j77])
		i = encodeVarintExecutorMessages(dAtA, i, uint64(j77))
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

func (m *TransferDataRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TransferDataRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TransferDataRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *TransferDataResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TransferDataResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TransferDataResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *ProtocolAssetDefinitionRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProtocolAssetDefinitionRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProtocolAssetDefinitionRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ProtocolAssetID != 0 {
		i = encodeVarintExecutorMessages(dAtA, i, uint64(m.ProtocolAssetID))
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

func (m *ProtocolAssetDefinitionResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProtocolAssetDefinitionResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProtocolAssetDefinitionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if m.ProtocolAsset != nil {
		{
			size, err := m.ProtocolAsset.MarshalToSizedBuffer(dAtA[:i])
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

func (m *ProtocolAssetListRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProtocolAssetListRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProtocolAssetListRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *ProtocolAssetListResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProtocolAssetListResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProtocolAssetListResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if len(m.ProtocolAssets) > 0 {
		for iNdEx := len(m.ProtocolAssets) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ProtocolAssets[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *ProtocolAssetList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProtocolAssetList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProtocolAssetList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if len(m.ProtocolAssets) > 0 {
		for iNdEx := len(m.ProtocolAssets) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ProtocolAssets[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
func (m *HistoricalOpenInterestsRequest) Size() (n int) {
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

func (m *HistoricalOpenInterestsResponse) Size() (n int) {
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
	if len(m.Interests) > 0 {
		for _, e := range m.Interests {
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

func (m *HistoricalUnipoolV3DataRequest) Size() (n int) {
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
	if m.Start != 0 {
		n += 1 + sovExecutorMessages(uint64(m.Start))
	}
	if m.End != 0 {
		n += 1 + sovExecutorMessages(uint64(m.End))
	}
	return n
}

func (m *HistoricalUnipoolV3DataResponse) Size() (n int) {
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
	if len(m.Events) > 0 {
		for _, e := range m.Events {
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

func (m *HistoricalProtocolAssetTransferRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
	}
	if m.ProtocolAsset != nil {
		l = m.ProtocolAsset.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Start != 0 {
		n += 1 + sovExecutorMessages(uint64(m.Start))
	}
	if m.Stop != 0 {
		n += 1 + sovExecutorMessages(uint64(m.Stop))
	}
	return n
}

func (m *HistoricalProtocolAssetTransferResponse) Size() (n int) {
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
	if len(m.Update) > 0 {
		for _, e := range m.Update {
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
	if len(m.Stats) > 0 {
		l = 0
		for _, e := range m.Stats {
			l += sovExecutorMessages(uint64(e))
		}
		n += 1 + sovExecutorMessages(uint64(l)) + l
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

func (m *UnipoolV3DataRequest) Size() (n int) {
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
	return n
}

func (m *UnipoolV3DataResponse) Size() (n int) {
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
	if len(m.Update) > 0 {
		for _, e := range m.Update {
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

func (m *ProtocolAssetTransferRequest) Size() (n int) {
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
	if m.ProtocolAsset != nil {
		l = m.ProtocolAsset.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *ProtocolAssetTransferResponse) Size() (n int) {
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
	if len(m.Update) > 0 {
		for _, e := range m.Update {
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

func (m *UnipoolV3DataIncrementalRefresh) Size() (n int) {
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
	if m.Update != nil {
		l = m.Update.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *ProtocolAssetDataIncrementalRefresh) Size() (n int) {
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
	if m.Update != nil {
		l = m.Update.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
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

func (m *AccountMovement) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Asset != nil {
		l = m.Asset.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Change != 0 {
		n += 9
	}
	if m.Type != 0 {
		n += 1 + sovExecutorMessages(uint64(m.Type))
	}
	l = len(m.Subtype)
	if l > 0 {
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	l = len(m.MovementID)
	if l > 0 {
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	if m.Time != nil {
		l = m.Time.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *AccountMovementRequest) Size() (n int) {
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
	if m.Type != 0 {
		n += 1 + sovExecutorMessages(uint64(m.Type))
	}
	if m.Filter != nil {
		l = m.Filter.Size()
		n += 1 + l + sovExecutorMessages(uint64(l))
	}
	return n
}

func (m *AccountMovementFilter) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
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

func (m *AccountMovementResponse) Size() (n int) {
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
	if len(m.Movements) > 0 {
		for _, e := range m.Movements {
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
	if m.FromID != nil {
		l = m.FromID.Size()
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
	if m.Type != 0 {
		n += 1 + sovExecutorMessages(uint64(m.Type))
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
	if m.ResponseID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.ResponseID))
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
	l = len(m.Tag)
	if l > 0 {
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

func (m *TransferDataRequest) Size() (n int) {
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

func (m *TransferDataResponse) Size() (n int) {
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
	if m.Success {
		n += 2
	}
	if m.RejectionReason != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RejectionReason))
	}
	return n
}

func (m *ProtocolAssetDefinitionRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RequestID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.RequestID))
	}
	if m.ProtocolAssetID != 0 {
		n += 1 + sovExecutorMessages(uint64(m.ProtocolAssetID))
	}
	return n
}

func (m *ProtocolAssetDefinitionResponse) Size() (n int) {
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
	if m.ProtocolAsset != nil {
		l = m.ProtocolAsset.Size()
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

func (m *ProtocolAssetListRequest) Size() (n int) {
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

func (m *ProtocolAssetListResponse) Size() (n int) {
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
	if len(m.ProtocolAssets) > 0 {
		for _, e := range m.ProtocolAssets {
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

func (m *ProtocolAssetList) Size() (n int) {
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
	if len(m.ProtocolAssets) > 0 {
		for _, e := range m.ProtocolAssets {
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

func sovExecutorMessages(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozExecutorMessages(x uint64) (n int) {
	return sovExecutorMessages(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *HistoricalOpenInterestsRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&HistoricalOpenInterestsRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`From:` + strings.Replace(fmt.Sprintf("%v", this.From), "Timestamp", "types.Timestamp", 1) + `,`,
		`To:` + strings.Replace(fmt.Sprintf("%v", this.To), "Timestamp", "types.Timestamp", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *HistoricalOpenInterestsResponse) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForInterests := "[]*Stat{"
	for _, f := range this.Interests {
		repeatedStringForInterests += strings.Replace(fmt.Sprintf("%v", f), "Stat", "models.Stat", 1) + ","
	}
	repeatedStringForInterests += "}"
	s := strings.Join([]string{`&HistoricalOpenInterestsResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Interests:` + repeatedStringForInterests + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
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
func (this *HistoricalUnipoolV3DataRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&HistoricalUnipoolV3DataRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`Start:` + fmt.Sprintf("%v", this.Start) + `,`,
		`End:` + fmt.Sprintf("%v", this.End) + `,`,
		`}`,
	}, "")
	return s
}
func (this *HistoricalUnipoolV3DataResponse) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForEvents := "[]*UPV3Update{"
	for _, f := range this.Events {
		repeatedStringForEvents += strings.Replace(fmt.Sprintf("%v", f), "UPV3Update", "models.UPV3Update", 1) + ","
	}
	repeatedStringForEvents += "}"
	s := strings.Join([]string{`&HistoricalUnipoolV3DataResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Events:` + repeatedStringForEvents + `,`,
		`SeqNum:` + fmt.Sprintf("%v", this.SeqNum) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *HistoricalProtocolAssetTransferRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&HistoricalProtocolAssetTransferRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ProtocolAsset:` + strings.Replace(fmt.Sprintf("%v", this.ProtocolAsset), "ProtocolAsset", "models.ProtocolAsset", 1) + `,`,
		`Start:` + fmt.Sprintf("%v", this.Start) + `,`,
		`Stop:` + fmt.Sprintf("%v", this.Stop) + `,`,
		`}`,
	}, "")
	return s
}
func (this *HistoricalProtocolAssetTransferResponse) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForUpdate := "[]*ProtocolAssetUpdate{"
	for _, f := range this.Update {
		repeatedStringForUpdate += strings.Replace(fmt.Sprintf("%v", f), "ProtocolAssetUpdate", "models.ProtocolAssetUpdate", 1) + ","
	}
	repeatedStringForUpdate += "}"
	s := strings.Join([]string{`&HistoricalProtocolAssetTransferResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Update:` + repeatedStringForUpdate + `,`,
		`SeqNum:` + fmt.Sprintf("%v", this.SeqNum) + `,`,
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
		`Stats:` + fmt.Sprintf("%v", this.Stats) + `,`,
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
func (this *UnipoolV3DataRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&UnipoolV3DataRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Subscribe:` + fmt.Sprintf("%v", this.Subscribe) + `,`,
		`Subscriber:` + strings.Replace(fmt.Sprintf("%v", this.Subscriber), "PID", "actor.PID", 1) + `,`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *UnipoolV3DataResponse) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForUpdate := "[]*UPV3Update{"
	for _, f := range this.Update {
		repeatedStringForUpdate += strings.Replace(fmt.Sprintf("%v", f), "UPV3Update", "models.UPV3Update", 1) + ","
	}
	repeatedStringForUpdate += "}"
	s := strings.Join([]string{`&UnipoolV3DataResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Update:` + repeatedStringForUpdate + `,`,
		`SeqNum:` + fmt.Sprintf("%v", this.SeqNum) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ProtocolAssetTransferRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ProtocolAssetTransferRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Subscribe:` + fmt.Sprintf("%v", this.Subscribe) + `,`,
		`Subscriber:` + strings.Replace(fmt.Sprintf("%v", this.Subscriber), "PID", "actor.PID", 1) + `,`,
		`ProtocolAsset:` + strings.Replace(fmt.Sprintf("%v", this.ProtocolAsset), "ProtocolAsset", "models.ProtocolAsset", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ProtocolAssetTransferResponse) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForUpdate := "[]*ProtocolAssetUpdate{"
	for _, f := range this.Update {
		repeatedStringForUpdate += strings.Replace(fmt.Sprintf("%v", f), "ProtocolAssetUpdate", "models.ProtocolAssetUpdate", 1) + ","
	}
	repeatedStringForUpdate += "}"
	s := strings.Join([]string{`&ProtocolAssetTransferResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Update:` + repeatedStringForUpdate + `,`,
		`SeqNum:` + fmt.Sprintf("%v", this.SeqNum) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *UnipoolV3DataIncrementalRefresh) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&UnipoolV3DataIncrementalRefresh{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`SeqNum:` + fmt.Sprintf("%v", this.SeqNum) + `,`,
		`Update:` + strings.Replace(fmt.Sprintf("%v", this.Update), "UPV3Update", "models.UPV3Update", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ProtocolAssetDataIncrementalRefresh) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ProtocolAssetDataIncrementalRefresh{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`SeqNum:` + fmt.Sprintf("%v", this.SeqNum) + `,`,
		`Update:` + strings.Replace(fmt.Sprintf("%v", this.Update), "ProtocolAssetUpdate", "models.ProtocolAssetUpdate", 1) + `,`,
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
func (this *AccountMovement) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AccountMovement{`,
		`Asset:` + strings.Replace(fmt.Sprintf("%v", this.Asset), "Asset", "models1.Asset", 1) + `,`,
		`Change:` + fmt.Sprintf("%v", this.Change) + `,`,
		`Type:` + fmt.Sprintf("%v", this.Type) + `,`,
		`Subtype:` + fmt.Sprintf("%v", this.Subtype) + `,`,
		`MovementID:` + fmt.Sprintf("%v", this.MovementID) + `,`,
		`Time:` + strings.Replace(fmt.Sprintf("%v", this.Time), "Timestamp", "types.Timestamp", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AccountMovementRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AccountMovementRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Account:` + strings.Replace(fmt.Sprintf("%v", this.Account), "Account", "models.Account", 1) + `,`,
		`Type:` + fmt.Sprintf("%v", this.Type) + `,`,
		`Filter:` + strings.Replace(this.Filter.String(), "AccountMovementFilter", "AccountMovementFilter", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AccountMovementFilter) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AccountMovementFilter{`,
		`Instrument:` + strings.Replace(fmt.Sprintf("%v", this.Instrument), "Instrument", "models.Instrument", 1) + `,`,
		`From:` + strings.Replace(fmt.Sprintf("%v", this.From), "Timestamp", "types.Timestamp", 1) + `,`,
		`To:` + strings.Replace(fmt.Sprintf("%v", this.To), "Timestamp", "types.Timestamp", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AccountMovementResponse) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForMovements := "[]*AccountMovement{"
	for _, f := range this.Movements {
		repeatedStringForMovements += strings.Replace(f.String(), "AccountMovement", "AccountMovement", 1) + ","
	}
	repeatedStringForMovements += "}"
	s := strings.Join([]string{`&AccountMovementResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`Movements:` + repeatedStringForMovements + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
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
		`FromID:` + strings.Replace(fmt.Sprintf("%v", this.FromID), "StringValue", "types.StringValue", 1) + `,`,
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
		`Type:` + fmt.Sprintf("%v", this.Type) + `,`,
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
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
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
		`Tag:` + fmt.Sprintf("%v", this.Tag) + `,`,
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
func (this *TransferDataRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&TransferDataRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Subscribe:` + fmt.Sprintf("%v", this.Subscribe) + `,`,
		`Subscriber:` + strings.Replace(fmt.Sprintf("%v", this.Subscriber), "PID", "actor.PID", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *TransferDataResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&TransferDataResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`SeqNum:` + fmt.Sprintf("%v", this.SeqNum) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ProtocolAssetDefinitionRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ProtocolAssetDefinitionRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ProtocolAssetID:` + fmt.Sprintf("%v", this.ProtocolAssetID) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ProtocolAssetDefinitionResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ProtocolAssetDefinitionResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`ProtocolAsset:` + strings.Replace(fmt.Sprintf("%v", this.ProtocolAsset), "ProtocolAsset", "models.ProtocolAsset", 1) + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ProtocolAssetListRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ProtocolAssetListRequest{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`Subscribe:` + fmt.Sprintf("%v", this.Subscribe) + `,`,
		`Subscriber:` + strings.Replace(fmt.Sprintf("%v", this.Subscriber), "PID", "actor.PID", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ProtocolAssetListResponse) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForProtocolAssets := "[]*ProtocolAsset{"
	for _, f := range this.ProtocolAssets {
		repeatedStringForProtocolAssets += strings.Replace(fmt.Sprintf("%v", f), "ProtocolAsset", "models.ProtocolAsset", 1) + ","
	}
	repeatedStringForProtocolAssets += "}"
	s := strings.Join([]string{`&ProtocolAssetListResponse{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`ProtocolAssets:` + repeatedStringForProtocolAssets + `,`,
		`Success:` + fmt.Sprintf("%v", this.Success) + `,`,
		`RejectionReason:` + fmt.Sprintf("%v", this.RejectionReason) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ProtocolAssetList) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForProtocolAssets := "[]*ProtocolAsset{"
	for _, f := range this.ProtocolAssets {
		repeatedStringForProtocolAssets += strings.Replace(fmt.Sprintf("%v", f), "ProtocolAsset", "models.ProtocolAsset", 1) + ","
	}
	repeatedStringForProtocolAssets += "}"
	s := strings.Join([]string{`&ProtocolAssetList{`,
		`RequestID:` + fmt.Sprintf("%v", this.RequestID) + `,`,
		`ResponseID:` + fmt.Sprintf("%v", this.ResponseID) + `,`,
		`ProtocolAssets:` + repeatedStringForProtocolAssets + `,`,
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
func (m *HistoricalOpenInterestsRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: HistoricalOpenInterestsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HistoricalOpenInterestsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *HistoricalOpenInterestsResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: HistoricalOpenInterestsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HistoricalOpenInterestsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Interests", wireType)
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
			m.Interests = append(m.Interests, &models.Stat{})
			if err := m.Interests[len(m.Interests)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *HistoricalUnipoolV3DataRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: HistoricalUnipoolV3DataRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HistoricalUnipoolV3DataRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Start", wireType)
			}
			m.Start = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Start |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field End", wireType)
			}
			m.End = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.End |= uint64(b&0x7F) << shift
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
func (m *HistoricalUnipoolV3DataResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: HistoricalUnipoolV3DataResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HistoricalUnipoolV3DataResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Events", wireType)
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
			m.Events = append(m.Events, &models.UPV3Update{})
			if err := m.Events[len(m.Events)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
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
		case 5:
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
		case 6:
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
func (m *HistoricalProtocolAssetTransferRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: HistoricalProtocolAssetTransferRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HistoricalProtocolAssetTransferRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field ProtocolAsset", wireType)
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
			if m.ProtocolAsset == nil {
				m.ProtocolAsset = &models.ProtocolAsset{}
			}
			if err := m.ProtocolAsset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Start", wireType)
			}
			m.Start = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Start |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stop", wireType)
			}
			m.Stop = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Stop |= uint64(b&0x7F) << shift
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
func (m *HistoricalProtocolAssetTransferResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: HistoricalProtocolAssetTransferResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HistoricalProtocolAssetTransferResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
			m.Update = append(m.Update, &models.ProtocolAssetUpdate{})
			if err := m.Update[len(m.Update)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
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
		case 5:
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
		case 6:
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
		case 6:
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
				m.Stats = append(m.Stats, v)
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
				if elementCount != 0 && len(m.Stats) == 0 {
					m.Stats = make([]models.StatType, 0, elementCount)
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
					m.Stats = append(m.Stats, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Stats", wireType)
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
func (m *UnipoolV3DataRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: UnipoolV3DataRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UnipoolV3DataRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *UnipoolV3DataResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: UnipoolV3DataResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UnipoolV3DataResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
			m.Update = append(m.Update, &models.UPV3Update{})
			if err := m.Update[len(m.Update)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
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
		case 5:
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
		case 6:
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
func (m *ProtocolAssetTransferRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: ProtocolAssetTransferRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProtocolAssetTransferRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field ProtocolAsset", wireType)
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
			if m.ProtocolAsset == nil {
				m.ProtocolAsset = &models.ProtocolAsset{}
			}
			if err := m.ProtocolAsset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *ProtocolAssetTransferResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: ProtocolAssetTransferResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProtocolAssetTransferResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
			m.Update = append(m.Update, &models.ProtocolAssetUpdate{})
			if err := m.Update[len(m.Update)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
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
		case 5:
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
		case 6:
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
func (m *UnipoolV3DataIncrementalRefresh) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: UnipoolV3DataIncrementalRefresh: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UnipoolV3DataIncrementalRefresh: illegal tag %d (wire type %d)", fieldNum, wire)
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
				m.Update = &models.UPV3Update{}
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
func (m *ProtocolAssetDataIncrementalRefresh) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: ProtocolAssetDataIncrementalRefresh: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProtocolAssetDataIncrementalRefresh: illegal tag %d (wire type %d)", fieldNum, wire)
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
				m.Update = &models.ProtocolAssetUpdate{}
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
func (m *AccountMovement) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: AccountMovement: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccountMovement: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
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
		case 2:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Change", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Change = float64(math.Float64frombits(v))
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= AccountMovementType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subtype", wireType)
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
			m.Subtype = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MovementID", wireType)
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
			m.MovementID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
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
			if m.Time == nil {
				m.Time = &types.Timestamp{}
			}
			if err := m.Time.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *AccountMovementRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: AccountMovementRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccountMovementRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= AccountMovementType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
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
				m.Filter = &AccountMovementFilter{}
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
func (m *AccountMovementFilter) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: AccountMovementFilter: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccountMovementFilter: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
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
		case 2:
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
		case 3:
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
func (m *AccountMovementResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: AccountMovementResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccountMovementResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Movements", wireType)
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
			m.Movements = append(m.Movements, &AccountMovement{})
			if err := m.Movements[len(m.Movements)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FromID", wireType)
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
			if m.FromID == nil {
				m.FromID = &types.StringValue{}
			}
			if err := m.FromID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= AccountMovementType(b&0x7F) << shift
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
				var v models.ExecutionInstruction
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowExecutorMessages
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= models.ExecutionInstruction(b&0x7F) << shift
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
					m.ExecutionInstructions = make([]models.ExecutionInstruction, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v models.ExecutionInstruction
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowExecutorMessages
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= models.ExecutionInstruction(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.ExecutionInstructions = append(m.ExecutionInstructions, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field ExecutionInstructions", wireType)
			}
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tag", wireType)
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
			m.Tag = string(dAtA[iNdEx:postIndex])
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
func (m *TransferDataRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: TransferDataRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TransferDataRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *TransferDataResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: TransferDataResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TransferDataResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *ProtocolAssetDefinitionRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: ProtocolAssetDefinitionRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProtocolAssetDefinitionRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field ProtocolAssetID", wireType)
			}
			m.ProtocolAssetID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExecutorMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ProtocolAssetID |= uint64(b&0x7F) << shift
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
func (m *ProtocolAssetDefinitionResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: ProtocolAssetDefinitionResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProtocolAssetDefinitionResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field ProtocolAsset", wireType)
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
			if m.ProtocolAsset == nil {
				m.ProtocolAsset = &models.ProtocolAsset{}
			}
			if err := m.ProtocolAsset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *ProtocolAssetListRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: ProtocolAssetListRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProtocolAssetListRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *ProtocolAssetListResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: ProtocolAssetListResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProtocolAssetListResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field ProtocolAssets", wireType)
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
			m.ProtocolAssets = append(m.ProtocolAssets, &models.ProtocolAsset{})
			if err := m.ProtocolAssets[len(m.ProtocolAssets)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *ProtocolAssetList) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: ProtocolAssetList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProtocolAssetList: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field ProtocolAssets", wireType)
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
			m.ProtocolAssets = append(m.ProtocolAssets, &models.ProtocolAsset{})
			if err := m.ProtocolAssets[len(m.ProtocolAssets)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
