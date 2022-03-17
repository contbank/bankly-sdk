package bankly

import (
	"github.com/contbank/grok"
	"net/http"
	"strings"
)

// errorIncomeReportList ...
var errorIncomeReportList = []Error{
	{
		ErrorKey:  "INVALID_CALENDAR_FOR_INCOME_REPORT",
		GrokError: ErrInvalidIncomeReportCalendar,
	},
	{
		ErrorKey:  "INVALID_PARAMETER_INCOME_REPORT",
		GrokError: ErrInvalidIncomeReportParameter,
	},
}

// FindIncomeReportError Find income report errors.
func FindIncomeReportError(code string, messages ...string) *grok.Error {
	code = verifyInvalidIncomeReportParameter(code, messages)

	for _, v := range errorCardList {
		if v.ErrorKey == code {
			return v.GrokError
		}
	}

	return grok.NewError(http.StatusConflict, messages...)
}

// verifyInvalidIncomeReportParameter Find the correspondent error message for income reports.
func verifyInvalidIncomeReportParameter(code string, messages []string) string {
	if code == "CALENDAR_NOT_ALLOWED" {
		for _, m := range messages {
			switch {
			case strings.Contains(strings.ToLower(m), "calendar informed is not allowed"):
				return "INVALID_CALENDAR_FOR_INCOME_REPORT"
			default:
				return "INVALID_PARAMETER_INCOME_REPORT"
			}
		}
	}
	return code
}