package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

	// Setup server
	e := echo.New()
	e.POST("/v1/order", createOrder)
	e.GET("/actuator/health", healthCheck)
	e.GET("/actuator/info", appInfo)

	// Start server
	healthStatus = "UP"
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", *port)); err != nil {
			log.Fatal("error serving http request", err)
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
	if err := e.Shutdown(sCtx); err != nil && err != http.ErrServerClosed {
		log.Fatal("found error while shutdown server")
	}
	log.Println("Bye.")
}

func createOrder(ctx echo.Context) error {
	log.Println("BEGIN - Processing Order")
	type order struct {
		ID string `json:"id"`
	}

	time.Sleep(5 * time.Second)
	log.Println("DONE - Processing Order")
	resp := order{
		ID: uuid.New().String(),
	}
	return ctx.JSON(http.StatusOK, resp)
}

func healthCheck(ctx echo.Context) error {
	type health struct {
		Status string `json:"status"`
	}

	resp := health{
		Status: healthStatus,
	}
	return ctx.JSON(http.StatusOK, resp)
}

func appInfo(ctx echo.Context) error {
	type info struct {
		Build string `json:"build"`
	}
	resp := info{
		Build: buildVersion,
	}
	return ctx.JSON(http.StatusOK, resp)
}
