package main

import (
	"fmt"
	"sort"
)

func main() {
	l := []int{1, 2, 3, 4}
	c := make(chan int)
	go func() {
		sort.Ints(l)
		c <- 1
	}()
	<-c
	fmt.Printf("%v\n", l)
}
