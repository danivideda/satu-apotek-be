package service

import (
	"net/http"
	"time"
)

func setOwnerCookiesBase(w http.ResponseWriter, sessionID string, csrfToken string, exp time.Time) {
	c := &http.Cookie{
		Name:     "owner_session",
		Value:    sessionID,
		Path:     "/",
		Expires:  exp,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	setOwnerCSRFCookiesBase(w, csrfToken, exp)
}

func setOwnerCSRFCookiesBase(w http.ResponseWriter, csrfToken string, exp time.Time) {
	csrf := &http.Cookie{
		Name:     "owner_csrf",
		Value:    csrfToken,
		Path:     "/",
		Expires:  exp,
		Secure:   false,
		HttpOnly: false,
	}
	http.SetCookie(w, csrf)
}

func SetOwnerCookies(w http.ResponseWriter, sessionID string, exp time.Time) {
	csrfToken, _ := NewCSRFToken(sessionID)
	setOwnerCookiesBase(w, sessionID, csrfToken, exp)
}

func DeleteOwnerCookies(w http.ResponseWriter) {
	setOwnerCookiesBase(w, "", "", time.Now())
}

func DeleteOwnerCSRFCookie(w http.ResponseWriter) {
	setOwnerCSRFCookiesBase(w, "", time.Now())
}

func SetPharmacyCookies(w http.ResponseWriter, sessionID string, exp time.Time) {
	c := &http.Cookie{
		Name:     "pharmacy_session",
		Value:    sessionID,
		Path:     "/",
		Expires:  exp,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

func SetUserCookies(w http.ResponseWriter, sessionID string, exp time.Time) {
	c := &http.Cookie{
		Name:     "user_session",
		Value:    sessionID,
		Path:     "/",
		Expires:  exp,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	csrfToken, _ := NewCSRFToken(sessionID)
	csrf := &http.Cookie{
		Name:     "user_csrf",
		Value:    csrfToken,
		Path:     "/",
		Expires:  exp,
		Secure:   false,
		HttpOnly: false,
	}
	http.SetCookie(w, csrf)
}
