package handler

import (
	"io"
	"net/http"

	"events/internal/application"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type movie struct {
	app *application.App
}

func NewMovie(app *application.App) *movie {
	return &movie{app: app}
}

func (m *movie) Create(w http.ResponseWriter, r *http.Request) {
	m.app.Log.Sugar().Infof("kafka address %s", m.app.Conf.Kafka)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{m.app.Conf.Kafka},
		Topic:   "movie-events",
	})
	defer writer.Close()

	content, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		m.app.Log.Error("Ошибка чтения тела запроса", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = writer.WriteMessages(r.Context(), kafka.Message{
		Value: []byte(content),
	})

	if err != nil {
		m.app.Log.Error("Ошибка при отправке сообщение в брокер", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
