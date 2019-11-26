# mikku

[![Actions Status](https://github.com/p1ass/mikku/workflows/Go%20tests/badge.svg)](https://github.com/p1ass/mikku/actions)
[![Actions Status](https://github.com/p1ass/mikku/workflows/Static%20check/badge.svg)](https://github.com/p1ass/mikku/actions)
[![GoDoc](https://godoc.org/github.com/p1ass/mikku?status.svg)](https://godoc.org/github.com/p1ass/mikku)
[![license](https://img.shields.io/badge/license-MIT-4183c4.svg)](https://github.com/p1ass/mikku/blob/master/LICENSE)

`mikku` is a CLI command tool to create GitHub release and PR supporting Semantic Versioning and Kubernetes.

## Getting Started

### Prepare GitHub access token

You need an OAuth2 access token. Generate [personal API token](https://github.com/settings/tokens) with *repo* scope.

### Set environment variable

- `MIKKU_GITHUB_ACCESS_TOKEN`: your OAuth2 access token.
- `MIKKU_GITHUB_OWNER`: repository owner or org name. 
    - Ex. `p1ass` when `p1ass/mikku`

```bash
$ export MIKKU_GITHUB_ACCESS_TOKEN=[YOUR_ACCESS_TOKEN]
$ export MIKKU_GITHUB_OWNER=[GITHUB_OWNER_NAME]
```

### Create a new GitHub release to bump patch version

When the latest tag name is `v1.2.3`, the below command bump to `v1.2.4`.

```bash
$ mikku release sample-repository patch
```

Note that `mikku` doesn't build and push a docker image, so you have to do it using CI service such as CircleCI.


### Create a pull request updating docker image tag in Kubernetes manifest file

Update image tag in Kubernetes manifest file existing in `MIKKU_MANIFEST_REPOSITORY` to the latest version.
```bash
$ export MIKKU_MANIFEST_REPOSITORY=sample-manifest-repository
$ export MIKKU_MANIFEST_FILEPATH=manifests/{{.Repository}}/deployment.yml
$ export MIKKU_DOCKER_IMAGE_NAME={{.Owner}}/{{.Repository}}

$ mikku pr sample-repository
```

## Commands

#### `mikku release <repository> <major | minor | patch | (version)>`

Create a tag and a GitHub release.
If you use `major`, `minor`, or `patch`, the latest tag name must be compatible with Semantic Versioning.

##### Arguments

- `major` : major version up
- `minor` : minor version up
- `path` : patch version up
- `version` : create tag with a given version. Ex. `v1.0.0`

##### Examples

```bash
$ mikku release sample-repository v1.0.0
$ mikku release sample-repository patch # v1.0.0 → v1.0.1
$ mikku release sample-repository minor # v1.0.1 → v1.1.0
$ mikku release sample-repository major # v1.1.0 → v2.0.0
```

#### `mikku pr [-m <manifest-repository>] [-p <path-to-manifest-file>]  [-i <image-name>] <repository>`

Create a pull request updating Docker image tag written in Kubernetes manifest file as below:

```yaml
spec:
    containers:
    - name: sample-repository-container
        image: p1ass/sample-repository:v1.0.0
        ↓ Replace
        image: p1ass/sample-repository:v1.0.1
```

##### Options

- `-m`
    - Specify a repository existing Kubernetes manifest file.
    - Optional. 
    - Default : `MIKKU_MANIFEST_REPOSITORY` environment variable.

- `-p` 
	- File path where the target docker image is written. 
    - Optional. 
    - Default : `MIKKU_MANIFEST_FILEPATH` environment variable.
    - You can use [text/template](https://golang.org/pkg/text/template/) in Go.
        - Support variable : `{{.Owner}}`, `{{.Repository}}`
    - Ex. `manifests/{{.Repository}}/deployment.yml`

- `-i`
	- Docker image name.
	- Optional. 
    - Default : `MIKKU_DOCKER_IMAGE_NAME` environment variable.
    - You can use [text/template](https://golang.org/pkg/text/template/) in Go.
        - Support variable : `{{.Owner}}`, `{{.Repository}}`
    - Ex. `asia.gcr.io/{{.Owner}}/{{.Repository}}`



##### Examples

```bash
$ export MIKKU_GITHUB_ACCESS_TOKEN=[YOUR_ACCESS_TOKEN]
$ export MIKKU_GITHUB_OWNER=p1ass

# Set environment variables or add options when executing commands
$ export MIKKU_MANIFEST_REPOSITORY=manifest-repository
$ export MIKKU_MANIFEST_FILEPATH=manifests/{{.Repository}}/deployment.yml
$ export MIKKU_DOCKER_IMAGE_NAME=asia.gcr.io/{{.Owner}}/{{.Repository}}

# The most simple case
# When the latest tag name is `v1.0.1`,
# replace p1ass/sample-repository:v1.0.0 existing in manifest-repository to p1ass/sample-repository:v1.0.1.
$ mikku pr sample-repository

# When the manifest file exists in the same repository
$ mikku pr -m sample-repository sample-repository

# Specify Kubernetes manifest file
$ mikku pr -p {{.Owner}}/{{.Repository}}/deployment.yml sample-repository

# Specify docker image name
$ mikku pr -i docker.p1ass.com/{{.Repository}} sample-repository
```
