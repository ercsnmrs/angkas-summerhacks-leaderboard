package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/driver"
)

type RedisService struct {
	Client *Client
	logger *slog.Logger
}

// Creates a new instance of the Open Loyalty Service client with the provided configuration.
func NewCacheService(client Client, logger *slog.Logger) *RedisService {
	svc := &RedisService{
		Client: &client,
		logger: logger,
	}

	return svc
}

func (c *RedisService) GetJWTToken(ctx context.Context) (string, error) {
	val, err := c.Client.Get(ctx, "jwt_token").Result()

	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}

	return val, nil
}

func (c *RedisService) GetJWTExpiry(ctx context.Context) (time.Duration, error) {
	ttl, err := c.Client.TTL(ctx, "jwt_token").Result()

	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	return ttl, nil
}

func (c *RedisService) GetRefreshToken(ctx context.Context) (string, error) {
	token, err := c.Client.Get(ctx, "refresh_token").Result()

	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("refresh token not found")
		}
		return "", err
	}

	return token, nil
}

func (c *RedisService) SetAuthenticationTokens(ctx context.Context, JWTToken string, refreshToken string) error {
	_, err := c.Client.Set(ctx, "jwt_token", JWTToken, time.Hour*23).Result()
	if err != nil {
		return err
	}
	_, err = c.Client.Set(ctx, "refresh_token", refreshToken, time.Hour*23).Result()
	if err != nil {
		return err
	}

	return nil
}
func (c *RedisService) CheckHighestNetEarnings(ctx context.Context, netEarnings float64, serviceZone string) (float64, error) {
	key := fmt.Sprintf("highest_net_earnings:%s", serviceZone)
	expiration := time.Hour * 26
	cachedNetEarnings, err := c.Client.Get(ctx, key).Float64()
	if err != nil {
		if err == redis.Nil {
			// Cache not available, store current netEarnings
			_, err = c.Client.Set(ctx, key, netEarnings, expiration).Result()
			if err != nil {
				return 0, err
			}
			return netEarnings, nil
		}
		return 0, err
	}

	if netEarnings > cachedNetEarnings {
		// Net earnings is higher than cached value, update cache
		_, err = c.Client.Set(ctx, key, netEarnings, expiration).Result()
		if err != nil {
			return 0, err
		}
		return netEarnings, nil
	}

	return cachedNetEarnings, nil
}

func (c *RedisService) SetDriverRating(ctx context.Context, driver driver.Driver) (err error) {
	driverKey := fmt.Sprintf("driver:%s", driver.DriverID)

	// Store driver information in Redis as a JSON string
	driverJSON, err := json.Marshal(driver)
	if err != nil {
		return fmt.Errorf("failed to marshal driver data: %v", err)
	}

	if err := c.Client.Set(ctx, driverKey, driverJSON, 0).Err(); err != nil {
		return fmt.Errorf("failed to set driver data in Redis: %v", err)
	}

	return nil
}

func (c *RedisService) GetDriverRating(ctx context.Context, driverID string) (string, error) {
	driverKey := fmt.Sprintf("driver:%s", driverID)
	val, err := c.Client.Get(ctx, driverKey).Result()

	if err != nil {
		if err == redis.Nil {
			return "nodata", nil
		}
		return "", err
	}

	return val, nil
}

func (c *RedisService) GetActiveLeaderboard(ctx context.Context, scope string) (*[]driver.Driver, error) {
	key := fmt.Sprintf("driver_leaderboard:%s", scope)

	count := 100
	driverIDs, err := c.Client.ZRevRange(ctx, key, 0, int64(count-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get top drivers from leaderboard: %v", err)
	}

	var drivers []driver.Driver
	for _, driverID := range driverIDs {
		driverKey := fmt.Sprintf("driver:%s", driverID)
		driverJSON, err := c.Client.Get(ctx, driverKey).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get driver data from Redis: %v", err)
		}

		var driver driver.Driver
		if err := json.Unmarshal([]byte(driverJSON), &driver); err != nil {
			return nil, fmt.Errorf("failed to unmarshal driver data: %v", err)
		}
		drivers = append(drivers, driver)
	}

	return &drivers, nil
}

func (c *RedisService) GetPreviousLeaderboard(ctx context.Context) (string, error) {
	val, err := c.Client.Get(ctx, "previous_top_bikers").Result()

	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}

	return val, nil
}

func (c *RedisService) RefreshLeaderboard(ctx context.Context, driver driver.Driver) error {
	driverKey := fmt.Sprintf("driver:%s", driver.DriverID)

	// Store driver information in Redis as a JSON string
	driverJSON, err := json.Marshal(driver)
	if err != nil {
		return fmt.Errorf("failed to marshal driver data: %v", err)
	}

	if err := c.Client.Set(ctx, driverKey, driverJSON, 0).Err(); err != nil {
		return fmt.Errorf("failed to set driver data in Redis: %v", err)
	}

	leaderboadKey := fmt.Sprintf("driver_leaderboard:%s", driver.ServiceZone)
	// Update the leaderboard
	if err := c.Client.ZAdd(ctx, leaderboadKey, redis.Z{
		Score:  driver.Rating.Average,
		Member: driver.DriverID,
	}).Err(); err != nil {
		return fmt.Errorf("failed to update leaderboard: %v", err)
	}

	return nil
}
