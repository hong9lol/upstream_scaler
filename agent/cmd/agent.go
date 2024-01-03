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

	hpaHanlder := manager.NewHPAHandler()
	go hpaHanlder.Start()

	statCollector := manager.NewStatCollector()
	go statCollector.Start()
	api.RunServer("localhost", 3000)
	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("All goroutines have finished")
}
