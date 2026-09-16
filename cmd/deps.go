package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var dependenciesCmd = &cobra.Command{
	Use:           "deps",
	SilenceErrors: true,
	Short:         "Verify CLI requirements.",
	Long:          "Ensure required programs are installed. Docker, AWS CLI, GH CLI.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}