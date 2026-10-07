package services

import (
	"github.com/juju/errors"

	"shipping_rate_cli/internal/application/abstractions"
	"shipping_rate_cli/internal/domain/entities"
	valueobjects "shipping_rate_cli/internal/domain/value-objects"
)

var _ abstractions.IParcelService = (*ParcelService)(nil)

type ParcelService struct {
	calculator     abstractions.IShippingRateCalculator
	zoneRepository abstractions.IZoneRepository
}

func NewParcelService(calculator abstractions.IShippingRateCalculator, zoneRepository abstractions.IZoneRepository) *ParcelService {
	return &ParcelService{calculator: calculator, zoneRepository: zoneRepository}
}

func (this *ParcelService) Calculate(parcel entities.Parcel) (*valueobjects.Money, error) {
	if err := parcel.Validate(); err != nil {
		return nil, errors.Annotate(err, "validating parcel failed")
	}

	zone := this.zoneRepository.FindByName(parcel.ZoneName)
	if zone == nil {
		return nil, errors.Errorf("finding zone failed: unknown zone %q", parcel.ZoneName)
	}

	rate, err := this.calculator.Calculate(parcel, *zone)
	if err != nil {
		return nil, errors.Annotate(err, "calculating shipping rate failed")
	}
	return rate, nil
}
