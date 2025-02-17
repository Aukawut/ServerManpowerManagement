package handlers

import (
	"database/sql"
	"fmt"

	"github.com/Aukawut/ServerManpowerManagement/config"
	"github.com/Aukawut/ServerManpowerManagement/model"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetUsers(c *fiber.Ctx) error {

	var users []model.UsersMaster

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

	rows, errQuery := db.Query(`SELECT [UHR_EmpCode]
      ,[UHR_OrgCode]
      ,[UHR_FullName_th]
      ,[UHR_FullName_en]
      ,[UHR_Department]
      ,[UHR_GroupDepartment]
      ,[UHR_POSITION]
      ,[UHR_GMail]
      ,[UHR_Sex]
      ,[UHR_StatusToUse]
      ,[UHR_OrgGroup]
      ,[UHR_OrgName]
  FROM [DB_MANPOWER_MGT].[dbo].[V_Users]`)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var user model.UsersMaster

		errScan := rows.Scan(
			&user.UHR_EmpCode,
			&user.UHR_OrgCode,
			&user.UHR_FullName_th,
			&user.UHR_FullName_en,
			&user.UHR_Department,
			&user.UHR_GroupDepartment,
			&user.UHR_POSITION,
			&user.UHR_GMail,
			&user.UHR_Sex,
			&user.UHR_StatusToUse,
			&user.UHR_OrgGroup,
			&user.UHR_OrgName,
		)

		if errScan != nil {
			fmt.Println("Error : ", errScan.Error())
		} else {
			users = append(users, user)
		}
	}

	defer rows.Close()

	if len(users) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"msg":     "",
			"status":  "Ok",
			"results": users,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "",
			"status":  "",
			"results": users,
		})
	}

}
