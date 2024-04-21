package types

import (
	"math/big"
	"time"
)

type Block struct {
	Number       uint64
	Hash         string
	ParentHash   string
	Timestamp    time.Time
	Transactions []Transaction
}

type Transaction struct {
	BlockHash   string
	BlockNumber uint64
	Hash        string
	From        string
	To          string
	Value       big.Int // ideally a decimal.Decimal but I cannot use external libraries for this exercise.
	// ... other fields omitted for the scope of this exercise
}
