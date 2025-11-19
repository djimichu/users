package user

import "errors"

var (
	ErrHashPass      = errors.New("failed hashing password")
	ErrIncorrectPass = errors.New("incorrect password")
)
