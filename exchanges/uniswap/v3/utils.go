package v3

import (
	"math/big"

	"gitlab.com/alphaticks/alpha-connect/models"
	"gitlab.com/alphaticks/gorderbook"
)

const (
	ETH_CLIENT_URL = "wss://eth-mainnet.alchemyapi.io/v2/hdodrT8DMC-Ow9rd6qIOcjpZgr5_Ixdg"
)

func ProcessUpdate(pool *gorderbook.UnipoolV3, update *models.UPV3Update) error {
	if i := update.Initialize; i != nil {
		sqrt := big.NewInt(1).SetBytes(i.SqrtPriceX96)
		pool.Initialize(
			sqrt,
			i.Tick,
		)
	}
	if m := update.Mint; m != nil {
		var own [20]byte
		copy(own[:], m.Owner)
		amount := big.NewInt(1).SetBytes(m.Amount)
		amount0 := big.NewInt(1).SetBytes(m.Amount0)
		amount1 := big.NewInt(1).SetBytes(m.Amount1)
		pool.Mint(
			own,
			m.TickLower,
			m.TickUpper,
			amount,
			amount0,
			amount1,
		)
	}
	if b := update.Burn; b != nil {
		var own [20]byte
		copy(own[:], b.Owner)
		amount := big.NewInt(1).SetBytes(b.Amount)
		amount0 := big.NewInt(1).SetBytes(b.Amount0)
		amount1 := big.NewInt(1).SetBytes(b.Amount1)
		pool.Burn(
			own,
			b.TickLower,
			b.TickUpper,
			amount,
			amount0,
			amount1,
		)
	}
	if s := update.Swap; s != nil {
		amount0 := big.NewInt(1).SetBytes(s.Amount0)
		amount1 := big.NewInt(1).SetBytes(s.Amount1)
		sqrt := big.NewInt(1).SetBytes(s.SqrtPriceX96)
		sqrtPrev := pool.GetSqrtPrice()
		if sqrt.Cmp(sqrtPrev) > 0 {
			sqrt.Add(sqrt, big.NewInt(1))
			amount0.Neg(amount0)
		} else {
			sqrt.Sub(sqrt, big.NewInt(1))
			amount1.Neg(amount1)
		}
		pool.Swap(
			amount0,
			amount1,
			sqrt,
			s.Tick,
		)
	}
	if c := update.Collect; c != nil {
		var own [20]byte
		copy(own[:], c.Owner)
		amount0 := big.NewInt(1).SetBytes(c.AmountRequested0)
		amount1 := big.NewInt(1).SetBytes(c.AmountRequested1)
		pool.Collect(
			own,
			c.TickLower,
			c.TickUpper,
			amount0,
			amount1,
		)
	}
	if f := update.Flash; f != nil {
		amount0 := big.NewInt(1).SetBytes(f.Amount0)
		amount1 := big.NewInt(1).SetBytes(f.Amount1)
		pool.Flash(
			amount0,
			amount1,
		)
	}
	return nil
}
