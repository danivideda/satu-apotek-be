package handler

import (
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type pharmacyHandler struct {
	repo repository.Repository
}

func (h *pharmacyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authOwner, err := middleware.AuthOwnerFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	// TODO:
	// 1. input necessary data about the Apotek
	// 2. Apotek Name
	var payload struct {
		ApotekName string `json:"apotek_name"`
	}
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	// TODO:
	// 1. Create new Apotek under the OwnerID relation in DB table
	// 2. Return ApotekID
	apotek, err := h.repo.Pharmacies.Create(ctx, authOwner.ID, payload.ApotekName)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := map[string]any{
		"owner_id":    authOwner.ID,
		"apotek_name": payload.ApotekName,
		"apotek_id":   apotek.ID,
	}
	if err := json.ResponseCreated(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *pharmacyHandler) GetByOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authOwner, err := middleware.AuthOwnerFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	pharmacies, err := h.repo.Pharmacies.GetByOwner(ctx, authOwner.ID)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	if err := json.ResponseOK(w, pharmacies); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}
