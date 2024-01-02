package main

import (
	"fmt"
	"sync"

	"github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/manager"
)

func main() {

	// init stat object to store resource metrics
	// Init()

	// restfull api to provide current stat
	// time.Sleep(10000 * time.Millisecond) // 10s

	// _db := db.NewDB()
	// _db.Tlogic()

	// resp, err := http.Get("http://127.0.0.1:5001/api/v1/hpa")
	// if err != nil {
	// 	panic(err)
	// }

	// defer resp.Body.Close()

	// // 결과 출력
	// data, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// type MetricsObject struct {
	// 	Name              string `json:"ame,omitempty"`
	// 	TargetUtilization int    `json:"target_utilization,omitempty"`
	// 	Type              string `json:"type,omitempty"`
	// }
	// type HPAObject struct {
	// 	Name        string          `json:"name,omitempty"`
	// 	Namespace   string          `json:"namespace,omitempty"`
	// 	MinReplicas int             `json:"min_replica,omitempty"`
	// 	MaxReplicas int             `json:"max_replica,omitempty"`
	// 	Target      string          `json:"target,omitempty"`
	// 	Metrics     []MetricsObject `json:"metrics,omitempty"`
	// }
	// a := []HPAObject{}
	// json.Unmarshal(data, &a)
	// fmt.Println(a[0].Name)
	// for {
	// 	time.Sleep(1000 * time.Millisecond) // 1s
	// 	// ret, _, _ := ExecuteRemoteCommand("media-mongodb-6575f756d7-t5ngn", "default", "head -1 /sys/fs/cgroup/cpu.stat")

	// 	// fmt.Println(ret)
	// 	// fmt.Println("Hello")
	// }
	var wg sync.WaitGroup
	wg.Add(2)
	hpaHanlder := manager.NewHPAHandler()
	go hpaHanlder.Start()

	statCollector := manager.NewStatCollector()
	go statCollector.Start()

	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("All goroutines have finished")
}
