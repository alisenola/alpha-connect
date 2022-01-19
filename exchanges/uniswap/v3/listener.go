package v3

import (
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/gogo/protobuf/types"
	extypes "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/go-graphql-client"
	"gitlab.com/alphaticks/gorderbook"
	"gitlab.com/alphaticks/xchanger/constants"
	uniswap "gitlab.com/alphaticks/xchanger/exchanges/uniswap/V3"
	xchangerUtils "gitlab.com/alphaticks/xchanger/utils"
)

type checkSockets struct{}
type postAggTrade struct{}

type InstrumentData struct {
	orderBook      []*gorderbook.UnipoolV3
	seqNum         uint64
	lastUpdateTime uint64
	lastHBTime     time.Time
	lastAggTradeTs uint64
}

type Listener struct {
	extypes.Listener
	poolWs         *uniswap.Websocket
	security       *models.Security
	dialerPool     *xchangerUtils.DialerPool
	instrumentData *InstrumentData
	logger         *log.Logger
	lastPingTime   time.Time
	uniExectutor   *actor.PID
	socketTicker   *time.Ticker
}

func NewListenerProducer(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Producer {
	return func() actor.Actor {
		return NewListener(security, dialerPool)
	}
}

func NewListener(security *models.Security, dialerPool *xchangerUtils.DialerPool) actor.Actor {
	return &Listener{
		poolWs:         nil,
		security:       security,
		dialerPool:     dialerPool,
		instrumentData: nil,
		logger:         nil,
		socketTicker:   nil,
	}
}

func (state *Listener) Receive(context actor.Context) {
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

	case *messages.UnipoolV3DataRequest:
		if err := state.OnUnipoolV3DataRequest(context); err != nil {
			state.logger.Error("error processing OnUnipoolV3DataRequest", log.Error(err))
			panic(err)
		}

	case *checkSockets:
		if err := state.checkSockets(context); err != nil {
			state.logger.Error("error checking socket", log.Error(err))
			panic(err)
		}
	}
}

func (state *Listener) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()),
		log.String("exchange", state.security.Exchange.Name),
		log.String("symbol", state.security.Symbol))

	state.lastPingTime = time.Now()
	state.uniExectutor = actor.NewPID(context.ActorSystem().Address(), "executor/"+constants.UNISWAPV3.Name+"_executor")

	state.instrumentData = &InstrumentData{
		orderBook:      nil,
		seqNum:         uint64(time.Now().UnixNano()),
		lastUpdateTime: 0,
		lastHBTime:     time.Now(),
		lastAggTradeTs: 0,
	}

	socketTicker := time.NewTicker(5 * time.Second)
	state.socketTicker = socketTicker
	go func(pid *actor.PID) {
		for {
			select {
			case _ = <-socketTicker.C:
				context.Send(pid, &checkSockets{})
			case <-time.After(10 * time.Second):
				// timer stopped, we leave
				return
			}
		}
	}(context.Self())

	return nil
}

func (state *Listener) OnUnipoolV3DataRequest(context actor.Context) error {
	if state.poolWs != nil {
		_ = state.poolWs.Disconnect()
	}
	ws := uniswap.NewWebsocket()

	ws.NewSubscriptionClient()
	if err := ws.Connect(); err != nil {
		return fmt.Errorf("error connecting to websocket: %v", err)
	}

	fmt.Println("LISTENER UNIPOOL DATA REQUEST")
	time.Sleep(10 * time.Second)
	future := context.RequestFuture(
		state.uniExectutor,
		&messages.UnipoolV3DataRequest{
			RequestID: uint64(time.Now().UnixNano()),
			Subscribe: false,
			Instrument: &models.Instrument{
				SecurityID: &types.UInt64Value{Value: state.security.SecurityID},
				Symbol:     &types.StringValue{Value: state.security.Symbol},
				Exchange:   state.security.Exchange,
			},
		},
		60*time.Second)
	res, err := future.Result()
	if err != nil {
		return fmt.Errorf("error getting pool snapshot %v", err)
	}
	msg, ok := res.(*messages.UnipoolV3DataResponse)
	if !ok {
		return fmt.Errorf("was expecting UnipoolV3DataResponse, got %s", reflect.TypeOf(msg).String())
	}
	if !msg.Success {
		return fmt.Errorf("error fetching the pool snapshot %s", msg.RejectionReason.String())
	}
	if msg.Snapshot == nil {
		return fmt.Errorf("pool has empty snapshot")
	}

	if _, err := ws.SubscribeTransactions(graphql.Int(msg.Snapshot.Timestamp.Seconds), graphql.ID(state.security.Symbol)); err != nil {
		return fmt.Errorf("error subscribing to the pool transactions: %v", err)
	}
	state.poolWs = ws

	unipoolV3 := gorderbook.NewUnipoolV3(
		int32(state.security.TakerFee.Value),
	)
	fmt.Println("INITIALIZING")
	unipoolV3.Initialize(big.NewInt(1).SetBytes(msg.Snapshot.SqrtPrice), msg.Snapshot.Tick)
	fmt.Println("SYNCING")
	unipoolV3.Sync(
		msg.Snapshot.Ticks,
		msg.Snapshot.Liquidity,
		msg.Snapshot.ProtocolFees_0,
		msg.Snapshot.ProtocolFees_1,
		msg.Snapshot.FeeGrowthGlobal_0X128,
		msg.Snapshot.FeeGrowthGlobal_1X128,
		msg.Snapshot.TotalValueLockedToken_0,
		msg.Snapshot.TotalValueLockedToken_1,
		msg.Snapshot.FeeTier,
		msg.Snapshot.Positions,
	)
	//TODO Make function in listener_utils in order to compute all the transactions inside transactions.Transactions
	//TODO Check for tokensOwed0 and tokensOwed1 from the position structure in uniswap contract in theGraph
	sync := false
	for !sync {
		if !state.poolWs.ReadMessage() {
			return fmt.Errorf("error reading the message %v", err)
		}
		transactions, ok := state.poolWs.Msg.Message.(*uniswap.Transactions)
		if !ok {
			return fmt.Errorf("incorrect message type %v", err)
		}
		if len(transactions.Transactions) == 0 {
			state.instrumentData.lastUpdateTime = uint64(ws.Msg.ClientTime.UnixNano() / 1000)
			sync = true
		}
		fmt.Println("PROCESSING", transactions.Transactions, len(transactions.Transactions))
		fmt.Printf("INSIDE %+v \n", transactions.Transactions[999])
		processTransactions(transactions, unipoolV3)
	}
	go func(ws *uniswap.Websocket, actor *actor.PID) {
		for ws.ReadMessage() {
			context.Send(actor, ws.Msg)
		}
	}(state.poolWs, context.Self())
	return nil
}

func (state *Listener) Clean(context actor.Context) error {
	if state.socketTicker != nil {
		state.socketTicker.Stop()
		state.socketTicker = nil
	}

	return nil
}

func (state *Listener) checkSockets(context actor.Context) error {

	return nil
}
