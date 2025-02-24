package model

type ManpowerEmployeeCode struct {
	EmployeeCode string `json:"employeeCode"`
}

type BodyAddManpower struct {
	Users    []ManpowerEmployeeCode `json:"users"`
	ActionBy string                 `json:"actionBy"`
	ManDate  string                 `json:"manDate"`
}

type BodyCloneManpower struct {
	ActionBy string `json:"actionBy"`
}

type SummaryManTermination struct {
	UHR_Department string
	PERSON         int
}

type ManpowerMaster struct {
	Id                  int
	DATE                string
	UHR_EmpCode         string
	UHR_OrgCode         string
	UHR_FullName_th     string
	UHR_FullName_en     string
	UHR_Department      string
	UHR_GroupDepartment string
	UHR_POSITION        interface{}
	UHR_GMail           interface{}
	UHR_Sex             string
	UHR_StatusToUse     string
	UHR_OrgGroup        string
	UHR_OrgName         string
	CREATED_AT          interface{}
	CREATED_BY          interface{}
	UPDATED_AT          interface{}
	UPDATED_BY          interface{}
}
