package go_salesforce_api_client_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	go_salesforce_api_client "github.com/MASA-JAPAN/go-salesforce-api-client"
)

func TestGetRecordCounts(t *testing.T) {
	// Mock Salesforce API response
	mockResponse := `{
		"sobjects": {
			"Account": {"count": 1500},
			"Contact": {"count": 2500},
			"Lead": {"count": 1000}
		}
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	client := go_salesforce_api_client.Client{
		AccessToken: "mock_access_token",
		InstanceURL: ts.URL,
	}

	objects := []string{"Account", "Contact", "Lead"}
	recordCounts, err := client.GetRecordCounts(objects)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var expected map[string]interface{}
	if err := json.Unmarshal([]byte(mockResponse), &expected); err != nil {
		t.Fatalf("error unmarshalling expected response: %v", err)
	}

	for _, obj := range objects {
		if recordCounts["sobjects"].(map[string]interface{})[obj].(map[string]interface{})["count"] != expected["sobjects"].(map[string]interface{})[obj].(map[string]interface{})["count"] {
			t.Errorf("expected %v for %s, got %v", expected["sobjects"].(map[string]interface{})[obj], obj, recordCounts["sobjects"].(map[string]interface{})[obj])
		}
	}
}
