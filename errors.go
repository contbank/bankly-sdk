package bankly

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/contbank/grok"
)

var (
	// ErrEntryNotFound ...
	ErrEntryNotFound = grok.NewError(http.StatusNotFound, "NOT_FOUND", "not found")
	// ErrDuplicateCompany ...
	ErrDuplicateCompany = grok.NewError(http.StatusConflict, "DUPLICATE_COMPANY", "duplicate company")
	// ErrInvalidToken ...
	ErrInvalidToken = grok.NewError(http.StatusConflict, "INVALID_TOKEN", "invalid token")
	// ErrInvalidBusinessSize ...
	ErrInvalidBusinessSize = grok.NewError(http.StatusBadRequest, "INVALID_BUSINESS_SIZE", "invalid business size")
	// ErrEmailAlreadyInUse ...
	ErrEmailAlreadyInUse = grok.NewError(http.StatusBadRequest, "EXISTS_EMAIL", "email already in use")
	// ErrPhoneAlreadyInUse ...
	ErrPhoneAlreadyInUse = grok.NewError(http.StatusBadRequest, "EXISTS_PHONE", "phone already in use")
	// ErrCustomerRegistrationCannotBeReplaced ...
	ErrCustomerRegistrationCannotBeReplaced = grok.NewError(http.StatusConflict, "CUSTOMER_CANNOT_BE_REPLACED", "customer registration cannot be replaced")
	// ErrAccountHolderNotExists ...
	ErrAccountHolderNotExists = grok.NewError(http.StatusBadRequest, "NOT_EXISTS_HOLDER", "account holder not exists")
	// ErrHolderAlreadyHaveAAccount ...
	ErrHolderAlreadyHaveAAccount = grok.NewError(http.StatusConflict, "EXISTS_HOLDER", "holder already have a account")
	// ErrInvalidCorrelationID ...
	ErrInvalidCorrelationID = grok.NewError(http.StatusBadRequest, "INVALID_CORRELATION_ID", "invalid correlation id")
	// ErrInvalidAmount ...
	ErrInvalidAmount = grok.NewError(http.StatusBadRequest, "INVALID_AMOUNT", "invalid amount")
	// ErrInsufficientBalance ...
	ErrInsufficientBalance = grok.NewError(http.StatusBadRequest, "INSUFFICIENT_BALANCE", "insufficient balance")
	// ErrInvalidAuthenticationCodeOrAccount ...
	ErrInvalidAuthenticationCodeOrAccount = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER", "invalid authentication code or account number")
	// ErrInvalidAccountNumber ...
	ErrInvalidAccountNumber = grok.NewError(http.StatusBadRequest, "INVALID_ACCOUNT_NUMBER", "invalid account number")
	// ErrOutOfServicePeriod ...
	ErrOutOfServicePeriod = grok.NewError(http.StatusBadRequest, "OUT_SERVICE_PERIOD", "out of service period")
	// ErrCashoutLimitNotEnough ...
	ErrCashoutLimitNotEnough = grok.NewError(http.StatusBadRequest, "CASHOUT_LIMIT_NOT_ENOUGH", "cashout limit not enough")
	// ErrInvalidParameter ...
	ErrInvalidParameter = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER", "invalid parameter")
	// ErrInvalidParameterLength ...
	ErrInvalidParameterLength = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER", "invalid parameter length")
	// ErrInvalidAddressNumberLength ...
	ErrInvalidAddressNumberLength = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER", "invalid address number length")
	// ErrInvalidRegisterNameLength ...
	ErrInvalidRegisterNameLength = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER", "invalid register name length")
	// ErrInvalidParameterSpecialCharacters ...
	ErrInvalidParameterSpecialCharacters = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER", "invalid parameter with special characters")
	// ErrInvalidSocialNameLength ...
	ErrInvalidSocialNameLength = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER", "invalid social name length")
	// ErrInvalidEmailLength ...
	ErrInvalidEmailLength = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER", "invalid email length")
	// ErrInvalidAPIEndpoint ...
	ErrInvalidAPIEndpoint = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER", "invalid api endpoint")
	// ErrMethodNotAllowed ...
	ErrMethodNotAllowed = grok.NewError(http.StatusMethodNotAllowed, "INVALID_PARAMETER", "method not allowed")
	// ErrSendDocumentAnalysis ...
	ErrSendDocumentAnalysis = grok.NewError(http.StatusMethodNotAllowed, "INVALID_PARAMETER", "send document analysis error")
	// ErrGetDocumentAnalysis ...
	ErrGetDocumentAnalysis = grok.NewError(http.StatusMethodNotAllowed, "INVALID_PARAMETER", "get document analysis error")
	// ErrScouterQuantity ...
	ErrScouterQuantity = grok.NewError(http.StatusUnprocessableEntity, "INVALID_PARAMETER", "max boleto amount per day reached")
	// ErrBoletoInvalidStatus ...
	ErrBoletoInvalidStatus = grok.NewError(http.StatusUnprocessableEntity, "INVALID_PARAMETER", "boleto was in an invalid status")
	// ErrBarcodeNotFound ...
	ErrBarcodeNotFound = grok.NewError(http.StatusNotFound, "INVALID_PARAMETER", "bar code not found")
	// ErrPaymentInvalidStatus ...
	ErrPaymentInvalidStatus = grok.NewError(http.StatusUnprocessableEntity, "INVALID_PARAMETER", "payment was in an invalid status")
	// ErrDefaultTransfers ...
	ErrDefaultTransfers = grok.NewError(http.StatusConflict, "INVALID_PARAMETER", "error transfers")
	// ErrDefaultFindTransfers ...
	ErrDefaultFindTransfers = grok.NewError(http.StatusConflict, "INVALID_PARAMETER", "error find transfers")
	// ErrDefaultPayment ...
	ErrDefaultPayment = grok.NewError(http.StatusInternalServerError, "INVALID_PARAMETER", "error payment")
	// ErrDefaultBusinessAccounts ...
	ErrDefaultBusinessAccounts = grok.NewError(http.StatusInternalServerError, "INVALID_PARAMETER", "error business accounts")
	// ErrDefaultCustomersAccounts ...
	ErrDefaultCustomersAccounts = grok.NewError(http.StatusInternalServerError, "INVALID_PARAMETER", "error customers accounts")
	// ErrDefaultBalance ...
	ErrDefaultBalance = grok.NewError(http.StatusInternalServerError, "INVALID_PARAMETER", "error balance")
	// ErrDefaultLogin ...
	ErrDefaultLogin = grok.NewError(http.StatusInternalServerError, "INVALID_PARAMETER", "error login")
	// ErrDefaultBank ...
	ErrDefaultBank = grok.NewError(http.StatusInternalServerError, "INVALID_PARAMETER", "error bank")
	// ErrDefaultBankStatements ...
	ErrDefaultBankStatements = grok.NewError(http.StatusInternalServerError, "INVALID_PARAMETER", "error bank statements")
	// ErrDefaultIncomeReport ...
	ErrDefaultIncomeReport = grok.NewError(http.StatusInternalServerError, "INVALID_PARAMETER", "error income report")
	//ErrDefaultBoletos ...
	ErrDefaultBoletos = grok.NewError(http.StatusInternalServerError, "INVALID_PARAMETER", "error bank boletos")
	// ErrDefaultFreshDesk ...
	ErrDefaultFreshDesk = grok.NewError(http.StatusInternalServerError, "INVALID_PARAMETER", "error in fresh desk api")
	// ErrFreshDeskTicketNotFound ...
	ErrFreshDeskTicketNotFound = grok.NewError(http.StatusNotFound, "INVALID_PARAMETER", "error in fresh desk ticket not found")
	// ErrInvalidRecipientBranch ...
	ErrInvalidRecipientBranch = grok.NewError(http.StatusConflict, "INVALID_PARAMETER", "invalid recipient branch number")
	// ErrInvalidRecipientAccount ...
	ErrInvalidRecipientAccount = grok.NewError(http.StatusConflict, "INVALID_PARAMETER", "invalid recipient account number")
	// ErrDefaultCard ...
	ErrDefaultCard = grok.NewError(http.StatusInternalServerError, "INVALID_PARAMETER", "error card")
	// ErrDefaultPix ...
	ErrDefaultPix = grok.NewError(http.StatusInternalServerError, "INVALID_PARAMETER", "error pix")
	// ErrKeyNotFound ...
	ErrKeyNotFound = grok.NewError(http.StatusNotFound, "INVALID_PARAMETER", "key not found")
	// ErrInvalidQrCodePayload ...
	ErrInvalidQrCodePayload = grok.NewError(http.StatusConflict, "INVALID_PARAMETER", "invalid qrcode payload")
	// ErrInvalidKeyType ...
	ErrInvalidKeyType = grok.NewError(http.StatusUnprocessableEntity, "INVALID_PARAMETER", "invalid key type")
	// ErrInvalidParameterPix ...
	ErrInvalidParameterPix = grok.NewError(http.StatusUnprocessableEntity, "INVALID_PARAMETER", "invalid parameter")
	// ErrInsufficientBalancePix ...
	ErrInsufficientBalancePix = grok.NewError(http.StatusConflict, "INVALID_PARAMETER", "invalid parameter")
	// ErrInvalidAccountType ...
	ErrInvalidAccountType = grok.NewError(http.StatusUnprocessableEntity, "INVALID_PARAMETER", "invalid parameter")
	// ErrCardActivate ...
	ErrCardActivate = grok.NewError(http.StatusNotModified, "INVALID_PARAMETER", "error card activate")
	// ErrCardStatusUpdate ...
	ErrCardStatusUpdate = grok.NewError(http.StatusNotModified, "INVALID_PARAMETER", "error update status card")
	// ErrCardPasswordUpdate ...
	ErrCardPasswordUpdate = grok.NewError(http.StatusNotModified, "INVALID_PARAMETER", "error update password card")
	// ErrInvalidPassword ...
	ErrInvalidPassword = grok.NewError(http.StatusUnauthorized, "INVALID_PARAMETER", "invalid password")
	// ErrInvalidCardName ...
	ErrInvalidCardName = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER", "invalid card name")
	// ErrInvalidIdentifier ...
	ErrInvalidIdentifier = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER", "invalid identifier")
	// ErrCardAlreadyActivated ...
	ErrCardAlreadyActivated = grok.NewError(http.StatusConflict, "INVALID_PARAMETER", "card already activated")
	// ErrOperationNotAllowedCardStatus ...
	ErrOperationNotAllowedCardStatus = grok.NewError(http.StatusMethodNotAllowed, "INVALID_PARAMETER", "operation not allowed for current card status")
	// ErrInvalidIncomeReportCalendar ...
	ErrInvalidIncomeReportCalendar = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER", "invalid income report calendar")
	// ErrInvalidIncomeReportParameter ...
	ErrInvalidIncomeReportParameter = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER", "invalid income report parameter")
	// ErrDefaultCancelCustomersAccounts ...
	ErrDefaultCancelCustomersAccounts = grok.NewError(http.StatusConflict, "INVALID_PARAMETER", "error cancel customers accounts")
	// ErrAccountNonZeroBalance ...
	ErrAccountNonZeroBalance = grok.NewError(http.StatusConflict, "INVALID_PARAMETER", "error account non zero balance")
	// ErrAccountAlreadyBeenCanceled ...
	ErrAccountAlreadyBeenCanceled = grok.NewError(http.StatusUnprocessableEntity, "INVALID_PARAMETER", "error account already been canceled")
	// ErrAccountNotFound ...
	ErrAccountNotFound = grok.NewError(http.StatusNotFound, "INVALID_PARAMETER", "error account not found")
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
	{
		ErrorKey:  "HOLDER_HAS_SOME_ACCOUNTS_WITH_NON_ZERO_BALANCE",
		GrokError: ErrAccountNonZeroBalance,
	},
	{
		ErrorKey:  "HOLDER_HAS_ALREADY_BEEN_CANCELED",
		GrokError: ErrAccountAlreadyBeenCanceled,
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
		GrokError: grok.NewError(http.StatusConflict, code, messages...),
	}
}

// FindErrorByErrorModel ..
func FindErrorByErrorModel(response ErrorModel) *Error {
	if response.Code != "" {
		return FindError(response.Code, response.Messages...)
	}
	return &Error{
		ErrorKey:  response.Key,
		GrokError: grok.NewError(http.StatusBadRequest, response.Key, response.Value),
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

// errorIncomeReportList ...
var errorIncomeReportList = []Error{
	{
		ErrorKey:  "INVALID_CALENDAR_FOR_INCOME_REPORT",
		GrokError: ErrInvalidIncomeReportCalendar,
	},
	{
		ErrorKey:  "INVALID_PARAMETER_INCOME_REPORT",
		GrokError: ErrInvalidIncomeReportParameter,
	},
}

// FindIncomeReportError Find income report errors.
func FindIncomeReportError(code string, messages ...string) *grok.Error {
	code = verifyInvalidIncomeReportParameter(code, messages)

	for _, v := range errorCardList {
		if v.ErrorKey == code {
			return v.GrokError
		}
	}

	return grok.NewError(http.StatusConflict, code, messages...)
}

// verifyInvalidIncomeReportParameter Find the correspondent error message for income reports.
func verifyInvalidIncomeReportParameter(code string, messages []string) string {
	if code == "CALENDAR_NOT_ALLOWED" {
		for _, m := range messages {
			switch {
			case strings.Contains(strings.ToLower(m), "calendar informed is not allowed"):
				return "INVALID_CALENDAR_FOR_INCOME_REPORT"
			default:
				return "INVALID_PARAMETER_INCOME_REPORT"
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

	return grok.NewError(http.StatusConflict, code, messages...)
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

	return grok.NewError(http.StatusConflict, code, messages...)
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
	return grok.NewError(http.StatusBadRequest, errorModel.Key, errorModel.Key+" - "+errorModel.Value)
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"Key: %s - Messages: %s",
		e.ErrorKey,
		strings.Join(e.GrokError.Messages, "\n"),
	)
}
