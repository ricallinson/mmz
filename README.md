# mmz

[![Build Status](https://travis-ci.org/ricallinson/mmz.svg?branch=master)](https://travis-ci.org/ricallinson/mmz)

__UNSTABLE__

Interface for configuring, logging and visualizing data for a Manzanita Micro Zilla controller.

## Install and Run on Raspberry Pi

    sudo apt-get install go
    export GOPATH=$HOME/Library/Go/gocode
    go install github.com/ricallinson/mmz
    $HOME/Library/Go/gocode/src/github.com/ricallinson/mmz
    $HOME/Library/Go/gocode/bin/mmz -hairball /dev/tty.usbserial

Where `/dev/tty.usbserial` is the location of your USB to RS-232 serial port adapter.

You should then see `The Manzanita Micro Zilla interface is now running on port '8080'.` printed to the console.

In a browser now open http://[raspberry-pi-ip-address]:8080/ to see the interface. For example http://192.168.0.98:8080/.

## Usage

Start the application with the following command;

    mmz -hairball /path/to/zilla/serial

For my setup the actual command is;

    mmz -hairball /dev/tty.usbserial

Then open a browser to http://localhost:8080/.

## Setup Development Environment

Requires a [Go](https://golang.org/dl/) environment.

    go get github.com/ricallinson/mmz
    go get github.com/goforgery/forgery2
    go get github.com/goforgery/mustache
    go get github.com/goforgery/static
    go get github.com/tarm/serial
    go get github.com/ricallinson/simplebdd
    go install github.com/ricallinson/mmz

## Test

    go test

## Code Coverage Report

    go test -covermode=count -coverprofile=count.out; go tool cover -html=count.out

## Links

* https://www.stewright.me/2013/05/install-and-run-raspbian-from-a-usb-flash-drive/

## Help

### Connect to Hairball

    screen /dev/tty.usbserial

Exit from Hairball `CTRL + A + \`.
