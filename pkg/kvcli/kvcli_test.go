package kvcli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockResponse represents a predefined response for the mock server
type MockResponse struct {
	StatusCode int
	Body       interface{}
	Headers    map[string]string
}

// SetupMockServer creates a test server that returns predefined responses
func SetupMockServer(t *testing.T, responses map[string]MockResponse) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method
		key := method + " " + path

		if response, ok := responses[key]; ok {
			// Set headers
			for k, v := range response.Headers {
				w.Header().Set(k, v)
			}

			// Set status code
			w.WriteHeader(response.StatusCode)

			// Write response body
			if response.Body != nil {
				switch body := response.Body.(type) {
				case string:
					w.Write([]byte(body))
				case []byte:
					w.Write(body)
				default:
					// Assume it's a JSON-serializable object
					json.NewEncoder(w).Encode(body)
				}
			}
		} else {
			t.Logf("No mock response found for %s", key)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Mock response not found"))
		}
	}))
}

func TestClientGet(t *testing.T) {
	// Define mock responses
	responses := map[string]MockResponse{
		"GET /v1/testkey": {
			StatusCode: http.StatusOK,
			Body:       "testvalue",
			Headers:    map[string]string{"Content-Type": "text/plain"},
		},
		"GET /v1/nonexistent": {
			StatusCode: http.StatusNotFound,
			Body:       "key not found",
			Headers:    map[string]string{"Content-Type": "text/plain"},
		},
	}

	server := SetupMockServer(t, responses)
	defer server.Close()

	client := NewClient(server.URL)

	// Test getting existing key
	value, err := client.Get("testkey")
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}
	if value != "testvalue" {
		t.Errorf("Expected value 'testvalue', got '%s'", value)
	}

	// Test getting non-existent key
	_, err = client.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent key, got nil")
	}
}

func TestClientSet(t *testing.T) {
	// Define mock responses
	responses := map[string]MockResponse{
		"PUT /v1/testkey": {
			StatusCode: http.StatusCreated,
			Headers:    map[string]string{"Content-Type": "text/plain"},
		},
		"PUT /v1/error": {
			StatusCode: http.StatusInternalServerError,
			Body:       "internal server error",
			Headers:    map[string]string{"Content-Type": "text/plain"},
		},
	}

	server := SetupMockServer(t, responses)
	defer server.Close()

	client := NewClient(server.URL)

	// Test successful set
	err := client.Set("testkey", "testvalue")
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	// Test set with server error
	err = client.Set("error", "value")
	if err == nil {
		t.Error("Expected error for server error, got nil")
	}
}

func TestClientDelete(t *testing.T) {
	// Define mock responses
	responses := map[string]MockResponse{
		"DELETE /v1/testkey": {
			StatusCode: http.StatusOK,
			Body:       "deleted",
			Headers:    map[string]string{"Content-Type": "text/plain"},
		},
		"DELETE /v1/nonexistent": {
			StatusCode: http.StatusNotFound,
			Body:       "key not found",
			Headers:    map[string]string{"Content-Type": "text/plain"},
		},
	}

	server := SetupMockServer(t, responses)
	defer server.Close()

	client := NewClient(server.URL)

	// Test successful delete
	err := client.Delete("testkey")
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	// Test delete non-existent key
	err = client.Delete("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent key, got nil")
	}
}

func TestClientBulkOperations(t *testing.T) {
	// Define mock data
	records := []struct {
		Key   string
		Value string
	}{
		{"bulk1", "value1"},
		{"bulk2", "value2"},
		{"bulk3", "value3"},
	}

	// Convert to JSON for response
	recordsJSON, _ := json.Marshal(records)

	// Define mock responses
	responses := map[string]MockResponse{
		"PUT /v1/": {
			StatusCode: http.StatusCreated,
			Headers:    map[string]string{"Content-Type": "application/json"},
		},
		"GET /v1/": {
			StatusCode: http.StatusOK,
			Body:       recordsJSON,
			Headers:    map[string]string{"Content-Type": "application/json"},
		},
		"DELETE /v1/": {
			StatusCode: http.StatusOK,
			Headers:    map[string]string{"Content-Type": "application/json"},
		},
	}

	server := SetupMockServer(t, responses)
	defer server.Close()

	client := NewClient(server.URL)

	// Create key-value pairs for bulk set
	kvPairs := make(map[string]string)
	for _, r := range records {
		kvPairs[r.Key] = r.Value
	}

	// Test bulk set
	err := client.BulkSet(kvPairs)
	if err != nil {
		t.Fatalf("Failed to bulk set: %v", err)
	}

	// Create key list for bulk get
	keys := make([]string, 0, len(records))
	for _, r := range records {
		keys = append(keys, r.Key)
	}

	// Test bulk get
	results, err := client.BulkGet(keys)
	if err != nil {
		t.Fatalf("Failed to bulk get: %v", err)
	}

	// Verify results
	if len(results) != len(records) {
		t.Errorf("Expected %d results, got %d", len(records), len(results))
	}

	for i, r := range records {
		if results[r.Key] != r.Value {
			t.Errorf("Entry %d: expected value '%s' for key '%s', got '%s'", i, r.Value, r.Key, results[r.Key])
		}
	}

	// Test bulk delete
	err = client.BulkDelete(keys)
	if err != nil {
		t.Fatalf("Failed to bulk delete: %v", err)
	}
}

func TestClientMetrics(t *testing.T) {
	// Define metrics response
	metrics := map[string]int64{
		"KeysStored":       10,
		"ValueBytesStored": 1024,
		"GetOps":           100,
		"SetOps":           50,
		"DelOps":           20,
	}

	// Define mock responses
	responses := map[string]MockResponse{
		"GET /metrics": {
			StatusCode: http.StatusOK,
			Body:       metrics,
			Headers:    map[string]string{"Content-Type": "application/json"},
		},
	}

	server := SetupMockServer(t, responses)
	defer server.Close()

	client := NewClient(server.URL)

	// Test getting metrics
	result, err := client.GetMetrics()
	if err != nil {
		t.Fatalf("Failed to get metrics: %v", err)
	}

	// Verify metrics
	if result.KeysStored != 10 {
		t.Errorf("Expected KeysStored to be 10, got %d", result.KeysStored)
	}
	if result.ValueBytesStored != 1024 {
		t.Errorf("Expected ValueBytesStored to be 1024, got %d", result.ValueBytesStored)
	}
	if result.GetOps != 100 {
		t.Errorf("Expected GetOps to be 100, got %d", result.GetOps)
	}
	if result.SetOps != 50 {
		t.Errorf("Expected SetOps to be 50, got %d", result.SetOps)
	}
	if result.DelOps != 20 {
		t.Errorf("Expected DelOps to be 20, got %d", result.DelOps)
	}
}