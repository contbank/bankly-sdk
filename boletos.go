package bankly

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

//Boletos ...
type Boletos struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

//NewBoletos ...
func NewBoletos(httpClient *http.Client, session Session) *Boletos {
	return &Boletos{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(session),
	}
}

//CreateBoleto ...
func (b *Boletos) CreateBoleto(ctx context.Context, model *BoletoRequest) (*BoletoResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	logrus.
		WithField("request", model).
		WithFields(fields).
		Info("creating boleto")

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
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
	endpoint := u.String()

	reqbyte, err := json.Marshal(model)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error encoding model to json")
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(reqbyte))

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating request")
		return nil, err
	}

	token, err := b.authentication.Token()

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error in authentication request")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
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

	fields["bankly_response_status_code"] = resp.StatusCode

	if resp.StatusCode == http.StatusAccepted {
		var body *BoletoResponse

		err = json.Unmarshal(respBody, &body)

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultBoletos
		}

		return body, nil
	}

	var bodyErr []*BoletoErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultBoletos
	}

	if len(bodyErr) > 0 {
		errModel := bodyErr[0]
		err = FindError(errModel.Code, errModel.Message)
		logrus.
			WithField("bankly_error", bodyErr).
			WithFields(fields).
			WithError(err).
			Error("bankly create boleto error")
		return nil, err
	}

	return nil, ErrDefaultBoletos
}

//FindBoleto ...
func (b *Boletos) FindBoleto(ctx context.Context, model *FindBoletoRequest) (*BoletoDetailedResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	logrus.
		WithField("request", model).
		WithFields(fields).
		Info("getting boleto")

	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, "branch")
	u.Path = path.Join(u.Path, model.Account.Branch)
	u.Path = path.Join(u.Path, "number")
	u.Path = path.Join(u.Path, model.Account.Number)
	u.Path = path.Join(u.Path, model.AuthenticationCode)
	endpoint := u.String()

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating request")
		return nil, err
	}

	token, err := b.authentication.Token()

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

	fields["bankly_response_status_code"] = resp.StatusCode

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
			Error("bankly find boleto error")
		return nil, err
	}

	return nil, ErrDefaultBoletos
}

//FilterBoleto ...
func (b *Boletos) FilterBoleto(ctx context.Context, date time.Time) (*FilterBoletoResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	logrus.
		WithField("filter_date", date).
		WithFields(fields).
		Info("filtering boletos")

	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error parsing api endpoint")
		return nil, err
	}
	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, "searchstatus")
	u.Path = path.Join(u.Path, url.QueryEscape(date.UTC().Format("2006-01-02T15:04:05")))
	endpoint := u.String()

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating request")
		return nil, err
	}

	token, err := b.authentication.Token()

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

	fields["bankly_response_status_code"] = resp.StatusCode

	if resp.StatusCode == http.StatusOK {
		var response *FilterBoletoResponse

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
			Error("bankly filteting boleto error")
		return nil, err
	}

	return nil, ErrDefaultBoletos
}

//FindBoletoByBarCode ...
func (b *Boletos) FindBoletoByBarCode(ctx context.Context, barcode string) (*BoletoDetailedResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	logrus.
		WithField("barcode", barcode).
		WithFields(fields).
		Info("finding boleto by barcode")

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

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating request")

		return nil, err
	}

	token, err := b.authentication.Token()

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

	fields["bankly_response_status_code"] = resp.StatusCode

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

//DownloadBoleto ...
func (b *Boletos) DownloadBoleto(ctx context.Context, authenticationCode string, w io.Writer) error {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	logrus.
		WithField("authentication_code", authenticationCode).
		WithFields(fields).
		Info("downloading boleto")

	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, authenticationCode)
	u.Path = path.Join(u.Path, "pdf")
	endpoint := u.String()

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating request")
		return err
	}

	token, err := b.authentication.Token()

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error in authentication request")
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", b.session.APIVersion)

	resp, err := b.httpClient.Do(req)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error performing the request")
		return err
	}

	fields["bankly_response_status_code"] = resp.StatusCode

	if resp.StatusCode == http.StatusOK {
		_, err := io.Copy(w, resp.Body)
		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error writting bytes to writer")
			return ErrDefaultBoletos
		}

		return nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return ErrEntryNotFound
	}

	return ErrDefaultBoletos
}

//CancelBoleto ...
func (b *Boletos) CancelBoleto(ctx context.Context, model *CancelBoletoRequest) error {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	logrus.
		WithField("request", model).
		WithFields(fields).
		Info("canceling boleto")

	err := grok.Validator.Struct(model)

	if err != nil {
		return grok.FromValidationErros(err)
	}

	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error parsing api endpoint")
		return err
	}

	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, "cancel")
	endpoint := u.String()

	reqbyte, err := json.Marshal(model)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error encoding model to json")
		return err
	}

	req, err := http.NewRequest("DELETE", endpoint, bytes.NewReader(reqbyte))

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating request")
		return err
	}

	token, err := b.authentication.Token()

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error in authentication request")
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", b.session.APIVersion)

	resp, err := b.httpClient.Do(req)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error performing the request")
		return err
	}

	defer resp.Body.Close()

	fields["bankly_response_status_code"] = resp.StatusCode

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return ErrDefaultBoletos
}

//SimulatePayment ...
func (b *Boletos) SimulatePayment(ctx context.Context, model *SimulatePaymentRequest) error {
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

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(reqbyte))

	if err != nil {
		return err
	}

	token, err := b.authentication.Token()

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", b.session.APIVersion)

	resp, err := b.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	var bodyErr []*BoletoErrorResponse

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
