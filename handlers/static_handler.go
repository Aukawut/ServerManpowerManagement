package handlers

import "github.com/gofiber/fiber/v2"

func DownloadFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	return c.SendFile("./public" + filename)
}
