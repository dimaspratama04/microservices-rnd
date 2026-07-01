package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"order-service/config"
	"order-service/handler"
	"order-service/repository"
	"order-service/usecase"
)

func SetupApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	app.Use(otelfiber.Middleware())

	// Init layers
	orderRepo := repository.NewOrderRepository(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)

	// Register handlers
	handler.NewOrderHandler(app, orderUsecase)

	return app
}

func main() {
	var db *gorm.DB
	if os.Getenv("GO_ENV") != "test" {
		db = config.InitDB()
	}

	tp := config.InitTracer()
	if tp != nil {
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()
	}

	app := SetupApp(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
}
