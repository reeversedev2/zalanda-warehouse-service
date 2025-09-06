package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/handlers"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/middleware"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/utils"
)

func Routes(app *fiber.App) {
	// Index
	app.Get("/", handlers.Index)

	// Authentication routes (public)
	app.Post("/api/auth/register", handlers.Register)
	app.Post("/api/auth/login", handlers.Login)

	// Protected routes (require authentication)
	app.Use("/api/protected", middleware.AuthMiddleware)
	app.Get("/api/protected/profile", handlers.GetProfile)

	// Products (require Picker role or above)
	app.Use("/api/products", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.Picker))
	app.Get("/api/products", handlers.ListProducts)
	app.Get("/api/products/:productId", handlers.ListProductById)
	app.Post("/api/products", handlers.CreateProduct)
	app.Patch("/api/products/:productId", handlers.UpdateProduct)

	// Dashboard (require Manager role)
	app.Use("/api/dashboard", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.Manager))
	// Dashboard related APIs can be added here
	// app.Get("/api/dashboard", handlers.Dashboard)

	// Analytics (require Manager role)
	app.Use("/api/analytics", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.Manager))
	app.Get("/api/analytics/product/status", handlers.GetProductStatusEvents)

	// Batch products (require Manager role)
	app.Use("/batch-products", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.Manager))
	app.Post("/batch-products", handlers.CreateBatchProducts)

	// Companies (require Manager role)
	app.Use("/api/company", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.Manager))
	app.Post("/api/company", handlers.CreateCompany)
	app.Get("/api/companies", handlers.ListCompanies)
	app.Put("/api/company/:companyId", handlers.UpdateCompany)

}
