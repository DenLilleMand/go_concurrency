package main

import (
	"fmt"
)

func gen_ints(ints chan int, val int) {
	for {
		ints <- val
	}
}

func main() {
	ints := make(chan int)
	ints2 := make(chan int)
	go gen_ints(ints, 1)
	go gen_ints(ints2, 2)

	for {
		select {
		case s := <-ints:
			fmt.Printf("s = %+v\n", s)
		case s1 := <-ints2:
			fmt.Printf("s1 = %+v\n", s1)
		}
	}
}
