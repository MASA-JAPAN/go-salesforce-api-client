package go_salesforce_api_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ToolingResponse represents the response structure from Salesforce tooling API
type ToolingResponse struct {
	TotalSize int                      `json:"totalSize"`
	Done      bool                     `json:"done"`
	Records   []map[string]interface{} `json:"records"`
}

// QueryToolingAPI executes a SOQL query against the Salesforce Tooling API
func (c *Client) QueryToolingAPI(soql string) (*ToolingResponse, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	encodedSoql := url.QueryEscape(soql)
	url := fmt.Sprintf("%s/services/data/v58.0/tooling/query/?q=%s", c.InstanceURL, encodedSoql)

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
		return nil, fmt.Errorf("failed to execute tooling query, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var queryResp ToolingResponse
	if err := json.Unmarshal(body, &queryResp); err != nil {
		return nil, err
	}

	return &queryResp, nil
}
