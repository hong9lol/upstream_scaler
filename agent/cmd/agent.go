package main

import (
	"fmt"
	"time"
)

// type stat map[string]map[string]int

// var objectStat stat

// func Init() stat {
// 	if objectStat == nil {
// 		objectStat = make(stat)
// 	}

// 	return objectStat
// }

func main() {

	// init stat object to store resource metrics
	// Init()

	// restfull api to provide current stat
	// time.Sleep(10000 * time.Millisecond) // 10s
	for {
		time.Sleep(1000 * time.Millisecond) // 1s
		fmt.Println("Hello")
	}
}
