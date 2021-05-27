package bankly_test

import (
	"os"
	"testing"

	"github.com/contbank/grok"

	"github.com/contbank/bankly-sdk"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DocumentAnalysisTestSuite struct {
	suite.Suite
	assert           *assert.Assertions
	session          *bankly.Session
	documentAnalysis *bankly.DocumentAnalysis
}

func TestDocumentAnalysisTestSuite(t *testing.T) {
	suite.Run(t, new(DocumentAnalysisTestSuite))
}

func (s *DocumentAnalysisTestSuite) SetupTest() {
	s.assert = assert.New(s.T())

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(os.Getenv("BANKLY_CLIENT_ID")),
		ClientSecret: bankly.String(os.Getenv("BANKLY_CLIENT_SECRET")),
	})

	s.assert.NoError(err)

	s.session = session
	s.documentAnalysis = bankly.NewDocumentAnalysis(*s.session)
}

func (s *DocumentAnalysisTestSuite) TestSendDocumentAnalysis() {

	documentNumber := grok.GeneratorCPF()
	request := bankly.DocumentAnalysisRequest{
		DocumentType: bankly.DocumentTypeSELFIE,
		DocumentSide: bankly.DocumentSideFront,
		Image:        getSelfieBase64(),
	}

	balance, err := s.documentAnalysis.SendDocumentAnalysis(documentNumber, request)

	s.assert.NoError(err)
	s.assert.NotNil(balance)
}

func (s *DocumentAnalysisTestSuite) TestFindDocumentAnalysis() {
	documentNumber := "48195413315"
	documentAnalysisToken := "EXZePni82uzF5ndyXInNXPN_PmUIeGcy"
	balance, err := s.documentAnalysis.FindDocumentAnalysis(documentNumber, documentAnalysisToken)

	s.assert.NoError(err)
	s.assert.NotNil(balance)
}
