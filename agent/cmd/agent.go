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
	controllerServiceName := "127.0.0.1"
	hpaHanlder := manager.NewHPAHandler()
	go hpaHanlder.Start(controllerServiceName)

	statCollector := manager.NewStatCollector()
	go statCollector.Start(controllerServiceName)
	api.RunServer("localhost", 3000)
	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("All goroutines have finished")
}
