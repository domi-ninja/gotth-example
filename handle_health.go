package app

import (
	"log"
	"net/http"
	"time"

	"domi.ninja/example-project/internal/db_generated"
	"github.com/google/uuid"
)

func (app *App) HandleHealth(w http.ResponseWriter, r *http.Request) {

	newHealthCheck := db_generated.HealthchecksCreateParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
	}

	newHealthCheckCreated, err := app.db.HealthchecksCreate(r.Context(), newHealthCheck)
	if err != nil {
		log.Print("error creating healthcheck: ", err)
		RespondWithError(w, http.StatusInternalServerError)
		return
	}

	newhealthCheckRead, err := app.db.HealthchecksRead(r.Context(), newHealthCheckCreated.ID)
	if err != nil {
		log.Print("error reading healthcheck: ", err)
		RespondWithError(w, http.StatusInternalServerError)
		return
	}

	if newHealthCheckCreated.ID != newhealthCheckRead.ID {
		log.Print("healthcheck id mismatch: %v != %v", newHealthCheckCreated.ID, newhealthCheckRead.ID)
		RespondWithError(w, http.StatusInternalServerError)
		return
	}

	// Basic health check - just return OK for now
	// In a production environment, you might want to check database, external services, etc.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"db":"ok"}`))
}
