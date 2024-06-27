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

func NewOrderHandler(logger hclog.Logger) *OrderHandler {
	store := data.NewOrderStore()
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

	order := r.Context().Value(middleware.RequestBody{}).(data.Order)
	order.Status = data.Initiated
	orderId, err := o.store.AddOrder(&order)

	if o.writeError(w, err); err != nil {
		return
	}

	addedOrder, err := o.store.GetOrder(orderId)
	if o.writeError(w, err); err != nil {
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
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

	orders, err := o.store.GetAllOrders()
	if o.writeError(w, err); err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	dataUtils.ToJSON(&orders, w)
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
	if o.writeError(w, err); err != nil {
		return
	}
	o.logger.Debug(fmt.Sprintf("handling Get order %#v", orderId))

	order, err := o.store.GetOrder(int64(orderId))
	if o.writeError(w, err); err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
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
	o.writeError(w, err)
	o.logger.Debug(fmt.Sprintf("handling PATCH order %#v", orderId))

	order := r.Context().Value(middleware.RequestBody{}).(data.Order)
	fmt.Println(order)

	err = o.store.UpdateOrderStatus(int64(orderId), order.Status)
	if o.writeError(w, err); err != nil {
		return
	}
	updatedOrder, err := o.store.GetOrder(int64(orderId))
	if o.writeError(w, err); err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
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
