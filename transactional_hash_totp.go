package bankly

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type TransactionalHashTOTP struct {
	httpClient BanklyHttpClient
}

//NewTransactionalHashTOTP ...
func NewTransactionalHashTOTP(newHttpClient BanklyHttpClient) *TransactionalHashTOTP {
	return &TransactionalHashTOTP{newHttpClient}
}

//TransactionalHash...
func (p *TransactionalHashTOTP) TransactionalHash(ctx context.Context, transactional TransactionalHashRequest, identifier string) (*TransactionalHash, error) {
	requestID, _ := ctx.Value("Request-Id").(string)

	if requestID == "" {
		ctx = GenerateNewRequestID(ctx)
	} else {
		ctx = context.WithValue(ctx, "Request-Id", ctx.Value("Request-Id").(string))
	}

	requestID = GetRequestID(ctx)

	fields := logrus.Fields{
		"request_id": requestID,
		"identifier": identifier,
		"object":     transactional,
	}
	url := "/totp"

	header := http.Header{}
	header.Add("x-bkly-user-id", identifier)

	resp, err := p.httpClient.Post(ctx, url, transactional, &header)
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

	response := new(TransactionalHash)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	return response, nil
}

//TransactionalHashValidate...
func (p *TransactionalHashTOTP) TransactionalHashValidate(ctx context.Context, transactional TransactionalHash, identifier string) (*TransactionalHashValidateResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)

	if requestID == "" {
		ctx = GenerateNewRequestID(ctx)
	} else {
		ctx = context.WithValue(ctx, "Request-Id", ctx.Value("Request-Id").(string))
	}

	requestID = GetRequestID(ctx)

	fields := logrus.Fields{
		"request_id": requestID,
		"object":     transactional,
	}
	url := "/totp"

	header := http.Header{}
	header.Add("x-bkly-user-id", identifier)

	resp, err := p.httpClient.Patch(ctx, url, transactional, nil, &header)
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

	response := new(TransactionalHashValidateResponse)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	return response, nil
}
