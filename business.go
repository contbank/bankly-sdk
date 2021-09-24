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

//Business ...
type Business struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

//NewBusiness ...
func NewBusiness(httpClient *http.Client, session Session) *Business {
	return &Business{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

//CreateBusiness ...
func (c *Business) CreateBusiness(ctx context.Context, businessRequest BusinessRequest) error {

	if businessRequest.BusinessType == BusinessTypeMEI {
		businessName := businessRequest.LegalRepresentative.RegisterName + " " + businessRequest.LegalRepresentative.Document
		businessRequest.BusinessName = businessName
	}

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	endpoint, err := c.getBusinessAPIEndpoint(requestID, businessRequest.Document, false, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error get business api endpoint")
		return err
	}

	reqbyte, err := json.Marshal(businessRequest)

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

	req = setRequestHeader(req, token, c.session.APIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusAccepted {
		return nil
	} else if resp.StatusCode == http.StatusInternalServerError {
		logrus.
			WithFields(fields).
			Error("internal server error - CreateBusiness")
		return ErrDefaultBusinessAccounts
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error unmarshal - CreateBusiness")
		return err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		logrus.
			WithFields(fields).
			Error("body error - CreateBusiness")
		return FindError(errModel.Code, errModel.Messages...)
	}

	logrus.
		WithFields(fields).
		Error("default error business accounts - CreateBusiness")

	return ErrDefaultBusinessAccounts
}

//UpdateBusiness ...
func (c *Business) UpdateBusiness(ctx context.Context,
	businessDocument string, businessUpdateRequest BusinessUpdateRequest) error {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	endpoint, err := c.getBusinessAPIEndpoint(requestID, businessDocument, false, nil)
	if err != nil {
		return err
	}

	reqbyte, err := json.Marshal(businessUpdateRequest)

	req, err := http.NewRequestWithContext(ctx, "PATCH", *endpoint, bytes.NewReader(reqbyte))
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

	req = setRequestHeader(req, token, c.session.APIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusAccepted {
		return nil
	} else if resp.StatusCode == http.StatusMethodNotAllowed {
		return ErrMethodNotAllowed
	} else if resp.StatusCode == http.StatusInternalServerError {
		logrus.
			WithFields(fields).
			Error("internal server error - UpdateBusiness")
		return ErrDefaultBusinessAccounts
	}

	var bodyErr *ErrorResponse

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
		Error("error default business accounts - UpdateBusiness")

	return ErrDefaultBusinessAccounts
}

//CreateBusinessAccount ...
func (c *Business) CreateBusinessAccount(ctx context.Context,
	businessAccountRequest BusinessAccountRequest) (*AccountResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	endpoint, err := c.getBusinessAPIEndpoint(requestID, businessAccountRequest.Document, true, nil)
	if err != nil {
		return nil, err
	}

	reqbyte, err := json.Marshal(businessAccountRequest)

	req, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(reqbyte))
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

	req = setRequestHeader(req, token, c.session.APIVersion)

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

	if resp.StatusCode == http.StatusCreated {
		var bodyResp *AccountResponse

		err = json.Unmarshal(respBody, &bodyResp)

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error unmarshal - CreateBusinessAccount")
			return nil, err
		}

		return bodyResp, nil
	} else if resp.StatusCode == http.StatusInternalServerError {
		logrus.
			WithFields(fields).
			Error("internal server error - CreateBusinessAccount")
		return nil, ErrDefaultBusinessAccounts
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
		Error("error default business accounts - CreateBusinessAccount")

	return nil, ErrDefaultBusinessAccounts
}

//FindBusiness ...
func (c *Business) FindBusiness(ctx context.Context, identifier string) (*BusinessResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id" : requestID,
		"identifier" : identifier,
	}

	resultLevel := ResultLevelDetailed
	endpoint, err := c.getBusinessAPIEndpoint(requestID, identifier, false, &resultLevel)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		return nil, err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		return nil, err
	}

	req = setRequestHeader(req, token, c.session.APIVersion)

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
		var response BusinessResponse

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
			Info("response with success - FindBusiness")

		return &response, nil
	} else if resp.StatusCode == http.StatusNotFound {
		logrus.
			WithFields(fields).
			WithError(ErrEntryNotFound).
			Error("entry not found - FindBusiness")
		return nil, ErrEntryNotFound
	} else if resp.StatusCode == http.StatusInternalServerError {
		logrus.
			WithFields(fields).
			Error("internal server error - FindBusiness")
		return nil, ErrDefaultBusinessAccounts
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error - FindBusiness")
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		logrus.
			WithFields(fields).
			Error("body error - FindBusiness")
		errModel := bodyErr.Errors[0]
		return nil, FindError(errModel.Code, errModel.Messages...)
	}

	logrus.
		WithFields(fields).
		Error("default error business accounts - FindBusiness")

	return nil, ErrDefaultBusinessAccounts
}

//FindBusinessAccounts ...
func (c *Business) FindBusinessAccounts(ctx context.Context, identifier string) ([]AccountResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id" : requestID,
		"identifier" : identifier,
	}

	endpoint, err := c.getBusinessAPIEndpoint(requestID, identifier, true, nil)
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

	req = setRequestHeader(req, token, c.session.APIVersion)

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
		var response []AccountResponse

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
			Info("response with success - FindBusiness")

		return response, nil
	} else if resp.StatusCode == http.StatusNotFound {
		logrus.
			WithFields(fields).
			WithError(ErrEntryNotFound).
			Error("error entry not found - FindBusiness")
		return nil, ErrEntryNotFound
	} else if resp.StatusCode == http.StatusInternalServerError {
		logrus.
			WithFields(fields).
			Error("internal server error - FindBusiness")
		return nil, ErrDefaultBusinessAccounts
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
		logrus.
			WithFields(fields).
			Error("body error - FindBusiness")
		errModel := bodyErr.Errors[0]
		return nil, FindError(errModel.Code, errModel.Messages...)
	}

	logrus.
		WithFields(fields).
		Error("error default business accounts - FindBusiness")

	return nil, ErrDefaultBusinessAccounts
}

// getBusinessAPIEndpoint
func (c *Business) getBusinessAPIEndpoint(requestID string, identifier string,
	isAccountPath bool, resultLevel *ResultLevel) (*string, error) {

	fields := logrus.Fields {
		"request_id" : requestID,
		"identifier" : identifier,
		"is_account_path" : isAccountPath,
		"result_level" : resultLevel,
	}

	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, BusinessPath)
	u.Path = path.Join(u.Path, grok.OnlyDigits(identifier))

	if isAccountPath == true {
		u.Path = path.Join(u.Path, AccountsPath)
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
