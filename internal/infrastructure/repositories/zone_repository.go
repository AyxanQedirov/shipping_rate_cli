package repositories

import (
	"shipping_rate_cli/internal/application/abstractions"
	"shipping_rate_cli/internal/domain/entities"
	"shipping_rate_cli/internal/domain/enums"
	valueobjects "shipping_rate_cli/internal/domain/value-objects"
)

var _ abstractions.IZoneRepository = (*ZoneRepository)(nil)

type ZoneRepository struct {
	zones map[string]entities.Zone
}

func NewZoneRepository() *ZoneRepository {
	return &ZoneRepository{
		zones: map[string]entities.Zone{
			"domestic":      {Name: "domestic", Surcharge: valueobjects.Money{Amount: 0, Currency: enums.EUR}},
			"eu":            {Name: "eu", Surcharge: valueobjects.Money{Amount: 8, Currency: enums.EUR}},
			"international": {Name: "international", Surcharge: valueobjects.Money{Amount: 15, Currency: enums.EUR}},
		},
	}
}

func (this *ZoneRepository) FindByName(name string) *entities.Zone {
	zone, ok := this.zones[name]
	if !ok {
		return nil
	}

	return &zone
}
