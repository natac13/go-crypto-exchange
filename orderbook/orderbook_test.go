// BEGIN: 1c2d3e4f5g6h
package orderbook

// END: 1c2d3e4f5g6hpackage main

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrderA := NewOrder(true, 5, 0)
	buyOrderB := NewOrder(true, 10, 0)
	buyOrderC := NewOrder(true, 25, 0)

	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	fmt.Println(l)

	l.DeleteOrder(buyOrderB)

	fmt.Println(l)
	fmt.Println("done")
}

func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrderbook()

	sellOrderA := NewOrder(false, 10, 0)
	sellOrderB := NewOrder(false, 5, 0)
	ob.PlaceLimitOrder(10_000, sellOrderA)
	ob.PlaceLimitOrder(9_000, sellOrderB)

	assert(t, len(ob.Asks()), 2)
	assert(t, len(ob.Orders), 2)
	assert(t, ob.Orders[sellOrderA.ID], sellOrderA)
	assert(t, ob.Orders[sellOrderB.ID], sellOrderB)
}

func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrderbook()

	sellOrderA := NewOrder(false, 20, 0)
	ob.PlaceLimitOrder(10_000, sellOrderA)

	buyOrderA := NewOrder(true, 10, 0)
	matches := ob.PlaceMarketOrder(buyOrderA)

	assert(t, len(matches), 1)
	assert(t, len(ob.asks), 1)
	assert(t, len(ob.bids), 0)
	assert(t, ob.AskTotalVolume(), 10.0)
	assert(t, matches[0].Ask, sellOrderA)
	assert(t, matches[0].Bid, buyOrderA)
	assert(t, matches[0].Price, 10_000.0)
	assert(t, buyOrderA.IsFilled(), true)

	fmt.Printf("%+v", matches)
}

func TestPlaceMarketOrderMultiFill(t *testing.T) {
	ob := NewOrderbook()

	buyOrderA := NewOrder(true, 5, 0) // fully filled
	buyOrderB := NewOrder(true, 8, 0) // partially filled
	buyOrderC := NewOrder(true, 1, 0) // un filled
	buyOrderD := NewOrder(true, 1, 0) // un filled

	ob.PlaceLimitOrder(5_000, buyOrderC)
	ob.PlaceLimitOrder(5_000, buyOrderD)
	ob.PlaceLimitOrder(9_000, buyOrderB)
	ob.PlaceLimitOrder(10_000, buyOrderA)

	// when we place a sell market order we want to fill the highest price first
	// theerfore we should be left with a order at 5_000 for 3

	assert(t, ob.BidTotalVolume(), 15.0)
	assert(t, len(ob.bids), 3)

	sellOrderA := NewOrder(false, 10, 0)
	matches := ob.PlaceMarketOrder(sellOrderA)

	assert(t, len(matches), 2)
	// need to make sure that the filled orders are removed from the orderbook
	// assert(t, ob.BidTotalVolume(), 5.0)
	// assert(t, sellOrderA.IsFilled(), true)
	assert(t, len(ob.bids), 2)
	// assert(t, len(ob.bids[0].Orders), 2)
}

func TestPlaceMarketOrderMultiFillWithReversedSamePriceBid(t *testing.T) {
	ob := NewOrderbook()

	buyOrderA := NewOrder(true, 5, 0)
	buyOrderB := NewOrder(true, 8, 0)
	buyOrderC := NewOrder(true, 10, 0)
	buyOrderD := NewOrder(true, 1, 0)

	ob.PlaceLimitOrder(5_000, buyOrderD)
	ob.PlaceLimitOrder(5_000, buyOrderC)
	ob.PlaceLimitOrder(9_000, buyOrderB)
	ob.PlaceLimitOrder(10_000, buyOrderA)

	// when we place a sell market order we want to fill the highest price first
	// theerfore we should be left with a order at 5_000 for 3

	assert(t, ob.BidTotalVolume(), 24.0)
	assert(t, len(ob.bids), 3)

	sellOrderA := NewOrder(false, 22, 0)
	matches := ob.PlaceMarketOrder(sellOrderA)

	assert(t, len(matches), 4)
	// need to make sure that the filled orders are removed from the orderbook
	assert(t, ob.BidTotalVolume(), 2.0)
	assert(t, sellOrderA.IsFilled(), true)
	assert(t, len(ob.bids), 1)
	assert(t, len(ob.bids[0].Orders), 1)

}

func TestCancelOrder(t *testing.T) {
	ob := NewOrderbook()

	buyOrderA := NewOrder(true, 5, 0)
	buyOrderB := NewOrder(true, 8, 0)
	buyOrderC := NewOrder(true, 10, 0)
	buyOrderD := NewOrder(true, 1, 0)

	ob.PlaceLimitOrder(5_000, buyOrderD)
	ob.PlaceLimitOrder(5_000, buyOrderC)
	ob.PlaceLimitOrder(9_000, buyOrderB)
	ob.PlaceLimitOrder(10_000, buyOrderA)

	assert(t, len(ob.bids), 3)
	assert(t, ob.BidTotalVolume(), 24.0)

	assert(t, len(ob.Orders), 4)
	ob.CancelOrder(buyOrderB)

	assert(t, len(ob.bids), 2)
	assert(t, ob.BidTotalVolume(), 16.0)

	assert(t, len(ob.Orders), 3)

	_, ok := ob.Orders[buyOrderB.ID]
	assert(t, ok, false)

	_, ok = ob.BidLimits[9_000]
	assert(t, ok, false)
}

func TestCancelOrderAsk(t *testing.T) {
	ob := NewOrderbook()

	sellOrderA := NewOrder(false, 5, 0)

	ob.PlaceLimitOrder(10_000, sellOrderA)

	assert(t, len(ob.asks), 1)
	assert(t, ob.AskTotalVolume(), 5.0)

	assert(t, len(ob.Orders), 1)
	ob.CancelOrder(sellOrderA)

	assert(t, len(ob.asks), 0)
	assert(t, ob.AskTotalVolume(), 0.0)

	_, ok := ob.Orders[sellOrderA.ID]
	assert(t, ok, false)

	_, ok = ob.AskLimits[10_000]
	assert(t, ok, false)
}
