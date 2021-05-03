package bankly

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	// InternalBankCode ...
	InternalBankCode string = "332"
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

// CreateExternalTransfer ...
func (t *Transfers) CreateExternalTransfer(correlationID string, model TransfersRequest) (*TransfersResponse, error) {
	logrus.
		WithFields(logrus.Fields{
			"correlation_id" : correlationID,
		}).
		Info("create external transfer")
	return t.createTransfers(correlationID, model)
}

// CreateInternalTransfer ...
func (t *Transfers) CreateInternalTransfer(correlationID string, model TransfersRequest) (*TransfersResponse, error) {
	logrus.
		WithFields(logrus.Fields{
			"correlation_id" : correlationID,
		}).
		Info("create internal transfer")
	model.Recipient.BankCode = InternalBankCode
	return t.createTransfers(correlationID, model)
}

// createTransfers ...
func (t *Transfers) createTransfers(correlationID string, model TransfersRequest) (*TransfersResponse, error) {
	err := Validator.Struct(model)
	if err != nil {
		logrus.
			WithError(err).
			Error("error validating model")
		return nil, err
	}

	endpoint, err := t.getTransferAPIEndpoint()
	if err != nil {
		logrus.
			WithError(err).
			Error("error getting api endpoint")
		return nil, err
	}

	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.
			WithError(err).
			Error("error marshal model")
		return nil, err
	}

	req, err := http.NewRequest("POST", *endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.
			WithError(err).
			Error("error new request")
		return nil, err
	}

	token, err := t.authentication.Token()
	if err != nil {
		logrus.
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", correlationID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.
			WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusAccepted {
		var body *TransfersResponse

		err = json.Unmarshal(respBody, &body)
		if err != nil {
			logrus.
				WithError(err).
				Error("error unmarshal")
			return nil, err
		}
		
		return body, nil
	}

	return nil, errors.New("error create transfers")
}

// FindTransfers ...
func (t *Transfers) FindTransfers() {

}

// FindTransfersByCode ...
func (t *Transfers) FindTransfersByCode() {

}

// getTransferAPIEndpoint
func (t *Transfers) getTransferAPIEndpoint() (*string, error) {
	u, err := url.Parse(t.session.APIEndpoint)
	if err != nil {
		logrus.
			WithError(err).
			Error("error api endpoint")
		return nil, err
	}
	u.Path = path.Join(u.Path, TransfersPath)
	endpoint := u.String()
	return &endpoint, nil
}
