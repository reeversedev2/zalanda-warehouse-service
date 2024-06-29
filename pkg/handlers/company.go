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
	return c.Status(200).JSON(company)
}

func ListCompanies(c *fiber.Ctx) error {
	companies := []models.Company{}

	database.DB.Db.Scopes(utils.Paginate(companies, &pagination.Pagination{
		Limit: c.QueryInt("limit", 10),
		Page:  c.QueryInt("page", 1),
		Sort:  c.Query("sort", "id desc"),
	}, database.DB.Db)).Find(&companies)

	return c.Status(200).JSON(companies)
}

func FindCompanyById(id int, company *models.Company) (*models.Company, error) {
	database.DB.Db.Find(&company, "id = ?", id)
	if company.ID == 0 {
		return nil, errors.New("company does not exist")
	}
	return company, nil
}

func FindCompanyByName(name string, company *models.Company) error {
	database.DB.Db.Find(&company, "name = ?", name)
	if company.ID != 0 {
		error := fmt.Sprintf("company exists already with name %v", name)
		return errors.New(error)
	}
	return nil
}
