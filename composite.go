package go_salesforce_api_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// CompositeResponse represents the generic Salesforce API response
type CompositeResponse struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Errors  []any  `json:"errors"`
}

// CreateRecords creates multiple Salesforce records
func (c *Client) CreateRecords(objectType string, records []map[string]interface{}) ([]CompositeResponse, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/composite/sobjects", c.InstanceURL)

	// Ensure each record has an "attributes" field
	for i := range records {
		records[i]["attributes"] = map[string]string{"type": objectType}
	}

	requestBody := map[string]interface{}{
		"allOrNone": true,
		"records":   records,
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
		return nil, fmt.Errorf("failed to  create records, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responses []CompositeResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		return nil, err
	}

	return responses, nil
}

// UpdateRecords updates multiple Salesforce records
func (c *Client) UpdateRecords(objectType string, records []map[string]interface{}) error {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/composite/sobjects", c.InstanceURL)

	for i := range records {
		records[i]["attributes"] = map[string]string{"type": objectType}
	}

	requestBody := map[string]interface{}{
		"allOrNone": true,
		"records":   records,
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

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to  update records, status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteRecords deletes multiple Salesforce records
func (c *Client) DeleteRecords(objectType string, recordIDs []string) error {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/composite", c.InstanceURL)

	compositeRequests := []map[string]interface{}{}
	for _, id := range recordIDs {
		compositeRequests = append(compositeRequests, map[string]interface{}{
			"method":      "DELETE",
			"url":         fmt.Sprintf("/services/data/v58.0/sobjects/%s/%s", objectType, id),
			"referenceId": id,
		})
	}

	requestBody := map[string]interface{}{
		"allOrNone":        true,
		"compositeRequest": compositeRequests,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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
		return fmt.Errorf("failed to bulk delete records, status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}
