package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/middleware"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type userHandler struct {
	repo repository.Repository
}

func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
		AppID    string `json:"app_id"`
	}
	if ok := parseAndValidateJSONPayload(w, r, &payload); !ok {
		return
	}

	// check if AppID belongs to Authd Owner
	authOwner, err := middleware.AuthOwnerFromCtx(ctx)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
	pharmacy, err := h.repo.Pharmacies.GetByAppID(ctx, payload.AppID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			json.ResponseBadRequest(w, r, fmt.Errorf("%w, app_id not found", err))
		default:
			json.ResponseInternalServerError(w, r, err)
		}
		return
	}
	// only pass if ownerID matches
	if pharmacy.OwnerID != authOwner.ID {
		json.ResponseForbidden(w, r, ErrAppIDNotAllowed)
		return
	}

	// create the user under this pharmacyID
	passwordHash, err := argon2id.CreateHash(payload.Password, argon2id.DefaultParams)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}
	user, err := h.repo.Users.Create(ctx, payload.Username, passwordHash, pharmacy.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateValue):
			json.ResponseBadRequest(w, r, fmt.Errorf("%w, username already exists", err))
		default:
			json.ResponseInternalServerError(w, r, err)
		}
		return
	}

	res := map[string]any{
		"username": user.Username,
	}

	if err := json.WriteResponse(w, http.StatusCreated, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}

func (h *userHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if err := json.ResponseOK(w, "from profile"); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}
