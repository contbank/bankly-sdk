package bankly

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type ErrorHandler func(log *logrus.Entry, resp *http.Response) error

type TokenProvider interface {
	Token(ctx context.Context) (string, error)
}

type BanklyHttpClient interface {
	NewRequest(ctx context.Context, method string, url string, body interface{}, query map[string]string, header *http.Header) (*http.Request, error)
	Request(ctx context.Context, method string, url string, body interface{}, query map[string]string, header *http.Header) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
	Post(ctx context.Context, url string, body interface{}, header *http.Header) (*http.Response, error)
	Delete(ctx context.Context, url string, body interface{}, header *http.Header) (*http.Response, error)
	Patch(ctx context.Context, url string, body interface{}, query map[string]string, header *http.Header) (*http.Response, error)
	Put(ctx context.Context, url string, body interface{}, header *http.Header) (*http.Response, error)
	Get(ctx context.Context, url string, query map[string]string, header *http.Header) (*http.Response, error)
	SetErrorHandler(handler ErrorHandler)
}

type apiClient struct {
	Session        Session
	HttpClient     *http.Client
	Authentication TokenProvider
	errorHandler   ErrorHandler
}

func (c *apiClient) SetErrorHandler(handler ErrorHandler) {
	c.errorHandler = handler
}

//NewBanklyHttpClient ...
func NewBanklyHttpClient(session Session,
	httpClient *http.Client,
	authentication TokenProvider) BanklyHttpClient {
	return &apiClient{
		Session:        session,
		HttpClient:     httpClient,
		Authentication: authentication,
	}
}

func (c *apiClient) NewRequest(ctx context.Context, method string, url string, body interface{}, query map[string]string, header *http.Header) (*http.Request, error) {
	var bodyReader io.Reader
	log := logrus.WithFields(initLog(ctx))
	endpoint, err := c.getEndpointAPI(log, url)
	if err != nil {
		return nil, err
	}

	if query != nil {
		endpoint = buildQueryParams(endpoint, query)
	}

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			log.WithError(err).Error("error marshal body request")
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bodyReader)
	if err != nil {
		log.WithError(err).Error("error new request")
		return nil, err
	}

	token, err := c.Authentication.Token(ctx)
	if err != nil {
		log.WithError(err).Error("error authentication")
		return nil, err
	}

	req = setRequestHeader(req, token, c.Session.APIVersion, header)
	return req, nil
}

func (c *apiClient) Request(ctx context.Context, method string, url string, body interface{}, query map[string]string, header *http.Header) (*http.Response, error) {
	log := logrus.WithFields(initLog(ctx))
	req, err := c.NewRequest(ctx, method, url, body, query, header)
	if err != nil {
		log.WithError(err).Error("error http client")
		return nil, err
	}
	return c.Do(req)
}

func (c *apiClient) Do(req *http.Request) (*http.Response, error) {
	log := logrus.WithFields(initLog(req.Context()))
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("error http client")
		return nil, err
	}

	return handleResponse(resp, log, c.errorHandler)
}

func (c *apiClient) Post(ctx context.Context, url string, body interface{}, header *http.Header) (*http.Response, error) {
	return c.Request(ctx, http.MethodPost, url, body, nil, header)
}

func (c *apiClient) Delete(ctx context.Context, url string, body interface{}, header *http.Header) (*http.Response, error) {
	return c.Request(ctx, http.MethodDelete, url, body, nil, header)

}

func (c *apiClient) Patch(ctx context.Context, url string, body interface{}, query map[string]string, header *http.Header) (*http.Response, error) {
	return c.Request(ctx, http.MethodPatch, url, body, query, header)
}

func (c *apiClient) Put(ctx context.Context, url string, body interface{}, header *http.Header) (*http.Response, error) {
	return c.Request(ctx, http.MethodPut, url, body, nil, header)
}

func (c *apiClient) Get(ctx context.Context, url string, query map[string]string, header *http.Header) (*http.Response, error) {
	return c.Request(ctx, http.MethodGet, url, nil, query, header)
}

func handleResponse(resp *http.Response, log *logrus.Entry, handler ErrorHandler) (*http.Response, error) {

	if resp != nil {
		log.WithField("http response", resp.StatusCode).
			Info("handle response - status code")
	}

	switch {
	case resp.StatusCode == http.StatusOK:
		return resp, nil
	case resp.StatusCode == http.StatusAccepted:
		return resp, nil
	case resp.StatusCode == http.StatusNoContent:
		return resp, nil
	case resp.StatusCode == http.StatusCreated:
		return resp, nil
	case resp.StatusCode == http.StatusNotFound:
		return nil, ErrEntryNotFound
	}

	if handler != nil {
		return nil, handler(log, resp)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	return nil, grok.NewError(resp.StatusCode, "DEFAULT_ERROR", string(respBody))
}

func (c *apiClient) getEndpointAPI(log *logrus.Entry, relativePath string) (string, error) {
	u, err := url.Parse(c.Session.APIEndpoint)
	if err != nil {
		log.WithError(err).Error("error parsing api endpoint")
		return "", err
	}

	u.Path = path.Join(u.Path, relativePath)
	endpoint := u.String()
	log.WithField("endpoint", endpoint).Info("get endpoint success")
	return endpoint, nil
}

func initLog(ctx context.Context) logrus.Fields {
	requestID, _ := ctx.Value("Request-Id").(string)
	return logrus.Fields{
		"request_id": requestID,
	}
}

func buildQueryParams(endpoint string, queryParams map[string]string) string {
	if len(queryParams) == 0 {
		return endpoint
	}
	query := url.Values{}
	for key, value := range queryParams {
		query.Set(key, value)
	}
	return endpoint + "?" + query.Encode()
}
