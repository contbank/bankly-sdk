package bankly_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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
	s3session  		 *session.Session
	documentAnalysis *bankly.DocumentAnalysis
	s3manager  	     *bankly.S3Manager
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

	s.s3session = session.Must(session.NewSession(&aws.Config{
		Region:           bankly.String(os.Getenv("AWS_S3_REGION")),
		Endpoint:         bankly.String(os.Getenv("AWS_S3_ENDPOINT")),
		S3ForcePathStyle: aws.Bool(true),
	}))
	s.s3manager = bankly.NewS3Manager(s.s3session)

	s.documentAnalysis = bankly.NewDocumentAnalysis(*s.bankSession, *s.s3manager)
}

func (s *DocumentAnalysisTestSuite) TestDocumentAnalysisBucketExists(){
	buckets, error := s.s3manager.ListBuckets()

	exists := false
	for _, bucketName := range buckets {
		if bucketName == bankly.DocumentAnalysisBucket {
			exists = true
		}
	}

	s.assert.NoError(error)
	s.assert.True(exists)
}

func (s *DocumentAnalysisTestSuite) TestSendDocumentAnalysis() {

	docType := bankly.DocumentTypeSELFIE
	docSide := bankly.DocumentSideFront
	documentNumber := grok.GeneratorCPF()

	request := bankly.DocumentAnalysisRequest{
		DocumentType : docType,
		DocumentSide : docSide,
		//// URLImage : "test_images/selfie1.jpeg",
		URLImage : "https://s3-us-west-2.amazonaws.com/temp.documentanalysis/contbank.png",
	}

	response, err := s.documentAnalysis.SendDocumentAnalysis(documentNumber, request)

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
	s.assert.Equal(docType, bankly.DocumentType(respDocAnalysis.DocumentType))
	s.assert.Equal(docSide, bankly.DocumentSide(respDocAnalysis.DocumentSide))
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
	s.assert.Equal(docType, bankly.DocumentType(respDocAnalysis.DocumentType))
	s.assert.Equal(docSide, bankly.DocumentSide(respDocAnalysis.DocumentSide))
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
	s.assert.Equal(docType, bankly.DocumentType(respDocAnalysis.DocumentType))
	s.assert.Equal(docSide, bankly.DocumentSide(respDocAnalysis.DocumentSide))
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
	s.assert.Equal(docType, bankly.DocumentType(respDocAnalysis.DocumentType))
	s.assert.Equal(docSide, bankly.DocumentSide(respDocAnalysis.DocumentSide))
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
	s.assert.Equal(docType, bankly.DocumentType(respDocAnalysis.DocumentType))
	s.assert.Equal(docSide, bankly.DocumentSide(respDocAnalysis.DocumentSide))
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
	s.assert.Equal(docType, bankly.DocumentType(respDocAnalysis.DocumentType))
	s.assert.Equal(docSide, bankly.DocumentSide(respDocAnalysis.DocumentSide))
	s.assert.Equal(documentNumber, respDocAnalysis.DocumentNumber)
}

func (s *DocumentAnalysisTestSuite) createDocumentAnalysis(documentNumber string,
	docType bankly.DocumentType, docSide bankly.DocumentSide) *bankly.DocumentAnalysisResponse {

	request := bankly.DocumentAnalysisRequest{
		DocumentType : docType,
		DocumentSide : docSide,
		URLImage : "test_images/selfie1.jpeg",
	}
	resp, err := s.documentAnalysis.SendDocumentAnalysis(documentNumber, request)
	s.assert.NoError(err)
	s.assert.NotNil(resp)
	s.assert.NotNil(resp.Token)

	return resp
}
