package bankly

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

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
func NewTransfers(httpClient *http.Client, session Session) *Transfers {
	return &Transfers{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(session),
	}
}

// CreateTransfer ...
func (t *Transfers) CreateTransfer(ctx context.Context, correlationID string, model TransfersRequest) (*TransferByCodeResponse, error) {
	logrus.
		WithFields(logrus.Fields{
			"correlation_id": correlationID,
		}).
		Info("create transfer")
	return t.createTransferOperation(ctx, correlationID, model)
}

// CreateInternalTransfer ...
func (t *Transfers) CreateInternalTransfer(ctx context.Context, correlationID string, model TransfersRequest) (*TransferByCodeResponse, error) {
	logrus.
		WithFields(logrus.Fields{
			"correlation_id": correlationID,
		}).
		Info("create internal transfer")
	// TODO quando transação interna, necessário validar algo? limite de transação é maior do que quando externa?
	model.Recipient.BankCode = InternalBankCode
	return t.createTransferOperation(ctx, correlationID, model)
}

// CreateExternalTransfer ...
func (t *Transfers) CreateExternalTransfer(ctx context.Context, requestID string, model TransfersRequest) (*TransferByCodeResponse, error) {
	logrus.
		WithFields(logrus.Fields{
			"request_id": requestID,
		}).
		Info("create external transfer")
	return t.createTransferOperation(ctx, requestID, model)
}

// createTransferOperation ...
func (t *Transfers) createTransferOperation(ctx context.Context, requestID string, model TransfersRequest) (*TransferByCodeResponse, error) {

	fields := logrus.Fields{
		"request_id": requestID,
		"model":      model,
	}

	err := grok.Validator.Struct(model)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := t.getTransferAPIEndpoint(requestID, nil, nil, nil, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error getting api endpoint")
		return nil, err
	}

	reqbyte, err := json.Marshal(model)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error marshal model")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(reqbyte))
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error new request")
		return nil, err
	}

	token, err := t.authentication.Token(ctx)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", requestID)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.
			WithFields(fields).
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
				WithFields(fields).
				WithError(err).
				Error("error unmarshal")
			return nil, err
		}

		return body, nil
	}

	var bodyErr *TransferErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error - createTransferOperation")
		return nil, err
	}

	if bodyErr != nil && (len(bodyErr.Errors) > 0 || bodyErr.Code != "") {
		logrus.
			WithFields(fields).
			Error("body error - createTransferOperation")
		return nil, FindTransferError(*bodyErr)
	}

	logrus.
		WithFields(fields).
		Error("default error transfer - createTransferOperation")

	return nil, ErrDefaultTransfers
}

// FindTransfers ...
func (t *Transfers) FindTransfers(ctx context.Context, requestID *string,
	branch *string, account *string, pageSize *int, nextPage *string) (*TransfersResponse, error) {

	if requestID == nil {
		return nil, ErrInvalidCorrelationID
	} else if branch == nil || account == nil {
		return nil, ErrInvalidAccountNumber
	}

	fields := logrus.Fields{
		"request_id": requestID,
		"branch":     branch,
		"account":    account,
		"page_size":  pageSize,
		"next_page":  nextPage,
	}

	endpoint, err := t.getTransferAPIEndpoint(*requestID, nil, branch, account, pageSize)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error transfer api endpoint")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		return nil, err
	}

	token, err := t.authentication.Token(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", *requestID)

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
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error unmarshal")
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
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error unmarshal transfer error response")
		return nil, err
	}

	if bodyErr != nil && (len(bodyErr.Errors) > 0 || bodyErr.Code != "") {
		logrus.
			WithFields(fields).
			Error("body error")
		return nil, FindTransferError(*bodyErr)
	}

	logrus.
		WithFields(fields).
		WithError(err).
		Error("default error transfer - FindTransfers")

	return nil, ErrDefaultFindTransfers
}

// FindTransfersByCode ...
func (t *Transfers) FindTransfersByCode(ctx context.Context, requestID *string,
	authenticationCode *string, branch *string, account *string) (*TransferByCodeResponse, error) {

	if requestID == nil {
		return nil, ErrInvalidCorrelationID
	} else if authenticationCode == nil || branch == nil || account == nil {
		return nil, ErrInvalidAuthenticationCodeOrAccount
	}

	fields := logrus.Fields{
		"request_id":          requestID,
		"authentication_code": authenticationCode,
		"branch":              branch,
		"account":             account,
	}

	endpoint, err := t.getTransferAPIEndpoint(*requestID, authenticationCode, branch, account, nil)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		return nil, err
	}

	token, err := t.authentication.Token(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", t.session.APIVersion)
	req.Header.Add("x-correlation-id", *requestID)

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
			logrus.
				WithFields(fields).
				WithError(err).
				Error("error unmarshal")
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
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	if bodyErr != nil && (len(bodyErr.Errors) > 0 || bodyErr.Code != "") {
		logrus.
			WithFields(fields).
			Error("body error")
		return nil, FindTransferError(*bodyErr)
	}

	logrus.
		WithFields(fields).
		Error("default error transfer - FindTransfersByCode")

	return nil, ErrDefaultFindTransfers
}

// getTransferAPIEndpoint
func (t *Transfers) getTransferAPIEndpoint(correlationID string,
	authenticationCode *string, branch *string, account *string, pageSize *int) (*string, error) {

	u, err := url.Parse(t.session.APIEndpoint)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"correlation_id":      correlationID,
				"authentication_code": authenticationCode,
				"branch":              branch,
				"account":             account,
				"page_size":           pageSize,
			}).
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
