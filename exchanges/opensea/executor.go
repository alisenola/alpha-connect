package opensea

import (
	goContext "context"
	"encoding/json"
	"fmt"
	"gitlab.com/alphaticks/alpha-connect/utils"
	gorderbook "gitlab.com/alphaticks/gorderbook/gorderbook.models"
	"gitlab.com/alphaticks/xchanger/constants"
	"gitlab.com/alphaticks/xchanger/exchanges"
	"math/big"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"time"

	registry "gitlab.com/alphaticks/alpha-public-registry-grpc"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	models2 "gitlab.com/alphaticks/xchanger/models"
	xutils "gitlab.com/alphaticks/xchanger/utils"

	extype "gitlab.com/alphaticks/alpha-connect/exchanges/types"
	"gitlab.com/alphaticks/alpha-connect/jobs"
	opensea "gitlab.com/alphaticks/xchanger/exchanges/opensea"
)

type QueryRunner struct {
	pid       *actor.PID
	rateLimit *exchanges.RateLimit
}

type Executor struct {
	extype.BaseExecutor
	queryRunners             []*QueryRunner
	marketableProtocolAssets map[uint64]*models.MarketableProtocolAsset
	credentials              *models2.APICredentials
	dialerPool               *xutils.DialerPool
	logger                   *log.Logger
	registry                 registry.PublicRegistryClient
}

func NewExecutor(registry registry.PublicRegistryClient, dialerPool *xutils.DialerPool, credentials *models2.APICredentials) actor.Actor {
	return &Executor{
		dialerPool:  dialerPool,
		registry:    registry,
		credentials: credentials,
	}
}

func (state *Executor) getQueryRunner() *QueryRunner {
	sort.Slice(state.queryRunners, func(i, j int) bool {
		return rand.Uint64()%2 == 0
	})

	for _, q := range state.queryRunners {
		if !q.rateLimit.IsRateLimited() {
			return q
		}
	}
	return nil
}

func (state *Executor) Receive(context actor.Context) {
	extype.ReceiveExecutor(state, context)
}

func (state *Executor) Initialize(context actor.Context) error {
	state.logger = log.New(
		log.InfoLevel,
		"",
		log.String("ID", context.Self().Id),
		log.String("type", reflect.TypeOf(*state).String()))

	dialers := state.dialerPool.GetDialers()
	for _, dialer := range dialers {
		client := &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 1024,
				TLSHandshakeTimeout: 10 * time.Second,
				DialContext:         dialer.DialContext,
			},
			Timeout: 10 * time.Second,
		}
		props := actor.PropsFromProducer(func() actor.Actor {
			return jobs.NewHTTPQuery(client)
		})
		state.queryRunners = append(state.queryRunners, &QueryRunner{
			pid:       context.Spawn(props),
			rateLimit: exchanges.NewRateLimit(100, time.Minute),
		})
	}
	return state.UpdateMarketableProtocolAssetList(context)
}

func (state *Executor) UpdateMarketableProtocolAssetList(context actor.Context) error {
	assets := make([]*models.MarketableProtocolAsset, 0)
	reg := state.registry
	ctx, cancel := goContext.WithTimeout(goContext.Background(), 10*time.Second)
	defer cancel()
	//TODO add matic and solana protocols
	filter := registry.MarketableProtocolAssetFilter{
		MarketId: []uint32{constants.OPENSEA.ID},
	}
	in := registry.MarketableProtocolAssetsRequest{
		Filter: &filter,
	}

	//Get marketable protocol assets
	res, err := reg.MarketableProtocolAssets(ctx, &in)
	if err != nil {
		return fmt.Errorf("error querying marketable protocol asset list: %v", err)
	}
	marketables := res.MarketableProtocolAssets
	var protocolAssetIds []uint64
	for _, marketable := range marketables {
		protocolAssetIds = append(protocolAssetIds, marketable.ProtocolAssetId)
	}

	//Get protocol assets from the corresponding marketable protocol assets
	r, err := reg.ProtocolAssets(ctx, &registry.ProtocolAssetsRequest{
		Filter: &registry.ProtocolAssetFilter{
			ProtocolAssetId: protocolAssetIds,
		},
	})
	if err != nil {
		return fmt.Errorf("error querying protocol asset list: %v", err)
	}
	protocolAssets := r.ProtocolAssets
	idToProtocolAsset := make(map[uint64]*registry.ProtocolAsset)
	for _, asset := range protocolAssets {
		idToProtocolAsset[asset.ProtocolAssetId] = asset
	}
	//Create the marketable assets from the protocol assets and the marketable protocol assets
	for _, marketable := range marketables {
		if addr, ok := idToProtocolAsset[marketable.ProtocolAssetId].Meta["address"]; !ok || len(addr) < 2 {
			state.logger.Warn("invalid protocol asset address")
			continue
		}
		_, ok := big.NewInt(1).SetString(idToProtocolAsset[marketable.ProtocolAssetId].Meta["address"][2:], 16)
		if !ok {
			state.logger.Warn("invalid protocol asset address", log.Error(err))
			continue
		}
		as, ok := constants.GetAssetByID(idToProtocolAsset[marketable.ProtocolAssetId].AssetId)
		if !ok {
			state.logger.Warn(fmt.Sprintf("error getting asset with id %d", idToProtocolAsset[marketable.ProtocolAssetId].AssetId))
			continue
		}
		ch, ok := constants.GetChainByID(idToProtocolAsset[marketable.ProtocolAssetId].ChainId)
		if !ok {
			state.logger.Warn(fmt.Sprintf("error getting chain with id %d", idToProtocolAsset[marketable.ProtocolAssetId].ChainId))
			continue
		}
		assets = append(
			assets,
			&models.MarketableProtocolAsset{
				MarketableProtocolAssetID: marketable.MarketableProtocolAssetId,
				ProtocolAsset: &models.ProtocolAsset{
					ProtocolAssetID: idToProtocolAsset[marketable.ProtocolAssetId].ProtocolAssetId,
					Protocol: &models2.Protocol{
						ID:   constants.ERC721.ID,
						Name: "ERC-721",
					},
					Asset: &models2.Asset{
						Name:   as.Name,
						Symbol: as.Symbol,
						ID:     as.ID,
					},
					Chain: &models2.Chain{
						ID:   ch.ID,
						Name: ch.Name,
						Type: ch.Type,
					},
					Meta: idToProtocolAsset[marketable.ProtocolAssetId].Meta,
				},
				Market: &constants.OPENSEA,
			},
		)
	}
	state.marketableProtocolAssets = make(map[uint64]*models.MarketableProtocolAsset)
	for _, a := range assets {
		state.marketableProtocolAssets[a.MarketableProtocolAssetID] = a
	}
	context.Send(context.Parent(), &messages.MarketableProtocolAssetList{
		ResponseID:               uint64(time.Now().UnixNano()),
		MarketableProtocolAssets: assets,
		Success:                  true,
	})

	return nil
}

func (state *Executor) OnMarketableProtocolAssetListRequest(context actor.Context) error {
	req := context.Message().(*messages.MarketableProtocolAssetListRequest)
	passets := make([]*models.MarketableProtocolAsset, len(state.marketableProtocolAssets))
	i := 0
	for _, v := range state.marketableProtocolAssets {
		passets[i] = v
		i += 1
	}
	context.Respond(&messages.MarketableProtocolAssetList{
		RequestID:                req.RequestID,
		ResponseID:               uint64(time.Now().UnixNano()),
		Success:                  true,
		MarketableProtocolAssets: passets,
	})
	return nil
}

func (state *Executor) OnHistoricalSalesRequest(context actor.Context) error {
	req := context.Message().(*messages.HistoricalSalesRequest)
	msg := &messages.HistoricalSalesResponse{
		RequestID:  req.RequestID,
		ResponseID: uint64(time.Now().UnixNano()),
		Success:    false,
	}
	if state.credentials.APIKey == "" {
		msg.RejectionReason = messages.UnsupportedRequest
		context.Respond(msg)
		return nil
	}
	pAsset := state.marketableProtocolAssets[req.MarketableProtocolAssetID]
	if pAsset == nil {
		msg.RejectionReason = messages.UnknownProtocolAsset
		context.Respond(msg)
		return nil
	}

	qr := state.getQueryRunner()
	if qr == nil {
		msg.RejectionReason = messages.RateLimitExceeded
		context.Respond(msg)
		return nil
	}

	params := opensea.NewGetEventsParams()
	add := pAsset.ProtocolAsset.Meta["address"]
	params.SetAssetContractAddress(add)
	params.SetEventType("successful")
	if req.To != nil {
		params.SetOccurredBefore(uint64(req.To.Seconds))
	}
	if req.From != nil {
		params.SetOccurredAfter(uint64(req.From.Seconds))
	}
	r, weight, err := opensea.GetEvents(params, state.credentials.APIKey)
	if err != nil {
		msg.RejectionReason = messages.UnsupportedOrderCharacteristic
		context.Respond(msg)
		return nil
	}
	qr.rateLimit.Request(weight)

	//Global variables
	cursor := ""
	done := false
	var sales []*models.Sale
	sender := context.Sender() //keep copy of sender

	var processFuture func(res interface{}, err error)
	processFuture = func(res interface{}, err error) {
		if err != nil {
			msg.RejectionReason = messages.HTTPError
			context.Respond(msg)
			return
		}

		resp := res.(*jobs.PerformQueryResponse)
		if resp.StatusCode != 200 {
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				err := fmt.Errorf("%d %s", resp.StatusCode, string(resp.Response))
				state.logger.Warn("http client error", log.Error(err))
				msg.RejectionReason = messages.HTTPError
				context.Respond(msg)
				return
			} else if resp.StatusCode >= 500 {
				err := fmt.Errorf("%d %s", resp.StatusCode, string(resp.Response))
				state.logger.Warn("http server error", log.Error(err))
				msg.RejectionReason = messages.HTTPError
				context.Respond(msg)
				return
			}
			return
		}

		var events opensea.EventsResponse
		if err := json.Unmarshal(resp.Response, &events); err != nil {
			msg.RejectionReason = messages.ExchangeAPIError
			context.Respond(msg)
			return
		}
		for _, e := range events.AssetEvents {
			var from [20]byte
			var to [20]byte
			var tokenID [32]byte
			var price [32]byte
			f, ok := big.NewInt(1).SetString(e.Transaction.FromAccount.Address[2:], 16)
			if !ok {
				state.logger.Warn("incorrect address format", log.String("address", e.Transaction.FromAccount.Address))
				continue
			}
			t, ok := big.NewInt(1).SetString(e.Transaction.ToAccount.Address[2:], 16)
			if !ok {
				state.logger.Warn("incorrect address format", log.String("address", e.Transaction.ToAccount.Address))
				continue
			}
			token, ok := big.NewInt(1).SetString(e.Asset.TokenId, 10)
			if !ok {
				state.logger.Warn("incorrect tokenID format", log.String("tokenID", e.Asset.TokenId))
				continue
			}
			i, err := strconv.ParseInt(e.Transaction.BlockNumber, 10, 64)
			if err != nil {
				state.logger.Warn("incorrect block number format", log.String("block number", e.Transaction.BlockNumber))
				continue
			}
			p, ok := big.NewInt(1).SetString(e.TotalPrice, 10)
			if !ok {
				state.logger.Warn("incorrect price format", log.String("price", e.TotalPrice))
				continue
			}
			tim, err := time.Parse("2006-01-02T15:04:05", e.Transaction.Timestamp)
			if err != nil {
				state.logger.Warn("incorrect timestamp format", log.String("ts", e.Transaction.Timestamp))
				continue
			}
			f.FillBytes(from[:])
			t.FillBytes(to[:])
			token.FillBytes(tokenID[:])
			p.FillBytes(price[:])
			ts := utils.NanoToTimestamp(uint64(tim.UnixNano()))
			sales = append(sales, &models.Sale{
				Transfer: &gorderbook.AssetTransfer{
					From:    from[:],
					To:      to[:],
					TokenId: tokenID[:],
				},
				Block:     uint64(i),
				Price:     price[:],
				Timestamp: ts,
				Id:        uint64(e.Transaction.Id),
			})
		}
		cursor = events.Next
		done = cursor == ""
		if !done {
			params.SetCursor(cursor)
			r, weight, err = opensea.GetEvents(params, state.credentials.APIKey)
			if err != nil {
				msg.RejectionReason = messages.UnsupportedOrderCharacteristic
				context.Respond(msg)
				return
			}
			go func() {
				qr.rateLimit.WaitRequest(weight)
				fut := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: r}, 15*time.Second)
				qr.rateLimit.Request(weight)
				context.AwaitFuture(fut, processFuture)
				return
			}()
		} else {
			msg.Success = true
			msg.Sale = sales
			msg.SeqNum = uint64(time.Now().UnixNano())
			context.Send(sender, msg)
		}
	}
	future := context.RequestFuture(qr.pid, &jobs.PerformHTTPQueryRequest{Request: r}, 10*time.Minute)
	context.AwaitFuture(future, processFuture)
	return nil
}

func (state *Executor) Clean(context actor.Context) error {
	return nil
}

func (state *Executor) GetLogger() *log.Logger {
	return state.logger
}
