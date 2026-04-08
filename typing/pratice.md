# Go 初學者 14 天學習 SOP — 以 eventconsumer 專案為目標

> 目標：14 天內具備完整理解 eventconsumer 專案的能力
> 學習資源順序：[Go Tour](https://go.dev/tour/) → [Go by Example](https://gobyexample.com/) → 直接讀 eventconsumer 原始碼

---

## 第一階段：Go 語言基礎（Day 1–4）

### Day 1 — 環境 + 型別 + 控制流

| 項目 | 學習內容 | 對應專案位置 |
|---|---|---|
| 環境 | 安裝 Go、`go run`、`go build`、module 概念 | `go.mod` |
| 變數宣告 | `var`, `:=`, `const` 群組, `_` blank identifier | `cmd/app/main.go` L22-31 |
| 基本型別 | `string`, `int`, `int64`, `uint64`, `bool`, `[]byte` | 散見各處 |
| 控制流 | `if err != nil`, `for`, `for range`, `switch/case` | `cmd/app/config.go` L72-86 |

**驗收**：能寫出一個程式讀取 CLI 參數，用 `switch` 判斷輸入，用 `for` 迭代 slice 並印出結果

---

### Day 2 — Struct + Method + 指標

| 項目 | 學習內容 | 對應專案位置 |
|---|---|---|
| Struct 定義 | 欄位、初始化、巢狀 struct | `cmd/app/config.go` L15-32（`AppConfig` 包 `BrokerConfig`）|
| Struct Tag | `mapstructure:"host"`, `xml:"table,attr"`, `bson:"CP_CODE"` | `internal/parserxml/parser.go` L8-21, `internal/dao/issue.go` L13-23 |
| Method | 值接收器 vs 指標接收器 | `config.go` L35（值）vs `couponissue.go` L72（指標）|
| 指標 | `*T`, `&T`, nil check, dereference | `main.go` L105-114（`issue == nil` / `*issue`）|

**驗收**：定義一個 struct 帶 tag，寫值/指標接收器各一個 method，理解何時用哪種

---

### Day 3 — 函式進階 + 錯誤處理

| 項目 | 學習內容 | 對應專案位置 |
|---|---|---|
| 多回傳值 | `(value, error)` 慣例 | 所有函式皆遵循 |
| Error wrapping | `fmt.Errorf("...: %w", err)` | `main.go` L90 |
| Sentinel error | `var errX = errors.New(...)` + `errors.Is()` | `main.go` L30, L191 |
| Variadic func | `func Foo(fields ...Field)` | `internal/pipeline/logger.go` L155 |
| Closure | 匿名函式當 callback / goroutine body | `main.go` L234-237 |
| Defer | `defer f.Close()`, defer + 匿名函式 | `main.go` L91-95 |

**驗收**：寫一個函式鏈 A→B→C，每層用 `%w` 包裝錯誤，最外層用 `errors.Is()` 判斷 sentinel error

---

### Day 4 — Slice / Array / Map + 標準庫工具

| 項目 | 學習內容 | 對應專案位置 |
|---|---|---|
| Slice | `make([]T, len, cap)`, `append`, `for range` | `internal/pipeline/types.go` L22-27 |
| 固定大小 Array | `[10]chan RoutedMessage` 語意差異 | `internal/pipeline/system.go` L19 |
| Map | `make(map[K]V)`, lookup `v, ok := m[k]` | `system.go` L20, L117 |
| `[]byte` 深拷貝 | `append([]byte(nil), src...)` | `main.go` L277 |
| `strings.Builder` | 高效字串拼接 | `internal/pipeline/logger.go` L277-289 |
| `strconv` | `ParseInt`, `Sprintf` | `main.go` L289 |
| `flag` | CLI 參數解析 | `main.go` L347-353 |

**驗收**：寫一個程式用 map 統計文字檔詞頻，結果用 `strings.Builder` 組合輸出

---

## 第二階段：Interface + 測試（Day 5–6）

### Day 5 — Interface 與多型

| 項目 | 學習內容 | 對應專案位置 |
|---|---|---|
| Interface 定義 | 方法集合，隱式實作（不需 `implements`）| `internal/pipeline/types.go` L39-42（`Processor`）|
| 多實作 | 同一 interface 由不同 struct 實作 | `AsyncLogger` / `ZapLogger` 都實作 `Logger` |
| Interface 作為參數 | 依賴抽象而非具體型別 | `system.go` L28（接收 `Logger`）|
| DAO interface | 隔離資料層 | `internal/dao/issue.go` L29-31（`IssueDAO`）|
| Type switch | `switch v := x.(type)` | `internal/pipeline/logger.go` L296-312 |
| `any` 型別 | `interface{}` 的別名 | `Field.value any` |

**驗收**：定義一個 `Formatter` interface，寫 JSON / Text 兩種實作，main 裡透過 interface 呼叫

---

### Day 6 — 單元測試

| 項目 | 學習內容 | 對應專案位置 |
|---|---|---|
| `testing` 套件 | `func TestXxx(t *testing.T)` | `cmd/app/config_test.go` |
| `t.Fatalf` / `t.Errorf` | 斷言模式 `if got != want` | `internal/parserxml/parser_test.go` |
| `t.TempDir()` | 測試用暫存目錄（自動清理）| `config_test.go` L11 |
| 測試負面路徑 | 預期 error 的測試 | `internal/parserxml/router_test.go` L49-58 |
| `go test ./...` | 執行所有測試 | — |

**驗收**：為 Day 5 的 `Formatter` 寫正向 + 反向測試各一個，`go test` 全過

---

## 第三階段：並行模型（Day 7–10）⚠️ 核心重點

### Day 7 — Goroutine + Channel 基礎

| 項目 | 學習內容 | 對應專案位置 |
|---|---|---|
| `go func()` | 啟動 goroutine | `main.go` L78-85 |
| Buffered channel | `make(chan T, size)` | `system.go` L51 |
| Channel 方向 | `<-chan T`（唯讀）| `system.go` L131（`workCh <-chan RoutedMessage`）|
| `for range` channel | close 時自動結束迴圈 | `system.go` L102, `logger.go` L119 |
| `close(ch)` | 通知下游結束 | `system.go` L88-90 |

**驗收**：寫一個 producer → buffered channel → consumer 模型，producer `close` 後 consumer 自動結束

---

### Day 8 — select 多路複用

| 項目 | 學習內容 | 對應專案位置 |
|---|---|---|
| `select` 基本 | 同時等待多個 channel | `main.go` L264-271（讀 subscription + ctx cancel）|
| `select` + `default` | 非阻塞寫入 / drop | `logger.go` L246-249（log buffer 滿就 drop）|
| `select` + timer | 帶延遲的重連等待 | `main.go` L299-307（`waitForBrokerReconnect`）|
| Context cancel 偵測 | `case <-ctx.Done()` | `system.go` L81-82 |

**驗收**：寫一個程式用 `select` 同時監聽 data channel 和 timeout，模擬 5 秒無資料就印 warning

---

### Day 9 — sync 套件

| 項目 | 學習內容 | 對應專案位置 |
|---|---|---|
| `sync.WaitGroup` | `Add` / `Done` / `Wait` 等待所有 goroutine | `system.go` L23 |
| `sync.Once` | 確保 `close(ch)` 只執行一次 | `system.go` L24, L88-90 |
| `atomic.Int64` | 無鎖計數器 | `logger.go` L53（drop 計數）|

**驗收**：啟動 5 個 worker goroutine，用 `WaitGroup` 等待全部完成，用 `atomic` 累加計數，主程式印總和

---

### Day 10 — context.Context

| 項目 | 學習內容 | 對應專案位置 |
|---|---|---|
| `context.Background()` | 根 context | `main.go` L73 |
| `context.WithCancel` | 手動取消 | `main.go` L73-74 |
| `context.WithTimeout` | 帶超時 | `internal/mongodb/client.go` L54 |
| Context 傳播 | 所有函式第一個參數帶 `ctx` | 貫穿整個專案 |
| `ctx.Done()` / `ctx.Err()` | 偵測取消 | `main.go` L131, `system.go` L81 |

**驗收**：寫一個程式 `WithCancel` → 傳 ctx 給 goroutine → 主程式 3 秒後 `cancel()` → goroutine 偵測 `ctx.Done()` 退出

---

## 第四階段：專案特有知識（Day 11–13）

### Day 11 — XML 解析（兩種模式）

| 項目 | 學習內容 | 對應專案位置 |
|---|---|---|
| `xml.Unmarshal` | 一次性反序列化（struct tag 映射）| `internal/parserxml/parser.go` L33-39 |
| `xml.NewDecoder` + Token 串流 | 逐 token 解析，提前 return 省效能 | `internal/parserxml/router.go` L27-80 |
| `io.EOF` 判斷 | token 讀完的結束條件 | `router.go` L38-40 |

**驗收**：拿 `test/mq-XML20260302_cleaned.txt` 的樣本資料，分別用兩種方式解析並印出 table 名稱和欄位

---

### Day 12 — OS Signal + Graceful Shutdown 全流程

| 項目 | 學習內容 | 對應專案位置 |
|---|---|---|
| `os/signal` | `signal.Notify`, `signal.Stop` | `main.go` L75-77 |
| Shutdown 順序 | signal → `cancel()` → `CloseInput()` → `Wait()` → `logger.Close()` | `main.go` L78-85, L125-128 |
| Defer 順序 | LIFO — 後宣告先執行 | `defer cancel()` 在 `defer logger.Close()` 之後宣告，所以 cancel 先跑 |

**驗收**：畫出 `main.go` 中 6 個 `defer` 的實際執行順序，並解釋為什麼這個順序能確保資源正確釋放

---

### Day 13 — 專案架構串讀

把所有知識串起來，走讀 `run()` 函式全流程：

```
┌─ parseRunOptions()           ← flag 解析
│
├─ loadAppConfig / loadMongoConfig  ← viper + struct tag 反序列化
│
├─ NewLogger()                 ← interface 工廠, goroutine 非同步寫入
│
├─ context.WithCancel          ← 取消傳播根
│   └─ go signal handler       ← goroutine + select
│
├─ mongodb.New()               ← context.WithTimeout ping
│
├─ NewSystem()                 ← map[string]Processor + [10]Processor
│   ├─ Start()                 ← goroutine × N（runRouter + runWorker）
│   │   ├─ runRouter           ← for range rawCh → ParseRoute → 分派到 tableChans / issueChans
│   │   └─ runWorker × 15     ← for range workCh → Parse → Process
│   └─ channel close 鏈        ← CloseInput → rawCh close → worker chan close → goroutine 結束
│
├─ for ctx.Err() == nil        ← 重連迴圈
│   ├─ connectBrokerSubscription  ← net.Dialer + cleanup closure
│   └─ consumeSubscription     ← select { ctx.Done / subscription.C }
│       └─ Submit → rawCh      ← select 非阻塞寫入
│
└─ defer 鏈倒序清理
```

**驗收**：能對著程式碼口述上面這張圖的每一步，解釋資料從 broker 進來到 worker 處理完的完整路徑

---

## 第五階段：驗收整合（Day 14）

### Day 14 — 實作一個迷你版 pipeline

自己從零寫一個簡化版，驗證所有知識是否串通：

1. **main** — `flag` 解析參數、`context.WithCancel`、`signal.Notify` 處理 Ctrl+C
2. **producer** — goroutine 模擬每 100ms 送一筆 `RawMessage` 到 buffered channel
3. **router** — goroutine 用 `for range` 讀 channel，用 `xml.NewDecoder` token 解析取 table 名，分派到對應 worker channel
4. **worker** — 定義 `Processor` interface，寫 2 個實作，各自用 goroutine 消費 worker channel
5. **shutdown** — `close` raw channel → close worker channels → `WaitGroup.Wait()` → 印統計

**最終驗收 checklist**：

- [ ] `go build` 編譯成功
- [ ] `go test ./...` 全過
- [ ] Ctrl+C 後 graceful shutdown（不丟資料、不 panic）
- [ ] 能畫出自己程式的 goroutine 數量與 channel 連接關係圖

---

## 時間分配建議

| 天數 | 階段 | 每日投入 | 重點 |
|---|---|---|---|
| Day 1–4 | 基礎語法 | 3–4 hr | **動手寫 > 讀文件**，每天一個驗收練習 |
| Day 5–6 | Interface + 測試 | 3–4 hr | 理解「隱式實作」是 Go 的核心設計 |
| Day 7–10 | 並行模型 | 4–5 hr | **最關鍵**，多畫 goroutine/channel 流向圖 |
| Day 11–13 | 專案走讀 | 3–4 hr | 邊讀邊對照前面學的語法 |
| Day 14 | 整合實作 | 5–6 hr | 從零寫迷你版，確認融會貫通 |

---

## 推薦學習資源（依順序）

1. [Go Tour](https://go.dev/tour/) — 互動式基礎教學
2. [Go by Example](https://gobyexample.com/) — 每個語法點一個可執行範例
3. [Effective Go](https://go.dev/doc/effective_go) — 官方慣用寫法指南
4. eventconsumer 原始碼 — 學完前兩個資源後直接讀專案