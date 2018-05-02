<h1 align="center">
  <br>
  <img src="https://cdn.rawgit.com/pedrolopesme/call-it/3280fc97/call-it.png" alt="Call It" width="200">
  <br>
  Call It
  <br>
  <br>
</h1>

<h4 align="center">A CLI program to benchmark URL responses across multiple requests</h4>

<p align="center">
  <a href="https://travis-ci.org/pedrolopesme/call-it"> <img src="https://api.travis-ci.org/pedrolopesme/call-it.svg?branch=master" /></a>
  <a href="https://goreportcard.com/report/github.com/pedrolopesme/call-it"> <img src="https://goreportcard.com/badge/github.com/pedrolopesme/call-it" /></a>
  <a href="https://codeclimate.com/github/pedrolopesme/call-it/maintainability"> <img src="https://api.codeclimate.com/v1/badges/e7854e559e20c9e250de/maintainability" /></a>
</p>
<br>
 
### Makefile

This project provides a Makefile with all common operations need to develop, test and build call-it.

* build: generates binaries
* test: runs all tests
* clean: removes binaries
* run: executes main func
* fmt: runs gofmt for all go files


### Running tests

Tests were write using [Testify](github.com/stretchr/testify/assert). In order to run them, just type:

```shell
$ make test
```


### Building

This project uses [DEP](https://golang.github.io/dep/docs/installation.html)
as package manager. After installed, you'll need to:

```shell
$ dep ensure
```

### Credits

Call It logo was created by Flat Icons, released under Flaticon Basic License.

### License

[MIT](LICENSE.md)