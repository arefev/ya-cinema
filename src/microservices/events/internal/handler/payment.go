package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"events/internal/application"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type payment struct {
	app *application.App
}

func NewPayment(app *application.App) *payment {
	return &payment{app: app}
}

func (p *payment) Create(w http.ResponseWriter, r *http.Request) {
	p.app.Log.Sugar().Infof("kafka address %s", p.app.Conf.Kafka)

	w.Header().Set("Content-Type", "application/json")

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{p.app.Conf.Kafka},
		Topic:   "payment-events",
	})
	defer writer.Close()

	content, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		p.app.Log.Error("Ошибка чтения тела запроса", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	err = writer.WriteMessages(r.Context(), kafka.Message{
		Value: []byte(content),
	})

	if err != nil {
		p.app.Log.Error("Ошибка при отправке сообщение в брокер", zap.Error(err))
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
			"id": "payment-1-viewed",
			"type": "payment",
			"timestamp": "2023-01-15T14:30:00Z",
			"payload": map[string]any{},
		},
	})
}
