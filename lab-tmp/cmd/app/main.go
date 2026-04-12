package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(4)
	go worker(1, &wg)
	go worker(2, &wg)
	go worker(3, &wg)
	go worker(4, &wg)

	wg.Wait()
	fmt.Println("Hello, World!")
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d: Working...\n", id)
}
