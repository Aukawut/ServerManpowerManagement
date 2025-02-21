package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	SetupUserRoutes(app)
	SetupAuthRoutes(app)
	SetupDepartmentRoutes(app)
	SetupOrganizeRoutes(app)
	SetupReportRoutes(app)
	SetupManPowerRoutes(app)
}
