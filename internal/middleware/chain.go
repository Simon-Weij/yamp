package middleware

import (
	"net/http"
	"slices"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, mw := range slices.Backward(middlewares) {
		h = mw(h)
	}
	return h
}

type Chainer struct {
	middlewares []Middleware
}

func New(middlewares ...Middleware) Chainer {
	return Chainer{middlewares: middlewares}
}

func (c Chainer) Then(h http.Handler) http.Handler {
	for _, mw := range slices.Backward(c.middlewares) {
		h = mw(h)
	}
	return h
}
