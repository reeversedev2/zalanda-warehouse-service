package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/database"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/producer"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/router"
)

func main() {
	database.ConnectDB()
	producer.StartConnect()

	app := fiber.New()
	router.Routes(app)

	app.Use(cors.New())

	app.Listen(":3000")
}
