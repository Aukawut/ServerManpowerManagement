package routes

import (
	"github.com/Aukawut/ServerManpowerManagement/handlers"
	"github.com/Aukawut/ServerManpowerManagement/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupOrganizeRoutes(app *fiber.App) {
	// Route group
	user := app.Group("/organize")

	// Route childen
	user.Get("/name", middleware.DecodeToken, handlers.GetOrganizeNameMaster)
	user.Get("/group", middleware.DecodeToken, handlers.GetOrganizeGroup)

}
