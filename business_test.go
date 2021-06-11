package bankly_test

import (
	"os"
	"testing"
	"time"

	bankly "github.com/contbank/bankly-sdk"
	"github.com/contbank/grok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BusinessTestSuite struct {
	suite.Suite
	assert   *assert.Assertions
	session  *bankly.Session
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

func (s *BusinessTestSuite) TestCreateBusiness_TypeEI_SizeME() {
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeEI, bankly.BusinessSizeME)

	err := s.business.CreateBusiness(businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)
}

func (s *BusinessTestSuite) TestCreateBusiness_TypeMEI_SizeMEI() {
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeMEI, bankly.BusinessSizeMEI)

	err := s.business.CreateBusiness(businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)
}

func (s *BusinessTestSuite) TestCreateBusiness_TypeEIRELI_SizeEPP() {
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeEIRELI, bankly.BusinessSizeEPP)

	err := s.business.CreateBusiness(businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)
}

func (s *BusinessTestSuite) TestCreateBusinessErrorInvalidTypeMEIAndSizeEPP() {
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeMEI, bankly.BusinessSizeEPP)

	err := s.business.CreateBusiness(businessRequest)

	s.assert.Error(err)
	s.assert.Equal(err, bankly.ErrInvalidBusinessSize)
}

func (s *BusinessTestSuite) TestCreateBusinessErrorInvalidTypeMEIAndSizeME() {
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeMEI, bankly.BusinessSizeME)

	err := s.business.CreateBusiness(businessRequest)

	s.assert.Error(err)
	s.assert.Equal(err, bankly.ErrInvalidBusinessSize)
}

func (s *BusinessTestSuite) TestUpdateBusinessName() {

	// TODO O PATCH da Business não está funcionando na Bankly. Aguardando retorno deles para ativar este teste.
	s.T().Skip("O PATCH da Business não está funcionando na Bankly. Aguardando retorno deles para ativar este teste.")

	// create business
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeMEI, bankly.BusinessSizeMEI)
	err := s.business.CreateBusiness(businessRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	// update business
	newBusinessName := "NOVO EMAIL DA EMPRESA VIA UPDATE REQUEST"
	businessUpdateRequest := &bankly.BusinessUpdateRequest{
		BusinessName: newBusinessName,
	}
	err = s.business.UpdateBusiness(businessRequest.Document, *businessUpdateRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	// find updated business
	updatedAccount, err := s.business.FindBusiness(businessRequest.Document)
	s.assert.NoError(err)
	s.assert.NotNil(updatedAccount)
	s.assert.Equal(businessRequest.Document, updatedAccount.Document)
	s.assert.Equal(newBusinessName, updatedAccount.BusinessName)
}

func (s *BusinessTestSuite) TestUpdateBusinessEmailAndBusinessTypeAndBusinessType() {

	// TODO O PATCH da Business não está funcionando na Bankly. Aguardando retorno deles para ativar este teste.
	s.T().Skip("O PATCH da Business não está funcionando na Bankly. Aguardando retorno deles para ativar este teste.")

	// create business
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeMEI, bankly.BusinessSizeMEI)
	err := s.business.CreateBusiness(businessRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	// update business
	newEmail := "novo_email_" + businessRequest.Document + "@contbank.com"
	businessUpdateRequest := &bankly.BusinessUpdateRequest{
		BusinessEmail: newEmail,
		BusinessType:  bankly.BusinessTypeEI,
		BusinessSize:  bankly.BusinessSizeME,
	}
	err = s.business.UpdateBusiness(businessRequest.Document, *businessUpdateRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	// find business updated
	updatedAccount, err := s.business.FindBusiness("59619372000143")
	s.assert.NoError(err)
	s.assert.NotNil(updatedAccount)
	s.assert.Equal(businessRequest.Document, updatedAccount.Document)
	s.assert.Equal(bankly.BusinessTypeEI, updatedAccount.BusinessType)
	s.assert.Equal(bankly.BusinessSizeME, updatedAccount.BusinessSize)
	s.assert.Equal(newEmail, updatedAccount.BusinessEmail)

}

func (s *BusinessTestSuite) TestCreateBusinessAccount() {
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeEI, bankly.BusinessSizeME)

	err := s.business.CreateBusiness(businessRequest)
	s.assert.NoError(err)
	s.assert.Nil(err)

	businessAccountRequest := bankly.BusinessAccountRequest{
		Document:    businessRequest.Document,
		AccountType: bankly.PaymentAccount,
	}

	time.Sleep(time.Second)

	account, err := s.business.CreateBusinessAccount(businessAccountRequest)
	s.assert.NoError(err)
	s.assert.NotNil(account)
}

func (s *BusinessTestSuite) TestCreateBusinessAccountErrorDoesntHaveAnApprovedRegistration() {
	businessAccountRequest := bankly.BusinessAccountRequest{
		Document:    bankly.GeneratorCNPJ(),
		AccountType: bankly.PaymentAccount,
	}
	account, err := s.business.CreateBusinessAccount(businessAccountRequest)

	s.assert.Error(err)
	s.assert.Equal(err, bankly.ErrAccountHolderNotExists)
	s.assert.Nil(account)
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

func (s *BusinessTestSuite) TestBusinessName_TypeMEI() {
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeMEI, bankly.BusinessSizeMEI)
	businessName := businessRequest.LegalRepresentative.RegisterName + " " + businessRequest.LegalRepresentative.Document
	err := s.business.CreateBusiness(businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)

	time.Sleep(time.Millisecond)

	account, err := s.business.FindBusiness(businessRequest.Document)

	s.assert.NoError(err)
	s.assert.NotNil(account)
	s.assert.Equal(businessName, account.BusinessName)
}

func (s *BusinessTestSuite) TestBusinessName_TypeEI() {
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeEI, bankly.BusinessSizeEPP)
	err := s.business.CreateBusiness(businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)

	time.Sleep(time.Millisecond)

	account, err := s.business.FindBusiness(businessRequest.Document)

	s.assert.NoError(err)
	s.assert.NotNil(account)
	s.assert.Equal(businessRequest.BusinessName, account.BusinessName)
}

func (s *BusinessTestSuite) TestBusinessName_TypeEIRELI() {
	businessRequest := createBusinessRequest(grok.GeneratorCNPJ(), bankly.BusinessTypeEIRELI, bankly.BusinessSizeEPP)
	err := s.business.CreateBusiness(businessRequest)

	s.assert.NoError(err)
	s.assert.Nil(err)

	time.Sleep(time.Millisecond)

	account, err := s.business.FindBusiness(businessRequest.Document)

	s.assert.NoError(err)
	s.assert.NotNil(account)
	s.assert.Equal(businessRequest.BusinessName, account.BusinessName)
}

func createBusinessRequest(document string, businessType bankly.BusinessType, businessSize bankly.BusinessSize) bankly.BusinessRequest {

	randomName := bankly.RandStringBytes(10)
	email := "email_de_teste_" + randomName + "@contbank.com"

	return bankly.BusinessRequest{
		Document:            bankly.OnlyDigits(document),
		BusinessName:        "Nome da Empresa " + randomName,
		TradingName:         "Nome Fantasia " + randomName,
		BusinessEmail:       email,
		BusinessType:        businessType,
		BusinessSize:        businessSize,
		BusinessAddress:     createAddress(),
		LegalRepresentative: createLegalRepresentative(),
	}
}

func createLegalRepresentative() *bankly.LegalRepresentative {
	randSurname := randStringBytes(10)
	return &bankly.LegalRepresentative{
		Document:     grok.GeneratorCPF(),
		RegisterName: "Nome do Representante Legal " + randSurname,
		Phone: &bankly.Phone{
			CountryCode: "55",
			Number:      bankly.GeneratorCellphone(),
		},
		Address:    createAddress(),
		BirthDate:  time.Date(1990, time.June, 29, 0, 0, 0, 0, time.UTC),
		MotherName: "Nome da Mãe do Representante Legal " + randSurname,
		Email:      "email_legal_representative_" + randSurname + "@contbank.com",
	}
}

func createAddress() *bankly.Address {
	return &bankly.Address{
		ZipCode:        "05410900",
		City:           "São Paulo",
		AddressLine:    "Rua Qualquer XYZ",
		BuildingNumber: "2000",
		Neighborhood:   "Pinheiros",
		State:          "SP",
		Country:        "BR",
	}
}
