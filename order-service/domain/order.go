package domain

import (
	"time"

	"gorm.io/gorm"
)

type OrderModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	ProductID uint           `json:"product_id"`
	Quantity  int            `json:"quantity"`
	Total     float64        `gorm:"type:decimal(12,2)" json:"total"`
	Status    string         `json:"status"`
	Invoice   *InvoiceModel  `gorm:"foreignKey:OrderID" json:"invoice,omitempty"`
}

type OrderRepository interface {
	FindAll() ([]OrderModel, error)
	FindByID(id string) (OrderModel, error)
	FindByInvoiceID(invoiceID string) (OrderModel, error)
	Create(order *OrderModel) error
	Update(order *OrderModel) error
	UpdateStatus(id string, status string) error
	UpdateOrderStatusByInvoiceID(invoiceID string, status string) error
	Delete(id string) error
}

