package main

import (
	"currency-service/server"
	"net"
	"os"
	"os/signal"

	protos "currency-service/protos"

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
	cs := server.NewCurrency(logger)

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
