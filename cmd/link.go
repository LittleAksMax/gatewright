package cmd

import (
	"errors"

	"github.com/LittleAksMax/gatewright/internal/config"
	"github.com/spf13/cobra"
)

const (
	ghArg  string = "gh"
	awsArg string = "aws"
)

func newLinkCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:           "link",
		SilenceErrors: true,
		Short:         "Link Github (gh) or AWS (aws) to make creating required policies.",
		Long:          "Create relevant Github Actions workflows, Makefile (with relevant targets) a plain Dockerfile. It creates the relevant AWS Policies and ECR repo. It also sets up the GitHub repo with the assumed conditions.",
		ValidArgs: []cobra.Completion{ghArg, awsArg},
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("not implemented")
		},
	}
}
