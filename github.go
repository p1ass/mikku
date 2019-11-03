package main

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

// GitHubClient is a interface for calling GitHub API
type GitHubClient interface {
	CreateRelease(ctx context.Context, owner, repo string, release *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error)
}

// Client is client implementing GitHubClient
type Client struct {
	owner string
	GitHubClient
}

// NewClient returns a pointer of Client
// If accessToken is empty, you can't make any changes to the repository
func NewClient(owner, accessToken string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return newClient(owner, client.Repositories)
}

func newClient(owner string, githubCli GitHubClient) *Client {
	return &Client{
		owner:        owner,
		GitHubClient: githubCli,
	}
}

const releaseTemplate = `
## Changelog
- test (#10) by @p1ass
`

// CreateReleaseByTagName creates GitHub release with a given tag
func (cli *Client) CreateReleaseByTagName(repo, tagName string) (*github.RepositoryRelease, error) {
	body, err := generateReleaseBody()
	if err != nil {
		return nil, fmt.Errorf("generate release body: %w", err)
	}

	ctx := context.Background()
	release, _, err := cli.CreateRelease(ctx, cli.owner, repo, &github.RepositoryRelease{
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
