package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type myClaim struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func runScript2() {
	claims := myClaim{
		Username: "danivideda",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "satu-apotek-api",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Second)),
		},
	}
	key := "secret_key"
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims = claims
	s, err := t.SignedString([]byte(key))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)
}
