package go_salesforce_api_client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateRecords(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode([]CompositeResponse{{ID: "000000000000000000", Success: true}}); err != nil {
			t.Errorf("Failed to encode: %s", err)
		}
	}))
	defer server.Close()

	client := &Client{AccessToken: "mock_token", InstanceURL: server.URL}
	records := []map[string]interface{}{{"Name": "Test Record"}}
	_, err := client.CreateRecords("Account", records)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestUpdateRecords(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &Client{AccessToken: "mock_token", InstanceURL: server.URL}
	records := []map[string]interface{}{{"Id": "000000000000000000", "Name": "Updated Record"}}
	err := client.UpdateRecords("Account", records)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteRecords(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{AccessToken: "mock_token", InstanceURL: server.URL}
	err := client.DeleteRecords("Account", []string{"000000000000000000"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
