package main

import (
	"bytes"
	"encoding/json"
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

	DB.AutoMigrate(&Order{})
}

func main() {
	app := fiber.New()

	// Initialize DB if not in test mode
	if os.Getenv("GO_ENV") != "test" {
		initDB()
	}

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
}
