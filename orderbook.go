package main

import (
	"fmt"
	"sort"
	"time"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

type Order struct {
	Size      float64
	Bid       bool
	Limit     *Limit
	Timestamp int64
}

type Orders []*Order

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }

func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Size:      size,
		Bid:       bid,
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("[size: %.2f]", o.Size)
}

type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit

type ByBestAsk struct{ Limits }

func (a ByBestAsk) Len() int      { return len(a.Limits) }
func (a ByBestAsk) Swap(i, j int) { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }

// When sorting asks, we want the lowest price first (ascending)
func (a ByBestAsk) Less(i, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }

type ByBestBid struct{ Limits }

func (a ByBestBid) Len() int      { return len(a.Limits) }
func (a ByBestBid) Swap(i, j int) { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }

// When sorting bids, we want the highest price first (descending)
func (a ByBestBid) Less(i, j int) bool { return a.Limits[i].Price > a.Limits[j].Price }

func NewLimit(price float64) *Limit {
	return &Limit{Price: price, Orders: []*Order{}}
}

func (l *Limit) String() string {
	return fmt.Sprintf("{price: %.2f, totalVolume: %.2f, orders: %v}", l.Price, l.TotalVolume, l.Orders)
}

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (l *Limit) DeleteOrder(o *Order) {
	for i, order := range l.Orders {
		if order == o {
			l.Orders = append(l.Orders[:i], l.Orders[i+1:]...)
			// return
		}
	}
	o.Limit = nil
	l.TotalVolume -= o.Size

	sort.Sort(l.Orders)
}

type Orderbook struct {
	Asks []*Limit
	Bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		Asks:      []*Limit{},
		Bids:      []*Limit{},
		AskLimits: map[float64]*Limit{},
		BidLimits: map[float64]*Limit{},
	}
}

func (ob *Orderbook) PlaceOrder(price float64, o *Order) []Match {
	// 1 Try to match with existing orders
	// matching logic *********
	// if completetly filled the size will be zero
	// if the size is not zero, then we need to add the order to the books to be filled later

	// 2. add the rest of the order to the books
	if o.Size > 0.0 {
		ob.add(price, o)
	}

	return []Match{}
}

// BUY 10 BTC @ 15K
func (ob *Orderbook) add(price float64, o *Order) {
	var limit *Limit

	if o.Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AskLimits[price]
	}

	if limit == nil {
		limit = NewLimit(price)
		if o.Bid {
			ob.BidLimits[price] = limit
			ob.Bids = append(ob.Bids, limit)
		} else {
			ob.AskLimits[price] = limit
			ob.Asks = append(ob.Asks, limit)
		}
	}

	limit.AddOrder(o)
	// o.Limit = limit
}
