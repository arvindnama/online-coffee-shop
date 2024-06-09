package server

import (
	"context"

	protos "github.com/arvindnama/golang-microservices/currency-service/protos"

	"github.com/arvindnama/golang-microservices/currency-service/data"

	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	logger hclog.Logger
	er     *data.ExchangeRates
}

func NewCurrency(logger hclog.Logger, er *data.ExchangeRates) *Currency {
	return &Currency{logger, er}
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
	cr, err := c.er.GetRate(req.GetBase().String(), req.GetDestination().String())
	if err != nil {
		return nil, err
	}
	return &protos.RateResponse{Rate: cr}, nil

}
