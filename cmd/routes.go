package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/handlers"
)

func setupRoutes(app *fiber.App) {
	// Index
	app.Get("/", handlers.Index)

	// Products
	app.Get("/product/:productId", handlers.ListProductById)
	app.Get("/products", handlers.ListProducts)
	app.Post("/product", handlers.CreateProduct)
	app.Patch("/product/:productId", handlers.UpdateProduct)

	// Batch products
	app.Post("/batch-products", handlers.CreateBatchProducts)
}
