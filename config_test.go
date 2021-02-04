package mikku

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfig_validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		GitHubAccessToken string
		GitHubOwner       string
		wantErr           error
	}{
		{
			name:              "all fields are filled",
			GitHubAccessToken: "github-access-token",
			GitHubOwner:       "github-owner",
			wantErr:           nil,
		},
		{
			name:              "empty github access token",
			GitHubAccessToken: "",
			GitHubOwner:       "github-owner",
			wantErr:           errEmptyGitHubAccessToken,
		},
		{
			name:              "empty github owner",
			GitHubAccessToken: "github-access-token",
			GitHubOwner:       "",
			wantErr:           errEmptyGitHubOwner,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				GitHubAccessToken: tt.GitHubAccessToken,
				GitHubOwner:       tt.GitHubOwner,
			}
			err := cfg.validate()
			if (tt.wantErr == nil && err != nil) || (tt.wantErr != nil && !errors.As(err, &tt.wantErr)) {
				t.Errorf("Config.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
