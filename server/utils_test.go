// BEGIN: x9c3f8d4b3e6
package server

import (
	"math/big"
	"testing"
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
			want: 0,
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

// END: x9c3f8d4b3e6
