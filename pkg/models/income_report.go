package models

// IncomeReportRequest ...
type IncomeReportRequest struct {
	Account string `validate:"required" form:"account"`
	Year    string `validate:"required" form:"year"`
}

// IncomeReportResponse ...
type IncomeReportResponse struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	IncomeFile  string `json:"incomeFile"`
}
