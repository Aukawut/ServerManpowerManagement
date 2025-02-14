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
