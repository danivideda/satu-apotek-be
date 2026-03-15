package handler

import (
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type ownerHandler struct {
	repo repository.Repository
}

func (h *ownerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authOwner, err := middleware.AuthOwnerFromCtx(ctx)
	if err != nil {
		json.ResponseUnauthorized(w, r, err)
		return
	}
	var ownerID int64 = authOwner.OwnerID
	owners, err := h.repo.Owners.GetByID(ctx, ownerID)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	if err := json.WriteResponse(w, http.StatusOK, owners); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}
