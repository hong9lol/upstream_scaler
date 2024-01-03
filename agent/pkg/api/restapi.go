package api

import (
	"net/http"
	"strconv"
)

var requestHandler *RequestHandler

func init() {
	requestHandler = NewRequestHandler()
}

func RunServer(addr string, port int) {
	http.Handle("/api/v1/metrics", http.HandlerFunc(requestHandler.GetMetrics))
	http.ListenAndServe(addr+":"+strconv.Itoa(port), nil)
}
