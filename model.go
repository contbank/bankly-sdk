package bankly

import "time"

const (
	//LoginPath ..
	LoginPath = "connect/token"
	//CustomersPath ..
	CustomersPath = "customers"
	//AccountsPath ...
	AccountsPath = "accounts"
	//TransfersPath ...
	TransfersPath = "fund-transfers"
	//BusinessPath ...
	BusinessPath = "business"
	//BanksPath ...
	BanksPath = "banklist"
)

//AuthenticationResponse ...
type AuthenticationResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

//ErrorLoginResponse ...
type ErrorLoginResponse struct {
	Message string `json:"error"`
}

//CustomersRequest ...
type CustomersRequest struct {
	Documment    string    `validate:"required,cpf" json:"documentNumber,omitempty"`
	RegisterName string    `validate:"required" json:"registerName,omitempty"`
	SocialName   string    `json:"socialName,omitempty"`
	Phone        *Phone    `validate:"required,dive" json:"phone,omitempty"`
	Address      *Address  `validate:"required,dive" json:"address,omitempty"`
	BirthDate    time.Time `validate:"required" json:"birthDate,omitempty"`
	MotherName   string    `validate:"required" json:"motherName,omitempty"`
	Email        string    `validate:"required" json:"email,omitempty"`
}

//CustomersAccountRequest ...
type CustomersAccountRequest struct {
	AccountType AccountType `validate:"required" json:"accountType"`
}

//AccountType ...
type AccountType string

const (
	//PaymentAccount ...
	PaymentAccount AccountType = "PAYMENT_ACCOUNT"
)

//CustomersResponse ...
type CustomersResponse struct {
	Phone                      Phone     `json:"phone"`
	Address                    Address   `json:"address"`
	Email                      string    `json:"email"`
	MotherName                 string    `json:"motherName"`
	BirthDate                  time.Time `json:"birthDate"`
	IsPoliticallyExposedPerson bool      `json:"isPoliticallyExposedPerson"`
	Reasons                    []string  `json:"reasons"`
	DocumentNumber             string    `json:"documentNumber"`
	RegisterName               string    `json:"registerName"`
	SocialName                 string    `json:"socialName"`
	Status                     string    `json:"status"`
	Profile                    string    `json:"profile"`
	CreatedAt                  time.Time `json:"createdAt"`
	UpdatedAt                  time.Time `json:"updatedAt"`
}

//Phone ...
type Phone struct {
	CountryCode string `validate:"required" json:"countryCode,omitempty"`
	Number      string `validate:"required" json:"number,omitempty"`
}

//Address ...
type Address struct {
	ZipCode        string `validate:"required" json:"zipCode,omitempty"`
	AddressLine    string `validate:"required" json:"addressLine,omitempty"`
	BuildingNumber string `validate:"required" json:"buildingNumber,omitempty"`
	Complement     string `json:"complement,omitempty"`
	Neighborhood   string `validate:"required" json:"neighborhood,omitempty"`
	City           string `validate:"required" json:"city,omitempty"`
	State          string `validate:"required" json:"state,omitempty"`
	Country        string `validate:"required" json:"country,omitempty"`
}

//ErrorResponse ...
type ErrorResponse struct {
	Errors    []ErrorModel `json:"errors,omitempty"`
	Reference string       `json:"reference,omitempty"`
}

//ErrorModel ...
type ErrorModel struct {
	Code         string   `json:"code,omitempty"`
	PropertyName string   `json:"propertyName,omitempty"`
	Messages     []string `json:"messages,omitempty"`
}

//AccountResponse ...
type AccountResponse struct {
	Balance *BalanceRespone `json:"balance,omitempty"`
	Status  string          `json:"status,omitempty"`
	Branch  string          `json:"branch,omitempty"`
	Number  string          `json:"number,omitempty"`
}

//BalanceRespone ...
type BalanceRespone struct {
	InProcess BalanceValue `json:"inProcess,omitempty"`
	Available BalanceValue `json:"available,omitempty"`
	Blocked   BalanceValue `json:"blocked,omitempty"`
}

//BalanceValue ...
type BalanceValue struct {
	Amount   float64 `json:"amount,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

//TransfersRequest ...
type TransfersRequest struct {
	Amount      int64            `validate:"required" json:"amount"`
	Sender      SenderRequest    `validate:"required,dive" json:"sender"`
	Recipient   RecipientRequest `validate:"required,dive" json:"recipient"`
	Description string           `json:"description"`
}

//SenderRequest ...
type SenderRequest struct {
	Branch   string `validate:"required" json:"branch"`
	Account  string `validate:"required" json:"account"`
	Document string `validate:"required" json:"document"`
	Name     string `validate:"required" json:"name"`
}

//RecipientRequest ...
type RecipientRequest struct {
	TransfersType TransfersType `validate:"required" json:"accountType"`
	BankCode      string        `validate:"required" json:"bankCode"`
	Branch        string        `validate:"required" json:"branch"`
	Account       string        `validate:"required" json:"account"`
	Document      string        `validate:"required" json:"document"`
	Name          string        `validate:"required" json:"name"`
}

//TransfersType ...
type TransfersType string

const (
	//CheckingAccount Conta corrente
	CheckingAccount TransfersType = "CHECKING"
	//SavingsAccount Conta Poupança
	SavingsAccount TransfersType = "SAVINGS"
)

//TransfersResponse ...
type TransfersResponse struct {
	AuthenticationCode string `json:"authenticationCode"`
}

//BusinessRequest ...
type BusinessRequest struct {
	Document            string              `validate:"required,cnpj" json:"documentNumber,omitempty"`
	BusinessName        string              `validate:"required" json:"businessName,omitempty"`
	TradingName         string              `json:"tradingName,omitempty"`
	BusinessEmail       string              `json:"businessEmail,omitempty"`
	BusinessType        BusinessType        `validate:"required" json:"businessType"`
	BusinessSize        BusinessSize        `validate:"required" json:"businessSize"`
	BusinessAddress     *Address            `validate:"required,dive" json:"businessAddress,omitempty"`
	LegalRepresentative LegalRepresentative `validate:"required,dive" json:"legalRepresentative,omitempty"`
}

//BusinessType ...
type BusinessType string

const (
	BusinessTypeMEI    BusinessType = "MEI"
	BusinessTypeEI     BusinessType = "EI"
	BusinessTypeEIRELI BusinessType = "EIRELI"
)

//BusinessSize ...
type BusinessSize string

const (
	BusinessSizeMEI BusinessSize = "MEI"
	BusinessSizeME  BusinessSize = "ME"
	BusinessSizeEPP BusinessSize = "EPP"
)

//BusinessSize ...
type ResultLevel string

const (
	ResultLevelBasic    ResultLevel = "BASIC"
	ResultLevelDetailed ResultLevel = "DETAILED"
)

//LegalRepresentative ...
type LegalRepresentative struct {
	Documment    string    `validate:"required,cpf" json:"documentNumber,omitempty"`
	RegisterName string    `validate:"required" json:"registerName,omitempty"`
	SocialName   string    `json:"socialName,omitempty"`
	Phone        *Phone    `validate:"required,dive" json:"phone,omitempty"`
	Address      *Address  `validate:"required,dive" json:"address,omitempty"`
	BirthDate    time.Time `validate:"required" json:"birthDate,omitempty"`
	MotherName   string    `validate:"required" json:"motherName,omitempty"`
	Email        string    `validate:"required" json:"email,omitempty"`
}

//BusinessResponse ...
type BusinessResponse struct {
	ResultLevel  ResultLevel  `json:"resultLevel,omitempty"`
	Document     string       `json:"documentNumber,omitempty"`
	BusinessName string       `json:"businessName,omitempty"`
	TradingName  string       `json:"tradingName,omitempty"`
	Status       string       `json:"status,omitempty"`
	BusinessType BusinessType `json:"businessType"`
	BusinessSize BusinessSize `json:"businessSize"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

//BusinessAccountRequest ...
type BusinessAccountRequest struct {
	Document    string      `validate:"required,cnpj" json:"documentNumber,omitempty"`
	AccountType AccountType `validate:"required" json:"accountType"`
}

//FilterBankListRequest ...
type FilterBankListRequest struct {
	IDs      []string
	Name     *string
	Product  *string
	Page     *int
	PageSize *int
}

//BankDataResponse ...
type BankDataResponse struct {
	Name        string   `json:"name,omitempty"`
	ISPB        string   `json:"ispb,omitempty"`
	Code        string   `json:"code,omitempty"`
	ShortName   string   `json:"shortName,omitempty"`
	IsSPIDirect bool     `json:"isSPIDirect,omitempty"`
	Products    []string `json:"products,omitempty"`
}
