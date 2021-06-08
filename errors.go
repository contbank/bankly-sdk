package bankly

import (
	"net/http"

	"github.com/contbank/grok"
)

const ScouterQuantityCode = "SCOUTER_QUANTITY"

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
	// ErrInvalidCorrelationId ...
	ErrInvalidCorrelationId = grok.NewError(http.StatusBadRequest, "invalid correlation id")
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
)

// BanklyError ...
type BanklyError ErrorModel

// Error ..
type Error struct {
	banklyError BanklyError
	grokError   *grok.Error
}

var errorList = []Error{
	{
		banklyError: BanklyError{Code: "INVALID_PERSONAL_BUSINESS_SIZE"},
		grokError:   ErrInvalidBusinessSize,
	},
	{
		banklyError: BanklyError{Code: "EMAIL_ALREADY_IN_USE"},
		grokError:   ErrEmailAlreadyInUse,
	},
	{
		banklyError: BanklyError{Code: "PHONE_ALREADY_IN_USE"},
		grokError:   ErrPhoneAlreadyInUse,
	},
	{
		banklyError: BanklyError{Code: "CUSTOMER_REGISTRATION_CANNOT_BE_REPLACED"},
		grokError:   ErrCustomerRegistrationCannotBeReplaced,
	},
	{
		banklyError: BanklyError{Code: "ACCOUNT_HOLDER_NOT_EXISTS"},
		grokError:   ErrAccountHolderNotExists,
	},
	{
		banklyError: BanklyError{Code: "HOLDER_ALREADY_HAVE_A_ACCOUNT"},
		grokError:   ErrHolderAlreadyHaveAAccount,
	},
	{
		banklyError: BanklyError{Code: "INVALID_PARAMETER"},
		grokError:   ErrInvalidParameter,
	},
	{
		banklyError: BanklyError{Code: ScouterQuantityCode},
		grokError:   ErrScouterQuantity,
	},
	{
		banklyError: BanklyError{Code: "BANKSLIP_SETTLEMENT_STATUS_VALIDATE"},
		grokError:   ErrBoletoInvalidStatus,
	},
	{
		banklyError: BanklyError{Code: "BAR_CODE_NOT_FOUND"},
		grokError:   ErrBarcodeNotFound,
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
		grokError:           ErrInvalidCorrelationId,
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
}

// FindError ..
func FindError(errorModel ErrorModel) *grok.Error {
	for _, v := range errorList {
		if v.banklyError.Code == errorModel.Code {
			return v.grokError
		}
	}

	var message string
	if len(errorModel.Messages) > 0 {
		message = errorModel.Messages[0]
	}

	return grok.NewError(http.StatusBadRequest, errorModel.Code+" - "+message)
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
