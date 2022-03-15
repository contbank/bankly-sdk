package models

import "time"

const (
	// PaymentPath ...
	PaymentPath = "bill-payment"
)

// ValidatePaymentRequest ...
type ValidatePaymentRequest struct {
	Code string `validate:"required" json:"code,omitempty"`
}

// ValidatePaymentResponse ...
type ValidatePaymentResponse struct {
	ID                string         `json:"id,omitempty"`
	Assignor          string         `json:"assignor,omitempty"`
	Code              string         `json:"code,omitempty"`
	Digitable         string         `json:"digitable,omitempty"`
	Amount            float64        `json:"amount,omitempty"`
	OriginalAmount    float64        `json:"originalAmount,omitempty"`
	MinAmount         float64        `json:"minAmount,omitempty"`
	MaxAmount         float64        `json:"maxAmount,omitempty"`
	AllowChangeAmount bool           `json:"allowChangeAmount,omitempty"`
	DueDate           string         `json:"dueDate,omitempty"`
	SettleDate        string         `json:"settleDate,omitempty"`
	NextSettle        bool           `json:"nextSettle,omitempty"`
	Payer             *PaymentPayer  `json:"payer,omitempty"`
	Recipient         *PaymentPayer  `json:"recipient,omitempty"`
	BusinessHours     *BusinessHours `json:"businessHours,omitempty"`
	Charges           *Charges       `json:"charges,omitempty"`
}

// PaymentPayer ...
type PaymentPayer struct {
	Name           string `json:"name,omitempty"`
	DocumentNumber string `json:"documentNumber,omitempty"`
}

// BusinessHours ...
type BusinessHours struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

// Charges ...
type Charges struct {
	InterestAmountCalculated float64 `json:"interestAmountCalculated,omitempty"`
	FineAmountCalculated     float64 `json:"fineAmountCalculated,omitempty"`
	DiscountAmount           float64 `json:"discountAmount,omitempty"`
}

// ConfirmPaymentRequest ...
type ConfirmPaymentRequest struct {
	ID          string  `validate:"required" json:"id,omitempty"`
	Amount      float64 `validate:"required" json:"amount,omitempty"`
	Description *string `json:"description,omitempty"`
	BankBranch  string  `validate:"required" json:"bankBranch,omitempty"`
	BankAccount string  `validate:"required" json:"bankAccount,omitempty"`
}

// ConfirmPaymentResponse ...
type ConfirmPaymentResponse struct {
	AuthenticationCode string    `json:"authenticationCode,omitempty"`
	SettledDate        time.Time `json:"settledDate,omitempty"`
}

// FilterPaymentsRequest ...
type FilterPaymentsRequest struct {
	BankBranch  string `validate:"required"`
	BankAccount string `validate:"required"`
	PageSize    int    `validate:"required"`
	PageToken   *string
}

// PaymentResponse ...
type PaymentResponse struct {
	AuthenticationCode string     `json:"authenticationCode,omitempty"`
	Status             string     `json:"status,omitempty"`
	Digitable          string     `json:"digitable,omitempty"`
	Description        *string    `json:"description,omitempty"`
	BankBranch         string     `json:"bankBranch,omitempty"`
	BankAccount        string     `json:"bankAccount,omitempty"`
	RecipientDocument  string     `json:"recipientDocument,omitempty"`
	RecipientName      string     `json:"recipientName,omitempty"`
	Amount             float64    `json:"amount,omitempty"`
	OriginalAmount     float64    `json:"originalAmount,omitempty"`
	Assignor           string    `json:"assignor,omitempty"`
	Charges            *Charges  `json:"charges,omitempty"`
	SettleDate         time.Time `json:"settleDate,omitempty"`
	PaymentDate        time.Time  `json:"paymentDate,omitempty"`
	ConfirmedAt        time.Time  `json:"confirmedAt,omitempty"`
	DueDate            *time.Time `json:"dueDate,omitempty"`
	CompanyKey         *string    `json:"companyKey,omitempty"`
	DocumentNumber     *string    `json:"documentNumber,omitempty"`
}

// FilterPaymentsResponse ...
type FilterPaymentsResponse struct {
	NextPageToken string             `json:"nextPage,omitempty"`
	Data          []*PaymentResponse `json:"data,omitempty"`
}

// DetailPaymentRequest ...
type DetailPaymentRequest struct {
	BankBranch         string `validate:"required"`
	BankAccount        string `validate:"required"`
	AuthenticationCode string `validate:"required"`
}