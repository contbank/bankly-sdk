package models

import "time"

const (
	// BoletosPath ...
	BoletosPath = "bankslip"
)

//BoletoType ...
type BoletoType string

const (
	Deposit BoletoType = "Deposit"
	Levy    BoletoType = "Levy"
)

// BoletoRequest ...
type BoletoRequest struct {
	Alias       *string    `json:"alias,omitempty"`
	Document    string     `validate:"required,cnpjcpf" json:"documentNumber,omitempty"`
	Amount      float64    `validate:"required" json:"amount,omitempty"`
	DueDate     time.Time  `validate:"required" json:"dueDate,omitempty"`
	EmissionFee bool       `json:"emissionFee,omitempty"`
	Type        BoletoType `validate:"required" json:"type,omitempty"`
	Account     *Account   `validate:"required" json:"account,omitempty"`
	Payer       *Payer     `validate:"required" json:"payer,omitempty"`
}

// SimulatePaymentRequest ...
type SimulatePaymentRequest struct {
	AuthenticationCode string   `validate:"required" json:"authenticationCode,omitempty"`
	Account            *Account `validate:"required" json:"account,omitempty"`
}

// Account ...
type Account struct {
	Number string `validate:"required" json:"number,omitempty"`
	Branch string `validate:"required" json:"branch,omitempty"`
}

// Payer ...
type Payer struct {
	Name      string         `validate:"required" json:"name,omitempty"`
	TradeName string         `json:"tradeName,omitempty"`
	Document  string         `validate:"required,cnpjcpf" json:"document,omitempty"`
	Address   *BoletoAddress `validate:"required" json:"address,omitempty"`
}

type BoletoAddress struct {
	AddressLine string `validate:"required" json:"addressLine,omitempty"`
	ZipCode     string `validate:"required" json:"zipCode,omitempty"`
	State       string `validate:"required" json:"state,omitempty"`
	City        string `validate:"required" json:"city,omitempty"`
}

// BoletoResponse ...
type BoletoResponse struct {
	AuthenticationCode string   `json:"authenticationCode,omitempty"`
	Account            *Account `json:"account,omitempty"`
}

// BoletoAmount ...
type BoletoAmount struct {
	Value    float64 `json:"value,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

// BoletoPayment ...
type BoletoPayment struct {
	ID             string    `json:"id,omitempty"`
	Amount         float64   `json:"amount,omitempty"`
	PaymentChannel string    `json:"paymentChannel,omitempty"`
	PaidOutDate    time.Time `json:"paidOutDate,omitempty"`
}

// BoletoDetailedResponse ...
type BoletoDetailedResponse struct {
	Alias              *string          `json:"alias,omitempty"`
	AuthenticationCode string           `json:"authenticationCode,omitempty"`
	Digitable          string           `json:"digitable,omitempty"`
	Status             string           `json:"status,omitempty"`
	Document           string           `json:"documentNumber,omitempty"`
	DueDate            time.Time        `json:"dueDate,omitempty"`
	EmissionFee        bool             `json:"emissionFee,omitempty"`
	OurNumber       string           `json:"ourNumber,omitempty"`
	Type            BoletoType       `json:"type,omitempty"`
	Amount          *BoletoAmount    `json:"amount,omitempty"`
	Account         *Account         `json:"account,omitempty"`
	Payer           *Payer           `json:"payer,omitempty"`
	RecipientFinal  *Payer           `json:"recipientFinal,omitempty"`
	RecipientOrigin *Payer           `json:"recipientOrigin,omitempty"`
	Payments        []*BoletoPayment `json:"payments,omitempty"`

	// API is returning error for this field
	// EmissionDate time.Time `json:"emissionDate,omitempty"`

	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// FilterBoletoData ...
type FilterBoletoData struct {
	Alias              *string          `json:"alias,omitempty"`
	AuthenticationCode string           `json:"authenticationCode,omitempty"`
	Barcode            string           `json:"barcode,omitempty"`
	Digitable          string           `json:"digitable,omitempty"`
	Status             string           `json:"status,omitempty"`
	DueDate            time.Time        `json:"dueDate,omitempty"`
	Amount             *BoletoAmount    `json:"amount,omitempty"`
	Payer              *Payer           `json:"payer,omitempty"`
	RecipientFinal     *Payer           `json:"recipientFinal,omitempty"`
	RecipientOrigin    *Payer           `json:"recipientOrigin,omitempty"`
	Payments           []*BoletoPayment `json:"payments,omitempty"`
	// API is returning error for this field
	// EmissionDate time.Time `json:"emissionDate,omitempty"`
}

// FilterBoletoResponse ...
type FilterBoletoResponse struct {
	NextPageToken string             `json:"nextPageToken,omitempty"`
	Data          []FilterBoletoData `json:"data,omitempty"`
}

// FindBoletoRequest ...
type FindBoletoRequest struct {
	AuthenticationCode string   `validate:"required" json:"authenticationCode,omitempty"`
	Account            *Account `validate:"required" json:"account,omitempty"`
}

// CancelBoletoRequest ...
type CancelBoletoRequest struct {
	AuthenticationCode string   `validate:"required" json:"authenticationCode,omitempty"`
	Account            *Account `validate:"required" json:"account,omitempty"`
}
