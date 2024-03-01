package bankly

import (
	"context"
	"encoding/json"
	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

// IncomeReport ...
type IncomeReport struct {
	session        Session
	authentication *Authentication
	httpClient     *http.Client
}

// NewIncomeReport ...
func NewIncomeReport(httpClient *http.Client, session Session) *IncomeReport {
	return &IncomeReport{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

// GetIncomeReport ...
func (c *IncomeReport) GetIncomeReport(ctx context.Context,
	model *IncomeReportRequest) (*IncomeReportResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"object":     model,
	}

	// model validator
	if err := grok.Validator.Struct(model); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error income report validator")
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, AccountsPath)
	u.Path = path.Join(u.Path, model.Account)
	u.Path = path.Join(u.Path, IncomeReportPath)

	q := u.Query()

	q.Set("calendar", grok.OnlyDigits(model.Year))

	u.RawQuery = q.Encode()
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error request income report")
		return nil, err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error in authentication request")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", "2.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error performing the request")
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var response *IncomeReportResponse

		respBody, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error unmarshal income report")
			return nil, err
		}

		return response, nil
	}

	return nil, IncomeReportErrorHandler(fields, resp)
}

// IncomeReportErrorHandler ...
func IncomeReportErrorHandler(fields logrus.Fields, resp *http.Response) error {
	var bodyErr *ErrorResponse

	respBody, _ := ioutil.ReadAll(resp.Body)

	err := json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return ErrDefaultIncomeReport
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
