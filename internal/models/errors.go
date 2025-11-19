package models

import "errors"

var (
	ErrShortName = errors.New("your name too short")
	ErrLongName  = errors.New("your name too long")
	ErrShortPass = errors.New("your password too short")
	ErrLongPass  = errors.New("your password too long")
)
