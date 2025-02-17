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

func SummaryHeadCountByDeptAndSex(c *fiber.Ctx) error {

	var results []model.ReportHeadCountByDeptAndSex

	_, ok := c.Locals("user").(jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid user",
		})
	}

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer db.Close()

	rows, errQuery := db.Query(`SELECT LTRIM(RTRIM([UHR_Department])) as DEPARTMENT,FEMALE,MALE,OTHER FROM [dbo].[V_SexByDepartment] ORDER BY FEMALE DESC`)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var result model.ReportHeadCountByDeptAndSex

		errScan := rows.Scan(
			&result.DEPARTMENT,
			&result.FEMALE,
			&result.MALE,
			&result.OTHER,
		)

		if errScan != nil {
			fmt.Println(errScan.Error())
		} else {
			results = append(results, result)
		}
	}

	defer rows.Close()

	if len(results) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"msg":     "",
			"status":  "Ok",
			"results": results,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "",
			"status":  "",
			"results": results,
		})
	}

}

func SummaryHeadCountByDept(c *fiber.Ctx) error {

	var results []model.ReportHeadCountByDept

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

	rows, errQuery := db.Query(`SELECT LTRIM(RTRIM([DEPARTMENT])) as DEPARTMENT,HEAD_COUNT FROM [dbo].[V_SummaryHeadCountDepartment] ORDER BY HEAD_COUNT DESC`)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var result model.ReportHeadCountByDept

		errScan := rows.Scan(
			&result.DEPARTMENT,
			&result.HEAD_COUNT,
		)

		if errScan != nil {
			fmt.Println("Error : ", errScan.Error())
		} else {
			results = append(results, result)
		}
	}

	defer rows.Close()

	if len(results) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"msg":     "",
			"status":  "Ok",
			"results": results,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "",
			"status":  "",
			"results": results,
		})
	}

}

func SummaryHeadCountByPosition(c *fiber.Ctx) error {

	var results []model.ReportHeadCountByPosition

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

	rows, errQuery := db.Query(`SELECT LTRIM(RTRIM([POSITION])) as [POSITION]
      ,[HEAD_COUNT] FROM [DB_MANPOWER_MGT].[dbo].[V_SummaryHeadCountPosition] ORDER BY [HEAD_COUNT] DESC`)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var result model.ReportHeadCountByPosition

		errScan := rows.Scan(
			&result.POSITION,
			&result.HEAD_COUNT,
		)

		if errScan != nil {
			fmt.Println("Error : ", errScan.Error())
		} else {
			results = append(results, result)
		}
	}

	defer rows.Close()

	if len(results) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"msg":     "",
			"status":  "Ok",
			"results": results,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "",
			"status":  "",
			"results": results,
		})
	}

}
