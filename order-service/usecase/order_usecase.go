package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"order-service/domain"
	"os"
)

type orderUsecase struct {
	repo domain.OrderRepository
}

func NewOrderUsecase(repo domain.OrderRepository) domain.OrderUsecase {
	return &orderUsecase{repo}
}

func (u *orderUsecase) GetOrders() ([]domain.Order, error) {
	return u.repo.FindAll()
}

func (u *orderUsecase) GetOrderByID(id string) (domain.Order, error) {
	return u.repo.FindByID(id)
}

func (u *orderUsecase) CreateOrder(order *domain.Order) error {
	// Validate Product Existence
	productSvcURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productSvcURL == "" {
		productSvcURL = "http://localhost:8081"
	}

	productResp, err := http.Get(fmt.Sprintf("%s/products/%d", productSvcURL, order.ProductID))
	if err != nil || productResp.StatusCode != 200 {
		return errors.New("product not found or product service unavailable")
	}

	order.Status = "Pending"
	if err := u.repo.Create(order); err != nil {
		return err
	}

	// Call Payment Service
	paymentSvcURL := os.Getenv("PAYMENT_SERVICE_URL")
	if paymentSvcURL != "" {
		go func() {
			paymentPayload := map[string]interface{}{
				"order_id": order.ID,
				"amount":   order.Total,
			}
			body, _ := json.Marshal(paymentPayload)
			http.Post(paymentSvcURL+"/payments", "application/json", bytes.NewBuffer(body))
		}()
	}

	return nil
}

func (u *orderUsecase) UpdateOrder(id string, req *domain.Order) error {
	order, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Retain original keys
	req.ID = order.ID
	req.CreatedAt = order.CreatedAt

	return u.repo.Update(req)
}

func (u *orderUsecase) UpdateOrderStatus(id string, status string) error {
	return u.repo.UpdateStatus(id, status)
}

func (u *orderUsecase) DeleteOrder(id string) error {
	return u.repo.Delete(id)
}
