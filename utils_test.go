// BEGIN: x9c3f8d4b3e6
package main

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestWeiToEth(t *testing.T) {
	testCases := []struct {
		name string
		wei  *big.Int
		want float64
	}{
		{
			name: "1 wei",
			wei:  big.NewInt(1),
			want: 0.000000000000000001,
		},
		{
			name: "1 ether",
			wei:  big.NewInt(1000000000000000000),
			want: 1,
		},
		// {
		// 	name: "10 ether",
		// 	wei:  big.NewInt(int64(10000000000000000000)),
		// 	want: 10,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := weiToEth(tc.wei)
			if got != tc.want {
				t.Errorf("weiToEth(%v) = %v; want %v", tc.wei, got, tc.want)
			}
		})
	}
}

func TestTransferETH(t *testing.T) {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		t.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA("YOUR_PRIVATE_KEY_HERE")
	if err != nil {
		t.Fatalf("Failed to parse private key: %v", err)
	}

	toAddress := common.HexToAddress("YOUR_TO_ADDRESS_HERE")
	amount := orderSizeToWei(1.0)

	err = transferETH(client, privateKey, toAddress, amount)
	if err != nil {
		t.Fatalf("Failed to transfer ETH: %v", err)
	}
}

// END: x9c3f8d4b3e6
