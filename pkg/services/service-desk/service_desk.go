package service_desk

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/contbank/bankly-sdk/pkg/errors"
	"github.com/contbank/bankly-sdk/pkg/models"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	// DefaultPriority ...
	DefaultPriority = 2
	// DefaultProductID ...
	DefaultProductID = 47000002195
	// DefaultGroupID ...
	DefaultGroupID = 47000657625
	// DefaultTicketType ...
	DefaultTicketType = "Bankly Service Desk - API"
	// DefaultTicketRequest ...
	DefaultTicketRequest = "Criacao de Conta PJ - Bankly"
	// CompanyKey ...
	CompanyKey = "CONTBANK"

	// SAAccountType ...
	SAAccountType = "Tipo - S.A"
	// LTDAAccountType ...
	LTDAAccountType = "Tipo - LTDA"

	// ApprovedReason ...
	ApprovedReason = "APROVADO"
	// InAnalysisReason ...
	InAnalysisReason = "EM ANALISE"
	// MissingInfoReason ...
	MissingInfoReason = "FALTA_INFO"
	// ReprovedReason ...
	ReprovedReason = "REPROVADO"

	// OpenStatus ...
	OpenStatus = 2
	// PendingStatus ...
	PendingStatus = 3
	// ResolvedStatus ...
	ResolvedStatus = 4
	// ClosedStatus ...
	ClosedStatus = 5
	// AnsweredStatus ...
	AnsweredStatus = 6
)

// ServiceDesk ...
type ServiceDesk struct {
	httpClient *http.Client
	session    *ServiceDeskSession
}

// NewServiceDesk ...
func NewServiceDesk(httpClient *http.Client, session *ServiceDeskSession) *ServiceDesk {
	return &ServiceDesk{
		httpClient: httpClient,
		session:    session,
	}
}

// CreateTicket ...
func (s *ServiceDesk) CreateTicket(ctx context.Context,
	model *models.CreateTicketRequest) (*models.CreateTicketResponse, error) {

	endpoint, err := url.Parse(s.session.APIEndpoint)
	if err != nil {
		logrus.WithError(err).Error("error while parsing url")
		return nil, err
	}

	endpoint.Path = path.Join(endpoint.Path, "api/v2")
	endpoint.Path = path.Join(endpoint.Path, "tickets")

	formData := map[string][]string{
		"group_id":                       {strconv.Itoa(model.GroupID)},
		"product_id":                     {strconv.Itoa(model.ProductID)},
		"type":                           {model.Type},
		"description":                    {model.Description},
		"subject":                        {model.Subject},
		"email":                          {model.Email},
		"cc_emails[]":                    model.CopyEmails,
		"priority":                       {strconv.Itoa(model.Priority)},
		"status":                         {strconv.Itoa(model.Status)},
		"custom_fields[cf_solicitacao]":  {model.Request},
		"custom_fields[cf_solicitacao2]": {model.Request2},
		"custom_fields[cf_companykey]":   {model.CompanyKey},
		"custom_fields[cf_razao_social]": {model.RazaoSocial},
		"custom_fields[cf_cnpj_empresa]": {model.CNPJ},
		"custom_fields[cf_cpf]":          {"-"},
		"custom_fields[cf_celular]":      {strconv.Itoa(model.Cellphone)},
		"custom_fields[cf_comprastransaesestornos_4_ltimos_dgitos_do_carto]": {strconv.Itoa(model.Last4Digits)},
	}

	body, contentType, err := generateFormDataBody(formData, model.Attachments)
	if err != nil {
		logrus.WithError(err).Error("error while creating body data")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint.String(), body)
	if err != nil {
		logrus.WithError(err).Error("error while creating request")
		return nil, err
	}

	req.Header.Set("Authorization", getAuthorization(s.session.APIKey))
	req.Header.Set("Content-Type", contentType)

	res, err := s.httpClient.Do(req)
	if err != nil {
		logrus.WithError(err).Error("error while performing the request")
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		logrus.Errorf("invalid status code: %d", res.StatusCode)
		return nil, errors.ErrDefaultFreshDesk
	}

	var result *models.CreateTicketResponse

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		logrus.WithError(err).Error("error while reading the response")
		return nil, err
	}

	return result, nil
}

// GetTicket ...
func (s *ServiceDesk) GetTicket(ctx context.Context, id int) (*models.GetTicketResponse, error) {

	endpoint, err := url.Parse(s.session.APIEndpoint)
	if err != nil {
		logrus.WithError(err).Error("error while parsing url")
		return nil, err
	}

	endpoint.Path = path.Join(endpoint.Path, "api/v2")
	endpoint.Path = path.Join(endpoint.Path, "tickets")
	endpoint.Path = path.Join(endpoint.Path, strconv.Itoa(id))

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint.String(), nil)
	if err != nil {
		logrus.WithError(err).Error("error while creating request")
		return nil, err
	}

	req.Header.Set("Authorization", getAuthorization(s.session.APIKey))

	res, err := s.httpClient.Do(req)
	if err != nil {
		logrus.WithError(err).Error("error while performing the request")
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, errors.ErrFreshDeskTicketNotFound
	}

	if res.StatusCode != http.StatusOK {
		logrus.Errorf("invalid status code: %d", res.StatusCode)
		return nil, errors.ErrDefaultFreshDesk
	}

	var result *models.GetTicketResponse

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		logrus.WithError(err).Error("error while reading the response")
		return nil, err
	}

	return result, nil
}

// FilterTickets ...
func (s *ServiceDesk) FilterTickets(ctx context.Context,
	model *models.FilterTicketsRequest) (*models.FilterTicketsResponse, error) {

	endpoint, err := url.Parse(s.session.APIEndpoint)
	if err != nil {
		logrus.WithError(err).Error("error while parsing url")
		return nil, err
	}

	endpoint.Path = path.Join(endpoint.Path, "api/v2")
	endpoint.Path = path.Join(endpoint.Path, "search/tickets")

	q := endpoint.Query()
	q.Set("query", fmt.Sprintf("\"status:%d\"", model.Status))
	endpoint.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint.String(), nil)
	if err != nil {
		logrus.WithError(err).Error("error while creating request")
		return nil, err
	}

	req.Header.Set("Authorization", getAuthorization(s.session.APIKey))

	res, err := s.httpClient.Do(req)
	if err != nil {
		logrus.WithError(err).Error("error while performing the request")
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logrus.Errorf("invalid status code: %d", res.StatusCode)
		return nil, errors.ErrDefaultFreshDesk
	}

	var result *models.FilterTicketsResponse

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		logrus.WithError(err).Error("error while reading the response")
		return nil, err
	}

	return result, nil
}

func getAuthorization(apiKey string) string {
	key := fmt.Sprintf("%s:X", apiKey)
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(key)))
}

func generateFormDataBody(values map[string][]string, attachmentPaths []string) (*bytes.Buffer, string, error) {
	var b bytes.Buffer

	w := multipart.NewWriter(&b)
	defer w.Close()

	for key, items := range values {
		for _, val := range items {
			var (
				fw  io.Writer
				err error
			)

			if fw, err = w.CreateFormField(key); err != nil {
				return nil, "", err
			}

			if _, err = fw.Write([]byte(val)); err != nil {
				return nil, "", err
			}
		}
	}

	for _, p := range attachmentPaths {
		var (
			fw         io.Writer
			attachment *os.File
			err        error
		)

		if attachment, err = os.Open(p); err != nil {
			return nil, "", err
		}

		defer attachment.Close()

		if fw, err = w.CreateFormFile("attachments[]", attachment.Name()); err != nil {
			return nil, "", err
		}

		if _, err = io.Copy(fw, attachment); err != nil {
			return nil, "", err
		}

	}

	return &b, w.FormDataContentType(), nil
}