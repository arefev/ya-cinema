package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"events/internal/application"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type user struct {
	app *application.App
}

func NewUser(app *application.App) *user {
	return &user{app: app}
}

func (u *user) Create(w http.ResponseWriter, r *http.Request) {
	u.app.Log.Sugar().Infof("kafka address %s", u.app.Conf.Kafka)

	w.Header().Set("Content-Type", "application/json")

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{u.app.Conf.Kafka},
		Topic:   "user-events",
	})
	defer writer.Close()

	content, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		u.app.Log.Error("Ошибка чтения тела запроса", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	err = writer.WriteMessages(r.Context(), kafka.Message{
		Value: []byte(content),
	})

	if err != nil {
		u.app.Log.Error("Ошибка при отправке сообщение в брокер", zap.Error(err))
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
			"id": "user-1-viewed",
			"type": "user",
			"timestamp": "2023-01-15T14:30:00Z",
			"payload": map[string]any{},
		},
	})
}
