# releases

[![Travis CI](https://travis-ci.org/genuinetools/releases.svg?branch=master)](https://travis-ci.org/genuinetools/releases)

Server to show latest GitHub Releases for a set of repositories.

## Installation

#### Binaries

- **darwin** [386](https://github.com/genuinetools/releases/releases/download/v0.0.3/releases-darwin-386) / [amd64](https://github.com/genuinetools/releases/releases/download/v0.0.3/releases-darwin-amd64)
- **freebsd** [386](https://github.com/genuinetools/releases/releases/download/v0.0.3/releases-freebsd-386) / [amd64](https://github.com/genuinetools/releases/releases/download/v0.0.3/releases-freebsd-amd64)
- **linux** [386](https://github.com/genuinetools/releases/releases/download/v0.0.3/releases-linux-386) / [amd64](https://github.com/genuinetools/releases/releases/download/v0.0.3/releases-linux-amd64) / [arm](https://github.com/genuinetools/releases/releases/download/v0.0.3/releases-linux-arm) / [arm64](https://github.com/genuinetools/releases/releases/download/v0.0.3/releases-linux-arm64)
- **solaris** [amd64](https://github.com/genuinetools/releases/releases/download/v0.0.3/releases-solaris-amd64)
- **windows** [386](https://github.com/genuinetools/releases/releases/download/v0.0.3/releases-windows-386) / [amd64](https://github.com/genuinetools/releases/releases/download/v0.0.3/releases-windows-amd64)

#### Via Go

```bash
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
 Version: v0.0.3
 Build: 488407e

  -d    run in debug mode
  -interval string
        interval on which to refetch release data (default "1h")
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
