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
				_ = os.Setenv("MIKKU_MANIFEST_REPOSITORY", "MIKKU_MANIFEST_REPOSITORY")
				_ = os.Setenv("MIKKU_MANIFEST_FILEPATH", "MIKKU_MANIFEST_FILEPATH")
				_ = os.Setenv("MIKKU_DOCKER_IMAGE_NAME", "MIKKU_DOCKER_IMAGE_NAME")
				return func() {
					_ = os.Unsetenv("MIKKU_GITHUB_ACCESS_TOKEN")
					_ = os.Unsetenv("MIKKU_GITHUB_OWNER")
					_ = os.Unsetenv("MIKKU_MANIFEST_REPOSITORY")
					_ = os.Unsetenv("MIKKU_MANIFEST_FILEPATH")
					_ = os.Unsetenv("MIKKU_DOCKER_IMAGE_NAME")
				}
			},
			want: &Config{
				GitHubAccessToken:  "MIKKU_GITHUB_ACCESS_TOKEN",
				GitHubOwner:        "MIKKU_GITHUB_OWNER",
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

			got, err := ReadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("ReadConfig() diff=%s", cmp.Diff(got, tt.want))
			}
		})
	}
}
