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
	user.Get("/headcount/department/sex/:date", middleware.DecodeToken, handlers.SummaryHeadCountByDeptAndSex)
	user.Get("/headcount/department/utype/:date", middleware.DecodeToken, handlers.SummaryHeadByUserTypeAndDept)
	user.Get("/headcount/department/:date", middleware.DecodeToken, handlers.SummaryHeadCountByDept)
	user.Get("/headcount/position/:date", middleware.DecodeToken, handlers.SummaryHeadCountByPosition)
	user.Get("/headcount/sex/:date", middleware.DecodeToken, handlers.SummaryHeadCountSex)
	user.Get("/headcount/:start/:end", middleware.DecodeToken, handlers.SummaryManpowerByDate)
	user.Get("/manpower/position/group/:start/:end/:department/:utype", middleware.DecodeToken, handlers.SummaryManpowerByGroupPosition)

}
