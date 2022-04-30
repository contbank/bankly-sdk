package bankly

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

//Client ...
type Client struct {
	config     Config
	httpClient *http.Client
}

//NewClient ...
func NewClient(config Config, httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		config:     config,
	}
}

func (a *Client) Register(ctx context.Context) (*ClientRegisterResponse, error) {
	u, err := url.Parse(*a.config.LoginEndpoint)

	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, ClientPath)
	endpoint := u.String()

	banklyRequest := &ClientRegisterBanklyRequest{
		GrantTypes: []string{
			"client_credentials",
		},
		ResponseTypes: []string{
			"access_token",
		},
		TokenEndpointAuthMethod: "tls_client_auth",
		TLSClientAuthSubjectDn:  a.config.Certificate.SubjectDn,
		CompanyKey:              *a.config.CompanyKey,
		Scope:                   *a.config.Scopes,
	}

	reqbyte, err := json.Marshal(banklyRequest)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		return nil, err
	}

	resp, err := a.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var response *ClientRegisterResponse

		respBody, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			return nil, err
		}

		return response, nil
	}

	return nil, ErrDefaultLogin
}
