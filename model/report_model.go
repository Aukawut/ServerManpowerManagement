package model

type ReportHeadCountByDeptAndSex struct {
	DEPARTMENT string
	FEMALE     int
	MALE       int
	OTHER      int
}

type ReportHeadCountByDept struct {
	DEPARTMENT string
	HEAD_COUNT int
}

type ReportHeadCountByPosition struct {
	POSITION   string
	HEAD_COUNT int
}

type ReportHeadCountByDeptAndUsersType struct {
	DEPARTMENT string
	INDIRECT   int
	DIRECT     int
	SGA        int
	OTHER      int
	TOTAL      int
}

type ReportHeadCountBySex struct {
	UHR_Sex string
	AMOUNT  int
}
