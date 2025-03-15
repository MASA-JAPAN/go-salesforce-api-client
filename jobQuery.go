package go_salesforce_api_client

import (
	"bytes"
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

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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

	if resp.StatusCode != http.StatusCreated {
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

	req, err := http.NewRequest("GET", url, nil)
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

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
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

	req, err := http.NewRequest("DELETE", url, nil)
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

func (c *Client) GetJobQueryResultPages(jobID string) ([]string, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/jobs/query/%s/resultPages", c.InstanceURL, jobID)

	req, err := http.NewRequest("GET", url, nil)
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
		return nil, fmt.Errorf("failed to retrieve job query result pages, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var resultPages []string
	if err := json.Unmarshal(body, &resultPages); err != nil {
		return nil, err
	}

	return resultPages, nil
}

// GetFilteredJobQueries retrieves job queries with filtering options
func (c *Client) GetFilteredJobQueries(isPkChunkingEnabled bool, jobType, concurrencyMode, queryLocator string) ([]JobQueryResponse, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/jobs/query?isPkChunkingEnabled=%t&jobType=%s&concurrencyMode=%s&queryLocator=%s",
		c.InstanceURL, isPkChunkingEnabled, jobType, concurrencyMode, queryLocator)

	req, err := http.NewRequest("GET", url, nil)
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
		return nil, fmt.Errorf("failed to retrieve filtered job queries, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jobResponses []JobQueryResponse
	if err := json.Unmarshal(body, &jobResponses); err != nil {
		return nil, err
	}

	return jobResponses, nil
}
