package fbinance

import (
	"gitlab.com/alphaticks/alpha-connect/exchanges/tests"
	"gitlab.com/alphaticks/xchanger/exchanges/fbinance"
	"testing"
)

func TestAccountListener(t *testing.T) {
	fbinance.EnableTestNet()
	tests.MarketFill(t, tests.AccntTest{
		// TODO
	})
}
