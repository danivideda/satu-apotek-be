package handler

import (
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/jwt"
	"github.com/danivideda/satu-apotek-be/internal/store"
)

type userHandler struct {
	store store.Storage
}

// TODO user create handler
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

	user, err := h.store.Users.Create(ctx, payload.Username, passwordHash, payload.PharmacyID)
	if err != nil {
		json.ResponseBadRequest(w, r, err)
		return
	}

	token, exp, err := newAuthToken(user.ID.String(), jwt.RoleUser)
	if err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	if err := setCookies(w, token.RefreshToken, *exp, jwt.RoleUser); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}

	res := struct {
		Username string `json:"username"`
		PharmacyID string `json:"pharmacy_id"`
		*authToken
	}{
		Username: user.Username,
		PharmacyID: user.PharmacyID.String(),
		authToken: token,
	}

	if err := json.WriteResponse(w, http.StatusCreated, res); err != nil {
		json.ResponseInternalServerError(w, r, err)
		return
	}
}
