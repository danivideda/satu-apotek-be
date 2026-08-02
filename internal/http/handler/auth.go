package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/danivideda/satu-apotek-be/internal/service"
)

type authHandler struct {
	repo repository.Repository
}

var (
	ownerSessionTTL = env.GetString("OWNER_SESSION_TTL", "168h")
	userSessionTTL  = env.GetString("USER_SESSION_TTL", "168h")
)

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

	// Create owner and insert owner's session into database
	passwordHash, err := argon2id.CreateHash(payload.Password, argon2id.DefaultParams)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
	ownerID, ownerSessionID, exp, err := h.repo.Owners.Create(ctx, payload.Username, payload.Email, passwordHash)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	service.SetOwnerCookies(w, ownerSessionID, exp)

	res := map[string]any{
		"owner_id": ownerID,
		// "owner_session": ownerSessionID, **SESSION SHOULD NOT BE INCLUDED IN ANY JSON PAYLOAD***
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

	// Get owner password
	owner, err := h.repo.Owners.GetByEmail(ctx, payload.Email)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	// check password hash to match
	match, err := argon2id.ComparePasswordAndHash(payload.Password, owner.PasswordHash)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	if !match {
		json.ResponseBadRequest(w, r, ErrInvalidPassword)
		return
	}

	ttl, err := time.ParseDuration(ownerSessionTTL)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
	ownerSession, err := h.repo.OwnerSessions.Create(ctx, owner.ID, time.Now().Add(ttl))
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
	service.SetOwnerCookies(w, ownerSession.ID.String(), ownerSession.ExpiresAt.Time)
	h.repo.CacheStore.OwnerSessions.SetDefault(ownerSession.ID.String(), owner.ID)

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

	deletedOwnerSession, err := h.repo.OwnerSessions.Delete(ctx, authOwner.SessionID)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	service.DeleteOwnerCookies(w)

	res := map[string]string{
		"deleted_session": deletedOwnerSession.ID.String(),
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
	// When CSRF cookie is missing, it means the CSRF protection middleware is unable to validate previous request,
	// hence, it deletes the CSRF cookie with `Set-Cookie` sent from the server.
	_, err = r.Cookie("owner_csrf")
	if err != nil {
		fmt.Println(err)
		service.SetOwnerCookies(w, authOwner.SessionID, authOwner.SessionExp)
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

	// Get user
	user, err := h.repo.Users.GetByID(ctx, payload.UserID)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	// check password hash to match
	match, err := argon2id.ComparePasswordAndHash(payload.Password, user.PasswordHash)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	if !match {
		json.ResponseBadRequest(w, r, ErrInvalidPassword)
		return
	}

	// Check if current Pharmacy session have the User that tried to login
	authPharmacy, err := middleware.AuthPharmacyFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
	if !service.UserExistsInPharmacy(authPharmacy.Users, user.ID) {
		json.ResponseForbidden(w, r, fmt.Errorf("user doesn't belong to current authd pharmacy"))
		return
	}

	ttl, err := time.ParseDuration(userSessionTTL)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
	userSession, err := h.repo.UserSessions.Create(ctx, user.ID, time.Now().Add(ttl))
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
	service.SetUserCookies(w, userSession.ID.String(), userSession.ExpiresAt.Time)
	h.repo.CacheStore.UserSessions.SetDefault(userSession.ID.String(), user.ID)

	if err := json.ResponseNoContent(w); err != nil {
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