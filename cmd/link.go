package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LittleAksMax/gatewright/internal/config"
	"github.com/LittleAksMax/gatewright/internal/requests"
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
		ValidArgs:     []cobra.Completion{ghArg, awsArg},
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}
			if args[0] != ghArg && args[0] != awsArg {
				return fmt.Errorf("invalid argument %q: expected %q or %q", args[0], ghArg, awsArg)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case ghArg:
				return handleGithubLinkFlow(cmd.Context(), cfg)
			case awsArg:
				return handleAWSLinkFlow(cmd.Context(), cfg)
			default:
				return fmt.Errorf("unrecognised argument: %s", args[0])
			}
		},
	}
}

func handleGithubLinkFlow(parent context.Context, cfg *config.Config) error {
	deviceAuth, err := requests.PostGithubDeviceCode(parent, cfg.GithubCfg.GetClientID())
	if err != nil {
		return err
	}

	// Ask user to authorise device with provided details
	ctx, cancel := context.WithTimeout(parent, time.Duration(deviceAuth.ExpiresIn)*time.Second)
	defer cancel()

	interval := time.Duration(deviceAuth.Interval) * time.Second
	timer := time.NewTimer(interval)
	defer timer.Stop()

	fmt.Printf("Copy this code to clipboard: %s\n", deviceAuth.UserCode)
	fmt.Printf("Sign in with the user code at this link: %s\n", deviceAuth.VerificationURI)

	var token requests.GithubPollResult
	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				fmt.Println("Timed out authorising device")
				return nil
			}
			return fmt.Errorf("device authorisation canceled: %w", ctx.Err())
		case <-timer.C:
		}

		var err error
		token, err = requests.PostGithubPollAccessToken(ctx, cfg.GithubCfg.GetClientID(), deviceAuth.DeviceCode)
		if err == nil {
			break
		}

		switch {
		case ctx.Err() != nil:
			continue // let the select above report timeout/cancel
		case errors.Is(err, requests.ErrAuthorizationPending):
			// keep polling
		case errors.Is(err, requests.ErrSlowDown):
			interval += 5 * time.Second // required by GitHub's device flow spec
		default:
			return err // expired_token, access_denied, network failure, etc.
		}
		timer.Reset(interval)
	}

	// Store the credentials
	cfg.GithubCfg.Auth = &config.GithubAuthConfig{
		AccessToken:           token.AccessToken,
		RefreshToken:          token.RefreshToken,
		AccessTokenExpiresAt:  time.Now().Add(time.Second * time.Duration(token.AccessTokenExpiresIn)),
		RefreshTokenExpiresAt: time.Now().Add(time.Second * time.Duration(token.RefreshTokenExpiresIn)),
	}
	if err := cfg.GithubCfg.Auth.CommitToPersistent(); err != nil {
		return err
	}

	return nil
}

func handleAWSLinkFlow(ctx context.Context, cfg *config.Config) error {
	return errors.New("not implemented")
}
