package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/natac13/go-crypto-exchange/client"
	"github.com/natac13/go-crypto-exchange/server"
)

const (
	maxOrders = 3
)

var (
	tick   = 5 * time.Second
	myAsks = make(map[float64]int64)
	myBids = make(map[float64]int64)
)

func marketOrderPlacer(c *client.Client) {
	ticker := time.NewTicker(tick)
	for {
		marketSellOrder := &client.PlaceMarketOrderParams{
			UserID: 9,
			Bid:    false,
			Size:   1,
		}

		sellRes, err := c.PlaceMarketOrder(marketSellOrder)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("placed market order from the client: => ", sellRes.OrderID)

		marketBuyOrder := &client.PlaceMarketOrderParams{
			UserID: 9,
			Bid:    true,
			Size:   1,
		}

		buyRes, err := c.PlaceMarketOrder(marketBuyOrder)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("placed market order from the client: => ", buyRes.OrderID)

		<-ticker.C
	}
}

func makeMarketSimple(c *client.Client) {
	ticker := time.NewTicker(tick)
	for {
		// place the ask
		// get the best ask and best bid, and calculate the spread
		bestAsk, err := c.GetBestAsk()
		if err != nil {
			log.Println(err)
		}
		bestBid, err := c.GetBestBid()
		if err != nil {
			log.Println(err)
		}
		spread := math.Abs(bestAsk - bestBid)
		fmt.Println("spread: => ", spread)

		// place 2 orders to tighten the spread
		if len(myBids) < maxOrders {
			bidLimit := &client.PlaceLimitOrderParams{
				UserID: 7,
				Bid:    true,
				Price:  bestBid + 100,
				Size:   1,
			}
			bidRes, err := c.PlaceLimitOrder(bidLimit)
			if err != nil {
				log.Println(err)
			}
			myBids[bidLimit.Price] = bidRes.OrderID
			fmt.Print("bid order id: => ", bidRes)
		}
		if len(myAsks) < maxOrders {
			askLimit := &client.PlaceLimitOrderParams{
				UserID: 7,
				Bid:    false,
				Price:  bestAsk - 100,
				Size:   1,
			}
			askRes, err := c.PlaceLimitOrder(askLimit)
			if err != nil {
				log.Println(err)
			}
			myAsks[askLimit.Price] = askRes.OrderID
			fmt.Print("ask order id: => ", askRes)
		}

		fmt.Println("best ask: => ", bestAsk)
		fmt.Println("best bid: => ", bestBid)

		<-ticker.C
	}
}

func seedMarket(c *client.Client) error {
	ask := &client.PlaceLimitOrderParams{
		UserID: 9,
		Bid:    false,
		Price:  10_000,
		Size:   10,
	}

	bid := &client.PlaceLimitOrderParams{
		UserID: 8,
		Bid:    true,
		Price:  9_000,
		Size:   10,
	}

	_, err := c.PlaceLimitOrder(ask)
	if err != nil {
		return err
	}

	_, err = c.PlaceLimitOrder(bid)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	go server.StartServer()

	time.Sleep(1 * time.Second)

	c := client.NewClient()

	if err := seedMarket(c); err != nil {
		panic(err)
	}

	go makeMarketSimple(c)
	time.Sleep(1 * time.Second)
	go marketOrderPlacer(c)

	// for {

	// randAskPrice := rand.Intn(10_000)
	// limitParams := &clientPkg.PlaceLimitOrderParams{
	// 	UserID: 8,
	// 	Bid:    false,
	// 	Price:  float64(randAskPrice),
	// 	Size:   5,
	// }

	// _, err := client.PlaceLimitOrder(limitParams)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// randPrice := rand.Intn(10_000)
	// otherLimitOrderParams := &clientPkg.PlaceLimitOrderParams{
	// 	UserID: 9,
	// 	Bid:    true,
	// 	Price:  float64(randPrice),
	// 	Size:   5,
	// }

	// _, err = client.PlaceLimitOrder(otherLimitOrderParams)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// bestBid, err := client.GetBestBid()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("best bid: => ", bestBid)

	// bestAsk, err := client.GetBestAsk()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("best ask: => ", bestAsk)

	// fmt.Println("placed limit order from the client: => ", res.OrderID)

	// marketParams := &clientPkg.PlaceMarketOrderParams{
	// 	UserID: 7,
	// 	Bid:    true,
	// 	Size:   5,
	// }

	// _, err = client.PlaceMarketOrder(marketParams)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("placed market order from the client: => ", res.OrderID)

	time.Sleep(1 * time.Second)
	// }

	select {}

}
