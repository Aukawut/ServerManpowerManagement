package routes

import (
	"github.com/Aukawut/ServerManpowerManagement/handlers"
	"github.com/Aukawut/ServerManpowerManagement/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupHeadCountRoutes(app *fiber.App) {
	headcount := app.Group("/headcount")

	headcount.Get("/options/year", middleware.DecodeToken, handlers.GetYearOptions)
	headcount.Get("/options/moth", middleware.DecodeToken, handlers.GetMonthOptions)
	headcount.Post("/:month/:year", middleware.DecodeToken, handlers.SyncExcelFile)
	headcount.Get("/files/:month/:year", middleware.DecodeToken, handlers.GetHeadCountFile)
	headcount.Get("/cost/report/:uuid", middleware.DecodeToken, handlers.GetHeadCountAndCostByUuid)
	headcount.Delete("/:uuid", middleware.DecodeToken, handlers.DeleteHeadCountData)

}
