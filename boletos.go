package bankly

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

// Boletos ...
type Boletos struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

// NewBoletos ...
func NewBoletos(httpClient *http.Client, session Session) *Boletos {
	return &Boletos{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

// CreateBankslip
func (b *Boletos) CreateBankslip(ctx context.Context, model *BoletoRequest) (*BoletoResponse, error) {

	// api version
	if model.APIVersion == nil {
		model.APIVersion = aws.String(b.session.APIEndpoint)
	}

	fields := logrus.Fields{
		"request_id":  GetRequestID(ctx),
		"api_version": model.APIVersion,
		"object":      model,
	}

	// validator
	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(b.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, BoletosPath)
	endpoint := u.String()

	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error encoding model to json")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error creating request")
		return nil, err
	}

	token, err := b.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error in authentication request")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", *model.APIVersion)

	// call bankly
	resp, err := b.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error performing the request")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusAccepted {
		var body *BoletoResponse

		err = json.Unmarshal(respBody, &body)
		if err != nil {
			logrus.WithFields(fields).
				WithError(err).Error("error decoding json response")
			return nil, ErrDefaultBoletos
		}

		return body, nil
	}

	var bodyErr []*ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error decoding json response")
		return nil, ErrDefaultBoletos
	}

	if len(bodyErr) > 0 {
		errModel := bodyErr[0]
		err = FindError(errModel.Code, errModel.Message)
		logrus.WithField("bankly_error", bodyErr).WithFields(fields).
			WithError(err).Error("bankly create boleto error")
		return nil, err
	}

	return nil, ErrDefaultBoletos
}

// FindBankslip ...
func (b *Boletos) FindBankslip(ctx context.Context, model *FindBoletoRequest) (*BoletoDetailedResponse, error) {

	// api version
	if model.APIVersion == nil {
		model.APIVersion = aws.String(b.session.APIEndpoint)
	}

	fields := logrus.Fields{
		"request_id":  GetRequestID(ctx),
		"api_version": model.APIVersion,
		"object":      model,
	}

	// validator
	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(b.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, "branch")
	u.Path = path.Join(u.Path, model.Account.Branch)
	u.Path = path.Join(u.Path, "number")
	u.Path = path.Join(u.Path, model.Account.Number)
	u.Path = path.Join(u.Path, model.AuthenticationCode)
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error creating request")
		return nil, err
	}

	token, err := b.authentication.Token(ctx)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error in authentication request")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("api-version", *model.APIVersion)

	resp, err := b.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error performing the request")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *BoletoDetailedResponse

		err = json.Unmarshal(respBody, &response)
		if err != nil {
			logrus.WithFields(fields).
				WithError(err).Error("error decoding json response")
			return nil, ErrDefaultBoletos
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error decoding json response")
		return nil, ErrDefaultBoletos
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err = FindError(errModel.Code, errModel.Messages...)
		logrus.WithField("bankly_error", bodyErr).WithFields(fields).
			WithError(err).Error("bankly find boleto error")
		return nil, err
	}

	return nil, ErrDefaultBoletos
}

// DownloadBankslip ...
func (b *Boletos) DownloadBankslip(ctx context.Context, authenticationCode string, apiVersion *string, w io.Writer) error {

	// api version
	if apiVersion == nil {
		apiVersion = aws.String(b.session.APIEndpoint)
	}

	fields := logrus.Fields{
		"request_id":  grok.GetRequestID(ctx),
		"api_version": apiVersion,
		"object":      authenticationCode,
	}

	if apiVersion == nil {
		return ErrInvalidAPIVersion
	}

	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, authenticationCode)
	u.Path = path.Join(u.Path, "pdf")
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error creating request")
		return err
	}

	token, err := b.authentication.Token(ctx)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error in authentication request")
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", *apiVersion)

	resp, err := b.httpClient.Do(req)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error performing the request")
		return err
	}

	if resp.StatusCode == http.StatusOK {
		_, err := io.Copy(w, resp.Body)
		if err != nil {
			logrus.WithFields(fields).
				WithError(err).Error("error writting bytes to writer")
			return ErrDefaultBoletos
		}

		return nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return ErrEntryNotFound
	}

	return ErrDefaultBoletos
}

// CancelBankslip ...
func (b *Boletos) CancelBankslip(ctx context.Context, model *CancelBoletoRequest) error {

	// api version
	if model.APIVersion == nil {
		model.APIVersion = aws.String(b.session.APIEndpoint)
	}

	fields := logrus.Fields{
		"request_id":  GetRequestID(ctx),
		"api_version": model.APIVersion,
		"object":      model,
	}

	// validator
	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error cancel boleto validator")
		return grok.FromValidationErros(err)
	}

	u, err := url.Parse(b.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error parsing api endpoint")
		return err
	}
	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, "cancel")
	endpoint := u.String()

	// mashal
	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error encoding model to json")
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error creating request")
		return err
	}

	token, err := b.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error in authentication request")
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", *model.APIVersion)

	// call bankly
	resp, err := b.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error performing the request")
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	logrus.WithFields(fields).Error("error default - cancel boletos")

	return ErrDefaultBoletos
}

// SandboxSimulateBankslipPayment ...
func (b *Boletos) SandboxSimulateBankslipPayment(ctx *context.Context, model *SandboxSimulateBankslipPaymentRequest) error {

	// api version
	if model.APIVersion == nil {
		model.APIVersion = aws.String(b.session.APIEndpoint)
	}

	fields := logrus.Fields{
		"request_id":  GetRequestID(*ctx),
		"api_version": model.APIVersion,
		"object":      model,
	}

	// validator
	if err := grok.Validator.Struct(model); err != nil {
		return grok.FromValidationErros(err)
	}

	u, err := url.Parse(b.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error parsing api endpoint")
		return err
	}
	u.Path = path.Join(u.Path, BoletosSettledPath)
	endpoint := u.String()

	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error encoding model to json")
		return err
	}

	req, err := http.NewRequestWithContext(*ctx, "POST", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error creating request")
		return err
	}

	token, err := b.authentication.Token(*ctx)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error in authentication request")
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", *model.APIVersion)

	resp, err := b.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error performing the request")
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK ||
		resp.StatusCode == http.StatusCreated ||
		resp.StatusCode == http.StatusAccepted {
		return nil
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	var bodyErr []*ErrorResponse

	if err := json.Unmarshal(respBody, &bodyErr); err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error decoding json response")
		return ErrDefaultBoletos
	}

	if len(bodyErr) > 0 {
		errModel := bodyErr[0]
		if err := FindError(errModel.Code, errModel.Message); err != nil {
			logrus.WithField("bankly_error", bodyErr).WithFields(fields).
				WithError(err).Error("bankly create boleto error")
			return err
		}
	}

	return ErrDefaultBoletos
}

/*
// FilterBoleto ...
func (b *Boletos) FilterBoleto(ctx context.Context, date time.Time) (*FilterBoletoResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error parsing api endpoint")
		return nil, err
	}
	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, "searchstatus")
	u.Path = path.Join(u.Path, url.QueryEscape(date.UTC().Format("2006-01-02")))
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error creating request")
		return nil, err
	}

	token, err := b.authentication.Token(ctx)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error in authentication request")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("api-version", b.session.APIVersion)

	resp, err := b.httpClient.Do(req)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error performing the request")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *FilterBoletoResponse
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			logrus.WithFields(fields).
				WithError(err).Error("error decoding json response")
			return nil, ErrDefaultBoletos
		}
		return response, nil
	} else if resp.StatusCode == http.StatusNotFound {
		logrus.WithFields(fields).Info("not found")
		return nil, ErrBoletoNotFound
	}

	var bodyErr *ErrorResponse
	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error decoding json response")
		return nil, ErrDefaultBoletos
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err = FindError(errModel.Code, errModel.Messages...)
		logrus.WithFields(fields).WithField("bankly_error", bodyErr).
			WithError(err).Error("bankly filteting boleto error")
		return nil, err
	}

	return nil, ErrDefaultBoletos
}
*/

/*
// FindBoletoByBarCode ...
func (b *Boletos) FindBoletoByBarCode(ctx context.Context, barcode string) (*BoletoDetailedResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, barcode)
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating request")

		return nil, err
	}

	token, err := b.authentication.Token(ctx)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error in authentication request")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("api-version", b.session.APIVersion)

	resp, err := b.httpClient.Do(req)

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
		var response *BoletoDetailedResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultBoletos
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
		return nil, ErrDefaultBoletos
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err = FindError(errModel.Code, errModel.Messages...)
		logrus.
			WithField("bankly_error", bodyErr).
			WithFields(fields).
			WithError(err).
			Error("bankly find boleto by barcode error")
		return nil, err
	}

	return nil, ErrDefaultBoletos
}
*/

/*
// SimulatePayment ...
func (b *Boletos) SimulatePayment(ctx context.Context, model *SimulatePaymentRequest) error {

	// api version
	if model.APIVersion == nil {
		model.APIVersion = aws.String(b.session.APIEndpoint)
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		return grok.FromValidationErros(err)
	}

	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, "settlementpayment")
	endpoint := u.String()

	reqbyte, err := json.Marshal(model)

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqbyte))

	if err != nil {
		return err
	}

	token, err := b.authentication.Token(ctx)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", *model.APIVersion)

	resp, err := b.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	var bodyErr []*ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		return err
	}

	if len(bodyErr) > 0 {
		err := bodyErr[0]
		return FindError(err.Code, err.Message)
	}

	return ErrDefaultBoletos
}
*/
