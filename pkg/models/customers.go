package models

import (
	"time"
)

const (
	// CustomersPath ..
	CustomersPath = "customers"
)

// CustomerStatus
type CustomerStatus string

const (
	// CustomerStatusPendingApproval
	CustomerStatusPendingApproval CustomerStatus = "PENDING_APPROVAL"
	// CustomerStatusApproved
	CustomerStatusApproved CustomerStatus = "APPROVED"
	// CustomerStatusReproved
	CustomerStatusReproved CustomerStatus = "REPROVED"
	// CustomerStatusCanceled
	CustomerStatusCanceled CustomerStatus = "CANCELED"
	// CustomerStatusBlacklisted
	CustomerStatusBlacklisted CustomerStatus = "BLACKLISTED"
)

// CustomersRequest ...
type CustomersRequest struct {
	Document     string    		  `validate:"required,cpf" json:"documentNumber,omitempty"`
	RegisterName string    		  `validate:"required" json:"registerName,omitempty"`
	SocialName   string    `json:"socialName,omitempty"`
	Phone        *Phone    `validate:"required,dive" json:"phone,omitempty"`
	Address      *Address  `validate:"required,dive" json:"address,omitempty"`
	BirthDate    time.Time `validate:"required" json:"birthDate,omitempty"`
	MotherName   string    		  `validate:"required" json:"motherName,omitempty"`
	Email        string    		  `validate:"required" json:"email,omitempty"`
}

// CustomersAccountRequest ...
type CustomersAccountRequest struct {
	AccountType AccountType `validate:"required" json:"accountType"`
}

// CustomersResponse ...
type CustomersResponse struct {
	DocumentNumber  string          `json:"documentNumber"`
	RegisterName    string          `json:"registerName"`
	SocialName      string          `json:"socialName"`
	Email      string         `json:"email"`
	Phone      Phone          `json:"phone"`
	Address    Address        `json:"address"`
	MotherName string         `json:"motherName"`
	IsPoliticallyExposedPerson bool `json:"isPoliticallyExposedPerson"`
	Reasons    []string       `json:"reasons"`
	Status     CustomerStatus `json:"status"`
	Profile    string         `json:"profile"`
}

// CustomerUpdateRequest ...
type CustomerUpdateRequest struct {
	RegisterName string    `validate:"required" json:"registerName,omitempty"`
	SocialName   string    `json:"socialName,omitempty"`
	Phone        *Phone    `validate:"required,dive" json:"phone,omitempty"`
	Address      *Address  `validate:"required,dive" json:"address,omitempty"`
	BirthDate    time.Time `validate:"required" json:"birthDate,omitempty"`
	MotherName   string    `validate:"required" json:"motherName,omitempty"`
	Email        string    `validate:"required" json:"email,omitempty"`
}