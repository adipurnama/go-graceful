package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
)

var buildVersion string
var healthStatus string

func main() {
	// port
	waitTimeS := flag.Int("wait", 10, "OUT_OF_SERVICE wait time in seconds")
	shutdownTimeS := flag.Int("shutdown", 10, "Shutdown timeout")
	port := flag.Int64("port", 8081, "Insert http port for server")
	flag.Parse()

	buildVersion = os.Getenv("APP_GIT_BUILD_VERSION")

	router := http.NewServeMux()
	router.HandleFunc("/v1/order", createOrder)
	router.HandleFunc("/actuator/health", healthCheck)
	router.HandleFunc("/actuator/info", appInfo)

	s := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: router,
	}

	healthStatus = "UP"
	go func() {
		log.Printf("serving http service at port %d \n", *port)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("error serving http request")
		}
	}()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quitCh
	log.Println("quit signal received..")

	// wait for LB see us as OUT_OF_SERVICE
	healthStatus = "OUT_OF_SERVICE"
	time.Sleep(time.Duration(*waitTimeS) * time.Second)

	// then shutdown server with timeout
	log.Println("Shutting down server..")
	sCtx, cancel := context.WithTimeout(context.Background(), time.Duration(*shutdownTimeS)*time.Second)
	defer cancel()
	if err := s.Shutdown(sCtx); err != nil && err != http.ErrServerClosed {
		log.Fatal("found error while shutdown server")
	}
	log.Println("Bye.")
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("BEGIN - Processing Order")
	type order struct {
		ID string `json:"id"`
	}

	time.Sleep(5 * time.Second)
	log.Println("DONE - Processing Order")
	resp := order{
		ID: uuid.New().String(),
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	type health struct {
		Status string `json:"status"`
	}

	resp := health{
		Status: healthStatus,
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func appInfo(w http.ResponseWriter, r *http.Request) {
	type info struct {
		Build string `json:"build"`
	}
	resp := info{
		Build: buildVersion,
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
