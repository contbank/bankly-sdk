package bankly

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// IncomeReport ...
type IncomeReport struct {
	httpClient BanklyHttpClient
}

// NewIncomeReport ...
func NewIncomeReport(newHttpClient BanklyHttpClient) *IncomeReport {
	newHttpClient.SetErrorHandler(CardErrorHandler)
	return &IncomeReport{newHttpClient}
}

// GetIncomeReport ...
func (c *IncomeReport) GetIncomeReport(ctx context.Context, model *IncomeReportRequest) (*IncomeReportResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"object":     model,
	}

	// model validator
	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	// endpoint url
	url := fmt.Sprintf("/accounts/%s/income-report/print", model.Account)
	fields["url"] = url

	// query parameters
	query := make(map[string]string)
	query["calendar"] = grok.OnlyDigits(model.Year)

	// do request to bankly
	resp, err := c.httpClient.Get(ctx, url, query, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	// response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	var incomeReportResponse *IncomeReportResponse
	err = json.Unmarshal(respBody, &incomeReportResponse)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultIncomeReport
	}

	defer resp.Body.Close()
	return incomeReportResponse, nil
}

// IncomeReportErrorHandler ...
func IncomeReportErrorHandler(fields logrus.Fields, resp *http.Response) error {
	var bodyErr *ErrorResponse

	respBody, _ := ioutil.ReadAll(resp.Body)

	err := json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return ErrDefaultCard
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err := FindIncomeReportError(errModel.Code, errModel.Messages...)
		fields["bankly_error"] = bodyErr
		logrus.WithFields(fields).WithError(err).Error("bankly get income report error")
		return err
	}

	return ErrDefaultIncomeReport
}
