package server

import (
	"context"

	protos "github.com/arvindnama/golang-microservices/currency-service/protos"

	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	logger hclog.Logger
}

func NewCurrency(logger hclog.Logger) *Currency {
	return &Currency{logger}
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

	return &protos.RateResponse{Rate: 100.5}, nil

}
