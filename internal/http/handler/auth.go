package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/jwt"
	"github.com/danivideda/satu-apotek-be/internal/store"
	"github.com/jackc/pgx/v5"
)

type authHandler struct {
	store store.Storage
}

type authToken struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

const (
	ownerCookieName = "owner_auth_token"
	userCookieName  = "user_auth_token"

	ownerCookiePath = "/v1/auth/owners/refresh"
	userCookiePath  = "/v1/auth/users/refresh"

	ownerAccessTokenName = "owner_access_token"
	userAccessTokenName = "user_access_token"
)

func refresh(refreshToken *http.Cookie, accessTokenName string, h *authHandler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, err := jwt.ValidateRefreshToken(refreshToken.Value)
	if err != nil {
		json.ResponseUnauthorized(w, r, err)
		return
	}

	if err := validateSession(ctx, h.store, claims.SessionID); err != nil {
		if errors.Is(err, ErrRevokedAuthToken) {
			json.ResponseUnauthorized(w, r, err)
			return
		} else {
			json.ResponseInternalServerError(w, r, err)
			return
		}
	}

	accessToken, err := jwt.NewAccessToken(claims.ID, claims.Role, claims.SessionID)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := map[string]string{
		accessTokenName: accessToken,
	}
	if err := json.WriteResponse(w, http.StatusOK, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func hashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func verifyPassword(password string, hash []byte) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, string(hash))
}

func newAuthToken(id string, role jwt.RoleClaims) (*authToken, *time.Time, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, nil, err
	}

	exp, refreshToken, err := jwt.NewRefreshToken(id, role, sessionID)
	if err != nil {
		return nil, nil, err
	}
	accessToken, err := jwt.NewAccessToken(id, role, sessionID)
	if err != nil {
		return nil, nil, err
	}

	token := authToken{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}

	return &token, exp, nil
}

func generateSessionID() (string, error) {
	b := make([]byte, 32) // 32 bytes for a 256-bit ID
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func setCookies(w http.ResponseWriter, refreshToken string, exp time.Time, role jwt.RoleClaims) error {
	var cookieName string
	switch role {
	case jwt.RoleOwner:
		cookieName = ownerCookieName
	case jwt.RoleUser:
		cookieName = userCookieName
	default:
		return ErrWrongRole
	}

	var cookiePath string
	switch role {
	case jwt.RoleOwner:
		cookiePath = ownerCookiePath
	case jwt.RoleUser:
		cookiePath = userCookiePath
	default:
		return ErrWrongRole
	}

	authCookie := http.Cookie{
		Name:     cookieName,
		Value:    refreshToken,
		Path:     cookiePath,
		Expires:  exp,
		Secure:   false, // Set to true for HTTPS
		HttpOnly: true,  // Prevent client-side script access
	}
	http.SetCookie(w, &authCookie)

	return nil
}

func validateSession(ctx context.Context, store store.Storage, sessionID string) error {
	_, err := store.RevokedSessions.GetBySessionID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// session id not found, so the session is still valid
			return nil
		} else {
			return err
		}
	}
	// session id found, so the session is already revoked
	return ErrRevokedAuthToken
}
