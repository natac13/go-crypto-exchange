// BEGIN: 1c2d3e4f5g6h
package main

// END: 1c2d3e4f5g6hpackage main

import (
	"fmt"
	"testing"
)

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrderA := NewOrder(true, 5)
	buyOrderB := NewOrder(true, 10)
	buyOrderC := NewOrder(true, 25)

	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	fmt.Println(l)

	l.DeleteOrder(buyOrderB)

	fmt.Println(l)
	fmt.Println("done")
}

func TestOrderbook(t *testing.T) {
	ob := NewOrderbook()

	// Add some bids and asks
	buyOrderA := NewOrder(true, 10)
	buyOrderB := NewOrder(true, 2_000)
	ob.PlaceOrder(18_000, buyOrderA)
	ob.PlaceOrder(19_000, buyOrderB)

	for _, limit := range ob.Bids {
		fmt.Println(limit)
		fmt.Printf("TotalVolume: %.2f\n", limit.TotalVolume)
		fmt.Printf("%+v\n", limit.Orders)
	}

	fmt.Printf("%+v", ob)

}
