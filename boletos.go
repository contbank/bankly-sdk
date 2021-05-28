package bankly

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/contbank/grok"
)

//Boletos ...
type Boletos struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

//NewBoletos ...
func NewBoletos(session Session) *Boletos {
	return &Boletos{
		session: session,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authentication: NewAuthentication(session),
	}
}

//CreateBoleto ...
func (b *Boletos) CreateBoleto(model *BoletoRequest) (*BoletoResponse, error) {
	err := grok.Validator.Struct(model)

	if err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, BoletosPath)
	endpoint := u.String()

	reqbyte, err := json.Marshal(model)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(reqbyte))

	if err != nil {
		return nil, err
	}

	token, err := b.authentication.Token()

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", b.session.APIVersion)

	resp, err := b.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusAccepted {
		var body *BoletoResponse

		err = json.Unmarshal(respBody, &body)

		if err != nil {
			return nil, err
		}

		return body, nil
	}

	var bodyErr []*BoletoErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		return nil, err
	}

	if len(bodyErr) > 0 {
		err := bodyErr[0]

		if err.Code == ScouterQuantityKey {
			return nil, ErrScouterQuantity
		}

		if err.Message != nil {
			return nil, errors.New(*err.Message)
		}
	}

	return nil, errors.New("error create boletos")
}

//FindBoleto ...
func (b *Boletos) FindBoleto(model *FindBoletoRequest) (*BoletoDetailedResponse, error) {
	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, "branch")
	u.Path = path.Join(u.Path, model.Account.Branch)
	u.Path = path.Join(u.Path, "number")
	u.Path = path.Join(u.Path, model.Account.Number)
	u.Path = path.Join(u.Path, model.AuthenticationCode)
	endpoint := u.String()

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	token, err := b.authentication.Token()

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("api-version", b.session.APIVersion)

	resp, err := b.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *BoletoDetailedResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			return nil, err
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("not found")
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		return nil, err
	}

	if bodyErr.Errors != nil {
		return nil, errors.New(bodyErr.Errors[0].Messages[0])
	}

	return nil, errors.New("error find boleto")
}

//FilterBoleto ...
func (b *Boletos) FilterBoleto(date time.Time) (*FilterBoletoResponse, error) {
	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, "searchstatus")
	u.Path = path.Join(u.Path, url.QueryEscape(date.UTC().Format("2006-01-02T15:04:05")))
	endpoint := u.String()

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	token, err := b.authentication.Token()

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("api-version", b.session.APIVersion)

	resp, err := b.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *FilterBoletoResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			return nil, err
		}

		return response, nil
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		return nil, err
	}

	if bodyErr.Errors != nil {
		return nil, errors.New(bodyErr.Errors[0].Messages[0])
	}

	return nil, errors.New("error find boleto")
}

//FindBoletoByBarCode ...
func (b *Boletos) FindBoletoByBarCode(barcode string) (*BoletoDetailedResponse, error) {
	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, barcode)
	endpoint := u.String()

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	token, err := b.authentication.Token()

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("api-version", b.session.APIVersion)

	resp, err := b.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response *BoletoDetailedResponse

		err = json.Unmarshal(respBody, &response)

		if err != nil {
			return nil, err
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("not found")
	}

	var bodyErr *ErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		return nil, err
	}

	if bodyErr.Errors != nil {
		return nil, errors.New(bodyErr.Errors[0].Messages[0])
	}

	return nil, errors.New("error find boleto")
}

//DownloadBoleto ...
func (b *Boletos) DownloadBoleto(authenticationCode string, w io.Writer) error {
	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, authenticationCode)
	u.Path = path.Join(u.Path, "pdf")
	endpoint := u.String()

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return err
	}

	token, err := b.authentication.Token()

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", b.session.APIVersion)

	resp, err := b.httpClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		_, err := io.Copy(w, resp.Body)
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("not found")
	}

	return errors.New("error download boleto")
}

//CancelBoleto ...
func (b *Boletos) CancelBoleto(model *CancelBoletoRequest) error {
	err := grok.Validator.Struct(model)

	if err != nil {
		return grok.FromValidationErros(err)
	}

	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, "cancel")
	endpoint := u.String()

	reqbyte, err := json.Marshal(model)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", endpoint, bytes.NewReader(reqbyte))

	if err != nil {
		return err
	}

	token, err := b.authentication.Token()

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", b.session.APIVersion)

	resp, err := b.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return errors.New("error cancel boleto")
}

//PayBoleto simulates a boleto beeing payed
func (b *Boletos) PayBoleto(model *PayBoletoRequest) error {
	err := grok.Validator.Struct(model)

	if err != nil {
		return grok.FromValidationErros(err)
	}

	u, err := url.Parse(b.session.APIEndpoint)

	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, BoletosPath)
	u.Path = path.Join(u.Path, "settlementpayment")
	endpoint := u.String()

	reqbyte, err := json.Marshal(model)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(reqbyte))

	if err != nil {
		return err
	}

	token, err := b.authentication.Token()

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("api-version", b.session.APIVersion)

	resp, err := b.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	var bodyErr []*BoletoErrorResponse

	err = json.Unmarshal(respBody, &bodyErr)

	if err != nil {
		return err
	}

	if len(bodyErr) > 0 {
		err := bodyErr[0]

		if err.Code == BoletoInvalidStatusKey {
			return ErrBoletoInvalidStatus
		}

		if err.Message != nil {
			return errors.New(*err.Message)
		}
	}

	return errors.New("error pay boleto")
}
