package abstractions

import (
	"shipping_rate_cli/internal/domain/entities"
	valueobjects "shipping_rate_cli/internal/domain/value-objects"
)

// IShippingRateCalculator is the pricing strategy injected into ParcelService.
type IShippingRateCalculator interface {
	Calculate(parcel entities.Parcel, zone entities.Zone) (*valueobjects.Money, error)
}
