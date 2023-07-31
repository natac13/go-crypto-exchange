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
	tick = 2 * time.Second
)

func marketOrderPlacer(c *client.Client) {
	ticker := time.NewTicker(tick)
	for {
		marketSellOrder := &client.PlaceMarketOrderParams{
			UserID: 9,
			Bid:    false,
			Size:   0.2,
		}

		_, err := c.PlaceMarketOrder(marketSellOrder)
		if err != nil {
			log.Fatal(err)
		}

		marketBuyOrder := &client.PlaceMarketOrderParams{
			UserID: 9,
			Bid:    true,
			Size:   0.2,
		}

		_, err = c.PlaceMarketOrder(marketBuyOrder)
		if err != nil {
			log.Fatal(err)
		}

		<-ticker.C
	}
}

func makeMarketSimple(c *client.Client) {
	ticker := time.NewTicker(tick)
	for {
		orders, err := c.GetOrders(7)
		if err != nil {
			log.Println(err)
		}

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
		fmt.Println("exchange spread: => ", spread)

		// place 2 orders to tighten the spread
		fmt.Println("===============================")
		if len(orders.Bids) < maxOrders {
			fmt.Println("placing bid")
			bidLimit := &client.PlaceLimitOrderParams{
				UserID: 7,
				Bid:    true,
				Price:  bestBid + 100,
				Size:   5,
			}
			_, err := c.PlaceLimitOrder(bidLimit)
			if err != nil {
				log.Println(err)
			}
		}
		if len(orders.Asks) < maxOrders {
			askLimit := &client.PlaceLimitOrderParams{
				UserID: 7,
				Bid:    false,
				Price:  bestAsk - 100,
				Size:   5,
			}
			_, err := c.PlaceLimitOrder(askLimit)
			if err != nil {
				log.Println(err)
			}
		}

		fmt.Println("best ask: => ", bestAsk)
		fmt.Println("best bid: => ", bestBid)

		<-ticker.C
	}
}

func seedMarket(c *client.Client) error {
	ask := &client.PlaceLimitOrderParams{
		UserID: 8,
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
	marketOrderPlacer(c)

	select {}

}
