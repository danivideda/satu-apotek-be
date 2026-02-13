package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type authHandler struct {
	repo  repository.Repository
}

var (
	ownerSessionTTL = env.GetString("OWNER_SESSION_TTL", "168h")
)

func (h *authHandler) OwnerRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload of username, email, and password
	var payload struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	// Create owner and insert owner's session into database
	passwordHash, err := argon2id.CreateHash(payload.Password, argon2id.DefaultParams)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
	ownerID, ownerSessionID, err := h.repo.Owners.Create(ctx, payload.Username, payload.Email, passwordHash)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	result := struct {
		OwnerID        int64  `json:"owner_id"`
		OwnerSessionID string `json:"owner_session_id"`
	}{OwnerID: ownerID, OwnerSessionID: ownerSessionID}

	if err := json.ResponseCreated(w, result); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *authHandler) OwnerLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload of username and password
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	// Get owner password
	owner, err := h.repo.Owners.GetByUsername(ctx, payload.Username)
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

	h.setOwnerSessionCookie(w, ownerSession.ID.String())

	if err := json.ResponseNoContent(w); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *authHandler) OwnerRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload struct {
		SessionID string `json:"session_id"`
	}
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}
	newOwnerSession, err := h.repo.OwnerSessions.Update(ctx, payload.SessionID, time.Now().Add(7*24*time.Hour))
	if err != nil {
		json.ResponseNotFound(w, r, err)
		return
	}

	if err := json.ResponseOK(w, newOwnerSession.ID); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *authHandler) OwnerLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload struct {
		SessionID string `json:"session_id"`
	}
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	deletedOwnerSession, err := h.repo.OwnerSessions.Delete(ctx, payload.SessionID)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := map[string]string{
		"deleted_session": deletedOwnerSession.ID.String(),
	}
	if err := json.ResponseOK(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *authHandler) generateSessionID() (string, error) {
	b := make([]byte, 32) // 32 bytes for a 256-bit ID
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (h *authHandler) setOwnerSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "owner_session",
		Value:    sessionID,
		Expires:  time.Now().Add(1 * time.Hour),
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	})
}
