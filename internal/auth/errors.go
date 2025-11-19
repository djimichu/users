package auth

import "errors"

var (
	ErrMissHeader    = errors.New("missing auth header")
	ErrInvalidHeader = errors.New("invalid auth header")
)
