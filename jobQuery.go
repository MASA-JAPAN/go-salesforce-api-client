package go_salesforce_api_client

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// JobQueryResponse represents the response from Salesforce Bulk Query API
type JobQueryResponse struct {
	ID     string `json:"id"`
	State  string `json:"state"`
	Object string `json:"object"`
}

// JobQueryResult represents a single row in the job query results.
type JobQueryResult map[string]string

// CreateJobQuery initiates a Bulk Query Job in Salesforce
func (c *Client) CreateJobQuery(query string) (*JobQueryResponse, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/jobs/query", c.InstanceURL)

	requestBody := map[string]interface{}{
		"operation":   "query",
		"query":       query,
		"contentType": "CSV",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create job query, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jobResponse JobQueryResponse
	if err := json.Unmarshal(body, &jobResponse); err != nil {
		return nil, err
	}

	return &jobResponse, nil
}

// GetJobQuery retrieves the status and details of a Bulk Query Job in Salesforce
func (c *Client) GetJobQuery(jobID string) (*JobQueryResponse, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/jobs/query/%s", c.InstanceURL, jobID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to retrieve job query, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jobResponse JobQueryResponse
	if err := json.Unmarshal(body, &jobResponse); err != nil {
		return nil, err
	}

	return &jobResponse, nil
}

// GetJobQueryResults retrieves the job query results using pagination and maxRecords
func (c *Client) GetJobQueryResults(jobID, queryLocator string, maxRecords int) (string, string, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return "", "", errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/jobs/query/%s/results?maxRecords=%d", c.InstanceURL, jobID, maxRecords)
	if queryLocator != "" {
		url += fmt.Sprintf("&locator=%s", queryLocator)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("failed to retrieve job query results, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// Convert response body to string since Bulk API returns CSV, not JSON
	responseData := string(body)

	// Extract next locator from headers
	nextLocator := resp.Header.Get("Sforce-Locator")

	return responseData, nextLocator, nil
}

// GetJobQueryResultsParsed retrieves job query results and converts them into a structured format
func (c *Client) GetJobQueryResultsParsed(jobID, queryLocator string, maxRecords int) ([]JobQueryResult, string, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, "", errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/jobs/query/%s/results?maxRecords=%d", c.InstanceURL, jobID, maxRecords)
	if queryLocator != "" {
		url += fmt.Sprintf("&locator=%s", queryLocator)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("failed to retrieve job query results, status: %d, response: %s", resp.StatusCode, string(body))
	}

	// Parse CSV response
	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, "", err
	}

	if len(records) < 1 {
		return nil, "", errors.New("empty CSV response")
	}

	// Extract headers
	headers := records[0]

	// Convert to slice of maps
	var results []JobQueryResult
	for _, row := range records[1:] {
		entry := make(JobQueryResult)
		for i, value := range row {
			entry[headers[i]] = value
		}
		results = append(results, entry)
	}

	// Extract next locator from headers
	nextLocator := resp.Header.Get("Sforce-Locator")

	return results, nextLocator, nil
}

func (c *Client) AbortJobQuery(jobID string) error {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/jobs/query/%s", c.InstanceURL, jobID)

	requestBody := map[string]string{
		"state": "Aborted",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to abort job query, status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteJobQuery permanently removes a Bulk Query Job in Salesforce
func (c *Client) DeleteJobQuery(jobID string) error {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/jobs/query/%s", c.InstanceURL, jobID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete job query, status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}
