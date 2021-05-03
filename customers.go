package bankly

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

//Customers ...
type Customers struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

//NewCustomers ...
func NewCustomers(session Session) *Customers {
	return &Customers{
		session: session,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authentication: NewAuthentication(session),
	}
}

//CreateRegistration ...
func (c *Customers) CreateRegistration(customer CustomersRequest) error {
	err := Validator.Struct(customer)

	if err != nil {
		return err
	}

	u, err := url.Parse(c.session.APIEndpoint)

	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, CustomersPath)
	u.Path = path.Join(u.Path, customer.Document)
	endpoint := u.String()

	reqbyte, err := json.Marshal(customer)

	if err != nil {
		return err
	}

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

	if resp.StatusCode == http.StatusAccepted {
		return nil
	}

	var bodyErr *ErrorResponse

	respBody, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		return err
	}

	if bodyErr.Errors != nil {
		return FindError(bodyErr.Errors[0])
	}

	return errors.New("error create registration")
}

//FindRegistration ...
func (c *Customers) FindRegistration(document string) (*CustomersResponse, error) {
	u, err := url.Parse(c.session.APIEndpoint)

	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, CustomersPath)
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
		var response *CustomersResponse

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
		return nil, FindError(bodyErr.Errors[0])
	}

	return nil, errors.New("error find registration")

}

//CreateAccount ...
func (c *Customers) CreateAccount(document string, accountType AccountType) (*AccountResponse, error) {
	u, err := url.Parse(c.session.APIEndpoint)

	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, CustomersPath)
	u.Path = path.Join(u.Path, document)
	u.Path = path.Join(u.Path, AccountsPath)
	endpoint := u.String()

	model := &CustomersAccountRequest{
		AccountType: accountType,
	}

	reqbyte, err := json.Marshal(model)

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
		return nil, FindError(bodyErr.Errors[0])
	}

	return nil, errors.New("error create account")
}

//FindAccounts ...
func (c *Customers) FindAccounts(document string) ([]AccountResponse, error) {
	u, err := url.Parse(c.session.APIEndpoint)

	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, CustomersPath)
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
		return nil, FindError(bodyErr.Errors[0])
	}

	return nil, errors.New("error find accounts")
}
