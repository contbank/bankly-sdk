package bankly

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

//Business ...
type Business struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

//NewBusiness ...
func NewBusiness(session Session) *Business {
	return &Business{
		session: session,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authentication: NewAuthentication(session),
	}
}

//CreateBusiness ...
func (c *Business) CreateBusiness(businessRequest BusinessRequest) error {

	endpoint, err := c.getBusinessAPIEndpoint(businessRequest.Document, false)
	if err != nil {
		return err
	}

	reqbyte, err := json.Marshal(businessRequest)

	req, err := http.NewRequest("PUT", *endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		return err
	}

	token, err := c.authentication.Token()
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", c.session.APIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusAccepted {
		return nil
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return err
	}

	if bodyErr.Errors != nil {
		return FindError(bodyErr.Errors[0])
	}
	return errors.New("error create business")
}

//UpdateBusiness ...
func (c *Business) UpdateBusiness(businessDocument string, businessUpdateRequest BusinessUpdateRequest) error {

	endpoint, err := c.getBusinessAPIEndpoint(businessDocument, false)
	if err != nil {
		return err
	}

	reqbyte, err := json.Marshal(businessUpdateRequest)

	req, err := http.NewRequest("PATCH", *endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		return err
	}

	token, err := c.authentication.Token()
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", c.session.APIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusAccepted {
		return nil
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return err
	}

	if bodyErr.Errors != nil {
		return FindError(bodyErr.Errors[0])
	}
	return errors.New("error updating business")
}

//CreateBusinessAccount ...
func (c *Business) CreateBusinessAccount(businessAccountRequest BusinessAccountRequest) (*AccountResponse, error) {

	endpoint, err := c.getBusinessAPIEndpoint(businessAccountRequest.Document, true)
	if err != nil {
		return nil, err
	}

	reqbyte, err := json.Marshal(businessAccountRequest)

	req, err := http.NewRequest("POST", *endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.
			WithError(err).
			Error("error new request")
		return nil, err
	}

	token, err := c.authentication.Token()
	if err != nil {
		logrus.
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", c.session.APIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.
			WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusCreated {
		var bodyResp *AccountResponse

		err = json.Unmarshal(respBody, &bodyResp)
		if err != nil {
			logrus.
				WithError(err).
				Error("error unmarshal")
			return nil, err
		}

		return bodyResp, nil
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.
			WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	if bodyErr.Errors != nil {
		return nil, FindError(bodyErr.Errors[0])
	}
	return nil, errors.New("error create business account")
}

//FindBusiness ...
func (c *Business) FindBusiness(document string) (*BusinessResponse, error) {

	endpoint, err := c.getBusinessAPIEndpoint(document, false)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", *endpoint, nil)
	if err != nil {
		return nil, err
	}

	token, err := c.authentication.Token()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", c.session.APIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response BusinessResponse

		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if bodyErr.Errors != nil {
		return nil, FindError(bodyErr.Errors[0])
	}

	return nil, errors.New("error find business")
}

//FindBusinessAccounts ...
func (c *Business) FindBusinessAccounts(document string) ([]AccountResponse, error) {

	endpoint, err := c.getBusinessAPIEndpoint(document, true)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", *endpoint, nil)
	if err != nil {
		return nil, err
	}

	token, err := c.authentication.Token()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", c.session.APIVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response []AccountResponse

		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return nil, err
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if bodyErr.Errors != nil {
		return nil, FindError(bodyErr.Errors[0])
	}

	return nil, errors.New("error find business accounts")
}

// getBusinessAPIEndpoint
func (c *Business) getBusinessAPIEndpoint(document string, isAccountPath bool) (*string, error) {
	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.
			WithError(err).
			Error("error api endpoint")
		return nil, err
	}
	u.Path = path.Join(u.Path, BusinessPath)
	u.Path = path.Join(u.Path, grok.OnlyDigits(document))
	if isAccountPath == true {
		u.Path = path.Join(u.Path, AccountsPath)
	}
	endpoint := u.String()
	return &endpoint, nil
}