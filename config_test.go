package mikku

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name    string
		setEnv  func() func()
		want    *Config
		wantErr bool
	}{
		{
			name: "no environment variable set",
			setEnv: func() func() {
				return func() {}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "GitHub access token doesn't be set",
			setEnv: func() func() {
				_ = os.Setenv("MIKKU_GITHUB_OWNER", "MIKKU_GITHUB_OWNER")
				return func() {
					_ = os.Unsetenv("MIKKU_GITHUB_OWNER")
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "GitHub owner doesn't be set",
			setEnv: func() func() {
				_ = os.Setenv("MIKKU_GITHUB_ACCESS_TOKEN", "MIKKU_GITHUB_ACCESS_TOKEN")
				return func() {
					_ = os.Unsetenv("MIKKU_GITHUB_ACCESS_TOKEN")
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "All env config be set",
			setEnv: func() func() {
				_ = os.Setenv("MIKKU_GITHUB_ACCESS_TOKEN", "MIKKU_GITHUB_ACCESS_TOKEN")
				_ = os.Setenv("MIKKU_GITHUB_OWNER", "MIKKU_GITHUB_OWNER")
				return func() {
					_ = os.Unsetenv("MIKKU_GITHUB_ACCESS_TOKEN")
					_ = os.Unsetenv("MIKKU_GITHUB_OWNER")
				}
			},
			want: &Config{
				GitHubAccessToken: "MIKKU_GITHUB_ACCESS_TOKEN",
				GitHubOwner:       "MIKKU_GITHUB_OWNER",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.setEnv()()

			got, err := readConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("readConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("readConfig() diff=%s", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestReadPRConfig(t *testing.T) {
	tests := []struct {
		name    string
		setEnv  func() func()
		want    *PRConfig
		wantErr bool
	}{
		{
			name: "no environment variable set",
			setEnv: func() func() {
				return func() {}
			},
			want:    &PRConfig{},
			wantErr: false,
		},
		{
			name: "All env config be set",
			setEnv: func() func() {
				_ = os.Setenv("MIKKU_MANIFEST_REPOSITORY", "MIKKU_MANIFEST_REPOSITORY")
				_ = os.Setenv("MIKKU_MANIFEST_FILEPATH", "MIKKU_MANIFEST_FILEPATH")
				_ = os.Setenv("MIKKU_DOCKER_IMAGE_NAME", "MIKKU_DOCKER_IMAGE_NAME")
				return func() {
					_ = os.Unsetenv("MIKKU_MANIFEST_REPOSITORY")
					_ = os.Unsetenv("MIKKU_MANIFEST_FILEPATH")
					_ = os.Unsetenv("MIKKU_DOCKER_IMAGE_NAME")
				}
			},
			want: &PRConfig{
				ManifestRepository: "MIKKU_MANIFEST_REPOSITORY",
				ManifestFilepath:   "MIKKU_MANIFEST_FILEPATH",
				DockerImageName:    "MIKKU_DOCKER_IMAGE_NAME",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.setEnv()()

			got, err := readPRConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("readConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("readConfig() diff=%s", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestPRConfig_overrideConfig(t *testing.T) {
	t.Parallel()

	type args struct {
		manifestRepo       string
		pathToManifestFile string
		imageName          string
	}
	tests := []struct {
		name string
		cfg  *PRConfig
		args args
		want *PRConfig
	}{
		{
			name: "override manifest repository",
			cfg: &PRConfig{
				ManifestRepository: "ManifestRepository",
				ManifestFilepath:   "ManifestFilepath",
				DockerImageName:    "DockerImageName",
			},
			args: args{
				manifestRepo:       "overrideManifestRepo",
				pathToManifestFile: "",
				imageName:          "",
			},
			want: &PRConfig{
				ManifestRepository: "overrideManifestRepo",
				ManifestFilepath:   "ManifestFilepath",
				DockerImageName:    "DockerImageName",
			},
		},
		{
			name: "override file path",
			cfg: &PRConfig{
				ManifestRepository: "ManifestRepository",
				ManifestFilepath:   "ManifestFilepath",
				DockerImageName:    "DockerImageName",
			},
			args: args{
				manifestRepo:       "",
				pathToManifestFile: "overridePathToManifestFile",
				imageName:          "",
			},
			want: &PRConfig{
				ManifestRepository: "ManifestRepository",
				ManifestFilepath:   "overridePathToManifestFile",
				DockerImageName:    "DockerImageName",
			},
		},
		{
			name: "override docker image",
			cfg: &PRConfig{
				ManifestRepository: "ManifestRepository",
				ManifestFilepath:   "ManifestFilepath",
				DockerImageName:    "DockerImageName",
			},
			args: args{
				manifestRepo:       "",
				pathToManifestFile: "",
				imageName:          "overrideDockerImageName",
			},
			want: &PRConfig{
				ManifestRepository: "ManifestRepository",
				ManifestFilepath:   "ManifestFilepath",
				DockerImageName:    "overrideDockerImageName",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cfg.overrideConfig(tt.args.manifestRepo, tt.args.pathToManifestFile, tt.args.imageName); !cmp.Equal(tt.cfg, tt.want) {
				t.Errorf("Config.overrideConfig() = %v, want %v", tt.cfg, tt.want)
			}
		})
	}
}
