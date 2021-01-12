# API Frame

A Go module containing framework code for building microservices

## Installation

Since this is a private module, additional auth steps must be completed to `go get` it.

1. Ensure your [.netrc file](https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html) is configured with an access token to use for gitlab.innovationup.stream.

```
machine gitlab.innovationup.stream
login <gitlab email>
password <gitlab personal access token>
```

2. Set the GOPRIVATE env var when running `go get` so Golang bypasses module proxy servers and downloads directly from gitlab.

```shell
$ go env -w GOPRIVATE=gitlab.innovationup.stream 
$ go get gitlab.innovationup.stream/innovation-upstream/api-frame
```
