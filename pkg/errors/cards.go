package errors

import (
	"github.com/contbank/grok"
	"net/http"
	"strings"
)

var errorCardList = []Error{
	{
		ErrorKey:  "INVALID_CARD_PASSWORD",
		GrokError: ErrInvalidPassword,
	},
	{
		ErrorKey:  "OPERATION_NOT_ALLOWED_FOR_CURRENT_CARD_STATUS",
		GrokError: ErrOperationNotAllowedCardStatus,
	},
	{
		ErrorKey:  "CARD_ALREADY_ACTIVATED",
		GrokError: ErrCardAlreadyActivated,
	},
	{
		ErrorKey:  "CARD_SERVICE_UNAVAILABLE",
		GrokError: ErrCardServiceUnavailable,
	},
	{
		ErrorKey:  "INVALID_CARD_NAME_EMPTY",
		GrokError: ErrInvalidCardName,
	},
	{
		ErrorKey:  "INVALID_DOCUMENT_NUMBER_EMPTY",
		GrokError: ErrInvalidIdentifier,
	},
	{
		ErrorKey:  "INVALID_PARAMETER_CARD",
		GrokError: ErrInvalidParameter,
	},
}

// FindCardError Find cards errors.
func FindCardError(code string, messages ...string) *grok.Error {
	code = verifyInvalidCardParameter(code, messages)

	for _, v := range errorCardList {
		if v.ErrorKey == code {
			return v.GrokError
		}
	}

	return grok.NewError(http.StatusConflict, messages...)
}

// verifyInvalidCardParameter Find the correspondent error message for Cards.
func verifyInvalidCardParameter(code string, messages []string) string {
	if code == "INVALID_PARAMETER" {
		for _, m := range messages {
			switch {
			case strings.Contains(strings.ToLower(m), "card name"):
				return "INVALID_CARD_NAME_EMPTY"
			case strings.Contains(strings.ToLower(m), "document number"):
				return "INVALID_DOCUMENT_NUMBER_EMPTY"
			default:
				return "INVALID_PARAMETER_CARD"
			}
		}
	} else if code == "009" {
		return "OPERATION_NOT_ALLOWED_FOR_CURRENT_CARD_STATUS"
	} else if code == "011" {
		return "INVALID_CARD_PASSWORD"
	} else if code == "021" {
		return "CARD_ALREADY_ACTIVATED"
	} else if code == "300" {
		return "CARD_SERVICE_UNAVAILABLE"
	}
	return code
}

// CardErrorResponse ...
type CardErrorResponse struct {
	ErrorKey string `json:"errorKey,omitempty"`
	CodeMessageErrorResponse
}

// ErrorCard ...
type ErrorCard struct {
	Code         string
	Messages     []string
	Metadata     interface{}
	PropertyName string
	Reasons      []interface{}
}

// BanklyCardError ...
type BanklyCardError struct {
	ErrorsCard ErrorCard
}