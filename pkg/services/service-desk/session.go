package bankly

import (
	utils "github.com/contbank/bankly-sdk/pkg/utils"
	"github.com/contbank/grok"
	"os"
)

// ServiceDeskConfig ...
type ServiceDeskConfig struct {
	APIEndpoint *string
	APIKey      *string
}

// ServiceDeskSession ...
type ServiceDeskSession struct {
	APIEndpoint string
	APIKey      string
}

// NewServiceDeskSession ...
func NewServiceDeskSession(config ServiceDeskConfig) (*ServiceDeskSession, error) {
	err := grok.Validator.Struct(config)

	if err != nil {
		return nil, grok.FromValidationErros(err)
	}

	if config.APIEndpoint == nil {
		config.APIEndpoint = utils.String("https://meuacesso.freshdesk.com")
	}

	if config.APIKey == nil {
		config.APIKey = utils.String(os.Getenv("FRESH_DESK_API_KEY"))
	}

	var session = &ServiceDeskSession{
		APIEndpoint: *config.APIEndpoint,
		APIKey:      *config.APIKey,
	}

	return session, nil
}