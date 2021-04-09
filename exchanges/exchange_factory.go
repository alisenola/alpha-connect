package exchanges

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/account"
	"gitlab.com/alphaticks/alpha-connect/exchanges/binance"
	"gitlab.com/alphaticks/alpha-connect/exchanges/bitfinex"
	"gitlab.com/alphaticks/alpha-connect/exchanges/bithumb"
	"gitlab.com/alphaticks/alpha-connect/exchanges/bithumbg"
	"gitlab.com/alphaticks/alpha-connect/exchanges/bitmex"
	"gitlab.com/alphaticks/alpha-connect/exchanges/bitstamp"
	"gitlab.com/alphaticks/alpha-connect/exchanges/bitz"
	"gitlab.com/alphaticks/alpha-connect/exchanges/bybiti"
	"gitlab.com/alphaticks/alpha-connect/exchanges/bybitl"
	"gitlab.com/alphaticks/alpha-connect/exchanges/coinbasepro"
	"gitlab.com/alphaticks/alpha-connect/exchanges/cryptofacilities"
	"gitlab.com/alphaticks/alpha-connect/exchanges/deribit"
	"gitlab.com/alphaticks/alpha-connect/exchanges/fbinance"
	"gitlab.com/alphaticks/alpha-connect/exchanges/ftx"
	"gitlab.com/alphaticks/alpha-connect/exchanges/gemini"
	"gitlab.com/alphaticks/alpha-connect/exchanges/hitbtc"
	"gitlab.com/alphaticks/alpha-connect/exchanges/huobi"
	"gitlab.com/alphaticks/alpha-connect/exchanges/huobif"
	"gitlab.com/alphaticks/alpha-connect/exchanges/huobip"
	"gitlab.com/alphaticks/alpha-connect/exchanges/kraken"
	"gitlab.com/alphaticks/alpha-connect/exchanges/okcoin"
	"gitlab.com/alphaticks/alpha-connect/exchanges/okex"
	"gitlab.com/alphaticks/alpha-connect/exchanges/upbit"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/xchanger/constants"
	models2 "gitlab.com/alphaticks/xchanger/models"
	"gitlab.com/alphaticks/xchanger/utils"
)

// TODO sample size config
var Portfolio = account.NewPortfolio(300)

func NewAccount(accountInfo *models.Account) (*account.Account, error) {
	if accnt := Portfolio.GetAccount(accountInfo.AccountID); accnt != nil {
		return accnt, nil
	}
	switch accountInfo.Exchange.ID {
	case constants.BITMEX.ID:
		accnt, err := account.NewAccount(accountInfo, &constants.BITCOIN, 1./0.00000001)
		if err != nil {
			return nil, err
		}
		Portfolio.AddAccount(accnt)
		return accnt, nil
	case constants.FBINANCE.ID:
		accnt, err := account.NewAccount(accountInfo, &constants.TETHER, 100000000)
		if err != nil {
			return nil, err
		}
		Portfolio.AddAccount(accnt)
		return accnt, nil
	default:
		return nil, fmt.Errorf("unsupported exchange %s", accountInfo.Exchange.Name)
	}
}

func NewAccountListenerProducer(account *account.Account) actor.Producer {
	switch account.Exchange.ID {
	case constants.BITMEX.ID:
		return func() actor.Actor { return bitmex.NewAccountListener(account) }
	case constants.FBINANCE.ID:
		return func() actor.Actor { return fbinance.NewAccountListener(account) }
	default:
		return nil
	}
}

func NewPaperAccountListenerProducer(account *account.Account) actor.Producer {
	switch account.Exchange.ID {
	case constants.BITMEX.ID:
		return func() actor.Actor { return bitmex.NewPaperAccountListener(account) }
	default:
		return nil
	}
}

func NewInstrumentListenerProducer(security *models.Security, dialerPool *utils.DialerPool) actor.Producer {
	switch security.Exchange.ID {
	case constants.BINANCE.ID:
		return func() actor.Actor { return binance.NewListener(security, dialerPool) }
	case constants.BITFINEX.ID:
		return func() actor.Actor { return bitfinex.NewListener(security, dialerPool) }
	case constants.BITMEX.ID:
		return func() actor.Actor { return bitmex.NewListener(security, dialerPool) }
	case constants.BITSTAMP.ID:
		return func() actor.Actor { return bitstamp.NewListenerL3(security, dialerPool) }
	case constants.BITZ.ID:
		return func() actor.Actor { return bitz.NewListener(security, dialerPool) }
	case constants.COINBASEPRO.ID:
		return func() actor.Actor { return coinbasepro.NewListener(security, dialerPool) }
	case constants.CRYPTOFACILITIES.ID:
		return func() actor.Actor { return cryptofacilities.NewListener(security, dialerPool) }
	case constants.FBINANCE.ID:
		return func() actor.Actor { return fbinance.NewListener(security, dialerPool) }
	case constants.FTX.ID:
		return func() actor.Actor { return ftx.NewListener(security, dialerPool) }
	case constants.GEMINI.ID:
		return func() actor.Actor { return gemini.NewListener(security, dialerPool) }
	case constants.HITBTC.ID:
		return func() actor.Actor { return hitbtc.NewListener(security, dialerPool) }
	case constants.KRAKEN.ID:
		return func() actor.Actor { return kraken.NewListener(security, dialerPool) }
	case constants.OKCOIN.ID:
		return func() actor.Actor { return okcoin.NewListener(security, dialerPool) }
	case constants.OKEX.ID:
		return func() actor.Actor { return okex.NewListener(security, dialerPool) }
	case constants.DERIBIT.ID:
		return func() actor.Actor { return deribit.NewListener(security, dialerPool) }
	case constants.HUOBI.ID:
		return func() actor.Actor { return huobi.NewListener(security, dialerPool) }
	case constants.HUOBIP.ID:
		return func() actor.Actor { return huobip.NewListener(security, dialerPool) }
	case constants.HUOBIF.ID:
		return func() actor.Actor { return huobif.NewListener(security, dialerPool) }
	case constants.BYBITI.ID:
		return func() actor.Actor { return bybiti.NewListener(security, dialerPool) }
	case constants.BYBITL.ID:
		return func() actor.Actor { return bybitl.NewListener(security, dialerPool) }
	case constants.UPBIT.ID:
		return func() actor.Actor { return upbit.NewListener(security, dialerPool) }
	case constants.BITHUMB.ID:
		return func() actor.Actor { return bithumb.NewListener(security, dialerPool) }
	case constants.BITHUMBG.ID:
		return func() actor.Actor { return bithumbg.NewListener(security, dialerPool) }
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
	case constants.GEMINI.ID:
		return func() actor.Actor { return gemini.NewExecutor() }
	case constants.HITBTC.ID:
		return func() actor.Actor { return hitbtc.NewExecutor() }
	case constants.KRAKEN.ID:
		return func() actor.Actor { return kraken.NewExecutor() }
	case constants.OKCOIN.ID:
		return func() actor.Actor { return okcoin.NewExecutor() }
	case constants.OKEX.ID:
		return func() actor.Actor { return okex.NewExecutor() }
	case constants.DERIBIT.ID:
		return func() actor.Actor { return deribit.NewExecutor() }
	case constants.HUOBI.ID:
		return func() actor.Actor { return huobi.NewExecutor() }
	case constants.HUOBIP.ID:
		return func() actor.Actor { return huobip.NewExecutor() }
	case constants.HUOBIF.ID:
		return func() actor.Actor { return huobif.NewExecutor() }
	case constants.BYBITI.ID:
		return func() actor.Actor { return bybiti.NewExecutor() }
	case constants.BYBITL.ID:
		return func() actor.Actor { return bybitl.NewExecutor() }
	case constants.UPBIT.ID:
		return func() actor.Actor { return upbit.NewExecutor() }
	case constants.BITHUMB.ID:
		return func() actor.Actor { return bithumb.NewExecutor() }
	case constants.BITHUMBG.ID:
		return func() actor.Actor { return bithumbg.NewExecutor() }
		/*
			case constants.BITTREX:
				return func() actor.Actor { return bittrex.NewExecutor() }
		*/
	default:
		return nil
	}
}
