# lab3

This lab is a small standalone Go project that reads a line of text and normalizes it by removing all whitespace.

## Run

```bash
go run ./cmd/app
```

## Behavior

- Input `1234` returns `1234`
- Input `12 34` also returns `1234`

## Structure

- `cmd/app`: application entrypoint
- `internal/app`: input and output flow
- `internal/service`: reusable normalization logic
- `test`: integration-test placeholder
