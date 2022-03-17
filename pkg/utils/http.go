package bankly

import (
	"bytes"
	"context"
	"encoding/json"
	errors "github.com/contbank/bankly-sdk/pkg/errors"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"

	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

type ErrorHandler func(fields logrus.Fields, resp *http.Response) error

const (
	GET   = "GET"
	POST  = "POST"
	PATCH = "PATCH"
)

type BanklyHttpClient struct {
	Session    	   bankly.Session
	HttpClient     *http.Client
	Authentication *bankly.Authentication
	ErrorHandler   ErrorHandler
}

// NewBanklyHttpClient ...
func NewBanklyHttpClient(session bankly.Session,
	httpClient *http.Client,
	authentication *bankly.Authentication) *BanklyHttpClient {
	return &BanklyHttpClient{
		Session:        session,
		HttpClient:     httpClient,
		Authentication: authentication,
	}
}

// Post ...
func (client *BanklyHttpClient) Post(ctx context.Context, url string, body interface{}, header *http.Header) (*http.Response, error) {
	fields := initLog(ctx)
	data, err := json.Marshal(body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error marshal body request")
		return nil, err
	}

	endpoint, _ := client.getEndpointAPI(fields, url)

	req, err := http.NewRequestWithContext(ctx, POST, endpoint, bytes.NewReader(data))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error new request")
		return nil, err
	}

	token, err := client.Authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error authentication")
		return nil, err
	}

	req = SetRequestHeader(req, token, client.Session.APIVersion, header)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error http client")
		return nil, err
	}

	return handleResponse(resp, fields, client.ErrorHandler)
}

// Get ...
func (client *BanklyHttpClient) Get(ctx context.Context, url string, query map[string]string, header *http.Header) (*http.Response, error) {
	fields := initLog(ctx)

	endpoint, _ := client.getEndpointAPI(fields, url)

	if query != nil {
		endpoint = buildQueryParams(endpoint, query)
	}

	req, err := http.NewRequestWithContext(ctx, GET, endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error new request")
		return nil, err
	}

	token, err := client.Authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error authentication")
		return nil, err
	}

	req = SetRequestHeader(req, token, client.Session.APIVersion, header)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error http client")
		return nil, err
	}

	return handleResponse(resp, fields, client.ErrorHandler)
}

// Patch ...
func (client *BanklyHttpClient) Patch(ctx context.Context, url string, body interface{}, header *http.Header) (*http.Response, error) {
	fields := initLog(ctx)
	data, err := json.Marshal(body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error marshal body request")
		return nil, err
	}

	endpoint, _ := client.getEndpointAPI(fields, url)

	req, err := http.NewRequestWithContext(ctx, PATCH, endpoint, bytes.NewReader(data))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error new request")
		return nil, err
	}

	token, err := client.Authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error authentication")
		return nil, err
	}

	req = SetRequestHeader(req, token, client.Session.APIVersion, header)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error http client")
		return nil, err
	}

	return handleResponse(resp, fields, client.ErrorHandler)
}

// handleResponse ...
func handleResponse(resp *http.Response, fields logrus.Fields, handler ErrorHandler) (*http.Response, error) {

	switch {
	case resp.StatusCode == http.StatusOK:
		return resp, nil
	case resp.StatusCode == http.StatusAccepted:
		return resp, nil
	case resp.StatusCode == http.StatusNoContent:
		return resp, nil
	case resp.StatusCode == http.StatusNotFound:
		return nil, errors.ErrEntryNotFound
	case resp.StatusCode == http.StatusForbidden:
		return nil, errors.ErrServiceForbidden
	case resp.StatusCode == http.StatusGatewayTimeout:
		return nil, errors.ErrGatewayTimeout
	}

	if handler != nil {
		return nil, handler(fields, resp)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	return nil, grok.NewError(resp.StatusCode, string(respBody))
}

// getEndpointAPI ...
func (client *BanklyHttpClient) getEndpointAPI(fields logrus.Fields, URLpath string) (string, error) {
	u, err := url.Parse(client.Session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error parsing api endpoint")
		return "", err
	}

	u.Path = path.Join(u.Path, URLpath)
	endpoint := u.String()
	fields["endpoint"] = endpoint
	logrus.WithFields(fields).Info("get endpoint success")
	return endpoint, nil
}

// initLog ...
func initLog(ctx context.Context) logrus.Fields {
	return logrus.Fields{
		"request_id": GetRequestID(ctx),
	}
}

// buildQueryParams ...
func buildQueryParams(endpoint string, query map[string]string) string {
	endpoint = endpoint + "?"
	for key, value := range query {
		endpoint += key + "=" + value + "&"
	}
	return endpoint
}