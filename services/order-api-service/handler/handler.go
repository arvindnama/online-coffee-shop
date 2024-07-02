package handler

import (
	"fmt"
	"net/http"
	"strconv"

	dataUtils "github.com/arvindnama/golang-microservices/libs/utils/data-utils"
	"github.com/arvindnama/golang-microservices/order-service/data"
	"github.com/arvindnama/golang-microservices/order-service/middleware"
	"github.com/hashicorp/go-hclog"
)

type OrderHandler struct {
	logger hclog.Logger
	store  data.OrderDatabase
}

func NewOrderHandler(logger hclog.Logger, store data.OrderDatabase) *OrderHandler {
	return &OrderHandler{logger, store}
}

// swagger:route POST /orders orders createOrder
// Creates an order
// responses:
//
// 201: OrderResponse
// 401: ErrorResponse
// 500: ErrorResponse
func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	o.logger.Debug("handling create order")

	w.Header().Add("Content-Type", "application/json")

	order := r.Context().Value(middleware.RequestBody{}).(data.Order)
	order.Status = data.Initiated
	orderId, err := o.store.AddOrder(r.Context(), &order)

	if o.writeError(w, err); err != nil {
		return
	}

	addedOrder, err := o.store.GetOrder(orderId)
	if o.writeError(w, err); err != nil {
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = dataUtils.ToJSON(addedOrder, w)

	if o.writeError(w, err); err != nil {
		return
	}
}

// swagger:route GET /orders orders getAllOrders
// Gets all the registered orders from the database
// responses:
//
// 200: OrdersResponse
// 401: ErrorResponse
// 500: ErrorResponse
func (o *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	o.logger.Debug("handling Get All orders")
	w.Header().Add("Content-Type", "application/json")
	pageNo, err := strconv.Atoi(r.URL.Query().Get("page_no"))
	if err != nil {
		pageNo = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 10
	}

	o.logger.Debug("PageSize", pageSize, "pageno", pageNo)
	orders, hasMore, err := o.store.GetAllOrders(pageNo, pageSize)
	if o.writeError(w, err); err != nil {
		return
	}

	ordersPaginated := &data.OrderPaginated{
		Content:  orders,
		PageNo:   pageNo,
		PageSize: pageSize,
		HasMore:  hasMore,
	}
	dataUtils.ToJSON(&ordersPaginated, w)
}

// swagger:route GET /orders/{id} orders getOrder
// Gets the order with orderId from the database
// responses:
//
// 200: OrderResponse
// 401: ErrorResponse
// 500: ErrorResponse
func (o *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderId, err := strconv.Atoi(r.PathValue("id"))
	w.Header().Add("Content-Type", "application/json")
	if o.writeError(w, err); err != nil {
		return
	}
	o.logger.Debug(fmt.Sprintf("handling Get order %#v", orderId))

	order, err := o.store.GetOrder(int64(orderId))
	if o.writeError(w, err); err != nil {
		return
	}

	dataUtils.ToJSON(order, w)
}

// swagger:route PATCH /orders/{id} orders patchOrder
// Updates the status of the Order
// responses:
//
// 201: OrderResponse
// 401: ErrorResponse
// 500: ErrorResponse
func (o *OrderHandler) PatchOrder(w http.ResponseWriter, r *http.Request) {
	orderId, err := strconv.Atoi(r.PathValue("id"))
	w.Header().Add("Content-Type", "application/json")
	if o.writeError(w, err); err != nil {
		return
	}
	o.logger.Debug(fmt.Sprintf("handling PATCH order %#v", orderId))

	order := r.Context().Value(middleware.RequestBody{}).(data.Order)

	err = o.store.UpdateOrderStatus(r.Context(), int64(orderId), order.Status)
	if o.writeError(w, err); err != nil {
		return
	}
	updatedOrder, err := o.store.GetOrder(int64(orderId))
	if o.writeError(w, err); err != nil {
		return
	}

	dataUtils.ToJSON(updatedOrder, w)
}

func (o *OrderHandler) writeError(w http.ResponseWriter, err error) {
	if err != nil {
		o.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		dataUtils.ToJSON(&data.ValidationError{
			Messages: []string{err.Error()},
		}, w)
	}
}
