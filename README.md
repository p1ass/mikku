# mikku

[![Actions Status](https://github.com/p1ass/mikku/workflows/Go%20tests/badge.svg)](https://github.com/p1ass/mikku/actions)
[![Actions Status](https://github.com/p1ass/mikku/workflows/Static%20check/badge.svg)](https://github.com/p1ass/mikku/actions)
[![Release](https://img.shields.io/github/v/release/p1ass/mikku.svg)](https://img.shields.io/github/v/release/p1ass/mikku.svg)
[![license](https://img.shields.io/badge/license-MIT-4183c4.svg)](https://github.com/p1ass/mikku/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/p1ass/mikku)](https://goreportcard.com/report/github.com/p1ass/mikku)

`mikku` is a CLI tool to help version management and Kubernetes manifest updates.

## Features
- Create GitHub releases with bumping Semantic Versioning tag 
	
## Installation

### From GitHub release

If you use Windows or Linux, replace `windows_amd64` or `linux_amd64` instead of `darwin_amd64`.

```bash
$ VERSION=1.0.1
$ curl -O -L https://github.com/p1ass/mikku/releases/download/v${VERSION}/mikku_${VERSION}_darwin_amd64.tar.gz
$ tar -zxvf mikku_${VERSION}_darwin_amd64.tar.gz
$ chmod a+x mikku
$ mv mikku /usr/local/bin/mikku
$ mikku --help
```


Binaries are available on GitHub releases. [p1ass/mikku/releases](https://github.com/p1ass/mikku/releases)

### go get

```bash
$ go install github.com/p1ass/mikku/cmd/mikku@v1.0.1
$ mikku --help
```

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


## Screenshots

### Release

```bash
$ mikku release sample-repository v1.0.0
```

![changelog](images/changelog.png)


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

## For developers

### Build


```bash
$ go build -o mikku cmd/mikku/main.go
```

### Tests

```bash
go test -v ./...
```

## LICENCE

MIT
