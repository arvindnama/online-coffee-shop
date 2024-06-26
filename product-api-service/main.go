package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	currencyClient "github.com/arvindnama/golang-microservices/currency-service/protos"
	"github.com/arvindnama/golang-microservices/product-api-service/data"
	"github.com/arvindnama/golang-microservices/product-api-service/handlers"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-openapi/runtime/middleware"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the service")
var currencyServiceAddress = env.String("CS_ADDRESS", false, "localhost:9092", "currency service address")

func main() {
	env.Parse()

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "Product Api Service",
		Level: hclog.Debug,
	})

	stdLogger := logger.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	conn, err := grpc.NewClient(
		*currencyServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	cc := currencyClient.NewCurrencyClient(conn)

	pDB := data.NewProductsDB(logger, cc)
	validation := data.NewValidation()
	productsHandler := handlers.NewProducts(logger, validation, pDB)

	serveMux := mux.NewRouter()
	getRouter := serveMux.Methods("GET").Subrouter()
	getRouter.HandleFunc("/products", productsHandler.GetAllProducts).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/products", productsHandler.GetAllProducts)
	getRouter.HandleFunc("/products/{id:[0-9]+}", productsHandler.GetProduct).Queries("currency", "{[A-Z]{3}}")
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
		gorillaHandlers.AllowedOrigins([]string{"*"}),
	)

	server := &http.Server{
		Addr:         *bindAddress,
		Handler:      corsHandler(serveMux),
		ErrorLog:     stdLogger,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting Http Server at %#v\n", *bindAddress))
		err := server.ListenAndServe()
		if err != nil {
			logger.Error("Error Starting Http Server", err)
		}
	}()

	// Create a channel using make
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	signal.Notify(channel, os.Kill)

	sig := <-channel
	logger.Info("\nSignal received", sig)

	timeoutCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeoutCtx)
}
