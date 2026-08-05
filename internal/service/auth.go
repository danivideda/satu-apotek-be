package service

import (
	"net/http"
	"time"
)

// Owner sessions
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

// User sessions
func setUserCookiesBase(w http.ResponseWriter, sessionID string, csrfToken string, exp time.Time) {
	c := &http.Cookie{
		Name:     "user_session",
		Value:    sessionID,
		Path:     "/",
		Expires:  exp,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	setUserCSRFCookiesBase(w, csrfToken, exp)
}

func setUserCSRFCookiesBase(w http.ResponseWriter, csrfToken string, exp time.Time) {
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

func SetUserCookies(w http.ResponseWriter, sessionID string, exp time.Time) {
	csrfToken, _ := NewCSRFToken(sessionID)
	setUserCookiesBase(w, sessionID, csrfToken, exp)
}

func DeleteUserCookies(w http.ResponseWriter) {
	setUserCookiesBase(w, "", "", time.Now())
}

func DeleteUserCSRFCookie(w http.ResponseWriter) {
	setUserCSRFCookiesBase(w, "", time.Now())
}