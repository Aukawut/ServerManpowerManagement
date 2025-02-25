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
	user.Get("/system/authen", middleware.DecodeToken, handlers.GetUsersAuthen)
	user.Get("/manpower/:date", middleware.DecodeToken, handlers.GetManpowerByDate)
	user.Put("/auth/active/:employeeCode", middleware.DecodeToken, handlers.ActiveUser)
	user.Post("/auth/insert", middleware.DecodeToken, handlers.InsertAuthenUser)
	user.Delete("/auth/delete/:code", middleware.DecodeToken, handlers.DeleteAuthenUser)

}
