package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func (wa *WineAPI) StartService(ctx context.Context) (context.Context, error) {
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	if wa.config.WineDataFileCSV != "" {
		err := wa.wineData.LoadData()
		if err != nil {
			fmt.Println("Failed to load CSV data: ", err)
		}
	}

	router := mux.NewRouter().StrictSlash(true)

	wa.metrics.metricsLogger(ctx)
	router.Use(wa.metrics.metricsMiddleware)

	router.HandleFunc("/wine", wa.getAllWine).Methods("GET")
	router.HandleFunc("/wine/{id}", wa.getOneWine).Methods("GET")
	router.HandleFunc("/wine", wa.createWine).Methods("PUT")
	router.HandleFunc("/status", wa.getStatus).Methods("GET")

	serviceAddress := fmt.Sprintf("%v:%v", "0.0.0.0", wa.config.Port)

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
		fmt.Printf("WineAPI started. Press ctrl-c to stop.\n")
		<-c
		err := srv.Shutdown(ctx)
		if err != nil {
			fmt.Println(err)
		}
		cancel()
	}()

	return ctx, nil
}
