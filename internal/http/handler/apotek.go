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

	authOwner, err := middleware.AuthOwnerFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	// TODO: 
	// 1. input necessary data about the Apotek
	// 2. Apotek Name, address, description (optional)
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
	res := map[string]any{
		"owner_id":    authOwner.ID,
		"owner_session":  authOwner.SessionID,
		"apotek_name": payload.ApotekName,
	}
	if err := json.ResponseOK(w, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}
