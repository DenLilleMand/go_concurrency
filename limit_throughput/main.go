package main

import (
	"fmt"
	"strconv"
	"time"
)

var (
	MaxOutstanding = 4 // Maximum requests being processed at any one time
	sem            = make(chan int, MaxOutstanding)
)

type Request struct {
	Name string
}

func handle(r *Request) {
	sem <- 1         // Put something into the buffer while we wait
	time.Sleep(2000) // process the request
	fmt.Println(r.Name)
	<-sem
}

func serve(queue chan *Request) {
	for {
		req := <-queue
		go handle(req)
	}

}

func main() {
	requests := []Request{}
	for i := 0; i < 10; i++ {
		requests = append(requests, Request{strconv.Itoa(i)})
	}
	queue := make(chan *Request)
	go serve(queue)
	for _, r := range requests {
		req := r // It is a OK way to avoid the &r bug in a for-loop
		queue <- &req
	}
	// It is a server of some sort, so it will just run forever, would be nice
	// to listen for Unix signals i guess but thats not the point
	for {
	}

}
