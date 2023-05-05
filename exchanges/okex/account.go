package okex

import (
	"fmt"
	"math"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/account"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"
	"gitlab.com/alphaticks/xchanger"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges/okex"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

type checkSocket struct{}
type checkAccount struct{}
type refreshKey struct{}
type refreshMarkPrices struct{}

type AccountListener struct {
	account                 *account.Account
	readOnly                bool
	seqNum                  uint64
	okexExecutor            *actor.PID
	ws                      *okex.Websocket
	executorManager         *actor.PID
	logger                  *log.Logger
	registry                registry.StaticClient
	checkAccountTicker      *time.Ticker
	checkSocketTicker       *time.Ticker
	refreshKeyTicker        *time.Ticker
	refreshMarkPricesTicker *time.Ticker
	lastPingTime            time.Time
	securities              map[uint64]*models.Security
	symbolToSec             map[string]*models.Security
	client                  *http.Client
	db                      *gorm.DB
}

func NewAccountListenerProducer(account *account.Account, registry registry.StaticClient, db *gorm.DB, client *http.Client, strict bool) actor.Producer {
	return func() actor.Actor {
		return NewAccountListener(account, registry, db, client, strict)
	}
}

func NewAccountListener(account *account.Account, registry registry.StaticClient, db *gorm.DB, client *http.Client, readOnly bool) actor.Actor {
	return &AccountListener{
		account:         account,
		readOnly:        readOnly,
		registry:        registry,
		seqNum:          0,
		ws:              nil,
		executorManager: nil,
		logger:          nil,
		db:              db,
		client:          client,
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
			// Attention, no panic in restarting or infinite loop
		}
		state.logger.Info("actor restarting")

	case *messages.AccountDataRequest:
		if err := state.OnAccountDataRequest(context); err != nil {
			state.logger.Error("error processing OnAccountDataRequest", log.Error(err))
			panic(err)
		}

	case *messages.PositionsRequest:
		if err := state.OnPositionsRequest(context); err != nil {
			state.logger.Error("error processing OnPositionListRequest", log.Error(err))
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

	case *messages.AccountInformationRequest:
		if err := state.OnAccountInformationRequest(context); err != nil {
			state.logger.Error("error processing OnAccountInformationRequest", log.Error(err))
			panic(err)
		}

	case *messages.AccountMovementRequest:
		if err := state.OnAccountMovementRequest(context); err != nil {
			state.logger.Error("error processing OnAccountMovementRequest", log.Error(err))
			panic(err)
		}

	case *messages.TradeCaptureReportRequest:
		if err := state.OnTradeCaptureReportRequest(context); err != nil {
			state.logger.Error("error processing OnTradeCaptureReportRequest", log.Error(err))
			panic(err)
		}

	case *messages.NewOrderSingleRequest:
		if err := state.OnNewOrderSingleRequest(context); err != nil {
			state.logger.Error("error processing OnNewOrderSingleRequest", log.Error(err))
			panic(err)
		}

	case *xchanger.WebsocketMessage:
		if err := state.onWebsocketMessage(context); err != nil {
			state.logger.Error("error processing onWebsocketMessage", log.Error(err))
			panic(err)
		}

	case *checkSocket:
		if err := state.checkSocket(context); err != nil {
			state.logger.Error("error checking socket", log.Error(err))
			panic(err)
		}

	case *checkAccount:
		if err := state.checkAccount(context); err != nil {
			state.logger.Error("error checking account", log.Error(err))
			panic(err)
		}

	case *refreshMarkPrices:
		if err := state.refreshMarkPrices(context); err != nil {
			state.logger.Error("error refreshing mark prices", log.Error(err))
			panic(err)
		}
	}
}

func (state *AccountListener) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))
	state.okexExecutor = actor.NewPID(context.ActorSystem().Address(), "executor/exchanges/"+constants.OKEX.Name+"_executor")
	if state.client == nil {
		state.client = &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 1024,
				TLSHandshakeTimeout: 10 * time.Second,
			},
			Timeout: 10 * time.Second,
		}
	} else {
		fmt.Println("USING CUSTOM CLIENT")
	}
	// Request securities
	executor := actor.NewPID(context.ActorSystem().Address(), "executor")
	res, err := context.RequestFuture(executor, &messages.SecurityListRequest{}, 10*time.Second).Result()

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
	// TODO filtering should be done by the executor, when specifying exchange in the request
	var filteredSecurities []*models.Security
	for _, s := range securityList.Securities {
		if s.Exchange.ID == state.account.Exchange.ID {
			filteredSecurities = append(filteredSecurities, s)
		}
	}

	/*
		// Then fetch fees
		res, err = context.RequestFuture(state.krakenfExecutor, &messages.AccountInformationRequest{
			Account: state.account.Account,
		}, 10*time.Second).Result()

		if err != nil {
			return fmt.Errorf("error getting account information from executor: %v", err)
		}

		information, ok := res.(*messages.AccountInformationResponse)
		if !ok {
			return fmt.Errorf("was expecting AccountInformationResponse, got %s", reflect.TypeOf(res).String())
		}

		if !information.Success {
			return fmt.Errorf("error fetching account information: %s", information.RejectionReason.String())
		}
	*/

	state.securities = make(map[uint64]*models.Security)
	state.symbolToSec = make(map[string]*models.Security)
	for _, sec := range filteredSecurities {
		state.securities[sec.SecurityID] = sec
		state.symbolToSec[sec.Symbol] = sec
	}
	state.seqNum = 0

	if err := state.subscribeAccount(context); err != nil {
		return fmt.Errorf("error subscribing to account: %v", err)
	}

	checkAccountTicker := time.NewTicker(1 * time.Minute)
	state.checkAccountTicker = checkAccountTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-checkAccountTicker.C:
				context.Send(pid, &checkAccount{})
			case <-time.After(2 * time.Minute):
				if state.checkAccountTicker != checkAccountTicker {
					return
				}
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
				if state.checkSocketTicker != checkSocketTicker {
					return
				}
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
				if state.refreshKeyTicker != refreshKeyTicker {
					return
				}
			}
		}
	}(context.Self())

	refreshMarkPricesTicker := time.NewTicker(5 * time.Second)
	state.refreshMarkPricesTicker = refreshMarkPricesTicker
	go func(pid *actor.PID) {
		for {
			select {
			case <-refreshMarkPricesTicker.C:
				context.Send(pid, &refreshMarkPrices{})
			case <-time.After(31 * time.Second):
				if state.refreshMarkPricesTicker != refreshMarkPricesTicker {
					return
				}
			}
		}
	}(context.Self())

	return nil
}

// TODO
func (state *AccountListener) Clean(context actor.Context) error {
	if state.ws != nil {
		if err := state.ws.Disconnect(); err != nil {
			state.logger.Info("error disconnecting socket", log.Error(err))
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

	if state.refreshMarkPricesTicker != nil {
		state.refreshMarkPricesTicker.Stop()
		state.refreshMarkPricesTicker = nil
	}

	if !state.readOnly {
		for _, sec := range state.securities {
			context.Request(state.okexExecutor, &messages.OrderMassCancelRequest{
				Account: state.account.Account,
				Filter: &messages.OrderFilter{
					Instrument: &models.Instrument{
						SecurityID: &wrapperspb.UInt64Value{Value: sec.SecurityID},
						Symbol:     &wrapperspb.StringValue{Value: sec.Symbol},
						Exchange:   sec.Exchange,
					},
				},
			})
		}
	}

	return nil
}

func (state *AccountListener) OnAccountDataRequest(context actor.Context) error {
	msg := context.Message().(*messages.AccountDataRequest)
	res := &messages.AccountDataResponse{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Securities: state.account.GetSecurities(),
		Orders:     state.account.GetOrders(nil),
		Positions:  state.account.GetPositions(),
		Balances:   state.account.GetBalances(),
		SeqNum:     state.seqNum,
	}

	makerFee := state.account.GetMakerFee()
	takerFee := state.account.GetTakerFee()
	if makerFee != nil {
		res.MakerFee = &wrapperspb.DoubleValue{Value: *makerFee}
	}
	if takerFee != nil {
		res.TakerFee = &wrapperspb.DoubleValue{Value: *takerFee}
	}

	context.Respond(res)

	return nil
}

func (state *AccountListener) OnPositionsRequest(context actor.Context) error {
	msg := context.Message().(*messages.PositionsRequest)
	// TODO FILTER
	positions := state.account.GetPositions()
	context.Respond(&messages.PositionList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Positions:  positions,
	})
	return nil
}

func (state *AccountListener) OnBalancesRequest(context actor.Context) error {
	msg := context.Message().(*messages.BalancesRequest)
	// TODO FILTER
	balances := state.account.GetBalances()
	context.Respond(&messages.BalanceList{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
		Balances:   balances,
	})
	return nil
}

func (state *AccountListener) OnOrderStatusRequest(context actor.Context) error {
	req := context.Message().(*messages.OrderStatusRequest)
	if req.Filter != nil && req.Filter.Instrument != nil {
		req.Filter.Instrument.Symbol.Value = strings.ToUpper(req.Filter.Instrument.Symbol.Value)
	}
	orders := state.account.GetOrders(req.Filter)

	context.Respond(&messages.OrderList{
		RequestID: req.RequestID,
		Success:   true,
		Orders:    orders,
	})
	return nil
}

func (state *AccountListener) OnAccountInformationRequest(context actor.Context) error {
	context.Forward(state.okexExecutor)
	return nil
}

func (state *AccountListener) OnAccountMovementRequest(context actor.Context) error {
	context.Forward(state.okexExecutor)
	return nil
}

func (state *AccountListener) OnTradeCaptureReportRequest(context actor.Context) error {
	context.Forward(state.okexExecutor)
	return nil
}

func (state *AccountListener) OnNewOrderSingleRequest(context actor.Context) error {
	req := context.Message().(*messages.NewOrderSingleRequest)
	req.Account = state.account.Account
	if req.Expire != nil && req.Expire.AsTime().Before(time.Now()) {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       req.RequestID,
			Success:         false,
			RejectionReason: messages.RejectionReason_RequestExpired,
		})
		return nil
	}
	// Check order quantity
	order := &models.Order{
		OrderID:               "",
		ClientOrderID:         req.Order.ClientOrderID,
		Instrument:            req.Order.Instrument,
		OrderStatus:           models.OrderStatus_PendingNew,
		OrderType:             req.Order.OrderType,
		Side:                  req.Order.OrderSide,
		TimeInForce:           req.Order.TimeInForce,
		LeavesQuantity:        req.Order.Quantity,
		Price:                 req.Order.Price,
		CumQuantity:           0,
		ExecutionInstructions: req.Order.ExecutionInstructions,
		Tag:                   req.Order.Tag,
	}
	report, res := state.account.NewOrder(order)
	if res != nil {
		context.Respond(&messages.NewOrderSingleResponse{
			RequestID:       req.RequestID,
			Success:         false,
			RejectionReason: *res,
		})
	} else {
		sender := context.Sender()
		reqResponse := &messages.NewOrderSingleResponse{
			RequestID: req.RequestID,
			Success:   false,
		}

		if report != nil {
			report.Instrument.Symbol.Value = strings.ToLower(report.Instrument.Symbol.Value)

			params, rej := buildPostOrderRequest(report.Instrument.Symbol.Value, req.Order)
			if rej != nil {
				reqResponse.RejectionReason = *rej
				context.Send(sender, reqResponse)
				return nil
			}

			if state.ws != nil {
				_ = state.ws.Disconnect()
			}
			start := time.Now()

			ws := okex.NewWebsocket()
			if err := ws.ConnectPrivate(nil); err != nil {
				return fmt.Errorf("error connecting to the websocket: " + err.Error())
			}

			if err := ws.Login(req.Account.ApiCredentials, req.Account.ApiCredentials.AccountID); err != nil {
				return err
			}

			if !ws.ReadMessage() {
				return ws.Err
			}
			log, ok := ws.Msg.Message.(okex.WSLoginResponse)
			if !ok {
				return fmt.Errorf("error converting message to WSLogin")
			}
			if log.Msg != "" {
				return fmt.Errorf("error login following accounts:" + log.Msg)
			}

			if err := ws.PlaceOrder(params); err != nil {
				return err
			}

			if !ws.ReadMessage() {
				return ws.Err
			}

			orders, ok := ws.Msg.Message.([]okex.WSPlaceOrder)
			if !ok {
				return fmt.Errorf("error converting message to []WSPlaceOrder")
			}
			orderId := orders[0].OrdId

			reqResponse.NetworkRtt = durationpb.New(time.Since(start))
			reqResponse.Success = true
			reqResponse.OrderID = orderId
			fmt.Println("NEW SUCCESS", reqResponse.OrderID)
			context.Send(sender, reqResponse)

		} else {
			if req.ResponseType == messages.ResponseType_Result {
				reqResponse.Success = true
				context.Send(sender, reqResponse)
			}
		}
	}

	return nil
}

func (state *AccountListener) onWebsocketMessage(context actor.Context) error {
	msg := context.Message().(*xchanger.WebsocketMessage)
	if state.ws == nil || msg.WSID != state.ws.ID {
		return nil
	}
	state.lastPingTime = time.Now()

	if msg.Message == nil {
		return fmt.Errorf("received nil message")
	}

	return nil
}

func (state *AccountListener) subscribeAccount(context actor.Context) error {
	if state.ws != nil {
		_ = state.ws.Disconnect()
	}

	ws := okex.NewWebsocket()
	// TODO Dialer
	if err := ws.ConnectPrivate(nil); err != nil {
		return fmt.Errorf("error connecting to krakenf websocket: %v", err)
	}

	if err := ws.Login(state.account.ApiCredentials, state.account.ApiCredentials.AccountID); err != nil {
		return fmt.Errorf("error getting challenge: %v", err)
	}
	if !ws.ReadMessage() {
		return fmt.Errorf("error reading WSInfo response")
	}

	var balanceArgs []map[string]string
	balanceArg := okex.NewBalanceAndPositionRequest("balance_and_position")
	balanceArgs = append(balanceArgs, balanceArg)
	if err := ws.PrivateSubscribe(balanceArgs); err != nil {
		return fmt.Errorf("error subscribing to balance and position: %v", err)
	}

	var positionArgs []map[string]string
	positionArg := okex.NewPositionRequest("positions", "FUTURES")
	positionArgs = append(positionArgs, positionArg)
	if err := ws.PrivateSubscribe(positionArgs); err != nil {
		return fmt.Errorf("error subscribing to open positions: %v", err)
	}

	var orderArgs []map[string]string
	orderArg := okex.NewOrderRequest("orders", "FUTURES")
	orderArg.SetInstFamily("BTC-USD")
	if err := ws.PrivateSubscribe(orderArgs); err != nil {
		return fmt.Errorf("error subscribing to orders: %v", err)
	}

	// Get balances, positions, accounts
	var balances []okex.WSBalanceAndPosition
	var positions []okex.WSPosition
	var orders []okex.WSOrder

	ready := false
	for !ready {
		if !ws.ReadMessage() {
			return fmt.Errorf("error reading ws message")
		}
		fmt.Println("cccccccccccccccc", ws.Msg)
		switch msg := ws.Msg.Message.(type) {
		case []okex.WSBalanceAndPosition:
			balances = msg
		case []okex.WSPosition:
			positions = msg
		case []okex.WSOrder:
			orders = msg
		default:
			context.Send(context.Self(), ws.Msg)
		}
		ready = balances != nil && positions != nil && orders != nil
	}

	morders := make([]*models.Order, len(orders))
	for i, o := range orders {
		mo := WSOrderToModel(&o)
		sec, ok := state.symbolToSec[strings.ToLower(o.InstId)]
		if !ok {
			fmt.Println(state.symbolToSec)
			return fmt.Errorf("got order for unknown symbol %s", o.InstId)
		}
		mo.Instrument.SecurityID = wrapperspb.UInt64(sec.SecurityID)
		morders[i] = mo
	}
	mpositions := make([]*models.Position, len(positions))
	for i, p := range positions {
		mp := WSPositionToModel(&p)
		sec, ok := state.symbolToSec[strings.ToLower(p.InstId)]
		if !ok {
			fmt.Println(state.symbolToSec)
			return fmt.Errorf("got position for unknown symbol %s", p.InstId)
		}
		mp.Instrument.SecurityID = wrapperspb.UInt64(sec.SecurityID)
		mpositions[i] = mp
	}
	var mbalances []*models.Balance
	for _, v := range balances[0].PosData {
		// Get asset from symbol
		asset := SymbolToAsset(v.InstId)
		if asset == nil {
			return fmt.Errorf("unknown asset %s", v.InstId)
		}
		mbalances = append(mbalances, &models.Balance{
			Asset:    asset,
			Quantity: v.Pos,
		})
	}

	msecurities := make([]*models.Security, len(state.securities))
	i := 0
	for _, s := range state.securities {
		msecurities[i] = s
		i += 1
	}

	var mf, tf = 0., 0.
	if err := state.account.Sync(msecurities, morders, mpositions, mbalances, &mf, &tf); err != nil {
		return fmt.Errorf("error syncing account: %v", err)
	}

	go func(ws *okex.Websocket, pid *actor.PID) {
		for ws.ReadMessage() {
			context.Send(pid, ws.Msg)
		}
	}(ws, context.Self())
	state.ws = ws

	return nil
}

func (state *AccountListener) checkSocket(context actor.Context) error {

	if time.Since(state.lastPingTime) > 5*time.Second {
		_ = state.ws.Ping()
		state.lastPingTime = time.Now()
	}

	if state.ws.Err != nil || !state.ws.Connected {
		if state.ws.Err != nil {
			state.logger.Info("error on socket", log.Error(state.ws.Err))
		}
		if err := state.subscribeAccount(context); err != nil {
			return fmt.Errorf("error subscribing to account: %v", err)
		}
	}

	return nil
}

func (state *AccountListener) refreshMarkPrices(context actor.Context) error {
	future := context.RequestFuture(state.okexExecutor, &messages.PositionsRequest{
		RequestID: 0,
		Account:   state.account.Account,
	}, 10*time.Second)
	context.ReenterAfter(future, func(res interface{}, err error) {
		if err != nil {
			state.logger.Error("error updating mark price", log.Error(err))
		}
		pos := res.(*messages.PositionList)
		if !pos.Success {
			state.logger.Error("error updating mark price", log.String("rejection", pos.RejectionReason.String()))
		}
		for _, p := range pos.Positions {
			if p.MarkPrice != nil {
				state.account.UpdateMarkPrice(p.Instrument.SecurityID.Value, p.MarkPrice.Value)
			}
			if p.MaxNotionalValue != nil {
				state.account.UpdateMaxNotionalValue(p.Instrument.SecurityID.Value, p.MaxNotionalValue.Value)
			}
		}
	})
	return nil
}

func (state *AccountListener) checkAccount(context actor.Context) error {
	state.account.CleanOrders()

	pos1 := state.account.GetPositions()
	// Fetch positions
	res, err := context.RequestFuture(state.okexExecutor, &messages.PositionsRequest{
		Instrument: nil,
		Account:    state.account.Account,
	}, 10*time.Second).Result()

	if err != nil {
		return fmt.Errorf("error getting positions from executor: %v", err)
	}

	positionList, ok := res.(*messages.PositionList)
	if !ok {
		return fmt.Errorf("was expecting PositionList, got %s", reflect.TypeOf(res).String())
	}

	if !positionList.Success {
		return fmt.Errorf("error getting positions: %s", positionList.RejectionReason.String())
	}

	var pos2 []*models.Position
	for _, p := range positionList.Positions {
		if p.Quantity != 0 {
			pos2 = append(pos2, p)
		}
	}
	if len(pos1) != len(pos2) {
		return fmt.Errorf("different number of positions: %d %d", len(pos1), len(pos2))
	}

	// sort
	sort.Slice(pos1, func(i, j int) bool {
		return pos1[i].Instrument.SecurityID.Value < pos1[j].Instrument.SecurityID.Value
	})
	sort.Slice(pos2, func(i, j int) bool {
		return pos2[i].Instrument.SecurityID.Value < pos2[j].Instrument.SecurityID.Value
	})

	for i := range pos1 {
		//lp := math.Ceil(1. / state.securities[pos1[i].Instrument.SecurityID.Value].RoundLot.Value)
		diff := math.Abs(pos1[i].Quantity-pos2[i].Quantity) / math.Abs(pos1[i].Quantity+pos2[i].Quantity)
		if diff > 0.01 {
			return fmt.Errorf("different position quantity: %f %f", pos1[i].Cost, pos2[i].Cost)
		}
		diff = math.Abs(pos1[i].Cost-pos2[i].Cost) / math.Abs(pos1[i].Cost+pos2[i].Cost)
		if diff > 0.01 {
			return fmt.Errorf("different position cost: %f %f", pos1[i].Cost, pos2[i].Cost)
		}
	}

	// Update
	for _, p := range pos2 {
		if p.MarkPrice != nil {
			state.account.UpdateMarkPrice(p.Instrument.SecurityID.Value, p.MarkPrice.Value)
		}
		if p.MaxNotionalValue != nil {
			state.account.UpdateMaxNotionalValue(p.Instrument.SecurityID.Value, p.MaxNotionalValue.Value)
		}
	}

	if err := state.account.CheckExpiration(); err != nil {
		return fmt.Errorf("error checking expired orders: %v", err)
	}

	// Fetch balances
	res, err = context.RequestFuture(state.okexExecutor, &messages.BalancesRequest{
		Account: state.account.Account,
	}, 10*time.Second).Result()

	if err != nil {
		return fmt.Errorf("error getting balances from executor: %v", err)
	}

	balanceList, ok := res.(*messages.BalanceList)
	if !ok {
		return fmt.Errorf("was expecting BalanceList, got %s", reflect.TypeOf(res).String())
	}

	if !balanceList.Success {
		return fmt.Errorf("error getting balances: %s", balanceList.RejectionReason.String())
	}

	balanceMap1 := make(map[uint32]float64)
	balanceMap2 := make(map[uint32]float64)
	for _, b := range balanceList.Balances {
		balanceMap1[b.Asset.ID] = b.Quantity
	}
	for _, b := range state.account.GetBalances() {
		balanceMap2[b.Asset.ID] = b.Quantity
	}

	for k, b1 := range balanceMap1 {
		b2 := balanceMap2[k]
		diff := math.Abs(b1-b2) / math.Abs(b1+b2)
		if diff > 0.01 {
			return fmt.Errorf("different margin amount: %f %f", state.account.GetMargin(nil), balanceList.Balances[0].Quantity)
		}
	}

	return nil
}
