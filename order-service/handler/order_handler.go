package handler

import (
	"order-service/domain"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	usecase domain.OrderUsecase
}

func NewOrderHandler(app *fiber.App, usecase domain.OrderUsecase) {
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
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if len(orders) == 0 {
		return c.JSON(fiber.Map{
			"message": "No orders found",
			"data":    []domain.Order{},
		})
	}
	return c.JSON(fiber.Map{
		"message": "Orders retrieved successfully",
		"data":    orders,
	})
}

func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")
	order, err := h.usecase.GetOrderByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}
	return c.JSON(order)
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	order := new(domain.Order)
	if err := c.BodyParser(order); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.usecase.CreateOrder(order); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Order created successfully",
		"data":    order,
	})
}

func (h *OrderHandler) UpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	order := new(domain.Order)
	if err := c.BodyParser(order); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.usecase.UpdateOrder(id, order); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Order updated successfully",
		"data":    order,
	})
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	body := struct {
		Status string `json:"status"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.usecase.UpdateOrderStatus(id, body.Status); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Order status updated successfully",
		"status":  body.Status,
	})
}

func (h *OrderHandler) DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.usecase.DeleteOrder(id); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	return c.JSON(fiber.Map{
		"message": "Order deleted successfully",
		"id":      id,
	})
}
