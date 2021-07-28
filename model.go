package bankly

import (
	"os"
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
	// PaymentPath ...
	PaymentPath = "bill-payment"
	// BanksPath ...
	BanksPath = "banklist"
	//BankStatementsPath ...
	BankStatementsPath = "events"
	// DocumentAnalysisPath ...
	DocumentAnalysisPath = "/document-analysis"
)

const (
	// InternalBankCode ...
	InternalBankCode string = "332"
)

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

type DocumentAnalysisStatus string

const (
	// DocumentAnalysisStatusAnalyzing...
	DocumentAnalysisStatusAnalyzing DocumentAnalysisStatus = "ANALYZING"
	// DocumentAnalysisStatusAnalysisCompleted ...
	DocumentAnalysisStatusAnalysisCompleted DocumentAnalysisStatus = "ANALYSIS_COMPLETED"
	// DocumentAnalysisStatusUnexpectedError ...
	DocumentAnalysisStatusUnexpectedError DocumentAnalysisStatus = "UNEXPECTED_ERROR"
	// DocumentAnalysisStatusForbiddenWord ...
	DocumentAnalysisStatusForbiddenWord DocumentAnalysisStatus = "FORBIDDEN_WORD"
	// DocumentAnalysisStatusDataRecused ...
	DocumentAnalysisStatusDataRecused DocumentAnalysisStatus = "DATA_RECUSED"
	// DocumentAnalysisStatusPhotoRecused ...
	DocumentAnalysisStatusPhotoRecused DocumentAnalysisStatus = "PHOTO_RECUSED"
)

type DetailsStatus string

const (
	// DetailsStatusLivenessFound ...
	DetailsStatusLivenessFound DetailsStatus = "LIVENESS_FOUND"
	// DetailsStatusNoLiveness ...
	DetailsStatusNoLiveness DetailsStatus = "NO_LIVENESS"
	// DetailsStatusCouldNotDetectFace ...
	DetailsStatusCouldNotDetectFace DetailsStatus = "COULD_NOT_DETECT_FACE"
	// DetailsStatusDetectFace ...
	DetailsStatusDetectFace DetailsStatus = "DETECTED_FACE"
	// DetailsStatusManyFacesDetected ...
	DetailsStatusManyFacesDetected DetailsStatus = "MANY_FACES_DETECTED"
	// DetailsStatusHasFaceMatch ...
	DetailsStatusHasFaceMatch DetailsStatus = "HAS_FACE_MATCH"
	// DetailsStatusUnMatchedDocument ...
	DetailsStatusUnMatchedDocument DetailsStatus = "UNMATCHED_DOCUMENT"
	// DetailsStatusNoDocumentFound ...
	DetailsStatusNoDocumentFound DetailsStatus = "NO_DOCUMENT_FOUND"
	// DetailsStatusNoInfoFound ...
	DetailsStatusNoInfoFound DetailsStatus = "NO_INFO_FOUND"
	// DetailsStatusDetectedDocument ...
	DetailsStatusDetectedDocument DetailsStatus = "DETECTED_DOCUMENT"
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

// CustomersResponse ...
type CustomersResponse struct {
	DocumentNumber             string         `json:"documentNumber"`
	RegisterName               string         `json:"registerName"`
	SocialName                 string         `json:"socialName"`
	Email                      string         `json:"email"`
	Phone                      Phone          `json:"phone"`
	Address                    Address        `json:"address"`
	MotherName                 string         `json:"motherName"`
	BirthDate                  string         `json:"birthDate"`
	IsPoliticallyExposedPerson bool           `json:"isPoliticallyExposedPerson"`
	Reasons                    []string       `json:"reasons"`
	Status                     CustomerStatus `json:"status"`
	Profile                    string         `json:"profile"`
	CreatedAt                  string         `json:"createdAt"`
	UpdatedAt                  string         `json:"updatedAt"`
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

type BoletoAddress struct {
	AddressLine string `validate:"required" json:"addressLine,omitempty"`
	ZipCode     string `validate:"required" json:"zipCode,omitempty"`
	State       string `validate:"required" json:"state,omitempty"`
	City        string `validate:"required" json:"city,omitempty"`
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

// ErrorResponse ...
type ErrorResponse struct {
	Errors    []ErrorModel `json:"errors,omitempty"`
	Reference string       `json:"reference,omitempty"`
}

//BoletoErrorResponse ...
type BoletoErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

//PaymentErrorResponse ...
type PaymentErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
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
	Bank    *BankData       `json:"bank,omitempty"`
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
	Status             TransfersStatus    `json:"status"`
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
	ResultLevelBasic      ResultLevel = "BASIC"
	ResultLevelDetailed   ResultLevel = "DETAILED"
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
	Document    string     `validate:"required,cnpjcpf" json:"documentNumber,omitempty"`
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

//SimulatePaymentRequest ...
type SimulatePaymentRequest struct {
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

// BankData ...
type BankData struct {
	ISPB string `json:"ispb,omitempty"`
	Name string `json:"name,omitempty"`
	Code string `json:"compe,omitempty"`
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

// DocumentAnalysisRequest ...
type DocumentAnalysisRequest struct {
	Document     string       `validate:"required" json:"document,omitempty"`
	DocumentType DocumentType `validate:"required" json:"document_type,omitempty"`
	DocumentSide DocumentSide `validate:"required" json:"document_side,omitempty"`
	ImageFile    os.File      `validate:"required" json:"image_file,omitempty"`
}

// DocumentAnalysisRequestedResponse ...
type DocumentAnalysisRequestedResponse struct {
	Token string `json:"token,omitempty"`
}

// DocumentAnalysisResponse ...
type DocumentAnalysisResponse struct {
	DocumentNumber  string                 `json:"document_number,omitempty"`
	Token           string                 `json:"token,omitempty"`
	Status          DocumentAnalysisStatus `json:"status,omitempty"`
	DocumentType    string                 `json:"document_type,omitempty"`
	DocumentSide    string                 `json:"document_side,omitempty"`
	FaceMatch       *FaceMatch             `json:"face_match,omitempty"`
	FaceDetails     *FaceDetails           `json:"face_details,omitempty"`
	DocumentDetails *DocumentDetails       `json:"document_details,omitempty"`
	Liveness        *Liveness              `json:"liveness,omitempty"`
	AnalyzedAt      string                 `json:"analyzed_at,omitempty"`
}

// BanklyDocumentAnalysisResponse ...
type BanklyDocumentAnalysisResponse struct {
	Token           string                 `json:"token,omitempty"`
	Status          DocumentAnalysisStatus `json:"status,omitempty"`
	DocumentType    string                 `json:"documentType,omitempty"`
	DocumentSide    string                 `json:"documentSide,omitempty"`
	FaceMatch       *FaceMatch             `json:"faceMatch,omitempty"`
	FaceDetails     *FaceDetails           `json:"faceDetails,omitempty"`
	DocumentDetails *DocumentDetails       `json:"documentDetails,omitempty"`
	Liveness        *Liveness              `json:"liveness,omitempty"`
	AnalyzedAt      string                 `json:"analyzedAt,omitempty"`
}

type FaceMatch struct {
	Status     DetailsStatus `json:"status,omitempty"`
	Similarity float32       `json:"similarity,omitempty"`
	Confidence float32       `json:"confidence,omitempty"`
}

type FaceDetails struct {
	Status     DetailsStatus `json:"status,omitempty"`
	Confidence float32       `json:"confidence,omitempty"`
	AgeRange   *AgeRange     `json:"ageRange,omitempty"`
	Gender     *Gender       `json:"gender,omitempty"`
	Sunglasses *Sunglasses   `json:"sunglasses,omitempty"`
	EyesOpen   *EyesOpen     `json:"eyesOpen,omitempty"`
	Emotions   []*Emotions   `json:"emotions,omitempty"`
}

type AgeRange struct {
	Low  int `json:"low,omitempty"`
	High int `json:"high,omitempty"`
}

type Gender struct {
	Value      string  `json:"value,omitempty"`
	Confidence float32 `json:"confidence,omitempty"`
}

type Sunglasses struct {
	Value      bool    `json:"value,omitempty"`
	Confidence float32 `json:"confidence,omitempty"`
}

type EyesOpen struct {
	Value      bool    `json:"value,omitempty"`
	Confidence float32 `json:"confidence,omitempty"`
}

type Emotions struct {
	Value      string  `json:"value,omitempty"`
	Confidence float32 `json:"confidence,omitempty"`
}

type DocumentDetails struct {
	Status                          DetailsStatus `json:"status,omitempty"`
	IdentifiedDocumentType          string        `json:"identifiedDocumentType,omitempty"`
	IdNumber                        string        `json:"idNumber,omitempty"`
	CpfNumber                       string        `json:"cpfNumber,omitempty"`
	BirthDate                       string        `json:"birthDate,omitempty"`
	FatherName                      string        `json:"fatherName,omitempty"`
	MotherName                      string        `json:"motherName,omitempty"`
	RegisterName                    string        `json:"registerName,omitempty"`
	ValidDate                       string        `json:"validDate,omitempty"`
	DriveLicenseCategory            string        `json:"driveLicenseCategory,omitempty"`
	DriveLicenseNumber              string        `json:"driveLicenseNumber,omitempty"`
	DriveLicenseFirstQualifyingDate string        `json:"driveLicenseFirstQualifyingDate,omitempty"`
	FederativeUnit                  string        `json:"federativeUnit,omitempty"`
	IssuedBy                        string        `json:"issuedBy,omitempty"`
	IssuePlace                      string        `json:"issuePlace,omitempty"`
	IssueDate                       string        `json:"issueDate,omitempty"`
}

type Liveness struct {
	Status     DetailsStatus `json:"status,omitempty"`
	Confidence float32       `json:"confidence,omitempty"`
}

// ValidatePaymentRequest ...
type ValidatePaymentRequest struct {
	Code string `validate:"required" json:"code,omitempty"`
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

// ValidatePaymentResponse ...
type ValidatePaymentResponse struct {
	ID                string         `json:"id,omitempty"`
	Assignor          string         `json:"assignor,omitempty"`
	Code              string         `json:"code,omitempty"`
	Digitable         string         `json:"digitable,omitempty"`
	Amount            float64        `json:"amount,omitempty"`
	OriginalAmount    float64        `json:"originalAmount,omitempty"`
	MaxAmount         float64        `json:"maxAmount,omitempty"`
	NextSettle        bool           `json:"nextSettle,omitempty"`
	AllowChangeAmount bool           `json:"allowChangeAmount,omitempty"`
	DueDate           string         `json:"dueDate,omitempty"`
	SettleDate        string         `json:"settleDate,omitempty"`
	Payer             *PaymentPayer  `json:"payer,omitempty"`
	Recipient         *PaymentPayer  `json:"recipient,omitempty"`
	BusinessHours     *BusinessHours `json:"businessHours,omitempty"`
	Charges           *Charges       `json:"charges,omitempty"`
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
	AuthenticationCode string    `json:"authenticationCode,omitempty"`
	Status             string    `json:"status,omitempty"`
	Digitable          string    `json:"digitable,omitempty"`
	Description        *string   `json:"description,omitempty"`
	BankBranch         string    `json:"bankBranch,omitempty"`
	BankAccount        string    `json:"bankAccount,omitempty"`
	RecipientDocument  string    `json:"recipientDocument,omitempty"`
	RecipientName      string    `json:"recipientName,omitempty"`
	Amount             float64   `json:"amount,omitempty"`
	OriginalAmount     float64   `json:"originalAmount,omitempty"`
	Asignor            float64   `json:"asignor,omitempty"`
	Charges            *Charges  `json:"charges,omitempty"`
	SettleDate         time.Time `json:"settleDate,omitempty"`
	PaymentDate        time.Time `json:"paymentDate,omitempty"`
	ConfirmedAt        time.Time `json:"confirmedAt,omitempty"`
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

// ParseDocumentAnalysisResponse ....
func ParseDocumentAnalysisResponse(documentNumber string, banklyResponse *BanklyDocumentAnalysisResponse) *DocumentAnalysisResponse {
	if banklyResponse == nil {
		return nil
	}
	return &DocumentAnalysisResponse{
		DocumentNumber:  documentNumber,
		Token:           banklyResponse.Token,
		Status:          banklyResponse.Status,
		DocumentType:    banklyResponse.DocumentType,
		DocumentSide:    banklyResponse.DocumentSide,
		FaceMatch:       banklyResponse.FaceMatch,
		FaceDetails:     banklyResponse.FaceDetails,
		DocumentDetails: banklyResponse.DocumentDetails,
		Liveness:        banklyResponse.Liveness,
		AnalyzedAt:      banklyResponse.AnalyzedAt,
	}
}
