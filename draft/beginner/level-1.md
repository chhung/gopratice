```go
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// 1. 常數與自定義型別 (Enum 模擬)
const (
	MaxRetries = 3
	AppName    = "GoDrill-V1"
)

type Status int

const (
	StatusIdle Status = iota // iota 會自動累加 0, 1, 2...
	StatusRunning
	StatusError
)

// 2. 結構體 (Struct) 與 介面 (Interface)
type Operatable interface {
	Execute(val float64) (float64, error)
}

type Calculator struct {
	ID      string
	Version int
	Status  Status
	LastRun time.Time
}

// 3. 指標接收者 (Method Receiver) - 會修改原始物件
func (c *Calculator) Reset() {
	c.Status = StatusIdle
	fmt.Printf("[System] %s has been reset.\n", c.ID)
}

// 4. 核心邏輯：多重回傳與型別操作
func (c *Calculator) Execute(a, b float64, op string) (float64, error) {
	c.LastRun = time.Now()
	
	switch op {
	case "+":
		return a + b, nil
	case "/":
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("invalid operator: %s", op)
	}
}

func main() {
	// 5. 變數宣告 (Short Declaration & Explicit Type)
	var baseCount int = 100
	factor := 1.5
	message := "Initializing " + AppName

	// 6. 複合型別：切片 (Slice) 與 映射 (Map)
	operators := []string{"+", "/", "*"} // Slice
	history := make(map[string]float64)  // Map

	// 7. 結構體實例化 (取指標)
	calc := &Calculator{
		ID:      "CORE-01",
		Version: 1,
		Status:  StatusRunning,
	}

	fmt.Println(message)

	// 8. 迴圈、判斷與錯誤處理
	for i, op := range operators {
		// 強制型別轉換：int -> float64
		res, err := calc.Execute(float64(baseCount), factor, op)

		if err != nil {
			log.Printf("[Warn] Task %d (%s) failed: %v", i, op, err)
			continue
		}

		history[op] = res
		fmt.Printf("Step %d: %d %s %.2f = %.2f\n", i, baseCount, op, factor, res)
	}

	// 9. 併發練習 (The Hardcore Part)
	fmt.Println("\n--- Starting Concurrent Tasks ---")
	
	var wg sync.WaitGroup
	dataChan := make(chan string, 5) // 帶緩衝的通道

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		// 啟動 Goroutine
		go func(workerID int) {
			defer wg.Done()
			time.Sleep(time.Millisecond * 500) // 模擬耗時
			dataChan <- fmt.Sprintf("Worker %d: Task Complete", workerID)
		}(i) // 傳入參數避免閉包捕獲問題
	}

	// 10. 監聽與通道關閉
	go func() {
		wg.Wait()
		close(dataChan)
	}()

	// 11. 從 Channel 讀取結果
	for msg := range dataChan {
		fmt.Println("[Async]", msg)
	}

	// 12. 延遲執行與檔案操作 (I/O)
	defer fmt.Println("\n[Exit] All processes finished cleanly.")

	f, err := os.Create("run_log.txt")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("Last calculation at: %v\n", calc.LastRun))
}
```