package bankly

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/contbank/grok"
)

//Payment ...
type Payment struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

//NewPayment ...
func NewPayment(session Session) *Payment {
	return &Payment{
		session: session,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authentication: NewAuthentication(session),
	}
}

// ValidatePayment ...
func (p *Payment) ValidatePayment(correlationID string, model *ValidatePaymentRequest) (*ValidatePaymentResponse, error) {

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(p.session.APIEndpoint)

	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, PaymentPath)
	u.Path = path.Join(u.Path, "validate")
	endpoint := u.String()

	reqbyte, err := json.Marshal(model)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(reqbyte))

	if err != nil {
		return nil, err
	}

	token, err := p.authentication.Token()

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", p.session.APIVersion)
	req.Header.Add("x-correlation-id", correlationID)

	resp, err := p.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *ValidatePaymentResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			return nil, err
		}

		return response, nil
	}

	var bodyErr *PaymentErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if bodyErr.Code != "" {
		return nil, FindError(ErrorModel{
			Code:     bodyErr.Code,
			Messages: []string{bodyErr.Message},
		})
	}

	return nil, ErrDefaultPayment
}

// ConfirmPayment ...
func (p *Payment) ConfirmPayment(correlationID string, model *ConfirmPaymentRequest) (*ConfirmPaymentResponse, error) {

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(p.session.APIEndpoint)

	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, PaymentPath)
	u.Path = path.Join(u.Path, "confirm")
	endpoint := u.String()

	reqbyte, err := json.Marshal(model)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(reqbyte))

	if err != nil {
		return nil, err
	}

	token, err := p.authentication.Token()

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", p.session.APIVersion)
	req.Header.Add("x-correlation-id", correlationID)

	resp, err := p.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *ConfirmPaymentResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			return nil, err
		}

		return response, nil
	}
	var bodyErr *PaymentErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if bodyErr.Code != "" {
		return nil, FindError(ErrorModel{
			Code:     bodyErr.Code,
			Messages: []string{bodyErr.Message},
		})
	}

	return nil, ErrDefaultPayment
}

// FilterPayments ...
func (p *Payment) FilterPayments(correlationID string, model *FilterPaymentsRequest) (*FilterPaymentsResponse, error) {

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(p.session.APIEndpoint)

	if err != nil {
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

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	token, err := p.authentication.Token()

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", p.session.APIVersion)
	req.Header.Add("x-correlation-id", correlationID)

	resp, err := p.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *FilterPaymentsResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			return nil, err
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr *PaymentErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if bodyErr.Code != "" {
		return nil, FindError(ErrorModel{
			Code:     bodyErr.Code,
			Messages: []string{bodyErr.Message},
		})
	}

	return nil, ErrDefaultPayment
}

// DetailPayment ...
func (p *Payment) DetailPayment(correlationID string, model *DetailPaymentRequest) (*PaymentResponse, error) {

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(p.session.APIEndpoint)

	if err != nil {
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

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	token, err := p.authentication.Token()

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", p.session.APIVersion)
	req.Header.Add("x-correlation-id", correlationID)

	resp, err := p.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *PaymentResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			return nil, err
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr *PaymentErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if bodyErr.Code != "" {
		return nil, FindError(ErrorModel{
			Code:     bodyErr.Code,
			Messages: []string{bodyErr.Message},
		})
	}

	return nil, ErrDefaultPayment
}
