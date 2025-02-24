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
			"msg": "Invalid User data",
		})
	}

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error Connection: " + err.Error())
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
			fmt.Println("Error Scan: ", errScan.Error())
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

func GetUsersAuthen(c *fiber.Ctx) error {

	var users []model.UsersAuthentication

	_, ok := c.Locals("user").(jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid JsonWebToken",
		})
	}

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error Connection: " + err.Error())
	}

	defer db.Close()

	rows, errQuery := db.Query(`SELECT a.EMPLOYEE_CODE,u.UHR_FullName_th,u.UHR_POSITION,u.UHR_OrgName,
u.UHR_OrgCode,r.ROLE_ID,r.ROLE_NAME,a.ACTIVE,u.UHR_OrgGroup,u.UHR_Department FROM TBL_USERS a
LEFT JOIN V_Users u ON a.EMPLOYEE_CODE COLLATE Thai_CI_AS = 
u.UHR_EmpCode  COLLATE Thai_CI_AS
LEFT JOIN TBL_ROLES r ON a.ROLE_ID = r.ROLE_ID`)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var user model.UsersAuthentication

		errScan := rows.Scan(
			&user.EMPLOYEE_CODE,
			&user.UHR_FullName_th,
			&user.UHR_POSITION,
			&user.UHR_OrgName,
			&user.UHR_OrgCode,
			&user.ROLE_ID,
			&user.ROLE_NAME,
			&user.ACTIVE,
			&user.UHR_OrgGroup,
			&user.UHR_Department,
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

func GetManpowerByDate(c *fiber.Ctx) error {

	var users []model.UsersMaster
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
		fmt.Println("Error Create connection: " + err.Error())
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
      ,[UHR_OrgName] FROM (SELECT * FROM TBL_MANPOWER) a WHERE a.[DATE] = @date`, sql.Named("date", date))

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

func ActiveUser(c *fiber.Ctx) error {

	employeeCode := c.Params("employeeCode")

	var req model.ActiveUserBody

	fmt.Println(req.ActionBy)

	// แปลงข้อมูล JSON ที่รับมา
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Invalid request body",
		})
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

	_, errQuery := db.Exec(`UPDATE [dbo].[TBL_USERS] SET [ACTIVE] = @active,[UPDATED_AT] = GETDATE(),[UPDATED_BY] = @by WHERE [EMPLOYEE_CODE] = @code`,
		sql.Named("active", req.Active),
		sql.Named("code", employeeCode),
		sql.Named("by", req.ActionBy),
	)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	return c.JSON(fiber.Map{"err": false, "msg": "Updated", "status": "Ok"})

}

func IsDuplicatedUserAuthen(db *sql.DB, code string) bool {
	var empCode string

	err := db.QueryRow(`SELECT [EMPLOYEE_CODE] FROM [dbo].[TBL_USERS] WHERE [EMPLOYEE_CODE] = @code`, sql.Named("code", code)).Scan(&empCode)
	if err != nil {
		return false
	} else {

		return true
	}
}

func InsertAuthenUser(c *fiber.Ctx) error {

	var req model.InsertAuthenUserBody

	// แปลงข้อมูล JSON ที่รับมา
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Invalid request body",
		})
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

	var code string

	errorQuery := db.QueryRow(`SELECT [UHR_EmpCode] FROM [DB_MANPOWER_MGT].[dbo].[V_Users] 
	WHERE UHR_StatusToUse = 'ENABLE' AND UHR_EmpCode = @code`, sql.Named("code", req.EmployeeCode)).Scan(&code)

	if errorQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": "User isn't found!"})

	} else if IsDuplicatedUserAuthen(db, code) {

		return c.JSON(fiber.Map{"err": true, "msg": "User duplicated!"})
	}

	_, errorInsert := db.Exec(`INSERT INTO [dbo].[TBL_USERS] ([EMPLOYEE_CODE],[ROLE_ID],[ACTIVE],[CREATED_AT],[CREATED_BY]) 
	VALUES (@code,@role,@active,GETDATE(),@by)`,
		sql.Named("code", req.EmployeeCode),
		sql.Named("role", req.Role),
		sql.Named("active", req.Active),
		sql.Named("by", req.ActionBy),
	)

	if errorInsert != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errorInsert.Error()})
	}

	return c.JSON(fiber.Map{"err": false, "msg": "User authen inserted", "status": "Ok"})

}
