package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"text/template"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

const (
	baseBranch      = "master"
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

// GitHubClient is a interface for calling GitHub API
type GitHubClient interface {
	CreateRelease(ctx context.Context, owner, repo string, release *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error)
	GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error)
}

// GitHubService calls GitHub API through GitHubClient
type GitHubService struct {
	owner string
	cli   GitHubClient
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

	return newGitHubService(owner, client.Repositories)
}

func newGitHubService(owner string, githubCli GitHubClient) *GitHubService {
	return &GitHubService{
		owner: owner,
		cli:   githubCli,
	}
}

// CreateReleaseByTagName creates GitHub release with a given tag
func (s *GitHubService) CreateReleaseByTagName(repo, tagName string) (*github.RepositoryRelease, error) {
	body, err := generateReleaseBody()
	if err != nil {
		return nil, fmt.Errorf("generate release body: %w", err)
	}

	ctx := context.Background()
	release, _, err := s.cli.CreateRelease(ctx, s.owner, repo, &github.RepositoryRelease{
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

// GetLatestRelease gets the latest release
func (s *GitHubService) GetLatestRelease(repo string) (*github.RepositoryRelease, error) {
	ctx := context.Background()
	release, resp, err := s.cli.GetLatestRelease(ctx, s.owner, repo)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%s: %w", repo, ErrReleaseNotFound)
		}
		return nil, fmt.Errorf("call getting the latest release API: %w", err)
	}
	return release, nil
}
