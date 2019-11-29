package mikku

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

const (
	baseBranch  = "master"
	listPerPage = 10
)

var (
	// errReleaseNotFound represents error that the release does not found
	errReleaseNotFound = errors.New("release not found")

	errFileORRepoNotFound = errors.New("file or repository not found")
	errContentIsDirectory = errors.New("content is directory, not file")
)

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

// gitHubRepositoriesClient is a interface for calling GitHub API about repositories
type gitHubRepositoriesClient interface {
	CreateRelease(ctx context.Context, owner, repo string, release *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error)
	GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error)

	GetContents(ctx context.Context, owner, repo, path string, opt *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
	UpdateFile(ctx context.Context, owner, repo, path string, opt *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
}

// gitHubPullRequestsClient is a interface for calling GitHub API about pull requests
type gitHubPullRequestsClient interface {
	List(ctx context.Context, owner string, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)

	Create(ctx context.Context, owner string, repo string, pull *github.NewPullRequest) (*github.PullRequest, *github.Response, error)
}

type gitHubGitClient interface {
	GetRef(ctx context.Context, owner string, repo string, ref string) (*github.Reference, *github.Response, error)
	CreateRef(ctx context.Context, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error)
}

// githubClient handles application logic using GitHub API
type githubClient struct {
	owner   string
	repoCli gitHubRepositoriesClient
	prCli   gitHubPullRequestsClient
	gitCli  gitHubGitClient
}

// newGitHubClientUsingEnv returns a pointer of githubClient
// If accessToken is empty, you can't make any changes to the repository
func newGitHubClientUsingEnv(owner, accessToken string) *githubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return newGitHubClient(owner, client.Repositories, client.PullRequests, client.Git)
}

func newGitHubClient(owner string, repoCli gitHubRepositoriesClient, prCli gitHubPullRequestsClient, gitCli gitHubGitClient) *githubClient {
	return &githubClient{
		owner:   owner,
		repoCli: repoCli,
		prCli:   prCli,
		gitCli:  gitCli,
	}
}

func (s *githubClient) getLastPublishedAndCurrentTag(repo string) (time.Time, string, error) {
	after := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	tag := ""
	release, err := s.getLatestRelease(repo)
	if err != nil {
		return after, "", fmt.Errorf("get latest release: %w", err)
	}

	after = release.GetPublishedAt().Time
	tag = release.GetTagName()
	return after, tag, nil
}

// createRelease creates GitHub release with a given tag
func (s *githubClient) createRelease(repo, tagName, body string) (*github.RepositoryRelease, error) {
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

// getLatestRelease gets the latest release
func (s *githubClient) getLatestRelease(repo string) (*github.RepositoryRelease, error) {
	ctx := context.Background()
	release, resp, err := s.repoCli.GetLatestRelease(ctx, s.owner, repo)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%s: %w", repo, errReleaseNotFound)
		}
		return nil, fmt.Errorf("call getting the latest release API: %w", err)
	}
	return release, nil
}

func (s *githubClient) getMergedPRsAfter(repo string, after time.Time) ([]*github.PullRequest, error) {
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

// getFile fetches the file from GitHub and return content and hash
func (s *githubClient) getFile(repo, filePath string) (content string, hash string, err error) {
	ctx := context.Background()
	file, _, resp, err := s.repoCli.GetContents(ctx, s.owner, repo, filePath, &github.RepositoryContentGetOptions{
		Ref: baseBranch,
	})
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return "", "", fmt.Errorf("%s: %w", filePath, errFileORRepoNotFound)
		}
		return "", "", fmt.Errorf("call getting contents api: %w", err)
	}

	if file == nil {
		return "", "", errContentIsDirectory
	}

	content, err = file.GetContent()
	if err != nil {
		return "", "", fmt.Errorf("failed to decode encoded file: %w", err)
	}

	hash = file.GetSHA()
	return content, hash, nil
}

// pushFile pushes the updated file
func (s *githubClient) pushFile(repo, filePath, branch, commitMessage, commitSHA string, body []byte) error {
	ctx := context.Background()
	_, _, err := s.repoCli.UpdateFile(ctx, s.owner, repo, filePath, &github.RepositoryContentFileOptions{
		Message: github.String(commitMessage),
		Content: body,
		SHA:     github.String(commitSHA),
		Branch:  github.String(branch),
		Committer: &github.CommitAuthor{
			Name:  github.String("mikku"),
			Email: github.String("mikku@p1ass.com"),
		},
	})
	if err != nil {
		return fmt.Errorf("call updating file api: %w", err)
	}

	return nil
}

// createPullRequest creates a pull request
func (s *githubClient) createPullRequest(repo, branch, title, body string) (*github.PullRequest, error) {
	ctx := context.Background()
	pr, _, err := s.prCli.Create(ctx, s.owner, repo, &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(branch),
		Base:  github.String(baseBranch),
		Body:  github.String(body),
	})
	if err != nil {
		return nil, fmt.Errorf("call create pull request api: %w", err)
	}

	return pr, nil
}

// createBranch creates new branch
func (s *githubClient) createBranch(repo, branch string) error {
	ctx := context.Background()
	ref, _, err := s.gitCli.GetRef(ctx, s.owner, repo, fmt.Sprintf("heads/%s", baseBranch))
	if err != nil {
		return fmt.Errorf("call getting reference api :%w", err)
	}

	ref.Ref = github.String("refs/heads/" + branch)
	if _, _, err := s.gitCli.CreateRef(ctx, s.owner, repo, ref); err != nil {
		return fmt.Errorf("call creating reference api: %w", err)
	}

	return nil
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
