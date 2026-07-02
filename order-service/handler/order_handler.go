package handler

import (
	"order-service/dto"
	"order-service/usecase"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	usecase usecase.OrderUsecase
}

func NewOrderHandler(app *fiber.App, usecase usecase.OrderUsecase) {
	handler := &OrderHandler{usecase}

	app.Get("/orders", handler.GetOrders)
	app.Get("/orders/:id", handler.GetOrderByID)
	app.Post("/orders", handler.CreateOrder)
	app.Put("/orders/:id", handler.UpdateOrder)
	app.Patch("/orders/:id/status", handler.UpdateOrderStatus)
	app.Delete("/orders/:id", handler.DeleteOrder)
}

func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {
	orders, err := h.usecase.GetOrders()
	if err != nil {
		return c.Status(500).JSON(dto.APIResponse{Message: err.Error()})
	}
	if len(orders) == 0 {
		return c.JSON(dto.APIResponse{
			Message: "No orders found",
			Data:    []dto.OrderAPIResponse{},
		})
	}
	return c.JSON(dto.APIResponse{
		Message: "Orders retrieved successfully",
		Data:    orders,
	})
}

func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")
	order, err := h.usecase.GetOrderByID(id)
	if err != nil {
		return c.Status(404).JSON(dto.APIResponse{Message: "Order not found"})
	}
	return c.JSON(dto.APIResponse{
		Message: "Order retrieved successfully",
		Data:    order,
	})
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var payload dto.OrderAPIRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(dto.APIResponse{Message: "Invalid request body"})
	}

	order, err := h.usecase.CreateOrder(&payload)
	if err != nil {
		return c.Status(400).JSON(dto.APIResponse{Message: err.Error()})
	}

	return c.Status(201).JSON(dto.APIResponse{
		Message: "Order created successfully",
		Data:    order,
	})
}

func (h *OrderHandler) UpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload dto.OrderAPIRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(dto.APIResponse{Message: "Invalid request body"})
	}

	order, err := h.usecase.UpdateOrder(id, &payload)
	if err != nil {
		return c.Status(404).JSON(dto.APIResponse{Message: err.Error()})
	}

	return c.JSON(dto.APIResponse{
		Message: "Order updated successfully",
		Data:    order,
	})
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	body := struct {
		Status string `json:"status"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(dto.APIResponse{Message: "Invalid request body"})
	}

	if err := h.usecase.UpdateOrderStatus(id, body.Status); err != nil {
		return c.Status(404).JSON(dto.APIResponse{Message: err.Error()})
	}

	return c.JSON(dto.APIResponse{
		Message: "Order status updated successfully",
		Data:    fiber.Map{"status": body.Status},
	})
}

func (h *OrderHandler) DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.usecase.DeleteOrder(id); err != nil {
		return c.Status(404).JSON(dto.APIResponse{Message: "Order not found"})
	}

	return c.JSON(dto.APIResponse{
		Message: "Order deleted successfully",
		Data:    fiber.Map{"id": id},
	})
}
