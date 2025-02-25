package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetUpStatic(app *fiber.App) {
	// Route group
	static := app.Group("/public")

	static.Static("/", "./public")

}
