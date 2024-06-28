package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	dataUtils "github.com/arvindnama/golang-microservices/libs/utils/data-utils"
	"github.com/arvindnama/golang-microservices/order-service/config"
	"github.com/arvindnama/golang-microservices/order-service/handler"
	"github.com/arvindnama/golang-microservices/order-service/middleware"
	"github.com/hashicorp/go-hclog"
)

func main() {

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "Order Api Service",
		Level: hclog.LevelFromString(config.Env.LogLevel),
	})
	stdLogger := logger.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	v := dataUtils.NewValidation(nil)
	m := middleware.NewMiddleware(logger, v)
	h := handler.NewOrderHandler(logger)

	router := http.NewServeMux()

	loadRoutes(m, h, router)

	stack := middleware.CreateMiddlewareStack(
		m.Logging,
		m.AllowCors,
	)
	bindAddress := config.Env.Address
	server := &http.Server{
		Addr:     bindAddress,
		Handler:  stack(router),
		ErrorLog: stdLogger,
	}
	go func() {

		logger.Info(fmt.Sprintf("Starting Http Server at %v\n", bindAddress))
		err := server.ListenAndServe()
		if err != nil {
			logger.Error("Error Starting Http Server", err)
		}
	}()

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	signal.Notify(channel, os.Kill)

	sig := <-channel
	logger.Info("signal Received", sig)
	timeoutCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeoutCtx)

}
