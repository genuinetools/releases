# releases

[![Travis CI](https://img.shields.io/travis/genuinetools/releases.svg?style=for-the-badge)](https://travis-ci.org/genuinetools/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/genuinetools/releases)
[![Github All Releases](https://img.shields.io/github/downloads/genuinetools/releases/total.svg?style=for-the-badge)](https://github.com/genuinetools/releases/releases)

Server to show latest GitHub Releases for a set of repositories.

 * [Installation](README.md#installation)
      * [Binaries](README.md#binaries)
      * [Via Go](README.md#via-go)
      * [Running with Docker](README.md#running-with-docker)
 * [Usage](README.md#usage)

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
          _
 _ __ ___| | ___  __ _ ___  ___  ___
| '__/ _ \ |/ _ \/ _` / __|/ _ \/ __|
| | |  __/ |  __/ (_| \__ \  __/\__ \
|_|  \___|_|\___|\__,_|___/\___||___/

 Server to show latest GitHub Releases for a set of repositories.
 Version: v0.0.5
 Build: 442907b

  -d    run in debug mode
  -interval duration
        interval on which to refetch release data (default 1h0m0s)
  -nouser
        do not include your user
  -orgs value
        organizations to include
  -p int
        port for the server to listen on (shorthand) (default 8080)
  -port int
        port for the server to listen on (default 8080)
  -token string
        GitHub API token (or env var GITHUB_TOKEN)
  -url string
        GitHub Enterprise URL
  -v    print version and exit (shorthand)
  -version
        print version and exit
```
