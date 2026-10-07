package main

import (
	"fmt"
	"os"

	"github.com/juju/errors"

	"shipping_rate_cli/internal/application/services"
	"shipping_rate_cli/internal/infrastructure/repositories"
	"shipping_rate_cli/internal/presentation/commands"
)

func main() {
	repos := repositories.NewRepositories()
	appServices := services.NewServices(repos)
	command := commands.NewCommand(appServices.ParcelService)

	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", errors.Annotate(err, "executing CLI command failed"))
		os.Exit(1)
	}
}
