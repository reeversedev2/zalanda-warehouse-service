package handlers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/database"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/models"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/pagination"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/utils"
)

// Create a new company
func CreateCompany(c *fiber.Ctx) error {
	company := new(models.Company)
	if err := c.BodyParser(company); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err := FindCompanyByName(company.Name, &models.Company{})
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	database.DB.Db.Create(&company)

	// message := amqp091.Publishing{
	// 	ContentType: "text/plain",
	// 	Body:        []byte(c.Query("msg")),
	// }

	// // Attempt to publish a message to the queue.
	// if err := broker.BrokerCh.Publish(
	// 	"",              // exchange
	// 	"QueueService1", // queue name
	// 	false,           // mandatory
	// 	false,           // immediate
	// 	message,         // message to publish
	// ); err != nil {
	// 	return err
	// }

	return c.Status(fiber.StatusOK).JSON(company)
}

// Show all companies
func ListCompanies(c *fiber.Ctx) error {
	companies := []models.Company{}

	database.DB.Db.Scopes(utils.Paginate(companies, &pagination.Pagination{
		Limit: c.QueryInt("limit", 10),
		Page:  c.QueryInt("page", 1),
		Sort:  c.Query("sort", "id desc"),
	}, database.DB.Db)).Find(&companies)

	return c.Status(200).JSON(companies)
}

// Update the existing company
func UpdateCompany(c *fiber.Ctx) error {
	company := new(models.Company)
	companyId := c.Params("companyId")

	if err := c.BodyParser(company); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	database.DB.Db.Where("id=?", companyId).Updates(&company)

	return c.Status(fiber.StatusOK).JSON(company)
}

// Search company by Company ID
func FindCompanyById(id int, company *models.Company) (*models.Company, error) {
	database.DB.Db.Find(&company, "id = ?", id)
	if company.ID == 0 {
		return nil, errors.New("company does not exist")
	}
	return company, nil
}

// Search Company by Company Name
func FindCompanyByName(name string, company *models.Company) error {
	database.DB.Db.Find(&company, "name = ?", name)
	if company.ID != 0 {
		error := fmt.Sprintf("company exists already with name %v", name)
		return errors.New(error)
	}
	return nil
}
