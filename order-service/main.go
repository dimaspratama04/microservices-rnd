package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
	Status    string  `json:"status"`
}

var DB *gorm.DB

func initDB() {
	var err error
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = "host=localhost user=user password=password dbname=order_db port=5432 sslmode=disable"
	}
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

func SetupApp() *fiber.App {
	app := fiber.New()

	app.Get("/orders", func(c *fiber.Ctx) error {
		var orders []Order
		DB.Find(&orders)
		if len(orders) == 0 {
			return c.JSON(fiber.Map{
				"message": "No orders found",
				"data":    []Order{},
			})
		}
		return c.JSON(fiber.Map{
			"message": "Orders retrieved successfully",
			"data":    orders,
		})
	})

	app.Get("/orders/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var order Order
		if err := DB.First(&order, id).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.JSON(order)
	})

	app.Post("/orders", func(c *fiber.Ctx) error {
		order := new(Order)
		if err := c.BodyParser(order); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Validate Product Existence
		productSvcURL := os.Getenv("PRODUCT_SERVICE_URL")
		if productSvcURL == "" {
			productSvcURL = "http://localhost:8081"
		}

		productResp, err := http.Get(fmt.Sprintf("%s/products/%d", productSvcURL, order.ProductID))
		if err != nil || productResp.StatusCode != 200 {
			return c.Status(400).JSON(fiber.Map{
				"error": "Product not found or product service unavailable",
			})
		}

		order.Status = "Pending"
		DB.Create(&order)

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

		return c.Status(201).JSON(fiber.Map{
			"message": "Order created successfully",
			"data":    order,
		})
	})

	app.Put("/orders/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var order Order
		if err := DB.First(&order, id).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		if err := c.BodyParser(&order); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}
		DB.Save(&order)
		return c.JSON(fiber.Map{
			"message": "Order updated successfully",
			"data":    order,
		})
	})

	app.Delete("/orders/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var order Order
		if err := DB.First(&order, id).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		DB.Delete(&order)
		return c.JSON(fiber.Map{
			"message": "Order deleted successfully",
			"id":      id,
		})
	})

	return app
}

func main() {
	// Initialize DB if not in test mode
	if os.Getenv("GO_ENV") != "test" {
		initDB()
	}

	app := SetupApp()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
}
