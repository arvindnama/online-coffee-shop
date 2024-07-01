package data

import (
	"context"
	"fmt"
	"regexp"

	protos "github.com/arvindnama/golang-microservices/currency-service/protos"
	dbUtils "github.com/arvindnama/golang-microservices/libs/utils/db-utils"
	"github.com/arvindnama/golang-microservices/product-api-service/config"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ProductsStore struct {
	currencySvc   protos.CurrencyClient
	logger        hclog.Logger
	rates         map[string]float64
	subRateClient protos.Currency_SubscribeRatesClient
	dbConn        *gorm.DB
}

func New(logger hclog.Logger, currencySvc protos.CurrencyClient) (*ProductsStore, error) {
	dbConn, err := dbUtils.NewGormDbConnection(&config.ENV.DBConfig, logger)

	if err != nil {
		return nil, err
	}

	store := &ProductsStore{
		logger:        logger,
		currencySvc:   currencySvc,
		rates:         make(map[string]float64),
		subRateClient: nil,
		dbConn:        dbConn,
	}

	go store.handleRateUpdate()
	return store, nil
}

func (store *ProductsStore) GetProducts(currency string) (Products, error) {

	var products Products

	store.dbConn.Find(&products)
	err := store.updateProductRate(currency, products)

	return products, err

}

func (store *ProductsStore) updateProductRate(currency string, products Products) error {

	var rate float64 = 1.00
	var err error
	if currency != "" {
		rate, err = store.getRate(currency)

		store.logger.Debug("Currency", currency, "rate", rate)
		if err != nil {
			return err
		}
	}

	for _, prod := range products {
		prod.Price = prod.Price * rate
	}

	return nil

}

func (store *ProductsStore) AddProduct(prod *Product) {
	store.dbConn.Create(prod)
}

func (store *ProductsStore) UpdateProduct(id int, prod *Product) error {
	_, err := store.GetProductById(id, "")

	if err != nil {
		return err
	}

	prod.ID = id
	store.dbConn.Save(prod)
	return nil
}

func (store *ProductsStore) GetProductById(id int, currency string) (*Product, error) {

	var products Products

	store.dbConn.Where("id=?", id).First(&products)

	if len(products) == 0 {
		return nil, ErrPrdNotFound
	}

	err := store.updateProductRate(currency, products)

	return products[0], err
}

func (store *ProductsStore) DeleteProduct(id int) error {

	store.dbConn.Delete(&Product{}, id)
	return nil
}

func (store *ProductsStore) getRate(currency string) (float64, error) {

	if cr, ok := store.rates[currency]; ok {
		return cr, nil
	}

	req := &protos.RateRequest{
		Base:        protos.Currencies_EUR,
		Destination: protos.Currencies(protos.Currencies_value[currency]),
	}
	resp, err := store.currencySvc.GetRate(context.Background(), req)

	if err != nil {
		if s, ok := status.FromError(err); ok {
			md := s.Details()[0].(*protos.RateRequest)
			if s.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf(
					"unable to get rate from currency server base: %s & dest: %s cannot be same",
					md.Base.String(),
					md.Destination.String(),
				)
			}
			return -1, fmt.Errorf(
				"unable to get rate from currency server base: %s , dest: %s",
				md.Base.String(),
				md.Destination.String(),
			)
		}
		return 0, err
	}

	store.logger.Debug(
		"gRPC currency client GetRate",
		"src", protos.Currencies_EUR,
		"dest", currency,
		"rate", resp.Rate,
	)

	// now subscribe to the rate change
	store.subRateClient.SendMsg(req)

	store.rates[currency] = resp.Rate

	return resp.Rate, nil
}

func (store *ProductsStore) handleRateUpdate() {
	sub, err := store.currencySvc.SubscribeRates(context.Background())

	store.subRateClient = sub

	if err != nil {
		store.logger.Error("Unable to subscribe to CSvc")
		return
	}

	for {
		srr, err := sub.Recv()

		if gerr := srr.GetError(); gerr != nil {
			serr := status.FromProto(gerr)

			if serr.Code() == codes.AlreadyExists {
				store.logger.Error(
					"Cannot subscribe to for rate updates more than once", "error", serr,
				)
			}
		}

		if streamingRateResp := srr.GetRateResponse(); streamingRateResp != nil {
			if err != nil {
				store.logger.Error("Error retrieving Rated from CSvc")
				return
			}

			store.logger.Debug(
				"New Rate received",
				"base", streamingRateResp.GetBase(),
				"dest", streamingRateResp.GetDestination(),
				"new rate", streamingRateResp.Rate,
			)
			store.rates[streamingRateResp.Destination.String()] = streamingRateResp.Rate
		}
	}
}

func ValidateSKU(fl validator.FieldLevel) bool {
	// sku format: xxxx-xxxx-xxxx
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}
