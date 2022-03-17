package bankly

import (
	"context"
	"encoding/json"
	"fmt"
	errors "github.com/contbank/bankly-sdk/pkg/errors"
	models "github.com/contbank/bankly-sdk/pkg/models"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/contbank/grok"

	"github.com/sirupsen/logrus"
)

type Card struct {
	httpClient utils.BanklyHttpClient
}

//NewCard ...
func NewCard(newHttpClient utils.BanklyHttpClient) *Card {
	newHttpClient.ErrorHandler = CardErrorHandler
	return &Card{newHttpClient}
}

// GetCardsByIdentifier ...
func (c *Card) GetCardsByIdentifier(ctx context.Context, identifier string) ([]models.CardResponse, error) {
	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
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

	var cardsResponseDTO []models.CardResponseDTO

	err = json.Unmarshal(respBody, &cardsResponseDTO)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, errors.ErrDefaultCard
	}

	var cards []models.CardResponse
	for _, crd := range cardsResponseDTO {
		cards = append(cards, *models.ParseResponseCard(&crd))
	}

	defer resp.Body.Close()
	return cards, nil
}

// GetCardByProxy ...
func (c *Card) GetCardByProxy(ctx context.Context, proxy string) (*models.CardResponse, error) {
	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
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

	var cardResponseDTO *models.CardResponseDTO
	err = json.Unmarshal(respBody, &cardResponseDTO)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, errors.ErrDefaultCard
	}

	defer resp.Body.Close()
	return models.ParseResponseCard(cardResponseDTO), nil
}

// GetCardByActivateCode ...
func (c *Card) GetCardByActivateCode(ctx context.Context, activateCode string) ([]models.CardResponse, error) {
	fields := logrus.Fields{
		"request_id":    utils.GetRequestID(ctx),
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

	var cardsResponseDTO []models.CardResponseDTO

	err = json.Unmarshal(respBody, &cardsResponseDTO)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, errors.ErrDefaultCard
	}

	var cards []models.CardResponse
	for _, crd := range cardsResponseDTO {
		cards = append(cards, *models.ParseResponseCard(&crd))
	}

	defer resp.Body.Close()
	return cards, nil
}

// GetNextStatusByProxy ...
func (c *Card) GetNextStatusByProxy(ctx context.Context, proxy string) ([]models.CardNextStatus, error) {
	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
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
		return []models.CardNextStatus{}, nil
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding body response")
		return nil, err
	}

	var cardNextStatus []models.CardNextStatus
	err = json.Unmarshal(respBody, &cardNextStatus)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, errors.ErrDefaultCard
	}

	defer resp.Body.Close()
	return cardNextStatus, nil
}

// GetCardByAccount ...
func (c *Card) GetCardByAccount(ctx context.Context, accountNumber, accountBranch,
	identifier string) ([]models.CardResponse, error) {
	fields := logrus.Fields{
		"request_id":     utils.GetRequestID(ctx),
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

	var cardsResponseDTO []models.CardResponseDTO
	err = json.Unmarshal(respBody, &cardsResponseDTO)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, errors.ErrDefaultCard
	}

	var cards []models.CardResponse

	for _, crd := range cardsResponseDTO {
		cards = append(cards, *models.ParseResponseCard(&crd))
	}

	defer resp.Body.Close()
	return cards, nil
}

// CreateCard ...
func (c *Card) CreateCard(ctx context.Context, cardDTO *models.CardCreateDTO) (*models.CardCreateResponse, error) {
	cardLog := *cardDTO
	cardLog.CardData.Password = ""

	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
		"object":     cardLog,
	}

	if cardDTO == nil {
		logrus.WithFields(fields).WithError(errors.ErrInvalidParameter).Error("error invalid parameter")
		return nil, errors.ErrInvalidParameter
	}

	url := fmt.Sprintf("cards/%s", strings.ToLower(string(cardDTO.CardType)))
	fields["url"] = url

	body := models.CardCreateRequest{
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

	var response *models.CardCreateResponse

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, errors.ErrDefaultCard
	}

	defer resp.Body.Close()
	return response, nil
}

// UpdateStatusCardByProxy ...
func (c *Card) UpdateStatusCardByProxy(ctx context.Context, proxy *string,
	cardUpdateStatusDTO *models.CardUpdateStatusDTO) error {

	cardLog := *cardUpdateStatusDTO
	cardLog.Password = ""

	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
		"proxy":      proxy,
		"object":     cardLog,
	}

	if proxy == nil || cardUpdateStatusDTO == nil {
		logrus.WithFields(fields).WithError(errors.ErrInvalidParameter).Error("error invalid parameter")
		return errors.ErrInvalidParameter
	}

	url := fmt.Sprintf("cards/%s/status", *proxy)
	fields["url"] = url

	response, err := c.httpClient.Patch(ctx, url, cardUpdateStatusDTO, nil)

	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return err
	} else if response != nil &&
		response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted {
		logrus.WithFields(fields).WithError(errors.ErrCardStatusUpdate).Error("error update status card")
		return err
	}

	logrus.WithFields(fields).Info("card status updated with success")
	return nil
}

// ActivateCardByProxy
func (c *Card) ActivateCardByProxy(ctx context.Context, proxy *string,
	cardActivateDTO *models.CardActivateDTO) error {

	cardLog := *cardActivateDTO
	cardLog.Password = ""

	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
		"proxy":      proxy,
		"object":     cardLog,
	}

	if proxy == nil || cardActivateDTO == nil {
		logrus.WithFields(fields).WithError(errors.ErrInvalidParameter).Error("error invalid parameter")
		return errors.ErrInvalidParameter
	}

	url := fmt.Sprintf("cards/%s/activate", *proxy)
	fields["url"] = url

	response, err := c.httpClient.Patch(ctx, url, cardActivateDTO, nil)

	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return err
	} else if response != nil &&
		response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted {
		logrus.WithFields(fields).WithError(errors.ErrCardActivate).Error(errors.ErrCardActivate)
		return errors.ErrCardActivate
	}

	logrus.WithFields(fields).Info("card activated with success")
	return nil
}

// UpdatePasswordByProxy
func (c *Card) UpdatePasswordByProxy(ctx context.Context, proxy *string,
	cardUpdatePasswordDTO *models.CardUpdatePasswordDTO) error {

	cardLog := *cardUpdatePasswordDTO
	cardLog.Password = ""

	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
		"proxy":      proxy,
		"object":     cardLog,
	}

	if proxy == nil || cardUpdatePasswordDTO == nil {
		logrus.WithFields(fields).WithError(errors.ErrInvalidParameter).Error("error invalid parameter")
		return errors.ErrInvalidParameter
	}

	url := fmt.Sprintf("cards/%s/password", *proxy)
	fields["url"] = url

	resp, err := c.httpClient.Patch(ctx, url, cardUpdatePasswordDTO, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error(err.Error())
		return err
	} else if resp != nil &&
		resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		logrus.WithFields(fields).WithError(errors.ErrCardPasswordUpdate).Error(errors.ErrCardPasswordUpdate.Error())
		return err
	}

	logrus.WithFields(fields).Info("card password updated with success")
	return nil
}

// GetTransactionsByProxy ...
func (c *Card) GetTransactionsByProxy(ctx context.Context, proxy *string,
	page, startDate, endDate, pageSize string) (*models.CardTransactionsResponse, error) {

	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
		"proxy":      proxy,
		"page":       page,
		"page_size":  pageSize,
		"start_date": startDate,
		"end_date":   endDate,
	}

	if proxy == nil {
		logrus.WithFields(fields).WithError(errors.ErrInvalidParameter).Error("error invalid parameter")
		return nil, errors.ErrInvalidParameter
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

	var response *models.CardTransactionsResponse

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, errors.ErrDefaultCard
	}

	defer resp.Body.Close()
	return response, nil
}

// GetPCIByProxy
func (c *Card) GetPCIByProxy(ctx context.Context, proxy *string, cardPCIDTO *models.CardPCIDTO) (*models.CardPCIResponse, error) {
	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
		"proxy":      proxy,
	}

	if proxy == nil || cardPCIDTO == nil {
		logrus.WithFields(fields).WithError(errors.ErrInvalidParameter).Error("error invalid parameter")
		return nil, errors.ErrInvalidParameter
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
		var response *models.CardPCIResponse

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
		return nil, errors.ErrEntryNotFound
	}

	var bodyErr *errors.ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return nil, errors.FindError(errModel.Code, errModel.Messages...)
	}

	logrus.WithFields(fields).
		Error("error default card response - FindRegistration")

	return nil, errors.ErrDefaultCard
}

// GetTrackingByProxy
func (c *Card) GetTrackingByProxy(ctx context.Context, proxy *string) (*models.CardTrackingResponse, error) {
	fields := logrus.Fields{
		"request_id": utils.GetRequestID(ctx),
		"proxy":      proxy,
	}

	if proxy == nil {
		logrus.WithFields(fields).WithError(errors.ErrInvalidParameter).Error("error invalid parameter")
		return nil, errors.ErrInvalidParameter
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

	if resp.StatusCode == http.StatusOK {
		var response *models.CardTrackingResponse

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
		return nil, errors.ErrEntryNotFound
	}

	var bodyErr *errors.ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return nil, errors.FindError(errModel.Code, errModel.Messages...)
	}

	logrus.WithFields(fields).
		Error("error default card response - FindRegistration")

	return nil, errors.ErrDefaultCard
}

// CardErrorHandler ...
func CardErrorHandler(fields logrus.Fields, resp *http.Response) error {
	var bodyErr *errors.ErrorResponse

	respBody, _ := ioutil.ReadAll(resp.Body)

	err := json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return errors.ErrDefaultCard
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		err := errors.FindCardError(errModel.Code, errModel.Messages...)

		fields["bankly_error"] = bodyErr
		logrus.WithFields(fields).WithError(err).Error("bankly get card error")

		return err
	}

	return errors.ErrDefaultCard
}
