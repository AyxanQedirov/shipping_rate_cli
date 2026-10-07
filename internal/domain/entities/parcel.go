package entities

import (
	"math"

	"github.com/juju/errors"
)

// Parcel contains weight in kilograms and dimensions in centimetres.
type Parcel struct {
	Weight   float64
	Length   float64
	Width    float64
	Height   float64
	ZoneName string
}

func (this Parcel) Validate() error {
	if this.Weight <= 0 || math.IsNaN(this.Weight) || math.IsInf(this.Weight, 0) {
		return errors.New("validating weight failed: weight must be a finite number greater than zero")
	}

	if this.Length <= 0 || math.IsNaN(this.Length) || math.IsInf(this.Length, 0) {
		return errors.New("validating length failed: length must be a finite number greater than zero")
	}

	if this.Width <= 0 || math.IsNaN(this.Width) || math.IsInf(this.Width, 0) {
		return errors.New("validating width failed: width must be a finite number greater than zero")
	}

	if this.Height <= 0 || math.IsNaN(this.Height) || math.IsInf(this.Height, 0) {
		return errors.New("validating height failed: height must be a finite number greater than zero")
	}

	if this.ZoneName == "" {
		return errors.New("validating zone name failed: zone name can not be empty")
	}

	return nil
}
