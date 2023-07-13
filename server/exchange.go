package server

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/natac13/go-crypto-exchange/orderbook"
)

type Exchange struct {
	users      map[int64]*User
	orders     map[int64]*int64
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
		orders:     make(map[int64]*int64),
		Client:     client,
	}, nil
}
