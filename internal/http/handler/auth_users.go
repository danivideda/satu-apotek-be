package handler

import (
	"errors"
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/jwt"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
)

// TODO user login handler
func (h *authHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
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

	owner, err := h.store.Owners.GetByUsername(ctx, payload.Username)
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

	if err := setCookies(w, token.RefreshToken, *exp, jwt.RoleUser); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := struct {
		Message     string `json:"message"`
		AccessToken string `json:"access_token"`
	}{
		Message:     "Login success",
		AccessToken: token.AccessToken,
	}

	if err := json.WriteResponse(w, http.StatusOK, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

// TODO user logout handler
func (h *authHandler) UserLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims := middleware.AuthClaimsFromContext(ctx)
	if claims == nil {
		json.ResponseInternalServerError(w, r, ErrInvalidAuthToken)
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
	revokedSession, err := h.store.RevokedSessions.Create(ctx, claims.SessionID)
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

// TODO user refresh handler
func (h *authHandler) UserRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie(userCookieName)
	if err != nil {
		json.ResponseUnauthorized(w, r, err)
		return
	}

	refresh(refreshToken, userAccessTokenName, h, w, r)
}