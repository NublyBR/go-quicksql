package quicksql

import "errors"

var (
	ErrInvalidIdentifier = errors.New("invalid identifier")
	ErrNil               = errors.New("function does not accept nil")
)
