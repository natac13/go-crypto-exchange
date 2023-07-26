package server

import (
	"crypto/ecdsa"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/natac13/go-crypto-exchange/orderbook"
)

type Exchange struct {
	users map[int64]*User
	// map user id to their orders
	mu         sync.RWMutex
	Orders     map[int64][]*orderbook.Order
	orderbooks map[Market]*orderbook.Orderbook
	PrivateKey *ecdsa.PrivateKey
	Client     *ethclient.Client
}

func NewExchange(privateKey string, client *ethclient.Client) (*Exchange, error) {
	orderbooks := make(map[Market]*orderbook.Orderbook)
	orderbooks[MarketETH] = orderbook.NewOrderbook()

	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	// publicAddress := crypto.PubkeyToAddress(pk.PublicKey)

	return &Exchange{
		orderbooks: orderbooks,
		PrivateKey: pk,
		users:      make(map[int64]*User),
		Orders:     make(map[int64][]*orderbook.Order),
		Client:     client,
	}, nil
}
