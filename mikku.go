package mikku

import (
	"errors"
	"fmt"
	"os"
	"strings"
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

// PullRequest is the entry point of `mikku pr` command
func PullRequest(repo, manifestRepo, pathToManifestFile, imageName string) error {
	cfg, err := ReadConfig()
	if err != nil {
		return fmt.Errorf("release: %w", err)
	}

	svc := NewGitHubService(cfg.GitHubOwner, cfg.GitHubAccessToken)

	manifest, hash, err := svc.GetFile(manifestRepo, pathToManifestFile)
	if err != nil {
		return fmt.Errorf("failed to get manifest file: %w", err)
	}

	release, err := svc.getLatestRelease(repo)
	if err != nil {
		return fmt.Errorf("failed to get latest release: %w", err)

	}
	tag := release.GetTagName()

	currentTag, err := getCurrentTag(manifest, imageName)
	if err != nil {
		return fmt.Errorf("failed to get current tag in yaml file: %w", err)
	}
	replacedFile := strings.ReplaceAll(manifest, imageName+":"+currentTag, imageName+":"+tag)

	branch := fmt.Sprintf("bump-%s-to-%s", imageName, tag)

	if err := svc.CreateBranch(manifestRepo, branch); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	commitMessage := fmt.Sprintf("Bump %s to %s", imageName, tag)

	if err := svc.PushFile(manifestRepo, pathToManifestFile, branch, commitMessage, hash, []byte(replacedFile)); err != nil {
		return fmt.Errorf("failed to push updated the manifest file: %w", err)
	}

	title := fmt.Sprintf("Bump %s to %s", imageName, tag)
	body := fmt.Sprintf("Bump %s to %s", imageName, tag)
	pr, err := svc.CreatePullRequest(manifestRepo, branch, title, body)
	if err != nil {
		return fmt.Errorf("failed to create a pull request: %w", err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Pull request created. %s\n", pr.GetHTMLURL())

	return nil
}
