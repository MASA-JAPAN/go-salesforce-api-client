package go_salesforce_api_client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuery(t *testing.T) {
	mockResponse := QueryResponse{
		TotalSize: 1,
		Done:      true,
		Records:   []map[string]any{{"Id": "000000000000000000", "Name": "Test Record"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := &Client{AccessToken: "mock_token", InstanceURL: server.URL}
	query := "SELECT Id, Name FROM Account"
	resp, err := client.Query(query)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.TotalSize != mockResponse.TotalSize {
		t.Errorf("Expected TotalSize %d, got %d", mockResponse.TotalSize, resp.TotalSize)
	}

	if resp.Done != mockResponse.Done {
		t.Errorf("Expected Done %t, got %t", mockResponse.Done, resp.Done)
	}

	if len(resp.Records) != len(mockResponse.Records) {
		t.Errorf("Expected %d records, got %d", len(mockResponse.Records), len(resp.Records))
	}

}
