package models

import (
	"github.com/contbank/grok"
	"time"
)

type CardType string

// CardResponse ...
const (
	// VirtualCardType cartao virtual
	VirtualCardType CardType = "VIRTUAL"
	// PhysicalCardType cartao fisico
	PhysicalCardType CardType = "PHYSICAL"
	// MultipleCardType cartao multiplo
	MultipleCardType CardType = "MULTIPLE"
)

type CardResponse struct {
	Created          string              `json:"created"`
	CompanyKey       string              `json:"companyKey"`
	DocumentNumber   string              `json:"documentNumber"`
	ActivateCode     string              `json:"activateCode"`
	BankAgency       string              `json:"bankAgency"`
	BankAccount      string              `json:"bankAccount"`
	LastFourDigits   string              `json:"lastFourDigits"`
	Proxy            string              `json:"proxy"`
	Name             string              `json:"name"`
	Alias            string              `json:"alias"`
	CardType         CardType            `json:"cardType"`
	Status           string              `json:"status"`
	PhysicalBinds    []CardBind          `json:"physicalBinds"`
	VirtualBind      CardBind            `json:"virtualBind"`
	AllowContactless bool                `json:"allowContactless"`
	Address          CardAddress         `json:"address"`
	HistoryStatus    []CardHistoryStatus `json:"historyStatus"`
	ActivatedAt      *time.Time          `json:"activatedAt"`
	LastUpdatedAt    time.Time           `json:"lastUpdatedAt"`
	IsActivated      bool                `json:"isActivated"`
	IsLocked         bool                `json:"isLocked"`
	IsCanceled       bool                `json:"isCanceled"`
	IsBuilding       bool                `json:"isBuilding"`
	IsFirtual        bool                `json:"isFirtual"`
	IsPos            bool                `json:"isPos"`
	SettlementDay    int16               `json:"settlementDay"`
}

// CardAddress ...
type CardAddress struct {
	ZipCode      string  `validate:"required" json:"zipCode,omitempty"`
	Address      string  `validate:"required" json:"address,omitempty"`
	Number       string  `validate:"required" json:"number,omitempty"`
	Complement   *string `json:"complement,omitempty"`
	Neighborhood string  `validate:"required" json:"neighborhood,omitempty"`
	City         string  `validate:"required" json:"city,omitempty"`
	State        string  `validate:"required" json:"state,omitempty"`
	Country      string  `validate:"required" json:"country,omitempty"`
}

type CardResponseDTO struct {
	Created          string              `json:"created"`
	CompanyKey       string              `json:"companyKey"`
	DocumentNumber   string              `json:"documentNumber"`
	ActivateCode     string              `json:"activateCode"`
	BankAgency       string              `json:"bankAgency"`
	BankAccount      string              `json:"bankAccount"`
	LastFourDigits   string              `json:"lastFourDigits"`
	Proxy            string              `json:"proxy"`
	Name             string              `json:"name"`
	Alias            string              `json:"alias"`
	CardType         CardType            `json:"cardType"`
	Status           string              `json:"status"`
	PhysicalBinds    []CardBind          `json:"physicalBinds"`
	VirtualBind      CardBind            `json:"virtualBind"`
	AllowContactless bool                `json:"allowContactless"`
	Address          CardAddress         `json:"address"`
	HistoryStatus    []CardHistoryStatus `json:"historyStatus"`
	ActivatedAt      *time.Time          `json:"activatedAt"`
	LastUpdatedAt    time.Time           `json:"lastUpdatedAt"`
	IsActivated      bool                `json:"isActivated"`
	IsLocked         bool                `json:"isLocked"`
	IsCanceled       bool                `json:"isCanceled"`
	IsBuilding       bool                `json:"isBuilding"`
	IsFirtual        bool                `json:"isFirtual"`
	IsPre            bool                `json:"isPre"`
	IsPos            bool                `json:"isPos"`
	IsDebit          bool                `json:"isDebit"`
	PaymentDay       int16               `json:"paymentDay"`
}

type CardBind struct {
	Proxy   string    `json:"proxy"`
	Created time.Time `json:"created"`
}

type CardHistoryStatus struct {
	Modified time.Time `json:"modified"`
	Value    string    `json:"value"`
}

type CardNextStatus struct {
	Value        string `json:"value"`
	IsDefinitive bool   `json:"isDefinitive"`
}

type CardCreateDTO struct {
	CardType CardType `json:"cardType"`
	CardData CardCreateRequest
}

type CardUpdateStatusDTO struct {
	Status           string `json:"status"`
	Password         string `json:"password"`
	UpdateCardBinded bool   `json:"updateCardBinded"`
}

type CardActivateDTO struct {
	Password     string `json:"password"`
	ActivateCode string `json:"activateCode"`
}

type CardUpdatePasswordDTO struct {
	Password string `json:"password"`
}

type CardCreateRequest struct {
	DocumentNumber string      `json:"documentNumber"`
	CardName       string      `json:"cardName"`
	Alias          string      `json:"alias"`
	BankAgency     string      `json:"bankAgency"`
	BankAccount    string      `json:"bankAccount"`
	ProgramId      int16       `json:"programId,omitempty"`
	Password string      `json:"password"`
	Address  CardAddress `json:"address"`
}

type CardCreateResponse struct {
	Proxy        string `json:"proxy"`
	ActivateCode string `json:"activateCode"`
}

type CardTransactionsResponse struct {
	Account struct {
		Number string `json:"number"`
		Agency string `json:"agency"`
	} `json:"account"`
	Amount struct {
		Value  float64 `json:"value"`
		Local  float64 `json:"local"`
		Net    float64 `json:"net"`
		Iof    float64 `json:"iof"`
		Markup float64 `json:"markup"`
	} `json:"amount"`
	Merchant struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		MCC  string `json:"mcc"`
		City string `json:"city"`
	} `json:"merchant"`
	AuthorizationCode    string `json:"authorizationCode"`
	CountryCode          string `json:"countryCode"`
	CurrencyCode         string `json:"currencyCode"`
	EntryMode            string `json:"entryMode"`
	Status               string `json:"status"`
	TransactionTimestamp string `json:"transactionTimestamp"`
	TransactionType      string `json:"transactionType"`
}

// CardPasswordDTO ...
type CardPCIDTO struct {
	Password string `json:"password"`
}

// CardPCIResponse ...
type CardPCIResponse struct {
	CardNumber     string `json:"cardNumber"`
	Cvv            string `json:"cvv"`
	ExpirationDate string `json:"expirationDate"`
}

// CardTrackingResponse ...
type CardTrackingResponse struct {
	CreatedDate           time.Time             `json:"createdDate,omitempty"`
	Name                  string                `json:"name,omitempty"`
	Alias                 string                `json:"alias,omitempty"`
	EstimatedDeliveryDate time.Time             `json:"estimatedDeliveryDate,omitempty"`
	Function         string                `json:"function,omitempty"`
	ExternalTracking CardExternalTracking  `json:"externalTracking,omitempty"`
	Address          []CardTrackingAddress `json:"address,omitempty"`
	Status           []CardTrackingStatus  `json:"status,omitempty"`
	Finalized        []Finalized           `json:"finalized,omitempty"`
}

// CardExternalTracking ...
type CardExternalTracking struct {
	Code    string `json:"code,omitempty"`
	Partner string `json:"partner,omitempty"`
}

// CardTrackingStatus ...
type CardTrackingStatus struct {
	CreatedDate time.Time `json:"createdDate,omitempty"`
	Type        string    `json:"type,omitempty"`
	Reason      string    `json:"reason,omitempty"`
}

// CardTrackingAddress ...
type CardTrackingAddress struct {
	ZipCode      string `json:"zipCode,omitempty"`
	Address      string `json:"address,omitempty"`
	Number       string `json:"number,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
	Complement   string `json:"complement,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	Country      string `json:"country,omitempty"`
	IsActive     bool   `json:"isActive,omitempty"`
}

// Finalized ...
type Finalized []struct {
	CreatedDate      time.Time `json:"createdDate,omitempty"`
	RecipientName    string    `json:"recipientName,omitempty"`
	RecipientKinship string    `json:"recipientKinship,omitempty"`
	DocumentNumber   string    `json:"documentNumber,omitempty"`
	Attempts         int       `json:"attempts,omitempty"`
}

// ParseResponseCard ...
func ParseResponseCard(cardResponseDTO *CardResponseDTO) *CardResponse {
	return &CardResponse{
		Created:          cardResponseDTO.Created,
		CompanyKey:       cardResponseDTO.CompanyKey,
		DocumentNumber:   grok.OnlyDigits(cardResponseDTO.DocumentNumber),
		ActivateCode:     cardResponseDTO.ActivateCode,
		BankAgency:       grok.OnlyLettersOrDigits(cardResponseDTO.BankAgency),
		BankAccount:      grok.OnlyLettersOrDigits(cardResponseDTO.BankAccount),
		LastFourDigits:   cardResponseDTO.LastFourDigits,
		Proxy:            cardResponseDTO.Proxy,
		Name:             grok.ToTitle(cardResponseDTO.Name),
		Alias:            grok.ToTitle(cardResponseDTO.Alias),
		CardType:         cardResponseDTO.CardType,
		Status:           cardResponseDTO.Status,
		PhysicalBinds:    cardResponseDTO.PhysicalBinds,
		VirtualBind:      cardResponseDTO.VirtualBind,
		AllowContactless: cardResponseDTO.AllowContactless,
		Address:          cardResponseDTO.Address,
		HistoryStatus:    cardResponseDTO.HistoryStatus,
		ActivatedAt:      cardResponseDTO.ActivatedAt,
		LastUpdatedAt:    cardResponseDTO.LastUpdatedAt,
		IsActivated:      cardResponseDTO.IsActivated,
		IsLocked:         cardResponseDTO.IsLocked,
		IsCanceled:       cardResponseDTO.IsCanceled,
		IsBuilding:       cardResponseDTO.IsBuilding,
		IsFirtual:        cardResponseDTO.IsFirtual,
		IsPos:            cardResponseDTO.IsPos,
		SettlementDay:    cardResponseDTO.PaymentDay,
	}
}