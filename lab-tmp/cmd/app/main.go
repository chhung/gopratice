package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	test()

	//chan_test()
	fmt.Println("Hello, World!")
}

func test() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop() // 確保程式結束時清理資源

	fmt.Println("伺服器運行中，按 Ctrl+C 結束...")

	// 2. 阻塞在這裡，直到接收到 SIGINT 或 SIGTERM
	<-ctx.Done()

	// 3. 收到訊號後的清理工作
	fmt.Println("接收到關閉訊號，正在安全退出...")
}

func chan_test() {
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
}

func worker(id int, wg *sync.WaitGroup, input chan<- string) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second * 1)
		input <- fmt.Sprintf("[%s] Worker %d: Working...",
			time.Now().Format("2006-01-02 15:04:05"), id)
	}
}
