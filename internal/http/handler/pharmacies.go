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

type pharmacyJSON struct {
	AppID     any `json:"app_id"`
	Name      any `json:"name"`
	CreatedAt any `json:"created_at"`
	UpdatedAt any `json:"updated_at"`
}

func (h *pharmacyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authOwner, err := middleware.AuthOwnerFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	var payload struct {
		ApotekName string `json:"apotek_name"`
	}
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	// TODO:
	// 1. Create new Apotek under the OwnerID relation in DB table
	// 2. Return AppID
	apotek, err := h.repo.Pharmacies.Create(ctx, authOwner.ID, payload.ApotekName)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := pharmacyJSON{
		AppID:     apotek.AppID,
		Name:      apotek.Name,
		CreatedAt: apotek.CreatedAt,
		UpdatedAt: apotek.UpdatedAt,
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

	res := []pharmacyJSON{}
	for _, pharmacy := range *pharmacies {
		item := pharmacyJSON{
			AppID:     pharmacy.AppID,
			Name:      pharmacy.Name,
			CreatedAt: pharmacy.CreatedAt,
			UpdatedAt: pharmacy.UpdatedAt,
		}
		res = append(res, item)
	}
	if err := json.ResponseOK(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}
