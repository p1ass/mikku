package mikku

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var mikkuVersion string

var commandRelease = &cli.Command{
	Name:    "release",
	Aliases: []string{"r"},
	Usage:   "Create a tag and a GitHub release",
	UsageText: `
	mikku release <repository> <major | minor | patch | (version)>

	Create a tag and a GitHub release.
	If you execute mikku release <major, minor, or patch>, the latest tag name must be
	compatible with Semantic Versioning.

	- major : major version up Ex. v1.0.0 → v1.0.1
	- minor : minor version up Ex. v1.0.1 → v1.1.0
	- path : patch version up Ex. v1.1.0 → v2.0.0
	- version : create tag with a given version Ex. v1.0.0
	`,
	Action: doRelease,
}

var commandPullRequest = &cli.Command{
	Name:  "pr",
	Usage: "Create a pull request updating Docker image tag written in Kubernetes manifest file",
	UsageText: `
	mikku pr [options...] <repository>

	Create a pull request updating Docker image tag written in Kubernetes manifest file.
	`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "manifest",
			Aliases: []string{"m"},
			Usage:   "Repository existing Kubernetes manifest file. It overrides MIKKU_MANIFEST_REPOSITORY environment variable",
		},
		&cli.StringFlag{
			Name:    "path",
			Aliases: []string{"p"},
			Usage:   "File path where the target docker image is written. It overrides MIKKU_MANIFEST_FILEPATH environment variable",
		},
		&cli.StringFlag{
			Name:    "image",
			Aliases: []string{"i"},
			Usage:   "Docker image name. It overrides MIKKU_DOCKER_IMAGE_NAME environment variable",
		},
	},
	Action: doPullRequest,
}

func doRelease(c *cli.Context) error {
	return nil
}

func doPullRequest(c *cli.Context) error {
	return nil
}

// Run runs commands depending on the given argument
func Run(args []string) {
	app := &cli.App{
		Name:  "mikku",
		Usage: "Bump Semantic Versioning tag, create GitHub release and update Kubernetes manifest file",
		Authors: []*cli.Author{
			{
				Name:  "p1ass",
				Email: "contact@p1ass.com",
			},
		},
		Version: mikkuVersion,
		Commands: []*cli.Command{
			commandRelease,
			commandPullRequest,
		},
	}

	if err := app.Run(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
