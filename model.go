package bankly

import (
	"time"
)

const (
	// LoginPath ..
	LoginPath = "connect/token"
	// CustomersPath ..
	CustomersPath = "customers"
	// AccountsPath ...
	AccountsPath = "accounts"
	// TransfersPath ...
	TransfersPath = "fund-transfers"
	// BusinessPath ...
	BusinessPath = "business"
	// BoletosPath ...
	BoletosPath = "bankslip"
	// BanksPath ...
	BanksPath = "banklist"
	// DocumentAnalysisPath ...
	DocumentAnalysisPath = "/document-analysis"
)

const (
	// InternalBankCode ...
	InternalBankCode string = "332"
)

// AuthenticationResponse ...
type AuthenticationResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// ErrorLoginResponse ...
type ErrorLoginResponse struct {
	Message string `json:"error"`
}

// CustomersRequest ...
type CustomersRequest struct {
	Document     string    `validate:"required,cpf" json:"documentNumber,omitempty"`
	RegisterName string    `validate:"required" json:"registerName,omitempty"`
	SocialName   string    `json:"socialName,omitempty"`
	Phone        *Phone    `validate:"required,dive" json:"phone,omitempty"`
	Address      *Address  `validate:"required,dive" json:"address,omitempty"`
	BirthDate    time.Time `validate:"required" json:"birthDate,omitempty"`
	MotherName   string    `validate:"required" json:"motherName,omitempty"`
	Email        string    `validate:"required" json:"email,omitempty"`
}

// CustomersAccountRequest ...
type CustomersAccountRequest struct {
	AccountType AccountType `validate:"required" json:"accountType"`
}

// AccountType ...
type AccountType string

const (
	// PaymentAccount ...
	PaymentAccount AccountType = "PAYMENT_ACCOUNT"
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

// CustomersResponse ...
type CustomersResponse struct {
	DocumentNumber             string    `json:"documentNumber"`
	RegisterName               string    `json:"registerName"`
	SocialName                 string    `json:"socialName"`
	Email                      string    `json:"email"`
	Phone                      Phone     `json:"phone"`
	Address                    Address   `json:"address"`
	MotherName                 string    `json:"motherName"`
	BirthDate                  time.Time `json:"birthDate"`
	IsPoliticallyExposedPerson bool      `json:"isPoliticallyExposedPerson"`
	Reasons                    []string  `json:"reasons"`
	Status                     string    `json:"status"`
	Profile                    string    `json:"profile"`
	CreatedAt                  time.Time `json:"createdAt"`
	UpdatedAt                  time.Time `json:"updatedAt"`
}

// Phone ...
type Phone struct {
	CountryCode string `validate:"required" json:"countryCode,omitempty"`
	Number      string `validate:"required" json:"number,omitempty"`
}

// Address ...
type Address struct {
	ZipCode        string  `validate:"required" json:"zipCode,omitempty"`
	AddressLine    string  `validate:"required" json:"addressLine,omitempty"`
	BuildingNumber string  `validate:"required" json:"buildingNumber,omitempty"`
	Complement     *string `json:"complement,omitempty"`
	Neighborhood   string  `validate:"required" json:"neighborhood,omitempty"`
	City           string  `validate:"required" json:"city,omitempty"`
	State          string  `validate:"required" json:"state,omitempty"`
	Country        string  `validate:"required" json:"country,omitempty"`
}

// Account ...
type Account struct {
	Number string `validate:"required" json:"number,omitempty"`
	Branch string `validate:"required" json:"branch,omitempty"`
}

// Payer ...
type Payer struct {
	Name      string   `validate:"required" json:"name,omitempty"`
	TradeName string   `validate:"required" json:"tradeName,omitempty"`
	Document  string   `validate:"required,cpfcnpj" json:"document,omitempty"`
	Address   *Address `validate:"required" json:"address,omitempty"`
}

// ErrorResponse ...
type ErrorResponse struct {
	Errors    []ErrorModel `json:"errors,omitempty"`
	Reference string       `json:"reference,omitempty"`
}

//BoletoErrorResponse ...
type BoletoErrorResponse struct {
	Code    string  `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

// TransferErrorResponse ...
type TransferErrorResponse struct {
	Code            string               `json:"code,omitempty"`
	Message         string               `json:"message,omitempty"`
	Layer           string               `json:"layer,omitempty"`
	ApplicationName string               `json:"applicationName,omitempty"`
	Errors          []KeyValueErrorModel `json:"errors,omitempty"`
}

// ErrorModel ...
type ErrorModel struct {
	Code         string   `json:"code,omitempty"`
	PropertyName string   `json:"propertyName,omitempty"`
	Messages     []string `json:"messages,omitempty"`
}

// KeyValueErrorModel ...
type KeyValueErrorModel struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// AccountResponse ...
type AccountResponse struct {
	Balance *BalanceRespone `json:"balance,omitempty"`
	Status  string          `json:"status,omitempty"`
	Branch  string          `json:"branch,omitempty"`
	Number  string          `json:"number,omitempty"`
	Bank    *Bank           `json:"bank,omitempty"`
}

// BalanceRespone ...
type BalanceRespone struct {
	InProcess BalanceValue `json:"inProcess,omitempty"`
	Available BalanceValue `json:"available,omitempty"`
	Blocked   BalanceValue `json:"blocked,omitempty"`
}

// BalanceValue ...
type BalanceValue struct {
	Amount   float64 `json:"amount,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

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

// TransfersAccountType ...
type TransfersAccountType string

const (
	//CheckingAccount Conta corrente
	CheckingAccount TransfersAccountType = "CHECKING"
	//SavingsAccount Conta Poupança
	SavingsAccount TransfersAccountType = "SAVINGS"
)

// TransfersResponse ...
type TransfersResponse struct {
	ContinuationToken string                   `json:"continuationToken"`
	Data              []TransferByCodeResponse `json:"data"`
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
	Status             string             `json:"status"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
}

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
	BusinessName        string               `validate:"required" json:"businessName,omitempty"`
	TradingName         string               `json:"tradingName,omitempty"`
	BusinessEmail       string               `json:"businessEmail,omitempty"`
	BusinessType        BusinessType         `validate:"required" json:"businessType"`
	BusinessSize        BusinessSize         `validate:"required" json:"businessSize"`
	BusinessAddress     *Address             `validate:"required,dive" json:"businessAddress,omitempty"`
	LegalRepresentative *LegalRepresentative `validate:"required,dive" json:"legalRepresentative,omitempty"`
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

// BusinessSize ...
type ResultLevel string

const (
	ResultLevelBasic    ResultLevel = "BASIC"
	ResultLevelDetailed ResultLevel = "DETAILED"
	ResultLevelOnlyStatus ResultLevel = "ONLY_STATUS"
)

// LegalRepresentative ...
type LegalRepresentative struct {
	Document     string    `validate:"required,cpf" json:"documentNumber,omitempty"`
	RegisterName string    `validate:"required" json:"registerName,omitempty"`
	SocialName   string    `json:"socialName,omitempty"`
	Phone        *Phone    `validate:"required,dive" json:"phone,omitempty"`
	Address      *Address  `validate:"required,dive" json:"address,omitempty"`
	BirthDate    time.Time `validate:"required" json:"birthDate,omitempty"`
	MotherName   string    `validate:"required" json:"motherName,omitempty"`
	Email        string    `validate:"required" json:"email,omitempty"`
}

// BusinessResponse ...
type BusinessResponse struct {
	ResultLevel   ResultLevel  `json:"resultLevel,omitempty"`
	Document      string       `json:"documentNumber,omitempty"`
	BusinessName  string       `json:"businessName,omitempty"`
	TradingName   string       `json:"tradingName,omitempty"`
	BusinessEmail string       `json:"businessEmail,omitempty"`
	Status        string       `json:"status,omitempty"`
	BusinessType  BusinessType `json:"businessType"`
	BusinessSize  BusinessSize `json:"businessSize"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
}

//BusinessAccountRequest ...
type BusinessAccountRequest struct {
	Document    string      `validate:"required,cnpj" json:"documentNumber,omitempty"`
	AccountType AccountType `validate:"required" json:"accountType"`
}

//BoletoType ...
type BoletoType string

const (
	Deposit BoletoType = "Deposit"
	Levy    BoletoType = "Levy"
)

//BoletoRequest ...
type BoletoRequest struct {
	Alias       *string    `json:"alias,omitempty"`
	Document    string     `validate:"required,cpfcnpj" json:"documentNumber,omitempty"`
	Amount      float64    `validate:"required" json:"amount,omitempty"`
	DueDate     time.Time  `validate:"required" json:"dueDate,omitempty"`
	EmissionFee bool       `json:"emissionFee,omitempty"`
	Type        BoletoType `validate:"required" json:"type,omitempty"`
	Account     *Account   `validate:"required" json:"account,omitempty"`
	Payer       *Payer     `validate:"required" json:"payer,omitempty"`
}

//BoletoResponse ...
type BoletoResponse struct {
	AuthenticationCode string   `json:"authenticationCode,omitempty"`
	Account            *Account `json:"account,omitempty"`
}

//BoletoAmount ...
type BoletoAmount struct {
	Value    float64 `json:"value,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

//BoletoPayment ...
type BoletoPayment struct {
	ID             string    `json:"id,omitempty"`
	Amount         float64   `json:"amount,omitempty"`
	PaymentChannel string    `json:"paymentChannel,omitempty"`
	PaidOutDate    time.Time `json:"paidOutDate,omitempty"`
}

//BoletoDetailedResponse ...
type BoletoDetailedResponse struct {
	Alias              *string          `json:"alias,omitempty"`
	AuthenticationCode string           `json:"authenticationCode,omitempty"`
	Digitable          string           `json:"digitable,omitempty"`
	Status             string           `json:"status,omitempty"`
	Document           string           `json:"documentNumber,omitempty"`
	DueDate            time.Time        `json:"dueDate,omitempty"`
	EmissionFee        bool             `json:"emissionFee,omitempty"`
	OurNumber          string           `json:"ourNumber,omitempty"`
	Type               BoletoType       `json:"type,omitempty"`
	Amount             *BoletoAmount    `json:"amount,omitempty"`
	Account            *Account         `json:"account,omitempty"`
	Payer              *Payer           `json:"payer,omitempty"`
	RecipientFinal     *Payer           `json:"recipientFinal,omitempty"`
	RecipientOrigin    *Payer           `json:"recipientOrigin,omitempty"`
	Payments           []*BoletoPayment `json:"payments,omitempty"`

	// API is returning error for this field
	// EmissionDate time.Time `json:"emissionDate,omitempty"`

	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

//FilterBoletoData ...
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

//FilterBoletoResponse ...
type FilterBoletoResponse struct {
	NextPageToken string             `json:"nextPageToken,omitempty"`
	Data          []FilterBoletoData `json:"data,omitempty"`
}

//FindBoletoRequest ...
type FindBoletoRequest struct {
	AuthenticationCode string   `validate:"required" json:"authenticationCode,omitempty"`
	Account            *Account `validate:"required" json:"account,omitempty"`
}

//CancelBoletoRequest ...
type CancelBoletoRequest struct {
	AuthenticationCode string   `validate:"required" json:"authenticationCode,omitempty"`
	Account            *Account `validate:"required" json:"account,omitempty"`
}

//PayBoletoRequest ...
type PayBoletoRequest struct {
	AuthenticationCode string   `validate:"required" json:"authenticationCode,omitempty"`
	Account            *Account `validate:"required" json:"account,omitempty"`
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

// TransferRequest ...
type TransferRequest struct {
	Amount      float64    `validate:"required" json:"amount,omitempty"`
	Description string     `validate:"required" json:"description,omitempty"`
	Sender      *Sender    `validate:"required,dive" json:"sender,omitempty"`
	Recipient   *Recipient `validate:"required,dive" json:"recipient,omitempty"`
}

// Sender ...
type Sender struct {
	Branch   string `validate:"required" json:"branch,omitempty"`
	Account  string `validate:"required" json:"account,omitempty"`
	Document string `validate:"required" json:"document,omitempty"`
	Name     string `validate:"required" json:"name,omitempty"`
}

// SenderResponse ...
type SenderResponse struct {
	Document string           `json:"document,omitempty"`
	Name     string           `json:"name,omitempty"`
	Account  *AccountResponse `json:"account,omitempty"`
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

// RecipientResponse ...
type RecipientResponse struct {
	Document string           `json:"document,omitempty"`
	Name     string           `json:"name,omitempty"`
	Account  *AccountResponse `json:"account,omitempty"`
}

type DocumentType string

const (
	// DocumentTypeRG ...
	DocumentTypeRG DocumentType = "RG"
	// DocumentTypeCNH ...
	DocumentTypeCNH DocumentType = "CNH"
	// DocumentTypeSELFIE ...
	DocumentTypeSELFIE DocumentType = "SELFIE"
)

type DocumentSide string

const (
	// DocumentSideFront ...
	DocumentSideFront DocumentSide = "FRONT"
	// DocumentSideBack ...
	DocumentSideBack DocumentSide = "BACK"
)

// DocumentAnalysisRequest ...
type DocumentAnalysisRequest struct {
	DocumentType		DocumentType		`json:"documentType,omitempty"`
	DocumentSide		DocumentSide		`json:"documentSide,omitempty"`
	Image				string				`json:"image,omitempty"`
}

// DocumentAnalysisRequestedResponse ...
type DocumentAnalysisRequestedResponse struct {
	Token				string		`json:"token,omitempty"`
}

// DocumentAnalysisResponse ...
type DocumentAnalysisResponse struct {
	DocumentNumber		string				`json:"document_number,omitempty"`
	Token				string				`json:"token,omitempty"`
	Status				string				`json:"status,omitempty"`
	DocumentType		string				`json:"documentType,omitempty"`
	DocumentSide		string				`json:"documentSide,omitempty"`
	FaceMatch			*FaceMatch  		`json:"faceMatch,omitempty"`
	FaceDetails			*FaceDetails		`json:"faceDetails,omitempty"`
	DocumentDetails		*DocumentDetails	`json:"documentDetails,omitempty"`
	Liveness			*Liveness			`json:"liveness,omitempty"`
	AnalyzedAt			string				`json:"analyzedAt,omitempty"`
}

type FaceMatch struct {
	Status				string 				`json:"status,omitempty"`
	Similarity			int 				`json:"similarity,omitempty"`
	Confidence			int 				`json:"confidence,omitempty"`
}

type FaceDetails struct {
	Status				string 			  	`json:"status,omitempty"`
	Confidence			float32 			`json:"confidence,omitempty"`
	AgeRange			*AgeRange   	  	`json:"ageRange,omitempty"`
	Gender				*Gender     	  	`json:"gender,omitempty"`
	Sunglasses			*Sunglasses       	`json:"sunglasses,omitempty"`
	EyesOpen			*EyesOpen     	  	`json:"eyesOpen,omitempty"`
	Emotions			[]*Emotions       	`json:"emotions,omitempty"`
}

type AgeRange struct {
	Low		int		`json:"low,omitempty"`
	High	int		`json:"high,omitempty"`
}

type Gender struct {
	Value			string 		`json:"value,omitempty"`
	Confidence		float32 	`json:"confidence,omitempty"`
}

type Sunglasses struct {
	Value			bool 		`json:"value,omitempty"`
	Confidence		float32 	`json:"confidence,omitempty"`
}

type EyesOpen struct {
	Value			bool 		`json:"value,omitempty"`
	Confidence		float32 	`json:"confidence,omitempty"`
}

type Emotions struct {
	Value			string 		`json:"value,omitempty"`
	Confidence		float32 	`json:"confidence,omitempty"`
}

type DocumentDetails struct {
	Status								string 		`json:"status,omitempty"`
	IdentifiedDocumentType				string 		`json:"identifiedDocumentType,omitempty"`
	IdNumber							string 		`json:"idNumber,omitempty"`
	CpfNumber							string 		`json:"cpfNumber,omitempty"`
	BirthDate							string 		`json:"birthDate,omitempty"`
	FatherName							string 		`json:"fatherName,omitempty"`
	MotherName							string 		`json:"motherName,omitempty"`
	RegisterName						string 		`json:"registerName,omitempty"`
	ValidDate							string 		`json:"validDate,omitempty"`
	DriveLicenseCategory				string 		`json:"driveLicenseCategory,omitempty"`
	DriveLicenseNumber					string 		`json:"driveLicenseNumber,omitempty"`
	DriveLicenseFirstQualifyingDate		string 		`json:"driveLicenseFirstQualifyingDate,omitempty"`
	FederativeUnit						string 		`json:"federativeUnit,omitempty"`
	IssuedBy							string 		`json:"issuedBy,omitempty"`
	IssuePlace							string 		`json:"issuePlace,omitempty"`
	IssueDate							string 		`json:"issueDate,omitempty"`
}

type Liveness struct {
	Status			string 		`json:"status,omitempty"`
	Confidence		float32 	`json:"confidence,omitempty"`
}