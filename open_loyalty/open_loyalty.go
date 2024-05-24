package open_loyalty

import (
	"log/slog"
	"net/http"
	"time"
)

type (
	// TalonOneClient is a struct that represents the Talon.One API client.
	OpenLoyaltyClient struct {
		Client        *http.Client
		Logger        *slog.Logger
		Cache         CacheRepository
		BaseURL       string
		Username      string
		Password      string
		ClientTimeout *time.Duration
		StoreID       string
	}
	// Config holds the configuration values for the Talon.One API client.
	Config struct {
		URL           string
		Username      string
		Password      string
		ClientTimeout *time.Duration
		StoreID       string
	}

	Members struct {
		LoyaltyID string `xml:"loyaltyId"`
	}
)
