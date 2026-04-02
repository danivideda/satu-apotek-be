package handler

import (
	"errors"
)

var (
	ErrInvalidPassword  = errors.New("invalid password given")
	ErrInvalidAuthToken = errors.New("invalid auth token")
	ErrRevokedAuthToken = errors.New("auth session already revoked")
	ErrWrongRole = errors.New("role provided is wrong")

	ErrAppIDNotAllowed = errors.New("appID not allowed for current auth'd owner")
)
