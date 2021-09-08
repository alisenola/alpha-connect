package account

import (
	"fmt"
	"testing"
)

func TestPosition(t *testing.T) {
	pos := NewPosition(false, 10, 10, 1e8, 1, 0, 0)
	pos.Buy(35.5, 100, false)
	var tot float64
	for i := 0; i < 1000; i++ {
		_, real := pos.Sell(71, 0.1, false)
		tot += -float64(real) / 1e8
	}
	fmt.Println(tot)

	tot = 0.
	pos = NewPosition(false, 10, 10, 1e8, 1, 0, 0)
	pos.Buy(35.5, 100, false)
	_, real := pos.Sell(71, 100, false)
	tot = -float64(real) / 1e8
	fmt.Println(tot)
}
