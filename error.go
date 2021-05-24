package bankly

import (
	"errors"
)

const (
	// ScouterQuantityKey ...
	ScouterQuantityKey = "SCOUTER_QUANTITY"

	// BoletoInvalidStatusKey ...
	BoletoInvalidStatusKey = "BANKSLIP_SETTLEMENT_STATUS_VALIDATE"
)

var (
	// ErrScouterQuantity ...
	ErrScouterQuantity = errors.New(ScouterQuantityKey)

	// ErrBoletoInvalidStatus ...
	ErrBoletoInvalidStatus = errors.New(BoletoInvalidStatusKey)
)
