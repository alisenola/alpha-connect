package account

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFillCollector(t *testing.T) {
	fc := NewFillCollector(10000, 0., 100, []uint64{1, 2, 3})
	tss := time.Now()
	fc.AddFill(0, 1, 1.0, true, true, tss)
	fc.AddFill(0, 1, 1.0, false, true, tss)
	fc.AddFill(0, 1, 1.0, true, false, tss)
	fc.AddFill(0, 1, 1.0, false, false, tss)
	fc.AddFill(0, 2, 1.0, true, true, tss)
	fc.AddFill(0, 2, 1.0, false, true, tss)
	fc.AddFill(0, 2, 1.0, true, false, tss)
	fc.AddFill(0, 2, 1.0, false, false, tss)
	time.Sleep(500 * time.Millisecond)
	for i := 0; i < 10; i++ {
		fc.Collect(1, 2.0)
		time.Sleep(1 * time.Second)
	}
	mb, tb, ms, ts := fc.GetMoveAfterFill()
	for i := 0; i < 10; i++ {
		// Price increased, so move after fill of buy should be positive
		assert.InDelta(t, 1.0, mb[i], 0.1)
		assert.InDelta(t, 1.0, tb[i], 0.1)
		// Price increased, so move after fill of sell should be negative
		assert.InDelta(t, -0.5, ms[i], 0.1)
		assert.InDelta(t, -0.5, ts[i], 0.1)
	}
	fmt.Println(mb, tb, ms, ts)
}
