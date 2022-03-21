package bankly

import (
	"bytes"
	"context"
	"encoding/json"
	errors "github.com/contbank/bankly-sdk/pkg/errors"
	models "github.com/contbank/bankly-sdk/pkg/models"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

//Customers ...
type Customers struct {
	session        bankly.Session
	httpClient     *http.Client
	authentication *bankly.Authentication
}

//NewCustomers ...
func NewCustomers(httpClient *http.Client, session bankly.Session) *Customers {
	return &Customers{
		session:        session,
		httpClient:     httpClient,
		authentication: bankly.NewAuthentication(httpClient, session),
	}
}

//CreateRegistration ...
func (c *Customers) CreateRegistration(ctx context.Context, customer models.CustomersRequest) error {
	err := grok.Validator.Struct(customer)
	if err != nil {
		return grok.FromValidationErros(err)
	}

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"customer":   customer,
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

	req, err := http.NewRequestWithContext(ctx, "PUT", *endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error new request")
		return err
	}

	token, err := c.authentication.Token(ctx)
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

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		return nil
	}

	var bodyErr *errors.ErrorResponse

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
		return errors.FindError(errModel.Code, errModel.Messages...)
	}

	logrus.
		WithFields(fields).
		Error("error default customers accounts - CreateRegistration")

	return errors.ErrDefaultCustomersAccounts
}

//FindRegistration ...
func (c *Customers) FindRegistration(ctx context.Context, identifier string) (*models.CustomersResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"identifier": identifier,
	}

	resultLevel := models.ResultLevelDetailed
	endpoint, err := c.getCustomerAPIEndpoint(requestID, identifier, false, &resultLevel)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error new request")
		return nil, err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	req = utils.SetRequestHeader(req, token, c.session.APIVersion, nil)

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

	if resp.StatusCode == http.StatusOK {
		var response *models.CustomersResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error unmarshal")
			return nil, err
		}

		fields["response"] = response
		logrus.
			WithFields(fields).
			Info("response with success")

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.ErrEntryNotFound
	}

	var bodyErr *errors.ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return nil, errors.FindError(errModel.Code, errModel.Messages...)
	}

	logrus.
		WithFields(fields).
		Error("error default customers accounts - FindRegistration")

	return nil, errors.ErrDefaultCustomersAccounts
}

//UpdateRegistration ...
func (c *Customers) UpdateRegistration(ctx context.Context, document string, customerUpdateRequest models.CustomerUpdateRequest) error {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"document":   document,
	}

	method := "PUT"
	customerResponse, err := c.FindRegistration(ctx, document)
	if customerResponse.Status == models.CustomerStatusApproved {
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

	req, err := http.NewRequestWithContext(ctx, method, *endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		return err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		return err
	}

	req = utils.SetRequestHeader(req, token, c.session.APIVersion, nil)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusAccepted {
		return nil
	} else if resp.StatusCode == http.StatusMethodNotAllowed {
		return errors.ErrMethodNotAllowed
	}

	var bodyErr *errors.ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return errors.FindError(errModel.Code, errModel.Messages...)
	}

	logrus.
		WithFields(fields).
		Error("error default customers accounts - UpdateRegistration")

	return errors.ErrDefaultCustomersAccounts
}

//CreateAccount ...
func (c *Customers) CreateAccount(ctx context.Context, document string, accountType models.AccountType) (*models.AccountResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id":  requestID,
		"document":    document,
		"accountType": accountType,
	}

	endpoint, err := c.getCustomerAPIEndpoint(requestID, document, true, nil)
	if err != nil {
		return nil, err
	}

	model := &models.CustomersAccountRequest{
		AccountType: accountType,
	}

	reqbyte, err := json.Marshal(model)

	req, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		return nil, err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		return nil, err
	}

	req = utils.SetRequestHeader(req, token, c.session.APIVersion, nil)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusCreated {
		var bodyResp *models.AccountResponse

		err = json.Unmarshal(respBody, &bodyResp)

		if err != nil {
			return nil, err
		}

		return bodyResp, nil
	}

	var bodyErr *errors.ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return nil, errors.FindError(errModel.Code, errModel.Messages...)
	}

	logrus.
		WithFields(fields).
		Error("error default customers accounts - CreateAccount")

	return nil, errors.ErrDefaultCustomersAccounts
}

//FindAccounts ...
func (c *Customers) FindAccounts(ctx context.Context, document string) ([]models.AccountResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"document":   document,
	}

	endpoint, err := c.getCustomerAPIEndpoint(requestID, document, true, nil)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error new request")
		return nil, err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	req = utils.SetRequestHeader(req, token, c.session.APIVersion, nil)

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

	if resp.StatusCode == http.StatusOK {
		var response []models.AccountResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error unmarshal")
			return nil, err
		}

		fields["response"] = response
		logrus.
			WithFields(fields).
			Info("response with success")

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.ErrEntryNotFound
	}

	var bodyErr *errors.ErrorResponse

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
		return nil, errors.FindError(errModel.Code, errModel.Messages...)
	}

	logrus.
		WithFields(fields).
		Error("error default customers accounts - FindAccounts")

	return nil, errors.ErrDefaultCustomersAccounts
}

// getCustomerAPIEndpoint
func (c *Customers) getCustomerAPIEndpoint(requestID string, identifier string,
	isAccountPath bool, resultLevel *models.ResultLevel) (*string, error) {

	fields := logrus.Fields{
		"request_id":      requestID,
		"identifier":      identifier,
		"is_account_path": isAccountPath,
		"result_level":    resultLevel,
	}

	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, models.CustomersPath)
	u.Path = path.Join(u.Path, grok.OnlyDigits(identifier))

	if isAccountPath == true {
		u.Path = path.Join(u.Path, models.AccountsPath)
	}

	if resultLevel != nil {
		q := u.Query()
		q.Set("resultLevel", string(*resultLevel))
		u.RawQuery = q.Encode()
	}

	endpoint := u.String()

	fields["endpoint"] = endpoint
	logrus.
		WithFields(fields).
		Info("get endpoint success")

	return &endpoint, nil
}
