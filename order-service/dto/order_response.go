package dto

import (
	"order-service/domain"
	"time"
)

type OrderAPIResponse struct {
	InvoiceID string     `json:"invoice_id,omitempty"`
	Quantity  int        `json:"quantity"`
	Total     float64    `json:"total"`
	Status    string     `json:"status"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type APIResponse struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func MapToOrderDTO(order *domain.OrderModel) OrderAPIResponse {
	invoiceID := ""
	if order.Invoice != nil && order.Invoice.InvoiceID != nil {
		invoiceID = *order.Invoice.InvoiceID
	}

	return OrderAPIResponse{
		Quantity:  order.Quantity,
		Total:     order.Total,
		Status:    order.Status,
		InvoiceID: invoiceID,
		CreatedAt: &order.CreatedAt,
		UpdatedAt: &order.UpdatedAt,
	}
}

func MapToOrderDTOs(orders []domain.OrderModel) []OrderAPIResponse {
	dtos := make([]OrderAPIResponse, 0, len(orders))
	for _, o := range orders {
		oCopy := o
		dtos = append(dtos, MapToOrderDTO(&oCopy))
	}
	return dtos
}
