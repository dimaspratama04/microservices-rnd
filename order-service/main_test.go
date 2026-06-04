package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&Order{})
}

func TestOrderCRUD(t *testing.T) {
	os.Setenv("GO_ENV", "test")
	setupTestDB()

	app := fiber.New()
	
	// Register routes (manually for test to avoid app.Listen)
	app.Get("/orders", func(c *fiber.Ctx) error {
		var orders []Order
		DB.Find(&orders)
		return c.JSON(orders)
	})
	app.Get("/orders/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var order Order
		if err := DB.First(&order, id).Error; err != nil {
			return c.Status(404).SendString("Order not found")
		}
		return c.JSON(order)
	})
	app.Post("/orders", func(c *fiber.Ctx) error {
		order := new(Order)
		if err := c.BodyParser(order); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		order.Status = "Pending"
		DB.Create(&order)
		return c.Status(201).JSON(order)
	})
	app.Put("/orders/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var order Order
		if err := DB.First(&order, id).Error; err != nil {
			return c.Status(404).SendString("Order not found")
		}
		if err := c.BodyParser(&order); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		DB.Save(&order)
		return c.JSON(order)
	})
	app.Delete("/orders/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var order Order
		if err := DB.First(&order, id).Error; err != nil {
			return c.Status(404).SendString("Order not found")
		}
		DB.Delete(&order)
		return c.SendString("Order deleted")
	})

	// 1. Create
	newOrder := Order{ProductID: 1, Quantity: 2, Total: 100.0}
	body, _ := json.Marshal(newOrder)
	req := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)

	var createdOrder Order
	json.NewDecoder(resp.Body).Decode(&createdOrder)
	assert.Equal(t, uint(1), createdOrder.ProductID)
	assert.Equal(t, "Pending", createdOrder.Status)

	// 2. Read All
	req = httptest.NewRequest("GET", "/orders", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)

	// 3. Read One
	req = httptest.NewRequest("GET", "/orders/1", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)

	// 4. Update
	createdOrder.Status = "Paid"
	body, _ = json.Marshal(createdOrder)
	req = httptest.NewRequest("PUT", "/orders/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)

	var updatedOrder Order
	json.NewDecoder(resp.Body).Decode(&updatedOrder)
	assert.Equal(t, "Paid", updatedOrder.Status)

	// 5. Delete
	req = httptest.NewRequest("DELETE", "/orders/1", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)

	// 6. Verify Delete
	req = httptest.NewRequest("GET", "/orders/1", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, 404, resp.StatusCode)
}
