package bankly

import (
	"bytes"
	"context"
	"encoding/json"
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
	Session        Session
	HttpClient     *http.Client
	Authentication *Authentication
	errorHandler   ErrorHandler
}

//NewBanklyHttpClient ...
func NewBanklyHttpClient(session Session,
	httpClient *http.Client,
	authentication *Authentication) *BanklyHttpClient {
	return &BanklyHttpClient{
		Session:        session,
		HttpClient:     httpClient,
		Authentication: authentication,
	}
}

func (client *BanklyHttpClient) Post(ctx context.Context, url string, body interface{}, header *http.Header) (*http.Response, error) {
	fields := initLog(ctx)
	data, err := json.Marshal(body)
	if err != nil {
		logErrorWithFields(fields, err, "error marshal body request", nil)
		return nil, err
	}

	endpoint, _ := client.getEndpointAPI(fields, url)

	req, err := http.NewRequestWithContext(ctx, POST, endpoint, bytes.NewReader(data))
	if err != nil {
		logErrorWithFields(fields, err, "error new request", nil)
		return nil, err
	}

	token, err := client.Authentication.Token(ctx)
	if err != nil {
		logErrorWithFields(fields, err, "error authentication", nil)
		return nil, err
	}

	req = setRequestHeader(req, token, client.Session.APIVersion, header)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		logErrorWithFields(fields, err, "error http client", nil)
		return nil, err
	}

	return handleResponse(resp, fields, client.errorHandler)
}

func (client *BanklyHttpClient) Get(ctx context.Context, url string, query map[string]string, header *http.Header) (*http.Response, error) {
	fields := initLog(ctx)

	endpoint, _ := client.getEndpointAPI(fields, url)

	if query != nil {
		endpoint = buildQueryParams(endpoint, query)
	}

	req, err := http.NewRequestWithContext(ctx, GET, endpoint, nil)
	if err != nil {
		logErrorWithFields(fields, err, "error new request", nil)
		return nil, err
	}

	token, err := client.Authentication.Token(ctx)
	if err != nil {
		logErrorWithFields(fields, err, "error authentication", nil)
		return nil, err
	}

	req = setRequestHeader(req, token, client.Session.APIVersion, header)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		logErrorWithFields(fields, err, "error http client", nil)
		return nil, err
	}

	return handleResponse(resp, fields, client.errorHandler)
}

func (client *BanklyHttpClient) Patch(ctx context.Context, url string, body interface{}, header *http.Header) (*http.Response, error) {
	fields := initLog(ctx)
	data, err := json.Marshal(body)
	if err != nil {
		logErrorWithFields(fields, err, "error marshal body request", nil)
		return nil, err
	}

	endpoint, _ := client.getEndpointAPI(fields, url)

	req, err := http.NewRequestWithContext(ctx, PATCH, endpoint, bytes.NewReader(data))
	if err != nil {
		logErrorWithFields(fields, err, "error new request", nil)
		return nil, err
	}

	token, err := client.Authentication.Token(ctx)
	if err != nil {
		logErrorWithFields(fields, err, "error authentication", nil)
		return nil, err
	}

	req = setRequestHeader(req, token, client.Session.APIVersion, header)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		logErrorWithFields(fields, err, "error http client", nil)
		return nil, err
	}

	return handleResponse(resp, fields, client.errorHandler)
}

func handleResponse(resp *http.Response, fields logrus.Fields, handler ErrorHandler) (*http.Response, error) {

	switch {
	case resp.StatusCode == http.StatusOK:
		return resp, nil
	case resp.StatusCode == http.StatusAccepted:
		return resp, nil
	case resp.StatusCode == http.StatusNoContent:
		return resp, nil
	case resp.StatusCode == http.StatusNotFound:
		return nil, ErrEntryNotFound
	}

	if handler != nil {
		return nil, handler(fields, resp)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	return nil, grok.NewError(resp.StatusCode, string(respBody))
}

func (client *BanklyHttpClient) getEndpointAPI(fields logrus.Fields, URLpath string) (string, error) {
	u, err := url.Parse(client.Session.APIEndpoint)
	if err != nil {
		logErrorWithFields(fields, err, "error parsing api endpoint", nil)
		return "", err
	}

	u.Path = path.Join(u.Path, URLpath)
	endpoint := u.String()
	fields["endpoint"] = endpoint
	logInfoWithFields(fields, "get endpoint sucess")
	return endpoint, nil
}

func initLog(ctx context.Context) logrus.Fields {
	requestID, _ := ctx.Value("Request-Id").(string)
	return logrus.Fields{
		"request_id": requestID,
	}
}

func buildQueryParams(endpoint string, query map[string]string) string {
	endpoint = endpoint + "?"
	for key, value := range query {
		endpoint += key + "=" + value + "&"
	}
	return endpoint
}
