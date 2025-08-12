package app

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithText(contentType string, w http.ResponseWriter, code int, text string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	w.Write([]byte(text))
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// renders generic error message for a given http status code, does not leak internal error details
func respondWithError(w http.ResponseWriter, code int) {
	errMsg := "error"
	if code < 400 {
		log.Fatalf("Argument error, think about what you are doing here.")
	}

	if code == http.StatusUnauthorized {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"Restricted\"")
		errMsg = "Unauthorized"
	}

	if code == http.StatusForbidden {
		errMsg = "Forbidden"
	}

	if code == http.StatusNotFound {
		errMsg = "Not Found"
	}

	if code == http.StatusInternalServerError {
		errMsg = "Internal Server Error"
	}

	if code == http.StatusBadRequest {
		errMsg = "Bad Request"
	}

	if code == http.StatusUnauthorized {
		errMsg = "Unauthorized"
	}

	w.WriteHeader(code)
	w.Write([]byte(errMsg))
}
