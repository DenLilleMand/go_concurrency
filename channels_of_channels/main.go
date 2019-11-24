package main

import (
	"fmt"
	"time"
)

type Request struct {
	Name   string
	Result chan int
}

func handle(r *Request) {
	time.Sleep(1000) // process the request
	fmt.Printf("Processed request:%v\n", r.Name)
	r.Result <- 1
}

func main() {
	req := Request{"Initial request", make(chan int)}
	go handle(&req)
	fmt.Printf("Result: %v\n", <-req.Result)

}
