package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/danivideda/satu-apotek-be/internal/http/auth"
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/danivideda/satu-apotek-be/internal/http/response"
	"github.com/danivideda/satu-apotek-be/internal/store"
)

type authHandler struct {
	store store.Storage
}

type authToken struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func (h *authHandler) RegisterOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type registerOwnerPayload struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var payload registerOwnerPayload
	if err := json.Read(w, r, &payload); err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	// CreateHash returns an Argon2id hash of a plain-text password using the
	// provided algorithm parameters. The returned hash follows the format used
	// by the Argon2 reference C implementation and looks like this:
	// $argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
	hashedPassword, err := argon2id.CreateHash(payload.Password, argon2id.DefaultParams)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	param := dbsqlc.CreateOwnerParams{
		Username:     payload.Username,
		Email:        payload.Email,
		PasswordHash: []byte(hashedPassword),
	}
	owner, err := h.store.Owners.CreateOwner(ctx, param)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	token, exp, err := newAuthToken(owner.ID)
	if err != nil {
		response.InternalServerErrorResponse(w, r, err)
		return
	}

	setCookies(w, token.RefreshToken, *exp)

	res := struct {
		*dbsqlc.CreateOwnerRow
		*authToken
	}{
		CreateOwnerRow: owner,
		authToken:      token,
	}

	if err := json.WriteResponse(w, http.StatusCreated, res); err != nil {
		response.InternalServerErrorResponse(w, r, err)
		return
	}
}

func (h *authHandler) LoginOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type loginOwnerPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var payload loginOwnerPayload
	if err := json.Read(w, r, &payload); err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	owner, err := h.store.Owners.GetOwnerByUsername(ctx, payload.Username)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	// ComparePasswordAndHash performs a constant-time comparison between a
	// plain-text password and Argon2id hash, using the parameters and salt
	// contained in the hash. It returns true if they match, otherwise it returns
	// false.
	match, err := argon2id.ComparePasswordAndHash(payload.Password, string(owner.PasswordHash))
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	if !match {
		response.BadRequestResponse(w, r, ErrInvalidPassword)
		return
	}

	token, exp, err := newAuthToken(owner.ID)
	if err != nil {
		response.InternalServerErrorResponse(w, r, err)
		return
	}

	setCookies(w, token.RefreshToken, *exp)

	res := struct {
		Username    string `json:"username"`
		Message     string `json:"message"`
		AccessToken string `json:"access_token"`
	}{
		Username:    payload.Username,
		Message:     "Login success",
		AccessToken: token.AccessToken,
	}

	if err := json.WriteResponse(w, http.StatusOK, res); err != nil {
		response.InternalServerErrorResponse(w, r, err)
		return
	}
}

func newAuthToken(id int32) (*authToken, *time.Time, error) {
	idString := strconv.Itoa(int(id))
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, nil, err
	}

	exp, refreshToken, err := auth.NewRefreshToken(idString, auth.RoleOwner, sessionID)
	if err != nil {
		return nil, nil, err
	}
	accessToken, err := auth.NewAccessToken(idString, auth.RoleOwner, sessionID)
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

func setCookies(w http.ResponseWriter, refreshToken string, exp time.Time) {
	authCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/v1/auth/refresh",
		Expires:  exp,
		Secure:   false, // Set to true for HTTPS
		HttpOnly: true,  // Prevent client-side script access
	}
	http.SetCookie(w, &authCookie)
}

func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	claims, err := auth.ValidateRefreshToken(refreshToken.Value)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	accessToken, err := auth.NewAccessToken(claims.ID, auth.RoleOwner, claims.SessionID)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	res := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: accessToken,
	}
	if err := json.WriteResponse(w, http.StatusOK, res); err != nil {
		response.InternalServerErrorResponse(w, r, err)
		return
	}
}

func (h *authHandler) LogoutOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	claims, ok := ctx.Value(middleware.AuthClaimsCtx).(*auth.AuthClaims)
	if !ok {
		response.InternalServerErrorResponse(w, r, ErrInvalidAuthToken)
		return
	}

	if err := json.WriteResponse(w, http.StatusOK, claims.SessionID); err != nil {
		response.InternalServerErrorResponse(w, r, err)
		return
	}
}
