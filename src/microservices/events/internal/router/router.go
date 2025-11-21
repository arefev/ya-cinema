package router

import (
	"encoding/json"
	"events/internal/application"
	"events/internal/handler"
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

	movieHandler := handler.NewMovie(app)
	userHandler := handler.NewUser(app)
	paymentHandler := handler.NewPayment(app)

	r.Post("/api/events/movie", movieHandler.Create)
	r.Post("/api/events/user", userHandler.Create)
	r.Post("/api/events/payment", paymentHandler.Create)

	return r
}
