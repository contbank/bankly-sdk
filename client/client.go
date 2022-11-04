package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/contbank/bankly-sdk"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type TokenProvider interface {
	Token(ctx context.Context) (string, error)
}

type Client interface {
	Request(ctx context.Context, method string, endpoint string, payload interface{}) (*BanklyResponse, error)
	Do(request *http.Request) (*BanklyResponse, error)
	NewRequest(ctx context.Context, method string, endpoint string, payload interface{}) (*http.Request, error)
	Get(ctx context.Context, endpoint string) (*BanklyResponse, error)
	Post(ctx context.Context, endpoint string, payload interface{}) (*BanklyResponse, error)
	Put(ctx context.Context, endpoint string, payload interface{}) (*BanklyResponse, error)
	Patch(ctx context.Context, endpoint string, payload interface{}) (*BanklyResponse, error)
}

type client struct {
	session        bankly.Session
	httpClient     *http.Client
	authentication TokenProvider
}

// NewClient ...
func NewClient(httpClient *http.Client, session bankly.Session) Client {
	return &client{
		session:        session,
		httpClient:     httpClient,
		authentication: bankly.NewAuthentication(httpClient, session),
	}
}

type BanklyResponse struct {
	*http.Response
	Body []byte
}

func (r BanklyResponse) Json(object interface{}) error {
	return json.Unmarshal(r.Body, object)
}

func (c client) Request(ctx context.Context, method string, endpoint string, payload interface{}) (*BanklyResponse, error) {
	request, err := c.NewRequest(ctx, method, endpoint, payload)
	if err != nil {
		return nil, err
	}
	return c.Do(request)
}

func (c client) Do(request *http.Request) (*BanklyResponse, error) {
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.WithError(err).Error("error closing response body")
		}
	}(response.Body)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &BanklyResponse{
		Response: response,
		Body:     body,
	}, nil
}

func (c client) NewRequest(ctx context.Context, method string, endpoint string, payload interface{}) (*http.Request, error) {
	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.WithError(err).Error("error parsing api endpoint")
		return nil, err
	}
	u.Path = path.Join(u.Path, endpoint)
	banklyEndpoint := u.String()

	var requestBytes []byte
	if payload != nil {
		requestBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequestWithContext(ctx, method, banklyEndpoint, bytes.NewReader(requestBytes))
	if err != nil {
		return nil, err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", token)
	request.Header.Add("Content-type", "application/json")
	request.Header.Add("api-version", c.session.APIVersion)

	return request, nil
}

func (c client) Get(ctx context.Context, endpoint string) (*BanklyResponse, error) {
	return c.Request(ctx, http.MethodGet, endpoint, nil)
}

func (c client) Post(ctx context.Context, endpoint string, payload interface{}) (*BanklyResponse, error) {
	return c.Request(ctx, http.MethodPost, endpoint, payload)
}

func (c client) Put(ctx context.Context, endpoint string, payload interface{}) (*BanklyResponse, error) {
	return c.Request(ctx, http.MethodPut, endpoint, payload)
}

func (c client) Patch(ctx context.Context, endpoint string, payload interface{}) (*BanklyResponse, error) {
	return c.Request(ctx, http.MethodPatch, endpoint, payload)
}
