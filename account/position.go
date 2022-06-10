package account

import (
	"gitlab.com/alphaticks/alpha-connect/models"
	"math"
)

type Position struct {
	inverse         bool
	tickPrecision   float64
	lotPrecision    float64
	marginPrecision float64
	multiplier      float64
	makerFee        float64
	takerFee        float64
	cost            int64
	rawSize         int64
	cross           bool
}

func NewPosition(inverse bool, tp, lp, mp, mul, maker, taker float64) *Position {
	return &Position{
		inverse:         inverse,
		tickPrecision:   tp,
		lotPrecision:    lp,
		marginPrecision: mp,
		multiplier:      mul,
		makerFee:        maker,
		takerFee:        taker,
		cost:            0,
		rawSize:         0,
		cross:           false,
	}
}

func (pos *Position) UpdateFees(makerFee, takerFee float64) {
	pos.makerFee = makerFee
	pos.takerFee = takerFee
}

func (pos *Position) Buy(price, quantity float64, taker bool) (int64, int64) {
	if pos.inverse {
		price = 1. / price
	}
	rawFillQuantity := int64(math.Round(quantity * pos.lotPrecision))
	rawPrice := int64(math.Round(price * pos.multiplier * pos.marginPrecision))
	var fee int64
	if taker {
		fee = int64(math.Floor(quantity * price * pos.marginPrecision * math.Abs(pos.multiplier) * pos.takerFee))
	} else {
		fee = int64(math.Floor(quantity * price * pos.marginPrecision * math.Abs(pos.multiplier) * math.Abs(pos.makerFee)))
		if pos.makerFee < 0 {
			fee = -fee
		}
	}

	// -10 of pos, 100 of cost, you buy 5 at 5, you close 5 & cost of 25, what's your realized pnl ?
	// - 5 of pos,  75 of cost,
	// Realized = PNL - unrealized PNL at the trade price

	// pnl =

	var realizedCost int64 = 0
	if pos.rawSize < 0 {
		// We are closing our positions from c.size to c.size + size
		closedSize := rawFillQuantity
		// we don't close 'size' if we go over 0 and re-open longs
		if -pos.rawSize < closedSize {
			closedSize = -pos.rawSize
		}

		closedMarginValue := int64(math.Round((float64(closedSize) / pos.lotPrecision) * float64(rawPrice)))
		contractMarginValue := int64(math.Round((float64(-pos.rawSize) / pos.lotPrecision) * float64(rawPrice)))
		unrealizedCost := pos.cost + contractMarginValue
		realizedCost = int64(math.Round((float64(closedSize) / float64(-pos.rawSize)) * float64(unrealizedCost)))

		// Remove closed from cost
		pos.cost += closedMarginValue
		// Remove realized cost
		pos.cost -= realizedCost

		rawFillQuantity -= closedSize
		pos.rawSize += closedSize
	}
	if rawFillQuantity > 0 {
		// We are opening a position
		// math.Ceil doesn't work: -21581.95 -> -21581 but should have been -21582
		// math.Round doesn't work: -21561.017 -> -21561 but should have been -21562
		// math.Floor doesn't work: -21522.73 -> -21523 but should have been -21522
		// Raw opened margin value
		openedMarginValue := int64(math.Round((float64(rawFillQuantity) / pos.lotPrecision) * float64(rawPrice)))
		pos.cost += openedMarginValue
		pos.rawSize += rawFillQuantity
	}

	return fee, realizedCost
}

func (pos *Position) Sell(price, quantity float64, taker bool) (int64, int64) {
	if pos.inverse {
		price = 1. / price
	}

	rawFillQuantity := int64(math.Round(quantity * pos.lotPrecision))
	rawPrice := int64(math.Round(price * pos.multiplier * pos.marginPrecision))
	var fee, realizedCost int64
	if taker {
		fee = int64(math.Floor(quantity * price * pos.marginPrecision * math.Abs(pos.multiplier) * pos.takerFee))
	} else {
		fee = int64(math.Floor(quantity * price * pos.marginPrecision * math.Abs(pos.multiplier) * math.Abs(pos.makerFee)))
		if pos.makerFee < 0 {
			fee = -fee
		}
	}
	if pos.rawSize > 0 {
		// We are closing our position from c.size to c.size + size
		closedSize := rawFillQuantity
		if pos.rawSize < closedSize {
			closedSize = pos.rawSize
		}
		// math.Floor doesn't work
		closedMarginValue := int64(math.Round((float64(closedSize) / pos.lotPrecision) * float64(rawPrice)))    //closedSize * rawPrice
		contractMarginValue := int64(math.Round((float64(pos.rawSize) / pos.lotPrecision) * float64(rawPrice))) //pos.rawSize * rawPrice
		unrealizedCost := pos.cost - contractMarginValue
		realizedCost = int64(math.Round((float64(closedSize) / float64(pos.rawSize)) * float64(unrealizedCost)))

		// Transfer cost
		pos.cost -= closedMarginValue
		// Remove realized cost
		pos.cost -= realizedCost
		pos.rawSize -= closedSize
		rawFillQuantity -= closedSize
	}
	if rawFillQuantity > 0 {
		// We are opening a position
		openedMarginValue := int64(math.Round((float64(rawFillQuantity) / pos.lotPrecision) * float64(rawPrice)))
		pos.cost -= openedMarginValue
		pos.rawSize -= rawFillQuantity
	}

	return fee, realizedCost
}

func (pos *Position) Funding(markPrice, fee float64) int64 {
	size := float64(pos.rawSize) / pos.lotPrecision
	return int64(math.Floor(size * markPrice * pos.marginPrecision * math.Abs(pos.multiplier) * fee))
}

func (pos *Position) GetPosition() *models.Position {
	if pos.rawSize != 0 {
		return &models.Position{
			Quantity: float64(pos.rawSize) / pos.lotPrecision,
			Cost:     float64(pos.cost) / pos.marginPrecision,
			Cross:    pos.cross,
		}
	} else {
		return nil
	}
}

func (pos *Position) Sync(cost, size float64) {
	pos.cost = int64(math.Round(cost * pos.marginPrecision))
	pos.rawSize = int64(math.Round(pos.lotPrecision * size))
}
