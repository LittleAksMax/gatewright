package cmd

import (
	"errors"
	"time"

	"github.com/LittleAksMax/gatewright/internal/config"
	"github.com/spf13/cobra"
)

type CommandRunEFunc = func(cmd *cobra.Command, args []string) error

func requireLinkedAWSAndGithub(cfg *config.Config) CommandRunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		if !cfg.AWSCfg.IsLinked() {
			return errors.New("aws not linked, use `gwt link aws` to link")
		}
		if !cfg.GithubCfg.IsLinked() {
			return errors.New("github not linked, use `gwt link gh` to link")
		}

		// Expired credentials
		if cfg.GithubCfg.Auth.RefreshTokenExpiresAt.Before(time.Now()) {
			return errors.New("github credentials expired, use `gwt link gh` to link")
		}

		return nil
	}
}
