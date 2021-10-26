package bankly_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/contbank/bankly-sdk"
	"github.com/contbank/grok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FreshDeskTestSuite struct {
	suite.Suite
	assert   *assert.Assertions
	httpMock *grok.HTTPClientMock
	service  *bankly.ServiceDesk
}

func TestFreshDeskTestSuite(t *testing.T) {
	suite.Run(t, new(FreshDeskTestSuite))
}

func (s *FreshDeskTestSuite) SetupSuite() {
	s.assert = assert.New(s.T())

	session, err := bankly.NewServiceDeskSession(bankly.ServiceDeskConfig{
		APIEndpoint: bankly.String(""),
		APIKey:      bankly.String(""),
	})

	s.assert.NoError(err)

	s.httpMock = grok.NewHTTPClientMock()
	s.service = bankly.NewServiceDesk(s.httpMock.Client(), session)
}

func (s *FreshDeskTestSuite) TestCreateTicket() {
	ctx := context.Background()

	s.setupCreateTicket(randomNumber())

	req := &bankly.CreateTicketRequest{
		GroupID:     bankly.DefaultGroupID,
		ProductID:   bankly.DefaultProductID,
		Type:        bankly.DefaultTicketType,
		Priority:    bankly.DefaultPriority,
		Status:      bankly.OpenStatus,
		Subject:     "API TEST",
		Email:       "email@contbank.com",
		CompanyKey:  bankly.CompanyKey,
		Request:     bankly.DefaultTicketRequest,
		Request2:    bankly.SAAccountType,
		CNPJ:        "48374325000160",
		RazaoSocial: "NOME EMPRESA TEST",
		Cellphone:   11911112222,
		Attachments: []string{"go.sum"},
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

	res, err := s.service.FilterTickets(ctx, &bankly.FilterTicketsRequest{
		Status: status,
	})

	s.assert.NoError(err)
	s.assert.Equal(1, res.Total)
	s.assert.Len(res.Results, 1)
}

func (s *FreshDeskTestSuite) setupCreateTicket(id int) {
	bodyResponse, _ := json.Marshal(&bankly.CreateTicketResponse{ID: id})

	response := &grok.MockedResponseResult{
		Status: http.StatusCreated,
		Body:   string(bodyResponse),
	}

	s.httpMock.AddMock("/api/v2/tickets", response)
}

func (s *FreshDeskTestSuite) setupGetTicket(id, status int) {
	bodyResponse, _ := json.Marshal(&bankly.GetTicketResponse{
		ID:     id,
		Status: status,
	})

	response := &grok.MockedResponseResult{
		Status: http.StatusOK,
		Body:   string(bodyResponse),
	}

	s.httpMock.AddMock(fmt.Sprintf("/api/v2/tickets/%d", id), response)
}

func (s *FreshDeskTestSuite) setupFilterTickets(id, status int) {
	bodyResponse, _ := json.Marshal(&bankly.FilterTicketsResponse{
		Total:   1,
		Results: []*bankly.GetTicketResponse{{ID: id, Status: status}},
	})

	response := &grok.MockedResponseResult{
		Status: http.StatusOK,
		Body:   string(bodyResponse),
	}

	s.httpMock.AddMock("/api/v2/search/tickets?query=%22status%3A2%22", response)
}

func randomNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(1000)
}
