package main

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v28/github"
)

func TestGitHubClient_CreateReleaseByTagName(t *testing.T) {
	t.Parallel()

	type args struct {
		repo    string
		tagName string
	}
	tests := []struct {
		name     string
		args     args
		injector func(*MockGitHubClient) *MockGitHubClient
		want     *github.RepositoryRelease
		wantErr  bool
	}{
		{
			name: "create v1.0.0 release",
			args: args{
				repo:    "test-repo",
				tagName: "v1.0.0",
			},
			injector: func(cli *MockGitHubClient) *MockGitHubClient {
				cli.EXPECT().CreateRelease(gomock.Any(), "test-owner", "test-repo", &github.RepositoryRelease{
					TagName: github.String("v1.0.0"),
					Name:    github.String("v1.0.0"),
					Body: github.String(`
## Changelog
- test (#10) by @p1ass
`),
				}).Return(&github.RepositoryRelease{
					TagName:         github.String("v1.0.0"),
					TargetCommitish: github.String("TargetCommitish"),
					Name:            github.String("v1.0.0"),
					Body: github.String(`
## Changelog
- test (#10) by @p1ass
`),
				}, nil, nil)
				return cli
			},
			want: &github.RepositoryRelease{
				TagName:         github.String("v1.0.0"),
				TargetCommitish: github.String("TargetCommitish"),
				Name:            github.String("v1.0.0"),
				Body: github.String(`
## Changelog
- test (#10) by @p1ass
`),
			},
			wantErr: false,
		},
		{
			name: "create release API failed",
			args: args{
				repo:    "test-repo",
				tagName: "v1.0.0",
			},
			injector: func(cli *MockGitHubClient) *MockGitHubClient {
				cli.EXPECT().CreateRelease(gomock.Any(), "test-owner", "test-repo", &github.RepositoryRelease{
					TagName: github.String("v1.0.0"),
					Name:    github.String("v1.0.0"),
					Body: github.String(`
## Changelog
- test (#10) by @p1ass
`),
				}).Return(nil, nil, fmt.Errorf("error has occured"))
				return cli
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			cli := NewMockGitHubClient(ctrl)
			cli = tt.injector(cli)

			c := newClient("test-owner", cli)

			got, err := c.CreateReleaseByTagName(tt.args.repo, tt.args.tagName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Client.CreateRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateReleaseBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "no error",
			want: `
## Changelog
- test (#10) by @p1ass
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateReleaseBody()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateReleaseBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateReleaseBody() = %v, want %v", got, tt.want)
			}
		})
	}
}
