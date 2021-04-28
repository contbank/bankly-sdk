package bankly_test

import (
	"github.com/contbank/grok"
	"os"
	"testing"
	"time"

	"github.com/contbank/bankly-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BusinessTestSuite struct {
	suite.Suite
	assert    *assert.Assertions
	session   *bankly.Session
	business *bankly.Business
}

func TestBusinessTestSuite(t *testing.T) {
	suite.Run(t, new(BusinessTestSuite))
}

func (s *BusinessTestSuite) SetupTest() {
	s.assert = assert.New(s.T())

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(os.Getenv("BANKLY_CLIENT_ID")),
		ClientSecret: bankly.String(os.Getenv("BANKLY_CLIENT_SECRET")),
	})

	s.assert.NoError(err)

	s.session = session
	s.business = bankly.NewBusiness(*s.session)
}

func (s *BusinessTestSuite) TestCreateBusiness() {
	businessRequest := s.createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeEI, bankly.BusinessSizeME)

	err := s.business.CreateBusiness(businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)
}

func (s *BusinessTestSuite) TestCreateBusinessErrorInvalidTypeMEIAndSizeEPP() {
	businessRequest := s.createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeMEI, bankly.BusinessSizeEPP)

	err := s.business.CreateBusiness(businessRequest)

	s.assert.Error(err)
	s.assert.EqualError(err, "Invalid business size to personal business")
}

func (s *BusinessTestSuite) TestCreateBusinessErrorInvalidTypeMEIAndSizeME() {
	businessRequest := s.createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeMEI, bankly.BusinessSizeME)

	err := s.business.CreateBusiness(businessRequest)

	s.assert.Error(err)
	s.assert.EqualError(err, "Invalid business size to personal business")
}

func (s *BusinessTestSuite) TestCreateBusinessAccount() {
	businessRequest := s.createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeEI, bankly.BusinessSizeME)

	err := s.business.CreateBusiness(businessRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	businessAccountRequest := bankly.BusinessAccountRequest {
		Document: businessRequest.Document,
		AccountType: bankly.PaymentAccount,
	}

	account, err := s.business.CreateBusinessAccount(businessAccountRequest)
	s.assert.NoError(err)
	s.assert.NotNil(account)
}

func (s *BusinessTestSuite) TestCreateBusinessAccountErrorDoesntHaveAnApprovedRegistration() {
	businessAccountRequest := bankly.BusinessAccountRequest {
		Document: grok.GeneratorCNPJ(),
		AccountType: bankly.PaymentAccount,
	}

	account, err := s.business.CreateBusinessAccount(businessAccountRequest)

	s.assert.Error(err)
	s.assert.EqualError(err, "Account holder does not exist or does not have an approved registration yet")
	s.assert.Nil(account)
}

func (s *BusinessTestSuite) createBusinessRequest(document string,
	businessType bankly.BusinessType, businessSize bankly.BusinessSize) bankly.BusinessRequest {

	randomName := randStringBytes(10)
	email := "email_de_teste_" + randomName + "@contbank.com"

	return bankly.BusinessRequest {
		Document: grok.OnlyDigits(document),
		BusinessName: "Nome da Empresa " + randomName,
		TradingName: "Nome Fantasia " + randomName,
		BusinessEmail: email,
		BusinessType: businessType,
		BusinessSize: businessSize,
		BusinessAddress: &bankly.Address{
			ZipCode:        "05410900",
			City:           "São Paulo",
			AddressLine:    "Rua Qualquer XYZ",
			BuildingNumber: "2000",
			Neighborhood:   "Pinheiros",
			State:          "SP",
			Country:        "BR",
		},
		LegalRepresentative: s.createLegalRepresentative(),
	}
}

func (s *BusinessTestSuite) createLegalRepresentative() bankly.LegalRepresentative {
	randSurname := randStringBytes(10)

	return bankly.LegalRepresentative {
		Documment: grok.GeneratorCPF(),
		RegisterName: "Nome do Representante Legal " + randSurname,
		Phone: &bankly.Phone {
			CountryCode: "55",
			Number:      grok.GeneratorCellphone(),
		},
		Address: &bankly.Address {
			ZipCode:        "05410900",
			City:           "São Paulo",
			AddressLine:    "Rua Qualquer XYZ",
			BuildingNumber: "2000",
			Neighborhood:   "Pinheiros",
			State:          "SP",
			Country:        "BR",
		},
		BirthDate:    time.Date(1990, time.June, 29, 0, 0, 0, 0, time.UTC),
		MotherName:   "Nome da Mãe do Representante Legal " + randSurname,
		Email:        "email_legal_representative_" + randSurname + "@contbank.com",
	}
}

func (s *BusinessTestSuite) TestFindBusiness() {
	account, err := s.business.FindBusiness("59619372000143")

	s.assert.NoError(err)
	s.assert.NotNil(account)
}

func (s *BusinessTestSuite) TestFindBusinessAccounts() {
	account, err := s.business.FindBusinessAccounts("59619372000143")

	s.assert.NoError(err)
	s.assert.NotNil(account)
}