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

func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	o.logger.Debug("handling create order")

	order := r.Context().Value(middleware.RequestBody{}).(data.Order)
	order.Status = data.Initiated
	orderId, err := o.store.AddOrder(&order)
	writeError(w, err)

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")

	addedOrder, err := o.store.GetOrder(orderId)
	writeError(w, err)
	err = dataUtils.ToJSON(addedOrder, w)
	writeError(w, err)

}

func (o *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	o.logger.Debug("handling Get All orders")

	orders, err := o.store.GetAllOrders()
	writeError(w, err)
	w.WriteHeader(http.StatusOK)
	dataUtils.ToJSON(&orders, w)
}

func (o *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderId, err := strconv.Atoi(r.PathValue("id"))
	writeError(w, err)
	o.logger.Debug(fmt.Sprintf("handling Get order %#v", orderId))

	order, err := o.store.GetOrder(int64(orderId))
	writeError(w, err)

	dataUtils.ToJSON(order, w)
}

func (o *OrderHandler) PatchOrder(w http.ResponseWriter, r *http.Request) {
	orderId, err := strconv.Atoi(r.PathValue("id"))
	writeError(w, err)
	o.logger.Debug(fmt.Sprintf("handling PATCH order %#v", orderId))

	order := r.Context().Value(middleware.RequestBody{}).(data.Order)
	fmt.Println(order)

	err = o.store.UpdateOrderStatus(int64(orderId), order.Status)
	writeError(w, err)
	updatedOrder, err := o.store.GetOrder(int64(orderId))
	writeError(w, err)
	dataUtils.ToJSON(updatedOrder, w)
}

func writeError(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(http.StatusCreated)
		dataUtils.ToJSON(&data.ValidationError{
			Messages: []string{err.Error()},
		}, w)
	}
}
