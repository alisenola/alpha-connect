package account

import "testing"

func TestPositionSessionPnL(t *testing.T) {
	// Without fees
	pos := NewPosition(0, false, 1, 1, 1, 1, 0, 0)
	pos.Buy(100, 1, false)
	if pos.sessionPnL != 0 {
		t.Fatalf("was expecting 0, got %d", pos.sessionPnL)
	}
	pos.Sell(110, 1, false)
	if pos.sessionPnL != 10 {
		t.Fatalf("was expecting 10, got %d", pos.sessionPnL)
	}
	pos = NewPosition(0, false, 1, 1, 1, 1, 0, 0)
	pos.Buy(110, 1, false)
	if pos.sessionPnL != 0 {
		t.Fatalf("was expecting 0, got %d", pos.sessionPnL)
	}
	pos.Sell(100, 1, false)
	if pos.sessionPnL != -10 {
		t.Fatalf("was expecting -10, got %d", pos.sessionPnL)
	}
	// With fees
	pos = NewPosition(0, false, 1, 1, 1, 1, 0.01, 0)
	pos.Buy(100, 1, false)
	if pos.sessionPnL != -1 {
		t.Fatalf("was expecting -1, got %d", pos.sessionPnL)
	}
	pos.Sell(100, 1, false)
	if pos.sessionPnL != -2 {
		t.Fatalf("was expecting -2, got %d", pos.sessionPnL)
	}
}
