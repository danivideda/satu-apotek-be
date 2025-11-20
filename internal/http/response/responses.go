package response

import (
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"net/http"
)

func OK(w http.ResponseWriter, r *http.Request, data any) error {
	return json.WriteResponse(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, r *http.Request, data any) error {
	return json.WriteResponse(w, http.StatusCreated, data)
}
