package handler

import (
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

func (h *apotekHandler) CreateCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload struct {
		ApotekID string `json:"apotek_id"`
	}
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	mycode := "ABCDEF"
	apotekCode, err := h.repo.ApotekCode.Create(ctx, payload.ApotekID, mycode)
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
