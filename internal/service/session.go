package service

import (
	"net/http"
	"time"
)

type Session struct {
	Owners     OwnerSessionService
	Users      UserSessionService
	Pharmacies PharmacySessionService
}

type OwnerSessionService interface {
	SetCookies(w http.ResponseWriter, sessionID string, exp time.Time)
	DeleteCookies(w http.ResponseWriter)
	DeleteCSRFCookie(w http.ResponseWriter)
}

type UserSessionService interface {
	SetCookies(w http.ResponseWriter, sessionID string, exp time.Time)
	DeleteCookies(w http.ResponseWriter)
	DeleteCSRFCookie(w http.ResponseWriter)
}

type PharmacySessionService interface {
	SetCookies(w http.ResponseWriter, sessionID string, exp time.Time)
}

type ownerSessionService struct{}
type userSessionService struct{}
type pharmacySessionService struct{}

func NewSessionService() Session {
	return Session{
		Owners:     &ownerSessionService{},
		Users:      &userSessionService{},
		Pharmacies: &pharmacySessionService{},
	}
}

// Owner sessions
func (s *ownerSessionService) setCookiesBase(w http.ResponseWriter, sessionID string, csrfToken string, exp time.Time) {
	c := &http.Cookie{
		Name:     "owner_session",
		Value:    sessionID,
		Path:     "/",
		Expires:  exp,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	s.setCSRFCookiesBase(w, csrfToken, exp)
}

func (s *ownerSessionService) setCSRFCookiesBase(w http.ResponseWriter, csrfToken string, exp time.Time) {
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

func (s *ownerSessionService) SetCookies(w http.ResponseWriter, sessionID string, exp time.Time) {
	csrfToken, _ := NewCSRFToken(sessionID)
	s.setCookiesBase(w, sessionID, csrfToken, exp)
}

func (s *ownerSessionService) DeleteCookies(w http.ResponseWriter) {
	s.setCookiesBase(w, "", "", time.Now())
}

func (s *ownerSessionService) DeleteCSRFCookie(w http.ResponseWriter) {
	s.setCSRFCookiesBase(w, "", time.Now())
}

func (s *pharmacySessionService) SetCookies(w http.ResponseWriter, sessionID string, exp time.Time) {
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
func (s *userSessionService) setCookiesBase(w http.ResponseWriter, sessionID string, csrfToken string, exp time.Time) {
	c := &http.Cookie{
		Name:     "user_session",
		Value:    sessionID,
		Path:     "/",
		Expires:  exp,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	s.setCSRFCookiesBase(w, csrfToken, exp)
}

func (s *userSessionService) setCSRFCookiesBase(w http.ResponseWriter, csrfToken string, exp time.Time) {
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

func (s *userSessionService) SetCookies(w http.ResponseWriter, sessionID string, exp time.Time) {
	csrfToken, _ := NewCSRFToken(sessionID)
	s.setCookiesBase(w, sessionID, csrfToken, exp)
}

func (s *userSessionService) DeleteCookies(w http.ResponseWriter) {
	s.setCookiesBase(w, "", "", time.Now())
}

func (s *userSessionService) DeleteCSRFCookie(w http.ResponseWriter) {
	s.setCSRFCookiesBase(w, "", time.Now())
}
