package router

import (
	"net/http"
	"proxy/internal/application"
	"proxy/internal/handler"

	"github.com/go-chi/chi/v5"
)

func movies(app *application.App) http.Handler {
	r := chi.NewRouter()

	movieHandler := handler.NewProxyHandler(app, app.Conf.MoviesServiceUrl)

	if app.Conf.GradualMigration {
		movieHandler.UseGradualMigration()
	}

	r.HandleFunc("/", movieHandler.Proxy)

	return r
}
