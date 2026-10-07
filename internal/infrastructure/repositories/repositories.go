package repositories

import "shipping_rate_cli/internal/application/abstractions"

type Repositories struct {
	ZoneRepository abstractions.IZoneRepository
}

func NewRepositories() *Repositories {
	return &Repositories{
		ZoneRepository: NewZoneRepository(),
	}
}
