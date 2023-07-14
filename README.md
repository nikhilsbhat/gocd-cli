# GoCD commandline interface

[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsbhat/gocd-cli)](https://goreportcard.com/report/github.com/nikhilsbhat/gocd-cli)
[![shields](https://img.shields.io/badge/license-MIT-blue)](https://github.com/nikhilsbhat/gocd-cli/blob/main/LICENSE)
[![shields](https://godoc.org/github.com/nikhilsbhat/gocd-cli?status.svg)](https://godoc.org/github.com/nikhilsbhat/gocd-cli)
[![shields](https://img.shields.io/github/v/tag/nikhilsbhat/gocd-cli.svg)](https://github.com/nikhilsbhat/gocd-cli/tags)
[![shields](https://img.shields.io/github/downloads/nikhilsbhat/gocd-cli/total.svg)](https://github.com/nikhilsbhat/gocd-cli/releases)

command-line interface for `GoCD` that helps in interacting with [GoCD](https://www.gocd.org/) server.

## Introduction

GoCD has user interface from where all the work that this CLI does can be operated, but this cli targets admins who manage GoCD.
By providing cli equivalent support of the UI.

This interacts with `GoCD` server's api to encrypt, decrypt secrets get list of pipelines, create config-repos and many more.

This cli uses GoCD golang [sdk](https://github.com/nikhilsbhat/gocd-sdk-go). If you find bug with CLI, probably that bug would at the SDK.

## Requirements

* [Go](https://golang.org/dl/) 1.17 or above . Installing go can be found [here](https://golang.org/doc/install).
* Basic understanding of CI/CD server [GoCD](https://www.gocd.org/) and GoCD golang [sdk](https://github.com/nikhilsbhat/gocd-sdk-go).

## Authorization

The authorization configuration used for GoCD can be cached locally so that it can be used for future operations.

The command `auth-config` will do the work.

```shell
# Running the below command should cache configurations under $HOME/.gocd/auth_config.yaml.
gocd-cli auth-config store --server-url <gocd-url> --username <username> --password <password>

# User creds cached can be validated using below command.
gocd-cli who-am-i
# The response to the above command should be:
# user: admin

# Once we have authorization configurations cached, we do not need to pass the credentials every time we invoke the cli.
gocd-cli environment list
```

## Documentation

Updated documentation on all available commands and flags can be found [here](https://github.com/nikhilsbhat/gocd-cli/blob/main/docs/doc/gocd-cli.md).

## Installation

* Recommend installing released versions. Release binaries are available on the [releases](https://github.com/nikhilsbhat/gocd-cli/releases) page and docker from [here](https://hub.docker.com/repository/docker/basnik/gocd-cli).
* Can always build it locally by running `go build` against cloned repo.

### Note

* The command `gocd-cli pipeline validate-syntax` would use GoCD's plugin binary to validate the pipeline syntax.
* Since the plugins are jars, it is expected to have Java installed, on the machine from which the command would be executed.
