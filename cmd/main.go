package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func main() {
	// port
	port := flag.Int64("port", 8081, "Insert http port for server")
	flag.Parse()

	router := http.NewServeMux()
	router.HandleFunc("/v1/order", createOrder)

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
