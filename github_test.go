package main

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

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
			defer ctrl.Finish()
			cli := NewMockGitHubClient(ctrl)
			cli = tt.injector(cli)

			s := newGitHubService("test-owner", cli)

			got, err := s.CreateReleaseByTagName(tt.args.repo, tt.args.tagName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitHubService.CreateRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GitHubService.CreateRelease() = %v, want %v", got, tt.want)
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

func TestGitHubService_GetLatestRelease(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		repo     string
		injector func(*MockGitHubClient) *MockGitHubClient
		want     *github.RepositoryRelease
		wantErr  error
	}{{
		name: "release found",
		repo: "test-repo",
		injector: func(cli *MockGitHubClient) *MockGitHubClient {
			cli.EXPECT().GetLatestRelease(gomock.Any(), "test-owner", "test-repo").Return(
				&github.RepositoryRelease{
					TagName:         github.String("v1.0.0"),
					TargetCommitish: github.String("master"),
					Name:            github.String("v1.0.0"),
					Body:            github.String("body"),
					Draft:           github.Bool(false),
					Prerelease:      github.Bool(false),
					ID:              github.Int64(19355655),
					CreatedAt:       &github.Timestamp{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)},
					PublishedAt:     &github.Timestamp{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)},
					URL:             github.String("https://api.github.com/repos/test-repo/test-owner/releases/19355655"),
					HTMLURL:         github.String("https://github.com/test-repo/test-owner/releases/tag/v1.0.0"),
					NodeID:          github.String("MDc6UmVsZWFzZTE5MzU1NjU1"),
				},
				&github.Response{
					Response: &http.Response{
						StatusCode: http.StatusNotFound,
					},
				},
				nil)
			return cli
		},
		want: &github.RepositoryRelease{
			TagName:         github.String("v1.0.0"),
			TargetCommitish: github.String("master"),
			Name:            github.String("v1.0.0"),
			Body:            github.String("body"),
			Draft:           github.Bool(false),
			Prerelease:      github.Bool(false),
			ID:              github.Int64(19355655),
			CreatedAt:       &github.Timestamp{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)},
			PublishedAt:     &github.Timestamp{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)},
			URL:             github.String("https://api.github.com/repos/test-repo/test-owner/releases/19355655"),
			HTMLURL:         github.String("https://github.com/test-repo/test-owner/releases/tag/v1.0.0"),
			NodeID:          github.String("MDc6UmVsZWFzZTE5MzU1NjU1"),
		},
		wantErr: nil,
	},
		{
			name: "release not found",
			repo: "test-repo",
			injector: func(cli *MockGitHubClient) *MockGitHubClient {
				cli.EXPECT().GetLatestRelease(gomock.Any(), "test-owner", "test-repo").Return(nil, &github.Response{
					Response: &http.Response{
						StatusCode: http.StatusNotFound,
					},
				}, fmt.Errorf("404 not found"))
				return cli
			},
			want:    nil,
			wantErr: ErrReleaseNotFound,
		},
		{
			name: "unhandled error",
			repo: "test-repo",
			injector: func(cli *MockGitHubClient) *MockGitHubClient {
				cli.EXPECT().GetLatestRelease(gomock.Any(), "test-owner", "test-repo").Return(nil, &github.Response{
					Response: &http.Response{
						Request: &http.Request{
							Method: "GET",
						},
						StatusCode: http.StatusInternalServerError,
					},
				}, &github.ErrorResponse{})
				return cli
			},
			want:    nil,
			wantErr: &github.ErrorResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cli := NewMockGitHubClient(ctrl)
			cli = tt.injector(cli)

			s := newGitHubService("test-owner", cli)

			got, err := s.GetLatestRelease(tt.repo)
			fmt.Printf("%#v\n", got)
			if (tt.wantErr == nil && err != nil) || (tt.wantErr != nil && !errors.As(err, &tt.wantErr)) {
				t.Errorf("GitHubService.GetLatestRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GitHubService.GetLatestRelease()diff=%s", cmp.Diff(got, tt.want))
			}
		})
	}
}
