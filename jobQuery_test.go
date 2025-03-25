package go_salesforce_api_client_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	go_salesforce_api_client "github.com/MASA-JAPAN/go-salesforce-api-client"
)

func TestCreateJobQuery(t *testing.T) {
	mockResponse := go_salesforce_api_client.JobQueryResponse{
		ID:     "7501X00000XXXXXQAQ",
		State:  "Open",
		Object: "Account",
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{AccessToken: "test_token", InstanceURL: server.URL}
	resp, err := client.CreateJobQuery("SELECT Id FROM Account")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if resp.ID != mockResponse.ID {
		t.Errorf("Expected ID %s, got: %s", mockResponse.ID, resp.ID)
	}
}

func TestGetJobQuery(t *testing.T) {
	mockResponse := go_salesforce_api_client.JobQueryResponse{
		ID:     "7501X00000XXXXXQAQ",
		State:  "InProgress",
		Object: "Account",
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{AccessToken: "test_token", InstanceURL: server.URL}
	resp, err := client.GetJobQuery("7501X00000XXXXXQAQ")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if resp.State != mockResponse.State {
		t.Errorf("Expected State %s, got: %s", mockResponse.State, resp.State)
	}
}

func TestGetJobQueryResultsParsed(t *testing.T) {
	mockCSV := "Id,Name\n001ABC,Acme Corp\n002DEF,Global Inc"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Sforce-Locator", "")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, mockCSV)
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{AccessToken: "test_token", InstanceURL: server.URL}
	results, locator, err := client.GetJobQueryResultsParsed("job_id", "", 1000)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got: %d", len(results))
	}
	if locator != "" {
		t.Errorf("Expected empty locator, got: %s", locator)
	}
}

func TestAbortJobQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH, got: %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{AccessToken: "test_token", InstanceURL: server.URL}
	err := client.AbortJobQuery("job_id")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteJobQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got: %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{AccessToken: "test_token", InstanceURL: server.URL}
	err := client.DeleteJobQuery("job_id")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestGetJobQueryResultsParsed_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "")
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{
		AccessToken: "test_token",
		InstanceURL: server.URL,
	}

	_, _, err := client.GetJobQueryResultsParsed("dummy_job_id", "", 1000)
	if err == nil {
		t.Error("Expected error for empty CSV response, got nil")
	}
}

func TestGetJobQueryResultsParsed_InvalidCSV(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "bad,data\nnot,enough,columns")
	}))
	defer server.Close()

	client := &go_salesforce_api_client.Client{
		AccessToken: "test_token",
		InstanceURL: server.URL,
	}

	_, _, err := client.GetJobQueryResultsParsed("dummy_job_id", "", 1000)
	if err == nil {
		t.Error("Expected CSV parse error, got nil")
	}
}
