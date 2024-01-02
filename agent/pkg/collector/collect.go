package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ResponseData struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func collect() {
	for {
		fmt.Println("Hello")
		resp, err := http.Get("http://localhost:8080/api/v2.0/stats/summary/")
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		// 결과 출력
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var responseData ResponseData
		if err := json.Unmarshal(body, &responseData); err != nil {
			panic(err)
		}

		fmt.Printf("UserID: %d\n", responseData.UserID)
		fmt.Printf("ID: %d\n", responseData.ID)
		fmt.Printf("Title: %s\n", responseData.Title)
		fmt.Printf("Completed: %v\n", responseData.Completed)
		time.Sleep(1000 * time.Millisecond) // 1s
	}
}
