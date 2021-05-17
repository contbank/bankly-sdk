package bankly

import (
	"errors"
	"os"
)

//Config ...
type Config struct {
	LoginEndpoint *string
	APIEndpoint   *string
	ClientID      *string `validate:"required"`
	ClientSecret  *string `validate:"required"`
	APIVersion    *string
	Cache         *Redis
}

//Redis ...
type Redis struct {
	Endpoint string
	Port     string
	User     string
	Pass     string
}

//Session ...
type Session struct {
	LoginEndpoint string
	APIEndpoint   string
	ClientID      string
	ClientSecret  string
	APIVersion    string
	Cache         *Redis
}

//NewSession ...
func NewSession(config Config) (*Session, error) {
	err := Validator.Struct(config)

	if err != nil {
		return nil, err
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

	err = validateCredentials()
	if err != nil {
		return nil, err
	}

	var session = &Session{
		LoginEndpoint: *config.LoginEndpoint,
		APIEndpoint:   *config.APIEndpoint,
		ClientID:      *config.ClientID,
		ClientSecret:  *config.ClientSecret,
		APIVersion:    *config.APIVersion,
	}

	return session, nil
}

func validateCredentials() error {
	if os.Getenv("BANKLY_CLIENT_ID") == "" || os.Getenv("BANKLY_CLIENT_SECRET") == "" {
		return errors.New("Invalid client id or client secret")
	}
	return nil
}
