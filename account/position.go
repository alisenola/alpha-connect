package account

import (
	"gitlab.com/alphaticks/alphac/models"
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

func (pos *Position) UpdateFees(makerFee, takerFee float64) {
	pos.makerFee = makerFee
	pos.takerFee = takerFee
}

func (pos *Position) Buy(price, quantity float64, taker bool) (int64, int64) {
	if pos.inverse {
		price = 1. / price
	}
	rawFillQuantity := int64(quantity * pos.lotPrecision)
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
	var realizedCost int64 = 0
	if pos.rawSize < 0 {
		// We are closing our positions from c.size to c.size + size
		closedSize := rawFillQuantity
		// we don't close 'size' if we go over 0 and re-open longs
		if -pos.rawSize < closedSize {
			closedSize = -pos.rawSize
		}

		contractMarginValue := int64(math.Round(float64(-pos.rawSize*rawPrice) / pos.lotPrecision))
		closedMarginValue := int64(math.Round(float64(closedSize*rawPrice) / pos.lotPrecision))

		unrealizedCost := pos.cost + contractMarginValue
		realizedCost = int64(math.Round(float64(closedSize*unrealizedCost) / float64(-pos.rawSize)))

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
		openedMarginValue := int64(math.Round(float64(rawFillQuantity*rawPrice) / pos.lotPrecision))
		pos.cost += openedMarginValue
		pos.rawSize += rawFillQuantity
	}

	return fee, realizedCost
}

func (pos *Position) Sell(price, quantity float64, taker bool) (int64, int64) {
	if pos.inverse {
		price = 1. / price
	}

	rawFillQuantity := int64(quantity * pos.lotPrecision)
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
		closedMarginValue := int64(math.Round(float64(closedSize*rawPrice) / pos.lotPrecision))    //closedSize * rawPrice
		contractMarginValue := int64(math.Round(float64(pos.rawSize*rawPrice) / pos.lotPrecision)) //pos.rawSize * rawPrice

		unrealizedCost := pos.cost - contractMarginValue
		realizedCost = int64(math.Round(float64(closedSize*unrealizedCost) / float64(pos.rawSize)))

		// Transfer cost
		pos.cost -= closedMarginValue
		// Remove realized cost
		pos.cost -= realizedCost
		pos.rawSize -= closedSize
		rawFillQuantity -= closedSize
	}
	if rawFillQuantity > 0 {
		// We are opening a position
		openedMarginValue := int64(math.Round(float64(rawFillQuantity*rawPrice) / pos.lotPrecision))
		pos.cost -= openedMarginValue
		pos.rawSize -= rawFillQuantity
	}

	return fee, realizedCost
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
