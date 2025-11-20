package router

import (
	"encoding/json"
	"events/internal/application"
	"net/http"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

func New(app *application.App) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chi_middleware.Logger)

	app.Log.Info("Server started")

	r.Get("/api/events/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"status": true})
	})

	r.Mount("/api/events/movie", movies(app))

	return r
}
