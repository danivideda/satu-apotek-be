package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

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

	ownerID, ownerSessionID, err := h.repo.Owners.Create(ctx, payload.Username, payload.Email, passwordHash)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	token, exp, err := newAuthToken(strconv.Itoa(int(ownerID)), jwt.RoleOwner)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	if err := setCookies(w, token.RefreshToken, *exp, jwt.RoleOwner); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := map[string]string{
		"username": "temp-username-for-refactoring",
		ownerAccessTokenName: token.AccessToken,
		"owner_session_id": ownerSessionID,
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

	token, exp, err := newAuthToken(strconv.Itoa(int(owner.ID)), jwt.RoleOwner)
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
	idInt, _ := strconv.Atoi(claims.ID)
	ttl, _ := time.ParseDuration(ownerSessionTTL)
	ownerSession, err := h.repo.OwnerSessions.Create(ctx, int64(idInt), time.Now().Add(ttl))
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := map[string]string{
		"session_id": ownerSession.ID.String(),
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
