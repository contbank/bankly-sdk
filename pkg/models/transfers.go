package bankly

import (
	"time"
)

const (
	// TransfersPath ...
	TransfersPath = "fund-transfers"
	// InternalBankCode ...
	InternalBankCode string = "332"
)

// TransfersAccountType ...
type TransfersAccountType string

const (
	//CheckingAccount Conta corrente
	CheckingAccount TransfersAccountType = "CHECKING"
	//SavingsAccount Conta Poupança
	SavingsAccount TransfersAccountType = "SAVINGS"
)

// TransfersStatus
type TransfersStatus string

const (
	// Created
	TransfersStatusCreated TransfersStatus = "CREATED"
	// InProcess
	TransfersStatusInProcess TransfersStatus = "IN_PROCESS"
	// Approved
	TransfersStatusApproved TransfersStatus = "APPROVED"
	// Reproved
	TransfersStatusReproved TransfersStatus = "REPROVED"
	// Done
	TransfersStatusDone TransfersStatus = "DONE"
	// Undone
	TransfersStatusUndone TransfersStatus = "UNDONE"
	// Canceled
	TransfersStatusCanceled TransfersStatus = "CANCELED"
)

// TransfersRequest ...
type TransfersRequest struct {
	Amount      int64            `validate:"required" json:"amount"`
	Sender      SenderRequest    `validate:"required,dive" json:"sender"`
	Recipient   RecipientRequest `validate:"required,dive" json:"recipient"`
	Description string           `json:"description"`
}

// SenderRequest ...
type SenderRequest struct {
	Branch   string `validate:"required" json:"branch"`
	Account  string `validate:"required" json:"account"`
	Document string `validate:"required" json:"document"`
	Name     string `validate:"required" json:"name"`
}

// RecipientRequest ...
type RecipientRequest struct {
	TransfersAccountType TransfersAccountType `validate:"required" json:"accountType"`
	BankCode             string               `validate:"required" json:"bankCode"`
	Branch               string               `validate:"required" json:"branch"`
	Account              string               `validate:"required" json:"account"`
	Document             string               `validate:"required" json:"document"`
	Name                 string               `validate:"required" json:"name"`
}

// TransferRequest ...
type TransferRequest struct {
	Amount      float64    `validate:"required" json:"amount,omitempty"`
	Description string     `validate:"required" json:"description,omitempty"`
	Sender      *Sender    `validate:"required,dive" json:"sender,omitempty"`
	Recipient   *Recipient `validate:"required,dive" json:"recipient,omitempty"`
}

// TransfersResponse ...
type TransfersResponse struct {
	ContinuationToken string                   `json:"continuationToken"`
	Data              []TransferByCodeResponse `json:"data"`
}

// Sender ...
type Sender struct {
	Branch   string `validate:"required" json:"branch,omitempty"`
	Account  string `validate:"required" json:"account,omitempty"`
	Document string `validate:"required" json:"document,omitempty"`
	Name     string `validate:"required" json:"name,omitempty"`
}

// Recipient ...
type Recipient struct {
	BankCode    string                `validate:"required" json:"bankCode,omitempty"`
	Branch      string                `validate:"required" json:"branch,omitempty"`
	Account     string                `validate:"required" json:"account,omitempty"`
	Document    string                `validate:"required" json:"document,omitempty"`
	Name        string                `validate:"required" json:"name,omitempty"`
	AccountType *TransfersAccountType `validate:"required,dive" json:"accountType,omitempty"`
}

// TransferByCodeResponse ...
type TransferByCodeResponse struct {
	CompanyKey         string             `json:"companyKey"`
	AuthenticationCode string             `json:"authenticationCode"`
	Amount             float64            `json:"amount"`
	CorrelationId      string             `json:"correlationId"`
	Sender             *SenderResponse    `json:"sender"`
	Recipient          *RecipientResponse `json:"recipient"`
	Channel            string             `json:"channel"`
	Operation          string             `json:"operation"`
	Identifier         string             `json:"identifier"`
	Status             TransfersStatus    `json:"status"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
}

// AccountResponse ...
type AccountResponse struct {
	Balance *BalanceRespone `json:"balance,omitempty"`
	Status  string          `json:"status,omitempty"`
	Branch  string          `json:"branch,omitempty"`
	Number  string          `json:"number,omitempty"`
	Bank    *BankData       `json:"bank,omitempty"`
}

// SenderResponse ...
type SenderResponse struct {
	Document string           `json:"document,omitempty"`
	Name     string           `json:"name,omitempty"`
	Account  *AccountResponse `json:"account,omitempty"`
}

// RecipientResponse ...
type RecipientResponse struct {
	Document string           `json:"document,omitempty"`
	Name     string           `json:"name,omitempty"`
	Account  *AccountResponse `json:"account,omitempty"`
}
