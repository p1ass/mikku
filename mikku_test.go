package mikku

import (
	"strings"
	"testing"

	"github.com/google/go-github/v28/github"
)

func Test_generateReleaseBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tag     string
		prs     []*github.PullRequest
		want    string
		wantErr bool
	}{
		{
			name: "No Pull Requests",
			tag:  "v1.0.0",
			prs:  []*github.PullRequest{},
			want: `
# v1.0.0
## Changelog
`,
			wantErr: false,
		},
		{
			name: "Nil Pull Requests",
			tag:  "v1.0.0",
			prs:  nil,
			want: `
# v1.0.0
## Changelog
`,
			wantErr: false,
		},
		{
			name: "One Pull Request",
			tag:  "v1.0.0",
			prs: []*github.PullRequest{
				{
					Number:  github.Int(1),
					Title:   github.String("Pull Request Title"),
					HTMLURL: github.String("https://github.com/test-owner/test-repo/pull/1"),
					User: &github.User{
						Login:   github.String("test-owner"),
						HTMLURL: github.String("https://github.com/test-owner"),
					},
				},
			},
			want: `
# v1.0.0
## Changelog
- Pull Request Title ([#1](https://github.com/test-owner/test-repo/pull/1)) by [@test-owner](https://github.com/test-owner)
`,
			wantErr: false,
		},
		{
			name: "Two Pull Requests",
			tag:  "v1.0.0",
			prs: []*github.PullRequest{
				{
					Number:  github.Int(2),
					Title:   github.String("Second Pull Request Title"),
					HTMLURL: github.String("https://github.com/test-owner/test-repo/pull/2"),
					User: &github.User{
						Login:   github.String("test-owner"),
						HTMLURL: github.String("https://github.com/test-owner"),
					},
				},
				{
					Number:  github.Int(1),
					Title:   github.String("First Pull Request Title"),
					HTMLURL: github.String("https://github.com/test-owner/test-repo/pull/1"),
					User: &github.User{
						Login:   github.String("test-owner"),
						HTMLURL: github.String("https://github.com/test-owner"),
					},
				},
			},
			want: `
# v1.0.0
## Changelog
- Second Pull Request Title ([#2](https://github.com/test-owner/test-repo/pull/2)) by [@test-owner](https://github.com/test-owner)
- First Pull Request Title ([#1](https://github.com/test-owner/test-repo/pull/1)) by [@test-owner](https://github.com/test-owner)
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateReleaseBody(tt.tag, tt.prs)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateReleaseBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if strings.EqualFold(got, tt.want) {
				t.Errorf("generateReleaseBody() = %v, want %v", got, tt.want)
			}
		})
	}
}
