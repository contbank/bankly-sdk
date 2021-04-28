package bankly

import (
	"encoding/json"
	"errors"
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
func NewAuthentication(session Session) *Authentication {
	return &Authentication{
		session: session,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (a *Authentication) login() (*AuthenticationResponse, error) {
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
	}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(formData.Encode()))

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

		return nil, errors.New(bodyErr.Message)
	}

	return nil, errors.New("error login")
}

//Token ...
func (a Authentication) Token() (string, error) {
	response, err := a.login()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", response.TokenType, response.AccessToken), nil
}
