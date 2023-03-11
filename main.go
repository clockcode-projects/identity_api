// Copyright (C) 2023 Fernando Silva
// SPDX-License-Identifier: BSD-3-Clause
//
// Portions of the source code from this file are reused from
// building-microservices-youtube repository by Nic Jackson
// Copyright (c) 2020 Nicholas Jackson
// SPDX-License-Identifier: MIT
// Origin: https://github.com/nicholasjackson/building-microservices-youtube/blob/main/product-api/main.go

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/clockcode-projects/identity_api/controllers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	logger := log.New(os.Stdout, "[Identity API] - ", log.LstdFlags)

	logger.Println("| Starting Identity API")
	logger.Println("| ")
	logger.Println("| Initializing environment variables...")

	// Try to load environment variables
	variablesError := godotenv.Load()
	if variablesError != nil {
		logger.Println("| Not possible to load variables. Required settings not present. Please, verify and try again.")
		logger.Fatal("| Error: ", variablesError)
		os.Exit(1)
	}

	app_environment := os.Getenv("APP_ENVIRONMENT")
	logger.Printf("| Application environment is set to: %s", app_environment)

	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_sslm := os.Getenv("DB_SSLM")
	db_cert := os.Getenv("DB_CERT")

	// Parse connecting string
	switch app_environment {
	case "development":
		db_name = db_name + "_development"
		db_user = db_user + "_development"
	case "staging":
		db_name = db_name + "_staging"
		db_user = db_user + "_staging"
	case "production":
		db_name = db_name + "_production"
		db_user = db_user + "_production"
	default:
		logger.Fatal("| Error: Only legal environment names are development, staging and production. Please, verify and try again.")
		os.Exit(1)
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db_user, db_pass, db_host, db_port, db_name)

	if db_sslm != "" {
		connectionString = connectionString + "?sslmode=" + db_sslm
	}

	if db_cert != "" {
		connectionString = connectionString + "&sslrootcert=" + db_cert
	}

	// Setting application context
	applicationContext := context.Background()

	// Instantiate pgxpool as dataContext
	dataContext, dbError := pgxpool.New(applicationContext, connectionString)
	if dbError != nil {
		logger.Println("| It was not possible to create connection pool. Please, verify and try again.")
		logger.Fatal("| Error: ", dbError)
		os.Exit(1)
	}

	dbConnectError := dataContext.Ping(applicationContext)
	if dbConnectError != nil {
		logger.Println("| It was not possible to connect to the database. Please, verify and try again.")
		logger.Fatal("| Error: ", dbConnectError)
		os.Exit(1)
	}
	defer dataContext.Close()

	// Try to get database version via SQL query
	var databaseServerVersion string
	checkVersionError := dataContext.QueryRow(applicationContext, "SELECT version()").Scan(&databaseServerVersion)
	if checkVersionError != nil {
		logger.Println("| It was not possible to check server version. Please, verify and try again.")
		logger.Printf("| Error: %v", checkVersionError)
	}
	logger.Printf("| Database server is using version: %s", databaseServerVersion)

	// API configuration
	apiVersion := "v1"
	apiPrefix := "/api/" + apiVersion

	// Declaration of endpoints with dependency injection
	discoveryController := controllers.NewDiscoveryController(logger)

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
