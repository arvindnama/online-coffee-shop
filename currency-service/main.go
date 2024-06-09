package main

import (
	"net"
	"os"
	"os/signal"

	"github.com/arvindnama/golang-microservices/currency-service/data"
	protos "github.com/arvindnama/golang-microservices/currency-service/protos"
	"github.com/arvindnama/golang-microservices/currency-service/server"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9091", "Bind address for the service")

func main() {
	env.Parse()

	logger := hclog.Default()
	gs := grpc.NewServer()
	er, err := data.NewExchangeRates(logger)

	if err != nil {
		logger.Error("Unable to create exchange rate", err)
		os.Exit(1)
	}
	cs := server.NewCurrency(logger, er)

	protos.RegisterCurrencyServer(gs, cs)

	lis, err := net.Listen("tcp", *bindAddress)
	if err != nil {
		logger.Error("Unable to listen", err)
		os.Exit(1)
	}

	go func() {
		reflection.Register(gs)
		logger.Info("Starting server on port ", *bindAddress)
		gs.Serve(lis)
	}()

	channel := make(chan os.Signal, 1)

	signal.Notify(channel, os.Interrupt)
	signal.Notify(channel, os.Kill)

	sig := <-channel

	logger.Info("Signal received", sig)
	lis.Close()
}
