package bankly

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

//Balance ...
type Balance struct {
	session        Session
	authentication *Authentication
	httpClient     *http.Client
}

//NewBalance ...
func NewBalance(httpClient *http.Client, session Session) *Balance {
	return &Balance{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(session),
	}
}

//Balance ...
func (c *Balance) Balance(account string) (*AccountResponse, error) {
	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, AccountsPath)
	u.Path = path.Join(u.Path, account)

	q := u.Query()

	q.Set("includeBalance", "true")

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
		var response *AccountResponse

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
		errModel := bodyErr.Errors[0]
		return nil, FindError(errModel.Code, errModel.Messages...)
	}

	return nil, ErrDefaultBalance
}
