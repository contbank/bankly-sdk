package bankly_test

import (
	"context"
	models "github.com/contbank/bankly-sdk/pkg/models"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	incomeReport "github.com/contbank/bankly-sdk/pkg/services/income-report"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IncomeReportTestSuite struct {
	suite.Suite
	assert           *assert.Assertions
	ctx              context.Context
	session          *bankly.Session
	bankIncomeReport *incomeReport.IncomeReport
}

func TestIncomeReportTestSuite(t *testing.T) {
	suite.Run(t, new(IncomeReportTestSuite))
}

func (s *IncomeReportTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session, err := bankly.NewSession(bankly.Config{
		ClientID:     bankly.String(*utils.GetEnvBanklyClientID()),
		ClientSecret: bankly.String(*utils.GetEnvBanklyClientSecret()),
		Scopes:       bankly.String("income.report.read"),
	})

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session

	newHttpClient := utils.BanklyHttpClient{
		Session:        *session,
		HttpClient:     httpClient,
		Authentication: bankly.NewAuthentication(httpClient, *session),
	}

	s.bankIncomeReport = incomeReport.NewIncomeReport(newHttpClient)
}

// TestIncomeReport_SUCCESS ...
func (s *IncomeReportTestSuite) TestIncomeReport_SUCCESS() {
	report, err := s.bankIncomeReport.GetIncomeReport(s.ctx,
		&models.IncomeReportRequest{
			Account: "184152",
			Year:    "2021",
		})

	s.assert.NoError(err)
	s.assert.NotNil(report)
	s.assert.NotEmpty(report.FileName)
	s.assert.NotEmpty(report.IncomeFile)
}

// TestIncomeReport_INVALID_ACCOUNT_NUMBER ...
func (s *IncomeReportTestSuite) TestIncomeReport_INVALID_ACCOUNT_NUMBER() {
	report, err := s.bankIncomeReport.GetIncomeReport(s.ctx,
		&models.IncomeReportRequest{
			Account: "100000",
			Year:    "2021",
		})

	s.assert.Error(err)
	s.assert.Nil(report)
}
