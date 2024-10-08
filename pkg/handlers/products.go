package handlers

import (
	"errors"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/cache"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/database"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/models"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/pagination"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/producer"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/utils"
)

// Search Product by Product name
func FindProductByName(productName string, product *models.Product) error {
	database.DB.Db.Find(&product, "name = ?", productName)
	if product.ID != 0 {
		error := fmt.Sprintf("%v already exists", productName)
		return errors.New(error)
	}
	return nil
}

// Search Product by Product ID
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

// List all products
func ListProducts(c *fiber.Ctx) error {
	products := []models.Product{}

	database.DB.Db.Scopes(utils.Paginate(products, &pagination.Pagination{
		Limit: c.QueryInt("limit", 10),
		Page:  c.QueryInt("page", 1),
		// get sort value from the URL params
		Sort: c.Query("sort", fmt.Sprintf("id %s", c.Query("sort_by"))),
	}, database.DB.Db)).Find(&products)

	return c.Status(200).JSON(products)
}

// Create a new Product
func CreateProduct(c *fiber.Ctx) error {
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var productModal models.Product
	productExistsErr := FindProductByName(product.Name, &productModal)
	if productExistsErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": productExistsErr.Error(),
		})
	}

	var company models.Company
	foundCompany, companyErr := FindCompanyById(product.CompanyID, &company)
	if companyErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": companyErr.Error(),
		})
	}

	database.DB.Db.Create(&product)

	// product's company field will have company info
	product.Company = *foundCompany

	return c.Status(200).JSON(product)
}

// Creates fake products for testing and initial seeding
func CreateBatchProducts(c *fiber.Ctx) error {
	for i := 0; i < 1000; i++ {
		database.DB.Db.Create(&models.Product{
			Name:     gofakeit.ProductName(),
			Price:    gofakeit.Price(100, 1000),
			Category: gofakeit.RandomString([]string{"Electronics", "Clothing", "Food", "Furniture", "Books"}),
			Expire:   gofakeit.Date().String(),
			Status:   gofakeit.RandomString([]string{"NEW", "RETURNED", "DAMAGED", "REFURBISHED"}),
			Image:    gofakeit.URL(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Batch products created",
	})

}

// Update existing Product
func UpdateProduct(c *fiber.Ctx) error {
	product := new(models.Product)
	productId := c.Params("productId")

	// find product if it exists
	foundErr := database.DB.Db.First(&product, productId).Error
	if foundErr != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	// parse the request body
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// update the product
	result := database.DB.Db.Where("id=?", productId).Updates(&product)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}

	// update the product status in Redis via RabbitMQ
	err := UpdateAnalytics(utils.Message{
		"product": fmt.Sprintf("%s:%s", productId, product.Status),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(product)
}

// Update product status in Redis via RabbitMQ
func UpdateAnalytics(msg utils.Message) error {
	channelRabbitMQ, err := producer.GetChannel()
	if err != nil {
		return err
	}

	serializedMsg, err := utils.SerializeToBytes(msg)
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(serializedMsg),
	}

	// Attempt to publish a message to the queue.
	if err := channelRabbitMQ.Publish(
		"",                  // exchange
		"ProductsDashboard", // queue name
		false,               // mandatory
		false,               // immediate
		message,             // message to publish
	); err != nil {
		return err
	}
	return nil
}

func GetProductStatusEvents(c *fiber.Ctx) error {
	redis := cache.NewRedis()
	packed, err := redis.RedisClient.LRange("product_status:packed", 0, -1).Result()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	new, err := redis.RedisClient.LRange("product_status:new", 0, -1).Result()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// return the arrays of packed and new products with their product IDs
	return c.Status(200).JSON(fiber.Map{
		"packed": packed,
		"new":    new,
	})
}
