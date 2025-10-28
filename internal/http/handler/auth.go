package handler

import (
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"golang.org/x/crypto/bcrypt"
)

type registerOwnerPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) RegisterOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload registerOwnerPayload
	if err := json.Read(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	param := dbsqlc.CreateOwnerParams{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
	}
	owner, err := h.store.Owners.CreateOwner(ctx, param)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := json.WriteResponse(w, http.StatusCreated, owner); err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}
}
