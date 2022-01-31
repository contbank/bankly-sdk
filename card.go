package bankly

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/contbank/grok"

	"github.com/sirupsen/logrus"
)

type Card struct {
	httpClient BanklyHttpClient
}

//NewCard ...
func NewCard(newHttpClient BanklyHttpClient) *Card {
	newHttpClient.errorHandler = CardErrorHandler
	return &Card{newHttpClient}
}

// GetCardsByIdentifier ...
func (c *Card) GetCardsByIdentifier(ctx context.Context, identifier string) ([]CardResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"identifier": identifier,
	}

	url := "cards/document/" + grok.OnlyDigits(identifier)

	resp, err := c.httpClient.Get(ctx, url, nil, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	var cardsResponseDTO []CardResponseDTO

	err = json.Unmarshal(respBody, &cardsResponseDTO)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultCard
	}

	var cards []CardResponse
	for _, crd := range cardsResponseDTO {
		cards = append(cards, *ParseResponseCard(&crd))
	}

	defer resp.Body.Close()
	return cards, nil
}

// GetCardByProxy ...
func (c *Card) GetCardByProxy(ctx context.Context, proxy string) (*CardResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"proxy":      proxy,
	}

	url := "cards/" + proxy

	resp, err := c.httpClient.Get(ctx, url, nil, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	var cardResponseDTO *CardResponseDTO
	err = json.Unmarshal(respBody, &cardResponseDTO)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultCard
	}

	defer resp.Body.Close()
	return ParseResponseCard(cardResponseDTO), nil
}

// GetCardByActivateCode ...
func (c *Card) GetCardByActivateCode(ctx context.Context, activateCode string) ([]CardResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id":    requestID,
		"activate_code": activateCode,
	}

	url := "cards/activateCode/" + activateCode

	resp, err := c.httpClient.Get(ctx, url, nil, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	var cardsResponseDTO []CardResponseDTO

	err = json.Unmarshal(respBody, &cardsResponseDTO)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultCard
	}

	var cards []CardResponse
	for _, crd := range cardsResponseDTO {
		cards = append(cards, *ParseResponseCard(&crd))
	}

	defer resp.Body.Close()
	return cards, nil
}

// GetNextStatusByProxy ...
func (c *Card) GetNextStatusByProxy(ctx context.Context, proxy string) ([]CardNextStatus, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"proxy":      proxy,
	}

	url := "cards/" + proxy + "/nextStatus"

	resp, err := c.httpClient.Get(ctx, url, nil, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	if resp.StatusCode == 204 {
		return []CardNextStatus{}, nil
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	var cardNextStatus []CardNextStatus
	err = json.Unmarshal(respBody, &cardNextStatus)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultCard
	}

	defer resp.Body.Close()
	return cardNextStatus, nil
}

// GetCardByAccount ...
func (c *Card) GetCardByAccount(ctx context.Context, accountNumber, accountBranch, identifier string) ([]CardResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id":     requestID,
		"identifier":     identifier,
		"account_branch": accountBranch,
		"account_number": accountNumber,
	}

	url := "cards/account/" + grok.OnlyLettersOrDigits(accountNumber)

	query := make(map[string]string)
	query["agency"] = grok.OnlyLettersOrDigits(accountBranch)
	query["documentNumber"] = grok.OnlyDigits(identifier)

	resp, err := c.httpClient.Get(ctx, url, query, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	var cardsResponseDTO []CardResponseDTO
	err = json.Unmarshal(respBody, &cardsResponseDTO)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultCard
	}

	var cards []CardResponse

	for _, crd := range cardsResponseDTO {
		cards = append(cards, *ParseResponseCard(&crd))
	}

	defer resp.Body.Close()
	return cards, nil
}

// CreateCard ...
func (c *Card) CreateCard(ctx context.Context, cardDTO CardCreateDTO) (*CardCreateResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	cardLog := cardDTO
	cardLog.CardData.Password = ""
	fields["object"] = cardLog

	url := "cards/" + strings.ToLower(string(cardDTO.CardType))

	body := CardCreateRequest{
		DocumentNumber: grok.OnlyDigits(cardDTO.CardData.DocumentNumber),
		CardName:       grok.ToTitle(cardDTO.CardData.CardName),
		Alias:          grok.ToTitle(cardDTO.CardData.Alias),
		BankAgency:     grok.OnlyLettersOrDigits(cardDTO.CardData.BankAgency),
		BankAccount:    grok.OnlyLettersOrDigits(cardDTO.CardData.BankAccount),
		ProgramId:      cardDTO.CardData.ProgramId,
		Password:       cardDTO.CardData.Password,
	}

	resp, err := c.httpClient.Post(ctx, url, body, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	var response *CardCreateResponse

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultCard
	}

	defer resp.Body.Close()
	return response, nil
}

// UpdateStatusCard ...
func (c *Card) UpdateStatusCard(ctx context.Context, proxy string, cardUpdateStatusDTO CardUpdateStatusDTO) (*http.Response, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	cardLog := cardUpdateStatusDTO
	cardLog.Password = ""
	fields["object"] = cardLog

	url := "cards/" + proxy + "/status"

	resp, err := c.httpClient.Patch(ctx, url, cardUpdateStatusDTO, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	return resp, nil
}

// ActivateCardByProxy
func (c *Card) ActivateCardByProxy(ctx context.Context, proxy string, cardActivateDTO CardActivateDTO) (*http.Response, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	cardLog := cardActivateDTO
	cardLog.Password = ""
	fields["object"] = cardLog

	url := "cards/" + proxy + "/activate"

	resp, err := c.httpClient.Patch(ctx, url, cardActivateDTO, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	return resp, nil
}

// AlterPasswordByProxy
func (c *Card) AlterPasswordByProxy(ctx context.Context, proxy string, cardAlterPasswordDTO CardAlterPasswordDTO) (*http.Response, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	cardLog := cardAlterPasswordDTO
	cardLog.Password = ""
	fields["object"] = cardLog

	url := "cards/" + proxy + "/password"

	resp, err := c.httpClient.Patch(ctx, url, cardAlterPasswordDTO, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	return resp, nil
}

// GetTransactionsByProxy ...
func (c *Card) GetTransactionsByProxy(ctx context.Context, proxy, page, startDate, endDate, pageSize string) (*CardTransactionsResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"proxy":      proxy,
	}

	url := "cards/" + proxy + "/transactions"

	query := make(map[string]string)
	query["page"] = page
	query["startDate"] = startDate
	query["endDate"] = endDate
	query["pageSize"] = pageSize

	resp, err := c.httpClient.Get(ctx, url, query, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	var response *CardTransactionsResponse

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultCard
	}

	defer resp.Body.Close()
	return response, nil
}

// GetPCIByProxy
func (c *Card) GetPCIByProxy(ctx context.Context, proxy string, cardPCIDTO CardPCIDTO) (*CardPCIResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"proxy":      proxy,
	}

	url := "cards/" + proxy + "/pci"

	resp, err := c.httpClient.Post(ctx, url, cardPCIDTO, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	var response *CardPCIResponse

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultCard
	}

	defer resp.Body.Close()
	return response, nil
}

//CardErrorHandler ...
func CardErrorHandler(fields logrus.Fields, resp *http.Response) error {
	var bodyErr *ErrorResponse
	respBody, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return ErrDefaultCard
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err := FindCardError(errModel.Code, errModel.Messages...)

		fields["bankly_error"] = bodyErr
		logrus.WithFields(fields).WithError(err).Error("bankly get card error")

		return err
	}
	return ErrDefaultCard
}
