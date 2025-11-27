package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/jwt"
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
		"expires":   apotekCode.ExpiresAt.Time,
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

	// Send back the long lived access token with Apotek ID + Apotek Session ID
	token, err := jwt.NewApotekToken(apotekCode.ApotekID.String())
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
	res := map[string]any{
		"apotek_id": apotekCode.ApotekID,
		"token": token,
	}
	if err := json.ResponseOK(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *apotekHandler) generateApotekCode() (string, error) {
	b := make([]byte, 3) // 3 bytes for 6-character hex string code
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	encoded := hex.EncodeToString(b)
	return encoded, nil
}

func (h *apotekHandler) newApotekToken(apotekID string) (string, error) {
	return jwt.NewApotekToken(apotekID)
}