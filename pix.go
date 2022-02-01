package bankly

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Pix struct {
	httpClient BanklyHttpClient
}

//NewPix ...
func NewPix(newHttpClient BanklyHttpClient) *Pix {
	return &Pix{newHttpClient}
}

//GetAddresskey ...
func (p *Pix) GetAddresskey(ctx context.Context, key string, currentIdentity string) (*PixAddressKeyResponse, error) {
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
func (p *Pix) QrCodeDecode(ctx context.Context, encode *PixQrCodeDecodeRequest) (*PixQrCodeDecodeResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"object":     encode,
	}

	url := "pix/qrcodes/decode"

	resp, err := p.httpClient.Post(ctx, url, encode, nil)
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
