package abstractions

import "shipping_rate_cli/internal/domain/entities"

type IZoneRepository interface {
	FindByName(name string) *entities.Zone
}
