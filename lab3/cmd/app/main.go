package main

import (
	"log"
	"os"

	"lab3/internal/app"
)

func main() {
	if err := app.Run(os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
