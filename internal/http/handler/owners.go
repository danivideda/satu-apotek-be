package handler

import (
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/patrickmn/go-cache"
)

type ownerHandler struct {
	repo  repository.Repository
	cache *cache.Cache
}

func (h *ownerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var ownerID int64 = 23
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
