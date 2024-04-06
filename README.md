# DEPRECATED

# API Frame

A Go module containing framework code for building microservices

## Installation

Since this is a private module, additional auth steps must be completed to `go get` it.

1. Ensure your [.netrc file](https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html) is configured with a [Gitlab access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#creating-a-personal-access-token) to use for github.com.

```
machine github.com
login <gitlab email>
password <gitlab personal access token>
```

2. Set the `GOPRIVATE` env var when running `go get` so Golang bypasses module proxy servers and downloads directly from gitlab.

```shell
$ go env -w GOPRIVATE=github.com 
$ go get github.com/innovation-upstream/api-frame
```
