package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/danivideda/satu-apotek-be/internal/service"
)

type authHandler struct {
	repo            repository.Repository
	authService     *service.Auth
	ownerSessionTTL time.Duration
	userSessionTTL  time.Duration
}

func newAuthHandler(repo repository.Repository, authService *service.Auth, ownerSessionTTL time.Duration, userSessionTTL time.Duration) *authHandler {
	return &authHandler{repo, authService, ownerSessionTTL, userSessionTTL}
}

func (h *authHandler) OwnerRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload of username, email, and password
	var payload struct {
		Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6,max=64,excludesall=\t\n "`
	}
	if ok := parseAndValidateJSONPayload(w, r, &payload); !ok {
		return
	}

	// Handle owner register
	if err := h.authService.Owner.Register(ctx, payload.Username, payload.Email, payload.Password); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := map[string]any{
		"username": payload.Username,
	}

	if err := json.ResponseCreated(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *authHandler) OwnerLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload of username and password
	var payload struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6,max=64,excludesall=\t\n "`
	}
	if ok := parseAndValidateJSONPayload(w, r, &payload); !ok {
		return
	}

	// Handle Owner Login
	ownerSession, err := h.authService.Owner.Login(ctx, payload.Email, payload.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			json.ResponseBadRequest(w, r, err)
		default:
			json.ResponseInternalServerError(w, r, err)
		}
		return
	}
	cookie.Owner.SetSession(w, ownerSession.ID, ownerSession.Exp)

	if err := json.ResponseNoContent(w); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *authHandler) OwnerLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authOwner, err := middleware.AuthOwnerFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	err = h.authService.Owner.Logout(ctx, authOwner.SessionID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidSession):
			json.ResponseUnauthorized(w, r, err)
		default:
			json.ResponseInternalServerError(w, r, err)
		}
		cookie.Owner.DeleteSession(w)
		return
	}

	res := map[string]string{
		"deleted_session": authOwner.SessionID,
	}
	if err := json.ResponseOK(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *authHandler) OwnerCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authOwner, err := middleware.AuthOwnerFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	// Resend CSRF Cookie if it's missing.
	// When CSRF cookie is missing, it means the CSRF protection middleware unable to validate previous request
	// and deletes the CSRF cookie completely.
	_, err = r.Cookie("owner_csrf")
	if err != nil {
		fmt.Println(err)
		cookie.Owner.SetSession(w, authOwner.SessionID, authOwner.SessionExp)
	}

	if err := json.ResponseNoContent(w); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *authHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload of userID and password
	var payload struct {
		UserID   int64  `json:"user_id"`
		Password string `json:"password"`
	}
	if ok := parseAndValidateJSONPayload(w, r, &payload); !ok {
		return
	}

	// Handle User Login
	authPharmacy, err := middleware.AuthPharmacyFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
	newUserSession, err := h.authService.User.Login(ctx, payload.UserID, payload.Password, authPharmacy.Users)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserForbidden):
			json.ResponseForbidden(w, r, err)
		default:
			json.ResponseInternalServerError(w, r, err)
		}
		return
	}

	cookie.User.SetSession(w, newUserSession.ID, newUserSession.Exp)

	if err := json.ResponseNoContent(w); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *authHandler) UserCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, err := middleware.AuthUserFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	if err := json.ResponseNoContent(w); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *authHandler) UserLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, err := middleware.AuthUserFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	deletedUserSession, err := h.repo.UserSessions.Delete(ctx, authUser.SessionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.repo.CacheStore.UserSessions.Delete(deletedUserSession.ID.String())
			cookie.User.DeleteSession(w)
			json.ResponseBadRequest(w, r, err)
			return
		} else {
			json.ResponseInternalServerError(w, r, err)
		}
		return
	}

	h.repo.CacheStore.UserSessions.Delete(deletedUserSession.ID.String())
	cookie.User.DeleteSession(w)

	res := map[string]string{
		"deleted_session": deletedUserSession.ID.String(),
	}
	if err := json.ResponseOK(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *authHandler) PharmacyCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, err := middleware.AuthPharmacyFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	if err := json.ResponseNoContent(w); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}
