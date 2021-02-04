package mikku

import (
	"errors"
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

var (
	errEmptyGitHubAccessToken = errors.New("should be set MIKKU_GITHUB_ACCESS_TOKEN")
	errEmptyGitHubOwner       = errors.New("should be set MIKKU_GITHUB_OWNER")
)

// Config represents config using all commands
type Config struct {
	GitHubAccessToken string `envconfig:"MIKKU_GITHUB_ACCESS_TOKEN" required:"true"`
	GitHubOwner       string `envconfig:"MIKKU_GITHUB_OWNER" required:"true"`
}

func (cfg *Config) validate() error {
	if cfg.GitHubAccessToken == "" {
		return errEmptyGitHubAccessToken
	}

	if cfg.GitHubOwner == "" {
		return errEmptyGitHubOwner
	}

	return nil
}

func readConfig() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("failed to read environment variables: %w", err)
	}
	return cfg, nil
}
