package income_report

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/contbank/bankly-sdk/pkg/errors"
	"github.com/contbank/bankly-sdk/pkg/models"
	"github.com/contbank/bankly-sdk/pkg/utils"
	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// IncomeReport ...
type IncomeReport struct {
	httpClient utils.BanklyHttpClient
}

// NewIncomeReport ...
func NewIncomeReport(newHttpClient utils.BanklyHttpClient) *IncomeReport {
	newHttpClient.ErrorHandler = IncomeReportErrorHandler
	return &IncomeReport{newHttpClient}
}

// GetIncomeReport ...
func (c *IncomeReport) GetIncomeReport(ctx context.Context,
	model *models.IncomeReportRequest) (*models.IncomeReportResponse, error) {

	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
		"object" : model,
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

	var incomeReportResponse *models.IncomeReportResponse
	err = json.Unmarshal(respBody, &incomeReportResponse)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, errors.ErrDefaultIncomeReport
	}

	defer resp.Body.Close()
	return incomeReportResponse, nil
}

// IncomeReportErrorHandler ...
func IncomeReportErrorHandler(fields logrus.Fields, resp *http.Response) error {
	var bodyErr *errors.ErrorResponse

	respBody, _ := ioutil.ReadAll(resp.Body)

	err := json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return errors.ErrDefaultCard
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err := errors.FindIncomeReportError(errModel.Code, errModel.Messages...)
		fields["bankly_error"] = bodyErr
		logrus.WithFields(fields).WithError(err).Error("bankly get income report error")
		return err
	}

	return errors.ErrDefaultIncomeReport
}