package services

import (
	"shipping_rate_cli/internal/application/abstractions"
	"shipping_rate_cli/internal/infrastructure/repositories"
)

type Services struct {
	ShippingRateCalculator abstractions.IShippingRateCalculator
	ParcelService          abstractions.IParcelService
}

func NewServices(repos *repositories.Repositories) *Services {
	calculator := NewShippingRateCalculator()

	return &Services{
		ShippingRateCalculator: calculator,
		ParcelService:          NewParcelService(calculator, repos.ZoneRepository),
	}
}
