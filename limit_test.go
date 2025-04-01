package go_salesforce_api_client_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	go_salesforce_api_client "github.com/MASA-JAPAN/go-salesforce-api-client"
)

func TestGetLimits(t *testing.T) {
	t.Parallel()
	// Mock Salesforce API response
	mockResponse := `{
		"DailyApiRequests": {
			"Max": 5000000,
			"Remaining": 4999990
		}
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponse)); err != nil {
				t.Errorf("Failed to write: %s", err)
		}
	}))
	defer ts.Close()

	client := go_salesforce_api_client.Client{
		AccessToken: "mock_access_token",
		InstanceURL: ts.URL,
	}

	limits, err := client.GetLimits()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var expected map[string]interface{}
	if err := json.Unmarshal([]byte(mockResponse), &expected); err != nil {
		t.Fatalf("error unmarshalling expected response: %v", err)
	}

	if limits["DailyApiRequests"].(map[string]interface{})["Remaining"] != expected["DailyApiRequests"].(map[string]interface{})["Remaining"] {
		t.Errorf("expected %v, got %v", expected, limits)
	}
}
