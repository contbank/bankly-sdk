package bankly

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
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

// CreateTransfer ...
func (t *Transfers) CreateTransfer(correlationID string, model TransfersRequest) (*TransferByCodeResponse, error) {
	logrus.
		WithFields(logrus.Fields{
			"correlation_id": correlationID,
		}).
		Info("create transfer")
	return t.createTransferOperation(correlationID, model)
}

// CreateInternalTransfer ...
func (t *Transfers) CreateInternalTransfer(correlationID string, model TransfersRequest) (*TransferByCodeResponse, error) {
	logrus.
		WithFields(logrus.Fields{
			"correlation_id": correlationID,
		}).
		Info("create internal transfer")
	// TODO quando transação interna, necessário validar algo? limite de transação é maior do que quando externa?
	model.Recipient.BankCode = InternalBankCode
	return t.createTransferOperation(correlationID, model)
}

// CreateExternalTransfer ...
func (t *Transfers) CreateExternalTransfer(correlationID string, model TransfersRequest) (*TransferByCodeResponse, error) {
	logrus.
		WithFields(logrus.Fields{
			"correlation_id": correlationID,
		}).
		Info("create external transfer")
	return t.createTransferOperation(correlationID, model)
}

// createTransferOperation ...
func (t *Transfers) createTransferOperation(correlationID string, model TransfersRequest) (*TransferByCodeResponse, error) {

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.
			WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := t.getTransferAPIEndpoint(nil, nil, nil, nil)
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
		var body *TransferByCodeResponse

		err = json.Unmarshal(respBody, &body)
		if err != nil {
			logrus.
				WithError(err).
				Error("error unmarshal")
			return nil, err
		}

		return body, nil
	}

	var bodyErr *TransferErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		return nil, FindTransferError(*bodyErr)
	}

	return nil, errors.New("error create transfer operation")
}

// FindTransfers ...
func (t *Transfers) FindTransfers(correlationID *string,
	branch *string, account *string, pageSize *int, nextPage *string) (*TransfersResponse, error) {

	if correlationID == nil {
		return nil, ErrInvalidCorrelationId
	} else if branch == nil || account == nil {
		return nil, ErrInvalidAccountNumber
	}

	endpoint, err := t.getTransferAPIEndpoint(nil, branch, account, pageSize)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", *endpoint, nil)
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
	req.Header.Add("x-correlation-id", *correlationID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response TransfersResponse

		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr *TransferErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		return nil, FindTransferError(*bodyErr)
	}

	return nil, errors.New("error find transfers")
}

// FindTransfersByCode ...
func (t *Transfers) FindTransfersByCode(correlationID *string,
	authenticationCode *string, branch *string, account *string) (*TransferByCodeResponse, error) {

	if correlationID == nil {
		return nil, ErrInvalidCorrelationId
	} else if authenticationCode == nil || branch == nil || account == nil {
		return nil, ErrInvalidAuthenticationCodeOrAccount
	}

	endpoint, err := t.getTransferAPIEndpoint(authenticationCode, branch, account, nil)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", *endpoint, nil)
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
	req.Header.Add("x-correlation-id", *correlationID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response TransferByCodeResponse

		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr *TransferErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		return nil, FindTransferError(*bodyErr)
	}

	return nil, errors.New("error find transfer by code")
}

// getTransferAPIEndpoint
func (t *Transfers) getTransferAPIEndpoint(
	authenticationCode *string, branch *string, account *string, pageSize *int) (*string, error) {

	u, err := url.Parse(t.session.APIEndpoint)
	if err != nil {
		logrus.
			WithError(err).
			Error("error api endpoint")
		return nil, err
	}
	u.Path = path.Join(u.Path, TransfersPath)

	if authenticationCode != nil {
		u.Path = path.Join(u.Path, *authenticationCode)
	}

	if branch != nil && account != nil {
		q := u.Query()
		q.Set("branch", *branch)
		q.Set("account", *account)
		if pageSize != nil {
			q.Set("pageSize", strconv.Itoa(*pageSize))
		}
		u.RawQuery = q.Encode()
	}

	endpoint := u.String()
	return &endpoint, nil
}
