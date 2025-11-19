package response

import (
	"github.com/danivideda/satu-apotek-be/internal/http/json"
	"log"
	"net/http"
)

func InternalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	json.WriteResponseErr(w, http.StatusInternalServerError, "the server encountered some problem")
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	json.WriteResponseErr(w, http.StatusNotFound, "not found")
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	json.WriteResponseErr(w, http.StatusBadRequest, "bad request")
}

func UnauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not authorized error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	json.WriteResponseErr(w, http.StatusUnauthorized, "unauthorized")
}

func OK(w http.ResponseWriter, r *http.Request, data any) error {
	return json.WriteResponse(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, r *http.Request, data any) error {
	return json.WriteResponse(w, http.StatusCreated, data)
}