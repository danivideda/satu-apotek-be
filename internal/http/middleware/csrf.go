package middleware

import (
	"errors"
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/service"
)

func (m *AppMiddleware) CSRFProtectionOwner(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// CSRF Token is not required on GET request
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		csrfToken := r.Header.Get("X-CSRF-Token")
		ownerSessionCookie, _ := r.Cookie("owner_session")
		err := service.VerifyCSRFToken(ownerSessionCookie.Value, csrfToken)
		if err != nil { 
			// If any error happen, immediately delete CSRF cookie so it can be refreshed in /check endpoint later
			service.DeleteOwnerCSRFCookie(w)
			if errors.Is(err, service.ErrMalformedCSRFToken) || errors.Is(err, service.ErrInvalidCSRFToken) {
				json.ResponseInvalidCSRFToken(w, r, err)
			} else {
				json.ResponseInternalServerError(w, r, err)
			}
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (m *AppMiddleware) CSRFProtectionUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// CSRF Token is not required on GET request
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		csrfToken := r.Header.Get("X-CSRF-Token")
		userSessionCookie, _ := r.Cookie("user_session")
		err := service.VerifyCSRFToken(userSessionCookie.Value, csrfToken)
		if err != nil {
			if errors.Is(err, service.ErrMalformedCSRFToken) || errors.Is(err, service.ErrInvalidCSRFToken) {
				json.ResponseInvalidCSRFToken(w, r, err)
			} else {
				json.ResponseInternalServerError(w, r, err)
			}
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
