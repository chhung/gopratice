package main

import "fmt"

/*
併發教學
goroutine
channel, 與goroutine的通訊管道, （unbuffered → buffered → range/close）
select, 多個channel的監聽器
sync.WaitGroup, 等待多個goroutine完成
worker pool 模式
*/
func main() {
	fmt.Println("Hello, World!")
}
