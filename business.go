package bankly

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Business ...
type Business struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

// NewBusiness ...
func NewBusiness(httpClient *http.Client, session Session) *Business {
	return &Business{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

// CreateBusinessRegistration ...
func (c *Business) CreateBusinessRegistration(ctx context.Context, model BusinessRequest) error {

	fields := logrus.Fields{
		"request_id": grok.GetRequestID(ctx),
		"object":     model,
	}

	// validator
	if err := grok.Validator.Struct(model); err != nil {
		logrus.WithFields(fields).Error("invalid model - model")
		return grok.FromValidationErros(err)
	}

	// error if is LTDA or SA (PJOTAO)
	if isCorporationBusiness(model.BusinessType) {
		logrus.WithFields(fields).
			WithError(ErrCorporationBusinessNotAllowed).
			Error("error corporation business not allowed")
		return ErrCorporationBusinessNotAllowed
	}

	model = normalizeBusinessNameMEI(model)
	businessRequest := ParseSimpleBusinessRequest(&model)

	// validator
	if err := grok.Validator.Struct(businessRequest); err != nil {
		logrus.WithFields(fields).Error("invalid model - business request")
		return grok.FromValidationErros(err)
	}

	fields["object"] = businessRequest

	// getting API endpoint URL
	endpoint, err := c.getBusinessAPIEndpoint(grok.GetRequestID(ctx), businessRequest.DocumentNumber, false, nil)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error get business api endpoint")
		return err
	}

	reqbyte, err := json.Marshal(businessRequest)

	req, err := http.NewRequestWithContext(ctx, "PUT", *endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error new request")
		return err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error authentication - CreateBusinessRegistration")
		return err
	}

	req = setRequestHeader(req, token, c.session.APIVersion, nil)

	logrus.WithFields(fields).Info("doing request - CreateBusinessRegistration")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error http client - CreateBusinessRegistration")
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusOK {
		return nil
	} else if resp.StatusCode == http.StatusInternalServerError {
		logrus.WithFields(fields).Error("internal server error - CreateBusinessRegistration")
		return ErrDefaultBusinessAccounts
	}

	var bodyErr *ErrorResponse
	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error unmarshal - CreateBusinessRegistration")
		return err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		logrus.WithFields(fields).Error("body error - CreateBusinessRegistration")
		return FindError(errModel.Code, errModel.Messages...)
	}

	logrus.WithFields(fields).Error("default error business accounts - CreateBusinessRegistration")
	return ErrDefaultBusinessAccounts
}

// CreateCorporationBusinessRequest ...
func (c *Business) CreateCorporationBusinessRequest(ctx context.Context,
	businessRequest CorporationBusinessRequest) error {

	fields := logrus.Fields{
		"request_id": grok.GetRequestID(ctx),
		"object":     businessRequest,
	}

	// validator
	if err := grok.Validator.Struct(businessRequest); err != nil {
		return grok.FromValidationErros(err)
	}

	// error if is MEI, EI or EIRELI (PJOTINHA)
	if !isCorporationBusiness(businessRequest.BusinessType) {
		return ErrSimpleBusinessNotAllowed
	}

	// corporation business size adjust
	businessRequest.BusinessSize =
		string(ParseCorporationBusinessSize(BusinessSize(businessRequest.BusinessSize)))

	fields["object"] = businessRequest

	// getting API endpoint URL
	endpoint, err := c.getCorporationBusinessAPIEndpoint(grok.GetRequestID(ctx), businessRequest.DocumentNumber)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error get corporation business api endpoint")
		return err
	}

	reqbyte, err := json.Marshal(businessRequest)

	req, err := http.NewRequestWithContext(ctx, "PUT", *endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error new request")
		return err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error authentication")
		return err
	}

	req = setRequestHeader(req, token, c.session.APIVersion, nil)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error http client")
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusOK {
		return nil
	} else if resp.StatusCode == http.StatusInternalServerError {
		logrus.WithFields(fields).Error("internal server error - CreateCorporationBusinessRegistration")
		return ErrDefaultCorporationBusinessAccounts
	}

	var bodyErr *ErrorResponse
	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error unmarshal - CreateCorporationBusinessRegistration")
		return err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		logrus.WithFields(fields).Error("body error - CreateCorporationBusinessRegistration")
		return FindError(errModel.Code, errModel.Messages...)
	}

	logrus.WithFields(fields).Error("default error corporation business accounts - CreateCorporationBusinessRegistration")
	return ErrDefaultCorporationBusinessAccounts
}

// UpdateBusiness ...
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

	req = setRequestHeader(req, token, c.session.APIVersion, nil)

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

// CreateBusinessAccount ...
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

	req = setRequestHeader(req, token, c.session.APIVersion, nil)

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

// FindBusiness ...
func (c *Business) FindBusiness(ctx context.Context, identifier string) (*BusinessResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"identifier": identifier,
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

	req = setRequestHeader(req, token, c.session.APIVersion, nil)

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

// FindBusinessAccounts ...
func (c *Business) FindBusinessAccounts(ctx context.Context, identifier string) ([]AccountResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"identifier": identifier,
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

	req = setRequestHeader(req, token, c.session.APIVersion, nil)

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

// CancelBusinessAccount ...
func (c *Business) CancelBusinessAccount(ctx context.Context, identifier string,
	cancelAccountRequest CancelAccountRequest) error {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"identifier": identifier,
	}

	// api endpoint
	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error api endpoint - cancel account")
		return err
	}
	u.Path = path.Join(u.Path, BusinessPath)
	u.Path = path.Join(u.Path, grok.OnlyDigits(identifier))
	u.Path = path.Join(u.Path, "cancel")
	endpoint := u.String()

	reqbyte, err := json.Marshal(cancelAccountRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		return err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		return err
	}

	req = setRequestHeader(req, token, c.session.APIVersion, nil)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusAccepted {
		return nil
	} else if resp.StatusCode == http.StatusNoContent {
		return nil
	} else if resp.StatusCode == http.StatusMethodNotAllowed {
		return ErrMethodNotAllowed
	} else if resp.StatusCode == http.StatusNotFound {
		return ErrAccountNotFound
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

	logrus.WithFields(fields).
		WithError(ErrDefaultCancelCustomersAccounts).Error("error default cancel customers accounts")

	return ErrDefaultCancelCustomersAccounts
}

// getBusinessAPIEndpoint ...
func (c *Business) getBusinessAPIEndpoint(requestID string, identifier string,
	isAccountPath bool, resultLevel *ResultLevel) (*string, error) {

	fields := logrus.Fields{
		"request_id":      requestID,
		"identifier":      identifier,
		"is_account_path": isAccountPath,
		"result_level":    resultLevel,
	}

	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, BusinessPath)
	u.Path = path.Join(u.Path, grok.OnlyDigits(identifier))

	if isAccountPath {
		u.Path = path.Join(u.Path, AccountsPath)
	}

	if resultLevel != nil {
		q := u.Query()
		q.Set("resultLevel", string(*resultLevel))
		u.RawQuery = q.Encode()
	}

	endpoint := u.String()

	fields["endpoint"] = endpoint
	logrus.WithFields(fields).
		Info("get endpoint success")

	return &endpoint, nil
} // getBusinessAPIEndpoint

// getCorporationBusinessAPIEndpoint ...
func (c *Business) getCorporationBusinessAPIEndpoint(requestID string, identifier string) (*string, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"identifier": identifier,
	}

	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CorporationBusinessPath)
	u.Path = path.Join(u.Path, grok.OnlyDigits(identifier))

	endpoint := u.String()

	fields["endpoint"] = endpoint
	logrus.WithFields(fields).
		Info("get endpoint success")

	return &endpoint, nil
}

// normalizeBusinessNameMEI Ajusta o nome da empresa quando MEI, incluindo o identifier do proprietário ao final
func normalizeBusinessNameMEI(businessRequest BusinessRequest) BusinessRequest {
	if businessRequest.BusinessType == BusinessTypeMEI && len(businessRequest.LegalRepresentatives) > 0 {
		// SLU
		if strings.Contains(strings.ToUpper(businessRequest.BusinessName), "LTDA") {
			return businessRequest
		}
		// MEI
		cpf := businessRequest.LegalRepresentatives[0].DocumentNumber
		businessName := businessRequest.LegalRepresentatives[0].RegisterName
		if !strings.Contains(businessName, cpf) {
			businessName = businessRequest.LegalRepresentatives[0].RegisterName + " " + cpf
		}
		businessRequest.BusinessName = businessName
	}
	return businessRequest
}

// isCorporationBusiness ...
func isCorporationBusiness(businessType BusinessType) bool {
	if businessType == BusinessTypeLTDA || businessType == BusinessTypeSA {
		return true
	}
	return false
}
