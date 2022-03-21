package bankly

import (
	"context"
	"encoding/json"
	errors "github.com/contbank/bankly-sdk/pkg/errors"
	models "github.com/contbank/bankly-sdk/pkg/models"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Pix struct {
	httpClient utils.BanklyHttpClient
}

//NewPix ...
func NewPix(newHttpClient utils.BanklyHttpClient) *Pix {
	newHttpClient.ErrorHandler = PixErrorHandler
	return &Pix{newHttpClient}
}

//GetAddresskey ...
func (p *Pix) GetAddresskey(ctx context.Context, key string, currentIdentity string) (*models.PixAddressKeyResponse, error) {

	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
		"key":        key,
	}

	url := "pix/entries/" + key

	header := http.Header{}
	header.Add("x-bkly-pix-user-id", currentIdentity)

	resp, err := p.httpClient.Get(ctx, url, nil, &header)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	response := new(models.PixAddressKeyResponse)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, errors.ErrDefaultPix
	}

	return response, nil
}

//CashOut ...
func (p *Pix) CashOut(ctx context.Context, pix *models.PixCashOutRequest) (*models.PixCashOutResponse, error) {

	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
		"object":     pix,
	}

	url := "pix/cash-out"

	resp, err := p.httpClient.Post(ctx, url, pix, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	response := new(models.PixCashOutResponse)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, errors.ErrDefaultPix
	}

	return response, nil
}

//QrCodeDecode ...
func (p *Pix) QrCodeDecode(ctx context.Context, encode *models.PixQrCodeDecodeRequest,
	currentIdentity string) (*models.PixQrCodeDecodeResponse, error) {

	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
		"object":     encode,
	}

	url := "pix/qrcodes/decode"

	header := http.Header{}
	header.Add("x-bkly-pix-user-id", currentIdentity)

	resp, err := p.httpClient.Post(ctx, url, encode, &header)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	response := new(models.PixQrCodeDecodeResponse)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, errors.ErrDefaultPix
	}

	return response, nil
}

func (p *Pix) GetCashOutByAuthenticationCode(ctx context.Context, accountNumber string,
	authenticationCode string) (*models.PixCashOutByAuthenticationCodeResponse, error) {

	fields := logrus.Fields{
		"request_id":          utils.GetRequestID(ctx),
		"authentication_code": authenticationCode,
	}

	url := "/pix/cash-out/accounts/" + accountNumber + "/authenticationcode/" + authenticationCode

	header := http.Header{}
	header.Add("x-correlation-id", utils.GetRequestID(ctx))

	resp, err := p.httpClient.Get(ctx, url, nil, &header)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	response := new(models.PixCashOutByAuthenticationCodeResponse)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, errors.ErrDefaultPix
	}

	return response, nil
}

//PixErrorHandler ...
func PixErrorHandler(fields logrus.Fields, resp *http.Response) error {
	var bodyErr *errors.ErrorResponse
	respBody, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return errors.ErrDefaultPix
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err := errors.FindPixError(errModel.Code, errModel.Messages...)

		fields["bankly_error"] = bodyErr
		logrus.WithFields(fields).WithError(err).Error("bankly get pix error")

		return err
	}
	return errors.ErrDefaultPix
}
