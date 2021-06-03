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
func (p *Payment) ValidatePayment(model *ValidatePaymentRequest) (*ValidatePaymentResponse, error) {

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

	req = setRequestHeader(req, token, p.session.APIVersion)

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

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if bodyErr.Errors != nil {
		return nil, FindError(bodyErr.Errors[0])
	}

	return nil, errors.New("error validate payment")
}

// ConfirmPayment ...
func (p *Payment) ConfirmPayment(model *ConfirmPaymentRequest) (*ConfirmPaymentResponse, error) {

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

	req = setRequestHeader(req, token, p.session.APIVersion)

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

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if bodyErr.Errors != nil {
		return nil, FindError(bodyErr.Errors[0])
	}

	return nil, errors.New("error validate payment")
}
