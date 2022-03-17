package bankly

import (
	"bytes"
	"context"
	"encoding/json"
	errors "github.com/contbank/bankly-sdk/pkg/errors"
	models "github.com/contbank/bankly-sdk/pkg/models"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

// Boletos ...
type Boletos struct {
	session    bankly.Session
	httpClient *http.Client
	authentication *bankly.Authentication
}

// NewBoletos ...
func NewBoletos(httpClient *http.Client, session bankly.Session) *Boletos {
	return &Boletos{
		session:        session,
		httpClient:     httpClient,
		authentication: bankly.NewAuthentication(httpClient, session),
	}
}

// CreateBoleto ...
func (b *Boletos) CreateBoleto(ctx context.Context, model *models.BoletoRequest) (*models.BoletoResponse, error) {

	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
	}

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, models.BoletosPath)
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
	req.Header.Add("api-version", b.session.APIVersion)

	resp, err := b.httpClient.Do(req)

	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error performing the request")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusAccepted {
		var body *models.BoletoResponse

		err = json.Unmarshal(respBody, &body)

		if err != nil {
			logrus.WithFields(fields).
				WithError(err).Error("error decoding json response")
			return nil, errors.ErrDefaultBoletos
		}

		return body, nil
	}

	var bodyErr []*errors.ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error decoding json response")
		return nil, errors.ErrDefaultBoletos
	}

	if len(bodyErr) > 0 {
		errModel := bodyErr[0]
		err = errors.FindError(errModel.Code, errModel.Message)
		logrus.
			WithField("bankly_error", bodyErr).
			WithFields(fields).
			WithError(err).
			Error("bankly create boleto error")
		return nil, err
	}

	return nil, errors.ErrDefaultBoletos
}

// FindBoleto ...
func (b *Boletos) FindBoleto(ctx context.Context, model *models.FindBoletoRequest) (*models.BoletoDetailedResponse, error) {
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

	u.Path = path.Join(u.Path, models.BoletosPath)
	u.Path = path.Join(u.Path, "branch")
	u.Path = path.Join(u.Path, model.Account.Branch)
	u.Path = path.Join(u.Path, "number")
	u.Path = path.Join(u.Path, model.Account.Number)
	u.Path = path.Join(u.Path, model.AuthenticationCode)
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
		var response *models.BoletoDetailedResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error decoding json response")
			return nil, errors.ErrDefaultBoletos
		}

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
			Error("error decoding json response")
		return nil, errors.ErrDefaultBoletos
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err = errors.FindError(errModel.Code, errModel.Messages...)
		logrus.
			WithField("bankly_error", bodyErr).
			WithFields(fields).
			WithError(err).
			Error("bankly find boleto error")
		return nil, err
	}

	return nil, errors.ErrDefaultBoletos
}

//FilterBoleto ...
func (b *Boletos) FilterBoleto(ctx context.Context, date time.Time) (*models.FilterBoletoResponse, error) {
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
	u.Path = path.Join(u.Path, models.BoletosPath)
	u.Path = path.Join(u.Path, "searchstatus")
	u.Path = path.Join(u.Path, url.QueryEscape(date.UTC().Format("2006-01-02T15:04:05")))
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
		var response *models.FilterBoletoResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error decoding json response")
			return nil, errors.ErrDefaultBoletos
		}

		return response, nil
	}

	var bodyErr *errors.ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error decoding json response")
		return nil, errors.ErrDefaultBoletos
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err = errors.FindError(errModel.Code, errModel.Messages...)
		logrus.
			WithField("bankly_error", bodyErr).
			WithFields(fields).
			WithError(err).
			Error("bankly filteting boleto error")
		return nil, err
	}

	return nil, errors.ErrDefaultBoletos
}

//FindBoletoByBarCode ...
func (b *Boletos) FindBoletoByBarCode(ctx context.Context, barcode string) (*models.BoletoDetailedResponse, error) {
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

	u.Path = path.Join(u.Path, models.BoletosPath)
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
		var response *models.BoletoDetailedResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error decoding json response")
			return nil, errors.ErrDefaultBoletos
		}

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
			Error("error decoding json response")
		return nil, errors.ErrDefaultBoletos
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err = errors.FindError(errModel.Code, errModel.Messages...)
		logrus.
			WithField("bankly_error", bodyErr).
			WithFields(fields).
			WithError(err).
			Error("bankly find boleto by barcode error")
		return nil, err
	}

	return nil, errors.ErrDefaultBoletos
}

//DownloadBoleto ...
func (b *Boletos) DownloadBoleto(ctx context.Context, authenticationCode string, w io.Writer) error {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, models.BoletosPath)
	u.Path = path.Join(u.Path, authenticationCode)
	u.Path = path.Join(u.Path, "pdf")
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating request")
		return err
	}

	token, err := b.authentication.Token(ctx)

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

	if resp.StatusCode == http.StatusOK {
		_, err := io.Copy(w, resp.Body)
		if err != nil {
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error writting bytes to writer")
			return errors.ErrDefaultBoletos
		}

		return nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return errors.ErrEntryNotFound
	}

	return errors.ErrDefaultBoletos
}

//CancelBoleto ...
func (b *Boletos) CancelBoleto(ctx context.Context, model *models.CancelBoletoRequest) error {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

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

	u.Path = path.Join(u.Path, models.BoletosPath)
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

	req, err := http.NewRequestWithContext(ctx, "DELETE", endpoint, bytes.NewReader(reqbyte))

	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating request")
		return err
	}

	token, err := b.authentication.Token(ctx)

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

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return errors.ErrDefaultBoletos
}

//SimulatePayment ...
func (b *Boletos) SimulatePayment(ctx context.Context, model *models.SimulatePaymentRequest) error {
	err := grok.Validator.Struct(model)

	if err != nil {
		return grok.FromValidationErros(err)
	}

	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, models.BoletosPath)
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

	var bodyErr []*errors.ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		return err
	}

	if len(bodyErr) > 0 {
		err := bodyErr[0]
		return errors.FindError(err.Code, err.Message)
	}

	return errors.ErrDefaultBoletos
}