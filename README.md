<h1 align="center">
  <br>
  <img src="assets/call-it.png" alt="Call It" width="200">
  <br>
  Call It
  <br>
  <br>
</h1>

<h4 align="center">🎨 Modern TUI for HTTP load testing with beautiful visualizations</h4>

<p align="center">
  <a href="https://travis-ci.org/pedrolopesme/call-it"> <img src="https://api.travis-ci.org/pedrolopesme/call-it.svg?branch=master" /></a>
  <a href="https://goreportcard.com/report/github.com/pedrolopesme/call-it"> <img src="https://goreportcard.com/badge/github.com/pedrolopesme/call-it" /></a>
  <a href="https://codeclimate.com/github/pedrolopesme/call-it/maintainability"> <img src="https://api.codeclimate.com/v1/badges/e7854e559e20c9e250de/maintainability" /></a>
</p>
<br>
 

### Installation

This project uses [go modules](https://blog.golang.org/using-go-modules)
as package manager. You can install the source code with:

```shell
$ go get -v github.com/pedrolopesme/call-it
```

After getting the source code you can use the following command to get the dependencies:

```shell
$ go get -v ./...
```

To get the binary into your path, you can run the following inside the source code's folder:

```shell
$ go install .
```

### Makefile

This project provides a Makefile with all common operations need to develop, test and build call-it.

* build: generates binaries
* test: runs all tests
* clean: removes binaries
* run: executes main func
* fmt: runs gofmt for all go files


### Running tests

Tests were write using [Testify](https://github.com/stretchr/testify). In order to run them, just type:

```shell
$ make test
```

### Credits

These are the main external packages that make up Call It:

| packages | description |
|---|---|
| **[cli](https://github.com/urfave/cli)** | **A simple, fast, and fun package for building command line apps in Go** |
| **[httpmock](https://github.com/jarcoal/httpmock/tree/v1)** | **HTTP mocking for Golang** |
| **[spinner](https://github.com/briandowns/spinner)** | **Go (golang) package for providing a terminal spinner/progress indicator with options** |
| **[tablewriter](https://github.com/olekukonko/tablewriter)** | **ASCII table in golang** |
| **[testify](https://github.com/stretchr/testify)** | **A toolkit with common assertions and mocks that plays nicely with the standard library** |
| **[goscritp](github.com/matryer/goscript)** | **Goscript: Runtime execution of Go code.** |


Call It logo was created by Flat Icons, released under Flaticon Basic License.


### License

[MIT](LICENSE.md)