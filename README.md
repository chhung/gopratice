# gopratice

This repository contains standalone Go practice labs.

## Labs

- `lab1`: a minimal application structure for simple operations
- `lab2`: a minimal HTTP service structure with message endpoints
- `lab3`: input normalization that removes all whitespace from a line of text
- `lab4`: 檔案讀寫，結構化寫入及讀取，測試大檔（超過4G）每行讀取速度
- `lab5`: 並發程式寫法
- `lab6`: 讓prometheus來收資料，做為監控使用
- `lab7`: 讀取redis
- `lab8`: 讀取及寫入activeMQ，不同的queue name，或TOPIC
- `lab9`: MONGODB CRUD
- `lab10`: viper + config
- `lab11`: parser xml
- `lab12`: 簡易的activeMQ連線。
- `lab13`: 簡易資料庫連線，展示解耦
- `lab14`: Graceful Shutdown

Each lab is organized as its own Go module with this baseline layout:

- `cmd/app`: program entrypoint
- `internal/app`: application wiring
- `internal/service`: reusable logic
- `internal/model`: domain types
- `test`: integration-test placeholder

## Notes

- `lab3` runs with `go run ./cmd/app` and treats `1234` and `12 34` as the same normalized input.
