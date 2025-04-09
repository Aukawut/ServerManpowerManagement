package model

type Month struct {
	MonthNo   int
	MonthName string
}

type ExcelCostAndheaCount struct {
	Category string
	Cell     string
	CellNo   int
}

type ExcelCostValue struct {
	Category    string
	LABOR       any
	ROFO_PEOPLE any
	OB_PEOPLE   any
	ACTUAL_COST any
	ROFO_COST   any
	OB_COST     any
	MONTH       string
	YEAR        string
	UUID        string
}

type HeadCountFile struct {
	Id         int
	FILE_NAME  string
	MONTH_NAME string
	MONTH      int
	YEAR       int
	UUID       string
	CREATED_AT string
	UPDATED_AT any
}

type SyncFileJsonBody struct {
	ActionBy string `json:"actionBy"`
}

type HeadCountAndCostReport struct {
	ClassName          string
	ACTUAL_PERSON      int
	ACTUAL_COST        float64
	ROFO_PEOPLE        int
	ROFO_COST          float64
	OB_PEOPLE          int
	OB_COST            float64
	ACTUAL_PEOPLE_PREV int
	ACTUAL_COST_PREV   float64
	COST_ACTUAL_ROFO   float64
	COST_ACTUAL_OB     float64
	COST_ACTUAL_PREV   float64
	PEOPLE_ACTUAL_OB   int
	PEOPLE_AC_PREV     int
	ClassNo            string
	LABOR              string
}

type YearOptions struct {
	Year int `json:"year"`
}

type MonthOptions struct {
	Month int `json:"month"`
}
