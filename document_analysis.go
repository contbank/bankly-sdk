package bankly

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

// DocumentAnalysis ...
type DocumentAnalysis struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

// NewDocumentAnalysis ...
func NewDocumentAnalysis(httpClient *http.Client, session Session) *DocumentAnalysis {
	return &DocumentAnalysis{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

// SendDocumentAnalysis ...
func (c *DocumentAnalysis) SendDocumentAnalysis(ctx context.Context, request DocumentAnalysisRequest) (*DocumentAnalysisResponse, error) {
	err := grok.Validator.Struct(request)
	if err != nil {
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := c.getDocumentAnalysisAPIEndpoint(request.Document, nil, nil, request.IsCorporationBusiness)
	if err != nil {
		return nil, err
	}

	payload, writer, err := createSendImagePayload(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", *endpoint, payload)
	if err != nil {
		logrus.
			WithError(err).
			Error("error new request")
		return nil, err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		logrus.
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("api-version", c.session.APIVersion)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.
			WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusAccepted {
		var bodyResp *DocumentAnalysisRequestedResponse

		err = json.Unmarshal(respBody, &bodyResp)
		if err != nil {
			logrus.
				WithError(err).
				Error("error unmarshal")
			return nil, err
		}

		response := &DocumentAnalysisResponse{
			DocumentNumber: request.Document,
			DocumentType:   string(request.DocumentType),
			DocumentSide:   string(request.DocumentSide),
			Token:          bodyResp.Token,
		}
		return response, nil
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.
			WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return nil, FindError(errModel.Code, errModel.Messages...)
	}

	return nil, ErrSendDocumentAnalysis
}

// FindDocumentAnalysis ...
func (c *DocumentAnalysis) FindDocumentAnalysis(ctx context.Context, documentNumber string,
	documentAnalysisToken string) (*DocumentAnalysisResponse, error) {
	
	resultLevel := ResultLevelDetailed
	endpoint, err := c.getDocumentAnalysisAPIEndpoint(documentNumber, &resultLevel, &documentAnalysisToken, nil)
	if err != nil {
		return nil, err
	} else if endpoint == nil {
		return nil, ErrInvalidAPIEndpoint
	}

	req, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		logrus.WithError(err).Error("error new request")
		return nil, err
	}

	token, err := c.authentication.Token(ctx)
	if err != nil {
		logrus.WithError(err).Error("error token")
		return nil, err
	}

	req = setRequestHeader(req, token, c.session.APIVersion, nil)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithError(err).Error("error request")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response []*BanklyDocumentAnalysisResponse

		err = json.Unmarshal(respBody, &response)
		if err != nil {
			logrus.WithError(err).Error("error unmarshal")
			return nil, err
		}

		return ParseDocumentAnalysisResponse(documentNumber, response[0]), nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)
	if err != nil {
		logrus.WithError(err).Error("error unmarshal")
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return nil, FindError(errModel.Code, errModel.Messages...)
	}

	return nil, ErrGetDocumentAnalysis
}

// getDocumentAnalysisAPIEndpoint ...
func (c *DocumentAnalysis) getDocumentAnalysisAPIEndpoint(document string, resultLevel *ResultLevel,
	documentAnalysisToken *string, isCorporationBusiness *bool) (*string, error) {

	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.
			WithError(err).
			Error("error api endpoint")
		return nil, err
	}
	u.Path = path.Join(u.Path, DocumentAnalysisPath)
	u.Path = path.Join(u.Path, grok.OnlyDigits(document))
	if isCorporationBusiness != nil && *isCorporationBusiness == true {
		u.Path = path.Join(u.Path, CorporationBusinessPath)
	}
	if documentAnalysisToken != nil {
		q := u.Query()
		q.Set("token", *documentAnalysisToken)
		u.RawQuery = q.Encode()
	}
	if resultLevel != nil {
		q := u.Query()
		q.Set("resultLevel", string(*resultLevel))
		u.RawQuery = q.Encode()
	}
	endpoint := u.String()
	return &endpoint, nil
}

// createSendImagePayload ...
func createSendImagePayload(request DocumentAnalysisRequest) (*bytes.Buffer, *multipart.Writer, error) {
	payload := &bytes.Buffer{}

	writer := multipart.NewWriter(payload)
	file, errFile := os.Open(request.ImageFile.Name())
	if errFile != nil {
		return nil, nil, errFile
	}
	defer file.Close()

	writerFormFile, errFormFile := createFormFile(writer, file)
	if errFormFile != nil {
		logrus.
			WithError(errFormFile).
			Error("error")
		return nil, nil, errFormFile
	}

	bFormFile, bErrorFormFile := ioutil.ReadFile(file.Name())
	if bErrorFormFile != nil {
		logrus.
			WithError(bErrorFormFile).
			Error("error")
		return nil, nil, bErrorFormFile
	}
	writerFormFile.Write(bFormFile)

	errTypeField := writer.WriteField("documentType", string(request.DocumentType))
	if errTypeField != nil {
		logrus.
			WithError(errTypeField).
			Error("error document type field")
		return nil, nil, errTypeField
	}

	errSideField := writer.WriteField("documentSide", string(request.DocumentSide))
	if errSideField != nil {
		logrus.
			WithError(errSideField).
			Error("error document side field")
		return nil, nil, errSideField
	}

	errClose := writer.Close()
	if errClose != nil {
		logrus.
			WithError(errClose).
			Error("error writer close")
		return nil, nil, errClose
	}

	return payload, writer, nil
}

func createFormFile(writer *multipart.Writer, file *os.File) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
		escapeQuotes("image"), escapeQuotes(file.Name())))

	contentType, errorContentType := getFileContentType(file)
	if errorContentType != nil {
		logrus.
			WithError(errorContentType).
			Error("error document type field")
		return nil, errorContentType
	}
	h.Set("Content-Type", contentType)
	return writer.CreatePart(h)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func getFileContentType(out *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
