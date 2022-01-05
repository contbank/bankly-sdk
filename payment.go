package bankly

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

//Payment ...
type Payment struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

//NewPayment ...
func NewPayment(httpClient *http.Client, session Session) *Payment {
	return &Payment{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

// ValidatePayment ...
func (p *Payment) ValidatePayment(ctx context.Context, correlationID string, model *ValidatePaymentRequest) (*ValidatePaymentResponse, error) {
	fields := logrus.Fields{
		"request_id": correlationID,
	}

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(p.session.APIEndpoint)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, PaymentPath)
	u.Path = path.Join(u.Path, "validate")
	endpoint := u.String()

	reqbyte, err := json.Marshal(model)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error encoding model to json")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqbyte))

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating request")
		return nil, err
	}

	token, err := p.authentication.Token(ctx)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error in authentication request")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", p.session.APIVersion)
	req.Header.Add("x-correlation-id", correlationID)

	resp, err := p.httpClient.Do(req)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error performing the request")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *ValidatePaymentResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultPayment
		}

		return response, nil
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultPayment
	}

	if bodyErr.Code != "" {
		err = FindError(bodyErr.Code, bodyErr.Message)
		logrus.
			WithField("bankly_error", bodyErr).
			WithFields(fields).
			WithError(err).
			Error("bankly validate payment error")
		return nil, err
	}

	return nil, ErrDefaultPayment
}

// ConfirmPayment ...
func (p *Payment) ConfirmPayment(ctx context.Context, correlationID string, model *ConfirmPaymentRequest) (*ConfirmPaymentResponse, error) {

	fields := logrus.Fields{
		"request_id": correlationID,
	}

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(p.session.APIEndpoint)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, PaymentPath)
	u.Path = path.Join(u.Path, "confirm")
	endpoint := u.String()

	reqbyte, err := json.Marshal(model)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error encoding model to json")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqbyte))

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating request")
		return nil, err
	}

	token, err := p.authentication.Token(ctx)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error in authentication request")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", p.session.APIVersion)
	req.Header.Add("x-correlation-id", correlationID)

	resp, err := p.httpClient.Do(req)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error performing the request")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *ConfirmPaymentResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultPayment
		}

		return response, nil
	}
	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultPayment
	}

	if bodyErr.Code != "" {
		err = FindError(bodyErr.Code, bodyErr.Message)
		logrus.
			WithField("bankly_error", bodyErr).
			WithFields(fields).
			WithError(err).
			Error("bankly confirm payment error")
		return nil, err
	}

	return nil, ErrDefaultPayment
}

// FilterPayments ...
func (p *Payment) FilterPayments(ctx context.Context, correlationID string, model *FilterPaymentsRequest) (*FilterPaymentsResponse, error) {

	fields := logrus.Fields{
		"request_id": correlationID,
	}

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(p.session.APIEndpoint)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, PaymentPath)

	q := u.Query()
	q.Set("bankAccount", model.BankAccount)
	q.Set("bankBranch", model.BankBranch)
	q.Set("pageSize", strconv.Itoa(model.PageSize))

	if model.PageToken != nil {
		q.Set("pageToken", *model.PageToken)
	}

	u.RawQuery = q.Encode()
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating request")
		return nil, err
	}

	token, err := p.authentication.Token(ctx)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error in authentication request")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", p.session.APIVersion)
	req.Header.Add("x-correlation-id", correlationID)

	resp, err := p.httpClient.Do(req)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error performing the request")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *FilterPaymentsResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultPayment
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
			Error("error decoding json response")
		return nil, ErrDefaultPayment
	}

	if bodyErr.Code != "" {
		err = FindError(bodyErr.Code, bodyErr.Message)
		logrus.
			WithField("bankly_error", bodyErr).
			WithFields(fields).
			WithError(err).
			Error("bankly filter payments error")
		return nil, err
	}

	return nil, ErrDefaultPayment
}

// DetailPayment ...
func (p *Payment) DetailPayment(ctx context.Context, correlationID string, model *DetailPaymentRequest) (*PaymentResponse, error) {

	fields := logrus.Fields{
		"request_id": correlationID,
	}

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(p.session.APIEndpoint)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, PaymentPath)
	u.Path = path.Join(u.Path, "detail")

	q := u.Query()
	q.Set("bankAccount", model.BankAccount)
	q.Set("bankBranch", model.BankBranch)
	q.Set("authenticationCode", model.AuthenticationCode)

	u.RawQuery = q.Encode()
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating request")
		return nil, err
	}

	token, err := p.authentication.Token(ctx)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error in authentication request")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", p.session.APIVersion)
	req.Header.Add("x-correlation-id", correlationID)

	resp, err := p.httpClient.Do(req)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error performing the request")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *PaymentResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultPayment
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
			Error("error decoding json response")
		return nil, ErrDefaultPayment
	}

	if bodyErr.Code != "" {
		err = FindError(bodyErr.Code, bodyErr.Message)
		logrus.
			WithField("bankly_error", bodyErr).
			WithFields(fields).
			WithError(err).
			Error("bankly get payment detail error")
		return nil, err
	}

	return nil, ErrDefaultPayment
}
