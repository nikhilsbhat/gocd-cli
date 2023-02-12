# GoCD commandline interface

[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsbhat/gocd-cli)](https://goreportcard.com/report/github.com/nikhilsbhat/gocd-cli)
[![shields](https://img.shields.io/badge/license-MIT-blue)](https://github.com/nikhilsbhat/gocd-cli/blob/master/LICENSE)
[![shields](https://godoc.org/github.com/nikhilsbhat/gocd-cli?status.svg)](https://godoc.org/github.com/nikhilsbhat/gocd-cli)

command-line interface for `GoCD` that helps in interacting with [GoCD](https://www.gocd.org/) server.

## Introduction

GoCD has user interface from where all the work that this CLI does can be operated, but this cli targets admins who manage GoCD.
By providing cli equivalent support of the UI.

This interacts with `GoCD` server's api to encrypt, decrypt secrets get list of pipelines, create config-repos and many more.

This cli uses GoCD golang [sdk](https://github.com/nikhilsbhat/gocd-sdk-go). If you find bug with CLI, probably that bug would at the SDK.

## Requirements

* [Go](https://golang.org/dl/) 1.17 or above . Installing go can be found [here](https://golang.org/doc/install).
* Basic understanding of CI/CD server [GoCD](https://www.gocd.org/) and GoCD golang [sdk](https://github.com/nikhilsbhat/gocd-sdk-go).


## Installation

* Recommend installing released versions. Release binaries are available on the [releases](https://github.com/nikhilsbhat/gocd-cli/releases) page and docker from [here](https://hub.docker.com/repository/docker/basnik/gocd-cli).
* Can always build it locally by running `go build` against cloned repo.
