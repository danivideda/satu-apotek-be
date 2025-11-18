package jwt

import "errors"

var (
	ErrClaimsInvalid = errors.New("invalid claims type")
)