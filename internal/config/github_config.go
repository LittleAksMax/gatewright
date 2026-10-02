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
	Auth     *GithubAuthConfig
	clientID string
}

func (g *GithubConfig) IsLinked() bool {
	return g.Auth != nil
}

func (g *GithubConfig) GetClientID() string {
	return g.clientID
}

func readGithubAuthConfigFromPersistent() *GithubAuthConfig {
	path, err := configPath("github")
	if err != nil {
		return nil
	}

	// Check if file exists
	if _, err := os.Stat(path); err != nil {
		// Only return error if error is different from file not existing
		return nil
	}

	// Ensure no nil-dereference
	gAuthToRead := GithubAuthConfig{}
	if err := readFile(path, &gAuthToRead); err != nil {
		return nil
	}
	return &gAuthToRead
}

func (gAuth *GithubAuthConfig) CommitToPersistent() error {
	path, err := configPath("github")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	return writeFile(path, gAuth) // from authconfig.go — sets 0600 on the file
}
