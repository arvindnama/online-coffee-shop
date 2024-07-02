package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/arvindnama/golang-microservices/currency-service/config"
	"github.com/arvindnama/golang-microservices/currency-service/data"
	"github.com/arvindnama/golang-microservices/currency-service/server"
	currency "github.com/arvindnama/golang-microservices/libs/grpc-protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "Currency service",
		Level: hclog.LevelFromString(config.Env.LogLevel),
	})
	gs := grpc.NewServer()
	er, err := data.NewExchangeRates(logger)

	if err != nil {
		logger.Error("Unable to create exchange rate", err)
		os.Exit(1)
	}
	cs := server.NewCurrency(logger, er)

	currency.RegisterCurrencyServer(gs, cs)

	bindAddress := config.Env.Address
	lis, err := net.Listen("tcp", bindAddress)
	if err != nil {
		logger.Error("Unable to listen", err)
		os.Exit(1)
	}

	go func() {
		reflection.Register(gs)
		logger.Info(fmt.Sprintf("Starting server on port %#v", bindAddress))
		gs.Serve(lis)
	}()

	channel := make(chan os.Signal, 1)

	signal.Notify(channel, os.Interrupt)
	signal.Notify(channel, os.Kill)

	sig := <-channel

	logger.Info("Signal received", sig)
	lis.Close()
}
