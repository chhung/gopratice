Golang 學習瓶頸：不知道內建函式 / 標準庫怎麼用  
這是很多新手常見的限速步驟，因為 Go 高度依賴標準庫，功能強大但需要熟悉。   

1. 最重要的內建函式（Built-in Functions）
這些不需要 import，可以直接使用： 

|  類型   | 內建函式  | 主要用途 | 建議練習 | 
|  ----  | ----  |  ----  |  ----  |
| Slice / Array  | make, len, cap, append, copy, delete | 建立與操作 slice、map | 每天練習 |
| 指標  | new | 建立指標變數 | 基礎練習 |
| 錯誤與流程 | panic, recover | 錯誤處理與恢復 | 了解即可 |
| Channel | close | 關閉 channel | 併發練習 |
| 複數 | complex, real, imag | 複數運算 | 很少用 |

2. 最常用標準庫包（Top 10，建議優先熟悉）

- `fmt` —— 格式化輸出、打印 (`Print`, `Printf`, `Println`)
- `strings` —— 字串處理 (`Contains`, `Split`, `Join`, `Trim`, `Replace`)
- `strconv` —— 字串與數字轉換 (`Atoi`, `Itoa`, `ParseInt`, `FormatBool`)
- `time` —— 時間處理 (`Now`, `Parse`, `Format`, `Since`, `Sleep`)
- `os` —— 作業系統相關（檔案、環境變數、命令列參數）
- `io` —— 資料串流處理 (`Copy`, `ReadAll`, `MultiReader`)
- `path/filepath` —— 跨平台路徑處理
- `encoding/json` —— JSON 序列化與反序列化
- `net/http` —— HTTP 客戶端與伺服器
- `log` —— 記錄日誌
- `context` —— 請求取消、超時控制（非常重要）

---
## Golang 學習瓶頸：Goroutine + Channel + Context + Signal + Graceful Shutdown

這幾個概念高度連結，是 Go 併發（Concurrency）的核心。把這五個東西當成「一套系統」來學，而不是分開背。

### 大圖（Big Picture）—— 它們如何一起運作？
- Goroutine：輕量執行緒（像小工人）
- Channel：工人之間傳遞訊息的管道
- Context：給工人「取消指令」的遙控器（最重要！）
- Signal：接收作業系統的關機訊號（Ctrl+C 或 SIGTERM）
- Graceful Shutdown：收到關機指令後，優雅地讓所有工人把手上工作做完再下班

正確流程：
收到 Signal → 取消 Context → 所有 Goroutine 看到 Context.Done() 就停工 → 用 WaitGroup 等所有人結束 → 程式安全結束。

### 各個關鍵概念 + 必背函式表

#### Goroutine（最簡單）
``` go
go func() { ... }()   // 啟動一個 goroutine
```
- 記住：永遠要確保每個 goroutine 會結束（否則記憶體洩漏）

#### Channel（傳訊管道）

| 函式 / 操作 | 用途 | 範例 |
|  ----  |  ----  |  ----  |
| make(chan T) | 建立無緩channel | ch := make(chan int) |
| make(chan T, 10) | 建立有緩衝 channel | - | 
| ch <- value | 發送 | - | 
| value := <-ch | 接收 | - | 
| close(ch) | 關閉（通知沒資料了）| 很重要！ |

#### Context（取消神器）

| 函式 | 用途 | 常見用法 |
| --- | --- | --- |
| context.Background() | 根 context | 起點 |
| context.WithCancel(ctx) | 可手動取消 | ctx, cancel := ... |
| context.WithTimeout(ctx, 5*time.Second) | 自動超時取消 | 最常用 |
| ctx.Done() | 收到取消訊號的 channel | select { case <-ctx.Done(): ... } |
| ctx.Err() | 檢查為什麼取消 | context.Canceled 或 context.DeadlineExceeded |

#### Signal（接收關機訊號）—— 2026 年推薦用法
```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()   // 一定要記得呼叫！
```

#### Graceful Shutdown 完整模板（直接抄去用）
```go
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

func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():  // 收到取消訊號
			fmt.Printf("Worker %d 收到關機指令，準備結束\n", id)
			return
		default:
			// 正常工作...
			fmt.Printf("Worker %d 工作中...\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	// 1. 建立可被 Signal 取消的 Context（最推薦寫法）
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	// 2. 啟動多個 goroutine
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(ctx, i, &wg)
	}

	fmt.Println("程式啟動中... 按 Ctrl+C 測試關機")

	// 3. 等待關機訊號（已經被 ctx 處理了）
	<-ctx.Done()
	fmt.Println("\n收到關機訊號，開始優雅關閉...")

	// 4. 給一點時間讓工人結束（可選 timeout）
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 5. 等待所有 goroutine 結束
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("所有 goroutine 已安全結束")
	case <-shutdownCtx.Done():
		fmt.Println("關機超時，強制結束")
	}
}
```

#### 常見錯誤 & 腦袋空白解法

| 錯誤情況 | 為什麼會發生 | 正確寫法提示 |
| ---- | --- | --- |
| Goroutine 永遠不結束沒檢查 | ctx.Done() | 一定用 select + ctx.Done() |
| Channel 沒關閉 | 記憶體洩漏 | 結束時 close(ch) |
| Context 沒傳遞 | 取消指令傳不到 | 所有函式第一個參數都要接 ctx |
| Signal 重複觸發 | 沒呼叫 stop() | defer stop() 一定要寫 |
| 腦袋空白不知道從哪開始 | 沒模板 | 直接複製上面完整範例改 |
