package server

import (
	"context"
	"io"
	"time"

	protos "github.com/arvindnama/golang-microservices/currency-service/protos"

	"github.com/arvindnama/golang-microservices/currency-service/data"

	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	logger hclog.Logger
	rates  *data.ExchangeRates
	subs   map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

func NewCurrency(logger hclog.Logger, rates *data.ExchangeRates) *Currency {
	//[learning] make method is use to `make` similar to alloco / calloc
	// alternatively in case of structs and map you use type{} to create
	// subs := map[protos.Currency_SubscribeRatesServer]*[]protos.RateRequest{}
	subs := make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)
	c := &Currency{logger, rates, subs}
	go c.handleRateUpdates()
	return c
}

func (c *Currency) GetRate(
	ctx context.Context,
	req *protos.RateRequest,
) (*protos.RateResponse, error) {

	c.logger.Info(
		"handle GetRate",
		"base",
		req.GetBase(),
		"destination",
		req.GetDestination(),
	)
	cr, err := c.rates.GetRate(req.GetBase().String(), req.GetDestination().String())
	if err != nil {
		return nil, err
	}
	return &protos.RateResponse{
			Base:        req.Base,
			Destination: req.Destination,
			Rate:        cr,
		},
		nil

}

func (c *Currency) handleRateUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.logger.Debug("Got Updated Rate")
		for sub, requests := range c.subs {
			for _, request := range requests {
				rate, err := c.rates.GetRate(
					request.GetBase().String(),
					request.GetDestination().String(),
				)

				if err != nil {
					c.logger.Error(
						"Unable to get rate for", "base", request.GetBase(), "dest", request.GetDestination(),
					)
					continue
				}

				c.logger.Debug(
					"Sending an update to client", "base", request.GetBase(), "dest", request.GetDestination(), "new rate", rate,
				)

				err = sub.Send(&protos.RateResponse{
					Base:        request.Base,
					Destination: request.Destination,
					Rate:        rate,
				})

				if err != nil {
					c.logger.Error(
						"Unable to send msg to client rate for", "base", request.GetBase(), "dest", request.GetDestination(),
					)
				}
			}
		}
	}
}

func (c *Currency) SubscribeRates(
	src protos.Currency_SubscribeRatesServer,
) error {
	for {
		rr, err := src.Recv()

		//[learning]: when client closes the stream, err with io.EOF
		if err == io.EOF {
			c.logger.Info("client closed the connection")
			return err
		}

		if err != nil {
			c.logger.Error("Unable to read from client", "error", err)
			return err
		}

		c.logger.Info(
			"Handle client request", "request_base", rr.GetBase(), "request_dest", rr.GetDestination(),
		)

		sub, ok := c.subs[src]

		if !ok {
			// [learning]: make cannot be used in create slice ?? (get more info)
			sub = []*protos.RateRequest{}
		}
		sub = append(sub, rr)
		c.subs[src] = sub
	}
	return nil
}
