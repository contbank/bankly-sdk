package bankly

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

// Bank ...
type Bank struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

//NewBank ...
func NewBank(session Session) *Bank {
	return &Bank{
		session: session,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authentication: NewAuthentication(session),
	}
}

//GetByID returns a list with all available financial instituitions
func (c *Bank) GetByID(id string) (*BankDataResponse, error) {
	u, err := url.Parse(c.session.APIEndpoint)

	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, BanksPath)
	u.Path = path.Join(u.Path, id)
	endpoint := u.String()

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	token, err := c.authentication.Token()

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("api-version", c.session.APIVersion)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *BankDataResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			return nil, err
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("not found")
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		return nil, err
	}

	if bodyErr.Errors != nil {
		return nil, errors.New(bodyErr.Errors[0].Messages[0])
	}

	return nil, errors.New("error bank")
}

//List returns a list with all available financial instituitions
func (c *Bank) List(filter *FilterBankListRequest) ([]*BankDataResponse, error) {
	u, err := url.Parse(c.session.APIEndpoint)

	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, BanksPath)

	q := u.Query()

	for _, id := range filter.IDs {
		q.Add("id", id)
	}

	if filter.Name != nil {
		q.Set("name", *filter.Name)
	}

	if filter.Product != nil {
		q.Set("product", *filter.Product)
	}

	if filter.Page != nil {
		q.Set("page", strconv.Itoa(*filter.Page))
	}

	if filter.Name != nil {
		q.Set("pageSize", strconv.Itoa(*filter.PageSize))
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

	req.Header.Add("Authorization", token)
	req.Header.Add("api-version", c.session.APIVersion)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response []*BankDataResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			return nil, err
		}

		return response, nil
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		return nil, err
	}

	if bodyErr.Errors != nil {
		return nil, errors.New(bodyErr.Errors[0].Messages[0])
	}

	return nil, errors.New("error bank")
}
