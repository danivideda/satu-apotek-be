package handler

import (
	"errors"
)

var (
	ErrInvalidPassword  = errors.New("invalid password given")
	ErrInvalidAuthToken = errors.New("auth token is invalid")
)
