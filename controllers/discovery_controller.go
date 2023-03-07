package controllers

import (
	"fmt"
	"log"
	"net/http"
)

type DiscoveryController struct {
	logger *log.Logger
}

// Same as NewDiscoveryController or NewDiscovery, but with explicit naming for the constructor
func DiscoveryControllerConstructor(logger *log.Logger) *DiscoveryController {
	return &DiscoveryController{logger: logger}
}

// Using default ServeHTTP signature for the API
func (controller *DiscoveryController) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	// Add the entpoint signature in the log for ease of identification
	controller.logger.Println("- [Discovery Controller]")

	// Write to the body of the response a message that the API is working
	fmt.Fprintf(responseWriter, "Identity API is working...")
}
