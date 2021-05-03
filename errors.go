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
)

type BanklyError ErrorModel

type Errors struct {
	banklyError BanklyError
	grokError *grok.Error
}

var errorsList = []Errors {
	Errors {
		banklyError: BanklyError { Code : "INVALID_PERSONAL_BUSINESS_SIZE" },
		grokError : ErrInvalidBusinessSize,
	},
	Errors {
		banklyError: BanklyError { Code : "EMAIL_ALREADY_IN_USE" },
		grokError : ErrEmailAlreadyInUse,
	},
	Errors {
		banklyError: BanklyError { Code : "PHONE_ALREADY_IN_USE" },
		grokError : ErrPhoneAlreadyInUse,
	},
	Errors {
		banklyError: BanklyError { Code : "CUSTOMER_REGISTRATION_CANNOT_BE_REPLACED" },
		grokError : ErrCustomerRegistrationCannotBeReplaced,
	},
	Errors {
		banklyError: BanklyError { Code : "ACCOUNT_HOLDER_NOT_EXISTS" },
		grokError : ErrAccountHolderNotExists,
	},
	Errors {
		banklyError: BanklyError { Code : "HOLDER_ALREADY_HAVE_A_ACCOUNT" },
		grokError : ErrHolderAlreadyHaveAAccount,
	},
}

func FindError(errorModel ErrorModel) *grok.Error {
	for _, v := range errorsList {
		if v.banklyError.Code == errorModel.Code {
			return v.grokError
		}
	}
	return grok.NewError(http.StatusBadRequest, errorModel.Code + " - " + errorModel.Messages[0])
}