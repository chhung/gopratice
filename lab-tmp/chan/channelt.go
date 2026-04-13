package main

import (
	"fmt"
)

func main() {
	ch_1 := make(chan string, 2)
	ch_1 <- "Hello, World!"
	fmt.Println(<-ch_1)
	fmt.Println("Hello, World!")
}
