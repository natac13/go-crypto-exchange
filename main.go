package main

import (
	"log"
	"time"

	clientPkg "github.com/natac13/go-crypto-exchange/client"
	"github.com/natac13/go-crypto-exchange/server"
)

func main() {
	go server.StartServer()

	time.Sleep(1 * time.Second)

	client := clientPkg.NewClient()

	for {

		limitParams := &clientPkg.PlaceLimitOrderParams{
			UserID: 8,
			Bid:    false,
			Price:  10_000,
			Size:   3,
		}

		_, err := client.PlaceLimitOrder(limitParams)
		if err != nil {
			log.Fatal(err)
		}
		otherLimitOrderParams := &clientPkg.PlaceLimitOrderParams{
			UserID: 9,
			Bid:    false,
			Price:  9_000,
			Size:   10,
		}

		_, err = client.PlaceLimitOrder(otherLimitOrderParams)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println("placed limit order from the client: => ", res.OrderID)

		marketParams := &clientPkg.PlaceMarketOrderParams{
			UserID: 7,
			Bid:    true,
			Size:   13,
		}

		_, err = client.PlaceMarketOrder(marketParams)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println("placed market order from the client: => ", res.OrderID)

		time.Sleep(1 * time.Second)
	}

	select {}

}
