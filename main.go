package main

import (
	"build-go-microservice/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	helloHandler := handlers.NewHello(logger)
	gbHandler := handlers.NewGoodbye(logger)

	serveMux := http.NewServeMux()

	serveMux.Handle("/", helloHandler)
	serveMux.Handle("/goodbye", gbHandler)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// Create a channel using make
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	signal.Notify(channel, os.Kill)

	sig := <-channel
	logger.Println("\nSignal received", sig)

	timeoutCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeoutCtx)
}
