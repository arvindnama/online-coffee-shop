package data

import (
	"context"
	"fmt"
	"regexp"

	protos "github.com/arvindnama/golang-microservices/currency-service/protos"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductsDB struct {
	currencySvc   protos.CurrencyClient
	logger        hclog.Logger
	products      []*Product
	rates         map[string]float64
	subRateClient protos.Currency_SubscribeRatesClient
}

// Product defines the structure for an API products
// swagger:model Product
type Product struct {
	//this if of the product
	//
	// required: true
	// min:1
	ID int `json:"id"`

	// the name for this product
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this product
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"gt=0,required"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"required,sku"`
	// json: "-" indicated field will be omitted

}
type Products []*Product

var PRODUCTS_SEED_DATA = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc-xyz-lmn",
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and string coffee without milk",
		Price:       1.99,
		SKU:         "xyz-ijk-abc",
	},
}

func NewProductsDB(logger hclog.Logger, currencySvc protos.CurrencyClient) *ProductsDB {
	pDB := &ProductsDB{
		logger:        logger,
		currencySvc:   currencySvc,
		products:      PRODUCTS_SEED_DATA,
		rates:         make(map[string]float64),
		subRateClient: nil,
	}

	go pDB.handleRateUpdate()
	return pDB
}

var ErrPrdNotFound = fmt.Errorf("product not found")

func (pDB *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return pDB.products, nil
	}

	rate, err := pDB.getRate(currency)

	if err != nil {
		return nil, err
	}

	pl := Products{}

	for _, prod := range pDB.products {
		//[learning]: dereferencing and initializing into another var , will clone the struct
		np := *prod

		np.Price = np.Price * rate
		pl = append(pl, &np)
	}
	return pl, nil

}

func (pDB *ProductsDB) AddProduct(prod *Product) {
	prod.ID = pDB.getNextId()
	pDB.products = append(pDB.products, prod)
}

func (pDB *ProductsDB) UpdateProduct(id int, prod *Product) error {
	pos, err := pDB.findProduct(id)

	if err != nil {
		return err
	}

	prod.ID = id
	pDB.products[pos] = prod
	return nil
}

func (pDB *ProductsDB) GetProductById(id int, currency string) (*Product, error) {
	pos, err := pDB.findProduct(id)

	if err != nil {
		return nil, err
	}

	rate, err := pDB.getRate(currency)

	if err != nil {
		return nil, err
	}

	np := *pDB.products[pos]
	np.Price = np.Price * rate

	return &np, nil
}

func (pDB *ProductsDB) DeleteProduct(id int) error {
	pos, err := pDB.findProduct(id)

	if err != nil {
		return err
	}

	pDB.products = append(pDB.products[:pos], pDB.products[pos+1])
	return nil
}

func (pDB *ProductsDB) findProduct(id int) (int, error) {
	for i, p := range pDB.products {
		if p.ID == id {
			return i, nil
		}
	}
	return -1, ErrPrdNotFound
}

func (pDB *ProductsDB) getNextId() int {
	lp := pDB.products[len(pDB.products)-1]
	return lp.ID + 1
}

func (pDB *ProductsDB) getRate(currency string) (float64, error) {

	if cr, ok := pDB.rates[currency]; ok {
		return cr, nil
	}

	req := &protos.RateRequest{
		Base:        protos.Currencies_EUR,
		Destination: protos.Currencies(protos.Currencies_value[currency]),
	}
	resp, err := pDB.currencySvc.GetRate(context.Background(), req)

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

	pDB.logger.Debug(
		"gRPC currency client GetRate",
		"src", protos.Currencies_EUR,
		"dest", currency,
		"rate", resp.Rate,
	)

	// now subscribe to the rate change
	pDB.subRateClient.SendMsg(req)

	pDB.rates[currency] = resp.Rate

	return resp.Rate, nil
}

func (pDb *ProductsDB) handleRateUpdate() {
	sub, err := pDb.currencySvc.SubscribeRates(context.Background())

	pDb.subRateClient = sub

	if err != nil {
		pDb.logger.Error("Unable to subscribe to CSvc")
		return
	}

	for {
		srr, err := sub.Recv()

		if gerr := srr.GetError(); gerr != nil {
			serr := status.FromProto(gerr)

			if serr.Code() == codes.AlreadyExists {
				pDb.logger.Error(
					"Cannot subscribe to for rate updates more than once", "error", serr,
				)
			}
		}

		if streamingRateResp := srr.GetRateResponse(); streamingRateResp != nil {
			if err != nil {
				pDb.logger.Error("Error retrieving Rated from CSvc")
				return
			}

			pDb.logger.Debug(
				"New Rate received",
				"base", streamingRateResp.GetBase(),
				"dest", streamingRateResp.GetDestination(),
				"new rate", streamingRateResp.Rate,
			)
			pDb.rates[streamingRateResp.Destination.String()] = streamingRateResp.Rate
		}
	}
}

func ValidateSKU(fl validator.FieldLevel) bool {
	// sku format: xxxx-xxxx-xxxx
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}
