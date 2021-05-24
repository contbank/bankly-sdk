package bankly

import (
	"reflect"
	"strings"

	"github.com/Nhanderu/brdoc"
	"github.com/go-playground/validator/v10"
)

var (
	// Validator ...
	Validator = NewValidator()
)

// NewValidator ...
func NewValidator() *validator.Validate {
	validate := validator.New()

	validate.RegisterTagNameFunc(JSONTagName)

	validate.RegisterValidation("cnpj", CNPJ)
	validate.RegisterValidation("cpf", CPF)
	validate.RegisterValidation("cpfcnpj", CPFCNPJ)

	return validate
}

//JSONTagName ...
func JSONTagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

//CNPJ ...
func CNPJ(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		s := strings.Replace(field.String(), ".", "", -1)
		s = strings.Replace(s, "-", "", -1)
		s = strings.Replace(s, "/", "", -1)
		return brdoc.IsCNPJ(s)
	default:
		return false
	}
}

//CPF ...
func CPF(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		s := strings.Replace(field.String(), ".", "", -1)
		s = strings.Replace(s, "-", "", -1)
		return brdoc.IsCPF(s)
	default:
		return false
	}
}

//CPFCNPJ ...
func CPFCNPJ(fl validator.FieldLevel) bool {
	return CPF(fl) || CNPJ(fl)
}
