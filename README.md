# Call it [![Build Status](https://travis-ci.org/pedrolopesme/call-it.svg?branch=master)](https://travis-ci.org/pedrolopesme/call-it) [![Build Status](https://goreportcard.com/badge/github.com/pedrolopesme/call-it)](https://goreportcard.com/report/github.com/pedrolopesme/call-it)
A simple program to benchmark URL responses across multiple requests

### Makefile

This project provides a Makefile with all common operations need to develop, test and build call-it.

* build: generates binaries
* test: runs all tests
* clean: removes binaries
* run: executes main func
* fmt: runs gofmt for all go files


### Building

This project uses [DEP](https://golang.github.io/dep/docs/installation.html)
as package manager. After installed, you'll need to:

```
$ dep ensure
```