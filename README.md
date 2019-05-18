# s3-grep

[![GitHub release](https://img.shields.io/github/release/dabdada/s3-grep.svg?style=flat-square)](https://github.com/dabdada/s3-grep/releases)
[![GoDoc](https://godoc.org/github.com/dabdada/s3-grep?status.svg)](https://godoc.org/github.com/dabdada/s3-grep)
[![MIT license](https://img.shields.io/github/license/dabdada/s3-grep.svg?style=flat-square)](https://github.com/dabdada/s3-grep/blob/master/LICENSE)


__This is not yet usable. Work in Progress__

Command Line Interface to grep through s3 buckets

## Installation

### Use go get

    $ go get -u github.com/dabdada/s3-grep

## Requirements

You need your AWS credentials configured ([How To](https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/setup-credentials.html)) and have the necassary access rights for the bucket to grep in.

## Development

You need to install Golang >= 1.11 ([Download](https://golang.org/dl/))

    $ git clone git@github.com:dabdada/s3-grep.git
    $ export GO111MODULE=on
