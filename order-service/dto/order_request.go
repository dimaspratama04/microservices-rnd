package dto

type OrderAPIRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}
