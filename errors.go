package bankly

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/contbank/grok"
)

var (
	// ErrEntryNotFound ...
	ErrEntryNotFound = grok.NewError(http.StatusNotFound, "not found")
	// ErrDuplicateCompany ...
	ErrDuplicateCompany = grok.NewError(http.StatusConflict, "duplicate company")
	// ErrInvalidToken ...
	ErrInvalidToken = grok.NewError(http.StatusConflict, "invalid token")
	// ErrInvalidBusinessSize ...
	ErrInvalidBusinessSize = grok.NewError(http.StatusBadRequest, "invalid business size")
	// ErrEmailAlreadyInUse ...
	ErrEmailAlreadyInUse = grok.NewError(http.StatusBadRequest, "email already in use")
	// ErrPhoneAlreadyInUse ...
	ErrPhoneAlreadyInUse = grok.NewError(http.StatusBadRequest, "phone already in use")
	// ErrCustomerRegistrationCannotBeReplaced ...
	ErrCustomerRegistrationCannotBeReplaced = grok.NewError(http.StatusConflict, "customer registration cannot be replaced")
	// ErrAccountHolderNotExists ...
	ErrAccountHolderNotExists = grok.NewError(http.StatusBadRequest, "account holder not exists")
	// ErrHolderAlreadyHaveAAccount ...
	ErrHolderAlreadyHaveAAccount = grok.NewError(http.StatusConflict, "holder already have a account")
	// ErrInvalidCorrelationID ...
	ErrInvalidCorrelationID = grok.NewError(http.StatusBadRequest, "invalid correlation id")
	// ErrInvalidAmount ...
	ErrInvalidAmount = grok.NewError(http.StatusBadRequest, "invalid amount")
	// ErrInsufficientBalance ...
	ErrInsufficientBalance = grok.NewError(http.StatusBadRequest, "insufficient balance")
	// ErrInvalidAuthenticationCodeOrAccount ...
	ErrInvalidAuthenticationCodeOrAccount = grok.NewError(http.StatusBadRequest, "invalid authentication code or account number")
	// ErrInvalidAccountNumber ...
	ErrInvalidAccountNumber = grok.NewError(http.StatusBadRequest, "invalid account number")
	// ErrOutOfServicePeriod ...
	ErrOutOfServicePeriod = grok.NewError(http.StatusBadRequest, "out of service period")
	// ErrCashoutLimitNotEnough ...
	ErrCashoutLimitNotEnough = grok.NewError(http.StatusBadRequest, "cashout limit not enough")
	// ErrInvalidParameter ...
	ErrInvalidParameter = grok.NewError(http.StatusBadRequest, "invalid parameter")
	// ErrInvalidParameterLength ...
	ErrInvalidParameterLength = grok.NewError(http.StatusBadRequest, "invalid parameter length")
	// ErrInvalidAddressNumberLength ...
	ErrInvalidAddressNumberLength = grok.NewError(http.StatusBadRequest, "invalid address number length")
	// ErrInvalidRegisterNameLength ...
	ErrInvalidRegisterNameLength = grok.NewError(http.StatusBadRequest, "invalid register name length")
	// ErrInvalidParameterSpecialCharacters ...
	ErrInvalidParameterSpecialCharacters = grok.NewError(http.StatusBadRequest, "invalid parameter with special characters")
	// ErrInvalidSocialNameLength ...
	ErrInvalidSocialNameLength = grok.NewError(http.StatusBadRequest, "invalid social name length")
	// ErrInvalidEmailLength ...
	ErrInvalidEmailLength = grok.NewError(http.StatusBadRequest, "invalid email length")
	// ErrInvalidAPIEndpoint ...
	ErrInvalidAPIEndpoint = grok.NewError(http.StatusBadRequest, "invalid api endpoint")
	// ErrMethodNotAllowed ...
	ErrMethodNotAllowed = grok.NewError(http.StatusMethodNotAllowed, "method not allowed")
	// ErrSendDocumentAnalysis ...
	ErrSendDocumentAnalysis = grok.NewError(http.StatusMethodNotAllowed, "send document analysis error")
	// ErrGetDocumentAnalysis ...
	ErrGetDocumentAnalysis = grok.NewError(http.StatusMethodNotAllowed, "get document analysis error")
	// ErrScouterQuantity ...
	ErrScouterQuantity = grok.NewError(http.StatusUnprocessableEntity, "max boleto amount per day reached")
	// ErrBoletoInvalidStatus ...
	ErrBoletoInvalidStatus = grok.NewError(http.StatusUnprocessableEntity, "boleto was in an invalid status")
	// ErrBarcodeNotFound ...
	ErrBarcodeNotFound = grok.NewError(http.StatusNotFound, "bar code not found")
	// ErrPaymentInvalidStatus ...
	ErrPaymentInvalidStatus = grok.NewError(http.StatusUnprocessableEntity, "payment was in an invalid status")
	// ErrDefaultTransfers ...
	ErrDefaultTransfers = grok.NewError(http.StatusConflict, "error transfers")
	// ErrDefaultFindTransfers ...
	ErrDefaultFindTransfers = grok.NewError(http.StatusConflict, "error find transfers")
	// ErrDefaultPayment ...
	ErrDefaultPayment = grok.NewError(http.StatusInternalServerError, "error payment")
	// ErrDefaultBusinessAccounts ...
	ErrDefaultBusinessAccounts = grok.NewError(http.StatusInternalServerError, "error business accounts")
	// ErrDefaultCustomersAccounts ...
	ErrDefaultCustomersAccounts = grok.NewError(http.StatusInternalServerError, "error customers accounts")
	// ErrDefaultBalance ...
	ErrDefaultBalance = grok.NewError(http.StatusInternalServerError, "error balance")
	// ErrDefaultLogin ...
	ErrDefaultLogin = grok.NewError(http.StatusInternalServerError, "error login")
	// ErrDefaultBank ...
	ErrDefaultBank = grok.NewError(http.StatusInternalServerError, "error bank")
	// ErrDefaultBankStatements ...
	ErrDefaultBankStatements = grok.NewError(http.StatusInternalServerError, "error bank statements")
	//ErrDefaultBoletos ...
	ErrDefaultBoletos = grok.NewError(http.StatusInternalServerError, "error bank boletos")
	//ErrClientIDClientSecret ...
	ErrClientIDClientSecret = grok.NewError(http.StatusInternalServerError, "error client id or client secret")
	// ErrDefaultFreshDesk ...
	ErrDefaultFreshDesk = grok.NewError(http.StatusInternalServerError, "error in fresh desk api")
	// ErrFreshDeskTicketNotFound ...
	ErrFreshDeskTicketNotFound = grok.NewError(http.StatusNotFound, "error in fresh desk ticket not found")
	// ErrInvalidRecipientBranch ...
	ErrInvalidRecipientBranch = grok.NewError(http.StatusConflict, "invalid recipient branch number")
	// ErrInvalidRecipientAccount ...
	ErrInvalidRecipientAccount = grok.NewError(http.StatusConflict, "invalid recipient account number")
	// ErrDefaultCard ...
	ErrDefaultCard = grok.NewError(http.StatusInternalServerError, "error card")
	// ErrDefaultPix ...
	ErrDefaultPix = grok.NewError(http.StatusInternalServerError, "error pix")
	// ErrKeyNotFound ...
	ErrKeyNotFound = grok.NewError(http.StatusNotFound, "key not found")
	// ErrInvalidQrCodePayload ...
	ErrInvalidQrCodePayload = grok.NewError(http.StatusConflict, "invalid qrcode payload")
	// ErrInvalidKeyType ...
	ErrInvalidKeyType = grok.NewError(http.StatusUnprocessableEntity, "invalid key type")
	// ErrInvalidParameterPix ...
	ErrInvalidParameterPix = grok.NewError(http.StatusUnprocessableEntity, "invalid parameter")
	// ErrInsufficientBalancePix ...
	ErrInsufficientBalancePix = grok.NewError(http.StatusConflict, "invalid parameter")
	// ErrInvalidAccountType ...
	ErrInvalidAccountType = grok.NewError(http.StatusUnprocessableEntity, "invalid parameter")
	// ErrCardActivate ...
	ErrCardActivate = grok.NewError(http.StatusNotModified, "error card activate")
	// ErrCardStatusUpdate ...
	ErrCardStatusUpdate = grok.NewError(http.StatusNotModified, "error update status card")
	// ErrCardPasswordUpdate ...
	ErrCardPasswordUpdate = grok.NewError(http.StatusNotModified, "error update password card")
	// ErrInvalidPassword ...
	ErrInvalidPassword = grok.NewError(http.StatusUnauthorized, "invalid password")
	// ErrInvalidCardName ...
	ErrInvalidCardName = grok.NewError(http.StatusBadRequest, "invalid card name")
	// ErrInvalidIdentifier ...
	ErrInvalidIdentifier = grok.NewError(http.StatusBadRequest, "invalid identifier")
	// ErrCardAlreadyActivated ...
	ErrCardAlreadyActivated = grok.NewError(http.StatusConflict, "card already activated")
	// ErrOperationNotAllowedCardStatus ...
	ErrOperationNotAllowedCardStatus = grok.NewError(http.StatusMethodNotAllowed, "operation not allowed for current card status")
)

// BanklyError ...
type BanklyError ErrorModel

// Error ..
type Error struct {
	ErrorKey  string
	GrokError *grok.Error
}

type ErrorCard struct {
	Code         string
	Messages     []string
	Metadata     interface{}
	PropertyName string
	Reasons      []interface{}
}

type BanklyCardError struct {
	ErrorsCard ErrorCard
}

var errorList = []Error{
	{
		ErrorKey:  "INVALID_PERSONAL_BUSINESS_SIZE",
		GrokError: ErrInvalidBusinessSize,
	},
	{
		ErrorKey:  "EMAIL_ALREADY_IN_USE",
		GrokError: ErrEmailAlreadyInUse,
	},
	{
		ErrorKey:  "PHONE_ALREADY_IN_USE",
		GrokError: ErrPhoneAlreadyInUse,
	},
	{
		ErrorKey:  "CUSTOMER_REGISTRATION_CANNOT_BE_REPLACED",
		GrokError: ErrCustomerRegistrationCannotBeReplaced,
	},
	{
		ErrorKey:  "ACCOUNT_HOLDER_NOT_EXISTS",
		GrokError: ErrAccountHolderNotExists,
	},
	{
		ErrorKey:  "HOLDER_ALREADY_HAVE_A_ACCOUNT",
		GrokError: ErrHolderAlreadyHaveAAccount,
	},
	{
		ErrorKey:  "SCOUTER_QUANTITY",
		GrokError: ErrScouterQuantity,
	},
	{
		ErrorKey:  "BANKSLIP_SETTLEMENT_STATUS_VALIDATE",
		GrokError: ErrBoletoInvalidStatus,
	},
	{
		ErrorKey:  "BAR_CODE_NOT_FOUND",
		GrokError: ErrBarcodeNotFound,
	},
	{
		ErrorKey:  "INVALID_PARAMETER",
		GrokError: ErrInvalidParameter,
	},
	{
		ErrorKey:  "INVALID_PARAMETER_LENGTH",
		GrokError: ErrInvalidParameterLength,
	},
	{
		ErrorKey:  "INVALID_PARAMETER_SPECIAL_CHARACTERS",
		GrokError: ErrInvalidParameterSpecialCharacters,
	},
	{
		ErrorKey:  "INVALID_ADDRESS_NUMBER_LENGTH",
		GrokError: ErrInvalidAddressNumberLength,
	},
	{
		ErrorKey:  "INVALID_REGISTER_NAME_LENGTH",
		GrokError: ErrInvalidRegisterNameLength,
	},
	{
		ErrorKey:  "INVALID_SOCIAL_NAME_LENGTH",
		GrokError: ErrInvalidSocialNameLength,
	},
	{
		ErrorKey:  "INVALID_EMAIL_LENGTH",
		GrokError: ErrInvalidEmailLength,
	},
}

// BanklyTransferError ..
type BanklyTransferError KeyValueErrorModel

// TransferError ..
type TransferError struct {
	banklyTransferError BanklyTransferError
	grokError           *grok.Error
}

var transferErrorList = []TransferError{
	{
		banklyTransferError: BanklyTransferError{Key: "x-correlation-id"},
		grokError:           ErrInvalidCorrelationID,
	},
	{
		banklyTransferError: BanklyTransferError{Key: "$.amount"},
		grokError:           ErrInvalidAmount,
	},
	{
		banklyTransferError: BanklyTransferError{Key: "INSUFFICIENT_BALANCE"},
		grokError:           ErrInsufficientBalance,
	},
	{
		banklyTransferError: BanklyTransferError{Key: "CASH_OUT_NOT_ALLOWED_OUT_OF_BUSINESS_PERIOD"},
		grokError:           ErrOutOfServicePeriod,
	},
	{
		banklyTransferError: BanklyTransferError{Key: "CASHOUT_LIMIT_NOT_ENOUGH"},
		grokError:           ErrCashoutLimitNotEnough,
	},
	{
		banklyTransferError: BanklyTransferError{Key: "Recipient.Branch"},
		grokError:           ErrInvalidRecipientBranch,
	},
	{
		banklyTransferError: BanklyTransferError{Key: "Recipient.Account"},
		grokError:           ErrInvalidRecipientAccount,
	},
}

// FindError Find errors.
func FindError(code string, messages ...string) *Error {
	code = verifyInvalidParameter(code, messages)

	for _, v := range errorList {
		if v.ErrorKey == code {
			return &v
		}
	}

	return &Error{
		ErrorKey:  code,
		GrokError: grok.NewError(http.StatusConflict, messages...),
	}
}

// FindErrorByErrorModel ..
func FindErrorByErrorModel(response ErrorModel) *Error {
	if response.Code != "" {
		return FindError(response.Code, response.Messages...)
	}
	return &Error{
		ErrorKey:  response.Key,
		GrokError: grok.NewError(http.StatusBadRequest, response.Value),
	}
}

// verifyInvalidParameter Find the correspondent error message.
func verifyInvalidParameter(code string, messages []string) string {
	if code == "INVALID_PARAMETER" {
		for _, m := range messages {
			if strings.Contains(strings.ToLower(m), "length of 'building number'") {
				return "INVALID_ADDRESS_NUMBER_LENGTH"
			} else if strings.Contains(strings.ToLower(m), "length of 'register name'") {
				return "INVALID_REGISTER_NAME_LENGTH"
			} else if strings.Contains(strings.ToLower(m), "length of 'social name'") {
				return "INVALID_SOCIAL_NAME_LENGTH"
			} else if strings.Contains(strings.ToLower(m), "length of 'email'") {
				return "INVALID_EMAIL_LENGTH"
			} else if strings.Contains(strings.ToLower(m), "not allowed to include numbers or special characters") {
				return "INVALID_PARAMETER_SPECIAL_CHARACTERS"
			} else if strings.Contains(strings.ToLower(m), "length of") {
				return "INVALID_PARAMETER_LENGTH"
			}
		}
	}
	return code
}

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
	}
	return code
}

var errorPixList = []Error{
	{
		ErrorKey:  "ENTRY_NOT_FOUND",
		GrokError: ErrKeyNotFound,
	},
	{
		ErrorKey:  "INVALID_QRCODE_PAYLOAD_CONTENT_TO_DECODE",
		GrokError: ErrInvalidQrCodePayload,
	},
	{
		ErrorKey:  "INVALID_KEY_TYPE",
		GrokError: ErrInvalidKeyType,
	},
	{
		ErrorKey:  "INVALID_PARAMETER_PIX",
		GrokError: ErrInvalidParameterPix,
	},
	{
		ErrorKey:  "INSUFFICIENT_BALANCE",
		GrokError: ErrInsufficientBalancePix,
	},
	{
		ErrorKey:  "INVALID_ACCOUNT_TYPE",
		GrokError: ErrInvalidAccountType,
	},
}

func verifyInvalidPixParameter(code string, messages []string) string {
	if code == "INVALID_PARAMETER" {
		for _, m := range messages {
			switch {
			case strings.Contains(strings.ToLower(m), "addressing key value does not match with addressing key type"):
				return "INVALID_KEY_TYPE"
			case strings.Contains(strings.ToLower(m), "sender.account.type"):
				return "INVALID_ACCOUNT_TYPE"
			default:
				return "INVALID_PARAMETER_PIX"
			}
		}
	}
	return code
}

// FindPixError
func FindPixError(code string, messages ...string) *grok.Error {
	code = verifyInvalidPixParameter(code, messages)

	for _, v := range errorPixList {
		if v.ErrorKey == code {
			return v.GrokError
		}
	}

	return grok.NewError(http.StatusConflict, messages...)
}

// ParseErr ..
func ParseErr(err error) (*Error, bool) {
	banklyErr, ok := err.(*Error)
	return banklyErr, ok
}

// FindTransferError ..
func FindTransferError(transferErrorResponse TransferErrorResponse) *grok.Error {
	// get the error code if errors list is null
	if len(transferErrorResponse.Errors) == 0 && transferErrorResponse.Code != "" {
		transferErrorResponse.Errors = []KeyValueErrorModel{
			{
				Key: transferErrorResponse.Code,
			},
		}
	}
	// checking the errors list
	errorModel := transferErrorResponse.Errors[0]
	for _, v := range transferErrorList {
		if v.banklyTransferError.Key == errorModel.Key {
			return v.grokError
		}
	}
	return grok.NewError(http.StatusBadRequest, errorModel.Key+" - "+errorModel.Value)
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"Key: %s - Messages: %s",
		e.ErrorKey,
		strings.Join(e.GrokError.Messages, "\n"),
	)
}
