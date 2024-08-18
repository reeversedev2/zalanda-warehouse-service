package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/handlers"
)

func Routes(app *fiber.App) {
	// Index
	app.Get("/", handlers.Index)

	// Products
	app.Get("/api/product/:productId", handlers.ListProductById)
	app.Get("/api/products", handlers.ListProducts)
	app.Post("/api/product", handlers.CreateProduct)
	app.Patch("/api/product/:productId", handlers.UpdateProduct)

	// Batch products
	app.Post("/batch-products", handlers.CreateBatchProducts)

	// Companies
	app.Post("/api/company", handlers.CreateCompany)
	app.Get("/api/companies", handlers.ListCompanies)
	app.Put("/api/company/:companyId", handlers.UpdateCompany)

	app.Get("/message/send", handlers.SendMessage)
	// Get status from Redis
	app.Get("/product/status", handlers.ReceiveMessage)
}
