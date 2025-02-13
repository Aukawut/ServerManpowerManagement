package main

import (
	"github.com/Aukawut/ServerManpowerManagement/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {

	// Crate Instance fiber
	app := fiber.New()

	// Register Routes
	routes.SetupUserRoutes(app)

	app.Listen(":5250")
}
