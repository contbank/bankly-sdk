package bankly

import (
	"github.com/contbank/grok"
	"net/http"
	"strings"
)

var errorPixList = []Error{
	{
		ErrorKey:  "ENTRY_NOT_FOUND",
		GrokError: ErrKeyNotFound,
	},
	{
		ErrorKey:  "INVALID_QRCODE_PAYLOAD_CONTENT_TO_DECODE",
		GrokError: ErrInvalidQrCodePayload,
	},
	{
		ErrorKey:  "INVALID_KEY_TYPE",
		GrokError: ErrInvalidKeyType,
	},
	{
		ErrorKey:  "INVALID_PARAMETER_PIX",
		GrokError: ErrInvalidParameterPix,
	},
	{
		ErrorKey:  "INSUFFICIENT_BALANCE",
		GrokError: ErrInsufficientBalancePix,
	},
	{
		ErrorKey:  "INVALID_ACCOUNT_TYPE",
		GrokError: ErrInvalidAccountType,
	},
}

func verifyInvalidPixParameter(code string, messages []string) string {
	if code == "INVALID_PARAMETER" {
		for _, m := range messages {
			switch {
			case strings.Contains(strings.ToLower(m), "addressing key value does not match with addressing key type"):
				return "INVALID_KEY_TYPE"
			case strings.Contains(strings.ToLower(m), "sender.account.type"):
				return "INVALID_ACCOUNT_TYPE"
			default:
				return "INVALID_PARAMETER_PIX"
			}
		}
	}
	return code
}

// FindPixError
func FindPixError(code string, messages ...string) *grok.Error {
	code = verifyInvalidPixParameter(code, messages)

	for _, v := range errorPixList {
		if v.ErrorKey == code {
			return v.GrokError
		}
	}

	return grok.NewError(http.StatusConflict, messages...)
}