package service

import (
	"net/http"
	"time"
)

func SetOwnerCookies(w http.ResponseWriter, sessionID string, exp time.Time) {
	c := &http.Cookie{
		Name:    "owner_session",
		Value:   sessionID,
		Path:    "/",
		Expires: exp,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	csrfToken, _ := NewCSRFToken(sessionID)
	csrf := &http.Cookie{
		Name:    "owner_csrf",
		Value:   csrfToken,
		Path:    "/",
		Expires: exp,
		Secure:   false,
		HttpOnly: false,
	}
	http.SetCookie(w, csrf)
}

func DeleteOwnerCookies(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:    "owner_session",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(1 * time.Second),
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	csrf := &http.Cookie{
		Name:    "owner_csrf",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(1 * time.Second),
		Secure:   false,
		HttpOnly: false,
	}
	http.SetCookie(w, csrf)
}

func SetPharmacyCookies(w http.ResponseWriter, sessionID string, exp time.Time) {
	c := &http.Cookie{
		Name:    "pharmacy_session",
		Value:   sessionID,
		Path:    "/",
		Expires: exp,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}