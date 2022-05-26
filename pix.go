package bankly

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Pix struct {
	httpClient BanklyHttpClient
}

//NewPix ...
func NewPix(newHttpClient BanklyHttpClient) *Pix {
	newHttpClient.errorHandler = PixErrorHandler
	return &Pix{newHttpClient}
}

//GetAddressKey ...
func (p *Pix) GetAddressKey(ctx context.Context, key string, currentIdentity string) (*PixAddressKeyResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
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

	response := new(PixAddressKeyResponse)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	return response, nil
}

//CashOut ...
func (p *Pix) CashOut(ctx context.Context, pix *PixCashOutRequest) (*PixCashOutResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
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

	response := new(PixCashOutResponse)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	return response, nil
}

//QrCodeDecode ...
func (p *Pix) QrCodeDecode(ctx context.Context, encode *PixQrCodeDecodeRequest, currentIdentity string) (*PixQrCodeDecodeResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
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

	response := new(PixQrCodeDecodeResponse)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	return response, nil
}

func (p *Pix) GetCashOutByAuthenticationCode(ctx context.Context, accountNumber string, authenticationCode string) (*PixCashOutByAuthenticationCodeResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id":          requestID,
		"authentication_code": authenticationCode,
	}

	url := "/pix/cash-out/accounts/" + accountNumber + "/authenticationcode/" + authenticationCode

	header := http.Header{}
	header.Add("x-correlation-id", requestID)

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

	response := new(PixCashOutByAuthenticationCodeResponse)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	return response, nil
}

func (p *Pix) CreateAddressKey(ctx context.Context, pix PixAddressKeyCreateRequest) (*PixAddressKeyCreateResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)

	if requestID == "" {
		ctx = GenerateNewRequestID(ctx)
	} else {
		ctx = context.WithValue(ctx, "Request-Id", ctx.Value("Request-Id").(string))
	}

	requestID = GetRequestID(ctx)

	fields := logrus.Fields{
		"request_id": requestID,
		"object":     pix,
	}
	url := "/pix/entries"

	header := http.Header{}
	header.Add("x-correlation-id", requestID)

	resp, err := p.httpClient.Post(ctx, url, pix, &header)
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

	response := new(PixAddressKeyCreateResponse)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	return response, nil
}

func (p *Pix) DeleteAddressKey(ctx context.Context, addressingKey string) (*http.Response, error) {
	requestID, _ := ctx.Value("Request-Id").(string)

	if requestID == "" {
		ctx = GenerateNewRequestID(ctx)
	} else {
		ctx = context.WithValue(ctx, "Request-Id", ctx.Value("Request-Id").(string))
	}

	requestID = GetRequestID(ctx)

	fields := logrus.Fields{
		"request_id":    requestID,
		"addressingKey": addressingKey,
	}

	header := http.Header{}
	header.Add("x-correlation-id", requestID)
	header.Add("x-bkly-pix-user-id", addressingKey)

	url := fmt.Sprintf("/pix/entries/%s", addressingKey)

	resp, err := p.httpClient.Delete(ctx, url, addressingKey, &header)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}

//PixErrorHandler ...
func PixErrorHandler(fields logrus.Fields, resp *http.Response) error {
	var bodyErr *ErrorResponse
	respBody, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return ErrDefaultPix
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err := FindPixError(errModel.Code, errModel.Messages...)

		fields["bankly_error"] = bodyErr
		logrus.WithFields(fields).WithError(err).Error("bankly get pix error")

		return err
	}
	return ErrDefaultPix
}