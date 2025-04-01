package go_salesforce_api_client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateRecord(t *testing.T) {
	mockResponse := SobjectResponse{ID: "1", Success: true}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(mockResponse); err != nil {
			t.Errorf("Failed to encode: %s", err)
		}
	}))
	defer server.Close()

	client := &Client{AccessToken: "mock_token", InstanceURL: server.URL}
	record := map[string]interface{}{"Name": "Test Record"}
	resp, err := client.CreateRecord("Account", record)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.ID != mockResponse.ID {
		t.Errorf("Expected ID %s, got %s", mockResponse.ID, resp.ID)
	}
}

func TestGetRecord(t *testing.T) {
	mockResponse := map[string]interface{}{"Id": "1", "Name": "Test Record"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(mockResponse); err != nil {
			t.Errorf("Failed to encode: %s", err)
		}
	}))
	defer server.Close()

	client := &Client{AccessToken: "mock_token", InstanceURL: server.URL}
	resp, err := client.GetRecord("Account", "1")
	fmt.Println(resp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestUpdateRecord(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &Client{AccessToken: "mock_token", InstanceURL: server.URL}
	updates := map[string]interface{}{"Name": "Updated Record"}
	err := client.UpdateRecord("Account", "1", updates)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteRecord(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &Client{AccessToken: "mock_token", InstanceURL: server.URL}
	err := client.DeleteRecord("Account", "1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDescribeSObject(t *testing.T) {
	mockResponse := map[string]interface{}{
		"name":      "Account",
		"label":     "Account",
		"keyPrefix": "001",
		"custom":    false,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(mockResponse); err != nil {
			t.Errorf("Failed to encode: %s", err)
		}
	}))
	defer server.Close()

	client := &Client{AccessToken: "mock_token", InstanceURL: server.URL}

	describe, err := client.DescribeSObject("Account")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if describe["name"].(string) != "Account" {
		t.Errorf("Expected name 'Account', got %s", describe["name"].(string))
	}
	if describe["label"].(string) != "Account" {
		t.Errorf("Expected label 'Account', got %s", describe["label"].(string))
	}
	if describe["keyPrefix"].(string) != "001" {
		t.Errorf("Expected keyPrefix '001', got %s", describe["keyPrefix"].(string))
	}
	if describe["custom"].(bool) != false {
		t.Errorf("Expected custom 'false', got %v", describe["custom"].(bool))
	}
}
