package bankly

import (
	"errors"
	"os"
	"time"

	"github.com/contbank/grok"
	"github.com/patrickmn/go-cache"
)

//Config ...
type Config struct {
	LoginEndpoint *string
	APIEndpoint   *string
	ClientID      *string `validate:"required"`
	ClientSecret  *string `validate:"required"`
	APIVersion    *string
	Cache         *cache.Cache
}

//Session ...
type Session struct {
	LoginEndpoint string
	APIEndpoint   string
	ClientID      string
	ClientSecret  string
	APIVersion    string
	Cache         cache.Cache
}

//NewSession ...
func NewSession(config Config) (*Session, error) {
	err := grok.Validator.Struct(config)

	if err != nil {
		return nil, grok.FromValidationErros(err)
	}

	if config.APIEndpoint == nil {
		config.APIEndpoint = String("https://api.sandbox.bankly.com.br")
	}

	if config.LoginEndpoint == nil {
		config.LoginEndpoint = String("https://login.sandbox.bankly.com.br")
	}

	if config.APIVersion == nil {
		config.APIVersion = String("1.0")
	}

	if config.ClientID == nil {
		config.ClientID = String(os.Getenv("BANKLY_CLIENT_ID"))
	}

	if config.ClientSecret == nil {
		config.ClientID = String(os.Getenv("BANKLY_CLIENT_SECRET"))
	}

	if *config.ClientID == "" || *config.ClientSecret == "" {
		return nil, errors.New("Invalid client id or client secret")
	}

	if config.Cache == nil {
		config.Cache = cache.New(10*time.Minute, 1*time.Second)
	}

	var session = &Session{
		LoginEndpoint: *config.LoginEndpoint,
		APIEndpoint:   *config.APIEndpoint,
		ClientID:      *config.ClientID,
		ClientSecret:  *config.ClientSecret,
		APIVersion:    *config.APIVersion,
		Cache:         *config.Cache,
	}

	return session, nil
}
