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

const (
	authOwnerCtx    = "AuthOwnerCtx"
	authPharmacyCtx = "AuthPharmacyCtx"
)

var (
	ownerSessionTTL    = env.GetString("OWNER_SESSION_TTL", "168h")
	pharmacySessionTTL = env.GetString("PHARMACY_SESSION_TTL", "168h")
)

type authOwner struct {
	ID        int64
	SessionID string
}

type authPharmacy struct {
	ID        int64
	SessionID string
	Users     []repository.UserCache
}

func (m *AppMiddleware) AuthOwner(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// 1. Get session from cookie
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
				json.ResponseInternalServerError(w, r, errors.New("type assertion failed, ownerID is not int64"))
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

func (m *AppMiddleware) AuthPharmacy(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sessionCookie, err := r.Cookie("pharmacy_session")
		if err != nil {
			json.ResponseUnauthorized(w, r, fmt.Errorf("%w: no pharmacy session", err))
			return
		}
		sessionID := sessionCookie.Value

		if val, found := m.repo.CacheStore.PharmacySessions.Get(sessionID); found {
			pharmacySessionCache, ok := val.(repository.PharmacySessionCacheValue)
			if !ok {
				json.ResponseInternalServerError(w, r, errors.New("type assertion failed, incorrect Pharmacy Session Cache Value form"))
				return
			}

			// Pass the Cache value to authPharmacyCtx
			authPharmacy := authPharmacy{
				ID:        pharmacySessionCache.PharmacyID,
				SessionID: sessionID,
				Users:     pharmacySessionCache.Users,
			}
			ctx := context.WithValue(ctx, authPharmacyCtx, authPharmacy)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		ttl, err := time.ParseDuration(pharmacySessionTTL)
		if err != nil {
			json.ResponseInternalServerError(w, r, err)
			return
		}

		pharmacySession, err := m.repo.PharmacySessions.Update(ctx, sessionID, time.Now().Add(ttl))
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				json.ResponseUnauthorized(w, r, err)
			} else {
				json.ResponseInternalServerError(w, r, err)
			}
			return
		}

		// get Users[] that's associated with PharmacyID
		usersCache, err := service.GetUsersFromPharmacyID(ctx, m.repo, pharmacySession.PharmacyID)
		if err != nil {
			json.ResponseInternalServerError(w, r, err)
			return
		}

		// update sessionID
		sessionID = pharmacySession.ID.String()
		m.repo.CacheStore.PharmacySessions.SetDefault(sessionID, repository.PharmacySessionCacheValue{
			PharmacyID: pharmacySession.PharmacyID,
			Users:      *usersCache,
		})
		service.SetPharmacyCookies(w, sessionID, pharmacySession.ExpiresAt.Time)

		authPharmacy := authPharmacy{
			ID:        pharmacySession.PharmacyID,
			SessionID: sessionID,
			Users:     *usersCache,
		}
		ctx = context.WithValue(ctx, authPharmacyCtx, authPharmacy)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func AuthOwnerFromCtx(ctx context.Context) (*authOwner, error) {
	authOwner, ok := ctx.Value(authOwnerCtx).(authOwner)
	if !ok {
		return nil, errors.New("AuthOwnerCtx type assertion missmatch")
	}
	return &authOwner, nil
}

func AuthPharmacyFromCtx(ctx context.Context) (*authPharmacy, error) {
	authPharmacy, ok := ctx.Value(authPharmacyCtx).(authPharmacy)
	if !ok {
		return nil, errors.New("AuthPharmacyCtx type assertion missmatch")
	}
	return &authPharmacy, nil
}
