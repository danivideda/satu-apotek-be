package handler

import (
	"log"
	"net/http"

	"github.com/danivideda/satu-apotek-be/internal/http/json"
)

func internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	json.WriteResponseErr(w, http.StatusInternalServerError, "the server encountered some problem")
}

func notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	json.WriteResponseErr(w, http.StatusNotFound, "not found")
}

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	json.WriteResponseErr(w, http.StatusBadRequest, "bad request")
}
