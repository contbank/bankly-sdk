package bankly

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

//Authentication ...
type Authentication struct {
	session    Session
	httpClient *http.Client
}

//NewAuthentication ...
func NewAuthentication(httpClient *http.Client, session Session) *Authentication {
	return &Authentication{
		session:    session,
		httpClient: httpClient,
	}
}

func (a *Authentication) login(ctx context.Context) (*AuthenticationResponse, error) {
	u, err := url.Parse(a.session.LoginEndpoint)

	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, LoginPath)
	endpoint := u.String()

	formData := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {a.session.ClientID},
		"client_secret": {a.session.ClientSecret},
		"scope":         {a.session.Scope},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(formData.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-type", "application/x-www-form-urlencoded")

	resp, err := a.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var response *AuthenticationResponse

		respBody, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			return nil, err
		}

		return response, nil
	}
	if resp.StatusCode == http.StatusBadRequest {
		var bodyErr *ErrorLoginResponse

		respBody, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(respBody, &bodyErr)

		if err != nil {
			return nil, err
		}

		return nil, FindError("400", bodyErr.Message)
	}

	return nil, ErrDefaultLogin
}

//Token ...
func (a Authentication) Token(ctx context.Context) (string, error) {
	if token, found := a.session.Cache.Get("token"); found {
		return token.(string), nil
	}

	response, err := a.login(ctx)

	if err != nil {
		return "", err
	}

	a.session.Cache.Set("token", fmt.Sprintf("%s %s", response.TokenType, response.AccessToken), time.Second*time.Duration(int64(response.ExpiresIn-10)))

	return fmt.Sprintf("%s %s", response.TokenType, response.AccessToken), nil
}
