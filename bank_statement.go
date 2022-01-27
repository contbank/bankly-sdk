package bankly

import (
	"context"
	"encoding/json"
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
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

//NewBankStatement ...
func NewBankStatement(httpClient *http.Client, session Session) *BankStatement {
	return &BankStatement{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

// FilterBankStatements ...
func (c *BankStatement) FilterBankStatements(ctx context.Context, model *FilterBankStatementRequest) ([]*Statement, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.
			WithError(err).
			WithFields(fields).
			Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, BankStatementsPath)

	q := u.Query()

	q.Set("branch", model.Branch)
	q.Set("account", model.Account)
	q.Set("includeDetails", strconv.FormatBool(model.IncludeDetails))
	q.Set("page", strconv.Itoa(int(model.Page)))
	q.Set("pageSize", strconv.Itoa(int(model.PageSize)))

	if model.BeginDateTime != nil {
		q.Set("beginDateTime", model.BeginDateTime.UTC().Format("2006-01-02T15:04:05"))
	}

	if model.EndDateTime != nil {
		q.Set("endDateTime", model.EndDateTime.UTC().Format("2006-01-02T15:04:05"))
	}

	if len(model.CardProxy) > 0 {
		for _, c := range model.CardProxy {
			q.Add("cardProxy", c)
		}
	}

	u.RawQuery = q.Encode()
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		logrus.
			WithError(err).
			WithFields(fields).
			Error("error creating request")
		return nil, err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		logrus.
			WithError(err).
			WithFields(fields).
			Error("error in authentication request")
		return nil, err
	}

	req = setRequestHeader(req, token, c.session.APIVersion, nil)

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
		var response []*Statement

		err = json.Unmarshal(respBody, &response)
		if err != nil {
			logrus.
				WithError(err).
				WithFields(fields).
				Error("error decoding json response")
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
		logrus.
			WithError(err).
			WithFields(fields).
			Error("error decoding json response")
		return nil, ErrDefaultBankStatements
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err = FindErrorByErrorModel(errModel)

		logrus.
			WithError(err).
			WithFields(fields).
			WithField("bankly_error", bodyErr).
			Error("bankly filter banks statements")

		return nil, err
	}

	return nil, ErrDefaultBankStatements
}
