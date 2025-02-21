package model

type Department struct {
	DHR_DName string
	DHR_DDes  string
}

type UserGroupPosition struct {
	DATE            string
	PHR_PGroupCode  interface{}
	UHR_EmpCode     interface{}
	UHR_FullName_en string
	UHR_FullName_th string
	UHR_Department  string
	UHR_POSITION    string
}

type DepartmentOfUsers struct {
	UHR_Department string
}
