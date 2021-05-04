package bankly

import (
	"github.com/contbank/grok"
	"net/http"
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
	// ErrInvalidCorrelationId
	ErrInvalidCorrelationId = grok.NewError(http.StatusBadRequest, "invalid correlation id")
	// ErrInvalidAmount
	ErrInvalidAmount = grok.NewError(http.StatusBadRequest, "invalid amount")
)

type BanklyError ErrorModel

type Error struct {
	banklyError BanklyError
	grokError 	*grok.Error
}

var errorList = []Error {
	Error {
		banklyError: BanklyError { Code : "INVALID_PERSONAL_BUSINESS_SIZE" },
		grokError : ErrInvalidBusinessSize,
	},
	Error {
		banklyError: BanklyError { Code : "EMAIL_ALREADY_IN_USE" },
		grokError : ErrEmailAlreadyInUse,
	},
	Error {
		banklyError: BanklyError { Code : "PHONE_ALREADY_IN_USE" },
		grokError : ErrPhoneAlreadyInUse,
	},
	Error {
		banklyError: BanklyError { Code : "CUSTOMER_REGISTRATION_CANNOT_BE_REPLACED" },
		grokError : ErrCustomerRegistrationCannotBeReplaced,
	},
	Error {
		banklyError: BanklyError { Code : "ACCOUNT_HOLDER_NOT_EXISTS" },
		grokError : ErrAccountHolderNotExists,
	},
	Error {
		banklyError: BanklyError { Code : "HOLDER_ALREADY_HAVE_A_ACCOUNT" },
		grokError : ErrHolderAlreadyHaveAAccount,
	},
}

type BanklyTransferError KeyValueErrorModel

type TransferError struct {
	banklyTransferError BanklyTransferError
	grokError 			*grok.Error
}

var transferErrorList = []TransferError {
	TransferError {
		banklyTransferError	: BanklyTransferError { Key : "x-correlation-id" },
		grokError 			: ErrInvalidCorrelationId,
	},
	TransferError {
		banklyTransferError	: BanklyTransferError { Key : "$.amount" },
		grokError 			: ErrInvalidAmount,
	},
}

func FindError(errorModel ErrorModel) *grok.Error {
	for _, v := range errorList {
		if v.banklyError.Code == errorModel.Code {
			return v.grokError
		}
	}
	return grok.NewError(http.StatusBadRequest, errorModel.Code + " - " + errorModel.Messages[0])
}

func FindTransferError(errorModel KeyValueErrorModel) *grok.Error {
	for _, v := range transferErrorList {
		if v.banklyTransferError.Key == errorModel.Key {
			return v.grokError
		}
	}
	return grok.NewError(http.StatusBadRequest, errorModel.Key + " - " + errorModel.Value)
}