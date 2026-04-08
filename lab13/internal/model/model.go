package model

import "time"

// OrderDao 是聚合根，一張訂單包含多個品項
type OrderDao struct {
	ID         string
	CustomerNo string
	Items      []OrderItemDao
}

// OrderItemDao 訂單明細
type OrderItemDao struct {
	ItemNo string
}

type ItemDao struct {
	ID   string
	Name string
}

type CustomerDao struct {
	ID   string
	Name string
}

// NewOrderDao 根據客編和商品編號 slice 建立一筆 OrderDao（聚合根）
func NewOrderDao(customerNo string, items []string) OrderDao {
	orderID := time.Now().Format("20060102150405")

	orderItems := make([]OrderItemDao, 0, len(items))
	for _, item := range items {
		orderItems = append(orderItems, OrderItemDao{ItemNo: item})
	}

	return OrderDao{
		ID:         orderID,
		CustomerNo: customerNo,
		Items:      orderItems,
	}
}
