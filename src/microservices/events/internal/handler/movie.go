package handler

import (
	"encoding/json"
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

	w.Header().Set("Content-Type", "application/json")

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{m.app.Conf.Kafka},
		Topic:   "movie-events",
	})
	defer writer.Close()

	content, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		m.app.Log.Error("Ошибка чтения тела запроса", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	err = writer.WriteMessages(r.Context(), kafka.Message{
		Value: []byte(content),
	})

	if err != nil {
		m.app.Log.Error("Ошибка при отправке сообщение в брокер", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"partition": 0,
		"offset": 42,
		"event": map[string]any{
			"id": "movie-1-viewed",
			"type": "movie",
			"timestamp": "2023-01-15T14:30:00Z",
			"payload": map[string]any{},
		},
	})
}
