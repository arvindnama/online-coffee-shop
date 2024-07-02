package routes

import (
	"net/http"

	"github.com/arvindnama/golang-microservices/order-service/handler"
	"github.com/arvindnama/golang-microservices/order-service/middleware"
	openApiMiddleware "github.com/go-openapi/runtime/middleware"
)

func LoadRoutes(
	m *middleware.Middleware,
	h *handler.OrderHandler,
	router *http.ServeMux,
) {

	requestBodyMiddleware := middleware.CreateMiddlewareStack(
		m.IsAuthenticated,
		m.DeserializeRequestBody,
		m.ValidateRequestBody,
	)

	requestBodyNoValidationMiddleware := middleware.CreateMiddlewareStack(
		m.IsAuthenticated,
		m.DeserializeRequestBody,
	)

	router.HandleFunc("GET /orders/{id}", h.GetOrder)
	router.HandleFunc("GET /orders", h.GetAllOrders)
	router.Handle("POST /orders", requestBodyMiddleware(http.HandlerFunc(h.CreateOrder)))
	router.Handle("PATCH /orders/{id}", requestBodyNoValidationMiddleware(http.HandlerFunc(h.PatchOrder)))

	// documentation

	ops := openApiMiddleware.RedocOpts{SpecURL: "/swagger.yaml"}

	redocGetDocHandler := openApiMiddleware.Redoc(ops, nil)
	router.Handle("GET /docs", redocGetDocHandler)
	router.Handle("GET /swagger.yaml", http.FileServer(http.Dir("./")))
}
