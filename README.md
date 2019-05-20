# s3-grep

[![GitHub release](https://img.shields.io/github/release/dabdada/s3-grep.svg?style=flat-square)](https://github.com/dabdada/s3-grep/releases)
[![GoDoc](https://godoc.org/github.com/dabdada/s3-grep?status.svg)](https://godoc.org/github.com/dabdada/s3-grep)
[![MIT license](https://img.shields.io/github/license/dabdada/s3-grep.svg?style=flat-square)](https://github.com/dabdada/s3-grep/blob/master/LICENSE)


Command Line Interface to grep through s3 buckets

## Installation

### Use go get

    $ go get -u github.com/dabdada/s3-grep

## Usage

    $ s3-grep search_query --profile <some_configured_aws_profile> --bucket <some_bucket_name>

### Example:

    $ s3-grep hello --profile dev --bucket myBucket

If you want to search for multiple seperated words you need to enclose the search query in quotation marks:

    $ s3-grep "hello world" --profile dev --bucket myBucket

## Requirements

You need your AWS credentials configured ([How To](https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/setup-credentials.html)) and have the necassary access rights for the bucket to grep in.

## Development

You need to install Golang >= 1.11 ([Download](https://golang.org/dl/))

    $ git clone git@github.com:dabdada/s3-grep.git
    $ export GO111MODULE=on

Build with `go build` and run tests with `go test ./...`
