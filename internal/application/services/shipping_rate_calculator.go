package services

import (
	"github.com/juju/errors"

	"shipping_rate_cli/internal/application/abstractions"
	"shipping_rate_cli/internal/domain/entities"
	"shipping_rate_cli/internal/domain/enums"
	valueobjects "shipping_rate_cli/internal/domain/value-objects"
)

var _ abstractions.IShippingRateCalculator = (*ShippingRateCalculator)(nil)

type ShippingRateCalculator struct {
}

func NewShippingRateCalculator() *ShippingRateCalculator {
	return &ShippingRateCalculator{}
}

func (this *ShippingRateCalculator) Calculate(parcel entities.Parcel, zone entities.Zone) (*valueobjects.Money, error) {
	baseRate := this.calculateBaseRate(parcel)

	if err := baseRate.Add(zone.Surcharge); err != nil {
		return nil, errors.Annotate(err, "adding zone surcharge failed")
	}

	return baseRate, nil
}

func (this *ShippingRateCalculator) calculateBaseRate(parcel entities.Parcel) *valueobjects.Money {
	if parcel.Weight <= 1 {
		return &valueobjects.Money{Amount: 5, Currency: enums.EUR}
	}

	if parcel.Weight <= 5 {
		return &valueobjects.Money{Amount: 10, Currency: enums.EUR}
	}

	return &valueobjects.Money{Amount: 20, Currency: enums.EUR}
}
