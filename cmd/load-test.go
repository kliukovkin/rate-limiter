package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {
	n := flag.Int("n", 1, "amount of workers")
	rps := flag.Int("rps", 10, "requests per second")
	c := makeWorkerPool(*n)
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		for i := 0; i < *rps; i++ {
			c <- 1
		}
	}
	close(c)
}

func worker(c chan int) {
	client := &http.Client{}
	for range c {
		req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
		req.Header.Set("UserID", "1")
		resp, _ := client.Do(req)
		fmt.Println(resp.StatusCode)
	}
}

func makeWorkerPool(n int) chan int {
	c := make(chan int)
	for i := 0; i < n; i++ {
		go worker(c)
	}
	return c
}
