package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/handlers"
	"github.com/reeversedev2/zalanda-warehouse-service/pkg/utils"
)


func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header required",
		})
	}

	// extracting token from "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	token := tokenParts[1]
	claims, err := handlers.ValidateJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	c.Locals("user", claims)
	return c.Next()
}

func RoleMiddleware(requiredRole utils.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*handlers.Claims)

		// Convert user role string to enum
		userRole := utils.StringToRole(user.Role)
		if userRole == -1 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Invalid user role",
			})
		}

		// Check if user has sufficient permissions (higher enum value = more permissions)
		if userRole < requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}



