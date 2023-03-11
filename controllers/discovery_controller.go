// Copyright (C) 2023 Fernando Silva
// SPDX-License-Identifier: BSD-3-Clause

package controllers

import (
	"fmt"
	"log"
	"net/http"
)

type DiscoveryController struct {
	logger *log.Logger
}

// Constructor
func NewDiscoveryController(logger *log.Logger) *DiscoveryController {
	return &DiscoveryController{logger: logger}
}

// API Entrypoint
func (controller *DiscoveryController) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	controller.logger.Println("| Request served for Discovery Controller")
	fmt.Fprintf(responseWriter, "Identity API is working...")
}
