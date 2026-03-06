package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/danivideda/satu-apotek-be/internal/service"
)

const AuthSessionOwnerCtx = "AuthOwnerSessionCtx"

var (
	ownerSessionTTL = env.GetString("OWNER_SESSION_TTL", "168h")
)

func (m *AppMiddleware) AuthSessionOwner(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sessionCookie, err := r.Cookie("owner_session")
		if err != nil {
			json.ResponseUnauthorized(w, r, err)
			return
		}
		sessionID := sessionCookie.Value

		// 2. Check if session exist in cache. If exist, pass the request.
		if val, found := m.repo.CacheStore.OwnerSessions.Get(sessionID); found {
			fmt.Println("Session value:", sessionID)
			ownerID := val
			ctx := context.WithValue(ctx, AuthSessionOwnerCtx, ownerID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// 3. Check if session exist in DB. If exist, renew the session_id and expires_at value. If not exist
		// or expired, then the session is invalid
		ttl, err := time.ParseDuration(ownerSessionTTL)
		if err != nil {
			json.ResponseInternalServerError(w, r, err)
			return
		}
		ownerSession, err := m.repo.OwnerSessions.Update(ctx, sessionID, time.Now().Add(ttl))
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				json.ResponseUnauthorized(w, r, err)
			} else {
				json.ResponseInternalServerError(w, r, err)
			}
			return
		}
		sessionID = ownerSession.ID.String()
		m.repo.CacheStore.OwnerSessions.SetDefault(sessionID, ownerSession.OwnerID)
		service.SetOwnerSessionCookie(w, sessionID, ownerSession.ExpiresAt.Time)

		fmt.Println("Session value:", sessionID)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
