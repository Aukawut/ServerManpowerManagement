package routes

import (
	"github.com/Aukawut/ServerManpowerManagement/handlers"
	"github.com/Aukawut/ServerManpowerManagement/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupManPowerRoutes(app *fiber.App) {
	// Route group
	user := app.Group("/manpower")

	// Route childen
	user.Get("/", middleware.DecodeToken, handlers.GetUsers)
	user.Get("/position/group/:start/:end", middleware.DecodeToken, handlers.GetManpowerByGroupPositionAndDate)
	user.Post("/", middleware.DecodeToken, handlers.AddManpower)
	user.Get("/:date", middleware.DecodeToken, handlers.GetManpowerByDate)
	user.Get("/group/category/:start/:end", middleware.DecodeToken, handlers.SummaryManpowerByGroupCategory)
	user.Get("/fliter/:department/:orgName/:orgGroup/:startDate/:endDate", middleware.DecodeToken, handlers.GetFliterManpowerMaster)
	user.Delete("/delete/:id/:action", middleware.DecodeToken, handlers.DeleteManpowerById)

}
