package mikku

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/google/go-github/v28/github"
)

const (
	releaseBodyTemplate = `
## Changelog
{{ range $i, $pr := .PullRequests }}
- {{ $pr.Title }} (#{{ $pr.Number }}) by @{{ $pr.User.Login }}{{ end }}
`
)

var (
	errInvalidSemanticVersioningTag = errors.New("invalid semantic versioning tag")
)

// Release is the entry point of `mikku release` command
func Release(repo string, version string) error {
	cfg, err := ReadConfig()
	if err != nil {
		return fmt.Errorf("release: %w", err)
	}

	svc := NewGitHubService(cfg.GitHubOwner, cfg.GitHubAccessToken)

	isFirstRelease := false

	after, currentTag, err := svc.getLastPublishedAndCurrentTag(repo)
	if err != nil {
		if errors.Is(err, ErrReleaseNotFound) {
			isFirstRelease = true
			_, _ = fmt.Fprintf(os.Stdout, "Release not found. First Release...\n")

		} else {
			return fmt.Errorf("failed to get latest published date or tag: %w", err)
		}
	}

	newTag, err := determineNewTag(version, currentTag)
	if err != nil {
		if errors.Is(err, errInvalidSemanticVersioningTag) && isFirstRelease {
			_, _ = fmt.Fprintf(os.Stderr, "ERROR: You must specify the tag because of the first release.\n")
		}
		return fmt.Errorf("failed to determine new tag: %w", err)
	}

	prs, err := svc.getMergedPRsAfter(repo, after)
	if err != nil {
		return fmt.Errorf("get pull requests: %w", err)
	}

	body, err := generateReleaseBody(prs)
	if err != nil {
		return fmt.Errorf("failed to generate release body: %w", err)
	}

	newRelease, err := svc.CreateReleaseByTagName(repo, newTag, body)
	if err != nil {
		return fmt.Errorf("failed to create release: %w", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "Release was created.\n")
	_, _ = fmt.Fprintf(os.Stdout, *newRelease.HTMLURL+"\n")

	return nil
}

func generateReleaseBody(prs []*github.PullRequest) (string, error) {
	tmpl, err := template.New("body").Parse(releaseBodyTemplate)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}

	buff := bytes.NewBuffer([]byte{})

	body := map[string]interface{}{"PullRequests": prs}

	if err := tmpl.Execute(buff, body); err != nil {
		return "", fmt.Errorf("template execute error: %w", err)
	}
	return buff.String(), nil
}
