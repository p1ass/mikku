package mikku

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config represents config using all commands
type Config struct {
	GitHubAccessToken string `envconfig:"MIKKU_GITHUB_ACCESS_TOKEN" required:"true"`
	GitHubOwner       string `envconfig:"MIKKU_GITHUB_OWNER" required:"true"`
}

// PRConfig represents config using pr command
type PRConfig struct {
	ManifestRepository string `envconfig:"MIKKU_MANIFEST_REPOSITORY"`
	ManifestFilepath   string `envconfig:"MIKKU_MANIFEST_FILEPATH"`
	DockerImageName    string `envconfig:"MIKKU_DOCKER_IMAGE_NAME"`
}

func (cfg *PRConfig) overrideConfig(manifestRepo, pathToManifestFile, imageName string) {
	if manifestRepo != "" {
		cfg.ManifestRepository = manifestRepo
	}
	if pathToManifestFile != "" {
		cfg.ManifestFilepath = pathToManifestFile
	}
	if imageName != "" {
		cfg.DockerImageName = imageName
	}
}

func readConfig() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("failed to read environment variables: %w", err)
	}
	return cfg, nil
}

func readPRConfig() (*PRConfig, error) {
	cfg := &PRConfig{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("failed to read environment variables: %w", err)
	}
	return cfg, nil
}
