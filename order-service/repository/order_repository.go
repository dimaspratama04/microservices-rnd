package repository

import (
	"fmt"
	"order-service/domain"
	"time"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) FindAll() ([]domain.OrderModel, error) {
	var orders []domain.OrderModel
	err := r.db.Preload("Invoice").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindByID(id string) (domain.OrderModel, error) {
	var order domain.OrderModel
	err := r.db.Preload("Invoice").First(&order, "id = ?", id).Error
	return order, err
}

func (r *orderRepository) FindByInvoiceID(invoiceID string) (domain.OrderModel, error) {
	var order domain.OrderModel
	err := r.db.Preload("Invoice").Joins("Invoice").Where("\"Invoice\".invoice_id = ?", invoiceID).First(&order).Error
	return order, err
}

func (r *orderRepository) Create(order *domain.OrderModel) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if order.Invoice != nil {
			var lastInvoice domain.InvoiceModel
			currentYear := time.Now().Format("2006")

			err := tx.Where("invoice_id LIKE ?", "INV-"+currentYear+"-%").Order("created_at desc").First(&lastInvoice).Error

			var sequence int = 1
			if err == nil && lastInvoice.InvoiceID != nil && *lastInvoice.InvoiceID != "" {
				fmt.Sscanf(*lastInvoice.InvoiceID, "INV-"+currentYear+"-%04d", &sequence)
				sequence++
			}

			newInvoiceID := fmt.Sprintf("INV-%s-%04d", currentYear, sequence)
			order.Invoice.InvoiceID = &newInvoiceID
		}

		return tx.Create(order).Error
	})
}

func (r *orderRepository) Update(order *domain.OrderModel) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&domain.OrderModel{}).Where("id = ?", id).Update("status", status).Error
}

func (r *orderRepository) UpdateOrderStatusByInvoiceID(invoiceID string, status string) error {
	var order domain.OrderModel
	if err := r.db.Joins("Invoice").Where("\"Invoice\".invoice_id = ?", invoiceID).First(&order).Error; err != nil {
		return err
	}
	return r.db.Model(&order).Update("status", status).Error
}

func (r *orderRepository) Delete(id string) error {
	var order domain.OrderModel
	if err := r.db.First(&order, "id = ?", id).Error; err != nil {
		return err
	}
	return r.db.Delete(&order).Error
}
