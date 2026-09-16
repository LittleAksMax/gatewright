package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:           "init",
	SilenceErrors: true,
	Short:         "Initialise a project to be gatewright-managed.",
	Long:          "Create relevant Github Actions workflows, Makefile (with relevant targets) a plain Dockerfile. It creates the relevant AWS Policies and ECR repo. It also sets up the GitHub repo with the assumed conditions.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}