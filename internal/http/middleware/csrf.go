package middleware

import (
	"errors"
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/service"
)

func (m *AppMiddleware) OwnerCSRFProtection(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// CSRF Token is not required on GET request
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		csrfToken := r.Header.Get("X-CSRF-Token")
		ownerSessionCookie, _ := r.Cookie("owner_session")
		ok, err := service.VerifyCSRFToken(ownerSessionCookie.Value, csrfToken)
		if err != nil {
			if errors.Is(err, service.ErrMalformedCSRFToken) {
				json.ResponseInvalidCSRFToken(w, r, err)
			} else {
				json.ResponseInternalServerError(w, r, err)
			}
			return
		}
		if !ok {
			json.ResponseInvalidCSRFToken(w, r, service.ErrInvalidCSRFToken)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
