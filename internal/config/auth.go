package config

import (
	"errors"

	"github.com/spf13/cobra"
)

type CommandRunEFunc func(cmd *cobra.Command, args []string) error

func RequireAWSAndGithub(cfg *Config, runFunc CommandRunEFunc) CommandRunEFunc {
	return RequireAWS(cfg, RequireGithub(cfg, runFunc))
}

func RequireAWS(cfg *Config, runFunc CommandRunEFunc) CommandRunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		if !cfg.IsAWSLinked() {
			return errors.New("must link AWS with `gwt link aws`")
		}
		return runFunc(cmd, args)
	}
}

func RequireGithub(cfg *Config, runFunc CommandRunEFunc) CommandRunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		if !cfg.IsGHLinked() {
			return errors.New("must link Github with `gwt link gh`")
		}
		return runFunc(cmd, args)
	}
}
