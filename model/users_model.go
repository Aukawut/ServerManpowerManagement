package model

type Users struct {
	EMPLOYEE_CODE string
	ROLE_ID       int
	CREATED_AT    interface{}
	UPDATED_AT    interface{}
	CREATED_BY    interface{}
	UPDATED_BY    interface{}
	ACTIVE        string
}

type UsersInfoLogin struct {
	UHR_EmpCode    string
	UHR_Position   string
	UHR_Department string
	UHR_FullName   string
	UHR_Sex        string
	RoleName       string
}

type UserJsonWebToken struct {
	Department    string
	Employee_code string
	Exp           interface{}
	Position      string
	Role          string
}

type UsersMaster struct {
	UHR_EmpCode         string
	UHR_OrgCode         string
	UHR_FullName_th     string
	UHR_FullName_en     string
	UHR_Department      string
	UHR_GroupDepartment string
	UHR_POSITION        string
	UHR_GMail           string
	UHR_Sex             string
	UHR_StatusToUse     string
	UHR_OrgGroup        string
	UHR_OrgName         string
}

type UsersAuthentication struct {
	EMPLOYEE_CODE   string
	UHR_FullName_th string
	UHR_POSITION    string
	UHR_OrgGroup    string
	UHR_OrgName     string
	UHR_OrgCode     string
	ROLE_ID         int
	ROLE_NAME       string
	ACTIVE          string
	UHR_Department  string
}

type ActiveUserBody struct {
	ActionBy string `json:"actionBy"`
	Active   string `json:"active"`
}

type InsertAuthenUserBody struct {
	EmployeeCode string `json:"employeeCode"`
	Active       string `json:"active"`
	Role         int    `json:"role"`
	ActionBy     string `json:"actionBy"`
}
