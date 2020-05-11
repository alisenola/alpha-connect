package listener

import (
	"gitlab.com/alphaticks/xchanger/exchanges"
)

type GetInstrumentsRequest struct {
	RequestID int64
	Exchange  string
}

type GetInstrumentsResponse struct {
	RequestID   int64
	Instruments []*exchanges.Instrument
}

type SetupInstrumentsRequest struct {
	RequestID   int64
	Instruments []*exchanges.Instrument
}

type SetupInstrumentsResponse struct {
	RequestID int64
}

type RestartExchangeListenerRequest struct {
	RequestID int64
	Exchange  string
}
