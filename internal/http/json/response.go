package json

import (
	"log"
	"net/http"
)

func ResponseOK(w http.ResponseWriter, data any) error {
	return WriteResponse(w, http.StatusOK, data)
}

func ResponseCreated(w http.ResponseWriter, data any) error {
	return WriteResponse(w, http.StatusCreated, data)
}

func ResponseNoContent(w http.ResponseWriter) error {
	return WriteResponse(w, http.StatusNoContent, nil)
}

func ResponseInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error, %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteResponseErr(w, http.StatusInternalServerError, "the server encountered some problem")
}

func ResponseNotFound(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found, %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteResponseErr(w, http.StatusNotFound, "not found")
}

func ResponseBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request, %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteResponseErr(w, http.StatusBadRequest, "bad request")
}

func ResponseUnauthorized(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not authorized, %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteResponseErr(w, http.StatusUnauthorized, "unauthorized")
}

func ResponseForbidden(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("forbidden, %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteResponseErr(w, http.StatusForbidden, "forbidden")
}

func ResponseInvalidCSRFToken(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("forbidden, %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteResponseErr(w, http.StatusForbidden, "invalid CSRF token")
}