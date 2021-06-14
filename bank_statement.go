package bankly

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/contbank/grok"
)

//BankStatement ...
type BankStatement struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

//NewBankStatement ...
func NewBankStatement(session Session) *BankStatement {
	return &BankStatement{
		session: session,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authentication: NewAuthentication(session),
	}
}

// FilterBankStatements ...
func (c *BankStatement) FilterBankStatements(model *FilterBankStatementRequest) ([]*Statement, error) {

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
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

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	token, err := c.authentication.Token()
	if err != nil {
		return nil, err
	}

	req = setRequestHeader(req, token, c.session.APIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response []*Statement

		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return nil, err
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		return nil, FindError(bodyErr.Errors[0])
	}

	return nil, ErrDefaultBankStatements
}
