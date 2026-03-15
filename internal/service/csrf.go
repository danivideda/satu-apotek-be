package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/danivideda/satu-apotek-be/internal/env"
)

var CSRF_SECRET = env.GetString("CSRF_SECRET", "csrf-secret-key")

func NewCSRFToken(sessionID string) (string, error) {
	randomVal := rand.Text()
	message := fmt.Sprintf("%s!%s", sessionID, randomVal)
	hash := hmac.New(sha256.New, []byte(CSRF_SECRET))
	_, err := hash.Write([]byte(message))
	if err != nil {
		return "", err
	}
	newHmac := hash.Sum(nil)

	csrfToken := fmt.Sprintf("%s.%s", hex.EncodeToString([]byte(newHmac)), randomVal)
	return csrfToken, nil
}

func VerifyCSRFToken(sessionID, csrfToken string) (bool, error) {
	csrf := strings.Split(csrfToken, ".")
	if len(csrf) < 2 {
		return false, ErrMalformedCSRFToken
	}
	hmacFromRequest := csrf[0]
	randomVal := csrf[1]

	message := fmt.Sprintf("%s!%s", sessionID, randomVal)
	hash := hmac.New(sha256.New, []byte(CSRF_SECRET))
	_, err := hash.Write([]byte(message))
	if err != nil {
		return false, err
	}

	bytesArrayHmacFromRequest, err := hex.DecodeString(hmacFromRequest)
	if err != nil {
		return false, err
	}
	expectedHmac := hash.Sum(nil)
	isValid := hmac.Equal(bytesArrayHmacFromRequest, expectedHmac)
	return isValid, nil
}