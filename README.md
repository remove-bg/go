# Go

[![CircleCI](https://circleci.com/gh/remove-bg/go.svg?style=shield)](https://circleci.com/gh/remove-bg/go)

*Under development*

## Development

Prerequisites:

- `go 1.12`
- [`dep`](https://golang.github.io/dep/)

Getting started:

```
git clone git@github.com:remove-bg/go.git $GOPATH/github.com/remove-bg/go
cd $GOPATH/github.com/remove-bg/go
bin/setup
bin/test
```

To build & try out locally:

```
go build -o removebg main.go
./removebg --help
```
