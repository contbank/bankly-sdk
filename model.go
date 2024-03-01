package bankly

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/contbank/grok"
)

const (
	// LoginPath ...
	LoginPath = "connect/token"
	// LoginMtlsPath ...
	LoginMtlsPath = "oauth2/token"
	// ClientPath ...
	ClientPath = "oauth2/register"
	// CustomersPath ..
	CustomersPath = "customers"
	// AccountsPath ...
	AccountsPath = "accounts"
	// IncomeReportPath ...
	IncomeReportPath = "income-report"
	// TransfersPath ...
	TransfersPath = "fund-transfers"
	// BusinessPath ...
	BusinessPath = "business"
	// CorporationBusinessPath ...
	CorporationBusinessPath = "corporation-business"
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
	// CardDocumentPath ...
	CardDocumentPath = "cards/document"
	// BoletosSettledPath ...
	BoletosSettledPath = "/bankslip/settlementpayment"
)

type PixType string

const (
	// PixCNPJ ...
	PixCNPJ PixType = "CNPJ"
	// PixCPF ...
	PixCPF PixType = "CPF"
	// PixEMAIL ...
	PixEMAIL PixType = "EMAIL"
	// PixPHONE ...
	PixPHONE PixType = "PHONE"
	//  PixEVP ...
	PixEVP PixType = "EVP"
)

type PixClaimType string

const (
	Portability PixClaimType = "PORTABILITY"
	Ownership   PixClaimType = "OWNERSHIP"
)

type StatusClaim string

const (
	Open              StatusClaim = "OPEN"
	WaitingResolution StatusClaim = "WAITING_RESOLUTION"
	Confirmed         StatusClaim = "CONFIRMED"
	CanceledClaim     StatusClaim = "CANCELED"
	CompletedClaim    StatusClaim = "COMPLETED"
)

type CancelReason string

const (
	ClaimerRequest   CancelReason = "CLAIMER_REQUEST"
	DonorRequest     CancelReason = "DONOR_REQUEST"
	AccountClosure   CancelReason = "ACCOUNT_CLOSURE"
	Fraud            CancelReason = "FRAUD"
	DefaultOperation CancelReason = "DEFAULT_OPERATION"
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
	// DocumentTypeARTICLES_OF_INCORPORATION ...
	DocumentTypeARTICLES_OF_INCORPORATION DocumentType = "ARTICLES_OF_INCORPORATION"
	// DocumentTypeLAST_CONTRACT_CHANGE ...
	DocumentTypeLAST_CONTRACT_CHANGE DocumentType = "LAST_CONTRACT_CHANGE"
	// DocumentTypeBALANCE_SHEET ...
	DocumentTypeBALANCE_SHEET DocumentType = "BALANCE_SHEET"
	// DocumentTypePOWER_OF_ATTORNEY ...
	DocumentTypePOWER_OF_ATTORNEY DocumentType = "POWER_OF_ATTORNEY"
)

type DocumentSide string

const (
	// DocumentSideFront ...
	DocumentSideFront DocumentSide = "FRONT"
	// DocumentSideBack ...
	DocumentSideBack DocumentSide = "BACK"
)

type Provider string

const (
	// DocumentProviderUnicoCheck ...
	DocumentProviderUnicoCheck Provider = "UNICO_CHECK"
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

// DeclaredIncome Faixa de renda declarada (PF)
type DeclaredIncome string

const (
	// LESS_THAN_ONE_THOUSAND Inferior a mil.
	LESS_THAN_ONE_THOUSAND DeclaredIncome = "LESS_THAN_ONE_THOUSAND"
	// FROM_ONE_THOUSAND_TO_TWO_THOUSAND De mil a dois mil.
	FROM_ONE_THOUSAND_TO_TWO_THOUSAND DeclaredIncome = "FROM_ONE_THOUSAND_TO_TWO_THOUSAND"
	// FROM_TWO_THOUSAND_TO_THREE_THOUSAND De 2 mil a 3 mil.
	FROM_TWO_THOUSAND_TO_THREE_THOUSAND DeclaredIncome = "FROM_TWO_THOUSAND_TO_THREE_THOUSAND"
	// FROM_THREE_THOUSAND_TO_FIVE_THOUSAND De 3 mil a 5 mil.
	FROM_THREE_THOUSAND_TO_FIVE_THOUSAND DeclaredIncome = "FROM_THREE_THOUSAND_TO_FIVE_THOUSAND"
	// FROM_FIVE_THOUSAND_TO_TEN_THOUSAND De 5 mil a 10 mil.
	FROM_FIVE_THOUSAND_TO_TEN_THOUSAND DeclaredIncome = "FROM_FIVE_THOUSAND_TO_TEN_THOUSAND"
	// FROM_TEN_THOUSAND_TO_TWENTY_THOUSAND De 10 mil a 20 mil.
	FROM_TEN_THOUSAND_TO_TWENTY_THOUSAND DeclaredIncome = "FROM_TEN_THOUSAND_TO_TWENTY_THOUSAND"
	// OVER_TWENTY_THOUSAND Acima de 20 mil.
	OVER_TWENTY_THOUSAND DeclaredIncome = "OVER_TWENTY_THOUSAND"
)

type PoliticallyExposedPersonLevel string

const (
	// NONE nenhuma exposição política
	NONE PoliticallyExposedPersonLevel = "NONE"
	// SELF a própria pessoa é exposta politicamente
	SELF PoliticallyExposedPersonLevel = "SELF"
	// RELATED algum parente ou pessoa relacionada é exposta politicamente
	RELATED PoliticallyExposedPersonLevel = "RELATED"
)

// DeclaredAnnualBilling Faixa de faturamento anual da empresa (PJ)
type DeclaredAnnualBilling string

const (
	// UP_TO_FIFTY_THOUSAND Até 50 mil
	UP_TO_FIFTY_THOUSAND DeclaredAnnualBilling = "UP_TO_FIFTY_THOUSAND"
	// MORE_THAN_FIFTY_THOUSAND_UP_TO_ONE_HUNDRED_THOUSAND De 50 mil a 100 mil.
	MORE_THAN_FIFTY_THOUSAND_UP_TO_ONE_HUNDRED_THOUSAND DeclaredAnnualBilling = "MORE_THAN_FIFTY_THOUSAND_UP_TO_ONE_HUNDRED_THOUSAND"
	// MORE_THAN_ONE_HUNDRED_THOUSAND_UP_TO_TWO_HUNDRED_AND_FIFTY_THOUSAND De 100 mil a 250 mil.
	MORE_THAN_ONE_HUNDRED_THOUSAND_UP_TO_TWO_HUNDRED_AND_FIFTY_THOUSAND DeclaredAnnualBilling = "MORE_THAN_ONE_HUNDRED_THOUSAND_UP_TO_TWO_HUNDRED_AND_FIFTY_THOUSAND"
	// MORE_THAN_TWO_HUNDRED_AND_FIFTY_THOUSAND_UP_TO_FIVE_HUNDRED_THOUSAND De 250 mil a 500 mil.
	MORE_THAN_TWO_HUNDRED_AND_FIFTY_THOUSAND_UP_TO_FIVE_HUNDRED_THOUSAND DeclaredAnnualBilling = "MORE_THAN_TWO_HUNDRED_AND_FIFTY_THOUSAND_UP_TO_FIVE_HUNDRED_THOUSAND"
	// MORE_THAN_FIVE_HUNDRED_THOUSAND_UP_TO_ONE_MILLION De 500 mil a 1 milhão.
	MORE_THAN_FIVE_HUNDRED_THOUSAND_UP_TO_ONE_MILLION DeclaredAnnualBilling = "MORE_THAN_FIVE_HUNDRED_THOUSAND_UP_TO_ONE_MILLION"
	// MORE_THAN_ONE_MILLION_UP_TO_TWO_MILLION_AND_FIVE_HUNDRED_THOUSAND De 1 milhão a 2 milhões e 500 mil.
	MORE_THAN_ONE_MILLION_UP_TO_TWO_MILLION_AND_FIVE_HUNDRED_THOUSAND DeclaredAnnualBilling = "MORE_THAN_ONE_MILLION_UP_TO_TWO_MILLION_AND_FIVE_HUNDRED_THOUSAND"
	// MORE_THAN_TWO_MILLION_AND_FIVE_HUNDRED_THOUSAND_UP_TO_FIVE_MILLION De 2 milhões e 500 mil a 5 milhões.
	MORE_THAN_TWO_MILLION_AND_FIVE_HUNDRED_THOUSAND_UP_TO_FIVE_MILLION DeclaredAnnualBilling = "MORE_THAN_TWO_MILLION_AND_FIVE_HUNDRED_THOUSAND_UP_TO_FIVE_MILLION"
	// MORE_THAN_FIVE_MILLION_UP_TO_TEN_MILLION De 5 milhões a 10 milhões.
	MORE_THAN_FIVE_MILLION_UP_TO_TEN_MILLION DeclaredAnnualBilling = "MORE_THAN_FIVE_MILLION_UP_TO_TEN_MILLION"
	// MORE_THAN_TEN_MILLION_UP_TO_TWENTY_FIVE_MILLION De 10 milhões a 25 milhões.
	MORE_THAN_TEN_MILLION_UP_TO_TWENTY_FIVE_MILLION DeclaredAnnualBilling = "MORE_THAN_TEN_MILLION_UP_TO_TWENTY_FIVE_MILLION"
	// MORE_THAN_TWENTY_FIVE_MILLION_UP_TO_FIFTY_MILLION De 25 milhões a 50 milhões.
	MORE_THAN_TWENTY_FIVE_MILLION_UP_TO_FIFTY_MILLION DeclaredAnnualBilling = "MORE_THAN_TWENTY_FIVE_MILLION_UP_TO_FIFTY_MILLION"
	// MORE_THAN_FIFTY_MILLION_UP_TO_ONE_HUNDRED_MILLION De 50 milhões a 100 milhões.
	MORE_THAN_FIFTY_MILLION_UP_TO_ONE_HUNDRED_MILLION DeclaredAnnualBilling = "MORE_THAN_FIFTY_MILLION_UP_TO_ONE_HUNDRED_MILLION"
	// MORE_THAN_ONE_HUNDRED_MILLION_UP_TO_TWO_HUNDRED_AND_FIFTY_MILLION De 100 milhões a 250 milhões.
	MORE_THAN_ONE_HUNDRED_MILLION_UP_TO_TWO_HUNDRED_AND_FIFTY_MILLION DeclaredAnnualBilling = "MORE_THAN_ONE_HUNDRED_MILLION_UP_TO_TWO_HUNDRED_AND_FIFTY_MILLION"
	// MORE_THAN_TWO_HUNDRED_AND_FIFTY_MILLION_UP_TO_FIVE_HUNDRED_MILLION De 250 milhões a 500 milhões.
	MORE_THAN_TWO_HUNDRED_AND_FIFTY_MILLION_UP_TO_FIVE_HUNDRED_MILLION DeclaredAnnualBilling = "MORE_THAN_TWO_HUNDRED_AND_FIFTY_MILLION_UP_TO_FIVE_HUNDRED_MILLION"
	// MORE_THAN_FIVE_HUNDRED_MILLION Mais de 500 milhões.
	MORE_THAN_FIVE_HUNDRED_MILLION DeclaredAnnualBilling = "MORE_THAN_FIVE_HUNDRED_MILLION"
	// NOT_DECLARED Não declarada.
	NOT_DECLARED DeclaredAnnualBilling = "NOT_DECLARED"
	// EXEMPT_COMPANY Empresa isenta.
	EXEMPT_COMPANY DeclaredAnnualBilling = "EXEMPT_COMPANY"
	// INACTIVE_COMPANY Empresa inativa.
	INACTIVE_COMPANY DeclaredAnnualBilling = "INACTIVE_COMPANY"
)

type DuplicateCardStatus string

const (
	// LOST_MY_CARD Cartão perdido
	LOST_MY_CARD DuplicateCardStatus = "LostMyCard"
	// CARD_WAS_STOLEN Cartão roubado
	CARD_WAS_STOLEN DuplicateCardStatus = "CardWasStolen"
	// CARD_WAS_DAMAGED Cartão danificado
	CARD_WAS_DAMAGED DuplicateCardStatus = "CardWasDamaged"
	// CARD_NOT_DELIVERED Cartão não entregue
	CARD_NOT_DELIVERED DuplicateCardStatus = "CardNotDelivered"
	// UNRECOGNIZED_ONLINE_PURCHASE Compra online não reconhecida
	UNRECOGNIZED_ONLINE_PURCHASE DuplicateCardStatus = "UnrecognizedOnlinePurchase"
)

// DocumentAnalysisUnicoCheckRequest ...
type DocumentAnalysisUnicoCheckRequest struct {
	Document         string                   `validate:"required" json:"document,omitempty"`
	DocumentType     DocumentType             `validate:"required" json:"document_type,omitempty"`
	DocumentSide     DocumentSide             `validate:"required" json:"document_side,omitempty"`
	ImageFile        os.File                  `validate:"required" json:"image,omitempty"`
	Provider         string                   `validate:"required" json:"provider,omitempty"`
	ProviderMetaData *ProviderMetadataRequest `json:"provider_metadata,omitempty"`
}

// ProviderMetadataRequest ...
type ProviderMetadataRequest struct {
	IsLastDocument bool   `validate:"required" json:"is_last_document,omitempty"`
	Encrypted      string `validate:"required" json:"encrypted,omitempty"`
}

// ParseProviderMetadaRequest ...
func ParseProviderMetadaRequest(r *ProviderMetadataRequest) *ProviderMetadata {
	return &ProviderMetadata{
		IsLastDocument: r.IsLastDocument,
		Encrypted:      r.Encrypted,
	}
}

// DocumentAnalysisUnicoCheck ...
type DocumentAnalysisUnicoCheck struct {
	Document         string            `validate:"required" json:"document,omitempty"`
	DocumentType     DocumentType      `validate:"required" json:"documentType,omitempty"`
	DocumentSide     DocumentSide      `validate:"required" json:"documentSide,omitempty"`
	ImageFile        os.File           `validate:"required" json:"image,omitempty"`
	Provider         string            `validate:"required" json:"provider,omitempty"`
	ProviderMetaData *ProviderMetadata `json:"providerMetadata,omitempty"`
}

// ProviderMetadata ...
type ProviderMetadata struct {
	IsLastDocument bool   `json:"isLastDocument,omitempty"`
	Encrypted      string `json:"encrypted,omitempty"`
}

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
	Document                 string                   `validate:"required,cpf" json:"documentNumber,omitempty"`
	RegisterName             string                   `validate:"required" json:"registerName,omitempty"`
	SocialName               string                   `json:"socialName,omitempty"`
	Email                    string                   `validate:"required" json:"email,omitempty"`
	Phone                    *Phone                   `validate:"required,dive" json:"phone,omitempty"`
	Address                  *Address                 `validate:"required,dive" json:"address,omitempty"`
	BirthDate                time.Time                `validate:"required" json:"birthDate,omitempty"`
	MotherName               string                   `validate:"required" json:"motherName,omitempty"`
	DeclaredIncome           DeclaredIncome           `json:"declaredIncome,omitempty"` // deprecated
	AssertedIncome           AssertedIncome           `validate:"required" json:"assertedIncome,omitempty"`
	Occupation               string                   `validate:"required" json:"occupation,omitempty"`
	PoliticallyExposedPerson PoliticallyExposedPerson `validate:"required,dive" json:"pep,omitempty"`
	Documentation            Documentation            `validate:"required,dive" json:"documentation,omitempty"`
}

// AssertedIncome ...
type AssertedIncome struct {
	Value    float64 `json:"value,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

// Documentation
type Documentation struct {
	DocumentNumber string `json:"documentNumber,omitempty"`
	TokenSelfie    string `validate:"required" json:"selfie,omitempty"`
	TokenCardFront string `validate:"required" json:"idCardFront,omitempty"`
	TokenCardBack  string `validate:"required" json:"idCardBack,omitempty"`
}

// CorporationBusinessDocumentation
type CorporationBusinessDocumentation struct {
	ArticlesOfIncorporation *string       `json:"articlesOfIncorporation,omitempty"`
	LastContractChange      *string       `json:"lastContractChange,omitempty"`
	BalanceSheet            *string       `json:"balanceSheet,omitempty"`
	LegalRepresentative     Documentation `json:"legalRepresentative,omitempty"`
	PowerOfAttorney         *string       `json:"powerOfAttorney,omitempty"`
}

// PoliticallyExposedPerson ...
type PoliticallyExposedPerson struct {
	Level PoliticallyExposedPersonLevel `validate:"required" json:"level,omitempty"`
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
	IsPoliticallyExposedPerson bool           `json:"isPoliticallyExposedPerson"`
	Reasons                    []string       `json:"reasons"`
	Status                     CustomerStatus `json:"status"`
	Profile                    string         `json:"profile"`
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
	IsActive     bool    `json:"isActive,omitempty"`
}

// BoletoAddress ...
type BoletoAddress struct {
	ZipCode      string `validate:"required" json:"zipCode,omitempty"`
	AddressLine  string `validate:"required" json:"addressLine,omitempty"`
	Neighborhood string `validate:"required" json:"neighborhood,omitempty"`
	State        string `validate:"required" json:"state,omitempty"`
	City         string `validate:"required" json:"city,omitempty"`
}

// Account ...
type Account struct {
	Number      string               `validate:"required" json:"number,omitempty"`
	Branch      string               `validate:"required" json:"branch,omitempty"`
	AccountType TransfersAccountType `json:"type,omitempty"`
}

// Payer ...
type Payer struct {
	Name      string         `validate:"required" json:"name,omitempty"`
	TradeName string         `json:"tradeName,omitempty"`
	Document  string         `validate:"required,cnpjcpf" json:"document,omitempty"`
	Address   *BoletoAddress `validate:"required" json:"address,omitempty"`
}

// BoletoPayer ...
type BoletoPayer struct {
	Name      string         `validate:"required,max=50" json:"name,omitempty"`
	TradeName string         `validate:"max=80" json:"tradeName,omitempty"`
	Document  string         `validate:"required,cnpjcpf" json:"document,omitempty"`
	Address   *BoletoAddress `validate:"required" json:"address,omitempty"`
}

// ErrorResponse ...
type ErrorResponse struct {
	Errors    []ErrorModel `json:"errors,omitempty"`
	Title     string       `json:"title,omitempty"`
	Status    int32        `json:"status,omitempty"`
	TraceId   string       `json:"traceId,omitempty"`
	Reference string       `json:"reference,omitempty"`
	CodeMessageErrorResponse
}

// CardErrorResponse ...
type CardErrorResponse struct {
	ErrorKey string `json:"errorKey,omitempty"`
	CodeMessageErrorResponse
}

// CodeMessageErrorResponse ...
type CodeMessageErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// TransferErrorResponse ...
type TransferErrorResponse struct {
	Layer           string               `json:"layer,omitempty"`
	ApplicationName string               `json:"applicationName,omitempty"`
	Errors          []KeyValueErrorModel `json:"errors,omitempty"`
	CodeMessageErrorResponse
}

// ErrorModel ...
type ErrorModel struct {
	Code         string   `json:"code,omitempty"`
	PropertyName string   `json:"propertyName,omitempty"`
	Messages     []string `json:"messages,omitempty"`
	KeyValueErrorModel
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

// IncomeReportRequest ...
type IncomeReportRequest struct {
	Account string `validate:"required" form:"account"`
	Year    string `validate:"required" form:"year"`
}

// IncomeReportResponse ...
type IncomeReportResponse struct {
	Data  IncomeReportData   `json:"data"`
	Links []IncomeReportLink `json:"links"`
}

// IncomeReportData ...
type IncomeReportData struct {
	Holder  *IncomeReportHolder  `json:"holder"`
	Account *IncomeReportAccount `json:"account"`
	Payers  []IncomeReportPayers `json:"payers"`
}

// IncomeReportHolder ...
type IncomeReportHolder struct {
	Type      string                `json:"type"`
	CreatedAt time.Time             `json:"createdAt"`
	Name      string                `json:"name"`
	Document  *IncomeReportDocument `json:"document"`
}

// IncomeReportDocument ...
type IncomeReportDocument struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

// IncomeReportSource ...
type IncomeReportSource struct {
	Name     string                `json:"name"`
	Document *IncomeReportDocument `json:"document"`
}

// IncomeReportPayers ...
type IncomeReportPayers struct {
	Source     *IncomeReportSource   `json:"source"`
	Currencies *IncomeReportCurrency `json:"currencies"`
}

// IncomeReportCurrency ...
type IncomeReportCurrency struct {
	BRL *IncomeReportBRL `json:"brl"`
}

// IncomeReportBRL ...
type IncomeReportBRL struct {
	Balances []IncomeReportBalance `json:"balances"`
}

// IncomeReportBalance ...
type IncomeReportBalance struct {
	Year   int                 `json:"year"`
	Amount *IncomeReportAmount `json:"amount"`
}

// IncomeReportAmount ...
type IncomeReportAmount struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

// IncomeReportAccount ...
type IncomeReportAccount struct {
	Branch string `json:"branch"`
	Number string `json:"number"`
}

// IncomeReportLink ...
type IncomeReportLink struct {
	URL    string `json:"url"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
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
	DocumentNumber        string                 `json:"documentNumber,omitempty"`
	BusinessName          string                 `json:"businessName,omitempty"`
	TradingName           string                 `json:"tradingName,omitempty"`
	BusinessEmail         string                 `json:"businessEmail,omitempty"`
	BusinessType          BusinessType           `json:"businessType"`
	BusinessSize          BusinessSize           `json:"businessSize"`
	BusinessAddress       *Address               `json:"businessAddress,omitempty"`
	DeclaredAnnualBilling DeclaredAnnualBilling  `json:"declaredAnnualBilling"`
	LegalRepresentatives  []*LegalRepresentative `json:"legalRepresentatives,omitempty"`
}

// SimpleBusinessRequest ...
type SimpleBusinessRequest struct {
	DocumentNumber        string                     `validate:"required,cnpj" json:"documentNumber,omitempty"`
	BusinessName          string                     `validate:"required" json:"businessName,omitempty"`
	TradingName           string                     `json:"tradingName,omitempty"`
	BusinessEmail         string                     `validate:"required" json:"businessEmail,omitempty"`
	BusinessType          BusinessType               `validate:"required" json:"businessType"`
	BusinessSize          BusinessSize               `validate:"required" json:"businessSize"`
	BusinessAddress       *Address                   `validate:"required,dive" json:"businessAddress,omitempty"`
	DeclaredAnnualBilling DeclaredAnnualBilling      `validate:"required" json:"declaredAnnualBilling"`
	LegalRepresentative   *SimpleLegalRepresentative `validate:"required,dive" json:"legalRepresentative,omitempty"`
}

// CorporationBusinessRequest ...
type CorporationBusinessRequest struct {
	DocumentNumber        string                            `validate:"required,cnpj" json:"documentNumber,omitempty"`
	BusinessName          string                            `validate:"required" json:"businessName,omitempty"`
	TradingName           string                            `json:"tradingName,omitempty"`
	BusinessEmail         string                            `validate:"required" json:"businessEmail,omitempty"`
	Phone                 *Phone                            `validate:"required,dive" json:"phone,omitempty"`
	BusinessAddress       *Address                          `validate:"required,dive" json:"businessAddress,omitempty"`
	BusinessType          BusinessType                      `validate:"required" json:"businessType"`
	BusinessSize          string                            `validate:"required" json:"businessSize"`
	CNAECode              *string                           `validate:"required" json:"cnaeCode,omitempty"`
	LegalNature           *string                           `validate:"required" json:"legalNature,omitempty"`
	OpeningDate           *time.Time                        `validate:"required" json:"openingDate,omitempty"`
	AnnualBilling         *float64                          `json:"annualBilling,omitempty"` //deprecated
	DeclaredAnnualBilling DeclaredAnnualBilling             `validate:"required" json:"declaredAnnualBilling,omitempty"`
	LegalRepresentatives  []*CorporationLegalRepresentative `validate:"required,dive" json:"legalRepresentatives,omitempty"`
	Owners                []*Owners                         `json:"owners,omitempty"`
	Documentation         CorporationBusinessDocumentation  `validate:"required,dive" json:"documentation,omitempty"`
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
	RegisterName             string                   `validate:"required" json:"registerName,omitempty"`
	SocialName               string                   `json:"socialName,omitempty"`
	Phone                    *Phone                   `validate:"required,dive" json:"phone,omitempty"`
	Address                  *Address                 `validate:"required,dive" json:"address,omitempty"`
	BirthDate                time.Time                `validate:"required" json:"birthDate,omitempty"`
	MotherName               string                   `validate:"required" json:"motherName,omitempty"`
	Email                    string                   `validate:"required" json:"email,omitempty"`
	PoliticallyExposedPerson PoliticallyExposedPerson `json:"pep,omitempty"`
	Occupation               string                   `json:"occupation,omitempty"`
	AssertedIncome           AssertedIncome           `json:"assertedIncome,omitempty"`
}

// CancelAccountReason ...
type CancelAccountReason string

const (
	// HolderRequest ...
	CancelAccountHolderRequest CancelAccountReason = "HOLDER_REQUEST"
	// CommercialDisagreement ...
	CancelAccountCommercialDisagreement CancelAccountReason = "COMMERCIAL_DISAGREEMENT"
)

// CancelAccountRequest ...
type CancelAccountRequest struct {
	Reason CancelAccountReason `validate:"required" json:"reason"`
}

// BusinessType ...
type BusinessType string

const (
	BusinessTypeMEI    BusinessType = "MEI"
	BusinessTypeEI     BusinessType = "EI"
	BusinessTypeEIRELI BusinessType = "EIRELI"
	BusinessTypeLTDA   BusinessType = "LTDA"
	BusinessTypeSA     BusinessType = "SA"
)

// BusinessSize ...
type BusinessSize string

const (
	BusinessSizeMEI BusinessSize = "MEI"
	BusinessSizeME  BusinessSize = "ME"
	BusinessSizeEPP BusinessSize = "EPP"
)

// CorporationBusinessSize ...
type CorporationBusinessSize string

const (
	CorporationBusinessSizeSMALL  CorporationBusinessSize = "SMALL"
	CorporationBusinessSizeMEDIUM CorporationBusinessSize = "MEDIUM"
	CorporationBusinessSizeLARGE  CorporationBusinessSize = "LARGE"
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
	DocumentNumber             string                      `validate:"required,cnpjcpf" json:"documentNumber,omitempty"`
	Document                   LegalRepresentativeDocument `json:"document,omitempty"`
	RegisterName               string                      `json:"registerName,omitempty"`
	SocialName                 string                      `json:"socialName,omitempty"`
	BusinessName               string                      `json:"businessName,omitempty"`
	TradingName                string                      `json:"tradingName,omitempty"`
	BusinessType               BusinessType                `json:"businessType,omitempty"`
	Email                      string                      `validate:"required" json:"email,omitempty"`
	Phone                      *Phone                      `validate:"required,dive" json:"phone,omitempty"`
	Address                    *Address                    `validate:"required,dive" json:"address,omitempty"`
	BirthDate                  time.Time                   `json:"birthDate,omitempty"`
	MotherName                 string                      `json:"motherName,omitempty"`
	IncomeDeclared             *float64                    `json:"incomeDeclared,omitempty"` //deprecated
	DeclaredIncome             DeclaredIncome              `json:"declaredIncome,omitempty"`
	DeclaredAnnualBilling      DeclaredAnnualBilling       `json:"declaredAnnualBilling,omitempty"`
	Occupation                 string                      `json:"occupation,omitempty"`
	ParticipationPercentage    *int                        `json:"participationPercentage,omitempty"`
	IsPoliticallyExposedPerson *bool                       `json:"isPoliticallyExposedPerson,omitempty"` //deprecated
	PoliticallyExposedPerson   PoliticallyExposedPerson    `json:"pep,omitempty"`
	MemberQualification        *string                     `json:"memberQualification,omitempty"`
	PowerOfAttorney            *PowerOfAttorney            `json:"powerOfAttorney,omitempty"`
	Documentation              *Documentation              `json:"documentation,omitempty"`
}

// CorporationLegalRepresentative ...
type CorporationLegalRepresentative struct {
	DocumentNumber             string                   `validate:"required,cnpjcpf" json:"documentNumber,omitempty"` //// corporation nao
	RegisterName               string                   `validate:"required" json:"registerName,omitempty"`
	SocialName                 string                   `json:"socialName,omitempty"`
	BirthDate                  time.Time                `validate:"required" json:"birthDate,omitempty"`
	MotherName                 string                   `validate:"required" json:"motherName,omitempty"`
	Email                      string                   `validate:"required" json:"email,omitempty"`
	DeclaredIncome             DeclaredIncome           `validate:"required" json:"declaredIncome,omitempty"`
	Occupation                 string                   `validate:"required" json:"occupation,omitempty"`
	IsPoliticallyExposedPerson *bool                    `json:"isPoliticallyExposedPerson,omitempty"` //deprecated
	PoliticallyExposedPerson   PoliticallyExposedPerson `validate:"required,dive" json:"pep,omitempty"`
	Phone                      *Phone                   `validate:"required,dive" json:"phone,omitempty"`
	Address                    *Address                 `validate:"required,dive" json:"address,omitempty"`
	MemberQualification        *string                  `json:"memberQualification,omitempty"`
	PowerOfAttorney            *PowerOfAttorney         `json:"powerOfAttorney,omitempty"`
}

// PowerOfAttorney ...
type PowerOfAttorney struct {
	StartDate time.Time `json:"startDate,omitempty"`
	EndDate   time.Time `json:"endDate,omitempty"`
}

// Owners ...
type Owners struct {
	LegalRepresentative
}

// SimpleLegalRepresentative ...
type SimpleLegalRepresentative struct {
	DocumentNumber           string                      `validate:"required,cnpjcpf" json:"documentNumber,omitempty"` //deprecated
	Document                 LegalRepresentativeDocument `validate:"required,dive" json:"document,omitempty"`
	RegisterName             string                      `validate:"required" json:"registerName,omitempty"`
	SocialName               string                      `json:"socialName,omitempty"`
	Email                    string                      `validate:"required" json:"email,omitempty"`
	Phone                    *Phone                      `validate:"required,dive" json:"phone,omitempty"`
	Address                  *Address                    `validate:"required,dive" json:"address,omitempty"`
	BirthDate                time.Time                   `validate:"required" json:"birthDate,omitempty"`
	MotherName               string                      `validate:"required" json:"motherName,omitempty"`
	DeclaredIncome           DeclaredIncome              `validate:"required" json:"declaredIncome,omitempty"`
	Occupation               string                      `validate:"required" json:"occupation,omitempty"`
	PoliticallyExposedPerson PoliticallyExposedPerson    `validate:"required,dive" json:"pep,omitempty"`
	Documentation            Documentation               `validate:"required,dive" json:"documentation,omitempty"`
}

// LegalRepresentativeDocument ...
type LegalRepresentativeDocument struct {
	Value string `validate:"required,cnpjcpf" json:"value,omitempty"`
	Type  string `validate:"required" json:"type,omitempty"`
}

// BusinessResponse ...
type BusinessResponse struct {
	ResultLevel   ResultLevel  `json:"resultLevel,omitempty"`
	Document      string       `json:"documentNumber,omitempty"`
	BusinessName  string       `json:"businessName,omitempty"`
	TradingName   string       `json:"tradingName,omitempty"`
	BusinessEmail string       `json:"businessEmail,omitempty"`
	Status        string       `json:"status,omitempty"`
	Reasons       []string     `json:"reasons,omitempty"`
	BusinessType  BusinessType `json:"businessType"`
	BusinessSize  BusinessSize `json:"businessSize"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
}

// BusinessAccountRequest ...
type BusinessAccountRequest struct {
	Document    string      `validate:"required,cnpj" json:"documentNumber,omitempty"`
	AccountType AccountType `validate:"required" json:"accountType"`
}

// BoletoType ...
type BoletoType string

const (
	Deposit BoletoType = "Deposit" // depósito (quando emitente é o próprio pagador)
	Levy    BoletoType = "Levy"    // cobrança (quando emitente é diferente do pagador)
	Invoice BoletoType = "Invoice" // cartão de crédito
)

// BoletoRequest ...
type BoletoRequest struct {
	APIVersion   *string          `validate:"required" json:"api_version,omitempty"`
	Account      *Account         `validate:"required" json:"account,omitempty"`
	Document     string           `validate:"required,cnpjcpf" json:"documentNumber,omitempty"`
	Amount       float64          `validate:"required" json:"amount,omitempty"`
	DueDate      time.Time        `validate:"required" json:"dueDate,omitempty"`
	Type         BoletoType       `validate:"required" json:"type,omitempty"`
	Alias        *string          `json:"alias,omitempty"`
	Payer        *BoletoPayer     `json:"payer,omitempty"`
	Fine         *BoletoFine      `json:"fine,omitempty"`
	Interest     *BoletoInterest  `json:"interest,omitempty"`
	Discount     *BoletoDiscounts `json:"discount,omitempty"`  // api-version 2.0
	Discounts    *BoletoDiscounts `json:"discounts,omitempty"` // deprecated // api-version 1.0
	ClosePayment time.Time        `json:"closePayment,omitempty"`
}

// BoletoResponse ...
type BoletoResponse struct {
	AuthenticationCode string   `json:"authenticationCode,omitempty"`
	Account            *Account `json:"account,omitempty"`
}

// BoletoInterestType ...
type BoletoInterestType string

const (
	FixedAmountInterestType BoletoInterestType = "FixedAmount"
	PercentInterestType     BoletoInterestType = "Percent"
)

// BoletoInterest ...
type BoletoInterest struct {
	StartDate time.Time          `validate:"required" json:"startDate,omitempty"`
	Value     float64            `validate:"required" json:"value,omitempty"`
	Type      BoletoInterestType `validate:"required" json:"type,omitempty"`
}

// BoletoFineType ...
type BoletoFineType string

const (
	FixedAmountFineType BoletoFineType = "FixedAmount"
	PercentFineType     BoletoFineType = "Percent"
)

// BoletoFine ...
type BoletoFine struct {
	StartDate time.Time      `validate:"required" json:"startDate,omitempty"`
	Value     float64        `validate:"required" json:"value,omitempty"`
	Type      BoletoFineType `validate:"required" json:"type,omitempty"`
}

// BoletoDiscountsType ...
type BoletoDiscountsType string

const (
	FixedAmountDiscountType  BoletoDiscountsType = "FixedAmountUntilLimitDate"
	FixedPercentDiscountType BoletoDiscountsType = "FixedPercentUntilLimitDate"
)

// BoletoDiscounts ...
type BoletoDiscounts struct {
	LimitDate time.Time           `json:"limitDate"`
	Value     float64             `validate:"required" json:"value,omitempty"`
	Type      BoletoDiscountsType `validate:"required" json:"type,omitempty"`
}

// BoletoAmount ...
type BoletoAmount struct {
	Value    float64 `json:"value,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

type PaymentChannel string

const (
	AgencyPaymentChannel               PaymentChannel = "Agency"
	SelfServiceTerminalPaymentChannel  PaymentChannel = "SelfServiceTerminal"
	InternetBankingPaymentChannel      PaymentChannel = "InternetBanking"
	CorrespondentBankingPaymentChannel PaymentChannel = "CorrespondentBanking"
	CallCenterPaymentChannel           PaymentChannel = "CallCenter"
	EletronicFilePaymentChannel        PaymentChannel = "EletronicFile"
	DDAPaymentChannel                  PaymentChannel = "DDA"
)

// BoletoPayment ...
type BoletoPayment struct {
	ID             string         `json:"id,omitempty"`
	Amount         float64        `json:"amount,omitempty"`
	PaymentChannel PaymentChannel `json:"paymentChannel,omitempty"`
	PaidOutDate    time.Time      `json:"paidOutDate,omitempty"`
}

// BoletoDetailedResponse ...
type BoletoDetailedResponse struct {
	Alias              *string          `json:"alias,omitempty"`
	AuthenticationCode string           `json:"authenticationCode,omitempty"`
	Barcode            string           `json:"barcode,omitempty"`
	Digitable          string           `json:"digitable,omitempty"`
	Status             string           `json:"status,omitempty"`
	Document           string           `json:"document,omitempty"`
	DueDate            time.Time        `json:"dueDate,omitempty"`
	ClosePayment       *time.Time       `json:"closePayment,omitempty"`
	EmissionDate       *time.Time       `json:"emissionDate,omitempty"`
	OurNumber          string           `json:"ourNumber,omitempty"`
	Type               BoletoType       `json:"type,omitempty"`
	Amount             *BoletoAmount    `json:"amount,omitempty"`
	Account            *Account         `json:"account,omitempty"`
	Payer              *Payer           `json:"payer,omitempty"`
	RecipientFinal     *Payer           `json:"recipientFinal,omitempty"`
	RecipientOrigin    *Payer           `json:"recipientOrigin,omitempty"`
	Payments           []*BoletoPayment `json:"payments,omitempty"`
	Fine               *BoletoFine      `json:"fine,omitempty"`
	Interest           *BoletoInterest  `json:"interest,omitempty"`
	Discount           *BoletoDiscounts `json:"discount,omitempty"`
	Discounts          *BoletoDiscounts `json:"discounts,omitempty"` // deprecated (api-version 1.0)
	UpdatedAt          time.Time        `json:"updatedAt,omitempty"`
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
	APIVersion         *string  `validate:"required" json:"api_version,omitempty"`
	AuthenticationCode string   `validate:"required" json:"authenticationCode,omitempty"`
	Account            *Account `validate:"required" json:"account,omitempty"`
}

// CancelBoletoRequest ...
type CancelBoletoRequest struct {
	APIVersion         *string  `validate:"required" json:"api_version,omitempty"`
	AuthenticationCode string   `validate:"required" json:"authenticationCode,omitempty"`
	Account            *Account `validate:"required" json:"account,omitempty"`
}

// SimulatePaymentRequest ...
type SimulatePaymentRequest struct {
	APIVersion         *string  `validate:"required" json:"api_version,omitempty"`
	AuthenticationCode string   `validate:"required" json:"authenticationCode,omitempty"`
	Account            *Account `validate:"required" json:"account,omitempty"`
}

// FilterBankListRequest ...
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

// BankDataResponse ...
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

type EventStatus string

const (
	// Active sucesso. evento realizado.
	Active EventStatus = "ACTIVE"
	// Canceled erro. evento não realizado por algum motivo.
	Canceled EventStatus = "CANCELED"
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
	Status         EventStatus            `json:"status,omitempty"`
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
	Status         *EventStatus
	EventName      []string
}

// DocumentAnalysisRequest ...
type DocumentAnalysisRequest struct {
	Document              string       `validate:"required" json:"document,omitempty"`
	DocumentType          DocumentType `validate:"required" json:"document_type,omitempty"`
	DocumentSide          DocumentSide `validate:"required" json:"document_side,omitempty"`
	ImageFile             os.File      `validate:"required" json:"image_file,omitempty"`
	IsCorporationBusiness *bool        `json:"is_corporation_business,omitempty"`
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
	Assignor           string     `json:"assignor,omitempty"`
	Charges            *Charges   `json:"charges,omitempty"`
	SettleDate         time.Time  `json:"settleDate,omitempty"`
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

// CreateTicketRequest ...
type CreateTicketRequest struct {
	GroupID     int      `json:"group_id"`
	ProductID   int      `json:"product_id"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Subject     string   `json:"subject"`
	Email       string   `json:"email"`
	CopyEmails  []string `json:"copy_emails"`
	Priority    int      `json:"priority"`
	Status      int      `json:"status"`
	Attachments []string `json:"attachments"`
	Request     string   `json:"cf_solicitacao"`
	Request2    string   `json:"cf_solicitacao2"`
	CompanyKey  string   `json:"cf_companykey"`
	RazaoSocial string   `json:"cf_razao_social"`
	CNPJ        string   `json:"cf_cnpj_empresa"`
	Cellphone   int      `json:"cf_celular"`
	Last4Digits int      `json:"cf_comprastransaesestornos_4_ltimos_dgitos_do_carto"`
}

// CreateTicketResponse ...
type CreateTicketResponse struct {
	ID int `json:"id"`
}

// GetTicketResponse ...
type GetTicketResponse struct {
	ID            int    `json:"id"`
	Status        int    `json:"status"`
	Reason        string `json:"cf_sdkycpjbankly"`
	AccountBranch string `json:"cf_sdkycpjbanklyagencia"`
	AccountNumber string `json:"cf_sdkycpjbanklyconta"`
}

// FilterTicketsRequest ...
type FilterTicketsRequest struct {
	Status int
}

// FilterTicketsResponse ...
type FilterTicketsResponse struct {
	Total   int                  `json:"total"`
	Results []*GetTicketResponse `json:"results"`
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

type CardTrackingResponse struct {
	CreatedDate           time.Time            `json:"createdDate,omitempty"`
	Name                  string               `json:"name,omitempty"`
	Alias                 string               `json:"alias,omitempty"`
	EstimatedDeliveryDate time.Time            `json:"estimatedDeliveryDate,omitempty"`
	Function              string               `json:"function,omitempty"`
	ExternalTracking      CardExternalTracking `json:"externalTracking,omitempty"`
	Address               []CardAddress        `json:"address,omitempty"`
	Status                []CardTrackingStatus `json:"status,omitempty"`
	Finalized             []Finalized          `json:"finalized,omitempty"`
}

type CardExternalTracking struct {
	Code    string `json:"code,omitempty"`
	Partner string `json:"partner,omitempty"`
}

type CardTrackingStatus struct {
	CreatedDate time.Time `json:"createdDate,omitempty"`
	Type        string    `json:"type,omitempty"`
	Reason      string    `json:"reason,omitempty"`
}

type Finalized struct {
	RecipientName    string    `json:"recipientName,omitempty"`
	RecipientKinship string    `json:"recipientKinship,omitempty"`
	Identifier       string    `json:"documentNumber,omitempty"`
	Attempts         int       `json:"attempts,omitempty"`
	CreatedDate      time.Time `json:"createdDate,omitempty"`
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

type CardDuplicateDTO struct {
	Status         DuplicateCardStatus `json:"status"`
	DocumentNumber string              `json:"documentNumber"`
	Description    string              `json:"description"`
	Password       string              `json:"password"`
	Address        CardAddress         `json:"address"`
}

type CardDuplicateResponse struct {
	Proxy        string `json:"proxy"`
	ActivateCode string `json:"activateCode"`
}

type CardActivateDTO struct {
	Password     string `json:"password"`
	ActivateCode string `json:"activateCode"`
}

type CardContactlessDTO struct {
	Active bool `json:"active"`
}

type CardUpdatePasswordDTO struct {
	Password string `validate:"required,numeric,min=4,max=4" json:"password"`
}

type CardCreateRequest struct {
	DocumentNumber string      `json:"documentNumber"`
	CardName       string      `json:"cardName"`
	Alias          string      `json:"alias"`
	BankAgency     string      `json:"bankAgency"`
	BankAccount    string      `json:"bankAccount"`
	ProgramId      int16       `json:"programId,omitempty"`
	Password       string      `json:"password"`
	Address        CardAddress `json:"address"`
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
	Type  PixType `json:"type"`
	Value string  `json:"value"`
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
	DocumentType   string                    `json:"documentType"`
	DocumentNumber string                    `json:"documentNumber"`
	Name           string                    `json:"name"`
}

type PixCashOutRecipientResponse struct {
	Account        PixCashOutAccountResponse `json:"account"`
	Bank           PixCashOutBankResponse    `json:"bank"`
	DocumentType   string                    `json:"documentType"`
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

// PixAdditionalDataValue ...
type PixAdditionalDataValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// PixAddress ...
type PixAddress struct {
	AddressLine string `json:"addressLine"`
	State       string `json:"state"`
	City        string `json:"city"`
	ZipCode     string `json:"zipCode"`
}

// PixPayer ...
type PixPayer struct {
	Name           string     `json:"name"`
	DocumentNumber string     `json:"documentNumber"`
	Type           string     `json:"type"`
	Address        PixAddress `json:"address"`
}

// PixQrCodeStaticRequest ...
type PixQrCodeStaticRequest struct {
	AddressingKey     PixTypeValue             `validate:"required" json:"addressingKey"`
	Amount            float64                  `json:"amount"`
	RecipientName     string                   `validate:"required" json:"recipientName"`
	PixQrCodeLocation PixQrCodeLocation        `validate:"required" json:"location"`
	ConciliationID    string                   `json:"conciliationId"`
	CategoryCode      *string                  `json:"categoryCode"`
	AdditionalData    []PixAdditionalDataValue `json:"additionalData,omitempty"`
}

// PixQrCodeDynamicRequest ...
type PixQrCodeDynamicRequest struct {
	AddressingKey    PixTypeValue             `json:"addressingKey"`
	Amount           float64                  `json:"amount"`
	RecipientName    string                   `json:"recipientName"`
	ConciliationID   string                   `json:"conciliationId"`
	AdditionalData   []PixAdditionalDataValue `json:"additionalData,omitempty"`
	SinglePayment    bool                     `json:"singlePayment"`
	Payer            PixPayer                 `json:"payer"`
	ChangeAmountType string                   `json:"changeAmountType"`
	ExpiresAt        string                   `json:"expiresAt"`
	PayerRequestText string                   `json:"payerRequestText,omitempty"`
}

// PixQrCodeResponse ...
type PixQrCodeResponse struct {
	EncodedValue string `json:"encodedValue"`
}

type PixQrCodeDecodeRequest struct {
	EncodedValue string `json:"encodedValue"`
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

type PixQrCodeBankResponse struct {
	Name string `json:"name"`
	Ispb string `json:"ispb"`
}

type PixQrCodePaymentResponse struct {
	BaseValue       float64 `json:"baseValue"`
	InterestValue   float64 `json:"interestValue"`
	PenaltyValue    float64 `json:"penaltyValue"`
	DiscountValue   float64 `json:"discountValue"`
	ReductionValue  float64 `json:"reductionValue"`
	TotalValue      float64 `json:"totalValue"`
	DueDate         string  `json:"dueDate"`
	ChangeValue     float64 `json:"changeValue"`
	WithdrawalValue float64 `json:"withdrawalValue"`
}

type PixQrCodeLocation struct {
	City    string `json:"city"`
	ZipCode string `json:"zipCode"`
}

type PixQrCodeLocationResponse struct {
	City    string `json:"city"`
	ZipCode string `json:"zipCode"`
}

type PixQrCodeDecodeResponse struct {
	EndToEndID         string                      `json:"endToEndId"`
	ConciliationID     string                      `json:"conciliationId"`
	AddressingKey      PixTypeValue                `json:"addressingKey"`
	QrCodeType         string                      `json:"qrCodeType"`
	Payer              PixPayer                    `json:"payer"`
	Holder             PixHolder                   `json:"holder"`
	Bank               PixQrCodeBankResponse       `json:"bank"`
	Payment            PixQrCodePaymentResponse    `json:"payment"`
	Location           PixQrCodeLocationResponse   `json:"location"`
	WithdrawalDetail   PixQrCodeWithdrawalDetail   `json:"withdrawalDetail"`
	ChangeAmountDetail PixQrCodeChangeAmountDetail `json:"changeAmountDetail"`
	QrCodePurpose      string                      `json:"qrCodePurpose"`
}

type PixQrCodeWithdrawalDetail struct {
	ChangeAmountType string `json:"changeAmountType"`
	AgentType        string `json:"agentType"`
	ProviderISPB     string `json:"providerIspb"`
}

type PixQrCodeChangeAmountDetail struct {
	ChangeAmountType string `json:"changeAmountType"`
	AgentType        string `json:"agentType"`
}

type PixCashOutByAuthenticationCodeResponse struct {
	CompanyKey         string                      `json:"companyKey"`
	AuthenticationCode string                      `json:"authenticationCode"`
	EndToEndID         string                      `json:"endToEndId"`
	InitializationType string                      `json:"initializationType"`
	Amount             float64                     `json:"amount"`
	Description        string                      `json:"description"`
	CorrelationID      string                      `json:"correlationId"`
	Sender             PixCashOutSenderResponse    `json:"sender"`
	Recipient          PixCashOutRecipientResponse `json:"recipient"`
	Channel            string                      `json:"channel"`
	Status             TransfersStatus             `json:"status"`
	Type               string                      `json:"type"`
	CreatedAt          time.Time                   `json:"createdAt"`
	UpdatedAt          time.Time                   `json:"updatedAt"`
}

// Pix Request
type PixAddressKeyCreateRequest struct {
	AddressingKey PixTypeValue `json:"addressingKey"`
	Account       Account      `json:"account"`
}

// Claim Request
type PixClaimRequest struct {
	Type          PixClaimType `json:"type"`
	AddressingKey PixTypeValue `json:"addressingKey"`
	Claimer       Claimer      `json:"claimer"`
}

// Claimer
type Claimer struct {
	Branch string      `json:"branch"`
	Number string      `json:"number"`
	Bank   BankClaimer `json:"bank"`
}

// BankClaimer
type BankClaimer struct {
	Name string `json:"name"`
	Ispb string `json:"ispb"`
}

// PostPixClaimerResponse
type PixClaimResponse struct {
	ClaimId             string       `json:"claimId"`
	Type                PixClaimType `json:"type"`
	AddressingKey       PixTypeValue `json:"addressingKey"`
	Claimer             Claimer      `json:"claimer"`
	Status              StatusClaim  `json:"status"`
	CreatedAt           string       `json:"createdAt"`
	ResolutionLimitDate string       `json:"resolutionLimitDate"`
	ConclusionLimitDate string       `json:"conclusionLimitDate"`
}

// PixClaimConfirmResponse
type PixClaimConfirmResponse struct {
	PixClaimResponse
	Donor          Claimer     `json:"donor"`
	PreviousStatus StatusClaim `json:"previousStatus"`
	ConfirmReason  string      `json:"confirmReason"`
	ConfirmedBy    string      `json:"confirmedBy"`
	UpdatedAt      string      `json:"updatedAt"`
	ConfirmedAt    string      `json:"confirmedAt"`
}

// PixClaimCompleteResponse
type PixClaimCompleteResponse struct {
	PixClaimConfirmResponse
	CompletedAt string `json:"completedAt"`
}

// PixClaimCompleteResponse
type PixClaimCancelResponse struct {
	PixClaimResponse
	Donor          Claimer      `json:"donor"`
	PreviousStatus StatusClaim  `json:"previousStatus"`
	UpdatedAt      string       `json:"updatedAt"`
	CancelReason   CancelReason `json:"cancelReason"`
	CanceledBy     string       `json:"canceledBy"`
	CanceledAt     string       `json:"canceledAt"`
}

// PixClaimCancelReason
type PixClaimCancelReason struct {
	Reason string `json:"reason"`
}

// PixClaimConfirmReason
type PixClaimConfirmReason struct {
	Reason string `json:"reason"`
}

// PostPixResponse
type PixAddressKeyCreateResponse struct {
	AddressingKey PixTypeValue `json:"addressingKey"`
	Account       struct {
		Branch string `json:"branch"`
		Number string `json:"number"`
		Type   string `json:"type"`
		Holder struct {
			Type           string `json:"type"`
			DocumentNumber string `json:"documentNumber"`
			Name           string `json:"name"`
			SocialName     string `json:"socialName"`
		} `json:"holder"`
		Bank struct {
			Ispb  string `json:"ispb"`
			Compe string `json:"compe"`
			Name  string `json:"name"`
		} `json:"bank"`
	} `json:"account"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	OwnedAt   time.Time `json:"ownedAt"`
}

// TransactionalHashResponse
type TransactionalHash struct {
	Hash string `json:"hash"`
	Code string `json:"code"`
}

// TransactionalHashValidate
type TransactionalHashValidateResponse struct {
	Hash                string `json:"hash"`
	ExpirationInSeconds int    `json:"expirationInSeconds"`
}

// TransactionalHashRequest
type TransactionalHashRequest struct {
	Context   string                `json:"context"`
	Operation string                `json:"operation"`
	Data      TransactionalHashData `json:"data"`
}

// TransactionalHashData
type TransactionalHashData struct {
	AddressingKey TransactionalHashAddressingKey `json:"addressingKey"`
}

// TransactionalHashAddressingKey
type TransactionalHashAddressingKey struct {
	Type  PixType `json:"type"`
	Value string  `json:"value"`
}

// DeleteAddressKeyResponse
type DeleteAddressKeyResponse struct {
	EndToEndID    string `json:"endToEndId"`
	AddressingKey struct {
		Type string `json:"type"`
	} `json:"addressingKey"`
	Holder struct {
		Type       string `json:"type"`
		Name       string `json:"name"`
		SocialName string `json:"socialName"`
		Document   string `json:"document"`
	} `json:"holder"`

	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	OwnedAt   time.Time `json:"ownedAt"`
}

type AddressingKeyValue string

type ClientRegisterBanklyRequest struct {
	GrantTypes              []string `json:"grant_types"`
	TLSClientAuthSubjectDn  string   `json:"tls_client_auth_subject_dn"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	ResponseTypes           []string `json:"response_types"`
	CompanyKey              string   `json:"company_key"`
	Scope                   string   `json:"scope"`
}

type ClientRegisterResponse struct {
	GrantTypes                       []string `json:"grant_types"`
	SubjectType                      string   `json:"subject_type"`
	TLSClientAuthSubjectDn           string   `json:"tls_client_auth_subject_dn"`
	RegistrationClientURI            string   `json:"registration_client_uri"`
	CompanyKey                       string   `json:"company_key"`
	RegistrationAccessTokenExpiresIn int64    `json:"registration_access_token_expires_in"`
	RegistrationAccessToken          string   `json:"registration_access_token"`
	ClientID                         string   `json:"client_id"`
	TokenEndpointAuthMethod          string   `json:"token_endpoint_auth_method"`
	RequireProofKey                  bool     `json:"require_proof_key"`
	Scope                            string   `json:"scope"`
	TokenEndpointAuthMethods         []string `json:"token_endpoint_auth_methods"`
	ClientIDIssuedAt                 int64    `json:"client_id_issued_at"`
	AccessTokenTTL                   int64    `json:"access_token_ttl"`
	ResponseTypes                    []string `json:"response_types"`
}

// Certificate ...
type Certificate struct {
	Certificate      string `json:"certificate"`
	CertificateChain string `json:"certificateChain"`
	SubjectDn        string `json:"subjectDn"`
	PrivateKey       string `json:"privateKey"`
	Passphrase       string `json:"passphrase"`
	UUID             string `json:"uuid"`
	ClientID         string `json:"client_id"`
}

// SandboxSimulateBankslipPaymentRequest ...
type SandboxSimulateBankslipPaymentRequest struct {
	APIVersion         *string `validate:"required" json:"api_version,omitempty"`
	AuthenticationCode string  `validate:"required" json:"authenticationCode"`
	Account            Account `validate:"required" json:"account"`
}

// ParseSimpleBusinessRequest ...
func ParseSimpleBusinessRequest(r *BusinessRequest) *SimpleBusinessRequest {
	if r == nil {
		return nil
	}
	resp := &SimpleBusinessRequest{
		DocumentNumber:        r.DocumentNumber,
		BusinessName:          r.BusinessName,
		TradingName:           r.TradingName,
		BusinessEmail:         r.BusinessEmail,
		BusinessType:          r.BusinessType,
		BusinessSize:          r.BusinessSize,
		BusinessAddress:       r.BusinessAddress,
		DeclaredAnnualBilling: r.DeclaredAnnualBilling,
	}
	if len(r.LegalRepresentatives) > 0 {
		resp.LegalRepresentative = ParseSimpleLegalRepresentative(r.LegalRepresentatives[0])
	}
	newBusinessName := NormalizeNameWithoutSpecialCharacters(aws.String(resp.BusinessName))
	if newBusinessName != nil {
		resp.BusinessName = *newBusinessName
	}
	newTradingName := NormalizeNameWithoutSpecialCharacters(aws.String(resp.TradingName))
	if newTradingName != nil {
		resp.TradingName = *newTradingName
	}
	return resp
}

// ParseSimpleLegalRepresentative ...
func ParseSimpleLegalRepresentative(rep *LegalRepresentative) *SimpleLegalRepresentative {
	if rep == nil {
		return nil
	}
	response := &SimpleLegalRepresentative{
		DocumentNumber: rep.DocumentNumber,
		Document: LegalRepresentativeDocument{
			Value: rep.DocumentNumber,
			Type:  "CPF",
		},
		RegisterName:             rep.RegisterName,
		SocialName:               rep.SocialName,
		Phone:                    rep.Phone,
		Address:                  rep.Address,
		BirthDate:                rep.BirthDate,
		MotherName:               rep.MotherName,
		Email:                    rep.Email,
		DeclaredIncome:           rep.DeclaredIncome,
		Occupation:               rep.Occupation,
		PoliticallyExposedPerson: rep.PoliticallyExposedPerson,
	}
	if rep.Documentation != nil {
		response.Documentation = *rep.Documentation
	}
	return response
}

// ParseCorporationBusinessSize ...
func ParseCorporationBusinessSize(size BusinessSize) CorporationBusinessSize {
	switch size {
	case BusinessSizeMEI:
		return CorporationBusinessSizeSMALL
	case BusinessSizeME:
		return CorporationBusinessSizeMEDIUM
	case BusinessSizeEPP:
		return CorporationBusinessSizeLARGE
	default:
		return CorporationBusinessSizeMEDIUM
	}
}
