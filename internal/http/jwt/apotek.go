package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ApotekClaims struct {
	ApotekID int64 `json:"apotek_id"`
	jwt.RegisteredClaims
}

func NewApotekToken(apotekID int64) (string, error) {
	exp := time.Now().Add(24 * 7 * time.Hour)
	claims := ApotekClaims{
		ApotekID: apotekID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "satu-apotek-api",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(refreshTokenKey))
}
