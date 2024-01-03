package api

import (
	"fmt"
	"net/http"
)

const (
	CPU     string = "cpu"
	MEMMORY string = "memory"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type RequestHandler struct{}

func NewRequestHandler() *RequestHandler {
	return new(RequestHandler)
}

func (rh *RequestHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Metrics")

	w.Write([]byte("hello"))
}
