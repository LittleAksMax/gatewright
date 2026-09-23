package cmd

import (
	"errors"

	"github.com/LittleAksMax/gatewright/internal/config"
	"github.com/spf13/cobra"
)

func newSetupHostCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:           "setup",
		SilenceErrors: true,
		Short:         "Prepare a host for deployments.",
		Long:          "Create a required SSH key on the host and a user that will be the deployment user for CD. Also creates the AWS IAM policies attached to the EC2 instance that are required for all deployments.",
		RunE: config.RequireAWS(cfg, func(cmd *cobra.Command, args []string) error {
			return errors.New("not implemented")
		}),
	}
}
