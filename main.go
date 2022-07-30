package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Config struct {
	MetricsTick      time.Duration
	Port             int
	WineDataFileCSV  string
	WineDBFileSQLite string
	MaxRecords       int
}

type Metrics struct {
	keysStored int
	getOps     int
	putOps     int
	delOps     int
}

type Status struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Ts      string `json:"ts"`
}

type Kvd struct {
	config  *Config
	mutex   *sync.RWMutex
	store   map[string]byte
	metrics Metrics
	status  Status
}

var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

var ErrorNoSuchKey = errors.New("no such key")

func (kvd *Kvd) Init(c *Config) {
	kvd.config = c
	// if _, err := .wineData.InitDB(c); err != nil {
	// 	fmt.Println("Error: could not initialize wine data DB")
	// }
	kvd.mutex = &sync.RWMutex{}
	// kvd.metrics.Init(&wa.wineData, c)
}

func Get(key string) (string, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func Put(key string, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()

	return nil
}

func Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()

	return nil
}

func (kvd *Kvd) getStatus(w http.ResponseWriter, r *http.Request) {
	// wa.wineData.status.Ts = time.Now().Format(time.RFC3339)
	w.Header().Add("Content-Type", "application/json")
	// err := json.NewEncoder(w).Encode(wa.wineData.status)
	// if err != nil {
	// 	return
	// }
}

// keyValuePutHandler expects to be called with a PUT request for
// the "/v1/key/{key}" resource.
func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Retrieve "key" from the request
	key := vars["key"]

	value, err := io.ReadAll(r.Body) // The request body has our value
	defer r.Body.Close()

	if err != nil { // If we have an error, report it
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	err = Put(key, string(value)) // Store the value as a string
	if err != nil {               // If we have an error, report it
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // All good! Return StatusCreated
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Retrieve "key" from the request
	key := vars["key"]

	value, err := Get(key) // Get value for key
	if errors.Is(err, ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value)) // Write the value to the response
}

func helloMuxHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello gorilla/mux!\n"))
}

func (kvd *Kvd) StartService(ctx context.Context) (context.Context, error) {

func main() {
	r := mux.NewRouter()

	// Register keyValuePutHandler as the handler function for PUT
	// requests matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	r.HandleFunc("/status", getStatus).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}