package cookie

import (
	"fmt"
	"net/http"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/http/csrf"
)

type cookie struct {
	Owner    OwnerCookie
	User     UserCookie
	Pharmacy PharmacyCookie
}

type OwnerCookie interface {
	SetSession(w http.ResponseWriter, sessionID string, exp time.Time)
	DeleteSession(w http.ResponseWriter)
	DeleteCSRF(w http.ResponseWriter)
}

type UserCookie interface {
	SetSession(w http.ResponseWriter, sessionID string, exp time.Time)
	DeleteSession(w http.ResponseWriter)
	DeleteCSRF(w http.ResponseWriter)
}

type PharmacyCookie interface {
	SetSession(w http.ResponseWriter, sessionID string, exp time.Time)
}

type ownerCookie struct{}
type userCookie struct{}
type pharmacyCookie struct{}

func New() cookie {
	return cookie{
		Owner:    &ownerCookie{},
		User:     &userCookie{},
		Pharmacy: &pharmacyCookie{},
	}
}

// Owner sessions
func (c *ownerCookie) setCookiesBase(w http.ResponseWriter, sessionID string, csrfToken string, exp time.Time) {
	httpCookie := &http.Cookie{
		Name:     "owner_session",
		Value:    sessionID,
		Path:     "/",
		Expires:  exp,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, httpCookie)

	c.setCSRFCookiesBase(w, csrfToken, exp)
}

func (c *ownerCookie) setCSRFCookiesBase(w http.ResponseWriter, csrfToken string, exp time.Time) {
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

func (c *ownerCookie) SetSession(w http.ResponseWriter, sessionID string, exp time.Time) {
	csrfToken, _ := csrf.NewCSRFToken(sessionID)
	c.setCookiesBase(w, sessionID, csrfToken, exp)
}

func (c *ownerCookie) DeleteSession(w http.ResponseWriter) {
	c.setCookiesBase(w, "", "", time.Now())
}

func (c *ownerCookie) DeleteCSRF(w http.ResponseWriter) {
	c.setCSRFCookiesBase(w, "", time.Now())
}

func (c *pharmacyCookie) SetSession(w http.ResponseWriter, sessionID string, exp time.Time) {
	fmt.Println("Runs here")
	httpCookie := &http.Cookie{
		Name:     "pharmacy_session",
		Value:    sessionID,
		Path:     "/",
		Expires:  exp,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, httpCookie)
}

// User sessions
func (c *userCookie) setCookiesBase(w http.ResponseWriter, sessionID string, csrfToken string, exp time.Time) {
	httpCookie := &http.Cookie{
		Name:     "user_session",
		Value:    sessionID,
		Path:     "/",
		Expires:  exp,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, httpCookie)

	c.setCSRFCookiesBase(w, csrfToken, exp)
}

func (c *userCookie) setCSRFCookiesBase(w http.ResponseWriter, csrfToken string, exp time.Time) {
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

func (c *userCookie) SetSession(w http.ResponseWriter, sessionID string, exp time.Time) {
	csrfToken, _ := csrf.NewCSRFToken(sessionID)
	c.setCookiesBase(w, sessionID, csrfToken, exp)
}

func (c *userCookie) DeleteSession(w http.ResponseWriter) {
	c.setCookiesBase(w, "", "", time.Now())
}

func (c *userCookie) DeleteCSRF(w http.ResponseWriter) {
	c.setCSRFCookiesBase(w, "", time.Now())
}
