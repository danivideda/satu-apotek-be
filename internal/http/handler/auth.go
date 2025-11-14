package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/jwt"
)

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
		badRequestResponse(w, r, err)
		return
	}

	// CreateHash returns an Argon2id hash of a plain-text password using the
	// provided algorithm parameters. The returned hash follows the format used
	// by the Argon2 reference C implementation and looks like this:
	// $argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
	hashedPassword, err := argon2id.CreateHash(payload.Password, argon2id.DefaultParams)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	param := dbsqlc.CreateOwnerParams{
		Username: payload.Username,
		Email:    payload.Email,
		PasswordHash: []byte(hashedPassword),
	}
	owner, err := h.store.Owners.CreateOwner(ctx, param)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	token, exp, err := newAuthToken(owner.ID)
	if err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}

	authCookie := http.Cookie{
		Name:    "refresh_token",
		Value:   token.RefreshToken,
		Path:    "/v1/refresh",
		Domain:  "localhost", // Or your actual domain
		Expires: *exp,
		Secure:  false, // Set to true for HTTPS
		HttpOnly: true, // Prevent client-side script access
	}
	http.SetCookie(w, &authCookie)

	response := struct {
		*dbsqlc.CreateOwnerRow
		*authToken
	}{
		CreateOwnerRow: owner,
		authToken:      token,
	}

	if err := json.WriteResponse(w, http.StatusCreated, response); err != nil {
		internalServerErrorResponse(w, r, err)
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
		badRequestResponse(w, r, err)
		return
	}

	owner, err := h.store.Owners.GetOwnerByUsername(ctx, payload.Username)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	// ComparePasswordAndHash performs a constant-time comparison between a
	// plain-text password and Argon2id hash, using the parameters and salt
	// contained in the hash. It returns true if they match, otherwise it returns
	// false.
	match, err := argon2id.ComparePasswordAndHash(payload.Password, string(owner.PasswordHash))
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if !match {
		badRequestResponse(w, r, ErrInvalidPassword)
		return
	}

	token, exp, err := newAuthToken(owner.ID)
	if err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}

	authCookie := http.Cookie{
		Name:    "refresh_token",
		Value:   token.RefreshToken,
		Path:    "/v1/refresh",
		Domain:  "localhost", // Or your actual domain
		Expires: *exp,
		Secure:  false, // Set to true for HTTPS
		HttpOnly: true, // Prevent client-side script access
	}
	http.SetCookie(w, &authCookie)

	response := struct {
		Username string `json:"username"`
		Message  string `json:"message"`
		*authToken
	}{
		Username:  payload.Username,
		Message:   "Login success",
		authToken: token,
	}

	if err := json.WriteResponse(w, http.StatusOK, response); err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}
}

func newAuthToken(id int32) (*authToken, *time.Time, error) {
	ownerID := strconv.Itoa(int(id))
	exp, refreshToken, err := jwt.NewRefreshToken(ownerID)
	if err != nil {
		return nil, nil, err
	}
	accessToken, err := jwt.NewAccessToken(ownerID)
	if err != nil {
		return nil, nil, err
	}

	token := authToken{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}

	return &token, exp, nil
}

func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := json.WriteResponse(w, http.StatusOK, refreshToken.String()); err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}
}