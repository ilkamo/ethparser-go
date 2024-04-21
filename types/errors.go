package types

import "errors"

var (
	ErrAddressNotFound = errors.New("address not found")
	ErrAlreadyRunning  = errors.New("parser is already running")
)
