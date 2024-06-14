package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/database"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/models"
)

func ListProductById(c *fiber.Ctx) error {
	productId := c.Params("productId")
	product := models.Product{}

	database.DB.Db.Where("id=?", productId).First(&product)

	if product.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "No product was found",
		})
	}

	return c.Status(200).JSON(product)
}

func ListProducts(c *fiber.Ctx) error {
	products := []models.Product{}
	database.DB.Db.Find(&products)

	return c.Status(200).JSON(products)
}

func CreateProduct(c *fiber.Ctx) error {
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	database.DB.Db.Create(&product)

	return c.Status(200).JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	product := new(models.Product)
	productId := c.Params("productId")

	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	database.DB.Db.Where("id=?", productId).Updates(&product)

	return c.Status(200).JSON(product)

}
