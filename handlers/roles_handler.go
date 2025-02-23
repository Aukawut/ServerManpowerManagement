package handlers

import (
	"database/sql"
	"fmt"

	"github.com/Aukawut/ServerManpowerManagement/config"
	"github.com/Aukawut/ServerManpowerManagement/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetRoles(c *fiber.Ctx) error {
	var roles []model.Roles

	_, ok := c.Locals("user").(jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid user data",
		})
	}

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	rows, errQuery := db.Query(`SELECT  [ROLE_ID]
      ,[ROLE_NAME]
      ,[CREATED_AT]
      ,[UPDATED_AT]
      ,[CREATED_BY]
      ,[UPDATED_BY]
  FROM [DB_MANPOWER_MGT].[dbo].[TBL_ROLES]`)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var role model.Roles

		errScan := rows.Scan(
			&role.ROLE_ID,
			&role.ROLE_NAME,
			&role.CREATED_AT,
			&role.UPDATED_AT,
			&role.CREATED_BY,
			&role.UPDATED_BY,
		)

		if errScan != nil {
			fmt.Println("Error : ", errScan.Error())
		} else {
			roles = append(roles, role)
		}
	}

	defer rows.Close()

	if len(roles) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"msg":     "",
			"status":  "Ok",
			"results": roles,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "",
			"status":  "",
			"results": roles,
		})
	}
}
