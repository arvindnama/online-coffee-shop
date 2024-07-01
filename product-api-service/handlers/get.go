package handlers

import (
	"net/http"
	"strconv"

	dataUtils "github.com/arvindnama/golang-microservices/libs/utils/data-utils"
	"github.com/arvindnama/golang-microservices/product-api-service/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//
//	200: ProductsResponse
//	501: ErrorResponse
func (p *Products) GetAllProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Info("Handle GET products")

	currency := r.URL.Query().Get("currency")
	pageNoQ := r.URL.Query().Get("page_no")
	pageSizeQ := r.URL.Query().Get("page_size")

	pageNo := 1
	pageSize := 10

	if pageNoQ != "" {
		if val, err := strconv.Atoi(pageNoQ); err == nil {
			pageNo = val
		}
	}

	if pageSizeQ != "" {
		if val, err := strconv.Atoi(pageSizeQ); err == nil {
			pageSize = val
		}
	}

	p.l.Debug("currency", currency)
	p.l.Debug("page_no", pageNo)
	p.l.Debug("page_size", pageSize)
	products, hasMore, err := p.pDB.GetProducts(currency, pageNo, pageSize)

	productsPaginatedResponse := &data.ProductsPaginatedResponse{
		Content:  products,
		PageNo:   pageNo,
		PageSize: pageSize,
		HasMore:  hasMore,
	}
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		dataUtils.ToJSON(&GenericError{Message: err.Error()}, rw)
	}

	rw.Header().Add("Content-Type", "application/json")
	err = dataUtils.ToJSON(productsPaginatedResponse, rw)

	if err != nil {
		http.Error(rw, "Unable to marshal produces", http.StatusInternalServerError)
	}

}

// swagger:route GET /products/{id} products listProduct
// Returns a product
// responses:
//
//	200: ProductResponse
//	404: ErrorResponse
//	501: ErrorResponse
func (p *Products) GetProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Info("Handle GET products")
	rw.Header().Add("Content-Type", "application/json")

	currency := r.URL.Query().Get("currency")
	id := getProductID(r)
	product, err := p.pDB.GetProductById(id, currency)
	if err != nil {
		p.l.Error("[ERROR] fetching product", err)

		switch err {
		case data.ErrPrdNotFound:
			rw.WriteHeader(http.StatusNotFound)
			dataUtils.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		default:
			rw.WriteHeader(http.StatusInternalServerError)
			dataUtils.ToJSON(&GenericError{Message: err.Error()}, rw)
		}

	}

	err = dataUtils.ToJSON(product, rw)

	if err != nil {
		p.l.Error("[ERROR] serializing error")
		dataUtils.ToJSON(&GenericError{Message: err.Error()}, rw)
	}

}
