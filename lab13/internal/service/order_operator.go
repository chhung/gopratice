// 這個檔案裡面，實作的應該是跟訂單有關的行為
// 是符合使用者操作抽像化，而不是底層對資料庫操作
// 也就是封裝過後的行為
package order

import (
	"lab13/internal/model"
	"lab13/internal/repository"
)

// Service 定義訂單相關的操作介面，外部只認識這個 interface
type Service interface {
	AddtoCart(no string)
	RemoveFromCart(no string)
	CheckOut() error
}

// NewOrderService 建立訂單 service，回傳 interface 隱藏內部實作
func NewOrderService(customerNo string) Service {
	return &orderInfo{
		customerNo: customerNo,
		repo:       repository.OrderRepository{},
	}
}

type orderInfo struct {
	customerNo string
	items      []string
	repo       repository.OrderRepository
}

func (in *orderInfo) AddtoCart(no string) {
	in.items = append(in.items, no)
}

func (in *orderInfo) RemoveFromCart(no string) {
	for i, item := range in.items {
		if item == no {
			in.items = append(in.items[:i], in.items[i+1:]...)
			return
		}
	}
}

func (in *orderInfo) CheckOut() error {
	// 用 model 建立聚合根，再交給 repository 寫入
	o := model.NewOrderDao(in.customerNo, in.items)
	return in.repo.Save(o)
}
