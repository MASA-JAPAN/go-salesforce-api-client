package go_salesforce_api_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// CountResponse represents the response structure from Salesforce record count API
type CountResponse map[string]interface{}

// GetRecordCounts retrieves the record count for specified Salesforce objects
func (c *Client) GetRecordCounts(objects []string) (CountResponse, error) {
	if c.AccessToken == "" || c.InstanceURL == "" {
		return nil, errors.New("missing authentication details")
	}

	sObjectsParam := ""
	for i, obj := range objects {
		if i > 0 {
			sObjectsParam += ","
		}
		sObjectsParam += obj
	}

	url := fmt.Sprintf("%s/services/data/v58.0/limits/recordCount?sObjects=%s", c.InstanceURL, sObjectsParam)

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
		return nil, fmt.Errorf("failed to retrieve record counts, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var counts CountResponse
	if err := json.Unmarshal(body, &counts); err != nil {
		return nil, err
	}

	return counts, nil
}
