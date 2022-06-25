package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://192.168.0.194:30007/prime", nil)
	if err != nil {
		fmt.Println("error creating request", err)
		return
	}
	req.Header.Add("Num", "50000000")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error sending http request", err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	fmt.Println("response:", string(body))
}
