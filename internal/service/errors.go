package service

import (
	"errors"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrTypeAssertionFailed = errors.New("type assertion failed")
	ErrInvalidSession      = errors.New("session invalid")
)
