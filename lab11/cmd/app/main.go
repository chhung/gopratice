package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"lab11/internal/app"
)

func main() {
	file := flag.String("file", "", "XML filename under resources/xmldata/ (e.g. coupondt.txt)")
	flag.Parse()

	if *file == "" {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/app -file <filename>")
		fmt.Fprintln(os.Stderr, "example: go run ./cmd/app -file coupondt.txt")
		os.Exit(1)
	}

	filePath := filepath.Join("resources", "xmldata", *file)

	if err := app.Run(filePath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
