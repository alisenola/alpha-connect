package exchanges

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alphac/exchanges/binance"
	"gitlab.com/alphaticks/alphac/exchanges/bitfinex"
	"gitlab.com/alphaticks/alphac/exchanges/bitmex"
	"gitlab.com/alphaticks/alphac/exchanges/bitstamp"
	"gitlab.com/alphaticks/alphac/exchanges/bitz"
	"gitlab.com/alphaticks/alphac/exchanges/coinbasepro"
	"gitlab.com/alphaticks/alphac/exchanges/cryptofacilities"
	"gitlab.com/alphaticks/alphac/exchanges/fbinance"
	"gitlab.com/alphaticks/alphac/exchanges/ftx"
	"gitlab.com/alphaticks/alphac/exchanges/gemini"
	"gitlab.com/alphaticks/alphac/exchanges/hitbtc"
	"gitlab.com/alphaticks/alphac/exchanges/huobi"
	"gitlab.com/alphaticks/alphac/exchanges/kraken"
	"gitlab.com/alphaticks/alphac/exchanges/okcoin"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/xchanger/constants"
	models2 "gitlab.com/alphaticks/xchanger/models"
)

func NewInstrumentListenerProducer(security *models.Security) actor.Producer {
	switch security.Exchange.ID {
	case constants.BINANCE.ID:
		return func() actor.Actor { return binance.NewListener(security) }
	case constants.BITFINEX.ID:
		return func() actor.Actor { return bitfinex.NewListener(security) }
	case constants.BITMEX.ID:
		return func() actor.Actor { return bitmex.NewListener(security) }
	case constants.BITSTAMP.ID:
		return func() actor.Actor { return bitstamp.NewListener(security) }
	case constants.BITZ.ID:
		return func() actor.Actor { return bitz.NewListener(security) }
	case constants.COINBASEPRO.ID:
		return func() actor.Actor { return coinbasepro.NewListener(security) }
	case constants.CRYPTOFACILITIES.ID:
		return func() actor.Actor { return cryptofacilities.NewListener(security) }
	case constants.FBINANCE.ID:
		return func() actor.Actor { return fbinance.NewListener(security) }
	case constants.FTX.ID:
		return func() actor.Actor { return ftx.NewListener(security) }
	case constants.HUOBI.ID:
		return func() actor.Actor { return huobi.NewListener(security) }
	case constants.GEMINI.ID:
		return func() actor.Actor { return gemini.NewListener(security) }
	case constants.HITBTC.ID:
		return func() actor.Actor { return hitbtc.NewListener(security) }
	case constants.KRAKEN.ID:
		return func() actor.Actor { return kraken.NewListener(security) }
	case constants.OKCOIN.ID:
		return func() actor.Actor { return okcoin.NewListener(security) }
		/*
			case constants.BITTREX:
			return func() actor.Actor { return bittrex.NewListener(instrument) }



		*/
	default:
		return nil
	}
}

func NewExchangeExecutorProducer(exchange *models2.Exchange) actor.Producer {
	switch exchange.ID {
	case constants.BINANCE.ID:
		return func() actor.Actor { return binance.NewExecutor() }
	case constants.BITFINEX.ID:
		return func() actor.Actor { return bitfinex.NewExecutor() }
	case constants.BITMEX.ID:
		return func() actor.Actor { return bitmex.NewExecutor() }
	case constants.BITSTAMP.ID:
		return func() actor.Actor { return bitstamp.NewExecutor() }
	case constants.BITZ.ID:
		return func() actor.Actor { return bitz.NewExecutor() }
	case constants.COINBASEPRO.ID:
		return func() actor.Actor { return coinbasepro.NewExecutor() }
	case constants.CRYPTOFACILITIES.ID:
		return func() actor.Actor { return cryptofacilities.NewExecutor() }
	case constants.FBINANCE.ID:
		return func() actor.Actor { return fbinance.NewExecutor() }
	case constants.FTX.ID:
		return func() actor.Actor { return ftx.NewExecutor() }
	case constants.HUOBI.ID:
		return func() actor.Actor { return huobi.NewExecutor() }
	case constants.GEMINI.ID:
		return func() actor.Actor { return gemini.NewExecutor() }
	case constants.HITBTC.ID:
		return func() actor.Actor { return hitbtc.NewExecutor() }
	case constants.KRAKEN.ID:
		return func() actor.Actor { return kraken.NewExecutor() }
	case constants.OKCOIN.ID:
		return func() actor.Actor { return okcoin.NewExecutor() }

		/*

			case constants.BITTREX:
				return func() actor.Actor { return bittrex.NewExecutor() }



		*/
	default:
		return nil
	}
}
