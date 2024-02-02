package manager

import (
	"io"
	"net/http"
	"time"

	"github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/config"
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

	// notify to controller: new agent started
	resp, err := http.Get("http://" + controllerServiceName + ".upstream-system.svc.cluster.local:5001/api/v1/agent")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	for {
		resp, err := http.Get("http://" + controllerServiceName + ".upstream-system.svc.cluster.local:5001/api/v1/hpas")
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

		time.Sleep(1000 * config.HPA_UPDATE_INTERVAL_SENCOND * time.Millisecond) // 15s
	}
}
