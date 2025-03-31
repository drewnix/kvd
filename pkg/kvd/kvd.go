package kvd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// Status represents the current state of the KVD server
type Status struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Ts      string `json:"ts"`
	Version string `json:"version"`
}

// Config represents the configuration for the KVD server
type Config struct {
	Port       int
	MaxRecords int
	Host       string
	LogLevel   string
}

// Kvd represents the KVD server instance
type Kvd struct {
	config *Config
	db     DB
	status Status
	logger *log.Logger
}

// Record represents a key-value pair
type Record struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// DefaultConfig returns a default configuration for the KVD server
func DefaultConfig() *Config {
	return &Config{
		Port:       8080,
		MaxRecords: 10000,
		Host:       "0.0.0.0",
		LogLevel:   "info",
	}
}

// Init initializes the KVD server
func (kvd *Kvd) Init(c *Config) error {
	kvd.config = c
	if kvd.config == nil {
		kvd.config = DefaultConfig()
	}
	
	// Initialize logger
	kvd.logger = log.New(os.Stdout, "KVD: ", log.LstdFlags)
	
	// Initialize database
	if err := kvd.db.Init(); err != nil {
		kvd.logger.Printf("Error initializing database: %v", err)
		return fmt.Errorf("could not initialize database: %w", err)
	}
	
	// Set server status
	kvd.status = Status{
		Status:  "ok",
		Message: "Server initialized",
		Version: "1.0.0",
	}
	
	return nil
}

// toJSON converts an object to JSON bytes
func (kvd *Kvd) toJSON(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize: %w", err)
	}
	return b.Bytes(), nil
}

// statusHandler handles requests for server status
func (kvd *Kvd) statusHandler(w http.ResponseWriter, r *http.Request) {
	kvd.status.Ts = time.Now().Format(time.RFC3339)
	w.Header().Set("Content-Type", "application/json")
	
	if err := json.NewEncoder(w).Encode(kvd.status); err != nil {
		kvd.logger.Printf("Error encoding status response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// metricsHandler handles requests for server metrics
func (kvd *Kvd) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	if err := json.NewEncoder(w).Encode(kvd.db.metrics); err != nil {
		kvd.logger.Printf("Error encoding metrics response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// keyGetHandler handles requests to get a single key
func (kvd *Kvd) keyGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := kvd.db.Get(key)
	if errors.Is(err, ErrKeyNotFound) || errors.Is(err, ErrInvalidKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		kvd.logger.Printf("Error getting key %s: %v", key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(value)); err != nil {
		kvd.logger.Printf("Error writing response: %v", err)
	}
}

// keyDeleteHandler handles requests to delete a single key
func (kvd *Kvd) keyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	// Get the value before deleting it
	value, err := kvd.db.Get(key)
	if errors.Is(err, ErrKeyNotFound) || errors.Is(err, ErrInvalidKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		kvd.logger.Printf("Error getting key %s: %v", key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := kvd.db.Delete(key); err != nil {
		kvd.logger.Printf("Error deleting key %s: %v", key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(value)); err != nil {
		kvd.logger.Printf("Error writing response: %v", err)
	}
}

// keyPutHandler handles requests to set a single key
func (kvd *Kvd) keyPutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		http.Error(w, "Key cannot be empty", http.StatusBadRequest)
		return
	}

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		kvd.logger.Printf("Error reading request body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := kvd.db.Set(key, string(value)); err != nil {
		kvd.logger.Printf("Error setting key %s: %v", key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// keyManyPutHandler handles bulk key set operations
func (kvd *Kvd) keyManyPutHandler(w http.ResponseWriter, r *http.Request) {
	var records []Record

	if r.Body == nil {
		http.Error(w, "Request body is required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	
	if err != nil {
		kvd.logger.Printf("Error reading request body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &records); err != nil {
		kvd.logger.Printf("Error unmarshaling request: %v", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if len(records) == 0 {
		http.Error(w, "No records provided", http.StatusBadRequest)
		return
	}

	if err := kvd.db.BulkSet(records); err != nil {
		kvd.logger.Printf("Error in bulk set operation: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// keyManyGetHandler handles bulk key get operations
func (kvd *Kvd) keyManyGetHandler(w http.ResponseWriter, r *http.Request) {
	var keys []string

	if r.Body == nil {
		http.Error(w, "Request body is required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	
	if err != nil {
		kvd.logger.Printf("Error reading request body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &keys); err != nil {
		kvd.logger.Printf("Error unmarshaling request: %v", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if len(keys) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
		return
	}

	records, err := kvd.db.BulkGet(keys)
	if err != nil {
		kvd.logger.Printf("Error in bulk get operation: %v", err)
		status := http.StatusInternalServerError
		if errors.Is(err, ErrKeyNotFound) || errors.Is(err, ErrInvalidKey) {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(records); err != nil {
		kvd.logger.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// keyManyDeletesHandler handles bulk key delete operations
func (kvd *Kvd) keyManyDeletesHandler(w http.ResponseWriter, r *http.Request) {
	var keys []string

	if r.Body == nil {
		http.Error(w, "Request body is required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	
	if err != nil {
		kvd.logger.Printf("Error reading request body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &keys); err != nil {
		kvd.logger.Printf("Error unmarshaling request: %v", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if len(keys) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := kvd.db.BulkDelete(keys); err != nil {
		kvd.logger.Printf("Error in bulk delete operation: %v", err)
		status := http.StatusInternalServerError
		if errors.Is(err, ErrKeyNotFound) || errors.Is(err, ErrInvalidKey) {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// StartService starts the KVD HTTP server
func (kvd *Kvd) StartService(ctx context.Context) (context.Context, error) {
	ctx, cancel := context.WithCancel(ctx)
	
	// Set up signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Create and configure router
	router := mux.NewRouter().StrictSlash(true)

	// API routes
	router.HandleFunc("/v1/", kvd.keyManyPutHandler).Methods(http.MethodPut)
	router.HandleFunc("/v1/", kvd.keyManyGetHandler).Methods(http.MethodGet)
	router.HandleFunc("/v1/", kvd.keyManyDeletesHandler).Methods(http.MethodDelete)
	router.HandleFunc("/v1/{key}", kvd.keyPutHandler).Methods(http.MethodPut)
	router.HandleFunc("/v1/{key}", kvd.keyGetHandler).Methods(http.MethodGet)
	router.HandleFunc("/v1/{key}", kvd.keyDeleteHandler).Methods(http.MethodDelete)
	
	// Admin routes
	router.HandleFunc("/status", kvd.statusHandler).Methods(http.MethodGet)
	router.HandleFunc("/metrics", kvd.metricsHandler).Methods(http.MethodGet)

	// Configure server address
	host := kvd.config.Host
	if host == "" {
		host = "0.0.0.0"
	}
	serviceAddress := fmt.Sprintf("%s:%d", host, kvd.config.Port)

	// Create HTTP server
	srv := &http.Server{
		Addr:         serviceAddress,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	// Start HTTP server
	go func() {
		kvd.logger.Printf("Starting KVD server on %s", serviceAddress)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			kvd.logger.Printf("HTTP server error: %v", err)
		}
		cancel()
	}()

	// Handle graceful shutdown
	go func() {
		kvd.logger.Println("KVD started. Press Ctrl+C to stop.")
		<-signalChan
		kvd.logger.Println("Shutdown signal received, stopping server...")
		
		// Create a shutdown timeout context
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
		defer shutdownCancel()
		
		if err := srv.Shutdown(shutdownCtx); err != nil {
			kvd.logger.Printf("Server shutdown error: %v", err)
		}
		
		kvd.logger.Println("Server stopped")
		cancel()
	}()

	return ctx, nil
}
