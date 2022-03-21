package bankly_test

import (
	"context"
	models "github.com/contbank/bankly-sdk/pkg/models"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	bankStatements "github.com/contbank/bankly-sdk/pkg/services/bank-statements"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BankStatementTestSuite struct {
	suite.Suite
	assert  *assert.Assertions
	ctx     context.Context
	session *bankly.Session
	bank    *bankStatements.BankStatement
}

func TestBankStatementTestSuite(t *testing.T) {
	suite.Run(t, new(BankStatementTestSuite))
}

func (s *BankStatementTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(*utils.GetEnvBanklyClientID()),
		ClientSecret: bankly.String(*utils.GetEnvBanklyClientSecret()),
		Scopes:       bankly.String("events.read"),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.bank = bankStatements.NewBankStatement(httpClient, *s.session)
}

func (s *BankStatementTestSuite) TestFilterBankStatements() {
	// TODO corrigir este teste. Pode ser que não tenha esta conta.
	s.T().Skip("Criar a conta e depois dar um filter.")

	endTime := time.Now().Add(-24 * time.Hour)
	req := &models.FilterBankStatementRequest{
		Branch:         "0001",
		Account:        "184039",
		Page:           1,
		PageSize:       1,
		IncludeDetails: true,
		EndDateTime:    &endTime,
		CardProxy:      []string{"123", "456"},
	}

	r, err := s.bank.FilterBankStatements(s.ctx, req)

	s.assert.NoError(err)
	s.assert.NotEmpty(r)
}

func (s *BankStatementTestSuite) TestFilterBankStatements_ActiveStatus() {
	// TODO corrigir este teste. Pode ser que não tenha esta conta.
	s.T().Skip("Criar a conta e depois dar um filter.")

	endTime := time.Now().Add(-24 * time.Hour)
	req := &models.FilterBankStatementRequest{
		Branch:         "0001",
		Account:        "184039",
		Page:           1,
		PageSize:       1,
		IncludeDetails: true,
		EndDateTime:    &endTime,
		CardProxy:      []string{"123", "456"},
	}

	status := models.Active
	req.Status = &status

	r, err := s.bank.FilterBankStatements(s.ctx, req)

	s.assert.NoError(err)
	s.assert.NotEmpty(r)
}

func (s *BankStatementTestSuite) TestFilterBankStatements_InvalidPageSizeError() {
	endTime := time.Now().Add(-24 * time.Hour)
	req := &models.FilterBankStatementRequest{
		Branch:         "0001",
		Account:        "184039",
		Page:           1,
		PageSize:       500,
		IncludeDetails: true,
		EndDateTime:    &endTime,
		CardProxy:      []string{"123", "456"},
	}

	r, err := s.bank.FilterBankStatements(s.ctx, req)

	s.assert.Error(err)
	s.assert.Nil(r)
}
