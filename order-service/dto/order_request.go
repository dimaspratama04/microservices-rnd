package dto

type CreateOrderAPIRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}
