package account

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFillCollector(t *testing.T) {
	fc := NewFillCollector(10000, 0., []uint64{1, 2, 3})
	tss := time.Now().UnixMilli()
	fc.AddFill(1, 1.0, true, true, tss)
	fc.AddFill(1, 1.0, false, true, tss)
	fc.AddFill(1, 1.0, true, false, tss)
	fc.AddFill(1, 1.0, false, false, tss)
	fc.AddFill(2, 1.0, true, true, tss)
	fc.AddFill(2, 1.0, false, true, tss)
	fc.AddFill(2, 1.0, true, false, tss)
	fc.AddFill(2, 1.0, false, false, tss)
	time.Sleep(500 * time.Millisecond)
	for i := 0; i < 10; i++ {
		fc.Collect(1, 2.0)
		time.Sleep(1 * time.Second)
	}
	mb, tb, ms, ts := fc.GetMoveAfterFill()
	for i := 0; i < 10; i++ {
		// mb should be 1.0
		assert.InDelta(t, 1.0, mb[i], 0.1)
		assert.InDelta(t, 1.0, tb[i], 0.1)
		assert.InDelta(t, -0.5, ms[i], 0.1)
		assert.InDelta(t, -0.5, ts[i], 0.1)
	}
	fmt.Println(mb, tb, ms, ts)
}
