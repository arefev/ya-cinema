package router

import (
	"events/internal/application"
	"events/internal/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func movies(app *application.App) http.Handler {
	r := chi.NewRouter()

	movieHandler := handler.NewMovie(app)

	r.Post("/", movieHandler.Create)

	return r
}
