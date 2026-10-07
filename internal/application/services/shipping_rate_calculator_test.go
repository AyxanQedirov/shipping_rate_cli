package services_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/suite"

	"shipping_rate_cli/internal/application/services"
	"shipping_rate_cli/internal/domain/entities"
	"shipping_rate_cli/internal/domain/enums"
	valueobjects "shipping_rate_cli/internal/domain/value-objects"
)

type ShippingRateCalculatorSuite struct {
	suite.Suite
	calculator *services.ShippingRateCalculator
}

func (this *ShippingRateCalculatorSuite) SetupTest() {
	this.calculator = services.NewShippingRateCalculator()
}

func (this *ShippingRateCalculatorSuite) Test_Calculate_RetursMoney_InCaseOfAllWeightTiersAndZones() {
	//Arrange
	testCases := []struct {
		name                      string
		weight                    float64
		zoneName                  string
		surcharge, expectedAmount float64
	}{
		{"DomesticBelow1Kg", 0.5, "domestic", 0, 5},
		{"DomesticExactly1Kg", 1, "domestic", 0, 5},
		{"DomesticAbove1Kg", math.Nextafter(1, math.Inf(1)), "domestic", 0, 10},
		{"DomesticBetween1And5Kg", 3, "domestic", 0, 10},
		{"DomesticExactly5Kg", 5, "domestic", 0, 10},
		{"DomesticAbove5Kg", math.Nextafter(5, math.Inf(1)), "domestic", 0, 20},
		{"DomesticHeavyParcel", 10, "domestic", 0, 20},
		{"EuBelow1Kg", 0.5, "eu", 8, 13},
		{"EuExactly1Kg", 1, "eu", 8, 13},
		{"EuAbove1Kg", math.Nextafter(1, math.Inf(1)), "eu", 8, 18},
		{"EuBetween1And5Kg", 3, "eu", 8, 18},
		{"EuExactly5Kg", 5, "eu", 8, 18},
		{"EuAbove5Kg", math.Nextafter(5, math.Inf(1)), "eu", 8, 28},
		{"EuHeavyParcel", 10, "eu", 8, 28},
		{"InternationalBelow1Kg", 0.5, "international", 15, 20},
		{"InternationalExactly1Kg", 1, "international", 15, 20},
		{"InternationalAbove1Kg", math.Nextafter(1, math.Inf(1)), "international", 15, 25},
		{"InternationalBetween1And5Kg", 3, "international", 15, 25},
		{"InternationalExactly5Kg", 5, "international", 15, 25},
		{"InternationalAbove5Kg", math.Nextafter(5, math.Inf(1)), "international", 15, 35},
		{"InternationalHeavyParcel", 10, "international", 15, 35},
	}

	for _, testCase := range testCases {
		this.Run(testCase.name, func() {
			parcel := entities.Parcel{Weight: testCase.weight, Length: 30, Width: 20, Height: 10, ZoneName: testCase.zoneName}
			zone := entities.Zone{Name: testCase.zoneName, Surcharge: valueobjects.Money{Amount: testCase.surcharge, Currency: enums.EUR}}

			//Act
			rate, err := this.calculator.Calculate(parcel, zone)

			//Assert
			this.Require().NoError(err)
			this.Require().NotNil(rate)
			this.Equal(testCase.expectedAmount, rate.Amount)
			this.Equal(enums.EUR, rate.Currency)
		})
	}
}

func (this *ShippingRateCalculatorSuite) Test_Calculate_RetursError_InCaseOfDifferentCurrencies() {
	//Arrange
	parcel := entities.Parcel{Weight: 3, ZoneName: "eu"}
	zone := entities.Zone{Name: "eu", Surcharge: valueobjects.Money{Amount: 8, Currency: enums.Currency("USD")}}

	//Act
	rate, err := this.calculator.Calculate(parcel, zone)

	//Assert
	this.Nil(rate)
	this.EqualError(err, "adding zone surcharge failed: adding money failed: currencies are not the same and multicurrency is not supported")
}

func TestShippingRateCalculatorSuite(t *testing.T) {
	suite.Run(t, new(ShippingRateCalculatorSuite))
}
