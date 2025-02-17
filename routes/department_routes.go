package routes

import (
	"github.com/Aukawut/ServerManpowerManagement/handlers"
	"github.com/Aukawut/ServerManpowerManagement/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupDepartmentRoutes(app *fiber.App) {
	user := app.Group("/department")
	user.Get("/", middleware.DecodeToken, handlers.GetDepartment)
	user.Get("/users", middleware.DecodeToken, handlers.GetDepartmentOfUsers)

}
