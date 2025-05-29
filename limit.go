package go_salesforce_api_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// LimitsResponse represents the response structure from Salesforce limits API
type LimitsResponse map[string]interface{}

// GetLimits retrieves the API usage limits from Salesforce
func (c *Client) GetLimits() (LimitsResponse, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	url := fmt.Sprintf("%s/services/data/v58.0/limits", c.InstanceURL)

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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to retrieve limits, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var limits LimitsResponse
	if err := json.Unmarshal(body, &limits); err != nil {
		return nil, err
	}

	return limits, nil
}
