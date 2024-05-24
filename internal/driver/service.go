package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/trip"
)

// Service represents Tier service.
type Service struct {
	cache    CacheRepository
	provider ProviderService
	logger   *slog.Logger
}

// cacheRepository manages redis or any nosql storage operations
type CacheRepository interface {
	GetDriverRating(ctx context.Context, id string) (driver string, err error)
	SetDriverRating(ctx context.Context, driver Driver) (err error)
	CheckHighestNetEarnings(ctx context.Context, netEarning float64, serviceZone string) (highestNetEarnings float64, err error)
}

// providerService manages external service operations
type ProviderService interface {
	ImportDriverRating(ctx context.Context, list []Driver) (err error)
}

// NewService returns new tier service.
func NewDriverService(c CacheRepository, p ProviderService, l *slog.Logger) *Service {
	return &Service{
		cache:    c,
		provider: p,
		logger:   l,
	}
}

func (s Service) WriteToCSV(ctx context.Context, userID string, rating RFM, tier string) error {
	// Open the CSV file for writing
	file, err := os.OpenFile("tier_data.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new line with the tier data
	line := fmt.Sprintf("%s,%f,%f,%f,%s\n", userID, rating.Recency, rating.Frequency, rating.Monetary, tier)

	// Write the line to the CSV file
	_, err = file.WriteString(line)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetDriver(ctx context.Context, driverID string) (Driver, error) {
	d := Driver{}

	// Get the tier from the cache
	c, err := s.cache.GetDriverRating(ctx, driverID)
	if err != nil {
		return d, err
	}

	if c != "nodata" {
		// Convert the tier to a Driver struct from the cache
		err = json.Unmarshal([]byte(c), &d)
		if err != nil {
			return d, err
		}
	}

	// Convert the tier to a DriverRating struct from the cache

	return d, nil
}

func (s Service) UpdateUserRating(ctx context.Context, trip trip.Event) (Driver, error) {
	driverID := trip.DriverID
	newDriver := Driver{}

	// Get the tier from the cache
	c, err := s.cache.GetDriverRating(ctx, driverID)
	if err != nil {
		return newDriver, err
	}

	var driver Driver
	if c != "nodata" {
		// Convert the tier to a Driver struct from the cache
		err = json.Unmarshal([]byte(c), &driver)
		if err != nil {
			return newDriver, err
		}
	}

	// Update the driver rating
	newNetIncome := trip.Price.DriverEarnings + driver.NetIncome

	// Check Highest Net Earnings for the service zone, if yes replae
	netEarnings, err := s.cache.CheckHighestNetEarnings(ctx, newNetIncome, driver.ServiceZone)
	if err != nil {
		return newDriver, err
	}

	newDriver = Driver{
		DriverID:                     driverID,
		LastCompletedTripDate:        time.Now(),
		NetIncome:                    newNetIncome,
		NumberOfCompletedTrips:       driver.NumberOfCompletedTrips + 1,
		UniqueDateWithCompletedTrips: driver.UniqueDateWithCompletedTrips,
		ServiceZone:                  driver.ServiceZone,
		Rating: Rating{
			RFM: RFM{
				Recency:   driver.CalculateRecency(),
				Frequency: driver.CalculateFrequency(),
				Monetary:  driver.CalculateMonetary(netEarnings),
			},
		},
	}

	newDriver.Rating.Average = newDriver.CalculateAverage()

	return newDriver, nil
}
