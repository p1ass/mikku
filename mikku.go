package mikku

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/google/go-github/v28/github"
)

const (
	releaseBodyTemplate = `
# {{ .Tag }}
## Changelog
{{ range $i, $pr := .PullRequests }}
- {{ $pr.Title }} ([#{{ $pr.Number }}]({{ $pr.HTMLURL}})) by [@{{ $pr.User.Login }}]({{ $pr.User.HTMLURL }})
{{ end }}
`
)

func generateReleaseBody(tag string, prs []*github.PullRequest) (string, error) {
	tmpl, err := template.New("body").Parse(releaseBodyTemplate)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}

	buff := bytes.NewBuffer([]byte{})

	body := map[string]interface{}{"Tag": tag, "PullRequests": prs}

	if err := tmpl.Execute(buff, body); err != nil {
		return "", fmt.Errorf("template execute error: %w", err)
	}
	return buff.String(), nil
}
