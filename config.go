package mikku

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config represents config values read from the environment variables
type Config struct {
	GitHubAccessToken  string `envconfig:"MIKKU_GITHUB_ACCESS_TOKEN" required:"true"`
	GitHubOwner        string `envconfig:"MIKKU_GITHUB_OWNER" required:"true"`
	ManifestRepository string `envconfig:"MIKKU_MANIFEST_REPOSITORY"`
	ManifestFilepath   string `envconfig:"MIKKU_MANIFEST_FILEPATH"`
	DockerImageName    string `envconfig:"MIKKU_DOCKER_IMAGE_NAME"`
}

// ReadConfig reads config values from the environment variables
func ReadConfig() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("failed to read environment variables: %w", err)
	}
	return cfg, nil
}
