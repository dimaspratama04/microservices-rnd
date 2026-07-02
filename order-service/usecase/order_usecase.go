package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"order-service/domain"
	"order-service/dto"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type OrderUsecase interface {
	GetOrders() ([]dto.OrderAPIResponse, error)
	GetOrderByID(id string) (dto.OrderAPIResponse, error)
	GetOrderByInvoiceID(invoiceID string) (dto.OrderAPIResponse, error)
	CreateOrder(ctx context.Context, req *dto.CreateOrderAPIRequest) (dto.OrderAPIResponse, error)
	UpdateOrder(id string, req *dto.CreateOrderAPIRequest) (dto.OrderAPIResponse, error)
	UpdateOrderStatus(id string, status string) error
	UpdateOrderStatusByInvoiceID(invoiceID string, status string) error
	DeleteOrder(id string) error
}

type orderUsecase struct {
	repo domain.OrderRepository
}

func NewOrderUsecase(repo domain.OrderRepository) OrderUsecase {
	return &orderUsecase{repo}
}

func (u *orderUsecase) GetOrders() ([]dto.OrderAPIResponse, error) {
	orders, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return dto.MapToOrderDTOs(orders), nil
}

func (u *orderUsecase) GetOrderByID(id string) (dto.OrderAPIResponse, error) {
	order, err := u.repo.FindByID(id)
	if err != nil {
		return dto.OrderAPIResponse{}, err
	}
	return dto.MapToOrderDTO(&order), nil
}

func (u *orderUsecase) GetOrderByInvoiceID(invoiceID string) (dto.OrderAPIResponse, error) {
	order, err := u.repo.FindByInvoiceID(invoiceID)
	if err != nil {
		return dto.OrderAPIResponse{}, err
	}
	return dto.MapToOrderDTO(&order), nil
}

func (u *orderUsecase) CreateOrder(ctx context.Context, req *dto.CreateOrderAPIRequest) (dto.OrderAPIResponse, error) {
	order := &domain.OrderModel{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	productSvcURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productSvcURL == "" {
		productSvcURL = "http://localhost:8081"
	}

	reqProduct, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/products/%d", productSvcURL, order.ProductID), nil)
	if err != nil {
		return dto.OrderAPIResponse{}, err
	}

	// Inject OTEL trace context
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(reqProduct.Header))

	productResp, err := http.DefaultClient.Do(reqProduct)
	if err != nil || productResp.StatusCode != 200 {
		return dto.OrderAPIResponse{}, errors.New("product not found or product service unavailable")
	}
	defer productResp.Body.Close()

	var productData struct {
		Data struct {
			Price float64 `json:"price"`
		} `json:"data"`
	}
	if err := json.NewDecoder(productResp.Body).Decode(&productData); err != nil {
		return dto.OrderAPIResponse{}, errors.New("failed to decode product response")
	}

	order.Total = float64(order.Quantity) * productData.Data.Price
	order.Status = "Pending"
	order.Invoice = &domain.InvoiceModel{}
	if err := u.repo.Create(order); err != nil {
		return dto.OrderAPIResponse{}, err
	}

	return dto.MapToOrderDTO(order), nil
}

func (u *orderUsecase) UpdateOrder(id string, req *dto.CreateOrderAPIRequest) (dto.OrderAPIResponse, error) {
	order, err := u.repo.FindByID(id)
	if err != nil {
		return dto.OrderAPIResponse{}, err
	}

	order.ProductID = req.ProductID
	order.Quantity = req.Quantity

	if err := u.repo.Update(&order); err != nil {
		return dto.OrderAPIResponse{}, err
	}

	return dto.MapToOrderDTO(&order), nil
}

func (u *orderUsecase) UpdateOrderStatus(id string, status string) error {
	return u.repo.UpdateStatus(id, status)
}

func (u *orderUsecase) UpdateOrderStatusByInvoiceID(invoiceID string, status string) error {
	return u.repo.UpdateOrderStatusByInvoiceID(invoiceID, status)
}

func (u *orderUsecase) DeleteOrder(id string) error {
	return u.repo.Delete(id)
}
