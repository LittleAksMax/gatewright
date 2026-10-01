package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type GithubAuthConfig struct {
	AccessTokenExpiresAt  time.Time `cfg:"accessTokenExpiration"`
	RefreshTokenExpiresAt time.Time `cfg:"refreshTokenExpiration"`
	AccessToken           string    `cfg:"accessToken"`
	RefreshToken          string    `cfg:"refreshToken"`
}

type GithubConfig struct {
	Auth     *GithubAuthConfig // TODO: fetch from config files
	clientID string
}

func (g *GithubConfig) IsLinked() bool {
	return g.Auth != nil
}

func (g *GithubConfig) GetClientID() string {
	return g.clientID
}

func Read() (GithubConfig, error) {
	return GithubConfig{}, nil
}

func (gAuth *GithubAuthConfig) Commit() error {
	path, err := configPath("github")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	return writeFile(path, gAuth) // from authconfig.go — sets 0600 on the file
}
