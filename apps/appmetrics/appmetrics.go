package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	testCycle := 50
	succeedCycle := 0
	var totalTime int64
	totalTime = 0
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "http://172.18.0.2:30007/prime", nil)
	if err != nil {
		fmt.Println("error creating request", err)
		return
	}

	req.Header.Add("Num", "50000000")
	for i := 0; i < testCycle; i++ {
		start := time.Now()
		resp, err := client.Do(req)
		timeElapsed := time.Since(start).Microseconds()
		if err != nil {
			fmt.Println("error sending http request", err)
			continue
		}
		if resp.StatusCode != 200 {
			fmt.Println("error sending http request", resp.StatusCode)
			continue
		}
		succeedCycle++
		totalTime += timeElapsed
	}
	successRate := float32(succeedCycle) / float32(testCycle)
	fmt.Printf("Success Rate = %f %%\n", successRate*100)
	if succeedCycle != 0 {

		averageLatency := totalTime / int64(succeedCycle)
		fmt.Printf("Average Latency = %d ms\n", averageLatency/1000)
	}

}
