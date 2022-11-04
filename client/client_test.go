package client

import (
	"context"
	"github.com/contbank/bankly-sdk"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestHttpClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

type TestModel struct {
	Status string `json:"status"`
}

type MockToken struct {
	TheToken string
	Error    error
}

func (t MockToken) Token(ctx context.Context) (string, error) {
	return t.TheToken, t.Error
}

func newTestClient(httpClient *http.Client, tokenProvider TokenProvider) client {
	testCache := cache.New(cache.NoExpiration, cache.NoExpiration)
	return client{
		session: bankly.Session{
			LoginEndpoint: "http://test/login",
			APIEndpoint:   "http://test/",
			ClientID:      "ClientID",
			ClientSecret:  "ClientSecret",
			APIVersion:    "APIVersion",
			Cache:         *testCache,
			Scopes:        "Scopes",
			Mtls:          false,
		},
		httpClient:     httpClient,
		authentication: tokenProvider,
	}
}

func TestClient_Do(t *testing.T) {
	body := `{"status": "ok"}`
	request, err := http.NewRequest(http.MethodPost, "http://test/endpoint", ioutil.NopCloser(strings.NewReader(body)))
	require.Nil(t, err)
	var actualRequest *http.Request
	httpClient := NewTestHttpClient(func(req *http.Request) *http.Response {
		actualRequest = req
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(body)),
		}
	})
	testClient := newTestClient(httpClient, MockToken{TheToken: "token"})
	banklyResponse, err := testClient.Do(request)
	require.Nil(t, err)
	requestBody, _ := ioutil.ReadAll(request.Body)
	assert.Equal(t, *request, *actualRequest)
	assert.Equal(t, `{"status": "ok"}`, string(requestBody))
	assert.Equal(t, banklyResponse.StatusCode, http.StatusOK)
	assert.Equal(t, banklyResponse.Body, []byte(body))
}

func TestClient_NewRequest(t *testing.T) {
	ctx := context.Background()
	testClient := newTestClient(&http.Client{}, MockToken{TheToken: "token"})
	request, err := testClient.NewRequest(ctx, http.MethodPost, "/endpoint/resource", TestModel{"ok"})
	require.Nil(t, err)
	assert.Equal(t, "http://test/endpoint/resource", request.URL.String())
	assert.Equal(t, request.Header.Get("Authorization"), "token")
	assert.Equal(t, request.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, request.Header.Get("api-version"), "APIVersion")
}

func TestClient_Request(t *testing.T) {
	payload := TestModel{"ok"}
	ctx := context.Background()
	body := `{"status":"ok"}`
	var actualRequest *http.Request
	httpClient := NewTestHttpClient(func(req *http.Request) *http.Response {
		actualRequest = req
		return &http.Response{Body: ioutil.NopCloser(strings.NewReader(body))}
	})
	testClient := newTestClient(httpClient, MockToken{TheToken: "token"})

	_, err := testClient.Request(ctx, http.MethodPost, "/endpoint", payload)
	require.Nil(t, err)
	assert.Equal(t, http.MethodPost, actualRequest.Method)
	assert.Equal(t, "http://test/endpoint", actualRequest.URL.String())
	actualRequestBody, _ := ioutil.ReadAll(actualRequest.Body)
	assert.Equal(t, body, string(actualRequestBody))

	_, err = testClient.Post(ctx, "/endpoint", payload)
	assert.Equal(t, http.MethodPost, actualRequest.Method)

	_, err = testClient.Patch(ctx, "/endpoint", payload)
	assert.Equal(t, http.MethodPatch, actualRequest.Method)

	_, err = testClient.Put(ctx, "/endpoint", payload)
	assert.Equal(t, http.MethodPut, actualRequest.Method)

	_, err = testClient.Get(ctx, "/endpoint")
	assert.Equal(t, http.MethodGet, actualRequest.Method)
	actualRequestBody, _ = ioutil.ReadAll(actualRequest.Body)
	assert.Equal(t, "", string(actualRequestBody))
}
