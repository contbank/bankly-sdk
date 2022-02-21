package bankly

import (
	"context"
	"encoding/json"
	"fmt"
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
	httpClient     *http.Client
	authentication *Authentication
}

// NewIncomeReport ...
func NewIncomeReport(httpClient *http.Client, session Session) *IncomeReport {
	return &IncomeReport {
		session : session,
		httpClient : httpClient,
		authentication : NewAuthentication(httpClient, session),
	}
}

// GetIncomeReport ...
func (c *IncomeReport) GetIncomeReport(ctx context.Context, model *IncomeReportRequest) (*IncomeReportResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"object" : model,
	}

	// model validator
	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	// endpoint url
	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.WithError(err).
			WithFields(fields).Error("error parsing api endpoint")
		return nil, err
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("%s/%s/income-report/print", AccountsPath, model.Account))

	// parameters
	q := u.Query()
	q.Set("calendar", model.Year)

	u.RawQuery = q.Encode()
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("error creating request")
		return nil, err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("error in authentication request")
		return nil, err
	}

	// request header
	req = setRequestHeader(req, token, c.session.APIVersion, nil)

	// do bankly request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.
			WithError(err).
			WithFields(fields).
			Error("error performing the request")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *IncomeReportResponse

		err = json.Unmarshal(respBody, &response)
		if err != nil {
			logrus.WithError(err).
				WithFields(fields).Error("error decoding json response")
			return nil, ErrDefaultBankStatements
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.WithError(err).
			WithFields(fields).Error("error decoding json response")
		return nil, ErrDefaultIncomeReport
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err = FindErrorByErrorModel(errModel)

		logrus.WithFields(fields).WithField("bankly_error", bodyErr).
			WithError(err).Error("error bankly report income")

		return nil, err
	}

	return nil, ErrDefaultIncomeReport
}