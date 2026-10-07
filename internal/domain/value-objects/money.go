package valueobjects

import (
	"github.com/juju/errors"

	"shipping_rate_cli/internal/domain/enums"
)

type Money struct {
	Amount   float64
	Currency enums.Currency
}

func (this *Money) Add(another Money) error {
	if this.Currency != another.Currency {
		return errors.New("adding money failed: currencies are not the same and multicurrency is not supported")
	}

	this.Amount = this.Amount + another.Amount

	return nil
}
