# Confluent Cloud CLI

[![Build Status](https://semaphoreci.com/api/v1/projects/accef4bb-d1db-491f-b22e-0d438211c888/1992525/shields_badge.svg)](https://semaphoreci.com/confluent/cli)
![Release](release.svg)
[![codecov](https://codecov.io/gh/confluentinc/cli/branch/master/graph/badge.svg?token=67t1cdciLU)](https://codecov.io/gh/confluentinc/cli)

This is the v2 Confluent *Cloud CLI*. It also serves as the backbone for the Confluent "*Converged CLI*" efforts.
In particular, the repository also contains all of the code for the on-prem "*Confluent CLI*", which is also built
as part of the repo's build process.

 * [Install](#install)
    + [One Liner](#one-liner)
      - [Install Dir](#install-dir)
      - [Install Version](#install-version)
    + [Binary Tarball from S3](#binary-tarball-from-s3)
    + [Building From Source](#building-from-source)
  * [Developing](#developing)
    + [Go Version](#go-version)
    + [File Layout](#file-layout)
    + [Build Other Platforms](#build-other-platforms)
    + [URLS](#urls)
  * [Installers](#installers)
    + [Documentation](#documentation)
  * [Testing](#testing)
    + [Integration Tests](#integration-tests)

## Install

The CLI has pre-built binaries for macOS, Linux, and Windows, on both i386 and x86_64 architectures.

You can download a tarball with the binaries. These are both on Github releases and in S3.

### One Liner

The simplest way to install cross platform is with this one-liner:

    curl -sL https://cnfl.io/ccloud-cli | sh

Or for the on-prem binary (while these are separate binaries):

    curl -sL https://cnfl.io/cli | sh

It'll install in `./bin` by default.

#### Install Dir

You can also install to a specific directory. For example, install to `/usr/local/bin` by running:

    curl -sL https://cnfl.io/ccloud-cli | sudo sh -s -- -b /usr/local/bin

#### Install Version

You can list all available versions:

    curl -sL https://cnfl.io/ccloud-cli | sh -s -- -l

    curl -sL https://cnfl.io/cli | sh -s -- -l

And install a particular version if you desire:

    curl -sL https://cnfl.io/ccloud-cli | sudo sh -s -- v0.64.0

This downloads a binary tarball from S3 compiled for your distro and installs it.

### Binary Tarball from S3

You can download a binary tarball from S3.

To list all available versions:

    curl -s "https://s3-us-west-2.amazonaws.com/confluent.cloud?prefix=ccloud-cli/archives/&delimiter=/" | tidy -xml --wrap 100 -i - 2>/dev/null

To list all available packages for a version:

    VERSION=v0.95.0 # or latest
    curl -s "https://s3-us-west-2.amazonaws.com/confluent.cloud?prefix=ccloud-cli/archives/${VERSION#v}/&delimiter=/" | tidy -xml --wrap 100 -i - 2>/dev/null

To download a tarball for your OS and architecture:

    VERSION=v0.95.0 # or latest
    OS=darwin
    ARCH=amd64
    FILE=ccloud_${VERSION}_${OS}_${ARCH}.tar.gz
    curl -s https://s3-us-west-2.amazonaws.com/confluent.cloud/ccloud-cli/archives/${VERSION#v}/${FILE} -o ${FILE}

To install the CLI:

    tar -xzvf ${FILE}
    sudo mv ccloud/ccloud /usr/local/bin

### Building From Source

```
$ make deps
$ make build
$ dist/ccloud/ccloud_$(go env GOOS)_$(go env GOARCH)/ccloud -h # for cloud CLI
$ dist/confluent/confluent_$(go env GOOS)_$(go env GOARCH)/confluent -h # for on-prem Confluent CLI
```

## Developing

This repo requires golang 1.12. We recommend you use `goenv` to manage your go versions.
There's a `.go-version` file in this repo with the exact version we use (and test against in CI).

### Go Version

Fortunately `goenv` supports 1.12 already. Unfortunately this is only in their 2.0 branch which
is still in beta, meaning that its not in brew yet. So we have to build from source.

If you already have it via brew, uninstall it first:

    brew uninstall goenv

Now clone the repo and update your shell profile:

    git clone https://github.com/syndbg/goenv.git $GOPATH/src/github.com/syndbg/goenv
    echo 'export GOENV_ROOT="$GOPATH/src/github.com/syndbg/goenv"' >> ~/.bash_profile
    echo 'export PATH="$GOENV_ROOT/bin:$PATH"' >> ~/.bash_profile
    echo 'eval "$(goenv init -)"' >> ~/.bash_profile

### File Layout

This repo follows the [golang standard project layout](https://github.com/golang-standards/project-layout).

Here's the basic file structure:

* `cmd/confluent/main.go` - entrypoint for the CLI binary
* `internal/cmd/command.go` - bootstraps the root `confluent` CLI command
* `internal/cmd/<command>/<command>.go` - defines each command we support
* `internal/pkg/sdk/<resource>/<resource>.go` - a thin wrapper around `ccloud-sdk-go` to add logging and typed errors
   TODO: if we add logging and typed errors to the SDK, we might be able to drop the pkg/sdk stuff entirely.

Things under `internal/cmd` are commands, things under `internal/pkg` are packages to be used by commands.

When you add a new command or resource, assuming its already in the SDK, you generally just need to create
* `internal/cmd/<command>/<command>.go` (and test)
* `internal/pkg/sdk/<resource>/<resource>.go` (and test)

### Build Other Platforms

If you have a need to build a binary for a platform that is not the current one, use the following to target a different `.goreleaser-*` file matching the destined platform.

    make build-go GORELEASER_SUFFIX=-linux.yml   # build linux
    make build-go GORELEASER_SUFFIX=-mac.yml     # build mac
    make build-go GORELEASER_SUFFIX=-windows.yml # build windows

### URLS
Use the `login` command with the `--url` option to point to a different development environment

## Installers

This repo contains installers for [install-ccloud.sh](./install-ccloud.sh) and
[install-confluent.sh](./install-confluent.sh). These were based on installers
generated by [godownloader](https://github.com/goreleaser/godownloader) and
manually modified to download from S3 instead of GitHub. In turn, godownloader
relies on portable shell functions from [shlib](https://github.com/client9/shlib).

Although they've been manually modified, they're fairly clean and simple bash scripts.
The major modifications include
* reworked `github_release` into `s3_release` function
* updated `TARBALL_URL` and `CHECKSUM_URL` to point to S3 instead of GitHub API
* added a new `-l` flag to list versions from S3, since we can't link to our (private) github repo
* extracted a `BINARY` variable instead of having binary names hardcoded in `execute`
* updated version/tag handling of the `v` prefix; its expected in GitHub and inconsistently used in S3
* updated the usage message, logging, and file comments a bit

### Documentation

The CLI command [reference docs](https://docs.confluent.io/current/cloud/cli/command-reference/index.html)
are programmatically generated from the Cobra commands in this repo.

Just run:

    $ make docs

Cheat sheet:
```
	cli := &cobra.Command{
		Use:               cliName,
		Short:             "This is a short description",
		Long:              "This is a longer synopsis",
	}
```

## Testing

The CLI is tested with a combination of unit tests and integration tests
(backed by mocks). These are both contained within this repo.

We also have end-to-end system tests for
* ccloud-only functionality - [cc-system-tests](https://github.com/confluentinc/cc-system-tests/blob/master/test/cli_test.go)
* on-prem-only functionality - [muckrake](https://github.com/confluentinc/muckrake) (TODO: fix link to CLI tests)

Unit tests exist in `_test.go` files alongside the main source code files.

### Integration Tests

The [./test](./test) directory contains the integration tests. These build a CLI
binary and invoke commands on it. These CLI integration tests roughly follow this
[pattern](http://lucapette.me/writing-integration-tests-for-a-go-cli-application):

1. table tests for quickly testing a variety of CLI commands
1. golden files are expected output fixtures for spec compliance testing
1. http test server for stubbing the Confluent Platform Control Plane API

Read the [CLITest](./test/cli_test.go) configuration to get a better idea
about how to write and configure your own integration tests.

You can run just the integration tests with

    make test TEST_ARGS="./test/... -v"

You can update the golden files from the current output with

    make test TEST_ARGS="./test/... -update"

You can skip rebuilding the CLI if it already exists in `dist` with

    make test TEST_ARGS="./test/... -no-rebuild"

You can mix and match these flags. To update the golden files without rebuilding, and log verbosely

    make test TEST_ARGS="./test/... -update -no-rebuild -v"
