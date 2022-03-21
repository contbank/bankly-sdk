package bankly

import (
	"context"
	"encoding/json"
	errors "github.com/contbank/bankly-sdk/pkg/errors"
	models "github.com/contbank/bankly-sdk/pkg/models"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

//BankStatement ...
type BankStatement struct {
	session        bankly.Session
	httpClient     *http.Client
	authentication *bankly.Authentication
}

//NewBankStatement ...
func NewBankStatement(httpClient *http.Client, session bankly.Session) *BankStatement {
	return &BankStatement{
		session:        session,
		httpClient:     httpClient,
		authentication: bankly.NewAuthentication(httpClient, session),
	}
}

// FilterBankStatements ...
func (c *BankStatement) FilterBankStatements(ctx context.Context,
	model *models.FilterBankStatementRequest) ([]*models.Statement, error) {

	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
		"object":     model,
	}

	// validator
	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, models.BankStatementsPath)

	q := u.Query()
	q.Set("branch", model.Branch)
	q.Set("account", model.Account)
	q.Set("includeDetails", strconv.FormatBool(model.IncludeDetails))
	q.Set("page", strconv.Itoa(int(model.Page)))
	q.Set("pageSize", strconv.Itoa(int(model.PageSize)))

	// include begin datetime parameter
	if model.BeginDateTime != nil {
		q.Set("beginDateTime", model.BeginDateTime.UTC().Format("2006-01-02T15:04:05"))
	}

	// include end datetime parameter
	if model.EndDateTime != nil {
		q.Set("endDateTime", model.EndDateTime.UTC().Format("2006-01-02T15:04:05"))
	}

	// include card proxies parameters
	if len(model.CardProxy) > 0 {
		for _, c := range model.CardProxy {
			q.Add("cardProxy", c)
		}
	}

	u.RawQuery = q.Encode()
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error creating request")
		return nil, err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error in authentication request")
		return nil, err
	}

	req = utils.SetRequestHeader(req, token, c.session.APIVersion, nil)

	// call bankly
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error performing the request")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response []*models.Statement

		err = json.Unmarshal(respBody, &response)
		if err != nil {
			logrus.WithFields(fields).
				WithError(err).Error("error decoding json response")
			return nil, errors.ErrDefaultBankStatements
		}

		// event status filter
		if model.Status != nil {
			filtered := []*models.Statement{}
			funk.ForEach(response, func(stt *models.Statement) {
				if model.Status != nil && stt.Status == *model.Status {
					filtered = append(filtered, stt)
				}
			})
			response = filtered
		}

		return response, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, errors.ErrEntryNotFound
	}

	var bodyErr *errors.ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error decoding json response")
		return nil, errors.ErrDefaultBankStatements
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err = errors.FindErrorByErrorModel(errModel)

		logrus.WithFields(fields).WithField("bankly_error", bodyErr).
			WithError(err).Error("bankly filter banks statements")
		return nil, err
	}

	return nil, errors.ErrDefaultBankStatements
}
