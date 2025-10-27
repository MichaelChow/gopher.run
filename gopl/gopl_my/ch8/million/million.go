package main

import (
	"fmt"
	"sync"
)

func Million() {
	wg := sync.WaitGroup{}
	wg.Add(100000)
	for i := 0; i < 100000; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Print(i)
			// time.Sleep(1 * time.Second)
		}(i)
	}
	wg.Wait()
}

func r() {
	ch := make(chan int, 3)
	go func() {
		ch <- 1
		fmt.Print(1)
	}()
	go func() {
		ch <- 2
		fmt.Print(2)
	}()
	go func() {
		ch <- 3
		fmt.Print(3)
	}()
}

func main() {
	// Million()
	r()
}
