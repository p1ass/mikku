package mikku

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/google/go-github/v32/github"
)

const (
	releaseBodyTemplate = `
## Changelog
{{ range $i, $pr := .PullRequests }}
- {{ $pr.Title }} (#{{ $pr.Number }}) by @{{ $pr.User.Login }}{{ end }}
`
)

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
