package exchanges

import (
	"github.com/asynkron/protoactor-go/actor"
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
	"gitlab.com/alphaticks/alpha-connect/exchanges/bybits"
	"gitlab.com/alphaticks/alpha-connect/exchanges/coinbasepro"
	"gitlab.com/alphaticks/alpha-connect/exchanges/cryptofacilities"
	"gitlab.com/alphaticks/alpha-connect/exchanges/deribit"
	"gitlab.com/alphaticks/alpha-connect/exchanges/dydx"
	"gitlab.com/alphaticks/alpha-connect/exchanges/fbinance"
	"gitlab.com/alphaticks/alpha-connect/exchanges/ftx"
	"gitlab.com/alphaticks/alpha-connect/exchanges/ftxus"
	"gitlab.com/alphaticks/alpha-connect/exchanges/gate"
	"gitlab.com/alphaticks/alpha-connect/exchanges/gemini"
	"gitlab.com/alphaticks/alpha-connect/exchanges/hitbtc"
	"gitlab.com/alphaticks/alpha-connect/exchanges/huobi"
	"gitlab.com/alphaticks/alpha-connect/exchanges/huobif"
	"gitlab.com/alphaticks/alpha-connect/exchanges/huobil"
	"gitlab.com/alphaticks/alpha-connect/exchanges/huobip"
	"gitlab.com/alphaticks/alpha-connect/exchanges/kraken"
	"gitlab.com/alphaticks/alpha-connect/exchanges/okcoin"
	"gitlab.com/alphaticks/alpha-connect/exchanges/okex"
	"gitlab.com/alphaticks/alpha-connect/exchanges/okexp"
	"gitlab.com/alphaticks/alpha-connect/exchanges/opensea"
	"gitlab.com/alphaticks/alpha-connect/exchanges/types"
	v3 "gitlab.com/alphaticks/alpha-connect/exchanges/uniswap/v3"
	"gitlab.com/alphaticks/alpha-connect/exchanges/upbit"
	"gitlab.com/alphaticks/alpha-connect/models"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger/constants"
	models2 "gitlab.com/alphaticks/xchanger/models"
	"gitlab.com/alphaticks/xchanger/utils"
	"gorm.io/gorm"
)

// TODO sample size config
var Portfolio = account.NewPortfolio(300)

func NewAccount(accountInfo *models.Account) (*account.Account, error) {
	if accnt := Portfolio.GetAccount(accountInfo.Name); accnt != nil {
		return accnt, nil
	}
	accnt, err := account.NewAccount(accountInfo)
	if err != nil {
		return nil, err
	}
	Portfolio.AddAccount(accnt)
	return accnt, nil
}

func NewAccountListenerProducer(account *account.Account, registry registry.PublicRegistryClient, db *gorm.DB, readOnly bool) actor.Producer {
	switch account.Exchange.ID {
	case constants.BITMEX.ID:
		return func() actor.Actor { return bitmex.NewAccountListener(account) }
	case constants.BINANCE.ID:
		return func() actor.Actor { return binance.NewAccountListener(account, nil, nil) }
	case constants.FBINANCE.ID:
		return func() actor.Actor { return fbinance.NewAccountListener(account, registry, db, readOnly) }
	case constants.FTX.ID:
		return func() actor.Actor { return ftx.NewAccountListener(account, registry, db, readOnly) }
	case constants.FTXUS.ID:
		return func() actor.Actor { return ftxus.NewAccountListener(account) }
	case constants.DYDX.ID:
		return func() actor.Actor { return dydx.NewAccountListener(account, nil, nil) }
	case constants.BYBITL.ID:
		return func() actor.Actor { return bybitl.NewAccountListener(account, registry, db, readOnly) }
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

func NewInstrumentListenerProducer(securityID uint64, exchangeID uint32, dialerPool *utils.DialerPool) actor.Producer {
	switch exchangeID {
	case constants.BINANCE.ID:
		return func() actor.Actor { return binance.NewListener(securityID, dialerPool) }
	case constants.BITFINEX.ID:
		return func() actor.Actor { return bitfinex.NewListener(securityID, dialerPool) }
	case constants.BITMEX.ID:
		return func() actor.Actor { return bitmex.NewListener(securityID, dialerPool) }
	case constants.BITSTAMP.ID:
		return func() actor.Actor { return bitstamp.NewListenerL3(securityID, dialerPool) }
	case constants.BITZ.ID:
		return func() actor.Actor { return bitz.NewListener(securityID, dialerPool) }
	case constants.COINBASEPRO.ID:
		return func() actor.Actor { return coinbasepro.NewListener(securityID, dialerPool) }
	case constants.CRYPTOFACILITIES.ID:
		return func() actor.Actor { return cryptofacilities.NewListener(securityID, dialerPool) }
	case constants.FBINANCE.ID:
		return func() actor.Actor { return fbinance.NewListener(securityID, dialerPool) }
	case constants.FTX.ID:
		return func() actor.Actor { return ftx.NewListener(securityID, dialerPool) }
	case constants.FTXUS.ID:
		return func() actor.Actor { return ftxus.NewListener(securityID, dialerPool) }
	case constants.GEMINI.ID:
		return func() actor.Actor { return gemini.NewListener(securityID, dialerPool) }
	case constants.HITBTC.ID:
		return func() actor.Actor { return hitbtc.NewListener(securityID, dialerPool) }
	case constants.KRAKEN.ID:
		return func() actor.Actor { return kraken.NewListener(securityID, dialerPool) }
	case constants.OKCOIN.ID:
		return func() actor.Actor { return okcoin.NewListener(securityID, dialerPool) }
	case constants.OKEX.ID:
		return func() actor.Actor { return okex.NewListener(securityID, dialerPool) }
	case constants.DERIBIT.ID:
		return func() actor.Actor { return deribit.NewListener(securityID, dialerPool) }
	case constants.HUOBI.ID:
		return func() actor.Actor { return huobi.NewListener(securityID, dialerPool) }
	case constants.HUOBIP.ID:
		return func() actor.Actor { return huobip.NewListener(securityID, dialerPool) }
	case constants.HUOBIF.ID:
		return func() actor.Actor { return huobif.NewListener(securityID, dialerPool) }
	case constants.HUOBIL.ID:
		return func() actor.Actor { return huobil.NewListener(securityID, dialerPool) }
	case constants.BYBITI.ID:
		return func() actor.Actor { return bybiti.NewListener(securityID, dialerPool) }
	case constants.BYBITL.ID:
		return func() actor.Actor { return bybitl.NewListener(securityID, dialerPool) }
	case constants.BYBITS.ID:
		return func() actor.Actor { return bybits.NewListener(securityID, dialerPool) }
	case constants.UPBIT.ID:
		return func() actor.Actor { return upbit.NewListener(securityID, dialerPool) }
	case constants.BITHUMB.ID:
		return func() actor.Actor { return bithumb.NewListener(securityID, dialerPool) }
	case constants.BITHUMBG.ID:
		return func() actor.Actor { return bithumbg.NewListener(securityID, dialerPool) }
	case constants.DYDX.ID:
		return func() actor.Actor { return dydx.NewListener(securityID, dialerPool) }
	case constants.OKEXP.ID:
		return func() actor.Actor { return okexp.NewListener(securityID, dialerPool) }
	case constants.GATE.ID:
		return func() actor.Actor { return gate.NewListener(securityID, dialerPool) }
	case constants.UNISWAPV3.ID:
		return func() actor.Actor { return v3.NewListener(securityID, dialerPool) }
		/*
			case constants.BITTREX:
			return func() actor.Actor { return bittrex.NewListener(instrument) }
		*/
	default:
		return nil
	}
}

func NewExchangeExecutorProducer(exchange *models2.Exchange, config *types.ExecutorConfig) actor.Producer {
	switch exchange.ID {
	case constants.BINANCE.ID:
		return func() actor.Actor { return binance.NewExecutor(config) }
	case constants.BITFINEX.ID:
		return func() actor.Actor { return bitfinex.NewExecutor(config) }
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
		return func() actor.Actor { return fbinance.NewExecutor(config) }
	case constants.FTX.ID:
		return func() actor.Actor { return ftx.NewExecutor(config) }
	case constants.FTXUS.ID:
		return func() actor.Actor { return ftxus.NewExecutor(config) }
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
		return func() actor.Actor { return huobip.NewExecutor(config) }
	case constants.HUOBIF.ID:
		return func() actor.Actor { return huobif.NewExecutor() }
	case constants.HUOBIL.ID:
		return func() actor.Actor { return huobil.NewExecutor(config) }
	case constants.BYBITI.ID:
		return func() actor.Actor { return bybiti.NewExecutor(config) }
	case constants.BYBITL.ID:
		return func() actor.Actor { return bybitl.NewExecutor(config) }
	case constants.BYBITS.ID:
		return func() actor.Actor { return bybits.NewExecutor(config) }
	case constants.UPBIT.ID:
		return func() actor.Actor { return upbit.NewExecutor() }
	case constants.BITHUMB.ID:
		return func() actor.Actor { return bithumb.NewExecutor() }
	case constants.BITHUMBG.ID:
		return func() actor.Actor { return bithumbg.NewExecutor() }
	case constants.DYDX.ID:
		return func() actor.Actor { return dydx.NewExecutor(config) }
	case constants.OKEXP.ID:
		return func() actor.Actor { return okexp.NewExecutor(config) }
	case constants.GATE.ID:
		return func() actor.Actor { return gate.NewExecutor(config) }
	case constants.UNISWAPV3.ID:
		return func() actor.Actor { return v3.NewExecutor(config) }
	case constants.OPENSEA.ID:
		return func() actor.Actor { return opensea.NewExecutor(config) }
		/*
			case constants.BITTREX:
				return func() actor.Actor { return bittrex.NewExecutor() }
		*/
	default:
		return nil
	}
}
