package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var hostCmd = &cobra.Command{
	Use:           "host",
	SilenceErrors: true,
	Short:         "Prepare a host for deployment.",
	Long:          "Create a required SSH key on the host and a user that will be the deployment user for CD. Also creates the AWS IAM policies attached to the EC2 instance that are required for all deployments.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}