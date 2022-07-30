package kvd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type Status struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Ts      string `json:"ts"`
}

type Config struct {
	Port       int
	MaxRecords int
}

type Kvd struct {
	config *Config
	db     DB
	status Status
}

func (kvd *Kvd) Init(c *Config) {
	kvd.config = c
	if err := kvd.db.Init(); err != nil {
		fmt.Println("Error: could not initialize wine data DB")
	}
	kvd.status.Status = "ok"
}

func (kvd *Kvd) statusHandler(w http.ResponseWriter, r *http.Request) {
	kvd.status.Ts = time.Now().Format(time.RFC3339)
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(kvd.status)
	if err != nil {
		return
	}
}

func (kvd *Kvd) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(kvd.db.metrics)
	if err != nil {
		return
	}
}

func (kvd *Kvd) keyGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Retrieve "key" from the request
	key := vars["key"]

	value, err := kvd.db.Get(key) // Get value for key
	if errors.Is(err, ErrorInvalidKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value)) // Write the value to the response
}

func (kvd *Kvd) keyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Retrieve "key" from the request
	key := vars["key"]

	value, err := kvd.db.Get(key) // Get value for key
	if errors.Is(err, ErrorInvalidKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deleteErr := kvd.db.Delete(key)
	if deleteErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value)) // Write the value to the response
}

// keyValuePutHandler expects to be called with a PUT request for
// the "/v1/key/{key}" resource.
func (kvd *Kvd) keyPutHandler(w http.ResponseWriter, r *http.Request) {
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

	err = kvd.db.Set(key, string(value)) // Store the value as a string
	if err != nil {                      // If we have an error, report it
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // All good! Return StatusCreated
}

func (kvd *Kvd) StartService(ctx context.Context) (context.Context, error) {
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter().StrictSlash(true)

	//router.HandleFunc("/status", kvd.getStatus).Methods("GET")
	router.HandleFunc("/v1/{key}", kvd.keyPutHandler).Methods("PUT")
	router.HandleFunc("/v1/{key}", kvd.keyGetHandler).Methods("GET")
	router.HandleFunc("/v1/{key}", kvd.keyDeleteHandler).Methods("DELETE")
	router.HandleFunc("/status", kvd.statusHandler).Methods("GET")
	router.HandleFunc("/metrics", kvd.metricsHandler).Methods("GET")

	serviceAddress := fmt.Sprintf("%v:%v", "0.0.0.0", kvd.config.Port)

	srv := &http.Server{
		Addr:         serviceAddress,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	go func() {
		fmt.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Printf("Kvd started. Press ctrl-c to stop.\n")
		<-c
		err := srv.Shutdown(ctx)
		if err != nil {
			fmt.Println(err)
		}
		cancel()
	}()

	return ctx, nil
}
