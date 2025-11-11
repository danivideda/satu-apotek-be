package handler

import (
	"net/http"
	"strconv"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/http/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authToken struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func (h *Handler) RegisterOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type registerOwnerPayload struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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

	token, err := newAuthToken(owner.ID)
	if err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}

	responseData := struct {
		*dbsqlc.CreateOwnerRow
		*authToken
	}{
		CreateOwnerRow: owner,
		authToken:      token,
	}

	if err := json.WriteResponse(w, http.StatusCreated, responseData); err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}
}

func (h *Handler) LoginOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type loginOwnerPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

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

	token, err := newAuthToken(owner.ID)
	if err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}
	response := struct {
		Username string `json:"username"`
		Message  string `json:"message"`
		*authToken
	}{
		Username:  payload.Username,
		Message:   "Login success",
		authToken: token,
	}

	if err := json.WriteResponse(w, http.StatusOK, response); err != nil {
		internalServerErrorResponse(w, r, err)
		return
	}
}

func newAuthToken(id int32) (*authToken, error) {
	ownerID := strconv.Itoa(int(id))
	refreshToken, err := jwt.NewRefreshToken(ownerID)
	if err != nil {
		return nil, err
	}
	accessToken, err := jwt.NewAccessToken(ownerID)
	if err != nil {
		return nil, err
	}

	token := authToken{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}

	return &token, nil
}
