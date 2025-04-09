package handlers

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"strings"

	"database/sql"

	"github.com/Aukawut/ServerManpowerManagement/config"
	"github.com/Aukawut/ServerManpowerManagement/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

func GetHeadCountFile(c *fiber.Ctx) error {
	connString := config.LoadDatabaseConfig()

	month := c.Params("month")
	year := c.Params("year")

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error Open Connection: " + err.Error())
	}

	defer db.Close()

	files := []model.HeadCountFile{}

	rows, err := db.Query(`SELECT 
    [Id],
    FILE_NAME,
    CASE 
        WHEN MONTH = 1 THEN 'Jan'
        WHEN MONTH = 2 THEN 'Feb'
        WHEN MONTH = 3 THEN 'Mar'
        WHEN MONTH = 4 THEN 'Apr'
        WHEN MONTH = 5 THEN 'May'
        WHEN MONTH = 6 THEN 'Jun'
        WHEN MONTH = 7 THEN 'Jul'
        WHEN MONTH = 8 THEN 'Aug'
        WHEN MONTH = 9 THEN 'Sep'
        WHEN MONTH = 10 THEN 'Oct'
        WHEN MONTH = 11 THEN 'Nov'
        WHEN MONTH = 12 THEN 'Dec'
    END AS [MONTH_NAME],
    MONTH,
    YEAR,
    UUID,
    CREATED_AT,
    UPDATED_AT 
FROM 
    TBL_HEAD_FILES 
WHERE 
    ACTIVE = 'Y' AND [MONTH] = @month AND [YEAR] = @year`, sql.Named("month", month), sql.Named("year", year))
	if err != nil {
		return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
	}
	if rows.Next() {
		// # Model
		file := model.HeadCountFile{}

		// # Scan Value to struct
		err := rows.Scan(
			&file.Id,
			&file.FILE_NAME,
			&file.MONTH_NAME,
			&file.MONTH,
			&file.YEAR,
			&file.UUID,
			&file.CREATED_AT,
			&file.UPDATED_AT,
		)

		if err != nil {
			return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
		} else {
			files = append(files, file)
		}
	}

	if len(files) > 0 {
		return c.JSON(fiber.Map{"err": false, "msg": "", "status": "Ok", "results": files})
	} else {
		return c.JSON(fiber.Map{"err": true, "msg": ""})
	}

}

func CheckEmptyString(str string) string {

	if str == "" {
		return "0"
	} else {
		trim := strings.TrimSpace(str)
		if strings.Contains(trim, ",") {
			trim = strings.ReplaceAll(trim, ",", "")
		}
		return trim
	}
}

func CheckSheetName(sheetName string, sheets []string) bool {
	for _, item := range sheets {
		if sheetName == item {
			return true
		}
	}
	return false
}

func ReadExcelheadCount(sheetName string, uuid string, month string, year string, fileName string) ([]model.ExcelCostValue, error) {

	excelValuePr := []model.ExcelCostValue{}
	excelValueTemp := []model.ExcelCostValue{}

	path := fmt.Sprintf(`D:\2. Development Application\1. Prospira Company\23. Man Power Management\Document\HC & People cost 2025\%s`, fileName)
	headCountExcel := []model.ExcelCostAndheaCount{
		{Category: "#1", Cell: "B", CellNo: 7},
		{Category: "#2", Cell: "B", CellNo: 10},
		{Category: "#3", Cell: "B", CellNo: 13},
		{Category: "#4", Cell: "B", CellNo: 16},
		{Category: "#5", Cell: "B", CellNo: 19},
		{Category: "#6", Cell: "B", CellNo: 22},
		{Category: "#7", Cell: "B", CellNo: 25},
	}

	// Open Excel
	f, errEx := excelize.OpenFile(path)

	if errEx != nil {
		fmt.Println("Failed Open Excel file:" + errEx.Error())
		return excelValuePr, errEx
	}

	defer f.Close()

	sheetList := f.GetSheetList()

	// Check SheetName
	sheetOk := CheckSheetName(sheetName, sheetList)

	if !sheetOk {
		return excelValuePr, nil
	}

	for index := range headCountExcel {

		valuePermanent := model.ExcelCostValue{}
		valueTemp := model.ExcelCostValue{}

		// <-- Actual -->
		actualCostLocationPr := fmt.Sprintf(`D%d`, headCountExcel[index].CellNo+1)
		actualCostLocationTp := fmt.Sprintf(`D%d`, headCountExcel[index].CellNo+2)

		// <--- ROFO -->
		rofoPeoplePr := fmt.Sprintf(`F%d`, headCountExcel[index].CellNo+1)
		rofoPeopleTp := fmt.Sprintf(`F%d`, headCountExcel[index].CellNo+2)
		rofoCostPr := fmt.Sprintf(`G%d`, headCountExcel[index].CellNo+1)
		rofoCostTp := fmt.Sprintf(`G%d`, headCountExcel[index].CellNo+2)

		// <-- OB -->
		obPeoplePr := fmt.Sprintf(`O%d`, headCountExcel[index].CellNo+1)
		obPeopleTp := fmt.Sprintf(`O%d`, headCountExcel[index].CellNo+2)
		obCostPr := fmt.Sprintf(`P%d`, headCountExcel[index].CellNo+1)
		obCostTp := fmt.Sprintf(`P%d`, headCountExcel[index].CellNo+2)

		// Permanent
		valueRofoPeoplePr, _ := f.GetCellValue(sheetName, rofoPeoplePr)
		valueObPeoplePr, _ := f.GetCellValue(sheetName, obPeoplePr)
		valueActualCostPr, _ := f.GetCellValue(sheetName, actualCostLocationPr)
		valueRofoCostPr, _ := f.GetCellValue(sheetName, rofoCostPr)
		valueObCostPr, _ := f.GetCellValue(sheetName, obCostPr)

		// Temporary
		valueRofoPeopleTemp, _ := f.GetCellValue(sheetName, rofoPeopleTp)
		valueObPeopleTemp, _ := f.GetCellValue(sheetName, obPeopleTp)
		valueActualCostTemp, _ := f.GetCellValue(sheetName, actualCostLocationTp)
		valueRofoCostTemp, _ := f.GetCellValue(sheetName, rofoCostTp)
		valueObCostTemp, _ := f.GetCellValue(sheetName, obCostTp)

		// <---- Permanent Append --->
		valuePermanent.Category = headCountExcel[index].Category
		valuePermanent.LABOR = "Permanent"
		valuePermanent.ROFO_PEOPLE = CheckEmptyString(valueRofoPeoplePr)
		valuePermanent.OB_PEOPLE = CheckEmptyString(valueObPeoplePr)
		valuePermanent.ACTUAL_COST = CheckEmptyString(valueActualCostPr)
		valuePermanent.ROFO_COST = CheckEmptyString(valueRofoCostPr)
		valuePermanent.OB_COST = CheckEmptyString(valueObCostPr)
		valuePermanent.MONTH = month
		valuePermanent.YEAR = year
		valuePermanent.UUID = uuid

		// <---- Permanent Append --->
		valueTemp.Category = headCountExcel[index].Category
		valueTemp.LABOR = "Temporary"
		valueTemp.ROFO_PEOPLE = CheckEmptyString(valueRofoPeopleTemp)
		valueTemp.OB_PEOPLE = CheckEmptyString(valueObPeopleTemp)
		valueTemp.ACTUAL_COST = CheckEmptyString(valueActualCostTemp)
		valueTemp.ROFO_COST = CheckEmptyString(valueRofoCostTemp)
		valueTemp.OB_COST = CheckEmptyString(valueObCostTemp)
		valueTemp.MONTH = month
		valueTemp.YEAR = year
		valueTemp.UUID = uuid

		excelValuePr = append(excelValuePr, valuePermanent)
		excelValueTemp = append(excelValueTemp, valueTemp)

	}

	if len(excelValuePr) > 0 && len(excelValueTemp) > 0 {
		combined := append(excelValuePr, excelValueTemp...)
		return combined, nil
	} else {
		return excelValuePr, nil
	}

}

func ConvertSheetNameFormat(year string, month string) (string, error) {

	months := []model.Month{
		{MonthNo: 1, MonthName: "Jan"},
		{MonthNo: 2, MonthName: "Feb"},
		{MonthNo: 3, MonthName: "Mar"},
		{MonthNo: 4, MonthName: "Apr"},
		{MonthNo: 5, MonthName: "May"},
		{MonthNo: 6, MonthName: "Jun"},
		{MonthNo: 7, MonthName: "Jul"},
		{MonthNo: 8, MonthName: "Aug"},
		{MonthNo: 9, MonthName: "Sep"},
		{MonthNo: 10, MonthName: "Oct"},
		{MonthNo: 11, MonthName: "Nov"},
		{MonthNo: 12, MonthName: "Dec"},
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		return "", fmt.Errorf("invalid month: %s", month)
	}
	y := year

	if len(y) == 4 {
		y = year[2:]
	}

	monthName := months[monthInt-1].MonthName
	return monthName + y, nil
}

func DuplicatedSync(month string, year string) (bool, string) {
	uuid := ``
	connString := config.LoadDatabaseConfig()
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error Open Connect : " + err.Error())
	}

	defer db.Close()

	errQuery := db.QueryRow(fmt.Sprintf(`SELECT [UUID] FROM TBL_HEAD_FILES WHERE ACTIVE = 'Y' 
AND [MONTH] = '%s' AND [YEAR] = '%s'`, month, year)).Scan(&uuid)
	if errQuery != nil {
		uuid = ""
		return false, ""
	}

	if uuid != "" {
		return true, uuid

	} else {
		return false, ""
	}

}

func ConvertHeadCountFileName(year string, month string) string {
	months := []model.Month{
		{MonthNo: 1, MonthName: "Jan"},
		{MonthNo: 2, MonthName: "Feb"},
		{MonthNo: 3, MonthName: "Mar"},
		{MonthNo: 4, MonthName: "Apr"},
		{MonthNo: 5, MonthName: "May"},
		{MonthNo: 6, MonthName: "Jun"},
		{MonthNo: 7, MonthName: "Jul"},
		{MonthNo: 8, MonthName: "Aug"},
		{MonthNo: 9, MonthName: "Sep"},
		{MonthNo: 10, MonthName: "Oct"},
		{MonthNo: 11, MonthName: "Nov"},
		{MonthNo: 12, MonthName: "Dec"},
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		return ""
	}

	monthName := months[monthInt-1].MonthName

	fileName := fmt.Sprintf(`HC %s %s.xlsx`, monthName, year)

	return fileName

}

func VerifyFileSync(fileNameCost string, year string, month string, sheetName string) bool {
	// Head count file

	fileName := ConvertHeadCountFileName(year, month)

	pathCost := fmt.Sprintf(`D:\2. Development Application\1. Prospira Company\23. Man Power Management\Document\HC & People cost 2025\%s`, fileNameCost)
	pathHeadCountRaw := fmt.Sprintf(`D:\2. Development Application\1. Prospira Company\23. Man Power Management\Document\HC People 2025 (Monthly)\%s`, fileName)

	// Check if pathHeadCountRaw exists
	if _, err := os.Stat(pathHeadCountRaw); os.IsNotExist(err) {
		fmt.Printf("File does not exist: %s\n", pathHeadCountRaw)
		return false
	} else {
		// Check Sheet Name File Cost
		// # Open Excel
		f, errEx := excelize.OpenFile(pathCost)

		if errEx != nil {
			fmt.Println("Failed to open Excel file:" + errEx.Error())

		}

		defer f.Close()

		sheetList := f.GetSheetList()

		// Check SheetName
		sheetOk := CheckSheetName(sheetName, sheetList)

		if sheetOk {
			return true
		}

	}
	return false
}

func UnActiveHeadCountFile(uuid string, db *sql.DB) bool {
	_, errUpdate := db.Query(`UPDATE [dbo].[TBL_HEAD_FILES] SET [ACTIVE] = 'N' WHERE [UUID] = @uuid`, sql.Named("uuid", uuid))
	_, errUpdateRaw := db.Query(`UPDATE [dbo].[TBL_HEADCOUNT] SET [ACTIVE] = 'N' WHERE [UUID] = @uuid`, sql.Named("uuid", uuid))

	if errUpdate != nil || errUpdateRaw != nil {
		return false
	}
	defer db.Close()
	return true

}

func UpdateCostHeadCount(oldUuid string, data []model.ExcelCostValue, id string, fileName string, year string, month string, actionBy string) (bool, error) {
	connString := config.LoadDatabaseConfig()
	db, err := sql.Open("sqlserver", connString)

	createAt := ``

	errScan := db.QueryRow(`SELECT CREATED_AT FROM  [dbo].[TBL_HEAD_FILES] WHERE [UUID] = @uuid`, sql.Named("uuid", oldUuid)).Scan(&createAt)

	if err != nil {
		fmt.Println("Error Open Connect: " + err.Error())
	}

	if errScan != nil {
		createAt = ""
	}

	fmt.Println("createAt", createAt)

	_, errUpdate := db.Query(`UPDATE [dbo].[TBL_HEAD_FILES] SET [ACTIVE] = 'N' WHERE [UUID] = @uuid`, sql.Named("uuid", oldUuid))

	if errUpdate != nil {

		return false, errUpdate
	}

	//insert
	_, errInsert := db.Query(`INSERT INTO [TBL_HEAD_FILES] ([FILE_NAME],[ACTIVE],[MONTH],[YEAR],[UUID],[CREATED_BY],[CREATED_AT],[UPDATED_AT]) 
	VALUES (@fileName,'Y',@month,@year,@uuid,@actionBy,@createdAt,GETDATE())`,
		sql.Named("fileName", fileName),
		sql.Named("month", month),
		sql.Named("year", year),
		sql.Named("uuid", id),
		sql.Named("createdAt", createAt),
		sql.Named("actionBy", actionBy),
	)

	if errInsert != nil {
		return false, errInsert
	}

	for _, item := range data {

		_, err := db.Query(`INSERT INTO [dbo].[TBL_PEOPLE_COST] ([CATEGORY_CLASS]
      ,[LABOR]
      ,[ROFO_PEOPLE]
      ,[OB_PEOPLE]
      ,[ACTUAL_COST]
      ,[ROFO_COST]
      ,[OB_COST]
      ,[MONTH]
      ,[YEAR]
      ,[CREATED_AT]
      ,[CREATED_BY]
      ,[UUID]) VALUES (@categoryClass,@labor,@rofoPeople,@obPeople,@actualCost,@rofoCost,@obCost,@month,@year,@createDate,@actionBy,@uuid)`,
			sql.Named("categoryClass", item.Category),
			sql.Named("labor", item.LABOR),
			sql.Named("rofoPeople", item.ROFO_PEOPLE),
			sql.Named("obPeople", item.OB_PEOPLE),
			sql.Named("actualCost", item.ACTUAL_COST),
			sql.Named("rofoCost", item.ROFO_COST),
			sql.Named("obCost", item.OB_COST),
			sql.Named("month", item.MONTH),
			sql.Named("year", item.YEAR),
			sql.Named("uuid", id),
			sql.Named("createDate", createAt),
			sql.Named("actionBy", actionBy),
		)

		if err != nil {
			fmt.Println(err.Error())
			return false, err
		}

	}

	_, errInsertRaw := InsertRawDataHeadCount(id, fileName, actionBy, db)

	if errInsertRaw != nil {
		return false, errInsertRaw
	}

	defer db.Close()

	return true, nil
}

func InsertRawDataHeadCount(uuid string, fileName string, actionBy string, db *sql.DB) (bool, error) {
	path := fmt.Sprintf(`D:\2. Development Application\1. Prospira Company\23. Man Power Management\Document\HC People 2025 (Monthly)\%s`, fileName)

	// Open Excel
	f, errEx := excelize.OpenFile(path)
	if errEx != nil {
		fmt.Println("Failed to open Excel file:", errEx)
		return false, errEx
	}
	defer f.Close()

	sheetName := `database Feb`

	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Println("Failed to get rows:", err)
		return false, err
	}

	for index := 2; index <= len(rows); index++ {
		fmt.Println(index)
		employeeCode, _ := f.GetCellValue(sheetName, fmt.Sprintf(`B%d`, index))
		thaiName, _ := f.GetCellValue(sheetName, fmt.Sprintf(`C%d`, index))
		engName, _ := f.GetCellValue(sheetName, fmt.Sprintf(`D%d`, index))
		orgCode, _ := f.GetCellValue(sheetName, fmt.Sprintf(`E%d`, index))
		orgName, _ := f.GetCellValue(sheetName, fmt.Sprintf(`F%d`, index))
		position, _ := f.GetCellValue(sheetName, fmt.Sprintf(`G%d`, index))
		costCenter, _ := f.GetCellValue(sheetName, fmt.Sprintf(`H%d`, index))
		startDateRaw, _ := f.GetCellValue(sheetName, fmt.Sprintf(`I%d`, index))
		costClass, _ := f.GetCellValue(sheetName, fmt.Sprintf(`J%d`, index))

		// Trim and validate date string
		startDateRaw = strings.TrimSpace(startDateRaw)
		var startDateFormat sql.NullString

		if startDateRaw != "" {
			t, err := time.Parse("2/1/2006", startDateRaw)
			if err != nil {

				startDateFormat.Valid = false
			} else {
				startDateFormat.String = t.Format("2006-01-02")
				startDateFormat.Valid = true
			}
		}

		_, errInsert := db.Exec(`INSERT INTO [dbo].[TBL_HEADCOUNT] 
			([FILE_NAME],[EMPLOYEE_ID],[THAI_NAME],[ENG_NAME],[ORG_CODE],[ORG_NAME],[POSITION],
			[COST_CENTER],[START_DATE],[COST],[UUID],[CREATED_AT],[CREATED_BY],[UPDATED_BY]) 
			VALUES (@fileName,@empId,@thaiName,@engName,@orgCode,@orgName,@position,@costCenter,
			@startDate,@cost,@uuid,GETDATE(),@actionBy,@actionBy)`,
			sql.Named("fileName", fileName),
			sql.Named("empId", employeeCode),
			sql.Named("thaiName", thaiName),
			sql.Named("engName", engName),
			sql.Named("orgCode", orgCode),
			sql.Named("orgName", orgName),
			sql.Named("position", position),
			sql.Named("costCenter", costCenter),
			sql.Named("startDate", startDateFormat), // handles NULLs correctly
			sql.Named("cost", costClass),
			sql.Named("uuid", uuid),
			sql.Named("actionBy", actionBy),
		)

		if errInsert != nil {
			fmt.Printf("Insert failed on row %d: %v\n", index, errInsert)
			return false, errInsert
		}
	}

	return true, nil
}

func InsertCostHeadCount(data []model.ExcelCostValue, id string, fileName string, year string, month string, actionBy string) (bool, error) {
	fmt.Println("data", len(data))
	connString := config.LoadDatabaseConfig()
	db, err := sql.Open("sqlserver", connString)

	if err != nil {
		fmt.Println("Error Open Connect: " + err.Error())
	}

	//Insert File

	_, errInsert := db.Query(`INSERT INTO [TBL_HEAD_FILES] ([FILE_NAME],[ACTIVE],[MONTH],[YEAR],[UUID],[CREATED_BY],[CREATED_AT]) 
	VALUES (@fileName,'Y',@month,@year,@uuid,@actionBy,GETDATE())`,
		sql.Named("fileName", fileName),
		sql.Named("month", month),
		sql.Named("year", year),
		sql.Named("uuid", id),
		sql.Named("actionBy", actionBy),
	)

	if errInsert != nil {
		return false, errInsert
	}

	for _, item := range data {
		_, err := db.Query(`INSERT INTO [dbo].[TBL_PEOPLE_COST] ([CATEGORY_CLASS]
      ,[LABOR]
      ,[ROFO_PEOPLE]
      ,[OB_PEOPLE]
      ,[ACTUAL_COST]
      ,[ROFO_COST]
      ,[OB_COST]
      ,[MONTH]
      ,[YEAR]
      ,[CREATED_AT]
      ,[CREATED_BY]
      ,[UUID]) VALUES (@categoryClass,@labor,@rofoPeople,@obPeople,@actualCost,@rofoCost,@obCost,@month,@year,GETDATE(),@actionBy,@uuid)`,
			sql.Named("categoryClass", item.Category),
			sql.Named("labor", item.LABOR),
			sql.Named("rofoPeople", item.ROFO_PEOPLE),
			sql.Named("obPeople", item.OB_PEOPLE),
			sql.Named("actualCost", item.ACTUAL_COST),
			sql.Named("rofoCost", item.ROFO_COST),
			sql.Named("obCost", item.OB_COST),
			sql.Named("month", item.MONTH),
			sql.Named("year", item.YEAR),
			sql.Named("uuid", id),
			sql.Named("actionBy", actionBy),
		)

		if err != nil {
			fmt.Println(err.Error())
			return false, err
		}

	}

	_, errInsertRaw := InsertRawDataHeadCount(id, fileName, actionBy, db)

	if errInsertRaw != nil {
		return false, errInsertRaw
	}

	defer db.Close()

	return true, nil
}

func SyncExcelFile(c *fiber.Ctx) error {
	month := c.Params("month")
	year := c.Params("year")

	var req model.SyncFileJsonBody

	// แปลงข้อมูล JSON ที่รับมา
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Invalid request body",
		})
	}

	// File Cost
	fileNameCost := `HC & People Cost Feb 2025.xlsx`

	// UUID
	id := uuid.New()

	// Duplicated Checking
	duplicated, uuidPrev := DuplicatedSync(month, year)

	// Convert Sheet Name Format
	sheetName, _ := ConvertSheetNameFormat(year, month)

	// Check File is Exits
	fileVerify := VerifyFileSync(fileNameCost, year, month, sheetName)

	// If Ok
	if fileVerify {

		// Read Excel Cost
		data, _ := ReadExcelheadCount(sheetName, id.String(), month, year, fileNameCost)

		fileName := ConvertHeadCountFileName(year, month)

		if duplicated {
			// Update Status active to 'N' and Insert New Record
			updated, errUpdate := UpdateCostHeadCount(uuidPrev, data, id.String(), fileName, year, month, req.ActionBy)

			fmt.Println("Update")

			if updated {
				return c.JSON(fiber.Map{"err": false, "msg": "Sync data successfully!", "status": "Ok"})
			} else {

				return c.JSON(fiber.Map{"err": true, "status": "Ok", "msg": errUpdate.Error()})
			}

		} else {
			// Insert
			fmt.Println("Insert")
			inserted, errInsert := InsertCostHeadCount(data, id.String(), fileName, year, month, req.ActionBy)

			if inserted {

				return c.JSON(fiber.Map{"err": false, "msg": "Sync data successfully!", "status": "Ok"})
			} else {
				return c.JSON(fiber.Map{"err": false, "msg": errInsert.Error(), "status": "Ok"})

			}
		}

	} else {
		// Return Error to Client side
		return c.JSON(fiber.Map{"err": true, "msg": "File is not verify!"})
	}
}

func GetHeadCountAndCostByUuid(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	results := []model.HeadCountAndCostReport{}

	stmt := fmt.Sprintf(`SELECT ClassName,ACTUAL_PERSON,ACTUAL_COST,ROFO_PEOPLE,ROFO_COST,OB_PEOPLE,OB_COST,ACTUAL_PEOPLE_PREV,
ACTUAL_COST_PREV,COST_ACTUAL_ROFO,COST_ACTUAL_OB,COST_ACTUAL_PREV,PEOPLE_ACTUAL_OB,PEOPLE_AC_PREV,
ClassNo,LABOR FROM func_GetPeopleCostAnalysisByUUID ('%s')`, uuid)

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)

	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	rows, errQuery := db.Query(stmt)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		result := model.HeadCountAndCostReport{}
		errScan := rows.Scan(
			&result.ClassName,
			&result.ACTUAL_PERSON,
			&result.ACTUAL_COST,
			&result.ROFO_PEOPLE,
			&result.ROFO_COST,
			&result.OB_PEOPLE,
			&result.OB_COST,
			&result.ACTUAL_PEOPLE_PREV,
			&result.ACTUAL_COST_PREV,
			&result.COST_ACTUAL_ROFO,
			&result.COST_ACTUAL_OB,
			&result.COST_ACTUAL_PREV,
			&result.PEOPLE_ACTUAL_OB,
			&result.PEOPLE_AC_PREV,
			&result.ClassNo,
			&result.LABOR,
		)

		if errScan != nil {
			return c.JSON(fiber.Map{"err": true, "msg": errScan.Error()})
		} else {
			results = append(results, result)
		}
	}

	defer db.Close()

	return c.JSON(fiber.Map{"err": false, "status": "Ok", "results": results})

}

func GetYearOptions(c *fiber.Ctx) error {
	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)

	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	results := []model.YearOptions{}

	rows, errQuery := db.Query(`SELECT [YEAR]
  FROM [DB_MANPOWER_MGT].[dbo].[TBL_HEAD_FILES] WHERE ACTIVE = 'Y' GROUP BY [YEAR]`)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		result := model.YearOptions{}
		err := rows.Scan(
			&result.Year,
		)

		if err != nil {
			return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
		} else {
			results = append(results, result)
		}

	}
	defer rows.Close()

	return c.JSON(fiber.Map{"err": false, "results": results, "status": "Ok"})
}

func GetMonthOptions(c *fiber.Ctx) error {
	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)

	if err != nil {
		fmt.Println("Error Create connection: " + err.Error())
	}

	defer db.Close()

	results := []model.MonthOptions{}

	rows, errQuery := db.Query(`SELECT [YEAR]
  FROM [DB_MANPOWER_MGT].[dbo].[TBL_HEAD_FILES] WHERE ACTIVE = 'Y' GROUP BY [YEAR]`)

	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		result := model.MonthOptions{}
		err := rows.Scan(
			&result.Month,
		)

		if err != nil {
			return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
		} else {
			results = append(results, result)
		}

	}
	defer rows.Close()

	return c.JSON(fiber.Map{"err": false, "results": results, "status": "Ok"})
}

func DeleteHeadCountData(c *fiber.Ctx) error {
	connString := config.LoadDatabaseConfig()

	uuid := c.Params("uuid")
	db, err := sql.Open("sqlserver", connString)

	if err != nil {
		fmt.Println("Error Create connection: " + err.Error())
	}

	defer db.Close()

	_, errDelete := db.Exec(`UPDATE [dbo].[TBL_HEAD_FILES] SET [ACTIVE] = 'N' WHERE [UUID] = @uuid`, sql.Named("uuid", uuid))

	if errDelete != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errDelete.Error()})
	}

	defer db.Close()

	return c.JSON(fiber.Map{"err": false, "status": "Ok", "msg": "Headcount deleted!"})
}
