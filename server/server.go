package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/natac13/go-crypto-exchange/orderbook"
)

const (
	MarketETH Market = "ETH"

	MarketOrder OrderType = "MARKET"
	LimitOrder  OrderType = "LIMIT"

	// dont ever do this
	// user 0 is the exchange
	exchangePrivateKey = "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d"
)

type (
	OrderType string
	Market    string

	PlaceOrderRequest struct {
		UserID int64     `json:"userId"`
		Market Market    `json:"market"`
		Price  float64   `json:"price"`
		Size   float64   `json:"size"`
		Bid    bool      `json:"bid"`
		Type   OrderType `json:"type"` // market or limit
	}

	PlaceOrderResponse struct {
		OrderID int64  `json:"orderId"`
		Message string `json:"message"`
	}

	MatchedOrder struct {
		Price      float64 `json:"price"`
		SizeFilled float64 `json:"sizeFilled"`
		ID         int64   `json:"id"`
	}

	Order struct {
		ID        int64   `json:"id"`
		Price     float64 `json:"price"`
		Size      float64 `json:"size"`
		Bid       bool    `json:"bid"`
		Timestamp int64   `json:"timestamp"`
	}

	OrderbookResponse struct {
		Market         Market   `json:"market"`
		Asks           []*Order `json:"asks"`
		Bids           []*Order `json:"bids"`
		TotalBidVolume float64  `json:"totalBidVolume"`
		TotalAskVolume float64  `json:"totalAskVolume"`
	}
)

func seedUsers(ex *Exchange) error {
	// user 9
	// pkStr0 := "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d" // exchange
	pkStr1 := "6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1"
	pkStr2 := "6370fd033278c143179d81c5526140625662b8daa446c22ee2d73db3707e620c"
	pkStr3 := "646f1ce2fdad0e6deeeb5c7e8e5543bdde65e86029e2fd9fc169899c440a7913"
	pkStr4 := "add53f9a7e588d003326d1cbf9e4a43c061aadd9bc938c843a79e7b4fd2ad743"
	pkStr5 := "395df67f0c2d2d9fe1ad08d1bc8b6627011959b79c53d7dd6a3536a33ab8a4fd"
	pkStr6 := "e485d098507f54e7733a205420dfddbe58db035fa577fc294ebd14db90767a52"
	pkStr7 := "a453611d9419d0e56f499079478fd72c37b251a94bfde4d19872c44cf65386e3"
	pkStr8 := "829e924fdf021ba3dbbc4225edfece9aca04b929d6e75613329ca6f1d31c0bb4" // seller
	pkStr9 := "b0057716d5917badaf911b193b12b910811c1497b5bada8d7711f758981c3773" // buyer

	userData := []struct {
		pkStr string
		id    int64
	}{
		// {pkStr0, 0},
		{pkStr1, 1},
		{pkStr2, 2},
		{pkStr3, 3},
		{pkStr4, 4},
		{pkStr5, 5},
		{pkStr6, 6},
		{pkStr7, 7},
		{pkStr8, 8}, // seller
		{pkStr9, 9}, // buyer
	}

	for _, data := range userData {
		user, err := NewUser(data.pkStr, data.id)
		if err != nil {
			return err
		}
		ex.users[user.ID] = user
		pubKey := user.PublicKey
		pubAddress := crypto.PubkeyToAddress(*pubKey)
		balance, err := ex.Client.BalanceAt(context.Background(), pubAddress, nil)

		if err != nil {
			return err
		}

		fmt.Printf("user %d balance: %f\n", user.ID, weiToEth(balance))
	}
	return nil

}

func StartServer() {
	e := echo.New()
	e.HTTPErrorHandler = httpErrorHandler

	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}

	ex, err := NewExchange(exchangePrivateKey, client)
	if err != nil {
		log.Fatal(err)
	}

	err = seedUsers(ex)
	if err != nil {
		log.Fatal(err)
	}

	e.POST("/order", ex.handlePlaceOrder)
	e.DELETE("/order/:id", ex.handleCancelOrder)

	e.GET("/book/:market", ex.handleGetBook)
	e.GET("/book/:market/bids", ex.handleGetAllBids)
	e.GET("/book/:market/asks", ex.handleGetAllAsks)
	e.GET("/book/:market/best-bid", ex.handleGetBestBid)
	e.GET("/book/:market/best-ask", ex.handleGetBestAsk)

	e.Start(":3000")
}

func httpErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.JSON(code, err.Error())
}

func (ex *Exchange) handlePlaceMarketOrder(market Market, order *orderbook.Order) ([]orderbook.Match, []*MatchedOrder, error) {

	ob, ok := ex.orderbooks[market]

	if !ok {
		return nil, nil, fmt.Errorf("market not found")
	}
	matches := ob.PlaceMarketOrder(order)
	matchedOrders := make([]*MatchedOrder, len(matches))

	isBid := order.Bid

	totalSizeFilled := 0.0
	sumPrice := 0.0
	for i, match := range matches {
		id := match.Bid.ID
		if isBid {
			id = match.Ask.ID
		}
		matchedOrders[i] = &MatchedOrder{
			SizeFilled: match.SizeFilled,
			Price:      match.Price,
			ID:         id,
		}
		totalSizeFilled += match.SizeFilled
		sumPrice += match.Price * match.SizeFilled
	}

	avgPrice := sumPrice / totalSizeFilled

	log.Printf("filled MARKET order => id: {%d} bid: {%v} size filled: {%.2f} @ average price: {%.2f}", order.ID, order.Bid, totalSizeFilled, avgPrice)
	return matches, matchedOrders, nil
}

func (ex *Exchange) handlePlaceLimitOrder(market Market, price float64, order *orderbook.Order) error {
	ob, ok := ex.orderbooks[market]
	if !ok {
		return fmt.Errorf("market not found")
	}
	// transfer from the user to the exchange.
	// I don't think they really do this do to gas costs
	// they likey just keep track of the balances
	ob.PlaceLimitOrder(price, order)

	log.Printf("new LIMIT order => bid: {%v}  price: {%.2f}, size: {%.2f}", order.Bid, order.Limit.Price, order.Size)
	return nil
}

func (ex *Exchange) handlePlaceOrder(c echo.Context) error {
	var placeOrderData PlaceOrderRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return err
	}

	market := Market(placeOrderData.Market)
	order := orderbook.NewOrder(placeOrderData.Bid, placeOrderData.Size, placeOrderData.UserID)

	if placeOrderData.Type == LimitOrder {
		if err := ex.handlePlaceLimitOrder(market, placeOrderData.Price, order); err != nil {
			return err
		}
	}

	if placeOrderData.Type == MarketOrder {
		matches, _, err := ex.handlePlaceMarketOrder(market, order)
		if err != nil {
			return err
		}
		if err := ex.handleMatches(matches); err != nil {
			return err
		}

	}
	res := PlaceOrderResponse{
		OrderID: order.ID,
		Message: "order placed",
	}

	return c.JSON(http.StatusOK, res)
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

type PriceResponse struct {
	Price float64 `json:"price"`
}

func (ex *Exchange) handleGetBestBid(c echo.Context) error {
	market := Market(c.Param("market"))
	ob := ex.orderbooks[market]

	if len(ob.Bids()) == 0 {
		return fmt.Errorf("the bids are empty")
	}
	bestBidPrice := ob.Bids()[0].Price

	pr := PriceResponse{
		Price: bestBidPrice,
	}
	return c.JSON(http.StatusOK, pr)
}

func (ex *Exchange) handleGetBestAsk(c echo.Context) error {
	market := Market(c.Param("market"))
	ob := ex.orderbooks[market]

	if len(ob.Asks()) == 0 {
		return fmt.Errorf("the asks are empty")
	}
	bestAskPrice := ob.Asks()[0].Price

	pr := PriceResponse{
		Price: bestAskPrice,
	}
	return c.JSON(http.StatusOK, pr)
}

func (ex *Exchange) handleGetAllBids(c echo.Context) error {
	market := Market(c.Param("market"))
	ob := ex.orderbooks[market]

	bids := make([]*PriceResponse, len(ob.Bids()))
	for i, limit := range ob.Bids() {
		bids[i] = &PriceResponse{
			Price: limit.Price,
		}
	}

	return c.JSON(http.StatusOK, bids)
}

func (ex *Exchange) handleGetAllAsks(c echo.Context) error {
	market := Market(c.Param("market"))
	ob := ex.orderbooks[market]

	asks := make([]*PriceResponse, len(ob.Asks()))
	for i, limit := range ob.Asks() {
		asks[i] = &PriceResponse{
			Price: limit.Price,
		}
	}

	return c.JSON(http.StatusOK, asks)
}

func (ex *Exchange) handleCancelOrder(c echo.Context) error {
	idStr := c.Param("id")
	id, ok := strconv.Atoi(idStr)
	if ok != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"msg": "invalid id"})
	}

	ob := ex.orderbooks[MarketETH]
	order := ob.Orders[int64(id)]
	ob.CancelOrder(order)

	log.Println("order deleted, id: ", idStr, "market: ", "ETH-USD")

	return c.JSON(http.StatusOK, map[string]interface{}{"msg": "order deleted"})
}

func (ex *Exchange) handleMatches(matches []orderbook.Match) error {
	for _, match := range matches {
		fromUser, ok := ex.users[match.Ask.UserID]
		if !ok {
			return fmt.Errorf("ask user not found, ID: %d", match.Ask.UserID)
		}

		toUser, ok := ex.users[match.Bid.UserID]
		if !ok {
			return fmt.Errorf("bid user not found, ID: %d", match.Bid.UserID)
		}

		toAddress := crypto.PubkeyToAddress(toUser.PrivateKey.PublicKey)

		amount := orderSizeToWei(match.SizeFilled)
		if err := transferETH(ex.Client, fromUser.PrivateKey, toAddress, amount); err != nil {
			return err
		}
	}
	return nil
}
