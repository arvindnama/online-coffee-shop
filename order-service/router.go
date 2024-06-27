package main

import (
	"net/http"

	"github.com/arvindnama/golang-microservices/order-service/handler"
	"github.com/arvindnama/golang-microservices/order-service/middleware"
)

func loadRoutes(
	m *middleware.Middleware,
	h *handler.OrderHandler,
	router *http.ServeMux,
) {

	requestBodyMiddleware := middleware.CreateMiddlewareStack(
		m.DeserializeRequestBody,
		m.ValidateRequestBody,
	)
	router.HandleFunc("GET /orders/{id}", h.GetOrder)
	router.HandleFunc("GET /orders", h.GetAllOrders)
	router.Handle("POST /orders", requestBodyMiddleware(http.HandlerFunc(h.CreateOrder)))
	router.Handle("PATCH /orders/{id}", requestBodyMiddleware(http.HandlerFunc(h.PatchOrder)))
}
