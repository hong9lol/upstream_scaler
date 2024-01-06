package manager

import (
	"io"
	"net/http"
	"time"

	database "github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/db"
)

type HPAHandler struct {
	controllerServiceName string
}

func NewHPAHandler() *HPAHandler {
	return &HPAHandler{}
}

func (h *HPAHandler) Start(controllerServiceName string) {
	h.controllerServiceName = controllerServiceName
	db := database.GetInstance()

	for {
		resp, err := http.Get("http://" + controllerServiceName + ":5001/api/v1/hpa")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// 결과 출력
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		// update hpa
		db.UpdateHPA(data)

		time.Sleep(1000 * 10 * time.Millisecond) // 10s
	}
}
