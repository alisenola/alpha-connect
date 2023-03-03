package exchanges

import (
	"github.com/asynkron/protoactor-go/actor"
	"gitlab.com/alphaticks/alpha-connect/account"
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/exchanges/binance"
	"gitlab.com/alphaticks/alpha-connect/exchanges/bitfinex"
	"gitlab.com/alphaticks/alpha-connect/exchanges/bithumb"
	"gitlab.com/alphaticks/alpha-connect/exchanges/bithumbg"
	"gitlab.com/alphaticks/alpha-connect/exchanges/bitmex"
	"gitlab.com/alphaticks/alpha-connect/exchanges/bitstamp"
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
	"gitlab.com/alphaticks/alpha-connect/exchanges/gatef"
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
	v3 "gitlab.com/alphaticks/alpha-connect/exchanges/uniswap/v3"
	"gitlab.com/alphaticks/alpha-connect/exchanges/upbit"
	"gitlab.com/alphaticks/alpha-connect/models"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	tickstore_types "gitlab.com/alphaticks/tickstore-types"
	"gitlab.com/alphaticks/xchanger/constants"
	models2 "gitlab.com/alphaticks/xchanger/models"
	"gitlab.com/alphaticks/xchanger/utils"
	xutils "gitlab.com/alphaticks/xchanger/utils"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

// TODO sample size config
var portfolios = make(map[string]*account.Portfolio)
var accounts = make(map[string]*account.Account)
var portfolioMx = &sync.Mutex{}

func GetAccount(name string) *account.Account {
	portfolioMx.Lock()
	defer portfolioMx.Unlock()
	return accounts[name]
}

func GetPortfolio(ID string) *account.Portfolio {
	portfolioMx.Lock()
	defer portfolioMx.Unlock()
	portfolio, ok := portfolios[ID]
	if !ok {
		portfolio = account.NewPortfolio(100)
		portfolios[ID] = portfolio
	}
	return portfolio
}

func getPortfolio(ID string) *account.Portfolio {
	portfolio, ok := portfolios[ID]
	if !ok {
		portfolio = account.NewPortfolio(100)
		portfolios[ID] = portfolio
	}
	return portfolio
}

func NewAccount(accountInfo *models.Account, fillCollector *account.FillCollector) (*account.Account, error) {
	portfolioMx.Lock()
	defer portfolioMx.Unlock()
	acc, ok := accounts[accountInfo.Name]
	if ok {
		return acc, nil
	}
	accnt, err := account.NewAccount(accountInfo, fillCollector)
	if err != nil {
		return nil, err
	}
	accounts[accnt.Name] = accnt
	getPortfolio(accountInfo.Portfolio).AddAccount(accnt)
	return accnt, nil
}

func NewAccountListenerProducer(account *account.Account, registry registry.StaticClient, db *gorm.DB, client *http.Client, readOnly bool) actor.Producer {
	switch account.Exchange.ID {
	case constants.BITMEX.ID:
		return func() actor.Actor { return bitmex.NewAccountListener(account) }
	case constants.BINANCE.ID:
		return func() actor.Actor { return binance.NewAccountListener(account, nil, nil) }
	case constants.FBINANCE.ID:
		return func() actor.Actor { return fbinance.NewAccountListener(account, registry, db, client, readOnly) }
	case constants.FTX.ID:
		return func() actor.Actor { return ftx.NewAccountListener(account, registry, db, readOnly) }
	case constants.FTXUS.ID:
		return func() actor.Actor { return ftxus.NewAccountListener(account) }
	case constants.DYDX.ID:
		return func() actor.Actor { return dydx.NewAccountListener(account, nil, nil) }
	case constants.BYBITL.ID:
		return func() actor.Actor { return bybitl.NewAccountListener(account, registry, db, client, readOnly) }
	default:
		return nil
	}
}

func NewAccountReconcileProducer(accountCfg config.Account, account *models.Account, registry registry.StaticClient, store tickstore_types.TickstoreClient, db *gorm.DB) actor.Producer {
	switch account.Exchange.ID {
	case constants.BINANCE.ID:
		return binance.NewAccountReconcileProducer(accountCfg, account, registry, store, db)
	case constants.FBINANCE.ID:
		return fbinance.NewAccountReconcileProducer(accountCfg, account, registry, store, db)
	case constants.FTX.ID:
		return ftx.NewAccountReconcileProducer(account, registry, db)
	case constants.BYBITL.ID:
		return bybitl.NewAccountReconcileProducer(account, registry, db)
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

func NewInstrumentListenerProducer(securityID uint64, exchangeID uint32, dialerPool *utils.DialerPool, wsPool *utils.WebsocketPool) actor.Producer {
	switch exchangeID {
	case constants.BINANCE.ID:
		return func() actor.Actor { return binance.NewListener(securityID, dialerPool) }
	case constants.BITFINEX.ID:
		return func() actor.Actor { return bitfinex.NewListener(securityID, dialerPool) }
	case constants.BITMEX.ID:
		return func() actor.Actor { return bitmex.NewListener(securityID, dialerPool) }
	case constants.BITSTAMP.ID:
		return func() actor.Actor { return bitstamp.NewListenerL3(securityID, dialerPool) }
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
		return func() actor.Actor { return hitbtc.NewListener(securityID, wsPool) }
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
	case constants.GATEF.ID:
		return func() actor.Actor { return gatef.NewListener(securityID, dialerPool) }
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

func NewExchangeExecutorProducer(exchange *models2.Exchange, config *config.Config, dialerPool *xutils.DialerPool, registry registry.StaticClient, accountClients map[string]*http.Client) actor.Producer {
	switch exchange.ID {
	case constants.BINANCE.ID:
		return func() actor.Actor { return binance.NewExecutor(dialerPool, registry) }
	case constants.BITFINEX.ID:
		return func() actor.Actor { return bitfinex.NewExecutor(dialerPool, registry) }
	case constants.BITMEX.ID:
		return func() actor.Actor { return bitmex.NewExecutor() }
	case constants.BITSTAMP.ID:
		return func() actor.Actor { return bitstamp.NewExecutor() }
	case constants.COINBASEPRO.ID:
		return func() actor.Actor { return coinbasepro.NewExecutor() }
	case constants.CRYPTOFACILITIES.ID:
		return func() actor.Actor { return cryptofacilities.NewExecutor() }
	case constants.FBINANCE.ID:
		return func() actor.Actor { return fbinance.NewExecutor(config, dialerPool, registry, accountClients) }
	case constants.FTX.ID:
		return func() actor.Actor { return ftx.NewExecutor(dialerPool, registry) }
	case constants.FTXUS.ID:
		return func() actor.Actor { return ftxus.NewExecutor(dialerPool, registry) }
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
		return func() actor.Actor { return huobip.NewExecutor(dialerPool, registry) }
	case constants.HUOBIF.ID:
		return func() actor.Actor { return huobif.NewExecutor() }
	case constants.HUOBIL.ID:
		return func() actor.Actor { return huobil.NewExecutor(dialerPool, registry) }
	case constants.BYBITI.ID:
		return func() actor.Actor { return bybiti.NewExecutor(dialerPool, registry) }
	case constants.BYBITL.ID:
		return func() actor.Actor { return bybitl.NewExecutor(dialerPool, registry, accountClients) }
	case constants.BYBITS.ID:
		return func() actor.Actor { return bybits.NewExecutor(dialerPool, registry) }
	case constants.UPBIT.ID:
		return func() actor.Actor { return upbit.NewExecutor() }
	case constants.BITHUMB.ID:
		return func() actor.Actor { return bithumb.NewExecutor() }
	case constants.BITHUMBG.ID:
		return func() actor.Actor { return bithumbg.NewExecutor() }
	case constants.DYDX.ID:
		return func() actor.Actor { return dydx.NewExecutor(dialerPool, registry) }
	case constants.OKEXP.ID:
		return func() actor.Actor { return okexp.NewExecutor(dialerPool, registry) }
	case constants.GATE.ID:
		return func() actor.Actor { return gate.NewExecutor(dialerPool, registry) }
	case constants.GATEF.ID:
		return func() actor.Actor { return gatef.NewExecutor(dialerPool, registry) }
	case constants.UNISWAPV3.ID:
		return func() actor.Actor { return v3.NewExecutor(dialerPool, registry) }
	case constants.OPENSEA.ID:
		return func() actor.Actor { return opensea.NewExecutor(config, dialerPool, registry) }
		/*
			case constants.BITTREX:
				return func() actor.Actor { return bittrex.NewExecutor() }
		*/
	default:
		return nil
	}
}
