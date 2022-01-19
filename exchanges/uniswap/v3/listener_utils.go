package v3

import (
	"math/big"

	"gitlab.com/alphaticks/gorderbook"
	uniswap "gitlab.com/alphaticks/xchanger/exchanges/uniswap/V3"
)

func processTransactions(transactions *uniswap.Transactions, unipoolV3 *gorderbook.UnipoolV3) {
	for _, t := range transactions.Transactions {
		orderedT := t.OrderTransaction()
		for _, oT := range orderedT.Transaction {
			if m, ok := oT.(*uniswap.Mint); ok {
				owner, ok := big.NewInt(1).SetString(m.Owner[2:], 16)
				if !ok {
					continue
				}
				var ownerID [20]byte
				copy(ownerID[:], owner.Bytes())
				unipoolV3.Mint(ownerID, m.TickLower, m.TickUpper, m.Amount, m.Amount0, m.Amount1)
			}
			if b, ok := oT.(*uniswap.Burn); ok {
				owner, ok := big.NewInt(1).SetString(b.Owner[2:], 16)
				if !ok {
					continue
				}
				var ownerID [20]byte
				copy(ownerID[:], owner.Bytes())
				unipoolV3.Burn(ownerID, b.TickLower, b.TickUpper, b.Amount, b.Amount0, b.Amount1)
			}
			if s, ok := oT.(*uniswap.Swap); ok {
				limit := s.SqrtPriceX96
				if s.Amount0.Cmp(big.NewInt(0)) > 0 {
					limit.Add(limit, big.NewInt(-1))
				} else {
					limit.Add(limit, big.NewInt(1))
				}
				unipoolV3.Swap(s.Amount0, s.Amount1, limit, s.Tick)
			}
			if c, ok := oT.(*uniswap.Collect); ok {
				owner, ok := big.NewInt(1).SetString(c.Owner[2:], 16)
				if !ok {
					continue
				}
				var ownerID [20]byte
				copy(ownerID[:], owner.Bytes())
				unipoolV3.Collect(ownerID, c.TickLower, c.TickUpper, c.Amount0, c.Amount1)
			}
			if f, ok := oT.(*uniswap.Flash); ok {
				unipoolV3.Flash(f.Amount0, f.Amount1)
			}
		}
	}
}
