package go_salesforce_api_client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestAuthenticatePassword(t *testing.T) {
	t.Parallel()
	mockResponse := Client{
		AccessToken: "mock_access_token",
		InstanceURL: "https://mock.instance.url",
		TokenType:   "Bearer",
		IssuedAt:    "1234567890",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse form: %s", err)
		}

		if r.Form.Get("grant_type") != "password" {
			t.Errorf("Expected grant_type=password, got %s", r.Form.Get("grant_type"))
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(mockResponse); err != nil {
			t.Errorf("Failed to encode: %s", err)
		}
	}))
	defer server.Close()

	auth := Auth{
		ClientID:     "mock_client_id",
		ClientSecret: "mock_client_secret",
		Username:     "mock_user",
		Password:     "mock_pass",
		TokenURL:     server.URL,
	}

	client, err := auth.AuthenticatePassword()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client.AccessToken != mockResponse.AccessToken {
		t.Errorf("Expected AccessToken %s, got %s", mockResponse.AccessToken, client.AccessToken)
	}
}

func TestAuthenticateClientCredentials(t *testing.T) {
	t.Parallel()
	mockResponse := Client{
		AccessToken: "mock_access_token",
		InstanceURL: "https://mock.instance.url",
		TokenType:   "Bearer",
		IssuedAt:    "1234567890",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		defer r.Body.Close()

		values, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}

		if values.Get("grant_type") != "client_credentials" {
			t.Errorf("Expected grant_type=client_credentials, got %s", values.Get("grant_type"))
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(mockResponse); err != nil {
			t.Errorf("Failed to encode: %s", err)
		}
	}))
	defer server.Close()

	auth := Auth{
		ClientID:     "mock_client_id",
		ClientSecret: "mock_client_secret",
		TokenURL:     server.URL,
	}

	client, err := auth.AuthenticateClientCredentials()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client.AccessToken != mockResponse.AccessToken {
		t.Errorf("Expected AccessToken %s, got %s", mockResponse.AccessToken, client.AccessToken)
	}
}
