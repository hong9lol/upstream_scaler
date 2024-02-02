package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/api"
	"github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/config"
	"github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/manager"
)

const controllerServiceName = "upstream-controller"

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	hpaHanlder := manager.NewHPAHandler()
	go hpaHanlder.Start(controllerServiceName)

	// wait for first hpa update
	time.Sleep(1000 * config.HPA_UPDATE_INTERVAL_SENCOND * time.Millisecond) // 15s

	statCollector := manager.NewStatCollector()
	go statCollector.Start(controllerServiceName)
	api.RunServer("0.0.0.0", 3001)
	wg.Wait() // Wait for all goroutines to finish

	fmt.Println("All goroutines have finished")
}

// import (
// 	"fmt"

// 	"github.com/google/cadvisor/container"
// 	"github.com/google/cadvisor/container/libcontainer"
// )

// var (
// 	// Metrics to be ignored.
// 	// Tcp metrics are ignored by default.
// 	ignoreMetrics = container.MetricSet{
// 		container.MemoryNumaMetrics:              struct{}{},
// 		container.NetworkTcpUsageMetrics:         struct{}{},
// 		container.NetworkUdpUsageMetrics:         struct{}{},
// 		container.NetworkAdvancedTcpUsageMetrics: struct{}{},
// 		container.ProcessSchedulerMetrics:        struct{}{},
// 		container.ProcessMetrics:                 struct{}{},
// 		container.HugetlbUsageMetrics:            struct{}{},
// 		container.ReferencedMemoryMetrics:        struct{}{},
// 		container.CPUTopologyMetrics:             struct{}{},
// 		container.ResctrlMetrics:                 struct{}{},
// 		container.CPUSetMetrics:                  struct{}{},
// 	}

// 	// Metrics to be enabled.  Used only if non-empty.
// 	enableMetrics = container.MetricSet{}
// )

// func main() {
// 	fmt.Println("Cadvisor test")
// 	var includedMetrics container.MetricSet
// 	includedMetrics = enableMetrics
// 	ret, _ := libcontainer.GetCgroupSubsystems(includedMetrics)
// 	fmt.Println(ret)
// }

// import (
// 	"k8s.io/kubernetes/pkg/kubelet/server"
// )

//func main() {
// server.ListenAndServePodResources()

//}
