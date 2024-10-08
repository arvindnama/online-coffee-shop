package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/arvindnama/golang-microservices/product-images-service/config"
	"github.com/arvindnama/golang-microservices/product-images-service/handlers"

	"github.com/arvindnama/golang-microservices/product-images-service/files"

	goHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
)

func main() {

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "Product Images Service",
		Level: hclog.LevelFromString(config.Env.LogLevel),
	})

	stgLogger := logger.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	localStorage, err := files.NewLocalStorage(config.Env.ImagePath, 1024*1000*5)

	filesHandler := handlers.NewFilesHandler(logger, localStorage)

	if err != nil {
		logger.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	serverMux := mux.NewRouter()

	corsHandler := goHandlers.CORS(
		goHandlers.AllowedOrigins([]string{"http://localhost:3000"}),
	)
	postRouter := serverMux.NewRoute().Methods(http.MethodPost).Subrouter()

	postRouter.HandleFunc(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+.[a-z]{3}}",
		filesHandler.UploadREST,
	)
	postRouter.HandleFunc(
		"/",
		filesHandler.UploadMultipart,
	)

	getRouter := serverMux.NewRoute().Methods(http.MethodGet).Subrouter()

	gm := handlers.NewGzipHandler(logger)
	getRouter.Use(gm.GzipMiddleware)
	getRouter.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir(config.Env.ImagePath))),
	)

	bindAddress := config.Env.Address
	server := http.Server{
		Addr:         bindAddress,
		Handler:      corsHandler(serverMux),
		ErrorLog:     stgLogger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Info("Starting server", "bind_address", bindAddress)
		err := server.ListenAndServe()
		if err != nil {
			logger.Error("Unable to start server", err)
			os.Exit(1)
		}
	}()

	channel := make(chan os.Signal, 1)

	signal.Notify(channel, os.Interrupt)
	signal.Notify(channel, os.Kill)

	sig := <-channel

	logger.Info("Shutting down the server with ", "signal", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)

}
