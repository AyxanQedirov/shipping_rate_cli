package entities

import "shipping_rate_cli/internal/domain/value-objects"

type Zone struct {
	Name      string
	Surcharge valueobjects.Money
}
