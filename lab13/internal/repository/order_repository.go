package repository

import (
	"fmt"

	"lab13/cmd/internal/model"
)

// OrderRepository 負責 Order 聚合根的持久化操作
type OrderRepository struct{}

// Save 將一筆訂單（含多個明細）寫入資料庫
func (r *OrderRepository) Save(order model.OrderDao) error {
	// 每個 item 對應一筆 insert
	for _, item := range order.Items {
		// TODO: 替換成實際的 DB insert
		fmt.Printf("INSERT INTO orders (id, customer_no, item_no) VALUES ('%s', '%s', '%s')\n",
			order.ID, order.CustomerNo, item.ItemNo)
	}
	return nil
}
