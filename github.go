package mikku

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

const (
	baseBranch      = "master"
	listPerPage     = 10
	releaseTemplate = `
## Changelog
- test (#10) by @p1ass
`
)

var (
	// ErrReleaseNotFound represents error that the release does not found
	ErrReleaseNotFound = errors.New("release not found")
)

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

// GitHubRepositoriesClient is a interface for calling GitHub API about repositories
type GitHubRepositoriesClient interface {
	CreateRelease(ctx context.Context, owner, repo string, release *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error)
	GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error)
}

// GitHubPullRequestsClient is a interface for calling GitHub API about pull requests
type GitHubPullRequestsClient interface {
	List(ctx context.Context, owner string, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
}

// GitHubService handles application logic using GitHub API
type GitHubService struct {
	owner   string
	repoCli GitHubRepositoriesClient
	prCli   GitHubPullRequestsClient
}

// NewGitHubService returns a pointer of GitHubService
// If accessToken is empty, you can't make any changes to the repository
func NewGitHubService(owner, accessToken string) *GitHubService {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return newGitHubService(owner, client.Repositories, client.PullRequests)
}

func newGitHubService(owner string, githubCli GitHubRepositoriesClient, prCli GitHubPullRequestsClient) *GitHubService {
	return &GitHubService{
		owner:   owner,
		repoCli: githubCli,
		prCli:   prCli,
	}
}

// CreateReleaseByTagName creates GitHub release with a given tag
func (s *GitHubService) CreateReleaseByTagName(repo, tagName string) (*github.RepositoryRelease, error) {
	body, err := generateReleaseBody()
	if err != nil {
		return nil, fmt.Errorf("generate release body: %w", err)
	}

	ctx := context.Background()
	release, _, err := s.repoCli.CreateRelease(ctx, s.owner, repo, &github.RepositoryRelease{
		TagName: github.String(tagName),
		Name:    github.String(tagName),
		Body:    github.String(body),
	})
	if err != nil {
		return nil, fmt.Errorf("call creating release API: %w", err)
	}
	return release, nil
}

func generateReleaseBody() (string, error) {
	tmpl, err := template.New("body").Parse(releaseTemplate)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}

	buff := bytes.NewBuffer([]byte{})

	if err := tmpl.Execute(buff, nil); err != nil {
		return "", fmt.Errorf("template execute error: %w", err)
	}
	return buff.String(), nil
}

// getLatestRelease gets the latest release
func (s *GitHubService) getLatestRelease(repo string) (*github.RepositoryRelease, error) {
	ctx := context.Background()
	release, resp, err := s.repoCli.GetLatestRelease(ctx, s.owner, repo)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%s: %w", repo, ErrReleaseNotFound)
		}
		return nil, fmt.Errorf("call getting the latest release API: %w", err)
	}
	return release, nil
}

// GetMergedPRsAfterLatestRelease gets pull requests which are merged after the latest release
func (s *GitHubService) GetMergedPRsAfterLatestRelease(repo string) ([]*github.PullRequest, error) {
	release, err := s.getLatestRelease(repo)
	if err != nil {
		return nil, fmt.Errorf("get latest release: %w", err)
	}
	prs, err := s.getMergedPRsAfter(repo, release.CreatedAt.Time)
	if err != nil {
		return nil, fmt.Errorf("get pull requests: %w", err)
	}
	return prs, nil
}

func (s *GitHubService) getMergedPRsAfter(repo string, after time.Time) ([]*github.PullRequest, error) {
	opt := &github.PullRequestListOptions{
		State:       "closed",
		Base:        baseBranch,
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: listPerPage},
	}

	var prList []*github.PullRequest
	for {
		ctx := context.Background()
		prs, resp, err := s.prCli.List(ctx, s.owner, repo, opt)
		if err != nil {
			return nil, fmt.Errorf("call listing pull requests API: %w", err)
		}

		extractedPR, done := extractMergedPRsAfter(prs, after)
		prList = append(prList, extractedPR...)
		if done {
			break
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return prList, nil
}

// extractMergedPRsAfter extract merged PRs after a given time
// Return PRs and boolean whether we got all PRs we want
func extractMergedPRsAfter(prs []*github.PullRequest, after time.Time) ([]*github.PullRequest, bool) {
	var prList []*github.PullRequest
	done := false
	for _, pr := range prs {
		if pr.MergedAt != nil && pr.MergedAt.After(after) {
			prList = append(prList, pr)
		}
		if pr.UpdatedAt != nil && !pr.UpdatedAt.After(after) {
			done = true
			break
		}
	}
	return prList, done
}
