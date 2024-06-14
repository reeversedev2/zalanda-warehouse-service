package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/database"
)

func main() {
	database.ConnectDB()

	app := fiber.New()

	setupRoutes(app)

	app.Listen(":3000")
}
