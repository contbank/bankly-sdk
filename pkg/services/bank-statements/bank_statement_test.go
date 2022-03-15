package bank_statements_test

import (
	"context"
	"github.com/contbank/bankly-sdk/pkg/models"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	"github.com/contbank/bankly-sdk/pkg/services/bank-statements"
	"github.com/contbank/bankly-sdk/pkg/utils"
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
	session *authentication.Session
	bank    *bank_statements.BankStatement
}

func TestBankStatementTestSuite(t *testing.T) {
	suite.Run(t, new(BankStatementTestSuite))
}

func (s *BankStatementTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session, err := authentication.NewSession(authentication.Config{
		ClientID :     utils.String(*utils.GetEnvBanklyClientID()),
		ClientSecret : utils.String(*utils.GetEnvBanklyClientSecret()),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.bank = bank_statements.NewBankStatement(httpClient, *s.session)
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