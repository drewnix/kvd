// Package kvcli provides a client library for interacting with the KVD server
package kvcli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/drewnix/kvd/pkg/kvd"
)

// Client represents a KVD client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Metrics represents the metrics returned by the KVD server
type Metrics struct {
	KeysStored       int64 `json:"KeysStored"`
	ValueBytesStored int64 `json:"ValueBytesStored"`
	GetOps           int64 `json:"GetOps"`
	SetOps           int64 `json:"SetOps"`
	DelOps           int64 `json:"DelOps"`
}

// NewClient creates a new KVD client
func NewClient(serverURL string) *Client {
	return &Client{
		baseURL: serverURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Get retrieves a value for a given key
func (c *Client) Get(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key cannot be empty")
	}

	url := fmt.Sprintf("%s/v1/%s", c.baseURL, key)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("server returned error: %s (status: %d)", body, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// BulkGet retrieves multiple key-value pairs
func (c *Client) BulkGet(keys []string) (map[string]string, error) {
	if len(keys) == 0 {
		return make(map[string]string), nil
	}

	url := fmt.Sprintf("%s/v1/", c.baseURL)
	jsonData, err := json.Marshal(keys)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal keys: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned error: %s (status: %d)", body, resp.StatusCode)
	}

	var records []kvd.Record
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to map for easier usage
	result := make(map[string]string, len(records))
	for _, record := range records {
		result[record.Key] = record.Value
	}

	return result, nil
}

// Set sets a value for a given key
func (c *Client) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	url := fmt.Sprintf("%s/v1/%s", c.baseURL, key)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(value))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned error: %s (status: %d)", body, resp.StatusCode)
	}

	return nil
}

// BulkSet sets multiple key-value pairs
func (c *Client) BulkSet(kvPairs map[string]string) error {
	if len(kvPairs) == 0 {
		return nil
	}

	url := fmt.Sprintf("%s/v1/", c.baseURL)
	
	// Convert map to records
	records := make([]kvd.Record, 0, len(kvPairs))
	for key, value := range kvPairs {
		records = append(records, kvd.Record{
			Key:   key,
			Value: value,
		})
	}

	jsonData, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("failed to marshal records: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned error: %s (status: %d)", body, resp.StatusCode)
	}

	return nil
}

// Delete removes a key-value pair
func (c *Client) Delete(key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	url := fmt.Sprintf("%s/v1/%s", c.baseURL, key)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned error: %s (status: %d)", body, resp.StatusCode)
	}

	return nil
}

// BulkDelete removes multiple key-value pairs
func (c *Client) BulkDelete(keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	url := fmt.Sprintf("%s/v1/", c.baseURL)
	jsonData, err := json.Marshal(keys)
	if err != nil {
		return fmt.Errorf("failed to marshal keys: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned error: %s (status: %d)", body, resp.StatusCode)
	}

	return nil
}

// GetMetrics retrieves metrics from the server
func (c *Client) GetMetrics() (*Metrics, error) {
	url := fmt.Sprintf("%s/metrics", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned error: %s (status: %d)", body, resp.StatusCode)
	}

	var metrics Metrics
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &metrics, nil
}

// WithTimeout returns a context with a timeout
func (c *Client) WithTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}
