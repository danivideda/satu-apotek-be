package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
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

	csrfToken := fmt.Sprintf("%s.%s", base64.StdEncoding.EncodeToString([]byte(newHmac)), randomVal)
	return csrfToken, nil
}

func VerifyCSRFToken(sessionID, csrfToken string) error {
	csrf := strings.Split(csrfToken, ".")
	if len(csrf) < 2 {
		return ErrMalformedCSRFToken
	}
	hmacFromRequest := csrf[0]
	randomVal := csrf[1]

	message := fmt.Sprintf("%s!%s", sessionID, randomVal)
	hash := hmac.New(sha256.New, []byte(CSRF_SECRET))
	_, err := hash.Write([]byte(message))
	if err != nil {
		return err
	}
	expectedHmac := hash.Sum(nil)

	bytesArrayHmacFromRequest, err := base64.StdEncoding.DecodeString(hmacFromRequest)
	if err != nil {
		return ErrInvalidCSRFToken
	}
	isValid := hmac.Equal(bytesArrayHmacFromRequest, expectedHmac)
	if isValid {
		return nil
	} else {
		return ErrInvalidCSRFToken
	}
}