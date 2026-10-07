package commands

import (
	"fmt"
	"strings"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"shipping_rate_cli/internal/application/abstractions"
	"shipping_rate_cli/internal/domain/entities"
)

const (
	flagWeight = "weight"
	flagLength = "length"
	flagWidth  = "width"
	flagHeight = "height"
	flagZone   = "zone"
)

func NewCalculateCommand(service abstractions.IParcelService) *cobra.Command {
	var parcel entities.Parcel
	var zoneName string

	command := &cobra.Command{
		Use:           "calculate",
		Short:         "Calculate a parcel's shipping cost in EUR",
		Example:       "shipping-rate calculate --weight 3 --length 30 --width 20 --height 10 --zone eu",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			parcel.ZoneName = strings.ToLower(zoneName)

			rate, err := service.Calculate(parcel)
			if err != nil {
				return errors.Annotate(err, "calculating parcel shipping rate failed")
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s %.2f\n", rate.Currency, rate.Amount)
			return errors.Annotate(err, "writing shipping rate failed")
		},
	}

	flags := command.Flags()
	flags.Float64Var(&parcel.Weight, flagWeight, 0, "Parcel weight in kilograms")
	flags.Float64Var(&parcel.Length, flagLength, 0, "Parcel length in centimetres")
	flags.Float64Var(&parcel.Width, flagWidth, 0, "Parcel width in centimetres")
	flags.Float64Var(&parcel.Height, flagHeight, 0, "Parcel height in centimetres")
	flags.StringVar(&zoneName, flagZone, "", "Destination zone: domestic, eu, international")

	for _, name := range []string{flagWeight, flagLength, flagWidth, flagHeight, flagZone} {
		if err := command.MarkFlagRequired(name); err != nil {
			panic(errors.Annotate(err, "marking required flag failed"))
		}
	}

	return command
}
