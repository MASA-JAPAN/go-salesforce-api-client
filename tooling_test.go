package go_salesforce_api_client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueryToolingAPI(t *testing.T) {
	t.Parallel()
	mockResponse := ToolingResponse{
		TotalSize: 1,
		Done:      true,
		Records:   []map[string]interface{}{{"Id": "000000000000000000", "Name": "Test Metadata"}},
	}

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

	client := &Client{
		AccessToken: "mock_token",
		InstanceURL: server.URL,
	}

	soql := "SELECT Id, Name FROM MetadataComponent"
	resp, err := client.QueryToolingAPI(soql)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.TotalSize != mockResponse.TotalSize {
		t.Errorf("Expected TotalSize %d, got %d", mockResponse.TotalSize, resp.TotalSize)
	}

	if !resp.Done {
		t.Errorf("Expected Done to be true, got false")
	}

	if len(resp.Records) != len(mockResponse.Records) {
		t.Errorf("Expected %d records, got %d", len(mockResponse.Records), len(resp.Records))
	}

	if resp.Records[0]["Id"] != mockResponse.Records[0]["Id"] {
		t.Errorf("Expected record ID %s, got %s", mockResponse.Records[0]["Id"], resp.Records[0]["Id"])
	}
}

func TestCreateCustomField(t *testing.T) {
	t.Parallel()
	mockResponse := map[string]interface{}{
		"id":      "a1B3t000000XYZ",
		"success": true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(mockResponse); err != nil {
			t.Errorf("Failed to encode: %s", err)
		}
	}))
	defer server.Close()

	client := &Client{
		AccessToken: "mock_token",
		InstanceURL: server.URL,
	}

	fieldData := CustomField{
		FullName: "Account.Custom_Field__c",
		Metadata: struct {
			Label  string `json:"label"`
			Type   string `json:"type"`
			Length int    `json:"length,omitempty"`
		}{
			Label: "Custom Field",
			Type:  "Text",
		},
	}

	resp, err := client.CreateCustomField(fieldData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp["id"] != mockResponse["id"] {
		t.Errorf("Expected ID %s, got %s", mockResponse["id"], resp["id"])
	}
	if resp["success"] != mockResponse["success"] {
		t.Errorf("Expected success %t, got %t", mockResponse["success"], resp["success"])
	}
}
