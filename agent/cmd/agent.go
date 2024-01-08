package main

import (
	"fmt"
	"sync"

	"github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/api"
	"github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/manager"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	controllerServiceName := "upstream-controller"
	hpaHanlder := manager.NewHPAHandler()
	go hpaHanlder.Start(controllerServiceName)

	statCollector := manager.NewStatCollector()
	go statCollector.Start(controllerServiceName)
	api.RunServer("0.0.0.0", 3001)
	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("All goroutines have finished")
}
