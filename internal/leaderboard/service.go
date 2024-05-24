package leaderboard

import (
	"context"
	"log/slog"

	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/driver"
	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/trip"
)

// Service represents Tier service.
type Service struct {
	cache  CacheRepository
	user   UserRepository
	logger *slog.Logger
}

// cacheRepository manages redis or any nosql storage operations
type CacheRepository interface {
	GetActiveLeaderboard(ctx context.Context, scope string) (*[]driver.Driver, error)
	RefreshLeaderboard(ctx context.Context, user driver.Driver) error
}

type UserRepository interface {
	UpdateUserRating(ctx context.Context, trip trip.Event) (driver.Driver, error)
}

// NewService returns new tier service.
func NewLeaderboardService(c CacheRepository, u UserRepository, l *slog.Logger) *Service {
	return &Service{
		user:   u,
		cache:  c,
		logger: l,
	}
}

func (s Service) GetLeaderboard(ctx context.Context, scope string) (Leaderboard, error) {
	leaders := Leaderboard{}

	// Get the tier from the cache
	list, err := s.cache.GetActiveLeaderboard(ctx, scope)
	if err != nil {
		return leaders, err
	}

	return Leaderboard{
		Drivers: *list,
	}, nil
}

func (s Service) UpdateLeaderboard(ctx context.Context, trip trip.Event) error {
	s.logger.Info("updating leaderboard...")

	// get the driver from the cache
	user, err := s.user.UpdateUserRating(ctx, trip)
	if err != nil {
		return err
	}
	// update cache
	err = s.cache.RefreshLeaderboard(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
