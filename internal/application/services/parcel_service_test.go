package services_test

import (
	"testing"

	"github.com/juju/errors"
	"github.com/stretchr/testify/suite"

	"shipping_rate_cli/internal/application/services"
	"shipping_rate_cli/internal/domain/entities"
	"shipping_rate_cli/internal/domain/enums"
	valueobjects "shipping_rate_cli/internal/domain/value-objects"
	"shipping_rate_cli/internal/infrastructure/repositories"
)

type ParcelServiceSuite struct {
	suite.Suite
	service    *services.ParcelService
	repository *repositories.MockZoneRepository
	calculator *services.MockShippingRateCalculator
	parcel     entities.Parcel
	zone       entities.Zone
}

func (this *ParcelServiceSuite) SetupTest() {
	this.repository = repositories.NewMockZoneRepository(this.T())
	this.calculator = services.NewMockShippingRateCalculator(this.T())

	this.service = services.NewParcelService(this.calculator, this.repository)

	this.parcel = entities.Parcel{Weight: 3, Length: 30, Width: 20, Height: 10, ZoneName: "eu"}
	this.zone = entities.Zone{Name: "eu", Surcharge: valueobjects.Money{Amount: 8, Currency: enums.EUR}}
}

func (this *ParcelServiceSuite) Test_Calculate_RetursMoney_InCaseOfValidParcelAndKnownZone() {
	//Arrange
	expected := &valueobjects.Money{Amount: 18, Currency: enums.EUR}

	this.repository.
		EXPECT().
		FindByName("eu").
		Return(&this.zone).
		Once()

	this.calculator.
		EXPECT().
		Calculate(this.parcel, this.zone).
		Return(expected, nil).
		Once()

	//Act
	rate, err := this.service.Calculate(this.parcel)

	//Assert
	this.NoError(err)
	this.Same(expected, rate)
}

func (this *ParcelServiceSuite) Test_Calculate_RetursError_InCaseOfInvalidParcel() {
	//Arrange
	this.parcel.Weight = 0

	//Act
	rate, err := this.service.Calculate(this.parcel)

	//Assert
	this.Nil(rate)
	this.EqualError(err, "validating parcel failed: validating weight failed: weight must be a finite number greater than zero")
}

func (this *ParcelServiceSuite) Test_Calculate_RetursError_InCaseOfUnknownZone() {
	//Arrange
	this.parcel.ZoneName = "unknown"

	this.repository.
		EXPECT().
		FindByName("unknown").
		Return(nil).
		Once()

	//Act
	rate, err := this.service.Calculate(this.parcel)

	//Assert
	this.Nil(rate)
	this.EqualError(err, "finding zone failed: unknown zone \"unknown\"")
}

func (this *ParcelServiceSuite) Test_Calculate_RetursAnnotatedError_InCaseOfCalculatorFailure() {
	//Arrange
	cause := errors.New("pricing failed")

	this.repository.
		EXPECT().
		FindByName("eu").
		Return(&this.zone).
		Once()

	this.calculator.
		EXPECT().
		Calculate(this.parcel, this.zone).
		Return(nil, cause).
		Once()

	//Act
	rate, err := this.service.Calculate(this.parcel)

	//Assert
	this.Nil(rate)
	this.EqualError(err, "calculating shipping rate failed: pricing failed")
	this.Same(cause, errors.Cause(err), "annotation must preserve the original error")
}

func TestParcelServiceSuite(t *testing.T) {
	suite.Run(t, new(ParcelServiceSuite))
}
