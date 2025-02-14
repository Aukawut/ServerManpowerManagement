package routes

import (
	"github.com/Aukawut/ServerManpowerManagement/handlers"
	"github.com/Aukawut/ServerManpowerManagement/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	// Route group
	user := app.Group("/users")

	// Route childen
	user.Get("/", middleware.DecodeToken, handlers.GetUsers)

}
