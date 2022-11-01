package bankly

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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
	fields := logrus.Fields{
		"request_id": GetRequestID(ctx),
		"identifier": grok.OnlyDigits(identifier),
	}

	url := fmt.Sprintf("cards/document/%s", grok.OnlyDigits(identifier))
	fields["url"] = url

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
	fields := logrus.Fields{
		"request_id": GetRequestID(ctx),
		"proxy":      proxy,
	}

	url := fmt.Sprintf("cards/%s", proxy)
	fields["url"] = url

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
	fields := logrus.Fields{
		"request_id":    GetRequestID(ctx),
		"activate_code": activateCode,
	}

	url := fmt.Sprintf("cards/activateCode/%s", activateCode)
	fields["url"] = url

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
	fields := logrus.Fields{
		"request_id": GetRequestID(ctx),
		"proxy":      proxy,
	}

	url := fmt.Sprintf("cards/%s/nextStatus", proxy)
	fields["url"] = url

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
func (c *Card) GetCardByAccount(ctx context.Context, accountNumber, accountBranch,
	identifier string) ([]CardResponse, error) {
	fields := logrus.Fields{
		"request_id":     GetRequestID(ctx),
		"identifier":     grok.OnlyDigits(identifier),
		"account_branch": accountBranch,
		"account_number": accountNumber,
	}

	url := fmt.Sprintf("cards/account/%s", grok.OnlyLettersOrDigits(accountNumber))
	fields["url"] = url

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
func (c *Card) CreateCard(ctx context.Context, cardDTO *CardCreateDTO) (*CardCreateResponse, error) {
	cardLog := *cardDTO
	cardLog.CardData.Password = ""

	fields := logrus.Fields{
		"request_id": GetRequestID(ctx),
		"object":     cardLog,
	}

	if cardDTO == nil {
		logrus.WithFields(fields).WithError(ErrInvalidParameter).Error("error invalid parameter")
		return nil, ErrInvalidParameter
	}

	url := fmt.Sprintf("cards/%s", strings.ToLower(string(cardDTO.CardType)))
	fields["url"] = url

	body := CardCreateRequest{
		DocumentNumber: grok.OnlyDigits(cardDTO.CardData.DocumentNumber),
		CardName:       grok.ToTitle(cardDTO.CardData.CardName),
		Alias:          grok.ToTitle(cardDTO.CardData.Alias),
		BankAgency:     grok.OnlyLettersOrDigits(cardDTO.CardData.BankAgency),
		BankAccount:    grok.OnlyLettersOrDigits(cardDTO.CardData.BankAccount),
		ProgramId:      cardDTO.CardData.ProgramId,
		Password:       cardDTO.CardData.Password,
		Address:        cardDTO.CardData.Address,
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

// UpdateStatusCardByProxy ...
func (c *Card) UpdateStatusCardByProxy(ctx context.Context, proxy *string,
	cardUpdateStatusDTO *CardUpdateStatusDTO) error {

	cardLog := *cardUpdateStatusDTO
	cardLog.Password = ""

	fields := logrus.Fields{
		"request_id": GetRequestID(ctx),
		"proxy":      proxy,
		"object":     cardLog,
	}

	if proxy == nil || cardUpdateStatusDTO == nil {
		logrus.WithFields(fields).WithError(ErrInvalidParameter).Error("error invalid parameter")
		return ErrInvalidParameter
	}

	url := fmt.Sprintf("cards/%s/status", *proxy)
	fields["url"] = url

	response, err := c.httpClient.Patch(ctx, url, cardUpdateStatusDTO, nil, nil)

	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return err
	} else if response != nil &&
		response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted {
		logrus.WithFields(fields).WithError(ErrCardStatusUpdate).Error("error update status card")
		return err
	}

	logrus.WithFields(fields).Info("card status updated with success")
	return nil
}

// ActivateCardByProxy ...
func (c *Card) ActivateCardByProxy(ctx context.Context, proxy *string,
	cardActivateDTO *CardActivateDTO) error {

	cardLog := *cardActivateDTO
	cardLog.Password = ""

	fields := logrus.Fields{
		"request_id": GetRequestID(ctx),
		"proxy":      proxy,
		"object":     cardLog,
	}

	if proxy == nil || cardActivateDTO == nil {
		logrus.WithFields(fields).WithError(ErrInvalidParameter).Error("error invalid parameter")
		return ErrInvalidParameter
	}

	url := fmt.Sprintf("cards/%s/activate", *proxy)
	fields["url"] = url

	response, err := c.httpClient.Patch(ctx, url, cardActivateDTO, nil, nil)

	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return err
	} else if response != nil &&
		response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted {
		logrus.WithFields(fields).WithError(ErrCardActivate).Error(ErrCardActivate)
		return ErrCardActivate
	}

	logrus.WithFields(fields).Info("card activated with success")
	return nil
}

// ContactlessCardByProxy ...
func (c *Card) ContactlessCardByProxy(ctx context.Context, proxy *string,
	cardContactlessDTO *CardContactlessDTO) error {
	fields := logrus.Fields{
		"request_id": GetRequestID(ctx),
		"proxy":      proxy,
		"object":     cardContactlessDTO,
	}

	if proxy == nil || cardContactlessDTO == nil {
		logrus.WithFields(fields).WithError(ErrInvalidParameter).Error("error invalid parameter")
		return ErrInvalidParameter
	}

	// endpoint url
	url := fmt.Sprintf("cards/%s/contactless", *proxy)
	fields["url"] = url

	logrus.WithFields(fields).Info("processing contactless")

	query := make(map[string]string)
	query["allowContactless"] = strconv.FormatBool(cardContactlessDTO.Active)

	response, err := c.httpClient.Patch(ctx, url, nil, query, nil)

	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return err
	} else if response != nil &&
		response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted {
		logrus.WithFields(fields).WithError(ErrCardActivate).Error(ErrCardActivate)
		return ErrCardActivate
	}

	logrus.WithFields(fields).Info("update card contactless successfully")
	return nil
}

// UpdatePasswordByProxy ...
func (c *Card) UpdatePasswordByProxy(ctx context.Context, proxy string, model CardUpdatePasswordDTO) error {
	fields := logrus.Fields{
		"request_id": GetRequestID(ctx),
		"proxy":      proxy,
	}

	url := fmt.Sprintf("cards/%s/password", proxy)
	fields["url"] = url

	logrus.WithFields(fields).Info("updating card password at bankly")

	resp, err := c.httpClient.Patch(ctx, url, model, nil, nil)
	if err != nil {
		logrus.WithFields(fields).WithField("error_key", "ERROR-CARD-0002").
			WithError(err).Error(err.Error())
		return err
	} else if resp != nil &&
		resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		logrus.WithFields(fields).WithField("error_key", "ERROR-CARD-0003").
			WithError(ErrCardPasswordUpdate).Error(ErrCardPasswordUpdate.Error())
		return err
	}

	logrus.WithFields(fields).Info("card password updated with success")
	return nil
}

// GetTransactionsByProxy ...
func (c *Card) GetTransactionsByProxy(ctx context.Context, proxy *string,
	page, startDate, endDate, pageSize string) (*CardTransactionsResponse, error) {

	fields := logrus.Fields{
		"request_id": GetRequestID(ctx),
		"proxy":      proxy,
		"page":       page,
		"page_size":  pageSize,
		"start_date": startDate,
		"end_date":   endDate,
	}

	if proxy == nil {
		logrus.WithFields(fields).WithError(ErrInvalidParameter).Error("error invalid parameter")
		return nil, ErrInvalidParameter
	}

	url := fmt.Sprintf("cards/%s/transactions", *proxy)
	fields["url"] = url

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
func (c *Card) GetPCIByProxy(ctx context.Context, proxy *string, cardPCIDTO *CardPCIDTO) (*CardPCIResponse, error) {
	fields := logrus.Fields{
		"request_id": GetRequestID(ctx),
		"proxy":      proxy,
	}

	if proxy == nil || cardPCIDTO == nil {
		logrus.WithFields(fields).WithError(ErrInvalidParameter).Error("error invalid parameter")
		return nil, ErrInvalidParameter
	}

	url := fmt.Sprintf("cards/%s/pci", *proxy)
	fields["url"] = url

	resp, err := c.httpClient.Post(ctx, url, cardPCIDTO, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *CardPCIResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			logrus.WithFields(fields).WithError(err).Error("error unmarshal")
			return nil, err
		}

		fields["response"] = response
		logrus.WithFields(fields).Info("response with success")

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

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return nil, FindError(errModel.Code, errModel.Messages...)
	}

	logrus.WithFields(fields).
		Error("error default card response - FindRegistration")

	return nil, ErrDefaultCard
}

// GetTrackingByProxy ...
func (c *Card) GetTrackingByProxy(ctx context.Context, proxy *string) (*CardTrackingResponse, error) {
	fields := logrus.Fields{
		"request_id": GetRequestID(ctx),
		"proxy":      proxy,
	}

	if proxy == nil {
		logrus.WithFields(fields).WithError(ErrInvalidParameter).Error("error invalid parameter")
		return nil, ErrInvalidParameter
	}

	url := fmt.Sprintf("cards/%s/tracking", *proxy)
	fields["url"] = url

	resp, err := c.httpClient.Get(ctx, url, nil, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if respBody != nil && resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
		var response *CardTrackingResponse

		err = json.Unmarshal(respBody, &response)
		if err != nil {
			logrus.WithFields(fields).WithError(err).Error("error unmarshal - card tracking")
			return nil, err
		}

		fields["response"] = response
		logrus.WithFields(fields).Info("response with success")

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		logrus.WithFields(fields).WithError(ErrEntryNotFound).Error("entry not found - card tracking")
		return nil, ErrEntryNotFound
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return nil, FindError(errModel.Code, errModel.Messages...)
	}

	logrus.WithFields(fields).WithError(ErrDefaultCard).Error("error response - card tracking")
	return nil, ErrDefaultCard
}

// CardErrorHandler ...
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
		logrus.WithFields(fields).WithError(err).Error("bankly card error")

		return err
	}

	return ErrDefaultCard
}
