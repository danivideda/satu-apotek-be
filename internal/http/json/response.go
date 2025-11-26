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

func ResponseNoContent(w http.ResponseWriter, data any) error {
	return WriteResponse(w, http.StatusNoContent, data)
}

func ResponseInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteResponseErr(w, http.StatusInternalServerError, "the server encountered some problem")
}

func ResponseNotFound(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteResponseErr(w, http.StatusNotFound, "not found")
}

func ResponseBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteResponseErr(w, http.StatusBadRequest, "bad request")
}

func ResponseUnauthorized(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not authorized error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	WriteResponseErr(w, http.StatusUnauthorized, "unauthorized")
}