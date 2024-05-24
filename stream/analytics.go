package stream

import (
	"context"

	"log/slog"

	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/kafka"
	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/trip"
)

// Service is an interface that sends events to an analytics engine.
type Service interface {
	SendFakeTripData(ctx context.Context, trip trip.Event) error
}

type KafkaService struct {
	logger   *slog.Logger
	producer kafka.Writer
}
