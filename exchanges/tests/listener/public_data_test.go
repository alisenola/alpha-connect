package listener

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/alpha-connect/models/messages"
	"gitlab.com/alphaticks/xchanger/constants"
	"reflect"
	"testing"
	"time"
)

func TestBybitiLiquidations(t *testing.T) {
	t.Parallel()
	as, executor, clean := StartExecutor(&constants.BYBITI)
	defer clean()

	res, err := as.Root.RequestFuture(executor, &messages.HistoricalLiquidationsRequest{Instrument: &models.Instrument{
		SecurityID: &types.UInt64Value{Value: 7374647908427501521},
	}}, 10*time.Second).Result()
	if err != nil {
		t.Fatal(err)
	}
	liquidations, ok := res.(*messages.HistoricalLiquidationsResponse)
	if !ok {
		t.Fatalf("was expecting *messages.SecurityList, got %s", reflect.TypeOf(res).String())
	}
	if !liquidations.Success {
		t.Fatal(liquidations.RejectionReason.String())
	}
	fmt.Println(liquidations.Liquidations)
}
