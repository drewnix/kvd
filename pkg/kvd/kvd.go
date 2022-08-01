package kvd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

type Record struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func (kvd *Kvd) Init(c *Config) {
	kvd.config = c
	if err := kvd.db.Init(); err != nil {
		fmt.Println("Error: could not initialize wine data DB")
	}
	kvd.status.Status = "ok"
}

func (kvd *Kvd) toJSON(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize : %q", err)
	}
	return b.Bytes(), nil
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
	if errors.Is(err, ErrInvalidKey) {
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
	if errors.Is(err, ErrInvalidKey) {
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

func (kvd *Kvd) keyManyPutHandler(w http.ResponseWriter, r *http.Request) {
	var records []Record = make([]Record, 0)

	if r.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := r.Body.Close(); err != nil {
		fmt.Println("Could not close body IO: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &records); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(err); err != nil {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
	}

	err = kvd.db.BulkSet(records) // Store the value as a string
	if err != nil {               // If we have an error, report it
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // All good! Return StatusCreated
}

func (kvd *Kvd) keyManyGetHandler(w http.ResponseWriter, r *http.Request) {
	//var records []Record = make([]Record, 0)
	var query []string

	if r.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := r.Body.Close(); err != nil {
		fmt.Println("Could not close body IO: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &query); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
	}

	recs, err := kvd.db.BulkGet(query) // Store the value as a string
	if err != nil {                    // If we have an error, report it
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}
	data, err := kvd.toJSON(recs)
	if err != nil {
		fmt.Println("Could not convert to JSON: ", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func (kvd *Kvd) keyManyDeletesHandler(w http.ResponseWriter, r *http.Request) {
	var query []string

	if r.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := r.Body.Close(); err != nil {
		fmt.Println("Could not close body IO: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &query); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(err); err != nil {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
	}

	err = kvd.db.BulkDelete(query) // Store the value as a string
	if err != nil {                // If we have an error, report it
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) // All good! Return StatusOK
}

func (kvd *Kvd) StartService(ctx context.Context) (context.Context, error) {
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/v1/", kvd.keyManyPutHandler).Methods("PUT")
	router.HandleFunc("/v1/", kvd.keyManyGetHandler).Methods("GET")
	router.HandleFunc("/v1/", kvd.keyManyDeletesHandler).Methods("DELETE")
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
