package go_salesforce_api_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// SobjectResponse represents the generic Salesforce API response
type SobjectResponse struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Errors  []any  `json:"errors"`
}

// CreateRecord creates a new Salesforce record
func (c *Client) CreateRecord(objectType string, record map[string]interface{}) (*SobjectResponse, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/sobjects/%s/", c.InstanceURL, objectType)

	jsonData, err := json.Marshal(record)
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

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create record, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sfResp SobjectResponse
	if err := json.Unmarshal(body, &sfResp); err != nil {
		return nil, err
	}

	return &sfResp, nil
}

// GetRecord retrieves a Salesforce record by ID
func (c *Client) GetRecord(objectType, recordID string) (map[string]interface{}, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/sobjects/%s/%s", c.InstanceURL, objectType, recordID)

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
		return nil, fmt.Errorf("failed to retrieve record, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var record map[string]interface{}
	if err := json.Unmarshal(body, &record); err != nil {
		return nil, err
	}

	return record, nil
}

// UpdateRecord updates a Salesforce record by ID
func (c *Client) UpdateRecord(objectType, recordID string, updates map[string]interface{}) error {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/sobjects/%s/%s", c.InstanceURL, objectType, recordID)

	jsonData, err := json.Marshal(updates)
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

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update record, status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteRecord deletes a Salesforce record by ID
func (c *Client) DeleteRecord(objectType, recordID string) error {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/sobjects/%s/%s", c.InstanceURL, objectType, recordID)

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
		return fmt.Errorf("failed to delete record, status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DescribeSObject retrieves metadata for a given Salesforce object
func (c *Client) DescribeSObject(objectType string) (map[string]interface{}, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/sobjects/%s/describe", c.InstanceURL, objectType)

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
		return nil, fmt.Errorf("failed to retrieve SObject description, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var describe map[string]interface{}
	if err := json.Unmarshal(body, &describe); err != nil {
		return nil, err
	}

	return describe, nil
}
