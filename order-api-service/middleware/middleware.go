package middleware

import (
	"net/http"

	dataUtils "github.com/arvindnama/golang-microservices/libs/utils/data-utils"
	"github.com/hashicorp/go-hclog"
)

type Middleware struct {
	logger     hclog.Logger
	validation *dataUtils.Validation
}

type MiddlewareHandler func(next http.Handler) http.Handler

func NewMiddleware(logger hclog.Logger, validation *dataUtils.Validation) *Middleware {
	return &Middleware{
		logger,
		validation,
	}
}

func CreateMiddlewareStack(handlers ...MiddlewareHandler) MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		for i := len(handlers) - 1; i >= 0; i-- {
			next = handlers[i](next)
		}
		return next
	}
}
