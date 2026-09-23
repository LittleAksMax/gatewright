package cmd

import (
	"errors"

	"github.com/LittleAksMax/gatewright/internal/config"
	"github.com/spf13/cobra"
)

func newInitCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:           "init",
		SilenceErrors: true,
		Short:         "Initialise a project to be gatewright-managed.",
		Long:          "Create relevant Github Actions workflows, Makefile (with relevant targets) a plain Dockerfile. It creates the relevant AWS Policies and ECR repo. It also sets up the GitHub repo with the assumed conditions.",
		PreRunE: requireLinkedAWSAndGithub(cfg),
		RunE: config.RequireAWSAndGithub(cfg, func(cmd *cobra.Command, args []string) error {
			return errors.New("not implemented")
		}),
	}
}
