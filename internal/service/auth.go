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
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

func DeleteOwnerSessionCookie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:    "owner_session",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(1 * time.Second),
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}