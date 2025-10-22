package handler

import (
	"log"
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) GetAllOwners(w http.ResponseWriter, r *http.Request) {
	owners, err := h.store.ListOwners(r.Context())
	if err != nil {
		log.Printf("Error: %s", err)
	}

	if err := json.ResponseJSON(w, http.StatusCreated, owners); err != nil {
		log.Printf("Error occured: %s", err)
	}
}

type CreateOwnerPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) CreateOwner(w http.ResponseWriter, r *http.Request) {
	var payload CreateOwnerPayload
	if err := json.ReadJSON(w, r, &payload); err != nil {
		log.Printf("error response: %s", err.Error())
		json.ResponseJSONError(w, http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error response: %s", err.Error())
		json.ResponseJSONError(w, http.StatusBadRequest, err.Error())
	}

	param := dbsqlc.CreateOwnerParams{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
	}
	owner, err := h.store.CreateOwner(r.Context(), param)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	if err := json.ResponseJSON(w, http.StatusCreated, owner); err != nil {
		log.Printf("error response: %s", err.Error())
		json.ResponseJSONError(w, http.StatusBadRequest, err.Error())
	}
}
