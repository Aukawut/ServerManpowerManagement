package routes

import (
	"github.com/Aukawut/ServerManpowerManagement/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetUpStatic(app *fiber.App) {
	// Route group
	static := app.Group("/public")

	static.Static("/", "./public")
	static.Get("/download/:filename", handlers.DownloadFile)

}
