# releases

[![Travis CI](https://img.shields.io/travis/genuinetools/releases.svg?style=for-the-badge)](https://travis-ci.org/genuinetools/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/genuinetools/releases)
[![Github All Releases](https://img.shields.io/github/downloads/genuinetools/releases/total.svg?style=for-the-badge)](https://github.com/genuinetools/releases/releases)

Server to show latest GitHub Releases for a set of repositories.

<!-- toc -->

<!-- tocstop -->

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/genuinetools/releases/releases).

#### Via Go

```console
$ go get github.com/genuinetools/releases
```

#### Running with Docker

```console
$ docker run -d --restart always \
    --name releases \
    -p 127.0.0.1:8080:8080 \
    -e GITHUB_TOKEN="<token>" \
    r.j3ss.co/releases --org genuinetools
```

## Usage

```console
$ releases -h
releases -  Server to show latest GitHub Releases for a set of repositories.

Usage: releases <command>

Flags:

  --token                GitHub API token (or env var GITHUB_TOKEN)
  --update-release-body  update the body message for the release as well (default: false)
  --url                  GitHub Enterprise URL (default: <none>)
  -d                     enable debug logging (default: false)
  --interval             interval on which to refetch release data (default: 1h0m0s)
  --nouser               do not include your user (default: false)
  --orgs                 organizations to include (default: [])
  -p, --port             port for the server to listen on (default: 8080)

Commands:

  version  Show the version information.
```
