package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/kavehmz/prime"
)

func main() {
	fmt.Println("Appedge pod starts!")
	fmt.Println("My Ip is ", os.Getenv("APP_SERVICE_SERVICE_HOST"))
	fmt.Println("My Port is ", os.Getenv("APP_SERVICE_SERVICE_PORT"))
	http.HandleFunc("/prime", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received new request!")
		numStr := r.Header["Num"][0]
		reqNum, err := strconv.ParseUint(numStr, 10, 64)
		if err != nil {
			fmt.Println("failed to convert " + numStr)
			fmt.Fprintf(w, "failed to convert")
			return
		}
		p := prime.SieveOfEratosthenes(reqNum)
		fmt.Println("Prime number: ", len(p))
		fmt.Fprintf(w, "Number of primes: %d", len(p))
		//  := r.Header["ResponseMB"]
	})
	http.ListenAndServe(":8080", nil)
}
