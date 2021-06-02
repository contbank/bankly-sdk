package bankly_test

import (
	"os"
	"testing"
	"time"

	"github.com/contbank/grok"

	"github.com/contbank/bankly-sdk"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DocumentAnalysisTestSuite struct {
	suite.Suite
	assert           *assert.Assertions
	bankSession      *bankly.Session
	documentAnalysis *bankly.DocumentAnalysis
}

func TestDocumentAnalysisTestSuite(t *testing.T) {
	suite.Run(t, new(DocumentAnalysisTestSuite))
}

func (s *DocumentAnalysisTestSuite) SetupTest() {
	s.assert = assert.New(s.T())

	newSession, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(os.Getenv("BANKLY_CLIENT_ID")),
		ClientSecret: bankly.String(os.Getenv("BANKLY_CLIENT_SECRET")),
	})

	s.assert.NoError(err)

	s.bankSession = newSession

	s.documentAnalysis = bankly.NewDocumentAnalysis(*s.bankSession)
}

func (s *DocumentAnalysisTestSuite) TestSendDocumentAnalysis() {
	docType := bankly.DocumentTypeSELFIE
	docSide := bankly.DocumentSideFront
	documentNumber := grok.GeneratorCPF()

	imageFile, errFile := os.Open("test_images/selfie1.jpeg")
	s.assert.NoError(errFile)

	request := bankly.DocumentAnalysisRequest{
		Document : documentNumber,
		DocumentType : docType,
		DocumentSide : docSide,
		ImageFile : *imageFile,
	}

	response, err := s.documentAnalysis.SendDocumentAnalysis(request)

	s.assert.NoError(err)
	s.assert.NotNil(response)
	s.assert.Equal(docType, bankly.DocumentType(response.DocumentType))
	s.assert.Equal(docSide, bankly.DocumentSide(response.DocumentSide))
	s.assert.Equal(documentNumber, response.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysis_SELFIE_FRONT() {
	// create document analysis
	docType := bankly.DocumentTypeSELFIE
	docSide := bankly.DocumentSideFront
	documentNumber := grok.GeneratorCPF()
	resp := s.createDocumentAnalysis(documentNumber, docType, docSide)

	time.Sleep(time.Millisecond)

	// find document analysis
	respDocAnalysis, errDocAnalysis := s.documentAnalysis.FindDocumentAnalysis(documentNumber, resp.Token)

	s.assert.NoError(errDocAnalysis)
	s.assert.NotNil(respDocAnalysis)
	s.assert.NotEmpty(respDocAnalysis.Token)
	s.assert.NotEmpty(respDocAnalysis.Status)
	s.assert.Equal(documentNumber, respDocAnalysis.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysis_SELFIE_BACK() {
	// create document analysis
	docType := bankly.DocumentTypeSELFIE
	docSide := bankly.DocumentSideBack
	documentNumber := grok.GeneratorCPF()
	resp := s.createDocumentAnalysis(documentNumber, docType, docSide)

	time.Sleep(time.Millisecond)

	// find document analysis
	respDocAnalysis, errDocAnalysis := s.documentAnalysis.FindDocumentAnalysis(documentNumber, resp.Token)

	s.assert.NoError(errDocAnalysis)
	s.assert.NotNil(respDocAnalysis)
	s.assert.NotEmpty(respDocAnalysis.Token)
	s.assert.NotEmpty(respDocAnalysis.Status)
	s.assert.Equal(documentNumber, respDocAnalysis.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysis_CNH_FRONT() {
	// create document analysis
	docType := bankly.DocumentTypeCNH
	docSide := bankly.DocumentSideFront
	documentNumber := grok.GeneratorCPF()
	resp := s.createDocumentAnalysis(documentNumber, docType, docSide)

	time.Sleep(time.Millisecond)

	// find document analysis
	respDocAnalysis, errDocAnalysis := s.documentAnalysis.FindDocumentAnalysis(documentNumber, resp.Token)

	s.assert.NoError(errDocAnalysis)
	s.assert.NotNil(respDocAnalysis)
	s.assert.NotEmpty(respDocAnalysis.Token)
	s.assert.NotEmpty(respDocAnalysis.Status)
	s.assert.Equal(documentNumber, respDocAnalysis.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysis_CNH_BACK() {
	// create document analysis
	docType := bankly.DocumentTypeCNH
	docSide := bankly.DocumentSideBack
	documentNumber := grok.GeneratorCPF()
	resp := s.createDocumentAnalysis(documentNumber, docType, docSide)

	time.Sleep(time.Millisecond)

	// find document analysis
	respDocAnalysis, errDocAnalysis := s.documentAnalysis.FindDocumentAnalysis(documentNumber, resp.Token)

	s.assert.NoError(errDocAnalysis)
	s.assert.NotNil(respDocAnalysis)
	s.assert.NotEmpty(respDocAnalysis.Token)
	s.assert.NotEmpty(respDocAnalysis.Status)
	s.assert.Equal(documentNumber, respDocAnalysis.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysis_RG_FRONT() {
	// create document analysis
	docType := bankly.DocumentTypeRG
	docSide := bankly.DocumentSideFront
	documentNumber := grok.GeneratorCPF()
	resp := s.createDocumentAnalysis(documentNumber, docType, docSide)

	time.Sleep(time.Millisecond)

	// find document analysis
	respDocAnalysis, errDocAnalysis := s.documentAnalysis.FindDocumentAnalysis(documentNumber, resp.Token)

	s.assert.NoError(errDocAnalysis)
	s.assert.NotNil(respDocAnalysis)
	s.assert.NotEmpty(respDocAnalysis.Token)
	s.assert.NotEmpty(respDocAnalysis.Status)
	s.assert.Equal(documentNumber, respDocAnalysis.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysis_RG_BACK() {
	// create document analysis
	docType := bankly.DocumentTypeRG
	docSide := bankly.DocumentSideBack
	documentNumber := grok.GeneratorCPF()
	resp := s.createDocumentAnalysis(documentNumber, docType, docSide)

	time.Sleep(time.Millisecond)

	// find document analysis
	respDocAnalysis, errDocAnalysis := s.documentAnalysis.FindDocumentAnalysis(documentNumber, resp.Token)

	s.assert.NoError(errDocAnalysis)
	s.assert.NotNil(respDocAnalysis)
	s.assert.NotEmpty(respDocAnalysis.Token)
	s.assert.NotEmpty(respDocAnalysis.Status)
	s.assert.Equal(documentNumber, respDocAnalysis.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) createDocumentAnalysis(documentNumber string,
	docType bankly.DocumentType, docSide bankly.DocumentSide) *bankly.DocumentAnalysisResponse {

	imageFile, errFile := os.Open("test_images/selfie1.jpeg")
	s.assert.NoError(errFile)

	request := bankly.DocumentAnalysisRequest {
		Document : documentNumber,
		DocumentType : docType,
		DocumentSide : docSide,
		ImageFile : *imageFile,
	}

	resp, err := s.documentAnalysis.SendDocumentAnalysis(request)

	s.assert.NoError(err)
	s.assert.NotNil(resp)
	s.assert.NotNil(resp.Token)
	s.assert.Equal(docType, bankly.DocumentType(resp.DocumentType))
	s.assert.Equal(docSide, bankly.DocumentSide(resp.DocumentSide))
	s.assert.Equal(documentNumber, resp.DocumentNumber)

	return resp
}
