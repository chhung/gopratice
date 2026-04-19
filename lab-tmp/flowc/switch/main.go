package main

import (
	"fmt"
)

func main() {
	switch i := 4; i {
	case 1:
		fmt.Println("one")
		fallthrough
	case 2:
		fmt.Println("two")

	case 3:
		fmt.Println("three")
	default:
		fmt.Println("default")
	}

	fmt.Println("hello")
}
