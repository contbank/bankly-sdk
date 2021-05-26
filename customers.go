package bankly

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
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

	endpoint, err := c.getCustomerAPIEndpoint(customer.Document, false, nil)
	if err != nil {
		return err
	}

	reqbyte, err := json.Marshal(customer)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", *endpoint, bytes.NewReader(reqbyte))
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

	resultLevel := ResultLevelDetailed
	endpoint, err := c.getCustomerAPIEndpoint(document, false, &resultLevel)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", *endpoint, nil)
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
		var response *CustomersResponse

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

	if bodyErr.Errors != nil {
		return nil, FindError(bodyErr.Errors[0])
	}

	return nil, errors.New("error find registration")
}

//UpdateRegistration ...
func (c *Customers) UpdateRegistration(document string, customerUpdateRequest CustomerUpdateRequest) error {

	method := "PUT"
	customerResponse, err := c.FindRegistration(document)
	if customerResponse.Status == string(CustomerStatusApproved) {
		method = "PATCH"
	}

	endpoint, err := c.getCustomerAPIEndpoint(document, false, nil)
	if err != nil {
		return err
	}

	reqbyte, err := json.Marshal(customerUpdateRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, *endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		return err
	}

	token, err := c.authentication.Token()
	if err != nil {
		return err
	}

	req = setRequestHeader(req, token, c.session.APIVersion)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusAccepted {
		return nil
	} else if resp.StatusCode == http.StatusMethodNotAllowed {
		return ErrMethodNotAllowed
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return err
	}

	if bodyErr.Errors != nil {
		return FindError(bodyErr.Errors[0])
	}
	return errors.New("error updating customer")
}

//CreateAccount ...
func (c *Customers) CreateAccount(document string, accountType AccountType) (*AccountResponse, error) {

	endpoint, err := c.getCustomerAPIEndpoint(document, true, nil)
	if err != nil {
		return nil, err
	}

	model := &CustomersAccountRequest{
		AccountType: accountType,
	}

	reqbyte, err := json.Marshal(model)

	req, err := http.NewRequest("POST", *endpoint, bytes.NewReader(reqbyte))
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

	endpoint, err := c.getCustomerAPIEndpoint(document, true, nil)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", *endpoint, nil)
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
		var response []AccountResponse

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

	if bodyErr.Errors != nil {
		return nil, FindError(bodyErr.Errors[0])
	}

	return nil, errors.New("error find accounts")
}

// getCustomerAPIEndpoint
func (c *Customers) getCustomerAPIEndpoint(document string, isAccountPath bool, resultLevel *ResultLevel) (*string, error) {
	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.
			WithError(err).
			Error("error api endpoint")
		return nil, err
	}
	u.Path = path.Join(u.Path, CustomersPath)
	u.Path = path.Join(u.Path, grok.OnlyDigits(document))
	if isAccountPath == true {
		u.Path = path.Join(u.Path, AccountsPath)
	}
	if resultLevel != nil {
		q := u.Query()
		q.Set("resultLevel", string(*resultLevel))
		u.RawQuery = q.Encode()
	}
	endpoint := u.String()
	return &endpoint, nil
}

// setRequestHeader
func setRequestHeader(request *http.Request, token string, apiVersion string) (*http.Request) {
	request.Header.Add("Authorization", token)
	request.Header.Add("Content-type", "application/json")
	request.Header.Add("api-version", apiVersion)
	return request
}