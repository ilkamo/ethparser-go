package ethereum

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasEthNumberPrefix(t *testing.T) {
	testCases := []struct {
		number   string
		expected bool
	}{
		{
			number:   "0x0",
			expected: true,
		},
		{
			number:   "0x123",
			expected: true,
		},
		{
			number:   "0x98a7d9b8314c0000",
			expected: true,
		},
		{
			number:   "123",
			expected: false,
		},
		{
			number:   "0x",
			expected: true,
		},
		{
			number:   "x0",
			expected: false,
		},
		{
			number:   "xx",
			expected: false,
		},
		{
			number:   "0",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("chech %s eth number prefix", tc.number), func(t *testing.T) {
			got := HasEthNumberPrefix(tc.number)
			if got != tc.expected {
				t.Errorf("expected %t, got %t", tc.expected, got)
			}
		})
	}
}

func TestUint64FromEthNumber(t *testing.T) {
	testCases := []struct {
		number   string
		expected uint64
		wantErr  bool
	}{
		{
			number:   "0x0",
			expected: 0,
			wantErr:  false,
		},
		{
			number:   "0x123",
			expected: 291,
			wantErr:  false,
		},
		{
			number:   "0x98a7d9b8314c0000",
			expected: 11000000000000000000,
			wantErr:  false,
		},
		{
			number:  "123",
			wantErr: true,
		},
		{
			number:  "0x",
			wantErr: true,
		},
		{
			number:  "0xg",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("decode %s", tc.number), func(t *testing.T) {
			got, err := Uint64FromEthNumber(tc.number)
			if err != nil && !tc.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, got)
			}
		})
	}
}

func TestEthNumberFromUnit64(t *testing.T) {
	testCases := []struct {
		number   uint64
		expected string
	}{
		{
			number:   0,
			expected: "0x0",
		},
		{
			number:   291,
			expected: "0x123",
		},
		{
			number:   11000000000000000000,
			expected: "0x98a7d9b8314c0000",
		},
		{
			number:   123,
			expected: "0x7b",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("encode %d", tc.number), func(t *testing.T) {
			got := EthNumberFromUnit64(tc.number)
			if got != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestBigIntFromEthNumber(t *testing.T) {
	testCases := []struct {
		number   string
		expected *big.Int
		wantErr  bool
	}{
		{
			number:   "0x0",
			expected: big.NewInt(0),
			wantErr:  false,
		},
		{
			number:   "0x123",
			expected: big.NewInt(291),
			wantErr:  false,
		},
		{
			number:   "0x98a7d9b8314c0000",
			expected: mustBigIntFromString("11000000000000000000"),
			wantErr:  false,
		},
		{
			number:  "123",
			wantErr: true,
		},
		{
			number:  "0x",
			wantErr: true,
		},
		{
			number:  "0xg",
			wantErr: true,
		},
		{
			number:   "0x197D4DF19D605767337E9F14D3EEC8920E400000000000000",
			expected: mustBigIntFromString("10000000000000000000000000000000000000000000000000000000000"),
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("decode %s", tc.number), func(t *testing.T) {
			got, err := BigIntFromEthNumber(tc.number)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected.String(), got.String())
			}
		})
	}
}

func mustBigIntFromString(s string) *big.Int {
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("could not decode hex number to string")
	}

	return i
}
