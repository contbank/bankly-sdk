package bankly

import (
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
	Scope         string
}

//ServiceDeskConfig ...
type ServiceDeskConfig struct {
	APIEndpoint *string
	APIKey      *string
}

//ServiceDeskSession ...
type ServiceDeskSession struct {
	APIEndpoint string
	APIKey      string
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
		return nil, ErrClientIDClientSecret
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

//NewServiceDeskSession ...
func NewServiceDeskSession(config ServiceDeskConfig) (*ServiceDeskSession, error) {
	err := grok.Validator.Struct(config)

	if err != nil {
		return nil, grok.FromValidationErros(err)
	}

	if config.APIEndpoint == nil {
		config.APIEndpoint = String("https://meuacesso.freshdesk.com")
	}

	if config.APIKey == nil {
		config.APIKey = String(os.Getenv("FRESH_DESK_API_KEY"))
	}

	var session = &ServiceDeskSession{
		APIEndpoint: *config.APIEndpoint,
		APIKey:      *config.APIKey,
	}

	return session, nil
}
