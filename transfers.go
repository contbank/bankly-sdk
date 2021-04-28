package bankly

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

//Transfers ...
type Transfers struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

//NewTransfers ...
func NewTransfers(session Session) *Transfers {
	return &Transfers{
		session: session,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authentication: NewAuthentication(session),
	}
}

//CreateTransfers ...
func (t *Transfers) CreateTransfers(correlationID string, model TransfersRequest) (*TransfersResponse, error) {
	err := Validator.Struct(model)

	if err != nil {
		return nil, err
	}

	u, err := url.Parse(t.session.APIEndpoint)

	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, TransfersPath)
	endpoint := u.String()

	reqbyte, err := json.Marshal(model)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(reqbyte))

	if err != nil {
		return nil, err
	}

	token, err := t.authentication.Token()

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", correlationID)

	resp, err := t.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusAccepted {
		var body *TransfersResponse

		err = json.Unmarshal(respBody, &body)

		if err != nil {
			return nil, err
		}

		return body, nil
	}

	return nil, errors.New("error create transfers")
}

//FindTransfers ...
func (t *Transfers) FindTransfers() {

}

//FindTransfersByCode ...
func (t *Transfers) FindTransfersByCode() {

}
