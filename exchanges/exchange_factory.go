package exchanges

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alphac/exchanges/binance"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/xchanger/constants"
	models2 "gitlab.com/alphaticks/xchanger/models"
)

func NewInstrumentListenerProducer(security *models.Security) actor.Producer {
	switch security.Exchange.ID {
	case constants.BINANCE.ID:
		return func() actor.Actor { return binance.NewListener(security) }
		/*
			case constants.BITFINEX:
				return func() actor.Actor { return bitfinex.NewListener(instrument) }
			case constants.BITMEX:
				return func() actor.Actor { return bitmex.NewListener(instrument) }
			case constants.BITSTAMP:
				return func() actor.Actor { return bitstamp.NewListener(instrument) }
			case constants.BITTREX:
				return func() actor.Actor { return bittrex.NewListener(instrument) }
			case constants.COINBASEPRO:
				return func() actor.Actor { return coinbasepro.NewListener(instrument) }
			case constants.CRYPTOFACILITIES:
				return func() actor.Actor { return cryptofacilities.NewListener(instrument) }
			case constants.FBINANCE:
				return func() actor.Actor { return fbinance.NewListener(instrument) }
			case constants.GEMINI:
				return func() actor.Actor { return gemini.NewListener(instrument) }
			case constants.HITBTC:
				return func() actor.Actor { return hitbtc.NewListener(instrument) }
			case constants.KRAKEN:
				return func() actor.Actor { return kraken.NewListener(instrument) }
			case constants.OKCOIN:
				return func() actor.Actor { return okcoin.NewListener(instrument) }
			case constants.BITZ:
				return func() actor.Actor { return bitz.NewListener(instrument) }
			case constants.HUOBI:
				return func() actor.Actor { return huobi.NewListener(instrument) }
			case constants.FTX:
				return func() actor.Actor { return ftx.NewListener(instrument) }

		*/
	default:
		return nil
	}
}

func NewExchangeExecutorProducer(exchange *models2.Exchange) actor.Producer {
	switch exchange.ID {
	case constants.BINANCE.ID:
		return func() actor.Actor { return binance.NewExecutor() }
		/*
			case constants.BITFINEX:
				return func() actor.Actor { return bitfinex.NewExecutor() }
			case constants.BITMEX:
				return func() actor.Actor { return bitmex.NewExecutor() }
			case constants.BITSTAMP:
				return func() actor.Actor { return bitstamp.NewExecutor() }
			case constants.COINBASEPRO:
				return func() actor.Actor { return coinbasepro.NewExecutor() }
			case constants.KRAKEN:
				return func() actor.Actor { return kraken.NewExecutor() }
			case constants.CRYPTOFACILITIES:
				return func() actor.Actor { return cryptofacilities.NewExecutor() }
			case constants.FBINANCE:
				return func() actor.Actor { return fbinance.NewExecutor() }
			case constants.OKCOIN:
				return func() actor.Actor { return okcoin.NewExecutor() }
			case constants.GEMINI:
				return func() actor.Actor { return gemini.NewExecutor() }
			case constants.HITBTC:
				return func() actor.Actor { return hitbtc.NewExecutor() }
			case constants.BITTREX:
				return func() actor.Actor { return bittrex.NewExecutor() }
			case constants.BITZ:
				return func() actor.Actor { return bitz.NewExecutor() }
			case constants.HUOBI:
				return func() actor.Actor { return huobi.NewExecutor() }
			case constants.FTX:
				return func() actor.Actor { return ftx.NewExecutor() }

		*/
	default:
		return nil
	}
}
