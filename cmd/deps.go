package cmd

import (
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
			if args[0] == ghArg {
				return handleGithubLinkFlow(cmd.Context(), cfg)
			}
			if args[1] == awsArg {
				return handleAWSLinkFlow(cmd.Context(), cfg)
			}
			return nil // should be unreachable
		},
	}
}
