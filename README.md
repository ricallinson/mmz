# mmz

[![Build Status](https://travis-ci.org/ricallinson/mmz.svg?branch=master)](https://travis-ci.org/ricallinson/mmz)

__UNSTABLE__

Interface for configuring, logging and visualizing data for a Manzanita Micro Zilla controller.

## Install

Requires a [Go](https://golang.org/dl/) environment.

    go get github.com/ricallinson/mmz
    go get github.com/goforgery/forgery2
    go get github.com/tarm/serial
    go get github.com/ricallinson/simplebdd
    go install github.com/ricallinson/mmz

## Usage

Start the application with the following command;

    mmz -serial /path/to/zilla/serial

Then open a browser to http://localhost:8080/.

## Test

    go test

## Code Coverage

To be run per module;

    go test -coverprofile=coverage.out; go tool cover -html=coverage.out -o=coverage.html
    open coverage.html
