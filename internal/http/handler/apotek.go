package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type apotekHandler struct {
	repo repository.Repository
}

func (h *apotekHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload struct {
		Name string `json:"name"`
	}
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	auth := middleware.AuthClaimsFromContext(ctx)
	if auth == nil {
		json.ResponseInternalServerError(w, r, ErrInvalidAuthToken)
		return
	}

	pharmacy, err := h.repo.Pharmacies.Create(ctx, auth.ID, payload.Name)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	if err := json.ResponseCreated(w, map[string]any{"pharmacy": pharmacy}); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *apotekHandler) CreateOrUpdateCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload struct {
		ApotekID string `json:"apotek_id"`
	}
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	newCode, err := h.generateApotekCode()
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
	// apotekCode, err := h.repo.ApotekCode.Create(ctx, payload.ApotekID, mycode)
	apotekCode, err := h.repo.ApotekCode.Upsert(ctx, payload.ApotekID, newCode)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	res := map[string]any{
		"apotek_id": apotekCode.ApotekID,
		"code":      apotekCode.Code,
		"expires":   apotekCode.Expires.Time,
	}
	if err := json.ResponseCreated(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *apotekHandler) VerifyCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload struct {
		Code string `json:"code"`
	}
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	// Check Apotek Code in the repository
	apotekCode, err := h.repo.ApotekCode.GetByCode(ctx, payload.Code)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	// TODO: Send back the long lived access token with Apotek ID + Apotek Session ID
	res := map[string]any{
		"apotek_id": apotekCode.ApotekID,
		"token": "This-is-the-JWT-token",
	}
	if err := json.ResponseOK(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *apotekHandler) generateApotekCode() (string, error) {
	b := make([]byte, 3) // 32 bytes for a 256-bit ID
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	encoded := hex.EncodeToString(b)
	return encoded, nil
}