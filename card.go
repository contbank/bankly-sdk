package bankly

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Card struct {
	httpClient NewHttpClient
}

//NewCard ...
func NewCard(newHttpClient NewHttpClient) *Card {
	return &Card{newHttpClient}
}

//Cards ...
func (c *Card) GetCardsByIdentifier(ctx context.Context, identifier string) ([]CardResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"identifier": identifier,
	}

	url := "cards/document/" + identifier

	resp, err := c.httpClient.Get(ctx, url, nil)
	if err != nil {
		logErrorWithFields(fields, err, err.Error(), nil)
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding body response", nil)
		return nil, err
	}

	var cardsResponseDTO []CardResponseDTO
	err = json.Unmarshal(respBody, &cardsResponseDTO)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding json response", nil)
		return nil, ErrDefaultCard
	}

	var cards []CardResponse
	for _, crd := range cardsResponseDTO {
		cards = append(cards, *parseResponseCard(&crd))
	}

	defer resp.Body.Close()
	return cards, nil
}

func (c *Card) GetCardByProxy(ctx context.Context, proxy string) (*CardResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"identifier": proxy,
	}

	url := "cards/" + proxy

	resp, err := c.httpClient.Get(ctx, url, nil)
	if err != nil {
		logErrorWithFields(fields, err, err.Error(), nil)
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding body response", nil)
		return nil, err
	}

	var cardResponseDTO *CardResponseDTO
	err = json.Unmarshal(respBody, &cardResponseDTO)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding json response", nil)
		return nil, ErrDefaultCard
	}

	defer resp.Body.Close()
	return parseResponseCard(cardResponseDTO), nil
}

func (c *Card) GetNextStatusByProxy(ctx context.Context, proxy string) ([]CardNextStatus, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"identifier": proxy,
	}

	url := "cards/" + proxy + "/nextStatus"

	resp, err := c.httpClient.Get(ctx, url, nil)
	if err != nil {
		logErrorWithFields(fields, err, err.Error(), nil)
		return nil, err
	}

	if resp.StatusCode == 204 {
		return []CardNextStatus{}, nil
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding body response", nil)
		return nil, err
	}

	var cardNextStatus []CardNextStatus
	err = json.Unmarshal(respBody, &cardNextStatus)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding json response", nil)
		return nil, ErrDefaultCard
	}

	defer resp.Body.Close()
	return cardNextStatus, nil
}

func (c *Card) GetCardByAccount(ctx context.Context, bankAccount, bankAgency, documentNumber string) ([]CardResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"identifier": bankAccount,
	}

	url := "cards/account/" + bankAccount
	query := make(map[string]string)

	query["agency"] = bankAgency
	query["documentNumber"] = documentNumber

	resp, err := c.httpClient.Get(ctx, url, query)
	if err != nil {
		logErrorWithFields(fields, err, err.Error(), nil)
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding body response", nil)
		return nil, err
	}

	var cardsResponseDTO []CardResponseDTO
	err = json.Unmarshal(respBody, &cardsResponseDTO)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding json response", nil)
		return nil, ErrDefaultCard
	}

	var cards []CardResponse
	for _, crd := range cardsResponseDTO {
		cards = append(cards, *parseResponseCard(&crd))
	}

	defer resp.Body.Close()
	return cards, nil
}

func (c *Card) CreateCard(ctx context.Context, cardDTO CardCreateDTO) (*CardCreateResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	url := "cards/"
	switch cardDTO.CardType {
	case VirtualCardType:
		url = url + VirtualCardType
	case PhysicalCardType:
		url = url + PhysicalCardType
	}

	body := CardCreateRequest{
		DocumentNumber: cardDTO.CardData.DocumentNumber,
		CardName:       cardDTO.CardData.CardName,
		Alias:          cardDTO.CardData.Alias,
		BankAgency:     cardDTO.CardData.BankAgency,
		BankAccount:    cardDTO.CardData.BankAccount,
		ProgramId:      cardDTO.CardData.ProgramId,
		Password:       cardDTO.CardData.Password,
	}

	resp, err := c.httpClient.Post(ctx, url, body)
	if err != nil {
		logErrorWithFields(fields, err, err.Error(), nil)
		return nil, err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	var response *CardCreateResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding json response", nil)
		return nil, ErrDefaultCard
	}

	defer resp.Body.Close()
	return response, nil
}

func (c *Card) UpdateStatusCard(ctx context.Context, proxy string, cardUpdateStatusDTO CardUpdateStatusDTO) (*http.Response, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	url := "cards/" + proxy + "/status"

	resp, err := c.httpClient.Patch(ctx, url, cardUpdateStatusDTO)
	if err != nil {
		logErrorWithFields(fields, err, err.Error(), nil)
		return nil, err
	}

	return resp, nil
}

func (c *Card) GetTransactionsByProxy(ctx context.Context, proxy, page, startDate, endDate, pageSize string) (*CardTransactionsResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	url := "cards/" + proxy + "/transactions"
	
	query := make(map[string]string)
	query["page"] = page
	query["startDate"] = startDate
	query["endDate"] = endDate
	query["pageSize"] = pageSize

	resp, err := c.httpClient.Get(ctx, url, query)
	if err != nil {
		logErrorWithFields(fields, err, err.Error(), nil)
		return nil, err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	var response *CardTransactionsResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logErrorWithFields(fields, err, "error decoding json response", nil)
		return nil, ErrDefaultCard
	}

	defer resp.Body.Close()
	return response, nil
}

func logErrorWithFields(fields logrus.Fields, err error, msg string, hasField map[string]interface{}) {
	if hasField != nil {
		for prop, value := range hasField {
			logrus.
				WithField(prop, value).
				WithFields(fields).
				WithError(err).
				Error(msg)
		}
	} else {
		logrus.
			WithFields(fields).
			WithError(err).
			Error(msg)
	}
}

func logInfoWithFields(fields logrus.Fields, msg string) {
	logrus.
		WithFields(fields).
		Info(msg)
}

func parseResponseCard(cardResponseDTO *CardResponseDTO) *CardResponse {
	return &CardResponse{
		Created:          cardResponseDTO.Created,
		CompanyKey:       cardResponseDTO.CompanyKey,
		DocumentNumber:   cardResponseDTO.DocumentNumber,
		ActivateCode:     cardResponseDTO.ActivateCode,
		BankAgency:       cardResponseDTO.BankAgency,
		BankAccount:      cardResponseDTO.BankAccount,
		LastFourDigits:   cardResponseDTO.LastFourDigits,
		Proxy:            cardResponseDTO.Proxy,
		Name:             cardResponseDTO.Name,
		Alias:            cardResponseDTO.Alias,
		CardType:         cardResponseDTO.CardType,
		Status:           cardResponseDTO.Status,
		PhysicalBinds:    cardResponseDTO.PhysicalBinds,
		VirtualBind:      cardResponseDTO.VirtualBind,
		AllowContactless: cardResponseDTO.AllowContactless,
		Address:          cardResponseDTO.Address,
		HistoryStatus:    cardResponseDTO.HistoryStatus,
		ActivatedAt:      cardResponseDTO.ActivatedAt,
		LastUpdatedAt:    cardResponseDTO.LastUpdatedAt,
		IsFirtual:        cardResponseDTO.IsFirtual,
		IsPos:            cardResponseDTO.IsPos,
		SettlementDay:    cardResponseDTO.PaymentDay,
	}
}
