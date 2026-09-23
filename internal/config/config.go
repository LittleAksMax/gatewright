package config

type Config struct {
	awsCfg *AWSConfig
	ghCfg  *GithubConfig
}

func (c *Config) IsAWSLinked() bool {
	return c.awsCfg != nil
}

func (c *Config) IsGHLinked() bool {
	return c.ghCfg != nil
}
