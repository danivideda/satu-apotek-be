package service

import (
	"net/http"
	"time"
)

func SetOwnerSessionCookie(w http.ResponseWriter, sessionID string, exp time.Time) {
	c := &http.Cookie{
		Name:    "owner_session",
		Value:   sessionID,
		Path:    "/",
		Expires: exp,
		// MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}