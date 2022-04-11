package bybitl

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/account"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/bybitl"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"reflect"
	"time"
)

type checkSocket struct{}
type checkAccount struct{}
type refreshKey struct{}

type AccountListener struct {
	account            *account.Account
	seqNum             uint64
	bybitlExecutor     *actor.PID
	ws                 *bybitl.Websocket
	logger             *log.Logger
	checkAccountTicker *time.Ticker
	checkSocketTicker  *time.Ticker
	refreshKeyTicker   *time.Ticker
	lastPingTime       time.Time
	securities         map[uint64]*models.Security
	client             *http.Client
	txs                *mongo.Collection
	execs              *mongo.Collection
	//TODO Needed?
	//reconciler *actor.PID
}

func NewAccountListenerProducer(account *account.Account, txs, execs *mongo.Collection) actor.Producer {
	return func() actor.Actor {
		return NewAccountListener(account, txs, execs)
	}
}

func NewAccountListener(account *account.Account, txs, execs *mongo.Collection) actor.Actor {
	return &AccountListener{
		account: account,
		seqNum:  0,
		ws:      nil,
		logger:  nil,
		txs:     txs,
		execs:   execs,
	}
}

func (state *AccountListener) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.logger.Error("error initializing", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor started")
	case *actor.Stopping:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error stopping", log.Error(err))
			panic(err)
		}
		state.logger.Info("actor stopping")
	case *actor.Stopped:
		state.logger.Info("actor stopped")
	case *actor.Restarting:
		if err := state.Clean(context); err != nil {
			state.logger.Error("error restarting", log.Error(err))
		}
		state.logger.Info("actor restarting")
	case *messages.BalancesRequest:
		if err := state.OnBalancesRequest(context); err != nil {
			state.logger.Error("error processing OnBalancesRequest", log.Error(err))
			panic(err)
		}
	}
}

func (state *AccountListener) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
	)
	state.bybitlExecutor = actor.NewPID(context.ActorSystem().Address(), "executor/"+constants.BYBITL.Name+"_executor")
	state.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	if err := state.subscribeAccount(context); err != nil {
		return fmt.Errorf("error subscribing to account: %v", err)
	}

	//Fetch the securities
	fmt.Println("Fetching the securities")
	ex := actor.NewPID(context.ActorSystem().Address(), "executor")
	res, err := context.RequestFuture(ex, &messages.SecurityListRequest{}, 10*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error getting securities: %v", err)
	}
	securityList, ok := res.(*messages.SecurityList)
	if !ok {
		return fmt.Errorf("was expecting *messages.SecurityList, got %s", reflect.TypeOf(res).String())
	}
	if !securityList.Success {
		return fmt.Errorf("error getting securities: %s", securityList.RejectionReason.String())
	}
	var filteredSecurities []*models.Security
	for _, s := range securityList.Securities {
		if s.Exchange.ID == state.account.Exchange.ID {
			filteredSecurities = append(filteredSecurities, s)
		}
	}

	//Fetch the current balance
	fmt.Println("Fetching the balance")
	res, err = context.RequestFuture(state.bybitlExecutor, &messages.BalancesRequest{
		Account: state.account.Account,
	}, 10*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error getting balances from executor: %v", err)
	}
	balances, ok := res.(*messages.BalanceList)
	if !ok {
		return fmt.Errorf("was expecting *messages.BalanceList, got %s", reflect.TypeOf(res).String())
	}
	if !balances.Success {
		return fmt.Errorf("error getting balances: %s", balances.RejectionReason.String())
	}

	//Fetch the positions
	fmt.Println("Fetching the positions")
	res, err = context.RequestFuture(state.bybitlExecutor, &messages.PositionsRequest{
		Instrument: nil,
		Account:    state.account.Account,
	}, 10*time.Second).Result()
	if err != nil {
		return fmt.Errorf("error getting positions from executor: %v", err)
	}
	positions, ok := res.(*messages.PositionList)
	if !ok {
		return fmt.Errorf("was expecting *messages.PositionsList, got %s", reflect.TypeOf(res).String())
	}
	if !positions.Success {
		return fmt.Errorf("error getting positions: %s", positions.RejectionReason.String())
	}

	//Fetch the orders
	fmt.Println("Fetching the orders")
	var ords []*models.Order
	securityMap := make(map[uint64]*models.Security)
	for _, fs := range filteredSecurities {
		res, err = context.RequestFuture(state.bybitlExecutor, &messages.OrderStatusRequest{
			Account: state.account.Account,
			Filter: &messages.OrderFilter{
				Instrument: &models.Instrument{
					SecurityID: &types.UInt64Value{Value: fs.SecurityID},
					Symbol:     &types.StringValue{Value: fs.Symbol},
					Exchange:   fs.Exchange,
				},
			},
		}, 10*time.Second).Result()
		if err != nil {
			return fmt.Errorf("error getting orders from executor: %v", err)
		}
		orders, ok := res.(*messages.OrderList)
		if !ok {
			return fmt.Errorf("was expecting *messages.OrderList, got %s", reflect.TypeOf(res).String())
		}
		if !orders.Success {
			return fmt.Errorf("error getting orders: %s", orders.RejectionReason.String())
		}
		ords = append(ords, orders.Orders...)
		securityMap[fs.SecurityID] = fs
	}
	state.securities = securityMap
	state.seqNum = 0

	//Sync account
	makerFee := 0.0001
	takerFee := 0.0006
	if err := state.account.Sync(filteredSecurities, ords, positions.Positions, balances.Balances, &makerFee, &takerFee); err != nil {
		return fmt.Errorf("error syncing account: %v", err)
	}

	checkAccountTicker := time.NewTicker(5 * time.Minute)
	state.checkAccountTicker = checkAccountTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-checkAccountTicker.C:
				context.Send(pid, &checkAccount{})
			case <-time.After(6 * time.Minute):
				// timer stopped, we leave
				return
			}
		}
	}(context.Self())

	checkSocketTicker := time.NewTicker(5 * time.Second)
	state.checkSocketTicker = checkSocketTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-checkSocketTicker.C:
				context.Send(pid, &checkSocket{})
			case <-time.After(10 * time.Second):
				// timer stopped, we leave
				return
			}
		}
	}(context.Self())

	refreshKeyTicker := time.NewTicker(30 * time.Minute)
	state.refreshKeyTicker = refreshKeyTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-refreshKeyTicker.C:
				context.Send(pid, &refreshKey{})
			case <-time.After(31 * time.Minute):
				// timer stopped, we leave
				return
			}
		}
	}(context.Self())

	return nil
}

func (state *AccountListener) subscribeAccount(context actor.Context) error {
	if state.ws != nil {
		_ = state.ws.Disconnect()
	}
	
}

func (state *AccountListener) Clean(context actor.Context) error {
	if state.ws != nil {
		if err := state.ws.Disconnect(); err != nil {
			state.logger.Info("error disconnecting websocket", log.Error(err))
		}
	}

	if state.checkAccountTicker != nil {
		state.checkAccountTicker.Stop()
		state.checkAccountTicker = nil
	}

	if state.checkSocketTicker != nil {
		state.checkSocketTicker.Stop()
		state.checkSocketTicker = nil
	}

	if state.refreshKeyTicker != nil {
		state.refreshKeyTicker.Stop()
		state.refreshKeyTicker = nil
	}

	return nil
}
