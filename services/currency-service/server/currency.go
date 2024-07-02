package server

import (
	"context"
	"io"
	"time"

	currency "github.com/arvindnama/golang-microservices/libs/grpc-protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/arvindnama/golang-microservices/currency-service/data"

	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	logger hclog.Logger
	rates  *data.ExchangeRates
	subs   map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest
}

func NewCurrency(logger hclog.Logger, rates *data.ExchangeRates) *Currency {
	//[learning] make method is use to `make` similar to alloco / calloc
	// alternatively in case of structs and map you use type{} to create
	// subs := map[protos.Currency_SubscribeRatesServer]*[]protos.RateRequest{}
	subs := make(map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest)
	c := &Currency{logger, rates, subs}
	go c.handleRateUpdates()
	return c
}

func (c *Currency) GetRate(
	ctx context.Context,
	req *currency.RateRequest,
) (*currency.RateResponse, error) {

	c.logger.Info(
		"handle GetRate",
		"base", req.GetBase(),
		"destination", req.GetDestination(),
	)

	if req.Base == req.Destination {
		status := status.Newf(
			codes.InvalidArgument,
			"Base currency %s cannot be same as destination %s",
			req.GetBase().String(),
			req.GetDestination().String(),
		)
		status, wde := status.WithDetails(req)

		//[learning]: This rear code block, just being defensive
		if wde != nil {
			c.logger.Error(
				"handle GetRate",
				"base", req.GetBase(),
				"destination", req.GetDestination(),
				wde,
			)
			return nil, wde
		}

		return nil, status.Err()
	}

	cr, err := c.rates.GetRate(req.GetBase().String(), req.GetDestination().String())
	if err != nil {
		c.logger.Error(
			"handle GetRate",
			"base", req.GetBase(),
			"destination", req.GetDestination(),
			err,
		)
		return nil, err
	}
	return &currency.RateResponse{
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

				err = sub.Send(
					&currency.StreamingRateResponse{
						Message: &currency.StreamingRateResponse_RateResponse{
							RateResponse: &currency.RateResponse{
								Base:        request.Base,
								Destination: request.Destination,
								Rate:        rate,
							},
						},
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
	src currency.Currency_SubscribeRatesServer,
) error {
	for {
		rr, err := src.Recv()

		//[learning]: when client closes the stream, err with io.EOF
		if err == io.EOF {
			c.logger.Info("client closed the connection")
			break
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
			sub = []*currency.RateRequest{}
		}

		var validationStatus *status.Status
		// check for duplicate subscriptions
		for _, req := range sub {
			if rr.Base == req.Base && rr.Destination == req.Destination {
				// duplicate subscriptions
				validationStatus = status.Newf(
					codes.AlreadyExists,
					"Unable to subscribe as it already exists",
				)

				validationStatus, _ = validationStatus.WithDetails(rr)
				break
			}
		}

		if validationStatus != nil {
			src.Send(&currency.StreamingRateResponse{
				Message: &currency.StreamingRateResponse_Error{
					Error: validationStatus.Proto(),
				},
			})
			break
		}

		sub = append(sub, rr)
		c.subs[src] = sub
	}
	return nil
}
