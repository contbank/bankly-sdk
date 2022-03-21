package bankly_test

import (
	"context"
	models "github.com/contbank/bankly-sdk/pkg/models"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	documentAnalysis "github.com/contbank/bankly-sdk/pkg/services/document-analysis"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/contbank/grok"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DocumentAnalysisTestSuite struct {
	suite.Suite
	assert           *assert.Assertions
	ctx              context.Context
	bankSession      *bankly.Session
	documentAnalysis *documentAnalysis.DocumentAnalysis
}

func TestDocumentAnalysisTestSuite(t *testing.T) {
	suite.Run(t, new(DocumentAnalysisTestSuite))
}

func (s *DocumentAnalysisTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(*utils.GetEnvBanklyClientID()),
		ClientSecret: bankly.String(*utils.GetEnvBanklyClientSecret()),
		Scopes:       bankly.String("kyc.document.write kyc.document.read"),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.bankSession = session

	s.documentAnalysis = documentAnalysis.NewDocumentAnalysis(httpClient, *s.bankSession)
}

func (s *DocumentAnalysisTestSuite) TestSendDocumentAnalysis() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de docs. Mockar este teste.")

	docType := models.DocumentTypeSELFIE
	docSide := models.DocumentSideFront
	documentNumber := grok.GeneratorCPF()

	imageFile, errFile := os.Open("../../test_files/contbank.png")
	s.assert.NoError(errFile)

	request := models.DocumentAnalysisRequest{
		Document:     documentNumber,
		DocumentType: docType,
		DocumentSide: docSide,
		ImageFile:    *imageFile,
	}

	response, err := s.documentAnalysis.SendDocumentAnalysis(s.ctx, request)

	logrus.
		WithFields(logrus.Fields{
			"request":  request,
			"response": response,
		}).
		Info("document sent")

	s.assert.NoError(err)
	s.assert.NotNil(response)
	s.assert.Equal(docType, models.DocumentType(response.DocumentType))
	s.assert.Equal(docSide, models.DocumentSide(response.DocumentSide))
	s.assert.Equal(documentNumber, response.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysis_SELFIE_FRONT() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de docs. Mockar este teste.")

	// create document analysis
	docType := models.DocumentTypeSELFIE
	docSide := models.DocumentSideFront
	documentNumber := grok.GeneratorCPF()
	resp := s.createDocumentAnalysis(documentNumber, docType, docSide)

	time.Sleep(time.Millisecond)

	// find document analysis
	respDocAnalysis, errDocAnalysis := s.documentAnalysis.FindDocumentAnalysis(s.ctx, documentNumber, resp.Token)

	s.assert.NoError(errDocAnalysis)
	s.assert.NotNil(respDocAnalysis)
	s.assert.NotEmpty(respDocAnalysis.Token)
	s.assert.NotEmpty(respDocAnalysis.Status)
	s.assert.Equal(docType, models.DocumentType(respDocAnalysis.DocumentType))
	s.assert.Equal(docSide, models.DocumentSide(respDocAnalysis.DocumentSide))
	s.assert.Equal(documentNumber, respDocAnalysis.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysis_SELFIE_BACK() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de docs. Mockar este teste.")

	// create document analysis
	docType := models.DocumentTypeSELFIE
	docSide := models.DocumentSideBack
	documentNumber := grok.GeneratorCPF()
	resp := s.createDocumentAnalysis(documentNumber, docType, docSide)

	time.Sleep(time.Millisecond)

	// find document analysis
	respDocAnalysis, errDocAnalysis := s.documentAnalysis.FindDocumentAnalysis(s.ctx, documentNumber, resp.Token)

	s.assert.NoError(errDocAnalysis)
	s.assert.NotNil(respDocAnalysis)
	s.assert.NotEmpty(respDocAnalysis.Token)
	s.assert.NotEmpty(respDocAnalysis.Status)
	s.assert.Equal(docType, models.DocumentType(respDocAnalysis.DocumentType))
	s.assert.Equal(docSide, models.DocumentSide(respDocAnalysis.DocumentSide))
	s.assert.Equal(documentNumber, respDocAnalysis.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysis_CNH_FRONT() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de docs. Mockar este teste.")

	// create document analysis
	docType := models.DocumentTypeCNH
	docSide := models.DocumentSideFront
	documentNumber := grok.GeneratorCPF()
	resp := s.createDocumentAnalysis(documentNumber, docType, docSide)

	time.Sleep(time.Millisecond)

	// find document analysis
	respDocAnalysis, errDocAnalysis := s.documentAnalysis.FindDocumentAnalysis(s.ctx, documentNumber, resp.Token)

	s.assert.NoError(errDocAnalysis)
	s.assert.NotNil(respDocAnalysis)
	s.assert.NotEmpty(respDocAnalysis.Token)
	s.assert.NotEmpty(respDocAnalysis.Status)
	s.assert.Equal(docType, models.DocumentType(respDocAnalysis.DocumentType))
	s.assert.Equal(docSide, models.DocumentSide(respDocAnalysis.DocumentSide))
	s.assert.Equal(documentNumber, respDocAnalysis.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysis_CNH_BACK() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de docs. Mockar este teste.")

	// create document analysis
	docType := models.DocumentTypeCNH
	docSide := models.DocumentSideBack
	documentNumber := grok.GeneratorCPF()
	resp := s.createDocumentAnalysis(documentNumber, docType, docSide)

	time.Sleep(time.Millisecond)

	// find document analysis
	respDocAnalysis, errDocAnalysis := s.documentAnalysis.FindDocumentAnalysis(s.ctx, documentNumber, resp.Token)

	s.assert.NoError(errDocAnalysis)
	s.assert.NotNil(respDocAnalysis)
	s.assert.NotEmpty(respDocAnalysis.Token)
	s.assert.NotEmpty(respDocAnalysis.Status)
	s.assert.Equal(docType, models.DocumentType(respDocAnalysis.DocumentType))
	s.assert.Equal(docSide, models.DocumentSide(respDocAnalysis.DocumentSide))
	s.assert.Equal(documentNumber, respDocAnalysis.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysis_RG_FRONT() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de docs. Mockar este teste.")

	// create document analysis
	docType := models.DocumentTypeRG
	docSide := models.DocumentSideFront
	documentNumber := grok.GeneratorCPF()
	resp := s.createDocumentAnalysis(documentNumber, docType, docSide)

	time.Sleep(time.Millisecond)

	// find document analysis
	respDocAnalysis, errDocAnalysis := s.documentAnalysis.FindDocumentAnalysis(s.ctx, documentNumber, resp.Token)

	s.assert.NoError(errDocAnalysis)
	s.assert.NotNil(respDocAnalysis)
	s.assert.NotEmpty(respDocAnalysis.Token)
	s.assert.NotEmpty(respDocAnalysis.Status)
	s.assert.Equal(docType, models.DocumentType(respDocAnalysis.DocumentType))
	s.assert.Equal(docSide, models.DocumentSide(respDocAnalysis.DocumentSide))
	s.assert.Equal(documentNumber, respDocAnalysis.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysis_RG_BACK() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de docs. Mockar este teste.")

	// create document analysis
	docType := models.DocumentTypeRG
	docSide := models.DocumentSideBack
	documentNumber := grok.GeneratorCPF()
	resp := s.createDocumentAnalysis(documentNumber, docType, docSide)
	s.assert.NotNil(resp)
	s.assert.NotEmpty(resp.Token)

	time.Sleep(time.Millisecond)

	// find document analysis
	respDocAnalysis, errDocAnalysis := s.documentAnalysis.FindDocumentAnalysis(s.ctx, documentNumber, resp.Token)

	s.assert.NoError(errDocAnalysis)
	s.assert.NotNil(respDocAnalysis)
	s.assert.NotEmpty(respDocAnalysis.Token)
	s.assert.NotEmpty(respDocAnalysis.Status)
	s.assert.Equal(docType, models.DocumentType(respDocAnalysis.DocumentType))
	s.assert.Equal(docSide, models.DocumentSide(respDocAnalysis.DocumentSide))
	s.assert.Equal(documentNumber, respDocAnalysis.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysisError_INVALID_DOCUMENT() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de docs. Mockar este teste.")

	respDocAnalysis, errDocAnalysis := s.documentAnalysis.FindDocumentAnalysis(s.ctx, grok.GeneratorCPF(), "TOKEN")

	s.assert.Error(errDocAnalysis)
	s.assert.Nil(respDocAnalysis)
	s.assert.Contains(errDocAnalysis.Error(), "not found")
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysisError_INVALID_TOKEN() {
	// TODO Mockar teste
	s.T().Skip("Bankly está retornando erro 403 no endpoint de docs. Mockar este teste.")

	// create document analysis
	docType := models.DocumentTypeRG
	docSide := models.DocumentSideBack
	documentNumber := grok.GeneratorCPF()
	resp := s.createDocumentAnalysis(documentNumber, docType, docSide)
	s.assert.NotNil(resp)
	s.assert.NotEmpty(resp.Token)

	time.Sleep(time.Millisecond)

	// find document analysis
	respDocAnalysis, errDocAnalysis := s.documentAnalysis.FindDocumentAnalysis(s.ctx, documentNumber, "INVALID_TOKEN")

	s.assert.Error(errDocAnalysis)
	s.assert.Nil(respDocAnalysis)
	s.assert.Contains(errDocAnalysis.Error(), "not found")
}

func (s *DocumentAnalysisTestSuite) createDocumentAnalysis(documentNumber string,
	docType models.DocumentType, docSide models.DocumentSide) *models.DocumentAnalysisResponse {

	imageFile, errFile := os.Open("../../test_files/contbank.png")
	s.assert.NoError(errFile)

	request := models.DocumentAnalysisRequest{
		Document:     documentNumber,
		DocumentType: docType,
		DocumentSide: docSide,
		ImageFile:    *imageFile,
	}

	resp, err := s.documentAnalysis.SendDocumentAnalysis(s.ctx, request)

	s.assert.NoError(err)
	s.assert.NotNil(resp)
	s.assert.NotNil(resp.Token)
	s.assert.Equal(docType, models.DocumentType(resp.DocumentType))
	s.assert.Equal(docSide, models.DocumentSide(resp.DocumentSide))
	s.assert.Equal(documentNumber, resp.DocumentNumber)

	return resp
}
