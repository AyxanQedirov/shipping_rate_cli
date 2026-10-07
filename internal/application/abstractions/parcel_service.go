package abstractions

import (
	"shipping_rate_cli/internal/domain/entities"
	"shipping_rate_cli/internal/domain/value-objects"
)

type IParcelService interface {
	Calculate(parcel entities.Parcel) (*valueobjects.Money, error)
}
