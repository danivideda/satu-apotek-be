package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/danivideda/satu-apotek-be/internal/service"
)

const authOwnerCtx = "AuthOwnerSessionCtx"

type authOwner struct {
	ID        int64
	SessionID string
}

func AuthOwnerFromCtx(ctx context.Context) (*authOwner, error) {
	authOwner, ok := ctx.Value(authOwnerCtx).(authOwner)
	if !ok {
		return nil, errors.New("AuthOwnerCtx type assertion missmatch")
	}
	return &authOwner, nil
}

var (
	ownerSessionTTL = env.GetString("OWNER_SESSION_TTL", "168h")
)

func (m *AppMiddleware) AuthOwner(next http.Handler) http.Handler {
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
			ownerID, ok := val.(int64)
			if !ok {
				json.ResponseInternalServerError(w, r, errors.New("incorrect type assertion of OwnerID"))
				return
			}

			authOwner := authOwner{
				ID:        ownerID,
				SessionID: sessionID,
			}
			ctx := context.WithValue(ctx, authOwnerCtx, authOwner)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// 3. Check if session exist in DB. If exist, renew the session_id and expires_at value. If not exist
		// or is expired, then the session is invalid
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
		service.SetOwnerCookies(w, sessionID, ownerSession.ExpiresAt.Time)

		authOwner := authOwner{
			ID:        ownerSession.OwnerID,
			SessionID: sessionID,
		}
		ctx = context.WithValue(ctx, authOwnerCtx, authOwner)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
