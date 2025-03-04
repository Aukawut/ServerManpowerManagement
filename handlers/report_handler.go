package handlers

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/Aukawut/ServerManpowerManagement/config"
	"github.com/Aukawut/ServerManpowerManagement/model"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func SummaryHeadCountByDeptAndSex(c *fiber.Ctx) error {

	var results []model.ReportHeadCountByDeptAndSex

	date := c.Params("date")

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

	stmt := fmt.Sprintf(`SELECT LTRIM(RTRIM(UHR_Department)) as [DEPARTMENT],FEMALE,MALE,OTHER 
FROM [dbo].[func_GetSummaryDeptManpowerByDate]('%s') 
ORDER BY FEMALE DESC`, date)

	rows, errQuery := db.Query(stmt)

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
	date := c.Params("date")

	_, ok := c.Locals("user").(jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid User data",
		})
	}

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error Creating Connection: " + err.Error())
	}

	defer db.Close()

	rows, errQuery := db.Query(`SELECT LTRIM(RTRIM([DEPARTMENT])) as DEPARTMENT,HEAD_COUNT 
FROM (SELECT * FROM func_SummaryHeadCountDepartment(@date))a
ORDER BY HEAD_COUNT DESC`, sql.Named("date", date))

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
			fmt.Println("Error Scan error : ", errScan.Error())
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
	date := c.Params("date")

	_, ok := c.Locals("user").(jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid User data",
		})
	}

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error Creating connection: " + err.Error())
	}

	defer db.Close()

	rows, errQuery := db.Query(`SELECT LTRIM(RTRIM([POSITION])) as [POSITION]
      ,[HEAD_COUNT] FROM (
	  SELECT * FROM func_SummaryHeadCountPosition(@date)
	  ) a ORDER BY [HEAD_COUNT] DESC`, sql.Named("date", date))

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
			fmt.Println("Error Scan : ", errScan.Error())
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

func SummaryHeadByUserTypeAndDept(c *fiber.Ctx) error {

	var results []model.ReportHeadCountByDeptAndUsersType
	date := c.Params("date")

	_, ok := c.Locals("user").(jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": "JWT Invalid user data",
		})
	}

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error Creating connection: " + err.Error())
	}

	defer db.Close()

	rows, errQuery := db.Query(`
	SELECT LTRIM(RTRIM(UHR_Department)) as DEPARTMENT,Indirect as [INDIRECT],
	Direct as [DIRECT], [SGA],[OTHER],(ISNULL([INDIRECT],0) + ISNULL([DIRECT],0)
	+ ISNULL([SGA],0) + ISNULL([OTHER],0)) as [TOTAL]
	FROM (SELECT * FROM func_SummaryUserByUserTypeAndDept(@date)) at`, sql.Named("date", date))

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var result model.ReportHeadCountByDeptAndUsersType

		errScan := rows.Scan(
			&result.DEPARTMENT,
			&result.INDIRECT,
			&result.DIRECT,
			&result.SGA,
			&result.OTHER,
			&result.TOTAL,
		)

		if errScan != nil {
			fmt.Println("Error Scan : ", errScan.Error())
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

func SummaryHeadCountSex(c *fiber.Ctx) error {

	var results []model.ReportHeadCountBySex
	date := c.Params("date")
	_, ok := c.Locals("user").(jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid users data",
		})
	}

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error Open Connection: " + err.Error())
	}

	defer db.Close()

	rows, errQuery := db.Query(`
	  SELECT 
	  	CASE 
		WHEN UHR_Sex = 'F' THEN 'FEMALE' 
		WHEN UHR_Sex = 'M' THEN 'MALE'
		ELSE 'N/A' END AS [UHR_Sex],COUNT(*) as [AMOUNT]
  FROM (SELECT * FROM TBL_MANPOWER) a 
  WHERE UHR_StatusToUse = 'ENABLE' AND a.[DATE] = @date GROUP BY UHR_Sex ORDER BY COUNT(*) DESC`, sql.Named("date", date))

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var result model.ReportHeadCountBySex

		errScan := rows.Scan(
			&result.UHR_Sex,
			&result.AMOUNT,
		)

		if errScan != nil {
			fmt.Println("Error Scan value: ", errScan.Error())
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

func SummaryManpowerByGroupCategory(c *fiber.Ctx) error {

	var results []model.ManpowerByGroupCategory
	startDate := c.Params("start")
	endDate := c.Params("end")
	_, ok := c.Locals("user").(jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid Json data",
		})
	}

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error Connection: " + err.Error())
	}

	defer db.Close()

	rows, errQuery := db.Query(`
	SELECT m.[DATE],m.Daily_Operator,m.JP,m.Temporary,m.Permanent,m.Daily_Operator + m.JP + m.Permanent + m.Temporary as [Total] FROM (
	SELECT  a.*,(b.Person - (Daily_Operator + JP + Temporary)) as [Permanent] FROM (
	SELECT * 
	FROM dbo.Func_GetManpowerSummaryByGroupCategory(@start, @end)
) a LEFT JOIN (
	SELECT [DATE],COUNT(*) as [Person] FROM [dbo].TBL_MANPOWER 
	WHERE UHR_StatusToUse = 'ENABLE' GROUP BY [DATE]
) b ON a.DATE = b.DATE
  ) m`, sql.Named("start", startDate), sql.Named("end", endDate))

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var result model.ManpowerByGroupCategory

		errScan := rows.Scan(
			&result.DATE,
			&result.Daily_Operator,
			&result.JP,
			&result.Temporary,
			&result.Permanent,
			&result.Total,
		)

		if errScan != nil {
			fmt.Println("Error Scan Value : ", errScan.Error())
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

func SummaryManpowerByDate(c *fiber.Ctx) error {

	var results []model.ManpowerByDate
	startDate := c.Params("start")
	endDate := c.Params("end")
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

	stmt := fmt.Sprintf(`SELECT   m.[DATE] as [Date],
			COUNT(*) as [Person]
        FROM TBL_MANPOWER m
        LEFT JOIN V_Position p 
            ON m.UHR_POSITION COLLATE Thai_CI_AS = p.PHR_PName COLLATE Thai_CI_AS
        WHERE m.UHR_StatusToUse = 'ENABLE'
          AND m.[DATE] BETWEEN  '%s' AND '%s'
        GROUP BY m.[DATE] ORDER BY [Date],[Person] DESC`, startDate, endDate)

	rows, errQuery := db.Query(stmt)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var result model.ManpowerByDate

		errScan := rows.Scan(
			&result.Date,
			&result.Person,
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

func SummaryManpowerByGroupPosition(c *fiber.Ctx) error {

	var results []model.ManpowerByGroupPosition
	startDate := c.Params("start")
	endDate := c.Params("end")
	utypeParam := c.Params("utype")
	departmentParam := c.Params("department")

	department, _ := url.QueryUnescape(departmentParam)
	utype, err := url.QueryUnescape(utypeParam)

	if err != nil {
		fmt.Println("Error Convert String URL:", err)
		return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
	}

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

	stmt := fmt.Sprintf(`SELECT ISNULL(p.PHR_PGroupCode, 'N/A') AS PHR_PGroupCode,
			COUNT(*) as [Person]
        FROM TBL_MANPOWER m
        LEFT JOIN V_Position p 
            ON m.UHR_POSITION COLLATE Thai_CI_AS = p.PHR_PName COLLATE Thai_CI_AS
        WHERE m.UHR_StatusToUse = 'ENABLE'
          AND m.[DATE] BETWEEN  '%s' AND '%s'`, startDate, endDate)

	if department != "All" {
		stmt += fmt.Sprintf(` AND m.UHR_Department = '%s'`, department)
	}

	if utype != "All" {
		stmt += fmt.Sprintf(` AND m.UHR_OrgName = '%s'`, utype)
	}

	stmt += ` GROUP BY p.PHR_PGroupCode ORDER BY [Person] DESC`
	rows, errQuery := db.Query(stmt)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var result model.ManpowerByGroupPosition

		errScan := rows.Scan(

			&result.PHR_PGroupCode,
			&result.Person,
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
