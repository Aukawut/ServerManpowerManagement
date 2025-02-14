package middleware

import (
	"strings"

	"github.com/Aukawut/ServerManpowerManagement/handlers"
	"github.com/gofiber/fiber/v2"
)

func DecodeToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": "Authorization header must start with 'Bearer '",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	decoded, errToken := handlers.VerifyToken(token) // Ensure VerifyToken is implemented correctly

	if errToken != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errToken.Error(),
		})
	}

	// Store decoded claims in context for downstream handlers
	c.Locals("user", decoded)
	return c.Next()
}
