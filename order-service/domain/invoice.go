package domain

import (
	"time"

	"gorm.io/gorm"
)

type InvoiceModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	OrderID   uint           `gorm:"uniqueIndex" json:"order_id"`
	InvoiceID *string        `gorm:"uniqueIndex;type:varchar(50);default:null" json:"invoice_id"`
}
