package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	database "github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/db"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type RequestHandler struct{}

func NewRequestHandler() *RequestHandler {
	return new(RequestHandler)
}

func (rh *RequestHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Metrics")
	deployment := strings.Split(r.URL.Path, "api/v1/metrics/")[1]
	db := database.GetInstance()
	stat, err := db.GetStat(deployment)
	if err != nil {
		fmt.Println("Can not find the any pod of " + deployment + " in this node")
	}
	res, err := json.Marshal(stat)
	if err != nil {
		return
	}
	w.Write(res)
}
