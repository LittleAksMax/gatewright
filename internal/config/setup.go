package config

func CreateConfig() *Config {
	return &Config{
		AWSCfg: AWSConfig{},
		GithubCfg: GithubConfig{
			clientID: "Iv23liMjsglWwuqqgwdZ",
			Auth:     readGithubAuthConfigFromPersistent(), // read later
		},
	}
}
