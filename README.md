![goreporter](./logo.png)

# goreporter [![Version Status](https://img.shields.io/badge/v2.0-stable-orange.svg)](https://github.com/wgliang/goreporter/releases/tag/version2.0.0)

[![Current Release](https://img.shields.io/github/release/wgliang/goreporter.svg)](https://github.com/wgliang/goreporter/releases/latest)
[![Build Status](https://travis-ci.org/wgliang/goreporter.svg?branch=master)](https://travis-ci.org/wgliang/goreporter)
[![GoDoc](https://godoc.org/github.com/wgliang/goreporter?status.svg)](https://godoc.org/github.com/wgliang/goreporter)
[![License](https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

A Golang tool that does static analysis, unit testing, code review and generate code quality report. This is a tool that concurrently runs a whole bunch of those linters and normalises their output to a report:

<!-- MarkdownTOC -->

- [Supported linters](#supported-linters)
- [Supported template](#supported-template)
- [Todo List](#todo-list)
- [Installing](#installing)
	- [Requirements](#requirements)
- [Run it](#run-it)
- [Credits](#credits)
- [Quickstart](#quickstart)
- [Example](#example)
- [Report-example](#report-example)
- [Credits](#credits)

<!-- /MarkdownTOC -->

## Supported linters

- [unittest](https://github.com/wgliang/goreporter/tree/master/linters/unittest) - Golang unit test status.
- [deadcode](https://github.com/tsenart/deadcode) - Finds unused code.
- [gocyclo](https://github.com/alecthomas/gocyclo) - Computes the cyclomatic complexity of functions.
- [varcheck](https://github.com/opennota/check) - Find unused global variables and constants.
- [structcheck](https://github.com/opennota/check) - Find unused struct fields.
- [aligncheck](https://github.com/opennota/check) - Warn about un-optimally aligned structures.
- [errcheck](https://github.com/kisielk/errcheck) - Check that error return values are used.
- [copycode(dupl)](https://github.com/mibk/dupl) - Reports potentially duplicated code.
- [gosimple](https://github.com/dominikh/go-tools/tree/master/cmd/gosimple) - Report simplifications in code.
- [staticcheck](https://github.com/dominikh/go-tools/tree/master/cmd/staticcheck) - Statically detect bugs, both obvious and subtle ones.
- [godepgraph](https://github.com/kisielk/godepgraph) - Godepgraph is a program for generating a dependency graph of Go packages.
- [misspell](https://github.com/client9/misspell) - Correct commonly misspelled English words... quickly.

## Supported template

- html template file which can be loaded via `-t <file>`.

## Todo List

- Develop more linter, such as code line count statistics, function statistics, project information layer structure display, etc.
- There are still many shortcomings in the presentation of the report, showing that the model is not perfect, and hope that someone can help optimize it.

## Installing

### Requirements

- [Go](https://golang.org/dl/) 1.6+
- [Graphviz](http://www.graphviz.org/Download..php)

There are two options for installing goreporter.

- 1. Install a stable version, eg. `go get -u github.com/wgliang/goreporter/tree/version-1.0.0`.
   I will generally only tag a new stable version when it has passed the Travis regression tests. 
     The downside is that the binary will be called `goreporter.version-1.0.0`.

- 2. Install from HEAD with: `go get -u github.com/wgliang/goreporter`.
   This has the downside that changes to goreporter may break.

## Quickstart

Install goreporter (see above).

## Run it:

```
$ goreporter -p [projectRelativePath] -r [reportPath] -e [exceptPackagesName] -f [json/html]  {-t templatePathIfHtml}
```

- -p Must be a valid golang project path.
- -r Save the path to the report.
- -e Exceptional packages (multiple separated by commas, for example: "linters/aligncheck,linters/cyclo" ).
- -f report formate json OR html.
- -t Template path,if not specified, the default template will be used.

By default, the default template is used to generate reports in html format.

## Example

```
$ goreporter -p ../goreporter -r ../goreporter -t ./templates/template.html
```
you can see result detail:[online-example-report](http://fiisio.me/pages/goreporter-report.html)

example:github.com/wgliang/logcool

![github.com/wgliang/logcool](./doc/github-com-wgliang-goreporter-logcool.png)

## Credits

Templates is designed by liufeifei

Logo is designed by [xuri](https://github.com/Luxurioust)
