package executor

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	exchanges2 "gitlab.com/alphaticks/alphac/messages/exchanges"
	"gitlab.com/alphaticks/xchanger/exchanges"
)

type SetExchangeExecutorCredentials struct {
}

type GetExchangeExecutorRequest struct {
	RequestID int64
	Exchange  string
}

type GetExchangeExecutorResponse struct {
	RequestID int64
	PID       *actor.PID
}

type RestartExchangeExecutorRequest struct {
	RequestID int64
	Exchange  string
}

type RestartExchangeExecutorResponse struct {
	RequestID int64
	PID       *actor.PID
}

type GetInstrumentsRequest struct {
	RequestID int64
	Exchange  string
}

type GetInstrumentsResponse struct {
	RequestID   int64
	Error       error
	Instruments []*exchanges.Instrument
}

type GetOpenOrdersRequest struct {
	RequestID  int64
	Instrument *exchanges.Instrument
}

type OpenOrdersRequest struct {
	RequestID  int64
	Instrument *exchanges.Instrument
}

type CloseOrdersRequest struct {
	RequestID  int64
	Instrument *exchanges.Instrument
}

type CloseAllOrdersRequest struct {
	RequestID int64
}

type CloseAllOrdersResponse struct {
	RequestID int64
	Error     error
}

type GetOrderBookL2Request struct {
	RequestID  int64
	Instrument *exchanges.Instrument
}

type GetOrderBookL2Response struct {
	RequestID int64
	Error     error
	Snapshot  *exchanges2.OBL2Snapshot
}

type GetOrderBookL3Request struct {
	RequestID  int64
	Instrument *exchanges.Instrument
}

type GetOrderBookL3Response struct {
	RequestID int64
	Error     error
	Snapshot  *exchanges2.OBL3Snapshot
}

type GetOrdersRequest struct {
	Instrument *exchanges.Instrument
}

type GetOpenPositionsRequest struct {
	Request  int64
	Error    error
	Exchange string
}

type GetOpenPositionsResponse struct {
}
