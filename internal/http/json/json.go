package json

import (
	"encoding/json"
	"net/http"
	"slices"
)

func Write(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	noContentStatus := []int{
		http.StatusContinue,
		http.StatusSwitchingProtocols,
		http.StatusNoContent,
		http.StatusResetContent,
		http.StatusNotModified,
	}

	if !slices.Contains(noContentStatus, status) {
		return json.NewEncoder(w).Encode(data)
	}

	return nil
}

// Read from ResponseWriter body
func Read(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576 // 1 Mega Bytes
	http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func WriteResponseErr(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return Write(w, status, &envelope{Error: message})
}

func WriteResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return Write(w, status, &envelope{Data: data})
}