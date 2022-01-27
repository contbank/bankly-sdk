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
	GET   = "GET"
	POST  = "POST"
	PATCH = "PATCH"
)

type NewHttpClient struct {
	Session        Session
	HttpClient     *http.Client
	Authentication *Authentication
}

func (client *NewHttpClient) Post(ctx context.Context, url string, body interface{}) (*http.Response, error) {
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

	req = setRequestHeader(req, token, client.Session.APIVersion)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error http client")
		return nil, err
	}

	return handleResponse(resp, fields)
}

func (client *NewHttpClient) Get(ctx context.Context, url string, query map[string]string) (*http.Response, error) {
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

	req = setRequestHeader(req, token, client.Session.APIVersion)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error http client")
		return nil, err
	}

	return handleResponse(resp, fields)
}

func (client *NewHttpClient) Patch(ctx context.Context, url string, body interface{}) (*http.Response, error) {
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

	req = setRequestHeader(req, token, client.Session.APIVersion)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error http client")
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
	case resp.StatusCode == http.StatusNoContent:
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
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultCard
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err := FindCardError(errModel.Code, errModel.Messages...)

		var hasField = make(map[string]interface{})
		hasField["bankly_error"] = bodyErr
		logrus.WithFields(fields).WithError(err).Error("error bankly get card")

		return nil, err
	}
	return nil, ErrDefaultCard
}

func (client *NewHttpClient) getEndpointAPI(fields logrus.Fields, URLpath string) (string, error) {
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
