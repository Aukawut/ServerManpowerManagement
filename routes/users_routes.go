package routes

import (
	"github.com/Aukawut/ServerManpowerManagement/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	user := app.Group("/users")
	user.Get("/", handlers.GetUsers)

}
