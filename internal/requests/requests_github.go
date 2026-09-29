package requests

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

const githubDeviceCodeURL = "https://github.com/login/device/code"
const githubOAuthURL = "https://github.com/login/oauth"
const githubAccessTokenURL = githubOAuthURL + "/access_token"

var (
	ErrAuthorizationPending = errors.New("authorization pending")
	ErrSlowDown             = errors.New("slow down")
)

type GithubDeviceCodeResult struct {
	DeviceCode      string
	UserCode        string
	VerificationURI string
	ExpiresIn       int64
	Interval        int64
}

func PostGithubDeviceCode(ctx context.Context, clientID string) (GithubDeviceCodeResult, error) {
	responseBody, err := PostRequest(ctx, githubDeviceCodeURL, map[string]any{
		"client_id": clientID,
	})
	if err != nil {
		return GithubDeviceCodeResult{}, fmt.Errorf("requesting GitHub device code: %w", err)
	}

	values, err := url.ParseQuery(string(responseBody))
	if err != nil {
		return GithubDeviceCodeResult{}, fmt.Errorf("parsing GitHub device code response: %w", err)
	}
	if responseErr := checkURLValuesError(values); responseErr != nil {
		return GithubDeviceCodeResult{}, fmt.Errorf("GitHub device code request failed: %w", responseErr)
	}

	if responseErr := requireURLValues(values, "device_code", "user_code", "verification_uri", "expires_in", "interval"); responseErr != nil {
		return GithubDeviceCodeResult{}, fmt.Errorf("GitHub device code response invalid: %w", responseErr)
	}

	expiresIn, err := strconv.ParseInt(values.Get("expires_in"), 10, 64)
	if err != nil {
		return GithubDeviceCodeResult{}, fmt.Errorf("invalid GitHub expires_in value: %w", err)
	}
	interval, err := strconv.ParseInt(values.Get("interval"), 10, 64)
	if err != nil {
		return GithubDeviceCodeResult{}, fmt.Errorf("invalid GitHub interval value: %w", err)
	}

	return GithubDeviceCodeResult{
		DeviceCode:      values.Get("device_code"),
		UserCode:        values.Get("user_code"),
		VerificationURI: values.Get("verification_uri"),
		ExpiresIn:       expiresIn,
		Interval:        interval,
	}, nil
}

type GithubPollResult struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresIn  int
	RefreshTokenExpiresIn int
}

func PostGithubPollAccessToken(ctx context.Context, clientID string, deviceCode string) (GithubPollResult, error) {
	responseBody, err := PostRequest(ctx, githubAccessTokenURL, map[string]any{
		"client_id":   clientID,
		"device_code": deviceCode,
		"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
	})
	if err != nil {
		return GithubPollResult{}, err
	}

	responseValues, err := url.ParseQuery(string(responseBody))
	if err != nil {
		return GithubPollResult{}, fmt.Errorf("parsing GitHub access token response: %w", err)
	}

	if responseErr := checkURLValuesError(responseValues); responseErr != nil {
		switch responseValues.Get("error") {
		case "authorization_pending":
			return GithubPollResult{}, fmt.Errorf("%w: %w", ErrAuthorizationPending, responseErr)
		case "slow_down":
			return GithubPollResult{}, fmt.Errorf("%w: %w", ErrSlowDown, responseErr)
		}
		return GithubPollResult{}, fmt.Errorf("GitHub access token request failed: %w", responseErr)
	}
	if responseErr := requireURLValues(responseValues, "access_token", "expires_in", "refresh_token", "refresh_token_expires_in", "scope", "token_type"); responseErr != nil {
		return GithubPollResult{}, fmt.Errorf("GitHub access token response invalid: %w", responseErr)
	}

	accessTokenExpiresIn, err := strconv.Atoi(responseValues.Get("expires_in"))
	if err != nil {
		return GithubPollResult{}, fmt.Errorf("invalid GitHub access token expiry: %w", err)
	}
	refreshTokenExpiresIn, err := strconv.Atoi(responseValues.Get("refresh_token_expires_in"))
	if err != nil {
		return GithubPollResult{}, fmt.Errorf("invalid GitHub refresh token expiry: %w", err)
	}

	return GithubPollResult{
		AccessTokenExpiresIn:  accessTokenExpiresIn,
		RefreshTokenExpiresIn: refreshTokenExpiresIn,
		AccessToken:           responseValues.Get("access_token"),
		RefreshToken:          responseValues.Get("refresh_token"),
	}, nil
}
