package handler

import (
	"errors"
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	Owner    *ownerHandler
	Auth     *authHandler
	User     *userHandler
	Pharmacy *pharmacyHandler
}

func New(repo repository.Repository) Handler {
	return Handler{
		Auth:     &authHandler{repo: repo},
		Owner:    &ownerHandler{repo: repo},
		User:     &userHandler{repo: repo},
		Pharmacy: &pharmacyHandler{repo: repo},
	}
}

var validate = validator.New(validator.WithRequiredStructEnabled())

// Accepts pointer value to a payload struct, which then get parsed and validated.
//
// If error happens it will return false and handle the request automatically
func parseAndValidateJSONPayload(w http.ResponseWriter, r *http.Request, payloadPointer any) bool {
	if err := json.Read(w, r, payloadPointer); err != nil {
		json.ResponseBadRequest(w, r, err)
		return false
	}
	if err := validate.Struct(payloadPointer); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			json.ResponseBadRequest(w, r, validationErrors)
		} else {
			json.ResponseInternalServerError(w, r, err)
		}
		return false
	}
	return true
}