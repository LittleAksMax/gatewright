package cmd

import (
	"errors"

	"github.com/LittleAksMax/gatewright/internal/config"
	"github.com/spf13/cobra"
)

type CommandRunEFunc = func(cmd *cobra.Command, args []string) error

func requireLinkedAWSAndGithub(cfg *config.Config) CommandRunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		if cfg.IsAWSLinked() {
			return errors.New("aws not linked, use `gwt link aws` to link")
		}
		if cfg.IsGHLinked() {
			return errors.New("github not linked, use `gwt link gh` to link")
		}
		return nil
	}
}
