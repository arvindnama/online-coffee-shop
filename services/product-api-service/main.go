package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	currencyClient "github.com/arvindnama/golang-microservices/currency-service/protos"
	dataUtils "github.com/arvindnama/golang-microservices/libs/utils/data-utils"
	"github.com/arvindnama/golang-microservices/product-api-service/config"
	"github.com/arvindnama/golang-microservices/product-api-service/data"
	"github.com/arvindnama/golang-microservices/product-api-service/handlers"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-openapi/runtime/middleware"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "Product Api Service",
		Level: hclog.LevelFromString(config.ENV.LogLevel),
	})

	stdLogger := logger.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	conn, err := grpc.NewClient(
		config.ENV.CurrencyServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	cc := currencyClient.NewCurrencyClient(conn)

	pDB, err := data.New(logger, cc)
	if err != nil {
		panic(err)
	}
	cv := []*dataUtils.CustomValidator{
		{
			Field:     "sku",
			Validator: data.ValidateSKU,
		},
	}
	validation := dataUtils.NewValidation(cv)
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

	bindAddress := config.ENV.Address
	server := &http.Server{
		Addr:         bindAddress,
		Handler:      corsHandler(serveMux),
		ErrorLog:     stdLogger,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting Http Server at %#v\n", bindAddress))
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
