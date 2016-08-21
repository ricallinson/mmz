# mmz

[![Build Status](https://travis-ci.org/ricallinson/mmz.svg?branch=master)](https://travis-ci.org/ricallinson/mmz)

__UNSTABLE__

Interface for configuring, logging and visualizing data for a Manzanita Micro Zilla controller.

## Install and Run on Raspberry Pi

    sudo apt-get install go
    export GOPATH=$HOME/Library/Go/gocode
    go get github.com/ricallinson/mmz
    ~/Library/Go/gocode/src/github.com/ricallinson/mmz
    ~/Library/Go/gocode/bin/mmz -hairball /dev/usbXXXX

Where `/dev/usbXXXX` is the location of your USB to RS-232 serial port adapter.

You should then see `The Manzanita Micro Zilla interface is now running on port '8080'.` printed to the screen.

In a browser open http://[raspberry-pi-ip-address]:8080/. For example http://192.168.0.98:8080/.

## Setup Development Environment

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
