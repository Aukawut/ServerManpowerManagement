package handlers

import (
	"github.com/Aukawut/ServerManpowerManagement/model"
	"github.com/gofiber/fiber/v2"
)

func LoginDomain(c *fiber.Ctx) error {

	req := model.BodyLoginDomain{}
	var err error

	if err = c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": err.Error(),
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Username and Password is required!",
		})
	}

	return c.JSON(fiber.Map{
		"err": false,
		"msg": "Login success!",
	})

}
