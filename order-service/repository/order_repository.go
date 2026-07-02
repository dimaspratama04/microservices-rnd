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

func (r *orderRepository) FindAll() ([]domain.OrderModel, error) {
	var orders []domain.OrderModel
	err := r.db.Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindByID(id string) (domain.OrderModel, error) {
	var order domain.OrderModel
	err := r.db.First(&order, id).Error
	return order, err
}

func (r *orderRepository) Create(order *domain.OrderModel) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) Update(order *domain.OrderModel) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&domain.OrderModel{}).Where("id = ?", id).Update("status", status).Error
}

func (r *orderRepository) Delete(id string) error {
	var order domain.OrderModel
	if err := r.db.First(&order, id).Error; err != nil {
		return err
	}
	return r.db.Delete(&order).Error
}
