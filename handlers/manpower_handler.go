package handlers

import (
	"database/sql"
	"fmt"
	"strconv"

	"net/url"

	"github.com/Aukawut/ServerManpowerManagement/config"
	"github.com/Aukawut/ServerManpowerManagement/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type CountUser struct {
	AMONT int
}

func GetManpowerTerminations(c *fiber.Ctx) error {

	var users []model.SummaryManTermination

	var date = c.Params("date")

	fmt.Println(date)

	_, ok := c.Locals("user").(jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid Users data",
		})
	}

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error Create connection: " + err.Error())
	}

	defer db.Close()

	rows, errQuery := db.Query(`SELECT 
      [UHR_Department],
	  COUNT(*) as [PERSON]
    
  	  FROM [DB_MANPOWER_MGT].[dbo].[V_AllUserPSTH] WHERE  UHR_LastDate IS NOT NULL AND UHR_Department IS NOT NULL 
 	  AND UHR_StatusToUse = 'DISABLE'
     --AND CONVERT(VARCHAR(10),[UHR_LastDate],120) = @date
	 GROUP BY [UHR_Department] ORDER BY  COUNT(*) DESC
`, sql.Named("date", date))

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var user model.SummaryManTermination

		errScan := rows.Scan(
			&user.UHR_Department,
			&user.PERSON,
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

func CheckDuplicated(date string, code string) bool {
	// โหลดค่าสำหรับเชื่อมต่อฐานข้อมูล
	connString := config.LoadDatabaseConfig()

	// เปิดการเชื่อมต่อฐานข้อมูล
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error creating Connection: " + err.Error())
		return true
	}

	var duplicateEmp string
	stmt := `SELECT  UHR_EmpCode FROM [DB_MANPOWER_MGT].[dbo].[TBL_MANPOWER] WHERE DATE = @date AND UHR_EmpCode  = @code`

	errDuplicate := db.QueryRow(stmt).Scan(&duplicateEmp)

	if errDuplicate != nil {
		return false
	}

	if duplicateEmp != "" {
		//Update
		return true
	} else {
		return false
	}

}

func CheckUserOk(req model.BodyAddManpower) (bool, string) {

	// โหลดค่าสำหรับเชื่อมต่อฐานข้อมูล
	connString := config.LoadDatabaseConfig()

	// เปิดการเชื่อมต่อฐานข้อมูล
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error Connection: " + err.Error())
		return false, ""
	}

	defer db.Close()

	var rowsError string

	html := `<style>
    .table {
    
      width : 100%
    }
  </style>
  <table class="table">
    <thead>
      <tr>
        <th style="text-align:center;background: #1677FF;color:#fff;padding : 3px;">แถวที่</th>
        <th style="text-align:center;background: #1677FF;color:#fff;padding : 3px;">ข้อมูล</th>
      </tr>
    </thead>
    <tbody>`

	// Check User not found
	for index, u := range req.Users {

		var countUser int
		stmt := fmt.Sprintf(`SELECT COUNT(*)  FROM V_Users WHERE UHR_StatusToUse = 'ENABLE' AND UHR_EmpCode = '%s'`, u.EmployeeCode)
		errorQuery := db.QueryRow(stmt).Scan(&countUser)

		if errorQuery != nil {
			return false, ""

		}

		if countUser == 0 {

			rowsError += fmt.Sprintf(` User row : %d isn't found!`, index+2)
			html += fmt.Sprintf(`<tr>
        <td style="text-align:center;padding : 3px;">%d</td>
        <td style="text-align:center;padding : 3px;color:red;">%s</td>
      </tr>`, index+2, u.EmployeeCode)
		}
	}

	html += ` </tbody>
  </table>`

	if rowsError != "" {
		fmt.Println("rowsError", html)
		return true, html
	} else {
		fmt.Println("not Error", "")
		return false, ""
	}

}

func InsertManpower(employeeCode string, date string, actionBy string) bool {
	connString := config.LoadDatabaseConfig()

	// เปิดการเชื่อมต่อฐานข้อมูล
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error Creating connection: " + err.Error())
		return false
	}

	stmt := fmt.Sprintf(`INSERT INTO [dbo].[TBL_MANPOWER] ([DATE],[UHR_EmpCode],[UHR_OrgCode],[UHR_FullName_th],[UHR_FullName_en],[UHR_Department],[UHR_GroupDepartment],[UHR_POSITION]
           ,[UHR_GMail]
           ,[UHR_Sex]
           ,[UHR_StatusToUse]
           ,[UHR_OrgGroup]
           ,[UHR_OrgName]
           ,[CREATED_AT]
           ,[CREATED_BY])
     SELECT '%s'
           ,[UHR_EmpCode],[UHR_OrgCode],[UHR_FullName_th],[UHR_FullName_en],[UHR_Department],[UHR_GroupDepartment],[UHR_POSITION],[UHR_GMail],[UHR_Sex],[UHR_StatusToUse],[UHR_OrgGroup]
           ,[UHR_OrgName]
           ,GETDATE()
           ,'%s' FROM V_Users WHERE UHR_EmpCode = '%s'`, date, actionBy, employeeCode)

	_, errInsert := db.Exec(stmt)

	defer db.Close()
	return errInsert == nil

}

func JoinViewUpdateManpower(employeeCode string, date string, actionBy string) bool {
	connString := config.LoadDatabaseConfig()

	// เปิดการเชื่อมต่อฐานข้อมูล
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error Creating connection: " + err.Error())
		return false
	}

	defer db.Close()

	stmt := fmt.Sprintf(`
	UPDATE a
	SET a.[UHR_EmpCode] = b.UHR_EmpCode
		  ,a.[UHR_OrgCode] = b.[UHR_OrgCode]
		  ,a.[UHR_FullName_th] = b.[UHR_FullName_th]
		  ,a.[UHR_FullName_en] = b.[UHR_FullName_en] 
		  ,a.[UHR_Department] = b.[UHR_Department]
		  ,a.[UHR_GroupDepartment] = b.[UHR_GroupDepartment]
		  ,a.[UHR_POSITION] = b.[UHR_POSITION]
		  ,a.[UHR_GMail] = b.[UHR_GMail]
		  ,a.[UHR_Sex] = b.[UHR_Sex]
		  ,a.[UHR_StatusToUse] = b.[UHR_StatusToUse]
		  ,a.[UHR_OrgGroup] = b.[UHR_OrgGroup]
		  ,a.[UHR_OrgName] = b.[UHR_OrgName]
		  ,a.[UPDATED_BY] = '%s'
		  ,a.UPDATED_AT = GETDATE()
FROM [DB_MANPOWER_MGT].[dbo].[TBL_MANPOWER] a
LEFT JOIN V_Users b ON a.[UHR_EmpCode] COLLATE  Thai_CI_AS = b.[UHR_EmpCode] COLLATE Thai_CI_AS
WHERE  a.[DATE] = '%s' AND a.UHR_EmpCode = '%s'`, actionBy, date, employeeCode)

	_, errUpdate := db.Query(stmt)
	return errUpdate == nil

}

type UserDuplicate struct {
	UHR_EmpCode string
}

func CheckUserExists(db *sql.DB, date, empCode string) (bool, error) {
	stmt := fmt.Sprintf(`SELECT 1 FROM [DB_MANPOWER_MGT].[dbo].[TBL_MANPOWER] 
		WHERE DATE = '%s' AND UHR_EmpCode = '%s'`, date, empCode)
	results, err := db.Query(stmt)
	if err != nil {
		return false, err
	}
	defer results.Close()

	return results.Next(), nil
}

func processUsers(req model.BodyAddManpower, db *sql.DB) (int, int) {
	inserted, updated := 0, 0

	for _, u := range req.Users {

		exists, err := CheckUserExists(db, req.ManDate, u.EmployeeCode)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if exists {
			// Update

			if JoinViewUpdateManpower(u.EmployeeCode, req.ManDate, req.ActionBy) {
				updated++
			}
		} else {
			// Insert

			if InsertManpower(u.EmployeeCode, req.ManDate, req.ActionBy) {
				inserted++
			}
		}
	}

	return inserted, updated
}

func CheckUserDuplicated(req model.BodyAddManpower) (bool, string) {
	// โหลดค่าสำหรับเชื่อมต่อฐานข้อมูล
	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("SQL Connect error : " + err.Error())
		return false, ""
	}
	defer db.Close()

	inserted, updated := processUsers(req, db)

	return true, fmt.Sprintf("Inserted: %d, Updated: %d", inserted, updated)
}

func AddManpower(c *fiber.Ctx) error {
	var req model.BodyAddManpower

	// แปลงข้อมูล JSON ที่รับมา
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Invalid request body",
		})
	}

	// ตรวจสอบ JWT Token
	_, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid user data",
		})
	}

	// โหลดค่าสำหรับเชื่อมต่อฐานข้อมูล
	connString := config.LoadDatabaseConfig()

	// เปิดการเชื่อมต่อฐานข้อมูล
	db, err := sql.Open("sqlserver", connString)

	if err != nil {
		fmt.Println("Error create connection: " + err.Error())
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Database connection failed",
		})
	}

	defer db.Close()

	// Check User isn't found
	_, html := CheckUserOk(req)

	//Return Html to Client side
	if html != "" {
		return c.JSON(fiber.Map{
			"err":  true,
			"html": html,
			"msg":  "User not found",
		})
	} else {

		//Insert Update and Verify
		addData, _ := CheckUserDuplicated(req)

		if addData {

			return c.JSON(fiber.Map{"err": false, "msg": "Inserted!", "status": "Ok"})
		} else {
			return c.JSON(fiber.Map{"err": true, "msg": "Something went wrong!"})

		}
	}

}

func GetFliterManpowerMaster(c *fiber.Ctx) error {

	var users []model.ManpowerMaster
	department := c.Params("department")
	startDate := c.Params("startDate")
	endDate := c.Params("endDate")
	orgName := c.Params("orgName")
	orgGroup := c.Params("orgGroup")
	departmentDecoded, _ := url.QueryUnescape(department)

	_, ok := c.Locals("user").(jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid user data",
		})
	}

	stmt := fmt.Sprintf(`
	SELECT [Id],[DATE],[UHR_EmpCode],[UHR_OrgCode],[UHR_FullName_th],[UHR_FullName_en],[UHR_Department],[UHR_GroupDepartment]
      ,[UHR_POSITION]
      ,[UHR_GMail]
      ,[UHR_Sex]
      ,[UHR_StatusToUse]
      ,[UHR_OrgGroup]
      ,[UHR_OrgName]
      ,[CREATED_AT]
      ,[CREATED_BY]
      ,[UPDATED_AT]
      ,[UPDATED_BY]
  FROM TBL_MANPOWER WHERE [DATE] BETWEEN '%s' AND '%s'`, startDate, endDate)

	if department != "All" {
		stmt += fmt.Sprintf(` AND [UHR_Department] = '%s'`, departmentDecoded)
	}

	if orgName != "All" {
		stmt += fmt.Sprintf(` AND [UHR_OrgName]= '%s'`, orgName)
	}

	if orgGroup != "All" {
		stmt += fmt.Sprintf(` AND [UHR_OrgGroup]= '%s'`, orgGroup)
	}

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error creating Connection: " + err.Error())
	}

	defer db.Close()

	results, errQuery := db.Query(
		stmt,
		sql.Named("start", startDate),
		sql.Named("end", endDate),
		sql.Named("department", department),
		sql.Named("orgName", orgName),
		sql.Named("orgGroup", orgGroup),
	)

	if errQuery != nil {

		fmt.Println(errQuery.Error())
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for results.Next() {
		var user model.ManpowerMaster

		errScan := results.Scan(
			&user.Id,
			&user.DATE,
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
			&user.CREATED_AT,
			&user.CREATED_BY,
			&user.UPDATED_AT,
			&user.UPDATED_BY,
		)

		if errScan != nil {
			fmt.Println("Error : ", errScan.Error())
		} else {
			users = append(users, user)
		}
	}

	defer results.Close()

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
			"msg":     "User not found",
			"status":  "",
			"results": users,
		})
	}

}

func GetManpowerByGroupPositionAndDate(c *fiber.Ctx) error {

	//Loading connection string
	connString := config.LoadDatabaseConfig()

	var users []model.UserGroupPosition

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	start := c.Params("start")
	end := c.Params("end")

	rows, errQuery := db.Query(`
		SELECT DATE,PHR_PGroupCode,UHR_EmpCode,UHR_FullName_en,UHR_FullName_th,UHR_Department,UHR_POSITION FROM [dbo].V_ManPowerGroupPosition
WHERE [DATE] BETWEEN @start AND @end
ORDER BY [DATE],PHR_PGroupCode DESC
	`, sql.Named("start", start), sql.Named("end", end))

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		var user model.UserGroupPosition

		errScan := rows.Scan(
			&user.DATE,
			&user.PHR_PGroupCode,
			&user.UHR_EmpCode,
			&user.UHR_FullName_en,
			&user.UHR_FullName_th,
			&user.UHR_Department,
			&user.UHR_POSITION,
		)

		if errScan != nil {
			return c.JSON(fiber.Map{"err": true, "msg": errScan.Error()})

		} else {
			// Appended value to departments
			users = append(users, user)
		}
	}

	defer rows.Close()

	if len(users) > 0 {

		return c.JSON(fiber.Map{"err": false, "msg": "", "results": users, "status": "Ok"})
	} else {
		return c.JSON(fiber.Map{"err": true, "msg": "Departments isn't found.", "results": users})

	}

}

func DeleteManpowerById(c *fiber.Ctx) error {

	id := c.Params("id")
	user := c.Params("action")

	//Loading connection string
	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	var manId int
	errRow := db.QueryRow("SELECT [Id] FROM [dbo].[TBL_MANPOWER] WHERE [Id] = @id",
		sql.Named("id", id)).Scan(&manId)

	if errRow != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errRow.Error()})
	}

	stmt := `EXEC [dbo].[sProc_DeleteManpowerAndBackup] @Id = @Id, @action = @Action`

	_, err = db.Exec(stmt, sql.Named("Id", id), sql.Named("Action", user))

	if err != nil {
		return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
	}
	return c.JSON(fiber.Map{"err": false, "msg": "Deleted !", "status": "Ok"})
}

func CloningManpowerByDate(c *fiber.Ctx) error {

	source := c.Params("source")
	start := c.Params("start")
	end := c.Params("end")
	duration := c.Params("duration")

	var req model.BodyCloneManpower

	// แปลงข้อมูล JSON ที่รับมา
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Invalid request body",
		})
	}

	index := 0

	//Loading connection string
	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error Connection: " + err.Error())
	}

	defer db.Close()
	i, err := strconv.Atoi(duration)

	if err != nil {

		return c.JSON(fiber.Map{"err": true, "msg": "Convert String to Int error."})
	}
	var duplicated int
	// Loop
	for index < i {
		var date string

		errQuery := db.QueryRow(`SELECT TOP 1 [DATE] FROM [dbo].[TBL_MANPOWER] WHERE [DATE] IN (
		SELECT CONVERT(DATE,DATEADD(DAY,@index,@start)))`,
			sql.Named("index", index),
			sql.Named("start", start)).Scan(&date)

		if errQuery != nil {
			fmt.Println(errQuery.Error())

		}
		if date != "" {
			duplicated++
		}

		index++

	}
	defer db.Close()

	if duplicated > 0 {
		return c.JSON(fiber.Map{"err": true, "msg": "Manpower is duplicated!"})
	} else {

		stmt := fmt.Sprintf(`
		DECLARE	@return_value int,
		@resultMessage nvarchar(100)

		SELECT	@resultMessage = N'''%s'''

		EXEC	@return_value = [dbo].[sProc_CloneManpowerDataByDate]
				@souceDate = '%s',
				@start = '%s',
				@end = '%s',
				@actionBy = '%s',
				@resultMessage = @resultMessage OUTPUT

		SELECT	@resultMessage as N'@resultMessage'

		SELECT	'Return Value' = @return_value
		`, "Success", source, start, end, req.ActionBy)

		_, err = db.Exec(stmt)

		if err != nil {
			return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
		}

		return c.JSON(fiber.Map{"err": false, "msg": "Manpower uploaded!", "status": "Ok"})
	}

}
