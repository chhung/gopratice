package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// === 1. 建立可被 signal 控制的 context ===
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// === 3. 啟動 worker pool ===
	wg := sync.WaitGroup{}
	wg.Add(4)
	// === 4. 模擬 producer（丟任務）===
	go worker(ctx, 1, &wg)
	go worker(ctx, 2, &wg)
	go worker(ctx, 3, &wg)
	go worker(ctx, 4, &wg)

	<-ctx.Done()
	fmt.Println("\n[MAIN] Received shutdown signal")

	wg.Wait()
	fmt.Println("Hello, World!")

}

func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: Stopping\n", id)
			return
		default:
			fmt.Printf("Worker %d: Working...\n", id)
		}
	}
}
