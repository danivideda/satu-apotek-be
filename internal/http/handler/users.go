package handler

import (
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type userHandler struct {
	repo repository.Repository
}

func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type createUserPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
		PharmacyID string `json:"pharmacy_id"`
	}

	var payload createUserPayload
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	passwordHash, err := hashPassword(payload.Password)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	user, err := h.repo.Users.Create(ctx, payload.Username, passwordHash, payload.PharmacyID)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	res := map[string]string{
		"username": user.Username,
		"pharmacy_id": user.PharmacyID.String(),
	}

	if err := json.WriteResponse(w, http.StatusCreated, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}
