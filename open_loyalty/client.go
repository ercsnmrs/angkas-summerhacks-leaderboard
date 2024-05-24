package open_loyalty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const (
	// Auth Endpoints
	loginCheckEndpoint   = "/api/admin/login_check"
	tokenRefreshEndpoint = "/api/token/refresh"

	// Import Endpoints
	importMembersEndpoint = "/api/%s/import/member"
)

type Client interface {
	LoginCheck(ctx context.Context) (*loginCheckResponse, error)
	TokenRefresh(ctx context.Context) (*refreshTokenResponse, error)
	ImportMembers(ctx context.Context) (*importMembersResponse, error)
}

type CacheRepository interface {
	GetJWTToken(ctx context.Context) (string, error)
	GetJWTExpiry(ctx context.Context) (time.Duration, error)
	SetAuthenticationTokens(ctx context.Context, JWTToken string, refreshToken string) error
	GetRefreshToken(ctx context.Context) (string, error)
}

// Creates a new instance of the Open Loyalty API client with the provided configuration.
func NewOpenLoyaltyClient(config *Config, c CacheRepository, l *slog.Logger) *OpenLoyaltyClient {
	client := &OpenLoyaltyClient{
		Client: &http.Client{
			Timeout: 30 * time.Second, // default
		},

		BaseURL:  config.URL,
		Username: config.Username,
		Password: config.Password,
		StoreID:  config.StoreID,
		Cache:    c,
		Logger:   l,
	}

	if config.ClientTimeout != nil {
		client.Client.Timeout = *config.ClientTimeout
	}

	return client
}

func (c *OpenLoyaltyClient) preparePublicAPIRequest(req http.Request) *http.Request {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return &req
}

// func (c *OpenLoyaltyClient) preparePrivateAPIRequest(ctx context.Context, req http.Request) (*http.Request, error) {
// 	token, err := c.getJWTToken(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", token))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Accept", "application/json")

// 	return &req, nil
// }

func (c *OpenLoyaltyClient) prepareImportMembersAPIRequest(ctx context.Context, req http.Request) (*http.Request, error) {
	token, err := c.getJWTToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", token))
	req.Header.Set("Accept", "application/json")

	return &req, nil
}

type (
	loginCheckRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	loginCheckResponse struct {
		JWTToken     string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	refreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	refreshTokenResponse struct {
		JWTToken     string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	importMembersRequest struct {
		File string `json:"file"`
	}

	importMembersResponse struct {
		ImportID string `json:"importId"`
	}
)

// Stores the JWT token in the cache and refreshes it if the TTL is less than 1 hour.
// TODO: Update implementation of REDIS Service
func (c *OpenLoyaltyClient) getJWTToken(ctx context.Context) (token string, err error) {
	// Get the JWT token from the cache
	key, err := c.Cache.GetJWTToken(ctx)
	if err != nil {
		return "", err
	}

	if key == "" {
		resp, err := c.LoginCheck(ctx)
		if err != nil {
			return "", err
		}

		return resp.JWTToken, nil
	}

	// Check the TTL of the token
	ttl, err := c.Cache.GetJWTExpiry(ctx)
	if err != nil {
		return "", err
	}

	if ttl == 0 {
		resp, err := c.LoginCheck(ctx)
		if err != nil {
			return "", err
		}

		return resp.JWTToken, nil
	}

	// If the TTL is less than 1 hour, refresh the token
	if time.Hour > ttl {
		// Get the Refresh token from the cache
		refreshKey, err := c.Cache.GetRefreshToken(ctx)
		if err != nil {
			return "", err
		}

		newToken, err := c.RefreshToken(ctx, refreshKey)
		if err != nil {
			return "", err
		}

		err = c.Cache.SetAuthenticationTokens(ctx, newToken.JWTToken, newToken.RefreshToken)
		if err != nil {
			return "", err
		}

	}

	return key, nil
}

// Performs a login check.
func (c *OpenLoyaltyClient) LoginCheck(ctx context.Context) (*loginCheckResponse, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, loginCheckEndpoint)

	// Create request payload
	payload := &loginCheckRequest{
		Username: c.Username,
		Password: c.Password,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req = c.preparePublicAPIRequest(*req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	// if resp.StatusCode != http.StatusOK {
	// 	return some error
	// }
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger.ErrorContext(ctx, "error reading response body", "error", err)
		return nil, err
	}
	var respData loginCheckResponse
	err = json.Unmarshal(rawBody, &respData)
	if err != nil {
		c.Logger.ErrorContext(ctx, "error unmarshalling response body", "error", err)
		return nil, err
	}

	// Store the JWT token in the cache
	err = c.Cache.SetAuthenticationTokens(ctx, respData.JWTToken, respData.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &respData, nil
}

// Performs a token refresh.
func (c *OpenLoyaltyClient) RefreshToken(ctx context.Context, token string) (*refreshTokenResponse, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, tokenRefreshEndpoint)

	// Create request payload
	payload := &refreshTokenRequest{
		RefreshToken: token,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req = c.preparePublicAPIRequest(*req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	// if resp.StatusCode != http.StatusOK {
	// 	return some error
	// }
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger.ErrorContext(ctx, "error reading response body", "error", err)
		return nil, err
	}
	var respData refreshTokenResponse
	err = json.Unmarshal(rawBody, &respData)
	if err != nil {
		c.Logger.ErrorContext(ctx, "error unmarshalling response body", "error", err)
		return nil, err
	}

	return &respData, nil
}

// Imports members.

func (c *OpenLoyaltyClient) ImportMembers(ctx context.Context, importRequest importMembersRequest) (*importMembersResponse, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, fmt.Sprintf(importMembersEndpoint, c.StoreID))
	filePath := importRequest.File

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("import[file]", filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file content to form part
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close writer
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req, err = c.prepareImportMembersAPIRequest(ctx, *req)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger.ErrorContext(ctx, "error reading response body", "error", err)
		return nil, err
	}

	var respData importMembersResponse
	err = json.Unmarshal(rawBody, &respData)
	if err != nil {
		c.Logger.ErrorContext(ctx, "error unmarshalling response body", "error", err)
		return nil, err
	}

	// Read response body
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error importing members: %s", resp.Status)
	}
	return &respData, nil
}
