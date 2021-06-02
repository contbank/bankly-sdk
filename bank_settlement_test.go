package bankly_test

import (
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

	s.session = session
	s.bank = bankly.NewBankStatement(*s.session)
	s.boletos = bankly.NewBoletos(*s.session)
}

func (s *BankStatementTestSuite) TestFilterBankStatements() {
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

func (s *BankStatementTestSuite) createBoleto(document string, account *bankly.Account) *bankly.BoletoRequest {
	return &bankly.BoletoRequest{
		Alias:       bankly.String(bankly.RandStringBytes(12)),
		Document:    document,
		Amount:      1547.55,
		DueDate:     time.Now().Add(48 * time.Hour),
		EmissionFee: false,
		Type:        bankly.Levy,
		Account:     account,
		Payer:       s.createPayer(),
	}
}

func (s *BankStatementTestSuite) createPayer() *bankly.Payer {
	return &bankly.Payer{
		Name:      bankly.RandStringBytes(9),
		TradeName: bankly.RandStringBytes(15),
		Document:  bankly.GeneratorCPF(),
		Address:   s.createAddress(),
	}
}

func (s *BankStatementTestSuite) createAddress() *bankly.Address {
	return &bankly.Address{
		ZipCode:        "03503030",
		City:           "São Paulo",
		AddressLine:    "Rua Fulano de Tal",
		BuildingNumber: "1000",
		Neighborhood:   "Chácara Califórnia",
		State:          "SP",
		Country:        "BR",
	}
}
