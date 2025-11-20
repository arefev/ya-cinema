package handler

import (
	"net/http"

	"events/internal/application"
)

type movie struct {
	app *application.App
}

func NewMovie(app *application.App) *movie {
	return &movie{app: app}
}

func (m *movie) Create(w http.ResponseWriter, r *http.Request) {

}
