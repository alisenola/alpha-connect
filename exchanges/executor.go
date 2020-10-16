package exchanges

import (
	"errors"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alphac/account"
	"gitlab.com/alphaticks/alphac/models"
	"gitlab.com/alphaticks/alphac/models/messages"
	"gitlab.com/alphaticks/alphac/utils"
	xchangerModels "gitlab.com/alphaticks/xchanger/models"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
	"reflect"
	"time"
)

// The executor routes all the request to the underlying exchange executor & listeners
// He is the main part of the whole software..
//
type Executor struct {
	exchanges         []*xchangerModels.Exchange
	accounts          []*account.Account
	accountPortfolios map[string]*actor.PID
	accountManagers   map[string]*actor.PID
	executors         map[uint32]*actor.PID       // A map from exchange ID to executor
	securities        map[uint64]*models.Security // A map from security ID to security
	instruments       map[uint64]*actor.PID       // A map from security ID to market manager
	slSubscribers     map[uint64]*actor.PID       // A map from request ID to security list subscribers
	execSubscribers   map[uint64]*actor.PID       // A map from request ID to execution report subscribers
	logger            *log.Logger
	paperTrading      bool
	dialerPool        *xchangerUtils.DialerPool
}

func NewExecutorProducer(exchanges []*xchangerModels.Exchange, accounts []*account.Account, paperTrading bool, dialerPool *xchangerUtils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewExecutor(exchanges, accounts, paperTrading, dialerPool)
	}
}

func NewExecutor(exchanges []*xchangerModels.Exchange, accounts []*account.Account, paperTrading bool, dialerPool *xchangerUtils.DialerPool) actor.Actor {
	return &Executor{
		exchanges:    exchanges,
		accounts:     accounts,
		logger:       nil,
		paperTrading: paperTrading,
		dialerPool:   dialerPool,
	}
}

func (state *Executor) Receive(context actor.Context) {
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
			// Attention, no panic in restarting or infinite loop
		}
		state.logger.Info("actor restarting")

	case *messages.MarketDataRequest:
		if err := state.OnMarketDataRequest(context); err != nil {
			state.logger.Error("error processing OnMarketDataRequest", log.Error(err))
			panic(err)
		}

	case *messages.SecurityDefinitionRequest:
		if err := state.OnSecurityDefinitionRequest(context); err != nil {
			state.logger.Error("error processing OnSecurityDefinitionRequest", log.Error(err))
			panic(err)
		}

	case *messages.SecurityListRequest:
		if err := state.OnSecurityListRequest(context); err != nil {
			state.logger.Error("error processing OnSecurityListRequest", log.Error(err))
			panic(err)
		}

	case *messages.SecurityList:
		if err := state.OnSecurityList(context); err != nil {
			state.logger.Error("error processing OnSecurityList", log.Error(err))
			panic(err)
		}

	case *messages.PositionsRequest:
		if err := state.OnPositionsRequest(context); err != nil {
			state.logger.Error("error processing OnPositionRequest", log.Error(err))
			panic(err)
		}

	case *messages.BalancesRequest:
		if err := state.OnBalancesRequest(context); err != nil {
			state.logger.Error("error processing OnBalancesRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderStatusRequest:
		if err := state.OnOrderStatusRequest(context); err != nil {
			state.logger.Error("error processing OnOrderStatusRequest", log.Error(err))
			panic(err)
		}

	case *messages.NewOrderSingleRequest:
		if err := state.OnNewOrderSingleRequest(context); err != nil {
			state.logger.Error("error processing OnNewOrderSingle", log.Error(err))
			panic(err)
		}

	case *messages.NewOrderBulkRequest:
		if err := state.OnNewOrderBulkRequest(context); err != nil {
			state.logger.Error("error processing OnNewOrderBulk", log.Error(err))
			panic(err)
		}

	case *messages.OrderReplaceRequest:
		if err := state.OnOrderReplaceRequest(context); err != nil {
			state.logger.Error("error processing OrderReplaceRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderBulkReplaceRequest:
		if err := state.OnOrderBulkReplaceRequest(context); err != nil {
			state.logger.Error("error processing OrderBulkReplaceRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderCancelRequest:
		if err := state.OnOrderCancelRequest(context); err != nil {
			state.logger.Error("error processing OnOrderCancelRequest", log.Error(err))
			panic(err)
		}

	case *messages.OrderMassCancelRequest:
		if err := state.OnOrderMassCancelRequest(context); err != nil {
			state.logger.Error("error processing OnOrderMassCancelRequest", log.Error(err))
			panic(err)
		}

	case *actor.Terminated:
		if err := state.OnTerminated(context); err != nil {
			state.logger.Error("error processing OnTerminated", log.Error(err))
			panic(err)
		}
	}
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	state.instruments = make(map[uint64]*actor.PID)
	state.slSubscribers = make(map[uint64]*actor.PID)

	// Spawn all exchange executors
	state.executors = make(map[uint32]*actor.PID)
	for _, exch := range state.exchanges {
		producer := NewExchangeExecutorProducer(exch)
		if producer == nil {
			return fmt.Errorf("unknown exchange %s", exch.Name)
		}
		props := actor.PropsFromProducer(producer)
		state.executors[exch.ID], _ = context.SpawnNamed(props, exch.Name+"_executor")
	}

	// Request securities for each one of them
	var futures []*actor.Future
	request := &messages.SecurityListRequest{
		RequestID: 0,
		Subscribe: true,
	}
	for _, pid := range state.executors {
		fut := context.RequestFuture(pid, request, 20*time.Second)
		futures = append(futures, fut)
	}

	state.securities = make(map[uint64]*models.Security)
	for _, fut := range futures {
		res, err := fut.Result()
		if err != nil {
			state.logger.Error("error fetching securities for one venue", log.Error(err))
			continue
		}
		response, ok := res.(*messages.SecurityList)
		if !ok {
			return fmt.Errorf("was expecting GetSecuritiesResponse, got %s", reflect.TypeOf(res).String())
		}
		if !response.Success {
			return errors.New(response.RejectionReason.String())
		}
		for _, s := range response.Securities {
			if sec2, ok := state.securities[s.SecurityID]; ok {
				return fmt.Errorf("got two securities with the same ID: %s %s", sec2.Symbol, s.Symbol)
			}
			state.securities[s.SecurityID] = s
		}
	}

	// Spawn all account listeners
	state.accountManagers = make(map[string]*actor.PID)
	for _, accnt := range state.accounts {
		producer := NewAccountManagerProducer(accnt, state.paperTrading)
		if producer == nil {
			return fmt.Errorf("unknown exchange %s", accnt.Exchange.Name)
		}
		props := actor.PropsFromProducer(producer).WithSupervisor(
			actor.NewExponentialBackoffStrategy(100*time.Second, time.Second))
		state.accountManagers[accnt.AccountID] = context.Spawn(props)
	}

	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) OnMarketDataRequest(context actor.Context) error {
	request := context.Message().(*messages.MarketDataRequest)
	if request.Instrument == nil || request.Instrument.SecurityID == nil {
		context.Respond(&messages.MarketDataResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: messages.UnknownSecurityID,
		})
		return nil
	}
	securityID := request.Instrument.SecurityID.Value
	security, ok := state.securities[securityID]
	if !ok {
		context.Respond(&messages.MarketDataResponse{
			RequestID:       request.RequestID,
			Success:         false,
			RejectionReason: messages.UnknownSecurityID,
		})
		return nil
	}
	if pid, ok := state.instruments[securityID]; ok {
		context.Forward(pid)
	} else {
		props := actor.PropsFromProducer(NewMarketDataManagerProducer(security, state.dialerPool)).WithSupervisor(
			utils.NewExponentialBackoffStrategy(100*time.Second, time.Second, time.Second))
		pid := context.Spawn(props)
		state.instruments[securityID] = pid
		context.Forward(pid)
	}

	return nil
}

func (state *Executor) OnSecurityDefinitionRequest(context actor.Context) error {
	request := context.Message().(*messages.SecurityDefinitionRequest)
	if request.Instrument == nil || request.Instrument.SecurityID == nil {
		context.Respond(&messages.SecurityDefinitionResponse{
			RequestID:       request.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Security:        nil,
			RejectionReason: messages.UnknownSecurityID,
		})
	}
	if sec, ok := state.securities[request.Instrument.SecurityID.Value]; ok {
		context.Respond(&messages.SecurityDefinitionResponse{
			RequestID:  request.RequestID,
			ResponseID: uint64(time.Now().UnixNano()),
			Security:   sec,
		})
	} else {
		context.Respond(&messages.SecurityDefinitionResponse{
			RequestID:       request.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Security:        nil,
			RejectionReason: messages.UnknownSecurityID,
		})
	}
	return nil
}

func (state *Executor) OnSecurityListRequest(context actor.Context) error {
	request := context.Message().(*messages.SecurityListRequest)
	var securities []*models.Security
	for _, v := range state.securities {
		securities = append(securities, v)
	}
	response := &messages.SecurityList{
		RequestID:  request.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Securities: securities,
		Success:    true,
	}
	if request.Subscribe {
		context.Watch(request.Subscriber)
		state.slSubscribers[request.RequestID] = request.Subscriber
	}
	context.Respond(response)

	return nil
}

func (state *Executor) OnSecurityList(context actor.Context) error {
	securityList := context.Message().(*messages.SecurityList)
	// Do nothing
	if len(securityList.Securities) == 0 {
		return nil
	}
	exchangeID := securityList.Securities[0].Exchange.ID
	// It has to come from one exchange only, so delete all the known securities from that exchange
	for k, v := range state.securities {
		if v.Exchange.ID == exchangeID {
			delete(state.securities, k)
		}
	}
	// re-add them
	for _, s := range securityList.Securities {
		state.securities[s.SecurityID] = s
	}

	var securities []*models.Security
	for _, v := range state.securities {
		securities = append(securities, v)
	}
	for k, v := range state.slSubscribers {
		securityList := &messages.SecurityList{
			RequestID:  k,
			ResponseID: uint64(time.Now().UnixNano()),
			Securities: securities,
			Success:    true,
		}
		context.Send(v, securityList)
	}

	return nil
}

func (state *Executor) OnPositionsRequest(context actor.Context) error {
	msg := context.Message().(*messages.PositionsRequest)
	if msg.Account == nil {
		context.Respond(&messages.PositionList{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.InvalidAccount,
			Positions:       nil,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.AccountID]
	if !ok {
		context.Respond(&messages.PositionList{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.InvalidAccount,
			Positions:       nil,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnBalancesRequest(context actor.Context) error {
	msg := context.Message().(*messages.BalancesRequest)
	if msg.Account == nil {
		context.Respond(&messages.BalanceList{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.InvalidAccount,
			Balances:        nil,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.AccountID]
	if !ok {
		context.Respond(&messages.BalanceList{
			RequestID:       msg.RequestID,
			ResponseID:      uint64(time.Now().UnixNano()),
			Success:         false,
			RejectionReason: messages.InvalidAccount,
			Balances:        nil,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnOrderStatusRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderStatusRequest)
	if msg.Account == nil {
		context.Respond(&messages.OrderList{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.AccountID]
	if !ok {
		context.Respond(&messages.OrderList{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnNewOrderSingleRequest(context actor.Context) error {
	msg := context.Message().(*messages.NewOrderSingleRequest)
	if msg.Account == nil {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	if msg.Order == nil {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidRequest,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.AccountID]
	if !ok {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnNewOrderBulkRequest(context actor.Context) error {
	msg := context.Message().(*messages.NewOrderBulkRequest)
	if msg.Account == nil {
		context.Respond(&messages.NewOrderBulkResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.AccountID]
	if !ok {
		context.Respond(&messages.NewOrderBulkResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	if len(msg.Orders) == 0 {
		context.Respond(&messages.NewOrderBulkResponse{
			RequestID: msg.RequestID,
			Success:   true,
			OrderIDs:  nil,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnOrderReplaceRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderReplaceRequest)
	if msg.Account == nil {
		context.Respond(&messages.OrderReplaceResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.AccountID]
	if !ok {
		context.Respond(&messages.OrderReplaceResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnOrderBulkReplaceRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderBulkReplaceRequest)
	if msg.Account == nil {
		context.Respond(&messages.OrderBulkReplaceResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	accountManager, ok := state.accountManagers[msg.Account.AccountID]
	if !ok {
		context.Respond(&messages.OrderBulkReplaceResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnOrderCancelRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderCancelRequest)
	if msg.Account == nil || msg.Instrument == nil || msg.Instrument.Symbol == nil || (msg.ClientOrderID == nil && msg.OrderID == nil) {
		context.Respond(&messages.OrderCancelResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.UnknownOrder,
		})
	}
	accountManager, ok := state.accountManagers[msg.Account.AccountID]
	if !ok {
		context.Respond(&messages.OrderCancelResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.UnknownOrder,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnOrderMassCancelRequest(context actor.Context) error {
	msg := context.Message().(*messages.OrderMassCancelRequest)
	if msg.Account == nil {
		context.Respond(&messages.OrderCancelResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
	}
	accountManager, ok := state.accountManagers[msg.Account.AccountID]
	if !ok {
		context.Respond(&messages.OrderCancelResponse{
			RequestID:       msg.RequestID,
			Success:         false,
			RejectionReason: messages.InvalidAccount,
		})
		return nil
	}
	context.Forward(accountManager)
	return nil
}

func (state *Executor) OnTerminated(context actor.Context) error {
	// Handle subscriber krash
	msg := context.Message().(*actor.Terminated)
	for k, v := range state.slSubscribers {
		if v.Id == msg.Who.Id {
			delete(state.slSubscribers, k)
		}
	}

	for k, v := range state.instruments {
		if v.Id == msg.Who.Id {
			delete(state.instruments, k)
		}
	}

	return nil
}
