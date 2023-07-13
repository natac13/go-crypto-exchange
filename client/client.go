package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/natac13/go-crypto-exchange/server"
)

const EndPoint = "http://localhost:3000"

type Client struct {
	*http.Client
}

func NewClient() *Client {
	return &Client{
		Client: http.DefaultClient,
	}
}

type PlaceLimitOrderParams struct {
	UserID int64   `json:"userId"`
	Bid    bool    `json:"bid"`
	Price  float64 `json:"price"`
	Size   float64 `json:"size"`
}

type PlaceMarketOrderParams struct {
	UserID int64   `json:"userId"`
	Bid    bool    `json:"bid"`
	Size   float64 `json:"size"`
}

func (c *Client) PlaceMarketOrder(p *PlaceMarketOrderParams) (*server.PlaceOrderResponse, error) {
	params := &server.PlaceOrderRequest{
		UserID: p.UserID,
		Type:   server.MarketOrder,
		Bid:    p.Bid,
		Size:   p.Size,
		Market: server.MarketETH,
	}
	body, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	e := EndPoint + "/order"
	req, err := http.NewRequest("POST", e, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	response, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var placeLimitOrderResponse server.PlaceOrderResponse
	if err := json.NewDecoder(response.Body).Decode(&placeLimitOrderResponse); err != nil {
		return nil, err
	}

	return &placeLimitOrderResponse, nil
}

func (c *Client) PlaceLimitOrder(p *PlaceLimitOrderParams) (*server.PlaceOrderResponse, error) {
	params := &server.PlaceOrderRequest{
		UserID: p.UserID,
		Type:   server.LimitOrder,
		Bid:    p.Bid,
		Size:   p.Size,
		Price:  p.Price,
		Market: server.MarketETH,
	}
	body, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	e := EndPoint + "/order"
	req, err := http.NewRequest("POST", e, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	response, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var placeLimitOrderResponse server.PlaceOrderResponse
	if err := json.NewDecoder(response.Body).Decode(&placeLimitOrderResponse); err != nil {
		return nil, err
	}

	return &placeLimitOrderResponse, nil
}

func (c *Client) CancelOrder(orderId int64) error {
	e := fmt.Sprintf("%s/order/%d", EndPoint, orderId)

	req, err := http.NewRequest("DELETE", e, nil)
	if err != nil {
		return err
	}

	response, err := c.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}
