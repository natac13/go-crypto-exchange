package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/natac13/go-crypto-exchange/orderbook"
)

func main() {
	e := echo.New()
	ex := NewExchange()

	e.GET("/book/:market", ex.handleGetBook)
	e.POST("/order", ex.handlePlaceOrder)
	e.DELETE("/order/:id", ex.handleCancelOrder)
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

	if placeOrderData.Type == LimitOrder {
		ob.PlaceLimitOrder(placeOrderData.Price, order)
		return c.JSON(http.StatusOK, map[string]interface{}{"msg": "limit order placed"})
	}

	if placeOrderData.Type == MarketOrder {
		matches := ob.PlaceMarketOrder(order)
		return c.JSON(http.StatusOK, map[string]interface{}{"msg": "market order placed", "matchLength": len(matches)})
	}

	return nil
}

type Order struct {
	ID        int64   `json:"id"`
	Price     float64 `json:"price"`
	Size      float64 `json:"size"`
	Bid       bool    `json:"bid"`
	Timestamp int64   `json:"timestamp"`
}

type OrderbookResponse struct {
	Market         Market   `json:"market"`
	Asks           []*Order `json:"asks"`
	Bids           []*Order `json:"bids"`
	TotalBidVolume float64  `json:"totalBidVolume"`
	TotalAskVolume float64  `json:"totalAskVolume"`
}

func (ex *Exchange) handleGetBook(c echo.Context) error {
	market := Market(c.Param("market"))
	ob, ok := ex.orderbooks[market]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"msg": "market not found"})
	}

	orderbookResponse := OrderbookResponse{
		Market:         market,
		Asks:           []*Order{},
		Bids:           []*Order{},
		TotalBidVolume: ob.BidTotalVolume(),
		TotalAskVolume: ob.AskTotalVolume(),
	}

	for _, limit := range ob.Asks() {
		for _, o := range limit.Orders {
			order := Order{
				ID:        o.ID,
				Price:     limit.Price,
				Size:      o.Size,
				Bid:       o.Bid,
				Timestamp: o.Timestamp,
			}
			orderbookResponse.Asks = append(orderbookResponse.Asks, &order)
		}
	}

	for _, limit := range ob.Bids() {
		for _, o := range limit.Orders {
			order := Order{
				ID:        o.ID,
				Price:     limit.Price,
				Size:      o.Size,
				Bid:       o.Bid,
				Timestamp: o.Timestamp,
			}
			orderbookResponse.Bids = append(orderbookResponse.Bids, &order)
		}
	}

	return c.JSON(http.StatusOK, orderbookResponse)
}

func (ex *Exchange) handleCancelOrder(c echo.Context) error {
	idStr := c.Param("id")
	id, ok := strconv.Atoi(idStr)
	if ok != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"msg": "invalid id"})
	}

	ob := ex.orderbooks[MarketETH]
	orderCancelled := false

	for _, limit := range ob.Asks() {
		for _, o := range limit.Orders {
			if o.ID == int64(id) {
				ob.CancelOrder(o)
				orderCancelled = true
				return c.JSON(http.StatusOK, map[string]interface{}{"msg": "order cancelled", "orderCancelled": orderCancelled})
			}
		}
	}

	for _, limit := range ob.Bids() {
		for _, o := range limit.Orders {
			if o.ID == int64(id) {
				ob.CancelOrder(o)
				orderCancelled = true
				return c.JSON(http.StatusOK, map[string]interface{}{"msg": "order cancelled", "orderCancelled": orderCancelled})
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"msg": "order cancelled", "orderCancelled": orderCancelled})
}
