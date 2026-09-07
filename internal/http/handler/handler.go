package handler

import (
	"errors"
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/config"
	c "github.com/danivideda/satu-apotek-be/internal/http/cookie"
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

var cookie = c.New()

func New(repo repository.Repository, cfg config.Config) Handler {
	return Handler{
		Auth:     newAuthHandler(repo, cfg.Auth.OwnerSessionTTL, cfg.Auth.UserSessionTTL),
		Owner:    newOwnerHandler(repo),
		User:     newUserHandler(repo),
		Pharmacy: newPharmacyHandler(repo, cfg.Auth.PharmacySessionTTL),
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
