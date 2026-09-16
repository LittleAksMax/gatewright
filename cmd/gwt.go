package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var gwt = &cobra.Command{
	Use:           "gwt",
	SilenceErrors: true,
	Short:         "Simple tool for quickstarting simple projects.",
	Long:          "Gatewright bootstraps and creates CI for new monorepo projects that deploy as a single Docker container using GH Actions, ECR, and EC2 via SSM.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("a subcommand is required")
	},
}

func Execute() {
	if err := gwt.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
}
