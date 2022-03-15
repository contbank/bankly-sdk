package models

import "time"

const (
	//BankStatementsPath ...
	BankStatementsPath = "events"
)

// Statement ...
type Statement struct {
	AggregateID    string                 `json:"aggregateId,omitempty"`
	Type           string                 `json:"type,omitempty"`
	Category       string                 `json:"category,omitempty"`
	DocumentNumber string                 `json:"documentNumber,omitempty"`
	Branch         string                 `json:"bankBranch,omitempty"`
	Account        string                 `json:"bankAccount,omitempty"`
	Amount         float64                `json:"amount,omitempty"`
	Name           string                 `json:"name,omitempty"`
	Timestamp      time.Time              `json:"timestamp,omitempty"`
	Data           map[string]interface{} `json:"data,omitempty"`
}

// FilterBankStatementRequest ...
type FilterBankStatementRequest struct {
	Branch         string `validate:"required"`
	Account        string `validate:"required"`
	IncludeDetails bool
	CardProxy      []string
	BeginDateTime  *time.Time
	EndDateTime    *time.Time
	Page           int64 `validate:"required"`
	PageSize       int64 `validate:"required"`
}