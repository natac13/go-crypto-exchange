package main

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/natac13/go-crypto-exchange/orderbook"
)

func main() {
	e := echo.New()
	ex := NewExchange()

	e.GET("/book/:market", ex.handleGetBook)
	e.POST("/order", ex.handlePlaceOrder)
	e.Start(":3000")

}

type Market string

const (
	MarketETH Market = "ETH"
	// MarketBTC Market = "BTC"
)

type Exchange struct {
	orderbooks map[Market]*orderbook.Orderbook
}

func NewExchange() *Exchange {
	orderbooks := make(map[Market]*orderbook.Orderbook)
	orderbooks[MarketETH] = orderbook.NewOrderbook()

	return &Exchange{
		orderbooks: orderbooks,
	}
}

type OrderType string

const (
	MarketOrder OrderType = "MARKET"
	LimitOrder  OrderType = "LIMIT"
)

type PlaceOrderRequest struct {
	Market Market    `json:"market"`
	Price  float64   `json:"price"`
	Size   float64   `json:"size"`
	Bid    bool      `json:"bid"`
	Type   OrderType `json:"type"` // market or limit
}

func (ex *Exchange) handlePlaceOrder(c echo.Context) error {
	var placeOrderData PlaceOrderRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return err
	}

	market := Market(placeOrderData.Market)
	ob, ok := ex.orderbooks[market]
	if !ok {
		panic("market not found")
	}
	order := orderbook.NewOrder(placeOrderData.Bid, placeOrderData.Size)

	ob.PlaceLimitOrder(placeOrderData.Price, order)

	return c.JSON(http.StatusOK, map[string]interface{}{"msg": "order placed"})
}

type Order struct {
	Price     float64 `json:"price"`
	Size      float64 `json:"size"`
	Bid       bool    `json:"bid"`
	Timestamp int64   `json:"timestamp"`
}

type OrderbookResponse struct {
	Market Market   `json:"market"`
	Asks   []*Order `json:"asks"`
	Bids   []*Order `json:"bids"`
}

func (ex *Exchange) handleGetBook(c echo.Context) error {
	market := Market(c.Param("market"))
	ob, ok := ex.orderbooks[market]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"msg": "market not found"})
	}

	orderbookResponse := OrderbookResponse{
		Market: market,
		Asks:   []*Order{},
		Bids:   []*Order{},
	}
	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			o := Order{
				Price:     limit.Price,
				Size:      order.Size,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}
			orderbookResponse.Asks = append(orderbookResponse.Asks, &o)
		}
	}

	for _, limit := range ob.Bids() {
		for _, order := range limit.Orders {
			o := Order{
				Price:     limit.Price,
				Size:      order.Size,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}
			orderbookResponse.Bids = append(orderbookResponse.Bids, &o)
		}
	}

	return c.JSON(http.StatusOK, orderbookResponse)
}
