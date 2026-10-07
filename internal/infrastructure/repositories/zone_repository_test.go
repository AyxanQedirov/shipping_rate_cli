package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"shipping_rate_cli/internal/domain/enums"
	"shipping_rate_cli/internal/infrastructure/repositories"
)

type ZoneRepositorySuite struct {
	suite.Suite
	repository *repositories.ZoneRepository
}

func (this *ZoneRepositorySuite) SetupTest() {
	this.repository = repositories.NewZoneRepository()
}

func (this *ZoneRepositorySuite) Test_FindByName_RetursZone_InCaseOfKnownZone() {
	//Arrange
	testCases := []struct {
		name      string
		zoneName  string
		surcharge float64
	}{
		{"DomesticZone", "domestic", 0},
		{"EUZone", "eu", 8},
		{"InternationalZone", "international", 15},
	}

	for _, testCase := range testCases {
		this.Run(testCase.name, func() {
			//Act
			zone := this.repository.FindByName(testCase.zoneName)

			//Assert
			this.Require().NotNil(zone)
			this.Equal(testCase.zoneName, zone.Name)
			this.Equal(testCase.surcharge, zone.Surcharge.Amount)
			this.Equal(enums.EUR, zone.Surcharge.Currency)
		})
	}
}

func (this *ZoneRepositorySuite) Test_FindByName_RetursNil_InCaseOfUnknownZone() {
	//Arrange
	testCases := []struct {
		name     string
		zoneName string
	}{
		{"UnknownZone", "unknown"},
		{"UppercaseZone", "EU"},
		{"EmptyZone", ""},
	}

	for _, testCase := range testCases {
		this.Run(testCase.name, func() {
			//Act
			zone := this.repository.FindByName(testCase.zoneName)

			//Assert
			this.Nil(zone)
		})
	}
}

func TestZoneRepositorySuite(t *testing.T) {
	suite.Run(t, new(ZoneRepositorySuite))
}
