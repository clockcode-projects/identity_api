package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/clockcode-projects/identity_api/controllers"
)

func main() {
	// Create a logger instance
	logger := log.New(os.Stdout, "[Identity API] - ", log.LstdFlags)

	// Header for the log output
	logger.Println("| Starting Identity API")

	// Setting application context
	applicationContext := context.Background()

	// API configuration
	apiVersion := "v1"
	apiPrefix := "/api/" + apiVersion

	// Declaration of endpoints with dependency injection
	discoveryController := controllers.DiscoveryControllerConstructor(logger)

	// Router
	router := http.NewServeMux()
	router.Handle(apiPrefix+"/discovery", discoveryController)

	// Server configuration
	server := &http.Server{
		Addr:         ":80",
		Handler:      router,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Non-blocking startup
	go func() {
		serverError := server.ListenAndServe()
		if serverError != nil {
			logger.Println("| It was not possible to start server. Verify your configuration and try again.")
			logger.Fatal(serverError)
		}
	}()

	// In case of server started successfully, print the port which it is listening for connections on.
	logger.Printf("| Server started using port %s.\n", strings.Split(server.Addr, ":")[1])

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	signal := <-signalChannel
	logger.Println("| ")
	logger.Println("| Application termination request received. Attempting graceful shutdown.")
	logger.Println("| Signal type: ", signal)
	logger.Println("| ")

	timeoutContext, _ := context.WithTimeout(applicationContext, 60*time.Second)
	server.Shutdown(timeoutContext)

}
