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
	session    Session
	httpClient *http.Client
}

//NewClient ...
func NewClient(httpClient *http.Client, session Session) *Client {
	return &Client{
		session:    session,
		httpClient: httpClient,
	}
}

func (a *Client) Register(ctx context.Context, request *ClientRegisterRequest) (*Session, error) {
	u, err := url.Parse(a.session.LoginEndpoint)

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
		TLSClientAuthSubjectDn:  request.TLSClientAuthSubjectDn,
		CompanyKey:              request.CompanyKey,
		Scope:                   a.session.Scopes,
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

		a.session.ClientID = response.ClientID

		return &a.session, nil
	}

	return nil, ErrDefaultLogin
}
