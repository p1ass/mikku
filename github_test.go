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

func TestGitHubClient_createRelease(t *testing.T) {
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

			s := newGitHubClient("test-owner", cli, nil, nil)

			got, err := s.createRelease(tt.args.repo, tt.args.tagName, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("githubClient.CreateRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("githubClient.CreateRelease() = %v, want %v", got, tt.want)
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
			wantErr: errReleaseNotFound,
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

			s := newGitHubClient("test-owner", cli, nil, nil)

			got, err := s.getLatestRelease(tt.repo)
			fmt.Printf("%#v\n", got)
			if (tt.wantErr == nil && err != nil) || (tt.wantErr != nil && !errors.As(err, &tt.wantErr)) {
				t.Errorf("githubClient.getLatestRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("githubClient.getLatestRelease()diff=%s", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGitHubService_getMergedPRsAfter(t *testing.T) {
	t.Parallel()

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

			s := newGitHubClient("test-owner", nil, cli, nil)
			got, err := s.getMergedPRsAfter(tt.repo, tt.after)
			if (err != nil) != tt.wantErr {
				t.Errorf("githubClient.getMergedPRsAfter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("githubClient.getMergedPRsAfter() diff=%s", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGitHubService_getFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		repo        string
		filePath    string
		injector    func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient
		wantContent string
		wantHash    string
		wantErr     error
	}{
		{
			name:     "success in getting file",
			repo:     "test-repo",
			filePath: "manifests/repo/deployment.yml",
			injector: func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient {
				cli.EXPECT().GetContents(gomock.Any(), "test-owner", "test-repo", "manifests/repo/deployment.yml", &github.RepositoryContentGetOptions{
					Ref: baseBranch,
				}).Return(&github.RepositoryContent{
					Encoding: github.String("base64"),
					Content:  github.String("dGVzdC1jb250ZW50"),
					SHA:      github.String("test-hash"),
				}, nil, &github.Response{
					Response: &http.Response{
						StatusCode: http.StatusCreated,
					},
				}, nil)
				return cli
			},
			wantContent: "test-content",
			wantHash:    "test-hash",
			wantErr:     nil,
		},
		{
			name:     "file not found",
			repo:     "test-repo",
			filePath: "manifests/repo/deployment.yml",
			injector: func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient {
				cli.EXPECT().GetContents(gomock.Any(), "test-owner", "test-repo", "manifests/repo/deployment.yml", &github.RepositoryContentGetOptions{
					Ref: baseBranch,
				}).Return(nil, nil, &github.Response{
					Response: &http.Response{
						StatusCode: http.StatusNotFound,
					},
				}, errors.New("file not found"))
				return cli
			},
			wantContent: "",
			wantHash:    "",
			wantErr:     errFileORRepoNotFound,
		},
		{
			name:     "unknown error",
			repo:     "test-repo",
			filePath: "manifests/repo/deployment.yml",
			injector: func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient {
				cli.EXPECT().GetContents(gomock.Any(), "test-owner", "test-repo", "manifests/repo/deployment.yml", &github.RepositoryContentGetOptions{
					Ref: baseBranch,
				}).Return(nil, nil, &github.Response{
					Response: &http.Response{
						StatusCode: http.StatusInternalServerError,
					},
				}, errors.New("unknown error"))
				return cli
			},
			wantContent: "",
			wantHash:    "",
			wantErr:     errors.New("unknown error"),
		},
		{
			name:     "getting content is directory, not file",
			repo:     "test-repo",
			filePath: "manifests/repo/deployment.yml",
			injector: func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient {
				cli.EXPECT().GetContents(gomock.Any(), "test-owner", "test-repo", "manifests/repo/deployment.yml", &github.RepositoryContentGetOptions{
					Ref: baseBranch,
				}).Return(nil, []*github.RepositoryContent{
					{
						Encoding: github.String("base64"),
						Content:  github.String("dGVzdC1jb250ZW50"),
					},
				}, &github.Response{
					Response: &http.Response{
						StatusCode: http.StatusCreated,
					},
				}, nil)
				return cli
			},
			wantContent: "",
			wantHash:    "",
			wantErr:     errContentIsDirectory,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cli := NewMockGitHubRepositoriesClient(ctrl)
			cli = tt.injector(cli)

			s := newGitHubClient("test-owner", cli, nil, nil)

			gotContent, gotHash, err := s.getFile(tt.repo, tt.filePath)
			if (tt.wantErr == nil && err != nil) || (tt.wantErr != nil && !errors.As(err, &tt.wantErr)) {
				t.Errorf("githubClient.getFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotContent != tt.wantContent {
				t.Errorf("githubClient.getFile() = %v, wantContent %v", gotContent, tt.wantContent)
			}
			if gotHash != tt.wantHash {
				t.Errorf("githubClient.getFile() = %v, wantHash %v", gotHash, tt.wantHash)
			}
		})
	}
}

func TestGitHubService_pushFile(t *testing.T) {
	t.Parallel()

	type args struct {
		repo          string
		filePath      string
		branch        string
		commitMessage string
		commitSHA     string
		body          []byte
	}
	tests := []struct {
		name     string
		args     args
		injector func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient

		wantErr bool
	}{
		{
			name: "success in pushing a file",
			args: args{
				repo:          "test-repo",
				filePath:      "test-file-path",
				branch:        "test-branch",
				commitMessage: "test-commit-message",
				commitSHA:     "test-commit-sha",
				body:          []byte("test-body"),
			},
			injector: func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient {
				cli.EXPECT().UpdateFile(gomock.Any(), "test-owner", "test-repo", "test-file-path",
					&github.RepositoryContentFileOptions{
						Message: github.String("test-commit-message"),
						Content: []byte("test-body"),
						SHA:     github.String("test-commit-sha"),
						Branch:  github.String("test-branch"),
						Committer: &github.CommitAuthor{
							Name:  github.String("mikku"),
							Email: github.String("mikku@p1ass.com"),
						},
					}).Return(&github.RepositoryContentResponse{}, &github.Response{
					Response: &http.Response{
						StatusCode: http.StatusCreated,
					},
				}, nil)
				return cli
			},
			wantErr: false,
		},
		{
			name: "unknown error",
			args: args{
				repo:          "test-repo",
				filePath:      "test-file-path",
				branch:        "test-branch",
				commitMessage: "test-commit-message",
				commitSHA:     "test-commit-sha",
				body:          []byte("test-body"),
			},
			injector: func(cli *MockGitHubRepositoriesClient) *MockGitHubRepositoriesClient {
				cli.EXPECT().UpdateFile(gomock.Any(), "test-owner", "test-repo", "test-file-path",
					&github.RepositoryContentFileOptions{
						Message: github.String("test-commit-message"),
						Content: []byte("test-body"),
						SHA:     github.String("test-commit-sha"),
						Branch:  github.String("test-branch"),
						Committer: &github.CommitAuthor{
							Name:  github.String("mikku"),
							Email: github.String("mikku@p1ass.com"),
						},
					}).Return(nil, &github.Response{
					Response: &http.Response{
						StatusCode: http.StatusInternalServerError,
					},
				}, errors.New("unknown error"))
				return cli
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cli := NewMockGitHubRepositoriesClient(ctrl)
			cli = tt.injector(cli)

			s := newGitHubClient("test-owner", cli, nil, nil)

			if err := s.pushFile(tt.args.repo, tt.args.filePath, tt.args.branch, tt.args.commitMessage, tt.args.commitSHA, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("githubClient.pushFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitHubService_createPullRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		repo   string
		branch string
		title  string
		body   string
	}
	tests := []struct {
		name     string
		args     args
		injector func(cli *MockGitHubPullRequestsClient) *MockGitHubPullRequestsClient
		want     *github.PullRequest
		wantErr  bool
	}{
		{
			name: "success in creating a pull request",
			args: args{
				repo:   "test-repo",
				branch: "test-branch",
				title:  "test-title",
				body:   "test-body",
			},
			injector: func(cli *MockGitHubPullRequestsClient) *MockGitHubPullRequestsClient {
				cli.EXPECT().Create(gomock.Any(), "test-owner", "test-repo", &github.NewPullRequest{
					Title: github.String("test-title"),
					Head:  github.String("test-branch"),
					Base:  github.String("master"),
					Body:  github.String("test-body"),
				}).Return(&github.PullRequest{
					ID: github.Int64(1),
				}, &github.Response{
					Response: &http.Response{
						StatusCode: http.StatusCreated,
					},
				}, nil)
				return cli
			},
			want: &github.PullRequest{
				ID: github.Int64(1),
			},
			wantErr: false,
		},
		{
			name: "failed to create a pull request",
			args: args{
				repo:   "test-repo",
				branch: "test-branch",
				title:  "test-title",
				body:   "test-body",
			},
			injector: func(cli *MockGitHubPullRequestsClient) *MockGitHubPullRequestsClient {
				cli.EXPECT().Create(gomock.Any(), "test-owner", "test-repo", &github.NewPullRequest{
					Title: github.String("test-title"),
					Head:  github.String("test-branch"),
					Base:  github.String("master"),
					Body:  github.String("test-body"),
				}).Return(nil, &github.Response{
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
			cli := NewMockGitHubPullRequestsClient(ctrl)
			cli = tt.injector(cli)

			s := newGitHubClient("test-owner", nil, cli, nil)

			got, err := s.createPullRequest(tt.args.repo, tt.args.branch, tt.args.title, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("githubClient.createPullRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("githubClient.createPullRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitHubService_createBranch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		repo     string
		branch   string
		injector func(cli *MockGitHubGitClient) *MockGitHubGitClient
		wantErr  bool
	}{
		{
			name:   "success in creating a branch",
			repo:   "test-repo",
			branch: "test-branch",
			injector: func(cli *MockGitHubGitClient) *MockGitHubGitClient {
				cli.EXPECT().GetRef(gomock.Any(), "test-owner", "test-repo", "heads/master").
					Return(&github.Reference{
						Ref: github.String("refs/heads/master"),
						Object: &github.GitObject{
							Type: github.String("commit"),
							SHA:  github.String("sha-hash"),
						},
					}, &github.Response{
						Response: &http.Response{
							StatusCode: http.StatusOK,
						},
					}, nil)
				cli.EXPECT().CreateRef(gomock.Any(), "test-owner", "test-repo", &github.Reference{
					Ref: github.String("refs/heads/test-branch"),
					Object: &github.GitObject{
						Type: github.String("commit"),
						SHA:  github.String("sha-hash"),
					},
				}).Return(&github.Reference{
					Ref: github.String("refs/heads/test-branch"),
					Object: &github.GitObject{
						Type: github.String("commit"),
						SHA:  github.String("sha-hash"),
					},
				}, &github.Response{Response: &http.Response{StatusCode: http.StatusOK}}, nil)
				return cli
			},
			wantErr: false,
		},
		{
			name:   "failed to create reference",
			repo:   "test-repo",
			branch: "test-branch",
			injector: func(cli *MockGitHubGitClient) *MockGitHubGitClient {
				cli.EXPECT().GetRef(gomock.Any(), "test-owner", "test-repo", "heads/master").
					Return(&github.Reference{
						Ref: github.String("refs/heads/master"),
						Object: &github.GitObject{
							Type: github.String("commit"),
							SHA:  github.String("sha-hash"),
						},
					}, &github.Response{
						Response: &http.Response{
							StatusCode: http.StatusOK,
						},
					}, nil)
				cli.EXPECT().CreateRef(gomock.Any(), "test-owner", "test-repo", &github.Reference{
					Ref: github.String("refs/heads/test-branch"),
					Object: &github.GitObject{
						Type: github.String("commit"),
						SHA:  github.String("sha-hash"),
					},
				}).Return(nil, &github.Response{
					Response: &http.Response{StatusCode: http.StatusInternalServerError},
				}, errors.New("some error"))
				return cli
			},
			wantErr: true,
		},
		{
			name:   "failed to get reference",
			repo:   "test-repo",
			branch: "test-branch",
			injector: func(cli *MockGitHubGitClient) *MockGitHubGitClient {
				cli.EXPECT().GetRef(gomock.Any(), "test-owner", "test-repo", "heads/master").
					Return(nil, &github.Response{
						Response: &http.Response{
							StatusCode: http.StatusInternalServerError,
						},
					}, errors.New("some error"))
				return cli
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cli := NewMockGitHubGitClient(ctrl)
			cli = tt.injector(cli)

			s := newGitHubClient("test-owner", nil, nil, cli)

			if err := s.createBranch(tt.repo, tt.branch); (err != nil) != tt.wantErr {
				t.Errorf("githubClient.createBranch() error = %v, wantErr %v", err, tt.wantErr)
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
