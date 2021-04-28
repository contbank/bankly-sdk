package bankly

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/contbank/grok"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

//Business ...
type Business struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

//NewBusiness ...
func NewBusiness(session Session) *Business {
	return &Business{
		session: session,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authentication: NewAuthentication(session),
	}
}

//CreateBusiness ...
func (c *Business) CreateBusiness(businessRequest BusinessRequest) error {
	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, BusinessPath)
	u.Path = path.Join(u.Path, grok.OnlyDigits(businessRequest.Document))
	endpoint := u.String()

	reqbyte, err := json.Marshal(businessRequest)

	req, err := http.NewRequest("PUT", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		return err
	}

	token, err := c.authentication.Token()
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", c.session.APIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusAccepted {
		return nil
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return err
	}

	if bodyErr.Errors != nil {
		return errors.New(bodyErr.Errors[0].Messages[0])
	}
	return errors.New("error create business")
}

//CreateBusinessAccount ...
func (c *Business) CreateBusinessAccount(businessAccountRequest BusinessAccountRequest) (*AccountResponse, error) {
	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, BusinessPath)
	u.Path = path.Join(u.Path, grok.OnlyDigits(businessAccountRequest.Document))
	u.Path = path.Join(u.Path, AccountsPath)
	endpoint := u.String()

	reqbyte, err := json.Marshal(businessAccountRequest)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		return nil, err
	}

	token, err := c.authentication.Token()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", c.session.APIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusCreated {
		var bodyResp *AccountResponse

		err = json.Unmarshal(respBody, &bodyResp)

		if err != nil {
			return nil, err
		}

		return bodyResp, nil
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if bodyErr.Errors != nil {
		return nil, errors.New(bodyErr.Errors[0].Messages[0])
	}
	return nil, errors.New("error create business account")
}

//FindBusiness ...
func (c *Business) FindBusiness(document string) (*BusinessResponse, error) {
	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, BusinessPath)
	u.Path = path.Join(u.Path, document)
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
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", c.session.APIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response BusinessResponse

		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return nil, err
		}

		return &response, nil
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

	return nil, errors.New("error find business")
}


//FindBusinessAccounts ...
func (c *Business) FindBusinessAccounts(document string) ([]AccountResponse, error) {
	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, BusinessPath)
	u.Path = path.Join(u.Path, document)
	u.Path = path.Join(u.Path, AccountsPath)
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
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", c.session.APIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response []AccountResponse

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

	return nil, errors.New("error find business accounts")
}