package bankly

import (
	"time"
)

const (
	// BusinessPath ...
	BusinessPath = "business"
)

// BusinessType ...
type BusinessType string

const (
	BusinessTypeMEI    BusinessType = "MEI"
	BusinessTypeEI     BusinessType = "EI"
	BusinessTypeEIRELI BusinessType = "EIRELI"
)

// BusinessSize ...
type BusinessSize string

const (
	BusinessSizeMEI BusinessSize = "MEI"
	BusinessSizeME  BusinessSize = "ME"
	BusinessSizeEPP BusinessSize = "EPP"
)

// BusinessRequest ...
type BusinessRequest struct {
	Document            string               `validate:"required,cnpj" json:"documentNumber,omitempty"`
	BusinessName        string               `validate:"required" json:"businessName,omitempty"`
	TradingName         string               `json:"tradingName,omitempty"`
	BusinessEmail       string               `json:"businessEmail,omitempty"`
	BusinessType        BusinessType         `validate:"required" json:"businessType"`
	BusinessSize        BusinessSize         `validate:"required" json:"businessSize"`
	BusinessAddress     *Address             `validate:"required,dive" json:"businessAddress,omitempty"`
	LegalRepresentative *LegalRepresentative `validate:"required,dive" json:"legalRepresentative,omitempty"`
}

// BusinessUpdateRequest ...
type BusinessUpdateRequest struct {
	BusinessName        string `validate:"required" json:"businessName,omitempty"`
	TradingName         string `json:"tradingName,omitempty"`
	BusinessEmail       string               `json:"businessEmail,omitempty"`
	BusinessType        BusinessType         `validate:"required" json:"businessType"`
	BusinessSize        BusinessSize         `validate:"required" json:"businessSize"`
	BusinessAddress     *Address             `validate:"required,dive" json:"businessAddress,omitempty"`
	LegalRepresentative *LegalRepresentative `validate:"required,dive" json:"legalRepresentative,omitempty"`
}

// BusinessResponse ...
type BusinessResponse struct {
	ResultLevel  ResultLevel  `json:"resultLevel,omitempty"`
	Document     string       `json:"documentNumber,omitempty"`
	BusinessName  string       `json:"businessName,omitempty"`
	TradingName   string       `json:"tradingName,omitempty"`
	BusinessEmail string       `json:"businessEmail,omitempty"`
	Status        string       `json:"status,omitempty"`
	Reasons      []string     `json:"reasons,omitempty"`
	BusinessType BusinessType `json:"businessType"`
	BusinessSize BusinessSize `json:"businessSize"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
}

//BusinessAccountRequest ...
type BusinessAccountRequest struct {
	Document    string      `validate:"required,cnpj" json:"documentNumber,omitempty"`
	AccountType AccountType `validate:"required" json:"accountType"`
}