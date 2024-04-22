package ethereum

import (
	"errors"
	"math/big"
	"strconv"
)

var ErrInvalidHexNumber = errors.New("invalid hex number")

func HasEthNumberPrefix(s string) bool {
	// When encoding quantities (integers, numbers): encode as hex, prefix with "0x",
	// the most compact representation (slight exception: zero should be represented as "0x0").
	// More info: https://ethereum.org/en/developers/docs/apis/json-rpc/#quantities-encoding

	if len(s) < 2 {
		return false
	}

	if s[0] != '0' || s[1] != 'x' {
		return false
	}

	return true
}

func BigIntFromEthNumber(s string) (big.Int, error) {
	if !HasEthNumberPrefix(s) {
		return big.Int{}, ErrInvalidHexNumber
	}

	var i big.Int
	if _, ok := i.SetString(s, 0); !ok {
		return big.Int{}, errors.New("could not decode hex number to string")
	}

	return i, nil
}

func Uint64FromEthNumber(s string) (uint64, error) {
	// When decoding quantities (integers, numbers): decode hex, skip prefix "0x".

	if !HasEthNumberPrefix(s) {
		return 0, ErrInvalidHexNumber
	}

	n, err := BigIntFromEthNumber(s)
	if err != nil {
		return 0, err
	}

	return n.Uint64(), nil
}

func EthNumberFromUnit64(n uint64) string {
	if n == 0 {
		return "0x0"
	}

	// Convert the int64 to a hexadecimal string.
	return "0x" + strconv.FormatUint(n, 16)
}
