package hexutils

import (
	"fmt"
	"testing"
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

func TestDecodeEthNumber(t *testing.T) {
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
			got, err := DecodeEthNumberToUint(tc.number)
			if err != nil && !tc.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, got)
			}
		})
	}
}

func TestEncodeToEthNumber(t *testing.T) {
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
			got := EncodeToEthNumber(tc.number)
			if got != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestDecodeEthNumberToStr(t *testing.T) {
	testCases := []struct {
		number   string
		expected string
		wantErr  bool
	}{
		{
			number:   "0x0",
			expected: "0",
			wantErr:  false,
		},
		{
			number:   "0x123",
			expected: "291",
			wantErr:  false,
		},
		{
			number:   "0x98a7d9b8314c0000",
			expected: "11000000000000000000",
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
			expected: "10000000000000000000000000000000000000000000000000000000000",
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("decode %s", tc.number), func(t *testing.T) {
			got, err := DecodeEthNumberToStr(tc.number)
			if err != nil && !tc.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}
