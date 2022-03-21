package bankly_test

import (
	"context"
	"encoding/json"
	"fmt"
	models "github.com/contbank/bankly-sdk/pkg/models"
	serviceDesk "github.com/contbank/bankly-sdk/pkg/services/service-desk"
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/contbank/grok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FreshDeskTestSuite struct {
	suite.Suite
	assert   *assert.Assertions
	httpMock *grok.HTTPClientMock
	service  *serviceDesk.ServiceDesk
}

func TestFreshDeskTestSuite(t *testing.T) {
	suite.Run(t, new(FreshDeskTestSuite))
}

func (s *FreshDeskTestSuite) SetupSuite() {
	s.assert = assert.New(s.T())

	session, err := serviceDesk.NewServiceDeskSession(serviceDesk.ServiceDeskConfig{
		APIEndpoint: utils.String(""),
		APIKey:      utils.String(""),
	})

	s.assert.NoError(err)

	s.httpMock = grok.NewHTTPClientMock()
	s.service = serviceDesk.NewServiceDesk(s.httpMock.Client(), session)
}

func (s *FreshDeskTestSuite) TestCreateTicket() {
	ctx := context.Background()

	s.setupCreateTicket(randomNumber())

	req := &models.CreateTicketRequest{
		GroupID:     serviceDesk.DefaultGroupID,
		ProductID:   serviceDesk.DefaultProductID,
		Type:        serviceDesk.DefaultTicketType,
		Priority:    serviceDesk.DefaultPriority,
		Status:      serviceDesk.OpenStatus,
		Subject:     "API TEST",
		Email:       "email@contbank.com",
		CompanyKey:  serviceDesk.CompanyKey,
		Request:     serviceDesk.DefaultTicketRequest,
		Request2:    serviceDesk.SAAccountType,
		CNPJ:        "48374325000160",
		RazaoSocial: "NOME EMPRESA TEST",
		Cellphone:   11911112222,
		Attachments: []string{"../../test_files/contbank.png"},
	}

	res, err := s.service.CreateTicket(ctx, req)

	s.assert.NoError(err)
	s.assert.NotEmpty(res.ID)
}

func (s *FreshDeskTestSuite) TestGetTicket() {
	ctx := context.Background()
	id := randomNumber()
	status := 2

	s.setupGetTicket(id, status)

	res, err := s.service.GetTicket(ctx, id)

	s.assert.NoError(err)
	s.assert.Equal(id, res.ID)
	s.assert.Equal(status, res.Status)
}

func (s *FreshDeskTestSuite) TestFilterTickets() {
	ctx := context.Background()
	id := randomNumber()
	status := 4

	s.setupFilterTickets(id, status)

	res, err := s.service.FilterTickets(ctx, &models.FilterTicketsRequest{
		Status: status,
	})

	s.assert.NoError(err)
	s.assert.Equal(1, res.Total)
	s.assert.Len(res.Results, 1)
}

func (s *FreshDeskTestSuite) setupCreateTicket(id int) {
	bodyResponse, _ := json.Marshal(&models.CreateTicketResponse{ID: id})

	response := &grok.MockedResponseResult{
		Status: http.StatusCreated,
		Body:   string(bodyResponse),
	}

	s.httpMock.AddMock("api/v2/tickets", response)
}

func (s *FreshDeskTestSuite) setupGetTicket(id, status int) {
	bodyResponse, _ := json.Marshal(&models.GetTicketResponse{
		ID:     id,
		Status: status,
	})

	response := &grok.MockedResponseResult{
		Status: http.StatusOK,
		Body:   string(bodyResponse),
	}

	s.httpMock.AddMock(fmt.Sprintf("api/v2/tickets/%d", id), response)
}

func (s *FreshDeskTestSuite) setupFilterTickets(id, status int) {
	bodyResponse, _ := json.Marshal(&models.FilterTicketsResponse{
		Total:   1,
		Results: []*models.GetTicketResponse{{ID: id, Status: status}},
	})

	response := &grok.MockedResponseResult{
		Status: http.StatusOK,
		Body:   string(bodyResponse),
	}

	s.httpMock.AddMock("api/v2/search/tickets?query=%22status%3A4%22", response)
}

func randomNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(1000)
}
