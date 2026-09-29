package config

func CreateConfig() *Config {
	return &Config{
		awsCfg: nil,
		ghCfg: &GithubConfig{
			ClientID: "Iv23liMjsglWwuqqgwdZ",
			Auth:     nil, // TODO: read properly
		},
	}
}
