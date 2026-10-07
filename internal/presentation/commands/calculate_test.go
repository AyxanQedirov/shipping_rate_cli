package commands_test

import (
	"io"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"

	"shipping_rate_cli/internal/application/services"
	"shipping_rate_cli/internal/domain/entities"
	"shipping_rate_cli/internal/domain/enums"
	valueobjects "shipping_rate_cli/internal/domain/value-objects"
	"shipping_rate_cli/internal/presentation/commands"
)

type CalculateCommandSuite struct {
	suite.Suite
	service *services.MockParcelService
	command *cobra.Command
}

func (this *CalculateCommandSuite) SetupTest() {
	this.service = services.NewMockParcelService(this.T())

	this.command = commands.NewCommand(this.service)
	this.command.SetOut(io.Discard)
	this.command.SetErr(io.Discard)
}

func (this *CalculateCommandSuite) Test_Execute_RetursNil_InCaseOfValidFlags() {
	//Arrange
	parcel := entities.Parcel{
		Weight:   3,
		Length:   30,
		Width:    20,
		Height:   10,
		ZoneName: "eu",
	}
	result := &valueobjects.Money{
		Amount:   18,
		Currency: enums.EUR,
	}

	testCases := []struct {
		name     string
		zoneName string
	}{
		{"LowercaseZone", "eu"},
		{"UppercaseZone", "EU"},
		{"MixedCaseZone", "Eu"},
	}

	for _, testCase := range testCases {
		this.Run(testCase.name, func() {
			this.service.
				EXPECT().
				Calculate(parcel).
				Return(result, nil).
				Once()

			this.command.SetArgs([]string{
				"calculate",
				"--weight", "3",
				"--length", "30",
				"--width", "20",
				"--height", "10",
				"--zone", testCase.zoneName,
			})

			//Act
			err := this.command.Execute()

			//Assert
			this.NoError(err)
		})
	}
}

func (this *CalculateCommandSuite) Test_Execute_RetursError_InCaseOfMissingRequiredFlags() {
	//Arrange
	this.command.SetArgs([]string{"calculate", "--weight", "3", "--zone", "eu"})

	//Act
	err := this.command.Execute()

	//Assert
	this.EqualError(err, `required flag(s) "height", "length", "width" not set`)
}

func TestCalculateCommandSuite(t *testing.T) {
	suite.Run(t, new(CalculateCommandSuite))
}
