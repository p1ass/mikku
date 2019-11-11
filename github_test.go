package mikku

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
		body    string
	}
	tests := []struct {
		name     string
		args     args
		injector func(*MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient
		want     *github.RepositoryRelease
		wantErr  bool
	}{
		{
			name: "create v1.0.0 release",
			args: args{
				repo:    "test-repo",
				tagName: "v1.0.0",
				body:    "## v1.0.0",
			},
			injector: func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient {
				cli.EXPECT().CreateRelease(gomock.Any(), "test-owner", "test-repo", &github.RepositoryRelease{
					TagName: github.String("v1.0.0"),
					Name:    github.String("v1.0.0"),
					Body:    github.String("## v1.0.0"),
				}).Return(&github.RepositoryRelease{
					TagName:         github.String("v1.0.0"),
					TargetCommitish: github.String("TargetCommitish"),
					Name:            github.String("v1.0.0"),
					Body:            github.String("## v1.0.0"),
				}, nil, nil)
				return cli
			},
			want: &github.RepositoryRelease{
				TagName:         github.String("v1.0.0"),
				TargetCommitish: github.String("TargetCommitish"),
				Name:            github.String("v1.0.0"),
				Body:            github.String("## v1.0.0"),
			},
			wantErr: false,
		},
		{
			name: "create release API failed",
			args: args{
				repo:    "test-repo",
				tagName: "v1.0.0",
				body:    "## v1.0.0",
			},
			injector: func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient {
				cli.EXPECT().CreateRelease(gomock.Any(), "test-owner", "test-repo", &github.RepositoryRelease{
					TagName: github.String("v1.0.0"),
					Name:    github.String("v1.0.0"),
					Body:    github.String("## v1.0.0"),
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
			cli := NewMockGitHubRepositoriesClient(ctrl)
			cli = tt.injector(cli)

			s := newGitHubService("test-owner", cli, nil)

			got, err := s.CreateReleaseByTagName(tt.args.repo, tt.args.tagName, tt.args.body)
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

func TestGitHubService_getLatestRelease(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		repo     string
		injector func(*MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient
		want     *github.RepositoryRelease
		wantErr  error
	}{{
		name: "release found",
		repo: "test-repo",
		injector: func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient {
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
			injector: func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient {
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
			injector: func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient {
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
			cli := NewMockGitHubRepositoriesClient(ctrl)
			cli = tt.injector(cli)

			s := newGitHubService("test-owner", cli, nil)

			got, err := s.getLatestRelease(tt.repo)
			fmt.Printf("%#v\n", got)
			if (tt.wantErr == nil && err != nil) || (tt.wantErr != nil && !errors.As(err, &tt.wantErr)) {
				t.Errorf("GitHubService.getLatestRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GitHubService.getLatestRelease()diff=%s", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGitHubService_GetMergedPRsAfterLatestRelease(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		repo         string
		repoInjector func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient
		prInjector   func(cli *MockGitHubPullRequestsClient) *MockGitHubPullRequestsClient
		want         []*github.PullRequest
		wantErr      bool
	}{
		{
			name: "success",
			repo: "test-repo",
			repoInjector: func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient {
				cli.EXPECT().GetLatestRelease(gomock.Any(), "test-owner", "test-repo").Return(&github.RepositoryRelease{
					PublishedAt: &github.Timestamp{
						time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}, &github.Response{}, nil)
				return cli
			},
			prInjector: func(cli *MockGitHubPullRequestsClient) *MockGitHubPullRequestsClient {
				cli.EXPECT().List(gomock.Any(), "test-owner", "test-repo", gomock.Any()).Return([]*github.PullRequest{
					{
						ID:        github.Int64(2),
						UpdatedAt: timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
						MergedAt:  timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
					},
					{
						ID:        github.Int64(1),
						UpdatedAt: timeToPointer(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)),
						MergedAt:  timeToPointer(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)),
					},
				}, &github.Response{}, nil)
				return cli
			},
			want: []*github.PullRequest{
				{
					ID:        github.Int64(2),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantErr: false,
		},
		{
			name: "failed to call getting the latest release API",
			repo: "test-repo",
			repoInjector: func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient {
				cli.EXPECT().GetLatestRelease(gomock.Any(), "test-owner", "test-repo").Return(nil,
					&github.Response{
						Response: &http.Response{
							StatusCode: http.StatusInternalServerError,
						},
					}, errors.New("some error"))
				return cli
			},
			prInjector: func(cli *MockGitHubPullRequestsClient) *MockGitHubPullRequestsClient {
				return cli
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "failed to call listing PR API",
			repo: "test-repo",
			repoInjector: func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient {
				cli.EXPECT().GetLatestRelease(gomock.Any(), "test-owner", "test-repo").Return(&github.RepositoryRelease{
					PublishedAt: &github.Timestamp{
						time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}, &github.Response{}, nil)
				return cli
			},
			prInjector: func(cli *MockGitHubPullRequestsClient) *MockGitHubPullRequestsClient {
				cli.EXPECT().List(gomock.Any(), "test-owner", "test-repo", gomock.Any()).Return(nil,
					&github.Response{
						Response: &http.Response{
							StatusCode: http.StatusInternalServerError,
						},
					}, errors.New("some error"))
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
			repoCli := NewMockGitHubRepositoriesClient(ctrl)
			repoCli = tt.repoInjector(repoCli)
			prCli := NewMockGitHubPullRequestsClient(ctrl)
			prCli = tt.prInjector(prCli)

			s := newGitHubService("test-owner", repoCli, prCli)

			got, err := s.GetMergedPRsAfterLatestRelease(tt.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitHubService.GetMergedPRsAfterLatestRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GitHubService.GetMergedPRsAfterLatestRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitHubService_getMergedPRsAfter(t *testing.T) {

	tests := []struct {
		name     string
		repo     string
		after    time.Time
		injector func(cli *MockGitHubPullRequestsClient) *MockGitHubPullRequestsClient
		want     []*github.PullRequest
		wantErr  bool
	}{
		{
			name:  "get all necessary PRs at once",
			repo:  "test-repo",
			after: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			injector: func(cli *MockGitHubPullRequestsClient) *MockGitHubPullRequestsClient {
				cli.EXPECT().List(gomock.Any(), "test-owner", "test-repo", gomock.Any()).Return([]*github.PullRequest{
					{
						ID:        github.Int64(2),
						UpdatedAt: timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
						MergedAt:  timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
					},
					{
						ID:        github.Int64(1),
						UpdatedAt: timeToPointer(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)),
						MergedAt:  timeToPointer(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)),
					},
				}, &github.Response{
					NextPage: 1,
				}, nil)
				return cli
			},
			want: []*github.PullRequest{
				{
					ID:        github.Int64(2),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantErr: false,
		},
		{
			name:  "get all necessary PRs twice",
			repo:  "test-repo",
			after: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			injector: func(cli *MockGitHubPullRequestsClient) *MockGitHubPullRequestsClient {
				cli.EXPECT().List(gomock.Any(), "test-owner", "test-repo", gomock.Any()).Return([]*github.PullRequest{
					{
						ID:        github.Int64(4),
						UpdatedAt: timeToPointer(time.Date(2019, 1, 4, 0, 0, 0, 0, time.UTC)),
						MergedAt:  timeToPointer(time.Date(2019, 1, 4, 0, 0, 0, 0, time.UTC)),
					},
					{
						ID:        github.Int64(3),
						UpdatedAt: timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
						MergedAt:  timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
					},
				}, &github.Response{
					NextPage: 1,
				}, nil)
				cli.EXPECT().List(gomock.Any(), "test-owner", "test-repo", gomock.Any()).Return([]*github.PullRequest{
					{
						ID:        github.Int64(2),
						UpdatedAt: timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
						MergedAt:  timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
					},
					{
						ID:        github.Int64(1),
						UpdatedAt: timeToPointer(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)),
						MergedAt:  timeToPointer(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)),
					},
				}, &github.Response{
					NextPage: 2,
				}, nil)
				return cli
			},
			want: []*github.PullRequest{
				{
					ID:        github.Int64(4),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 4, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 4, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:        github.Int64(3),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:        github.Int64(2),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantErr: false,
		},
		{
			name:  "get all PRs, so no more PR",
			repo:  "test-repo",
			after: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			injector: func(cli *MockGitHubPullRequestsClient) *MockGitHubPullRequestsClient {
				cli.EXPECT().List(gomock.Any(), "test-owner", "test-repo", gomock.Any()).Return([]*github.PullRequest{
					{
						ID:        github.Int64(4),
						UpdatedAt: timeToPointer(time.Date(2019, 1, 4, 0, 0, 0, 0, time.UTC)),
						MergedAt:  timeToPointer(time.Date(2019, 1, 4, 0, 0, 0, 0, time.UTC)),
					},
					{
						ID:        github.Int64(3),
						UpdatedAt: timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
						MergedAt:  timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
					},
				}, &github.Response{
					NextPage: 0,
				}, nil)
				return cli
			},
			want: []*github.PullRequest{
				{
					ID:        github.Int64(4),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 4, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 4, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:        github.Int64(3),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantErr: false,
		},
		{
			name:  "list PR API error",
			repo:  "test-repo",
			after: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			injector: func(cli *MockGitHubPullRequestsClient) *MockGitHubPullRequestsClient {
				cli.EXPECT().List(gomock.Any(), "test-owner", "test-repo", gomock.Any()).Return(nil,
					&github.Response{
						NextPage: 1,
					}, errors.New("some error"))
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
			cli := NewMockGitHubPullRequestsClient(ctrl)
			cli = tt.injector(cli)

			s := newGitHubService("test-owner", nil, cli)
			got, err := s.getMergedPRsAfter(tt.repo, tt.after)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitHubService.getMergedPRsAfter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GitHubService.getMergedPRsAfter() diff=%s", cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_extractMergedPRsAfter(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name     string
		prs      []*github.PullRequest
		after    time.Time
		want     []*github.PullRequest
		wantDone bool
	}{
		{
			name: "extract all prs",
			prs: []*github.PullRequest{
				{
					ID:        github.Int64(2),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:        github.Int64(1),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
			after: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			want: []*github.PullRequest{
				{
					ID:        github.Int64(2),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:        github.Int64(1),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantDone: false,
		},
		{
			name: "not extract unmerged prs",
			prs: []*github.PullRequest{
				{
					ID:        github.Int64(2),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:        github.Int64(1),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
					MergedAt:  nil,
				},
			},
			after: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			want: []*github.PullRequest{
				{
					ID:        github.Int64(2),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantDone: false,
		},
		{
			name: "not extract pr which mergedAt equals to a given after time",
			prs: []*github.PullRequest{
				{
					ID:        github.Int64(2),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:        github.Int64(1),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
			after: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			want: []*github.PullRequest{
				{
					ID:        github.Int64(2),
					UpdatedAt: timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
					MergedAt:  timeToPointer(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantDone: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := extractMergedPRsAfter(tt.prs, tt.after)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("extractMergedPRsAfter() diff=%s", cmp.Diff(got, tt.want))
			}
			if got1 != tt.wantDone {
				t.Errorf("extractMergedPRsAfter() got1 = %v, want %v", got1, tt.wantDone)
			}
		})
	}
}

func timeToPointer(t time.Time) *time.Time {
	return &t
}
