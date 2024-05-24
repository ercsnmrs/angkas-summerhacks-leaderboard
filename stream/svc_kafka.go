package stream

import (
	"context"
	"encoding/json"

	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/kafka"
	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/trip"

	"log/slog"
)

var (
	TripTopic = "trips"
)

func NewKafkaService(l *slog.Logger, producer kafka.Writer) *KafkaService {
	return &KafkaService{
		logger:   l,
		producer: producer,
	}
}

func (a *KafkaService) SendFakeTripData(ctx context.Context, trip trip.Event) error {
	id, err := json.Marshal(map[string]string{
		"trip_id": trip.TripRequestID,
	})
	if err != nil {
		a.logger.ErrorContext(ctx, "failed to marshal event", slog.Any("err", err))
		return err
	}

	event, err := json.Marshal(trip)
	if err != nil {
		a.logger.ErrorContext(ctx, "failed to marshal event", slog.Any("err", err))
		return err
	}

	a.logger.DebugContext(ctx, "trip event details", slog.Any("details", event))

	if err := a.producer.Produce(ctx, id, event, TripTopic); err != nil {
		a.logger.ErrorContext(ctx, "failed to produce redeem coupon event", slog.Any("err", err))
		return err
	}
	return nil
}
