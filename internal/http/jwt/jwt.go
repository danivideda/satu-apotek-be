package jwt

import (
	"time"

	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/golang-jwt/jwt/v5"
)

type roleClaims string

const (
	RoleOwner roleClaims = "owner"
	RoleUser  roleClaims = "user"
)

var (
	refreshTokenKey = env.GetString("JWT_REFRESH_SECRET", "refresh-token-secret")
	accessTokenKey  = env.GetString("JWT_ACCESS_SECRET", "access-token-secret")
	refreshTokenTTL = env.GetString("REFRESH_TTL", "168h")
	accessTokenTTL  = env.GetString("ACCESS_TTL", "5m")
)

type MyClaim struct {
	ID        string     `json:"id"`
	Role      roleClaims `json:"role"`
	SessionID string     `json:"sid"`
	jwt.RegisteredClaims
}

func NewRefreshToken(id string, role roleClaims, sessionID string) (*time.Time, string, error) {
	ttl, err := time.ParseDuration(refreshTokenTTL)
	if err != nil {
		return nil, "", err
	}
	exp := time.Now().Add(ttl)
	signed, err := generate(id, role, sessionID, exp, refreshTokenKey)
	return &exp, signed, err
}

func NewAccessToken(id string, role roleClaims, sessionID string) (string, error) {
	ttl, err := time.ParseDuration(accessTokenTTL)
	if err != nil {
		return "", err
	}
	exp := time.Now().Add(ttl)
	return generate(id, role, sessionID, exp, accessTokenKey)
}

func generate(id string, role roleClaims, sessionID string, exp time.Time, signKey string) (string, error) {
	claims := MyClaim{
		ID:        id,
		Role:      role,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "satu-apotek-api",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(signKey))
}

func ValidateRefreshToken(t string) (*jwt.Token, *MyClaim, error) {
	return validate(t, refreshTokenKey)
}

func ValidateAccessToken(t string) (*jwt.Token, *MyClaim, error) {
	return validate(t, accessTokenKey)
}

func validate(t, signKey string) (*jwt.Token, *MyClaim, error) {
	token, err := jwt.ParseWithClaims(t, &MyClaim{}, func(t *jwt.Token) (any, error) {
		return []byte(signKey), nil
	})
	if err != nil {
		return nil, nil, err
	}

	if claims, ok := token.Claims.(*MyClaim); !ok {
		return nil, nil, ErrInvalidClaims
	} else {
		return token, claims, nil
	}
}