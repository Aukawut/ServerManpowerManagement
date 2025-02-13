package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Aukawut/ServerManpowerManagement/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Crate Instance fiber
	app := fiber.New()

	// Register Routes
	routes.SetupRoutes(app)

	port := os.Getenv("PORT")

	if port == "" {
		port = "5555" // default ถ้าไม่มีใน .env
	}

	// Listen Server on Port 5250
	fmt.Println("Server running on port", port)
	log.Fatal(app.Listen(":" + port))

}
