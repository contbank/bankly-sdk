package bank_test

import (
	"context"
	"github.com/contbank/bankly-sdk/pkg/models"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	"github.com/contbank/bankly-sdk/pkg/services/bank"
	"github.com/contbank/bankly-sdk/pkg/utils"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BankTestSuite struct {
	suite.Suite
	assert  *assert.Assertions
	ctx     context.Context
	session *authentication.Session
	bank    *bank.Bank
}

func TestBankTestSuite(t *testing.T) {
	suite.Run(t, new(BankTestSuite))
}

func (s *BankTestSuite) SetupTest() {
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
	s.bank = bank.NewBank(httpClient, *s.session)
}

func (s *BankTestSuite) TestList() {
	req := &models.FilterBankListRequest{}

	banks, err := s.bank.List(s.ctx, req)

	s.assert.NoError(err)
	s.assert.NotEmpty(banks)
}

/*
func (s *BankTestSuite) TestGetByID() {
	bank, err := s.bank.GetByID("001")

	s.assert.NoError(err)
	s.assert.NotNil(bank)
}
*/
////this test is only returning the last id provided in the list
// func (s *BankTestSuite) TestListFilterIDs() {
// 	filter := &bankly.FilterBankListRequest{}
// 	banks, err := s.bank.List(filter)

// 	s.assert.NoError(err)
// 	s.assert.NotEmpty(banks)

// 	ids := make([]string, 2)

// 	for i, b := range banks[:2] {
// 		ids[i] = b.Code
// 	}

// 	filterID := &bankly.FilterBankListRequest{
// 		IDs: ids,
// 	}
// 	filterBanks, err := s.bank.List(filterID)

// 	s.assert.NoError(err)
// 	s.assert.Len(filterBanks, len(ids))
// }
