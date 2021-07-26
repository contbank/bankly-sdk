package bankly_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/contbank/bankly-sdk"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BankStatementTestSuite struct {
	suite.Suite
	assert  *assert.Assertions
	session *bankly.Session
	bank    *bankly.BankStatement
	boletos *bankly.Boletos
}

func TestBankStatementTestSuite(t *testing.T) {
	suite.Run(t, new(BankStatementTestSuite))
}

func (s *BankStatementTestSuite) SetupTest() {
	s.assert = assert.New(s.T())

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(os.Getenv("BANKLY_CLIENT_ID")),
		ClientSecret: bankly.String(os.Getenv("BANKLY_CLIENT_SECRET")),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.bank = bankly.NewBankStatement(httpClient, *s.session)
	s.boletos = bankly.NewBoletos(httpClient, *s.session)
}

func (s *BankStatementTestSuite) TestFilterBankStatements() {

	// TODO corrigir este teste. Pode ser que não tenha esta conta.
	s.T().Skip("Criar a conta e depois dar um filter.")

	endTime := time.Now().Add(-24 * time.Hour)
	req := &bankly.FilterBankStatementRequest{
		Branch:         "0001",
		Account:        "184039",
		Page:           1,
		PageSize:       1,
		IncludeDetails: true,
		EndDateTime:    &endTime,
		CardProxy:      []string{"123", "456"},
	}

	r, err := s.bank.FilterBankStatements(req)

	s.assert.NoError(err)
	s.assert.NotEmpty(r)
}
