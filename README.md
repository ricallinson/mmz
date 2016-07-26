# mmz

__UNSTABLE__

Interface for configuring, displaying and logging data for a Manzanita Micro Zilla controller.

## Install

    go get github.com/ricallinson/mmz
    go get github.com/goforgery/forgery2
    go get github.com/tarm/serial
    go get github.com/ricallinson/simplebdd
    go install github.com/ricallinson/mmz

## Usage

    mmz

## Test

    go test

## Code Coverage

To be run per module;

    go test -coverprofile=coverage.out; go tool cover -html=coverage.out -o=coverage.html
    open coverage.html
