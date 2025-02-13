package routes

import (
	"github.com/Aukawut/ServerManpowerManagement/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	user := app.Group("/auth")
	user.Post("/domain", handlers.LoginDomain)

}
