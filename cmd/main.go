package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/database"
)

func main() {
	database.ConnectDB()

	app := fiber.New()

	setupRoutes(app)

	app.Use(cors.New())

	app.Listen(":3000")
}
