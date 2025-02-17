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

func GetOrganizeNameMaster(c *fiber.Ctx) error {

	var organizeName []model.OrganizeNameMaster

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

	rows, errQuery := db.Query(`
  SELECT 
     [UHR_OrgName]
  FROM [DB_MANPOWER_MGT].[dbo].[V_Users] GROUP BY [UHR_OrgName]`)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var organize model.OrganizeNameMaster

		errScan := rows.Scan(
			&organize.UHR_OrgName,
		)

		if errScan != nil {
			fmt.Println("Error : ", errScan.Error())
		} else {
			organizeName = append(organizeName, organize)
		}
	}

	defer rows.Close()

	if len(organizeName) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"msg":     "",
			"status":  "Ok",
			"results": organizeName,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "",
			"status":  "",
			"results": organizeName,
		})
	}

}

func GetOrganizeGroup(c *fiber.Ctx) error {

	var organizeGroup []model.OrganizeGroupMaster

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

	rows, errQuery := db.Query(`
  SELECT 
     [UHR_OrgGroup]
  FROM [DB_MANPOWER_MGT].[dbo].[V_Users] GROUP BY [UHR_OrgGroup]`)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var orgGroup model.OrganizeGroupMaster

		errScan := rows.Scan(
			&orgGroup.UHR_OrgGroup,
		)

		if errScan != nil {
			fmt.Println("Error : ", errScan.Error())
		} else {
			organizeGroup = append(organizeGroup, orgGroup)
		}
	}

	defer rows.Close()

	if len(organizeGroup) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"msg":     "",
			"status":  "Ok",
			"results": organizeGroup,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "",
			"status":  "",
			"results": organizeGroup,
		})
	}

}
