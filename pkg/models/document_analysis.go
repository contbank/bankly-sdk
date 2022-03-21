package bankly

import "os"

const (
	// DocumentAnalysisPath ...
	DocumentAnalysisPath = "/document-analysis"
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
