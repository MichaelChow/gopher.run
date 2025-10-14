package main

import (
	"fmt"
	"sync"
)

func print12(wg *sync.WaitGroup) {
	ch := make(chan int)
	go func() {
		defer wg.Done()
		fmt.Print(1)
		ch <- 0
		for i := 1; i < 10; i++ {
			<-ch
			fmt.Print(1)
			ch <- 0
		}
		<-ch

	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			<-ch
			fmt.Print(2)
			ch <- 0
		}
	}()

}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	print12(&wg)
	wg.Wait()
}
