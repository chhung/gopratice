package main

import (
	"fmt"
	order "lab13/internal/service"
)

func main() {
	svc := order.NewOrderService("201805289534")
	svc.AddtoCart("ITEM-A")
	svc.AddtoCart("ITEM-B")
	if err := svc.CheckOut(); err != nil {
		fmt.Println("checkout failed:", err)
		return
	}
	fmt.Println("done")
}
