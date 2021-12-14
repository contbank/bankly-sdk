package bankly

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/sirupsen/logrus"
)

const (
	GET  = "GET"
	POST = "POST"
)

type NewHttpClient struct {
	Session        Session
	HttpClient     *http.Client
	Authentication *Authentication
}

func (client *NewHttpClient) Post(ctx context.Context, method string, url string, body interface{}) (*http.Response, error) {
	fields := initLog(ctx)
	data, err := json.Marshal(body)
	if err != nil {
		logErrorWithFields(fields, err, "error marshal body request", nil)
		return nil, err
	}

	endpoint, _ := client.getEndpointAPI(fields, url)

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(data))
	if err != nil {
		logErrorWithFields(fields, err, "error new request", nil)
		return nil, err
	}

	token, err := client.Authentication.Token(ctx)
	if err != nil {
		logErrorWithFields(fields, err, "error authentication", nil)
		return nil, err
	}

	req = setRequestHeader(req, token, client.Session.APIVersion)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		logErrorWithFields(fields, err, "error http client", nil)
		return nil, err
	}

	return handleResponse(resp, fields)
}

func (client *NewHttpClient) Get(ctx context.Context, method string, url string) (*http.Response, error) {
	fields := initLog(ctx)

	endpoint, _ := client.getEndpointAPI(fields, url)

	req, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		logErrorWithFields(fields, err, "error new request", nil)
		return nil, err
	}

	token, err := client.Authentication.Token(ctx)
	if err != nil {
		logErrorWithFields(fields, err, "error authentication", nil)
		return nil, err
	}

	req = setRequestHeader(req, token, client.Session.APIVersion)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		logErrorWithFields(fields, err, "error http client", nil)
		return nil, err
	}

	return handleResponse(resp, fields)
}

func handleResponse(resp *http.Response, fields logrus.Fields) (*http.Response, error) {

	switch {
	case resp.StatusCode == http.StatusOK:
		return resp, nil
	case resp.StatusCode == http.StatusAccepted:
		return resp, nil
	case resp.StatusCode == http.StatusNotFound:
		return nil, ErrEntryNotFound
	}

	return responseIsError(fields, resp)
}

func responseIsError(fields logrus.Fields, resp *http.Response) (*http.Response, error) {
	var bodyErr *ErrorResponse
	respBody, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding json response", nil)
		return nil, ErrDefaultCard
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err := FindCardError(errModel.Code, errModel.Messages...)

		var hasField = make(map[string]interface{})
		hasField["bankly_error"] = bodyErr
		logErrorWithFields(fields, err, "bankly get card error", hasField)

		return nil, err
	}
	return nil, ErrDefaultCard
}

func (client *NewHttpClient) getEndpointAPI(fields logrus.Fields, URLpath string) (string, error) {
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
