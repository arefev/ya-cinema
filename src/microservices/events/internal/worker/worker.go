package worker

import (
	"context"
	"errors"
	"events/internal/application"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type worker struct {
	app *application.App
}

func NewWorker(app *application.App) *worker {
	return &worker{app: app}
}

func (w *worker) Run(ctx context.Context) error {
	topics := []string{"movie-events", "payment-events", "user-events"}
	group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		ID:      "events-service",
		Brokers: []string{"kafka:9092"},
		Topics:  topics,
	})

	if err != nil {
		w.app.Log.Error("error creating consumer group", zap.Error(err))
		return fmt.Errorf("worker run failed: error creating consumer group: %w", err)
	}

	defer group.Close()

	for {
		gen, err := group.Next(ctx)
		if err != nil {
			w.app.Log.Error("error next consumer group generation", zap.Error(err))
			break
		}

		for _, topic := range topics {
			assignments := gen.Assignments[topic]
			for _, assignment := range assignments {
				partition, offset := assignment.ID, assignment.Offset
				gen.Start(func(ctx context.Context) {
					// create reader for this partition.
					reader := kafka.NewReader(kafka.ReaderConfig{
						Brokers:   []string{w.app.Conf.Kafka},
						Topic:     topic,
						Partition: partition,
					})
					defer reader.Close()

					// seek to the last committed offset for this partition.
					reader.SetOffset(offset)
					for {
						msg, err := reader.ReadMessage(ctx)
						if err != nil {
							if errors.Is(err, kafka.ErrGenerationEnded) {
								// generation has ended.  commit offsets.  in a real app,
								// offsets would be committed periodically.
								gen.CommitOffsets(map[string]map[int]int64{topic: {partition: offset + 1}})
								return
							}

							w.app.Log.Error("error reading message", zap.Error(err))
							return
						}

						w.app.Log.Info(
							"received message", 
							zap.String("topic", msg.Topic),
							zap.Int("partition", msg.Partition),
							zap.Int64("offset", msg.Offset),
							zap.String("message", string(msg.Value)),
						)
						offset = msg.Offset
					}
				})
			}
		}
	}

	return nil
}
