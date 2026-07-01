package domain

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	ProductID uint           `json:"product_id"`
	Quantity  int            `json:"quantity"`
	Total     float64        `json:"total"`
	Status    string         `json:"status"`
}

type OrderRepository interface {
	FindAll() ([]Order, error)
	FindByID(id string) (Order, error)
	Create(order *Order) error
	Update(order *Order) error
	UpdateStatus(id string, status string) error
	Delete(id string) error
}

type OrderUsecase interface {
	GetOrders() ([]Order, error)
	GetOrderByID(id string) (Order, error)
	CreateOrder(order *Order) error
	UpdateOrder(id string, order *Order) error
	UpdateOrderStatus(id string, status string) error
	DeleteOrder(id string) error
}
