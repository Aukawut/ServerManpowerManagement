package routes

import (
	"github.com/Aukawut/ServerManpowerManagement/handlers"
	"github.com/Aukawut/ServerManpowerManagement/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupReportRoutes(app *fiber.App) {
	// Route group
	user := app.Group("/report")

	// Route childen
	user.Get("/headcount/department/sex", middleware.DecodeToken, handlers.SummaryHeadCountByDeptAndSex)
	user.Get("/headcount/department", middleware.DecodeToken, handlers.SummaryHeadCountByDept)
	user.Get("/headcount/position", middleware.DecodeToken, handlers.SummaryHeadCountByPosition)

}
