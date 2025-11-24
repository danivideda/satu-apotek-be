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

	type createPharmacyPayload struct {
		Name string `json:"name"`
	}
	var payload createPharmacyPayload
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

func (h *apotekHandler) Connect(w http.ResponseWriter, r *http.Request) {
	// Get OTP
	type ConnectApotekPayload struct {
		OTP string `json:"otp"`
	}

	var payload ConnectApotekPayload
	if err := json.Read(w, r, &payload); err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	// Check OTP in the repository
	// TODO: Create apotek_otp_codes table in database

	// Send back the long lived access token with Apotek ID + Apotek Session ID
}