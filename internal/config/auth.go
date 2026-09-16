package config

import "github.com/spf13/cobra"

type CommandRunEFunc func(cmd *cobra.Command, args []string) error

func RequireAWS(cfg *Config, runFunc CommandRunEFunc) CommandRunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		// TODO: the actual gate itself
		return runFunc(cmd, args)
	}
}

func RequireGithub(cfg *Config, runFunc CommandRunEFunc) CommandRunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		// TODO: the actual gate itself
		return runFunc(cmd, args)
	}
}