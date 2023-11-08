package bankly

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/contbank/grok"

	"github.com/sirupsen/logrus"
)

type Pix struct {
	httpClient BanklyHttpClient
}

// NewPix ...
func NewPix(newHttpClient BanklyHttpClient) *Pix {
	newHttpClient.SetErrorHandler(PixErrorHandler)
	return &Pix{newHttpClient}
}

// GetAddressKeysByAccount ...
func (p *Pix) GetAddressKeysByAccount(ctx context.Context, accountNumber string, currentIdentity string) ([]*PixTypeValue, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id":       requestID,
		"account":          accountNumber,
		"current_identity": currentIdentity,
	}

	url := "accounts/" + grok.OnlyDigits(accountNumber) + "/addressing-keys"

	header := http.Header{}
	header.Add("x-bkly-pix-user-id", grok.OnlyDigits(currentIdentity))
	header.Add("x-correlation-id", requestID)

	resp, err := p.httpClient.Get(ctx, url, nil, &header)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	response := []*PixTypeValue{}

	if resp.StatusCode == http.StatusNoContent {
		logrus.WithFields(fields).Info("no data found")
		return response, nil
	}

	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	return response, nil
}

// GetAddressKey ...
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

// CashOut ...
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

// QrCodeDynamic ...
func (p *Pix) QrCodeDynamic(ctx context.Context, data *PixQrCodeDynamicRequest,
	currentIdentity string) (*PixQrCodeResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	log := logrus.WithFields(logrus.Fields{
		"request_id": requestID,
		"object":     data,
	})

	url := "pix/qrcodes/dynamic/payment"

	header := http.Header{}
	header.Add("x-bkly-pix-user-id", currentIdentity)
	header.Add("x-correlation-id", requestID)

	resp, err := p.httpClient.Post(ctx, url, data, &header)
	if err != nil {
		log.WithError(err).Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("error decoding body response")
		return nil, err
	}

	result := new(PixQrCodeResponse)
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		log.WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	return result, nil
}

// QrCodeDecode ...
func (p *Pix) QrCodeDecode(ctx context.Context, encode *PixQrCodeDecodeRequest,
	currentIdentity string) (*PixQrCodeDecodeResponse, error) {
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

// GetCashOutByAuthenticationCode ...
func (p *Pix) GetCashOutByAuthenticationCode(ctx context.Context, accountNumber string,
	authenticationCode string) (*PixCashOutByAuthenticationCodeResponse, error) {

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

// CreateAddressKey ...
func (p *Pix) CreateAddressKey(ctx context.Context, pix *PixAddressKeyCreateRequest) (*PixAddressKeyCreateResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)

	if requestID == "" {
		ctx = GenerateNewRequestID(ctx)
	} else {
		ctx = context.WithValue(ctx, "Request-Id", ctx.Value("Request-Id").(string))
	}

	requestID = GetRequestID(ctx)

	fields := logrus.Fields{
		"request_id": requestID,
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

// DeleteAddressKey ...
func (p *Pix) DeleteAddressKey(ctx context.Context, identifier, addressingKey string) error {
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
	header.Add("x-bkly-pix-user-id", identifier)

	url := fmt.Sprintf("/pix/entries/%s", addressingKey)

	resp, err := p.httpClient.Delete(ctx, url, addressingKey, &header)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return err
	}

	defer resp.Body.Close()

	return nil
}

// PixErrorHandler ...
func PixErrorHandler(log *logrus.Entry, resp *http.Response) error {
	var bodyErr *ErrorResponse
	respBody, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		log.WithError(err).Error("error decoding json response")
		return ErrDefaultPix
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err := FindPixError(errModel.Code, errModel.Messages...)

		log.WithField("bankly_error", bodyErr).WithError(err).Error("bankly get pix error")

		return err
	}
	return ErrDefaultPix
}

// CreatePixClaim ...
func (p *Pix) CreatePixClaim(ctx context.Context, pix *PixClaimRequest, documentNumber string) (*PixClaimResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	
	if requestID == "" {
		ctx = GenerateNewRequestID(ctx)
	} else {
		ctx = context.WithValue(ctx, "Request-Id", ctx.Value("Request-Id").(string))
	}

	requestID = GetRequestID(ctx)

	fields := logrus.Fields{
		"request_id": requestID,
		"x-bkly-pix-user-id": documentNumber,
	}

	url := "/pix/claims"

	header := http.Header{}
	header.Add("x-bkly-pix-user-id", grok.OnlyDigits(documentNumber))

	resp, err := p.httpClient.Post(ctx, url, pix, &header)
	if err != nil {
		log.Println("resposta -> ", err)
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	response := new(PixClaimResponse)
	log.Println(response)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	return response, nil
}

// GetPixClaim ...
func (p *Pix) GetPixClaim(ctx context.Context, accountNumber string, documentNumber string, claimsFrom string) ([]*PixClaimResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id":       requestID,
		"account":          accountNumber,
		"current_identity": documentNumber,
	}

	url := fmt.Sprintf("/pix/claims?documentNumber=%v&claimsFrom=%v", documentNumber, claimsFrom)

	header := http.Header{}
	header.Add("x-bkly-pix-user-id", grok.OnlyDigits(documentNumber))
	header.Add("x-correlation-id", requestID)

	resp, err := p.httpClient.Get(ctx, url, nil, &header)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	response := []*PixClaimResponse{}

	if resp.StatusCode == http.StatusNoContent {
		logrus.WithFields(fields).Info("no data found")
		return response, nil
	}

	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	return response, nil
}

// ConfirmPixClaim ...
func (p *Pix) ConfirmPixClaim(ctx context.Context, documentNumber string, claimId string, reason *PixClaimConfirmReason) (*PixClaimConfirmResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id":       requestID,
		"current_identity": documentNumber,
	}

	url := fmt.Sprintf("/pix/claims/%v/confirm", claimId)

	header := http.Header{}
	header.Add("x-bkly-pix-user-id", grok.OnlyDigits(documentNumber))

	resp, err := p.httpClient.Patch(ctx, url, reason, nil, &header)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	response := new(PixClaimConfirmResponse)

	if resp.StatusCode == http.StatusNoContent {
		logrus.WithFields(fields).Info("no data found")
		return response, nil
	}

	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	return response, nil
}

// CompletePixClaim ...
func (p *Pix) CompletePixClaim(ctx context.Context, documentNumber string, claimId string) (*PixClaimCompleteResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id":       requestID,
		"current_identity": documentNumber,
	}

	url := fmt.Sprintf("/pix/claims/%v/complete", claimId)

	header := http.Header{}
	header.Add("x-bkly-pix-user-id", grok.OnlyDigits(documentNumber))

	resp, err := p.httpClient.Patch(ctx, url, nil, nil, &header)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	response := new(PixClaimCompleteResponse)

	if resp.StatusCode == http.StatusNoContent {
		logrus.WithFields(fields).Info("no data found")
		return response, nil
	}

	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	return response, nil
}

// CancelPixClaim ...
func (p *Pix) CancelPixClaim(ctx context.Context, documentNumber string, claimId string, reason *PixClaimConfirmReason) (*PixClaimCancelResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id":       requestID,
		"current_identity": documentNumber,
	}

	url := fmt.Sprintf("/pix/claims/%v/cancel", claimId)

	header := http.Header{}
	header.Add("x-bkly-pix-user-id", grok.OnlyDigits(documentNumber))

	resp, err := p.httpClient.Patch(ctx, url, reason, nil, &header)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	response := new(PixClaimCancelResponse)

	if resp.StatusCode == http.StatusNoContent {
		logrus.WithFields(fields).Info("no data found")
		return response, nil
	}

	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	return response, nil
}