package v3

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"github.com/hasura/go-graphql-client"
	"gitlab.com/alphaticks/alpha-connect/enum"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/jobs"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/alpha-connect/utils"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/uniswap"
	xutils "gitlab.com/alphaticks/xchanger/utils"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"time"
)

type QueryRunner struct {
	pid       *actor.PID
}

type Executor struct {
	extypes.BaseExecutor
	securities  []*models.Security
	queryRunners []*QueryRunner
	dialerPool   *xutils.DialerPool
	logger      *log.Logger
}

func NewExecutor(dialerPool *xutils.DialerPool) actor.Actor {
	return &Executor{
		queryRunners: nil,
		logger:       nil,
		dialerPool:   dialerPool,
	}
}

func (state *Executor) getQueryRunner() *QueryRunner {
	sort.Slice(state.queryRunners, func(i, j int) bool {
		return rand.Uint64()%2 == 0
	})

	return state.queryRunners[0]
}

func (state *Executor) Receive(context actor.Context) {
	extypes.ReceiveExecutor(state, context)
}

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	dialers := state.dialerPool.GetDialers()
	for _, dialer := range dialers {
		httpClient := &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 1024,
				TLSHandshakeTimeout: 10 * time.Second,
				DialContext:         dialer.DialContext,
			},
			Timeout: 10 * time.Second,
		}
		uniClient := graphql.NewClient(uniswap.APIURL, httpClient)
		props := actor.PropsFromProducer(func() actor.Actor {
			return jobs.NewGraphQuery(uniClient)
		})
		state.queryRunners = append(state.queryRunners, &QueryRunner{
			pid:             context.Spawn(props),
		})
	}
	return state.UpdateSecurityList(context)
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) UpdateSecurityList(context actor.Context) error {

	// https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2

	var securities []*models.Security

	query := uniswap.Pools{}
	qr := state.getQueryRunner()
	if qr == nil {
		return fmt.Errorf("rate limited")
	}

	future := context.RequestFuture(qr.pid, &jobs.PerformGraphQueryRequest{Query: &query}, 10*time.Second)
	context.AwaitFuture(future, func(resp interface{}, err error) {
		for _, pool := range query.Pools {
			baseCurrency, ok := constants.GetAssetBySymbol(pool.Token0)
			if !ok {
				//state.logger.Info("unknown symbol " + pair.BaseCurrency + " for instrument " + pair.InstrumentID)
				continue
			}
			quoteCurrency, ok := constants.GetAssetBySymbol(pool.Token1)
			if !ok {
				//state.logger.Info("unknown symbol " + pair.QuoteCurrency + " for instrument " + pair.InstrumentID)
				continue
			}

			security := models.Security{}
			security.Symbol = fmt.Sprintf("%s/%s", pool.Token0.Symbol, pool.Token1.Symbol)
			security.Underlying = baseCurrency
			security.QuoteCurrency = quoteCurrency
			security.Status = models.Trading
			security.Exchange = &constants.UNISWAPV3
			security.IsInverse = false
			security.SecurityType = enum.SecurityType_CRYPTO_AMM
			security.SecuritySubType = &types.StringValue{Value: enum.SecuritySubType_UNIPOOLV3}
			security.SecurityID = utils.SecurityID(security.SecurityType, security.Symbol, security.Exchange.Name, security.MaturityDate)
			security.MinPriceIncrement = &types.DoubleValue{Value: pool.TickSize} // TODO in bps ?
			security.RoundLot = &types.DoubleValue{Value: pool.SizeIncrement} // TODO Token precision ?
			security.TakerFee = // TODO pool fees
			securities = append(securities, &security)
		}

		state.securities = securities

		context.Send(context.Parent(), &messages.SecurityList{
			ResponseID: uint64(time.Now().UnixNano()),
			Success:    true,
			Securities: state.securities})
	})

	return nil
}

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	msg := context.Message().(*messages.SecurityListRequest)
	context.Respond(&messages.SecurityList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: state.securities})

	return nil
}
