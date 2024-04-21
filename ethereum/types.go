package ethereum

import (
	"fmt"
	"time"

	"github.com/ilkamo/ethparser-go/hexutils"
	"github.com/ilkamo/ethparser-go/types"
)

// Block transport layer data structure.
type block struct {
	Number       string        `json:"number"`
	Hash         string        `json:"hash"`
	ParentHash   string        `json:"parentHash"`
	Timestamp    string        `json:"timestamp"`
	Transactions []transaction `json:"transactions"`
}

func (b block) ToBlock() (types.Block, error) {
	parsedNumber, err := hexutils.DecodeEthNumberToUint(b.Number)
	if err != nil {
		return types.Block{}, fmt.Errorf("could not decode block number: %w", err)
	}

	// The timestamp field in an Ethereum block is a 256-bit value representing the Unix timestamp
	// of when the block was mined. The Unix timestamp represents time as the number of
	// seconds elapsed since January 1, 1970, at 00:00:00 UTC
	elapsedSeconds, err := hexutils.DecodeEthNumberToUint(b.Timestamp)
	if err != nil {
		return types.Block{}, fmt.Errorf("could not decode block timestamp: %w", err)
	}

	parsedTimestamp := time.Unix(int64(elapsedSeconds), 0)
	if err != nil {
		return types.Block{}, fmt.Errorf("could not decode block unix time: %w", err)
	}

	transactions := make([]types.Transaction, len(b.Transactions))
	for i, t := range b.Transactions {
		tx, err := t.ToTransaction()
		if err != nil {
			return types.Block{}, err
		}
		transactions[i] = tx
	}

	return types.Block{
		Number:       parsedNumber,
		Hash:         b.Hash,
		ParentHash:   b.ParentHash,
		Timestamp:    parsedTimestamp,
		Transactions: transactions,
	}, nil
}

// Transaction transport layer data structure.
type transaction struct {
	BlockHash            string `json:"blockHash"`
	BlockNumber          string `json:"blockNumber"`
	From                 string `json:"from"`
	Gas                  string `json:"gas"`
	GasPrice             string `json:"gasPrice"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	Hash                 string `json:"hash"`
	Input                string `json:"input"`
	Nonce                string `json:"nonce"`
	To                   string `json:"to"`
	TransactionIndex     string `json:"transactionIndex"`
	Value                string `json:"value"`
	Type                 string `json:"type"`
	ChainId              string `json:"chainId"`
	V                    string `json:"v"`
	R                    string `json:"r"`
	S                    string `json:"s"`
	YParity              string `json:"yParity"`
}

func (t transaction) ToTransaction() (types.Transaction, error) {
	parsedNumber, err := hexutils.DecodeEthNumberToUint(t.BlockNumber)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("could not decode tx block number: %w", err)
	}

	parsedValue, err := hexutils.DecodeEthNumberToStr(t.Value)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("could not decode tx value: %w", err)
	}

	return types.Transaction{
		BlockHash:   t.BlockHash,
		BlockNumber: parsedNumber,
		Hash:        t.Hash,
		From:        t.From,
		To:          t.To,
		Value:       parsedValue,
	}, nil
}
