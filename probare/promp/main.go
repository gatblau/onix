package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	port := flag.String("p", "8080", "Specify port. Default is 8080")
	flag.Parse()
	r := mux.NewRouter()
	r.Handle("/", promhttp.Handler()).Methods("GET")
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", *port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  time.Second * 10,
		Handler:      r,
	}
	stop := make(chan os.Signal, 1)
	go func() {
		log.Info().Msgf("Prometheus Probe is listening on :%s", *port)
		if err := server.ListenAndServe(); err != nil {
			log.Info().Msgf("stopping the server: %v", err)
			os.Exit(1)
		}
	}()
	signal.Notify(stop, os.Interrupt)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("? I am shutting down due to an error: %v\n", err)
	}
}
