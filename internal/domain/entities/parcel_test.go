package entities_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/suite"

	"shipping_rate_cli/internal/domain/entities"
)

type ParcelSuite struct {
	suite.Suite
	parcel entities.Parcel
}

func (this *ParcelSuite) SetupTest() {
	this.parcel = entities.Parcel{Weight: 3, Length: 30, Width: 20, Height: 10, ZoneName: "eu"}
}

func (this *ParcelSuite) Test_Validate_RetursNil_InCaseOfValidParcel() {
	//Arrange
	parcel := this.parcel

	//Act
	err := parcel.Validate()

	//Assert
	this.NoError(err)
}

func (this *ParcelSuite) Test_Validate_RetursError_InCaseOfInvalidMeasurements() {
	//Arrange
	testCases := []struct {
		name                          string
		field                         string
		weight, length, width, height float64
	}{
		{"ZeroWeight", "weight", 0, 30, 20, 10},
		{"NegativeWeight", "weight", -1, 30, 20, 10},
		{"NaNWeight", "weight", math.NaN(), 30, 20, 10},
		{"PositiveInfinityWeight", "weight", math.Inf(1), 30, 20, 10},
		{"NegativeInfinityWeight", "weight", math.Inf(-1), 30, 20, 10},
		{"ZeroLength", "length", 3, 0, 20, 10},
		{"NegativeLength", "length", 3, -1, 20, 10},
		{"NaNLength", "length", 3, math.NaN(), 20, 10},
		{"PositiveInfinityLength", "length", 3, math.Inf(1), 20, 10},
		{"NegativeInfinityLength", "length", 3, math.Inf(-1), 20, 10},
		{"ZeroWidth", "width", 3, 30, 0, 10},
		{"NegativeWidth", "width", 3, 30, -1, 10},
		{"NaNWidth", "width", 3, 30, math.NaN(), 10},
		{"PositiveInfinityWidth", "width", 3, 30, math.Inf(1), 10},
		{"NegativeInfinityWidth", "width", 3, 30, math.Inf(-1), 10},
		{"ZeroHeight", "height", 3, 30, 20, 0},
		{"NegativeHeight", "height", 3, 30, 20, -1},
		{"NaNHeight", "height", 3, 30, 20, math.NaN()},
		{"PositiveInfinityHeight", "height", 3, 30, 20, math.Inf(1)},
		{"NegativeInfinityHeight", "height", 3, 30, 20, math.Inf(-1)},
	}

	for _, testCase := range testCases {
		this.Run(testCase.name, func() {
			parcel := entities.Parcel{
				Weight:   testCase.weight,
				Length:   testCase.length,
				Width:    testCase.width,
				Height:   testCase.height,
				ZoneName: "eu",
			}

			//Act
			err := parcel.Validate()

			//Assert
			this.EqualError(err, "validating "+testCase.field+" failed: "+testCase.field+" must be a finite number greater than zero")
		})
	}
}

func (this *ParcelSuite) Test_Validate_RetursError_InCaseOfEmptyZoneName() {
	//Arrange
	this.parcel.ZoneName = ""

	//Act
	err := this.parcel.Validate()

	//Assert
	this.EqualError(err, "validating zone name failed: zone name can not be empty")
}

func TestParcelSuite(t *testing.T) {
	suite.Run(t, new(ParcelSuite))
}
