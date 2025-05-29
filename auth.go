package go_salesforce_api_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client represents the Salesforce OAuth token response
type Client struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	TokenType   string `json:"token_type"`
	IssuedAt    string `json:"issued_at"`
}

// Auth handles authentication with Salesforce
type Auth struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	TokenURL     string
}

// AuthenticatePassword performs an OAuth login and retrieves an access token
func (a *Auth) AuthenticatePassword() (*Client, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", a.ClientID)
	data.Set("client_secret", a.ClientSecret)
	data.Set("username", a.Username)
	data.Set("password", a.Password)

	resp, err := http.Post(a.TokenURL, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to authenticate with Salesforce")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var client Client
	if err := json.Unmarshal(body, &client); err != nil {
		return nil, err
	}

	return &client, nil
}

// AuthenticateClientCredentials performs Client Credentials OAuth flow
func (a *Auth) AuthenticateClientCredentials() (*Client, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", a.ClientID)
	data.Set("client_secret", a.ClientSecret)

	resp, err := http.Post(a.TokenURL, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to authenticate with Salesforce, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var client Client
	if err := json.Unmarshal(body, &client); err != nil {
		return nil, err
	}

	return &client, nil
}
