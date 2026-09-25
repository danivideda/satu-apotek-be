package service

import (
	"errors"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrTypeAssertionFailed = errors.New("type assertion failed")
	ErrInvalidSession      = errors.New("session invalid")
	ErrUserForbidden       = errors.New("user doesn't belong in current pharmacy session")
	ErrPharmacyCodeExpired = errors.New("code is expired")
)
