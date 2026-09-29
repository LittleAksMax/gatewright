package config

import "time"

type GithubAuthConfig struct {
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	AccessToken           string
	RefreshToken          string
}

type GithubConfig struct {
	Auth     *GithubAuthConfig // TODO: fetch from config files
	ClientID string
}
