package cmd

import (
	"errors"

	"github.com/LittleAksMax/gatewright/internal/config"
	"github.com/spf13/cobra"
)

func newDependenciesCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:           "deps",
		SilenceErrors: true,
		Short:         "Verify CLI requirements.",
		Long:          "Ensure required programs are installed. Docker, AWS CLI, GH CLI.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("not implemented")
		},
	}
}
