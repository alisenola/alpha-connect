package types

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"time"
)

type Reconcile interface {
	actor.Actor
	GetLogger() *log.Logger
	Initialize(context actor.Context) error
	Clean(context actor.Context) error
	OnAccountMovementRequest(context actor.Context) error
	OnTradeCaptureReportRequest(context actor.Context) error
}

type BaseReconcile struct {
}

func (state *BaseReconcile) GetLogger() *log.Logger {
	panic("not implemented")
}

func (state *BaseReconcile) Initialize(context actor.Context) error {
	return nil
}

func (state *BaseReconcile) Clean(context actor.Context) error {
	return nil
}

func (state *BaseReconcile) GetTransactions() *mongo.Collection {
	return nil
}

func (state *BaseReconcile) OnAccountMovementRequest(context actor.Context) error {
	req := context.Message().(*messages.AccountMovementRequest)
	context.Respond(&messages.AccountMovementResponse{
		RequestID:       req.RequestID,
		ResponseID:      rand.Uint64(),
		Success:         false,
		RejectionReason: messages.RejectionReason_UnsupportedRequest,
	})
	return nil
}

func (state *BaseReconcile) OnTradeCaptureReportRequest(context actor.Context) error {
	msg := context.Message().(*messages.TradeCaptureReportRequest)
	res := &messages.TradeCaptureReport{
		RequestID:  msg.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    true,
	}
	defer context.Respond(res)
	// TODO
	return nil
	/*
		defer context.Respond(res)
		filter := bson.D{
			{Key: "account", Value: msg.Account.Name},
			{Key: "type", Value: "TRADE"},
		}
		if msg.Filter.From != nil {
			from := msg.Filter.From.AsTime()
			filter = append(filter, bson.E{Key: "time", Value: bson.D{{Key: "$gt", Value: from}}})
		}
		if msg.Filter.To != nil {
			to := msg.Filter.To.AsTime()
			filter = append(filter, bson.E{Key: "time", Value: bson.D{{Key: "$lt", Value: to}}})
		}
		cur, err := state.GetTransactions().Find(goContext.Background(), filter)
		if err != nil {
			res.Success = false
			res.RejectionReason = messages.RejectionReason_Other
			state.GetLogger().Error("error fetching trades", log.Error(err))
			return nil
		}
		defer cur.Close(goContext.Background())
		for cur.Next(goContext.Background()) {
			var tx Transaction
			if err := cur.Decode(tx); err != nil {
				state.GetLogger().Error("error decoding tx", log.Error(err))
				res.Success = false
				res.RejectionReason = messages.RejectionReason_Other
				return nil
			}
			side := models.Side_Buy
			if tx.Fill.Quantity < 0 {
				side = models.Side_Sell
			}
			secID, _ := strconv.ParseUint(tx.Fill.SecurityID, 10, 64)
			ts := timestamppb.New(tx.Time)
			var commission float64
			var commissionAsset *xchangerModels.Asset
			for _, m := range tx.Movements {
				if m.Reason == int32(messages.AccountMovementType_Commission) {
					commission = m.Quantity
					commissionAsset, _ = constants.GetAssetByID(m.AssetID)
				}
			}
			res.Trades = append(res.Trades, &models.TradeCapture{
				Side:            side,
				Type:            models.TradeType_Regular,
				Price:           tx.Fill.Price,
				Quantity:        math.Abs(tx.Fill.Quantity),
				Commission:      commission,
				CommissionAsset: commissionAsset,
				TradeID:         tx.ID,
				Instrument:      &models.Instrument{SecurityID: &wrapperspb.UInt64Value{Value: secID}},
				Trade_LinkID:    nil,
				OrderID:         nil,
				ClientOrderID:   nil,
				TransactionTime: ts,
			})
		}
		return nil
	*/
}

func ReconcileReceive(state Reconcile, context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		if err := state.Initialize(context); err != nil {
			state.GetLogger().Error("error initializing", log.Error(err))
			panic(err)
		}
		state.GetLogger().Info("actor started")

	case *actor.Stopping:
		if err := state.Clean(context); err != nil {
			state.GetLogger().Error("error stopping", log.Error(err))
			panic(err)
		}
		state.GetLogger().Info("actor stopping")

	case *actor.Stopped:
		state.GetLogger().Info("actor stopped")

	case *actor.Restarting:
		if err := state.Clean(context); err != nil {
			state.GetLogger().Error("error restarting", log.Error(err))
			// Attention, no panic in restarting or infinite loop
		}
		state.GetLogger().Info("actor restarting")

	case *messages.AccountMovementRequest:
		if err := state.OnAccountMovementRequest(context); err != nil {
			state.GetLogger().Error("error processing OnAccountMovementRequest", log.Error(err))
			panic(err)
		}

	case *messages.TradeCaptureReportRequest:
		if err := state.OnTradeCaptureReportRequest(context); err != nil {
			state.GetLogger().Error("error processing OnTradeCaptureReportRequest", log.Error(err))
			panic(err)
		}
	}
}
