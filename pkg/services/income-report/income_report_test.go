package income_report_test

import (
	"context"
	"github.com/contbank/bankly-sdk/pkg/models"
	"github.com/contbank/bankly-sdk/pkg/services/authentication"
	"github.com/contbank/bankly-sdk/pkg/services/income-report"
	"github.com/contbank/bankly-sdk/pkg/utils"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IncomeReportTestSuite struct {
	suite.Suite
	assert  		 *assert.Assertions
	ctx     		 context.Context
	session 		 *authentication.Session
	bankIncomeReport *income_report.IncomeReport
}

func TestIncomeReportTestSuite(t *testing.T) {
	suite.Run(t, new(IncomeReportTestSuite))
}

func (s *IncomeReportTestSuite) SetupTest() {
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

	newHttpClient := utils.BanklyHttpClient{
		Session:        *session,
		HttpClient:     httpClient,
		Authentication: authentication.NewAuthentication(httpClient, *session),
	}

	s.bankIncomeReport = income_report.NewIncomeReport(newHttpClient)
}

// TestIncomeReport_SUCCESS ...
func (s *IncomeReportTestSuite) TestIncomeReport_SUCCESS() {
	report, err := s.bankIncomeReport.GetIncomeReport(s.ctx,
		&models.IncomeReportRequest{
			Account: "184152",
			Year: "2021",
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
			Year: "2021",
		})

	s.assert.Error(err)
	s.assert.Nil(report)
}