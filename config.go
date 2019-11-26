package mikku

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/kelseyhightower/envconfig"
)

var (
	errEmptyGitHubAccessToken  = errors.New("should be set MIKKU_GITHUB_ACCESS_TOKEN")
	errEmptyGitHubOwner        = errors.New("should be set MIKKU_GITHUB_OWNER")
	errEmptyManifestRepository = errors.New("should be set MIKKU_MANIFEST_REPOSITORY or --manifest option")
	errEmptyManifestFilePath   = errors.New("should be set MIKKU_MANIFEST_FILEPATH or --path option")
	errEmptyDockerImageName    = errors.New("should be set MIKKU_DOCKER_IMAGE_NAME or --image option")
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

func (cfg *PRConfig) validate() error {
	if cfg.ManifestRepository == "" {
		return errEmptyManifestRepository
	}

	if cfg.ManifestFilepath == "" {
		return errEmptyManifestFilePath
	}

	if cfg.DockerImageName == "" {
		return errEmptyDockerImageName
	}
	return nil
}

func (cfg *PRConfig) embedRepoInfo(owner, repo string) error {
	info := map[string]interface{}{"Owner": owner, "Repository": repo}

	embedManifestFilepath, err := parse(cfg.ManifestFilepath, info)
	if err != nil {
		return fmt.Errorf("parse manifest filepath: %w", err)
	}
	cfg.ManifestFilepath = embedManifestFilepath

	embedDockerImageName, err := parse(cfg.DockerImageName, info)
	if err != nil {
		return fmt.Errorf("parse docker image name: %w", err)
	}
	cfg.DockerImageName = embedDockerImageName

	return nil
}

func parse(text string, info map[string]interface{}) (string, error) {
	tmpl, err := template.New("text").Parse(text)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}

	buff := bytes.NewBuffer(make([]byte, 0, 20))
	if err := tmpl.Execute(buff, info); err != nil {
		return "", fmt.Errorf("template execute error: %w", err)
	}

	return buff.String(), nil
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
