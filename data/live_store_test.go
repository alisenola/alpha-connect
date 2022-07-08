package data

import (
	"gitlab.com/alphaticks/alpha-connect/config"
	"gitlab.com/alphaticks/alpha-connect/tests"
	types "gitlab.com/alphaticks/tickstore-types"
	"gitlab.com/alphaticks/xchanger/constants"
	"testing"
)

func TestLiveQuery(t *testing.T) {
	cfg := config.Config{Exchanges: []string{constants.BINANCE.Name}}
	as, executor, clean := tests.StartExecutor(t, &cfg)
	defer clean()
	lt, err := NewLiveStore(as, executor)
	if err != nil {
		t.Fatal(err)
	}
	qs := types.NewQuerySettings(types.WithSelector(`SELECT OBPrice(orderbook, "0.2") WHERE exchange="^binance$" base="^BTC$" quote="^USDT$" GROUPBY base`))
	q, err := lt.NewQuery(qs)
	if err != nil {
		t.Fatal(err)
	}
	i := 0
	for i < 100 {
		for q.Next() {
			q.Read()
			i += 1
		}
		if q.Err() != nil {
			t.Fatal(q.Err())
		}
	}

	if err := q.Close(); err != nil {
		t.Fatal(err)
	}
}
