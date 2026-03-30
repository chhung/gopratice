```go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode"
)

// ─────────────────────────────────────────────
// 1. 常數集合 & 哨兵錯誤 (Sentinel Errors)
//    errors.New 建立只比較位址的不可變錯誤值
//    var ErrXxx = errors.New(...) 宣告在套件層級
// ─────────────────────────────────────────────

const (
	defaultPort     = "8080"
	maxTextLength   = 280
	shutdownTimeout = 5 * time.Second
)

var (
	ErrEmptyText   = errors.New("text is required")
	ErrTextTooLong = errors.New("text exceeds 280 characters")
)

// ─────────────────────────────────────────────
// 2. 介面 (Interface)：只定義行為約定，不綁定實作
//    使用介面型別儲存依賴，讓程式碼易於測試和替換
//    介面本身不能當接收者，只有具體型別才能
// ─────────────────────────────────────────────

type Repository interface {
	List(ctx context.Context) ([]Message, error)
	Create(ctx context.Context, text string) (Message, error)
}

// ─────────────────────────────────────────────
// 3. 結構體 Struct + struct tag
//    json:"field_name"    → 控制 JSON 序列化的欄位名稱
//    json:"...,omitempty" → 零值時省略欄位（不輸出）
// ─────────────────────────────────────────────

type Message struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// ─────────────────────────────────────────────
// 4. 執行緒安全的記憶體儲存庫
//    sync.RWMutex：多個讀者可同時讀，寫者獨佔
//    RLock / RUnlock → 讀取保護
//    Lock  / Unlock  → 寫入保護
// ─────────────────────────────────────────────

type InMemoryRepo struct {
	mu       sync.RWMutex
	messages []Message
	nextID   int
}

// 指標接收者：需要讀取或修改 repo 的內部欄位
func (r *InMemoryRepo) List(_ context.Context) ([]Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 回傳副本，防止外部直接修改底層 slice
	result := make([]Message, len(r.messages))
	copy(result, r.messages)
	return result, nil
}

func (r *InMemoryRepo) Create(_ context.Context, text string) (Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	msg := Message{ID: r.nextID, Text: text, CreatedAt: time.Now()}
	r.messages = append(r.messages, msg)
	return msg, nil
}

// ─────────────────────────────────────────────
// 5. Service 層：驗證輸入 + 錯誤包裝
//    fmt.Errorf("context: %w", err) → %w 保留原始錯誤
//    讓呼叫端可用 errors.Is / errors.As 穿透包裝層比對
// ─────────────────────────────────────────────

type MessageService struct {
	repo Repository // 以介面儲存，解耦具體實作
}

func (s *MessageService) Create(ctx context.Context, text string) (Message, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return Message{}, ErrEmptyText
	}
	if len(text) > maxTextLength {
		return Message{}, ErrTextTooLong
	}
	msg, err := s.repo.Create(ctx, text)
	if err != nil {
		return Message{}, fmt.Errorf("create message: %w", err)
	}
	return msg, nil
}

func (s *MessageService) List(ctx context.Context) ([]Message, error) {
	msgs, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}
	return msgs, nil
}

// ─────────────────────────────────────────────
// 6. 錯誤處理：errors.Is + errors.Join
//    errors.Is → 穿透所有 %w 包裝層，比對哨兵錯誤
//    errors.Join → 合併多個錯誤（Go 1.20+），資源清理時常用
// ─────────────────────────────────────────────

func demonstrateErrors() {
	wrapped := fmt.Errorf("service: %w", ErrEmptyText)
	fmt.Println(errors.Is(wrapped, ErrEmptyText)) // true，穿透包裝

	closeErr1 := fmt.Errorf("disconnect: %w", errors.New("connection reset"))
	closeErr2 := fmt.Errorf("unsubscribe: %w", errors.New("already closed"))
	combined := errors.Join(closeErr1, closeErr2)
	fmt.Println(combined)
}

// ─────────────────────────────────────────────
// 7. Context：傳遞截止時間與取消信號
//    WithTimeout  → 超過期限自動取消，適合 DB 查詢、HTTP 呼叫
//    WithCancel   → 手動呼叫 cancel() 取消，適合提早結束
//    defer cancel() → 提早完成時立即釋放資源（務必呼叫）
// ─────────────────────────────────────────────

func demonstrateContext(parent context.Context) {
	ctx, cancel := context.WithTimeout(parent, 3*time.Second)
	defer cancel()

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("operation completed before timeout")
	case <-ctx.Done():
		fmt.Println("timed out:", ctx.Err())
	}
}

// ─────────────────────────────────────────────
// 8. Goroutine + sync.WaitGroup + 帶緩衝的 Channel
//    make(chan T, N)      → 緩衝值 N，發送方在緩衝滿前不阻塞
//    wg.Add(1)           → 啟動 goroutine 前呼叫，避免競態
//    defer wg.Done()     → goroutine 結束時自動通知
//    go func(arg){}(val) → 傳入參數避免閉包捕獲迴圈變數
// ─────────────────────────────────────────────

func runWorkers(ctx context.Context) {
	jobs := []string{"ingest", "parse", "publish"}
	errCh := make(chan error, len(jobs)) // 緩衝 = job 數，goroutine 不會因接收慢而卡住

	var wg sync.WaitGroup
	for _, job := range jobs {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case <-time.After(100 * time.Millisecond):
				fmt.Printf("worker [%s] done\n", name)
				errCh <- nil
			}
		}(job)
	}

	// 等所有 worker 結束後關閉 channel，讓 range 退出
	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			fmt.Println("worker error:", err)
		}
	}
}

// ─────────────────────────────────────────────
// 9. Select + ctx.Done()：多路 Channel 等待
//    select 選擇第一個就緒的 case，無優先順序
//    case msg, ok := <-ch → ok = false 表示 channel 已關閉
// ─────────────────────────────────────────────

func listenLoop(ctx context.Context, msgCh <-chan string) error {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("shutting down listener")
			return nil
		case msg, ok := <-msgCh:
			if !ok {
				return fmt.Errorf("message channel closed unexpectedly")
			}
			fmt.Println("received:", msg)
		}
	}
}

// ─────────────────────────────────────────────
// 10. Defer：LIFO 順序延遲執行，確保資源清理
//     defer f.Close() → 函式結束時自動關閉
//     defer func(){...}() → 可在 defer 內執行複雜邏輯
// ─────────────────────────────────────────────

func openAndProcess(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	defer func() {
		fmt.Printf("cleanup complete for: %s\n", f.Name())
	}()

	fmt.Println("processing:", f.Name())
	return nil
}

// ─────────────────────────────────────────────
// 11. Recover / Panic：捕捉 panic 防止服務崩潰
//     recover() 只能在 defer 內生效
//     回傳非 nil 值表示有 panic 被捕捉，否則為 nil
// ─────────────────────────────────────────────

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				fmt.Fprintf(os.Stderr, "panic recovered: %v\n", recovered)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// ─────────────────────────────────────────────
// 12. 型別斷言 (Type Assertion)
//     value, ok := iface.(ConcreteType)
//     ok == false 時斷言失敗，不會 panic
//     若省略 ok（單值形式）且失敗，則直接 panic
// ─────────────────────────────────────────────

func describeError(err error) string {
	if pathErr, ok := err.(*os.PathError); ok {
		return fmt.Sprintf("path error on %s: %v", pathErr.Path, pathErr.Err)
	}
	return err.Error()
}

// ─────────────────────────────────────────────
// 13. 閉包 Handler Factory：依賴注入的慣用法
//     函式回傳函式，閉包捕捉外部的 svc 依賴
//     json.NewDecoder + DisallowUnknownFields → 嚴格 JSON 解碼
//     errors.Is 在 handler 層判斷錯誤，回應適當的 HTTP 狀態碼
// ─────────────────────────────────────────────

func createMessageHandler(svc *MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Text string `json:"text"`
		}
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&input); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		msg, err := svc.Create(r.Context(), input.Text)
		if err != nil {
			if errors.Is(err, ErrEmptyText) || errors.Is(err, ErrTextTooLong) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(msg)
	}
}

// ─────────────────────────────────────────────
// 14. 匿名結構體 Slice：不定義只用一次的具名型別
//     常見於路由表、測試案例表、設定對照表
// ─────────────────────────────────────────────

func registerRoutes(mux *http.ServeMux, svc *MessageService) {
	routes := []struct {
		pattern string
		handler http.HandlerFunc
	}{
		{
			"GET /healthz",
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			},
		},
		{"POST /messages", createMessageHandler(svc)},
	}

	for _, route := range routes {
		mux.HandleFunc(route.pattern, route.handler)
	}
}

// ─────────────────────────────────────────────
// 15. strings.Map + unicode：字元級別字串轉換
//     strings.Map 遍歷每個 rune，回傳 -1 表示刪除該字元
//     unicode.IsSpace 涵蓋全形空白、Tab、換行等所有空白字元
// ─────────────────────────────────────────────

func normalizeInput(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, input)
}

// ─────────────────────────────────────────────
// 16. filepath.Join：跨平台安全路徑組合
//     自動處理路徑分隔符（Linux / vs Windows \）
//     避免手動字串拼接造成的雙斜線或錯誤分隔符
// ─────────────────────────────────────────────

func buildFilePath(dir, sub, file string) string {
	return filepath.Join(dir, sub, file)
}

// ─────────────────────────────────────────────
// 17. Table-driven 模式：以匿名結構體 Slice 定義案例
//     for _, tc := range cases → 逐一執行
//     實際測試放 *_test.go，使用 t.Run + t.Fatalf
// ─────────────────────────────────────────────

func runTableExample() {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "removes spaces", input: "12 34", want: "1234"},
		{name: "keeps digits",   input: "1234",  want: "1234"},
		{name: "empty string",   input: "",      want: ""},
	}

	for _, tc := range cases {
		got := normalizeInput(tc.input)
		if got != tc.want {
			fmt.Printf("[FAIL] %s: normalizeInput(%q) = %q, want %q\n", tc.name, tc.input, got, tc.want)
		} else {
			fmt.Printf("[PASS] %s\n", tc.name)
		}
	}
}

// ─────────────────────────────────────────────
// 18. Flag 解析（命令列參數）
//    flag.String("name", default, "usage") → 回傳 *string
//    flag.Parse()  → 解析 os.Args[1:]，在 main() 最前面呼叫
//    *flagVar      → 解引用取得實際值
//
// func main() {
//     configPath := flag.String("config", "config.json", "path to config file")
//     flag.Parse()
//     if *configPath == "" {
//         fmt.Fprintln(os.Stderr, "config flag is required")
//         os.Exit(1)
//     }
// }
// ─────────────────────────────────────────────

// ─────────────────────────────────────────────
// 19. signal.NotifyContext + 優雅關閉 (Graceful Shutdown)
//     signal.NotifyContext → OS 信號觸發 ctx 取消
//     os.Interrupt  → Ctrl+C
//     syscall.SIGTERM → Docker stop / kill 指令
//     server.Shutdown(ctx) → 等待進行中的請求完成後再關閉
// ─────────────────────────────────────────────

func main() {
	// 捕捉 Ctrl+C 與 SIGTERM，轉為 context 取消
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	repo := &InMemoryRepo{}
	svc := &MessageService{repo: repo}

	mux := http.NewServeMux()
	registerRoutes(mux, svc)

	server := &http.Server{
		Addr:         ":" + defaultPort,
		Handler:      recoverMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 在獨立 goroutine 啟動 HTTP server
	serverErr := make(chan error, 1)
	go func() {
		fmt.Printf("server listening on %s\n", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("http server: %w", err)
		}
	}()

	// 阻塞，等待信號或 server 錯誤
	select {
	case <-ctx.Done():
		fmt.Println("shutdown signal received")
	case err := <-serverErr:
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}

	// 優雅關閉：給 server 5 秒內完成進行中的請求
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "shutdown error: %v\n", err)
	}

	fmt.Println("server stopped cleanly")
}
```
