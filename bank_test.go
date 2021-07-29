package bankly_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/contbank/bankly-sdk"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BankTestSuite struct {
	suite.Suite
	assert  *assert.Assertions
	ctx     context.Context
	session *bankly.Session
	bank    *bankly.Bank
}

func TestBankTestSuite(t *testing.T) {
	suite.Run(t, new(BankTestSuite))
}

func (s *BankTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(os.Getenv("BANKLY_CLIENT_ID")),
		ClientSecret: bankly.String(os.Getenv("BANKLY_CLIENT_SECRET")),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.bank = bankly.NewBank(httpClient, *s.session)
}

func (s *BankTestSuite) TestList() {
	req := &bankly.FilterBankListRequest{}

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
