package octoprint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new OctoPrint client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// PrinterState represents the current state of the printer
type PrinterState struct {
	Text  string `json:"text"`
	Flags struct {
		Operational bool `json:"operational"`
		Paused      bool `json:"paused"`
		Printing    bool `json:"printing"`
		Error       bool `json:"error"`
		Ready       bool `json:"ready"`
	} `json:"flags"`
}

// GetPrinterState gets the current printer state
func (c *Client) GetPrinterState() (*PrinterState, error) {
	req, err := c.newRequest("GET", "/api/printer", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		State PrinterState `json:"state"`
	}

	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.State, nil
}

// SetActiveSpool sets the active spool for a tool using our custom API
func (c *Client) SetActiveSpool(spoolID string, tool int) error {
	payload := map[string]interface{}{
		"command":  "set_spool",
		"spool_id": spoolID,
		"tool":     tool,
	}

	req, err := c.newRequest("POST", "/api/plugin/spoolman_api", payload)
	if err != nil {
		return err
	}

	var response map[string]interface{}
	if err := c.doRequest(req, &response); err != nil {
		return err
	}

	// Check if the API returned success
	if success, ok := response["success"].(bool); !ok || !success {
		if errMsg, ok := response["error"].(string); ok {
			return fmt.Errorf("API error: %s", errMsg)
		}
		return fmt.Errorf("failed to set spool")
	}

	return nil
}

// GetCurrentSpool gets the currently selected spool for a tool
func (c *Client) GetCurrentSpool(tool int) (string, error) {
	payload := map[string]interface{}{
		"command": "get_current_spool",
		"tool":    tool,
	}

	req, err := c.newRequest("POST", "/api/plugin/spoolman_api", payload)
	if err != nil {
		return "", err
	}

	var response struct {
		Success bool   `json:"success"`
		SpoolID string `json:"spool_id"`
		Error   string `json:"error,omitempty"`
	}

	if err := c.doRequest(req, &response); err != nil {
		return "", err
	}

	if !response.Success {
		return "", fmt.Errorf("API error: %s", response.Error)
	}

	return response.SpoolID, nil
}

// Helper methods
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) doRequest(req *http.Request, result interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}
