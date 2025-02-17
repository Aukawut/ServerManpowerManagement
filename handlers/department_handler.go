package handlers

import (
	"database/sql"
	"fmt"

	"github.com/Aukawut/ServerManpowerManagement/config"
	"github.com/Aukawut/ServerManpowerManagement/model"
	"github.com/gofiber/fiber/v2"
)

func GetDepartment(c *fiber.Ctx) error {

	//Loading connection string
	connString := config.LoadDatabaseConfig()

	var departments []model.Department

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	rows, errQuery := db.Query(`SELECT DHR_DName,DHR_DDes FROM [dbo].[HRS_Department] WHERE DHR_DSt = 1 ORDER BY DHR_DDes`)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var department model.Department

		errScan := rows.Scan(
			&department.DHR_DName,
			&department.DHR_DDes,
		)

		if errScan != nil {
			return c.JSON(fiber.Map{"err": true, "msg": errScan.Error()})

		} else {
			// Appended value to departments
			departments = append(departments, department)
		}
	}

	defer rows.Close()

	if len(departments) > 0 {

		return c.JSON(fiber.Map{"err": false, "msg": "", "results": ""})
	} else {
		return c.JSON(fiber.Map{"err": true, "msg": "Departments isn't found.", "results": departments})

	}

}

func GetDepartmentOfUsers(c *fiber.Ctx) error {

	//Loading connection string
	connString := config.LoadDatabaseConfig()

	var departments []model.DepartmentOfUsers

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	rows, errQuery := db.Query(`
	SELECT 
     [UHR_Department]
   FROM [DB_MANPOWER_MGT].[dbo].[V_Users] GROUP BY [UHR_Department]
	`)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var department model.DepartmentOfUsers

		errScan := rows.Scan(
			&department.UHR_Department,
		)

		if errScan != nil {
			return c.JSON(fiber.Map{"err": true, "msg": errScan.Error()})

		} else {
			// Appended value to departments
			departments = append(departments, department)
		}
	}

	defer rows.Close()

	if len(departments) > 0 {

		return c.JSON(fiber.Map{"err": false, "msg": "", "results": departments, "status": "Ok"})
	} else {
		return c.JSON(fiber.Map{"err": true, "msg": "Departments isn't found.", "results": ""})

	}

}
