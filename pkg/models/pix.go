package bankly

import (
	"time"
)

type InitializationType string

// InitializationType ...
const (
	Manual        InitializationType = "Manual"
	Key           InitializationType = "Key"
	StaticQrCode  InitializationType = "StaticQrCode"
	DynamicQrCode InitializationType = "DynamicQrCode"
)

type PixAddressKeyResponse struct {
	EndToEndID    string       `json:"endToEndId"`
	AddressingKey PixTypeValue `json:"addressingKey"`
	Holder        PixHolder    `json:"holder"`
	Status        string       `json:"status"`
	CreatedAt     time.Time    `json:"createdAt"`
	OwnedAt       time.Time    `json:"ownedAt"`
}

type PixTypeValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type PixHolder struct {
	Type       string       `json:"type"`
	Name       string       `json:"name"`
	SocialName string       `json:"socialName"`
	Document   PixTypeValue `json:"document"`
}

type PixCashOutRequest struct {
	Sender             PixCashOutSenderRequest    `json:"sender"`
	Recipient          PixCashOutRecipientRequest `json:"recipient"`
	Amount             float64                    `json:"amount"`
	Description        string                     `json:"description"`
	InitializationType InitializationType         `json:"initializationType"`
	EndToEndID         string                     `json:"endToEndId"`
}

type PixCashOutAccountRequest struct {
	Branch string `json:"branch"`
	Number string `json:"number"`
}

type PixCashOutBankRequest struct {
	Ispb string `json:"ispb"`
}

type PixCashOutSenderRequest struct {
	Account        PixCashOutAccountRequest `json:"account"`
	Bank           PixCashOutBankRequest    `json:"bank"`
	DocumentNumber string                   `json:"documentNumber"`
	Name           string                   `json:"name"`
}

type PixCashOutRecipientRequest struct {
	Account        PixCashOutAccountRequest `json:"account"`
	Bank           PixCashOutBankRequest    `json:"bank"`
	DocumentNumber string                   `json:"documentNumber"`
	Name           string                   `json:"name"`
}

type PixCashOutAccountResponse struct {
	Branch string `json:"branch"`
	Number string `json:"number"`
	Type   string `json:"type"`
}

type PixCashOutBankResponse struct {
	Ispb  string `json:"ispb"`
	Compe string `json:"compe"`
	Name  string `json:"name"`
}

type PixCashOutSenderResponse struct {
	Account        PixCashOutAccountResponse `json:"account"`
	Bank           PixCashOutBankResponse    `json:"bank"`
	DocumentNumber string                    `json:"documentNumber"`
	Name           string                    `json:"name"`
}

type PixCashOutRecipientResponse struct {
	Account        PixCashOutAccountResponse `json:"account"`
	Bank           PixCashOutBankResponse    `json:"bank"`
	DocumentNumber string                    `json:"documentNumber"`
	Name           string                    `json:"name"`
}

type PixCashOutResponse struct {
	Amount             float64                     `json:"amount"`
	Description        string                      `json:"description"`
	Sender             PixCashOutSenderResponse    `json:"sender"`
	Recipient          PixCashOutRecipientResponse `json:"recipient"`
	AuthenticationCode string                      `json:"authenticationCode"`
}

type PixQrCodeDecodeRequest struct {
	EncodedValue string `json:"encodedValue"`
}

type PixQrCodeBankResponse struct {
	Name string `json:"name"`
}

type PixQrCodePaymentResponse struct {
	BaseValue       float64 `json:"baseValue"`
	InterestValue   float64 `json:"interestValue"`
	PenaltyValue    float64 `json:"penaltyValue"`
	DiscountValue   float64 `json:"discountValue"`
	TotalValue      float64 `json:"totalValue"`
	DueDate         string  `json:"dueDate"`
	ChangeValue     float64 `json:"changeValue"`
	WithdrawalValue float64 `json:"withdrawalValue"`
}

type PixQrCodeLocationResponse struct {
	City    string `json:"city"`
	ZipCode string `json:"zipCode"`
}

type PixQrCodeDecodeResponse struct {
	EndToEndID     string                    `json:"endToEndId"`
	ConciliationID string                    `json:"conciliationId"`
	AddressingKey  PixTypeValue              `json:"addressingKey"`
	QrCodeType     string                    `json:"qrCodeType"`
	Holder         PixHolder                 `json:"holder"`
	Bank           PixQrCodeBankResponse     `json:"bank"`
	Payment        PixQrCodePaymentResponse  `json:"payment"`
	Location       PixQrCodeLocationResponse `json:"location"`
	QrCodePurpose  string                    `json:"qrCodePurpose"`
}

type PixCashOutByAuthenticationCodeResponse struct {
	CompanyKey         string                      `json:"companyKey"`
	AuthenticationCode string                      `json:"authenticationCode"`
	InitializationType string                      `json:"initializationType"`
	Amount             float64                     `json:"amount"`
	CorrelationID      string                      `json:"correlationId"`
	Sender             PixCashOutSenderResponse    `json:"sender"`
	Recipient          PixCashOutRecipientResponse `json:"recipient"`
	Channel            string                      `json:"channel"`
	Status             TransfersStatus             `json:"status"`
	Type               string                      `json:"type"`
	CreatedAt          time.Time                   `json:"createdAt"`
	UpdatedAt          time.Time                   `json:"updatedAt"`
}
