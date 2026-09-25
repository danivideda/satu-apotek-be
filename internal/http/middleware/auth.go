package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/danivideda/satu-apotek-be/internal/service"
)

const (
	authOwnerCtx    = "AuthOwnerCtx"
	authUserCtx     = "AuthUserCtx"
	authPharmacyCtx = "AuthPharmacyCtx"
)

type authOwner struct {
	ID         int64
	SessionID  string
	SessionExp time.Time
}

type authUser struct {
	repository.UserCacheValue
	SessionID  string
	SessionExp time.Time
}

type authPharmacy struct {
	ID    int64
	Name  string
	Users []repository.UserCacheValue
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
		sessionExp := sessionCookie.Expires

		// 2. Resolve Owner SessionID
		ownerSession, rotated, err := m.sessionSvc.ResolveOwner(ctx, sessionID)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrInvalidSession):
				cookie.Owner.DeleteSession(w)
				json.ResponseUnauthorized(w, r, err)
			default:
				json.ResponseInternalServerError(w, r, err)
			}
			return
		}

		// 2.1 Use the same SessionID and expiry if not rotated
		if !rotated {
			authOwner := authOwner{
				ID:         ownerSession.OwnerID,
				SessionID:  sessionID,
				SessionExp: sessionExp,
			}
			ctx = context.WithValue(ctx, authOwnerCtx, authOwner)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// 2.2 Otherwise, SessionID is rotated and need to be updated
		cookie.Owner.SetSession(w, ownerSession.ID, ownerSession.Exp)
		authOwner := authOwner{
			ID:         ownerSession.OwnerID,
			SessionID:  ownerSession.ID,
			SessionExp: ownerSession.Exp,
		}
		ctx = context.WithValue(ctx, authOwnerCtx, authOwner)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func (m *AppMiddleware) AuthUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// 1. Get session from cookie
		sessionCookie, err := r.Cookie("user_session")
		if err != nil {
			json.ResponseUnauthorized(w, r, err)
			return
		}
		sessionID := sessionCookie.Value
		sessionExp := sessionCookie.Expires

		// 2. Resolve User Session
		userSession, rotated, err := m.sessionSvc.ResolveUser(ctx, sessionID)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrInvalidCredentials):
				cookie.User.DeleteSession(w)
				json.ResponseUnauthorized(w, r, err)
			default:
				json.ResponseInternalServerError(w, r, err)
			}
			return
		}

		userCache := repository.UserCacheValue{
			ID:       userSession.UserID,
			Username: userSession.Username,
		}

		// 2.1 Use the same SessionID and expiry if not rotated
		if !rotated {
			authUser := authUser{
				UserCacheValue: userCache,
				SessionID:      sessionID,
				SessionExp:     sessionExp,
			}
			ctx = context.WithValue(ctx, authUserCtx, authUser)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// 2.2 Otherwise, SessionID is rotated and need to be updated
		cookie.User.SetSession(w, sessionID, userSession.Exp)
		authUser := authUser{
			UserCacheValue: userCache,
			SessionID:      userSession.ID,
			SessionExp:     userSession.Exp,
		}
		ctx = context.WithValue(ctx, authUserCtx, authUser)
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
			pharmacyCache, ok := val.(repository.PharmacyCacheValue)
			if !ok {
				json.ResponseInternalServerError(w, r, errors.New("type assertion failed, incorrect Pharmacy Session Cache Value form"))
				return
			}

			// Pass the Cache value to authPharmacyCtx
			authPharmacy := authPharmacy{
				ID:    pharmacyCache.PharmacyID,
				Name:  pharmacyCache.Name,
				Users: pharmacyCache.Users,
			}
			ctx := context.WithValue(ctx, authPharmacyCtx, authPharmacy)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		pharmacySession, err := m.repo.PharmacySessions.Update(ctx, sessionID, time.Now().Add(m.config.Auth.PharmacySessionTTL))
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				json.ResponseUnauthorized(w, r, err)
			} else {
				json.ResponseInternalServerError(w, r, err)
			}
			return
		}

		pharmacyDetail, err := m.repo.Pharmacies.GetByID(ctx, pharmacySession.PharmacyID)
		if err != nil {
			json.ResponseInternalServerError(w, r, err)
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
		m.repo.CacheStore.PharmacySessions.SetDefault(sessionID, repository.PharmacyCacheValue{
			PharmacyID: pharmacySession.PharmacyID,
			Name:       pharmacyDetail.Name,
			Users:      *usersCache,
		})
		cookie.Pharmacy.SetSession(w, sessionID, pharmacySession.ExpiresAt.Time)

		authPharmacy := authPharmacy{
			ID:    pharmacySession.PharmacyID,
			Name:  pharmacyDetail.Name,
			Users: *usersCache,
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

func AuthUserFromCtx(ctx context.Context) (*authUser, error) {
	authUser, ok := ctx.Value(authUserCtx).(authUser)
	if !ok {
		return nil, errors.New("AuthOwnerCtx type assertion missmatch")
	}
	return &authUser, nil
}

func AuthPharmacyFromCtx(ctx context.Context) (*authPharmacy, error) {
	authPharmacy, ok := ctx.Value(authPharmacyCtx).(authPharmacy)
	if !ok {
		return nil, errors.New("AuthPharmacyCtx type assertion missmatch")
	}
	return &authPharmacy, nil
}
