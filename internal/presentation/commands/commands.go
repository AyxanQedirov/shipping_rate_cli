package commands

import (
	"github.com/spf13/cobra"

	"shipping_rate_cli/internal/application/abstractions"
)

func NewCommand(service abstractions.IParcelService) *cobra.Command {
	command := &cobra.Command{
		Use:           "shipping-rate",
		Short:         "Calculate parcel shipping rates",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	command.AddCommand(NewCalculateCommand(service))
	return command
}
