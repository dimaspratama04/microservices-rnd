package repository

import (
	"order-service/domain"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) FindAll() ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindByID(id string) (domain.Order, error) {
	var order domain.Order
	err := r.db.First(&order, id).Error
	return order, err
}

func (r *orderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&domain.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *orderRepository) Delete(id string) error {
	var order domain.Order
	if err := r.db.First(&order, id).Error; err != nil {
		return err
	}
	return r.db.Delete(&order).Error
}
