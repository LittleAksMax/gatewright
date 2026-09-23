package cmd

import (
	"errors"

	"github.com/LittleAksMax/gatewright/internal/config"
	"github.com/spf13/cobra"
)

func createCommand(cfg *config.Config) *cobra.Command {
	root := &cobra.Command{
		Use:           "gwt",
		SilenceErrors: true,
		Short:         "Simple tool for quickstarting simple projects.",
		Long:          "Gatewright bootstraps and creates CI for new monorepo projects that deploy as a single Docker container using GH Actions, ECR, and EC2 via SSM.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("a subcommand is required")
		},
	}

	root.AddCommand(
		newInitCommand(cfg),
		newSetupHostCommand(cfg),
		newLinkCommand(cfg),
		newDependenciesCommand(cfg),
	)

	return root
}
