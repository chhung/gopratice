package app

import (
	"fmt"

	"lab1/internal/model"
	"lab1/internal/service"
)

func Run() {
	fmt.Println("Hello, World!")

	box := "A box of chocolates"
	money := 100
	fmt.Printf("I have %d dollars and a %s.\n", money, box)

	addOp := model.AddOperation{}
	subtractOp := model.SubtractOperation{}

	fmt.Printf("The sum of 5 and 3 is %d.\n", service.Execute(addOp, 5, 3))
	fmt.Printf("The result of add operation is %d.\n", service.Execute(addOp, 10, 20))
	fmt.Printf("The result of subtract operation is %d.\n", service.Execute(subtractOp, 10, 20))
}