package v3_test

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"reflect"
	"testing"
	"time"
)

func TestPoolData(t *testing.T) {
	as, executor, cleaner := tests.StartExecutor(t, &constants.UNISWAPV3, nil)
	defer cleaner()

	securityID := []uint64{
		3923210566889873515,
	}

	res, err := as.Root.RequestFuture(executor, &messages.SecurityListRequest{}, 30*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	securityList, ok := res.(*messages.SecurityList)
	if !ok {
		t.Fatalf("was expecting *messages.SecurityList, got %s", reflect.TypeOf(res).String())
	}
	if !securityList.Success {
		t.Fatal(securityList.RejectionReason.String())
	}
	fmt.Println("GOT SEC", len(securityList.Securities))

	var sec *models.Security
	for _, s := range securityList.Securities {
		for _, secID := range securityID {
			if secID == s.SecurityID {
				sec = s
			}
		}
	}
	if sec == nil {
		t.Fatal("security not found")
	}

	res, err = as.Root.RequestFuture(executor, &messages.HistoricalUnipoolV3DataRequest{
		RequestID: uint64(time.Now().UnixNano()),
		Instrument: &models.Instrument{
			SecurityID: &types.UInt64Value{Value: sec.SecurityID},
			Exchange:   sec.Exchange,
			Symbol:     &types.StringValue{Value: sec.Symbol},
		},
		Start: 14038263,
		End:   14040263,
	}, 50*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	updates, ok := res.(*messages.HistoricalUnipoolV3DataResponse)
	if !ok {
		t.Fatalf("was expecting *messages.HistoricalUnipoolV3DataResponse, got %s", reflect.TypeOf(res).String())
	}
	if !updates.Success {
		t.Fatal(updates.RejectionReason.String())
	}

	fmt.Println(updates.Events)
}
