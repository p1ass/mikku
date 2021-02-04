package mikku

import (
	"errors"
	"fmt"
	"os"
)

var (
	errInvalidSemanticVersioningTag = errors.New("invalid semantic versioning tag")
)

// Release is the entry point of `mikku release` command
func Release(repo string, bumpTyp string) error {
	cfg, err := readConfig()
	if err != nil {
		return fmt.Errorf("release: %w", err)
	}

	svc := newGitHubClientUsingEnv(cfg.GitHubOwner, cfg.GitHubAccessToken)

	isFirstRelease := false

	after, currentTag, err := svc.getLastPublishedAndCurrentTag(repo)
	if err != nil {
		if errors.Is(err, errReleaseNotFound) {
			isFirstRelease = true
			_, _ = fmt.Fprintf(os.Stdout, "Release not found. First Release...\n")

		} else {
			return fmt.Errorf("failed to get latest published date or tag: %w", err)
		}
	}

	newTag, err := determineNewTag(currentTag, bumpTyp)
	if err != nil {
		if errors.Is(err, errInvalidSemanticVersioningTag) && isFirstRelease {
			return fmt.Errorf("you must specify the tag because of the first release")
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

	newRelease, err := svc.createRelease(repo, newTag, body)
	if err != nil {
		return fmt.Errorf("failed to create release: %w", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "Release was created.\n")
	_, _ = fmt.Fprintf(os.Stdout, *newRelease.HTMLURL+"\n")

	return nil
}
