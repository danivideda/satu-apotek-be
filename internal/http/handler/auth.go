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

type loginOwnerPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) LoginOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload loginOwnerPayload
	if err := json.Read(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	owner, err := h.store.Owners.GetOwnerByUsername(ctx, payload.Username)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword(owner.Password, []byte(payload.Password)); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	response := struct {
		Username string `json:"username"`
		Message string `json:"message"`
	}{Username: payload.Username, Message: "Login success"}

	if err := json.WriteResponse(w, http.StatusOK, response); err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}
}
