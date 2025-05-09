package main

import (
	"fmt"
	"sync"
)

var counter int

func add(wg *sync.WaitGroup) {
	for i := 0; i < 1000; i++ {
		counter++
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go add(&wg)
	}

	wg.Wait()
	fmt.Println("Expected counter = 10000")
	fmt.Println("Actual counter =", counter)
}