package utils

import (
	"fmt"
	"testing"
)

func TestNewRingSet(t *testing.T) {
	rs := NewRingSet(100)
	for i := 0; i < 500; i++ {
		rs.Insert(i, nil)
	}
	fmt.Println(rs.tail, rs.head)
	fmt.Println(rs.buffer)
	for i := 400; i < 500; i++ {
		if !rs.Contains(i) {
			t.Fatalf("ring set does not contain inserted element %d", i)
		}
	}
}

func TestUpsert(t *testing.T) {
	rs := NewRingSet(100)
	for i := 0; i < 500; i++ {
		rs.Insert(i, 1)
	}
	for i := 0; i < 500; i++ {
		rs.Upsert(i, 2)
	}
	fmt.Println(rs.tail, rs.head)
	fmt.Println(rs.buffer)
	for i := 400; i < 500; i++ {
		if !rs.Contains(i) {
			t.Fatalf("ring set does not contain inserted element %d", i)
		}
		if rs.Get(i) != 2 {
			t.Fatalf("was expecting 2, got %d", rs.Get(i))
		}
	}
}
