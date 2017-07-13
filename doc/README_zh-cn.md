![goreporter](./logo.png)

# goreporter

[![Current Release](https://img.shields.io/github/release/wgliang/goreporter.svg)](https://github.com/wgliang/goreporter/releases/latest)
[![Build Status](https://travis-ci.org/wgliang/goreporter.svg?branch=master)](https://travis-ci.org/wgliang/goreporter)
[![GoDoc](https://godoc.org/github.com/wgliang/goreporter?status.svg)](https://godoc.org/github.com/wgliang/goreporter)
[![License](https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

GoReporter是一个Golang编写的集代码静态分析，单元测试，代码审核和生成报告功能的工具。它会并发运行检测项并将结果规范化为报告：

<!-- MarkdownTOC -->

- [检测项](#检测项)
- [模版](#模版)
- [路线图](#路线图)
- [安装](#安装)
	- [依赖](#依赖)
- [运行](#运行)
- [快速开始](#快速开始)
- [例子](#例子)
- [报告Demo](#报告Demo)
- [致谢](#致谢)

<!-- /MarkdownTOC -->

## 检测项

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

## 模版

- HTML模版文件可以通过选项`-t <模版文件>`设置.

## 路线图

- 开发更多检测项, 例如代码行数统计，函数统计，项目结构等
- 报告展示还有不足需要继续改善，评估模型还不够完美
- 安全
- SQL检查

## 安装

### 依赖

- [Go](https://golang.org/dl/) 1.6版本以上
- [Graphviz](http://www.graphviz.org/Download..php)

两种方式安装

- 1. 安装稳定版本，你可以到tag中下载

- 2. 安装最新版本: `go get -u github.com/wgliang/goreporter`

## 快速开始

安装

## 运行

```
$ goreporter -p [projectRelativePath] -r [reportPath] -e [exceptPackagesName] -f [json/html]  {-t templatePathIfHtml}
```

- -p 有效的相对路径
- -r 报告保存的地址
- -e 例外的包，多个包使用逗号分隔。例如: "linters/aligncheck,linters/cyclo" ).
- -f 生成报告的格式
- -t 模版路径，不设置会使用默认模版

默认会生出HTML格式的报告

## 例子

```
$ goreporter -p ../goreporter -r ../goreporter -t ./templates/template.html
```
你可以在此查看详细:[online-example-report](http://fiisio.me/pages/goreporter-report.html)

例子:github.com/wgliang/logcool

![github.com/wgliang/logcool](./doc/github-com-wgliang-goreporter-logcool.png)

## 致谢

项目Logo由 [Ri Xu](https://github.com/xuri) 处理
