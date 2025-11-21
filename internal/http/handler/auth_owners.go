package handler

import (
	"errors"
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/jwt"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
)

func (h *authHandler) OwnerRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type registerOwnerPayload struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var payload registerOwnerPayload
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	passwordHash, err := hashPassword(payload.Password)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	owner, err := h.repo.Owners.Create(ctx, payload.Username, payload.Email, passwordHash)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	token, exp, err := newAuthToken(owner.ID.String(), jwt.RoleOwner)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	if err := setCookies(w, token.RefreshToken, *exp, jwt.RoleOwner); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := map[string]string{
		"username": owner.Username,
		ownerAccessTokenName: token.AccessToken,
	}

	if err := json.WriteResponse(w, http.StatusCreated, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *authHandler) OwnerLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type loginOwnerPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var payload loginOwnerPayload
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	owner, err := h.repo.Owners.GetByUsername(ctx, payload.Username)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	match, err := verifyPassword(payload.Password, owner.PasswordHash)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	if !match {
		json.ResponseBadRequest(w, r, ErrInvalidPassword)
		return
	}

	token, exp, err := newAuthToken(owner.ID.String(), jwt.RoleOwner)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	if err := setCookies(w, token.RefreshToken, *exp, jwt.RoleOwner); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := map[string]string{
		"message": "Login success",
		ownerAccessTokenName: token.AccessToken,
	}

	if err := json.WriteResponse(w, http.StatusOK, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *authHandler) OwnerLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims := middleware.AuthClaimsFromContext(ctx)
	if claims == nil {
		json.ResponseInternalServerError(w, r, ErrInvalidAuthToken)
		return
	}

	if err := validateSession(ctx, h.repo, claims.SessionID); err != nil {
		if errors.Is(err, ErrRevokedAuthToken) {
			json.ResponseUnauthorized(w, r, err)
			return
		} else {
			json.ResponseInternalServerError(w, r, err)
			return
		}
	}
	revokedSession, err := h.repo.RevokedSessions.Create(ctx, claims.SessionID)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := map[string]string{
		"revoked_session_id": revokedSession.SessionID,
	}

	if err := json.ResponseOK(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *authHandler) OwnerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie(ownerCookieName)
	if err != nil {
		json.ResponseUnauthorized(w, r, err)
		return
	}

	refresh(refreshToken, ownerAccessTokenName, h, w, r)
}
