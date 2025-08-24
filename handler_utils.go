package app

import (
	"encoding/json"
	"log"
	"net/http"

	"domi.ninja/example-project/frontend/components"
	"domi.ninja/example-project/webhelp"
)

func RespondWithText(contentType string, w http.ResponseWriter, code int, text string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	w.Write([]byte(text))
}

func RespondWithJson(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// renders generic error message for a given http status code, does not leak internal error details
func RespondWithError(w http.ResponseWriter, code int) {
	errMsg := http.StatusText(code)

	// TODO consider catching some errors here and not returning the raw status code to make it harder
	// to hack the app

	w.WriteHeader(code)
	w.Write([]byte(errMsg))
}

func RespondWithHtmlError(w http.ResponseWriter, r *http.Request, code int, message string) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))

	err := webhelp.RenderHTML(r.Context(), w, components.Error(message))
	if err != nil {
		log.Print("error rendering posts: ", err)
		RespondWithError(w, http.StatusInternalServerError)
		return
	}
}
