package repository

import "errors"

var (
	ErrNotFoundUser = errors.New("user with this ID not exist")
)
