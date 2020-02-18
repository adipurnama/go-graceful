package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

var buildVersion string
var healthStatus string

func main() {
	// port
	port := flag.Int64("port", 8081, "Insert http port for server")
	flag.Parse()

	buildVersion = os.Getenv("APP_GIT_BUILD_VERSION")

	router := http.NewServeMux()
	router.HandleFunc("/v1/order", createOrder)
	router.HandleFunc("/actuator/health", healthCheck)
	router.HandleFunc("/actuator/info", appInfo)

	healthStatus = "UP"
	s := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: router,
	}

	log.Printf("serving http service at port %d \n", *port)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("error serving http request")
	}
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	type order struct {
		ID string `json:"id"`
	}

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
