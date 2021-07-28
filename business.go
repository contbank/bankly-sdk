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
		authentication: NewAuthentication(session),
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
		"request_id" : requestID,
		"request" : businessRequest,
	}

	endpoint, err := c.getBusinessAPIEndpoint(requestID, businessRequest.Document, false)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error get business api endpoint")
		return err
	}

	reqbyte, err := json.Marshal(businessRequest)

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
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	fields["bankly_response_status_code"] = resp.StatusCode

	if resp.StatusCode == http.StatusAccepted {
		return nil
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
		"request_id" : requestID,
		"request" : businessDocument,
	}

	endpoint, err := c.getBusinessAPIEndpoint(requestID, businessDocument, false)
	if err != nil {
		return err
	}

	reqbyte, err := json.Marshal(businessUpdateRequest)

	req, err := http.NewRequest("PATCH", *endpoint, bytes.NewReader(reqbyte))
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
		"request_id" : requestID,
		"request" : businessAccountRequest,
	}

	endpoint, err := c.getBusinessAPIEndpoint(requestID, businessAccountRequest.Document, true)
	if err != nil {
		return nil, err
	}

	reqbyte, err := json.Marshal(businessAccountRequest)

	req, err := http.NewRequest("POST", *endpoint, bytes.NewReader(reqbyte))
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

	if resp.StatusCode == http.StatusCreated {
		var bodyResp *AccountResponse

		err = json.Unmarshal(respBody, &bodyResp)
		fields["bankly_response"] = bodyResp

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error unmarshal")
			return nil, err
		}

		return bodyResp, nil
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
func (c *Business) FindBusiness(ctx context.Context, document string) (*BusinessResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id" : requestID,
		"document" : document,
	}

	endpoint, err := c.getBusinessAPIEndpoint(requestID, document, false)
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
		var response BusinessResponse

		err = json.Unmarshal(respBody, &response)
		fields["bankly_response"] = response

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error unmarshal")
			return nil, err
		}

		return &response, nil
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
func (c *Business) FindBusinessAccounts(ctx context.Context, document string) ([]AccountResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id" : requestID,
		"document" : document,
	}

	endpoint, err := c.getBusinessAPIEndpoint(requestID, document, true)
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
func (c *Business) getBusinessAPIEndpoint(requestID string,
	document string, isAccountPath bool) (*string, error) {
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
	u.Path = path.Join(u.Path, BusinessPath)
	u.Path = path.Join(u.Path, grok.OnlyDigits(document))
	if isAccountPath == true {
		u.Path = path.Join(u.Path, AccountsPath)
	}
	endpoint := u.String()
	return &endpoint, nil
}
