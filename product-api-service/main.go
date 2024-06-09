package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	currencyClient "github.com/arvindnama/golang-microservices/currency-service/protos"
	"github.com/arvindnama/golang-microservices/product-api-service/data"
	"github.com/arvindnama/golang-microservices/product-api-service/handlers"
	"google.golang.org/grpc"

	"github.com/go-openapi/runtime/middleware"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func getEnv(key string, defaultValue string) string {
	envValue := os.Getenv(key)
	if envValue == "" {
		return defaultValue
	}
	return envValue
}

func main() {

	logger := log.New(os.Stdout, "product-api-service", log.LstdFlags)
	validation := data.NewValidation()

	conn, err := grpc.Dial("localhost:9091", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	cc := currencyClient.NewCurrencyClient(conn)

	productsHandler := handlers.NewProducts(logger, validation, cc)

	serveMux := mux.NewRouter()
	getRouter := serveMux.Methods("GET").Subrouter()
	getRouter.HandleFunc("/products", productsHandler.GetAllProducts)
	getRouter.HandleFunc("/products/{id:[0-9]+}", productsHandler.GetProduct)

	putRouter := serveMux.Methods("PUT").Subrouter()
	putRouter.Use(productsHandler.MiddlewareValidateProduct)
	putRouter.HandleFunc("/products/{id:[0-9]+}", productsHandler.UpdateProduct)

	postRouter := serveMux.Methods("POST").Subrouter()
	postRouter.Use(productsHandler.MiddlewareValidateProduct)
	postRouter.HandleFunc("/products", productsHandler.AddProduct)

	deleteRouter := serveMux.Methods("DELETE").Subrouter()
	deleteRouter.HandleFunc("/product/{id:[0-9]+}", productsHandler.DeleteProduct)

	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}

	redocGetDocHandler := middleware.Redoc(ops, nil)

	getRouter.Handle("/docs", redocGetDocHandler)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	corsHandler := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"http://localhost:3000"}),
	)

	bindAddress := getEnv("BIND_ADDRESS", ":9090")
	server := &http.Server{
		Addr:         bindAddress,
		Handler:      corsHandler(serveMux),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		logger.Printf("Starting Http Server at %v\n", bindAddress)
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
