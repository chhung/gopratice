package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	inputCh := make(chan string, 5)

	var wg sync.WaitGroup
	wg.Add(1)
	go worker(1, &wg, inputCh)

	go func() {
		fmt.Println("waitting for goroutine.")
		wg.Wait()
		close(inputCh)
		fmt.Println("goroutine done.")
	}()

	time.Sleep(time.Second * 6)
	for input := range inputCh {
		fmt.Printf("[%s] Received input: %s\n",
			time.Now().Format("2006-01-02 15:04:05"), input)
	}

	fmt.Println("Hello, World!")
}

func worker(id int, wg *sync.WaitGroup, input chan<- string) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second * 1)
		input <- fmt.Sprintf("[%s] Worker %d: Working...",
			time.Now().Format("2006-01-02 15:04:05"), id)
	}
}
