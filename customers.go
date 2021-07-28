package bankly

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

//Customers ...
type Customers struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

//NewCustomers ...
func NewCustomers(httpClient *http.Client, session Session) *Customers {
	return &Customers{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(session),
	}
}

//CreateRegistration ...
func (c *Customers) CreateRegistration(ctx context.Context, customer CustomersRequest) error {
	err := grok.Validator.Struct(customer)
	if err != nil {
		return grok.FromValidationErros(err)
	}

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id" : requestID,
		"customer" : customer,
	}

	endpoint, err := c.getCustomerAPIEndpoint(requestID, customer.Document, false, nil)
	if err != nil {
		return err
	}

	reqbyte, err := json.Marshal(customer)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error marshal")
		return err
	}

	req, err := http.NewRequest("PUT", *endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error new request")
		return err
	}

	token, err := c.authentication.Token()
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error authentication")
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", c.session.APIVersion)

	resp, err := c.httpClient.Do(req)

	fields["bankly_request_host"] = req.URL.Host
	fields["bankly_request_path"] = req.URL.Path
	fields["bankly_request_header_api_version"] = req.Header.Get("api-version")

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return err
	}

	defer resp.Body.Close()

	fields["bankly_response_status_code"] = resp.StatusCode

	if resp.StatusCode == http.StatusAccepted {
		return nil
	}

	var bodyErr *ErrorResponse

	respBody, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error unmarshal")
		return err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return FindError(errModel.Code, errModel.Messages...)
	}

	logrus.
		WithFields(fields).
		Error("error default customers accounts - CreateRegistration")

	return ErrDefaultCustomersAccounts
}

//FindRegistration ...
func (c *Customers) FindRegistration(ctx context.Context, document string) (*CustomersResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id" : requestID,
		"document" : document,
	}

	resultLevel := ResultLevelDetailed
	endpoint, err := c.getCustomerAPIEndpoint(requestID, document, false, &resultLevel)
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

	fields["bankly_request_host"] = req.URL.Host
	fields["bankly_request_path"] = req.URL.Path
	fields["bankly_request_header_api_version"] = req.Header.Get("api-version")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	fields["bankly_response_status_code"] = resp.StatusCode

	if resp.StatusCode == http.StatusOK {
		var response *CustomersResponse

		err = json.Unmarshal(respBody, &response)
		fields["bankly_response"] = response

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

	logrus.
		WithFields(fields).
		Error("error default customers accounts - FindRegistration")

	return nil, ErrDefaultCustomersAccounts
}

//UpdateRegistration ...
func (c *Customers) UpdateRegistration(ctx context.Context, document string, customerUpdateRequest CustomerUpdateRequest) error {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id" : requestID,
		"document" : document,
		"customerUpdateRequest" : customerUpdateRequest,
	}

	method := "PUT"
	customerResponse, err := c.FindRegistration(ctx, document)
	if customerResponse.Status == CustomerStatusApproved {
		method = "PATCH"
	}

	endpoint, err := c.getCustomerAPIEndpoint(requestID, document, false, nil)
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

	fields["bankly_request_host"] = req.URL.Host
	fields["bankly_request_path"] = req.URL.Path
	fields["bankly_request_header_api_version"] = req.Header.Get("api-version")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	fields["bankly_response_status_code"] = resp.StatusCode

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

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return FindError(errModel.Code, errModel.Messages...)
	}

	logrus.
		WithFields(fields).
		Error("error default customers accounts - UpdateRegistration")

	return ErrDefaultCustomersAccounts
}

//CreateAccount ...
func (c *Customers) CreateAccount(ctx context.Context, document string, accountType AccountType) (*AccountResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id" : requestID,
		"document" : document,
		"accountType" : accountType,
	}

	endpoint, err := c.getCustomerAPIEndpoint(requestID, document, true, nil)
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

	fields["bankly_request_host"] = req.URL.Host
	fields["bankly_request_path"] = req.URL.Path
	fields["bankly_request_header_api_version"] = req.Header.Get("api-version")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	fields["bankly_response_status_code"] = resp.StatusCode

	if resp.StatusCode == http.StatusCreated {
		var bodyResp *AccountResponse

		err = json.Unmarshal(respBody, &bodyResp)
		fields["bankly_response"] = bodyResp

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

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return nil, FindError(errModel.Code, errModel.Messages...)
	}

	logrus.
		WithFields(fields).
		Error("error default customers accounts - CreateAccount")

	return nil, ErrDefaultCustomersAccounts
}

//FindAccounts ...
func (c *Customers) FindAccounts(ctx context.Context, document string) ([]AccountResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id" : requestID,
		"document" : document,
	}

	endpoint, err := c.getCustomerAPIEndpoint(requestID, document, true, nil)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", *endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error new request")
		return nil, err
	}

	token, err := c.authentication.Token()
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	req = setRequestHeader(req, token, c.session.APIVersion)

	fields["bankly_request_host"] = req.URL.Host
	fields["bankly_request_path"] = req.URL.Path
	fields["bankly_request_header_api_version"] = req.Header.Get("api-version")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	fields["bankly_response_status_code"] = resp.StatusCode

	if resp.StatusCode == http.StatusOK {
		var response []AccountResponse

		err = json.Unmarshal(respBody, &response)
		fields["bankly_response"] = response

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error unmarshal")
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
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return nil, FindError(errModel.Code, errModel.Messages...)
	}

	logrus.
		WithFields(fields).
		Error("error default customers accounts - FindAccounts")

	return nil, ErrDefaultCustomersAccounts
}

// getCustomerAPIEndpoint
func (c *Customers) getCustomerAPIEndpoint(requestID string, document string, isAccountPath bool, resultLevel *ResultLevel) (*string, error) {
	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"request_id" : requestID,
			}).
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
func setRequestHeader(request *http.Request, token string, apiVersion string) *http.Request {
	request.Header.Add("Authorization", token)
	request.Header.Add("Content-type", "application/json")
	request.Header.Add("api-version", apiVersion)
	return request
}
